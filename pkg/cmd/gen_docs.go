package cmd

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var genDocsCmd = &cobra.Command{
	Use:    "gen-docs",
	Hidden: true,
	Short:  "Generates Markdown docs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return GenDocs()
	},
}

const fmTemplate = `---
title: "%s"
hide_title: true
---
`

func GenDocs() error {
	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1))
	}

	err := doc.GenMarkdownTreeCustom(rootCmd, "./docs", filePrepender, func(s string) string {
		return s
	})
	return err
}

func init() {
	rootCmd.AddCommand(genDocsCmd)
}
