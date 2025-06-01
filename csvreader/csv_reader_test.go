package csvreader

import (
	"io"
	"strings"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestCSVReader_ValidSimple(t *testing.T) {
	csv := "id;name;age\n1;Alice;30\n2;Bob;25\n3;Carol;22\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.NoError(t, err)
	require.NotNil(t, batch)
	require.Equal(t, [][]string{{"1", "Alice", "30"}, {"2", "Bob", "25"}}, batch)
}

func TestCSVReader_EmptyFile(t *testing.T) {
	r := strings.NewReader("")
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batch)
}

func TestCSVReader_HeaderOnly(t *testing.T) {
	csv := "id;name;age\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batch)
}

func TestCSVReader_EmptyLines(t *testing.T) {
	csv := "id;name;age\n\n\n1;Alice;30\n\n2;Bob;25\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.NoError(t, err)
	require.NotNil(t, batch)
	require.Equal(t, [][]string{{"1", "Alice", "30"}, {"2", "Bob", "25"}}, batch)
}

func TestCSVReader_UnicodeEmoji(t *testing.T) {
	csv := "id;name\n1;😀\n2;Привет\n3;こんにちは\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.NoError(t, err)
	require.NotNil(t, batch)
	require.Equal(t, [][]string{{"1", "😀"}, {"2", "Привет"}}, batch)
}

func TestCSVReader_LongLines(t *testing.T) {
	csv := "id;name\n" + strings.Repeat("1;"+strings.Repeat("a", 1000)+"\n", 10)
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(5)
	require.NoError(t, err)
	require.NotNil(t, batch)
}

func TestCSVReader_DifferentSeparators(t *testing.T) {
	csv := "id,name,age\n1,Alice,30\n2,Bob,25\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batch)

	csvTab := "id\tage\n1\t30\n2\t25\n"
	rTab := strings.NewReader(csvTab)
	readerTab := NewCSVReader(rTab)
	batchTab, err := readerTab.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batchTab)
}

func TestCSVReader_QuotedAndEscaped(t *testing.T) {
	csv := "id;name\n1;\"Alice; the \"\"Great\"\"\"\n2;\"Bob\"\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.NoError(t, err)
	require.NotNil(t, batch)
}

func TestCSVReader_BOM(t *testing.T) {
	csv := "\uFEFFid;name\n1;Alice\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(1)
	require.NoError(t, err)
	require.NotNil(t, batch)
}

func TestCSVReader_ColumnMismatch(t *testing.T) {
	csv := "id;name\n1\n2;Bob;25\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batch)
}

func TestCSVReader_InvalidCSVFormat(t *testing.T) {
	csv := "id;name\n1;Alice\n2;Bob;25\n3;\"Carol\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batch)
}

func TestCSVReader_PartialBatchAtEOF(t *testing.T) {
	csv := "id;name\n1;Alice\n2;Bob\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(3)
	require.Error(t, err)
	require.Nil(t, batch)
}

func TestCSVReader_BatchFromEmptyLines(t *testing.T) {
	csv := "id;name\n\n\n\n"
	r := strings.NewReader(csv)
	reader := NewCSVReader(r)
	batch, err := reader.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batch)
}

func TestCSVReader_IOError(t *testing.T) {
	// Эмулируем io.ErrUnexpectedEOF через readerWithFunc
	readFunc := func(p []byte) (int, error) { return 0, io.ErrUnexpectedEOF }
	reader := NewCSVReader(readerWithFunc(readFunc))
	batch, err := reader.ReadBatch(2)
	require.Error(t, err)
	require.Nil(t, batch)
}

type readerWithFunc func([]byte) (int, error)

func (f readerWithFunc) Read(p []byte) (int, error) {
	return f(p)
}

