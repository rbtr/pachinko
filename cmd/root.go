/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// include plugins
	_ "github.com/rbtr/pachinko/plugin"
)

var cfgFile string
var rootCtx context.Context

// root represents the base command when called without any subcommands.
var root = &cobra.Command{
	Use: "pachinko",
	Long: `
             _   _     _
 ___ ___ ___| |_|_|___| |_ ___
| . | .'|  _|   | |   | '_| . |
|  _|__,|___|_|_|_|_|_|_,_|___|
|_|

pluggable media sorter`,
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(rootConfig)

	// bind root flags
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pachinko.yaml)")
	root.PersistentFlags().Bool("dry-run", false, "run pipeline as read only and do not make changes")
	root.PersistentFlags().StringP("log-level", "v", "info", "log verbosity (trace,debug,info,warn,error)")
	root.PersistentFlags().String("log-format", "text", "log format (text,json)")
	if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
		log.Fatal(err)
	}

	// init signal channels
	var cancel context.CancelFunc
	rootCtx, cancel = context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Debug("caught exit sigal, exiting")
		cancel()
		os.Exit(0)
	}()
}

func rootConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".pachinko")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
