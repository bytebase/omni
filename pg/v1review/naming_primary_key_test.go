package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingPrimaryKey(t *testing.T) {
	format := v1review.TemplateConfig{Format: "^pk_{{table}}_{{column_list}}$"}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		schema v1review.Schema
		want   []keyFinding
	}{
		{
			name:   "table constraint ok",
			sql:    "CREATE TABLE book (id int, name text, CONSTRAINT pk_book_id_name PRIMARY KEY (id, name));",
			config: format,
		},
		{
			name:   "table constraint mismatch, table loses its schema",
			sql:    `CREATE TABLE s.book (id int, "Name" text, CONSTRAINT book_pk PRIMARY KEY (id, "Name"));`,
			config: format,
			want:   []keyFinding{keyMismatch("book", "book_pk", "^pk_book_id_Name$", `CREATE TABLE s.book (id int, "Name" text, CONSTRAINT book_pk PRIMARY KEY (id, "Name"))`)},
		},
		{
			name:   "unnamed and column primary keys not judged, nor other constraints",
			sql:    "CREATE TABLE t (a int CONSTRAINT bad PRIMARY KEY, PRIMARY KEY (a), CONSTRAINT bad2 UNIQUE (a)); ALTER TABLE t ADD COLUMN b int CONSTRAINT bad3 PRIMARY KEY, ADD PRIMARY KEY (b), ADD CONSTRAINT bad4 UNIQUE (b);",
			config: format,
		},
		{
			name:   "alter judges its add constraint commands in order, whatever relation it alters",
			sql:    "ALTER TABLE t ADD CONSTRAINT c1 PRIMARY KEY (a), ADD CONSTRAINT c2 PRIMARY KEY (a, b); ALTER FOREIGN TABLE ft ADD CONSTRAINT c3 PRIMARY KEY (a);",
			config: format,
			want: []keyFinding{
				keyMismatch("t", "c1", "^pk_t_a$", "ALTER TABLE t ADD CONSTRAINT c1 PRIMARY KEY (a), ADD CONSTRAINT c2 PRIMARY KEY (a, b)"),
				keyMismatch("t", "c2", "^pk_t_a_b$", "ALTER TABLE t ADD CONSTRAINT c1 PRIMARY KEY (a), ADD CONSTRAINT c2 PRIMARY KEY (a, b)"),
				keyMismatch("ft", "c3", "^pk_ft_a$", "ALTER FOREIGN TABLE ft ADD CONSTRAINT c3 PRIMARY KEY (a)"),
			},
		},
		{
			name:   "primary key using index takes the snapshot index's columns",
			sql:    "ALTER TABLE tech_book ADD CONSTRAINT pk_tech_book_id_name PRIMARY KEY USING INDEX old_index; ALTER TABLE tech_book ADD CONSTRAINT pk_tech_book PRIMARY KEY USING INDEX old_index;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("tech_book", "pk_tech_book", "^pk_tech_book_id_name$", "ALTER TABLE tech_book ADD CONSTRAINT pk_tech_book PRIMARY KEY USING INDEX old_index")},
		},
		{
			name:   "primary key using index looks in the statement's schema",
			sql:    "ALTER TABLE other.tech_book ADD CONSTRAINT c1 PRIMARY KEY USING INDEX old_index; ALTER TABLE public.tech_book ADD CONSTRAINT c2 PRIMARY KEY USING INDEX old_index; ALTER TABLE nope.tech_book ADD CONSTRAINT c3 PRIMARY KEY USING INDEX old_index;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "c1", "^pk_tech_book_o$", "ALTER TABLE other.tech_book ADD CONSTRAINT c1 PRIMARY KEY USING INDEX old_index"),
				keyMismatch("tech_book", "c2", "^pk_tech_book_id_name$", "ALTER TABLE public.tech_book ADD CONSTRAINT c2 PRIMARY KEY USING INDEX old_index"),
				keyMismatch("tech_book", "c3", "^pk_tech_book_$", "ALTER TABLE nope.tech_book ADD CONSTRAINT c3 PRIMARY KEY USING INDEX old_index"),
			},
		},
		{
			name:   "primary key using index not found has an empty column list",
			sql:    `ALTER TABLE tech_book ADD CONSTRAINT c1 PRIMARY KEY USING INDEX nope; ALTER TABLE nope ADD CONSTRAINT c2 PRIMARY KEY USING INDEX old_index; ALTER TABLE tech_book ADD CONSTRAINT c3 PRIMARY KEY USING INDEX "OLD_INDEX"; ALTER TABLE parted_p0 ADD CONSTRAINT c4 PRIMARY KEY USING INDEX parted_idx; ALTER TABLE book ADD CONSTRAINT c5 PRIMARY KEY USING INDEX book_idx;`,
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "c1", "^pk_tech_book_$", "ALTER TABLE tech_book ADD CONSTRAINT c1 PRIMARY KEY USING INDEX nope"),
				keyMismatch("nope", "c2", "^pk_nope_$", "ALTER TABLE nope ADD CONSTRAINT c2 PRIMARY KEY USING INDEX old_index"),
				keyMismatch("tech_book", "c3", "^pk_tech_book_$", `ALTER TABLE tech_book ADD CONSTRAINT c3 PRIMARY KEY USING INDEX "OLD_INDEX"`),
				keyMismatch("parted_p0", "c4", "^pk_parted_p0_$", "ALTER TABLE parted_p0 ADD CONSTRAINT c4 PRIMARY KEY USING INDEX parted_idx"),
				keyMismatch("book", "c5", "^pk_book_$", "ALTER TABLE book ADD CONSTRAINT c5 PRIMARY KEY USING INDEX book_idx"),
			},
		},
		{
			name:   "primary key using index without a snapshot has an empty column list",
			sql:    "ALTER TABLE tech_book ADD CONSTRAINT pk_tech_book_ PRIMARY KEY USING INDEX old_index; ALTER TABLE tech_book ADD CONSTRAINT c PRIMARY KEY USING INDEX old_index;",
			config: format,
			want:   []keyFinding{keyMismatch("tech_book", "c", "^pk_tech_book_$", "ALTER TABLE tech_book ADD CONSTRAINT c PRIMARY KEY USING INDEX old_index")},
		},
		{
			name:   "primary key using index finds the table by object case sensitivity and prints the statement's spelling",
			sql:    "ALTER TABLE book ADD CONSTRAINT c PRIMARY KEY USING INDEX book_idx;",
			config: format,
			schema: keySchema(false),
			want:   []keyFinding{keyMismatch("book", "c", "^pk_book_c$", "ALTER TABLE book ADD CONSTRAINT c PRIMARY KEY USING INDEX book_idx")},
		},
		{
			name:   "create table judges a using index constraint twice, after every table constraint",
			sql:    "CREATE TABLE tech_book (a int, CONSTRAINT c1 PRIMARY KEY USING INDEX old_index, CONSTRAINT c2 PRIMARY KEY (a));",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "c1", "^pk_tech_book_$", "CREATE TABLE tech_book (a int, CONSTRAINT c1 PRIMARY KEY USING INDEX old_index, CONSTRAINT c2 PRIMARY KEY (a))"),
				keyMismatch("tech_book", "c2", "^pk_tech_book_a$", "CREATE TABLE tech_book (a int, CONSTRAINT c1 PRIMARY KEY USING INDEX old_index, CONSTRAINT c2 PRIMARY KEY (a))"),
				keyMismatch("tech_book", "c1", "^pk_tech_book_id_name$", "CREATE TABLE tech_book (a int, CONSTRAINT c1 PRIMARY KEY USING INDEX old_index, CONSTRAINT c2 PRIMARY KEY (a))"),
			},
		},
		{
			name:   "renamed constraint judged when the table's index is primary",
			sql:    "ALTER TABLE tech_book RENAME CONSTRAINT old_pk TO pk_tech_book_id_name; ALTER TABLE tech_book RENAME CONSTRAINT old_pk TO bad; ALTER TABLE tech_book RENAME CONSTRAINT old_uk TO bad2; ALTER TABLE tech_book RENAME CONSTRAINT old_index TO bad3; ALTER TABLE tech_book RENAME CONSTRAINT nope TO bad4; ALTER TABLE twin_b RENAME CONSTRAINT twin_pk TO bad5; ALTER TABLE twin_b RENAME CONSTRAINT old_pk TO bad6; ALTER TABLE parted_p0 RENAME CONSTRAINT parted_pk TO bad7; ALTER TABLE parted RENAME CONSTRAINT parted_pk TO bad8;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^pk_tech_book_id_name$", "ALTER TABLE tech_book RENAME CONSTRAINT old_pk TO bad"),
				keyMismatch("twin_b", "bad5", "^pk_twin_b_b$", "ALTER TABLE twin_b RENAME CONSTRAINT twin_pk TO bad5"),
				keyMismatch("parted", "bad8", "^pk_parted_k$", "ALTER TABLE parted RENAME CONSTRAINT parted_pk TO bad8"),
			},
		},
		{
			name:   "renamed constraint looks in the statement's schema",
			sql:    "ALTER TABLE other.tech_book RENAME CONSTRAINT old_uk TO bad; ALTER TABLE other.tech_book RENAME CONSTRAINT old_pk TO bad2; ALTER TABLE public.tech_book RENAME CONSTRAINT old_pk TO bad3; ALTER TABLE nope.tech_book RENAME CONSTRAINT old_pk TO bad4;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^pk_tech_book_o$", "ALTER TABLE other.tech_book RENAME CONSTRAINT old_uk TO bad"),
				keyMismatch("tech_book", "bad3", "^pk_tech_book_id_name$", "ALTER TABLE public.tech_book RENAME CONSTRAINT old_pk TO bad3"),
			},
		},
		{
			name:   "renamed constraint finds the table by object case sensitivity and prints the statement's spelling",
			sql:    `ALTER TABLE book RENAME CONSTRAINT book_pk TO bad; ALTER TABLE "PUBLIC"."BOOK" RENAME CONSTRAINT book_pk TO bad2; ALTER TABLE tech_book RENAME CONSTRAINT "OLD_PK" TO bad3; ALTER TABLE tech_book RENAME CONSTRAINT "Old_Pk" TO bad4;`,
			config: format,
			schema: keySchema(false),
			want: []keyFinding{
				keyMismatch("book", "bad", "^pk_book_a$", "ALTER TABLE book RENAME CONSTRAINT book_pk TO bad"),
				keyMismatch("BOOK", "bad2", "^pk_BOOK_a$", `ALTER TABLE "PUBLIC"."BOOK" RENAME CONSTRAINT book_pk TO bad2`),
				keyMismatch("tech_book", "bad4", "^pk_tech_book_z$", `ALTER TABLE tech_book RENAME CONSTRAINT "Old_Pk" TO bad4`),
			},
		},
		{
			name:   "table name case sensitive when objects are",
			sql:    "ALTER TABLE book RENAME CONSTRAINT book_pk TO bad;",
			config: format,
			schema: keySchema(true),
		},
		{
			name:   "renamed index judged when primary, with the snapshot's table",
			sql:    "ALTER INDEX old_pk RENAME TO pk_tech_book_id_name; ALTER INDEX old_pk RENAME TO bad; ALTER INDEX old_uk RENAME TO bad2; ALTER INDEX old_index RENAME TO bad3; ALTER INDEX nope RENAME TO bad4; ALTER INDEX book_pk RENAME TO bad5; ALTER INDEX parted_pk RENAME TO bad6;",
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^pk_tech_book_id_name$", "ALTER INDEX old_pk RENAME TO bad"),
				keyMismatch("Book", "bad5", "^pk_Book_a$", "ALTER INDEX book_pk RENAME TO bad5"),
				keyMismatch("parted", "bad6", "^pk_parted_k$", "ALTER INDEX parted_pk RENAME TO bad6"),
			},
		},
		{
			name:   "renamed index looks in the statement's schema",
			sql:    `ALTER INDEX other.old_uk RENAME TO bad; ALTER INDEX other.old_pk RENAME TO bad2; ALTER INDEX other.other_pk RENAME TO bad3; ALTER INDEX other_pk RENAME TO bad4; ALTER INDEX "Other".other_pk RENAME TO bad5;`,
			config: format,
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^pk_tech_book_o$", "ALTER INDEX other.old_uk RENAME TO bad"),
				keyMismatch("tech_book", "bad3", "^pk_tech_book_o$", "ALTER INDEX other.other_pk RENAME TO bad3"),
			},
		},
		{
			name:   "renamed index, two tables hold the name, the first in snapshot order is taken",
			sql:    "ALTER INDEX twin_pk RENAME TO bad;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("twin_a", "bad", "^pk_twin_a_a$", "ALTER INDEX twin_pk RENAME TO bad")},
		},
		{
			name:   "renames without a snapshot not judged",
			sql:    "ALTER INDEX old_pk RENAME TO bad; ALTER TABLE tech_book RENAME CONSTRAINT old_pk TO bad2;",
			config: format,
		},
		{
			name:   "other renames not judged",
			sql:    "ALTER TABLE old_pk RENAME TO bad; ALTER TABLE tech_book RENAME COLUMN old_pk TO bad; ALTER DOMAIN old_pk RENAME CONSTRAINT old_pk TO bad;",
			config: format,
			schema: keySchema(true),
		},
		{
			name:   "create foreign table, create schema element, and routine body not judged",
			sql:    "CREATE FOREIGN TABLE ft (a int, CONSTRAINT bad PRIMARY KEY (a)) SERVER srv; CREATE SCHEMA s CREATE TABLE t (a int, CONSTRAINT bad PRIMARY KEY (a)); CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC ALTER TABLE t ADD CONSTRAINT bad PRIMARY KEY (a); END;",
			config: format,
		},
		{
			name:   "value that keeps the pattern from compiling",
			sql:    `ALTER TABLE "a(" ADD CONSTRAINT "` + strings.Repeat("k", 20) + `" PRIMARY KEY (x); ALTER TABLE t ADD CONSTRAINT nope PRIMARY KEY (x);`,
			config: v1review.TemplateConfig{Format: "^pk_{{table}}$", MaxLength: 10},
			want: []keyFinding{
				{
					Code: v1review.PatternInvalid, Table: "a(", Name: strings.Repeat("k", 20),
					Detail: "error parsing regexp: missing closing ): `^pk_a($`",
					Text:   `ALTER TABLE "a(" ADD CONSTRAINT "` + strings.Repeat("k", 20) + `" PRIMARY KEY (x)`,
				},
				keyMismatch("t", "nope", "^pk_t$", "ALTER TABLE t ADD CONSTRAINT nope PRIMARY KEY (x)"),
			},
		},
		{
			name:   "length counts bytes of the truncated name",
			sql:    `ALTER TABLE t ADD CONSTRAINT "主键名" PRIMARY KEY (x), ADD CONSTRAINT bcdefghi PRIMARY KEY (x), ADD CONSTRAINT "` + strings.Repeat("a", 70) + `" PRIMARY KEY (x);`,
			config: v1review.TemplateConfig{Format: "^(a{63}|[^a]*)$", MaxLength: 8},
			want: []keyFinding{
				{Code: v1review.NameTooLong, Table: "t", Name: "主键名", Text: `ALTER TABLE t ADD CONSTRAINT "主键名" PRIMARY KEY (x), ADD CONSTRAINT bcdefghi PRIMARY KEY (x), ADD CONSTRAINT "` + strings.Repeat("a", 70) + `" PRIMARY KEY (x)`},
				{Code: v1review.NameTooLong, Table: "t", Name: strings.Repeat("a", 63), Text: `ALTER TABLE t ADD CONSTRAINT "主键名" PRIMARY KEY (x), ADD CONSTRAINT bcdefghi PRIMARY KEY (x), ADD CONSTRAINT "` + strings.Repeat("a", 70) + `" PRIMARY KEY (x)`},
			},
		},
		{
			name:   "pattern then length",
			sql:    "ALTER TABLE t ADD CONSTRAINT toolong PRIMARY KEY (x);",
			config: v1review.TemplateConfig{Format: "^pk$", MaxLength: 4},
			want: []keyFinding{
				keyMismatch("t", "toolong", "^pk$", "ALTER TABLE t ADD CONSTRAINT toolong PRIMARY KEY (x)"),
				{Code: v1review.NameTooLong, Table: "t", Name: "toolong", Text: "ALTER TABLE t ADD CONSTRAINT toolong PRIMARY KEY (x)"},
			},
		},
		{name: "empty config requires nothing", sql: "ALTER TABLE t ADD CONSTRAINT anything PRIMARY KEY (x);"},
		{
			name:   "anchor excludes a leading comment and the semicolon",
			sql:    "-- c\nALTER TABLE t\n  ADD CONSTRAINT bad PRIMARY KEY (x);",
			config: format,
			want:   []keyFinding{keyMismatch("t", "bad", "^pk_t_x$", "ALTER TABLE t\n  ADD CONSTRAINT bad PRIMARY KEY (x)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runKeyRule(t, NamingPrimaryKey, tt.sql, tt.config, tt.schema)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
