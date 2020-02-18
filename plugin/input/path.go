/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package input

import (
	"os"
	"path/filepath"

	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

// FilePathInput walks a directory [src-dir], pushing everything in
// in that directory tree in to the pipeline
type FilePathInput struct {
	// SrcDir the directory to ingest
	SrcDir string `mapstructure:"src-dir"`
}

// Init noop
func (*FilePathInput) Init() error {
	return nil
}

// Consume runs the directory ingestion and pushes the contents of the
// directory tree in to the pipeline
func (p *FilePathInput) Consume(sink chan<- types.Media) {
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
	Register("filepath", func() Input {
		return &FilePathInput{
			SrcDir: "./src",
		}
	})
}
