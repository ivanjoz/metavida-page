package db

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"

	"metavida/backend/config"
)

func DeployTables(cfg config.Config, tables ...TableDeployInterface) error {
	Configure(cfg)
	plans := make([]tableDeployPlan, 0, len(tables))
	for _, table := range tables {
		plans = append(plans, table.BuildDeployPlan())
	}

	log.Printf("fetching current Cloudflare D1 schema ...")
	currentSchema, err := fetchDeployDatabaseSchema(plans)
	if err != nil {
		return err
	}
	log.Printf("loaded %d target table schemas from Cloudflare D1", len(currentSchema.tables))

	for _, plan := range plans {
		log.Printf("checking table %q ...", plan.name)
		changeCount, err := deployTablePlan(plan, currentSchema)
		if err != nil {
			return fmt.Errorf("table %q: %w", plan.name, err)
		}
		if changeCount == 0 {
			log.Printf("table %q: up to date", plan.name)
			continue
		}
		log.Printf("table %q: applied %d schema changes", plan.name, changeCount)
	}
	return nil
}

func (t sqliteTable[T]) TableName() string { return t.name }

func (t sqliteTable[T]) deploySQL() []string {
	sql := []string{
		buildCreateTableSQL(t.deployPlan(), true),
	}
	seen := map[string]bool{}
	for _, index := range t.indexes {
		indexCols := t.indexColumnNames(index)
		if len(indexCols) == 0 {
			continue
		}
		indexName := fmt.Sprintf("%s__%s_idx", t.name, strings.Join(indexCols, "_"))
		if seen[indexName] {
			continue
		}
		seen[indexName] = true
		sql = append(sql, fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s);",
			quoteIdent(indexName), quoteIdent(t.name), quoteIdentList(indexCols)))
	}
	return sql
}

func (t sqliteTable[T]) deployPlan() tableDeployPlan {
	return tableDeployPlan{
		name:              t.name,
		columns:           t.columns,
		primaryKeyColumns: t.primaryKeyColumns(),
		indexes:           t.indexDeployPlans(),
	}
}

func (t sqliteTable[T]) indexDeployPlans() []indexDeployPlan {
	plans := []indexDeployPlan{}
	seen := map[string]bool{}
	for _, index := range t.indexes {
		indexCols := t.indexColumnNames(index)
		if len(indexCols) == 0 {
			continue
		}
		indexName := fmt.Sprintf("%s__%s_idx", t.name, strings.Join(indexCols, "_"))
		if seen[indexName] {
			continue
		}
		seen[indexName] = true
		plans = append(plans, indexDeployPlan{
			name:    indexName,
			columns: indexCols,
			sql: fmt.Sprintf("CREATE INDEX %s ON %s (%s);",
				quoteIdent(indexName), quoteIdent(t.name), quoteIdentList(indexCols)),
		})
	}
	return plans
}

func deployTablePlan(plan tableDeployPlan, currentSchema deployDatabaseSchema) (int, error) {
	changeCount := 0
	tableSchema := currentSchema.tables[plan.name]
	log.Printf("planned indexes for table %q: %s", plan.name, formatPlannedIndexes(plan.indexes))
	log.Printf("existing indexes for table %q: %s", plan.name, formatExistingIndexes(tableSchema.indexes))
	if !tableSchema.exists {
		log.Printf("Missing table %s, creating...", plan.name)
		if _, err := execD1(buildCreateTableSQL(plan, false), nil); err != nil {
			return changeCount, err
		}
		changeCount++
		log.Printf("Table %s created!", plan.name)
		tableSchema = deployTableSchema{
			exists:  true,
			columns: plannedColumnSet(plan),
			indexes: map[string][]string{},
		}
	}

	// Heal pre-existing tables by adding new nullable columns without rebuilding data.
	for _, column := range missingColumns(plan, tableSchema) {
		log.Printf("Missing column %s, creating...", column.Name)
		if _, err := execD1(buildAddColumnSQL(plan.name, column), nil); err != nil {
			return changeCount, err
		}
		changeCount++
		log.Printf("Column %s created!", column.Name)
	}

	for _, index := range missingIndexes(plan, tableSchema) {
		log.Printf("Missing index %s (%s), creating...", index.name, strings.Join(index.columns, ", "))
		if _, err := execD1(index.sql, nil); err != nil {
			return changeCount, err
		}
		changeCount++
		tableSchema.indexes[index.name] = index.columns
		log.Printf("Index %s created!", index.name)
	}
	for _, index := range changedIndexes(plan, tableSchema) {
		log.Printf("Index %s columns changed from (%s) to (%s), recreating...",
			index.name,
			strings.Join(tableSchema.indexes[index.name], ", "),
			strings.Join(index.columns, ", "))
		if _, err := execD1(fmt.Sprintf("DROP INDEX %s;", quoteIdent(index.name)), nil); err != nil {
			return changeCount, err
		}
		if _, err := execD1(index.sql, nil); err != nil {
			return changeCount, err
		}
		changeCount++
		tableSchema.indexes[index.name] = index.columns
		log.Printf("Index %s recreated!", index.name)
	}
	return changeCount, nil
}

