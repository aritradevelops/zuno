package app

import "github.com/aritradevelops/zuno/cmd/config"

type AppContext struct {
	Config *config.Config
}

var Ctx = &AppContext{}
