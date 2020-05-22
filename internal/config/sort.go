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
	"github.com/rbtr/pachinko/internal/pipeline"
	internalout "github.com/rbtr/pachinko/internal/plugin/output"
	internalpre "github.com/rbtr/pachinko/internal/plugin/processor/pre"
	"github.com/rbtr/pachinko/plugin/input"
	"github.com/rbtr/pachinko/plugin/output"
	"github.com/rbtr/pachinko/plugin/processor"
	"github.com/spf13/viper"
)

type Sort struct {
	Root       `mapstructure:",squash"`
	Pipeline   pipeline.Config                             `mapstructure:"pipeline"`
	Inputs     []map[string]interface{}                    `mapstructure:"inputs"`
	Outputs    []map[string]interface{}                    `mapstructure:"outputs"`
	Processors map[processor.Type][]map[string]interface{} `mapstructure:"processors"`
}

func (c *Sort) ConfigurePipeline(pipe *pipeline.Pipeline) error {
	if err := mapstructure.Decode(c.Pipeline, pipe); err != nil {
		return err
	}
	for _, p := range c.Inputs {
		if name, ok := p["name"]; ok {
			if initializer, ok := input.Registry[name.(string)]; ok {
				plugin := initializer()
				if err := mapstructure.Decode(p, plugin); err != nil {
					return err
				}
				if err := plugin.Init(c.ctx); err != nil {
					return err
				}
				pipe.WithInputs(plugin)
			}
		}
	}

	ocfg := output.Config{
		DryRun: c.DryRun,
	}

	for _, p := range c.Outputs {
		if name, ok := p["name"]; ok {
			if initializer, ok := output.Registry[name.(string)]; ok {
				plugin := initializer()
				if err := mapstructure.Decode(p, plugin); err != nil {
					return err
				}
				if err := plugin.Init(c.ctx, ocfg); err != nil {
					return err
				}
				pipe.WithOutputs(plugin)
			}
		}
	}

	deleter := &internalout.Deleter{}
	if err := deleter.Init(c.ctx, ocfg); err != nil {
		return err
	}
	pipe.WithOutputs(deleter)

	categorizer := internalpre.NewCategorizer()
	if err := categorizer.Init(c.ctx); err != nil {
		return err
	}
	pipe.WithProcessors(categorizer)

	for _, t := range processor.Types {
		for _, p := range c.Processors[t] {
			if name, ok := p["name"]; ok {
				if initializer, ok := processor.Registry[t][name.(string)]; ok {
					plugin := initializer()
					if err := mapstructure.Decode(p, plugin); err != nil {
						return err
					}
					if err := plugin.Init(c.ctx); err != nil {
						return err
					}
					pipe.WithProcessors(plugin)
				}
			}
		}
	}

	return nil
}

func NewSort(ctx context.Context) *Sort {
	cfg := &Sort{
		Processors: map[processor.Type][]map[string]interface{}{
			processor.Pre:   {},
			processor.Intra: {},
			processor.Post:  {},
		},
	}
	cfg.ctx = ctx
	return cfg
}

// LoadConfig loadconfig.
func LoadSort(ctx context.Context) (*Sort, error) {
	cfg := NewSort(ctx)
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	err := viper.Unmarshal(cfg)
	return cfg, err
}
