package domain

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"

	"github.com/aritradevelops/zuno/cmd/data"
	"github.com/aritradevelops/zuno/cmd/utils"

	"github.com/ettle/strcase"
)

type AddNewDomainData struct {
	Package    string
	FieldsType string
	Module     string
	FileName   string
}

type AddFieldsToDomainData struct {
	Module     string // Domain module name, PascalCase (e.g. "ProductVariant")
	FieldsType string // Payload / fields struct (e.g. "ProductVariantFields")
	FileName   string // product_variant_service.go
	Fields     []data.Field
}

// AddNewDomain adds a new domain
func AddNewDomain(packageName string, module string) error {
	data, err := prepareAddNewDomainData(packageName, module)
	if err != nil {
		return err
	}
	return utils.CreateFromTemplate(
		templates, "templates/new_domain.gotmpl",
		path.Join(pathToDomain, data.FileName), data,
	)
}

// AddFieldsToDomain adds fields to the domain
func AddFieldsToDomain(module string, fields []data.Field) error {
	data, err := prepareAddFieldsToDomainData(module, fields)
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToDomain, data.FileName)

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
				for _, field := range data.Fields {
					st.Fields.List = append(st.Fields.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(field.Name)},
						Type:  ast.NewIdent(field.GoType()),
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: field.ServiceTags(),
						},
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

func prepareAddNewDomainData(packageName, module string) (AddNewDomainData, error) {
	return AddNewDomainData{
		Package:    packageName,
		FieldsType: module + "Fields",
		Module:     module,
		FileName:   strcase.ToSnake(module) + ".go",
	}, nil
}

func prepareAddFieldsToDomainData(module string, fields []data.Field) (AddFieldsToDomainData, error) {
	data := AddFieldsToDomainData{
		Module:     module,
		FieldsType: module + "Fields",
		FileName:   strcase.ToSnake(module) + ".go",
		Fields:     fields,
	}
	return data, nil
}
