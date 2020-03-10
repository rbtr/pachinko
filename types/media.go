/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package types

import (
	"fmt"

	"github.com/rbtr/pachinko/types/metadata"
	"github.com/rbtr/pachinko/types/metadata/movie"
	"github.com/rbtr/pachinko/types/metadata/tv"
	"github.com/rbtr/pachinko/types/metadata/video"
)

// Media is the container struct for a file flowing through the entire pipeline
type Media struct {
	Identifiers     map[string]string
	SourcePath      string
	DestinationPath string
	Category        Category
	Type            metadata.MediaType
	VideoMetadata   video.Metadata
	TVMetadata      tv.Metadata
	MovieMetadata   movie.Metadata
}

// String string TODO
func (m *Media) String() string {
	if m.Type == tv.TV {
		return fmt.Sprintf("%s Season %d Episode %d", m.TVMetadata.Name, m.TVMetadata.Episode.Season.Number, m.TVMetadata.Episode.Number)
	}
	return m.SourcePath
}
