package schemaloader

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/require"
)

func TestLoadSchema_ValidAllTypes(t *testing.T) {
	yaml := []byte(`fields:\n  - name: id\n    type: int64\n  - name: name\n    type: string\n  - name: created_at\n    type: timestamp\n  - name: price\n    type: decimal\n    precision: 10\n    scale: 2\n  - name: active\n    type: bool\n  - name: data\n    type: bytes\n  - name: score\n    type: float64\n  - name: event_date\n    type: date\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_NullableVariants(t *testing.T) {
	yaml := []byte(`fields:\n  - name: n1\n    type: int64\n    nullable: true\n  - name: n2\n    type: string\n    nullable: false\n  - name: n3\n    type: float32\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_DecimalPrecisionScale(t *testing.T) {
	yaml := []byte(`fields:\n  - name: price\n    type: decimal\n    precision: 20\n    scale: 5\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_DecimalNoPrecision(t *testing.T) {
	yaml := []byte(`fields:\n  - name: price\n    type: decimal\n    scale: 2\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_DecimalInvalidPrecision(t *testing.T) {
	yaml := []byte(`fields:\n  - name: price\n    type: decimal\n    precision: -1\n    scale: 2\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_TimestampInvalidFormat(t *testing.T) {
	yaml := []byte(`fields:\n  - name: created_at\n    type: timestamp\n    format: badformat\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_UnknownType(t *testing.T) {
	yaml := []byte(`fields:\n  - name: foo\n    type: uuid\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_EmptySchema(t *testing.T) {
	yaml := []byte(`fields: []`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_DuplicateFieldNames(t *testing.T) {
	yaml := []byte(`fields:\n  - name: id\n    type: int64\n  - name: id\n    type: string\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_FieldNoName(t *testing.T) {
	yaml := []byte(`fields:\n  - type: int64\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_FieldNoType(t *testing.T) {
	yaml := []byte(`fields:\n  - name: foo\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_UnknownFieldOption(t *testing.T) {
	yaml := []byte(`fields:\n  - name: foo\n    type: int64\n    extra: bar\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_NonBoolNullable(t *testing.T) {
	yaml := []byte(`fields:\n  - name: foo\n    type: int64\n    nullable: "yes"\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_SyntaxError(t *testing.T) {
	yaml := []byte(`fields\n  - name: foo\n    type: int64\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_FieldEmptyName(t *testing.T) {
	yaml := []byte(`fields:\n  - name: \n    type: int64\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_FieldEmptyType(t *testing.T) {
	yaml := []byte(`fields:\n  - name: foo\n    type: \n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_TooManyFields(t *testing.T) {
	yaml := []byte("fields:\n" + strings.Repeat("  - name: f\n    type: int64\n", 1000))
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_LongFieldName(t *testing.T) {
	yaml := []byte(`fields:\n  - name: "` + strings.Repeat("a", 300) + `"\n    type: int64\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_ReservedWordFieldName(t *testing.T) {
	yaml := []byte(`fields:\n  - name: select\n    type: int64\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_AllFieldsNullable(t *testing.T) {
	yaml := []byte(`fields:\n  - name: a\n    type: int64\n    nullable: true\n  - name: b\n    type: string\n    nullable: true\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_AllFieldsRequired(t *testing.T) {
	yaml := []byte(`fields:\n  - name: a\n    type: int64\n    nullable: false\n  - name: b\n    type: string\n    nullable: false\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_FieldsCaseSensitive(t *testing.T) {
	yaml := []byte(`fields:\n  - name: ID\n    type: int64\n  - name: id\n    type: string\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}

func TestLoadSchema_FieldOrderDiffersFromCSV(t *testing.T) {
	yaml := []byte(`fields:\n  - name: b\n    type: string\n  - name: a\n    type: int64\n`)
	schema, err := LoadSchema(yaml)
	require.Error(t, err)
	require.Nil(t, schema)
}
