package oc

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

var ocBinary, _ = exec.LookPath("oc")

type ocCommand struct {
	template string
	cmd      *exec.Cmd
}

func NewCommand(args ...string) *ocCommand {
	return &ocCommand{
		cmd: exec.Command(ocBinary, args...),
	}
}

func (c *ocCommand) Namespace(name string) *ocCommand {
	c.cmd.Args = append(c.cmd.Args, []string{"-n", name}...)
	return c
}

func (c *ocCommand) Template(template string) *ocCommand {
	c.cmd.Args = append(c.cmd.Args, []string{"--template", template, "-o", "go-template"}...)
	return c
}

func (c *ocCommand) String() (string, error) {
	out, err := c.cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func GetRegistryHostFromImage(spec string) string {
	ref, _ := ParseDockerImageReference(spec)
	return ref.RegistryURL().Host
}

// DockerImageReference points to a Docker image.
type DockerImageReference struct {
	Registry  string
	Namespace string
	Name      string
	Tag       string
	ID        string
}

const (
	// DockerDefaultNamespace is the value for namespace when a single segment name is provided.
	DockerDefaultNamespace = "library"
	// DockerDefaultRegistry is the value for the registry when none was provided.
	DockerDefaultRegistry = "docker.io"
	// DockerDefaultV1Registry is the host name of the default v1 registry
	DockerDefaultV1Registry = "index." + DockerDefaultRegistry
	// DockerDefaultV2Registry is the host name of the default v2 registry
	DockerDefaultV2Registry = "registry-1." + DockerDefaultRegistry
)

// TODO remove (base, tag, id)
func parseRepositoryTag(repos string) (string, string, string) {
	n := strings.Index(repos, "@")
	if n >= 0 {
		parts := strings.Split(repos, "@")
		return parts[0], "", parts[1]
	}
	n = strings.LastIndex(repos, ":")
	if n < 0 {
		return repos, "", ""
	}
	if tag := repos[n+1:]; !strings.Contains(tag, "/") {
		return repos[:n], tag, ""
	}
	return repos, "", ""
}

func IsRegistryDockerHub(registry string) bool {
	switch registry {
	case DockerDefaultRegistry, DockerDefaultV1Registry, DockerDefaultV2Registry:
		return true
	default:
		return false
	}
}

// ParseDockerImageReference parses a Docker pull spec string into a
// DockerImageReference.
func ParseDockerImageReference(spec string) (DockerImageReference, error) {
	var ref DockerImageReference
	// TODO replace with docker version once docker/docker PR11109 is merged upstream
	stream, tag, id := parseRepositoryTag(spec)

	repoParts := strings.Split(stream, "/")
	switch len(repoParts) {
	case 2:
		if isRegistryName(repoParts[0]) {
			// registry/name
			ref.Registry = repoParts[0]
			if IsRegistryDockerHub(ref.Registry) {
				ref.Namespace = DockerDefaultNamespace
			}
			if len(repoParts[1]) == 0 {
				return ref, fmt.Errorf("the docker pull spec %q must be two or three segments separated by slashes", spec)
			}
			ref.Name = repoParts[1]
			ref.Tag = tag
			ref.ID = id
			break
		}
		// namespace/name
		ref.Namespace = repoParts[0]
		if len(repoParts[1]) == 0 {
			return ref, fmt.Errorf("the docker pull spec %q must be two or three segments separated by slashes", spec)
		}
		ref.Name = repoParts[1]
		ref.Tag = tag
		ref.ID = id
		break
	case 3:
		// registry/namespace/name
		ref.Registry = repoParts[0]
		ref.Namespace = repoParts[1]
		if len(repoParts[2]) == 0 {
			return ref, fmt.Errorf("the docker pull spec %q must be two or three segments separated by slashes", spec)
		}
		ref.Name = repoParts[2]
		ref.Tag = tag
		ref.ID = id
		break
	case 1:
		// name
		if len(repoParts[0]) == 0 {
			return ref, fmt.Errorf("the docker pull spec %q must be two or three segments separated by slashes", spec)
		}
		ref.Name = repoParts[0]
		ref.Tag = tag
		ref.ID = id
		break
	default:
		return ref, fmt.Errorf("the docker pull spec %q must be two or three segments separated by slashes", spec)
	}

	return ref, nil
}

func (r DockerImageReference) RegistryURL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   r.AsV2().Registry,
	}
}

func (r DockerImageReference) AsV2() DockerImageReference {
	switch r.Registry {
	case DockerDefaultV1Registry, DockerDefaultRegistry:
		r.Registry = DockerDefaultV2Registry
	}
	return r
}

func isRegistryName(str string) bool {
	switch {
	case strings.Contains(str, ":"),
		strings.Contains(str, "."),
		str == "localhost":
		return true
	}
	return false
}
