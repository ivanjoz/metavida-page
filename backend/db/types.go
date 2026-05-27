package db

import "reflect"

const TypeLocalIndex int8 = 0

type TableSchema struct {
	Name      string
	Keys      []Coln
	Partition Coln
	Indexes   []Index
}

type Index struct {
	Type     int8
	Keys     []Coln
	Cols     []Coln
	KeepPart bool
}

type ColumnStatement struct {
	Col      string
	Operator string
	Value    any
	Values   []any
	From     []ColumnStatement
	To       []ColumnStatement
}

type columnInfo struct {
	Name      string
	FieldName string
	FieldType reflect.Type
}

type Coln interface {
	GetInfo() columnInfo
	GetName() string
}

type columnBinder interface {
	Coln
	setInfo(columnInfo)
	setTableInfo(*TableInfo)
	setSchemaStruct(any)
}

type TableSchemaInterface[T any] interface {
	GetSchema() TableSchema
}

type TableBaseInterface[T any, E any] interface {
	GetBaseStruct() E
	GetTableStruct() T
}

type TableStructInterfaceQuery[T any, E any] interface {
	SetRefSlice(*[]E)
}

type TableDeployInterface interface {
	TableName() string
	BuildDeploySQL() []string
	BuildDeployPlan() tableDeployPlan
}

type tableDeployPlan struct {
	name              string
	columns           []columnInfo
	primaryKeyColumns []string
	indexes           []indexDeployPlan
}

type indexDeployPlan struct {
	name    string
	columns []string
	sql     string
}

type TableInfo struct {
	statements     []ColumnStatement
	columnsInclude []columnInfo
	columnsExclude []columnInfo
	orderBy        string
	limit          int32
	refSlice       any
}

type sqliteTable[T any] struct {
	name      string
	columns   []columnInfo
	columnMap map[string]columnInfo
	keys      []columnInfo
	partition *columnInfo
	indexes   []Index
}
