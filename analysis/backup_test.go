package analysis

import (
	"context"
	"errors"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func restoreRequest(schema, table string) RestoreRequest {
	db := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{Name: "", Tables: []*metadata.TableMetadata{{Name: "flat"}}},
		{Name: "public", Tables: []*metadata.TableMetadata{{
			Name: "t",
			Indexes: []*metadata.IndexMetadata{
				{Name: "pk", Primary: true},
				{Name: "PK"},
				{Name: "t_pk", Primary: true},
			},
			Partitions: []*metadata.TablePartitionMetadata{{Name: "t_1"}},
		}}},
	}}
	return RestoreRequest{SourceDatabase: "db", SourceSchema: schema, SourceTable: table, Catalog: Catalog{
		GetDatabase: func(_ context.Context, name string) (*Database, error) {
			switch name {
			case "db":
				return &Database{Name: name, Metadata: db, ObjectCaseSensitive: true}, nil
			case "down":
				return nil, errors.New("down")
			default:
				return nil, nil
			}
		},
	}}
}

func TestLookUpRestoreTableReadsTheSourceTable(t *testing.T) {
	table, err := LookUpRestoreTable(context.Background(), restoreRequest("", "t_1"), "public", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	if table.Schema != "public" || table.Table.GetName() != "t" || table.PrimaryKey.GetName() != "pk" {
		t.Fatalf("table = %+v", table)
	}
	// Case-insensitive index names let the later "PK" replace "pk".
	table, err = LookUpRestoreTable(context.Background(), restoreRequest("", "t"), "public", false, nil)
	if err != nil || table.PrimaryKey.GetName() != "t_pk" {
		t.Fatalf("table = %+v, err = %v", table, err)
	}
	if table, err := LookUpRestoreTable(context.Background(), restoreRequest("public", "flat"), "", true, nil); err != nil || table.Schema != "" {
		t.Fatalf("an engine without schemas reads the unnamed schema: %+v, %v", table, err)
	}
}

func TestLookUpRestoreTableErrors(t *testing.T) {
	for _, tc := range []struct {
		request       RestoreRequest
		defaultSchema string
		notSynced     func(string) error
		want          string
	}{
		{restoreRequest("", "t"), "public", nil, ""},
		{RestoreRequest{SourceDatabase: "db"}, "public", nil, "the catalog cannot look databases up"},
		{func() RestoreRequest { r := restoreRequest("", "t"); r.SourceDatabase = "down"; return r }(), "public", nil, "failed to get database metadata for down: down"},
		{func() RestoreRequest { r := restoreRequest("", "t"); r.SourceDatabase = "gone"; return r }(), "public", nil, `database metadata for database "gone" not found (database not synced)`},
		{func() RestoreRequest { r := restoreRequest("", "t"); r.SourceDatabase = "gone"; return r }(), "public", func(db string) error { return errors.New("no " + db) }, "no gone"},
		{restoreRequest("s1", "t"), "public", nil, "schema metadata not found for s1"},
		{restoreRequest("", "nope"), "public", nil, "table metadata not found for public.nope"},
		{restoreRequest("", "nope"), "", nil, "no table metadata for db.nope"},
	} {
		_, err := LookUpRestoreTable(context.Background(), tc.request, tc.defaultSchema, true, tc.notSynced)
		if got := errorText(err); got != tc.want {
			t.Errorf("err = %q, want %q", got, tc.want)
		}
	}
}

func errorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func TestInRange(t *testing.T) {
	from, to := Position{Line: 2, Column: 3}, Position{Line: 4, Column: 1}
	for _, tc := range []struct {
		start, end Position
		want       bool
	}{
		{from, to, true},
		{Position{Line: 2, Column: 2}, to, false},
		{from, Position{Line: 4, Column: 2}, false},
		{Position{Line: 3, Column: 0}, Position{Line: 3, Column: 99}, true},
	} {
		if got := InRange(tc.start, tc.end, from, to); got != tc.want {
			t.Errorf("InRange(%v, %v) = %v", tc.start, tc.end, got)
		}
	}
}

func TestTruncateCommentCountsCharacters(t *testing.T) {
	if got := TruncateComment("é😀ab", 2); got != "é😀..." {
		t.Fatalf("TruncateComment = %q", got)
	}
	if got := TruncateComment("ab", 2); got != "ab" {
		t.Fatalf("TruncateComment = %q", got)
	}
}
