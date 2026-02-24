package mongodb

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path"
	"zuno/cmd/data"

	"github.com/ettle/strcase"
)

type AddNewModelData struct {
	Package    string
	FieldsType string
	Module     string
	FileName   string
}
type AddFieldsToModelData struct {
	Module     string
	FieldsType string
	Fields     []data.Field
	FileName   string
}

func AddNewModel(packageName, module string) error {
	data, err := prepareAddNewModelData(packageName, module)
	if err != nil {
		return err
	}
	tmplContent, err := loadTemplate("new_model")
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToModel, data.FileName)
	tmpl, err := template.New(filePath).Parse(tmplContent)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// AddFieldsToModel adds fields to the model
func AddFieldsToModel(module string, fields []data.Field) error {
	data, err := prepareAddFieldsToModelData(module, fields)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToRepository, data.FileName)

	sourceFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "random.go", sourceFile, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Inspect(f, func(n ast.Node) bool {

		// 1️⃣ Find Fields struct
		ts, ok := n.(*ast.TypeSpec)
		if ok && ts.Name.Name == data.FieldsType {
			st, ok := ts.Type.(*ast.StructType)
			if ok {
				for _, field := range fields {
					st.Fields.List = append(st.Fields.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(field.Name)},
						Type:  ast.NewIdent(field.Type),
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: field.ModelTags(),
						},
					})
				}
			}
		}

		funcName := "toRepository"

		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != funcName {
			return true
		}

		ast.Inspect(fn.Body, func(n ast.Node) bool {

			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			// Look for repository.UserFields composite literal
			se, ok := cl.Type.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkgIdent, ok := se.X.(*ast.Ident)
			if !ok || pkgIdent.Name != "repository" {
				return true
			}

			if se.Sel.Name != data.FieldsType {
				return true
			}
			for _, field := range data.Fields {
				// Append new field
				cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
					Key: ast.NewIdent(field.Name),
					Value: &ast.SelectorExpr{
						X:   ast.NewIdent("u"),
						Sel: ast.NewIdent(field.Name),
					},
				})
			}

			return false
		})

		return true
	})

	// Reset file before writing
	if _, err := sourceFile.Seek(0, 0); err != nil {
		return err
	}

	if err := sourceFile.Truncate(0); err != nil {
		return err
	}

	// Print back the modified source
	if err := format.Node(sourceFile, fset, f); err != nil {
		return err
	}

	return nil
}
func prepareAddNewModelData(packageName, module string) (AddNewModelData, error) {
	return AddNewModelData{
		Package:    packageName,
		FieldsType: module + "Fields",
		Module:     module,
		FileName:   strcase.ToSnake(module) + "_model.go",
	}, nil
}

func prepareAddFieldsToModelData(module string, fields []data.Field) (AddFieldsToModelData, error) {
	return AddFieldsToModelData{
		Module:     module,
		FieldsType: module + "Fields",
		Fields:     fields,
		FileName:   strcase.ToSnake(module) + "_model.go",
	}, nil
}
