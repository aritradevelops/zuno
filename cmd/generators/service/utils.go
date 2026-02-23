package service

import (
	"embed"
	"fmt"
	"path/filepath"
	"slices"
)

const pathToService = "internal/service"

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

func getArticle(input string) string {
	vowels := []string{"a", "e", "i", "o", "u", "A", "E", "I", "O", "U"}
	if slices.Contains(vowels, string(input[0])) {
		return "an"
	}
	return "a"
}
