package oc

import "strings"

func GetDockerImageReference(isTagName string) (string, error) {
	var (
		namespace string
		err       error
	)
	parts := strings.SplitN(isTagName, "/", 2)
	if len(parts) == 2 {
		namespace, isTagName = parts[0], parts[1]
	} else {
		namespace, err = GetCurrentProjectName()
		if err != nil {
			return "", err
		}
	}
	return NewCommand("get", "istag", isTagName).
		Namespace(namespace).
		Template("{{ .image.dockerImageReference }}").
		String()
}
