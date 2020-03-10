/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package output

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

// FilepathMover is a file mover, it will move files from src to dest with
// some options like creating dirs or overwriting existing dests
type FilepathMover struct {
	CreateDirs bool `mapstructure:"create-dirs"`
	DryRun     bool `mapstructure:"dry-run"`
	Overwrite  bool `mapstructure:"overwrite"`
}

func (*FilepathMover) Init(context.Context) error {
	return nil
}

func (mv *FilepathMover) mkdir(dir string) error {
	if mv.DryRun {
		log.Infof("move_output: (DRY_RUN) mkdir %s", dir)
		return nil
	}
	return os.MkdirAll(dir, os.ModePerm)
}

// rename attempts to rename the file
// if the source and dest are on the same volume, this is preferred - it's fast
// and atomic and handled by the filesystem
// if src and dest are on different volumes, it will error with a cross-device
// link message
func (mv *FilepathMover) rename(src, dest string) error {
	if mv.DryRun {
		log.Infof("move_output: (DRY_RUN) rename %s -> %s", src, dest)
		return nil
	}
	return os.Rename(src, dest)
}

// move copies the file from src to dest
// this is slow as it actually copies the bits over from src to dest
// should only be used to move data between volumes since rename is always
// faster within the filesystem boundary
func (mv *FilepathMover) move(src, dest string) error {
	if mv.DryRun {
		log.Infof("move_output: (DRY_RUN) copy %s -> %s", src, dest)
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return errors.Errorf("error opening src: %s", err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return errors.Errorf("error opening dest: %s", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return errors.Errorf("error writing dest: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return errors.Errorf("error removing src: %s", err)
	}
	return nil
}

func (mv *FilepathMover) moveMedia(m types.Media) error {
	if m.DestinationPath == "" {
		return errors.New("move_output: no dest path")
	}
	dir, _ := filepath.Split(m.DestinationPath)
	// check for dest directory, create if doesn't exist and allowed
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if !mv.CreateDirs {
			return errors.Errorf("move_output: dest (%s) does not exist and will not be created", dir)
		}
		if err := mv.mkdir(dir); err != nil {
			return err
		}
	}
	// check for dest file
	if _, err := os.Stat(m.DestinationPath); !os.IsNotExist(err) {
		if !mv.Overwrite {
			return errors.Errorf("move_output: file (%s) already exists and will not be overwritten", m.DestinationPath)
		}
	}
	// move src to dest
	if err := mv.rename(m.SourcePath, m.DestinationPath); err != nil {
		// failed to rename - probably cross-device link so try to move
		return mv.move(m.SourcePath, m.DestinationPath)
	}
	return nil
}

// Receive implements the Plugin interface on the FilepathMover
func (mv *FilepathMover) Receive(c <-chan types.Media) {
	log.Trace("started mover output")
	for m := range c {
		log.Tracef("mover_output: received_input %#v", m)
		if err := mv.moveMedia(m); err != nil {
			log.Errorf("mover_output: %s", err)
		} else {
			log.Infof("move_output: moved %s -> %s", m.SourcePath, m.DestinationPath)
		}
	}
}

func init() {
	Register("path-mover", func() Output {
		return &FilepathMover{
			CreateDirs: true,
			DryRun:     true,
			Overwrite:  false,
		}
	})
}
