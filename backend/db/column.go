package db

type Col[T any, E any] struct {
	info         columnInfo
	schemaStruct *T
	tableInfo    *TableInfo
}

func (c Col[T, E]) GetInfo() columnInfo {
	return c.info
}

func (c Col[T, E]) GetName() string {
	return c.info.Name
}

func (c *Col[T, E]) setInfo(info columnInfo) {
	c.info = info
}

func (c *Col[T, E]) setTableInfo(tableInfo *TableInfo) {
	c.tableInfo = tableInfo
}

func (c *Col[T, E]) setSchemaStruct(schemaStruct any) {
	if typed, ok := schemaStruct.(*T); ok {
		c.schemaStruct = typed
	}
}

func (c *Col[T, E]) Equals(v E) *T {
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{Col: c.info.Name, Operator: "=", Value: v})
	return c.schemaStruct
}

func (c *Col[T, E]) In(values ...E) *T {
	items := make([]any, 0, len(values))
	for _, value := range values {
		items = append(items, value)
	}
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{Col: c.info.Name, Operator: "IN", Values: items})
	return c.schemaStruct
}

func (c *Col[T, E]) GreaterThan(v E) *T {
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{Col: c.info.Name, Operator: ">", Value: v})
	return c.schemaStruct
}

func (c *Col[T, E]) GreaterEqual(v E) *T {
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{Col: c.info.Name, Operator: ">=", Value: v})
	return c.schemaStruct
}

func (c *Col[T, E]) LessThan(v E) *T {
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{Col: c.info.Name, Operator: "<", Value: v})
	return c.schemaStruct
}

func (c *Col[T, E]) LessEqual(v E) *T {
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{Col: c.info.Name, Operator: "<=", Value: v})
	return c.schemaStruct
}

func (c *Col[T, E]) Between(from E, to E) *T {
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{
		Col:      c.info.Name,
		Operator: "BETWEEN",
		From:     []ColumnStatement{{Col: c.info.Name, Value: from}},
		To:       []ColumnStatement{{Col: c.info.Name, Value: to}},
	})
	return c.schemaStruct
}

type ColSlice[T any, E any] struct {
	info         columnInfo
	schemaStruct *T
	tableInfo    *TableInfo
}

func (c ColSlice[T, E]) GetInfo() columnInfo {
	return c.info
}

func (c ColSlice[T, E]) GetName() string {
	return c.info.Name
}

func (c *ColSlice[T, E]) setInfo(info columnInfo) {
	c.info = info
}

func (c *ColSlice[T, E]) setTableInfo(tableInfo *TableInfo) {
	c.tableInfo = tableInfo
}

func (c *ColSlice[T, E]) setSchemaStruct(schemaStruct any) {
	if typed, ok := schemaStruct.(*T); ok {
		c.schemaStruct = typed
	}
}

func (c *ColSlice[T, E]) Contains(v E) *T {
	c.tableInfo.statements = append(c.tableInfo.statements, ColumnStatement{Col: c.info.Name, Operator: "CONTAINS", Value: v})
	return c.schemaStruct
}
