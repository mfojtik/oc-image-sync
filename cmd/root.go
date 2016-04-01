// Copyright © 2016 Michal Fojtik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mfojtik/oc-image-sync/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "oc-image-sync",
	Short: "Synchronize Docker image between registries/image streams",
	Long: `Use this command to synchronize Docker image between multiple OpenShift
image streams.

Examples:

# Export the image into a tar file from image stream:
$ oc-image-sync export openshift/ruby:2.0 > ruby.tar

# Import the image into custom namespace:
$ oc-image-sync import joe/ruby:2.0 ./ruby.tar

# Move image between two OpenShift datacenters:
$ oc-image-sync export foo/bar:1.0 | oc-image-sync import foo/bar:1.0 --config prod.conf
`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if _, err := exec.LookPath("oc"); err != nil {
		fmt.Println("ERROR: The OpenShift 'oc' binary must be installed in the PATH.")
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&config.User.Token, "token", "", "OpenShift OAuth access token (obtain using: '$ oc whoami -t')")
	RootCmd.PersistentFlags().StringVar(&config.User.Path, "config", "", "Path to OpenShift config file (eg: './user.kubeconfig')")
	RootCmd.PersistentFlags().StringVar(&config.User.ServerAddress, "address", "", "Docker registry address")
	RootCmd.PersistentFlags().StringVar(&config.User.Username, "user", "", "Docker registry username")
	RootCmd.PersistentFlags().StringVar(&config.User.Password, "password", "", "Docker registry password")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".oc-image-sync") // name of config file (without extension)
	viper.AddConfigPath("$HOME")          // adding home directory as first search path
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
