package data

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

type Field struct {
	Name string
	Type string
}

func (f Field) HandlerTags() string {
	return fmt.Sprintf("`json:\"%s\" example:\"%s\"`", strcase.ToSnake(f.Name), f.GetIntuitiveExample())
}

func (f Field) ModelTags() string {
	return fmt.Sprintf("`bson:\"%s\"`", strcase.ToSnake(f.Name))
}

func (f Field) GetIntuitiveExample() string {
	switch f.Type {
	case "string":
		if f.Name == "email" {
			return "someone@example.com"
		}
		return "Demo"
	case "int":
		return "123"
	case "float64":
		return "123.45"
	case "bool":
		return "true"
	default:
		return "unknown"
	}
}

func (f Field) ServiceTags() string {
	return fmt.Sprintf("`validate:\"%s\"`", f.GetIntuitiveValidations())
}

func (f Field) GetIntuitiveValidations() string {
	validations := []string{"required"}
	switch f.Name {
	case "email":
		validations = append(validations, "email")
	case "password":
		validations = append(validations, "min=8")
	case "int", "float64":
		validations = append(validations, "numeric")
	case "bool":
		validations = append(validations, "boolean")
	}
	return strings.Join(validations, ",")
}
