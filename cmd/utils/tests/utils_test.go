package utils_test

import (
	"embed"
	"os"
	"testing"
	"zuno/cmd/utils"

	"github.com/stretchr/testify/assert"
)

//go:embed templates/*
var templateFS embed.FS

func TestCloneTemplates(t *testing.T) {
	err := utils.CloneTemplates(templateFS, "templates", "test-clone", nil)
	assert.NoError(t, err)
	stat, err := os.Stat("./test-clone/file.txt")
	assert.NoError(t, err)
	assert.False(t, stat.IsDir())
	assert.NoError(t, os.RemoveAll("test-clone"))
}

func TestCreateFromTemplate(t *testing.T) {
	err := utils.CreateFromTemplate(templateFS, "templates/new_demo.go.gotmpl", "test-clone/new_demo.go", nil)
	assert.NoError(t, err)
	stat, err := os.Stat("./test-clone/new_demo.go")
	assert.NoError(t, err)
	assert.False(t, stat.IsDir())
	assert.NoError(t, os.RemoveAll("test-clone"))
}
