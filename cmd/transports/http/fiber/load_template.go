package fiber

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed templates/*
var templates embed.FS

// load template loads the content of a template
// TODO: enable overwriting
func loadTemplate(identifier string) (string, error) {
	f, err := templates.ReadFile(filepath.Join("templates", identifier+".gotmpl"))
	if err != nil {
		return "", fmt.Errorf("failed to load template %s due to: %w", identifier, err)
	}
	return string(f), nil
}
