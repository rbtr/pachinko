/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package plugin

import (
	_ "github.com/rbtr/pachinko/plugin/input/path"
	_ "github.com/rbtr/pachinko/plugin/output/move"
	_ "github.com/rbtr/pachinko/plugin/output/stdout"
	_ "github.com/rbtr/pachinko/plugin/processor/tvdb"
	_ "github.com/rbtr/pachinko/plugin/processor/tvdest"
	_ "github.com/rbtr/pachinko/plugin/processor/tvid"
	_ "github.com/rbtr/pachinko/plugin/processor/tvmeta"
	_ "github.com/rbtr/pachinko/plugin/processor/vidid"
)
