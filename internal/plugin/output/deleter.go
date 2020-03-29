/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package output

import (
	"container/list"
	"context"
	"os"

	"github.com/rbtr/pachinko/plugin/output"
	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

// Deleter is a deleter output used to clean up chaff
type Deleter struct {
	DryRun bool
}

func (*Deleter) Init(context.Context) error {
	return nil
}

// Receive implements the Plugin interface on the Deleter
func (d *Deleter) Receive(c <-chan types.Item) {
	log.Trace("started deleter output")
	q := list.New()
	for m := range c {
		log.Infof("deleter_output: received_input %#v", m)
		if m.Delete {
			log.Infof("deleter_output: queueing %s", m.SourcePath)
			q.PushBack(m)
		}
	}
	for e := q.Front(); e != nil; e = e.Next() {
		i := (e.Value).(types.Item)
		log.Infof("deleter_output: deleting %s", i.SourcePath)
		if d.DryRun {
			continue
		}
		if err := os.Remove(i.SourcePath); err != nil {
			log.Error(err)
		}
	}
}

func NewDeleter(dryRun bool) output.Output {
	return &Deleter{DryRun: dryRun}
}
