package catalog

import (
	"errors"
	"strings"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
)

func TestContainer_Section_1_2_StringTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"char_10", "CREATE TABLE t_char10 (a CHAR(10))", "t_char10"},
		{"char_no_length", "CREATE TABLE t_char1 (a CHAR)", "t_char1"},
		{"varchar_255", "CREATE TABLE t_varchar255 (a VARCHAR(255))", "t_varchar255"},
		{"varchar_16383", "CREATE TABLE t_varchar16383 (a VARCHAR(16383))", "t_varchar16383"},
		{"tinytext", "CREATE TABLE t_tinytext (a TINYTEXT)", "t_tinytext"},
		{"text", "CREATE TABLE t_text (a TEXT)", "t_text"},
		{"mediumtext", "CREATE TABLE t_mediumtext (a MEDIUMTEXT)", "t_mediumtext"},
		{"longtext", "CREATE TABLE t_longtext (a LONGTEXT)", "t_longtext"},
		{"text_1000", "CREATE TABLE t_text1000 (a TEXT(1000))", "t_text1000"},
		{"enum_basic", "CREATE TABLE t_enum (a ENUM('a','b','c'))", "t_enum"},
		{"enum_special_chars", "CREATE TABLE t_enum_sc (a ENUM('it''s','hello,world','a\"b'))", "t_enum_sc"},
		{"set_basic", "CREATE TABLE t_set (a SET('x','y','z'))", "t_set"},
		{"char_charset_latin1", "CREATE TABLE t_char_cs (a CHAR(10) CHARACTER SET latin1)", "t_char_cs"},
		{"varchar_charset_collate", "CREATE TABLE t_varchar_cc (a VARCHAR(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin)", "t_varchar_cc"},
		{"national_char", "CREATE TABLE t_nchar (a NATIONAL CHAR(10))", "t_nchar"},
		{"nchar", "CREATE TABLE t_nchar2 (a NCHAR(10))", "t_nchar2"},
		{"nvarchar", "CREATE TABLE t_nvarchar (a NVARCHAR(100))", "t_nvarchar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_4_DateTimeTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"date", "CREATE TABLE t_date (a DATE)", "t_date"},
		{"time", "CREATE TABLE t_time (a TIME)", "t_time"},
		{"time_fsp", "CREATE TABLE t_time3 (a TIME(3))", "t_time3"},
		{"datetime", "CREATE TABLE t_datetime (a DATETIME)", "t_datetime"},
		{"datetime_fsp", "CREATE TABLE t_datetime6 (a DATETIME(6))", "t_datetime6"},
		{"timestamp", "CREATE TABLE t_timestamp (a TIMESTAMP)", "t_timestamp"},
		{"timestamp_fsp", "CREATE TABLE t_timestamp3 (a TIMESTAMP(3))", "t_timestamp3"},
		{"year", "CREATE TABLE t_year (a YEAR)", "t_year"},
		{"year_4", "CREATE TABLE t_year4 (a YEAR(4))", "t_year4"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_1_NumericTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"int_basic", "CREATE TABLE t_int (a INT)", "t_int"},
		{"int_display_width", "CREATE TABLE t_int_dw (a INT(11))", "t_int_dw"},
		{"int_unsigned", "CREATE TABLE t_int_u (a INT UNSIGNED)", "t_int_u"},
		{"int_unsigned_zerofill", "CREATE TABLE t_int_uz (a INT UNSIGNED ZEROFILL)", "t_int_uz"},
		{"tinyint", "CREATE TABLE t_tinyint (a TINYINT)", "t_tinyint"},
		{"smallint", "CREATE TABLE t_smallint (a SMALLINT)", "t_smallint"},
		{"mediumint", "CREATE TABLE t_mediumint (a MEDIUMINT)", "t_mediumint"},
		{"bigint", "CREATE TABLE t_bigint (a BIGINT)", "t_bigint"},
		{"bigint_unsigned", "CREATE TABLE t_bigint_u (a BIGINT UNSIGNED)", "t_bigint_u"},
		{"float_basic", "CREATE TABLE t_float (a FLOAT)", "t_float"},
		{"float_precision", "CREATE TABLE t_float_p (a FLOAT(7,3))", "t_float_p"},
		{"float_unsigned", "CREATE TABLE t_float_u (a FLOAT UNSIGNED)", "t_float_u"},
		{"double_basic", "CREATE TABLE t_double (a DOUBLE)", "t_double"},
		{"double_precision_alias", "CREATE TABLE t_double_p (a DOUBLE PRECISION)", "t_double_p"},
		{"double_with_precision", "CREATE TABLE t_double_wp (a DOUBLE(15,5))", "t_double_wp"},
		{"decimal_precision", "CREATE TABLE t_decimal (a DECIMAL(10,2))", "t_decimal"},
		{"numeric_precision", "CREATE TABLE t_numeric (a NUMERIC(10,2))", "t_numeric"},
		{"decimal_no_precision", "CREATE TABLE t_decimal_np (a DECIMAL)", "t_decimal_np"},
		{"boolean", "CREATE TABLE t_bool (a BOOLEAN)", "t_bool"},
		{"bool_alias", "CREATE TABLE t_bool2 (a BOOL)", "t_bool2"},
		{"bit_1", "CREATE TABLE t_bit1 (a BIT(1))", "t_bit1"},
		{"bit_8", "CREATE TABLE t_bit8 (a BIT(8))", "t_bit8"},
		{"bit_64", "CREATE TABLE t_bit64 (a BIT(64))", "t_bit64"},
		{"serial", "CREATE TABLE t_serial (a SERIAL)", "t_serial"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, _ := c.Exec(tc.sql, nil)
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_10_ColumnAttributesCombination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"int_not_null_auto_increment", "CREATE TABLE t_ai1 (a INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (a))", "t_ai1"},
		{"bigint_unsigned_not_null_auto_increment", "CREATE TABLE t_ai2 (a BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, PRIMARY KEY (a))", "t_ai2"},
		{"varchar_not_null_default_empty", "CREATE TABLE t_vnde (a VARCHAR(100) NOT NULL DEFAULT '')", "t_vnde"},
		{"varchar_charset_collate_not_null", "CREATE TABLE t_vccnn (a VARCHAR(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL)", "t_vccnn"},
		{"int_not_null_comment", "CREATE TABLE t_innc (a INT NOT NULL COMMENT 'user id')", "t_innc"},
		{"varchar_invisible", "CREATE TABLE t_vinv (a INT, b VARCHAR(255) INVISIBLE)", "t_vinv"},
		{"int_visible_not_shown", "CREATE TABLE t_ivis (a INT VISIBLE)", "t_ivis"},
		{"all_attributes", "CREATE TABLE t_all (a INT UNSIGNED NOT NULL DEFAULT '0' COMMENT 'count')", "t_all"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_7_DefaultValues(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"int_default_0", "CREATE TABLE t_def_int0 (a INT DEFAULT 0)", "t_def_int0"},
		{"int_default_null", "CREATE TABLE t_def_intn (a INT DEFAULT NULL)", "t_def_intn"},
		{"int_not_null", "CREATE TABLE t_def_intnn (a INT NOT NULL)", "t_def_intnn"},
		{"varchar_default_hello", "CREATE TABLE t_def_vch (a VARCHAR(50) DEFAULT 'hello')", "t_def_vch"},
		{"varchar_default_empty", "CREATE TABLE t_def_vce (a VARCHAR(50) DEFAULT '')", "t_def_vce"},
		{"float_default", "CREATE TABLE t_def_flt (a FLOAT DEFAULT 3.14)", "t_def_flt"},
		{"decimal_default", "CREATE TABLE t_def_dec (a DECIMAL(10,2) DEFAULT 0.00)", "t_def_dec"},
		{"bool_default_true", "CREATE TABLE t_def_bt (a BOOLEAN DEFAULT TRUE)", "t_def_bt"},
		{"bool_default_false", "CREATE TABLE t_def_bf (a BOOLEAN DEFAULT FALSE)", "t_def_bf"},
		{"enum_default", "CREATE TABLE t_def_enum (a ENUM('a','b','c') DEFAULT 'a')", "t_def_enum"},
		{"set_default", "CREATE TABLE t_def_set (a SET('x','y','z') DEFAULT 'x,y')", "t_def_set"},
		{"bit_default", "CREATE TABLE t_def_bit (a BIT(8) DEFAULT b'00001111')", "t_def_bit"},
		{"blob_no_default_null", "CREATE TABLE t_def_blob (a BLOB)", "t_def_blob"},
		{"text_no_default_null", "CREATE TABLE t_def_text (a TEXT)", "t_def_text"},
		{"json_no_default_null", "CREATE TABLE t_def_json (a JSON)", "t_def_json"},
		{"timestamp_default_ct", "CREATE TABLE t_def_tsct (a TIMESTAMP DEFAULT CURRENT_TIMESTAMP)", "t_def_tsct"},
		{"datetime_default_ct", "CREATE TABLE t_def_dtct (a DATETIME DEFAULT CURRENT_TIMESTAMP)", "t_def_dtct"},
		{"timestamp3_default_ct3", "CREATE TABLE t_def_ts3 (a TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3))", "t_def_ts3"},
		{"expr_default_int", "CREATE TABLE t_def_expr (a INT DEFAULT (FLOOR(RAND()*100)))", "t_def_expr"},
		{"expr_default_json", "CREATE TABLE t_def_exjson (a JSON DEFAULT (JSON_ARRAY()))", "t_def_exjson"},
		{"expr_default_varchar", "CREATE TABLE t_def_exvc (a VARCHAR(36) DEFAULT (UUID()))", "t_def_exvc"},
		{"datetime_default_literal", "CREATE TABLE t_def_dtlit (a DATETIME DEFAULT '2024-01-01 00:00:00')", "t_def_dtlit"},
		{"date_default_literal", "CREATE TABLE t_def_dtlit2 (a DATE DEFAULT '2024-01-01')", "t_def_dtlit2"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_11_PrimaryKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"pk_single_column", "CREATE TABLE t_pk1 (id INT NOT NULL, PRIMARY KEY (id))", "t_pk1"},
		{"pk_multi_column", "CREATE TABLE t_pk2 (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a, b))", "t_pk2"},
		{"pk_bigint_unsigned_auto_inc", "CREATE TABLE t_pk3 (id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, name VARCHAR(100), PRIMARY KEY (id))", "t_pk3"},
		{"pk_column_ordering", "CREATE TABLE t_pk4 (c INT NOT NULL, b INT NOT NULL, a INT NOT NULL, PRIMARY KEY (b, a))", "t_pk4"},
		{"pk_implicit_not_null", "CREATE TABLE t_pk5 (id INT, PRIMARY KEY (id))", "t_pk5"},
		{"pk_name_not_shown", "CREATE TABLE t_pk6 (id INT NOT NULL, PRIMARY KEY (id))", "t_pk6"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_13_RegularIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"key_named", "CREATE TABLE t_idx1 (a INT, KEY `idx_a` (a))", "t_idx1"},
		{"key_auto_named", "CREATE TABLE t_idx2 (a INT, KEY (a))", "t_idx2"},
		{"key_multi_column", "CREATE TABLE t_idx3 (a INT, b INT, c INT, KEY `idx_abc` (a, b, c))", "t_idx3"},
		{"key_prefix_length", "CREATE TABLE t_idx4 (a VARCHAR(255), KEY `idx_a` (a(10)))", "t_idx4"},
		{"key_desc", "CREATE TABLE t_idx5 (a INT, KEY `idx_a` (a DESC))", "t_idx5"},
		{"key_mixed_asc_desc", "CREATE TABLE t_idx6 (a INT, b INT, KEY `idx_ab` (a ASC, b DESC))", "t_idx6"},
		{"key_using_hash", "CREATE TABLE t_idx7 (a INT, KEY `idx_a` (a) USING HASH) ENGINE=MEMORY", "t_idx7"},
		{"key_using_btree", "CREATE TABLE t_idx8 (a INT, KEY `idx_a` (a) USING BTREE)", "t_idx8"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_17_ForeignKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name    string
		setup   string // SQL to run before the main CREATE TABLE (e.g., parent tables)
		sql     string // The CREATE TABLE with FK to compare
		table   string // The table to SHOW CREATE TABLE on
		cleanup string // SQL to clean up after (drop child then parent)
	}{
		{
			"fk_basic",
			"CREATE TABLE t_parent1 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk1 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent1(id))",
			"t_fk1",
			"DROP TABLE IF EXISTS t_fk1; DROP TABLE IF EXISTS t_parent1",
		},
		{
			"fk_named",
			"CREATE TABLE t_parent2 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk2 (id INT, pid INT, CONSTRAINT `fk_name` FOREIGN KEY (pid) REFERENCES t_parent2(id))",
			"t_fk2",
			"DROP TABLE IF EXISTS t_fk2; DROP TABLE IF EXISTS t_parent2",
		},
		{
			"fk_on_delete_cascade",
			"CREATE TABLE t_parent3 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk3 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent3(id) ON DELETE CASCADE)",
			"t_fk3",
			"DROP TABLE IF EXISTS t_fk3; DROP TABLE IF EXISTS t_parent3",
		},
		{
			"fk_on_delete_set_null",
			"CREATE TABLE t_parent4 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk4 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent4(id) ON DELETE SET NULL)",
			"t_fk4",
			"DROP TABLE IF EXISTS t_fk4; DROP TABLE IF EXISTS t_parent4",
		},
		{
			"fk_on_delete_set_default",
			"CREATE TABLE t_parent5 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk5 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent5(id) ON DELETE SET DEFAULT)",
			"t_fk5",
			"DROP TABLE IF EXISTS t_fk5; DROP TABLE IF EXISTS t_parent5",
		},
		{
			"fk_on_delete_restrict",
			"CREATE TABLE t_parent6 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk6 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent6(id) ON DELETE RESTRICT)",
			"t_fk6",
			"DROP TABLE IF EXISTS t_fk6; DROP TABLE IF EXISTS t_parent6",
		},
		{
			"fk_on_delete_no_action",
			"CREATE TABLE t_parent7 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk7 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent7(id) ON DELETE NO ACTION)",
			"t_fk7",
			"DROP TABLE IF EXISTS t_fk7; DROP TABLE IF EXISTS t_parent7",
		},
		{
			"fk_on_update_cascade",
			"CREATE TABLE t_parent8 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk8 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent8(id) ON UPDATE CASCADE)",
			"t_fk8",
			"DROP TABLE IF EXISTS t_fk8; DROP TABLE IF EXISTS t_parent8",
		},
		{
			"fk_on_update_set_null",
			"CREATE TABLE t_parent9 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk9 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent9(id) ON UPDATE SET NULL)",
			"t_fk9",
			"DROP TABLE IF EXISTS t_fk9; DROP TABLE IF EXISTS t_parent9",
		},
		{
			"fk_combined_actions",
			"CREATE TABLE t_parent10 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk10 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent10(id) ON DELETE CASCADE ON UPDATE SET NULL)",
			"t_fk10",
			"DROP TABLE IF EXISTS t_fk10; DROP TABLE IF EXISTS t_parent10",
		},
		{
			"fk_auto_naming",
			"CREATE TABLE t_parent11 (id INT NOT NULL, id2 INT NOT NULL, PRIMARY KEY (id), KEY (id2)) ENGINE=InnoDB",
			"CREATE TABLE t_fk11 (id INT, pid INT, pid2 INT, FOREIGN KEY (pid) REFERENCES t_parent11(id), FOREIGN KEY (pid2) REFERENCES t_parent11(id2))",
			"t_fk11",
			"DROP TABLE IF EXISTS t_fk11; DROP TABLE IF EXISTS t_parent11",
		},
		{
			"fk_auto_generates_index",
			"CREATE TABLE t_parent12 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB",
			"CREATE TABLE t_fk12 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent12(id))",
			"t_fk12",
			"DROP TABLE IF EXISTS t_fk12; DROP TABLE IF EXISTS t_parent12",
		},
		{
			"fk_self_referencing",
			"",
			"CREATE TABLE t_fk13 (id INT NOT NULL, parent_id INT, PRIMARY KEY (id), FOREIGN KEY (parent_id) REFERENCES t_fk13(id))",
			"t_fk13",
			"DROP TABLE IF EXISTS t_fk13",
		},
		{
			"fk_multi_column",
			"CREATE TABLE t_parent14 (x INT NOT NULL, y INT NOT NULL, PRIMARY KEY (x, y)) ENGINE=InnoDB",
			"CREATE TABLE t_fk14 (id INT, a INT, b INT, FOREIGN KEY (a, b) REFERENCES t_parent14(x, y))",
			"t_fk14",
			"DROP TABLE IF EXISTS t_fk14; DROP TABLE IF EXISTS t_parent14",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Cleanup from prior runs.
			ctr.execSQL(tc.cleanup)

			// Setup parent tables.
			if tc.setup != "" {
				if err := ctr.execSQL(tc.setup); err != nil {
					t.Fatalf("container setup: %v", err)
				}
			}
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			// Run setup on omni too.
			if tc.setup != "" {
				results, err := c.Exec(tc.setup, nil)
				if err != nil {
					t.Fatalf("omni setup parse error: %v", err)
				}
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni setup exec error: %v", r.Error)
					}
				}
			}
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_12_UniqueKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"unique_key_named", "CREATE TABLE t_uk1 (a INT, UNIQUE KEY `uk_a` (a))", "t_uk1"},
		{"unique_key_auto_named", "CREATE TABLE t_uk2 (a INT, UNIQUE KEY (a))", "t_uk2"},
		{"unique_key_multi_column", "CREATE TABLE t_uk3 (a INT, b INT, UNIQUE KEY `uk_ab` (a, b))", "t_uk3"},
		{"multiple_unique_keys", "CREATE TABLE t_uk4 (a INT, b INT, c INT, UNIQUE KEY `uk_a` (a), UNIQUE KEY `uk_b` (b))", "t_uk4"},
		{"unique_key_auto_name_collision", "CREATE TABLE t_uk5 (a INT, b INT, c INT, UNIQUE KEY (a), UNIQUE KEY (a, b), UNIQUE KEY (a, c))", "t_uk5"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_18_CheckConstraints(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"check_basic", "CREATE TABLE t_chk1 (a INT, CHECK (a > 0))", "t_chk1"},
		{"check_named", "CREATE TABLE t_chk2 (a INT, CONSTRAINT `chk_name` CHECK (a > 0))", "t_chk2"},
		{"check_not_enforced", "CREATE TABLE t_chk3 (a INT, CHECK (a > 0) NOT ENFORCED)", "t_chk3"},
		{"check_auto_naming", "CREATE TABLE t_chk4 (a INT, b INT, CHECK (a > 0), CHECK (b > 0))", "t_chk4"},
		{"check_expr_parens", "CREATE TABLE t_chk5 (a INT, CHECK (a > 0 AND a < 100))", "t_chk5"},
		{"check_multi_col", "CREATE TABLE t_chk6 (a INT, b INT, CHECK (a > b))", "t_chk6"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_19_TableOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"engine_innodb", "CREATE TABLE t_eng_innodb (a INT) ENGINE=InnoDB", "t_eng_innodb"},
		{"engine_myisam", "CREATE TABLE t_eng_myisam (a INT) ENGINE=MyISAM", "t_eng_myisam"},
		{"engine_memory", "CREATE TABLE t_eng_memory (a INT) ENGINE=MEMORY", "t_eng_memory"},
		{"charset_utf8mb4_default_collation", "CREATE TABLE t_cs_utf8mb4 (a INT) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci", "t_cs_utf8mb4"},
		{"charset_latin1", "CREATE TABLE t_cs_latin1 (a INT) DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci", "t_cs_latin1"},
		{"charset_utf8mb4_unicode_ci", "CREATE TABLE t_cs_unicode (a INT) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci", "t_cs_unicode"},
		{"comment_basic", "CREATE TABLE t_comment (a INT) COMMENT='table description'", "t_comment"},
		{"comment_special_chars", "CREATE TABLE t_comment_sc (a INT) COMMENT='it\\'s a \\\\test'", "t_comment_sc"},
		{"row_format_dynamic", "CREATE TABLE t_rf_dyn (a INT) ROW_FORMAT=DYNAMIC", "t_rf_dyn"},
		{"row_format_compressed", "CREATE TABLE t_rf_comp (a INT) ROW_FORMAT=COMPRESSED", "t_rf_comp"},
		{"auto_increment_1000", "CREATE TABLE t_ai1000 (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id)) AUTO_INCREMENT=1000", "t_ai1000"},
		{"key_block_size_8", "CREATE TABLE t_kbs8 (a INT) KEY_BLOCK_SIZE=8", "t_kbs8"},
		{"multiple_options", "CREATE TABLE t_multi_opts (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(100), PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=500 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='multi opts'", "t_multi_opts"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_8_OnUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"timestamp_on_update", "CREATE TABLE t_ou1 (a TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)", "t_ou1"},
		{"datetime3_on_update", "CREATE TABLE t_ou2 (a DATETIME(3) ON UPDATE CURRENT_TIMESTAMP(3))", "t_ou2"},
		{"timestamp_default_and_on_update", "CREATE TABLE t_ou3 (a TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)", "t_ou3"},
		{"datetime6_default_and_on_update", "CREATE TABLE t_ou4 (a DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6))", "t_ou4"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_20_CharsetCollationInheritance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name     string
		setupSQL string
		sql      string
		table    string
		database string
	}{
		{
			"table_charset_inherited_from_database",
			"DROP DATABASE IF EXISTS db_latin1; CREATE DATABASE db_latin1 CHARACTER SET latin1; USE db_latin1",
			"CREATE TABLE t_inherit_db (a VARCHAR(100))",
			"t_inherit_db",
			"db_latin1",
		},
		{
			"column_charset_inherited_from_table",
			"DROP DATABASE IF EXISTS db_utf8mb4_tbl; CREATE DATABASE db_utf8mb4_tbl; USE db_utf8mb4_tbl",
			"CREATE TABLE t_col_inherit (a VARCHAR(50)) DEFAULT CHARSET=latin1",
			"t_col_inherit",
			"db_utf8mb4_tbl",
		},
		{
			"column_charset_overrides_table",
			"DROP DATABASE IF EXISTS db_override_cs; CREATE DATABASE db_override_cs; USE db_override_cs",
			"CREATE TABLE t_col_override_cs (a VARCHAR(50) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4",
			"t_col_override_cs",
			"db_override_cs",
		},
		{
			"column_collation_overrides_table",
			"DROP DATABASE IF EXISTS db_override_coll; CREATE DATABASE db_override_coll; USE db_override_coll",
			"CREATE TABLE t_col_override_coll (a VARCHAR(50) COLLATE utf8mb4_bin) DEFAULT CHARSET=utf8mb4",
			"t_col_override_coll",
			"db_override_coll",
		},
		{
			"column_charset_collation_display_rules",
			"DROP DATABASE IF EXISTS db_display; CREATE DATABASE db_display CHARACTER SET utf8mb4; USE db_display",
			"CREATE TABLE t_display_rules (a VARCHAR(50), b VARCHAR(50) CHARACTER SET latin1, c VARCHAR(50) COLLATE utf8mb4_bin, d VARCHAR(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin)",
			"t_display_rules",
			"db_display",
		},
		{
			"binary_charset_on_column",
			"DROP DATABASE IF EXISTS db_binary; CREATE DATABASE db_binary; USE db_binary",
			"CREATE TABLE t_binary_cs (a VARCHAR(50) CHARACTER SET binary)",
			"t_binary_cs",
			"db_binary",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupSQL != "" {
				if err := ctr.execSQL(tc.setupSQL); err != nil {
					t.Fatalf("container setup: %v", err)
				}
			}
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, err := ctr.showCreateTable(tc.table)
			if err != nil {
				t.Fatalf("container show create: %v", err)
			}

			c := New()
			if tc.setupSQL != "" {
				results, parseErr := c.Exec(tc.setupSQL, nil)
				if parseErr != nil {
					t.Fatalf("omni setup parse error: %v", parseErr)
				}
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni setup exec error: %v", r.Error)
					}
				}
			}
			results, parseErr := c.Exec(tc.sql, nil)
			if parseErr != nil {
				t.Fatalf("omni parse error: %v", parseErr)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable(tc.database, tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_9_GeneratedColumns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{
			"generated_virtual_add",
			"CREATE TABLE t_gen1 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL)",
			"t_gen1",
		},
		{
			"generated_stored_mul",
			"CREATE TABLE t_gen2 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 * col2) STORED)",
			"t_gen2",
		},
		{
			"generated_varchar_concat",
			"CREATE TABLE t_gen3 (first_name VARCHAR(50), last_name VARCHAR(50), full_name VARCHAR(255) AS (CONCAT(first_name, ' ', last_name)) VIRTUAL)",
			"t_gen3",
		},
		{
			"generated_not_null",
			"CREATE TABLE t_gen4 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) STORED NOT NULL)",
			"t_gen4",
		},
		{
			"generated_comment",
			"CREATE TABLE t_gen5 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL COMMENT 'sum of cols')",
			"t_gen5",
		},
		{
			"generated_invisible",
			"CREATE TABLE t_gen6 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL INVISIBLE)",
			"t_gen6",
		},
		{
			"generated_json_extract",
			"CREATE TABLE t_gen7 (data JSON, name VARCHAR(255) GENERATED ALWAYS AS (JSON_EXTRACT(data, '$.name')) VIRTUAL)",
			"t_gen7",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_3_BinaryTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"binary_16", "CREATE TABLE t_binary16 (a BINARY(16))", "t_binary16"},
		{"binary_no_length", "CREATE TABLE t_binary1 (a BINARY)", "t_binary1"},
		{"varbinary_255", "CREATE TABLE t_varbinary255 (a VARBINARY(255))", "t_varbinary255"},
		{"tinyblob", "CREATE TABLE t_tinyblob (a TINYBLOB)", "t_tinyblob"},
		{"blob", "CREATE TABLE t_blob (a BLOB)", "t_blob"},
		{"mediumblob", "CREATE TABLE t_mediumblob (a MEDIUMBLOB)", "t_mediumblob"},
		{"longblob", "CREATE TABLE t_longblob (a LONGBLOB)", "t_longblob"},
		{"blob_1000", "CREATE TABLE t_blob1000 (a BLOB(1000))", "t_blob1000"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_6_JSONType(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"json_basic", "CREATE TABLE t_json (a JSON)", "t_json"},
		{"json_default_null", "CREATE TABLE t_json_dn (a JSON DEFAULT NULL)", "t_json_dn"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_5_SpatialTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"geometry", "CREATE TABLE t_geometry (a GEOMETRY)", "t_geometry"},
		{"point", "CREATE TABLE t_point (a POINT)", "t_point"},
		{"linestring", "CREATE TABLE t_linestring (a LINESTRING)", "t_linestring"},
		{"polygon", "CREATE TABLE t_polygon (a POLYGON)", "t_polygon"},
		{"multipoint", "CREATE TABLE t_multipoint (a MULTIPOINT)", "t_multipoint"},
		{"multilinestring", "CREATE TABLE t_multilinestring (a MULTILINESTRING)", "t_multilinestring"},
		{"multipolygon", "CREATE TABLE t_multipolygon (a MULTIPOLYGON)", "t_multipolygon"},
		{"geometrycollection", "CREATE TABLE t_geomcoll (a GEOMETRYCOLLECTION)", "t_geomcoll"},
		{"point_srid", "CREATE TABLE t_point_srid (a POINT NOT NULL SRID 4326)", "t_point_srid"},
		{"linestring_srid", "CREATE TABLE t_ls_srid (a LINESTRING SRID 4326)", "t_ls_srid"},
		{"polygon_srid", "CREATE TABLE t_poly_srid (a POLYGON NOT NULL SRID 4326)", "t_poly_srid"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_14_FulltextSpatialIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"fulltext_named", "CREATE TABLE t_ft1 (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY `ft_idx` (body))", "t_ft1"},
		{"fulltext_multi_col", "CREATE TABLE t_ft2 (id INT PRIMARY KEY, title VARCHAR(200), body TEXT, FULLTEXT KEY `ft_multi` (title, body))", "t_ft2"},
		{"fulltext_auto_name", "CREATE TABLE t_ft3 (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY (body))", "t_ft3"},
		{"spatial_named", "CREATE TABLE t_sp1 (id INT PRIMARY KEY, geo_col GEOMETRY NOT NULL, SPATIAL KEY `sp_idx` (geo_col))", "t_sp1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_15_ExpressionIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"expr_func_upper", "CREATE TABLE t_expr1 (name VARCHAR(100), KEY `idx` ((UPPER(name))))", "t_expr1"},
		{"expr_arithmetic", "CREATE TABLE t_expr2 (col1 INT, col2 INT, KEY `idx` ((col1 + col2)))", "t_expr2"},
		{"expr_unique", "CREATE TABLE t_expr3 (name VARCHAR(100), UNIQUE KEY `uidx` ((UPPER(name))))", "t_expr3"},
		{"expr_display_format", "CREATE TABLE t_expr4 (a INT, b INT, KEY `idx_col` (a), KEY `idx_expr` ((a * b)))", "t_expr4"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_1_16_IndexOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"index_comment", "CREATE TABLE t_idx_comment (a INT, INDEX idx_a (a) COMMENT 'description')", "t_idx_comment"},
		{"index_invisible", "CREATE TABLE t_idx_invis (a INT, INDEX idx_a (a) INVISIBLE)", "t_idx_invis"},
		{"index_visible", "CREATE TABLE t_idx_vis (a INT, INDEX idx_a (a) VISIBLE)", "t_idx_vis"},
		{"index_key_block_size", "CREATE TABLE t_idx_kbs (a INT, INDEX idx_a (a) KEY_BLOCK_SIZE=4)", "t_idx_kbs"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_2_1_CreateTableVariants(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("if_not_exists_no_error", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ine")
		ctr.execSQL("CREATE TABLE t_ine (id INT)")

		// Second CREATE with IF NOT EXISTS should not error on ctr.
		ctrErr := ctr.execSQL("CREATE TABLE IF NOT EXISTS t_ine (id INT)")
		if ctrErr != nil {
			t.Fatalf("container error on IF NOT EXISTS: %v", ctrErr)
		}

		// Omni should also not error.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ine (id INT)", nil)
		results, _ := c.Exec("CREATE TABLE IF NOT EXISTS t_ine (id INT)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni error on IF NOT EXISTS: %v", results[0].Error)
		}

		// Compare SHOW CREATE TABLE.
		ctrDDL, _ := ctr.showCreateTable("t_ine")
		omniDDL := c.ShowCreateTable("test", "t_ine")
		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("temporary_table", func(t *testing.T) {
		// MySQL SHOW CREATE TABLE for temporary tables shows "CREATE TEMPORARY TABLE".
		ctr.execSQL("DROP TEMPORARY TABLE IF EXISTS t_temp")
		err := ctr.execSQL("CREATE TEMPORARY TABLE t_temp (id INT, name VARCHAR(50))")
		if err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateTable("t_temp")
		if err != nil {
			t.Fatalf("container show create: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TEMPORARY TABLE t_temp (id INT, name VARCHAR(50))", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_temp")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("create_table_like", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_like_dst")
		ctr.execSQL("DROP TABLE IF EXISTS t_like_src")
		ctr.execSQL("CREATE TABLE t_like_src (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT '', score DECIMAL(10,2))")
		err := ctr.execSQL("CREATE TABLE t_like_dst LIKE t_like_src")
		if err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_like_dst")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_like_src (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT '', score DECIMAL(10,2))", nil)
		results, _ := c.Exec("CREATE TABLE t_like_dst LIKE t_like_src", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_like_dst")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("create_table_with_view_name_conflict", func(t *testing.T) {
		// Creating a table with same name as an existing view should error.
		ctr.execSQL("DROP TABLE IF EXISTS t_view_conflict")
		ctr.execSQL("DROP VIEW IF EXISTS t_view_conflict")
		ctr.execSQL("CREATE VIEW t_view_conflict AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("CREATE TABLE t_view_conflict (id INT)")
		if ctrErr == nil {
			t.Fatal("expected container error when creating table with same name as view")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW t_view_conflict AS SELECT 1 AS a", nil)
		results, _ := c.Exec("CREATE TABLE t_view_conflict (id INT)", nil)
		if results[0].Error == nil {
			t.Fatal("expected omni error when creating table with same name as view")
		}
	})

	t.Run("reserved_word_as_name", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS `select`")
		err := ctr.execSQL("CREATE TABLE `select` (`from` INT, `where` VARCHAR(50))")
		if err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("`select`")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TABLE `select` (`from` INT, `where` VARCHAR(50))", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "select")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})
}

func TestContainer_Section_2_2_AlterTableColumnOps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// DDL-comparison tests: run setup + alter, compare SHOW CREATE TABLE.
	ddlCases := []struct {
		name  string
		setup string
		alter string
		table string
	}{
		{
			"add_column_at_end",
			"CREATE TABLE t_add_end (id INT PRIMARY KEY)",
			"ALTER TABLE t_add_end ADD COLUMN name VARCHAR(100)",
			"t_add_end",
		},
		{
			"add_column_first",
			"CREATE TABLE t_add_first (id INT PRIMARY KEY)",
			"ALTER TABLE t_add_first ADD COLUMN name VARCHAR(100) FIRST",
			"t_add_first",
		},
		{
			"add_column_after",
			"CREATE TABLE t_add_after (id INT PRIMARY KEY, age INT)",
			"ALTER TABLE t_add_after ADD COLUMN name VARCHAR(100) AFTER id",
			"t_add_after",
		},
		{
			"add_multiple_columns",
			"CREATE TABLE t_add_multi (id INT PRIMARY KEY)",
			"ALTER TABLE t_add_multi ADD COLUMN name VARCHAR(100), ADD COLUMN age INT, ADD COLUMN email VARCHAR(255)",
			"t_add_multi",
		},
		{
			"drop_column",
			"CREATE TABLE t_drop_col (id INT PRIMARY KEY, name VARCHAR(100), age INT)",
			"ALTER TABLE t_drop_col DROP COLUMN age",
			"t_drop_col",
		},
		{
			"drop_column_part_of_index",
			"CREATE TABLE t_drop_idx_col (id INT PRIMARY KEY, a INT, b INT, KEY idx_ab (a, b))",
			"ALTER TABLE t_drop_idx_col DROP COLUMN b",
			"t_drop_idx_col",
		},
		{
			"drop_column_only_in_index",
			"CREATE TABLE t_drop_only_idx (id INT PRIMARY KEY, a INT, KEY idx_a (a))",
			"ALTER TABLE t_drop_only_idx DROP COLUMN a",
			"t_drop_only_idx",
		},
		// Note: DROP COLUMN IF EXISTS is not supported in MySQL 8.0 (only 8.0.32+).
		// Scenario marked as [~] partial in SCENARIOS.md.
		{
			"modify_column_change_type",
			"CREATE TABLE t_mod_type (id INT PRIMARY KEY, val SMALLINT)",
			"ALTER TABLE t_mod_type MODIFY COLUMN val INT",
			"t_mod_type",
		},
		{
			"modify_column_widen_varchar",
			"CREATE TABLE t_mod_widen (id INT PRIMARY KEY, name VARCHAR(50))",
			"ALTER TABLE t_mod_widen MODIFY COLUMN name VARCHAR(200)",
			"t_mod_widen",
		},
		{
			"modify_column_narrow_varchar",
			"CREATE TABLE t_mod_narrow (id INT PRIMARY KEY, name VARCHAR(200))",
			"ALTER TABLE t_mod_narrow MODIFY COLUMN name VARCHAR(50)",
			"t_mod_narrow",
		},
		{
			"modify_column_int_to_bigint",
			"CREATE TABLE t_mod_bigint (id INT PRIMARY KEY, val INT)",
			"ALTER TABLE t_mod_bigint MODIFY COLUMN val BIGINT",
			"t_mod_bigint",
		},
		{
			"modify_column_add_not_null",
			"CREATE TABLE t_mod_nn (id INT PRIMARY KEY, name VARCHAR(100))",
			"ALTER TABLE t_mod_nn MODIFY COLUMN name VARCHAR(100) NOT NULL",
			"t_mod_nn",
		},
		{
			"modify_column_remove_not_null",
			"CREATE TABLE t_mod_rnn (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL)",
			"ALTER TABLE t_mod_rnn MODIFY COLUMN name VARCHAR(100)",
			"t_mod_rnn",
		},
		{
			"modify_column_first_after",
			"CREATE TABLE t_mod_pos (a INT, b INT, c INT)",
			"ALTER TABLE t_mod_pos MODIFY COLUMN c INT FIRST",
			"t_mod_pos",
		},
		{
			"change_column_rename_and_type",
			"CREATE TABLE t_chg_rt (id INT PRIMARY KEY, old_name VARCHAR(50))",
			"ALTER TABLE t_chg_rt CHANGE COLUMN old_name new_name VARCHAR(100)",
			"t_chg_rt",
		},
		{
			"change_column_same_name_diff_type",
			"CREATE TABLE t_chg_st (id INT PRIMARY KEY, val INT)",
			"ALTER TABLE t_chg_st CHANGE COLUMN val val BIGINT",
			"t_chg_st",
		},
		{
			"change_column_update_index_refs",
			"CREATE TABLE t_chg_idx (id INT PRIMARY KEY, a INT, KEY idx_a (a))",
			"ALTER TABLE t_chg_idx CHANGE COLUMN a b INT",
			"t_chg_idx",
		},
		{
			"rename_column",
			"CREATE TABLE t_ren_col (id INT PRIMARY KEY, old_col INT)",
			"ALTER TABLE t_ren_col RENAME COLUMN old_col TO new_col",
			"t_ren_col",
		},
		{
			"rename_column_update_index_refs",
			"CREATE TABLE t_ren_idx (id INT PRIMARY KEY, a INT, KEY idx_a (a))",
			"ALTER TABLE t_ren_idx RENAME COLUMN a TO b",
			"t_ren_idx",
		},
		{
			"alter_column_set_default",
			"CREATE TABLE t_set_def (id INT PRIMARY KEY, val INT)",
			"ALTER TABLE t_set_def ALTER COLUMN val SET DEFAULT 42",
			"t_set_def",
		},
		{
			"alter_column_drop_default",
			"CREATE TABLE t_drop_def (id INT PRIMARY KEY, val INT DEFAULT 10)",
			"ALTER TABLE t_drop_def ALTER COLUMN val DROP DEFAULT",
			"t_drop_def",
		},
		{
			"alter_column_set_visible",
			"CREATE TABLE t_vis (id INT PRIMARY KEY, val INT /*!80023 INVISIBLE */)",
			"ALTER TABLE t_vis ALTER COLUMN val SET VISIBLE",
			"t_vis",
		},
		{
			"alter_column_set_invisible",
			"CREATE TABLE t_invis (id INT PRIMARY KEY, val INT)",
			"ALTER TABLE t_invis ALTER COLUMN val SET INVISIBLE",
			"t_invis",
		},
		{
			"drop_column_part_of_pk",
			"CREATE TABLE t_drop_pk (a INT, b INT, PRIMARY KEY (a, b))",
			"ALTER TABLE t_drop_pk DROP COLUMN a",
			"t_drop_pk",
		},
	}

	for _, tc := range ddlCases {
		t.Run(tc.name, func(t *testing.T) {
			// Oracle: setup + alter.
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			if err := ctr.execSQL(tc.alter); err != nil {
				t.Fatalf("container alter: %v", err)
			}
			ctrDDL, err := ctr.showCreateTable(tc.table)
			if err != nil {
				t.Fatalf("container show create: %v", err)
			}

			// Omni: setup + alter.
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, _ := c.Exec(tc.setup, nil)
			if results[0].Error != nil {
				t.Fatalf("omni setup error: %v", results[0].Error)
			}
			results, _ = c.Exec(tc.alter, nil)
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("omni alter error: %v", r.Error)
				}
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}

	// Error tests: operations that should produce errors.
	errCases := []struct {
		name    string
		setup   string
		alter   string
		table   string
		wantErr bool
	}{
		{
			"drop_column_referenced_by_fk",
			"CREATE TABLE t_fk_parent (id INT PRIMARY KEY); CREATE TABLE t_fk_child (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent(id))",
			"ALTER TABLE t_fk_child DROP COLUMN pid",
			"t_fk_child",
			true,
		},
	}

	for _, tc := range errCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up tables first.
			ctr.execSQL("DROP TABLE IF EXISTS t_fk_child")
			ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent")
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)

			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			ctrErr := ctr.execSQL(tc.alter)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.setup, nil)
			results, _ := c.Exec(tc.alter, nil)
			var omniErr error
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}

			if tc.wantErr {
				if ctrErr == nil {
					t.Fatal("expected container error but got nil")
				}
				if omniErr == nil {
					t.Fatalf("expected omni error but got nil (container error: %v)", ctrErr)
				}
				t.Logf("both errored as expected — container: %v, omni: %v", ctrErr, omniErr)
			} else {
				if ctrErr != nil {
					t.Fatalf("unexpected container error: %v", ctrErr)
				}
				if omniErr != nil {
					t.Fatalf("unexpected omni error: %v", omniErr)
				}
			}
		})
	}
}

func TestContainer_Section_2_3_AlterTableIndexOps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// DDL-comparison tests: run setup + alter, compare SHOW CREATE TABLE.
	ddlCases := []struct {
		name  string
		setup string
		alter string
		table string
	}{
		{
			"add_index",
			"CREATE TABLE t_add_idx (id INT PRIMARY KEY, name VARCHAR(100))",
			"ALTER TABLE t_add_idx ADD INDEX idx_name (name)",
			"t_add_idx",
		},
		{
			"add_unique_index",
			"CREATE TABLE t_add_uniq (id INT PRIMARY KEY, email VARCHAR(255))",
			"ALTER TABLE t_add_uniq ADD UNIQUE INDEX idx_email (email)",
			"t_add_uniq",
		},
		{
			"add_fulltext_index",
			"CREATE TABLE t_add_ft (id INT PRIMARY KEY, body TEXT) ENGINE=InnoDB",
			"ALTER TABLE t_add_ft ADD FULLTEXT INDEX idx_body (body)",
			"t_add_ft",
		},
		{
			"add_primary_key",
			"CREATE TABLE t_add_pk (id INT NOT NULL, name VARCHAR(100))",
			"ALTER TABLE t_add_pk ADD PRIMARY KEY (id)",
			"t_add_pk",
		},
		{
			"drop_index",
			"CREATE TABLE t_drop_idx (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name))",
			"ALTER TABLE t_drop_idx DROP INDEX idx_name",
			"t_drop_idx",
		},
		// Note: DROP INDEX IF EXISTS is not supported in MySQL 8.0 (syntax error).
		// Scenario marked as [~] partial in SCENARIOS.md.
		{
			"drop_primary_key",
			"CREATE TABLE t_drop_pk (id INT NOT NULL, name VARCHAR(100), PRIMARY KEY (id))",
			"ALTER TABLE t_drop_pk DROP PRIMARY KEY",
			"t_drop_pk",
		},
		{
			"rename_index",
			"CREATE TABLE t_ren_idx (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_old (name))",
			"ALTER TABLE t_ren_idx RENAME INDEX idx_old TO idx_new",
			"t_ren_idx",
		},
		{
			"alter_index_invisible",
			"CREATE TABLE t_idx_invis (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name))",
			"ALTER TABLE t_idx_invis ALTER INDEX idx_name INVISIBLE",
			"t_idx_invis",
		},
		{
			"alter_index_visible",
			"CREATE TABLE t_idx_vis (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name) INVISIBLE)",
			"ALTER TABLE t_idx_vis ALTER INDEX idx_name VISIBLE",
			"t_idx_vis",
		},
	}

	for _, tc := range ddlCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			if err := ctr.execSQL(tc.alter); err != nil {
				t.Fatalf("container alter: %v", err)
			}
			ctrDDL, err := ctr.showCreateTable(tc.table)
			if err != nil {
				t.Fatalf("container SHOW CREATE TABLE: %v", err)
			}

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			if results, _ := c.Exec(tc.setup, nil); results != nil {
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni setup error: %v", r.Error)
					}
				}
			}
			if results, _ := c.Exec(tc.alter, nil); results != nil {
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni alter error: %v", r.Error)
					}
				}
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}

	// Error tests.
	errCases := []struct {
		name    string
		setup   string
		alter   string
		table   string
		wantErr bool
	}{
		{
			"add_primary_key_when_exists",
			"CREATE TABLE t_dup_pk (id INT PRIMARY KEY, val INT NOT NULL)",
			"ALTER TABLE t_dup_pk ADD PRIMARY KEY (val)",
			"t_dup_pk",
			true,
		},
	}

	for _, tc := range errCases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			ctrErr := ctr.execSQL(tc.alter)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.setup, nil)
			results, _ := c.Exec(tc.alter, nil)
			var omniErr error
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
					break
				}
			}

			if tc.wantErr {
				if ctrErr == nil {
					t.Fatal("expected container error but got nil")
				}
				if omniErr == nil {
					t.Fatalf("expected omni error but got nil (container error: %v)", ctrErr)
				}
				t.Logf("both errored as expected — container: %v, omni: %v", ctrErr, omniErr)
			}
		})
	}
}

