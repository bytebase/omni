package analysis

import (
	"strings"
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		input string
		want  QueryType
	}{
		// SELECT / WITH
		{"SELECT 1", QueryTypeSelect},
		{"select * from t", QueryTypeSelect},
		{"WITH cte AS (SELECT 1) SELECT * FROM cte", QueryTypeSelect},
		{"((SELECT 1))", QueryTypeSelect},
		{"(SELECT 1) UNION SELECT 2", QueryTypeSelect},
		{"with cte as (select 1) select * from cte", QueryTypeSelect},
		// VALUES is a query primary (DorisParser.g4:1186), so its rows are a
		// SELECT-shaped read, not a DML keyword.
		{"VALUES(1)", QueryTypeSelect},
		{"values(1), (2)", QueryTypeSelect},

		// Info-schema / admin reads
		{"SHOW TABLES", QueryTypeSelectInfoSchema},
		{"show databases", QueryTypeSelectInfoSchema},
		{"DESCRIBE t", QueryTypeSelectInfoSchema},
		{"describe my_table", QueryTypeSelectInfoSchema},
		{"DESC t", QueryTypeSelectInfoSchema},
		{"desc my_table", QueryTypeSelectInfoSchema},
		{"EXPLAIN SELECT 1", QueryTypeSelectInfoSchema},
		{"explain select * from t", QueryTypeSelectInfoSchema},
		{"HELP", QueryTypeSelectInfoSchema},
		{"help topics", QueryTypeSelectInfoSchema},

		// DML
		{"INSERT INTO t VALUES (1)", QueryTypeDML},
		{"insert into t values (1)", QueryTypeDML},
		{"UPDATE t SET c=1", QueryTypeDML},
		{"update t set c=1 where id=2", QueryTypeDML},
		{"DELETE FROM t", QueryTypeDML},
		{"delete from t where id=1", QueryTypeDML},
		{"TRUNCATE TABLE t", QueryTypeDML},
		{"truncate table t", QueryTypeDML},
		{"MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE", QueryTypeDML},
		{"LOAD DATA INFILE 'f' INTO TABLE t", QueryTypeDML},
		{"EXPORT TABLE t TO 's3://bucket'", QueryTypeDML},
		{"COPY INTO t FROM 's3://bucket'", QueryTypeDML},

		// DDL / session control
		{"CREATE TABLE t (id INT)", QueryTypeDDL},
		{"create table t (id int)", QueryTypeDDL},
		{"ALTER TABLE t ADD COLUMN c INT", QueryTypeDDL},
		{"alter table t drop column c", QueryTypeDDL},
		{"DROP TABLE t", QueryTypeDDL},
		{"drop view v", QueryTypeDDL},
		{"GRANT SELECT ON t TO u", QueryTypeDDL},
		{"REVOKE INSERT ON t FROM u", QueryTypeDDL},
		{"SET x = 1", QueryTypeDDL},
		{"set names utf8", QueryTypeDDL},
		{"UNSET x", QueryTypeDDL},
		{"BEGIN", QueryTypeDDL},
		{"BEGIN TRANSACTION", QueryTypeDDL},
		{"COMMIT", QueryTypeDDL},
		{"ROLLBACK", QueryTypeDDL},
		{"USE db", QueryTypeDDL},
		{"use mydb", QueryTypeDDL},
		{"KILL 123", QueryTypeDDL},
		{"ADMIN SET FRONTEND CONFIG ('key'='value')", QueryTypeDDL},
		{"ANALYZE TABLE t", QueryTypeDDL},
		{"REFRESH CATALOG my_catalog", QueryTypeDDL},
		{"CANCEL LOAD FROM db WHERE LABEL = 'lbl'", QueryTypeDDL},
		{"RECOVER TABLE t", QueryTypeDDL},
		{"CLEAN LABEL 'lbl' FROM db", QueryTypeDDL},
		{"BACKUP SNAPSHOT snap TO repo", QueryTypeDDL},
		{"RESTORE SNAPSHOT snap FROM repo", QueryTypeDDL},
		{"LOCK TABLES t READ", QueryTypeDDL},
		{"UNLOCK TABLES", QueryTypeDDL},
		{"INSTALL PLUGIN FROM '/path/plugin.so'", QueryTypeDDL},
		{"UNINSTALL PLUGIN plugin_name", QueryTypeDDL},
		{"SYNC", QueryTypeDDL},
		{"PAUSE ROUTINE LOAD FOR job", QueryTypeDDL},
		{"RESUME ROUTINE LOAD FOR job", QueryTypeDDL},
		{"WARM UP CLUSTER c WITH TABLE t", QueryTypeDDL},

		// Comment skipping
		{"-- comment\nSELECT 1", QueryTypeSelect},
		{"/* block comment */ SELECT 1", QueryTypeSelect},
		// '#' is no comment in either engine, so this is not the INSERT.
		{"# hash comment\nINSERT INTO t VALUES (1)", QueryTypeUnknown},
		{"-- leading comment\nCREATE TABLE t (id INT)", QueryTypeDDL},

		// Unknown / edge cases
		{"", QueryTypeUnknown},
		{"FOOBAR", QueryTypeUnknown},
		{"42", QueryTypeUnknown},
		// Query-grouping parens are transparent: (SELECT 1) is a SELECT
		// (engine-verified; the paren dispatch parses it as one).
		{"(SELECT 1)", QueryTypeSelect},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := Classify(tc.input)
			if got != tc.want {
				t.Errorf("Classify(%q) = %s, want %s", tc.input, got, tc.want)
			}
		})
	}
}

