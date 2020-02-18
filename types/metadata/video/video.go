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

// Video video
const Video metadata.MediaCategory = "video"

// AudioChannels audiochannels
type AudioChannels struct {
	FullRange    int
	LimitedRange int
}

// String string
func (audio *AudioChannels) String() string {
	return fmt.Sprintf("%d.%d", audio.FullRange, audio.LimitedRange)
}

// Resolution resolution
type Resolution struct {
	Width, Height int
}

// String string
func (rez *Resolution) String() string {
	return fmt.Sprintf("%dx%d", rez.Width, rez.Height)
}

// Metadata metadata
type Metadata struct {
	Resolution    Resolution
	AudioChannels AudioChannels
}
