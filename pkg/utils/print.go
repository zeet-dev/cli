package utils

import (
	"fmt"
)

// DisplayArray turn an array into a multi-line string containing its items
func DisplayArray(arr []string) (out string) {
	for _, s := range arr {
		out += fmt.Sprintf("- %s\n", s)
	}

	return
}
