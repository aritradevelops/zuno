package repository

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"text/template"

	"github.com/ettle/strcase"
)

type AddNewRepositoryData struct {
	Package        string
	FieldsType     string
	Module         string
	RepositoryType string
	FileName       string
}

type RegisterNewRepositoryData struct {
	Module         string // Domain module name, PascalCase (e.g. "ProductVariant")
	RepositoryType string // RepositoryType struct type (e.g. "ProductVariantRepository")
	FileName       string
}

// AddNewRepository adds a new repository
func AddNewRepository(packageName string, module string) error {
	data, err := prepareAddNewRepositoryData(packageName, module)
	if err != nil {
		return err
	}
	tmplContent, err := loadTemplate("new_repository")
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToRepository, data.FileName)
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

		// ---- Modify `type Repositories struct { ... }`
		if ts, ok := n.(*ast.TypeSpec); ok && ts.Name.Name == "Repositories" {
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
					Type: ast.NewIdent(data.RepositoryType),
				})
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

func prepareAddNewRepositoryData(packageName, module string) (AddNewRepositoryData, error) {
	return AddNewRepositoryData{
		Package:        packageName,
		FieldsType:     module + "Fields",
		Module:         module,
		RepositoryType: module + "Repository",
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
