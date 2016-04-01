// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	"github.com/mfojtik/oc-image-sync/pkg/docker"
	"github.com/mfojtik/oc-image-sync/pkg/oc"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import the Docker image to the OpenShift image stream",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(-1)
		}
		user, password, email, err := oc.GetDockerAuth()
		if err != nil {
			fmt.Printf("ERROR: Unable to get Docker image reference for image stream %q: %v", args[0], err)
			os.Exit(1)
		}
		if err := docker.Login(oc.GetRegistryHostFromImage(ref), user, password, email); err != nil {
			fmt.Printf("ERROR: Unable to login to Docker: %v", err)
			os.Exit(1)
		}
		if err := docker.ImportImage(os.Stdin, ref); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to import image stream %q: %v", args[0], err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
