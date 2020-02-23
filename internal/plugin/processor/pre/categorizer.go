/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package pre

import (
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

var defaultCategoryFileExtensions = map[types.Category][]string{
	types.Archive: {
		"7z",
		"gz",
		"gzip",
		"rar",
		"tar",
		"zip",
	},
	types.Image: {
		"bmp",
		"gif",
		"heic",
		"jpeg",
		"jpg",
		"png",
		"tiff",
	},
	types.Subtitle: {
		"srt",
		"sub",
	},
	types.Text: {
		"info",
		"nfo",
		"txt",
		"website",
	},
	types.Video: {
		"avi",
		"divx",
		"m4v",
		"mkv",
		"mov",
		"mp4",
		"xvid",
	},
}

type FileCategorizer struct {
	CategoryFileExtensions  map[types.Category][]string `mapstructure:"file-extensions"`
	fileExtensionCategories map[string]types.Category
}

func (cat *FileCategorizer) Init() error {
	log.Trace("categorizer initializing")
	cat.fileExtensionCategories = map[string]types.Category{}
	// transpose the category/extension map for immediate lookups
	for k, v := range cat.CategoryFileExtensions {
		for _, vv := range v {
			if kk, ok := cat.fileExtensionCategories[vv]; ok {
				return errors.Errorf("duplicate filetype::category mapping: %s::%s already exists as %s::%s, ", vv, k, vv, kk)
			}
			cat.fileExtensionCategories[vv] = k
		}
	}
	log.Trace("categorizer initialized")
	return nil
}

func (cat *FileCategorizer) identify(m types.Media) types.Media {
	ext := path.Ext(m.SourcePath)

	category := types.Unknown
	if ext == "" {
		log.Debug("categorizer: no extension, unknown category")
	}

	trimmed := strings.Trim(ext, ".")
	log.Tracef("categorizer: lookup extension '%s'", trimmed)
	if c, ok := cat.fileExtensionCategories[trimmed]; ok {
		log.Debugf("categorizer: identified %s as %s", ext, c)
		category = c
	}
	m.Category = category
	return m
}

func (cat *FileCategorizer) Process(in <-chan types.Media, out chan<- types.Media) {
	log.Trace("started categorizer")
	for m := range in {
		log.Debugf("categorizer: received input: %v", m)
		out <- cat.identify(m)
	}
}

func (*FileCategorizer) Type() processor.Type {
	return processor.Pre
}

func NewCategorizer() *FileCategorizer {
	return &FileCategorizer{
		CategoryFileExtensions: defaultCategoryFileExtensions,
	}
}
