package fiber

import (
	"embed"
	"slices"
)

//go:embed templates/*
var templates embed.FS

const pathToHttpProvider = "internal/transports/http"
const pathToRoutes = "internal/transports/http/routes"
const pathToHandlers = "internal/transports/http/handler"

func getArticle(input string) string {
	vowels := []string{"a", "e", "i", "o", "u", "A", "E", "I", "O", "U"}
	if slices.Contains(vowels, string(input[0])) {
		return "an"
	}
	return "a"
}
