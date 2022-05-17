package cmdutil

import (
	"io"

	"github.com/fatih/color"
)

// ErrorWriter is an io.Writer interface which prints to os.Stderr in red
type ErrorWriter struct {
	Out io.Writer
}

func (w *ErrorWriter) Write(s []byte) (int, error) {
	return color.New(color.FgRed).Fprintln(w.Out, string(s[:]))
}
