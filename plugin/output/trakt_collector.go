/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package output

import (
	"context"
	"strconv"

	"github.com/rbtr/go-trakt"
	internaltrakt "github.com/rbtr/pachinko/internal/trakt"
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/movie"
	"github.com/rbtr/pachinko/types/metadata/tv"
	log "github.com/sirupsen/logrus"
)

var _ Output = (*TraktCollector)(nil)

type TraktCollector struct {
	Authfile string `mapstructure:"authfile"`

	client *internaltrakt.Trakt
}

// Init reads the authfile, creates a client, refreshes the
// credentials, and writes them back to the authfile. Any failures
// will return an error.
func (t *TraktCollector) Init(ctx context.Context) error {
	auth, err := internaltrakt.ReadAuthFile(t.Authfile)
	if err != nil {
		return err
	}
	if t.client, err = internaltrakt.NewTrakt(auth); err != nil {
		return err
	}
	auth, err = t.client.Refresh(ctx)
	if err != nil {
		return err
	}
	return internaltrakt.WriteAuthFile(t.Authfile, auth)
}

func (t *TraktCollector) collectTV(m types.Item) error {
	tvdbID, err := strconv.Atoi(m.Identifiers["tvdb"])
	log.Debugf("trakt_collector: collecting by tvdb id: %d", tvdbID)
	if err != nil {
		return err
	}
	resp, err := t.client.Collection(context.TODO(), &trakt.CollectionBody{
		Episodes: []trakt.Episode{
			{
				IDs: trakt.IDs{
					TVDB: tvdbID,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	log.Debugf("trakt_collector: added %d, updated %d, existing %d", resp.Added.Episodes, resp.Updated.Episodes, resp.Existing.Episodes)
	return nil
}

func (t *TraktCollector) collectMovie(m types.Item) error {
	tmdbID, err := strconv.Atoi(m.Identifiers["tmdb"])
	log.Debugf("trakt_collector: collecting by tmdb id: %d", tmdbID)
	if err != nil {
		return err
	}
	resp, err := t.client.Collection(context.TODO(), &trakt.CollectionBody{
		Movies: []trakt.Movie{
			{
				IDs: trakt.IDs{
					TMDb: tmdbID,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	log.Debugf("trakt_collector: added %d, updated %d, existing %d", resp.Added.Movies, resp.Updated.Movies, resp.Existing.Movies)
	return nil
}

func (t *TraktCollector) Receive(in <-chan types.Item) {
	log.Trace("started trakt_collector output")
	for m := range in {
		log.Tracef("trakt_collector: received_input %#v", m)
		if m.MediaType == tv.TV {
			log.Infof("trakt_collector: collecting tv")
			if err := t.collectTV(m); err != nil {
				log.Error(err)
			}
		}
		if m.MediaType == movie.Movie {
			log.Infof("trakt_collector: collecting movie")
			if err := t.collectMovie(m); err != nil {
				log.Error(err)
			}
		}
	}
}

func init() {
	Register("trakt-collector", func() Output {
		return &TraktCollector{
			Authfile: internaltrakt.DefaultAuthfile,
		}
	})
}
