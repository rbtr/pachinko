/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package path

import (
	"os"
	"path/filepath"

	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/plugin/input"
	log "github.com/sirupsen/logrus"
)

// pathinput TODO
type PathInput struct {
	SrcDir string `mapstructure:"src-dir"`
}

func (*PathInput) Init() error {
	return nil
}

// Consume TODO
func (p *PathInput) Consume(sink chan<- types.Media) {
	log.Tracef("started path_input at %s", p.SrcDir)
	count := 0
	if err := filepath.Walk(p.SrcDir, func(path string, info os.FileInfo, err error) error {
		log.Debugf("path_input: encountered %s", path)
		if err != nil {
			log.Error(err)
			return err
		}
		if info.IsDir() {
			log.Tracef("path_input: skipping %s, not a leaf node", path)
			return nil
		}
		log.Debugf("path_input: found file: %s", path)
		sink <- types.Media{SourcePath: path}
		count = count + 1
		return nil
	}); err != nil {
		log.Errorf("path_input: %s", err)
	}
	log.Debugf("path_input: ingested %d files", count)
}

func init() {
	input.Register("path", func() input.Input {
		return &PathInput{
			SrcDir: "./src",
		}
	})
}
