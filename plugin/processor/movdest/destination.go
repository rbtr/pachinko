/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package movdest

import (
	"fmt"
	"path"

	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/internal/types/metadata/movie"
	"github.com/rbtr/pachinko/plugin/processor"
	log "github.com/sirupsen/logrus"
)

type DestinationSolver struct {
	DestDir      string `mapstructure:"dest-dir"`
	MovieDirs    bool   `mapstructure:"movie-dirs"`
	MoviesPrefix string `mapstructure:"movie-prefix"`
	OutputFormat string `mapstructure:"format"`
}

func (*DestinationSolver) Init() error {
	return nil
}

func (d *DestinationSolver) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started movie_dest processor")
	for m := range in {
		log.Tracef("movie_destination: received input %v", m)
		if m.Type != movie.Movie {
			log.Debugf("movie_destination: %s, type %s != Movie, skipping", m.SourcePath, m.Type)
			continue
		}

		if d.MovieDirs {
			m.DestinationPath = path.Join(
				d.DestDir,
				d.MoviesPrefix,
				fmt.Sprintf("%s (%d)", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear),
				fmt.Sprintf("%s (%d)%s", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear, path.Ext(m.SourcePath)),
			)
		} else {
			m.DestinationPath = path.Join(
				d.DestDir,
				d.MoviesPrefix,
				fmt.Sprintf("%s (%d)%s", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear, path.Ext(m.SourcePath)),
			)
		}
		out <- m
	}
}

func init() {
	processor.Register("movie-dest", func() processor.Processor {
		return &DestinationSolver{
			DestDir:      "dest",
			MovieDirs:    true,
			MoviesPrefix: "movies",
			OutputFormat: "not-implemented",
		}
	})
}
