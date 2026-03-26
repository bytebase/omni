package catalog

import "testing"

func TestDiffSequence(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []SequenceDiffEntry)
	}{
		{
			name:    "standalone sequence added",
			fromSQL: "",
			toSQL:   "CREATE SEQUENCE my_seq;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "my_seq" {
					t.Fatalf("expected my_seq, got %s", e.Name)
				}
				if e.SchemaName != "public" {
					t.Fatalf("expected public, got %s", e.SchemaName)
				}
				if e.From != nil {
					t.Fatalf("expected From to be nil")
				}
				if e.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
			},
		},
		{
			name:    "standalone sequence dropped",
			fromSQL: "CREATE SEQUENCE my_seq;",
			toSQL:   "",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "my_seq" {
					t.Fatalf("expected my_seq, got %s", e.Name)
				}
				if e.From == nil {
					t.Fatalf("expected From to be non-nil")
				}
				if e.To != nil {
					t.Fatalf("expected To to be nil")
				}
			},
		},
		{
			name:    "sequence increment changed",
			fromSQL: "CREATE SEQUENCE my_seq INCREMENT BY 1;",
			toSQL:   "CREATE SEQUENCE my_seq INCREMENT BY 5;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.Increment != 1 {
					t.Fatalf("expected From.Increment=1, got %d", e.From.Increment)
				}
				if e.To.Increment != 5 {
					t.Fatalf("expected To.Increment=5, got %d", e.To.Increment)
				}
			},
		},
		{
			name:    "sequence min/max values changed",
			fromSQL: "CREATE SEQUENCE my_seq MINVALUE 1 MAXVALUE 1000;",
			toSQL:   "CREATE SEQUENCE my_seq MINVALUE 10 MAXVALUE 5000;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.MinValue != 1 {
					t.Fatalf("expected From.MinValue=1, got %d", e.From.MinValue)
				}
				if e.To.MinValue != 10 {
					t.Fatalf("expected To.MinValue=10, got %d", e.To.MinValue)
				}
				if e.From.MaxValue != 1000 {
					t.Fatalf("expected From.MaxValue=1000, got %d", e.From.MaxValue)
				}
				if e.To.MaxValue != 5000 {
					t.Fatalf("expected To.MaxValue=5000, got %d", e.To.MaxValue)
				}
			},
		},
		{
			name:    "sequence cycle flag changed",
			fromSQL: "CREATE SEQUENCE my_seq NO CYCLE;",
			toSQL:   "CREATE SEQUENCE my_seq CYCLE;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.Cycle != false {
					t.Fatalf("expected From.Cycle=false, got %v", e.From.Cycle)
				}
				if e.To.Cycle != true {
					t.Fatalf("expected To.Cycle=true, got %v", e.To.Cycle)
				}
			},
		},
		{
			name:    "sequence start value changed",
			fromSQL: "CREATE SEQUENCE my_seq START WITH 1;",
			toSQL:   "CREATE SEQUENCE my_seq START WITH 100;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.Start != 1 {
					t.Fatalf("expected From.Start=1, got %d", e.From.Start)
				}
				if e.To.Start != 100 {
					t.Fatalf("expected To.Start=100, got %d", e.To.Start)
				}
			},
		},
		{
			name:    "sequence cache value changed",
			fromSQL: "CREATE SEQUENCE my_seq CACHE 1;",
			toSQL:   "CREATE SEQUENCE my_seq CACHE 10;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.CacheValue != 1 {
					t.Fatalf("expected From.CacheValue=1, got %d", e.From.CacheValue)
				}
				if e.To.CacheValue != 10 {
					t.Fatalf("expected To.CacheValue=10, got %d", e.To.CacheValue)
				}
			},
		},
		{
			name:    "sequence type changed int4 to int8",
			fromSQL: "CREATE SEQUENCE my_seq AS integer;",
			toSQL:   "CREATE SEQUENCE my_seq AS bigint;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From == nil || e.To == nil {
					t.Fatalf("expected both From and To to be non-nil")
				}
			},
		},
		{
			name: "SERIAL/IDENTITY-owned sequence skipped",
			fromSQL: `CREATE TABLE t1 (id serial PRIMARY KEY);`,
			toSQL:   "",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				// The serial column creates an owned sequence (OwnerRelOID != 0).
				// It should be skipped, so no sequence diff entries.
				if len(entries) != 0 {
					t.Fatalf("expected 0 entries (owned sequences skipped), got %d", len(entries))
				}
			},
		},
		{
			name:    "sequence unchanged produces no entry",
			fromSQL: "CREATE SEQUENCE my_seq INCREMENT BY 1 START WITH 1 CACHE 1 NO CYCLE;",
			toSQL:   "CREATE SEQUENCE my_seq INCREMENT BY 1 START WITH 1 CACHE 1 NO CYCLE;",
			check: func(t *testing.T, entries []SequenceDiffEntry) {
				if len(entries) != 0 {
					t.Fatalf("expected 0 entries, got %d", len(entries))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var from, to *Catalog
			var err error
			if tt.fromSQL == "" {
				from = New()
			} else {
				from, err = LoadSQL(tt.fromSQL)
				if err != nil {
					t.Fatalf("LoadSQL(from): %v", err)
				}
			}
			if tt.toSQL == "" {
				to = New()
			} else {
				to, err = LoadSQL(tt.toSQL)
				if err != nil {
					t.Fatalf("LoadSQL(to): %v", err)
				}
			}
			entries := diffSequences(from, to)
			tt.check(t, entries)
		})
	}
}

func TestDiffSequenceSerialAbsorption(t *testing.T) {
	// FROM: independent sequence + table using nextval (how metadata DDL looks)
	from, err := LoadSQL(`
		CREATE SEQUENCE users_id_seq;
		CREATE TABLE users (id integer NOT NULL DEFAULT nextval('users_id_seq'));
		ALTER SEQUENCE users_id_seq OWNED BY users.id;
	`)
	if err != nil {
		t.Fatal(err)
	}

	// TO: SERIAL absorbs the sequence (how user SDL looks)
	to, err := LoadSQL(`
		CREATE TABLE users (id serial NOT NULL);
	`)
	if err != nil {
		t.Fatal(err)
	}

	diff := Diff(from, to)

	// The sequence should NOT appear as dropped — it was absorbed into SERIAL
	for _, seq := range diff.Sequences {
		if seq.Action == DiffDrop && seq.Name == "users_id_seq" {
			t.Errorf("users_id_seq should not be reported as dropped (absorbed into SERIAL)")
		}
	}
}

func TestDiffSequenceSerialAbsorptionWithoutOwnedBy(t *testing.T) {
	// FROM: independent sequence WITHOUT OWNED BY (metadata DDL might omit it)
	from, err := LoadSQL(`
		CREATE SEQUENCE users_id_seq;
		CREATE TABLE users (id integer NOT NULL DEFAULT nextval('users_id_seq'));
	`)
	if err != nil {
		t.Fatal(err)
	}

	// TO: SERIAL absorbs the sequence
	to, err := LoadSQL(`
		CREATE TABLE users (id serial NOT NULL);
	`)
	if err != nil {
		t.Fatal(err)
	}

	diff := Diff(from, to)

	// Without OWNED BY, the from sequence has OwnerRelOID=0, but to has it owned.
	// The sequence should NOT appear as dropped.
	for _, seq := range diff.Sequences {
		if seq.Action == DiffDrop && seq.Name == "users_id_seq" {
			t.Errorf("users_id_seq should not be reported as dropped (absorbed into SERIAL)")
		}
	}
}
