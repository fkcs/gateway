package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var projectRootEnv string

func SetProjectRootEnv(env string) {
	projectRootEnv = env
}

func Root() string {
	var root string
	root, ok := os.LookupEnv(projectRootEnv)
	if !ok || len(root) == 0 {
		errMsg := fmt.Sprintf("env %s is not set", projectRootEnv)
		panic(errMsg)
	}
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	return root
}

func AbsPath(relPath string) string {
	return filepath.ToSlash(Root() + relPath)
}
