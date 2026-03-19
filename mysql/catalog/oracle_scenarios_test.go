package catalog

import (
	"testing"
)

func TestOracle_Section_1_2_StringTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_4_DateTimeTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_1_NumericTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			results, _ := c.Exec(tc.sql, nil)
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_10_ColumnAttributesCombination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_7_DefaultValues(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_11_PrimaryKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_13_RegularIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_17_ForeignKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL(tc.cleanup)

			// Setup parent tables.
			if tc.setup != "" {
				if err := oracle.execSQL(tc.setup); err != nil {
					t.Fatalf("oracle setup: %v", err)
				}
			}
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_12_UniqueKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_18_CheckConstraints(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_19_TableOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_8_OnUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_20_CharsetCollationInheritance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
				if err := oracle.execSQL(tc.setupSQL); err != nil {
					t.Fatalf("oracle setup: %v", err)
				}
			}
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, err := oracle.showCreateTable(tc.table)
			if err != nil {
				t.Fatalf("oracle show create: %v", err)
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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_9_GeneratedColumns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_3_BinaryTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_6_JSONType(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_5_SpatialTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_14_FulltextSpatialIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_15_ExpressionIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_1_16_IndexOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
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
			oracle.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := oracle.execSQL(tc.sql); err != nil {
				t.Fatalf("oracle exec: %v", err)
			}
			oracleDDL, _ := oracle.showCreateTable(tc.table)

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

			if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
					oracleDDL, omniDDL)
			}
		})
	}
}

func TestOracle_Section_2_1_CreateTableVariants(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle, cleanup := startOracle(t)
	defer cleanup()

	t.Run("if_not_exists_no_error", func(t *testing.T) {
		oracle.execSQL("DROP TABLE IF EXISTS t_ine")
		oracle.execSQL("CREATE TABLE t_ine (id INT)")

		// Second CREATE with IF NOT EXISTS should not error on oracle.
		oracleErr := oracle.execSQL("CREATE TABLE IF NOT EXISTS t_ine (id INT)")
		if oracleErr != nil {
			t.Fatalf("oracle error on IF NOT EXISTS: %v", oracleErr)
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
		oracleDDL, _ := oracle.showCreateTable("t_ine")
		omniDDL := c.ShowCreateTable("test", "t_ine")
		if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
				oracleDDL, omniDDL)
		}
	})

	t.Run("temporary_table", func(t *testing.T) {
		// MySQL SHOW CREATE TABLE for temporary tables shows "CREATE TEMPORARY TABLE".
		oracle.execSQL("DROP TEMPORARY TABLE IF EXISTS t_temp")
		err := oracle.execSQL("CREATE TEMPORARY TABLE t_temp (id INT, name VARCHAR(50))")
		if err != nil {
			t.Fatalf("oracle exec: %v", err)
		}
		oracleDDL, err := oracle.showCreateTable("t_temp")
		if err != nil {
			t.Fatalf("oracle show create: %v", err)
		}

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TEMPORARY TABLE t_temp (id INT, name VARCHAR(50))", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_temp")

		if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
				oracleDDL, omniDDL)
		}
	})

	t.Run("create_table_like", func(t *testing.T) {
		oracle.execSQL("DROP TABLE IF EXISTS t_like_dst")
		oracle.execSQL("DROP TABLE IF EXISTS t_like_src")
		oracle.execSQL("CREATE TABLE t_like_src (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT '', score DECIMAL(10,2))")
		err := oracle.execSQL("CREATE TABLE t_like_dst LIKE t_like_src")
		if err != nil {
			t.Fatalf("oracle exec: %v", err)
		}
		oracleDDL, _ := oracle.showCreateTable("t_like_dst")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t_like_src (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL DEFAULT '', score DECIMAL(10,2))", nil)
		results, _ := c.Exec("CREATE TABLE t_like_dst LIKE t_like_src", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "t_like_dst")

		if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
				oracleDDL, omniDDL)
		}
	})

	t.Run("create_table_with_view_name_conflict", func(t *testing.T) {
		// Creating a table with same name as an existing view should error.
		oracle.execSQL("DROP TABLE IF EXISTS t_view_conflict")
		oracle.execSQL("DROP VIEW IF EXISTS t_view_conflict")
		oracle.execSQL("CREATE VIEW t_view_conflict AS SELECT 1 AS a")
		oracleErr := oracle.execSQL("CREATE TABLE t_view_conflict (id INT)")
		if oracleErr == nil {
			t.Fatal("expected oracle error when creating table with same name as view")
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
		oracle.execSQL("DROP TABLE IF EXISTS `select`")
		err := oracle.execSQL("CREATE TABLE `select` (`from` INT, `where` VARCHAR(50))")
		if err != nil {
			t.Fatalf("oracle exec: %v", err)
		}
		oracleDDL, _ := oracle.showCreateTable("`select`")

		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		results, _ := c.Exec("CREATE TABLE `select` (`from` INT, `where` VARCHAR(50))", nil)
		if results[0].Error != nil {
			t.Fatalf("omni exec error: %v", results[0].Error)
		}
		omniDDL := c.ShowCreateTable("test", "select")

		if normalizeWhitespace(oracleDDL) != normalizeWhitespace(omniDDL) {
			t.Errorf("mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s",
				oracleDDL, omniDDL)
		}
	})
}
