package iostreams

import (
	"github.com/fatih/color"
)

type ColorScheme struct{}

func NewColorScheme() *ColorScheme {
	return &ColorScheme{}
}

func (c *ColorScheme) Red(t string) string {
	return color.RedString(t)
}

func (c *ColorScheme) Redf(t string, args ...interface{}) string {
	return color.RedString(t, args)
}

func (c *ColorScheme) Green(t string) string {
	return color.GreenString(t)
}

func (c *ColorScheme) Greenf(t string, args ...interface{}) string {
	return color.GreenString(t, args)
}
