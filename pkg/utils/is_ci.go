package utils

import "os"

func IsCI() bool {
	return os.Getenv("CI") != ""
}
