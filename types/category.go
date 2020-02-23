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
