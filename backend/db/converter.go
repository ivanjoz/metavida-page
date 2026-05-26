package db

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func normalizeSQLiteValue(value reflect.Value) (any, error) {
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, nil
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return int64(1), nil
		}
		return int64(0), nil
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return value.Bytes(), nil
		}
		fallthrough
	case reflect.Struct, reflect.Map, reflect.Array:
		bytes, err := json.Marshal(value.Interface())
		if err != nil {
			return nil, err
		}
		return string(bytes), nil
	default:
		return value.Interface(), nil
	}
}

func normalizeAnySQLiteValue(value any) any {
	ref := reflect.ValueOf(value)
	if !ref.IsValid() {
		return nil
	}
	normalized, err := normalizeSQLiteValue(ref)
	if err != nil {
		return fmt.Sprint(value)
	}
	return normalized
}

func setSQLiteValue(field reflect.Value, value any) error {
	if field.Kind() == reflect.Pointer {
		if value == nil {
			return nil
		}
		field.Set(reflect.New(field.Type().Elem()))
		return setSQLiteValue(field.Elem(), value)
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(fmt.Sprint(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(toInt64(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(uint64(toInt64(value)))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(toFloat64(value))
	case reflect.Bool:
		field.SetBool(toInt64(value) != 0 || value == true)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			if s, ok := value.(string); ok {
				field.SetBytes([]byte(s))
			}
			return nil
		}
		if s, ok := value.(string); ok {
			return json.Unmarshal([]byte(s), field.Addr().Interface())
		}
	case reflect.Struct, reflect.Map, reflect.Array:
		if s, ok := value.(string); ok {
			return json.Unmarshal([]byte(s), field.Addr().Interface())
		}
	}
	return nil
}

func sqliteType(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "BLOB"
		}
		return "TEXT"
	case reflect.Struct, reflect.Map, reflect.Array:
		return "TEXT"
	default:
		return "TEXT"
	}
}

func toInt64(value any) int64 {
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int8:
		return int64(typed)
	case int16:
		return int64(typed)
	case int32:
		return int64(typed)
	case int64:
		return typed
	case uint:
		return int64(typed)
	case uint8:
		return int64(typed)
	case uint16:
		return int64(typed)
	case uint32:
		return int64(typed)
	case uint64:
		return int64(typed)
	case float64:
		return int64(typed)
	case json.Number:
		n, _ := typed.Int64()
		return n
	}
	return 0
}

func toFloat64(value any) float64 {
	switch typed := value.(type) {
	case float32:
		return float64(typed)
	case float64:
		return typed
	case json.Number:
		n, _ := typed.Float64()
		return n
	}
	return float64(toInt64(value))
}
