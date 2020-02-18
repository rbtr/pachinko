/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package movie

import (
	"github.com/rbtr/pachinko/types/metadata"
)

// Movie movie
const Movie metadata.MediaType = "movie"

// Metadata metadata
type Metadata struct {
	Title       string
	ReleaseYear int
}
