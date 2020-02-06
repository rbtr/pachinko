/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package vidid

import (
	"path"
	"sort"

	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/internal/types/metadata/video"
	"github.com/rbtr/pachinko/internal/util"
	"github.com/rbtr/pachinko/plugin/processor"
	log "github.com/sirupsen/logrus"
)

var defaultFileTypes = []string{
	".avi",
	".divx",
	".m4v",
	".mkv",
	".mov",
	".mp4",
	".xvid",
}

type VideoIdentifier struct {
	FileTypes []string `mapstructure:"file-types"`
}

func (vid *VideoIdentifier) Init() error {
	sort.Strings(vid.FileTypes)
	return nil
}

func (vid *VideoIdentifier) identify(m types.Media) types.Media {
	ext := path.Ext(m.SourcePath)
	if util.StringSliceContains(vid.FileTypes, ext) {
		log.Debugf("video_identifier: identified %s as video", m.SourcePath)
		m.Category = video.Video
	} else {
		log.Debugf("video_identifier: identified %s as not video", m.SourcePath)
	}
	return m
}

func (vid *VideoIdentifier) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started video_identifier processor")
	for m := range in {
		log.Tracef("video_identifier: received input: %v", m)
		out <- vid.identify(m)
	}
}

func init() {
	processor.Register("video-id", func() processor.Processor {
		return &VideoIdentifier{
			FileTypes: defaultFileTypes,
		}
	})
}
