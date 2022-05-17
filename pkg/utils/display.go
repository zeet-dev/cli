package utils

import (
	"fmt"
)

// DisplayArray turn an array into a multi-line string containing its items
func DisplayArray(arr []string) (out string) {
	for _, s := range arr {
		out += fmt.Sprintf("- https://%s\n", s)
	}

	return
}

func DisplayMap(m map[string]string) (out string) {
	for k, v := range m {
		out += fmt.Sprintf("%s=%v", k, v)
	}

	return
}
