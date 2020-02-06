/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tvdest

import (
	"fmt"
	"path"

	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/internal/types/metadata/tv"
	"github.com/rbtr/pachinko/plugin/processor"
	log "github.com/sirupsen/logrus"
)

type DestinationSolver struct {
	DestDir      string `mapstructure:"dest-dir"`
	EpisodeNames bool   `mapstructure:"episode-names"`
	TVPrefix     string `mapstructure:"tv-prefix"`
	SeasonDirs   bool   `mapstructure:"season-dirs"`
	OutputFormat string `mapstructure:"format"`
}

func (*DestinationSolver) Init() error {
	return nil
}

func (d *DestinationSolver) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started tv_destination processor")
	for m := range in {
		log.Tracef("tv_destination: received input %v", m)
		if m.Type != tv.TV {
			log.Debugf("tv_destination: %s, type %s != TV, skipping", m.SourcePath, m.Type)
			continue
		}

		filename := ""
		if d.EpisodeNames && m.TVMetadata.Episode.Title != "" {
			filename = fmt.Sprintf("%s S%0.2dE%0.2d %s%s",
				m.TVMetadata.Name,
				m.TVMetadata.Episode.Season.Number,
				m.TVMetadata.Episode.Number,
				m.TVMetadata.Episode.Title,
				path.Ext(m.SourcePath))
		} else {
			filename = fmt.Sprintf("%s S%0.2dE%0.2d%s",
				m.TVMetadata.Name,
				m.TVMetadata.Episode.Season.Number,
				m.TVMetadata.Episode.Number,
				path.Ext(m.SourcePath))
		}

		if d.SeasonDirs {
			// => .../tv/Mr Robot/Season 01/Mr Robot S01E01.mkv
			m.DestinationPath = path.Join(
				d.DestDir,
				d.TVPrefix,
				m.TVMetadata.Name,
				fmt.Sprintf("Season %0.2d", m.TVMetadata.Episode.Season.Number),
				filename,
			)
		} else {
			// => .../tv/Mr Robot/Mr Robot S01E01.mkv
			m.DestinationPath = path.Join(
				d.DestDir,
				d.TVPrefix,
				m.TVMetadata.Name,
				filename,
			)
		}
		out <- m
	}
}

func init() {
	processor.Register("tv-dest", func() processor.Processor {
		return &DestinationSolver{
			DestDir:      "dest",
			EpisodeNames: false,
			TVPrefix:     "tv",
			SeasonDirs:   true,
			OutputFormat: "not-implemented",
		}
	})
}
