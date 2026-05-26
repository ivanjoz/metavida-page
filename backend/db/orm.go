package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"metavida/backend/config"
)

type ORM[T any] interface {
	Init() error
	Insert(record T) error
	GetByID(id string) (*T, error)
	Select() QueryBuilder[T]
}

type QueryBuilder[T any] interface {
	Where(column string) QueryBuilder[T]
	Equals(value any) QueryBuilder[T]
	Limit(value int) QueryBuilder[T]
	Exec() ([]T, error)
}

func NewORM[T any](cfg config.Config, tableName string) (ORM[T], error) {
	switch cfg.CloudProvider {
	case "", "cloudflare":
		return NewD1ORM[T](cfg, tableName), nil
	default:
		return nil, fmt.Errorf("unsupported cloud provider %q", cfg.CloudProvider)
	}
}

type D1ORM[T any] struct {
	cfg       config.Config
	tableName string
	client    *http.Client
}

func NewD1ORM[T any](cfg config.Config, tableName string) *D1ORM[T] {
	return &D1ORM[T]{
		cfg:       cfg,
		tableName: tableName,
		client:    &http.Client{Timeout: 20 * time.Second},
	}
}

func (o *D1ORM[T]) Init() error {
	schema := CreateTableSQL[T](o.tableName)
	_, err := o.execSQL(schema, nil)
	return err
}

func (o *D1ORM[T]) Insert(record T) error {
	cols, values := ColumnsAndValues(record)
	if len(cols) == 0 {
		return errors.New("record has no db columns")
	}

	placeholders := make([]string, len(cols))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	sql := fmt.Sprintf(
		"INSERT OR REPLACE INTO %s (%s) VALUES (%s)",
		o.tableName,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
	_, err := o.execSQL(sql, values)
	return err
}

func (o *D1ORM[T]) GetByID(id string) (*T, error) {
	rows, err := o.Select().Where("id").Equals(id).Limit(1).Exec()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (o *D1ORM[T]) Select() QueryBuilder[T] {
	return &d1QueryBuilder[T]{orm: o, limit: 100}
}

type d1QueryBuilder[T any] struct {
	orm       *D1ORM[T]
	whereCol  string
	whereOp   string
	whereArgs []any
	limit     int
}

func (b *d1QueryBuilder[T]) Where(column string) QueryBuilder[T] {
	b.whereCol = column
	return b
}

func (b *d1QueryBuilder[T]) Equals(value any) QueryBuilder[T] {
	b.whereOp = "="
	b.whereArgs = []any{value}
	return b
}

func (b *d1QueryBuilder[T]) Limit(value int) QueryBuilder[T] {
	b.limit = value
	return b
}

func (b *d1QueryBuilder[T]) Exec() ([]T, error) {
	sql := fmt.Sprintf("SELECT * FROM %s", b.orm.tableName)
	args := b.whereArgs
	if b.whereCol != "" {
		sql += fmt.Sprintf(" WHERE %s %s ?", b.whereCol, b.whereOp)
	}
	if b.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", b.limit)
	}

	resp, err := b.orm.execSQL(sql, args)
	if err != nil {
		return nil, err
	}
	return DecodeD1Rows[T](resp)
}

type d1QueryRequest struct {
	SQL    string `json:"sql"`
	Params []any  `json:"params,omitempty"`
}

type D1Response struct {
	Success bool `json:"success"`
	Result  []struct {
		Results []map[string]any `json:"results"`
		Success bool             `json:"success"`
		Error   string           `json:"error"`
	} `json:"result"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (o *D1ORM[T]) execSQL(sql string, params []any) (D1Response, error) {
	var empty D1Response
	if o.cfg.CloudflareAccountID == "" || o.cfg.CloudflareAPIToken == "" || o.cfg.CloudflareDatabaseID == "" {
		return empty, errors.New("missing Cloudflare D1 credentials in credentials.json")
	}

	payload, err := json.Marshal(d1QueryRequest{SQL: sql, Params: params})
	if err != nil {
		return empty, err
	}

	url := fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/query",
		o.cfg.CloudflareAccountID,
		o.cfg.CloudflareDatabaseID,
	)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return empty, err
	}
	req.Header.Set("Authorization", "Bearer "+o.cfg.CloudflareAPIToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := o.client.Do(req)
	if err != nil {
		return empty, err
	}
	defer res.Body.Close()

	var out D1Response
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return empty, err
	}
	if !resOK(res.StatusCode) || !out.Success {
		if len(out.Errors) > 0 {
			return empty, errors.New(out.Errors[0].Message)
		}
		return empty, fmt.Errorf("cloudflare d1 query failed with status %d", res.StatusCode)
	}
	return out, nil
}

func resOK(status int) bool {
	return status >= 200 && status < 300
}

func CreateTableSQL[T any](table string) string {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	columns := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		col := columnName(field)
		if col == "" {
			continue
		}
		sqlType := sqliteType(field.Type)
		definition := fmt.Sprintf("%s %s", col, sqlType)
		if col == "id" {
			definition += " PRIMARY KEY"
		}
		columns = append(columns, definition)
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", table, strings.Join(columns, ", "))
}

func ColumnsAndValues(record any) ([]string, []any) {
	value := reflect.ValueOf(record)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	typ := value.Type()

	cols := make([]string, 0, value.NumField())
	values := make([]any, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		col := columnName(field)
		if col == "" {
			continue
		}
		cols = append(cols, col)
		values = append(values, value.Field(i).Interface())
	}
	return cols, values
}

func DecodeD1Rows[T any](resp D1Response) ([]T, error) {
	rows := []T{}
	if len(resp.Result) == 0 {
		return rows, nil
	}
	for _, raw := range resp.Result[0].Results {
		bytes, err := json.Marshal(raw)
		if err != nil {
			return nil, err
		}
		var row T
		if err := json.Unmarshal(bytes, &row); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func columnName(field reflect.StructField) string {
	if field.PkgPath != "" {
		return ""
	}
	tag := field.Tag.Get("db")
	if tag == "-" {
		return ""
	}
	if tag != "" {
		return strings.Split(tag, ",")[0]
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		return strings.Split(jsonTag, ",")[0]
	}
	return strings.ToLower(field.Name)
}

func sqliteType(typ reflect.Type) string {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Bool:
		return "INTEGER"
	default:
		return "TEXT"
	}
}
