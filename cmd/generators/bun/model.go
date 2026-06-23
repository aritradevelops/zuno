package bun

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/aritradevelops/zuno/cmd/data"
	"github.com/aritradevelops/zuno/cmd/utils"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewModelData struct {
	Package    string
	FieldsType string
	Module     string
	Table      string
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
	return utils.CreateFromTemplate(
		templates, "templates/new_model.gotmpl",
		path.Join(pathToModel, data.FileName), data,
	)
}

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

		// 1️⃣ Find Fields struct and append struct fields
		ts, ok := n.(*ast.TypeSpec)
		if ok && ts.Name.Name == data.FieldsType {
			st, ok := ts.Type.(*ast.StructType)
			if ok {
				for _, field := range data.Fields {
					st.Fields.List = append(st.Fields.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(field.Name)},
						Type:  ast.NewIdent(field.GoType()),
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: field.BunTags(),
						},
					})
				}
			}
		}

		// 2️⃣ Update toDomain function
		funcName := "toDomain"

		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == funcName {

			ast.Inspect(fn.Body, func(n ast.Node) bool {

				cl, ok := n.(*ast.CompositeLit)
				if !ok {
					return true
				}

				// Look for domain.UserFields composite literal
				se, ok := cl.Type.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkgIdent, ok := se.X.(*ast.Ident)
				if !ok || pkgIdent.Name != "domain" {
					return true
				}

				if se.Sel.Name != data.FieldsType {
					return true
				}

				for _, field := range data.Fields {
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
		}

		// 3️⃣ Update Fields() function
		fn2, ok := n.(*ast.FuncDecl)
		if ok && fn2.Name.Name == "Fields" {

			if fn2.Body == nil {
				return true
			}

			for _, stmt := range fn2.Body.List {

				ret, ok := stmt.(*ast.ReturnStmt)
				if !ok {
					continue
				}

				for _, expr := range ret.Results {

					cl, ok := expr.(*ast.CompositeLit)
					if !ok {
						continue
					}

					arr, ok := cl.Type.(*ast.ArrayType)
					if !ok {
						continue
					}

					ident, ok := arr.Elt.(*ast.Ident)
					if !ok || ident.Name != "string" {
						continue
					}

					// Track existing fields to avoid duplicates
					existing := map[string]bool{}

					for _, e := range cl.Elts {
						if lit, ok := e.(*ast.BasicLit); ok {
							name := strings.Trim(lit.Value, `"`)
							existing[name] = true
						}
					}

					for _, field := range fields {

						name := field.DbName()

						if existing[name] {
							continue
						}

						cl.Elts = append(cl.Elts, &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s"`, name),
						})
					}
				}
			}
		}

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
	pluralize := pluralize.NewClient()
	return AddNewModelData{
		Package:    packageName,
		FieldsType: module + "Fields",
		Module:     module,
		Table:      strcase.ToSnake(pluralize.Plural(module)),
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
