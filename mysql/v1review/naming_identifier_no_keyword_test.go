package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingIdentifierNoKeyword(t *testing.T) {
	// identFinding is a finding as the tests compare it: its code, name,
	// table, pattern and search text, and the text its range anchors.
	type identFinding struct {
		Code    v1review.Code
		Name    string
		Table   string
		Pattern string
		Search  string
		Text    string
	}
	// at is a finding placed from a node; searched is one placed by searching
	// the statement for the column name.
	at := func(name, text string) identFinding {
		return identFinding{Code: v1review.NameIsKeyword, Name: name, Text: text}
	}
	searched := func(name, search, text string) identFinding {
		return identFinding{Code: v1review.NameIsKeyword, Name: name, Search: search, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []identFinding
	}{
		{name: "no keyword", sql: "CREATE TABLE `t`(a int, id int, master_server_id int);\nCREATE INDEX idx_a ON t (a);"},
		{
			name: "table anchors the statement, with no search",
			sql:  "CREATE TABLE execute(id int)",
			want: []identFinding{at("execute", "CREATE TABLE execute(id int)")},
		},
		{
			name: "database qualifier not judged",
			sql:  "CREATE TABLE `select`.t(id int);\nCREATE TABLE `test`.`interval`(id int);",
			want: []identFinding{at("interval", "CREATE TABLE `test`.`interval`(id int)")},
		},
		{
			name: "column carries its search text and the statement's range",
			sql:  "CREATE TABLE t(\n  id int,\n  `current_timestamp` int\n);",
			want: []identFinding{searched("current_timestamp", "current_timestamp", "CREATE TABLE t(\n  id int,\n  `current_timestamp` int\n)")},
		},
		{
			name: "matched in any case and named in source case",
			sql:  "CREATE TABLE `Execute`(`SELECT` int, `Interval` int)",
			want: []identFinding{
				at("Execute", "CREATE TABLE `Execute`(`SELECT` int, `Interval` int)"),
				searched("SELECT", "SELECT", "CREATE TABLE `Execute`(`SELECT` int, `Interval` int)"),
				searched("Interval", "Interval", "CREATE TABLE `Execute`(`SELECT` int, `Interval` int)"),
			},
		},
		{
			name: "non-reserved keywords are keywords",
			sql:  "CREATE TABLE t(status int, name int, data int, comment int)",
			want: []identFinding{
				searched("status", "status", "CREATE TABLE t(status int, name int, data int, comment int)"),
				searched("name", "name", "CREATE TABLE t(status int, name int, data int, comment int)"),
				searched("data", "data", "CREATE TABLE t(status int, name int, data int, comment int)"),
				searched("comment", "comment", "CREATE TABLE t(status int, name int, data int, comment int)"),
			},
		},
		{name: "words the keyword list leaves out", sql: "CREATE TABLE t(absent int, adddate int, auto int, master_compression_algorithm int, selects int, `select ` int)"},
		{
			name: "edge backticks are trimmed from the name and kept in the search",
			sql:  "CREATE TABLE ```execute```(```select``` int, `from``` int, `a``b` int)",
			want: []identFinding{
				at("execute", "CREATE TABLE ```execute```(```select``` int, `from``` int, `a``b` int)"),
				searched("select", "`select`", "CREATE TABLE ```execute```(```select``` int, `from``` int, `a``b` int)"),
				searched("from", "from`", "CREATE TABLE ```execute```(```select``` int, `from``` int, `a``b` int)"),
			},
		},
		{
			name: "table constraints anchor themselves",
			sql:  "CREATE TABLE t(a int, CONSTRAINT `select` PRIMARY KEY (a), KEY `order` (a), UNIQUE `delete` (a), FULLTEXT KEY `by` (a), CONSTRAINT `limit` CHECK (a > 0))",
			want: []identFinding{
				at("select", "CONSTRAINT `select` PRIMARY KEY (a)"),
				at("order", "KEY `order` (a)"),
				at("delete", "UNIQUE `delete` (a)"),
				at("by", "FULLTEXT KEY `by` (a)"),
				at("limit", "CONSTRAINT `limit` CHECK (a > 0)"),
			},
		},
		{
			name: "foreign key judges its symbol and then its index name",
			sql:  "CREATE TABLE t(a int, CONSTRAINT `into` FOREIGN KEY `values` (a) REFERENCES `set` (`update`), CONSTRAINT fk_t_a FOREIGN KEY `execute` (a) REFERENCES p (id))",
			want: []identFinding{
				at("into", "CONSTRAINT `into` FOREIGN KEY `values` (a) REFERENCES `set` (`update`)"),
				at("values", "CONSTRAINT `into` FOREIGN KEY `values` (a) REFERENCES `set` (`update`)"),
				at("execute", "CONSTRAINT fk_t_a FOREIGN KEY `execute` (a) REFERENCES p (id)"),
			},
		},
		{
			name: "unique judges its index name alone, primary key no index name",
			sql:  "CREATE TABLE t(a int, CONSTRAINT `from` UNIQUE KEY `where` (a), PRIMARY KEY `insert` (a))",
			want: []identFinding{at("where", "CONSTRAINT `from` UNIQUE KEY `where` (a)")},
		},
		{name: "column constraints and references not judged", sql: "CREATE TABLE t(a int CONSTRAINT `select` CHECK (a > 0), b int REFERENCES `from` (`where`))"},
		{
			name: "table, then columns, then constraints",
			sql:  "CREATE TABLE `select` (CONSTRAINT `from` CHECK (a > 0), `where` int, `order` int)",
			want: []identFinding{
				at("select", "CREATE TABLE `select` (CONSTRAINT `from` CHECK (a > 0), `where` int, `order` int)"),
				searched("where", "where", "CREATE TABLE `select` (CONSTRAINT `from` CHECK (a > 0), `where` int, `order` int)"),
				searched("order", "order", "CREATE TABLE `select` (CONSTRAINT `from` CHECK (a > 0), `where` int, `order` int)"),
				at("from", "CONSTRAINT `from` CHECK (a > 0)"),
			},
		},
		{
			name: "create table like and as select",
			sql:  "CREATE TABLE `select` LIKE `from`;\nCREATE TABLE `where` (`order` int) AS SELECT 1 AS `group`;",
			want: []identFinding{
				at("select", "CREATE TABLE `select` LIKE `from`"),
				at("where", "CREATE TABLE `where` (`order` int) AS SELECT 1 AS `group`"),
				searched("order", "order", "CREATE TABLE `where` (`order` int) AS SELECT 1 AS `group`"),
			},
		},
		{
			name: "alter add column in both forms anchors the statement",
			sql:  "ALTER TABLE t ADD COLUMN `select` int,\n  ADD (`from` int, `where` int)",
			want: []identFinding{
				at("select", "ALTER TABLE t ADD COLUMN `select` int,\n  ADD (`from` int, `where` int)"),
				at("from", "ALTER TABLE t ADD COLUMN `select` int,\n  ADD (`from` int, `where` int)"),
				at("where", "ALTER TABLE t ADD COLUMN `select` int,\n  ADD (`from` int, `where` int)"),
			},
		},
		{
			name: "alter parenthesized add does not judge its constraints",
			sql:  "ALTER TABLE t ADD COLUMN (`select` int, CONSTRAINT `from` FOREIGN KEY `where` (a) REFERENCES p (id))",
			want: []identFinding{at("select", "ALTER TABLE t ADD COLUMN (`select` int, CONSTRAINT `from` FOREIGN KEY `where` (a) REFERENCES p (id))")},
		},
		{
			name: "alter renames judge the new name",
			sql:  "ALTER TABLE `select` RENAME TO `from`;\nALTER TABLE t RENAME AS db.`where`;\nALTER TABLE t RENAME COLUMN `order` TO `group`;\nALTER TABLE t CHANGE COLUMN `by` `having` int;\nALTER TABLE t RENAME INDEX `limit` TO `into`;",
			want: []identFinding{
				at("from", "ALTER TABLE `select` RENAME TO `from`"),
				at("where", "ALTER TABLE t RENAME AS db.`where`"),
				at("group", "ALTER TABLE t RENAME COLUMN `order` TO `group`"),
				at("having", "ALTER TABLE t CHANGE COLUMN `by` `having` int"),
				at("into", "ALTER TABLE t RENAME INDEX `limit` TO `into`"),
			},
		},
		{
			name: "alter add constraint and add index",
			sql:  "ALTER TABLE t ADD CONSTRAINT `select` FOREIGN KEY `from` (a) REFERENCES `where` (`order`), ADD INDEX `group` (a), ADD UNIQUE `by` (a), ADD CONSTRAINT `having` CHECK (a > 0)",
			want: []identFinding{
				at("select", "ALTER TABLE t ADD CONSTRAINT `select` FOREIGN KEY `from` (a) REFERENCES `where` (`order`), ADD INDEX `group` (a), ADD UNIQUE `by` (a), ADD CONSTRAINT `having` CHECK (a > 0)"),
				at("from", "ALTER TABLE t ADD CONSTRAINT `select` FOREIGN KEY `from` (a) REFERENCES `where` (`order`), ADD INDEX `group` (a), ADD UNIQUE `by` (a), ADD CONSTRAINT `having` CHECK (a > 0)"),
				at("group", "ALTER TABLE t ADD CONSTRAINT `select` FOREIGN KEY `from` (a) REFERENCES `where` (`order`), ADD INDEX `group` (a), ADD UNIQUE `by` (a), ADD CONSTRAINT `having` CHECK (a > 0)"),
				at("by", "ALTER TABLE t ADD CONSTRAINT `select` FOREIGN KEY `from` (a) REFERENCES `where` (`order`), ADD INDEX `group` (a), ADD UNIQUE `by` (a), ADD CONSTRAINT `having` CHECK (a > 0)"),
				at("having", "ALTER TABLE t ADD CONSTRAINT `select` FOREIGN KEY `from` (a) REFERENCES `where` (`order`), ADD INDEX `group` (a), ADD UNIQUE `by` (a), ADD CONSTRAINT `having` CHECK (a > 0)"),
			},
		},
		{
			name: "alter actions not judged",
			sql:  "ALTER TABLE `select` MODIFY COLUMN `from` int, DROP COLUMN `where`, DROP INDEX `order`, ALTER COLUMN `group` SET DEFAULT 1, ALTER INDEX `by` INVISIBLE;",
		},
		{
			name: "create index",
			sql:  "CREATE UNIQUE INDEX `select` ON `from` (`where`);\nCREATE FULLTEXT INDEX execute ON t (a);",
			want: []identFinding{
				at("select", "CREATE UNIQUE INDEX `select` ON `from` (`where`)"),
				at("execute", "CREATE FULLTEXT INDEX execute ON t (a)"),
			},
		},
		{
			name: "other statements not judged",
			sql:  "RENAME TABLE t TO `select`;\nCREATE VIEW `from` AS SELECT 1 AS `where`;\nCREATE DATABASE `order`;\nSELECT 1 AS `group` FROM `by` AS `having`;\nDROP TABLE `limit`;",
		},
		{
			name: "routine, trigger and event bodies not judged",
			sql:  "CREATE PROCEDURE p() BEGIN CREATE TABLE `select` (`from` int); END;\nCREATE TRIGGER `where` BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;\nCREATE EVENT e ON SCHEDULE EVERY 1 DAY DO CREATE TABLE `order` (id int);",
		},
		{
			name: "later statement keeps its own ranges",
			sql:  "CREATE TABLE t (id int);\n\n-- select\nCREATE TABLE u (\n  id int,\n  `select` int,\n  KEY `from` (id)\n);",
			want: []identFinding{
				searched("select", "select", "CREATE TABLE u (\n  id int,\n  `select` int,\n  KEY `from` (id)\n)"),
				at("from", "KEY `from` (id)"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got []identFinding
			for _, f := range NamingIdentifierNoKeyword(parse(t, tt.sql)) {
				got = append(got, identFinding{f.Code, f.Name, f.Table, f.Pattern, f.Search, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
