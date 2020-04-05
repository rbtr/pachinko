/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package output

import (
	"container/heap"
	"context"
	"os"

	"github.com/rbtr/pachinko/plugin/output"
	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

// Deleter is a deleter output used to clean up chaff
type Deleter struct {
	dryRun bool
}

// Init noop
func (d *Deleter) Init(ctx context.Context, cfg output.Config) error {
	d.dryRun = cfg.DryRun
	return nil
}

// Receive implements the Plugin interface on the Deleter
func (d *Deleter) Receive(c <-chan types.Item) {
	log.Trace("started deleter output")
	h := &stringHeap{}
	for m := range c {
		log.Debugf("deleter_output: received_input %#v", m)
		if m.Delete {
			log.Infof("deleter_output: queueing %s", m.SourcePath)
			heap.Push(h, m.SourcePath)
		}
	}
	for h.Len() > 0 {
		path := heap.Pop(h).(string)
		log.Infof("deleter_output: deleting %s", path)
		if d.dryRun {
			continue
		}
		if err := os.Remove(path); err != nil {
			log.Error(err)
		}
	}
}

type stringHeap []string

func (h stringHeap) Len() int           { return len(h) }
func (h stringHeap) Less(i, j int) bool { return len(h[i]) > len(h[j]) }
func (h stringHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *stringHeap) Push(x interface{}) {
	*h = append(*h, x.(string))
}

func (h *stringHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
