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

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the Docker image from the OpenShift image stream",
	Long: `Use this command to export Docker image from OpenShift image stream

Examples:

# Export the image into a tar file from image stream:
$ oc-image-sync export openshift/ruby:2.0 > ruby.tar

# Export the image and import it to local Docker:
$ oc-image-sync export openshift/ruby:2.0 | docker import -
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(-1)
		}
		ref, err := oc.GetDockerImageReference(args[0])
		if err != nil {
			fmt.Printf("ERROR: Unable to get Docker image reference for image stream %q: %v", args[0], err)
			os.Exit(1)
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
		if err := docker.ExportImage(ref); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to export image stream %q: %v", args[0], err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
