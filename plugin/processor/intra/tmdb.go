/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package intra

import (
	"time"

	api "github.com/cyruzin/golang-tmdb"
	"github.com/pkg/errors"
	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/movie"
	log "github.com/sirupsen/logrus"
)

// Client TODO
type TMDbClient struct {
	APIKey string `mapstructure:"api-key"`

	client *api.Client
}

func (c *TMDbClient) Init() error {
	var err error
	if c.client, err = api.Init(c.APIKey); err != nil {
		return err
	}
	return nil
}

// identify returns the ID of best match movie search result, or an error
func (c *TMDbClient) identify(m types.Media) (api.MovieDetails, error) {
	res, err := c.client.GetSearchMovies(m.MovieMetadata.Title, nil)
	if err != nil || res == nil {
		return api.MovieDetails{}, err
	}
	if res.TotalResults == 0 {
		return api.MovieDetails{}, errors.Errorf("tmdb_decorator: no results for tmdb search for %s", m.MovieMetadata.Title)
	}
	// TODO: ugh, why are the inputs and outputs of your library different types for the same field
	details, err := c.client.GetMovieDetails(int(res.Results[0].ID), nil)
	if err != nil {
		return api.MovieDetails{}, err
	}
	if details == nil {
		return api.MovieDetails{}, errors.Errorf("tmdb_decorator: movie details nil for id %d", res.Results[0].ID)
	}
	return *details, nil
}

func (c *TMDbClient) addTMDbMetadata(m types.Media) types.Media {
	movie, err := c.identify(m)
	if err != nil {
		log.Errorf("tmdb_decorator: error identifying movie: %s", err)
		return m
	}

	m.MovieMetadata.Title = movie.Title
	log.Debugf("tmdb_decorator: parsing release date: %s", movie.ReleaseDate)
	if p, err := time.Parse(time.RFC3339, movie.ReleaseDate); err != nil {
		log.Error(err)
	} else {
		m.MovieMetadata.ReleaseYear = p.Year()
	}
	log.Tracef("tmdb_decorator: populated %v from tmdb", m)
	return m
}

func (c *TMDbClient) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started tmdb_decorator processor")
	for m := range in {
		log.Tracef("tmdb_decorator: received input: %v", m)
		if m.Type != movie.Movie {
			log.Infof("tmdb_decorator: %s type %s != Movie, skipping", m.SourcePath, m.Type)
			continue
		}
		out <- c.addTMDbMetadata(m)
	}
}

func init() {
	processor.Register(processor.Intra, "tmdb", func() processor.Processor {
		return &TMDbClient{}
	})
}
