/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"github.com/rbtr/pachinko/internal/config"
	internaltrakt "github.com/rbtr/pachinko/internal/trakt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var traktFile string

// trakt represents the config command
var trakt = &cobra.Command{
	Use:   "trakt",
	Short: "Connect to Trakt",
	Long: `
Use this command to connect Pachinko to Trakt.

Trakt requires that you connect and authorize Pachinko before it can make any
changes on your behalf.

This command will print a URL and code. Open the URL in your browser, sign in
to Trakt (if you aren't signed in already), and enter the code.

Pachinko will write the authorized credentials out to the file specified in
the "--authfile" flag. The credential is stored in a JSON notation:
	{
		"access-token": "[access-token]",
		"client-id": "76a0c1e8d3331021f6e312115e27fe4c29f4ef23ef89a0a69143a62d136ab994",
		"client-secret": "fe8d1f0921413028f92428d2922e13a728e27d2f35b26e315cf3dde31228568d",
		"created": "[created-at]",
		"expires": "7776000",
		"refresh-token": "[refresh-token"
	}

This credential file is portable and can be moved around with your
Pachinko install. Pachinko will automatically refresh it every 80 days during
normal operations.

The token expires after 90 days. If Pachinko can't refresh the token before
it expires, you will need to rerun this to generate a new authorization.
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.TraceLevel)
		cfg, err := config.LoadTrakt(rootCtx)
		if err != nil {
			log.Fatal(err)
		}
		if err := cfg.Validate(); err != nil {
			log.Fatal(err)
		}
		auth, err := internaltrakt.ReadAuthFile(cfg.Authfile)
		if err != nil {
			log.Fatal(err)
		}
		client, err := internaltrakt.NewTrakt(auth)
		if err != nil {
			log.Fatal(err)
		}
		if auth, err = client.Authorize(rootCtx); err != nil {
			log.Fatal(err)
		}
		if err := internaltrakt.WriteAuthFile(cfg.Authfile, auth); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	root.AddCommand(trakt)
	trakt.Flags().StringVar(&traktFile, "authfile", internaltrakt.DefaultAuthfile, "where to save the trakt authorization credential")
	trakt.Flags().BoolP("overwrite", "f", false, "overwrite the authfile if it exists already")
	if err := viper.BindPFlags(trakt.Flags()); err != nil {
		log.Fatal(err)
	}
}