func TestQueryTypeString(t *testing.T) {
	tests := []struct {
		qt   QueryType
		want string
	}{
		{QueryTypeUnknown, "UNKNOWN"},
		{QueryTypeSelect, "SELECT"},
		{QueryTypeSelectInfoSchema, "SELECT_INFO_SCHEMA"},
		{QueryTypeDML, "DML"},
		{QueryTypeDDL, "DDL"},
		{QueryType(999), "UNKNOWN"},
	}
	for _, tc := range tests {
		got := tc.qt.String()
		if got != tc.want {
			t.Errorf("QueryType(%d).String() = %q, want %q", int(tc.qt), got, tc.want)
		}
	}
}

func TestClassifyStatementCTEPrefixedDML(t *testing.T) {
	// Doris prefixes INSERT, UPDATE and DELETE with a WITH clause
	// (DorisParser.g4:138,145,149). The leading keyword says SELECT, so only
	// the tree's real operation keeps a write from being permitted as a read.
	for _, tc := range []struct {
		statement string
		want      QueryType
	}{
		{"WITH c AS (SELECT 1 AS a) INSERT INTO t (a) SELECT a FROM c", QueryTypeDML},
		{"WITH c AS (SELECT 1 AS a) UPDATE t SET a = 2 WHERE a IN (SELECT a FROM c)", QueryTypeDML},
		{"WITH c AS (SELECT 1 AS a) DELETE FROM t WHERE a IN (SELECT a FROM c)", QueryTypeDML},
		{"EXPLAIN WITH c AS (SELECT 1 AS a) DELETE FROM t", QueryTypeDML},
		{"WITH c AS (SELECT 1 AS a) SELECT * FROM c", QueryTypeSelect},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			if got := ClassifyStatement(tc.statement); got != tc.want {
				t.Errorf("ClassifyStatement(%q) = %s, want %s", tc.statement, got, tc.want)
			}
			// The keyword classifier alone cannot tell: WITH is its first
			// token, so it answers SELECT for the forms that begin with one.
			// That fallback answer is what the AST path replaces.
			if strings.HasPrefix(tc.statement, "WITH") {
				if got := Classify(tc.statement); got != QueryTypeSelect {
					t.Errorf("Classify(%q) = %s, want the keyword fallback %s", tc.statement, got, QueryTypeSelect)
				}
			}
		})
	}
}

// A load-property statement is DML, as a LOAD is; before it parsed, the
// keyword fallback classified `DELETE ON 1` as DML.
func TestClassifyStatementLoadProperty(t *testing.T) {
	for _, sql := range []string{"DELETE ON 1", "ORDER BY a", "COLUMNS(a)"} {
		if got := ClassifyStatement(sql); got != QueryTypeDML {
			t.Errorf("ClassifyStatement(%q) = %s, want DML", sql, got)
		}
	}
}

// COPY INTO loads a table from a stage, so the tree answers DML where the
// keyword fallback does; the procedure statements are DDL, and SHOW CREATE
// PROCEDURE is the info-schema read every SHOW is.
func TestClassifyStatementCopyAndProcedures(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      QueryType
	}{
		{"COPY INTO t FROM @stage", QueryTypeDML},
		{"COPY INTO t(c) FROM (SELECT a FROM @stage)", QueryTypeDML},
		{"ALTER PROCEDURE a()", QueryTypeDDL},
		{"CREATE PROCEDURE a()", QueryTypeDDL},
		{"REPLACE PROCEDURE a()", QueryTypeDDL},
		{"DROP PROC a", QueryTypeDDL},
		{"SHOW CREATE PROCEDURE a", QueryTypeSelectInfoSchema},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			if got := ClassifyStatement(tc.statement); got != tc.want {
				t.Errorf("ClassifyStatement(%q) = %s, want %s", tc.statement, got, tc.want)
			}
		})
	}
}
