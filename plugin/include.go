/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package plugin

import (
	// blank includes
	_ "github.com/rbtr/pachinko/plugin/input"
	_ "github.com/rbtr/pachinko/plugin/output"
	_ "github.com/rbtr/pachinko/plugin/processor/intra"
	_ "github.com/rbtr/pachinko/plugin/processor/post"
	_ "github.com/rbtr/pachinko/plugin/processor/pre"
)
