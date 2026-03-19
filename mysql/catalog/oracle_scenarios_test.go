package catalog

import "testing"

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
