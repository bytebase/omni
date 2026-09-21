package review

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

func database(owner, searchPath string, schemas ...*metadata.SchemaMetadata) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Name: "db", Owner: owner, SearchPath: searchPath, Schemas: schemas}
}

func schema(name, owner string, tables ...*metadata.TableMetadata) *metadata.SchemaMetadata {
	return &metadata.SchemaMetadata{Name: name, Owner: owner, Tables: tables}
}

// table builds a table from "name type" column specs.
func table(name, owner string, columns ...string) *metadata.TableMetadata {
	t := &metadata.TableMetadata{Name: name, Owner: owner}
	for _, c := range columns {
		nameType := strings.SplitN(c, " ", 2)
		t.Columns = append(t.Columns, &metadata.ColumnMetadata{Name: nameType[0], Type: nameType[1]})
	}
	return t
}

func TestWalkThrough(t *testing.T) {
	public := func(owner string, tables ...*metadata.TableMetadata) *metadata.DatabaseSchemaMetadata {
		return database(owner, "", schema("public", owner, tables...))
	}
	type finding struct {
		statement int
		text      string
		message   string
		targets   []int
	}
	tests := []struct {
		name    string
		sql     string
		targets []review.Target
		rules   []review.Rule
		want    []finding
	}{
		{
			name:    "no schema, nothing to reject against",
			sql:     "ALTER TABLE nope ADD COLUMN c int;",
			targets: []review.Target{{}},
		},
		{
			name:    "the first rejected statement is reported and the walk stops",
			sql:     "ALTER TABLE t ADD COLUMN c int;\nALTER TABLE t2 ADD COLUMN c int;",
			targets: []review.Target{{Schema: public("alice")}},
			want:    []finding{{0, "ALTER TABLE t ADD COLUMN c int", `relation "t" does not exist`, []int{0}}},
		},
		{
			name:    "state carries from statement to statement",
			sql:     "CREATE TABLE n (id int);\nALTER TABLE n ADD COLUMN c int;\nDROP TABLE n;\nALTER TABLE n ADD COLUMN d int;\nALTER TABLE n ADD COLUMN e int;",
			targets: []review.Target{{Schema: public("alice", table("t", "alice", "id integer"))}},
			want:    []finding{{3, "ALTER TABLE n ADD COLUMN d int", `relation "n" does not exist`, []int{0}}},
		},
		{
			name:    "the synced schema is the starting state",
			sql:     "ALTER TABLE t ADD COLUMN id int;",
			targets: []review.Target{{Schema: public("alice", table("t", "alice", "id integer"))}},
			want:    []finding{{0, "ALTER TABLE t ADD COLUMN id int", `column "id" specified more than once`, []int{0}}},
		},
		{
			name:    "DML is not simulated",
			sql:     "INSERT INTO nope VALUES (1); UPDATE nope SET a = 1 WHERE 1 = 1;",
			targets: []review.Target{{Schema: public("alice")}},
		},
		{
			name: "search path from the synced setting",
			sql:  "ALTER TABLE t ADD COLUMN c int;",
			targets: []review.Target{
				{Schema: database("alice", `"$user", app`, schema("app", "alice", table("t", "alice", "id integer")))},
				{Schema: database("alice", `public`, schema("app", "alice", table("t", "alice", "id integer")))},
			},
			want: []finding{{0, "ALTER TABLE t ADD COLUMN c int", `relation "t" does not exist`, []int{1}}},
		},
		{
			name: "$user in the search path is the session user",
			sql:  "ALTER TABLE t ADD COLUMN c int;",
			targets: []review.Target{
				{Schema: database("app", `"$user", public`, schema("app", "app", table("t", "app", "id integer"))), SessionUser: "app"},
				{Schema: database("app", `"$user", public`, schema("app", "app", table("t", "app", "id integer")))},
			},
			want: []finding{{0, "ALTER TABLE t ADD COLUMN c int", `relation "t" does not exist`, []int{1}}},
		},
		{
			name:    "SET search_path in the change",
			sql:     "SET search_path TO app;\nALTER TABLE t ADD COLUMN c int;",
			targets: []review.Target{{Schema: database("alice", "public", schema("app", "alice", table("t", "alice", "id integer")))}},
		},
		{
			name: "equal targets merge, different targets do not",
			sql:  "DELETE FROM t;\nALTER TABLE t ADD COLUMN c int;",
			targets: []review.Target{
				{Schema: public("alice")},
				{Schema: public("alice", table("t", "alice", "id integer"))},
				{Schema: public("alice")},
			},
			want: []finding{
				{0, "DELETE FROM t", "DELETE FROM t has no WHERE clause", []int{0, 1, 2}},
				{1, "ALTER TABLE t ADD COLUMN c int", `relation "t" does not exist`, []int{0, 2}},
			},
		},
		{
			name:    "walk through switched off",
			sql:     "ALTER TABLE t ADD COLUMN c int;",
			targets: []review.Target{{Schema: public("alice")}},
			rules:   []review.Rule{review.RequireWhere},
		},
		{
			name:    "ownership: altering a table needs its owner",
			sql:     "ALTER TABLE t ADD COLUMN c int;\nALTER TABLE u ADD COLUMN c int;",
			targets: []review.Target{{Schema: public("alice", table("t", "alice", "id integer"), table("u", "bob", "id integer")), SessionUser: "bob"}},
			want:    []finding{{0, "ALTER TABLE t ADD COLUMN c int", "table public.t is owned by alice, but the change runs as bob", []int{0}}},
		},
		{
			name:    "ownership: no session user, no check",
			sql:     "ALTER TABLE t ADD COLUMN c int;",
			targets: []review.Target{{Schema: public("alice", table("t", "alice", "id integer"))}},
		},
		{
			name:    "ownership: creating in a schema and creating a schema",
			sql:     "CREATE TABLE v (id int);\nCREATE TABLE s.v (id int);\nCREATE SCHEMA s2;",
			targets: []review.Target{{Schema: database("alice", "", schema("public", "alice"), schema("s", "bob")), SessionUser: "bob"}},
			want: []finding{
				{0, "CREATE TABLE v (id int)", "schema public is owned by alice, but the change runs as bob", []int{0}},
				{2, "CREATE SCHEMA s2", "database db is owned by alice, but the change runs as bob", []int{0}},
			},
		},
		{
			name:    "ownership: pg_database_owner is the database owner",
			sql:     "CREATE TABLE v (id int);",
			targets: []review.Target{{Schema: database("bob", "", schema("public", pgDatabaseOwner)), SessionUser: "bob"}},
		},
		{
			name:    "ownership: SET ROLE changes who the change runs as",
			sql:     "SET ROLE alice;\nALTER TABLE t ADD COLUMN c int;\nRESET ROLE;\nALTER TABLE t ADD COLUMN d int;",
			targets: []review.Target{{Schema: public("alice", table("t", "alice", "id integer")), SessionUser: "bob"}},
			want:    []finding{{3, "ALTER TABLE t ADD COLUMN d int", "table public.t is owned by alice, but the change runs as bob", []int{0}}},
		},
		{
			name:    "ownership: a finding does not stop the walk",
			sql:     "DROP TABLE t;\nCREATE INDEX i ON u (id);\nALTER TABLE t ADD COLUMN c int;",
			targets: []review.Target{{Schema: public("bob", table("t", "alice", "id integer"), table("u", "carol", "id integer")), SessionUser: "bob"}},
			want: []finding{
				{0, "DROP TABLE t", "table public.t is owned by alice, but the change runs as bob", []int{0}},
				{1, "CREATE INDEX i ON u (id)", "table public.u is owned by carol, but the change runs as bob", []int{0}},
				{2, "ALTER TABLE t ADD COLUMN c int", `relation "t" does not exist`, []int{0}},
			},
		},
		{
			name:    "ownership: a table the change created is not checked",
			sql:     "CREATE TABLE n (id int);\nALTER TABLE n ADD COLUMN c int;",
			targets: []review.Target{{Schema: public("bob"), SessionUser: "bob"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := tt.rules
			if rules == nil {
				rules = []review.Rule{review.WalkThrough, review.RequireWhere}
			}
			result, err := Review(context.Background(), tt.sql, review.Options{Rules: rules}, tt.targets)
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			if len(result.Failures) != 0 {
				t.Errorf("Failures = %v, want none", result.Failures)
			}
			want := []review.Finding{}
			for _, f := range tt.want {
				rule := review.WalkThrough
				if strings.HasSuffix(f.message, "has no WHERE clause") {
					rule = review.RequireWhere
				}
				want = append(want, review.Finding{Rule: rule, Statement: f.statement, Range: span(t, tt.sql, f.text), Message: f.message, Targets: f.targets})
			}
			got := result.Findings
			if got == nil {
				got = []review.Finding{}
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseSearchPath(t *testing.T) {
	tests := []struct {
		setting string
		want    []string
	}{
		{"", nil},
		{`"$user", public`, []string{"$user", "public"}},
		{`App, "My Schema", 'it''s', pg_catalog`, []string{"app", "My Schema", "it's", "pg_catalog"}},
		{`a,,b`, []string{"a", "b"}},
		{`"a,b", c`, []string{"a,b", "c"}},
	}
	for _, tt := range tests {
		if got := parseSearchPath(tt.setting); !cmp.Equal(got, tt.want) {
			t.Errorf("parseSearchPath(%q) = %v, want %v", tt.setting, got, tt.want)
		}
	}
}

func TestGroupTargets(t *testing.T) {
	a := database("alice", "", schema("public", "alice"))
	b := database("bob", "", schema("public", "bob"))
	groups, err := groupTargets([]review.Target{
		{Schema: a}, {Schema: b}, {Schema: a, SessionUser: "x"}, {Schema: a}, {}, {}, {Schema: b, BackupDatabaseExists: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	var got [][]int
	for _, g := range groups {
		got = append(got, g.indices)
	}
	want := [][]int{{0, 3}, {1, 6}, {2}, {4, 5}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("groups mismatch (-want +got):\n%s", diff)
	}
}

func TestReviewCanceledDuringTargets(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	target := review.Target{Schema: database("alice", "", schema("public", "alice"))}
	// Cancel after the parse: the first walk observes the done context.
	cancel()
	if _, err := Review(ctx, "ALTER TABLE t ADD COLUMN c int;", review.Options{}, []review.Target{target}); err == nil {
		t.Fatal("Review with a canceled context returned no error")
	}
}
