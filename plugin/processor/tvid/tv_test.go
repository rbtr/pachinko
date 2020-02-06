/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tvid

import (
	"testing"

	internaltesting "github.com/rbtr/pachinko/internal/testing"
	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/internal/types/metadata/tv"
)

func TestTVID_identify(t *testing.T) {
	tvid := &TVID{MatcherStrings: defaultMatchers}
	_ = tvid.Init()
	for _, tt := range internaltesting.TV {
		tt := tt
		for _, in := range tt.Inputs {
			in := in
			t.Run(tt.Name+"::"+in, func(t *testing.T) {
				t.Parallel()
				min := types.Media{SourcePath: in}
				mout := tvid.identify(min)
				if mout.Type != tt.Want.Type {
					t.Errorf("got %s, want %s", mout.Type, tt.Want.Type)
				}
			})
		}
	}
	for _, tt := range internaltesting.NotTV {
		tt := tt
		for _, in := range tt.Inputs {
			in := in
			t.Run(tt.Name+"::"+in, func(t *testing.T) {
				t.Parallel()
				min := types.Media{SourcePath: in}
				mout := tvid.identify(min)
				if mout.Type == tv.TV {
					t.Errorf("got %s but shouldn't have", mout.Type)
				}
			})
		}
	}
}
