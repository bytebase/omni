package parser

import "testing"

// TestGeneratedColumnPersistent: PERSISTENT is MariaDB's native synonym for
// STORED (SHOW CREATE normalizes it to STORED), so it maps to Stored=true.
func TestGeneratedColumnPersistent(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT, b INT AS (a+1) PERSISTENT)",
		"CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a+1) PERSISTENT)",
	} {
		t.Run(sql, func(t *testing.T) {
			ct := parseCreateTable(t, sql)
			col := ct.Columns[1]
			if col.Generated == nil {
				t.Fatalf("column b: expected a generated column")
			}
			if !col.Generated.Stored {
				t.Errorf("PERSISTENT should map to Stored=true")
			}
		})
	}
}

// TestGeneratedColumnPersistentAlter: PERSISTENT via ALTER ADD COLUMN (shared path).
func TestGeneratedColumnPersistentAlter(t *testing.T) {
	if _, err := Parse("ALTER TABLE t ADD COLUMN c INT AS (a+1) PERSISTENT"); err != nil {
		t.Errorf("unexpected parse error: %v", err)
	}
}

// TestMariaDBTableOptions: MariaDB's parser accepts any `IDENT = value` table
// option (an unknown name is a semantic error 1911, not a parse error), so a
// parse-only check must accept the shape. Covers the page/transactional family.
func TestMariaDBTableOptions(t *testing.T) {
	ct := parseCreateTable(t, "CREATE TABLE t (a INT) PAGE_COMPRESSED=1 PAGE_COMPRESSION_LEVEL=5 TRANSACTIONAL=1 PAGE_CHECKSUM=1")
	got := map[string]string{}
	for _, o := range ct.Options {
		got[o.Name] = o.Value
	}
	for _, name := range []string{"PAGE_COMPRESSED", "PAGE_COMPRESSION_LEVEL", "TRANSACTIONAL", "PAGE_CHECKSUM"} {
		if _, ok := got[name]; !ok {
			t.Errorf("missing table option %s (got %v)", name, got)
		}
	}
	if got["PAGE_COMPRESSION_LEVEL"] != "5" {
		t.Errorf("PAGE_COMPRESSION_LEVEL value: got %q, want \"5\"", got["PAGE_COMPRESSION_LEVEL"])
	}
}

// TestTableOptionMixedWithStructured: generic options coexist with structured ones.
func TestTableOptionMixedWithStructured(t *testing.T) {
	if _, err := Parse("CREATE TABLE t (a INT) ENGINE=Aria PAGE_CHECKSUM=1 ROW_FORMAT=PAGE"); err != nil {
		t.Errorf("unexpected parse error: %v", err)
	}
}

// TestGeneratedColumnPersistentQuotedReject: backtick-quoted `PERSISTENT` is a
// quoted identifier, not the keyword — MariaDB 11.8.8 rejects it (1064). (STORED
// and VIRTUAL match the keyword token, so they were never affected.)
func TestGeneratedColumnPersistentQuotedReject(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT, b INT AS (a+1) `PERSISTENT`)",
		"CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a+1) `PERSISTENT`)",
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestTableOptionNoValueReject: a generic option with `=` but no value token is
// a syntax error (1064 vs MariaDB 11.8.8) — consumeOptionValue returns "" at
// EOF/non-value, so the fallback must confirm a value was actually consumed.
func TestTableOptionNoValueReject(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT) PAGE_COMPRESSED=",
		"CREATE TABLE t (a INT) FOOBAR=",
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestTableOptionQuotedNameAccept: a backtick-quoted option NAME is a valid
// quoted identifier to MariaDB's parser (validated semantically), so omni must
// keep accepting it — the quoting guard belongs on the keyword, not the name.
func TestTableOptionQuotedNameAccept(t *testing.T) {
	if _, err := Parse("CREATE TABLE t (a INT) `PAGE_COMPRESSED`=1"); err != nil {
		t.Errorf("unexpected parse error: %v", err)
	}
}

// TestTableOptionBareReject: a bare option name without `= value` is a syntax
// error (1064 vs MariaDB 11.8.8) — the real parse rule the generic fallback keeps.
func TestTableOptionBareReject(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT) PAGE_COMPRESSED",
		"CREATE TABLE t (a INT) PAGE_COMPRESSED PAGE_CHECKSUM",
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}