func plannedColumnSet(plan tableDeployPlan) map[string]bool {
	columns := map[string]bool{}
	for _, column := range plan.columns {
		columns[column.Name] = true
	}
	return columns
}

func missingColumns(plan tableDeployPlan, tableSchema deployTableSchema) []columnInfo {
	missing := []columnInfo{}
	for _, column := range plan.columns {
		if tableSchema.columns[column.Name] {
			continue
		}
		missing = append(missing, column)
	}
	return missing
}

func missingIndexes(plan tableDeployPlan, tableSchema deployTableSchema) []indexDeployPlan {
	missing := []indexDeployPlan{}
	for _, index := range plan.indexes {
		if _, exists := tableSchema.indexes[index.name]; exists {
			continue
		}
		missing = append(missing, index)
	}
	return missing
}

func changedIndexes(plan tableDeployPlan, tableSchema deployTableSchema) []indexDeployPlan {
	changed := []indexDeployPlan{}
	for _, index := range plan.indexes {
		existingColumns, exists := tableSchema.indexes[index.name]
		if !exists || slices.Equal(existingColumns, index.columns) {
			continue
		}
		changed = append(changed, index)
	}
	return changed
}

func formatPlannedIndexes(indexes []indexDeployPlan) string {
	if len(indexes) == 0 {
		return "(none)"
	}
	parts := make([]string, 0, len(indexes))
	for _, index := range indexes {
		parts = append(parts, fmt.Sprintf("%s(%s)", index.name, strings.Join(index.columns, ", ")))
	}
	sort.Strings(parts)
	return strings.Join(parts, "; ")
}

func formatExistingIndexes(indexes map[string][]string) string {
	if len(indexes) == 0 {
		return "(none)"
	}
	parts := make([]string, 0, len(indexes))
	for name, columns := range indexes {
		parts = append(parts, fmt.Sprintf("%s(%s)", name, strings.Join(columns, ", ")))
	}
	sort.Strings(parts)
	return strings.Join(parts, "; ")
}

func buildCreateTableSQL(plan tableDeployPlan, ifNotExists bool) string {
	defs := make([]string, 0, len(plan.columns)+1)
	for _, col := range plan.columns {
		defs = append(defs, deployColumnDefinition(col, plan.primaryKeyColumns))
	}
	if primaryKey := primaryKeySQL(plan.primaryKeyColumns); len(plan.primaryKeyColumns) > 1 {
		defs = append(defs, primaryKey)
	}
	createPrefix := "CREATE TABLE"
	if ifNotExists {
		createPrefix += " IF NOT EXISTS"
	}
	return fmt.Sprintf("%s %s (\n\t%s\n);", createPrefix, quoteIdent(plan.name), strings.Join(defs, ",\n\t"))
}

func buildAddColumnSQL(tableName string, column columnInfo) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;",
		quoteIdent(tableName), quoteIdent(column.Name), sqliteType(column.FieldType))
}

