# Шаг 1. Подготовительный этап

## 1.1. Требования к поддерживаемым типам данных и структурам
- **Flat-структуры** (без вложенности)
- **Nullable** (поддержка null/пустых значений)
- **Точные типы** (int32, int64, float32, float64, bool, string/utf8, date, timestamp, decimal, byte array)
- **Ограничения**:
  - Нет массивов, списков, вложенных структур на первом этапе
  - Все поля — скалярные

## 1.2. Формат пользовательской YAML-схемы
```yaml
fields:
  - name: id
    type: int64
    nullable: false
  - name: name
    type: string
    nullable: true
  - name: created_at
    type: timestamp
    nullable: false
  - name: price
    type: decimal
    precision: 10
    scale: 2
    nullable: true
```
- **Обязательные поля:** name, type
- **Опциональные:** nullable (по умолчанию false), precision/scale (для decimal), format (для date/timestamp)
- **Валидация:** уникальность имён, поддерживаемый тип, корректность параметров

## 1.3. Карта соответствия YAML-типов и Arrow
| YAML type   | Arrow type     | Примечания                       |
|------------|---------------|----------------------------------|
| int32      | arrow.PrimitiveTypes.Int32   |                  |
| int64      | arrow.PrimitiveTypes.Int64   |                  |
| float32    | arrow.PrimitiveTypes.Float32 |                  |
| float64    | arrow.PrimitiveTypes.Float64 |                  |
| bool       | arrow.FixedWidthTypes.Boolean|                  |
| string     | arrow.BinaryTypes.String     | UTF8             |
| date       | arrow.FixedWidthTypes.Date32 | format=YYYY-MM-DD|
| timestamp  | arrow.FixedWidthTypes.Timestamp| format=RFC3339 |
| decimal    | arrow.Decimal128Type         | precision/scale  |
| bytes      | arrow.BinaryTypes.Binary     |                  |

## 1.4. Тестовые сценарии для схемы, маппинга и ошибок

### Валидные сценарии
- Корректная схема с полями всех поддерживаемых типов (int32, int64, float32, float64, bool, string, date, timestamp, decimal, bytes)
- Nullable поля (явно true/false и по умолчанию)
- Decimal c корректными precision/scale
- Timestamp/date с валидным format (если поддерживается)
- Минимальная валидная схема (1 поле)
- Смешанные required/optional (nullable)
- Схема с максимальным допустимым количеством полей
- Decimal с максимальными значениями precision/scale

### Ошибочные сценарии
- Неуникальные имена полей (дублирование)
- Неизвестный/неподдерживаемый тип (например, "uuid", "array")
- Decimal без precision/scale или с невалидными значениями (отрицательные, слишком большие, scale > precision)
- Timestamp/date с невалидным/неподдерживаемым format
- Nullable не boolean (например, строка)
- Пустая схема (нет полей)
- Поле без имени или типа
- Дублирование типа (например, два поля с одинаковым именем и разными типами)
- Неизвестные/неподдерживаемые опции в поле (например, extra: foo)
- Некорректные значения параметров (precision не int, nullable не bool)
- YAML-синтаксис с ошибкой (некорректный отступ, отсутствие ":")
- Поле с пустым именем или типом
- Decimal с отсутствующим/нулевым precision или scale
- Слишком большое число полей (тест на лимиты)
- Слишком длинные имена полей (edge-case)
- Использование зарезервированных слов в имени поля (если есть ограничения)

### Edge-cases
- Все поля nullable (и все required)
- Все поля одного типа (например, только int64)
- Decimal с граничными значениями precision/scale
- Timestamp с разными форматами (если поддерживается)
- Поля с именами, отличающимися только регистром ("ID" vs "id")
- Схема с полями, порядок которых отличается от порядка в CSV (если порядок важен)

---

**Definition of Done:**
- Все требования и ограничения зафиксированы
- Формат YAML-схемы и карта типов определены
- Тестовые сценарии и edge-cases перечислены
