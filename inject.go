package gp

import (
	"fmt"
	"os"

	"github.com/cpunion/go-python/internal/env"
)

var ProjectRoot string

func init() {
	injectDebug := os.Getenv("GP_INJECT_DEBUG")
	if ProjectRoot == "" {
		if injectDebug != "" {
			panic("ProjectRoot is not set, compile with -ldflags '-X github.com/cpunion/go-python.ProjectRoot=/path/to/project'")
		}
		return
	}
	envs, err := env.ReadEnv(ProjectRoot)
	if err != nil {
		panic(fmt.Sprintf("Failed to read env: %s", err))
	}
	if injectDebug != "" {
		fmt.Fprintf(os.Stderr, "Injecting envs for project: %s\n", ProjectRoot)
		for key, value := range envs {
			fmt.Fprintf(os.Stderr, "  %s=%s\n", key, value)
		}
		fmt.Fprintf(os.Stderr, "End of envs\n")
	}
	for key, value := range envs {
		fmt.Fprintf(os.Stderr, "Injecting env: %s=%s\n", key, value)
		os.Setenv(key, value)
	}
}
