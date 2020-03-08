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

	"github.com/mitchellh/mapstructure"
	"github.com/rbtr/pachinko/plugin/input"
	"github.com/rbtr/pachinko/plugin/output"
	"github.com/rbtr/pachinko/plugin/processor"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Genconf struct {
	Root       `mapstructure:",squash"`
	Format     string                      `mapstructure:"format"`
	Inputs     []string                    `mapstructure:"inputs"`
	Outputs    []string                    `mapstructure:"outputs"`
	Processors map[processor.Type][]string `mapstructure:"processors"`
}

func (c *Genconf) DefaultConfig(p *Sort) error {
	if len(c.Inputs) == 0 && len(c.Outputs) == 0 && len(c.Processors) == 0 {
		// no plugins specified, dump configs for them all
		for k := range input.Registry {
			c.Inputs = append(c.Inputs, k)
		}
		for k := range output.Registry {
			c.Outputs = append(c.Outputs, k)
		}
		for _, t := range []processor.Type{processor.Pre, processor.Post, processor.Intra} {
			for k := range processor.Registry[t] {
				c.Processors[t] = append(c.Processors[t], k)
			}
		}
	}

	for _, name := range c.Inputs {
		log.Tracef("making default config for plugin %s", name)
		if initializer, ok := input.Registry[name]; ok {
			var out map[string]interface{}
			if err := mapstructure.Decode(initializer(), &out); err != nil {
				return err
			}
			out["name"] = name
			p.Inputs = append(p.Inputs, out)
		}
	}

	for _, name := range c.Outputs {
		log.Tracef("making default config for plugin %s", name)
		if initializer, ok := output.Registry[name]; ok {
			var out map[string]interface{}
			if err := mapstructure.Decode(initializer(), &out); err != nil {
				return err
			}
			out["name"] = name
			p.Outputs = append(p.Outputs, out)
		}
	}

	for t, names := range c.Processors {
		for _, name := range names {
			log.Tracef("making default config for plugin %s", name)
			if initializer, ok := processor.Registry[t][name]; ok {
				var out map[string]interface{}
				if err := mapstructure.Decode(initializer(), &out); err != nil {
					return err
				}
				out["name"] = name
				p.Processors[t] = append(p.Processors[t], out)
			}
		}
	}

	return nil
}

func LoadGenconf(ctx context.Context) (*Genconf, error) {
	cfg := &Genconf{}
	cfg.ctx = ctx
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	err := viper.Unmarshal(cfg)
	return cfg, err
}
