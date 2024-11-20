package gp

import (
	"fmt"
	"os"

	"github.com/cpunion/go-python/internal/env"
)

var ProjectRoot string

func init() {
	if ProjectRoot == "" {
		panic("ProjectRoot is not set, compile with -ldflags '-X github.com/cpunion/go-python.ProjectRoot=/path/to/project/.deps'")
	}
	envs, err := env.ReadEnv(ProjectRoot)
	if err != nil {
		panic(fmt.Sprintf("Failed to read env: %s", err))
	}
	for key, value := range envs {
		os.Setenv(key, value)
	}
}
