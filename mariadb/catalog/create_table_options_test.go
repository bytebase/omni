package catalog

import "testing"

// What the CREATE TABLE forms this engine newly parses leave in the catalog.
// A form can start parsing and still misreport forever — a SQL_TSI_YEAR column
// that lands as an unknown type, an ENUM whose hex member comes out as the text
// `0x61`, a TABLE_CHECKSUM that is dropped on the floor — and a parser test
// cannot see any of it.
//
// Every want below is a live mariadb:11.8.8's own SHOW CREATE TABLE, read back
// from the server: `year(4)` for every spelling and modifier of YEAR but
// YEAR(2), `enum('a','b')` for a member written 0x61, NOT NULL AUTO_INCREMENT
// with a unique key for SERIAL DEFAULT VALUE, and CHECKSUM=1 for
// TABLE_CHECKSUM, which reports under CHECKSUM's name.
func TestCreateTableOptionsCatalog(t *testing.T) {
	for _, tc := range []struct {
		name  string
		ddl   string
		table string
		want  string
	}{
		{
			// SQL_TSI_YEAR is YEAR and not an unknown type, and neither
			// UNSIGNED nor ZEROFILL is part of what a YEAR column reports.
			name:  "SQL_TSI_YEAR is the YEAR type",
			ddl:   "CREATE TABLE y (a SQL_TSI_YEAR, b SQL_TSI_YEAR(4), c YEAR ZEROFILL, d SQL_TSI_YEAR(4) UNSIGNED ZEROFILL)",
			table: "y",
			want: "CREATE TABLE `y` (\n" +
				"  `a` year(4) DEFAULT NULL,\n" +
				"  `b` year(4) DEFAULT NULL,\n" +
				"  `c` year(4) DEFAULT NULL,\n" +
				"  `d` year(4) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci",
		},
		{
			// A hex or bit member is the bytes it holds, so 0x61 is the member
			// `a` and 0x6162 is `ab` — not the text that wrote them.
			name:  "an ENUM or SET member written as a literal is its bytes",
			ddl:   "CREATE TABLE e (c1 ENUM(0x61, 'b'), c2 SET(0b01100001, 'b'), c3 ENUM(x'61'), c4 ENUM(0x6162))",
			table: "e",
			want: "CREATE TABLE `e` (\n" +
				"  `c1` enum('a','b') DEFAULT NULL,\n" +
				"  `c2` set('a','b') DEFAULT NULL,\n" +
				"  `c3` enum('a') DEFAULT NULL,\n" +
				"  `c4` enum('ab') DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci",
		},
		{
			// SERIAL DEFAULT VALUE is NOT NULL AUTO_INCREMENT UNIQUE, so the
			// column is not null and auto-incrementing and the table gains the
			// unique key.
			name:  "SERIAL DEFAULT VALUE is NOT NULL AUTO_INCREMENT UNIQUE",
			ddl:   "CREATE TABLE s (a INT SERIAL DEFAULT VALUE)",
			table: "s",
			want: "CREATE TABLE `s` (\n" +
				"  `a` int(11) NOT NULL AUTO_INCREMENT,\n" +
				"  UNIQUE KEY `a` (`a`)\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci",
		},
		{
			// A name that begins with a digit is the name, backquoted on the
			// way out whether or not it was on the way in.
			name:  "a digit-leading name is a name",
			ddl:   "CREATE TABLE `123` (123a1 INT, 1_2 INT, 12e INT)",
			table: "123",
			want: "CREATE TABLE `123` (\n" +
				"  `123a1` int(11) DEFAULT NULL,\n" +
				"  `1_2` int(11) DEFAULT NULL,\n" +
				"  `12e` int(11) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci",
		},
		{
			// TABLE_CHECKSUM reports under CHECKSUM, the other name for the
			// same option, and as a flag: any non-zero value reads back as 1
			// and a zero as no option at all.
			name:  "TABLE_CHECKSUM reports as CHECKSUM",
			ddl:   "CREATE TABLE ck (a INT) TABLE_CHECKSUM=1",
			table: "ck",
			want: "CREATE TABLE `ck` (\n" +
				"  `a` int(11) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1",
		},
		{
			name:  "TABLE_CHECKSUM is a flag, not a number",
			ddl:   "CREATE TABLE ck (a INT) TABLE_CHECKSUM=1.5",
			table: "ck",
			want: "CREATE TABLE `ck` (\n" +
				"  `a` int(11) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci CHECKSUM=1",
		},
		{
			name:  "a zero checksum is no option",
			ddl:   "CREATE TABLE ck (a INT) TABLE_CHECKSUM=0",
			table: "ck",
			want: "CREATE TABLE `ck` (\n" +
				"  `a` int(11) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, tc.ddl)
			if got := c.ShowCreateTable("testdb", tc.table); got != tc.want {
				t.Errorf("SHOW CREATE TABLE after %s:\ngot:\n%s\nwant:\n%s", tc.ddl, got, tc.want)
			}
		})
	}
}
