package repository

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"zuno/cmd/data"
	"zuno/cmd/utils"

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

type AddFieldsToRepositoryData struct {
	Module     string // Domain module name, PascalCase (e.g. "ProductVariant")
	FieldsType string // Payload / fields struct (e.g. "ProductVariantFields")
	FileName   string // product_variant_repository.go
}

// AddNewRepository adds a new repository
func AddNewRepository(packageName string, module string) error {
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

// AddFieldsToRepository adds fields to the repository
func AddFieldsToRepository(module string, fields []data.Field) error {
	data, err := prepareAddFieldsToRepositoryData(module)
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
					})
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

func prepareAddFieldsToRepositoryData(module string) (AddFieldsToRepositoryData, error) {
	data := AddFieldsToRepositoryData{
		Module:     module,
		FieldsType: module + "Fields",
		FileName:   strcase.ToSnake(module) + "_repository.go",
	}
	return data, nil
}
