package db

import (
	"reflect"
	"strings"
	"testing"
)

type testPayload struct {
	Name string `json:"name"`
}

type testRecord struct {
	TableStruct[testTable, testRecord]
	CompanyID int32       `db:"company_id"`
	ID        string      `db:"id"`
	Scores    []int32     `db:"scores"`
	Payload   testPayload `db:"payload"`
}

type testTable struct {
	TableStruct[testTable, testRecord]
	CompanyID Col[testTable, int32]
	ID        Col[testTable, string]
	Scores    ColSlice[testTable, int32]
	Payload   Col[testTable, testPayload]
}

func (testTable) GetSchema() TableSchema {
	table := Table[testRecord]()
	return TableSchema{
		Name:      "test_records",
		Partition: table.CompanyID,
		Keys:      []Coln{table.ID},
		Indexes: []Index{
			{Type: TypeLocalIndex, Keys: []Coln{table.Scores}},
		},
	}
}

func TestBuildDeploySQLUsesSQLiteTypesAndIndexes(t *testing.T) {
	sql := strings.Join(Table[testRecord]().BuildDeploySQL(), "\n")

	for _, expected := range []string{
		`CREATE TABLE IF NOT EXISTS "test_records"`,
		`"company_id" INTEGER`,
		`"id" TEXT`,
		`"scores" TEXT`,
		`"payload" TEXT`,
		`PRIMARY KEY ("company_id", "id")`,
		`CREATE INDEX IF NOT EXISTS "test_records__company_id_scores_idx"`,
	} {
		if !strings.Contains(sql, expected) {
			t.Fatalf("expected deploy SQL to contain %q, got:\n%s", expected, sql)
		}
	}
}

func TestNormalizeSQLiteValueSerializesComplexValues(t *testing.T) {
	record := testRecord{
		Scores:  []int32{1, 2, 3},
		Payload: testPayload{Name: "demo"},
	}
	value := reflect.ValueOf(record)

	scores, err := normalizeSQLiteValue(value.FieldByName("Scores"))
	if err != nil {
		t.Fatal(err)
	}
	if scores != `[1,2,3]` {
		t.Fatalf("expected []int32 JSON string, got %#v", scores)
	}

	payload, err := normalizeSQLiteValue(value.FieldByName("Payload"))
	if err != nil {
		t.Fatal(err)
	}
	if payload != `{"name":"demo"}` {
		t.Fatalf("expected object JSON string, got %#v", payload)
	}
}

type badRecord struct {
	TableStruct[badTable, badRecord]
	ID int64 `db:"id"`
}

type badTable struct {
	TableStruct[badTable, badRecord]
	ID Col[badTable, string]
}

func (badTable) GetSchema() TableSchema {
	table := Table[badRecord]()
	return TableSchema{Name: "bad_records", Keys: []Coln{table.ID}}
}

func TestTableDeclarationTypeMismatchPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected type mismatch panic")
		}
	}()
	_ = Table[badRecord]()
}
