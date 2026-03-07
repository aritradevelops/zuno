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

	"github.com/aritradevelops/zuno/cmd/utils"
	"github.com/aritradevelops/zuno/pkg/logger"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewRepositoryData struct {
	Package        string
	RepositoryType string
	Module         string
	Variable       string
	VariablePlural string
	FieldsType     string
	Readable       string
	ReadablePlural string
	FileName       string
}

type RegisterNewRepositoryData struct {
	Module         string // Domain module name, PascalCase (e.g. "ProductVariant")
	RepositoryType string // RepositoryType struct type (e.g. "ProductVariantRepository")
	FileName       string
}

func AddNewRepository(packageName, module string) error {
	data, err := prepareAddNewRepositoryData(packageName, module)
	if err != nil {
		return err
	}
	return utils.CreateFromTemplate(
		templates, "templates/new_repository.gotmpl",
		path.Join(pathToRepository, data.FileName), data,
	)
}

// RegisterNewService registers a new service into the Services wrapper / struct
func RegisterNewRepository(module string) error {
	data, err := prepareRegisterNewRepositoryData(module)
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

		// ---- Modify constructor `func NewRepositories(...) *Repositories`
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "NewRepositories" {
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				ret, ok := n.(*ast.ReturnStmt)
				if !ok || len(ret.Results) != 1 {
					return true
				}

				// Expect: return &Repositories{ ... }
				unary, ok := ret.Results[0].(*ast.UnaryExpr)
				if !ok || unary.Op != token.AND {
					return true
				}

				cl, ok := unary.X.(*ast.CompositeLit)
				if !ok {
					return true
				}

				// Prevent duplicate key
				for _, elt := range cl.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						if key, ok := kv.Key.(*ast.Ident); ok && key.Name == data.Module {
							logger.Info("found duplicate key")
							return false
						}
					}
				}

				cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
					Key: ast.NewIdent(data.Module),
					Value: &ast.CallExpr{
						Fun: ast.NewIdent(fmt.Sprintf("New%s", data.RepositoryType)),
						Args: []ast.Expr{
							ast.NewIdent("db"),
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

func prepareAddNewRepositoryData(packageName, module string) (AddNewRepositoryData, error) {
	pluralize := pluralize.NewClient()
	return AddNewRepositoryData{
		Package:        packageName,
		RepositoryType: module + "Repository",
		Module:         module,
		Variable:       strcase.ToCamel(module),
		VariablePlural: pluralize.Plural(strcase.ToCamel(module)),
		FieldsType:     module + "Fields",
		Readable:       strings.ReplaceAll(strcase.ToKebab(module), "-", " "),
		ReadablePlural: pluralize.Plural(strings.ReplaceAll(strcase.ToKebab(module), "-", " ")),
		FileName:       strcase.ToSnake(module) + "_repository.go",
	}, nil
}

func prepareRegisterNewRepositoryData(module string) (RegisterNewRepositoryData, error) {
	return RegisterNewRepositoryData{
		Module:         module,
		RepositoryType: module + "Repository",
		FileName:       "repositories.go",
	}, nil
}
