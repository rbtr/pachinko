/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package pipeline

import (
	"context"
	"sync"

	"github.com/rbtr/pachinko/plugin/input"
	"github.com/rbtr/pachinko/plugin/output"
	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Buffer int `mapstructure:"buffer"`
}

// Pipeline pipeline
type Pipeline struct {
	Config
	inputs     []input.Input
	processors []processor.Processor
	outputs    []output.Output
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Config: Config{Buffer: 1},
	}
}

func (p *Pipeline) runInputs(ctx context.Context, sink chan<- types.Item) {
	var wg sync.WaitGroup
	for _, input := range p.inputs {
		wg.Add(1)
		go func(_ context.Context, f func(chan<- types.Item), source chan<- types.Item) {
			defer wg.Done()
			f(source)
		}(ctx, input.Consume, sink)
	}
	wg.Wait()
	log.Debug("pipeline: inputs finished")
}

func (p *Pipeline) runProcessors(ctx context.Context, source, sink chan types.Item) {
	// this noop post-processor attaches the final internal input stream to the
	// external sink
	p.processors = processor.AppendFunc(p.processors, func(in <-chan types.Item, _ chan<- types.Item) {
		for m := range in {
			sink <- m
		}
	})
	var wg sync.WaitGroup
	in := source
	out := make(chan types.Item)
	for _, processor := range p.processors {
		wg.Add(1)
		go func(_ context.Context, f func(<-chan types.Item, chan<- types.Item), in <-chan types.Item, out chan<- types.Item) {
			defer wg.Done()
			f(in, out)
			close(out)
		}(ctx, processor.Process, in, out)
		in = out
		out = make(chan types.Item)
	}
	wg.Wait()
	log.Debug("pipeline: processors finished")
}

func (p *Pipeline) runOutputs(ctx context.Context, source chan types.Item) {
	sinks := []chan<- types.Item{}
	var wg sync.WaitGroup
	for _, output := range p.outputs {
		wg.Add(1)
		out := make(chan types.Item)
		go func(_ context.Context, f func(<-chan types.Item), in <-chan types.Item) {
			defer wg.Done()
			f(in)
		}(ctx, output.Receive, out)
		sinks = append(sinks, out)
	}

	wg.Add(1)
	go func(ctx context.Context, in <-chan types.Item, outs []chan<- types.Item) {
		var wgOut sync.WaitGroup
		defer wg.Done()
		for m := range in {
			for _, out := range outs {
				wgOut.Add(1)
				go func(_ context.Context, i types.Item, o chan<- types.Item) {
					defer wgOut.Done()
					o <- i
				}(ctx, m, out)
			}
		}
		wgOut.Wait()
		for _, o := range outs {
			close(o)
		}
	}(ctx, source, sinks)
	wg.Wait()
	log.Debug("pipeline: outputs finished")
}

func (p *Pipeline) WithInputs(inputs ...input.Input) {
	p.inputs = append(p.inputs, inputs...)
}

func (p *Pipeline) WithProcessors(processors ...processor.Processor) {
	p.processors = append(p.processors, processors...)
}

func (p *Pipeline) WithOutputs(outputs ...output.Output) {
	p.outputs = append(p.outputs, outputs...)
}

func (p *Pipeline) Run(ctx context.Context) error {
	log.Debug("running pipeline")

	var wg sync.WaitGroup
	in := make(chan types.Item, p.Buffer)
	out := make(chan types.Item, p.Buffer)

	wg.Add(1)
	go func(ctx context.Context, sink chan types.Item) {
		log.Trace("pipeline: executing input thread")
		defer wg.Done()
		p.runInputs(ctx, sink)
		log.Trace("pipeline: closing input chan")
		close(sink)
	}(ctx, in)

	wg.Add(1)
	go func(ctx context.Context, source, sink chan types.Item) {
		log.Trace("pipeline: executing processor thread")
		defer wg.Done()
		p.runProcessors(ctx, source, sink)
		log.Trace("pipeline: closing processor chan")
		close(sink)
	}(ctx, in, out)

	wg.Add(1)
	go func(ctx context.Context, source chan types.Item) {
		log.Trace("pipeline: executing output thread")
		defer wg.Done()
		p.runOutputs(ctx, source)
	}(ctx, out)

	log.Debug("pipeline: waiting for threads to finish")
	wg.Wait()
	log.Debug("pipeline: threads finished")

	return nil
}
