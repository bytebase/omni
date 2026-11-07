package catalog

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/tidb/parser"

	gomysql "github.com/go-sql-driver/mysql"
)

// These operations are accepted by MySQL but have different TiDB semantics.
// Keep the catalog's error code and its schema unchanged on native rejection.
func TestNativeTiDBDDLRejections(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	cases := []struct{ name, setup, sql string }{
		{"index_nondeterministic", "", "CREATE TABLE native_reject(a INT,INDEX k((a+RAND())))"},
		{"index_bare_column", "", "CREATE TABLE native_reject(a INT,INDEX k((a)))"},
		{"index_json_result", "", "CREATE TABLE native_reject(doc JSON,INDEX k((doc->'$.name')))"},
		{"index_text_result", "", "CREATE TABLE native_reject(doc JSON,INDEX k((doc->>'$.name')))"},
		{"index_unknown_expression_column", "", "CREATE TABLE native_reject(a INT,INDEX k((missing+1)))"},
		{"alter_index_nondeterministic", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD INDEX k((a+RAND()))"},
		{"standalone_index_nondeterministic", "CREATE TABLE native_reject(a INT)", "CREATE INDEX k ON native_reject((a+RAND()))"},
		{"index_text_function", "", "CREATE TABLE native_reject(a TEXT,INDEX k((LOWER(a))))"},
		{"index_text_concat", "", "CREATE TABLE native_reject(a TEXT,INDEX k((CONCAT(a,'x'))))"},
		{"index_text_without_prefix", "", "CREATE TABLE native_reject(a TEXT,KEY k(a))"},
		{"index_character_prefix_too_long", "", "CREATE TABLE native_reject(a VARCHAR(1000),KEY k(a(1000)))"},
		{"index_composite_too_long", "", "CREATE TABLE native_reject(a VARCHAR(500),b VARCHAR(500),KEY k(a,b))"},
		{"index_aggregate", "", "CREATE TABLE native_reject(a INT,INDEX k((SUM(a))))"},
		{"alter_index_bare_column", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD INDEX k((a))"},
		{"standalone_index_json_result", "CREATE TABLE native_reject(doc JSON)", "CREATE INDEX k ON native_reject((JSON_EXTRACT(doc,'$.name')))"},
		{"alter_index_text_result", "CREATE TABLE native_reject(doc JSON)", "ALTER TABLE native_reject ADD INDEX k((JSON_UNQUOTE(doc->'$.name')))"},
		{"default_function", "", "CREATE TABLE native_reject (a INT DEFAULT (FLOOR(RAND()*100)))"},
		{"default_substring_keyword", "", "CREATE TABLE native_reject(a BLOB DEFAULT (SUBSTRING('abcdef',1,3)))"},
		{"default_trim_keyword", "", "CREATE TABLE native_reject(a TEXT DEFAULT (TRIM('a')))"},
		{"default_substring_quoted", "", "CREATE TABLE native_reject(a BLOB DEFAULT (`SUBSTRING`('abcdef',1,3)))"},
		{"composite_index_column", "CREATE TABLE native_reject (a INT,b INT,KEY k(a,b))", "ALTER TABLE native_reject DROP COLUMN b"},
		{"composite_primary_column", "CREATE TABLE native_reject (a INT,b INT,PRIMARY KEY(a,b))", "ALTER TABLE native_reject DROP COLUMN a"},
		{"clustered_primary_key", "CREATE TABLE native_reject (a INT PRIMARY KEY)", "ALTER TABLE native_reject DROP PRIMARY KEY"},
		{"charset_conversion", "CREATE TABLE native_reject (a VARCHAR(20))", "ALTER TABLE native_reject DEFAULT CHARACTER SET latin1"},
		{"fulltext_index", "CREATE TABLE native_reject (a TEXT)", "CREATE FULLTEXT INDEX k ON native_reject(a)"},
		{"generated_primary_key", "", "CREATE TABLE native_reject (a INT,b INT AS(a+1),PRIMARY KEY(b))"},
		{"nondeterministic_generated", "", "CREATE TABLE native_reject (a DATETIME AS(NOW()))"},
		{"datetime_precision", "", "CREATE TABLE native_reject(a DATETIME(7))"},
		{"function_precision", "", "CREATE TABLE native_reject(a DATETIME DEFAULT NOW(7))"},
		{"default_precision", "", "CREATE TABLE native_reject(a DATETIME(6) DEFAULT NOW())"},
		{"on_update_precision", "", "CREATE TABLE native_reject(a DATETIME(6) ON UPDATE NOW())"},
		{"date_on_update", "", "CREATE TABLE native_reject(a DATE ON UPDATE NOW())"},
		{"alter_default_function", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ALTER a SET DEFAULT (FLOOR(RAND()*100))"},
		{"alter_default_precision", "CREATE TABLE native_reject(a DATETIME(6))", "ALTER TABLE native_reject ALTER a SET DEFAULT NOW()"},
		{"last_column", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject DROP a"},
		{"readonly_timestamp_setting", "", "SET explicit_defaults_for_timestamp=0"},
		{"range_maxvalue_order", "", "CREATE TABLE native_reject(a INT) PARTITION BY RANGE(a)(PARTITION p0 VALUES LESS THAN(MAXVALUE),PARTITION p1 VALUES LESS THAN(100))"},
		{"partition_noninteger_function", "", "CREATE TABLE native_reject(a VARCHAR(10),b VARCHAR(10)) PARTITION BY RANGE(CONCAT(a,b))(PARTITION p0 VALUES LESS THAN('m'),PARTITION p1 VALUES LESS THAN(MAXVALUE))"},
		{"partition_timestamp", "", "CREATE TABLE native_reject(a TIMESTAMP) PARTITION BY RANGE(a)(PARTITION p0 VALUES LESS THAN(100))"},
		{"partition_constant", "", "CREATE TABLE native_reject(a INT) PARTITION BY HASH(1) PARTITIONS 2"},
		{"partition_division", "", "CREATE TABLE native_reject(a INT) PARTITION BY HASH(a/2) PARTITIONS 2"},
		{"partition_timestamp_arithmetic", "", "CREATE TABLE native_reject(a TIMESTAMP) PARTITION BY HASH(a+1) PARTITIONS 2"},
		{"partition_unknown_column", "", "CREATE TABLE native_reject(a INT) PARTITION BY HASH(CONCAT(missing,a)) PARTITIONS 2"},
		{"foreign_key_set_null", "", "CREATE TABLE native_reject(id INT PRIMARY KEY,parent_id INT NOT NULL,FOREIGN KEY(parent_id) REFERENCES native_reject(id) ON DELETE SET NULL)"},
		{"foreign_key_virtual_child", "", "CREATE TABLE native_reject(id INT PRIMARY KEY,parent_id INT AS(id+1) VIRTUAL,FOREIGN KEY(parent_id) REFERENCES native_reject(id))"},
		{"foreign_key_virtual_parent", "", "CREATE TABLE native_reject(id INT,parent_id INT AS(id+1) VIRTUAL,KEY k(parent_id),FOREIGN KEY(id) REFERENCES native_reject(parent_id))"},
		{"decimal_precision", "", "CREATE TABLE native_reject(a DECIMAL(66,0))"},
		{"decimal_scale", "", "CREATE TABLE native_reject(a DECIMAL(40,31))"},
		{"decimal_scale_above_precision", "", "CREATE TABLE native_reject(a DECIMAL(5,6))"},
		{"varchar_strict_length", "", "CREATE TABLE native_reject(a VARCHAR(65536))"},
		{"explicit_null_primary_key", "", "CREATE TABLE native_reject(a INT NULL PRIMARY KEY)"},
		{"explicit_null_table_primary_key", "", "CREATE TABLE native_reject(a INT NULL,PRIMARY KEY(a))"},
		{"null_default_primary_key", "", "CREATE TABLE native_reject(a INT DEFAULT NULL PRIMARY KEY)"},
		{"null_default_not_null", "", "CREATE TABLE native_reject(a INT NOT NULL DEFAULT NULL)"},
		{"create_index_unknown_column", "", "CREATE TABLE native_reject(a INT,KEY k(missing))"},
		{"reserved_index_name", "", "CREATE TABLE native_reject(a INT,INDEX `PRIMARY`(a))"},
		{"reserved_alter_index_name", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD INDEX `PRIMARY`(a)"},
		{"reserved_standalone_index_name", "CREATE TABLE native_reject(a INT)", "CREATE INDEX `PRIMARY` ON native_reject(a)"},
		{"invalid_sql_mode", "", "SET sql_mode='not_a_mode'"},
		{"multi_add_after_new_column", "CREATE TABLE native_reject(a INT,b INT)", "ALTER TABLE native_reject ADD x INT AFTER a,ADD y INT AFTER x"},
		{"multi_add_index_new_column", "CREATE TABLE native_reject(a INT,b INT)", "ALTER TABLE native_reject ADD c INT,ADD INDEX k(c)"},
		{"multi_parenthesized_key_new_column", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD COLUMN(b INT,c INT,PRIMARY KEY(a),UNIQUE KEY u(b))"},
		{"multi_unique_new_column", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD x INT,ADD UNIQUE INDEX ux(x)"},
		{"multi_primary_new_column", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD x INT NOT NULL,ADD PRIMARY KEY(a,x)"},
		{"multi_change_index_new_name", "CREATE TABLE native_reject(x INT)", "ALTER TABLE native_reject CHANGE x y INT,ADD INDEX k(y)"},
		{"multi_change_after_new_name", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject CHANGE a a_new INT,ADD d INT AFTER a_new"},
		{"multi_foreign_key_new_column", "CREATE TABLE native_check_parent(id INT PRIMARY KEY); CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD parent_id INT,ADD CONSTRAINT fk FOREIGN KEY(parent_id) REFERENCES native_check_parent(id)"},
		{"foreign_key_existing_parent_checks_off", "SET foreign_key_checks=0; CREATE TABLE native_check_parent(a INT)", "CREATE TABLE native_reject(a INT,FOREIGN KEY(a) REFERENCES native_check_parent(missing))"},
		{"add_index_missing_column", "CREATE TABLE native_reject(a INT)", "ALTER TABLE native_reject ADD INDEX k(missing)"},
		{"create_index_missing_column", "CREATE TABLE native_reject(a INT)", "CREATE INDEX k ON native_reject(missing)"},
		{"alter_compression", "CREATE TABLE native_reject(a INT) COMPRESSION='ZLIB'", "ALTER TABLE native_reject COMPRESSION='LZ4'"},
		{"unknown_gipk_setting", "", "SET sql_generate_invisible_primary_key=ON"},
		{"unsupported_invisible_column", "", "CREATE TABLE native_reject(a INT INVISIBLE,b INT)"},
		{"unknown_gipk_visibility", "", "SET show_gipk_in_create_table_and_information_schema=ON"},
		{"global_check_setting_scope", "", "SET tidb_enable_check_constraint=ON"},
		{"session_check_setting_scope", "", "SET SESSION tidb_enable_check_constraint=ON"},
		{"unsupported_create_procedure", "", "CREATE PROCEDURE native_reject_procedure() SELECT 1"},
		{"unsupported_drop_procedure", "", "DROP PROCEDURE native_reject_procedure"},
		{"unsupported_alter_procedure", "", "ALTER PROCEDURE native_reject_procedure COMMENT 'updated'"},
		{"enabled_check_unknown_column", "SET GLOBAL tidb_enable_check_constraint=ON", "CREATE TABLE native_reject(a INT,CHECK(missing>0))"},
		{"enabled_check_subquery", "SET GLOBAL tidb_enable_check_constraint=ON", "CREATE TABLE native_reject(a INT,CHECK(a>(SELECT 1)))"},
		{"enabled_check_nondeterministic", "SET GLOBAL tidb_enable_check_constraint=ON", "CREATE TABLE native_reject(a INT,CHECK(a>RAND()))"},
		{"enabled_check_variable", "SET GLOBAL tidb_enable_check_constraint=ON", "CREATE TABLE native_reject(a INT,CHECK(a>@x))"},
		{"enabled_check_duplicate", "SET GLOBAL tidb_enable_check_constraint=ON", "CREATE TABLE native_reject(a INT,CONSTRAINT ck CHECK(a>0),CONSTRAINT ck CHECK(a<10))"},
		{"enabled_check_schema_duplicate", "SET GLOBAL tidb_enable_check_constraint=ON; CREATE TABLE native_check_parent(a INT,CONSTRAINT ck CHECK(a>0))", "CREATE TABLE native_reject(a INT,CONSTRAINT ck CHECK(a<10))"},
		{"enabled_check_inline_other_column", "SET GLOBAL tidb_enable_check_constraint=ON", "CREATE TABLE native_reject(a INT CHECK(a<b),b INT)"},
		{"enabled_check_auto_increment", "SET GLOBAL tidb_enable_check_constraint=ON", "CREATE TABLE native_reject(a INT PRIMARY KEY AUTO_INCREMENT,CHECK(a>0))"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ctr.execSQL("DROP TABLE IF EXISTS native_reject"); err != nil {
				t.Fatal(err)
			}
			if err := ctr.execSQL("DROP TABLE IF EXISTS native_check_parent"); err != nil {
				t.Fatal(err)
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
				if err != nil {
					t.Fatalf("catalog setup: %v %v", err, results)
				}
				for _, result := range results {
					if result.Error != nil {
						t.Fatalf("catalog setup statement %d: %v", result.Index, result.Error)
					}
				}
			}
			before := c.ShowCreateTable("test", "native_reject")
			beforeChecks, beforeFKChecks, beforeStrict := c.checkConstraints, c.foreignKeyChecks, c.strictDDL
			nativeErr := ctr.execSQL(tc.sql)
			var native *gomysql.MySQLError
			if !errors.As(nativeErr, &native) {
				t.Fatalf("native rejection: %v", nativeErr)
			}
			results, err := c.Exec(tc.sql, nil)
			modeledCode := 0
			var modeledErr error
			if err == nil && len(results) > 0 {
				modeledErr = results[0].Error
			} else {
				modeledErr = err
			}
			var modeled *Error
			var parseErr *parser.ParseError
			if errors.As(modeledErr, &modeled) {
				modeledCode = modeled.Code
			} else if errors.As(modeledErr, &parseErr) {
				modeledCode = 1064
			}
			if modeledCode != int(native.Number) {
				t.Errorf("native: %v; catalog: %v", nativeErr, modeledErr)
			}
			if after := c.ShowCreateTable("test", "native_reject"); after != before {
				t.Errorf("rejected DDL changed the catalog: before=%s after=%s", before, after)
			}
			if c.checkConstraints != beforeChecks || c.foreignKeyChecks != beforeFKChecks || c.strictDDL != beforeStrict {
				t.Error("rejected statement changed catalog settings")
			}
		})
	}
}

func TestNativeTiDBColumnReadback(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	for _, tc := range []struct{ name, sql string }{
		{"native_date", "CREATE TABLE native_date(a DATE DEFAULT (CURDATE()))"},
		{"native_datetime", "CREATE TABLE native_datetime(a DATETIME(6) DEFAULT NOW(6) ON UPDATE CURRENT_TIMESTAMP(6))"},
		{"native_hash", "CREATE TABLE native_hash(a INT,KEY k(a) USING HASH)"},
		{"native_key_implicit", "CREATE TABLE native_key_implicit(id INT PRIMARY KEY,v INT) PARTITION BY KEY() PARTITIONS 4"},
		{"native_year", "CREATE TABLE native_year(a YEAR(4))"},
		{"native_range", "CREATE TABLE native_range(a DATE) PARTITION BY RANGE(YEAR(a))(PARTITION p0 VALUES LESS THAN(2020),PARTITION p1 VALUES LESS THAN(MAXVALUE))"},
		{"native_subpartition", "CREATE TABLE native_subpartition(a DATE,b INT) PARTITION BY RANGE(YEAR(a)) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN(2020),PARTITION p1 VALUES LESS THAN(MAXVALUE))"},
		{"native_list_default", "CREATE TABLE native_list_default(a INT) PARTITION BY LIST(a)(PARTITION p0 VALUES IN(1,2),PARTITION pd VALUES IN(DEFAULT))"},
		{"native_hash_div", "CREATE TABLE native_hash_div(a INT) PARTITION BY HASH(a DIV 2) PARTITIONS 2"},
		{"native_compression", "CREATE TABLE native_compression(a INT) COMPRESSION='ZLIB' COMMENT='retained compression'"},
		{"native_compression_copy", "CREATE TABLE native_compression_parent(a INT) COMPRESSION='ZLIB'; CREATE TABLE native_compression_copy LIKE native_compression_parent"},
		{"native_zerofill_copy", "CREATE TABLE native_zerofill_parent(a INT(5) ZEROFILL,b BOOLEAN ZEROFILL,c DECIMAL(10,2) ZEROFILL); CREATE TABLE native_zerofill_copy LIKE native_zerofill_parent"},
		{"native_disabled_check", "CREATE TABLE native_disabled_check(a INT,CHECK(a>RAND()))"},
		{"native_unnamed_expression_indexes", "CREATE TABLE native_unnamed_expression_indexes(a INT,INDEX((a+1)),INDEX((a*2)))"},
		{"native_reserved_column_unique", "CREATE TABLE native_reserved_column_unique(`PRIMARY` INT UNIQUE)"},
		{"native_latin1_index", "CREATE TABLE native_latin1_index(a VARCHAR(1000) CHARACTER SET latin1,INDEX k((LOWER(a))))"},
		{"native_charset_widen", "CREATE TABLE native_charset_widen(a VARCHAR(100)) DEFAULT CHARSET=latin1; ALTER TABLE native_charset_widen DEFAULT CHARACTER SET utf8mb4"},
		{"native_text_sizes", "CREATE TABLE native_text_sizes(a TEXT(1),b TEXT(1000),c TEXT(70000),d TEXT(20000000))"},
		{"native_binary_sizes", "CREATE TABLE native_binary_sizes(a BLOB(1),b BLOB(1000),c BLOB(70000),d BLOB(20000000))"},
		{"native_varchar_nonstrict", "SET SESSION sql_mode=''; CREATE TABLE native_varchar_nonstrict(a VARCHAR(65536))"},
		{"native_index_nonstrict", "SET SESSION sql_mode=''; CREATE TABLE native_index_nonstrict(a VARCHAR(1000),KEY k(a))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("native acceptance: %v", err)
			}
			c := New()
			mustExec(t, c, "CREATE DATABASE test; USE test; "+tc.sql)
			want, err := ctr.showCreateTable(tc.name)
			if err != nil {
				t.Fatal(err)
			}
			if got := c.ShowCreateTable("test", tc.name); normalizeWhitespace(want) != normalizeWhitespace(got) {
				t.Errorf("native:\n%s\ncatalog:\n%s", want, got)
			}
		})
	}
}

func TestNativeTiDBColumnMetadata(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test", nil); err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		"CREATE TABLE native_types(a FLOAT(10),b FLOAT(25),c BOOL,d NUMERIC(10,2),e SERIAL,f YEAR,g INT UNSIGNED,h VARCHAR(20),i TIMESTAMP(6),z INT(5) ZEROFILL,nz NUMERIC(10,2) ZEROFILL,bz BOOLEAN ZEROFILL)",
		"ALTER TABLE native_types ADD j FLOAT(30), ADD k NUMERIC(5,2), ADD l BOOLEAN",
		"ALTER TABLE native_types MODIFY a FLOAT(25), CHANGE b bb FLOAT(10)",
	} {
		if err := ctr.execSQL(sql); err != nil {
			t.Fatal(err)
		}
		results, err := c.Exec(sql, nil)
		if err != nil || len(results) != 1 || results[0].Error != nil {
			t.Fatalf("catalog: %v %v", err, results)
		}
		tbl := c.GetDatabase("test").GetTable("native_types")
		rows, err := ctr.db.QueryContext(ctr.ctx, "SELECT COLUMN_NAME,DATA_TYPE,COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='native_types' ORDER BY ORDINAL_POSITION")
		if err != nil {
			t.Fatal(err)
		}
		count := 0
		for rows.Next() {
			var name, baseType, fullType string
			if err := rows.Scan(&name, &baseType, &fullType); err != nil {
				t.Fatal(err)
			}
			count++
			col := tbl.GetColumn(name)
			if col == nil || col.DataType != baseType || col.ColumnType != fullType {
				t.Errorf("after %s: column %s native %s / %s; catalog %+v", sql, name, baseType, fullType, col)
			}
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		rows.Close()
		if count != len(tbl.Columns) {
			t.Errorf("native has %d columns, catalog has %d", count, len(tbl.Columns))
		}
	}
}

func TestNativeTiDBEnabledCheckReadback(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	if err := ctr.execSQL("SET GLOBAL tidb_enable_check_constraint=ON"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := ctr.execSQL("SET GLOBAL tidb_enable_check_constraint=OFF"); err != nil {
			t.Error(err)
		}
	}()
	for _, statements := range [][]string{
		{"CREATE TABLE native_checks(a INT,CHECK(a>0))"},
		{"CREATE TABLE native_checks(a INT CHECK(a>0),b INT,CHECK(a<b))"},
		{"CREATE TABLE native_checks(a INT,CONSTRAINT native_checks_chk_1 CHECK(a>0),CHECK(a<10))"},
		{"CREATE TABLE native_checks(a INT)", "ALTER TABLE native_checks ADD CONSTRAINT ck CHECK(a>0)", "ALTER TABLE native_checks ALTER CHECK ck NOT ENFORCED", "ALTER TABLE native_checks ALTER CHECK ck ENFORCED", "ALTER TABLE native_checks DROP CHECK ck"},
	} {
		t.Run(statements[0], func(t *testing.T) {
			if err := ctr.execSQL("DROP TABLE IF EXISTS native_checks"); err != nil {
				t.Fatal(err)
			}
			c := New()
			c.SetCheckConstraints(true)
			mustExec(t, c, "CREATE DATABASE test; USE test")
			for _, sql := range statements {
				if err := ctr.execSQL(sql); err != nil {
					t.Fatal(err)
				}
				mustExec(t, c, sql)
				want, err := ctr.showCreateTable("native_checks")
				if err != nil {
					t.Fatal(err)
				}
				got := c.ShowCreateTable("test", "native_checks")
				if normalizeWhitespace(want) != normalizeWhitespace(got) {
					t.Errorf("after %s:\nnative: %s\ncatalog: %s", sql, want, got)
				}
			}
		})
	}
}
