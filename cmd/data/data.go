package data

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
)

type Field struct {
	Name      string
	RawType   string
	IsArray   bool
	IsPointer bool
	SqlInfo   *SqlInfo
}

type SqlInfo struct {
	Unique   bool
	Nullable bool
	Default  string
}

func (f Field) HandlerTags() string {
	return fmt.Sprintf("`json:\"%s\" example:\"%s\"`", strcase.ToSnake(f.Name), f.GetIntuitiveExample())
}

func (f Field) BsonTags() string {
	return fmt.Sprintf("`bson:\"%s\"`", strcase.ToSnake(f.Name))
}

func (f Field) DbName() string {
	return strcase.ToSnake(f.Name)
}

func (f Field) GoType() string {
	goType := f.RawType
	if f.IsArray {
		goType = "[]" + goType
	}
	if f.IsPointer {
		goType = "*" + goType
	}
	return goType
}

func (f Field) ColumnName() string {
	return strcase.ToSnake(f.Name)
}

func (f Field) SqlType() string {
	switch f.RawType {
	case "string":
		return "VARCHAR"
	case "int", "int32":
		return "INTEGER"
	case "int8", "int16":
		return "SMALLINT"
	case "int64":
		return "BIGINT"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TIMESTAMP WITH TIME ZONE"
	case "float32":
		return "REAL"
	case "float64":
		return "DOUBLE PRECISION"
	default:
		return "TEXT"
	}
}

func (f Field) BunTags() string {
	additional := ""
	if f.SqlInfo.Unique {
		additional += ",unique"
	}
	if f.SqlInfo.Nullable {
		additional += ",nullzero"
	}
	if f.SqlInfo.Default != "" {
		additional += ",default:" + f.SqlInfo.Default
	}
	if f.IsArray {
		additional += ",array"
	}

	return fmt.Sprintf("`bun:\"%s,type:%s%s\"`", strcase.ToSnake(f.Name), strings.ToLower(f.SqlType()), additional)
}

func (f Field) GetIntuitiveExample() string {
	switch f.RawType {
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
