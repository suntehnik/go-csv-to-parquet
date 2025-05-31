# Архитектура библиотеки CSV to Parquet (Go)

## 1. Общий обзор
Библиотека предназначена для потокового преобразования больших CSV-файлов в формат Parquet с поддержкой всех типов данных Parquet, работы с локальными файлами и S3, а также сжатия (Snappy, Gzip и др.).

---

## 2. Основные компоненты

### 2.1. API контракты

#### Основной публичный API
```go
// Основная точка входа
func Convert(opts Options) error

// Структура параметров
Options struct {
    CsvPath      string // путь к CSV (локально или s3://)
    ParquetPath  string // путь для записи Parquet (локально или s3://)
    SchemaPath   string // путь к YAML-схеме (опционально)
    Compression  string // тип сжатия ("snappy", "gzip", "none", ...)
    BatchSize    int    // размер batch для потоковой обработки (опционально)
}
```

#### Ошибки
- Все ошибки возвращаются через error с подробным описанием причины.

---

### 2.2. Внутренние компоненты

- **CSV Reader**: потоковое чтение CSV (локально или из S3)
- **Schema Detector**: автоматическое определение схемы по первым N строкам CSV
- **Schema Loader**: загрузка и валидация пользовательской схемы (YAML)
- **Type Mapper**: сопоставление типов CSV с типами Parquet
- **Parquet Writer**: потоковая запись данных в Parquet с поддержкой сжатия
- **S3 Adapter**: абстракция для работы с файлами на S3 (чтение и запись)
- **Error Handling**: генерация информативных ошибок

---

## 3. Структура проекта

```
/csv2parquet
    /internal
        csvreader.go       // потоковый CSV-ридер
        schemadetector.go  // детектор схемы
        schemaloader.go    // загрузчик пользовательской схемы
        typemapper.go      // сопоставление типов
        parquetwriter.go   // потоковый Parquet-райтер
        s3adapter.go       // работа с S3
        errors.go          // ошибки
    options.go             // структура Options и публичный API
    convert.go             // glue-код и orchestration
    main_test.go           // тесты
```

---

## 4. Используемые зависимости

- **CSV**: стандартная библиотека Go (`encoding/csv`)
- **Parquet**: [github.com/xitongsys/parquet-go](https://github.com/xitongsys/parquet-go) (или аналогичная, поддерживающая streaming)
- **S3**: [github.com/aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) (работа с S3)
- **YAML**: [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) (парсинг схемы)

---

## 5. Пример публичного API
```go
err := csv2parquet.Convert(csv2parquet.Options{
    CsvPath:     "s3://bucket/input.csv",
    ParquetPath: "output.parquet",
    SchemaPath:  "schema.yaml", // опционально
    Compression: "snappy",      // опционально
    BatchSize:   1000,           // опционально
})
if err != nil {
    // обработка ошибки
}
```

---

## 6. Потоковая обработка
- CSV читается по batch (например, 1000 строк за раз)
- Каждая batch преобразуется и записывается в Parquet
- Память не превышает размер batch + буферов

---

## 7. Расширяемость
- Возможность добавить поддержку других облачных хранилищ (GCS, Azure) через интерфейс StorageAdapter
- Легко расширяемая поддержка новых типов сжатия и форматов схемы

---

## 8. Тестирование
- Unit-тесты для всех компонентов
- Интеграционные тесты для сценариев локальных и S3 файлов
- Тесты на большие файлы и различные типы данных

---

## 9. Ограничения
- Нет поддержки пользовательских маппингов и фильтрации
- Нет логирования внутри библиотеки

---

## 10. Документация
- Примеры использования и описание API в README.md
- Описание формата схемы и поддерживаемых типов
