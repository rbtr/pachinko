/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

/*
Package output provides an interface for, and implementations of,
output plugins which consume data from the pipeline datastream at
the end of processing.
*/
package output

import (
	"context"

	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

// Config is common/general output tunables
type Config struct {
	DryRun bool
}

// Output is plugin interface to handle the result
type Output interface {
	Receive(<-chan types.Item)
	Init(context.Context, Config) error
}

var Registry map[string](func() Output) = map[string](func() Output){}

func Register(name string, initializer func() Output) {
	if _, ok := Registry[name]; ok {
		log.Fatalf("output registry already contains plugin named %s", name)
	}
	Registry[name] = initializer
}
