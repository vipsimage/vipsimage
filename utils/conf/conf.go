package conf

import (
	"os"
)

func Getenv(name, def string) string {
	env := os.Getenv(name)
	if env == "" {
		return def
	}

	return env
}
