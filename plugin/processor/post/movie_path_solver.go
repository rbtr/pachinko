/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package post

import (
	"fmt"
	"path"

	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/movie"
	log "github.com/sirupsen/logrus"
)

type MoviePathSolver struct {
	DestDir      string `mapstructure:"dest-dir"`
	MovieDirs    bool   `mapstructure:"movie-dirs"`
	MoviesPrefix string `mapstructure:"movie-prefix"`
	OutputFormat string `mapstructure:"format"`
}

func (*MoviePathSolver) Init() error {
	return nil
}

func (p *MoviePathSolver) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started movie_destination processor")
	for m := range in {
		log.Tracef("movie_destination: received input %#v", m)
		if m.Type != movie.Movie {
			log.Debugf("movie_destination: %s, type [%s] != Movie, skipping", m.SourcePath, m.Type)
		} else {
			log.Infof("movie_destination: solving dest for %s", m.SourcePath)
			if p.MovieDirs {
				m.DestinationPath = path.Join(
					p.DestDir,
					p.MoviesPrefix,
					fmt.Sprintf("%s (%d)", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear),
					fmt.Sprintf("%s (%d)%s", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear, path.Ext(m.SourcePath)),
				)
			} else {
				m.DestinationPath = path.Join(
					p.DestDir,
					p.MoviesPrefix,
					fmt.Sprintf("%s (%d)%s", m.MovieMetadata.Title, m.MovieMetadata.ReleaseYear, path.Ext(m.SourcePath)),
				)
			}
		}
		out <- m
	}
}

func init() {
	processor.Register(processor.Post, "movie-path-solver", func() processor.Processor {
		return &MoviePathSolver{
			DestDir:      "dest",
			MovieDirs:    true,
			MoviesPrefix: "movies",
			OutputFormat: "not-implemented",
		}
	})
}
