package service

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
	"zuno/cmd/data"
	"zuno/cmd/utils"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewServiceData struct {
	Package            string
	FieldsType         string
	Module             string
	ServiceType        string
	ServiceVariable    string
	RepositoryVariable string
	RepositoryType     string
	Readable           string
	ReadablePlural     string
	Variable           string
	VariablePlural     string
	FileName           string
}

type RegisterNewServiceData struct {
	Module      string // Domain module name, PascalCase (e.g. "ProductVariant")
	ServiceType string // ServiceType struct type (e.g. "ProductVariantService")
	FileName    string
}

type AddFieldsToServiceData struct {
	Module             string // Domain module name, PascalCase (e.g. "ProductVariant")
	FieldsType         string // Payload / fields struct (e.g. "ProductVariantFields")
	RepositoryFunction string // fromRepositoryProductVariant
	Variable           string // productVariant
	FileName           string // product_variant_service.go
}

// AddNewService adds a new service
func AddNewService(packageName string, module string) error {
	data, err := prepareAddNewServiceData(packageName, module)
	if err != nil {
		return err
	}
	return utils.CreateFromTemplate(
		templates, "templates/new_service.gotmpl",
		path.Join(pathToService, data.FileName), data,
	)
}

// RegisterNewService registers a new service into the Services wrapper / struct
func RegisterNewService(module string) error {
	data, err := prepareRegisterNewServiceData(module)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToService, data.FileName)

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

		// ---- Modify `type Services struct { ... }`
		if ts, ok := n.(*ast.TypeSpec); ok && ts.Name.Name == "Services" {
			if st, ok := ts.Type.(*ast.StructType); ok {

				// prevent duplicate fields
				for _, f := range st.Fields.List {
					if len(f.Names) > 0 && f.Names[0].Name == data.Module {
						return false
					}
				}

				st.Fields.List = append(st.Fields.List, &ast.Field{
					Names: []*ast.Ident{
						ast.NewIdent(data.Module),
					},
					Type: ast.NewIdent(data.ServiceType),
				})
			}
		}

		// ---- Modify constructor `func New(...) *Services`
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "New" {

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				ret, ok := n.(*ast.ReturnStmt)
				if !ok || len(ret.Results) != 1 {
					return true
				}

				// Expect: return &Services{ ... }
				unary, ok := ret.Results[0].(*ast.UnaryExpr)
				if !ok || unary.Op != token.AND {
					return true
				}

				cl, ok := unary.X.(*ast.CompositeLit)
				if !ok {
					return true
				}

				// Ensure it's Services{}
				ident, ok := cl.Type.(*ast.Ident)
				if !ok || ident.Name != "Services" {
					return true
				}

				// Prevent duplicate key
				for _, elt := range cl.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						if key, ok := kv.Key.(*ast.Ident); ok && key.Name == data.Module {
							return false
						}
					}
				}

				cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
					Key: ast.NewIdent(data.Module),
					Value: &ast.CallExpr{
						Fun: ast.NewIdent(fmt.Sprintf("New%s", data.ServiceType)),
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("repositories"),
								Sel: ast.NewIdent(data.Module),
							},
						},
					},
				})

				return false
			})
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

// AddFieldsToService adds fields to the service
func AddFieldsToService(module string, fields []data.Field) error {
	data, err := prepareAddFieldsToServiceData(module)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToService, data.FileName)

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
							Value: field.ServiceTags(),
						},
					})
				}
			}
		}

		// 2️⃣ Find fromService<Module>
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == data.RepositoryFunction {

			ast.Inspect(fn.Body, func(n ast.Node) bool {

				// Look for composite literal
				cl, ok := n.(*ast.CompositeLit)
				if !ok {
					return true
				}

				// Find FieldsType composite literal
				if ident, ok := cl.Type.(*ast.Ident); ok && ident.Name == data.FieldsType {
					for _, field := range fields {
						cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
							Key: ast.NewIdent(field.Name),
							Value: &ast.SelectorExpr{
								X:   ast.NewIdent(data.Variable),
								Sel: ast.NewIdent(field.Name),
							},
						})
					}
				}

				return true
			})
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

func prepareAddNewServiceData(packageName, module string) (AddNewServiceData, error) {
	pluralize := pluralize.NewClient()
	return AddNewServiceData{
		Package:            packageName,
		FieldsType:         module + "Fields",
		Module:             module,
		ServiceType:        module + "Service",
		ServiceVariable:    strcase.ToCamel(module) + "Service",
		RepositoryVariable: strcase.ToCamel(module) + "Repository",
		RepositoryType:     module + "Repository",
		Readable:           strings.ReplaceAll(strcase.ToKebab(module), "-", " "),
		ReadablePlural:     pluralize.Plural(strings.ReplaceAll(strcase.ToKebab(module), "-", " ")),
		Variable:           strcase.ToCamel(module),
		VariablePlural:     pluralize.Plural(strcase.ToCamel(module)),
		FileName:           strcase.ToSnake(module) + "_service.go",
	}, nil
}

func prepareRegisterNewServiceData(module string) (RegisterNewServiceData, error) {
	return RegisterNewServiceData{
		Module:      module,
		ServiceType: module + "Service",
		FileName:    "services.go",
	}, nil
}

func prepareAddFieldsToServiceData(module string) (AddFieldsToServiceData, error) {
	data := AddFieldsToServiceData{
		Module:             module,
		FieldsType:         module + "Fields",
		RepositoryFunction: "fromRepository" + module,
		Variable:           strcase.ToCamel(module),
		FileName:           strcase.ToSnake(module) + "_service.go",
	}
	return data, nil
}
