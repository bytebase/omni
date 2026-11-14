package catalog

import (
	"errors"
	"testing"
)

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
		{
			// LIKE copies TRANSACTIONAL and the engine-defined options and drops
			// PAGE_CHECKSUM.
			name:  "CREATE TABLE LIKE keeps TRANSACTIONAL and engine options",
			ddl:   "CREATE TABLE lsrc (a int) page_checksum=1 transactional=1 page_compressed=1; CREATE TABLE ldst LIKE lsrc",
			table: "ldst",
			want: "CREATE TABLE `ldst` (\n" +
				"  `a` int(11) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci TRANSACTIONAL=1 `page_compressed`=1",
		},
		{
			// PAGE_CHECKSUM survives a rename and nothing else that rebuilds
			// the table without naming it.
			name:  "PAGE_CHECKSUM is dropped by an ALTER that does not name it",
			ddl:   "CREATE TABLE pc (a int) page_checksum=1 transactional=0; ALTER TABLE pc RENAME TO pc2; ALTER TABLE pc2 ADD b int",
			table: "pc2",
			want: "CREATE TABLE `pc2` (\n" +
				"  `a` int(11) DEFAULT NULL,\n" +
				"  `b` int(11) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci TRANSACTIONAL=0",
		},
		{
			// Aria's handler keeps the checksum an ALTER does not name.
			name:  "Aria keeps PAGE_CHECKSUM through an ALTER",
			ddl:   "CREATE TABLE ar (a INT) ENGINE=Aria PAGE_CHECKSUM=0; ALTER TABLE ar ADD b INT",
			table: "ar",
			want: "CREATE TABLE `ar` (\n" +
				"  `a` int(11) DEFAULT NULL,\n" +
				"  `b` int(11) DEFAULT NULL\n" +
				") ENGINE=Aria DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci PAGE_CHECKSUM=0",
		},
		{
			// A later setting moves an engine-defined option to the end, a hex
			// value reads as its number, and DEFAULT removes it.
			name:  "engine-defined options keep MariaDB's order",
			ddl:   "CREATE TABLE eo (a int) PAGE_COMPRESSED='1' page_compression_level = 0x3; ALTER TABLE eo PAGE_COMPRESSED=yes; ALTER TABLE eo encryption_key_id=2, page_compression_level=default",
			table: "eo",
			want: "CREATE TABLE `eo` (\n" +
				"  `a` int(11) DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci `PAGE_COMPRESSED`=yes `encryption_key_id`=2",
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

// TestEngineDefinedOptionRefusals holds the catalog to MariaDB 11.8.8's answers
// for engine-defined options InnoDB does not define (1911), a quoted name
// included, and for SEQUENCE=1 on a table that is no sequence (4086).
func TestEngineDefinedOptionRefusals(t *testing.T) {
	for _, tc := range []struct {
		ddl  string
		code int
	}{
		{"CREATE TABLE r1 (a int) foo=1", 1911},
		{"CREATE TABLE r2 (a int) `PAGE_CHECKSUM`='x'", 1911},
		{"CREATE TABLE r3 (a int); ALTER TABLE r3 READ_ONLY = 's'", 1911},
		{"CREATE TABLE r4 (a int); ALTER TABLE r4 SEQUENCE 1", 4086},
		{"CREATE TABLE r5 (a int) foo=1 ENGINE=MyISAM", 1911},
		{"CREATE TABLE r6 (a int) ENGINE=Aria page_compressed=1", 1911},
	} {
		c := wtSetup(t)
		results, err := c.Exec(tc.ddl, nil)
		if err != nil {
			t.Fatalf("%s: %v", tc.ddl, err)
		}
		last := results[len(results)-1].Error
		var e *Error
		if !errors.As(last, &e) || e.Code != tc.code {
			t.Errorf("%s: got %v, want %d", tc.ddl, last, tc.code)
		}
	}
}

// TestEngineDefinedOptionOtherEngine takes an option of an engine the catalog
// does not know as written: CONNECT's TABLE_TYPE is CONNECT's to check.
func TestEngineDefinedOptionOtherEngine(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE ct (a INT) ENGINE=CONNECT TABLE_TYPE=CSV")
}
