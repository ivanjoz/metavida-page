package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func execQuery[T TableSchemaInterface[T], E TableBaseInterface[T, E]](schemaStruct *T, tableInfo *TableInfo) error {
	table := compileTable[T, E](schemaStruct)
	sql, args := table.selectSQL(tableInfo)
	resp, err := execD1(sql, args)
	if err != nil {
		return err
	}

	refSlice, ok := tableInfo.refSlice.(*[]E)
	if !ok || refSlice == nil {
		return errors.New("query missing result slice")
	}
	rows, err := table.decodeRows(resp)
	if err != nil {
		return err
	}
	*refSlice = append(*refSlice, rows...)
	return nil
}

func (t sqliteTable[T]) selectSQL(info *TableInfo) (string, []any) {
	columns := t.selectColumns(info)
	parts := []string{"SELECT " + quoteIdentList(columns), "FROM " + quoteIdent(t.name)}
	args := []any{}
	whereParts := []string{}
	for _, statement := range info.statements {
		clause, clauseArgs := buildWhereClause(statement)
		whereParts = append(whereParts, clause)
		args = append(args, clauseArgs...)
	}
	if len(whereParts) > 0 {
		parts = append(parts, "WHERE "+strings.Join(whereParts, " AND "))
	}
	if info.orderBy != "" && len(t.keys) > 0 {
		parts = append(parts, "ORDER BY "+quoteIdent(t.keys[len(t.keys)-1].Name)+" "+info.orderBy)
	}
	if info.limit > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", info.limit))
	}
	return strings.Join(parts, " ") + ";", args
}

func (t sqliteTable[T]) getByRecordKey(record *T) (*T, error) {
	info := &TableInfo{limit: 1}
	value := reflect.ValueOf(record).Elem()
	for _, col := range t.whereKeyColumns() {
		info.statements = append(info.statements, ColumnStatement{
			Col:      col.Name,
			Operator: "=",
			Value:    value.FieldByName(col.FieldName).Interface(),
		})
	}
	sql, args := t.selectSQL(info)
	resp, err := execD1(sql, args)
	if err != nil {
		return nil, err
	}
	rows, err := t.decodeRows(resp)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

func (t sqliteTable[T]) decodeRows(resp d1Response) ([]T, error) {
	rows := []T{}
	if len(resp.Result) == 0 {
		return rows, nil
	}
	for _, raw := range resp.Result[0].Results {
		var record T
		value := reflect.ValueOf(&record).Elem()
		for _, col := range t.columns {
			rawValue, ok := raw[col.Name]
			if !ok || rawValue == nil {
				continue
			}
			field := value.FieldByName(col.FieldName)
			if !field.CanSet() {
				continue
			}
			if err := setSQLiteValue(field, rawValue); err != nil {
				return nil, fmt.Errorf("column %s: %w", col.Name, err)
			}
		}
		rows = append(rows, record)
	}
	return rows, nil
}

func (t sqliteTable[T]) selectColumns(info *TableInfo) []string {
	if len(info.columnsInclude) > 0 {
		names := make([]string, 0, len(info.columnsInclude))
		for _, col := range info.columnsInclude {
			names = append(names, col.Name)
		}
		return names
	}
	excluded := map[string]bool{}
	for _, col := range info.columnsExclude {
		excluded[col.Name] = true
	}
	names := make([]string, 0, len(t.columns))
	for _, col := range t.columns {
		if !excluded[col.Name] {
			names = append(names, col.Name)
		}
	}
	return names
}

func buildWhereClause(statement ColumnStatement) (string, []any) {
	col := quoteIdent(statement.Col)
	switch strings.ToUpper(statement.Operator) {
	case "IN":
		placeholders := make([]string, len(statement.Values))
		args := make([]any, 0, len(statement.Values))
		for i, value := range statement.Values {
			placeholders[i] = "?"
			args = append(args, normalizeAnySQLiteValue(value))
		}
		return fmt.Sprintf("%s IN (%s)", col, strings.Join(placeholders, ", ")), args
	case "BETWEEN":
		return fmt.Sprintf("%s BETWEEN ? AND ?", col), []any{
			normalizeAnySQLiteValue(statement.From[0].Value),
			normalizeAnySQLiteValue(statement.To[0].Value),
		}
	case "CONTAINS":
		return fmt.Sprintf("%s LIKE ?", col), []any{"%" + fmt.Sprint(statement.Value) + "%"}
	default:
		op := statement.Operator
		if op == "" {
			op = "="
		}
		return fmt.Sprintf("%s %s ?", col, op), []any{normalizeAnySQLiteValue(statement.Value)}
	}
}
