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

	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	"github.com/rbtr/pachinko/types/metadata/tv"
	log "github.com/sirupsen/logrus"
)

var defaultMovieMatchers = []string{
	`(?i)\b([\s\w.-]*)[\s.\/-]+[\s.\/-]?(?:[\s\(.\/-]?(\d{4})[\s\).\/-]?)(?:\/|.[A-Za-z]{3})`, // matches "Name (YEAR)."
}

type MoviePreProcessor struct {
	MatcherStrings []string `mapstructure:"matchers"`
	Sanitize       bool     `mapstructure:"sanitize-name"`

	matchers []*regexp.Regexp
}

func (p *MoviePreProcessor) Init() error {
	for _, str := range p.MatcherStrings {
		r := regexp.MustCompile(str)
		p.matchers = append(p.matchers, r)
	}
	return nil
}

// extract uses the Movie regexp to extract metadata from the input
func (p *MoviePreProcessor) extractMetadata(m types.Media) types.Media {
	var title, year string
	for _, matcher := range p.matchers {
		subs := matcher.FindAllStringSubmatch(m.SourcePath, -1)
		if len(subs) == 0 {
			continue
		}
		if matches := subs[len(subs)-1]; matches != nil {
			if title == "" && strings.TrimSpace(matches[1]) != "" {
				log.Tracef("movie_path_metadata: %v extracting title %s from %s", matcher.String(), matches[1], m.SourcePath)
				title = strings.TrimSpace(matches[1])
				if p.Sanitize {
					title = sanitizer.ReplaceAllString(title, " ")
				}
			}
			if year == "" && strings.TrimSpace(matches[2]) != "" {
				log.Tracef("movie_path_metadata: %v extracting year %s from %s", matcher.String(), matches[2], m.SourcePath)
				year = strings.TrimSpace(matches[2])
			}
			if title != "" && year != "" {
				break
			}
		}
	}
	m.MovieMetadata.Title = title
	m.MovieMetadata.ReleaseYear, _ = strconv.Atoi(year)
	return m
}

// identify tests if the input is matched by any of the Movie regexp
func (p *MoviePreProcessor) identify(m types.Media) bool {
	for _, matcher := range p.matchers {
		if matcher.MatchString(m.SourcePath) {
			log.Tracef("movie_path_metadata: regexp %s matched %s", matcher, m.SourcePath)
			return true
		}
	}
	return false
}

func (p *MoviePreProcessor) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started movie_path_metadata processor")
	for m := range in {
		log.Tracef("movie_path_metadata: received input: %v", m)
		if m.Category != types.Video {
			log.Debugf("movie_path_metadata: %s category %s != video, skipping", m.SourcePath, m.Category)
			continue
		}
		if !p.identify(m) {
			log.Debugf("movie_path_metadata: %s type %s != TV, skipping", m.SourcePath, m.Type)
			continue
		}
		m.Type = tv.TV
		log.Debugf("movie_path_metadata: received input %v", m)
		out <- p.extractMetadata(m)
	}
}

func init() {
	processor.Register(processor.Pre, "movie", func() processor.Processor {
		return &MoviePreProcessor{
			MatcherStrings: defaultMovieMatchers,
			Sanitize:       true,
		}
	})
}
