package service

import (
	"embed"
)

const pathToService = "internal/service"

//go:embed templates/*
var templates embed.FS
