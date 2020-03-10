/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"github.com/rbtr/pachinko/internal/config"
	"github.com/rbtr/pachinko/internal/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sort represents the sort command
var sort = &cobra.Command{
	Use:   "sort",
	Short: "Run the sorting pipeline.",
	Long: `
Use this command to execute the sorting pipeline.

With no arguments, sort will load the config from $HOME/.pachinko.yaml.
  $ pachinko sort

If no config is provided, no plugins will be loaded and the pipeline will
not do anything useful.
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.TraceLevel)
		sortConf, err := config.LoadSort(rootCtx)
		if err != nil {
			log.Fatal(err)
		}
		if err := sortConf.Validate(); err != nil {
			log.Fatal(err)
		}

		p := pipeline.NewPipeline()
		if err := sortConf.ConfigurePipeline(p); err != nil {
			log.Fatal(err)
		}

		if err := p.Run(rootCtx); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	root.AddCommand(sort)
}
