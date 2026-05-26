package db

import (
	"reflect"
	"regexp"
	"strings"
)

func columnName(field reflect.StructField) string {
	tag := field.Tag.Get("db")
	if tag != "" {
		return strings.Split(tag, ",")[0]
	}
	return toSnakeCase(field.Name)
}

var snakeRegexp = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(value string) string {
	return strings.ToLower(snakeRegexp.ReplaceAllString(value, "${1}_${2}"))
}

func quoteIdent(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func quoteIdentList(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, quoteIdent(value))
	}
	return strings.Join(quoted, ", ")
}

func colSet(columns []Coln) map[string]bool {
	set := map[string]bool{}
	for _, col := range columns {
		set[col.GetName()] = true
	}
	return set
}

type infoCol columnInfo

func (c infoCol) GetInfo() columnInfo { return columnInfo(c) }
func (c infoCol) GetName() string     { return c.Name }

func infosToColn(columns []columnInfo) []Coln {
	out := make([]Coln, 0, len(columns))
	for _, col := range columns {
		out = append(out, infoCol(col))
	}
	return out
}
