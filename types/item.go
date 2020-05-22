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

type FileType int

const (
	Directory FileType = iota
	File
)

// Item is the container struct for a file flowing through the entire pipeline.
type Item struct {
	Category        Category
	Delete          bool
	DestinationPath string
	FileType        FileType
	Identifiers     map[string]string
	MediaType       metadata.MediaType
	MovieMetadata   movie.Metadata
	SourcePath      string
	TVMetadata      tv.Metadata
	VideoMetadata   video.Metadata
}

// String formats the Item struct.
func (m *Item) String() string {
	if m.MediaType == tv.TV {
		return fmt.Sprintf("%s Season %d Episode %d", m.TVMetadata.Name, m.TVMetadata.Episode.Season.Number, m.TVMetadata.Episode.Number)
	}
	return m.SourcePath
}
