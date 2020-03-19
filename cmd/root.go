/*
 * Copyright (c) 2020 Red Hat, Inc.
 */

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/robnester-rh/afore/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cfgFile string
var Verbose bool
var prowConfigPath, jobConfigPath, periodicJobName string

// run represents the func to execute for the Run step of a command
var run = func(cmd *cobra.Command, args []string) {
	for configType, path := range map[string]string{"prow-config-path": prowConfigPath, "job-config-path": jobConfigPath} {
		ok, err := utils.ValidateConfigPath(path)
		if !ok {
			utils.Logger.Errorw(fmt.Sprintf("Error validating path for %s, cannot continue", configType), configType, path, "error", err)
			os.Exit(-1)
		}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = utils.GenerateCommand("afore", run)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.afore.yaml)")
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "enable debug logging")

	rootCmd.Flags().StringVar(&prowConfigPath, "prow-config-path", "", "Path to the Prow config file")
	rootCmd.Flags().StringVar(&jobConfigPath, "job-config-path", "", "Path to the Prow job config file")
	rootCmd.Flags().StringVar(&periodicJobName, "job-name", "", "Name of the Periodic job to manually trigger")
	// need to add bundleImagespec location string path
	// indexIMage string vavlue
	// openshift version (validate for being greater than [4.x])

	requiredFlags := []string{"prow-config-path", "job-config-path", "job-name"}
	for _, flag := range requiredFlags {
		_ = rootCmd.MarkFlagRequired(flag)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			utils.Logger.Error(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".afore" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".afore")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		utils.Logger.Debugw("Using specified config file: ", "filename", viper.ConfigFileUsed())
	}
}

func initLogging() {
	logWriterSyncer, err := utils.GetLogWriter()
	if err != nil {
		os.Exit(-1)
	}
	logEncoder := utils.GetLogEncoder()
	consoleEncoder := utils.GetConsoleEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(logEncoder, logWriterSyncer, zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, os.Stdout, zap.ErrorLevel),
	)
	utils.Logger = zap.New(core).Sugar()
	if Verbose {
		core = zapcore.NewTee(
			zapcore.NewCore(logEncoder, logWriterSyncer, zap.DebugLevel),
			zapcore.NewCore(consoleEncoder, os.Stdout, zap.DebugLevel),
		)
		utils.Logger = zap.New(core, zap.AddCaller()).Sugar()
	}
}
