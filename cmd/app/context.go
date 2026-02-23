package app

import "zuno/cmd/config"

type AppContext struct {
	Config *config.Config
}

var Ctx = &AppContext{}
