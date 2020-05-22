/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package types

import "regexp"

// AudioChannels regexp constants.
var AudioChannels = map[string]*regexp.Regexp{
	"2.0": regexp.MustCompile(`2\.0`),
	"5.1": regexp.MustCompile(`5\.1`),
	"7.1": regexp.MustCompile(`7\.1`),
}

// AudioFormats regexp constants.
var AudioFormats = map[string]*regexp.Regexp{
	"aac": regexp.MustCompile("aac"),
}

// ColorFormats regexp constants.
var ColorFormats = map[string]*regexp.Regexp{
	"8 bit":  regexp.MustCompile(`8.bit`),
	"10 bit": regexp.MustCompile(`10.bit`),
}

// Resolutions regexp constants.
var Resolutions = map[string]*regexp.Regexp{
	"1080p": regexp.MustCompile(`\b1080p?`),
	"720p":  regexp.MustCompile(`\b720p?`),
	"480p":  regexp.MustCompile(`\b480p?`),
}

// Sources regexp constants.
var Sources = map[string]*regexp.Regexp{
	"bluray": regexp.MustCompile("bluray"),
	"dvd":    regexp.MustCompile("dvd"),
	"hdtv":   regexp.MustCompile("hdtv"),
}

// TVSeason regexp constants.
var TVSeason = map[string]*regexp.Regexp{
	"season": regexp.MustCompile("season|series"),
}

// VideoFormats regexp constants.
var VideoFormats = map[string]*regexp.Regexp{
	"hevc":  regexp.MustCompile("hevc"),
	"h.264": regexp.MustCompile(`h\.?264`),
	"h.265": regexp.MustCompile(`h\.?265`),
	"mov":   regexp.MustCompile(`\bmov\b`),
	"mpeg":  regexp.MustCompile("mpeg"),
	"x264":  regexp.MustCompile(`x\.?264`),
	"x265":  regexp.MustCompile(`x\.?265`),
}
