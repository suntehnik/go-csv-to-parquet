package csvreader

import (
	"io"
	"strings"
)

// CSVReader реализует чтение CSV и преобразование в строки для Arrow
// (минимально для компиляции и тестов)
type CSVReader struct {
	r io.Reader
}

func NewCSVReader(r io.Reader) *CSVReader {
	return &CSVReader{r: r}
}

// ReadBatch читает batch строк из CSV (расширенная реализация: BOM, пустые строки, проверка количества колонок)
func (c *CSVReader) ReadBatch(batchSize int) ([][]string, error) {
	if c.r == nil || batchSize <= 0 {
		return nil, ErrNoData
	}
	buf := make([]byte, 4096)
	n, err := c.r.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if n == 0 {
		return nil, ErrNoData
	}
	lines := strings.Split(string(buf[:n]), "\n")
	// Пропуск BOM и пустых строк в начале
	header := ""
	var headerCols []string
	rows := [][]string{}
	sep := ";"
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		if header == "" {
			if strings.HasPrefix(l, "\uFEFF") {
				l = strings.TrimPrefix(l, "\uFEFF")
			}
			header = l
			headerCols, err = parseCSVRow(l, sep)
			if err != nil {
				return nil, err
			}
			if len(headerCols) == 1 {
				return nil, ErrUnsupportedSeparator
			}
			continue
		}
		cols, err := parseCSVRow(l, sep)
		if err != nil {
			return nil, err
		}
		if len(cols) != len(headerCols) {
			return nil, ErrColMismatch
		}
		rows = append(rows, cols)
		if len(rows) == batchSize {
			break
		}
	}
	if len(headerCols) == 0 || len(rows) == 0 {
		return nil, ErrNoData
	}
	// Если batch меньше batchSize и не EOF, это PartialBatchAtEOF (ошибка)
	if len(rows) < batchSize && err != io.EOF && len(rows) > 0 {
		return nil, ErrPartialBatch
	}
	return rows, nil
}

var ErrNoData = io.EOF
var ErrColMismatch = io.ErrUnexpectedEOF
var ErrPartialBatch = io.ErrShortBuffer

// parseCSVRow разбирает строку CSV с поддержкой кавычек и экранирования по RFC4180
func parseCSVRow(line, sep string) ([]string, error) {
	var res []string
	var cur strings.Builder
	inQuotes := false
	i := 0
	for i < len(line) {
		c := line[i]
		if c == '"' {
			if inQuotes && i+1 < len(line) && line[i+1] == '"' {
				cur.WriteByte('"')
				i++ // пропустить экранированную кавычку
			} else {
				inQuotes = !inQuotes
			}
		} else if !inQuotes && strings.HasPrefix(line[i:], sep) {
			res = append(res, cur.String())
			cur.Reset()
			i += len(sep) - 1
		} else {
			cur.WriteByte(c)
		}
		i++
	}
	res = append(res, cur.String())
	if inQuotes {
		return nil, ErrInvalidQuotes
	}
	return res, nil
}

var ErrInvalidQuotes = io.ErrUnexpectedEOF
var ErrUnsupportedSeparator = io.ErrClosedPipe
