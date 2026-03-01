package bun

import (
	"embed"
)

const pathToRepository = "internal/adapters/bun"
const pathToModel = "internal/adapters/bun"
const pathToBunAdapter = "internal/adapters/bun"
const pathToMigration = "internal/adapters/bun/migrations"

//go:embed templates/*
var templates embed.FS
