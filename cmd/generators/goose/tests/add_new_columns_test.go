package tests

import (
	"testing"

	"github.com/aritradevelops/zuno/cmd/data"
	"github.com/aritradevelops/zuno/cmd/generators/goose"
	"github.com/stretchr/testify/assert"
)

func TestAddNewColumnsMigration(t *testing.T) {
	err := goose.AddNewColumnsMigration("test", "test", "./tmp/migrations", []data.Field{
		{
			Name:    "email",
			RawType: "string",
			SqlInfo: &data.SqlInfo{
				Nullable: false,
				Default:  "",
				Unique:   true,
			},
		},
	})
	assert.NoError(t, err)
}
