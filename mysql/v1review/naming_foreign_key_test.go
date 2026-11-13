package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingForeignKey(t *testing.T) {
	format := v1review.TemplateConfig{Format: "^$|^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}_{{referenced_column}}$"}
	mismatch := func(table, name, pattern, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		want   []located
	}{
		{name: "create ok", sql: "CREATE TABLE book (author_id INT, FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id));", config: format},
		{name: "unnamed matches an empty alternative", sql: "CREATE TABLE book (author_id INT, FOREIGN KEY (author_id) REFERENCES author (id));", config: format},
		{
			name:   "unnamed judged as the empty name",
			sql:    "CREATE TABLE book (author_id INT, FOREIGN KEY (author_id) REFERENCES author (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$"},
			want:   []located{mismatch("book", "", "^fk_book$", "FOREIGN KEY (author_id) REFERENCES author (id)")},
		},
		{
			name:   "create anchors the constraint",
			sql:    "CREATE TABLE book (\n  author_id INT,\n  CONSTRAINT fk1\n    FOREIGN KEY (author_id) REFERENCES author (id) ON DELETE CASCADE\n);",
			config: format,
			want: []located{mismatch("book", "fk1", "^$|^fk_book_author_id_author_id$",
				"CONSTRAINT fk1\n    FOREIGN KEY (author_id) REFERENCES author (id) ON DELETE CASCADE")},
		},
		{name: "constraint symbol wins over the index name", sql: "CREATE TABLE book (a INT, CONSTRAINT fk_book_a_author_id FOREIGN KEY fk1 (a) REFERENCES author (id));", config: format},
		{
			name:   "constraint symbol judged, not the index name",
			sql:    "CREATE TABLE book (a INT, CONSTRAINT fk1 FOREIGN KEY fk_book_a_author_id (a) REFERENCES author (id));",
			config: format,
			want: []located{mismatch("book", "fk1", "^$|^fk_book_a_author_id$",
				"CONSTRAINT fk1 FOREIGN KEY fk_book_a_author_id (a) REFERENCES author (id)")},
		},
		{
			name:   "column lists joined with an underscore, databases dropped",
			sql:    "CREATE TABLE Shop.book (a INT, b INT, FOREIGN KEY fk1 (a, b) REFERENCES Lib.author (x, y));",
			config: format,
			want:   []located{mismatch("book", "fk1", "^$|^fk_book_a_b_author_x_y$", "FOREIGN KEY fk1 (a, b) REFERENCES Lib.author (x, y)")},
		},
		{name: "column references clause not judged", sql: "CREATE TABLE book (author_id INT REFERENCES author (id));", config: v1review.TemplateConfig{Format: "^fk_$"}},
		{name: "other constraints not judged", sql: "CREATE TABLE book (a INT, PRIMARY KEY (a), UNIQUE KEY Bad (a), INDEX Bad2 (a), CONSTRAINT Bad3 CHECK (a > 0));", config: v1review.TemplateConfig{Format: "^fk_$"}},
		{
			name:   "alter anchors the statement",
			sql:    "ALTER TABLE tech_book\n  ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (id);",
			config: format,
			want: []located{mismatch("tech_book", "fk_author_id", "^$|^fk_tech_book_author_id_author_id$",
				"ALTER TABLE tech_book\n  ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (id)")},
		},
		{
			name:   "alter add foreign key without constraint",
			sql:    "ALTER TABLE db.t ADD FOREIGN KEY fk1 (a) REFERENCES p (id);",
			config: format,
			want:   []located{mismatch("t", "fk1", "^$|^fk_t_a_p_id$", "ALTER TABLE db.t ADD FOREIGN KEY fk1 (a) REFERENCES p (id)")},
		},
		{name: "alter add column list foreign key not judged", sql: "ALTER TABLE t ADD COLUMN (z INT, FOREIGN KEY bad (z) REFERENCES p (id));", config: v1review.TemplateConfig{Format: "^fk_$"}},
		{name: "alter drop foreign key not judged", sql: "ALTER TABLE t DROP FOREIGN KEY bad;", config: v1review.TemplateConfig{Format: "^fk_$"}},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN ALTER TABLE t ADD FOREIGN KEY bad (a) REFERENCES p (id); END;", config: v1review.TemplateConfig{Format: "^fk_$"}},
		{
			name:   "length counts bytes",
			sql:    "ALTER TABLE t ADD CONSTRAINT `外键名` FOREIGN KEY (a) REFERENCES p (id), ADD CONSTRAINT `外键` FOREIGN KEY (a) REFERENCES p (id);",
			config: v1review.TemplateConfig{MaxLength: 8},
			want: []located{{Code: v1review.NameTooLong, Name: "外键名", Table: "t",
				Text: "ALTER TABLE t ADD CONSTRAINT `外键名` FOREIGN KEY (a) REFERENCES p (id), ADD CONSTRAINT `外键` FOREIGN KEY (a) REFERENCES p (id)"}},
		},
		{
			name:   "pattern and length both reported, pattern first",
			sql:    "CREATE TABLE t (a INT, FOREIGN KEY toolong (a) REFERENCES p (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referenced_table}}$", MaxLength: 4},
			want: []located{
				mismatch("t", "toolong", "^fk_p$", "FOREIGN KEY toolong (a) REFERENCES p (id)"),
				{Code: v1review.NameTooLong, Name: "toolong", Table: "t", Text: "FOREIGN KEY toolong (a) REFERENCES p (id)"},
			},
		},
		{
			name:   "a name is not quoted into the pattern",
			sql:    "CREATE TABLE `a.c` (x INT, FOREIGN KEY fk_abc (x) REFERENCES p (id), FOREIGN KEY fk_a (x) REFERENCES p (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$"},
			want:   []located{mismatch("a.c", "fk_a", "^fk_a.c$", "FOREIGN KEY fk_a (x) REFERENCES p (id)")},
		},
		{
			name:   "a name that breaks the pattern",
			sql:    "CREATE TABLE `a(` (x INT, FOREIGN KEY toolong (x) REFERENCES p (id), FOREIGN KEY fk2 (x) REFERENCES `b)` (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}{{referenced_table}}$", MaxLength: 4},
			want: []located{
				{Code: v1review.PatternInvalid, Name: "toolong", Table: "a(", Text: "FOREIGN KEY toolong (x) REFERENCES p (id)"},
				mismatch("a(", "fk2", "^fk_a(b)$", "FOREIGN KEY fk2 (x) REFERENCES `b)` (id)"),
			},
		},
		{
			name:   "a format without tokens",
			sql:    "ALTER TABLE t ADD FOREIGN KEY fk1 (a) REFERENCES p (id), ADD FOREIGN KEY bad (a) REFERENCES p (id);",
			config: v1review.TemplateConfig{Format: "fk"},
			want: []located{mismatch("t", "bad", "fk",
				"ALTER TABLE t ADD FOREIGN KEY fk1 (a) REFERENCES p (id), ADD FOREIGN KEY bad (a) REFERENCES p (id)")},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE t (a INT, FOREIGN KEY Bad (a) REFERENCES p (id));"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamingForeignKey(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("NamingForeignKey: %v", err)
			}
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingForeignKeyPatternInvalid(t *testing.T) {
	sql := "ALTER TABLE t ADD FOREIGN KEY fk1 (a) REFERENCES p (id);"
	got, err := NamingForeignKey(parse(t, sql), v1review.TemplateConfig{Format: "^fk_{{referencing_table}}($"})
	if err != nil {
		t.Fatalf("NamingForeignKey: %v", err)
	}
	if len(got) != 1 || got[0].Code != v1review.PatternInvalid {
		t.Fatalf("findings = %+v, want one PatternInvalid", got)
	}
	if want := "error parsing regexp: missing closing ): `^fk_t($`"; got[0].Detail != want {
		t.Errorf("Detail = %q, want %q", got[0].Detail, want)
	}
}
