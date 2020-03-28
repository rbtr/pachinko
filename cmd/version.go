/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

// version command
var version = &cobra.Command{
	Use:   "version",
	Short: "the pachinko version",
	Long: `
Outputs the version of Pachinko.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	root.AddCommand(version)
}
