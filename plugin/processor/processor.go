/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

/*
Package processor provides an interface for, and implementatations of,
processor plugins which are intermediate steps in the pipeline that
are able to mutate data as it goes by in the stream.
*/
package processor

import (
	"github.com/rbtr/pachinko/internal/types"
	log "github.com/sirupsen/logrus"
)

type Processor interface {
	Process(<-chan types.Media, chan<- types.Media)
	Init() error
}

type ProcessorFunc func(<-chan types.Media, chan<- types.Media)

func (pf ProcessorFunc) Process(in <-chan types.Media, out chan<- types.Media) {
	pf(in, out)
}

func (ProcessorFunc) Init() error {
	return nil
}

func AppendFunc(ps []Processor, fs ...ProcessorFunc) []Processor {
	pfs := make([]Processor, len(fs))
	for i := range fs {
		pfs[i] = fs[i]
	}
	return append(ps, pfs...)
}

var Registry map[string](func() Processor) = map[string](func() Processor){}

func Register(name string, initializer func() Processor) {
	if _, ok := Registry[name]; ok {
		log.Fatalf("processor registry already contains plugin named %s", name)
	}
	Registry[name] = initializer
}
