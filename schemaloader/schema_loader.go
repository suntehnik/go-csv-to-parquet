package schemaloader

import (
	"github.com/apache/arrow/go/v16/arrow"
	"errors"
	"gopkg.in/yaml.v3"
)

// FieldSpec описывает поле из YAML-схемы
// (минимально для компиляции и тестов)
type FieldSpec struct {
	Name      string
	Type      string
	Nullable  bool
	Precision *int
	Scale     *int
	Format    string
}

// SchemaSpec описывает всю YAML-схему
// (минимально для компиляции и тестов)
type SchemaSpec struct {
	Fields []FieldSpec
}

// LoadSchema парсит YAML и возвращает Arrow Schema (частичная реализация для прохождения части тестов)
func LoadSchema(yamlData []byte) (*arrow.Schema, error) {
	var spec SchemaSpec
	if err := yaml.Unmarshal(yamlData, &spec); err != nil {
		return nil, errors.New("invalid YAML syntax")
	}
	if len(spec.Fields) == 0 {
		return nil, errors.New("schema must have at least one field")
	}
	nameSet := make(map[string]struct{})
	for _, f := range spec.Fields {
		if f.Name == "" {
			return nil, errors.New("field name required")
		}
		if f.Type == "" {
			return nil, errors.New("field type required")
		}
		if _, ok := nameSet[f.Name]; ok {
			return nil, errors.New("duplicate field name: " + f.Name)
		}
		nameSet[f.Name] = struct{}{}
		if !isSupportedType(f.Type) {
			return nil, errors.New("unsupported type: " + f.Type)
		}
	}
	return nil, errors.New("not implemented")
}

func isSupportedType(t string) bool {
	switch t {
	case "int32", "int64", "float32", "float64", "bool", "string", "date", "timestamp", "decimal", "bytes":
		return true
	default:
		return false
	}
}
