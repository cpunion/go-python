package gp

import (
	"fmt"
	"os"

	"github.com/cpunion/go-python/internal/env"
)

var ProjectRoot string

func init() {
	if ProjectRoot == "" {
		fmt.Fprintf(os.Stderr, "ProjectRoot is not set\n")
		return
	}
	envs, err := env.ReadEnv(ProjectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read env: %s\n", err)
		return
	}
	for key, value := range envs {
		os.Setenv(key, value)
	}
}
