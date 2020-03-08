/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

/*
Package input provides an interface for, and implementations of,
input plugins which feed data from various sources in to the
pipeline datastream.
*/
package input

import (
	"context"

	"github.com/rbtr/pachinko/types"
	log "github.com/sirupsen/logrus"
)

// input input
type Input interface {
	Consume(chan<- types.Media)
	Init(context.Context) error
}

var Registry map[string](func() Input) = map[string](func() Input){}

func Register(name string, initializer func() Input) {
	if _, ok := Registry[name]; ok {
		log.Fatalf("input registry already contains plugin named %s", name)
	}
	Registry[name] = initializer
}
