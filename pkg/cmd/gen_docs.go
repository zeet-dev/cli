package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func createGenDocsCmd() *cobra.Command {
	var dir string

	genDocsCmd := &cobra.Command{
		Use:    "gen-docs",
		Hidden: true,
		Short:  "Generates Markdown docs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return GenDocs(dir)
		},
	}

	genDocsCmd.Flags().StringVarP(&dir, "dir", "d", "./docs", "The directory to place the markdown docs in")

	return genDocsCmd
}

const fmTemplate = `---
title: "%s"
hide_title: true
---
`

func GenDocs(dir string) error {
	if err := writeDocs(rootCmd, dir); err != nil {
		return err
	}

	return nil
}

func writeDocs(cmd *cobra.Command, dir string) error {
	for _, c := range cmd.Commands() {
		fmt.Println(c.Name())
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() || c.Hidden {
			continue
		}
		if err := writeDocs(c, dir); err != nil {
			return err
		}
	}

	generated, err := genDoc(cmd)
	if err != nil {
		return err
	}

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".md"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(generated.Bytes()); err != nil {
		return err
	}

	return nil
}

// Code from cobra.GenMarkdown but without "See Also" section
func genDoc(cmd *cobra.Command) (*bytes.Buffer, error) {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString(fmt.Sprintf(fmTemplate, name))

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("### Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd, name); err != nil {
		return buf, err
	}

	return buf, nil
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("### Options\n\n```\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Options inherited from parent commands\n\n```\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(createGenDocsCmd())
}
