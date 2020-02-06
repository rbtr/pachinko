/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package move

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/plugin/output"
	log "github.com/sirupsen/logrus"
)

// MoverOutput is a file mover, it will move files from src to dest with
// some options like creating dirs or overwriting existing dests
type MoveOutput struct {
	CreateDirs bool `mapstructure:"create-dirs"`
	DryRun     bool `mapstructure:"dry-run"`
	Overwrite  bool `mapstructure:"overwrite"`
}

func (*MoveOutput) Init() error {
	return nil
}

func (mv *MoveOutput) mkdir(dir string) error {
	if mv.DryRun {
		log.Infof("move_output: (DRY_RUN) mkdir %s", dir)
		return nil
	}
	return os.MkdirAll(dir, os.ModePerm)
}

func (mv *MoveOutput) move(src, dest string) error {
	if mv.DryRun {
		log.Infof("move_output: (DRY_RUN) mv %s -> %s", src, dest)
		return nil
	}
	return os.Rename(src, dest)
}

func (mv *MoveOutput) moveMedia(m types.Media) error {
	if m.DestinationPath == "" {
		return errors.New("no dest path")
	}
	dir, _ := filepath.Split(m.DestinationPath)
	// check for dest directory, create if doesn't exist and allowed
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if !mv.CreateDirs {
			return errors.Errorf("dest (%s) does not exist and will not be created", dir)
		}
		if err := mv.mkdir(dir); err != nil {
			return err
		}
	}
	// check for dest file
	if _, err := os.Stat(m.DestinationPath); !os.IsNotExist(err) {
		if !mv.Overwrite {
			return errors.Errorf("file (%s) already exists and will not be overwritten", m.DestinationPath)
		}
	}
	// move src to dest
	return mv.move(m.SourcePath, m.DestinationPath)
}

// Receive implements the Plugin interface on the MoveOutput
func (mv *MoveOutput) Receive(c <-chan types.Media) {
	log.Trace("started mover output")
	for m := range c {
		log.Tracef("mover_output: received_input %v", m)
		if err := mv.moveMedia(m); err != nil {
			log.Errorf("mover_output: %s", err)
		} else {
			log.Infof("moved %s -> %s", m.SourcePath, m.DestinationPath)
		}
	}
}

func init() {
	output.Register("move", func() output.Output {
		return &MoveOutput{
			CreateDirs: true,
			DryRun:     true,
			Overwrite:  false,
		}
	})
}
