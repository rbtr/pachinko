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

// genconf represents the config command.
var genconf = &cobra.Command{
	Use:   "genconf",
	Short: "Generate pachinko configs",
	Long: `
Use this command to generate pachinko configs.

With no arguments, the config generated will contain the default
configs for every compiled plug-in.
  $ pachinko genconf > config.yaml

Note: the generated config is always in alphabetical order of keys.
	  If order matters to pipeline plug-in execution, it will need
	  to be reordered after generation.
  
The config can be output as either yaml (default) or toml.
  $ pachinko genconf -o toml > config.toml
  
To only generate stubs for a subset of plug-ins, pass the plug-in
names as a comma separated list to the flag of their type.
  $ pachinko genconf --inputs=path --processors=tvid,movid > config.yaml

The common flags (dry-run, logging) will be automatically set in the 
output config when they are used on the config command.
  $ pachinko genconf --log-level=debug > config.yaml

The config file can then be edited and to fully customize the plug-ins
and the pachinko pipeline.
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.TraceLevel)
		cfg, err := config.LoadGenconf(rootCtx)
		if err != nil {
			log.Fatal(err)
		}
		if err := cfg.Validate(); err != nil {
			log.Fatal(err)
		}
		sortCfg := config.NewSort(rootCtx)
		if err := cfg.DefaultConfig(sortCfg); err != nil {
			log.Fatal(err)
		}

		var out map[string]interface{}
		if err := mapstructure.Decode(sortCfg, &out); err != nil {
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
	root.AddCommand(genconf)
	genconf.Flags().StringP("format", "o", "yaml", "config output format")
	genconf.Flags().StringSlice("inputs", []string{}, "comma-separated list of input plugins")
	genconf.Flags().StringSlice("outputs", []string{}, "comma-separated list of output plugins")
	genconf.Flags().StringSlice("processors", []string{}, "comma-separated list of processor plugins")
	if err := viper.BindPFlags(genconf.Flags()); err != nil {
		log.Fatal(err)
	}
}
