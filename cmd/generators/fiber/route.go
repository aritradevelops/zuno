package fiber

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"zuno/cmd/utils"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewRouterData struct {
	Package         string
	Module          string
	HandlerVariable string
	HandlerType     string
	FileName        string
}
type RegisterNewRouterData struct {
	Module      string
	HandlerType string
	RoutePrefix string // URL path prefix, plural, kebab-case (e.g. "product-variants")
	FileName    string
}

// AddNewRouter adds a router with basic routes for a new module
func AddNewRouter(packageName, module string) error {
	data, err := prepareAddNewRouterData(packageName, module)

	if err != nil {
		return err
	}

	return utils.CreateFromTemplate(
		templates, "templates/new_router.gotmpl",
		path.Join(pathToRoutes, data.FileName), data,
	)
}

// RegisterNewRouter registers a new router for the given module under the base router
func RegisterNewRouter(module string) error {
	data, err := prepareRegisterNewRouterData(module)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToRoutes, data.FileName)
	fset := token.NewFileSet()

	sourceFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	f, err := parser.ParseFile(fset, "random.go", sourceFile, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if fn.Name.Name != "Register" {
			return true
		}

		// Build: app.Group("/users")
		appGroupCall := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("app"),
				Sel: ast.NewIdent("Group"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"/%s"`, data.RoutePrefix),
				},
			},
		}

		// Build: RegisterUserRoutes(app.Group(...), authMiddleware, handlers.User)
		newCall := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: ast.NewIdent(fmt.Sprintf("Register%sRoutes", data.Module)),
				Args: []ast.Expr{
					appGroupCall,
					ast.NewIdent("authMiddleware"),
					&ast.SelectorExpr{
						X:   ast.NewIdent("handlers"),
						Sel: ast.NewIdent(data.Module),
					},
				},
			},
		}

		// Append to function body
		fn.Body.List = append(fn.Body.List, newCall)

		return false
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

func prepareAddNewRouterData(packageName, module string) (AddNewRouterData, error) {
	return AddNewRouterData{
		Package:         packageName,
		Module:          module,
		HandlerVariable: strcase.ToCamel(module) + "Handler",
		HandlerType:     module + "Handler",
		FileName:        strcase.ToSnake(module) + "_routes.go",
	}, nil
}

func prepareRegisterNewRouterData(module string) (RegisterNewRouterData, error) {
	pluralize := pluralize.NewClient()
	return RegisterNewRouterData{
		Module:      module,
		HandlerType: module + "Handler",
		RoutePrefix: pluralize.Plural(strcase.ToKebab(module)),
		FileName:    "register.go",
	}, nil
}
