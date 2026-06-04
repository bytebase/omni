package parser

import "testing"

// The standalone BINARY modifier (CHAR(10) BINARY — shorthand for the charset's
// _bin collation). TiDB accepts it on character types (CHAR/VARCHAR/TEXT
// variants, ENUM, SET, NCHAR/NVARCHAR) in any order relative to CHARACTER SET /
// COLLATE, and rejects it on BLOB/VARBINARY. Verified on TiDB v8.5.0.

func TestCharBinaryModifierAccepted(t *testing.T) {
	cases := []string{
		`CREATE TABLE t (a CHAR(5) BINARY)`,
		`CREATE TABLE t (a CHAR BINARY)`,
		`CREATE TABLE t (a VARCHAR(5) BINARY)`,
		`CREATE TABLE t (a TEXT BINARY)`,
		`CREATE TABLE t (a TINYTEXT BINARY)`,
		`CREATE TABLE t (a LONGTEXT BINARY)`,
		`CREATE TABLE t (a ENUM('x','y') BINARY)`,
		`CREATE TABLE t (a SET('x') BINARY)`,
		`CREATE TABLE t (a NCHAR(5) BINARY)`,
		`CREATE TABLE t (a NVARCHAR(5) BINARY)`,
		// BINARY in any order relative to CHARACTER SET / COLLATE.
		`CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 BINARY)`,
		`CREATE TABLE t (a CHAR(5) BINARY CHARACTER SET utf8mb4)`,
		`CREATE TABLE t (a CHAR(5) BINARY COLLATE utf8mb4_bin)`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestBinaryTypeBinaryModifierRejected(t *testing.T) {
	// TiDB rejects the BINARY modifier on binary types (they have no charset).
	cases := []string{
		`CREATE TABLE t (a BLOB BINARY)`,
		`CREATE TABLE t (a TINYBLOB BINARY)`,
		`CREATE TABLE t (a VARBINARY(5) BINARY)`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) accepted BINARY on a binary type, but TiDB rejects it", sql)
			}
		})
	}
}

// BINARY and CHARACTER SET may appear in any order, but COLLATE comes LAST —
// TiDB rejects BINARY after COLLATE.
func TestBinaryCollateOrdering(t *testing.T) {
	accept := []string{
		`CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 BINARY COLLATE utf8mb4_bin)`,
		`CREATE TABLE t (a CHAR(5) BINARY COLLATE utf8mb4_bin)`,
		`CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)`,
	}
	reject := []string{
		`CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin BINARY)`,
		`CREATE TABLE t (a CHAR(5) COLLATE utf8mb4_bin BINARY)`,
	}
	for _, sql := range accept {
		t.Run("accept/"+sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
	for _, sql := range reject {
		t.Run("reject/"+sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) accepted BINARY after COLLATE, but TiDB rejects it", sql)
			}
		})
	}
}

// CHARACTER SET and the standalone BINARY modifier may each appear at most once;
// TiDB rejects repeats. `CHARACTER SET binary BINARY` is NOT a repeat — the
// first `binary` is the charset name, the second is the modifier.
func TestRepeatedCharsetBinaryRejected(t *testing.T) {
	reject := []string{
		`CREATE TABLE t (a CHAR(5) BINARY BINARY)`,
		`CREATE TABLE t (a CHAR(5) BINARY CHARACTER SET utf8mb4 BINARY)`,
		`CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 BINARY BINARY)`,
		`CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 CHARACTER SET latin1)`,
	}
	for _, sql := range reject {
		t.Run("reject/"+sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) accepted a repeated charset/BINARY modifier, but TiDB rejects it", sql)
			}
		})
	}
	// charset name `binary` followed by the BINARY modifier is valid.
	if _, err := Parse(`CREATE TABLE t (a CHAR(5) CHARACTER SET binary BINARY)`); err != nil {
		t.Fatalf("CHARACTER SET binary BINARY should parse: %v", err)
	}
}
