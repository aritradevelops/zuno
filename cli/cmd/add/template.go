package add

import (
	"embed"
	"io"
	"path"
	"text/template"
)

//go:embed templates/**
var templateFs embed.FS

func Execute(w io.WriteCloser, relPath string, data any) error {
	tmpl, err := template.New("").ParseFS(templateFs, relPath)
	if err != nil {
		return err
	}
	defer w.Close()
	return tmpl.ExecuteTemplate(w, path.Base(relPath), data)
}
