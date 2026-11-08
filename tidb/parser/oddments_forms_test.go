package parser

import "testing"

// TestRemoveTTLTail pins REMOVE TTL as an AlterTableSpecSingleOpt tail: it
// follows the spec list without a comma and nothing may follow it. Container-
// verified against TiDB v8.5.5.
func TestRemoveTTLTail(t *testing.T) {
	accept := []string{
		"ALTER TABLE t REMOVE TTL",
		"ALTER TABLE t TTL_ENABLE = 'OFF' REMOVE TTL",
		"ALTER TABLE t ADD COLUMN z INT REMOVE TTL",
		"ALTER TABLE t COMMENT 'x' REMOVE TTL",
	}
	reject := []string{
		// A comma before the tail, a clause after it, or a second tail is 1064.
		"ALTER TABLE t TTL_ENABLE = 'OFF', REMOVE TTL",
		"ALTER TABLE t REMOVE TTL, REMOVE TTL",
		"ALTER TABLE t REMOVE TTL REMOVE TTL",
		"ALTER TABLE t REMOVE TTL COMMENT 'x'",
	}
	for _, s := range accept {
		if _, err := Parse(s); err != nil {
			t.Errorf("Parse(%q) rejected: %v", s, err)
		}
	}
	for _, s := range reject {
		if _, err := Parse(s); err == nil {
			t.Errorf("Parse(%q) accepted; TiDB answers 1064", s)
		}
	}
}

// TestDefaultValueThenCollate pins that a column DEFAULT value is a signed
// literal or function, so a following COLLATE is its own column option in
// either order. Container-verified against TiDB v8.5.5.
func TestDefaultValueThenCollate(t *testing.T) {
	for _, s := range []string{
		"CREATE TABLE t (a VARCHAR(8) DEFAULT NULL COLLATE utf8mb4_unicode_ci)",
		"CREATE TABLE t (a VARCHAR(8) COLLATE utf8mb4_unicode_ci DEFAULT NULL)",
		"CREATE TABLE t (a VARCHAR(8) DEFAULT 'x' COLLATE utf8mb4_bin)",
	} {
		if _, err := Parse(s); err != nil {
			t.Errorf("Parse(%q) rejected: %v", s, err)
		}
	}
	// A full expression is still not a default: DEFAULT 1+1 is a 1064.
	if _, err := Parse("CREATE TABLE t (a INT DEFAULT 1+1)"); err == nil {
		t.Errorf("Parse(DEFAULT 1+1) accepted; TiDB answers 1064")
	}
}

// TestDefaultFuncQualifiedColumn pins DEFAULT(col) over a qualified column.
// Container-verified against TiDB v8.5.5.
func TestDefaultFuncQualifiedColumn(t *testing.T) {
	for _, s := range []string{
		"SELECT default(a) FROM t1",
		"SELECT default(t1.a) FROM t1, t2",
		"SELECT default(db.t1.a) FROM t1",
	} {
		if _, err := Parse(s); err != nil {
			t.Errorf("Parse(%q) rejected: %v", s, err)
		}
	}
	// A fourth qualifier is one too many.
	if _, err := Parse("SELECT default(w.db.t1.a) FROM t1"); err == nil {
		t.Errorf("Parse(default(w.db.t1.a)) accepted; want a syntax error")
	}
}
