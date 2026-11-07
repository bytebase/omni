package catalog

import (
	"errors"
	"testing"
)

func TestNativeMariaDBColumnSemantics(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	for _, tc := range []struct{ name, setup, sql string }{
		{"inline_reference", "CREATE TABLE native_parent(id INT PRIMARY KEY)", "CREATE TABLE native_feature(a INT REFERENCES native_parent(id))"},
		{"cross_column_check", "", "CREATE TABLE native_feature(a INT CHECK(a<b),b INT)"},
		{"unknown_check_column", "", "CREATE TABLE native_feature(a INT CHECK(z>0))"},
		{"duplicate_check_name", "", "CREATE TABLE native_feature(a INT,CONSTRAINT same_name CHECK(a>0),CONSTRAINT same_name CHECK(a<10))"},
		{"json_alias", "", "CREATE TABLE native_feature(a JSON)"},
		{"json_alias_with_check", "", "CREATE TABLE native_feature(a JSON CHECK(JSON_LENGTH(a)>0))"},
		{"virtual_current_time", "", "CREATE TABLE native_feature(a DATETIME AS(NOW()) VIRTUAL)"},
		{"stored_current_time", "", "CREATE TABLE native_feature(a DATETIME AS(NOW()) STORED)"},
		{"indexed_current_time", "", "CREATE TABLE native_feature(a DATETIME AS(NOW()) VIRTUAL,KEY(a))"},
		{"virtual_child_fk", "CREATE TABLE native_parent(id INT PRIMARY KEY)", "CREATE TABLE native_feature(a INT,b INT AS(a+1) VIRTUAL,FOREIGN KEY(b) REFERENCES native_parent(id))"},
		{"text_unique_key", "", "CREATE TABLE native_feature(a TEXT,UNIQUE KEY(a))"},
		{"inline_text_unique_key", "", "CREATE TABLE native_feature(a TEXT UNIQUE)"},
		{"alter_inline_text_unique_key", "CREATE TABLE native_feature(a INT)", "ALTER TABLE native_feature ADD b TEXT UNIQUE"},
		{"missing_fulltext_parser", "", "CREATE TABLE native_feature(a TEXT,FULLTEXT KEY(a) WITH PARSER ngram)"},
		{"alter_missing_fulltext_parser", "CREATE TABLE native_feature(a TEXT)", "ALTER TABLE native_feature ADD FULLTEXT KEY k(a) WITH PARSER ngram"},
		{"create_missing_fulltext_parser", "CREATE TABLE native_feature(a TEXT)", "CREATE FULLTEXT INDEX k ON native_feature(a) WITH PARSER ngram"},
		{"text_primary_key", "", "CREATE TABLE native_feature(a TEXT PRIMARY KEY)"},
		{"explicit_null_primary_key", "", "CREATE TABLE native_feature(a INT NULL PRIMARY KEY)"},
		{"binary_key_prefix", "", "CREATE TABLE native_feature(a BLOB,KEY(a))"},
		{"latin_key_prefix", "", "CREATE TABLE native_feature(a TEXT CHARACTER SET latin1,KEY(a))"},
		{"tiny_key_prefix", "", "CREATE TABLE native_feature(a TINYTEXT,KEY(a))"},
		{"text_key_prefix", "", "CREATE TABLE native_feature(a TEXT,KEY(a))"},
		{"datetime_default_precision", "", "CREATE TABLE native_feature(a DATETIME(3) DEFAULT NOW(6))"},
		{"sysdate_default", "", "CREATE TABLE native_feature(a DATETIME DEFAULT SYSDATE())"},
		{"utc_default", "", "CREATE TABLE native_feature(a DATETIME DEFAULT UTC_TIMESTAMP())"},
		{"year_two_digit", "", "CREATE TABLE native_feature(a YEAR(2))"},
		{"year_ignored_width", "", "CREATE TABLE native_feature(a YEAR(5))"},
		{"compression_option", "", "CREATE TABLE native_feature(a INT) COMPRESSION='ZLIB'"},
		{"alter_json_column", "CREATE TABLE native_feature(a INT)", "ALTER TABLE native_feature ADD b JSON"},
		{"alter_inline_reference", "CREATE TABLE native_parent(id INT PRIMARY KEY); CREATE TABLE native_feature(a INT)", "ALTER TABLE native_feature ADD b INT REFERENCES native_parent(id)"},
		{"alter_inline_check", "CREATE TABLE native_feature(a INT)", "ALTER TABLE native_feature ADD b INT CHECK(b>a)"},
		{"alter_blob_default", "CREATE TABLE native_feature(a BLOB)", "ALTER TABLE native_feature ALTER a SET DEFAULT 'abc'"},
		{"alter_text_index", "CREATE TABLE native_feature(a TEXT)", "CREATE INDEX k ON native_feature(a)"},
		{"alter_unique_text_index", "CREATE TABLE native_feature(a TEXT)", "ALTER TABLE native_feature ADD UNIQUE KEY k(a)"},
		{"alter_nondeterministic_index", "CREATE TABLE native_feature(a DATETIME AS(NOW()) VIRTUAL)", "ALTER TABLE native_feature ADD KEY k(a)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, sql := range []string{"DROP TABLE IF EXISTS native_feature", "DROP TABLE IF EXISTS native_parent"} {
				if err := ctr.execSQL(sql); err != nil {
					t.Fatal(err)
				}
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE test; USE test", nil); err != nil {
				t.Fatal(err)
			}
			if tc.setup != "" {
				if err := ctr.execSQL(tc.setup); err != nil {
					t.Fatal(err)
				}
				results, err := c.Exec(tc.setup, nil)
				if err != nil || results[0].Error != nil {
					t.Fatalf("setup: %v %v", err, results)
				}
			}
			nativeErr := ctr.execSQL(tc.sql)
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			modeledErr := results[0].Error
			if nativeErr != nil {
				var modeled *Error
				if !errors.As(modeledErr, &modeled) || modeled.Code != int(mariaOracleErrorCode(nativeErr)) {
					t.Fatalf("native: %v; catalog: %v", nativeErr, modeledErr)
				}
				return
			}
			if modeledErr != nil {
				native, _ := ctr.showCreateTable("native_feature")
				t.Fatalf("native accepts: %s; catalog: %v", native, modeledErr)
			}
			native, err := ctr.showCreateTable("native_feature")
			if err != nil {
				t.Fatal(err)
			}
			modeled := c.ShowCreateTable("test", "native_feature")
			if normalizeWhitespace(native) != normalizeWhitespace(modeled) {
				t.Errorf("native:\n%s\ncatalog:\n%s", native, modeled)
			}
			constraints, err := ctr.queryConstraints("test", "native_feature")
			if err != nil {
				t.Fatal(err)
			}
			tbl := c.GetDatabase("test").GetTable("native_feature")
			if len(constraints) != len(tbl.Constraints) {
				t.Fatalf("native constraints: %v; catalog: %v", constraints, tbl.Constraints)
			}
			for _, native := range constraints {
				found := false
				for _, modeled := range tbl.Constraints {
					if native.Name == modeled.Name {
						found = true
					}
				}
				if !found {
					t.Errorf("native constraint %q missing in catalog", native.Name)
				}
			}
		})
	}
}
