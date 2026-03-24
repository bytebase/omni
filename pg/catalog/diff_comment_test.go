package catalog

import (
	"testing"
)

func TestDiffComment(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []CommentDiffEntry)
	}{
		{
			name:    "comment added on table",
			fromSQL: `CREATE TABLE t1 (id int);`,
			toSQL: `CREATE TABLE t1 (id int);
				COMMENT ON TABLE t1 IS 'a table';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 'r' {
					t.Errorf("expected ObjType 'r', got %c", e.ObjType)
				}
				if e.To != "a table" {
					t.Errorf("expected To='a table', got %q", e.To)
				}
			},
		},
		{
			name: "comment dropped from table",
			fromSQL: `CREATE TABLE t1 (id int);
				COMMENT ON TABLE t1 IS 'a table';`,
			toSQL: `CREATE TABLE t1 (id int);`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Errorf("expected DiffDrop, got %d", e.Action)
				}
				if e.ObjType != 'r' {
					t.Errorf("expected ObjType 'r', got %c", e.ObjType)
				}
				if e.From != "a table" {
					t.Errorf("expected From='a table', got %q", e.From)
				}
			},
		},
		{
			name: "comment changed on table",
			fromSQL: `CREATE TABLE t1 (id int);
				COMMENT ON TABLE t1 IS 'old comment';`,
			toSQL: `CREATE TABLE t1 (id int);
				COMMENT ON TABLE t1 IS 'new comment';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Errorf("expected DiffModify, got %d", e.Action)
				}
				if e.From != "old comment" {
					t.Errorf("expected From='old comment', got %q", e.From)
				}
				if e.To != "new comment" {
					t.Errorf("expected To='new comment', got %q", e.To)
				}
			},
		},
		{
			name:    "comment on column added",
			fromSQL: `CREATE TABLE t1 (id int, name text);`,
			toSQL: `CREATE TABLE t1 (id int, name text);
				COMMENT ON COLUMN t1.name IS 'user name';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 'r' {
					t.Errorf("expected ObjType 'r', got %c", e.ObjType)
				}
				if e.SubID == 0 {
					t.Errorf("expected non-zero SubID for column comment")
				}
				if e.To != "user name" {
					t.Errorf("expected To='user name', got %q", e.To)
				}
			},
		},
		{
			name: "comment on column changed",
			fromSQL: `CREATE TABLE t1 (id int, name text);
				COMMENT ON COLUMN t1.name IS 'old';`,
			toSQL: `CREATE TABLE t1 (id int, name text);
				COMMENT ON COLUMN t1.name IS 'new';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Errorf("expected DiffModify, got %d", e.Action)
				}
				if e.From != "old" {
					t.Errorf("expected From='old', got %q", e.From)
				}
				if e.To != "new" {
					t.Errorf("expected To='new', got %q", e.To)
				}
			},
		},
		{
			name: "comment on column dropped",
			fromSQL: `CREATE TABLE t1 (id int, name text);
				COMMENT ON COLUMN t1.name IS 'user name';`,
			toSQL: `CREATE TABLE t1 (id int, name text);`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Errorf("expected DiffDrop, got %d", e.Action)
				}
				if e.From != "user name" {
					t.Errorf("expected From='user name', got %q", e.From)
				}
			},
		},
		{
			name:    "comment on index",
			fromSQL: `CREATE TABLE t1 (id int); CREATE INDEX idx1 ON t1 (id);`,
			toSQL: `CREATE TABLE t1 (id int); CREATE INDEX idx1 ON t1 (id);
				COMMENT ON INDEX idx1 IS 'fast lookup';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 'i' {
					t.Errorf("expected ObjType 'i', got %c", e.ObjType)
				}
				if e.To != "fast lookup" {
					t.Errorf("expected To='fast lookup', got %q", e.To)
				}
			},
		},
		{
			name: "comment on function",
			fromSQL: `CREATE FUNCTION add1(a int) RETURNS int LANGUAGE sql AS 'SELECT a + 1';`,
			toSQL: `CREATE FUNCTION add1(a int) RETURNS int LANGUAGE sql AS 'SELECT a + 1';
				COMMENT ON FUNCTION add1(int) IS 'adds one';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 'f' {
					t.Errorf("expected ObjType 'f', got %c", e.ObjType)
				}
				if e.To != "adds one" {
					t.Errorf("expected To='adds one', got %q", e.To)
				}
			},
		},
		{
			name:    "comment on schema",
			fromSQL: `CREATE SCHEMA myschema;`,
			toSQL: `CREATE SCHEMA myschema;
				COMMENT ON SCHEMA myschema IS 'my schema';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 'n' {
					t.Errorf("expected ObjType 'n', got %c", e.ObjType)
				}
				if e.ObjDescription != "myschema" {
					t.Errorf("expected ObjDescription='myschema', got %q", e.ObjDescription)
				}
				if e.To != "my schema" {
					t.Errorf("expected To='my schema', got %q", e.To)
				}
			},
		},
		{
			name:    "comment on type enum",
			fromSQL: `CREATE TYPE mood AS ENUM ('happy', 'sad');`,
			toSQL: `CREATE TYPE mood AS ENUM ('happy', 'sad');
				COMMENT ON TYPE mood IS 'mood enum';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 't' {
					t.Errorf("expected ObjType 't', got %c", e.ObjType)
				}
				if e.To != "mood enum" {
					t.Errorf("expected To='mood enum', got %q", e.To)
				}
			},
		},
		{
			name:    "comment on sequence",
			fromSQL: `CREATE SEQUENCE myseq;`,
			toSQL: `CREATE SEQUENCE myseq;
				COMMENT ON SEQUENCE myseq IS 'my sequence';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 's' {
					t.Errorf("expected ObjType 's', got %c", e.ObjType)
				}
				if e.To != "my sequence" {
					t.Errorf("expected To='my sequence', got %q", e.To)
				}
			},
		},
		{
			name: "comment on constraint",
			fromSQL: `CREATE TABLE t1 (id int CONSTRAINT t1_pk PRIMARY KEY);`,
			toSQL: `CREATE TABLE t1 (id int CONSTRAINT t1_pk PRIMARY KEY);
				COMMENT ON CONSTRAINT t1_pk ON t1 IS 'primary key';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 'c' {
					t.Errorf("expected ObjType 'c', got %c", e.ObjType)
				}
				if e.To != "primary key" {
					t.Errorf("expected To='primary key', got %q", e.To)
				}
			},
		},
		{
			name: "comment on trigger",
			fromSQL: `CREATE TABLE t1 (id int);
				CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
				CREATE TRIGGER trg1 BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trg_fn();`,
			toSQL: `CREATE TABLE t1 (id int);
				CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
				CREATE TRIGGER trg1 BEFORE INSERT ON t1 FOR EACH ROW EXECUTE FUNCTION trg_fn();
				COMMENT ON TRIGGER trg1 ON t1 IS 'audit trigger';`,
			check: func(t *testing.T, entries []CommentDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.ObjType != 'g' {
					t.Errorf("expected ObjType 'g', got %c", e.ObjType)
				}
				if e.To != "audit trigger" {
					t.Errorf("expected To='audit trigger', got %q", e.To)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, err := LoadSQL(tt.fromSQL)
			if err != nil {
				t.Fatal(err)
			}
			to, err := LoadSQL(tt.toSQL)
			if err != nil {
				t.Fatal(err)
			}
			entries := diffComments(from, to)
			tt.check(t, entries)
		})
	}
}
