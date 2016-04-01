package docker

import (
	"fmt"
	"os"
	"os/exec"

	dockerclient "github.com/fsouza/go-dockerclient"
)

func Client() *dockerclient.Client {
	c, err := dockerclient.NewClientFromEnv()
	if err != nil {
		fmt.Printf("ERROR: Unable to connect to Docker daemon: %v\n", err)
		os.Exit(-1)
	}
	return c
}

func Login(registry, user, password, email string) error {
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return err
	}
	var out []byte
	out, err = exec.Command(dockerPath, "login", "-u", user, "-p", password, "-e", email, registry).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v (%q): %q", err, registry, string(out))
	}
	return nil
}
