/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var excludeDirs = []string{"tmp", "docs"}

// templatizeCmd represents the templatize command
var templatizeCmd = &cobra.Command{
	Use:   "templatize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		exampleDir := filepath.Join(wd, "example")
		dumpDir := filepath.Join(wd, "cmd/templates/base")

		err = filepath.WalkDir(exampleDir, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() && slices.Contains(excludeDirs, filepath.Base(filePath)) {
				return filepath.SkipDir
			}

			// compute path relative to example/
			relPath, err := filepath.Rel(exampleDir, filePath)
			if err != nil {
				return err
			}

			targetPath := filepath.Join(dumpDir, relPath)

			if d.IsDir() {
				return os.MkdirAll(targetPath, 0755)
			}

			// read file
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			// replace goserve -> template var
			updated := strings.ReplaceAll(
				string(content),
				"goserve",
				"{{.PackageName}}",
			)

			// ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}

			// write as .gotmpl
			return os.WriteFile(
				targetPath+".gotmpl",
				[]byte(updated),
				0644,
			)
		})

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(templatizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templatizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templatizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
