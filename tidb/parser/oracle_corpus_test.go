package parser

import "testing"

func TestOracleCorpus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test in short mode")
	}
	oracle := startParserOracle(t)

	// Create a dummy table for SELECT statements that reference 't'.
	oracle.db.ExecContext(oracle.ctx, "CREATE TABLE IF NOT EXISTS t (id INT, val INT, score INT, status INT, name VARCHAR(100), dept VARCHAR(50), salary DECIMAL(10,2), d DATE, created_at DATETIME)")

	// --- 3.1 Integer Types × Modifiers (12 cases) ---
	t.Run("IntegerTypes", func(t *testing.T) {
		cases := []struct{ name, sql string }{
			{"INT baseline", "CREATE TABLE t1 (a INT)"},
			{"INT UNSIGNED", "CREATE TABLE t2 (a INT UNSIGNED)"},
			{"INT SIGNED", "CREATE TABLE t3 (a INT SIGNED)"},
			{"INT display width", "CREATE TABLE t4 (a INT(11))"},
			{"INT UNSIGNED ZEROFILL", "CREATE TABLE t5 (a INT UNSIGNED ZEROFILL)"},
			{"TINYINT SIGNED", "CREATE TABLE t6 (a TINYINT SIGNED)"},
			{"TINYINT UNSIGNED ZEROFILL", "CREATE TABLE t7 (a TINYINT UNSIGNED ZEROFILL)"},
			{"SMALLINT width UNSIGNED", "CREATE TABLE t8 (a SMALLINT(5) UNSIGNED)"},
			{"MEDIUMINT SIGNED", "CREATE TABLE t9 (a MEDIUMINT SIGNED)"},
			{"BIGINT width UNSIGNED", "CREATE TABLE t10 (a BIGINT(20) UNSIGNED)"},
			{"INT1 UNSIGNED", "CREATE TABLE t11 (a INT1 UNSIGNED)"},
			{"INT8 SIGNED", "CREATE TABLE t12 (a INT8 SIGNED)"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assertOracleMatch(t, oracle, tc.sql)
			})
		}
		// Cleanup
		for i := 1; i <= 12; i++ {
			oracle.db.ExecContext(oracle.ctx, "DROP TABLE IF EXISTS t"+itoa(i))
		}
	})

	// --- 3.2 Decimal & Float Types (16 cases) ---
	t.Run("DecimalFloatTypes", func(t *testing.T) {
		cases := []struct{ name, sql string }{
			{"DECIMAL bare", "CREATE TABLE t13 (a DECIMAL)"},
			{"DECIMAL single prec", "CREATE TABLE t14 (a DECIMAL(10))"},
			{"DECIMAL full prec", "CREATE TABLE t15 (a DECIMAL(10,2))"},
			{"DECIMAL UNSIGNED", "CREATE TABLE t16 (a DECIMAL(10,2) UNSIGNED)"},
			{"NUMERIC", "CREATE TABLE t17 (a NUMERIC(10,2))"},
			{"DEC", "CREATE TABLE t18 (a DEC(10,2))"},
			{"FIXED", "CREATE TABLE t19 (a FIXED(10,2))"},
			{"FLOAT bare", "CREATE TABLE t20 (a FLOAT)"},
			{"FLOAT with prec", "CREATE TABLE t21 (a FLOAT(10,2))"},
			{"FLOAT 24", "CREATE TABLE t22 (a FLOAT(24))"},
			{"FLOAT 25", "CREATE TABLE t23 (a FLOAT(25))"},
			{"DOUBLE", "CREATE TABLE t24 (a DOUBLE)"},
			{"DOUBLE PRECISION", "CREATE TABLE t25 (a DOUBLE PRECISION)"},
			{"REAL", "CREATE TABLE t26 (a REAL)"},
			{"FLOAT4", "CREATE TABLE t27 (a FLOAT4)"},
			{"FLOAT8", "CREATE TABLE t28 (a FLOAT8)"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assertOracleMatch(t, oracle, tc.sql)
			})
		}
		for i := 13; i <= 28; i++ {
			oracle.db.ExecContext(oracle.ctx, "DROP TABLE IF EXISTS t"+itoa(i))
		}
	})

	// --- 3.3 String & Binary Types (21 cases) ---
	t.Run("StringBinaryTypes", func(t *testing.T) {
		cases := []struct{ name, sql string }{
			{"CHAR", "CREATE TABLE t29 (a CHAR(10))"},
			{"CHAR with charset collate", "CREATE TABLE t30 (a CHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci)"},
			{"VARCHAR", "CREATE TABLE t31 (a VARCHAR(255))"},
			{"TEXT", "CREATE TABLE t32 (a TEXT)"},
			{"TEXT with length", "CREATE TABLE t33 (a TEXT(1000))"},
			{"TINYTEXT", "CREATE TABLE t34 (a TINYTEXT)"},
			{"MEDIUMTEXT charset", "CREATE TABLE t35 (a MEDIUMTEXT CHARACTER SET latin1)"},
			{"LONGTEXT", "CREATE TABLE t36 (a LONGTEXT)"},
			{"LONG", "CREATE TABLE t37 (a LONG)"},
			{"LONG VARCHAR", "CREATE TABLE t38 (a LONG VARCHAR)"},
			{"BINARY", "CREATE TABLE t39 (a BINARY(16))"},
			{"VARBINARY", "CREATE TABLE t40 (a VARBINARY(255))"},
			{"BLOB", "CREATE TABLE t41 (a BLOB)"},
			{"BLOB with length", "CREATE TABLE t42 (a BLOB(1000))"},
			{"TINYBLOB", "CREATE TABLE t43 (a TINYBLOB)"},
			{"MEDIUMBLOB", "CREATE TABLE t44 (a MEDIUMBLOB)"},
			{"LONGBLOB", "CREATE TABLE t45 (a LONGBLOB)"},
			{"LONG VARBINARY", "CREATE TABLE t46 (a LONG VARBINARY)"},
			{"NATIONAL CHAR", "CREATE TABLE t47 (a NATIONAL CHAR(10))"},
			{"NCHAR", "CREATE TABLE t48 (a NCHAR(10))"},
			{"NVARCHAR", "CREATE TABLE t49 (a NVARCHAR(100))"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assertOracleMatch(t, oracle, tc.sql)
			})
		}
		for i := 29; i <= 49; i++ {
			oracle.db.ExecContext(oracle.ctx, "DROP TABLE IF EXISTS t"+itoa(i))
		}
	})

	// --- 3.4 Date/Time, JSON, Spatial & Special (19 cases) ---
	t.Run("DateTimeJSONSpatial", func(t *testing.T) {
		cases := []struct{ name, sql string }{
			{"DATE", "CREATE TABLE t50 (a DATE)"},
			{"TIME", "CREATE TABLE t51 (a TIME)"},
			{"TIME fsp", "CREATE TABLE t52 (a TIME(3))"},
			{"DATETIME", "CREATE TABLE t53 (a DATETIME)"},
			{"DATETIME fsp", "CREATE TABLE t54 (a DATETIME(6))"},
			{"TIMESTAMP", "CREATE TABLE t55 (a TIMESTAMP)"},
			{"TIMESTAMP fsp", "CREATE TABLE t56 (a TIMESTAMP(3))"},
			{"YEAR", "CREATE TABLE t57 (a YEAR)"},
			{"BIT", "CREATE TABLE t58 (a BIT(8))"},
			{"BOOL", "CREATE TABLE t59 (a BOOL)"},
			{"JSON", "CREATE TABLE t60 (a JSON)"},
			{"SERIAL", "CREATE TABLE t61 (a SERIAL)"},
			{"ENUM", "CREATE TABLE t62 (a ENUM('a','b','c'))"},
			{"SET", "CREATE TABLE t63 (a SET('x','y','z'))"},
			{"GEOMETRY", "CREATE TABLE t64 (a GEOMETRY)"},
			{"POINT", "CREATE TABLE t65 (a POINT)"},
			{"LINESTRING", "CREATE TABLE t66 (a LINESTRING)"},
			{"POLYGON", "CREATE TABLE t67 (a POLYGON)"},
			{"GEOMETRYCOLLECTION", "CREATE TABLE t68 (a GEOMETRYCOLLECTION)"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assertOracleMatch(t, oracle, tc.sql)
			})
		}
		for i := 50; i <= 68; i++ {
			oracle.db.ExecContext(oracle.ctx, "DROP TABLE IF EXISTS t"+itoa(i))
		}
	})

	// --- 3.5 Window Functions (7 cases) ---
	t.Run("WindowFunctions", func(t *testing.T) {
		cases := []struct{ name, sql string }{
			{"RANK named window", "SELECT RANK() OVER w FROM t WINDOW w AS (ORDER BY id)"},
			{"ROW_NUMBER partition", "SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) FROM t"},
			{"LAG default offset", "SELECT LAG(val) OVER (ORDER BY id) FROM t"},
			{"LAG all args", "SELECT LAG(val, 2, 0) OVER (ORDER BY id) FROM t"},
			{"NTH_VALUE", "SELECT NTH_VALUE(val, 3) OVER (ORDER BY id) FROM t"},
			{"SUM window frame", "SELECT SUM(val) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t"},
			{"multiple window funcs", "SELECT DENSE_RANK() OVER (ORDER BY score), PERCENT_RANK() OVER (ORDER BY score) FROM t"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assertOracleMatch(t, oracle, tc.sql)
			})
		}
	})

	// --- 3.6 Interval Expressions (8 cases) ---
	t.Run("IntervalExpressions", func(t *testing.T) {
		cases := []struct {
			name, sql     string
			knownMismatch bool // true = omni/MySQL disagree; logged but not a test failure
		}{
			{"interval add DAY", "SELECT NOW() + INTERVAL 1 DAY", false},
			{"interval sub HOUR", "SELECT NOW() - INTERVAL 2 HOUR", false},
			{"DATE_ADD literal", "SELECT DATE_ADD('2024-01-01', INTERVAL 1 MONTH)", false},
			{"DATE_SUB MINUTE", "SELECT DATE_SUB(NOW(), INTERVAL 30 MINUTE)", false},
			{"TIMESTAMPADD", "SELECT TIMESTAMPADD(MINUTE, 30, NOW())", false},
			{"EXTRACT", "SELECT EXTRACT(HOUR FROM NOW())", false},
			{"interval in WHERE", "SELECT * FROM t WHERE created_at > NOW() - INTERVAL 7 DAY", false},
			{"interval in BETWEEN", "SELECT * FROM t WHERE created_at BETWEEN NOW() - INTERVAL 1 MONTH AND NOW()", false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.knownMismatch {
					assertOracleMismatch(t, oracle, tc.sql)
				} else {
					assertOracleMatch(t, oracle, tc.sql)
				}
			})
		}
	})

	// --- 3.7 Rejected SQL (8 cases) ---
	t.Run("RejectedSQL", func(t *testing.T) {
		cases := []struct {
			name, sql     string
			knownMismatch bool
		}{
			{"VARCHAR UNSIGNED", "CREATE TABLE t_rej1 (a VARCHAR UNSIGNED)", false},
			{"TEXT ZEROFILL", "CREATE TABLE t_rej2 (a TEXT ZEROFILL)", false},
			{"JSON UNSIGNED", "CREATE TABLE t_rej3 (a JSON UNSIGNED)", false},
			// MISMATCH: omni accepts unquoted reserved words as identifiers
			// but MySQL correctly rejects them. Parser is too lenient.
			{"reserved word table name", "CREATE TABLE select (a INT)", true},
			{"reserved word col name", "CREATE TABLE t_rej5 (select INT)", true},
			{"incomplete PARTITION BY", "ALTER TABLE t PARTITION BY", false},
			{"RANGE no partitions", "ALTER TABLE t PARTITION BY RANGE(id)", false},
			{"invalid interval unit", "SELECT DATE_ADD(d, INTERVAL 1 INVALID_UNIT) FROM t", false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.knownMismatch {
					assertOracleMismatch(t, oracle, tc.sql)
				} else {
					assertOracleMatch(t, oracle, tc.sql)
				}
			})
		}
		// Cleanup any tables that might have been created
		for i := 1; i <= 5; i++ {
			oracle.db.ExecContext(oracle.ctx, "DROP TABLE IF EXISTS t_rej"+itoa(i))
		}
	})
}

// itoa converts an int to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
