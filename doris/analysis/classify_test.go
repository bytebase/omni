package analysis

import (
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
		{"with cte as (select 1) select * from cte", QueryTypeSelect},

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
		{"# hash comment\nINSERT INTO t VALUES (1)", QueryTypeDML},
		{"-- leading comment\nCREATE TABLE t (id INT)", QueryTypeDDL},

		// Unknown / edge cases
		{"", QueryTypeUnknown},
		{"FOOBAR", QueryTypeUnknown},
		{"42", QueryTypeUnknown},
		{"(SELECT 1)", QueryTypeUnknown}, // leading '(' is not a keyword
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
