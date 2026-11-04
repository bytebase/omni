package catalog

import "testing"

// What the CREATE TABLE forms this engine newly parses leave in the catalog.
// A form can start parsing and still misreport forever — an ENUM whose hex
// member comes out as the text `0x61`, a column whose digit-leading name is lost
// — and a parser test cannot see any of it.
//
// Only two of this family's forms are OceanBase's: an ENUM or SET member written
// as a hex or bit literal, and a name that begins with a digit. SQL_TSI_YEAR,
// the numeric modifiers on YEAR, SERIAL DEFAULT VALUE and TABLE_CHECKSUM are not
// read here, because CE 4.3.5.6 answers 1064 for each of them -- see
// oceanbase/parser/create_table_options_test.go.
//
// This catalog renders MySQL 8.0's default collation rather than
// OceanBase's; that divergence predates this change and is its own to fix. What
// the new forms contribute — the member values and the names — was read from a
// live 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-04.
func TestCreateTableOptionsCatalog(t *testing.T) {
	for _, tc := range []struct {
		name  string
		ddl   string
		table string
		want  string
	}{
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
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci",
		},
		{
			// A name that begins with a digit is the name, backquoted on the
			// way out whether or not it was on the way in.
			name:  "a digit-leading name is a name",
			ddl:   "CREATE TABLE `123` (123a1 INT, 1_2 INT, 12e INT)",
			table: "123",
			want: "CREATE TABLE `123` (\n" +
				"  `123a1` int DEFAULT NULL,\n" +
				"  `1_2` int DEFAULT NULL,\n" +
				"  `12e` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci",
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
