/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package types

import (
	"regexp"
	"testing"
)

func matchHelper(t *testing.T, name string, matcher map[string]*regexp.Regexp) {
	t.Run(name, func(t *testing.T) {
		for n, r := range matcher {
			if !r.MatchString(n) {
				t.Errorf("failed to match %s", n)
			}
		}
	})
}

func TestMatchers(t *testing.T) {
	matchHelper(t, "AudioChannels", AudioChannels)
	matchHelper(t, "AudioFormats", AudioFormats)
	matchHelper(t, "ColorFormats", ColorFormats)
	matchHelper(t, "VideoFormats", VideoFormats)
	matchHelper(t, "Resolutions", Resolutions)
	matchHelper(t, "Sources", Sources)
}
