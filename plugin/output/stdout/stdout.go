/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package stdout

import (
	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/plugin/output"
	log "github.com/sirupsen/logrus"
)

// StdoutOutput is a noop logging output used for dry-runs and testing
type StdoutOutput struct{}

func (*StdoutOutput) Init() error {
	return nil
}

// Receive implements the Plugin interface on the StdoutOutput
func (stdr *StdoutOutput) Receive(c <-chan types.Media) {
	log.Trace("started stdout output")
	for m := range c {
		log.Tracef("stdout_output: received_input %v", m)
		if m.SourcePath != "" && m.DestinationPath != "" {
			log.Infof("stdout_output: %s -> %s", m.SourcePath, m.DestinationPath)
		}
	}
}

func init() {
	output.Register("stdout", func() output.Output {
		return &StdoutOutput{}
	})
}
