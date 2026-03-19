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
