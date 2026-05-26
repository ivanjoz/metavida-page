package db

import (
	"fmt"
	"reflect"
	"strings"
)

func initStructTable[T TableSchemaInterface[T], E any](schemaStruct *T) *T {
	tableInfo := &TableInfo{}
	bindColumns(schemaStruct, tableInfo, baseRecordType[E]())
	return schemaStruct
}

func compileTable[T TableSchemaInterface[T], E any](schemaStruct *T) sqliteTable[E] {
	recordType := baseRecordType[E]()
	bindColumns(schemaStruct, &TableInfo{}, recordType)
	schema := (*schemaStruct).GetSchema()
	tableName := schema.Name
	if tableName == "" {
		tableName = toSnakeCase(recordType.Name())
	}

	table := sqliteTable[E]{
		name:      tableName,
		columns:   collectColumns(schemaStruct),
		columnMap: map[string]columnInfo{},
		indexes:   schema.Indexes,
	}
	for _, col := range table.columns {
		table.columnMap[col.Name] = col
	}
	for _, key := range schema.Keys {
		info := key.GetInfo()
		if col, ok := table.columnMap[info.Name]; ok {
			table.keys = append(table.keys, col)
		} else {
			panic(fmt.Sprintf("schema key column %q is not declared on table %s", info.Name, tableName))
		}
	}
	if schema.Partition != nil {
		info := schema.Partition.GetInfo()
		if col, ok := table.columnMap[info.Name]; ok {
			table.partition = &col
		} else {
			panic(fmt.Sprintf("schema partition column %q is not declared on table %s", info.Name, tableName))
		}
	}
	return table
}

func bindColumns(schemaStruct any, tableInfo *TableInfo, recordType reflect.Type) {
	schemaValue := reflect.ValueOf(schemaStruct).Elem()
	schemaType := schemaValue.Type()
	recordFields := recordFieldMap(recordType)

	for i := 0; i < schemaValue.NumField(); i++ {
		field := schemaValue.Field(i)
		if !field.CanAddr() || !field.Addr().CanInterface() {
			continue
		}

		if tableStruct, ok := field.Addr().Interface().(interface {
			setQueryContext(any, *TableInfo)
		}); ok {
			tableStruct.setQueryContext(schemaStruct, tableInfo)
			continue
		}

		binder, ok := field.Addr().Interface().(columnBinder)
		if !ok {
			continue
		}
		structField := schemaType.Field(i)
		recordField, exists := recordFields[structField.Name]
		if !exists {
			panic(fmt.Sprintf("table field %s has no matching record field", structField.Name))
		}
		expectedTypeName := genericColumnValueTypeName(field.Type())
		if expectedTypeName != "" {
			validateColumnType(structField.Name, recordField.Type, expectedTypeName, strings.Contains(field.Type().String(), "ColSlice["))
		}
		binder.setInfo(columnInfo{
			Name:      columnName(recordField),
			FieldName: recordField.Name,
			FieldType: recordField.Type,
		})
		binder.setTableInfo(tableInfo)
		binder.setSchemaStruct(schemaStruct)
	}
}

func collectColumns(schemaStruct any) []columnInfo {
	value := reflect.ValueOf(schemaStruct).Elem()
	columns := []columnInfo{}
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).CanAddr() || !value.Field(i).Addr().CanInterface() {
			continue
		}
		col, ok := value.Field(i).Addr().Interface().(Coln)
		if !ok {
			continue
		}
		info := col.GetInfo()
		if info.Name != "" {
			columns = append(columns, info)
		}
	}
	return columns
}

func recordFieldMap(recordType reflect.Type) map[string]reflect.StructField {
	fields := map[string]reflect.StructField{}
	for i := 0; i < recordType.NumField(); i++ {
		field := recordType.Field(i)
		if field.Anonymous || field.PkgPath != "" || field.Tag.Get("db") == "-" {
			continue
		}
		fields[field.Name] = field
	}
	return fields
}

func baseRecordType[E any]() reflect.Type {
	t := reflect.TypeOf((*E)(nil)).Elem()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func genericColumnValueTypeName(t reflect.Type) string {
	name := t.String()
	open := strings.LastIndex(name, ",")
	close := strings.LastIndex(name, "]")
	if open < 0 || close < 0 || close <= open+1 {
		return ""
	}
	return strings.TrimSpace(name[open+1 : close])
}

func validateColumnType(fieldName string, actual reflect.Type, expectedTypeName string, isSliceColumn bool) {
	for actual.Kind() == reflect.Pointer {
		actual = actual.Elem()
	}
	if isSliceColumn {
		if actual.Kind() != reflect.Slice || !typeNameMatches(actual.Elem(), expectedTypeName) {
			panic(fmt.Sprintf("table field %s expects []%s record field, found %s", fieldName, expectedTypeName, actual))
		}
		return
	}
	if !typeNameMatches(actual, expectedTypeName) {
		panic(fmt.Sprintf("table field %s expects %s record field, found %s", fieldName, expectedTypeName, actual))
	}
}

func typeNameMatches(actual reflect.Type, expected string) bool {
	if actual.String() == expected {
		return true
	}
	if actual.PkgPath() != "" && actual.PkgPath()+"."+actual.Name() == expected {
		return true
	}
	return false
}
