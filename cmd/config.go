/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
	"github.com/rbtr/pachinko/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate pachinko configs",
	Long: `
Use this command to generate pachinko configs.

With no arguments, the config generated will contain the default
configs for every compiled plug-in.
  $ pachinko config > config.yaml

Note: the generated config is always in alphabetical order of keys.
	  If order matters to pipeline plug-in execution, it will need
	  to be reordered after generation.
  
The config can be output as either yaml (default) or toml.
  $ pachinko config -o toml > config.toml
  
To only generate stubs for a subset of plug-ins, pass the plug-in
names as a comma separated list to the flag of their type.
  $ pachinko config --inputs=path --processors=tvid,movid > config.yaml

The common flags (dry-run, logging) will be automatically set in the 
output config when they are used on the config command.
  $ pachinko config --log-level=debug > config.yaml

The config file can then be edited and to fully customize the plug-ins
and the pachinko pipeline.
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.TraceLevel)
		cfg, err := config.LoadCmdConfig()
		if err != nil {
			log.Fatal(err)
		}
		if err := cfg.Validate(); err != nil {
			log.Fatal(err)
		}
		pipelineCfg := config.NewCmdSort()
		if err := cfg.DefaultConfig(pipelineCfg); err != nil {
			log.Fatal(err)
		}

		var out map[string]interface{}
		if err := mapstructure.Decode(pipelineCfg, &out); err != nil {
			log.Fatal(err)
		}

		buf := new(bytes.Buffer)
		switch cfg.Format {
		case "toml":
			if err := toml.NewEncoder(buf).Encode(out); err != nil {
				log.Fatal(err)
			}
		default:
			if err := yaml.NewEncoder(buf).Encode(out); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println(buf)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("format", "o", "yaml", "config output format")
	configCmd.Flags().StringSlice("inputs", []string{}, "comma-separated list of input plugins")
	configCmd.Flags().StringSlice("outputs", []string{}, "comma-separated list of output plugins")
	configCmd.Flags().StringSlice("processors", []string{}, "comma-separated list of processor plugins")
	if err := viper.BindPFlags(configCmd.Flags()); err != nil {
		log.Fatal(err)
	}
}
