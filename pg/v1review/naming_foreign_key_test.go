package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// foreignKeyFinding is a NamingForeignKey finding as the tests compare it.
type foreignKeyFinding struct {
	Code    v1review.Code
	Table   string
	Name    string
	Pattern string
	Detail  string
	Text    string
}

func runNamingForeignKey(t *testing.T, sql string, cfg v1review.TemplateConfig) []foreignKeyFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := NamingForeignKey(list.Items, cfg)
	if err != nil {
		t.Fatalf("NamingForeignKey: %v", err)
	}
	var findings []foreignKeyFinding
	for _, f := range got {
		findings = append(findings, foreignKeyFinding{f.Code, f.Table, f.Name, f.Pattern, f.Detail, sql[f.Start:f.End]})
	}
	return findings
}

func TestNamingForeignKey(t *testing.T) {
	format := v1review.TemplateConfig{Format: "^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}_{{referenced_column}}$"}
	mismatch := func(table, name, pattern, text string) foreignKeyFinding {
		return foreignKeyFinding{Code: v1review.NameMismatch, Table: table, Name: name, Pattern: pattern, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		want   []foreignKeyFinding
	}{
		{
			name:   "table constraint ok",
			sql:    "CREATE TABLE book (id int, author_id int, CONSTRAINT fk_book_author_id_author_id FOREIGN KEY (author_id) REFERENCES author (id));",
			config: format,
		},
		{
			name:   "table constraint mismatch",
			sql:    "CREATE TABLE book (author_id int, CONSTRAINT fk_book_author FOREIGN KEY (author_id) REFERENCES author (id));",
			config: format,
			want: []foreignKeyFinding{mismatch("book", "fk_book_author", "^fk_book_author_id_author_id$",
				"CREATE TABLE book (author_id int, CONSTRAINT fk_book_author FOREIGN KEY (author_id) REFERENCES author (id))")},
		},
		{
			name:   "column lists join with underscore, tables lose their schema",
			sql:    `ALTER TABLE s.book ADD CONSTRAINT fk FOREIGN KEY (a, "B") REFERENCES "Lib"."Author" (x, "Y");`,
			config: format,
			want: []foreignKeyFinding{mismatch("book", "fk", "^fk_book_a_B_Author_x_Y$",
				`ALTER TABLE s.book ADD CONSTRAINT fk FOREIGN KEY (a, "B") REFERENCES "Lib"."Author" (x, "Y")`)},
		},
		{
			name:   "referenced columns left out expand to nothing",
			sql:    "ALTER TABLE book ADD CONSTRAINT fk_book_author_id_author_ FOREIGN KEY (author_id) REFERENCES author; ALTER TABLE book ADD CONSTRAINT fk_book_author_id_author_id FOREIGN KEY (author_id) REFERENCES author;",
			config: format,
			want: []foreignKeyFinding{mismatch("book", "fk_book_author_id_author_id", "^fk_book_author_id_author_$",
				"ALTER TABLE book ADD CONSTRAINT fk_book_author_id_author_id FOREIGN KEY (author_id) REFERENCES author")},
		},
		{
			name:   "column constraint takes its column as the referencing column",
			sql:    `CREATE TABLE book (author_id int CONSTRAINT fk_book_author_id_author_id REFERENCES author (id), "Editor" int CONSTRAINT fk_editor REFERENCES person (id));`,
			config: format,
			want: []foreignKeyFinding{mismatch("book", "fk_editor", "^fk_book_Editor_person_id$",
				`CREATE TABLE book (author_id int CONSTRAINT fk_book_author_id_author_id REFERENCES author (id), "Editor" int CONSTRAINT fk_editor REFERENCES person (id))`)},
		},
		{
			name:   "unnamed foreign key judged with an empty name",
			sql:    "CREATE TABLE book (author_id int REFERENCES author (id));",
			config: format,
			want: []foreignKeyFinding{mismatch("book", "", "^fk_book_author_id_author_id$",
				"CREATE TABLE book (author_id int REFERENCES author (id))")},
		},
		{
			name:   "create table judges table constraints before column constraints",
			sql:    "CREATE TABLE t (a int CONSTRAINT c1 REFERENCES p (x), CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x), b int CONSTRAINT c3 REFERENCES p (x), CONSTRAINT c4 FOREIGN KEY (b) REFERENCES p (x));",
			config: v1review.TemplateConfig{Format: "^fk$"},
			want: []foreignKeyFinding{
				mismatch("t", "c2", "^fk$", "CREATE TABLE t (a int CONSTRAINT c1 REFERENCES p (x), CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x), b int CONSTRAINT c3 REFERENCES p (x), CONSTRAINT c4 FOREIGN KEY (b) REFERENCES p (x))"),
				mismatch("t", "c4", "^fk$", "CREATE TABLE t (a int CONSTRAINT c1 REFERENCES p (x), CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x), b int CONSTRAINT c3 REFERENCES p (x), CONSTRAINT c4 FOREIGN KEY (b) REFERENCES p (x))"),
				mismatch("t", "c1", "^fk$", "CREATE TABLE t (a int CONSTRAINT c1 REFERENCES p (x), CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x), b int CONSTRAINT c3 REFERENCES p (x), CONSTRAINT c4 FOREIGN KEY (b) REFERENCES p (x))"),
				mismatch("t", "c3", "^fk$", "CREATE TABLE t (a int CONSTRAINT c1 REFERENCES p (x), CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x), b int CONSTRAINT c3 REFERENCES p (x), CONSTRAINT c4 FOREIGN KEY (b) REFERENCES p (x))"),
			},
		},
		{
			name:   "alter judges its commands in order",
			sql:    "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 REFERENCES p (x), ADD CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x) NOT VALID, ADD b int CONSTRAINT c3 REFERENCES p (x);",
			config: v1review.TemplateConfig{Format: "^fk$"},
			want: []foreignKeyFinding{
				mismatch("t", "c1", "^fk$", "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 REFERENCES p (x), ADD CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x) NOT VALID, ADD b int CONSTRAINT c3 REFERENCES p (x)"),
				mismatch("t", "c2", "^fk$", "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 REFERENCES p (x), ADD CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x) NOT VALID, ADD b int CONSTRAINT c3 REFERENCES p (x)"),
				mismatch("t", "c3", "^fk$", "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 REFERENCES p (x), ADD CONSTRAINT c2 FOREIGN KEY (a) REFERENCES p (x) NOT VALID, ADD b int CONSTRAINT c3 REFERENCES p (x)"),
			},
		},
		{
			name:   "alter foreign table judged",
			sql:    "ALTER FOREIGN TABLE ft ADD CONSTRAINT c FOREIGN KEY (a) REFERENCES p (x);",
			config: v1review.TemplateConfig{Format: "^fk$"},
			want:   []foreignKeyFinding{mismatch("ft", "c", "^fk$", "ALTER FOREIGN TABLE ft ADD CONSTRAINT c FOREIGN KEY (a) REFERENCES p (x)")},
		},
		{
			name:   "other constraints not judged",
			sql:    `CREATE TABLE t (a int CONSTRAINT "Bad" PRIMARY KEY, b int CONSTRAINT "Bad2" UNIQUE CONSTRAINT "Bad3" CHECK (b > 0), CONSTRAINT "Bad4" UNIQUE (a, b)); ALTER TABLE t ADD CONSTRAINT "Bad5" CHECK (a > 0); ALTER TABLE t RENAME CONSTRAINT c TO "Bad6"; ALTER TABLE t DROP CONSTRAINT "Bad7"; ALTER TABLE t VALIDATE CONSTRAINT "Bad8";`,
			config: format,
		},
		{
			name:   "create foreign table, create schema element, and routine body not judged",
			sql:    `CREATE FOREIGN TABLE ft (a int CONSTRAINT bad REFERENCES p (x)) SERVER srv; CREATE SCHEMA s CREATE TABLE t (a int CONSTRAINT bad REFERENCES p (x)); CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC ALTER TABLE t ADD CONSTRAINT bad FOREIGN KEY (a) REFERENCES p (x); END;`,
			config: format,
		},
		{
			name:   "unquoted names fold to lower",
			sql:    "ALTER TABLE Book ADD CONSTRAINT FK_Book_Author_Id_Author_Id FOREIGN KEY (Author_Id) REFERENCES Author (Id);",
			config: format,
		},
		{
			name:   "value with a metacharacter changes the pattern",
			sql:    `ALTER TABLE "a.c" ADD CONSTRAINT fk_abc_x_p_y FOREIGN KEY (x) REFERENCES p (y); ALTER TABLE "a|c" ADD CONSTRAINT zz_c_x_p_y FOREIGN KEY (x) REFERENCES p (y); ALTER TABLE "a+c" ADD CONSTRAINT "fk_a+c_x_p_y" FOREIGN KEY (x) REFERENCES p (y);`,
			config: format,
			want: []foreignKeyFinding{mismatch("a+c", "fk_a+c_x_p_y", "^fk_a+c_x_p_y$",
				`ALTER TABLE "a+c" ADD CONSTRAINT "fk_a+c_x_p_y" FOREIGN KEY (x) REFERENCES p (y)`)},
		},
		{
			name:   "value that keeps the pattern from compiling",
			sql:    `ALTER TABLE "a(" ADD CONSTRAINT "` + strings.Repeat("k", 20) + `" FOREIGN KEY (x) REFERENCES p (y); ALTER TABLE t ADD CONSTRAINT nope FOREIGN KEY (x) REFERENCES p (y);`,
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$", MaxLength: 10},
			want: []foreignKeyFinding{
				{
					Code: v1review.PatternInvalid, Table: "a(", Name: strings.Repeat("k", 20),
					Detail: "error parsing regexp: missing closing ): `^fk_a($`",
					Text:   `ALTER TABLE "a(" ADD CONSTRAINT "` + strings.Repeat("k", 20) + `" FOREIGN KEY (x) REFERENCES p (y)`,
				},
				mismatch("t", "nope", "^fk_t$", "ALTER TABLE t ADD CONSTRAINT nope FOREIGN KEY (x) REFERENCES p (y)"),
			},
		},
		{
			name:   "format that never compiles reports every foreign key",
			sql:    "CREATE TABLE t (a int REFERENCES p, CONSTRAINT c FOREIGN KEY (a) REFERENCES p);",
			config: v1review.TemplateConfig{Format: "("},
			want: []foreignKeyFinding{
				{Code: v1review.PatternInvalid, Table: "t", Name: "c", Detail: "error parsing regexp: missing closing ): `(`", Text: "CREATE TABLE t (a int REFERENCES p, CONSTRAINT c FOREIGN KEY (a) REFERENCES p)"},
				{Code: v1review.PatternInvalid, Table: "t", Name: "", Detail: "error parsing regexp: missing closing ): `(`", Text: "CREATE TABLE t (a int REFERENCES p, CONSTRAINT c FOREIGN KEY (a) REFERENCES p)"},
			},
		},
		{
			name:   "token with no value is left as written",
			sql:    "ALTER TABLE t ADD CONSTRAINT fk_t FOREIGN KEY (x) REFERENCES p (y);",
			config: v1review.TemplateConfig{Format: "^fk_{{table}}$"},
			want:   []foreignKeyFinding{mismatch("t", "fk_t", "^fk_{{table}}$", "ALTER TABLE t ADD CONSTRAINT fk_t FOREIGN KEY (x) REFERENCES p (y)")},
		},
		{
			name:   "length counts bytes",
			sql:    `ALTER TABLE t ADD CONSTRAINT "外键约" FOREIGN KEY (x) REFERENCES p (y), ADD CONSTRAINT abcdefgh FOREIGN KEY (x) REFERENCES p (y);`,
			config: v1review.TemplateConfig{MaxLength: 8},
			want: []foreignKeyFinding{{
				Code: v1review.NameTooLong, Table: "t", Name: "外键约",
				Text: `ALTER TABLE t ADD CONSTRAINT "外键约" FOREIGN KEY (x) REFERENCES p (y), ADD CONSTRAINT abcdefgh FOREIGN KEY (x) REFERENCES p (y)`,
			}},
		},
		{
			name:   "name truncated to 63 bytes before judging",
			sql:    `ALTER TABLE t ADD CONSTRAINT "` + strings.Repeat("a", 70) + `" FOREIGN KEY (x) REFERENCES p (y);`,
			config: v1review.TemplateConfig{Format: "^a{63}$", MaxLength: 63},
		},
		{
			name:   "pattern then length",
			sql:    "ALTER TABLE t ADD CONSTRAINT toolong FOREIGN KEY (x) REFERENCES p (y);",
			config: v1review.TemplateConfig{Format: "^fk$", MaxLength: 4},
			want: []foreignKeyFinding{
				mismatch("t", "toolong", "^fk$", "ALTER TABLE t ADD CONSTRAINT toolong FOREIGN KEY (x) REFERENCES p (y)"),
				{Code: v1review.NameTooLong, Table: "t", Name: "toolong", Text: "ALTER TABLE t ADD CONSTRAINT toolong FOREIGN KEY (x) REFERENCES p (y)"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE t (a int REFERENCES p);"},
		{
			name:   "anchor excludes a leading comment and the semicolon",
			sql:    "-- c\nALTER TABLE t\n  ADD CONSTRAINT bad FOREIGN KEY (x) REFERENCES p (y);",
			config: v1review.TemplateConfig{Format: "^fk$"},
			want:   []foreignKeyFinding{mismatch("t", "bad", "^fk$", "ALTER TABLE t\n  ADD CONSTRAINT bad FOREIGN KEY (x) REFERENCES p (y)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runNamingForeignKey(t, tt.sql, tt.config)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