func fetchDeployDatabaseSchema(plans []tableDeployPlan) (deployDatabaseSchema, error) {
	schema := deployDatabaseSchema{tables: map[string]deployTableSchema{}}
	for _, plan := range plans {
		log.Printf("fetching columns and indexes for table %q ...", plan.name)
		columns, err := fetchTableColumns(plan.name)
		if err != nil {
			return schema, err
		}
		indexes, err := fetchTableIndexes(plan.name)
		if err != nil {
			return schema, err
		}
		schema.tables[plan.name] = deployTableSchema{
			exists:  len(columns) > 0,
			columns: columns,
			indexes: indexes,
		}
	}
	return schema, nil
}

func fetchTableColumns(tableName string) (map[string]bool, error) {
	resp, err := execD1(fmt.Sprintf("PRAGMA table_info(%s);", quoteIdent(tableName)), nil)
	if err != nil {
		return nil, err
	}
	columns := map[string]bool{}
	for _, row := range d1Rows(resp) {
		name, ok := row["name"].(string)
		if ok && name != "" {
			columns[name] = true
		}
	}
	return columns, nil
}

func fetchTableIndexes(tableName string) (map[string][]string, error) {
	resp, err := execD1(fmt.Sprintf("PRAGMA index_list(%s);", quoteIdent(tableName)), nil)
	if err != nil {
		return nil, err
	}
	indexes := map[string][]string{}
	for _, row := range d1Rows(resp) {
		name, ok := row["name"].(string)
		if !ok || name == "" {
			log.Printf("index metadata row for table %q did not include a string name: %#v", tableName, row)
			continue
		}
		columns, err := fetchIndexColumns(name)
		if err != nil {
			return nil, err
		}
		indexes[name] = columns
	}
	return indexes, nil
}

func fetchIndexColumns(indexName string) ([]string, error) {
	resp, err := execD1(fmt.Sprintf("PRAGMA index_info(%s);", quoteIdent(indexName)), nil)
	if err != nil {
		return nil, err
	}
	columns := []string{}
	for _, row := range d1Rows(resp) {
		name, ok := row["name"].(string)
		if ok && name != "" {
			columns = append(columns, name)
		}
	}
	return columns, nil
}

func d1Rows(resp d1Response) []map[string]any {
	rows := []map[string]any{}
	for _, result := range resp.Result {
		rows = append(rows, result.Results...)
	}
	return rows
}

type deployDatabaseSchema struct {
	tables map[string]deployTableSchema
}

type deployTableSchema struct {
	exists  bool
	columns map[string]bool
	indexes map[string][]string
}

func (t sqliteTable[T]) primaryKeyColumns() []string {
	names := []string{}
	if t.partition != nil {
		names = append(names, t.partition.Name)
	}
	for _, key := range t.keys {
		names = append(names, key.Name)
	}
	return names
}

func (t sqliteTable[T]) whereKeyColumns() []columnInfo {
	columns := []columnInfo{}
	if t.partition != nil {
		columns = append(columns, *t.partition)
	}
	columns = append(columns, t.keys...)
	return columns
}

func (t sqliteTable[T]) keyColumnSet() map[string]bool {
	set := map[string]bool{}
	for _, col := range t.whereKeyColumns() {
		set[col.Name] = true
	}
	return set
}

func (t sqliteTable[T]) indexColumnNames(index Index) []string {
	names := []string{}
	if index.Type == TypeLocalIndex || index.KeepPart {
		if t.partition != nil {
			names = append(names, t.partition.Name)
		}
	}
	for _, key := range index.Keys {
		name := key.GetName()
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

func columnDefinitionSQL(column columnInfo) string {
	return fmt.Sprintf("%s %s", quoteIdent(column.Name), sqliteType(column.FieldType))
}

func primaryKeySQL(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	if len(keys) == 1 {
		return " PRIMARY KEY"
	}
	return fmt.Sprintf("PRIMARY KEY (%s)", quoteIdentList(keys))
}

func deployColumnDefinition(column columnInfo, primaryKeyColumns []string) string {
	definition := columnDefinitionSQL(column)
	if len(primaryKeyColumns) == 1 && primaryKeyColumns[0] == column.Name {
		definition += primaryKeySQL(primaryKeyColumns)
	}
	return definition
}
