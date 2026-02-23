package fiber

import (
	"slices"
)

const pathToRoutes = "internal/transports/http/routes"
const pathToHandlers = "internal/transports/http/handler"

func getArticle(input string) string {
	vowels := []string{"a", "e", "i", "o", "u", "A", "E", "I", "O", "U"}
	if slices.Contains(vowels, string(input[0])) {
		return "an"
	}
	return "a"
}
