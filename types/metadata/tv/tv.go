/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tv

import (
	"time"

	"github.com/rbtr/pachinko/types/metadata"
)

// TV tv
const TV metadata.MediaType = "tv"

// Season season
type Season struct {
	Title  string
	Number int
}

// Episode episode
type Episode struct {
	Title          string
	Number         int
	AbsoluteNumber int
	Season         Season
	AirDate        time.Time
}

// Metadata metadata
type Metadata struct {
	Name        string
	ReleaseYear int
	Episode
}
