package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Insert[T TableBaseInterface[E, T], E TableSchemaInterface[E]](records *[]T, columnsToExclude ...Coln) error {
	if records == nil || len(*records) == 0 {
		return nil
	}
	table := compileTable[E, T](new(E))
	excluded := colSet(columnsToExclude)
	columns := make([]columnInfo, 0, len(table.columns))
	for _, col := range table.columns {
		if !excluded[col.Name] {
			columns = append(columns, col)
		}
	}

	for i := range *records {
		sql, args, err := table.insertSQL(&(*records)[i], columns)
		if err != nil {
			return err
		}
		if _, err := execD1(sql, args); err != nil {
			return err
		}
	}
	return nil
}

func InsertOne[T TableBaseInterface[E, T], E TableSchemaInterface[E]](record T, columnsToExclude ...Coln) error {
	return Insert(&[]T{record}, columnsToExclude...)
}

func Update[T TableBaseInterface[E, T], E TableSchemaInterface[E]](records *[]T, columnsToInclude ...Coln) error {
	if len(columnsToInclude) == 0 {
		return errors.New("no columns included for update")
	}
	if records == nil || len(*records) == 0 {
		return nil
	}
	table := compileTable[E, T](new(E))
	columns := make([]columnInfo, 0, len(columnsToInclude))
	for _, col := range columnsToInclude {
		columns = append(columns, col.GetInfo())
	}
	for i := range *records {
		sql, args, err := table.updateSQL(&(*records)[i], columns)
		if err != nil {
			return err
		}
		if _, err := execD1(sql, args); err != nil {
			return err
		}
	}
	return nil
}

func UpdateOne[T TableBaseInterface[E, T], E TableSchemaInterface[E]](record T, columnsToInclude ...Coln) error {
	return Update(&[]T{record}, columnsToInclude...)
}

func UpdateExclude[T TableBaseInterface[E, T], E TableSchemaInterface[E]](records *[]T, columnsToExclude ...Coln) error {
	if records == nil || len(*records) == 0 {
		return nil
	}
	table := compileTable[E, T](new(E))
	excluded := colSet(columnsToExclude)
	columns := make([]columnInfo, 0, len(table.columns))
	keyCols := table.keyColumnSet()
	for _, col := range table.columns {
		if !excluded[col.Name] && !keyCols[col.Name] {
			columns = append(columns, col)
		}
	}
	return Update(records, infosToColn(columns)...)
}

func (t sqliteTable[T]) insertSQL(record *T, columns []columnInfo) (string, []any, error) {
	names := make([]string, 0, len(columns))
	placeholders := make([]string, 0, len(columns))
	args := make([]any, 0, len(columns))
	value := reflect.ValueOf(record).Elem()
	for _, col := range columns {
		field := value.FieldByName(col.FieldName)
		normalized, err := normalizeSQLiteValue(field)
		if err != nil {
			return "", nil, fmt.Errorf("column %s: %w", col.Name, err)
		}
		names = append(names, quoteIdent(col.Name))
		placeholders = append(placeholders, "?")
		args = append(args, normalized)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO UPDATE SET %s;",
		quoteIdent(t.name),
		strings.Join(names, ", "),
		strings.Join(placeholders, ", "),
		t.upsertAssignments(columns),
	)
	return sql, args, nil
}

func (t sqliteTable[T]) updateSQL(record *T, columns []columnInfo) (string, []any, error) {
	if len(columns) == 0 {
		return "", nil, errors.New("no update columns")
	}
	keyColumns := t.whereKeyColumns()
	if len(keyColumns) == 0 {
		return "", nil, errors.New("update requires partition or key columns")
	}
	value := reflect.ValueOf(record).Elem()
	setParts := make([]string, 0, len(columns))
	whereParts := make([]string, 0, len(keyColumns))
	args := make([]any, 0, len(columns)+len(keyColumns))
	for _, col := range columns {
		field := value.FieldByName(col.FieldName)
		normalized, err := normalizeSQLiteValue(field)
		if err != nil {
			return "", nil, fmt.Errorf("column %s: %w", col.Name, err)
		}
		setParts = append(setParts, fmt.Sprintf("%s = ?", quoteIdent(col.Name)))
		args = append(args, normalized)
	}
	for _, col := range keyColumns {
		field := value.FieldByName(col.FieldName)
		normalized, err := normalizeSQLiteValue(field)
		if err != nil {
			return "", nil, fmt.Errorf("key column %s: %w", col.Name, err)
		}
		whereParts = append(whereParts, fmt.Sprintf("%s = ?", quoteIdent(col.Name)))
		args = append(args, normalized)
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s;", quoteIdent(t.name), strings.Join(setParts, ", "), strings.Join(whereParts, " AND ")), args, nil
}

func (t sqliteTable[T]) upsertAssignments(columns []columnInfo) string {
	keyCols := t.keyColumnSet()
	assignments := []string{}
	for _, col := range columns {
		if keyCols[col.Name] {
			continue
		}
		assignments = append(assignments, fmt.Sprintf("%s = excluded.%s", quoteIdent(col.Name), quoteIdent(col.Name)))
	}
	if len(assignments) == 0 {
		return fmt.Sprintf("%s = %s", quoteIdent(columns[0].Name), quoteIdent(columns[0].Name))
	}
	return strings.Join(assignments, ", ")
}
