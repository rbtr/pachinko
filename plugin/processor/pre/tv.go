/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package tvmeta

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/tv"
	log "github.com/sirupsen/logrus"
)

var defaultTVMatchers = []string{
	`(?i)\b([\s\w'.-]*)[\s.\/-]+(?:\((\d+)\))?[\s.\/-]?(\d{1,3})[x-](\d{1,3})`,                             // matches 1x1 and 1/1 patterns
	`(?i)\b([\s\w'.-]*?)?(?:[\s\(.\/-](\d{4})[\s\).\/-])?[\s\w.-]?(?:s+(\d+))(?:\.|\s|-|_|x)*(?:e+(\d+))`,  // matches S00E00 patterns
	`(?i)\b([\s\w'.-]*)[\s.\/-]+(?:\((\d+)\))?[\s.\/-]?(?:season|series).?(\d+).?(?:episode)?[^\d(]?(\d+)`, // matches Season 00 patterns
}

var sanitizer = regexp.MustCompile(`[^'\w]`)

type TVPreProcessor struct {
	MatcherStrings []string `mapstructure:"matchers"`
	Sanitize       bool     `mapstructure:"sanitize-name"`

	matchers []*regexp.Regexp
}

func (p *TVPreProcessor) Init(context.Context) error {
	log.Trace("tv_path_metadata: initializing")
	for _, str := range p.MatcherStrings {
		r := regexp.MustCompile(str)
		p.matchers = append(p.matchers, r)
	}
	log.Tracef("tv_path_metadata: initialized %d matchers", len(p.matchers))
	return nil
}

// extract uses the TV regexp to extract metadata from the input.
func (p *TVPreProcessor) extractMetadata(m types.Item) types.Item {
	var show, year, season, episode string
	for _, matcher := range p.matchers {
		subs := matcher.FindAllStringSubmatch(m.SourcePath, -1)
		if len(subs) == 0 {
			continue
		}
		if matches := subs[len(subs)-1]; matches != nil {
			if show == "" && strings.TrimSpace(matches[1]) != "" {
				log.Tracef("tv_path_metadata: %v extracting show name %s from %s", matcher.String(), matches[1], m.SourcePath)
				show = strings.TrimSpace(matches[1])
				if p.Sanitize {
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

// identify tests if the input is matched by any of the TV regexp.
func (p *TVPreProcessor) identify(m types.Item) bool {
	for _, matcher := range p.matchers {
		if matcher.MatchString(m.SourcePath) {
			log.Tracef("tv_path_metadata: regexp %s matched %s", matcher, m.SourcePath)
			return true
		}
		log.Tracef("tv_path_metadata: regexp %s did not match %s", matcher, m.SourcePath)
	}
	log.Tracef("tv_path_metadata: %s did not match identifiers", m.SourcePath)
	return false
}

func (p *TVPreProcessor) Process(in <-chan types.Item, out chan<- types.Item) {
	log.Trace("started tv_path_metadata processor")
	for m := range in {
		log.Tracef("tv_path_metadata: received input: %#v", m)
		if m.Category == types.Video {
			log.Infof("tv_path_metadata: %s category == video, testing for TV", m.SourcePath)
			if p.identify(m) {
				log.Infof("tv_path_metadata: %s is TV", m.SourcePath)
				m.MediaType = tv.TV
			}
		} else {
			log.Debugf("tv_path_metadata: %s category [%s] != video, skipping", m.SourcePath, m.Category)
		}
		if m.MediaType == tv.TV {
			log.Infof("tv_path_metadata: extracting metadata for %v", m)
			m = p.extractMetadata(m)
		} else {
			log.Debugf("tv_path_metadata: %s type [%s] != TV, skipping", m.SourcePath, m.MediaType)
		}
		out <- m
	}
}

func init() {
	processor.Register(processor.Pre, "tv", func() processor.Processor {
		return &TVPreProcessor{
			MatcherStrings: defaultTVMatchers,
			Sanitize:       true,
		}
	})
}
