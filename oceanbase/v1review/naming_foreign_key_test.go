package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

func TestNamingForeignKey(t *testing.T) {
	template := v1review.TemplateConfig{Format: "^$|^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}_{{referenced_column}}$"}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		want   []judged
	}{
		{
			name:   "create ok",
			sql:    "CREATE TABLE book (id INT, author_id INT, FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id));",
			config: template,
		},
		{
			name:   "create anchors the constraint",
			sql:    "CREATE TABLE book (\n  id INT,\n  author_id INT,\n  FOREIGN KEY fk_book_author_id (author_id) REFERENCES author (id)\n);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "fk_book_author_id", Table: "book",
				Pattern: "^$|^fk_book_author_id_author_id$",
				Text:    "FOREIGN KEY fk_book_author_id (author_id) REFERENCES author (id)",
			}},
		},
		{
			name:   "constraint symbol wins over the name after FOREIGN KEY",
			sql:    "CREATE TABLE book (id INT, CONSTRAINT fk_bad FOREIGN KEY fk_book_id_author_id (id) REFERENCES author (id));",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "fk_bad", Table: "book",
				Pattern: "^$|^fk_book_id_author_id$",
				Text:    "CONSTRAINT fk_bad FOREIGN KEY fk_book_id_author_id (id) REFERENCES author (id)",
			}},
		},
		{
			name:   "unnamed foreign key judged as the empty name",
			sql:    "CREATE TABLE book (id INT, FOREIGN KEY (id) REFERENCES author (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$"},
			want: []judged{{
				Code: v1review.NameMismatch, Table: "book", Pattern: "^fk_book$",
				Text: "FOREIGN KEY (id) REFERENCES author (id)",
			}},
		},
		{
			name:   "columns joined with underscores, tables without their databases",
			sql:    "CREATE TABLE shop.book (a INT, b INT, CONSTRAINT x FOREIGN KEY (a, b) REFERENCES lib.author (c, d));",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "book",
				Pattern: "^$|^fk_book_a_b_author_c_d$",
				Text:    "CONSTRAINT x FOREIGN KEY (a, b) REFERENCES lib.author (c, d)",
			}},
		},
		{
			name:   "referenced columns left out expand to nothing",
			sql:    "CREATE TABLE book (a INT, CONSTRAINT x FOREIGN KEY (a) REFERENCES author);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "book",
				Pattern: "^$|^fk_book_a_author_$",
				Text:    "CONSTRAINT x FOREIGN KEY (a) REFERENCES author",
			}},
		},
		{
			name:   "column reference not judged",
			sql:    "CREATE TABLE book (a INT REFERENCES author (id));",
			config: v1review.TemplateConfig{Format: "^never$"},
		},
		{
			name:   "other constraints not judged",
			sql:    "CREATE TABLE book (a INT, PRIMARY KEY (a), UNIQUE KEY Bad (a), INDEX Bad2 (a), CONSTRAINT Bad3 CHECK (a > 0));",
			config: v1review.TemplateConfig{Format: "^never$"},
		},
		{
			name:   "alter anchors the statement",
			sql:    "ALTER TABLE book\n  ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (id);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "fk_author_id", Table: "book",
				Pattern: "^$|^fk_book_author_id_author_id$",
				Text:    "ALTER TABLE book\n  ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (id)",
			}},
		},
		{
			name:   "alter adds several",
			sql:    "ALTER TABLE book ADD FOREIGN KEY x (a) REFERENCES p (id), ADD INDEX Bad (a), ADD CONSTRAINT fk_book_b_p_id FOREIGN KEY (b) REFERENCES p (id), ADD FOREIGN KEY y (c) REFERENCES p (id);",
			config: template,
			want: []judged{
				{
					Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^fk_book_a_p_id$",
					Text: "ALTER TABLE book ADD FOREIGN KEY x (a) REFERENCES p (id), ADD INDEX Bad (a), ADD CONSTRAINT fk_book_b_p_id FOREIGN KEY (b) REFERENCES p (id), ADD FOREIGN KEY y (c) REFERENCES p (id)",
				},
				{
					Code: v1review.NameMismatch, Name: "y", Table: "book", Pattern: "^$|^fk_book_c_p_id$",
					Text: "ALTER TABLE book ADD FOREIGN KEY x (a) REFERENCES p (id), ADD INDEX Bad (a), ADD CONSTRAINT fk_book_b_p_id FOREIGN KEY (b) REFERENCES p (id), ADD FOREIGN KEY y (c) REFERENCES p (id)",
				},
			},
		},
		{name: "drop foreign key not judged", sql: "ALTER TABLE book DROP FOREIGN KEY Bad;", config: v1review.TemplateConfig{Format: "^never$"}},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN ALTER TABLE book ADD CONSTRAINT Bad FOREIGN KEY (a) REFERENCES p (id); END;",
			config: v1review.TemplateConfig{Format: "^never$"},
		},
		{
			name:   "value not quoted in the pattern",
			sql:    "CREATE TABLE `a.b` (id INT, CONSTRAINT fk_axb FOREIGN KEY (id) REFERENCES p (id), CONSTRAINT fk_ab FOREIGN KEY (id) REFERENCES p (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$"},
			want: []judged{{
				Code: v1review.NameMismatch, Name: "fk_ab", Table: "a.b", Pattern: "^fk_a.b$",
				Text: "CONSTRAINT fk_ab FOREIGN KEY (id) REFERENCES p (id)",
			}},
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (id INT, CONSTRAINT `外键名` FOREIGN KEY (id) REFERENCES p (id), CONSTRAINT abcdefgh FOREIGN KEY (id) REFERENCES p (id));",
			config: v1review.TemplateConfig{MaxLength: 8},
			want: []judged{{
				Code: v1review.NameTooLong, Name: "外键名", Table: "book",
				Text: "CONSTRAINT `外键名` FOREIGN KEY (id) REFERENCES p (id)",
			}},
		},
		{
			name:   "pattern and length both reported",
			sql:    "ALTER TABLE book ADD CONSTRAINT toolong FOREIGN KEY (a) REFERENCES p (id);",
			config: v1review.TemplateConfig{Format: "^fk_", MaxLength: 4},
			want: []judged{
				{Code: v1review.NameMismatch, Name: "toolong", Table: "book", Pattern: "^fk_", Text: "ALTER TABLE book ADD CONSTRAINT toolong FOREIGN KEY (a) REFERENCES p (id)"},
				{Code: v1review.NameTooLong, Name: "toolong", Table: "book", Text: "ALTER TABLE book ADD CONSTRAINT toolong FOREIGN KEY (a) REFERENCES p (id)"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (id INT, CONSTRAINT Bad FOREIGN KEY (id) REFERENCES p (id));"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judge(t, tt.sql, func(stmts []ast.Node) ([]v1review.Finding, error) { return NamingForeignKey(stmts, tt.config) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// A foreign key whose expanded pattern does not compile is reported once,
// with the compile error, and is not judged for length; the next one is
// judged on its own.
func TestNamingForeignKeyPatternInvalid(t *testing.T) {
	sql := "CREATE TABLE `a(b` (id INT, CONSTRAINT toolong FOREIGN KEY (id) REFERENCES p (id));\n" +
		"ALTER TABLE ok ADD CONSTRAINT toolong FOREIGN KEY (id) REFERENCES p (id);"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := NamingForeignKey(list.Items, v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$", MaxLength: 4})
	if err != nil {
		t.Fatalf("NamingForeignKey: %v", err)
	}
	constraint := "CONSTRAINT toolong FOREIGN KEY (id) REFERENCES p (id)"
	alter := "ALTER TABLE ok ADD CONSTRAINT toolong FOREIGN KEY (id) REFERENCES p (id)"
	want := []v1review.Finding{
		{
			Code: v1review.PatternInvalid, Name: "toolong", Table: "a(b",
			Detail: "error parsing regexp: missing closing ): `^fk_a(b$`",
			Start:  28, End: 28 + len(constraint),
		},
		{Code: v1review.NameMismatch, Name: "toolong", Table: "ok", Pattern: "^fk_ok$", Start: 84, End: 84 + len(alter)},
		{Code: v1review.NameTooLong, Name: "toolong", Table: "ok", Start: 84, End: 84 + len(alter)},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
