package db

import "errors"

func Merge[T TableBaseInterface[E, T], E TableSchemaInterface[E]](
	records *[]T,
	columnsToExcludeUpdate []Coln,
	onUpdateHandler func(prev, current *T) bool,
	onInsertHandler func(*T),
) error {
	if records == nil || len(*records) == 0 {
		return nil
	}
	if onUpdateHandler == nil || onInsertHandler == nil {
		return errors.New("Merge requires non-nil handlers")
	}
	table := compileTable[E, T](new(E))
	if len(table.keys) == 0 {
		return errors.New("Merge requires at least one key column")
	}

	toInsert := []T{}
	toUpdate := []T{}
	for i := range *records {
		current := &(*records)[i]
		previous, err := table.getByRecordKey(current)
		if err != nil {
			return err
		}
		if previous == nil {
			onInsertHandler(current)
			toInsert = append(toInsert, *current)
			continue
		}
		if onUpdateHandler(previous, current) {
			toUpdate = append(toUpdate, *current)
		}
	}
	if err := Insert(&toInsert); err != nil {
		return err
	}
	if err := UpdateExclude(&toUpdate, columnsToExcludeUpdate...); err != nil {
		return err
	}
	return nil
}
