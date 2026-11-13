package catalog

import "testing"

// TestAlterTableOptionText holds checkAlterTableOptions to values as OceanBase reads
// them, which it reads with the parser's own lexer: adjacent strings folded ('DA' 'Y'
// is DAY, which the parser's oracle tape TestOceanBaseTableOptionGrammarOracle records
// OceanBase taking), and every comment form read as the parser reads it, so what one
// holds separates nothing and an executable comment's contents are the value.
func TestAlterTableOptionText(t *testing.T) {
	c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE r (b DATETIME PRIMARY KEY) PARTITION BY RANGE COLUMNS (b) (PARTITION p VALUES LESS THAN ('2026-01-01'))")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (TIME_UNIT = 'DA' 'Y')", 0},
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (TIME_UNIT = 'DAY' /* = , */)", 0},
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (TIME_UNIT = 'DAY' TIME_UNIT = 'DAY')", 1210},
		{"ALTER TABLE r TABLE_MODE 'nor' 'mal'", 0},
		{"ALTER TABLE r TABLE_MODE 'no' 'pe'", 1235},
		{"ALTER TABLE r TABLE_MODE /*x*/ = 'normal'", 0},
		{"ALTER TABLE r STORAGE_FORMAT_VERSION /*x*/ = 4", 0},
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (TIME_UNIT = 'DAY' --\tcomment\n)", 0},
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (TIME_UNIT = 'DAY' # comment\n)", 0},
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (TIME_UNIT = /*! 'DAY' */)", 0},
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (/*!TIME_UNIT = 'DAY'*/)", 0},
		{"ALTER TABLE r DYNAMIC_PARTITION_POLICY (TIME_UNIT = 'DAY' /*, TIME_UNIT = 'DAY' */)", 0},
		{"ALTER TABLE r TABLE_MODE /*!'normal'*/", 0},
		{"ALTER TABLE r TABLE_MODE --\tx\n'normal'", 0},
		{"ALTER TABLE r STORAGE_FORMAT_VERSION /*!4*/", 0},
	} {
		if got := omniCode(c, tc.sql); got != tc.code {
			t.Errorf("%s: got %d, want %d", tc.sql, got, tc.code)
		}
	}
}
