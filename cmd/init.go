package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"zuno/pkg/logger"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	packageName string
	adapter     string
	transports  []string
)

type TemplateData struct {
	PackageName string
	Adapter     string
}

//go:embed templates/base/*
var templateFs embed.FS

var directCloneDirs = []string{"locales"}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		var fields []huh.Field

		// Only ask if flag NOT provided
		if !cmd.Flags().Changed("package") {
			fields = append(fields,
				huh.NewInput().
					Title("What's the package name gonna be?").
					Placeholder("github.com/<username>/<package-name>").
					Value(&packageName),
			)
		}

		if !cmd.Flags().Changed("adapter") {
			fields = append(fields,
				huh.NewSelect[string]().
					Title("Choose an adapter").
					Options(
						huh.NewOption("mongodb", "mongodb"),
						huh.NewOption("postgres", "postgres"),
					).
					Value(&adapter),
			)
		}

		if !cmd.Flags().Changed("transports") {
			fields = append(fields,
				huh.NewMultiSelect[string]().
					Title("Select transports").
					Options(
						huh.NewOption("http", "http"),
						huh.NewOption("grpc", "grpc"),
						huh.NewOption("ws", "ws"),
					).
					Value(&transports),
			)
		}

		// Run form only if something needs prompting
		if len(fields) > 0 {
			if err := huh.NewForm(
				huh.NewGroup(fields...),
			).Run(); err != nil {
				return err
			}
		}
		defer func() {
			logger.Info().Msgf(`
			Re-run with : zuno init --package=%s \
			--adapter=%s \
			--transports=%s`, packageName, adapter, strings.Join(transports, ","))
		}()

		// At this point ALL values are guaranteed
		// packageName, adapter, transports are set
		// 1️⃣ go mod init
		if err := runGoCmd("mod", "init", packageName); err != nil {
			return err
		}

		// 2️⃣ clone templates
		if err := clone(); err != nil {
			return err
		}

		// 3️⃣ go mod tidy
		if err := runGoCmd("mod", "tidy"); err != nil {
			return err
		}

		// 4️⃣ swag init
		if err := runCmd(
			"swag",
			"init",
			"-g",
			"./internal/transports/http/server.go",
		); err != nil {
			return fmt.Errorf("failed to run swag init: %w", err)
		}
		return nil
	},
}

func clone() error {
	data := TemplateData{
		PackageName: packageName,
		Adapter:     adapter,
	}

	return fs.WalkDir(templateFs, "templates/base", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(p, d.IsDir() && slices.Contains(directCloneDirs, filepath.Base(p)))

		// Skip directories
		if d.IsDir() {
			if slices.Contains(directCloneDirs, filepath.Base(p)) {
				if err := copyDirFromEmbed(p, d.Name()); err != nil {
					return err
				}
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .gotmpl files
		if !strings.HasSuffix(d.Name(), ".gotmpl") {
			return nil
		}

		// Read template file
		content, err := templateFs.ReadFile(p)
		if err != nil {
			return err
		}

		// Strip `.gotmpl`
		relPath := strings.TrimPrefix(p, "templates/base/")
		outPath := strings.TrimSuffix(relPath, ".gotmpl")

		// Parse template
		tpl, err := template.New(d.Name()).Parse(string(content))
		if err != nil {
			return err
		}

		// Ensure parent directories exist
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}

		// Create output file
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Execute template
		return tpl.Execute(f, data)
	})
}

func copyDirFromEmbed(srcDir, dstDir string) error {
	return fs.WalkDir(templateFs, srcDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel := strings.TrimPrefix(p, srcDir)
		out := strings.TrimSuffix(filepath.Join(dstDir, rel), ".gotmpl")

		if d.IsDir() {
			return os.MkdirAll(out, 0o755)
		}

		b, err := templateFs.ReadFile(p)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}

		return os.WriteFile(out, b, 0o644)
	})
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(
		&packageName,
		"package",
		"",
		"Go module package name",
	)

	initCmd.Flags().StringVar(
		&adapter,
		"adapter",
		"",
		"Database adapter (mongodb|postgres)",
	)

	initCmd.Flags().StringSliceVar(
		&transports,
		"transports",
		nil,
		"Transports to enable (http, grpc, ws)",
	)
}

func runGoCmd(args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmd(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