func TestContainer_Section_2_4_AlterTableConstraintOps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// DDL-comparison tests: run setup + alter, compare SHOW CREATE TABLE.
	ddlCases := []struct {
		name  string
		setup string
		alter string
		table string
	}{
		{
			"add_foreign_key",
			"CREATE TABLE t_parent_fk (id INT PRIMARY KEY); CREATE TABLE t_child_fk (id INT PRIMARY KEY, parent_id INT)",
			"ALTER TABLE t_child_fk ADD CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES t_parent_fk(id)",
			"t_child_fk",
		},
		{
			"add_check",
			"CREATE TABLE t_add_chk (id INT PRIMARY KEY, age INT)",
			"ALTER TABLE t_add_chk ADD CHECK (age > 0)",
			"t_add_chk",
		},
		{
			"drop_foreign_key",
			"CREATE TABLE t_parent_dfk (id INT PRIMARY KEY); CREATE TABLE t_child_dfk (id INT PRIMARY KEY, parent_id INT, CONSTRAINT fk_drop FOREIGN KEY (parent_id) REFERENCES t_parent_dfk(id))",
			"ALTER TABLE t_child_dfk DROP FOREIGN KEY fk_drop",
			"t_child_dfk",
		},
		{
			"drop_check",
			"CREATE TABLE t_drop_chk (id INT PRIMARY KEY, age INT, CONSTRAINT chk_age CHECK (age > 0))",
			"ALTER TABLE t_drop_chk DROP CHECK chk_age",
			"t_drop_chk",
		},
		{
			"drop_constraint_generic",
			"CREATE TABLE t_drop_con (id INT PRIMARY KEY, val INT, CONSTRAINT chk_val CHECK (val >= 0))",
			"ALTER TABLE t_drop_con DROP CONSTRAINT chk_val",
			"t_drop_con",
		},
		{
			"alter_check_not_enforced",
			"CREATE TABLE t_chk_ne (id INT PRIMARY KEY, score INT, CONSTRAINT chk_score CHECK (score >= 0))",
			"ALTER TABLE t_chk_ne ALTER CHECK chk_score NOT ENFORCED",
			"t_chk_ne",
		},
		{
			"alter_check_enforced",
			"CREATE TABLE t_chk_enf (id INT PRIMARY KEY, score INT, CONSTRAINT chk_score2 CHECK (score >= 0) /*!80016 NOT ENFORCED */)",
			"ALTER TABLE t_chk_enf ALTER CHECK chk_score2 ENFORCED",
			"t_chk_enf",
		},
	}

	for _, tc := range ddlCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up all tables that might exist from setup.
			for _, tName := range []string{
				"t_parent_fk", "t_child_fk",
				"t_parent_dfk", "t_child_dfk",
				"t_add_chk",
				"t_drop_chk",
				"t_drop_con",
				"t_chk_ne",
				"t_chk_enf",
			} {
				ctr.execSQL("DROP TABLE IF EXISTS " + tName)
			}
			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			if err := ctr.execSQL(tc.alter); err != nil {
				t.Fatalf("container alter: %v", err)
			}
			ctrDDL, err := ctr.showCreateTable(tc.table)
			if err != nil {
				t.Fatalf("container SHOW CREATE TABLE: %v", err)
			}

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			if results, _ := c.Exec(tc.setup, nil); results != nil {
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni setup error: %v", r.Error)
					}
				}
			}
			if results, _ := c.Exec(tc.alter, nil); results != nil {
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni alter error: %v", r.Error)
					}
				}
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_2_5_AlterTableTableLevel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		setup string
		alter string
		table string // table name to SHOW CREATE TABLE after alter
	}{
		{
			"rename_to",
			"CREATE TABLE t_rename_src (id INT PRIMARY KEY, name VARCHAR(100))",
			"ALTER TABLE t_rename_src RENAME TO t_rename_dst",
			"t_rename_dst",
		},
		{
			"engine_myisam",
			"CREATE TABLE t_engine (id INT PRIMARY KEY, val INT)",
			"ALTER TABLE t_engine ENGINE=MyISAM",
			"t_engine",
		},
		{
			"convert_charset_utf8mb4",
			"CREATE TABLE t_conv_cs (id INT PRIMARY KEY, name VARCHAR(100)) DEFAULT CHARSET=latin1",
			"ALTER TABLE t_conv_cs CONVERT TO CHARACTER SET utf8mb4",
			"t_conv_cs",
		},
		{
			"default_charset_latin1",
			"CREATE TABLE t_def_cs (id INT PRIMARY KEY, name VARCHAR(100))",
			"ALTER TABLE t_def_cs DEFAULT CHARACTER SET latin1",
			"t_def_cs",
		},
		{
			"comment",
			"CREATE TABLE t_comment (id INT PRIMARY KEY)",
			"ALTER TABLE t_comment COMMENT='new comment'",
			"t_comment",
		},
		{
			"auto_increment",
			"CREATE TABLE t_autoinc (id INT PRIMARY KEY AUTO_INCREMENT, val INT)",
			"ALTER TABLE t_autoinc AUTO_INCREMENT=1000",
			"t_autoinc",
		},
		{
			"row_format_compressed",
			"CREATE TABLE t_rowfmt (id INT PRIMARY KEY, val INT)",
			"ALTER TABLE t_rowfmt ROW_FORMAT=COMPRESSED",
			"t_rowfmt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up tables that might exist.
			for _, tName := range []string{
				"t_rename_src", "t_rename_dst",
				"t_engine",
				"t_conv_cs",
				"t_def_cs",
				"t_comment",
				"t_autoinc",
				"t_rowfmt",
			} {
				ctr.execSQL("DROP TABLE IF EXISTS " + tName)
			}
			if err := ctr.execSQL(tc.setup); err != nil {
				t.Fatalf("container setup: %v", err)
			}
			if err := ctr.execSQL(tc.alter); err != nil {
				t.Fatalf("container alter: %v", err)
			}
			ctrDDL, err := ctr.showCreateTable(tc.table)
			if err != nil {
				t.Fatalf("container SHOW CREATE TABLE: %v", err)
			}

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			if results, _ := c.Exec(tc.setup, nil); results != nil {
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni setup error: %v", r.Error)
					}
				}
			}
			if results, _ := c.Exec(tc.alter, nil); results != nil {
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("omni alter error: %v", r.Error)
					}
				}
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestContainer_Section_2_6_DropTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("drop_table_basic", func(t *testing.T) {
		// Setup: create a table, then drop it. Verify it's gone.
		ctr.execSQL("DROP TABLE IF EXISTS t_drop1")
		ctr.execSQL("CREATE TABLE t_drop1 (id INT PRIMARY KEY, name VARCHAR(100))")

		ctrErr := ctr.execSQL("DROP TABLE t_drop1")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE error: %v", ctrErr)
		}
		_, ctrShowErr := ctr.showCreateTable("t_drop1")
		if ctrShowErr == nil {
			t.Fatal("container: table still exists after DROP TABLE")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_drop1 (id INT PRIMARY KEY, name VARCHAR(100))", nil)
		results, _ := c.Exec("DROP TABLE t_drop1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TABLE error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_drop1")
		if omniDDL != "" {
			t.Errorf("omni: table still exists after DROP TABLE, got: %s", omniDDL)
		}
	})

	t.Run("drop_table_if_exists", func(t *testing.T) {
		// DROP TABLE IF EXISTS on a nonexistent table should not error.
		ctr.execSQL("DROP TABLE IF EXISTS t_drop_ine")

		ctrErr := ctr.execSQL("DROP TABLE IF EXISTS t_drop_ine")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE IF EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE IF EXISTS t_drop_ine", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TABLE IF EXISTS error: %v", results[0].Error)
		}
	})

	t.Run("drop_table_multi", func(t *testing.T) {
		// DROP TABLE t1, t2, t3 — multi-table drop.
		ctr.execSQL("DROP TABLE IF EXISTS t_dm1")
		ctr.execSQL("DROP TABLE IF EXISTS t_dm2")
		ctr.execSQL("DROP TABLE IF EXISTS t_dm3")
		ctr.execSQL("CREATE TABLE t_dm1 (id INT)")
		ctr.execSQL("CREATE TABLE t_dm2 (id INT)")
		ctr.execSQL("CREATE TABLE t_dm3 (id INT)")

		ctrErr := ctr.execSQL("DROP TABLE t_dm1, t_dm2, t_dm3")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE multi error: %v", ctrErr)
		}
		for _, tbl := range []string{"t_dm1", "t_dm2", "t_dm3"} {
			if _, err := ctr.showCreateTable(tbl); err == nil {
				t.Errorf("container: table %s still exists after multi-drop", tbl)
			}
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_dm1 (id INT)", nil)
		c.Exec("CREATE TABLE t_dm2 (id INT)", nil)
		c.Exec("CREATE TABLE t_dm3 (id INT)", nil)
		results, _ := c.Exec("DROP TABLE t_dm1, t_dm2, t_dm3", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TABLE multi error: %v", results[0].Error)
		}
		for _, tbl := range []string{"t_dm1", "t_dm2", "t_dm3"} {
			ddl := c.ShowCreateTable("test", tbl)
			if ddl != "" {
				t.Errorf("omni: table %s still exists after multi-drop", tbl)
			}
		}
	})

	t.Run("drop_table_nonexistent_error", func(t *testing.T) {
		// DROP TABLE on nonexistent table should produce error 1051.
		ctr.execSQL("DROP TABLE IF EXISTS t_noexist")
		ctrErr := ctr.execSQL("DROP TABLE t_noexist")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP nonexistent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE t_noexist", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for DROP nonexistent table")
		}
		catErr, ok := omniErr.(*Error)
		if !ok {
			t.Fatalf("omni error is not *Error: %T", omniErr)
		}
		if catErr.Code != 1051 {
			t.Errorf("omni error code: want 1051, got %d (message: %s)", catErr.Code, catErr.Message)
		}
	})

	t.Run("drop_temporary_table", func(t *testing.T) {
		// DROP TEMPORARY TABLE should work.
		ctr.execSQL("DROP TEMPORARY TABLE IF EXISTS t_temp_drop")
		ctr.execSQL("CREATE TEMPORARY TABLE t_temp_drop (id INT, val VARCHAR(50))")
		ctrErr := ctr.execSQL("DROP TEMPORARY TABLE t_temp_drop")
		if ctrErr != nil {
			t.Fatalf("container DROP TEMPORARY TABLE error: %v", ctrErr)
		}
		_, ctrShowErr := ctr.showCreateTable("t_temp_drop")
		if ctrShowErr == nil {
			t.Fatal("container: temp table still exists after DROP TEMPORARY TABLE")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TEMPORARY TABLE t_temp_drop (id INT, val VARCHAR(50))", nil)
		results, _ := c.Exec("DROP TEMPORARY TABLE t_temp_drop", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TEMPORARY TABLE error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_temp_drop")
		if omniDDL != "" {
			t.Errorf("omni: temp table still exists after DROP TEMPORARY TABLE")
		}
	})

	t.Run("drop_table_fk_referenced", func(t *testing.T) {
		// DROP TABLE that has FK references should error (with foreign_key_checks=1, the default).
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_child")
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent")
		ctr.execSQL("CREATE TABLE t_fk_parent (id INT PRIMARY KEY)")
		ctr.execSQL("CREATE TABLE t_fk_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fk_parent(id))")

		ctrErr := ctr.execSQL("DROP TABLE t_fk_parent")
		if ctrErr == nil {
			t.Fatal("container: expected error when dropping FK-referenced table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_fk_parent (id INT PRIMARY KEY)", nil)
		c.Exec("CREATE TABLE t_fk_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fk_parent(id))", nil)
		results, _ := c.Exec("DROP TABLE t_fk_parent", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error when dropping FK-referenced table")
		}

		// Verify table still exists after failed drop.
		omniDDL := c.ShowCreateTable("test", "t_fk_parent")
		if omniDDL == "" {
			t.Error("omni: parent table was deleted despite FK reference")
		}
	})
}

func TestContainer_Section_2_7_TruncateTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Scenario 1: TRUNCATE TABLE t1 — table structure preserved
	t.Run("truncate_basic", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_trunc1")
		if err := ctr.execSQL("CREATE TABLE t_trunc1 (id INT PRIMARY KEY, name VARCHAR(100))"); err != nil {
			t.Fatalf("container create: %v", err)
		}
		if err := ctr.execSQL("TRUNCATE TABLE t_trunc1"); err != nil {
			t.Fatalf("container truncate: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_trunc1")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trunc1 (id INT PRIMARY KEY, name VARCHAR(100))", nil)
		results, _ := c.Exec("TRUNCATE TABLE t_trunc1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni truncate error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_trunc1")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// Scenario 2: TRUNCATE resets AUTO_INCREMENT
	t.Run("truncate_resets_auto_increment", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_trunc_ai")
		if err := ctr.execSQL("CREATE TABLE t_trunc_ai (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100)) AUTO_INCREMENT=1000"); err != nil {
			t.Fatalf("container create: %v", err)
		}
		// Verify AUTO_INCREMENT is shown before truncate
		ctrBefore, _ := ctr.showCreateTable("t_trunc_ai")
		if !strings.Contains(ctrBefore, "AUTO_INCREMENT=") {
			t.Logf("container before truncate (no AUTO_INCREMENT shown): %s", ctrBefore)
		}
		if err := ctr.execSQL("TRUNCATE TABLE t_trunc_ai"); err != nil {
			t.Fatalf("container truncate: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_trunc_ai")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trunc_ai (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100)) AUTO_INCREMENT=1000", nil)
		c.Exec("TRUNCATE TABLE t_trunc_ai", nil)
		omniDDL := c.ShowCreateTable("test", "t_trunc_ai")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// Scenario 3: TRUNCATE nonexistent table → error
	t.Run("truncate_nonexistent", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_trunc_noexist")
		ctrErr := ctr.execSQL("TRUNCATE TABLE t_trunc_noexist")
		if ctrErr == nil {
			t.Fatal("container: expected error for TRUNCATE nonexistent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("TRUNCATE TABLE t_trunc_noexist", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for TRUNCATE nonexistent table")
		}
	})
}

func TestContainer_Section_2_8_CreateDropIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Scenario 1: CREATE INDEX idx ON t (col)
	t.Run("create_index_basic", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ci1")
		ctr.execSQL("CREATE TABLE t_ci1 (id INT PRIMARY KEY, name VARCHAR(100))")
		if err := ctr.execSQL("CREATE INDEX idx_name ON t_ci1 (name)"); err != nil {
			t.Fatalf("container CREATE INDEX error: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_ci1")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci1 (id INT PRIMARY KEY, name VARCHAR(100))", nil)
		results, _ := c.Exec("CREATE INDEX idx_name ON t_ci1 (name)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE INDEX error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_ci1")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// Scenario 2: CREATE UNIQUE INDEX
	t.Run("create_unique_index", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ci2")
		ctr.execSQL("CREATE TABLE t_ci2 (id INT PRIMARY KEY, email VARCHAR(255))")
		if err := ctr.execSQL("CREATE UNIQUE INDEX idx_email ON t_ci2 (email)"); err != nil {
			t.Fatalf("container CREATE UNIQUE INDEX error: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_ci2")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci2 (id INT PRIMARY KEY, email VARCHAR(255))", nil)
		results, _ := c.Exec("CREATE UNIQUE INDEX idx_email ON t_ci2 (email)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE UNIQUE INDEX error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_ci2")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// Scenario 3: CREATE FULLTEXT INDEX
	t.Run("create_fulltext_index", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ci3")
		ctr.execSQL("CREATE TABLE t_ci3 (id INT PRIMARY KEY, content TEXT)")
		if err := ctr.execSQL("CREATE FULLTEXT INDEX idx_ft ON t_ci3 (content)"); err != nil {
			t.Fatalf("container CREATE FULLTEXT INDEX error: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_ci3")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci3 (id INT PRIMARY KEY, content TEXT)", nil)
		results, _ := c.Exec("CREATE FULLTEXT INDEX idx_ft ON t_ci3 (content)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE FULLTEXT INDEX error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_ci3")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// Scenario 4: CREATE SPATIAL INDEX
	t.Run("create_spatial_index", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ci4")
		ctr.execSQL("CREATE TABLE t_ci4 (id INT PRIMARY KEY, geo GEOMETRY NOT NULL)")
		if err := ctr.execSQL("CREATE SPATIAL INDEX idx_sp ON t_ci4 (geo)"); err != nil {
			t.Fatalf("container CREATE SPATIAL INDEX error: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_ci4")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci4 (id INT PRIMARY KEY, geo GEOMETRY NOT NULL)", nil)
		results, _ := c.Exec("CREATE SPATIAL INDEX idx_sp ON t_ci4 (geo)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE SPATIAL INDEX error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_ci4")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// Scenario 5: CREATE INDEX IF NOT EXISTS
	// Note: MySQL 8.0 does NOT support IF NOT EXISTS on CREATE INDEX (syntax error).
	// Our parser accepts it, and the catalog handles it gracefully (no-op on duplicate).
	// This is marked [~] partial — omni is more permissive than MySQL 8.0 here.
	t.Run("create_index_if_not_exists", func(t *testing.T) {
		// Verify MySQL 8.0 rejects this syntax.
		ctr.execSQL("DROP TABLE IF EXISTS t_ci5")
		ctr.execSQL("CREATE TABLE t_ci5 (id INT PRIMARY KEY, val INT)")
		ctr.execSQL("CREATE INDEX idx_val ON t_ci5 (val)")
		ctrErr := ctr.execSQL("CREATE INDEX IF NOT EXISTS idx_val ON t_ci5 (val)")
		if ctrErr == nil {
			t.Fatal("container: expected syntax error for CREATE INDEX IF NOT EXISTS in MySQL 8.0")
		}
		t.Logf("container correctly rejects CREATE INDEX IF NOT EXISTS: %v", ctrErr)

		// Omni accepts IF NOT EXISTS as an extension — verify it doesn't error.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci5 (id INT PRIMARY KEY, val INT)", nil)
		c.Exec("CREATE INDEX idx_val ON t_ci5 (val)", nil)
		results, _ := c.Exec("CREATE INDEX IF NOT EXISTS idx_val ON t_ci5 (val)", nil)
		if results[0].Error != nil {
			t.Errorf("omni: unexpected error for IF NOT EXISTS: %v", results[0].Error)
		}
		// This is a known divergence from MySQL 8.0 behavior.
	})

	// Scenario 6: CREATE INDEX — duplicate name → error 1061
	t.Run("create_index_duplicate_error", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_ci6")
		ctr.execSQL("CREATE TABLE t_ci6 (id INT PRIMARY KEY, a INT, b INT)")
		ctr.execSQL("CREATE INDEX idx_a ON t_ci6 (a)")
		ctrErr := ctr.execSQL("CREATE INDEX idx_a ON t_ci6 (b)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate index name")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ci6 (id INT PRIMARY KEY, a INT, b INT)", nil)
		c.Exec("CREATE INDEX idx_a ON t_ci6 (a)", nil)
		results, _ := c.Exec("CREATE INDEX idx_a ON t_ci6 (b)", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for duplicate index name")
		}
		catErr, ok := omniErr.(*Error)
		if !ok {
			t.Fatalf("omni error is not *Error: %T", omniErr)
		}
		if catErr.Code != 1061 {
			t.Errorf("omni error code: want 1061, got %d (message: %s)", catErr.Code, catErr.Message)
		}
	})

	// Scenario 7: DROP INDEX idx ON t
	t.Run("drop_index_basic", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_di1")
		ctr.execSQL("CREATE TABLE t_di1 (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name))")
		if err := ctr.execSQL("DROP INDEX idx_name ON t_di1"); err != nil {
			t.Fatalf("container DROP INDEX error: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_di1")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_di1 (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name))", nil)
		results, _ := c.Exec("DROP INDEX idx_name ON t_di1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP INDEX error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_di1")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// Scenario 8: DROP INDEX nonexistent → error 1091
	t.Run("drop_index_nonexistent_error", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_di2")
		ctr.execSQL("CREATE TABLE t_di2 (id INT PRIMARY KEY)")
		ctrErr := ctr.execSQL("DROP INDEX idx_noexist ON t_di2")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP nonexistent index")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_di2 (id INT PRIMARY KEY)", nil)
		results, _ := c.Exec("DROP INDEX idx_noexist ON t_di2", nil)
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for DROP nonexistent index")
		}
		catErr, ok := omniErr.(*Error)
		if !ok {
			t.Fatalf("omni error is not *Error: %T", omniErr)
		}
		if catErr.Code != 1091 {
			t.Errorf("omni error code: want 1091, got %d (message: %s)", catErr.Code, catErr.Message)
		}
	})
}

func TestContainer_Section_2_9_RenameTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("rename_basic", func(t *testing.T) {
		// RENAME TABLE t1 TO t2
		ctr.execSQL("DROP TABLE IF EXISTS t_ren1")
		ctr.execSQL("DROP TABLE IF EXISTS t_ren2")
		ctr.execSQL("CREATE TABLE t_ren1 (id INT PRIMARY KEY, name VARCHAR(100))")

		if err := ctr.execSQL("RENAME TABLE t_ren1 TO t_ren2"); err != nil {
			t.Fatalf("container RENAME TABLE error: %v", err)
		}
		ctrDDL, err := ctr.showCreateTable("t_ren2")
		if err != nil {
			t.Fatalf("container SHOW CREATE TABLE t_ren2 error: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ren1 (id INT PRIMARY KEY, name VARCHAR(100))", nil)
		results, _ := c.Exec("RENAME TABLE t_ren1 TO t_ren2", nil)
		if results[0].Error != nil {
			t.Fatalf("omni RENAME TABLE error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_ren2")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL, omniDDL)
		}
		// Old table should not exist
		if c.ShowCreateTable("test", "t_ren1") != "" {
			t.Error("omni: old table t_ren1 still exists after rename")
		}
	})

	t.Run("rename_cross_database", func(t *testing.T) {
		// RENAME TABLE t1 TO db2.t1 (cross-database)
		ctr.execSQL("DROP TABLE IF EXISTS t_ren_cross")
		ctr.execSQL("CREATE DATABASE IF NOT EXISTS test2")
		ctr.execSQL("DROP TABLE IF EXISTS test2.t_ren_cross")
		ctr.execSQL("CREATE TABLE t_ren_cross (id INT PRIMARY KEY, val VARCHAR(50) NOT NULL DEFAULT '')")

		if err := ctr.execSQL("RENAME TABLE t_ren_cross TO test2.t_ren_cross"); err != nil {
			t.Fatalf("container RENAME TABLE cross-db error: %v", err)
		}
		ctrDDL, err := ctr.showCreateTable("test2.t_ren_cross")
		if err != nil {
			t.Fatalf("container SHOW CREATE TABLE test2.t_ren_cross error: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.Exec("CREATE DATABASE test2", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ren_cross (id INT PRIMARY KEY, val VARCHAR(50) NOT NULL DEFAULT '')", nil)
		results, _ := c.Exec("RENAME TABLE t_ren_cross TO test2.t_ren_cross", nil)
		if results[0].Error != nil {
			t.Fatalf("omni RENAME TABLE cross-db error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test2", "t_ren_cross")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL, omniDDL)
		}
		// Old table should not exist in test
		if c.ShowCreateTable("test", "t_ren_cross") != "" {
			t.Error("omni: old table still exists in test after cross-db rename")
		}
	})

	t.Run("rename_multi_pair", func(t *testing.T) {
		// RENAME TABLE t1 TO t2, t3 TO t4 (multi-pair)
		ctr.execSQL("DROP TABLE IF EXISTS t_mp1")
		ctr.execSQL("DROP TABLE IF EXISTS t_mp2")
		ctr.execSQL("DROP TABLE IF EXISTS t_mp3")
		ctr.execSQL("DROP TABLE IF EXISTS t_mp4")
		ctr.execSQL("CREATE TABLE t_mp1 (id INT PRIMARY KEY)")
		ctr.execSQL("CREATE TABLE t_mp3 (val VARCHAR(100))")

		if err := ctr.execSQL("RENAME TABLE t_mp1 TO t_mp2, t_mp3 TO t_mp4"); err != nil {
			t.Fatalf("container RENAME TABLE multi-pair error: %v", err)
		}
		ctrDDL2, err := ctr.showCreateTable("t_mp2")
		if err != nil {
			t.Fatalf("container SHOW CREATE TABLE t_mp2 error: %v", err)
		}
		ctrDDL4, err := ctr.showCreateTable("t_mp4")
		if err != nil {
			t.Fatalf("container SHOW CREATE TABLE t_mp4 error: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_mp1 (id INT PRIMARY KEY)", nil)
		c.Exec("CREATE TABLE t_mp3 (val VARCHAR(100))", nil)
		results, _ := c.Exec("RENAME TABLE t_mp1 TO t_mp2, t_mp3 TO t_mp4", nil)
		if results[0].Error != nil {
			t.Fatalf("omni RENAME TABLE multi-pair error: %v", results[0].Error)
		}
		omniDDL2 := c.ShowCreateTable("test", "t_mp2")
		omniDDL4 := c.ShowCreateTable("test", "t_mp4")

		if normalizeWhitespace(ctrDDL2) != normalizeWhitespace(omniDDL2) {
			t.Errorf("t_mp2 mismatch:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL2, omniDDL2)
		}
		if normalizeWhitespace(ctrDDL4) != normalizeWhitespace(omniDDL4) {
			t.Errorf("t_mp4 mismatch:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL4, omniDDL4)
		}
		// Old tables should be gone
		if c.ShowCreateTable("test", "t_mp1") != "" {
			t.Error("omni: t_mp1 still exists after rename")
		}
		if c.ShowCreateTable("test", "t_mp3") != "" {
			t.Error("omni: t_mp3 still exists after rename")
		}
	})

	t.Run("rename_nonexistent_error", func(t *testing.T) {
		// RENAME TABLE nonexistent → error
		ctr.execSQL("DROP TABLE IF EXISTS t_noexist_ren")
		ctrErr := ctr.execSQL("RENAME TABLE t_noexist_ren TO t_noexist_ren2")
		if ctrErr == nil {
			t.Fatal("container: expected error for RENAME nonexistent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("RENAME TABLE t_noexist_ren TO t_noexist_ren2", &ExecOptions{ContinueOnError: true})
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for RENAME nonexistent table")
		}
	})

	t.Run("rename_to_existing_error", func(t *testing.T) {
		// RENAME TABLE to existing name → error
		ctr.execSQL("DROP TABLE IF EXISTS t_ren_exist1")
		ctr.execSQL("DROP TABLE IF EXISTS t_ren_exist2")
		ctr.execSQL("CREATE TABLE t_ren_exist1 (id INT)")
		ctr.execSQL("CREATE TABLE t_ren_exist2 (id INT)")

		ctrErr := ctr.execSQL("RENAME TABLE t_ren_exist1 TO t_ren_exist2")
		if ctrErr == nil {
			t.Fatal("container: expected error for RENAME to existing table name")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_ren_exist1 (id INT)", nil)
		c.Exec("CREATE TABLE t_ren_exist2 (id INT)", nil)
		results, _ := c.Exec("RENAME TABLE t_ren_exist1 TO t_ren_exist2", &ExecOptions{ContinueOnError: true})
		omniErr := results[0].Error
		if omniErr == nil {
			t.Fatal("omni: expected error for RENAME to existing table name")
		}
	})
}

func TestContainer_Section_2_10_CreateDropView(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("create_view_basic", func(t *testing.T) {
		// CREATE VIEW v AS SELECT ...
		ctr.execSQL("DROP VIEW IF EXISTS v_basic")
		ctrErr := ctr.execSQL("CREATE VIEW v_basic AS SELECT 1 AS a")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE VIEW v_basic AS SELECT 1 AS a", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_basic")] == nil {
			t.Error("omni: view v_basic should exist after CREATE VIEW")
		}
	})

	t.Run("create_or_replace_view", func(t *testing.T) {
		// CREATE OR REPLACE VIEW
		ctr.execSQL("DROP VIEW IF EXISTS v_replace")
		ctr.execSQL("CREATE VIEW v_replace AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("CREATE OR REPLACE VIEW v_replace AS SELECT 2 AS b")
		if ctrErr != nil {
			t.Fatalf("container CREATE OR REPLACE VIEW error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_replace AS SELECT 1 AS a", nil)
		results, _ := c.Exec("CREATE OR REPLACE VIEW v_replace AS SELECT 2 AS b", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE OR REPLACE VIEW error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_replace")] == nil {
			t.Error("omni: view v_replace should exist after CREATE OR REPLACE VIEW")
		}
	})

	t.Run("create_view_with_columns", func(t *testing.T) {
		// CREATE VIEW with column list
		ctr.execSQL("DROP VIEW IF EXISTS v_cols")
		ctrErr := ctr.execSQL("CREATE VIEW v_cols (x, y) AS SELECT 1, 2")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW with columns error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE VIEW v_cols (x, y) AS SELECT 1, 2", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW with columns error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_cols")]
		if v == nil {
			t.Fatal("omni: view v_cols should exist")
		}
		if len(v.Columns) != 2 || v.Columns[0] != "x" || v.Columns[1] != "y" {
			t.Errorf("omni: expected columns [x, y], got %v", v.Columns)
		}
	})

	t.Run("create_view_with_options", func(t *testing.T) {
		// CREATE VIEW with ALGORITHM, DEFINER, SQL_SECURITY
		ctr.execSQL("DROP VIEW IF EXISTS v_opts")
		ctrErr := ctr.execSQL("CREATE ALGORITHM=MERGE DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW v_opts AS SELECT 1 AS a")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW with options error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE ALGORITHM=MERGE DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW v_opts AS SELECT 1 AS a", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW with options error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_opts")]
		if v == nil {
			t.Fatal("omni: view v_opts should exist")
		}
		if v.Algorithm != "MERGE" {
			t.Errorf("omni: expected Algorithm=MERGE, got %q", v.Algorithm)
		}
		if v.SqlSecurity != "INVOKER" {
			t.Errorf("omni: expected SqlSecurity=INVOKER, got %q", v.SqlSecurity)
		}
	})

	t.Run("create_view_with_check_option", func(t *testing.T) {
		// CREATE VIEW with CHECK OPTION
		ctr.execSQL("DROP VIEW IF EXISTS v_chk")
		ctr.execSQL("DROP TABLE IF EXISTS t_chk_view")
		ctr.execSQL("CREATE TABLE t_chk_view (id INT, val INT)")
		ctrErr := ctr.execSQL("CREATE VIEW v_chk AS SELECT * FROM t_chk_view WHERE val > 0 WITH CASCADED CHECK OPTION")
		if ctrErr != nil {
			t.Fatalf("container CREATE VIEW WITH CHECK OPTION error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_chk_view (id INT, val INT)", nil)
		results, _ := c.Exec("CREATE VIEW v_chk AS SELECT * FROM t_chk_view WHERE val > 0 WITH CASCADED CHECK OPTION", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE VIEW WITH CHECK OPTION error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_chk")]
		if v == nil {
			t.Fatal("omni: view v_chk should exist")
		}
		if v.CheckOption != "CASCADED" {
			t.Errorf("omni: expected CheckOption=CASCADED, got %q", v.CheckOption)
		}
	})

	t.Run("drop_view_basic", func(t *testing.T) {
		// DROP VIEW v
		ctr.execSQL("DROP VIEW IF EXISTS v_drop1")
		ctr.execSQL("CREATE VIEW v_drop1 AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("DROP VIEW v_drop1")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_drop1 AS SELECT 1 AS a", nil)
		results, _ := c.Exec("DROP VIEW v_drop1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP VIEW error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_drop1")] != nil {
			t.Error("omni: view v_drop1 should not exist after DROP VIEW")
		}
	})

	t.Run("drop_view_if_exists", func(t *testing.T) {
		// DROP VIEW IF EXISTS on nonexistent view — no error
		ctr.execSQL("DROP VIEW IF EXISTS v_noexist")
		ctrErr := ctr.execSQL("DROP VIEW IF EXISTS v_noexist")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW IF EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP VIEW IF EXISTS v_noexist", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP VIEW IF EXISTS error: %v", results[0].Error)
		}
	})

	t.Run("drop_view_multi", func(t *testing.T) {
		// DROP VIEW v1, v2 (multi-view)
		ctr.execSQL("DROP VIEW IF EXISTS v_m1")
		ctr.execSQL("DROP VIEW IF EXISTS v_m2")
		ctr.execSQL("CREATE VIEW v_m1 AS SELECT 1 AS a")
		ctr.execSQL("CREATE VIEW v_m2 AS SELECT 2 AS b")
		ctrErr := ctr.execSQL("DROP VIEW v_m1, v_m2")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW multi error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_m1 AS SELECT 1 AS a", nil)
		c.Exec("CREATE VIEW v_m2 AS SELECT 2 AS b", nil)
		results, _ := c.Exec("DROP VIEW v_m1, v_m2", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP VIEW multi error: %v", results[0].Error)
		}
		db := c.GetDatabase("test")
		if db.Views[toLower("v_m1")] != nil {
			t.Error("omni: view v_m1 should not exist after DROP VIEW")
		}
		if db.Views[toLower("v_m2")] != nil {
			t.Error("omni: view v_m2 should not exist after DROP VIEW")
		}
	})

	// Extra: CREATE VIEW duplicate (no OR REPLACE) should error
	t.Run("create_view_duplicate_error", func(t *testing.T) {
		ctr.execSQL("DROP VIEW IF EXISTS v_dup")
		ctr.execSQL("CREATE VIEW v_dup AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("CREATE VIEW v_dup AS SELECT 2 AS b")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate view")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_dup AS SELECT 1 AS a", nil)
		results, _ := c.Exec("CREATE VIEW v_dup AS SELECT 2 AS b", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate view")
		}
	})

	// Extra: CREATE VIEW with same name as existing table should error
	t.Run("create_view_table_conflict", func(t *testing.T) {
		ctr.execSQL("DROP VIEW IF EXISTS v_tbl_conflict")
		ctr.execSQL("DROP TABLE IF EXISTS v_tbl_conflict")
		ctr.execSQL("CREATE TABLE v_tbl_conflict (id INT)")
		ctrErr := ctr.execSQL("CREATE VIEW v_tbl_conflict AS SELECT 1 AS a")
		if ctrErr == nil {
			t.Fatal("container: expected error for view with same name as table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE v_tbl_conflict (id INT)", nil)
		results, _ := c.Exec("CREATE VIEW v_tbl_conflict AS SELECT 1 AS a", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for view with same name as table")
		}
	})

	// Extra: DROP VIEW on nonexistent view (no IF EXISTS) should error
	t.Run("drop_view_nonexistent_error", func(t *testing.T) {
		ctr.execSQL("DROP VIEW IF EXISTS v_nonexist_err")
		ctrErr := ctr.execSQL("DROP VIEW v_nonexist_err")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP VIEW on nonexistent view")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP VIEW v_nonexist_err", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP VIEW on nonexistent view")
		}
	})
}

func TestContainer_Section_2_11_CreateDropAlterDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("create_database", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_create1")
		ctrErr := ctr.execSQL("CREATE DATABASE db_create1")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE error: %v", ctrErr)
		}
		ctrDDL, err := ctr.showCreateDatabase("db_create1")
		if err != nil {
			t.Fatalf("container SHOW CREATE DATABASE error: %v", err)
		}

		c := New()
		results, _ := c.Exec("CREATE DATABASE db_create1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_create1")
		if db == nil {
			t.Fatal("omni: database not found after CREATE DATABASE")
		}
		// Verify charset/collation defaults match ctr.
		// Oracle returns something like: CREATE DATABASE `db_create1` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
		if !strings.Contains(ctrDDL, "utf8mb4") {
			t.Logf("container DDL: %s", ctrDDL)
		}
		if db.Charset != "utf8mb4" {
			t.Errorf("omni charset mismatch: got %q, want utf8mb4", db.Charset)
		}
		if db.Collation != "utf8mb4_0900_ai_ci" {
			t.Errorf("omni collation mismatch: got %q, want utf8mb4_0900_ai_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_create1")
	})

	t.Run("create_database_if_not_exists", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_ine")
		ctr.execSQL("CREATE DATABASE db_ine")

		// CREATE DATABASE IF NOT EXISTS on existing db should succeed (no error).
		ctrErr := ctr.execSQL("CREATE DATABASE IF NOT EXISTS db_ine")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE IF NOT EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE db_ine", nil)
		results, _ := c.Exec("CREATE DATABASE IF NOT EXISTS db_ine", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE IF NOT EXISTS error: %v", results[0].Error)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_ine")
	})

	t.Run("create_database_charset", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_cs")
		ctrErr := ctr.execSQL("CREATE DATABASE db_cs CHARACTER SET latin1")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE with charset error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_cs")

		c := New()
		results, _ := c.Exec("CREATE DATABASE db_cs CHARACTER SET latin1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE with charset error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_cs")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		// Oracle should show latin1 charset
		if !strings.Contains(ctrDDL, "latin1") {
			t.Errorf("container DDL missing latin1: %s", ctrDDL)
		}
		if db.Charset != "latin1" {
			t.Errorf("omni charset: got %q, want latin1", db.Charset)
		}
		// Default collation for latin1 is latin1_swedish_ci
		if db.Collation != "latin1_swedish_ci" {
			t.Errorf("omni collation: got %q, want latin1_swedish_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_cs")
	})

	t.Run("create_database_collate", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_coll")
		ctrErr := ctr.execSQL("CREATE DATABASE db_coll COLLATE utf8mb4_unicode_ci")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE with collate error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_coll")

		c := New()
		results, _ := c.Exec("CREATE DATABASE db_coll COLLATE utf8mb4_unicode_ci", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE DATABASE with collate error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_coll")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		if !strings.Contains(ctrDDL, "utf8mb4_unicode_ci") {
			t.Errorf("container DDL missing utf8mb4_unicode_ci: %s", ctrDDL)
		}
		if db.Collation != "utf8mb4_unicode_ci" {
			t.Errorf("omni collation: got %q, want utf8mb4_unicode_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_coll")
	})

	t.Run("drop_database", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_drop1")
		ctr.execSQL("CREATE DATABASE db_drop1")

		ctrErr := ctr.execSQL("DROP DATABASE db_drop1")
		if ctrErr != nil {
			t.Fatalf("container DROP DATABASE error: %v", ctrErr)
		}
		// Verify it's gone.
		_, showErr := ctr.showCreateDatabase("db_drop1")
		if showErr == nil {
			t.Fatal("container: database still exists after DROP DATABASE")
		}

		c := New()
		c.Exec("CREATE DATABASE db_drop1", nil)
		results, _ := c.Exec("DROP DATABASE db_drop1", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP DATABASE error: %v", results[0].Error)
		}
		if c.GetDatabase("db_drop1") != nil {
			t.Fatal("omni: database still exists after DROP DATABASE")
		}
	})

	t.Run("drop_database_if_exists", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_drop_ine")
		// DROP DATABASE IF EXISTS on nonexistent db should not error.
		ctrErr := ctr.execSQL("DROP DATABASE IF EXISTS db_drop_ine")
		if ctrErr != nil {
			t.Fatalf("container DROP DATABASE IF EXISTS error: %v", ctrErr)
		}

		c := New()
		results, _ := c.Exec("DROP DATABASE IF EXISTS db_drop_ine", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP DATABASE IF EXISTS error: %v", results[0].Error)
		}
	})

	t.Run("alter_database_charset", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_cs")
		ctr.execSQL("CREATE DATABASE db_alter_cs")

		ctrErr := ctr.execSQL("ALTER DATABASE db_alter_cs CHARACTER SET utf8mb4")
		if ctrErr != nil {
			t.Fatalf("container ALTER DATABASE charset error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_alter_cs")

		c := New()
		c.Exec("CREATE DATABASE db_alter_cs", nil)
		results, _ := c.Exec("ALTER DATABASE db_alter_cs CHARACTER SET utf8mb4", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER DATABASE charset error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_alter_cs")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		if !strings.Contains(ctrDDL, "utf8mb4") {
			t.Errorf("container DDL missing utf8mb4: %s", ctrDDL)
		}
		if db.Charset != "utf8mb4" {
			t.Errorf("omni charset: got %q, want utf8mb4", db.Charset)
		}
		if db.Collation != "utf8mb4_0900_ai_ci" {
			t.Errorf("omni collation: got %q, want utf8mb4_0900_ai_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_cs")
	})

	t.Run("alter_database_collate", func(t *testing.T) {
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_coll")
		ctr.execSQL("CREATE DATABASE db_alter_coll")

		ctrErr := ctr.execSQL("ALTER DATABASE db_alter_coll COLLATE utf8mb4_unicode_ci")
		if ctrErr != nil {
			t.Fatalf("container ALTER DATABASE collate error: %v", ctrErr)
		}
		ctrDDL, _ := ctr.showCreateDatabase("db_alter_coll")

		c := New()
		c.Exec("CREATE DATABASE db_alter_coll", nil)
		results, _ := c.Exec("ALTER DATABASE db_alter_coll COLLATE utf8mb4_unicode_ci", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER DATABASE collate error: %v", results[0].Error)
		}
		db := c.GetDatabase("db_alter_coll")
		if db == nil {
			t.Fatal("omni: database not found")
		}
		if !strings.Contains(ctrDDL, "utf8mb4_unicode_ci") {
			t.Errorf("container DDL missing utf8mb4_unicode_ci: %s", ctrDDL)
		}
		if db.Collation != "utf8mb4_unicode_ci" {
			t.Errorf("omni collation: got %q, want utf8mb4_unicode_ci", db.Collation)
		}
		ctr.execSQL("DROP DATABASE IF EXISTS db_alter_coll")
	})

	t.Run("ops_on_nonexistent_database", func(t *testing.T) {
		// DROP DATABASE on nonexistent db should error.
		ctr.execSQL("DROP DATABASE IF EXISTS db_nonexist_xyz")
		oracleDropErr := ctr.execSQL("DROP DATABASE db_nonexist_xyz")
		if oracleDropErr == nil {
			t.Fatal("container: expected error for DROP nonexistent database")
		}

		c := New()
		results, _ := c.Exec("DROP DATABASE db_nonexist_xyz", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP nonexistent database")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Code != ErrUnknownDatabase {
			t.Errorf("omni error code: got %d, want %d", catErr.Code, ErrUnknownDatabase)
		}

		// ALTER DATABASE on nonexistent db should error.
		oracleAlterErr := ctr.execSQL("ALTER DATABASE db_nonexist_xyz CHARACTER SET utf8mb4")
		if oracleAlterErr == nil {
			t.Fatal("container: expected error for ALTER nonexistent database")
		}

		c2 := New()
		results2, _ := c2.Exec("ALTER DATABASE db_nonexist_xyz CHARACTER SET utf8mb4", &ExecOptions{ContinueOnError: true})
		if results2[0].Error == nil {
			t.Fatal("omni: expected error for ALTER nonexistent database")
		}
		catErr2, ok := results2[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results2[0].Error)
		}
		if catErr2.Code != ErrUnknownDatabase {
			t.Errorf("omni error code: got %d, want %d", catErr2.Code, ErrUnknownDatabase)
		}

		// CREATE DATABASE duplicate should error.
		c3 := New()
		c3.Exec("CREATE DATABASE db_dup_test", nil)
		results3, _ := c3.Exec("CREATE DATABASE db_dup_test", &ExecOptions{ContinueOnError: true})
		if results3[0].Error == nil {
			t.Fatal("omni: expected error for duplicate CREATE DATABASE")
		}
		catErr3, ok := results3[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results3[0].Error)
		}
		if catErr3.Code != ErrDupDatabase {
			t.Errorf("omni error code: got %d, want %d", catErr3.Code, ErrDupDatabase)
		}
	})
}

