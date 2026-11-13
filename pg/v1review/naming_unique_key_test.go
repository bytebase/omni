package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingUniqueKey(t *testing.T) {
	format := v1review.TemplateConfig{Format: "^uk_{{table}}_{{column_list}}$"}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		schema v1review.Schema
		want   []keyFinding
	}{
		{
			name:   "create unique index ok",
			sql:    "CREATE UNIQUE INDEX uk_tech_book_id_name ON tech_book (id, name);",
			config: format,
		},
		{
			name:   "create unique index mismatch, plain index not judged",
			sql:    "CREATE UNIQUE INDEX tech_book_id_name ON tech_book (id, name); CREATE INDEX bad ON tech_book (id);",
			config: format,
			want:   []keyFinding{keyMismatch("tech_book", "tech_book_id_name", "^uk_tech_book_id_name$", "CREATE UNIQUE INDEX tech_book_id_name ON tech_book (id, name)")},
		},
		{
			name:   "unnamed unique index judged with an empty name",
			sql:    "CREATE UNIQUE INDEX ON t (a, lower(b));",
			config: format,
			want:   []keyFinding{keyMismatch("t", "", "^uk_t_a$", "CREATE UNIQUE INDEX ON t (a, lower(b))")},
		},
		{
			name:   "create table judges table constraints before column constraints, unnamed ones not at all",
			sql:    `CREATE TABLE s.t (a int CONSTRAINT c1 UNIQUE, CONSTRAINT c2 UNIQUE (a, "B"), b int UNIQUE CONSTRAINT c3 UNIQUE CONSTRAINT c4 NOT NULL, UNIQUE (b), CONSTRAINT c5 UNIQUE (b), CONSTRAINT c6 PRIMARY KEY (a), CONSTRAINT c7 CHECK (a > 0));`,
			config: format,
			want: []keyFinding{
				keyMismatch("t", "c2", "^uk_t_a_B$", `CREATE TABLE s.t (a int CONSTRAINT c1 UNIQUE, CONSTRAINT c2 UNIQUE (a, "B"), b int UNIQUE CONSTRAINT c3 UNIQUE CONSTRAINT c4 NOT NULL, UNIQUE (b), CONSTRAINT c5 UNIQUE (b), CONSTRAINT c6 PRIMARY KEY (a), CONSTRAINT c7 CHECK (a > 0))`),
				keyMismatch("t", "c5", "^uk_t_b$", `CREATE TABLE s.t (a int CONSTRAINT c1 UNIQUE, CONSTRAINT c2 UNIQUE (a, "B"), b int UNIQUE CONSTRAINT c3 UNIQUE CONSTRAINT c4 NOT NULL, UNIQUE (b), CONSTRAINT c5 UNIQUE (b), CONSTRAINT c6 PRIMARY KEY (a), CONSTRAINT c7 CHECK (a > 0))`),
				keyMismatch("t", "c1", "^uk_t_a$", `CREATE TABLE s.t (a int CONSTRAINT c1 UNIQUE, CONSTRAINT c2 UNIQUE (a, "B"), b int UNIQUE CONSTRAINT c3 UNIQUE CONSTRAINT c4 NOT NULL, UNIQUE (b), CONSTRAINT c5 UNIQUE (b), CONSTRAINT c6 PRIMARY KEY (a), CONSTRAINT c7 CHECK (a > 0))`),
				keyMismatch("t", "c3", "^uk_t_b$", `CREATE TABLE s.t (a int CONSTRAINT c1 UNIQUE, CONSTRAINT c2 UNIQUE (a, "B"), b int UNIQUE CONSTRAINT c3 UNIQUE CONSTRAINT c4 NOT NULL, UNIQUE (b), CONSTRAINT c5 UNIQUE (b), CONSTRAINT c6 PRIMARY KEY (a), CONSTRAINT c7 CHECK (a > 0))`),
			},
		},
		{
			name:   "alter judges its commands in order, whatever relation it alters",
			sql:    "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 UNIQUE, ADD CONSTRAINT c2 UNIQUE (a), ADD UNIQUE (a), ADD b int CONSTRAINT c3 UNIQUE, ADD CONSTRAINT c4 PRIMARY KEY (a); ALTER FOREIGN TABLE ft ADD CONSTRAINT c5 UNIQUE (a);",
			config: v1review.TemplateConfig{Format: "^uk$"},
			want: []keyFinding{
				keyMismatch("t", "c1", "^uk$", "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 UNIQUE, ADD CONSTRAINT c2 UNIQUE (a), ADD UNIQUE (a), ADD b int CONSTRAINT c3 UNIQUE, ADD CONSTRAINT c4 PRIMARY KEY (a)"),
				keyMismatch("t", "c2", "^uk$", "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 UNIQUE, ADD CONSTRAINT c2 UNIQUE (a), ADD UNIQUE (a), ADD b int CONSTRAINT c3 UNIQUE, ADD CONSTRAINT c4 PRIMARY KEY (a)"),
				keyMismatch("t", "c3", "^uk$", "ALTER TABLE t ADD COLUMN a int CONSTRAINT c1 UNIQUE, ADD CONSTRAINT c2 UNIQUE (a), ADD UNIQUE (a), ADD b int CONSTRAINT c3 UNIQUE, ADD CONSTRAINT c4 PRIMARY KEY (a)"),
				keyMismatch("ft", "c5", "^uk$", "ALTER FOREIGN TABLE ft ADD CONSTRAINT c5 UNIQUE (a)"),
			},
		},
		{
			name:   "unique using index takes the snapshot index's columns",
			sql:    "ALTER TABLE tech_book ADD CONSTRAINT uk_tech_book_id_name UNIQUE USING INDEX old_index; ALTER TABLE tech_book ADD CONSTRAINT uk_tech_book UNIQUE USING INDEX old_index;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("tech_book", "uk_tech_book", "^uk_tech_book_id_name$", "ALTER TABLE tech_book ADD CONSTRAINT uk_tech_book UNIQUE USING INDEX old_index")},
		},
		{
			name:   "unique using index not found has an empty column list",
			sql:    `ALTER TABLE tech_book ADD CONSTRAINT c1 UNIQUE USING INDEX nope; ALTER TABLE nope ADD CONSTRAINT c2 UNIQUE USING INDEX old_index; ALTER TABLE tech_book ADD CONSTRAINT c3 UNIQUE USING INDEX "OLD_INDEX"; ALTER TABLE parted_p0 ADD CONSTRAINT c4 UNIQUE USING INDEX parted_idx;`,
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "c1", "^uk_tech_book_$", "ALTER TABLE tech_book ADD CONSTRAINT c1 UNIQUE USING INDEX nope"),
				keyMismatch("nope", "c2", "^uk_nope_$", "ALTER TABLE nope ADD CONSTRAINT c2 UNIQUE USING INDEX old_index"),
				keyMismatch("tech_book", "c3", "^uk_tech_book_$", `ALTER TABLE tech_book ADD CONSTRAINT c3 UNIQUE USING INDEX "OLD_INDEX"`),
				keyMismatch("parted_p0", "c4", "^uk_parted_p0_$", "ALTER TABLE parted_p0 ADD CONSTRAINT c4 UNIQUE USING INDEX parted_idx"),
			},
		},
		{
			name:   "unique using index without a snapshot has an empty column list",
			sql:    "ALTER TABLE tech_book ADD CONSTRAINT uk_tech_book_ UNIQUE USING INDEX old_index; ALTER TABLE tech_book ADD CONSTRAINT c UNIQUE USING INDEX old_index;",
			config: format,
			want:   []keyFinding{keyMismatch("tech_book", "c", "^uk_tech_book_$", "ALTER TABLE tech_book ADD CONSTRAINT c UNIQUE USING INDEX old_index")},
		},
		{
			name:   "unique using index looks in public whatever schema the statement names",
			sql:    "ALTER TABLE other.tech_book ADD CONSTRAINT c UNIQUE USING INDEX old_index; ALTER TABLE other.tech_book ADD CONSTRAINT c2 UNIQUE USING INDEX other_idx;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "c", "^uk_tech_book_id_name$", "ALTER TABLE other.tech_book ADD CONSTRAINT c UNIQUE USING INDEX old_index"),
				keyMismatch("tech_book", "c2", "^uk_tech_book_$", "ALTER TABLE other.tech_book ADD CONSTRAINT c2 UNIQUE USING INDEX other_idx"),
			},
		},
		{
			name:   "unique using index in a create table",
			sql:    "CREATE TABLE tech_book (a int, CONSTRAINT c UNIQUE USING INDEX old_uk);",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("tech_book", "c", "^uk_tech_book_id_name$", "CREATE TABLE tech_book (a int, CONSTRAINT c UNIQUE USING INDEX old_uk)")},
		},
		{
			name:   "unique using index finds the table by object case sensitivity and prints the statement's spelling",
			sql:    `ALTER TABLE book ADD CONSTRAINT c UNIQUE USING INDEX book_idx;`,
			config: format,
			schema: keySchema(false),
			want:   []keyFinding{keyMismatch("book", "c", "^uk_book_c$", "ALTER TABLE book ADD CONSTRAINT c UNIQUE USING INDEX book_idx")},
		},
		{
			name:   "renamed index judged when unique and not primary, with the snapshot's table",
			sql:    "ALTER INDEX old_uk RENAME TO uk_tech_book_id_name; ALTER INDEX old_uk RENAME TO bad; ALTER INDEX old_index RENAME TO bad2; ALTER INDEX old_pk RENAME TO bad3; ALTER INDEX nope RENAME TO bad4; ALTER INDEX mv_uk RENAME TO bad5; ALTER INDEX book_uk RENAME TO bad6;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^uk_tech_book_id_name$", "ALTER INDEX old_uk RENAME TO bad"),
				keyMismatch("mv", "bad5", "^uk_mv_m$", "ALTER INDEX mv_uk RENAME TO bad5"),
				keyMismatch("Book", "bad6", "^uk_Book_b$", "ALTER INDEX book_uk RENAME TO bad6"),
			},
		},
		{
			name:   "renamed index looks in public whatever schema the statement names",
			sql:    "ALTER INDEX other.old_uk RENAME TO bad; ALTER INDEX other.other_uk RENAME TO bad2;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("tech_book", "bad", "^uk_tech_book_id_name$", "ALTER INDEX other.old_uk RENAME TO bad")},
		},
		{
			name:   "renamed index, two tables hold the name, the first in snapshot order is taken",
			sql:    "ALTER INDEX twin_uk RENAME TO bad;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("twin_a", "bad", "^uk_twin_a_a$", "ALTER INDEX twin_uk RENAME TO bad")},
		},
		{
			name:   "renamed constraint judged when the table's index is unique and not primary",
			sql:    "ALTER TABLE tech_book RENAME CONSTRAINT old_uk TO uk_tech_book_id_name; ALTER TABLE tech_book RENAME CONSTRAINT old_uk TO bad; ALTER TABLE tech_book RENAME CONSTRAINT old_index TO bad2; ALTER TABLE tech_book RENAME CONSTRAINT old_pk TO bad3; ALTER TABLE tech_book RENAME CONSTRAINT nope TO bad4; ALTER TABLE twin_b RENAME CONSTRAINT twin_uk TO bad5; ALTER TABLE twin_b RENAME CONSTRAINT old_uk TO bad6; ALTER TABLE parted_p0 RENAME CONSTRAINT parted_uk TO bad7;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^uk_tech_book_id_name$", "ALTER TABLE tech_book RENAME CONSTRAINT old_uk TO bad"),
				keyMismatch("twin_b", "bad5", "^uk_twin_b_b$", "ALTER TABLE twin_b RENAME CONSTRAINT twin_uk TO bad5"),
			},
		},
		{
			name:   "renamed constraint looks in the statement's schema",
			sql:    "ALTER TABLE other.tech_book RENAME CONSTRAINT old_index TO bad; ALTER TABLE other.tech_book RENAME CONSTRAINT old_uk TO bad2; ALTER TABLE public.tech_book RENAME CONSTRAINT old_uk TO bad3; ALTER TABLE nope.tech_book RENAME CONSTRAINT old_uk TO bad4;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^uk_tech_book_o$", "ALTER TABLE other.tech_book RENAME CONSTRAINT old_index TO bad"),
				keyMismatch("tech_book", "bad3", "^uk_tech_book_id_name$", "ALTER TABLE public.tech_book RENAME CONSTRAINT old_uk TO bad3"),
			},
		},
		{
			name:   "renamed constraint finds the table by object case sensitivity and prints the statement's spelling",
			sql:    `ALTER TABLE book RENAME CONSTRAINT book_uk TO bad; ALTER TABLE "PUBLIC"."BOOK" RENAME CONSTRAINT book_uk TO bad2; ALTER TABLE tech_book RENAME CONSTRAINT "OLD_UK" TO bad3; ALTER TABLE tech_book RENAME CONSTRAINT "Old_Uk" TO bad4;`,
			config: format,
			schema: keySchema(false),
			want: []keyFinding{
				keyMismatch("book", "bad", "^uk_book_b$", "ALTER TABLE book RENAME CONSTRAINT book_uk TO bad"),
				keyMismatch("BOOK", "bad2", "^uk_BOOK_b$", `ALTER TABLE "PUBLIC"."BOOK" RENAME CONSTRAINT book_uk TO bad2`),
				keyMismatch("tech_book", "bad4", "^uk_tech_book_y$", `ALTER TABLE tech_book RENAME CONSTRAINT "Old_Uk" TO bad4`),
			},
		},
		{
			name:   "table name case sensitive when objects are",
			sql:    "ALTER TABLE book RENAME CONSTRAINT book_uk TO bad; ALTER TABLE book ADD CONSTRAINT uk_book_ UNIQUE USING INDEX book_idx;",
			config: format,
			schema: keySchema(true),
		},
		{
			name:   "renames without a snapshot not judged",
			sql:    "ALTER INDEX old_uk RENAME TO bad; ALTER TABLE tech_book RENAME CONSTRAINT old_uk TO bad2;",
			config: format,
		},
		{
			name:   "other renames not judged",
			sql:    "ALTER TABLE old_uk RENAME TO bad; ALTER TABLE tech_book RENAME COLUMN old_uk TO bad; ALTER DOMAIN old_uk RENAME CONSTRAINT old_uk TO bad;",
			config: format,
			schema: keySchema(true),
		},
		{
			name:   "create foreign table, create schema element, and routine body not judged",
			sql:    "CREATE FOREIGN TABLE ft (a int CONSTRAINT bad UNIQUE, CONSTRAINT bad2 UNIQUE (a)) SERVER srv; CREATE SCHEMA s CREATE TABLE t (a int CONSTRAINT bad UNIQUE) CREATE UNIQUE INDEX bad ON t (a); CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC ALTER TABLE t ADD CONSTRAINT bad UNIQUE (a); END;",
			config: format,
		},
		{
			name:   "value that keeps the pattern from compiling",
			sql:    `ALTER TABLE "a(" ADD CONSTRAINT "` + strings.Repeat("k", 20) + `" UNIQUE (x); ALTER TABLE t ADD CONSTRAINT nope UNIQUE (x);`,
			config: v1review.TemplateConfig{Format: "^uk_{{table}}$", MaxLength: 10},
			want: []keyFinding{
				{
					Code: v1review.PatternInvalid, Table: "a(", Name: strings.Repeat("k", 20),
					Detail: "error parsing regexp: missing closing ): `^uk_a($`",
					Text:   `ALTER TABLE "a(" ADD CONSTRAINT "` + strings.Repeat("k", 20) + `" UNIQUE (x)`,
				},
				keyMismatch("t", "nope", "^uk_t$", "ALTER TABLE t ADD CONSTRAINT nope UNIQUE (x)"),
			},
		},
		{
			name:   "length counts bytes of the truncated name",
			sql:    `ALTER TABLE t ADD CONSTRAINT "唯一键" UNIQUE (x), ADD CONSTRAINT bcdefghi UNIQUE (x), ADD CONSTRAINT "` + strings.Repeat("a", 70) + `" UNIQUE (x);`,
			config: v1review.TemplateConfig{Format: "^(a{63}|[^a]*)$", MaxLength: 8},
			want: []keyFinding{
				{Code: v1review.NameTooLong, Table: "t", Name: "唯一键", Text: `ALTER TABLE t ADD CONSTRAINT "唯一键" UNIQUE (x), ADD CONSTRAINT bcdefghi UNIQUE (x), ADD CONSTRAINT "` + strings.Repeat("a", 70) + `" UNIQUE (x)`},
				{Code: v1review.NameTooLong, Table: "t", Name: strings.Repeat("a", 63), Text: `ALTER TABLE t ADD CONSTRAINT "唯一键" UNIQUE (x), ADD CONSTRAINT bcdefghi UNIQUE (x), ADD CONSTRAINT "` + strings.Repeat("a", 70) + `" UNIQUE (x)`},
			},
		},
		{
			name:   "pattern then length",
			sql:    "CREATE UNIQUE INDEX toolong ON t (x);",
			config: v1review.TemplateConfig{Format: "^uk$", MaxLength: 4},
			want: []keyFinding{
				keyMismatch("t", "toolong", "^uk$", "CREATE UNIQUE INDEX toolong ON t (x)"),
				{Code: v1review.NameTooLong, Table: "t", Name: "toolong", Text: "CREATE UNIQUE INDEX toolong ON t (x)"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE UNIQUE INDEX anything ON t (x);"},
		{
			name:   "anchor excludes a leading comment and the semicolon",
			sql:    "-- c\nALTER TABLE t\n  ADD CONSTRAINT bad UNIQUE (x);",
			config: format,
			want:   []keyFinding{keyMismatch("t", "bad", "^uk_t_x$", "ALTER TABLE t\n  ADD CONSTRAINT bad UNIQUE (x)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runKeyRule(t, NamingUniqueKey, tt.sql, tt.config, tt.schema)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
