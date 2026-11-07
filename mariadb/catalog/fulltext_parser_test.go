package catalog

import (
	"strings"
	"testing"
)

// Configured-plugin metadata fixtures complement the missing-plugin native
// probes. The upstream plugin suite pins the stored WITH PARSER spelling:
// https://github.com/MariaDB/server/blob/mariadb-11.8.8/mysql-test/suite/plugins/r/fulltext_plugin.result
func TestConfiguredMariaDBFulltextParser(t *testing.T) {
	for _, ddl := range []string{
		"CREATE TABLE t(a TEXT,FULLTEXT KEY k(a) WITH PARSER simple_parser)",
		"CREATE TABLE t(a TEXT); ALTER TABLE t ADD FULLTEXT KEY k(a) WITH PARSER simple_parser",
		"CREATE TABLE t(a TEXT); CREATE FULLTEXT INDEX k ON t(a) WITH PARSER simple_parser",
	} {
		t.Run(ddl, func(t *testing.T) {
			c := New()
			c.SetFulltextParsers("simple_parser")
			mustExec(t, c, "CREATE DATABASE test; USE test; "+ddl)
			mustExec(t, c, "CREATE TABLE copy LIKE t")
			for _, name := range []string{"t", "copy"} {
				tbl := c.GetDatabase("test").GetTable(name)
				if len(tbl.Indexes) != 1 || tbl.Indexes[0].ParserName != "simple_parser" {
					t.Fatalf("lost configured parser on %s: %+v", name, tbl.Indexes)
				}
				if sql := c.ShowCreateTable("test", name); !strings.Contains(sql, "WITH PARSER `simple_parser`") {
					t.Errorf("lost WITH PARSER in stored DDL: %s", sql)
				}
			}
		})
	}
}
