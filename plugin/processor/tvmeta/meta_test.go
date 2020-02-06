/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tvmeta

import (
	"testing"

	internaltesting "github.com/rbtr/pachinko/internal/testing"
	"github.com/rbtr/pachinko/internal/types"
)

func TestTVPathMetadata_extractMetadata(t *testing.T) {
	tests := internaltesting.TV
	tv := &TVPathMetadata{MatcherStrings: defaultMatchers, Sanitize: true}
	tv.Init()
	for _, tt := range tests {
		tt := tt
		for _, in := range tt.Inputs {
			in := in
			t.Run(tt.Name+"::"+in, func(t *testing.T) {
				t.Parallel()
				min := types.Media{SourcePath: in}
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
