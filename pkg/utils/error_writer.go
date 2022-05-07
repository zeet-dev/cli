package utils

import (
	"os"

	"github.com/fatih/color"
)

// ErrorWriter is an io.Writer interface which prints to os.Stderr in red
type ErrorWriter struct{}

func (w *ErrorWriter) Write(s []byte) (int, error) {
	return color.New(color.FgRed).Fprintln(os.Stderr, string(s[:]))
}
