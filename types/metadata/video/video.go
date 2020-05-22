/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package video

import (
	"fmt"

	"github.com/rbtr/pachinko/types/metadata"
)

// Video defines the Video category enum.
const Video metadata.MediaCategory = "video"

// AudioChannels contains video metadata.
type AudioChannels struct {
	FullRange    int
	LimitedRange int
}

// String formats the AudioChannels struct.
func (audio *AudioChannels) String() string {
	return fmt.Sprintf("%d.%d", audio.FullRange, audio.LimitedRange)
}

// Resolution contains video metadata.
type Resolution struct {
	Width, Height int
}

// String formats the Resolution struct.
func (rez *Resolution) String() string {
	return fmt.Sprintf("%dx%d", rez.Width, rez.Height)
}

// Metadata contains Video metadata.
type Metadata struct {
	Resolution    Resolution
	AudioChannels AudioChannels
}
