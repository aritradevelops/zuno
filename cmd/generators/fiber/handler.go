package fiber

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

type AddNewHandlerData struct {
	Package         string // Go module import root (e.g. "goserve")
	Module          string // Domain module name, PascalCase (e.g. "ProductVariant")
	Variable        string // Variable name for the module, camelCase(e.g productVariant)
	VariablePlural  string // Variable name for the module, plural, camelCase(e.g productVariants)
	HandlerType     string // Handler struct type (e.g. "ProductVariantHandler")
	ServiceType     string // Service interface type (e.g. "ProductVariantService")
	ServiceVariable string // Service field/variable name (e.g. "productVariantService")
	FieldsType      string // Payload / fields struct (e.g. "ProductVariantFields")
	RoutePrefix     string // URL path prefix, plural, kebab-case (e.g. "product-variants")
	RouteTag        string // Swagger tag, kebab-case singular (e.g. "product-variant")
	Readable        string // Human-readable singular (e.g. "product variant")
	ReadablePlural  string // Human-readable plural (e.g. "product variants")
	Article         string // a | an
	FileName        string // product_variant_handler.go
}

type RegisterNewHandlerData struct {
	Module      string // Domain module name, PascalCase (e.g. "ProductVariant")
	HandlerType string // Handler struct type (e.g. "ProductVariantHandler")
	FileName    string
}

type AddFieldsToHandlerData struct {
	Module          string // Domain module name, PascalCase (e.g. "ProductVariant")
	FieldsType      string // Payload / fields struct (e.g. "ProductVariantFields")
	ServiceFunction string // fromProductVariant
	Variable        string // productVariant
	FileName        string // product_variant_handler.go
}

// AddNewHandler adds default handlers for a new module
func AddNewHandler(packageName, moduleName string) error {
	data, err := prepareNewHandlerData(packageName, moduleName)

	if err != nil {
		return err
	}

	err = utils.CreateFromTemplate(
		templates, "templates/new_handler.gotmpl",
		path.Join(pathToHandlers, data.FileName), data,
	)
	return err
}

// RegisterNewHandler registers a new handler into the Handlers wrapper / struct
func RegisterNewHandler(module string) error {
	data, err := prepareRegisterNewHandlerData(module)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToHandlers, data.FileName)

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

		// ---- Modify `type Handlers struct { ... }`
		if ts, ok := n.(*ast.TypeSpec); ok && ts.Name.Name == "Handlers" {
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
					Type: &ast.StarExpr{
						X: ast.NewIdent(data.HandlerType),
					},
				})
			}
		}

		// ---- Modify constructor `func New(...) *Handlers`
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "New" {

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				ret, ok := n.(*ast.ReturnStmt)
				if !ok || len(ret.Results) != 1 {
					return true
				}

				// Expect: return &Handlers{ ... }
				unary, ok := ret.Results[0].(*ast.UnaryExpr)
				if !ok || unary.Op != token.AND {
					return true
				}

				cl, ok := unary.X.(*ast.CompositeLit)
				if !ok {
					return true
				}

				// Ensure it's Handlers{}
				ident, ok := cl.Type.(*ast.Ident)
				if !ok || ident.Name != "Handlers" {
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
						Fun: ast.NewIdent(fmt.Sprintf("New%s", data.HandlerType)),
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("services"),
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

// AddFieldsToHandler adds fields to the handler
func AddFieldsToHandler(module string, fields []data.Field) error {
	data, err := prepareAddFieldsToHandlerData(module)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToHandlers, data.FileName)

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
							Value: field.HandlerTags(),
						},
					})
				}
			}
		}

		// 2️⃣ Find fromService<Module>
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == data.ServiceFunction {

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

func prepareNewHandlerData(packageName string, module string) (AddNewHandlerData, error) {
	pluralize := pluralize.NewClient()

	data := AddNewHandlerData{
		Package:         packageName,
		Module:          module,
		Variable:        strcase.ToCamel(module),
		VariablePlural:  pluralize.Plural(strcase.ToCamel(module)),
		HandlerType:     module + "Handler",
		ServiceType:     module + "Service",
		ServiceVariable: strcase.ToCamel(module) + "Service",
		FieldsType:      module + "Fields",
		RoutePrefix:     pluralize.Plural(strcase.ToKebab(module)),
		RouteTag:        strcase.ToKebab(module),
		Readable:        strings.ReplaceAll(strcase.ToKebab(module), "-", " "),
		ReadablePlural:  pluralize.Plural(strings.ReplaceAll(strcase.ToKebab(module), "-", " ")),
		Article:         getArticle(module),
		FileName:        strcase.ToSnake(module) + "_handler.go",
	}
	return data, nil
}

func prepareRegisterNewHandlerData(module string) (RegisterNewHandlerData, error) {
	data := RegisterNewHandlerData{
		Module:      module,
		HandlerType: module + "Handler",
		FileName:    "handlers.go",
	}
	return data, nil
}

func prepareAddFieldsToHandlerData(module string) (AddFieldsToHandlerData, error) {
	data := AddFieldsToHandlerData{
		Module:          module,
		FieldsType:      module + "Fields",
		ServiceFunction: "fromService" + module,
		Variable:        strcase.ToCamel(module),
		FileName:        strcase.ToSnake(module) + "_handler.go",
	}
	return data, nil
}
