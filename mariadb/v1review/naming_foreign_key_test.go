package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingForeignKey(t *testing.T) {
	template := v1review.TemplateConfig{Format: "^$|^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}_{{referenced_column}}$"}
	mismatch := func(table, name, pattern, text string) anchored {
		return anchored{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern, Text: text}
	}
	tooLong := func(table, name, text string) anchored {
		return anchored{Code: v1review.NameTooLong, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		want   []anchored
	}{
		{
			name:   "create ok, named after FOREIGN KEY",
			sql:    "CREATE TABLE book (id INT, author_id INT, FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id));",
			config: template,
		},
		{
			name:   "create anchors the constraint",
			sql:    "CREATE TABLE book (\n  id INT,\n  author_id INT,\n  FOREIGN KEY fk_book_author_id (author_id)\n    REFERENCES author (id) ON DELETE CASCADE\n);",
			config: template,
			want: []anchored{mismatch("book", "fk_book_author_id", "^$|^fk_book_author_id_author_id$",
				"FOREIGN KEY fk_book_author_id (author_id)\n    REFERENCES author (id) ON DELETE CASCADE")},
		},
		{
			name:   "constraint symbol wins over the name after FOREIGN KEY",
			sql:    "CREATE TABLE book (author_id INT, CONSTRAINT fk1 FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id));",
			config: template,
			want: []anchored{mismatch("book", "fk1", "^$|^fk_book_author_id_author_id$",
				"CONSTRAINT fk1 FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id)")},
		},
		{
			name:   "unnamed foreign key judged with an empty name",
			sql:    "CREATE TABLE book (author_id INT, FOREIGN KEY (author_id) REFERENCES author (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$"},
			want:   []anchored{mismatch("book", "", "^fk_book$", "FOREIGN KEY (author_id) REFERENCES author (id)")},
		},
		{
			name:   "unnamed foreign key passes a pattern that takes an empty name",
			sql:    "CREATE TABLE book (author_id INT, FOREIGN KEY (author_id) REFERENCES author (id));",
			config: template,
		},
		{
			name:   "alter anchors the statement",
			sql:    "ALTER TABLE tech_book\n  ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (id);",
			config: template,
			want: []anchored{mismatch("tech_book", "fk_author_id", "^$|^fk_tech_book_author_id_author_id$",
				"ALTER TABLE tech_book\n  ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (id)")},
		},
		{
			name:   "alter add with if not exists",
			sql:    "ALTER TABLE tech_book ADD FOREIGN KEY IF NOT EXISTS fk1 (author_id) REFERENCES author (id);",
			config: template,
			want: []anchored{mismatch("tech_book", "fk1", "^$|^fk_tech_book_author_id_author_id$",
				"ALTER TABLE tech_book ADD FOREIGN KEY IF NOT EXISTS fk1 (author_id) REFERENCES author (id)")},
		},
		{
			name:   "qualifiers dropped, column lists joined",
			sql:    "ALTER TABLE Shop.book ADD CONSTRAINT fk1 FOREIGN KEY (a, b(3) DESC) REFERENCES Lib.author (x, y);",
			config: template,
			want: []anchored{mismatch("book", "fk1", "^$|^fk_book_a_b_author_x_y$",
				"ALTER TABLE Shop.book ADD CONSTRAINT fk1 FOREIGN KEY (a, b(3) DESC) REFERENCES Lib.author (x, y)")},
		},
		{
			name:   "references without a column list",
			sql:    "ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES author;",
			config: template,
			want:   []anchored{mismatch("book", "fk1", "^$|^fk_book_a_author_$", "ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES author")},
		},
		{
			name:   "several foreign keys, in written order",
			sql:    "ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x), ADD INDEX BadIdx (a), ADD FOREIGN KEY fk2 (b) REFERENCES q (y);",
			config: template,
			want: []anchored{
				mismatch("book", "fk1", "^$|^fk_book_a_p_x$", "ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x), ADD INDEX BadIdx (a), ADD FOREIGN KEY fk2 (b) REFERENCES q (y)"),
				mismatch("book", "fk2", "^$|^fk_book_b_q_y$", "ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x), ADD INDEX BadIdx (a), ADD FOREIGN KEY fk2 (b) REFERENCES q (y)"),
			},
		},
		{
			name:   "other constraints not judged",
			sql:    "CREATE TABLE book (id INT, CONSTRAINT Bad1 PRIMARY KEY (id), CONSTRAINT Bad2 UNIQUE (id), CONSTRAINT Bad3 CHECK (id > 0), INDEX Bad4 (id));",
			config: template,
		},
		{
			name:   "column references clause not judged",
			sql:    "CREATE TABLE book (author_id INT REFERENCES author (id));",
			config: v1review.TemplateConfig{Format: "^fk$"},
		},
		{
			name:   "alter add parenthesized form not judged",
			sql:    "ALTER TABLE book ADD (author_id INT, CONSTRAINT fk1 FOREIGN KEY (author_id) REFERENCES author (id));",
			config: template,
		},
		{name: "drop foreign key not judged", sql: "ALTER TABLE book DROP FOREIGN KEY fk1;", config: template},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x); END;",
			config: template,
		},
		{
			name:   "names judged as parsed, quotes removed",
			sql:    "ALTER TABLE `my book` ADD CONSTRAINT `fk_my book_a``b_p_x` FOREIGN KEY (`a``b`) REFERENCES `p` (`x`);",
			config: template,
		},
		{
			name:   "a metacharacter in a name changes the pattern",
			sql:    "ALTER TABLE `b.k` ADD CONSTRAINT fk_book_a_p_x FOREIGN KEY (a) REFERENCES p (x);",
			config: template,
			want:   []anchored{mismatch("b.k", "fk_book_a_p_x", "^$|^fk_b.k_a_p_x$", "ALTER TABLE `b.k` ADD CONSTRAINT fk_book_a_p_x FOREIGN KEY (a) REFERENCES p (x)")},
		},
		{
			name:   "a metacharacter in a name matches as one",
			sql:    "ALTER TABLE `b.k` ADD CONSTRAINT fk_bXk_a_p_x FOREIGN KEY (a) REFERENCES p (x);",
			config: template,
		},
		{
			name:   "a name that keeps the pattern from compiling",
			sql:    "ALTER TABLE `b(k` ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x);",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$", MaxLength: 1},
			want: []anchored{{
				Code: v1review.PatternInvalid, Name: "fk1", Table: "b(k",
				Detail: "error parsing regexp: missing closing ): `^fk_b(k$`",
				Text:   "ALTER TABLE `b(k` ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x)",
			}},
		},
		{
			name:   "a format that never compiles, once per foreign key",
			sql:    "CREATE TABLE book (a INT, FOREIGN KEY fk1 (a) REFERENCES p (x), FOREIGN KEY fk2 (a) REFERENCES q (y));",
			config: v1review.TemplateConfig{Format: "(fk"},
			want: []anchored{
				{Code: v1review.PatternInvalid, Name: "fk1", Table: "book", Detail: "error parsing regexp: missing closing ): `(fk`", Text: "FOREIGN KEY fk1 (a) REFERENCES p (x)"},
				{Code: v1review.PatternInvalid, Name: "fk2", Table: "book", Detail: "error parsing regexp: missing closing ): `(fk`", Text: "FOREIGN KEY fk2 (a) REFERENCES q (y)"},
			},
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (a INT, CONSTRAINT `外键` FOREIGN KEY (a) REFERENCES p (x), CONSTRAINT abcde FOREIGN KEY (a) REFERENCES p (x));",
			config: v1review.TemplateConfig{MaxLength: 5},
			want:   []anchored{tooLong("book", "外键", "CONSTRAINT `外键` FOREIGN KEY (a) REFERENCES p (x)")},
		},
		{
			name:   "mismatch before length",
			sql:    "ALTER TABLE book ADD CONSTRAINT toolong FOREIGN KEY (a) REFERENCES p (x);",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_column}}$", MaxLength: 4},
			want: []anchored{
				mismatch("book", "toolong", "^fk_a$", "ALTER TABLE book ADD CONSTRAINT toolong FOREIGN KEY (a) REFERENCES p (x)"),
				tooLong("book", "toolong", "ALTER TABLE book ADD CONSTRAINT toolong FOREIGN KEY (a) REFERENCES p (x)"),
			},
		},
		{
			name:   "unanchored format matches anywhere in the name",
			sql:    "ALTER TABLE book ADD CONSTRAINT my_fk_a_x FOREIGN KEY (a) REFERENCES p (x), ADD CONSTRAINT other FOREIGN KEY (a) REFERENCES p (x);",
			config: v1review.TemplateConfig{Format: "fk_{{referencing_column}}"},
			want: []anchored{mismatch("book", "other", "fk_a",
				"ALTER TABLE book ADD CONSTRAINT my_fk_a_x FOREIGN KEY (a) REFERENCES p (x), ADD CONSTRAINT other FOREIGN KEY (a) REFERENCES p (x)")},
		},
		{name: "empty config requires nothing", sql: "ALTER TABLE book ADD CONSTRAINT Bad FOREIGN KEY (a) REFERENCES p (x);"},
		{
			name:   "findings in statement order",
			sql:    "ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x);\nCREATE TABLE student (a INT, CONSTRAINT fk2 FOREIGN KEY (a) REFERENCES p (x));",
			config: template,
			want: []anchored{
				mismatch("book", "fk1", "^$|^fk_book_a_p_x$", "ALTER TABLE book ADD CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p (x)"),
				mismatch("student", "fk2", "^$|^fk_student_a_p_x$", "CONSTRAINT fk2 FOREIGN KEY (a) REFERENCES p (x)"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamingForeignKey(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("NamingForeignKey: %v", err)
			}
			if diff := cmp.Diff(tt.want, anchor(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
