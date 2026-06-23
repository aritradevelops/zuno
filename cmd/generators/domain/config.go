package domain

import (
	"embed"
)

const pathToDomain = "internal/domain"

//go:embed templates/*
var templates embed.FS
