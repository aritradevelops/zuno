package utils

import (
	"embed"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

func CreateFromTemplate(templateFs embed.FS, templateName string, filePath string, data any) error {
	tmplContent, err := templateFs.ReadFile(templateName)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath = path.Join(wd, filePath)
	tmpl, err := template.New(filePath).Parse(string(tmplContent))
	if err != nil {
		return err
	}
	// Ensure parent directories exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func CopyDirFromEmbed(templateFs embed.FS, srcDir, dstDir string) error {
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

func CloneTemplates(templateFs embed.FS, rootDir string, destDir string, data any, directCloneDirs ...string) error {
	return fs.WalkDir(templateFs, rootDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip directories and clone them directly
		if d.IsDir() {
			if slices.Contains(directCloneDirs, filepath.Base(p)) {
				if err := CopyDirFromEmbed(templateFs, p, d.Name()); err != nil {
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

		relPath := strings.TrimPrefix(p, rootDir)
		outPath := strings.TrimSuffix(filepath.Join(destDir, relPath), ".gotmpl")
		return CreateFromTemplate(templateFs, p, outPath, data)
	})
}

func RunCmd(bin string, verbose bool, args ...string) error {
	cmd := exec.Command(bin, args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
