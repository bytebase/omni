package catalog

import (
	"strings"
	"testing"
)

// TestBinaryModifierCollationResolution pins how the standalone BINARY modifier
// resolves to a collation, verified against TiDB v8.5.0:
//   - BINARY → {charset}_bin (using the effective charset)
//   - an explicit COLLATE wins over BINARY
//   - CHARACTER SET binary → collation `binary` (NOT binary_bin)
func TestBinaryModifierCollationResolution(t *testing.T) {
	c := New()
	results, err := c.Exec(`CREATE DATABASE testdb; USE testdb;
CREATE TABLE t (
  a CHAR(10) BINARY,
  b VARCHAR(10) CHARACTER SET latin1 BINARY,
  c VARCHAR(10) CHARACTER SET utf8mb4 BINARY COLLATE utf8mb4_unicode_ci,
  d ENUM('x','y') CHARACTER SET binary BINARY,
  e SET('x') CHARACTER SET binary BINARY
) DEFAULT CHARSET=utf8mb4;`, nil)
	if err != nil {
		t.Fatalf("exec parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("exec error on stmt %d: %v", r.Index, r.Error)
		}
	}

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}
	want := map[string]string{
		"a": "utf8mb4_bin",        // BINARY → effective (db) charset utf8mb4 → utf8mb4_bin
		"b": "latin1_bin",         // explicit charset latin1 → latin1_bin
		"c": "utf8mb4_unicode_ci", // explicit COLLATE wins over BINARY
		"d": "binary",             // CHARACTER SET binary → binary, not binary_bin
		"e": "binary",
	}
	for col, exp := range want {
		got := tbl.GetColumn(col)
		if got == nil {
			t.Errorf("column %s not found", col)
			continue
		}
		if strings.ToLower(got.Collation) != exp {
			t.Errorf("column %s: got collation %q, want %q", col, got.Collation, exp)
		}
	}
}

// TestCharsetBinaryConvertsToBinaryType pins that a non-ENUM/SET string column
// with CHARACTER SET binary (with or without the BINARY modifier) becomes a
// binary type with no charset/collation — in BOTH the CREATE and ALTER paths.
func TestCharsetBinaryConvertsToBinaryType(t *testing.T) {
	c := New()
	exec := func(ddl string) {
		t.Helper()
		results, err := c.Exec(ddl, nil)
		if err != nil {
			t.Fatalf("exec parse error: %v", err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatalf("exec error on stmt %d: %v", r.Index, r.Error)
			}
		}
	}
	exec(`CREATE DATABASE testdb; USE testdb;`)
	exec(`CREATE TABLE t (
  a CHAR(10) CHARACTER SET binary BINARY,
  b VARCHAR(10) CHARACTER SET binary,
  c TEXT CHARACTER SET binary
);`)
	// ALTER ADD must apply the same conversion as CREATE.
	exec(`ALTER TABLE t ADD d CHAR(10) CHARACTER SET binary BINARY;`)
	exec(`ALTER TABLE t ADD e CHAR(10) CHARACTER SET binary;`)

	tbl := c.GetDatabase("testdb").GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}
	wantType := map[string]string{"a": "binary", "b": "varbinary", "c": "blob", "d": "binary", "e": "binary"}
	for col, exp := range wantType {
		got := tbl.GetColumn(col)
		if got == nil {
			t.Errorf("column %s not found", col)
			continue
		}
		if strings.ToLower(got.DataType) != exp {
			t.Errorf("column %s: got type %q, want %q", col, got.DataType, exp)
		}
		if got.Charset != "" || got.Collation != "" {
			t.Errorf("column %s: binary type must have no charset/collation, got charset=%q collation=%q", col, got.Charset, got.Collation)
		}
	}
}
