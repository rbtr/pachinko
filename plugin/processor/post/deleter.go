/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package post

import (
	"context"
	"path"
	"regexp"
	"strings"

	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

type Deleter struct {
	Categories     []string `mapstructure:"categories"`
	Extensions     []string `mapstructure:"extensions"`
	Directories    bool     `mapstructure:"directories"`
	MatcherStrings []string `mapstructure:"matchers"`

	matchers []*regexp.Regexp
}

func (p *Deleter) Init(context.Context) error {
	log.Trace("deleter: initializing")
	for _, str := range p.MatcherStrings {
		r := regexp.MustCompile(str)
		p.matchers = append(p.matchers, r)
	}
	log.Tracef("deleter: initialized %d matchers, %d exts, directories = %t", len(p.matchers), len(p.Extensions), p.Directories)
	return nil
}

func (p *Deleter) matchRegexps(s string) bool {
	for _, matcher := range p.matchers {
		if matcher.MatchString(s) {
			return true
		}
	}
	return false
}

func (p *Deleter) matchExts(s string) bool {
	for _, ext := range p.Extensions {
		if s == ext {
			return true
		}
	}
	return false
}

func (p *Deleter) shouldDelete(m types.Item) bool {
	if m.FileType == types.Directory && p.Directories {
		return true
	}
	if p.matchExts(strings.Trim(path.Ext(m.SourcePath), ".")) {
		return true
	}
	if p.matchRegexps(m.SourcePath) {
		return true
	}
	return false
}

func (p *Deleter) Process(in <-chan types.Item, out chan<- types.Item) {
	log.Trace("started deleter processor")
	for m := range in {
		log.Tracef("deleter: received input %#v", m)
		if p.shouldDelete(m) {
			log.Infof("deleter: marking %s for delete", m.SourcePath)
			m.Delete = true
		}
		out <- m
	}
}

func init() {
	processor.Register(processor.Post, "deleter", func() processor.Processor {
		var defaultExtensions = []string{}
		defaultExtensions = append(defaultExtensions, types.ArchiveExtensions...)
		defaultExtensions = append(defaultExtensions, types.ExecutableExtensions...)
		defaultExtensions = append(defaultExtensions, types.ImageExtensions...)
		defaultExtensions = append(defaultExtensions, types.SubtitleExtensions...)
		defaultExtensions = append(defaultExtensions, types.TextExtensions...)
		return &Deleter{
			Extensions:     defaultExtensions,
			Directories:    true,
			MatcherStrings: []string{},
		}
	})
}
