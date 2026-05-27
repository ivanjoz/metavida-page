package db

type TableStruct[T TableSchemaInterface[T], E TableBaseInterface[T, E]] struct {
	schemaStruct *T
	tableInfo    *TableInfo
}

func (t TableStruct[T, E]) GetBaseStruct() E {
	return *new(E)
}

func (t TableStruct[T, E]) GetTableStruct() T {
	return *new(T)
}

func (t *TableStruct[T, E]) SetRefSlice(refSlice *[]E) {
	t.tableInfo.refSlice = refSlice
}

func (t *TableStruct[T, E]) Select(columns ...Coln) *T {
	for _, col := range columns {
		t.tableInfo.columnsInclude = append(t.tableInfo.columnsInclude, col.GetInfo())
	}
	return t.schemaStruct
}

func (t *TableStruct[T, E]) Exclude(columns ...Coln) *T {
	for _, col := range columns {
		t.tableInfo.columnsExclude = append(t.tableInfo.columnsExclude, col.GetInfo())
	}
	return t.schemaStruct
}

func (t *TableStruct[T, E]) Limit(limit int32) *T {
	t.tableInfo.limit = limit
	return t.schemaStruct
}

func (t *TableStruct[T, E]) OrderDesc() *T {
	t.tableInfo.orderBy = "DESC"
	return t.schemaStruct
}

func (t *TableStruct[T, E]) Exec() error {
	return execQuery[T, E](t.schemaStruct, t.tableInfo)
}

func (t *TableStruct[T, E]) Insert(records *[]E, columnsToExclude ...Coln) error {
	return Insert(records, columnsToExclude...)
}

func (t *TableStruct[T, E]) InsertOne(record E, columnsToExclude ...Coln) error {
	return InsertOne(record, columnsToExclude...)
}

func (t *TableStruct[T, E]) Update(records *[]E, columnsToInclude ...Coln) error {
	return Update(records, columnsToInclude...)
}

func (t *TableStruct[T, E]) UpdateOne(record E, columnsToInclude ...Coln) error {
	return UpdateOne(record, columnsToInclude...)
}

func (t *TableStruct[T, E]) TableName() string {
	return compileTable[T, E](new(T)).name
}

func (t *TableStruct[T, E]) BuildDeploySQL() []string {
	table := compileTable[T, E](new(T))
	return table.deploySQL()
}

func (t *TableStruct[T, E]) BuildDeployPlan() tableDeployPlan {
	table := compileTable[T, E](new(T))
	return table.deployPlan()
}

func (t *TableStruct[T, E]) setQueryContext(schemaStruct any, tableInfo *TableInfo) {
	if typed, ok := schemaStruct.(*T); ok {
		t.schemaStruct = typed
	}
	t.tableInfo = tableInfo
}

func Query[T TableBaseInterface[E, T], E TableSchemaInterface[E]](refSlice *[]T) *E {
	refTable := initStructTable[E, T](new(E))
	any(refTable).(TableStructInterfaceQuery[E, T]).SetRefSlice(refSlice)
	return refTable
}

func Table[T TableBaseInterface[E, T], E TableSchemaInterface[E]]() *E {
	return initStructTable[E, T](new(E))
}
