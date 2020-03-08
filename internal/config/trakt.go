/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package config

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Trakt struct {
	Root      `mapstructure:",squash"`
	Authfile  string `mapstructure:"authfile"`
	Overwrite bool   `mapstructure:"overwrite"`
}

func LoadTrakt(ctx context.Context) (*Trakt, error) {
	cfg := &Trakt{}
	cfg.ctx = ctx
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	err := viper.Unmarshal(cfg)
	return cfg, err
}

func (t *Trakt) Validate() error {
	if t.Authfile == "" {
		return errors.New("authfile must be set")
	}
	return nil
}
