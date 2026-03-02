package barebone

import (
	"embed"

	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/cmd/utils"
)

//go:embed all:*
var templateFs embed.FS
var directCloneDirs = []string{"locales"}

func Initialize(cfg *config.Config) error {
	return utils.CloneTemplates(templateFs, "templates", ".", cfg, directCloneDirs...)
}
