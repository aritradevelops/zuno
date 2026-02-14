package translation

import (
	"goserve/pkg/logger"
	"os"
	"path"

	contribi18n "github.com/gofiber/contrib/v3/i18n"
	"github.com/gofiber/fiber/v3"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var translator *contribi18n.I18n

func init() {
	cwd, _ := os.Getwd()
	localesPath := path.Join(cwd, "locales")
	logger.Info().Str("path", localesPath).Msg("loading locales")
	translator = contribi18n.New(&contribi18n.Config{
		RootPath:        localesPath,
		AcceptLanguages: []language.Tag{language.English, language.Bengali},
		DefaultLanguage: language.English,
	})
}

func Localize(c fiber.Ctx, id string, data ...any) string {
	var tmplData any
	if len(data) > 0 {
		tmplData = data[0]
	}
	msg, err := translator.Localize(c, &i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: tmplData,
	})
	if err != nil {
		logger.Warn().Str("key", id).Msg("Missing entry in locales")
		return id
	}
	return msg
}
