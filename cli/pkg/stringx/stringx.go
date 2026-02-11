package stringx

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/gobeam/stringy"
)

var pluralizer = pluralize.NewClient()

type Stringx struct {
	instance stringy.StringManipulation
	original string
}

func New(original string) Stringx {
	return Stringx{
		instance: stringy.New(original),
		original: original,
	}
}

func (n Stringx) CollectionName() string {
	return pluralizer.Plural(n.instance.CamelCase().Get())
}

func (n Stringx) RepositoryName() string {
	return fmt.Sprintf("%sRepository", n.instance.PascalCase().Get())
}

func (n Stringx) ModelName() string {
	return fmt.Sprintf("%sModel", n.instance.PascalCase().Get())
}

func (n Stringx) ModuleName() string {
	return n.instance.PascalCase().Get()
}

func (n Stringx) VariableName() string {
	return n.instance.CamelCase().Get()
}

func (n Stringx) VariableNamePlural() string {
	return pluralizer.Plural(n.instance.CamelCase().Get())
}

func (n Stringx) RepositoryFileName() string {
	return fmt.Sprintf("%s_repository.go", n.instance.SnakeCase().ToLower())
}
func (n Stringx) ModelFileName() string {
	return fmt.Sprintf("%s_model.go", n.instance.SnakeCase().ToLower())
}
func (n Stringx) SnakeCase() string {
	return n.instance.SnakeCase().ToLower()
}

func (n Stringx) BsonTag() string {
	return fmt.Sprintf("bson:\"%s\"", n.SnakeCase())
}

func (n Stringx) JsonTag() string {
	return fmt.Sprintf("json:\"%s\"", n.SnakeCase())
}

func (n Stringx) ValidateTag() string {
	if n.ModuleName() == "Email" {
		return "validate:\"required,email\""
	}
	return "validate:\"required\""
}
