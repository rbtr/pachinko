/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tvid

import (
	"regexp"

	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/internal/types/metadata/tv"
	"github.com/rbtr/pachinko/internal/types/metadata/video"
	"github.com/rbtr/pachinko/plugin/processor"
	log "github.com/sirupsen/logrus"
)

var defaultMatchers = []string{
	// matches "S00E00" patterns
	`(?i)(?:s+\d+)(?:\.|\s|-|_|x)*(e+\d+)`,
	// matches "Season 00" patterns
	`(?i)(?:season)+(?:(?:\.|\s|-|_|x|\/)\d+)`,
	// matches "1x1" and "1/1" patterns
	`(?i)(?:[^[:alpha:]\d])\d{1,3}[x-]\d{1,3}`,
}

type TVID struct {
	MatcherStrings []string `mapstructure:"matchers"`

	matchers []*regexp.Regexp
}

func (tvid *TVID) Init() error {
	for _, str := range tvid.MatcherStrings {
		r := regexp.MustCompile(str)
		tvid.matchers = append(tvid.matchers, r)
	}
	return nil
}

func (tvid *TVID) identify(m types.Media) types.Media {
	match := false
	for _, matcher := range tvid.matchers {
		if matcher.MatchString(m.SourcePath) {
			log.Tracef("tv_identifier regexp %s matched %s", matcher, m.SourcePath)
			match = true
			break
		}
	}
	if match {
		log.Debugf("tv_identifier identified %s as tv", m.SourcePath)
		m.Type = tv.TV
	} else {
		log.Debugf("tv_identifier identified %s as not tv", m.SourcePath)
	}
	return m
}

func (tvid *TVID) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started tv_identifier processor")
	for m := range in {
		log.Tracef("tv_identifier: received input: %v", m)
		if m.Category != video.Video {
			log.Debugf("tv_identifier: %s category %s != video, skipping", m.SourcePath, m.Category)
			continue
		}
		out <- tvid.identify(m)
	}
}

func init() {
	processor.Register("tv-id", func() processor.Processor {
		return &TVID{
			MatcherStrings: defaultMatchers,
		}
	})
}
