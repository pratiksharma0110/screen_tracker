package utils

import (
	"fmt"
	"os"
)

func GetEnv(env string) string {
	value := os.Getenv(env)
	if value == "" {
		fmt.Fprintf(os.Stderr, "Environment variable %s not set\n", env)
		os.Exit(1)
	}
	return value
}
