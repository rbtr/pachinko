/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tvmeta

import (
	"context"
	"testing"

	internaltesting "github.com/rbtr/pachinko/internal/testing"
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/movie"
)

func TestMoviePreProcessor_extractMetadata(t *testing.T) {
	tv := &MoviePreProcessor{MatcherStrings: defaultMovieMatchers, Sanitize: true}
	tv.Init(context.TODO())
	for _, tt := range internaltesting.Movies {
		tt := tt
		for _, in := range tt.Inputs {
			in := in
			t.Run(tt.Name+"::"+in, func(t *testing.T) {
				t.Parallel()
				min := types.Item{SourcePath: in}
				mout := tv.extractMetadata(min)
				if mout.TVMetadata.Name != tt.Want.TVMetadata.Name {
					t.Errorf("got %s, want %s", mout.TVMetadata.Name, tt.Want.TVMetadata.Name)
				}
				if mout.TVMetadata.ReleaseYear != tt.Want.TVMetadata.ReleaseYear {
					t.Errorf("got %d, want %d", mout.TVMetadata.ReleaseYear, tt.Want.TVMetadata.ReleaseYear)
				}
				if mout.TVMetadata.Episode.Season.Number != tt.Want.TVMetadata.Episode.Season.Number {
					t.Errorf("got %d, want %d", mout.TVMetadata.Episode.Season.Number, tt.Want.TVMetadata.Episode.Season.Number)
				}
				if mout.TVMetadata.Episode.Number != tt.Want.TVMetadata.Episode.Number {
					t.Errorf("got %d, want %d", mout.TVMetadata.Episode.Number, tt.Want.TVMetadata.Episode.Number)
				}
			})
		}
	}
}

func TestMoviePreProcessor_identify(t *testing.T) {
	p := &MoviePreProcessor{MatcherStrings: defaultMovieMatchers, Sanitize: true}
	_ = p.Init(context.TODO())
	for _, tt := range internaltesting.Movies {
		tt := tt
		for _, in := range tt.Inputs {
			in := in
			t.Run(tt.Name+"::"+in, func(t *testing.T) {
				t.Parallel()
				min := types.Item{SourcePath: in}
				mout := p.identify(min)
				if mout == (movie.Movie != tt.Want.MediaType) {
					t.Errorf("got %t, want %s", mout, tt.Want.MediaType)
				}
			})
		}
	}
	for _, tt := range internaltesting.NotMovies {
		tt := tt
		for _, in := range tt.Inputs {
			in := in
			t.Run(tt.Name+"::"+in, func(t *testing.T) {
				t.Parallel()
				min := types.Item{SourcePath: in}
				mout := p.identify(min)
				if mout == (movie.Movie != tt.Want.MediaType) {
					t.Errorf("got %t but shouldn't have", mout)
				}
			})
		}
	}
}
