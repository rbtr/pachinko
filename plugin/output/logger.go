/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package output

import (
	"context"

	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

// Logger is a noop logging output used for dry-runs and testing
type Logger struct{}

func (*Logger) Init(context.Context) error {
	return nil
}

// Receive implements the Plugin interface on the Logger
func (stdr *Logger) Receive(c <-chan types.Media) {
	log.Trace("started stdout output")
	for m := range c {
		log.Tracef("stdout_output: received_input %#v", m)
		if m.SourcePath != "" && m.DestinationPath != "" {
			log.Infof("stdout_output: %s -> %s", m.SourcePath, m.DestinationPath)
		}
	}
}

func init() {
	Register("stdout", func() Output {
		return &Logger{}
	})
}
