/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tvmeta

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/rbtr/pachinko/internal/types"
	"github.com/rbtr/pachinko/internal/types/metadata/tv"
	"github.com/rbtr/pachinko/plugin/processor"
	log "github.com/sirupsen/logrus"
)

var defaultMatchers = []string{
	`(?i)\b([\s\w.-]*)[\s.\/-]+(?:\((\d+)\))?[\s.\/-]?(\d{1,3})[x-](\d{1,3})`,
	`(?i)\b([\s\w.-]*?)?(?:[\s\(.\/-](\d{4})[\s\).\/-])?[\s\w.-]?(?:s+(\d+))(?:\.|\s|-|_|x)*(?:e+(\d+))`,
	`(?i)\b([\s\w.-]*)[\s.\/-]+(?:\((\d+)\))?[\s.\/-]?(?:season|series).?(\d+).?(?:episode)?[^\d(]?(\d+)`,
}

var sanitizer = regexp.MustCompile(`[^\w]`)

type TVPathMetadata struct {
	MatcherStrings []string `mapstructure:"matchers"`
	Sanitize       bool     `mapstructure:"sanitize-name"`

	matchers []*regexp.Regexp
}

func (tv *TVPathMetadata) Init() error {
	for _, str := range tv.MatcherStrings {
		r := regexp.MustCompile(str)
		tv.matchers = append(tv.matchers, r)
	}
	return nil
}

// Extract TODO
func (tv *TVPathMetadata) extractMetadata(m types.Media) types.Media {
	var show, year, season, episode string
	for _, matcher := range tv.matchers {
		subs := matcher.FindAllStringSubmatch(m.SourcePath, -1)
		if len(subs) == 0 {
			continue
		}
		if matches := subs[len(subs)-1]; matches != nil {
			if show == "" && strings.TrimSpace(matches[1]) != "" {
				log.Tracef("tv_path_metadata: %v extracting show name %s from %s", matcher.String(), matches[1], m.SourcePath)
				show = strings.TrimSpace(matches[1])
				if tv.Sanitize {
					show = sanitizer.ReplaceAllString(show, " ")
				}
			}
			if year == "" && strings.TrimSpace(matches[2]) != "" {
				log.Tracef("tv_path_metadata: %v extracting year %s from %s", matcher.String(), matches[2], m.SourcePath)
				year = strings.TrimSpace(matches[2])
			}
			if season == "" && strings.TrimSpace(matches[3]) != "" {
				log.Tracef("tv_path_metadata: %v extracting season number %s from %s", matcher.String(), matches[3], m.SourcePath)
				season = strings.TrimSpace(matches[3])
			}
			if episode == "" && strings.TrimSpace(matches[4]) != "" {
				log.Tracef("tv_path_metadata: %v extracting episode number %s from %s", matcher.String(), matches[4], m.SourcePath)
				episode = strings.TrimSpace(matches[4])
			}
			if show != "" && year != "" && season != "" && episode != "" {
				break
			}
		}
	}
	m.TVMetadata.Name = show
	m.TVMetadata.ReleaseYear, _ = strconv.Atoi(year)
	m.TVMetadata.Season.Number, _ = strconv.Atoi(season)
	m.TVMetadata.Episode.Number, _ = strconv.Atoi(episode)
	return m
}

func (tvp *TVPathMetadata) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started tv_path_metadata processor")
	for m := range in {
		log.Tracef("tv_path_metedata: received %v", m)
		if m.Type != tv.TV {
			log.Debugf("tv_path_metadata: %s type %s != TV, skipping", m.SourcePath, m.Type)
			continue
		}
		log.Debugf("tv_path_metadata: received input %v", m)
		out <- tvp.extractMetadata(m)
	}
}

func init() {
	processor.Register("tv-meta", func() processor.Processor {
		return &TVPathMetadata{
			MatcherStrings: defaultMatchers,
			Sanitize:       true,
		}
	})
}
