package redshift

import "testing"

func TestExtractChangedResources(t *testing.T) {
	const (
		currentDatabase = "dev"
		currentSchema   = "public"
	)
	summary, err := ExtractChangedResources(`
CREATE TABLE public.created(id INT);
ALTER TABLE public.changed ADD COLUMN name TEXT;
DROP TABLE public.removed;
INSERT INTO public.rows SELECT 1;
UPDATE public.rows SET id = 2;
DELETE FROM public.rows WHERE id = 2;
SET search_path TO analytics;
CREATE TABLE unqualified(id INT);
SELECT * INTO copied_rows FROM public.rows;
`, currentDatabase, currentSchema)
	if err != nil {
		t.Fatalf("ExtractChangedResources returned error: %v", err)
	}

	if summary.InsertCount != 1 {
		t.Fatalf("InsertCount = %d, want 1", summary.InsertCount)
	}
	if summary.UpdateCount != 1 {
		t.Fatalf("UpdateCount = %d, want 1", summary.UpdateCount)
	}
	if summary.DeleteCount != 1 {
		t.Fatalf("DeleteCount = %d, want 1", summary.DeleteCount)
	}
	if summary.DMLCount != 3 {
		t.Fatalf("DMLCount = %d, want 3", summary.DMLCount)
	}

	assertChangedTable(t, summary, currentDatabase, "public", "created", ChangeKindCreate, false)
	assertChangedTable(t, summary, currentDatabase, "public", "changed", ChangeKindAlter, true)
	assertChangedTable(t, summary, currentDatabase, "public", "removed", ChangeKindDrop, true)
	assertChangedTable(t, summary, currentDatabase, "public", "rows", ChangeKindDML, true)
	assertChangedTable(t, summary, currentDatabase, "analytics", "unqualified", ChangeKindCreate, false)
	assertChangedTable(t, summary, currentDatabase, "analytics", "copied_rows", ChangeKindCreate, false)
}

func assertChangedTable(t *testing.T, summary *ChangeSummary, database, schema, name string, kind ChangeKind, affected bool) {
	t.Helper()
	for _, table := range summary.Tables {
		if table.Database == database && table.Schema == schema && table.Name == name {
			if table.Kind != kind {
				t.Fatalf("%s.%s.%s kind = %q, want %q", database, schema, name, table.Kind, kind)
			}
			if table.Affected != affected {
				t.Fatalf("%s.%s.%s affected = %v, want %v", database, schema, name, table.Affected, affected)
			}
			return
		}
	}
	t.Fatalf("missing changed table %s.%s.%s in %#v", database, schema, name, summary.Tables)
}
