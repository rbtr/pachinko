/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package config

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Root struct {
	// nolint: structcheck
	ctx       context.Context
	DryRun    bool   `mapstructure:"dry-run"`
	LogLevel  string `mapstructure:"log-level"`
	LogFormat string `mapstructure:"log-format"`
}

func (c *Root) configLogger() {
	log.SetLevel(log.InfoLevel)
	switch c.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
	if c.LogLevel != "" {
		if lvl, err := log.ParseLevel(c.LogLevel); err != nil {
			log.Fatal(err)
		} else {
			log.SetLevel(lvl)
		}
	}
}

// Validate validate
func (c *Root) Validate() error {
	c.configLogger()
	log.Debugf("loaded config: %+v", *c)
	if c.DryRun {
		log.Warn("DRY RUN: no changes will be made")
	}
	return nil
}
