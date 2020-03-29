/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package types

type Category string

const (
	Archive  Category = "archive"
	Image    Category = "image"
	Subtitle Category = "subtitle"
	Text     Category = "text"
	Unknown  Category = ""
	Video    Category = "video"
)

var ArchiveExtensions = []string{
	"7z",
	"gz",
	"gzip",
	"rar",
	"tar",
	"zip",
}

var ImageExtensions = []string{
	"bmp",
	"gif",
	"heic",
	"jpeg",
	"jpg",
	"png",
	"tiff",
}

var SubtitleExtensions = []string{
	"srt",
	"sub",
}

var TextExtensions = []string{
	"info",
	"nfo",
	"txt",
	"website",
}

var VideoExtensions = []string{
	"avi",
	"divx",
	"m4v",
	"mkv",
	"mov",
	"mp4",
	"xvid",
}
