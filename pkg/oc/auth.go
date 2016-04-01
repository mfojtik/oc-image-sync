package oc

import "fmt"

func GetToken() (string, error) {
	return NewCommand("whoami", "-t").String()
}

func GetCurrentUser() (string, error) {
	return NewCommand("whoami").String()
}

func GetDockerAuth() (user, password, email string, err error) {
	user, err = GetCurrentUser()
	if err != nil {
		return
	}
	password, err = GetToken()
	if err != nil {
		return
	}
	email = fmt.Sprintf("%s@openshift.local", user)
	return
}

func GetCurrentProjectName() (string, error) {
	return NewCommand("project", "-q").String()
}
