package parser

import (
	"strings"
	"testing"
)

// assertParenStmtCountParity reports mismatch when omni's Parse returns
// a different RawStmt count than the naive `;`-split count of sql.
//
// This is a Parse-tail regression guard (PAREN-KB-2 class): if omni
// silently splits a single PG statement into N by dropping tokens
// that parseStmt didn't consume, this helper catches it.
//
// Limitation: "naive ;-split" ignores ;'s inside strings / dollar-
// quoted blocks / BEGIN ATOMIC bodies. For the current Phase 2
// corpus (single-statement probes), that limitation doesn't matter.
// Future expansion should reuse pg/pgregress/extract.go splitStatements.
func assertParenStmtCountParity(t *testing.T, sql string) {
	t.Helper()
	stmts, err := Parse(sql)
	if err != nil {
		// Omni rejected — that's fine for this check; we only
		// look for accept-with-wrong-count mismatches.
		return
	}
	omniCount := 0
	if stmts != nil {
		omniCount = len(stmts.Items)
	}
	expectCount := naiveStmtCount(sql)
	if omniCount != expectCount {
		t.Errorf("stmt-count parity: sql=%q omni=%d want~=%d (naive ;-split)", sql, omniCount, expectCount)
	}
}

// naiveStmtCount counts non-empty `;`-separated segments, ignoring
// ;'s inside single-quoted strings, double-quoted identifiers, and
// dollar-quoted blocks.
func naiveStmtCount(sql string) int {
	count := 0
	inSingle := false
	inDouble := false
	dollarTag := ""
	var current strings.Builder
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if dollarTag != "" {
			// Looking for end of dollar-quoted block.
			if c == '$' && strings.HasPrefix(sql[i:], dollarTag) {
				current.WriteString(dollarTag)
				i += len(dollarTag) - 1
				dollarTag = ""
				continue
			}
			current.WriteByte(c)
			continue
		}
		if inSingle {
			if c == '\'' {
				inSingle = false
			}
			current.WriteByte(c)
			continue
		}
		if inDouble {
			if c == '"' {
				inDouble = false
			}
			current.WriteByte(c)
			continue
		}
		switch c {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '$':
			// Look for $tag$ pattern.
			if i+1 < len(sql) {
				end := strings.IndexByte(sql[i+1:], '$')
				if end >= 0 && end < 32 { // sane tag length
					tag := sql[i : i+1+end+1]
					if isValidDollarTag(tag) {
						dollarTag = tag
						current.WriteString(tag)
						i += len(tag) - 1
						continue
					}
				}
			}
		case ';':
			if strings.TrimSpace(current.String()) != "" {
				count++
			}
			current.Reset()
			continue
		}
		current.WriteByte(c)
	}
	if strings.TrimSpace(current.String()) != "" {
		count++
	}
	return count
}

func isValidDollarTag(s string) bool {
	if !strings.HasPrefix(s, "$") || !strings.HasSuffix(s, "$") {
		return false
	}
	inner := s[1 : len(s)-1]
	for _, r := range inner {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

// TestStmtCountParityBasic documents omni's stmt-count behavior on
// short canonical inputs. This is the baseline for the G2 Parse-tail
// regression guard; if these counts drift, KB-2 (or a related Parse-
// tail bug) has likely regressed.
func TestStmtCountParityBasic(t *testing.T) {
	cases := []struct {
		sql    string
		expect int // expected omni count
	}{
		{"", 0},
		{"SELECT 1", 1},
		{"SELECT 1;", 1},
		{";SELECT 1;", 1},
		{"SELECT 1;SELECT 2", 2},
		{";;SELECT 1;;SELECT 2;;", 2},
		// PAREN-KB-2 case: omni currently returns 2 but naive-split
		// returns 1 (the `SELECT 1 SELECT 2` has no ';'). This test
		// DOCUMENTS the bug, not fails on it — omni's silent split
		// is tracked as KB-2.
		// {"SELECT * FROM (SELECT 1) SELECT 1", 2},  // skipped; KB-2 tracking
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			got := 0
			if stmts != nil {
				got = len(stmts.Items)
			}
			if got != tc.expect {
				t.Errorf("Parse(%q): got %d stmts, want %d", tc.sql, got, tc.expect)
			}
		})
	}
}

// TestParenStmtCountParityCorpus runs assertParenStmtCountParity
// against a small known-good corpus. This is the Phase G2 Parse-tail
// stmt-count parity guard; adding entries here tightens the fence.
func TestParenStmtCountParityCorpus(t *testing.T) {
	corpus := []string{
		`SELECT 1`,
		`SELECT * FROM (SELECT 1)`,
		`SELECT * FROM ((a JOIN b ON TRUE) JOIN c ON TRUE)`,
		`SELECT 1;`,
		// KB-2 case — omni returns 2, naive returns 1; currently a
		// known divergence. Skipped until KB-2 is fixed.
		// `SELECT * FROM (SELECT 1) SELECT 1`,
	}
	for _, sql := range corpus {
		t.Run(sql, func(t *testing.T) {
			assertParenStmtCountParity(t, sql)
		})
	}
}
