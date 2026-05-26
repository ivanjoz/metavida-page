package db

import (
	"fmt"
	"log"
	"strings"

	"metavida/backend/config"
)

func DeployTables(cfg config.Config, tables ...TableDeployInterface) error {
	Configure(cfg)
	for _, table := range tables {
		name := table.TableName()
		log.Printf("deploying table %q ...", name)
		for _, sql := range table.BuildDeploySQL() {
			if _, err := execD1(sql, nil); err != nil {
				return fmt.Errorf("table %q: %w", name, err)
			}
		}
		log.Printf("table %q: OK", name)
	}
	return nil
}

func (t sqliteTable[T]) TableName() string { return t.name }

func (t sqliteTable[T]) deploySQL() []string {
	defs := make([]string, 0, len(t.columns)+1)
	for _, col := range t.columns {
		def := fmt.Sprintf("%s %s", quoteIdent(col.Name), sqliteType(col.FieldType))
		if len(t.primaryKeyColumns()) == 1 && t.primaryKeyColumns()[0] == col.Name {
			def += " PRIMARY KEY"
		}
		defs = append(defs, def)
	}
	if keys := t.primaryKeyColumns(); len(keys) > 1 {
		defs = append(defs, fmt.Sprintf("PRIMARY KEY (%s)", quoteIdentList(keys)))
	}

	sql := []string{
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n);", quoteIdent(t.name), strings.Join(defs, ",\n\t")),
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
