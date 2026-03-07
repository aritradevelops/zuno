package goose

import (
	"embed"
	"fmt"
)

//go:embed templates/*
var templates embed.FS

func GetMigrationPathFromAdapter(adapter string) string {
	return fmt.Sprintf("internal/adapters/%s/migrations", adapter)
}

func GetSQLTypeFromRaw(rawType string) string {
	switch rawType {
	case "string":
		return "VARCHAR(255)"
	case "int", "int32":
		return "INTEGER"
	case "int8", "int16":
		return "SMALLINT"
	case "int64":
		return "BIGINT"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TIMESTAMP WITH TIME ZONE"
	case "float32":
		return "REAL"
	case "float64":
		return "DOUBLE PRECISION"
	default:
		return "TEXT"
	}
}
