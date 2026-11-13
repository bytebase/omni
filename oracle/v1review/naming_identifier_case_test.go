package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingIdentifierCase(t *testing.T) {
	upper, lower := v1review.CaseConfig{Upper: true}, v1review.CaseConfig{}
	mismatch := func(name, text string) finding { return finding{v1review.NameCaseMismatch, name, text} }
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.CaseConfig
		want   []finding
	}{
		{
			name:   "upper: quoted lower-case name",
			sql:    `CREATE TABLE "table"(id int, name varchar(255))`,
			config: upper,
			want:   []finding{mismatch("table", `"table"`)},
		},
		{
			name:   "upper: quoted mixed-case name",
			sql:    `CREATE TABLE "Tt"(id int)`,
			config: upper,
			want:   []finding{mismatch("Tt", `"Tt"`)},
		},
		{name: "upper: unquoted names fold to upper case", sql: "create table techBook(id int, Name varchar(255))", config: upper},
		{name: "upper: quoted upper-case name", sql: `CREATE TABLE "TT" ("ID" int)`, config: upper},
		{name: "upper: no letters", sql: `CREATE TABLE "123_$" (id int)`, config: upper},
		{name: "upper: non-ASCII unquoted name folds", sql: "CREATE TABLE été (id int)", config: upper},
		{
			name:   "upper: non-ASCII quoted name",
			sql:    `CREATE TABLE "été" (id int)`,
			config: upper,
			want:   []finding{mismatch("été", `"été"`)},
		},
		{name: "upper: schema qualifier not judged", sql: `CREATE TABLE "shop".t (id int)`, config: upper},
		{name: "upper: references not judged", sql: `SELECT "id" FROM "Tt" WHERE "id" = 1`, config: upper},
		{
			name:   "upper: declaration in a block, on its own line",
			sql:    "DECLARE\n  \"v_Count\" NUMBER;\nBEGIN\n  SELECT COUNT(1) INTO \"v_Count\" FROM \"Tt\";\nEND;",
			config: upper,
			want:   []finding{mismatch("v_Count", `"v_Count" NUMBER;`)},
		},
		{
			name:   "upper: under explain plan",
			sql:    `EXPLAIN PLAN FOR SELECT a AS "b" FROM t`,
			config: upper,
			want:   []finding{mismatch("b", `a AS "b"`)},
		},
		{
			name:   "lower: unquoted names fold to upper case",
			sql:    "create table tech_book(id int)",
			config: lower,
			want:   []finding{mismatch("TECH_BOOK", "tech_book"), mismatch("ID", "id int")},
		},
		{name: "lower: quoted lower-case names", sql: `CREATE TABLE "tech_book" ("id" int)`, config: lower},
		{
			name:   "lower: quoted mixed-case name",
			sql:    `CREATE TABLE "Tt" ("id" int)`,
			config: lower,
			want:   []finding{mismatch("Tt", `"Tt"`)},
		},
		{name: "lower: no letters", sql: `CREATE TABLE "123_$" ("数据" int)`, config: lower},
		{
			name:   "lower: rename target and alias",
			sql:    "ALTER TABLE \"t\" RENAME TO t2;\nSELECT 1 \"one\", 2 two FROM dual \"d\";",
			config: lower,
			want:   []finding{mismatch("T2", "RENAME TO t2"), mismatch("TWO", "2 two")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := anchored(t, tt.sql, NamingIdentifierCase(parse(t, tt.sql), tt.config))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
