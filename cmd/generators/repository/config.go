package repository

import (
	"embed"
)

const pathToRepository = "internal/repository"

//go:embed templates/*
var templates embed.FS
