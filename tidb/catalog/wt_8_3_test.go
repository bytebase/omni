package catalog

import (
	"strings"
	"testing"
)

// --- 4.3 SHOW CREATE TABLE Charset Rendering ---

func TestWalkThrough_8_3_CharsetSameAsTable(t *testing.T) {
	// Column charset same as table — CHARACTER SET not shown in column def
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100)) DEFAULT CHARSET=utf8mb4")
	got := c.ShowCreateTable("testdb", "t")
	// Column should NOT show CHARACTER SET since it matches table charset
	if strings.Contains(got, "CHARACTER SET") {
		t.Errorf("expected no CHARACTER SET in column def when charset matches table\ngot:\n%s", got)
	}
}

func TestWalkThrough_8_3_CharsetDiffersFromTable(t *testing.T) {
	// Column charset differs from table — CHARACTER SET shown
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4")
	got := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(got, "CHARACTER SET latin1") {
		t.Errorf("expected CHARACTER SET latin1 in column def when charset differs from table\ngot:\n%s", got)
	}
}

func TestWalkThrough_8_3_CollationNonDefaultSameAsTable(t *testing.T) {
	// Column collation is non-default for its charset but same as table — COLLATE shown
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100)) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	got := c.ShowCreateTable("testdb", "t")
	// The column inherits utf8mb4_unicode_ci from table, which is non-default for utf8mb4.
	// MySQL shows COLLATE in this case (collation inherited but non-default for charset).
	if !strings.Contains(got, "COLLATE utf8mb4_unicode_ci") {
		t.Errorf("expected COLLATE utf8mb4_unicode_ci in column def\ngot:\n%s", got)
	}
	// CHARACTER SET should NOT be shown since charset matches table
	lines := strings.Split(got, "\n")
	for _, line := range lines {
		if strings.Contains(line, "`name`") && strings.Contains(line, "CHARACTER SET") {
			t.Errorf("expected no CHARACTER SET in column def when charset matches table\nline: %s", line)
		}
	}
}

// The SHOW CREATE TABLE text in these is what TiDB v8.5.5 writes for the same
// statements (read on 2026-11-06).
func TestWalkThrough_8_3_CollationDiffersFromTable(t *testing.T) {
	// Column collation differs from the table's in the same charset: only
	// COLLATE is written.
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, name VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	got := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(got, "  `name` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL\n") {
		t.Errorf("expected `name` varchar(100) COLLATE utf8mb4_bin\ngot:\n%s", got)
	}
}

func TestWalkThrough_8_3_Utf8mb4DefaultCollation(t *testing.T) {
	// TiDB writes the table's collation for every charset but binary.
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT) DEFAULT CHARSET=utf8mb4")
	got := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(got, ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin") {
		t.Errorf("expected COLLATE=utf8mb4_bin in table options for utf8mb4\ngot:\n%s", got)
	}
}

func TestWalkThrough_8_3_Latin1DefaultCollation(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT) DEFAULT CHARSET=latin1")
	got := c.ShowCreateTable("testdb", "t")
	if !strings.Contains(got, ") ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin") {
		t.Errorf("expected COLLATE=latin1_bin in table options for latin1\ngot:\n%s", got)
	}
}

func TestWalkThrough_8_3_BinaryCharset(t *testing.T) {
	// BINARY charset on column — rendered as CHARACTER SET binary
	// MySQL 8.0 converts CHAR(N) CHARACTER SET binary → binary(N),
	// and VARCHAR(N) CHARACTER SET binary → varbinary(N).
	// But ENUM/SET with CHARACTER SET binary retains the charset annotation.

	// Sub-test 1: CHAR CHARACTER SET binary → binary(N) in MySQL 8.0
	t.Run("char_binary_converts", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t (id INT, data CHAR(100) CHARACTER SET binary) DEFAULT CHARSET=utf8mb4")
		got := c.ShowCreateTable("testdb", "t")
		// MySQL converts to binary(100) type — no CHARACTER SET annotation
		if !strings.Contains(got, "`data` binary(100)") {
			t.Errorf("expected CHAR(100) CHARACTER SET binary to render as binary(100)\ngot:\n%s", got)
		}
	})

	// Sub-test 2: ENUM with CHARACTER SET binary — shows CHARACTER SET binary
	t.Run("enum_charset_binary", func(t *testing.T) {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE t (id INT, status ENUM('a','b','c') CHARACTER SET binary) DEFAULT CHARSET=utf8mb4")
		got := c.ShowCreateTable("testdb", "t")
		if !strings.Contains(got, "`status` enum('a','b','c') CHARACTER SET binary COLLATE binary DEFAULT NULL") {
			t.Errorf("expected CHARACTER SET binary COLLATE binary in ENUM column def\ngot:\n%s", got)
		}
	})
}
