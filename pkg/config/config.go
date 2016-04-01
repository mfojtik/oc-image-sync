package config

import "github.com/fsouza/go-dockerclient"

type Config struct {
	Path  string
	Token string
	docker.AuthConfiguration
}

var User = &Config{}
