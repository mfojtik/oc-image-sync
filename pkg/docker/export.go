package docker

import (
	"os"

	dockerclient "github.com/fsouza/go-dockerclient"
)

func ExportImage(name string) error {
	exportOpts := dockerclient.ExportImageOptions{
		Name:         name,
		OutputStream: os.Stdout,
	}
	return Client().ExportImage(exportOpts)
}
