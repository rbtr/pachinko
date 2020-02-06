/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package config

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rbtr/pachinko/internal/pipeline"
	"github.com/rbtr/pachinko/plugin/input"
	"github.com/rbtr/pachinko/plugin/output"
	"github.com/rbtr/pachinko/plugin/processor"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RootConfig struct {
	DryRun    bool   `mapstructure:"dry-run"`
	LogLevel  string `mapstructure:"log-level"`
	LogFormat string `mapstructure:"log-format"`
}

type ConfigCmdConfig struct {
	RootConfig `mapstructure:",squash"`
	Format     string   `mapstructure:"format"`
	Inputs     []string `mapstructure:"inputs"`
	Outputs    []string `mapstructure:"outputs"`
	Processors []string `mapstructure:"processors"`
}

type SortCmdConfig struct {
	RootConfig `mapstructure:",squash"`
	Pipeline   pipeline.Config          `mapstructure:"pipeline"`
	Inputs     []map[string]interface{} `mapstructure:"inputs"`
	Outputs    []map[string]interface{} `mapstructure:"outputs"`
	Processors []map[string]interface{} `mapstructure:"processors"`
}

func (c *RootConfig) configLogger() {
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
func (c *RootConfig) Validate() error {
	c.configLogger()
	log.Debugf("loaded config: %+v", *c)
	if c.DryRun {
		log.Warn("DRY RUN: no changes will be made")
	}
	return nil
}

func (c *SortCmdConfig) ConfigurePipeline(pipe *pipeline.Pipeline) error {
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
				if err := plugin.Init(); err != nil {
					return err
				}
				pipe.WithInputs(plugin)
			}
		}
	}

	for _, p := range c.Outputs {
		if name, ok := p["name"]; ok {
			if initializer, ok := output.Registry[name.(string)]; ok {
				plugin := initializer()
				if err := mapstructure.Decode(p, plugin); err != nil {
					return err
				}
				if err := plugin.Init(); err != nil {
					return err
				}
				pipe.WithOutputs(plugin)
			}
		}
	}

	for _, p := range c.Processors {
		if name, ok := p["name"]; ok {
			if initializer, ok := processor.Registry[name.(string)]; ok {
				plugin := initializer()
				if err := mapstructure.Decode(p, plugin); err != nil {
					return err
				}
				if err := plugin.Init(); err != nil {
					return err
				}
				pipe.WithProcessors(plugin)
			}
		}
	}

	return nil
}

// LoadConfig loadconfig
func LoadSortCmdConfig() (*SortCmdConfig, error) {
	cfg := &SortCmdConfig{}
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	err := viper.Unmarshal(cfg)
	return cfg, err
}

func (c *ConfigCmdConfig) DefaultConfig(p *SortCmdConfig) error {
	if len(c.Inputs) == 0 && len(c.Outputs) == 0 && len(c.Processors) == 0 {
		// no plugins specified, dump configs for them all
		for k := range input.Registry {
			c.Inputs = append(c.Inputs, k)
		}
		for k := range output.Registry {
			c.Outputs = append(c.Outputs, k)
		}
		for k := range processor.Registry {
			c.Processors = append(c.Processors, k)
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

	for _, name := range c.Processors {
		log.Tracef("making default config for plugin %s", name)
		if initializer, ok := processor.Registry[name]; ok {
			var out map[string]interface{}
			if err := mapstructure.Decode(initializer(), &out); err != nil {
				return err
			}
			out["name"] = name
			p.Processors = append(p.Processors, out)
		}
	}

	return nil
}

func LoadConfigCmdConfig() (*ConfigCmdConfig, error) {
	cfg := &ConfigCmdConfig{}
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	err := viper.Unmarshal(cfg)
	return cfg, err
}