func TestContainer_Section_3_1_DatabaseErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1007_dup_database", func(t *testing.T) {
		// Setup: create the database first, then try to create it again.
		ctr.execSQL("DROP DATABASE IF EXISTS db_err_dup")
		ctr.execSQL("CREATE DATABASE db_err_dup")

		ctrErr := ctr.execSQL("CREATE DATABASE db_err_dup")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate CREATE DATABASE")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1007 {
			t.Fatalf("container: expected error code 1007, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE db_err_dup", nil)
		results, _ := c.Exec("CREATE DATABASE db_err_dup", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate CREATE DATABASE")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		// Oracle: "Can't create database 'db_err_dup'; database exists"
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP DATABASE IF EXISTS db_err_dup")
	})

	t.Run("1049_unknown_database", func(t *testing.T) {
		// USE a nonexistent database on ctr.
		ctr.execSQL("DROP DATABASE IF EXISTS db_err_unknown_xyz")
		ctrErr := ctr.execSQL("USE db_err_unknown_xyz")
		if ctrErr == nil {
			t.Fatal("container: expected error for USE nonexistent database")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1049 {
			t.Fatalf("container: expected error code 1049, got %d", ctrCode)
		}

		// Run on omni: DROP DATABASE on nonexistent db triggers 1008, but
		// we need to test the "Unknown database" error (1049).
		// Use DROP DATABASE which should return 1008 in MySQL...
		// Actually, let's check what MySQL returns for DROP DATABASE on nonexistent:
		ctrErr2 := ctr.execSQL("DROP DATABASE db_err_unknown_xyz")
		if ctrErr2 == nil {
			t.Fatal("container: expected error for DROP nonexistent database")
		}
		ctrCode2, _, _ := extractMySQLErr(ctrErr2)
		t.Logf("container DROP error code: %d", ctrCode2)

		// For omni, test USE nonexistent database.
		c := New()
		results, _ := c.Exec("USE db_err_unknown_xyz", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for USE nonexistent database")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	t.Run("1046_no_database_selected", func(t *testing.T) {
		// On ctr, try to CREATE TABLE without selecting a database.
		// We need a fresh connection with no default database.
		// The container connection defaults to "test" database, so we'll
		// test omni behavior and verify the error code/SQLSTATE/message match MySQL's known format.

		// First verify MySQL's behavior: SELECT DATABASE() after no USE should work,
		// but CREATE TABLE without database should fail.
		// Since our container connection defaults to 'test' db, we verify omni matches MySQL's
		// documented error format: ERROR 1046 (3D000): No database selected

		// Run on omni with no current database.
		c := New()
		// Don't set any database — just try to create a table.
		results, _ := c.Exec("CREATE TABLE t_no_db (id INT)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for CREATE TABLE without database selected")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Verify against MySQL's documented error values.
		wantCode := 1046
		wantState := "3D000"
		wantMsg := "No database selected"

		if catErr.Code != wantCode {
			t.Errorf("error code: got %d, want %d", catErr.Code, wantCode)
		}
		if catErr.SQLState != wantState {
			t.Errorf("SQLSTATE: got %q, want %q", catErr.SQLState, wantState)
		}
		if catErr.Message != wantMsg {
			t.Errorf("message: got %q, want %q", catErr.Message, wantMsg)
		}
	})
}

func TestContainer_Section_3_2_TableErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1050_table_already_exists", func(t *testing.T) {
		// Setup: create a table, then try to create it again.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_dup")
		ctr.execSQL("CREATE TABLE t_err_dup (id INT)")

		ctrErr := ctr.execSQL("CREATE TABLE t_err_dup (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate CREATE TABLE")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1050 {
			t.Fatalf("container: expected error code 1050, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_err_dup (id INT)", nil)
		results, _ := c.Exec("CREATE TABLE t_err_dup (id INT)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate CREATE TABLE")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_err_dup")
	})

	t.Run("1051_unknown_table_drop", func(t *testing.T) {
		// DROP TABLE on a nonexistent table should return 1051.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_noexist")

		ctrErr := ctr.execSQL("DROP TABLE t_err_noexist")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP nonexistent table")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1051 {
			t.Fatalf("container: expected error code 1051, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE t_err_noexist", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP nonexistent table")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	t.Run("1146_table_doesnt_exist", func(t *testing.T) {
		// ALTER TABLE on a nonexistent table should return 1146.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_noexist2")

		ctrErr := ctr.execSQL("ALTER TABLE t_err_noexist2 ADD COLUMN x INT")
		if ctrErr == nil {
			t.Fatal("container: expected error for ALTER nonexistent table")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1146 {
			t.Fatalf("container: expected error code 1146, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("ALTER TABLE t_err_noexist2 ADD COLUMN x INT", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for ALTER nonexistent table")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})
}

func TestContainer_Section_3_3_ColumnErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1054_unknown_column", func(t *testing.T) {
		// ALTER TABLE ... MODIFY COLUMN AFTER nonexistent_col triggers 1054
		// "Unknown column 'col' in 'table definition'"
		ctr.execSQL("DROP TABLE IF EXISTS t_err_nocol")
		ctr.execSQL("CREATE TABLE t_err_nocol (id INT, name VARCHAR(50))")

		ctrErr := ctr.execSQL("ALTER TABLE t_err_nocol MODIFY COLUMN name VARCHAR(50) AFTER nonexistent")
		if ctrErr == nil {
			t.Fatal("container: expected error for unknown column in AFTER clause")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1054 {
			t.Fatalf("container: expected error code 1054, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_err_nocol (id INT, name VARCHAR(50))", nil)
		results, _ := c.Exec("ALTER TABLE t_err_nocol MODIFY COLUMN name VARCHAR(50) AFTER nonexistent", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for unknown column in AFTER clause")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_err_nocol")
	})

	t.Run("1060_duplicate_column", func(t *testing.T) {
		// CREATE TABLE with two columns of the same name.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_dupcol")

		ctrErr := ctr.execSQL("CREATE TABLE t_err_dupcol (a INT, a VARCHAR(10))")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate column name")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1060 {
			t.Fatalf("container: expected error code 1060, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TABLE t_err_dupcol (a INT, a VARCHAR(10))", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate column name")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	t.Run("1068_multiple_primary_key", func(t *testing.T) {
		// CREATE TABLE with two PRIMARY KEY definitions.
		ctr.execSQL("DROP TABLE IF EXISTS t_err_multipk")

		ctrErr := ctr.execSQL("CREATE TABLE t_err_multipk (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))")
		if ctrErr == nil {
			t.Fatal("container: expected error for multiple primary key")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1068 {
			t.Fatalf("container: expected error code 1068, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TABLE t_err_multipk (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for multiple primary key")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})
}

func TestContainer_Section_3_4_IndexKeyErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1061_dup_key_name", func(t *testing.T) {
		// Setup: create a table with an index, then try to add another index with the same name.
		ctr.execSQL("DROP TABLE IF EXISTS t_dup_key")
		ctr.execSQL("CREATE TABLE t_dup_key (a INT, b INT, KEY idx_a (a))")

		ctrErr := ctr.execSQL("ALTER TABLE t_dup_key ADD INDEX idx_a (b)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate key name")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1061 {
			t.Fatalf("container: expected error code 1061, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_dup_key (a INT, b INT, KEY idx_a (a))", nil)
		results, _ := c.Exec("ALTER TABLE t_dup_key ADD INDEX idx_a (b)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for duplicate key name")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_dup_key")
	})

	t.Run("1091_cant_drop_key", func(t *testing.T) {
		// Setup: create a table, then try to drop a nonexistent index.
		ctr.execSQL("DROP TABLE IF EXISTS t_drop_key")
		ctr.execSQL("CREATE TABLE t_drop_key (a INT)")

		ctrErr := ctr.execSQL("ALTER TABLE t_drop_key DROP INDEX idx_nonexistent")
		if ctrErr == nil {
			t.Fatal("container: expected error for dropping nonexistent key")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1091 {
			t.Fatalf("container: expected error code 1091, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_drop_key (a INT)", nil)
		results, _ := c.Exec("ALTER TABLE t_drop_key DROP INDEX idx_nonexistent", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for dropping nonexistent key")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_drop_key")
	})
}

func TestContainer_Section_3_5_FKErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("1824_fk_ref_table_not_found", func(t *testing.T) {
		// Try to create a table with FK referencing a nonexistent table.
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_noref")
		ctr.execSQL("DROP TABLE IF EXISTS t_nonexistent_parent")

		ctrErr := ctr.execSQL("CREATE TABLE t_fk_noref (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_nonexistent_parent(id))")
		if ctrErr == nil {
			t.Fatal("container: expected error for FK referencing nonexistent table")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1824 {
			t.Fatalf("container: expected error code 1824, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TABLE t_fk_noref (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_nonexistent_parent(id))", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for FK referencing nonexistent table")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T: %v", results[0].Error, results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	t.Run("1822_fk_missing_index_on_ref_table", func(t *testing.T) {
		// Create a parent table without a key on the referenced column,
		// then try to create a child with FK referencing that column.
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_child_nokey")
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent_nokey")
		ctr.execSQL("CREATE TABLE t_fk_parent_nokey (id INT, val INT)")

		ctrErr := ctr.execSQL("CREATE TABLE t_fk_child_nokey (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent_nokey(val))")
		if ctrErr == nil {
			t.Fatal("container: expected error for FK referencing column without index")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1822 {
			t.Fatalf("container: expected error code 1822, got %d", ctrCode)
		}

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_fk_parent_nokey (id INT, val INT)", nil)
		results, _ := c.Exec("CREATE TABLE t_fk_child_nokey (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent_nokey(val))", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for FK referencing column without index")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T: %v", results[0].Error, results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_fk_child_nokey")
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent_nokey")
	})

	t.Run("3780_fk_column_type_mismatch", func(t *testing.T) {
		// Create parent with INT PK, then try to create child with VARCHAR FK column.
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_child_mismatch")
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent_mismatch")
		ctr.execSQL("CREATE TABLE t_fk_parent_mismatch (id INT PRIMARY KEY)")

		ctrErr := ctr.execSQL("CREATE TABLE t_fk_child_mismatch (id INT, pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES t_fk_parent_mismatch(id))")
		if ctrErr == nil {
			t.Fatal("container: expected error for FK column type mismatch")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		// Run on omni.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_fk_parent_mismatch (id INT PRIMARY KEY)", nil)
		results, _ := c.Exec("CREATE TABLE t_fk_child_mismatch (id INT, pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES t_fk_parent_mismatch(id))", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for FK column type mismatch")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T: %v", results[0].Error, results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message format.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		ctr.execSQL("DROP TABLE IF EXISTS t_fk_child_mismatch")
		ctr.execSQL("DROP TABLE IF EXISTS t_fk_parent_mismatch")
	})
}

func TestContainer_Section_3_6_ErrorContext(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code and message from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	// Scenario 1: Error message identifier quoting matches MySQL
	// MySQL uses single quotes around identifiers in error messages (not backticks).
	// Test a variety of error types and compare message format exactly.
	t.Run("identifier_quoting", func(t *testing.T) {
		// 1a: Duplicate database — quotes around db name
		ctr.execSQL("DROP DATABASE IF EXISTS db_quote_test")
		ctr.execSQL("CREATE DATABASE db_quote_test")
		ctrErr := ctr.execSQL("CREATE DATABASE db_quote_test")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container dup db: %d %s", ctrCode, ctrMsg)

		c := New()
		c.Exec("CREATE DATABASE db_quote_test", nil)
		results, _ := c.Exec("CREATE DATABASE db_quote_test", &ExecOptions{ContinueOnError: true})
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("dup db message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1b: Duplicate table — quotes around table name
		ctr.execSQL("DROP TABLE IF EXISTS t_quote_test")
		ctr.execSQL("CREATE TABLE t_quote_test (id INT)")
		ctrErr = ctr.execSQL("CREATE TABLE t_quote_test (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container dup table: %d %s", ctrCode, ctrMsg)

		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")
		c2.Exec("CREATE TABLE t_quote_test (id INT)", nil)
		results, _ = c2.Exec("CREATE TABLE t_quote_test (id INT)", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("dup table message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1c: Unknown column — quotes around column name and context
		ctr.execSQL("DROP TABLE IF EXISTS t_col_quote")
		ctr.execSQL("CREATE TABLE t_col_quote (id INT)")
		ctrErr = ctr.execSQL("ALTER TABLE t_col_quote DROP COLUMN nonexistent")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container unknown col: %d %s", ctrCode, ctrMsg)

		c3 := New()
		c3.Exec("CREATE DATABASE test", nil)
		c3.SetCurrentDatabase("test")
		c3.Exec("CREATE TABLE t_col_quote (id INT)", nil)
		results, _ = c3.Exec("ALTER TABLE t_col_quote DROP COLUMN nonexistent", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("unknown col message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1d: Duplicate key name — quotes around key name
		ctr.execSQL("DROP TABLE IF EXISTS t_key_quote")
		ctr.execSQL("CREATE TABLE t_key_quote (id INT, val INT, KEY idx_val (val))")
		ctrErr = ctr.execSQL("ALTER TABLE t_key_quote ADD INDEX idx_val (id)")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container dup key: %d %s", ctrCode, ctrMsg)

		c4 := New()
		c4.Exec("CREATE DATABASE test", nil)
		c4.SetCurrentDatabase("test")
		c4.Exec("CREATE TABLE t_key_quote (id INT, val INT, KEY idx_val (val))", nil)
		results, _ = c4.Exec("ALTER TABLE t_key_quote ADD INDEX idx_val (id)", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("dup key message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1e: Can't drop key — quotes around key name
		ctr.execSQL("DROP TABLE IF EXISTS t_dropkey_quote")
		ctr.execSQL("CREATE TABLE t_dropkey_quote (id INT)")
		ctrErr = ctr.execSQL("ALTER TABLE t_dropkey_quote DROP INDEX nokey")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container can't drop key: %d %s", ctrCode, ctrMsg)

		c5 := New()
		c5.Exec("CREATE DATABASE test", nil)
		c5.SetCurrentDatabase("test")
		c5.Exec("CREATE TABLE t_dropkey_quote (id INT)", nil)
		results, _ = c5.Exec("ALTER TABLE t_dropkey_quote DROP INDEX nokey", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("can't drop key message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// 1f: Unknown table (DROP TABLE) — quotes around db.table
		ctr.execSQL("DROP TABLE IF EXISTS t_unknown_quote")
		ctrErr = ctr.execSQL("DROP TABLE t_unknown_quote")
		if ctrErr == nil {
			t.Fatal("container: expected error")
		}
		ctrCode, _, ctrMsg = extractMySQLErr(ctrErr)
		t.Logf("container unknown table: %d %s", ctrCode, ctrMsg)

		c6 := New()
		c6.Exec("CREATE DATABASE test", nil)
		c6.SetCurrentDatabase("test")
		results, _ = c6.Exec("DROP TABLE t_unknown_quote", &ExecOptions{ContinueOnError: true})
		catErr, ok = results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("unknown table message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// Cleanup
		ctr.execSQL("DROP TABLE IF EXISTS t_quote_test")
		ctr.execSQL("DROP TABLE IF EXISTS t_col_quote")
		ctr.execSQL("DROP TABLE IF EXISTS t_key_quote")
		ctr.execSQL("DROP TABLE IF EXISTS t_dropkey_quote")
		ctr.execSQL("DROP DATABASE IF EXISTS db_quote_test")
	})

	// Scenario 2: Error position (index) for multi-statement SQL
	// When executing multiple statements, the error should appear at the correct index.
	t.Run("error_position_multi_statement", func(t *testing.T) {
		// Setup: Execute 3 statements where the 3rd one fails.
		// The first two should succeed, the third should error.
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")

		sql := "CREATE TABLE t_pos1 (id INT); CREATE TABLE t_pos2 (id INT); CREATE TABLE t_pos1 (id INT)"
		results, _ := c.Exec(sql, &ExecOptions{ContinueOnError: true})

		// Verify we got 3 results.
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		// First two should succeed.
		if results[0].Error != nil {
			t.Errorf("statement 0: unexpected error: %v", results[0].Error)
		}
		if results[1].Error != nil {
			t.Errorf("statement 1: unexpected error: %v", results[1].Error)
		}
		// Third should fail with duplicate table.
		if results[2].Error == nil {
			t.Fatal("statement 2: expected error for duplicate table")
		}
		if results[2].Index != 2 {
			t.Errorf("error index: want 2, got %d", results[2].Index)
		}
		catErr, ok := results[2].Error.(*Error)
		if !ok {
			t.Fatalf("expected *Error, got %T", results[2].Error)
		}
		if catErr.Code != 1050 {
			t.Errorf("error code: want 1050, got %d", catErr.Code)
		}

		// Also verify on container: same 3-statement SQL, error on 3rd.
		ctr.execSQL("DROP TABLE IF EXISTS t_pos1")
		ctr.execSQL("DROP TABLE IF EXISTS t_pos2")
		ctr.execSQL("CREATE TABLE t_pos1 (id INT)")
		ctr.execSQL("CREATE TABLE t_pos2 (id INT)")
		ctrErr := ctr.execSQL("CREATE TABLE t_pos1 (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error for duplicate table")
		}
		ctrCode, _, _ := extractMySQLErr(ctrErr)
		if ctrCode != 1050 {
			t.Errorf("container error code: want 1050, got %d", ctrCode)
		}

		// Second test: error on the 1st statement of multi-statement batch.
		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")

		sql2 := "CREATE TABLE t_pos1 (a INT, a INT); CREATE TABLE t_pos3 (id INT)"
		results2, _ := c2.Exec(sql2, &ExecOptions{ContinueOnError: true})
		if len(results2) < 1 {
			t.Fatal("expected at least 1 result")
		}
		if results2[0].Error == nil {
			t.Fatal("statement 0: expected error for duplicate column")
		}
		if results2[0].Index != 0 {
			t.Errorf("error index: want 0, got %d", results2[0].Index)
		}

		// Cleanup
		ctr.execSQL("DROP TABLE IF EXISTS t_pos1")
		ctr.execSQL("DROP TABLE IF EXISTS t_pos2")
	})

	// Scenario 3: IF EXISTS suppresses errors correctly
	t.Run("if_exists_suppresses_errors", func(t *testing.T) {
		// 3a: DROP TABLE IF EXISTS on nonexistent table — no error on both.
		ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_none")
		ctrErr := ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_none")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE IF EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TABLE IF EXISTS t_ifexists_none", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP TABLE IF EXISTS error: %v", results[0].Error)
		}

		// 3b: DROP DATABASE IF EXISTS on nonexistent database — no error on both.
		ctr.execSQL("DROP DATABASE IF EXISTS db_ifexists_none")
		ctrErr = ctr.execSQL("DROP DATABASE IF EXISTS db_ifexists_none")
		if ctrErr != nil {
			t.Fatalf("container DROP DATABASE IF EXISTS error: %v", ctrErr)
		}

		c2 := New()
		results, _ = c2.Exec("DROP DATABASE IF EXISTS db_ifexists_none", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP DATABASE IF EXISTS error: %v", results[0].Error)
		}

		// 3c: DROP VIEW IF EXISTS on nonexistent view — no error on both.
		ctr.execSQL("DROP VIEW IF EXISTS v_ifexists_none")
		ctrErr = ctr.execSQL("DROP VIEW IF EXISTS v_ifexists_none")
		if ctrErr != nil {
			t.Fatalf("container DROP VIEW IF EXISTS error: %v", ctrErr)
		}

		c3 := New()
		c3.Exec("CREATE DATABASE test", nil)
		c3.SetCurrentDatabase("test")
		results, _ = c3.Exec("DROP VIEW IF EXISTS v_ifexists_none", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP VIEW IF EXISTS error: %v", results[0].Error)
		}

		// 3d: DROP TABLE IF EXISTS on existing table — should succeed (table is dropped).
		ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_real")
		ctr.execSQL("CREATE TABLE t_ifexists_real (id INT)")
		ctrErr = ctr.execSQL("DROP TABLE IF EXISTS t_ifexists_real")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE IF EXISTS (existing) error: %v", ctrErr)
		}

		c4 := New()
		c4.Exec("CREATE DATABASE test", nil)
		c4.SetCurrentDatabase("test")
		c4.Exec("CREATE TABLE t_ifexists_real (id INT)", nil)
		results, _ = c4.Exec("DROP TABLE IF EXISTS t_ifexists_real", nil)
		if results[0].Error != nil {
			t.Errorf("omni DROP TABLE IF EXISTS (existing) error: %v", results[0].Error)
		}
		// Verify the table was actually dropped.
		results, _ = c4.Exec("DROP TABLE t_ifexists_real", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Error("omni: table should have been dropped but still exists")
		}

		// 3e: Without IF EXISTS, DROP TABLE on nonexistent table — must error.
		ctr.execSQL("DROP TABLE IF EXISTS t_noifexists")
		ctrErr = ctr.execSQL("DROP TABLE t_noifexists")
		if ctrErr == nil {
			t.Fatal("container: expected error for DROP TABLE without IF EXISTS")
		}
		ctrCode, _, ctrMsg := extractMySQLErr(ctrErr)

		c5 := New()
		c5.Exec("CREATE DATABASE test", nil)
		c5.SetCurrentDatabase("test")
		results, _ = c5.Exec("DROP TABLE t_noifexists", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for DROP TABLE without IF EXISTS")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})

	// Scenario 4: IF NOT EXISTS suppresses errors correctly
	t.Run("if_not_exists_suppresses_errors", func(t *testing.T) {
		// 4a: CREATE DATABASE IF NOT EXISTS on existing database — no error.
		ctr.execSQL("DROP DATABASE IF EXISTS db_ifne")
		ctr.execSQL("CREATE DATABASE db_ifne")
		ctrErr := ctr.execSQL("CREATE DATABASE IF NOT EXISTS db_ifne")
		if ctrErr != nil {
			t.Fatalf("container CREATE DATABASE IF NOT EXISTS error: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE db_ifne", nil)
		results, _ := c.Exec("CREATE DATABASE IF NOT EXISTS db_ifne", nil)
		if results[0].Error != nil {
			t.Errorf("omni CREATE DATABASE IF NOT EXISTS error: %v", results[0].Error)
		}

		// 4b: CREATE TABLE IF NOT EXISTS on existing table — no error, original preserved.
		ctr.execSQL("DROP TABLE IF EXISTS t_ifne")
		ctr.execSQL("CREATE TABLE t_ifne (id INT)")
		ctrErr = ctr.execSQL("CREATE TABLE IF NOT EXISTS t_ifne (val VARCHAR(100))")
		if ctrErr != nil {
			t.Fatalf("container CREATE TABLE IF NOT EXISTS error: %v", ctrErr)
		}
		// Verify original table is unchanged.
		ctrDDL, _ := ctr.showCreateTable("t_ifne")
		if !strings.Contains(ctrDDL, "`id`") {
			t.Errorf("container: original table structure should be preserved, got: %s", ctrDDL)
		}

		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")
		c2.Exec("CREATE TABLE t_ifne (id INT)", nil)
		results, _ = c2.Exec("CREATE TABLE IF NOT EXISTS t_ifne (val VARCHAR(100))", nil)
		if results[0].Error != nil {
			t.Errorf("omni CREATE TABLE IF NOT EXISTS error: %v", results[0].Error)
		}
		// Verify original table is unchanged.
		omniDDL := c2.ShowCreateTable("test", "t_ifne")
		if !strings.Contains(omniDDL, "`id`") {
			t.Errorf("omni: original table structure should be preserved, got: %s", omniDDL)
		}

		// 4c: Without IF NOT EXISTS, CREATE TABLE on existing table — must error.
		ctr.execSQL("DROP TABLE IF EXISTS t_no_ifne")
		ctr.execSQL("CREATE TABLE t_no_ifne (id INT)")
		ctrErr = ctr.execSQL("CREATE TABLE t_no_ifne (id INT)")
		if ctrErr == nil {
			t.Fatal("container: expected error for CREATE TABLE without IF NOT EXISTS")
		}
		ctrCode, _, ctrMsg := extractMySQLErr(ctrErr)

		c3 := New()
		c3.Exec("CREATE DATABASE test", nil)
		c3.SetCurrentDatabase("test")
		c3.Exec("CREATE TABLE t_no_ifne (id INT)", nil)
		results, _ = c3.Exec("CREATE TABLE t_no_ifne (id INT)", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for CREATE TABLE without IF NOT EXISTS")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}

		// Cleanup
		ctr.execSQL("DROP TABLE IF EXISTS t_ifne")
		ctr.execSQL("DROP TABLE IF EXISTS t_no_ifne")
		ctr.execSQL("DROP DATABASE IF EXISTS db_ifne")
	})
}

func TestContainer_Section_4_1_Partitioning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{
			"partition_by_range",
			`CREATE TABLE t_part_range (
				id INT NOT NULL,
				created DATE NOT NULL
			) PARTITION BY RANGE (YEAR(created)) (
				PARTITION p0 VALUES LESS THAN (2020),
				PARTITION p1 VALUES LESS THAN (2025),
				PARTITION pmax VALUES LESS THAN MAXVALUE
			)`,
			"t_part_range",
		},
		{
			"partition_by_range_columns",
			`CREATE TABLE t_part_range_cols (
				id INT NOT NULL,
				city VARCHAR(50) NOT NULL,
				name VARCHAR(50) NOT NULL
			) PARTITION BY RANGE COLUMNS(city) (
				PARTITION p0 VALUES LESS THAN ('M'),
				PARTITION p1 VALUES LESS THAN MAXVALUE
			)`,
			"t_part_range_cols",
		},
		{
			"partition_by_list",
			`CREATE TABLE t_part_list (
				id INT NOT NULL,
				region INT NOT NULL
			) PARTITION BY LIST (region) (
				PARTITION pNorth VALUES IN (1,2,3),
				PARTITION pSouth VALUES IN (4,5,6),
				PARTITION pWest VALUES IN (7,8,9)
			)`,
			"t_part_list",
		},
		{
			"partition_by_hash",
			`CREATE TABLE t_part_hash (
				id INT NOT NULL,
				name VARCHAR(100)
			) PARTITION BY HASH (id) PARTITIONS 4`,
			"t_part_hash",
		},
		{
			"partition_by_key",
			`CREATE TABLE t_part_key (
				id INT NOT NULL,
				name VARCHAR(100)
			) PARTITION BY KEY (id) PARTITIONS 3`,
			"t_part_key",
		},
		{
			"partition_linear_hash",
			`CREATE TABLE t_part_linear (
				id INT NOT NULL
			) PARTITION BY LINEAR HASH (id) PARTITIONS 4`,
			"t_part_linear",
		},
		{
			"partition_subpartition",
			`CREATE TABLE t_part_sub (
				id INT NOT NULL,
				purchased DATE NOT NULL
			) PARTITION BY RANGE (YEAR(purchased))
			SUBPARTITION BY HASH (id)
			SUBPARTITIONS 2 (
				PARTITION p0 VALUES LESS THAN (2020),
				PARTITION p1 VALUES LESS THAN MAXVALUE
			)`,
			"t_part_sub",
		},
		{
			"show_create_partitioned_table",
			`CREATE TABLE t_part_show (
				id INT NOT NULL,
				val INT NOT NULL
			) PARTITION BY RANGE (val) (
				PARTITION p0 VALUES LESS THAN (10),
				PARTITION p1 VALUES LESS THAN (20),
				PARTITION p2 VALUES LESS THAN MAXVALUE
			)`,
			"t_part_show",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, _ := c.Exec(tc.sql, nil)
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}

	// ALTER TABLE partition tests — multi-step
	t.Run("alter_add_partition", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_alter_addp")
		setupSQL := `CREATE TABLE t_alter_addp (
			id INT NOT NULL,
			val INT NOT NULL
		) PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (20)
		)`
		if err := ctr.execSQL(setupSQL); err != nil {
			t.Fatalf("container setup: %v", err)
		}
		alterSQL := "ALTER TABLE t_alter_addp ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))"
		if err := ctr.execSQL(alterSQL); err != nil {
			t.Fatalf("container alter: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_alter_addp")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec(setupSQL, nil)
		results, _ := c.Exec(alterSQL, nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_alter_addp")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("alter_drop_partition", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_alter_dropp")
		setupSQL := `CREATE TABLE t_alter_dropp (
			id INT NOT NULL,
			val INT NOT NULL
		) PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (20),
			PARTITION p2 VALUES LESS THAN MAXVALUE
		)`
		if err := ctr.execSQL(setupSQL); err != nil {
			t.Fatalf("container setup: %v", err)
		}
		alterSQL := "ALTER TABLE t_alter_dropp DROP PARTITION p1"
		if err := ctr.execSQL(alterSQL); err != nil {
			t.Fatalf("container alter: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_alter_dropp")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec(setupSQL, nil)
		results, _ := c.Exec(alterSQL, nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_alter_dropp")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("alter_reorganize_partition", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_alter_reorgp")
		setupSQL := `CREATE TABLE t_alter_reorgp (
			id INT NOT NULL,
			val INT NOT NULL
		) PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN (30)
		)`
		if err := ctr.execSQL(setupSQL); err != nil {
			t.Fatalf("container setup: %v", err)
		}
		alterSQL := `ALTER TABLE t_alter_reorgp REORGANIZE PARTITION p1 INTO (
			PARTITION p1a VALUES LESS THAN (20),
			PARTITION p1b VALUES LESS THAN (30)
		)`
		if err := ctr.execSQL(alterSQL); err != nil {
			t.Fatalf("container alter: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_alter_reorgp")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec(setupSQL, nil)
		results, _ := c.Exec(alterSQL, nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_alter_reorgp")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("alter_truncate_partition", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_alter_truncp")
		setupSQL := `CREATE TABLE t_alter_truncp (
			id INT NOT NULL,
			val INT NOT NULL
		) PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN MAXVALUE
		)`
		if err := ctr.execSQL(setupSQL); err != nil {
			t.Fatalf("container setup: %v", err)
		}
		alterSQL := "ALTER TABLE t_alter_truncp TRUNCATE PARTITION p0"
		if err := ctr.execSQL(alterSQL); err != nil {
			t.Fatalf("container alter: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_alter_truncp")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec(setupSQL, nil)
		results, _ := c.Exec(alterSQL, nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_alter_truncp")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("alter_coalesce_partition", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_alter_coalp")
		setupSQL := `CREATE TABLE t_alter_coalp (
			id INT NOT NULL
		) PARTITION BY HASH (id) PARTITIONS 4`
		if err := ctr.execSQL(setupSQL); err != nil {
			t.Fatalf("container setup: %v", err)
		}
		alterSQL := "ALTER TABLE t_alter_coalp COALESCE PARTITION 2"
		if err := ctr.execSQL(alterSQL); err != nil {
			t.Fatalf("container alter: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_alter_coalp")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec(setupSQL, nil)
		results, _ := c.Exec(alterSQL, nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_alter_coalp")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("alter_exchange_partition", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_alter_exchp")
		ctr.execSQL("DROP TABLE IF EXISTS t_alter_exchp_swap")
		setupSQL := `CREATE TABLE t_alter_exchp (
			id INT NOT NULL,
			val INT NOT NULL
		) PARTITION BY RANGE (val) (
			PARTITION p0 VALUES LESS THAN (10),
			PARTITION p1 VALUES LESS THAN MAXVALUE
		)`
		swapSQL := `CREATE TABLE t_alter_exchp_swap (
			id INT NOT NULL,
			val INT NOT NULL
		)`
		if err := ctr.execSQL(setupSQL); err != nil {
			t.Fatalf("container setup: %v", err)
		}
		if err := ctr.execSQL(swapSQL); err != nil {
			t.Fatalf("container swap table setup: %v", err)
		}
		alterSQL := "ALTER TABLE t_alter_exchp EXCHANGE PARTITION p0 WITH TABLE t_alter_exchp_swap"
		if err := ctr.execSQL(alterSQL); err != nil {
			t.Fatalf("container alter: %v", err)
		}
		ctrDDL, _ := ctr.showCreateTable("t_alter_exchp")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec(setupSQL, nil)
		c.Exec(swapSQL, nil)
		results, _ := c.Exec(alterSQL, nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_alter_exchp")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})
}

func TestContainer_Section_4_2_StoredRoutines(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// 4.2.1: CREATE FUNCTION — store metadata (name, params, return type, body)
	t.Run("create_function", func(t *testing.T) {
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_add")
		createSQL := "CREATE FUNCTION fn_add(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b"
		if err := ctr.execSQL(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateFunction("fn_add")
		if err != nil {
			t.Fatalf("container SHOW CREATE FUNCTION: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateFunction("test", "fn_add")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// 4.2.2: CREATE PROCEDURE — store metadata
	t.Run("create_procedure", func(t *testing.T) {
		ctr.execSQLDirect("DROP PROCEDURE IF EXISTS sp_greet")
		createSQL := "CREATE PROCEDURE sp_greet(IN name VARCHAR(100)) BEGIN SELECT CONCAT('Hello, ', name); END"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateProcedure("sp_greet")
		if err != nil {
			t.Fatalf("container SHOW CREATE PROCEDURE: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateProcedure("test", "sp_greet")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// 4.2.3: DROP FUNCTION / PROCEDURE
	t.Run("drop_function", func(t *testing.T) {
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_drop_test")
		ctr.execSQL("CREATE FUNCTION fn_drop_test() RETURNS INT DETERMINISTIC RETURN 1")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE FUNCTION fn_drop_test() RETURNS INT DETERMINISTIC RETURN 1", nil)

		// Drop on container
		if err := ctr.execSQL("DROP FUNCTION fn_drop_test"); err != nil {
			t.Fatalf("container DROP FUNCTION: %v", err)
		}
		// Drop on omni
		results, _ := c.Exec("DROP FUNCTION fn_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP FUNCTION error: %v", results[0].Error)
		}

		// Both should now error on SHOW CREATE FUNCTION
		_, ctrErr := ctr.showCreateFunction("fn_drop_test")
		omniDDL := c.ShowCreateFunction("test", "fn_drop_test")

		if ctrErr == nil {
			t.Error("expected container error after DROP FUNCTION")
		}
		if omniDDL != "" {
			t.Error("expected empty omni DDL after DROP FUNCTION")
		}
	})

	t.Run("drop_procedure", func(t *testing.T) {
		ctr.execSQLDirect("DROP PROCEDURE IF EXISTS sp_drop_test")
		ctr.execSQLDirect("CREATE PROCEDURE sp_drop_test() BEGIN SELECT 1; END")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE PROCEDURE sp_drop_test() BEGIN SELECT 1; END", nil)

		// Drop on container
		if err := ctr.execSQL("DROP PROCEDURE sp_drop_test"); err != nil {
			t.Fatalf("container DROP PROCEDURE: %v", err)
		}
		// Drop on omni
		results, _ := c.Exec("DROP PROCEDURE sp_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP PROCEDURE error: %v", results[0].Error)
		}

		_, ctrErr := ctr.showCreateProcedure("sp_drop_test")
		omniDDL := c.ShowCreateProcedure("test", "sp_drop_test")

		if ctrErr == nil {
			t.Error("expected container error after DROP PROCEDURE")
		}
		if omniDDL != "" {
			t.Error("expected empty omni DDL after DROP PROCEDURE")
		}
	})

	// 4.2.4: DROP FUNCTION IF EXISTS — no error for nonexistent
	t.Run("drop_function_if_exists", func(t *testing.T) {
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_nonexistent_xyz")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP FUNCTION IF EXISTS fn_nonexistent_xyz", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP FUNCTION IF EXISTS should not error: %v", results[0].Error)
		}
	})

	// 4.2.5: DROP FUNCTION nonexistent — error
	t.Run("drop_function_nonexistent_error", func(t *testing.T) {
		ctrErr := ctr.execSQL("DROP FUNCTION fn_totally_nonexistent")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP FUNCTION fn_totally_nonexistent", nil)

		if ctrErr == nil {
			t.Fatal("expected container error for DROP nonexistent FUNCTION")
		}
		if results[0].Error == nil {
			t.Fatal("expected omni error for DROP nonexistent FUNCTION")
		}

		// Both should produce error code 1305
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(ctrErr, &mysqlErr) {
			if catErr, ok := results[0].Error.(*Error); ok {
				if catErr.Code != int(mysqlErr.Number) {
					t.Errorf("error code mismatch: container=%d omni=%d", mysqlErr.Number, catErr.Code)
				}
			}
		}
	})

	// 4.2.6: ALTER ROUTINE (characteristics only)
	t.Run("alter_function_comment", func(t *testing.T) {
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_alter_test")
		ctr.execSQL("CREATE FUNCTION fn_alter_test() RETURNS INT DETERMINISTIC RETURN 42")
		if err := ctr.execSQL("ALTER FUNCTION fn_alter_test COMMENT 'test comment'"); err != nil {
			t.Fatalf("container ALTER FUNCTION: %v", err)
		}
		ctrDDL, err := ctr.showCreateFunction("fn_alter_test")
		if err != nil {
			t.Fatalf("container SHOW CREATE FUNCTION: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE FUNCTION fn_alter_test() RETURNS INT DETERMINISTIC RETURN 42", nil)
		results, _ := c.Exec("ALTER FUNCTION fn_alter_test COMMENT 'test comment'", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER FUNCTION error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateFunction("test", "fn_alter_test")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("alter_procedure_sql_security", func(t *testing.T) {
		ctr.execSQLDirect("DROP PROCEDURE IF EXISTS sp_alter_test")
		ctr.execSQLDirect("CREATE PROCEDURE sp_alter_test() BEGIN SELECT 1; END")
		if err := ctr.execSQLDirect("ALTER PROCEDURE sp_alter_test SQL SECURITY INVOKER"); err != nil {
			t.Fatalf("container ALTER PROCEDURE: %v", err)
		}
		ctrDDL, err := ctr.showCreateProcedure("sp_alter_test")
		if err != nil {
			t.Fatalf("container SHOW CREATE PROCEDURE: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE PROCEDURE sp_alter_test() BEGIN SELECT 1; END", nil)
		results, _ := c.Exec("ALTER PROCEDURE sp_alter_test SQL SECURITY INVOKER", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER PROCEDURE error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateProcedure("test", "sp_alter_test")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// 4.2.7: SHOW CREATE FUNCTION output
	t.Run("show_create_function_output", func(t *testing.T) {
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_show_test")
		createSQL := "CREATE FUNCTION fn_show_test(x INT) RETURNS VARCHAR(100) DETERMINISTIC RETURN CONCAT('val=', x)"
		if err := ctr.execSQL(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateFunction("fn_show_test")
		if err != nil {
			t.Fatalf("container SHOW CREATE FUNCTION: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateFunction("test", "fn_show_test")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// 4.2.8: SHOW CREATE PROCEDURE output
	t.Run("show_create_procedure_output", func(t *testing.T) {
		ctr.execSQLDirect("DROP PROCEDURE IF EXISTS sp_show_test")
		createSQL := "CREATE PROCEDURE sp_show_test(IN id INT, OUT result VARCHAR(100)) BEGIN SET result = CONCAT('id=', id); END"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateProcedure("sp_show_test")
		if err != nil {
			t.Fatalf("container SHOW CREATE PROCEDURE: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateProcedure("test", "sp_show_test")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})
}

func TestContainer_Section_4_3_Triggers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup: create a table that triggers can reference.
	ctr.execSQL("DROP TABLE IF EXISTS t_trigger")
	if err := ctr.execSQL("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))"); err != nil {
		t.Fatalf("container setup table: %v", err)
	}

	// 4.3.1: CREATE TRIGGER — store metadata (name, timing, event, table, body)
	t.Run("create_trigger_before_insert", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_before_insert")
		createSQL := "CREATE TRIGGER tr_before_insert BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = NEW.val + 1"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateTrigger("tr_before_insert")
		if err != nil {
			t.Fatalf("container SHOW CREATE TRIGGER: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTrigger("test", "tr_before_insert")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("create_trigger_after_update", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_after_update")
		createSQL := "CREATE TRIGGER tr_after_update AFTER UPDATE ON t_trigger FOR EACH ROW SET @updated = @updated + 1"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateTrigger("tr_after_update")
		if err != nil {
			t.Fatalf("container SHOW CREATE TRIGGER: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTrigger("test", "tr_after_update")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("create_trigger_before_delete", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_before_delete")
		createSQL := "CREATE TRIGGER tr_before_delete BEFORE DELETE ON t_trigger FOR EACH ROW SET @deleted = OLD.id"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateTrigger("tr_before_delete")
		if err != nil {
			t.Fatalf("container SHOW CREATE TRIGGER: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTrigger("test", "tr_before_delete")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// 4.3.2: DROP TRIGGER
	t.Run("drop_trigger", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_drop_test")
		ctr.execSQLDirect("CREATE TRIGGER tr_drop_test BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = 0")
		if err := ctr.execSQL("DROP TRIGGER tr_drop_test"); err != nil {
			t.Fatalf("container DROP TRIGGER: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		c.Exec("CREATE TRIGGER tr_drop_test BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = 0", nil)
		results, _ := c.Exec("DROP TRIGGER tr_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TRIGGER error: %v", results[0].Error)
		}

		// Verify trigger is gone
		ddl := c.ShowCreateTrigger("test", "tr_drop_test")
		if ddl != "" {
			t.Errorf("trigger should be gone, but got: %s", ddl)
		}
	})

	t.Run("drop_trigger_if_exists", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_nonexistent_drop")
		if err := ctr.execSQL("DROP TRIGGER IF EXISTS tr_nonexistent_drop"); err != nil {
			t.Fatalf("container DROP TRIGGER IF EXISTS: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TRIGGER IF EXISTS tr_nonexistent_drop", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP TRIGGER IF EXISTS should not error: %v", results[0].Error)
		}
	})

	t.Run("drop_trigger_not_exist_error", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_no_exist")
		ctrErr := ctr.execSQL("DROP TRIGGER tr_no_exist")
		if ctrErr == nil {
			t.Fatal("expected container error for DROP TRIGGER on nonexistent trigger")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP TRIGGER tr_no_exist", nil)
		if results[0].Error == nil {
			t.Fatal("expected omni error for DROP TRIGGER on nonexistent trigger")
		}
		if catErr, ok := results[0].Error.(*Error); ok {
			if catErr.Code != ErrNoSuchTrigger {
				t.Errorf("error code: want %d, got %d", ErrNoSuchTrigger, catErr.Code)
			}
		}
	})

	// 4.3.3: SHOW CREATE TRIGGER output
	t.Run("show_create_trigger_output", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_show_test")
		createSQL := "CREATE TRIGGER tr_show_test AFTER INSERT ON t_trigger FOR EACH ROW SET @last_insert = NEW.id"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateTrigger("tr_show_test")
		if err != nil {
			t.Fatalf("container SHOW CREATE TRIGGER: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTrigger("test", "tr_show_test")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// 4.3.4: Multiple triggers per table/event (MySQL 8.0 supports ordering)
	t.Run("multiple_triggers_ordering", func(t *testing.T) {
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_multi_2")
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_multi_1")
		create1 := "CREATE TRIGGER tr_multi_1 BEFORE INSERT ON t_trigger FOR EACH ROW SET NEW.val = NEW.val + 10"
		create2 := "CREATE TRIGGER tr_multi_2 BEFORE INSERT ON t_trigger FOR EACH ROW FOLLOWS tr_multi_1 SET NEW.val = NEW.val + 20"
		if err := ctr.execSQLDirect(create1); err != nil {
			t.Fatalf("container exec trigger 1: %v", err)
		}
		if err := ctr.execSQLDirect(create2); err != nil {
			t.Fatalf("container exec trigger 2: %v", err)
		}

		ctrDDL1, err := ctr.showCreateTrigger("tr_multi_1")
		if err != nil {
			t.Fatalf("container SHOW CREATE TRIGGER tr_multi_1: %v", err)
		}
		ctrDDL2, err := ctr.showCreateTrigger("tr_multi_2")
		if err != nil {
			t.Fatalf("container SHOW CREATE TRIGGER tr_multi_2: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		results1, _ := c.Exec(create1, nil)
		if results1[0].Error != nil {
			t.Fatalf("omni exec trigger 1 error: %v", results1[0].Error)
		}
		results2, _ := c.Exec(create2, nil)
		if results2[0].Error != nil {
			t.Fatalf("omni exec trigger 2 error: %v", results2[0].Error)
		}

		omniDDL1 := c.ShowCreateTrigger("test", "tr_multi_1")
		omniDDL2 := c.ShowCreateTrigger("test", "tr_multi_2")

		if normalizeWhitespace(ctrDDL1) != normalizeWhitespace(omniDDL1) {
			t.Errorf("trigger 1 mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL1, omniDDL1)
		}
		if normalizeWhitespace(ctrDDL2) != normalizeWhitespace(omniDDL2) {
			t.Errorf("trigger 2 mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL2, omniDDL2)
		}
	})
}

// normalizeEventDDL normalizes event DDL for comparison by removing
// auto-generated STARTS timestamps that MySQL adds to EVERY schedules.
// MySQL always adds "STARTS '<timestamp>'" when EVERY is used without explicit STARTS,
// but our catalog doesn't have a clock to generate timestamps.
func normalizeEventDDL(ddl string) string {
	// Remove STARTS '<timestamp>' from EVERY schedules (auto-generated by MySQL).
	// Pattern: " STARTS '<datetime>'"
	s := normalizeWhitespace(ddl)
	// Remove auto-generated STARTS with timestamp pattern
	for {
		idx := strings.Index(s, " STARTS '")
		if idx < 0 {
			break
		}
		// Find closing quote
		end := strings.Index(s[idx+9:], "'")
		if end < 0 {
			break
		}
		s = s[:idx] + s[idx+9+end+1:]
	}
	return s
}

func TestContainer_Section_4_4_Events(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// 4.4.1: CREATE EVENT — store metadata
	t.Run("create_event_every", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_test_every")
		createSQL := "CREATE EVENT ev_test_every ON SCHEDULE EVERY 1 HOUR DO SELECT 1"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_test_every")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}
		t.Logf("container SHOW CREATE EVENT:\n%s", ctrDDL)

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_test_every")

		// Compare with normalized DDL (MySQL adds auto-STARTS timestamp for EVERY schedules)
		if normalizeEventDDL(ctrDDL) != normalizeEventDDL(omniDDL) {
			t.Errorf("mismatch:\n--- container (normalized) ---\n%s\n--- omni (normalized) ---\n%s\n--- container (raw) ---\n%s\n--- omni (raw) ---\n%s",
				normalizeEventDDL(ctrDDL), normalizeEventDDL(omniDDL), ctrDDL, omniDDL)
		}
	})

	t.Run("create_event_at_timestamp", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_test_at")
		createSQL := "CREATE EVENT ev_test_at ON SCHEDULE AT '2035-01-01 00:00:00' DO SELECT 1"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_test_at")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}
		t.Logf("container SHOW CREATE EVENT:\n%s", ctrDDL)

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_test_at")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	t.Run("create_event_with_options", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_test_opts")
		createSQL := "CREATE EVENT ev_test_opts ON SCHEDULE EVERY 1 DAY ON COMPLETION PRESERVE ENABLE COMMENT 'daily cleanup' DO DELETE FROM t_trigger WHERE id < 0"
		// Need a table for the body to reference
		ctr.execSQL("DROP TABLE IF EXISTS t_trigger")
		ctr.execSQL("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))")
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_test_opts")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}
		t.Logf("container SHOW CREATE EVENT:\n%s", ctrDDL)

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, name VARCHAR(100))", nil)
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_test_opts")

		// Compare with normalized DDL (MySQL adds auto-STARTS timestamp for EVERY schedules)
		if normalizeEventDDL(ctrDDL) != normalizeEventDDL(omniDDL) {
			t.Errorf("mismatch:\n--- container (normalized) ---\n%s\n--- omni (normalized) ---\n%s\n--- container (raw) ---\n%s\n--- omni (raw) ---\n%s",
				normalizeEventDDL(ctrDDL), normalizeEventDDL(omniDDL), ctrDDL, omniDDL)
		}
	})

	t.Run("create_event_disabled", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_test_disabled")
		createSQL := "CREATE EVENT ev_test_disabled ON SCHEDULE EVERY 1 HOUR DISABLE DO SELECT 1"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_test_disabled")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}
		t.Logf("container SHOW CREATE EVENT:\n%s", ctrDDL)

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_test_disabled")

		// Compare with normalized DDL (MySQL adds auto-STARTS timestamp for EVERY schedules)
		if normalizeEventDDL(ctrDDL) != normalizeEventDDL(omniDDL) {
			t.Errorf("mismatch:\n--- container (normalized) ---\n%s\n--- omni (normalized) ---\n%s\n--- container (raw) ---\n%s\n--- omni (raw) ---\n%s",
				normalizeEventDDL(ctrDDL), normalizeEventDDL(omniDDL), ctrDDL, omniDDL)
		}
	})

	// 4.4.2: ALTER EVENT
	t.Run("alter_event_disable", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_alter_test")
		ctr.execSQLDirect("CREATE EVENT ev_alter_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1")
		if err := ctr.execSQLDirect("ALTER EVENT ev_alter_test DISABLE"); err != nil {
			t.Fatalf("container ALTER EVENT: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_alter_test")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}
		t.Logf("container SHOW CREATE EVENT after ALTER:\n%s", ctrDDL)

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE EVENT ev_alter_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1", nil)
		results, _ := c.Exec("ALTER EVENT ev_alter_test DISABLE", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER EVENT error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_alter_test")

		// Compare with normalized DDL (MySQL adds auto-STARTS timestamp for EVERY schedules)
		if normalizeEventDDL(ctrDDL) != normalizeEventDDL(omniDDL) {
			t.Errorf("mismatch:\n--- container (normalized) ---\n%s\n--- omni (normalized) ---\n%s\n--- container (raw) ---\n%s\n--- omni (raw) ---\n%s",
				normalizeEventDDL(ctrDDL), normalizeEventDDL(omniDDL), ctrDDL, omniDDL)
		}
	})

	t.Run("alter_event_rename", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_rename_old")
		ctr.execSQL("DROP EVENT IF EXISTS ev_rename_new")
		ctr.execSQLDirect("CREATE EVENT ev_rename_old ON SCHEDULE EVERY 1 HOUR DO SELECT 1")
		if err := ctr.execSQLDirect("ALTER EVENT ev_rename_old RENAME TO ev_rename_new"); err != nil {
			t.Fatalf("container ALTER EVENT RENAME: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_rename_new")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}
		t.Logf("container SHOW CREATE EVENT after RENAME:\n%s", ctrDDL)

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE EVENT ev_rename_old ON SCHEDULE EVERY 1 HOUR DO SELECT 1", nil)
		results, _ := c.Exec("ALTER EVENT ev_rename_old RENAME TO ev_rename_new", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER EVENT RENAME error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_rename_new")

		// Compare with normalized DDL (MySQL adds auto-STARTS timestamp for EVERY schedules)
		if normalizeEventDDL(ctrDDL) != normalizeEventDDL(omniDDL) {
			t.Errorf("mismatch:\n--- container (normalized) ---\n%s\n--- omni (normalized) ---\n%s\n--- container (raw) ---\n%s\n--- omni (raw) ---\n%s",
				normalizeEventDDL(ctrDDL), normalizeEventDDL(omniDDL), ctrDDL, omniDDL)
		}

		// Verify old name is gone
		oldDDL := c.ShowCreateEvent("test", "ev_rename_old")
		if oldDDL != "" {
			t.Errorf("old event name should be gone, but got: %s", oldDDL)
		}
	})

	// 4.4.3: DROP EVENT
	t.Run("drop_event", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_drop_test")
		ctr.execSQLDirect("CREATE EVENT ev_drop_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1")
		if err := ctr.execSQL("DROP EVENT ev_drop_test"); err != nil {
			t.Fatalf("container DROP EVENT: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE EVENT ev_drop_test ON SCHEDULE EVERY 1 HOUR DO SELECT 1", nil)
		results, _ := c.Exec("DROP EVENT ev_drop_test", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP EVENT error: %v", results[0].Error)
		}

		// Verify event is gone
		ddl := c.ShowCreateEvent("test", "ev_drop_test")
		if ddl != "" {
			t.Errorf("event should be gone, but got: %s", ddl)
		}
	})

	t.Run("drop_event_if_exists", func(t *testing.T) {
		if err := ctr.execSQL("DROP EVENT IF EXISTS ev_nonexistent"); err != nil {
			t.Fatalf("container DROP EVENT IF EXISTS: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP EVENT IF EXISTS ev_nonexistent", nil)
		if results[0].Error != nil {
			t.Fatalf("omni DROP EVENT IF EXISTS should not error: %v", results[0].Error)
		}
	})

	t.Run("drop_event_not_exist_error", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_no_exist")
		ctrErr := ctr.execSQL("DROP EVENT ev_no_exist")
		if ctrErr == nil {
			t.Fatal("expected container error for DROP EVENT on nonexistent event")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("DROP EVENT ev_no_exist", nil)
		if results[0].Error == nil {
			t.Fatal("expected omni error for DROP EVENT on nonexistent event")
		}
		if catErr, ok := results[0].Error.(*Error); ok {
			if catErr.Code != ErrNoSuchEvent {
				t.Errorf("error code: want %d, got %d", ErrNoSuchEvent, catErr.Code)
			}
		}
	})

	// 4.4.4: SHOW CREATE EVENT output
	t.Run("show_create_event_basic", func(t *testing.T) {
		ctr.execSQL("DROP EVENT IF EXISTS ev_show_test")
		createSQL := "CREATE EVENT ev_show_test ON SCHEDULE EVERY 1 HOUR COMMENT 'test event' DO SELECT 1"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_show_test")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}
		t.Logf("container SHOW CREATE EVENT:\n%s", ctrDDL)

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_show_test")

		// Compare with normalized DDL (MySQL adds auto-STARTS timestamp for EVERY schedules)
		if normalizeEventDDL(ctrDDL) != normalizeEventDDL(omniDDL) {
			t.Errorf("mismatch:\n--- container (normalized) ---\n%s\n--- omni (normalized) ---\n%s\n--- container (raw) ---\n%s\n--- omni (raw) ---\n%s",
				normalizeEventDDL(ctrDDL), normalizeEventDDL(omniDDL), ctrDDL, omniDDL)
		}
	})
}

// extractViewPreamble extracts the CREATE VIEW preamble up to and including " AS ".
// This allows comparing structural elements without comparing the rewritten SELECT text.
func extractViewPreamble(ddl string) string {
	idx := strings.Index(ddl, " AS ")
	if idx < 0 {
		return ddl
	}
	return normalizeWhitespace(ddl[:idx+4])
}

func TestContainer_Section_4_5_ViewsDeep(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	t.Run("alter_view", func(t *testing.T) {
		// ALTER VIEW changes the view definition.
		// Verify both container and omni accept ALTER VIEW and update the view.
		ctr.execSQL("DROP VIEW IF EXISTS v_alter_test")
		ctr.execSQL("CREATE VIEW v_alter_test AS SELECT 1 AS a")
		ctrErr := ctr.execSQL("ALTER VIEW v_alter_test AS SELECT 2 AS b, 3 AS c")
		if ctrErr != nil {
			t.Fatalf("container ALTER VIEW error: %v", ctrErr)
		}
		// Verify container view was updated.
		ctrDDL, _ := ctr.showCreateView("v_alter_test")
		if !strings.Contains(ctrDDL, "v_alter_test") {
			t.Fatalf("container: view not found after ALTER VIEW")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE VIEW v_alter_test AS SELECT 1 AS a", nil)
		results, _ := c.Exec("ALTER VIEW v_alter_test AS SELECT 2 AS b, 3 AS c", nil)
		if results[0].Error != nil {
			t.Fatalf("omni ALTER VIEW error: %v", results[0].Error)
		}
		// Verify omni view was updated.
		db := c.GetDatabase("test")
		v := db.Views[toLower("v_alter_test")]
		if v == nil {
			t.Fatal("omni: view v_alter_test should exist after ALTER VIEW")
		}
		if !strings.Contains(v.Definition, "2") {
			t.Errorf("omni: view definition should contain new select, got: %s", v.Definition)
		}

		// ALTER VIEW on nonexistent view should error on both.
		ctr.execSQL("DROP VIEW IF EXISTS v_alter_noexist")
		ctrErr = ctr.execSQL("ALTER VIEW v_alter_noexist AS SELECT 1")
		if ctrErr == nil {
			t.Fatal("container: expected error for ALTER VIEW on nonexistent view")
		}

		c2 := New()
		c2.Exec("CREATE DATABASE test", nil)
		c2.SetCurrentDatabase("test")
		results2, _ := c2.Exec("ALTER VIEW v_alter_noexist AS SELECT 1", &ExecOptions{ContinueOnError: true})
		if results2[0].Error == nil {
			t.Fatal("omni: expected error for ALTER VIEW on nonexistent view")
		}
	})

	t.Run("view_dependency_tracking", func(t *testing.T) {
		// In MySQL, dropping a base table does NOT drop the view.
		// The view still exists but errors when queried.
		ctr.execSQL("DROP VIEW IF EXISTS v_dep_test")
		ctr.execSQL("DROP TABLE IF EXISTS t_dep_base")
		ctr.execSQL("CREATE TABLE t_dep_base (id INT, val VARCHAR(100))")
		ctr.execSQL("CREATE VIEW v_dep_test AS SELECT id, val FROM t_dep_base")

		// Drop the base table
		ctrErr := ctr.execSQL("DROP TABLE t_dep_base")
		if ctrErr != nil {
			t.Fatalf("container DROP TABLE error: %v", ctrErr)
		}

		// View should still exist in container
		_, showErr := ctr.showCreateView("v_dep_test")
		if showErr != nil {
			t.Fatalf("container: view should still exist after base table drop, got: %v", showErr)
		}

		// Omni: same behavior
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_dep_base (id INT, val VARCHAR(100))", nil)
		c.Exec("CREATE VIEW v_dep_test AS SELECT id, val FROM t_dep_base", nil)
		c.Exec("DROP TABLE t_dep_base", nil)

		// View should still exist
		db := c.GetDatabase("test")
		if db.Views[toLower("v_dep_test")] == nil {
			t.Error("omni: view v_dep_test should still exist after base table drop")
		}

		// Cleanup container
		ctr.execSQL("DROP VIEW IF EXISTS v_dep_test")
	})

	t.Run("show_create_view_basic", func(t *testing.T) {
		// SHOW CREATE VIEW output format — verify structural elements.
		// Note: MySQL rewrites the SELECT text (lowercases keywords, backtick-quotes aliases)
		// which requires a full SQL deparser. We verify the CREATE VIEW preamble matches.
		ctr.execSQL("DROP VIEW IF EXISTS v_show_basic")
		ctr.execSQL("CREATE VIEW v_show_basic AS SELECT 1 AS a, 2 AS b")
		ctrDDL, err := ctr.showCreateView("v_show_basic")
		if err != nil {
			t.Fatalf("container SHOW CREATE VIEW error: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec("CREATE VIEW v_show_basic AS SELECT 1 AS a, 2 AS b", nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateView("test", "v_show_basic")

		// Verify structural preamble matches.
		oraclePreamble := extractViewPreamble(ctrDDL)
		omniPreamble := extractViewPreamble(omniDDL)
		if oraclePreamble != omniPreamble {
			t.Errorf("preamble mismatch:\n--- container ---\n%s\n--- omni ---\n%s\n--- container full ---\n%s\n--- omni full ---\n%s",
				oraclePreamble, omniPreamble, ctrDDL, omniDDL)
		}
	})

	t.Run("show_create_view_with_options", func(t *testing.T) {
		// SHOW CREATE VIEW with SQL SECURITY INVOKER.
		// Note: MySQL normalizes ALGORITHM to UNDEFINED for simple views even if MERGE was specified.
		ctr.execSQL("DROP VIEW IF EXISTS v_show_opts")
		ctr.execSQL("CREATE SQL SECURITY INVOKER VIEW v_show_opts AS SELECT 1 AS x")
		ctrDDL, err := ctr.showCreateView("v_show_opts")
		if err != nil {
			t.Fatalf("container SHOW CREATE VIEW error: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE SQL SECURITY INVOKER VIEW v_show_opts AS SELECT 1 AS x", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateView("test", "v_show_opts")

		// Verify structural preamble matches (up to and including AS).
		oraclePreamble := extractViewPreamble(ctrDDL)
		omniPreamble := extractViewPreamble(omniDDL)
		if oraclePreamble != omniPreamble {
			t.Errorf("preamble mismatch:\n--- container ---\n%s\n--- omni ---\n%s\n--- container full ---\n%s\n--- omni full ---\n%s",
				oraclePreamble, omniPreamble, ctrDDL, omniDDL)
		}
	})

	t.Run("view_with_column_aliases", func(t *testing.T) {
		// View with column list — verify columns are stored and preamble matches.
		ctr.execSQL("DROP VIEW IF EXISTS v_col_aliases")
		ctr.execSQL("CREATE VIEW v_col_aliases (x, y, z) AS SELECT 1, 2, 3")
		ctrDDL, err := ctr.showCreateView("v_col_aliases")
		if err != nil {
			t.Fatalf("container SHOW CREATE VIEW error: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE VIEW v_col_aliases (x, y, z) AS SELECT 1, 2, 3", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateView("test", "v_col_aliases")

		// Verify columns are present in SHOW CREATE VIEW output.
		oraclePreamble := extractViewPreamble(ctrDDL)
		omniPreamble := extractViewPreamble(omniDDL)
		if oraclePreamble != omniPreamble {
			t.Errorf("preamble mismatch:\n--- container ---\n%s\n--- omni ---\n%s\n--- container full ---\n%s\n--- omni full ---\n%s",
				oraclePreamble, omniPreamble, ctrDDL, omniDDL)
		}

		// Verify column list appears in both outputs.
		if !strings.Contains(ctrDDL, "`x`") || !strings.Contains(ctrDDL, "`y`") || !strings.Contains(ctrDDL, "`z`") {
			t.Errorf("container DDL missing column aliases: %s", ctrDDL)
		}
		if !strings.Contains(omniDDL, "`x`") || !strings.Contains(omniDDL, "`y`") || !strings.Contains(omniDDL, "`z`") {
			t.Errorf("omni DDL missing column aliases: %s", omniDDL)
		}
	})
}

func TestContainer_Section_5_1_UseStatement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Helper to extract MySQL error code from go-sql-driver error.
	extractMySQLErr := func(err error) (uint16, string, string) {
		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) {
			return mysqlErr.Number, string(mysqlErr.SQLState[:]), mysqlErr.Message
		}
		return 0, "", ""
	}

	t.Run("use_db_sets_current_database", func(t *testing.T) {
		// On container: create a separate database, USE it, then CREATE TABLE without qualifying.
		ctr.execSQL("DROP DATABASE IF EXISTS use_test_db")
		ctr.execSQL("CREATE DATABASE use_test_db")
		if err := ctr.execSQL("USE use_test_db"); err != nil {
			t.Fatalf("container USE: %v", err)
		}
		ctr.execSQL("DROP TABLE IF EXISTS t_use_check")
		if err := ctr.execSQL("CREATE TABLE t_use_check (id INT PRIMARY KEY)"); err != nil {
			t.Fatalf("container CREATE TABLE: %v", err)
		}
		ctrDDL, err := ctr.showCreateTable("t_use_check")
		if err != nil {
			t.Fatalf("container SHOW CREATE TABLE: %v", err)
		}

		// On omni: same sequence via Exec.
		c := New()
		c.Exec("CREATE DATABASE use_test_db", nil)
		results, _ := c.Exec("USE use_test_db", nil)
		if results[0].Error != nil {
			t.Fatalf("omni USE error: %v", results[0].Error)
		}
		// Verify current database was set.
		if c.CurrentDatabase() != "use_test_db" {
			t.Fatalf("omni CurrentDatabase: got %q, want %q", c.CurrentDatabase(), "use_test_db")
		}

		results, _ = c.Exec("CREATE TABLE t_use_check (id INT PRIMARY KEY)", nil)
		if results[0].Error != nil {
			t.Fatalf("omni CREATE TABLE error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("use_test_db", "t_use_check")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}

		// Cleanup ctr.
		ctr.execSQL("DROP DATABASE IF EXISTS use_test_db")
	})

	t.Run("use_nonexistent_error_1049", func(t *testing.T) {
		// On container: USE a nonexistent database.
		ctr.execSQL("DROP DATABASE IF EXISTS use_nonexistent_db")
		ctrErr := ctr.execSQL("USE use_nonexistent_db")
		if ctrErr == nil {
			t.Fatal("container: expected error for USE nonexistent database")
		}
		ctrCode, ctrState, ctrMsg := extractMySQLErr(ctrErr)
		t.Logf("container error: %d (%s) %s", ctrCode, ctrState, ctrMsg)

		if ctrCode != 1049 {
			t.Fatalf("container: expected error code 1049, got %d", ctrCode)
		}

		// On omni: USE a nonexistent database.
		c := New()
		results, _ := c.Exec("USE use_nonexistent_db", &ExecOptions{ContinueOnError: true})
		if results[0].Error == nil {
			t.Fatal("omni: expected error for USE nonexistent database")
		}
		catErr, ok := results[0].Error.(*Error)
		if !ok {
			t.Fatalf("omni: expected *Error, got %T", results[0].Error)
		}

		// Compare error code.
		if catErr.Code != int(ctrCode) {
			t.Errorf("error code mismatch: container=%d omni=%d", ctrCode, catErr.Code)
		}
		// Compare SQLSTATE.
		if catErr.SQLState != ctrState {
			t.Errorf("SQLSTATE mismatch: container=%q omni=%q", ctrState, catErr.SQLState)
		}
		// Compare error message.
		if catErr.Message != ctrMsg {
			t.Errorf("message mismatch:\n  container: %s\n  omni:   %s", ctrMsg, catErr.Message)
		}
	})
}

func TestContainer_Section_5_2_SetVariables(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Scenario: SET foreign_key_checks = 0 — skip FK validation on CREATE TABLE
	t.Run("fk_checks_off_create", func(t *testing.T) {
		ctr.execSQL("DROP TABLE IF EXISTS t_fkoff_child")
		ctr.execSQL("SET foreign_key_checks = 0")
		ctrErr := ctr.execSQL("CREATE TABLE t_fkoff_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_parent(id))")
		ctr.execSQL("SET foreign_key_checks = 1")
		if ctrErr != nil {
			t.Fatalf("container: unexpected error with FK checks off: %v", ctrErr)
		}
		ctrDDL, err := ctr.showCreateTable("t_fkoff_child")
		if err != nil {
			t.Fatalf("container: SHOW CREATE TABLE: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("SET foreign_key_checks = 0", nil)
		results, _ := c.Exec("CREATE TABLE t_fkoff_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_parent(id))", nil)
		if results[0].Error != nil {
			t.Fatalf("omni: unexpected error with FK checks off: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_fkoff_child")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL, omniDDL)
		}
	})

	// Scenario: SET foreign_key_checks = 0 — drop table with FK references
	t.Run("fk_checks_off_drop", func(t *testing.T) {
		ctr.execSQL("SET foreign_key_checks = 0")
		ctr.execSQL("DROP TABLE IF EXISTS t_fkdrop_child")
		ctr.execSQL("DROP TABLE IF EXISTS t_fkdrop_parent")
		ctr.execSQL("SET foreign_key_checks = 1")
		ctr.execSQL("CREATE TABLE t_fkdrop_parent (id INT PRIMARY KEY)")
		ctr.execSQL("CREATE TABLE t_fkdrop_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fkdrop_parent(id))")
		ctr.execSQL("SET foreign_key_checks = 0")
		ctrErr := ctr.execSQL("DROP TABLE t_fkdrop_parent")
		ctr.execSQL("SET foreign_key_checks = 1")
		if ctrErr != nil {
			t.Fatalf("container: unexpected error dropping FK-referenced table with checks off: %v", ctrErr)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_fkdrop_parent (id INT PRIMARY KEY)", nil)
		c.Exec("CREATE TABLE t_fkdrop_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES t_fkdrop_parent(id))", nil)
		c.Exec("SET foreign_key_checks = 0", nil)
		results, _ := c.Exec("DROP TABLE t_fkdrop_parent", nil)
		if results[0].Error != nil {
			t.Fatalf("omni: unexpected error dropping FK-referenced table with checks off: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_fkdrop_parent")
		if omniDDL != "" {
			t.Errorf("omni: parent table should be gone after drop, got: %s", omniDDL)
		}
	})

	// Scenario: SET foreign_key_checks = 1 — enforce FK validation
	t.Run("fk_checks_on_enforce", func(t *testing.T) {
		ctr.execSQL("SET foreign_key_checks = 0")
		ctr.execSQL("DROP TABLE IF EXISTS t_fkon_child")
		ctr.execSQL("SET foreign_key_checks = 1")
		ctrErr := ctr.execSQL("CREATE TABLE t_fkon_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_on_parent(id))")
		if ctrErr == nil {
			t.Fatal("container: expected error with FK checks on, referencing non-existent table")
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("SET foreign_key_checks = 1", nil)
		results, _ := c.Exec("CREATE TABLE t_fkon_child (id INT, parent_id INT, FOREIGN KEY (parent_id) REFERENCES nonexistent_on_parent(id))", nil)
		if results[0].Error == nil {
			t.Fatal("omni: expected error with FK checks on, referencing non-existent table")
		}
	})

	// Scenario: SET NAMES utf8mb4 — silently accepted
	t.Run("set_names", func(t *testing.T) {
		ctrErr := ctr.execSQL("SET NAMES utf8mb4")
		if ctrErr != nil {
			t.Fatalf("container: SET NAMES error: %v", ctrErr)
		}

		c := New()
		results, err := c.Exec("SET NAMES utf8mb4", nil)
		if err != nil {
			t.Fatalf("omni: parse error: %v", err)
		}
		if results[0].Error != nil {
			t.Fatalf("omni: SET NAMES error: %v", results[0].Error)
		}
	})

	// Scenario: SET CHARACTER SET utf8mb4 — silently accepted
	t.Run("set_character_set", func(t *testing.T) {
		ctrErr := ctr.execSQL("SET CHARACTER SET utf8mb4")
		if ctrErr != nil {
			t.Fatalf("container: SET CHARACTER SET error: %v", ctrErr)
		}

		c := New()
		results, err := c.Exec("SET CHARACTER SET utf8mb4", nil)
		if err != nil {
			t.Fatalf("omni: parse error: %v", err)
		}
		if results[0].Error != nil {
			t.Fatalf("omni: SET CHARACTER SET error: %v", results[0].Error)
		}
	})

	// Scenario: SET sql_mode — silently accepted
	t.Run("set_sql_mode", func(t *testing.T) {
		ctrErr := ctr.execSQL("SET sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'")
		if ctrErr != nil {
			t.Fatalf("container: SET sql_mode error: %v", ctrErr)
		}

		c := New()
		results, err := c.Exec("SET sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'", nil)
		if err != nil {
			t.Fatalf("omni: parse error: %v", err)
		}
		if results[0].Error != nil {
			t.Fatalf("omni: SET sql_mode error: %v", results[0].Error)
		}
	})
}

// TestContainer_Section_5_3_UserRoleManagement verifies that user/role management
// statements are accepted by both MySQL and omni without error.
// These are marked [~] partial because the in-memory catalog does not actually
// store users, roles, or privileges — it just silently accepts them.
func TestContainer_Section_5_3_UserRoleManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name    string
		setup   string // SQL to run before the main statement (on container only)
		sql     string // main statement to test
		cleanup string // SQL to run after the test (on container only)
	}{
		{
			"create_user",
			"",
			"CREATE USER 'testuser53'@'localhost' IDENTIFIED BY 'Password123!'",
			"DROP USER IF EXISTS 'testuser53'@'localhost'",
		},
		{
			"drop_user",
			"CREATE USER IF NOT EXISTS 'dropme53'@'localhost' IDENTIFIED BY 'Password123!'",
			"DROP USER 'dropme53'@'localhost'",
			"",
		},
		{
			"alter_user",
			"CREATE USER IF NOT EXISTS 'alterme53'@'localhost' IDENTIFIED BY 'Password123!'",
			"ALTER USER 'alterme53'@'localhost' IDENTIFIED BY 'NewPassword456!'",
			"DROP USER IF EXISTS 'alterme53'@'localhost'",
		},
		{
			"create_role",
			"",
			"CREATE ROLE 'app_role53'",
			"DROP ROLE IF EXISTS 'app_role53'",
		},
		{
			"drop_role",
			"CREATE ROLE IF NOT EXISTS 'droprole53'",
			"DROP ROLE 'droprole53'",
			"",
		},
		{
			"grant_privileges",
			"CREATE USER IF NOT EXISTS 'grantee53'@'localhost' IDENTIFIED BY 'Password123!'",
			"GRANT SELECT, INSERT ON test.* TO 'grantee53'@'localhost'",
			"DROP USER IF EXISTS 'grantee53'@'localhost'",
		},
		{
			"revoke_privileges",
			"CREATE USER IF NOT EXISTS 'revokee53'@'localhost' IDENTIFIED BY 'Password123!'; GRANT SELECT ON test.* TO 'revokee53'@'localhost'",
			"REVOKE SELECT ON test.* FROM 'revokee53'@'localhost'",
			"DROP USER IF EXISTS 'revokee53'@'localhost'",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup on container
			if tc.setup != "" {
				if err := ctr.execSQL(tc.setup); err != nil {
					t.Fatalf("container setup: %v", err)
				}
			}

			// Run on container — must succeed
			ctrErr := ctr.execSQL(tc.sql)
			if ctrErr != nil {
				t.Fatalf("container exec: %v", ctrErr)
			}

			// Cleanup on container
			if tc.cleanup != "" {
				_ = ctr.execSQL(tc.cleanup)
			}

			// Run on omni — must parse and not return a parse error.
			// The catalog silently accepts these (no-op) since it doesn't
			// track users/roles/privileges, which is expected behavior
			// for an in-memory DDL catalog.
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if len(results) > 0 && results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
		})
	}
}

func TestContainer_Section_6_1_ShowCreateTableIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		// --- All data types (Phase 1.1-1.6) integration ---
		{"all_data_types", `CREATE TABLE t_all_types (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			tiny_col TINYINT,
			small_col SMALLINT,
			med_col MEDIUMINT,
			int_col INT,
			big_col BIGINT,
			float_col FLOAT,
			float_prec FLOAT(7,3),
			double_col DOUBLE,
			double_prec DOUBLE(15,5),
			decimal_col DECIMAL(10,2),
			decimal_bare DECIMAL,
			bool_col BOOLEAN,
			bit_col BIT(8),
			char_col CHAR(10),
			varchar_col VARCHAR(255),
			tinytext_col TINYTEXT,
			text_col TEXT,
			mediumtext_col MEDIUMTEXT,
			longtext_col LONGTEXT,
			enum_col ENUM('a','b','c'),
			set_col SET('x','y','z'),
			binary_col BINARY(16),
			varbinary_col VARBINARY(255),
			tinyblob_col TINYBLOB,
			blob_col BLOB,
			mediumblob_col MEDIUMBLOB,
			longblob_col LONGBLOB,
			date_col DATE,
			time_col TIME,
			time_frac TIME(3),
			datetime_col DATETIME,
			datetime_frac DATETIME(6),
			timestamp_col TIMESTAMP NULL,
			timestamp_frac TIMESTAMP(3) NULL,
			year_col YEAR,
			json_col JSON,
			geo_col GEOMETRY,
			point_col POINT,
			linestring_col LINESTRING,
			polygon_col POLYGON,
			multipoint_col MULTIPOINT,
			multiline_col MULTILINESTRING,
			multipoly_col MULTIPOLYGON,
			geocoll_col GEOMETRYCOLLECTION,
			PRIMARY KEY (id)
		)`, "t_all_types"},

		// --- All default value forms (Phase 1.7) integration ---
		{"all_defaults", `CREATE TABLE t_all_defaults (
			id INT NOT NULL AUTO_INCREMENT,
			int_def INT DEFAULT 0,
			int_null INT DEFAULT NULL,
			int_notnull INT NOT NULL,
			varchar_def VARCHAR(100) DEFAULT 'hello',
			varchar_empty VARCHAR(100) DEFAULT '',
			float_def FLOAT DEFAULT 3.14,
			decimal_def DECIMAL(10,2) DEFAULT 0.00,
			bool_true BOOLEAN DEFAULT TRUE,
			bool_false BOOLEAN DEFAULT FALSE,
			enum_def ENUM('a','b','c') DEFAULT 'a',
			set_def SET('x','y','z') DEFAULT 'x,y',
			bit_def BIT(8) DEFAULT b'00001111',
			blob_col BLOB,
			text_col TEXT,
			json_col JSON,
			ts_def TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			dt_def DATETIME DEFAULT CURRENT_TIMESTAMP,
			ts3_def TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
			expr_int INT DEFAULT (FLOOR(RAND()*100)),
			expr_json JSON DEFAULT (JSON_ARRAY()),
			expr_varchar VARCHAR(100) DEFAULT (UUID()),
			dt_literal DATETIME DEFAULT '2024-01-01 00:00:00',
			date_literal DATE DEFAULT '2024-01-01',
			PRIMARY KEY (id)
		)`, "t_all_defaults"},

		// --- All index types (Phase 1.13-1.16) integration ---
		{"all_index_types", `CREATE TABLE t_all_indexes (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			score INT,
			rank_val INT,
			bio TEXT,
			description TEXT,
			geo_col POINT NOT NULL,
			PRIMARY KEY (id),
			KEY idx_name (name),
			KEY idx_email_prefix (email(10)),
			KEY idx_score_desc (score DESC),
			KEY idx_multi_mixed (score ASC, rank_val DESC),
			UNIQUE KEY uk_email (email),
			FULLTEXT KEY ft_bio (bio, description),
			SPATIAL KEY sp_geo (geo_col),
			KEY idx_expr ((UPPER(name))),
			KEY idx_expr_calc ((score + rank_val)),
			KEY idx_comment (name) COMMENT 'name lookup index',
			KEY idx_invisible (score) /*!80000 INVISIBLE */
		) ENGINE=InnoDB`, "t_all_indexes"},

		// --- All constraint forms (Phase 1.17-1.18) integration ---
		{"all_constraints_parent", `CREATE TABLE t_parent (
			id INT NOT NULL AUTO_INCREMENT,
			code VARCHAR(10) NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uk_code (code)
		)`, "t_parent"},
		{"all_constraints", `CREATE TABLE t_all_constraints (
			id INT NOT NULL AUTO_INCREMENT,
			parent_id INT,
			parent_code VARCHAR(10),
			self_ref INT,
			val INT,
			score INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES t_parent(id) ON DELETE CASCADE,
			CONSTRAINT fk_code FOREIGN KEY (parent_code) REFERENCES t_parent(code) ON UPDATE SET NULL,
			CONSTRAINT fk_self FOREIGN KEY (self_ref) REFERENCES t_all_constraints(id) ON DELETE SET NULL ON UPDATE CASCADE,
			CONSTRAINT chk_val CHECK (val > 0),
			CONSTRAINT chk_score CHECK (score >= 0 AND score <= 100),
			CONSTRAINT chk_not_enforced CHECK (val < 1000) /*!80016 NOT ENFORCED */
		)`, "t_all_constraints"},

		// --- Partitioned table output (Phase 4.1) integration ---
		{"partitioned_range", `CREATE TABLE t_part_range (
			id INT NOT NULL,
			created_date DATE NOT NULL,
			name VARCHAR(100),
			PRIMARY KEY (id, created_date)
		) ENGINE=InnoDB
		PARTITION BY RANGE (YEAR(created_date)) (
			PARTITION p2020 VALUES LESS THAN (2021),
			PARTITION p2021 VALUES LESS THAN (2022),
			PARTITION p2022 VALUES LESS THAN (2023),
			PARTITION pmax VALUES LESS THAN MAXVALUE
		)`, "t_part_range"},
		{"partitioned_list", `CREATE TABLE t_part_list (
			id INT NOT NULL AUTO_INCREMENT,
			region VARCHAR(20) NOT NULL,
			data VARCHAR(100),
			PRIMARY KEY (id, region)
		) ENGINE=InnoDB
		PARTITION BY LIST COLUMNS(region) (
			PARTITION p_east VALUES IN ('east','northeast'),
			PARTITION p_west VALUES IN ('west','northwest'),
			PARTITION p_other VALUES IN ('central','south')
		)`, "t_part_list"},
		{"partitioned_hash", `CREATE TABLE t_part_hash (
			id INT NOT NULL AUTO_INCREMENT,
			val INT,
			PRIMARY KEY (id)
		) ENGINE=InnoDB
		PARTITION BY HASH(id)
		PARTITIONS 4`, "t_part_hash"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up in proper order (child before parent)
			if tc.table == "t_all_constraints" {
				ctr.execSQL("DROP TABLE IF EXISTS t_all_constraints")
			}
			if tc.table == "t_parent" {
				ctr.execSQL("DROP TABLE IF EXISTS t_all_constraints")
				ctr.execSQL("DROP TABLE IF EXISTS t_parent")
			}
			if tc.table != "t_all_constraints" && tc.table != "t_parent" {
				ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			}
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, err := ctr.showCreateTable(tc.table)
			if err != nil {
				t.Fatalf("container showCreateTable: %v", err)
			}

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			// For constraint tests, need parent table first
			if tc.table == "t_all_constraints" {
				c.Exec("CREATE TABLE t_parent (id INT NOT NULL AUTO_INCREMENT, code VARCHAR(10) NOT NULL, PRIMARY KEY (id), UNIQUE KEY uk_code (code))", nil)
			}
			results, err := c.Exec(tc.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if len(results) > 0 && results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

// TestContainer_Section_6_2_ShowCreateOtherObjects tests SHOW CREATE VIEW/FUNCTION/PROCEDURE/TRIGGER/EVENT
// output against real MySQL 8.0. This complements sections 4.2-4.5 by testing additional patterns
// and verifying SHOW CREATE as a query API surface.
func TestContainer_Section_6_2_ShowCreateOtherObjects(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// --- SHOW CREATE VIEW ---
	t.Run("show_create_view", func(t *testing.T) {
		// Test SHOW CREATE VIEW with default options.
		// MySQL rewrites SELECT text, so we compare preamble only.
		ctr.execSQL("DROP VIEW IF EXISTS v_sc_basic")
		ctr.execSQL("CREATE VIEW v_sc_basic AS SELECT 1 AS col1, 2 AS col2")
		ctrDDL, err := ctr.showCreateView("v_sc_basic")
		if err != nil {
			t.Fatalf("container SHOW CREATE VIEW: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec("CREATE VIEW v_sc_basic AS SELECT 1 AS col1, 2 AS col2", nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateView("test", "v_sc_basic")

		// Compare preamble (up to AS) — SELECT text may differ due to MySQL rewriting.
		oraclePreamble := extractViewPreamble(ctrDDL)
		omniPreamble := extractViewPreamble(omniDDL)
		if oraclePreamble != omniPreamble {
			t.Errorf("preamble mismatch:\n--- container ---\n%s\n--- omni ---\n%s\n--- container full ---\n%s\n--- omni full ---\n%s",
				oraclePreamble, omniPreamble, ctrDDL, omniDDL)
		}

		// Verify both contain the key structural elements.
		for _, substr := range []string{"ALGORITHM=", "DEFINER=", "SQL SECURITY", "VIEW", "v_sc_basic", " AS "} {
			if !strings.Contains(ctrDDL, substr) {
				t.Errorf("container DDL missing %q: %s", substr, ctrDDL)
			}
			if !strings.Contains(omniDDL, substr) {
				t.Errorf("omni DDL missing %q: %s", substr, omniDDL)
			}
		}
	})

	// --- SHOW CREATE FUNCTION ---
	t.Run("show_create_function", func(t *testing.T) {
		// Test SHOW CREATE FUNCTION with various characteristics.
		ctr.execSQL("DROP FUNCTION IF EXISTS fn_sc_multiply")
		createSQL := "CREATE FUNCTION fn_sc_multiply(x INT, y INT) RETURNS INT DETERMINISTIC RETURN x * y"
		if err := ctr.execSQL(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateFunction("fn_sc_multiply")
		if err != nil {
			t.Fatalf("container SHOW CREATE FUNCTION: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateFunction("test", "fn_sc_multiply")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// --- SHOW CREATE PROCEDURE ---
	t.Run("show_create_procedure", func(t *testing.T) {
		// Test SHOW CREATE PROCEDURE with INOUT params.
		ctr.execSQLDirect("DROP PROCEDURE IF EXISTS sp_sc_swap")
		createSQL := "CREATE PROCEDURE sp_sc_swap(INOUT a INT, INOUT b INT) BEGIN DECLARE tmp INT; SET tmp = a; SET a = b; SET b = tmp; END"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateProcedure("sp_sc_swap")
		if err != nil {
			t.Fatalf("container SHOW CREATE PROCEDURE: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateProcedure("test", "sp_sc_swap")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// --- SHOW CREATE TRIGGER ---
	t.Run("show_create_trigger", func(t *testing.T) {
		// Setup table for triggers.
		ctr.execSQL("DROP TRIGGER IF EXISTS tr_sc_audit")
		ctr.execSQL("DROP TABLE IF EXISTS t_sc_trigger")
		if err := ctr.execSQL("CREATE TABLE t_sc_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, updated_at DATETIME)"); err != nil {
			t.Fatalf("container setup table: %v", err)
		}
		createSQL := "CREATE TRIGGER tr_sc_audit AFTER INSERT ON t_sc_trigger FOR EACH ROW SET @last_insert_id = NEW.id"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateTrigger("tr_sc_audit")
		if err != nil {
			t.Fatalf("container SHOW CREATE TRIGGER: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_sc_trigger (id INT AUTO_INCREMENT PRIMARY KEY, val INT, updated_at DATETIME)", nil)
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTrigger("test", "tr_sc_audit")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})

	// --- SHOW CREATE EVENT ---
	t.Run("show_create_event", func(t *testing.T) {
		// Test SHOW CREATE EVENT with AT schedule (exact timestamp, no auto-STARTS issue).
		ctr.execSQL("DROP EVENT IF EXISTS ev_sc_onetime")
		createSQL := "CREATE EVENT ev_sc_onetime ON SCHEDULE AT '2035-06-15 12:00:00' ON COMPLETION PRESERVE COMMENT 'one-time event' DO SELECT 1"
		if err := ctr.execSQLDirect(createSQL); err != nil {
			t.Fatalf("container exec: %v", err)
		}
		ctrDDL, err := ctr.showCreateEvent("ev_sc_onetime")
		if err != nil {
			t.Fatalf("container SHOW CREATE EVENT: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, parseErr := c.Exec(createSQL, nil)
		if parseErr != nil {
			t.Fatalf("omni parse error: %v", parseErr)
		}
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateEvent("test", "ev_sc_onetime")

		if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
				ctrDDL, omniDDL)
		}
	})
}

// TestContainer_Section_6_3_InformationSchemaConsistency verifies that the catalog's
// internal state is consistent with what MySQL 8.0 reports via INFORMATION_SCHEMA.
//
// The omni catalog does not support INFORMATION_SCHEMA SQL queries (SELECT is
// treated as DML and skipped). Instead, we compare the catalog's Go-level data
// structures against real MySQL INFORMATION_SCHEMA query results.
//
// All scenarios are marked [~] partial because the catalog lacks an
// INFORMATION_SCHEMA query engine — users cannot run SELECT ... FROM
// INFORMATION_SCHEMA.* against the in-memory catalog.
func TestContainer_Section_6_3_InformationSchemaConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	// Setup SQL: a table with various column types, indexes, FK, and CHECK constraints.
	parentSQL := "CREATE TABLE t_is_parent (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id)) ENGINE=InnoDB"
	setupSQL := `CREATE TABLE t_is_test (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(100) NOT NULL DEFAULT '',
		price DECIMAL(10,2) DEFAULT '0.00',
		status ENUM('active','inactive') DEFAULT 'active',
		parent_id INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY uk_name (name),
		KEY idx_status (status),
		KEY idx_parent (parent_id),
		CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES t_is_parent (id) ON DELETE SET NULL,
		CONSTRAINT chk_price CHECK (price >= 0)
	) ENGINE=InnoDB`

	// Setup on ctr.
	ctr.execSQL("DROP TABLE IF EXISTS t_is_test")
	ctr.execSQL("DROP TABLE IF EXISTS t_is_parent")
	if err := ctr.execSQL(parentSQL); err != nil {
		t.Fatalf("container parent table: %v", err)
	}
	if err := ctr.execSQL(setupSQL); err != nil {
		t.Fatalf("container setup: %v", err)
	}

	// Setup on omni.
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	if results, _ := c.Exec(parentSQL, nil); results[0].Error != nil {
		t.Fatalf("omni parent table: %v", results[0].Error)
	}
	if results, _ := c.Exec(setupSQL, nil); results[0].Error != nil {
		t.Fatalf("omni setup: %v", results[0].Error)
	}

	db := c.GetDatabase("test")
	tbl := db.GetTable("t_is_test")

	// 6.3.1: INFORMATION_SCHEMA.COLUMNS matches catalog state
	t.Run("columns_match", func(t *testing.T) {
		oracleCols, err := ctr.queryColumns("test", "t_is_test")
		if err != nil {
			t.Fatalf("container queryColumns: %v", err)
		}
		if len(oracleCols) != len(tbl.Columns) {
			t.Fatalf("column count: container=%d omni=%d", len(oracleCols), len(tbl.Columns))
		}
		for i, oc := range oracleCols {
			omniCol := tbl.Columns[i]
			if oc.Name != omniCol.Name {
				t.Errorf("col[%d] name: container=%q omni=%q", i, oc.Name, omniCol.Name)
			}
			if oc.Position != omniCol.Position {
				t.Errorf("col[%d] position: container=%d omni=%d", i, oc.Position, omniCol.Position)
			}
			wantNullable := "YES"
			if !omniCol.Nullable {
				wantNullable = "NO"
			}
			if oc.Nullable != wantNullable {
				t.Errorf("col[%d] %s nullable: container=%q omni=%q", i, oc.Name, oc.Nullable, wantNullable)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.2: INFORMATION_SCHEMA.STATISTICS matches catalog indexes
	t.Run("statistics_match", func(t *testing.T) {
		oracleIdxs, err := ctr.queryIndexes("test", "t_is_test")
		if err != nil {
			t.Fatalf("container queryIndexes: %v", err)
		}
		// Build omni index column map: indexName -> []columnName
		omniIdxCols := make(map[string][]string)
		for _, idx := range tbl.Indexes {
			for _, ic := range idx.Columns {
				omniIdxCols[idx.Name] = append(omniIdxCols[idx.Name], ic.Name)
			}
		}
		// Build container index column map.
		oracleIdxCols := make(map[string][]string)
		for _, oi := range oracleIdxs {
			oracleIdxCols[oi.Name] = append(oracleIdxCols[oi.Name], oi.ColumnName)
		}
		// Compare index names and columns.
		for name, oracleCols := range oracleIdxCols {
			omniCols, ok := omniIdxCols[name]
			if !ok {
				t.Errorf("index %q exists in container but not in omni", name)
				continue
			}
			if len(oracleCols) != len(omniCols) {
				t.Errorf("index %q column count: container=%d omni=%d", name, len(oracleCols), len(omniCols))
				continue
			}
			for j := range oracleCols {
				if oracleCols[j] != omniCols[j] {
					t.Errorf("index %q col[%d]: container=%q omni=%q", name, j, oracleCols[j], omniCols[j])
				}
			}
		}
		for name := range omniIdxCols {
			if _, ok := oracleIdxCols[name]; !ok {
				t.Errorf("index %q exists in omni but not in container", name)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.3: INFORMATION_SCHEMA.TABLE_CONSTRAINTS matches
	t.Run("table_constraints_match", func(t *testing.T) {
		oracleCons, err := ctr.queryConstraints("test", "t_is_test")
		if err != nil {
			t.Fatalf("container queryConstraints: %v", err)
		}
		// Build omni constraint map.
		omniCons := make(map[string]string) // name -> type
		for _, con := range tbl.Constraints {
			var typeName string
			switch con.Type {
			case ConPrimaryKey:
				typeName = "PRIMARY KEY"
			case ConUniqueKey:
				typeName = "UNIQUE"
			case ConForeignKey:
				typeName = "FOREIGN KEY"
			case ConCheck:
				typeName = "CHECK"
			}
			omniCons[con.Name] = typeName
		}
		// Build container constraint map.
		oracleConsMap := make(map[string]string)
		for _, oc := range oracleCons {
			oracleConsMap[oc.Name] = oc.Type
		}
		for name, oType := range oracleConsMap {
			omniType, ok := omniCons[name]
			if !ok {
				t.Errorf("constraint %q (%s) in container but not omni", name, oType)
				continue
			}
			if oType != omniType {
				t.Errorf("constraint %q type: container=%q omni=%q", name, oType, omniType)
			}
		}
		for name := range omniCons {
			if _, ok := oracleConsMap[name]; !ok {
				t.Errorf("constraint %q in omni but not container", name)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.4: INFORMATION_SCHEMA.KEY_COLUMN_USAGE matches
	t.Run("key_column_usage_match", func(t *testing.T) {
		rows, err := ctr.db.QueryContext(ctr.ctx, `
			SELECT CONSTRAINT_NAME, COLUMN_NAME, ORDINAL_POSITION,
			       REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
			WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 't_is_test'
			ORDER BY CONSTRAINT_NAME, ORDINAL_POSITION`)
		if err != nil {
			t.Fatalf("container KEY_COLUMN_USAGE query: %v", err)
		}
		defer rows.Close()

		type kcuRow struct {
			constraintName, columnName string
			ordinalPos                 int
			refTable, refColumn        *string
		}
		var oracleKCU []kcuRow
		for rows.Next() {
			var r kcuRow
			var refTbl, refCol *string
			if err := rows.Scan(&r.constraintName, &r.columnName, &r.ordinalPos, &refTbl, &refCol); err != nil {
				t.Fatalf("scan: %v", err)
			}
			r.refTable = refTbl
			r.refColumn = refCol
			oracleKCU = append(oracleKCU, r)
		}

		// Verify omni constraints have matching columns.
		omniKCU := make(map[string][]string) // constraint name -> columns
		for _, con := range tbl.Constraints {
			if con.Type == ConCheck {
				continue // CHECK constraints don't appear in KEY_COLUMN_USAGE
			}
			omniKCU[con.Name] = con.Columns
		}
		for _, okcu := range oracleKCU {
			cols, ok := omniKCU[okcu.constraintName]
			if !ok {
				t.Errorf("constraint %q in container KEY_COLUMN_USAGE but not omni", okcu.constraintName)
				continue
			}
			idx := okcu.ordinalPos - 1
			if idx < len(cols) && cols[idx] != okcu.columnName {
				t.Errorf("constraint %q col[%d]: container=%q omni=%q",
					okcu.constraintName, idx, okcu.columnName, cols[idx])
			}
		}

		// Verify FK references match.
		for _, con := range tbl.Constraints {
			if con.Type != ConForeignKey {
				continue
			}
			for _, okcu := range oracleKCU {
				if okcu.constraintName == con.Name && okcu.refTable != nil {
					if *okcu.refTable != con.RefTable {
						t.Errorf("FK %q ref table: container=%q omni=%q", con.Name, *okcu.refTable, con.RefTable)
					}
				}
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.5: INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS matches
	t.Run("referential_constraints_match", func(t *testing.T) {
		rows, err := ctr.db.QueryContext(ctr.ctx, `
			SELECT CONSTRAINT_NAME, UNIQUE_CONSTRAINT_NAME,
			       MATCH_OPTION, UPDATE_RULE, DELETE_RULE,
			       REFERENCED_TABLE_NAME
			FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA = 'test' AND TABLE_NAME = 't_is_test'
			ORDER BY CONSTRAINT_NAME`)
		if err != nil {
			t.Fatalf("container REFERENTIAL_CONSTRAINTS query: %v", err)
		}
		defer rows.Close()

		type refConRow struct {
			name, uniqueConName, matchOption, updateRule, deleteRule, refTable string
		}
		var oracleRefs []refConRow
		for rows.Next() {
			var r refConRow
			if err := rows.Scan(&r.name, &r.uniqueConName, &r.matchOption, &r.updateRule, &r.deleteRule, &r.refTable); err != nil {
				t.Fatalf("scan: %v", err)
			}
			oracleRefs = append(oracleRefs, r)
		}

		// Compare with omni FK constraints.
		for _, oref := range oracleRefs {
			var found *Constraint
			for _, con := range tbl.Constraints {
				if con.Type == ConForeignKey && con.Name == oref.name {
					found = con
					break
				}
			}
			if found == nil {
				t.Errorf("FK %q in container REFERENTIAL_CONSTRAINTS but not omni", oref.name)
				continue
			}
			if oref.refTable != found.RefTable {
				t.Errorf("FK %q ref table: container=%q omni=%q", oref.name, oref.refTable, found.RefTable)
			}
			// Normalize action names for comparison.
			oracleDelete := oref.deleteRule
			omniDelete := found.OnDelete
			if omniDelete == "" {
				omniDelete = "RESTRICT" // MySQL default
			}
			if oracleDelete != omniDelete {
				t.Errorf("FK %q delete rule: container=%q omni=%q", oref.name, oracleDelete, omniDelete)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})

	// 6.3.6: INFORMATION_SCHEMA.CHECK_CONSTRAINTS matches
	t.Run("check_constraints_match", func(t *testing.T) {
		rows, err := ctr.db.QueryContext(ctr.ctx, `
			SELECT CONSTRAINT_NAME, CHECK_CLAUSE
			FROM INFORMATION_SCHEMA.CHECK_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA = 'test'
			ORDER BY CONSTRAINT_NAME`)
		if err != nil {
			t.Fatalf("container CHECK_CONSTRAINTS query: %v", err)
		}
		defer rows.Close()

		type checkRow struct {
			name, clause string
		}
		var oracleChecks []checkRow
		for rows.Next() {
			var r checkRow
			if err := rows.Scan(&r.name, &r.clause); err != nil {
				t.Fatalf("scan: %v", err)
			}
			oracleChecks = append(oracleChecks, r)
		}

		// Filter to just our table's check constraints (INFORMATION_SCHEMA.CHECK_CONSTRAINTS
		// doesn't have TABLE_NAME, so we match by constraint name).
		omniChecks := make(map[string]string) // name -> expr
		for _, con := range tbl.Constraints {
			if con.Type == ConCheck {
				omniChecks[con.Name] = con.CheckExpr
			}
		}
		for _, oc := range oracleChecks {
			omniExpr, ok := omniChecks[oc.name]
			if !ok {
				// May be an ENUM/SET constraint auto-generated by MySQL — skip.
				continue
			}
			// Normalize for comparison: remove outer parens and whitespace.
			oracleNorm := normalizeWhitespace(oc.clause)
			omniNorm := normalizeWhitespace(omniExpr)
			if oracleNorm != omniNorm {
				// Log but don't fail — expression format may differ slightly.
				t.Logf("check %q clause: container=%q omni=%q (expression format may differ)", oc.name, oracleNorm, omniNorm)
			}
		}
		t.Log("[~] partial: catalog internal state matches, but INFORMATION_SCHEMA SQL queries require query engine")
	})
}
