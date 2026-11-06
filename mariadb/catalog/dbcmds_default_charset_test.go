package catalog

import "testing"

// TestDatabaseCharsetDefault: CHARACTER SET DEFAULT and COLLATE DEFAULT name the server's, which a
// table then inherits, rather than a charset called "default" (mariadb:11.8.8 takes both).
func TestDatabaseCharsetDefault(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE dd CHARACTER SET DEFAULT")
	mustExec(t, c, "CREATE DATABASE de COLLATE DEFAULT")
	mustExec(t, c, "CREATE DATABASE df CHARACTER SET latin1")
	mustExec(t, c, "ALTER DATABASE df CHARACTER SET DEFAULT")
	mustExec(t, c, "CREATE DATABASE dg CHARACTER SET latin1")
	mustExec(t, c, "ALTER DATABASE dg COLLATE DEFAULT")
	for _, name := range []string{"dd", "df"} {
		if got := c.GetDatabase(name).Charset; got != c.defaultCharset {
			t.Errorf("%s has charset %q, want the server's %q", name, got, c.defaultCharset)
		}
	}
	for _, name := range []string{"de", "dg"} {
		if got := c.GetDatabase(name).Collation; got != c.defaultCollation {
			t.Errorf("%s has collation %q, want the server's %q", name, got, c.defaultCollation)
		}
	}
	c.SetCurrentDatabase("dd")
	mustExec(t, c, "CREATE TABLE t (a VARCHAR(10))")
	if got := c.GetDatabase("dd").GetTable("t").Charset; got != c.defaultCharset {
		t.Errorf("a table in dd has charset %q, want %q", got, c.defaultCharset)
	}
}
