package mongodb

import (
	"embed"
)

const pathToRepository = "internal/adapters/mongodb"
const pathToModel = "internal/adapters/mongodb"
const pathToMongodbAdapter = "internal/adapters/mongodb"

//go:embed templates/*
var templates embed.FS
