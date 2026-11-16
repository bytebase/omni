package parser

import "testing"

// TestForUpdateOfReject pins that omni/mariadb rejects FOR UPDATE/SHARE ... OF,
// which MariaDB removed in 11.4 (the OF object list is MySQL-only). FOR UPDATE /
// FOR SHARE without OF are unaffected (FOR SHARE itself remains a separately
// tracked over-accept — see subtractiveDivergences).
func TestForUpdateOfReject(t *testing.T) {
	reject := []string{
		"SELECT * FROM t FOR UPDATE OF t",
		"SELECT * FROM t FOR SHARE OF t",
		"SELECT * FROM t FOR UPDATE OF t, t2",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}
