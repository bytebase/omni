package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

func TestTiDBCommentBasic(t *testing.T) {
	// /*T! ... */ — always inject inner SQL
	sql := "CREATE TABLE t (id BIGINT /*T! PRIMARY KEY */)"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse TiDB comment: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", list.Len())
	}
	// Verify PRIMARY KEY was injected (column has PK constraint)
	stmt := list.Items[0].(*ast.CreateTableStmt)
	if len(stmt.Columns) == 0 {
		t.Fatal("no columns parsed")
	}
	hasPK := false
	for _, c := range stmt.Columns[0].Constraints {
		if c.Type == ast.ColConstrPrimaryKey {
			hasPK = true
		}
	}
	if !hasPK {
		t.Error("PRIMARY KEY from /*T! */ comment was not injected")
	}
}

func TestTiDBCommentFeatureSupported(t *testing.T) {
	// /*T![auto_rand] ... */ — inject because auto_rand is supported in v8.5
	sql := "SELECT /*T![auto_rand] 1 AS */ col FROM t"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse supported feature comment: %v", err)
	}
}

func TestTiDBCommentFeatureUnsupported(t *testing.T) {
	// /*T![unsupported_feature_xyz] ... */ — skip as comment
	sql := "SELECT /*T![unsupported_feature_xyz] WEIRD_STUFF */ 1"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse unsupported feature comment: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", list.Len())
	}
}

func TestTiDBCommentMultiFeature(t *testing.T) {
	// /*T![ttl] ... */ — single supported feature
	sql := "SELECT /*T![ttl] 1 AS */ col FROM t"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("failed to parse single-feature comment: %v", err)
	}
}

func TestTiDBCommentPreservesMySQL(t *testing.T) {
	// Standard MySQL conditional comments still work
	sql := "SELECT /*!50000 1 */ + 1"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("MySQL conditional comment broke: %v", err)
	}
}

func TestTiDBCommentEmpty(t *testing.T) {
	// A segment with only a TiDB comment is not empty
	sql := "/*T! SELECT 1 */"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("TiDB comment-only segment failed: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement from TiDB comment, got %d", list.Len())
	}
}

func TestTiDBCommentUnclosed(t *testing.T) {
	// Unclosed TiDB comments must not panic. TiDB v8.5.5 tolerates EOF
	// inside a bang comment: the content is already lexed as SQL, so these
	// all parse (container-verified).
	tests := []struct {
		sql     string
		wantErr bool
	}{
		{"SELECT /*T! 1", false},              // SELECT 1
		{"SELECT /*T![ttl] 1", false},         // SELECT 1
		{"SELECT /*T![auto_rand] col", false}, // SELECT col
		{"/*T! SELECT 1", false},              // SELECT 1
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("parser panicked on unclosed TiDB comment: %v", r)
				}
			}()
			_, err := Parse(tt.sql)
			if tt.wantErr && err == nil {
				t.Error("expected parse error for unclosed TiDB comment, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

func TestScanTiDBFeatureIDs(t *testing.T) {
	tests := []struct {
		bracket    string
		wantIDs    []string
		wellFormed bool
	}{
		{"[auto_rand]", []string{"auto_rand"}, true},
		{"[auto_rand,clustered_index]", []string{"auto_rand", "clustered_index"}, true},
		{"[unsupported_xyz]", []string{"unsupported_xyz"}, true},
		{"[]", nil, false},         // empty segment
		{"[abc,]", nil, false},     // trailing empty segment
		{"[abc", nil, false},       // missing ]
		{"[a b]", nil, false},      // space is not an ident char
		{"[auto_rand,]", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.bracket, func(t *testing.T) {
			ids, after, wellFormed := scanTiDBFeatureIDs(tt.bracket, 0)
			if wellFormed != tt.wellFormed {
				t.Fatalf("wellFormed = %v, want %v", wellFormed, tt.wellFormed)
			}
			if !tt.wellFormed {
				if after != 0 {
					t.Fatalf("malformed list must reset: after = %d, want 0", after)
				}
				return
			}
			if after != len(tt.bracket) {
				t.Fatalf("after = %d, want %d", after, len(tt.bracket))
			}
			if len(ids) != len(tt.wantIDs) {
				t.Fatalf("ids = %v, want %v", ids, tt.wantIDs)
			}
			for i := range ids {
				if ids[i] != tt.wantIDs[i] {
					t.Fatalf("ids = %v, want %v", ids, tt.wantIDs)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Comment-lexer conformance tests, container-verified against TiDB v8.5.5
// (pingcap/tidb@sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26).
// TiDB semantics (pkg/parser/lexer.go startWithSlash/startWithStar):
//   - block comments do NOT nest: the first */ closes the comment;
//   - an unclosed regular block comment is a parse error;
//   - /*! and supported /*T! comments lex their content as SQL; EOF inside
//     one is tolerated (the construct just ends);
//   - /*!NNNNN consumes exactly 5 version digits or none (scanVersionDigits(5,5));
//   - /*T![...] feature IDs follow scanFeatureIDs: a malformed bracket list
//     resets and the content (starting at '[') is lexed as SQL; a well-formed
//     list with any unsupported ID makes the whole thing a regular comment.
// ---------------------------------------------------------------------------

func parseVerdict(t *testing.T, sql string) (nStmts int, err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parser panicked on %q: %v", sql, r)
		}
	}()
	list, perr := Parse(sql)
	if perr != nil {
		return 0, perr
	}
	return list.Len(), nil
}

func TestBlockCommentNonNesting(t *testing.T) {
	tests := []struct {
		sql      string
		wantErr  bool
		wantLen  int
	}{
		// TiDB v8.5.5: accept — comment ends at first */, "1" is the select item.
		{"SELECT /* a /* b */ 1", false, 1},
		// TiDB v8.5.5: accept, 1 row — the SELECT after the comment executes.
		// omni used to swallow it (0 statements parsed = silent statement loss).
		{"/* a /* b */ SELECT 1", false, 1},
		{"/* a /* b */ SELECT 1; SELECT 2", false, 2},
		{"SELECT /* a /* b */ 1; SELECT 2", false, 2},
		{"SELECT /* a /* b */ /* c */ 1", false, 1},
		// Unsupported-feature /*T! is a regular comment: also non-nesting.
		{"SELECT /*T![unsupported_feature_xyz] a /* b */ 1", false, 1},
		// TiDB v8.5.5: reject 1064 near "*/ SELECT 1".
		{"/*T![unsupported_feature_xyz] a /* b */ */ SELECT 1", true, 0},
		// Comment then empty statement then SELECT: TiDB executes SELECT 2.
		{"/* a /* b */ ; SELECT 2", false, 1},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			n, err := parseVerdict(t, tt.sql)
			if tt.wantErr && err == nil {
				t.Fatalf("want parse error, got %d stmts", n)
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("want %d stmts, got error: %v", tt.wantLen, err)
				}
				if n != tt.wantLen {
					t.Fatalf("want %d stmts, got %d", tt.wantLen, n)
				}
			}
		})
	}
}

func TestBlockCommentUnterminated(t *testing.T) {
	// TiDB v8.5.5 rejects all of these with 1064 (unclosed regular comment).
	for _, sql := range []string{
		"SELECT 1 /*",
		"SELECT 1 /* unterminated",
		"SELECT 1 /**",
		"SELECT /** / 1",
		"/* only a comment",
		"SELECT 1 /*M! x",
		"SELECT 1 /* x */; SELECT 2 /* unterminated",
		// Well-formed but unsupported feature list downgrades to a regular
		// comment, so unclosed is an error too.
		"SELECT /*T![unsupported_xyz] 1",
	} {
		t.Run(sql, func(t *testing.T) {
			if n, err := parseVerdict(t, sql); err == nil {
				t.Fatalf("want parse error for unclosed comment, got %d stmts", n)
			}
		})
	}
}

func TestBangCommentUnterminated(t *testing.T) {
	// TiDB v8.5.5 ACCEPTS an unclosed /*! or supported /*T! construct: the
	// content is already lexed as SQL, EOF just ends it. omni used to panic
	// (slice out of range) on /*! and to swallow /*T! to end of input.
	tests := []struct {
		sql     string
		wantLen int
	}{
		{"SELECT /*! 1", 1},
		{"SELECT /*!50000 1", 1},
		{"/*! SELECT 1", 1},
		{"SELECT /*T! 1", 1},
		{"SELECT /*T![auto_rand] 1", 1},
		{"/*T! SELECT 1", 1},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			n, err := parseVerdict(t, tt.sql)
			if err != nil {
				t.Fatalf("want accept (%d stmts), got error: %v", tt.wantLen, err)
			}
			if n != tt.wantLen {
				t.Fatalf("want %d stmts, got %d", tt.wantLen, n)
			}
		})
	}
	// Unclosed bang construct whose content ends inside a regular comment:
	// TiDB rejects (the inner /* c never closes).
	if n, err := parseVerdict(t, "SELECT /*! 1 /* c"); err == nil {
		t.Fatalf("want parse error (unclosed inner comment), got %d stmts", n)
	}
}

func TestBangCommentVersionDigits(t *testing.T) {
	// scanVersionDigits(5,5): exactly five digits are the version prefix;
	// fewer-than-five or the sixth-and-later digits are SQL content.
	tests := []struct {
		sql     string
		wantErr bool
	}{
		{"SELECT /*!123 */", false},     // content "123": SELECT 123
		{"SELECT /*!1234567 */", false}, // "12345" version, content "67"
		{"SELECT /*!523456 */", false},  // "52345" version, content "6"
		{"SELECT /*!12345 6 */", false}, // classic versioned comment
		{"SELECT /*!m 1 */", true},      // no digits: content "m 1" → SELECT m 1 rejects
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			n, err := parseVerdict(t, tt.sql)
			if tt.wantErr && err == nil {
				t.Fatalf("want parse error, got %d stmts", n)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want accept, got error: %v", err)
			}
		})
	}
}

func TestBangCommentContentScanning(t *testing.T) {
	tests := []struct {
		sql     string
		wantErr bool
	}{
		// */ hidden inside string literals, line comments, backtick idents,
		// and inner regular comments does not close the construct.
		{"SELECT /*! '*/' */", false},
		{"SELECT /*!50000 1 # */\n */", false},
		{"SELECT /*!50000 1 -- */\n */", false},
		{"SELECT /*! 1 AS `*/` */ FROM (SELECT 1) t", false},
		{"SELECT /*!50000 1 /* c */ + 1 */", false},
		{"SELECT /*T! 1 /* c */ + 1 */", false},
		// Nested bang markers do not stack: TiDB's inBangComment is a bool,
		// so the first unconsumed */ closes everything and the trailing */
		// is a syntax error. TiDB v8.5.5: reject 1064 near "/".
		{"SELECT /*!50000 1 /*!50000 + 1 */ */", true},
		{"SELECT /*T! 1 /*T! + 1 */ */", true},
		// Multi-statement inside a versioned comment splits as SQL.
		{"SELECT /*!50000 1; SELECT 2 */", false},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			n, err := parseVerdict(t, tt.sql)
			if tt.wantErr && err == nil {
				t.Fatalf("want parse error, got %d stmts", n)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want accept, got error: %v", err)
			}
		})
	}
}

func TestTiDBFeatureCommentFSM(t *testing.T) {
	tests := []struct {
		sql     string
		wantErr bool
	}{
		// Malformed bracket lists reset: content starts at '[' and is lexed
		// as SQL, which rejects. TiDB v8.5.5: 1064 near "[...".
		{"SELECT /*T![] 1 */", true},
		{"SELECT /*T![abc 1 */", true},
		{"SELECT /*T![abc,] 1 */", true},
		// No bracket at all: bang mode, content right after T!.
		{"SELECT /*T! 1 */", false},
		{"SELECT /*T!50000 1 */", true}, // content "50000 1" → SELECT 50000 1 rejects
		// Well-formed, all-supported list (v8.5.5 has 9 feature IDs).
		{"SELECT /*T![auto_rand,ttl] 1 + 1 */", false},
		{"SELECT /*T![force_inc] 1 */", false},
		{"SELECT /*T![affinity] 1 */", false},
		{"SELECT /*T![global_index] 1 */", false},
		// Lowercase /*t! is never special: regular comment, so the select
		// list vanishes and the bare SELECT rejects (TiDB: 1064 near "").
		{"SELECT /*t! 1 */", true},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			n, err := parseVerdict(t, tt.sql)
			if tt.wantErr && err == nil {
				t.Fatalf("want parse error, got %d stmts", n)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want accept, got error: %v", err)
			}
		})
	}
}

func TestGlobalIndexCommentPreserved(t *testing.T) {
	// v8.5.5 supports global_index; with a stale feature list the comment
	// body (NONCLUSTERED) was silently dropped from the AST.
	sql := "CREATE TABLE tgi (a BIGINT PRIMARY KEY /*T![global_index] NONCLUSTERED */)"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("want 1 stmt, got %d", list.Len())
	}
	// NONCLUSTERED sets the tri-state clustered flag to false; without the
	// comment the outfunc omits the field entirely.
	out := ast.NodeToString(list.Items[0])
	if !strings.Contains(out, ":clustered f") {
		t.Errorf("NONCLUSTERED from /*T![global_index] */ comment lost: %s", out)
	}
}

// Known, deliberate divergences from TiDB v8.5.5 (parked leniencies, both
// reject-only on the TiDB side — omni accepting is fail-safe because the
// engine remains the final validity gate):
//  1. TiDB's Scanner.Lex lookahead for AS OF / TO TIMESTAMP / MEMBER OF
//     corrupts inBangComment when the token before the closing */ is
//     AS/TO/MEMBER (getNextToken restores the reader but not the flag), so
//     TiDB rejects e.g. `SELECT /*! 1 AS */ x`. omni closes the construct
//     cleanly and accepts.
//  2. omni tolerates unclosed string literals (`SELECT '1`), so bang content
//     ending in an unclosed string (`SELECT /*! '1 */`) parses as a string
//     where TiDB rejects. Pre-existing lexer-wide leniency, tracked
//     separately from comment handling.
//  3. omni accepts `SELECT * / 1` (raw input, container-verified TiDB 1064),
//     so a trailing */ after a closed comment in select-item position
//     (`SELECT /* a /* b */ */ 1`) parses even though the token stream now
//     matches TiDB. Pre-existing select-grammar leniency.
func TestBangCommentParkedLeniencies(t *testing.T) {
	for _, sql := range []string{
		"SELECT /*! 1 AS */ x",      // TiDB: 1064 via lookahead state corruption
		"SELECT /*! '1 */",          // TiDB: 1064 via unclosed-string detection
		"SELECT /* a /* b */ */ 1",  // TiDB: 1064 near "/ 1" (SELECT * / 1 shape)
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := parseVerdict(t, sql); err != nil {
				t.Fatalf("parked leniency regressed to reject: %v", err)
			}
		})
	}
}

func TestBangCommentSpliceTokenBoundaries(t *testing.T) {
	// Stripped markers and the closing */ are token boundaries in TiDB, so
	// the splice must not butt adjacent content bytes together
	// (container-verified against TiDB v8.5.5).
	t.Run("stripped marker cannot reconstruct a comment opener", func(t *testing.T) {
		// TiDB lexes 2 / * 3 * / and rejects 1064; naive zero-byte marker
		// stripping produced "2/ *3" -> phantom /* -> accepted as SELECT 2.
		if n, err := parseVerdict(t, "SELECT /*! 2//*!50000*3*/ */"); err == nil {
			t.Fatalf("want parse error, got %d stmts", n)
		}
	})
	t.Run("closing */ separates content and suffix tokens", func(t *testing.T) {
		// TiDB: x aliased YZ (two tokens). A seamless splice merged them
		// into one column reference xYZ (silent wrong AST).
		list, err := Parse("SELECT /*!50000 x*/YZ FROM (SELECT 1 AS x) t")
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		out := ast.NodeToString(list.Items[0])
		if strings.Contains(out, "xYZ") {
			t.Errorf("content and suffix merged into one token: %s", out)
		}
	})
	t.Run("ident butted against closer stays separate", func(t *testing.T) {
		// TiDB: SELECT SELE CT parses (SELE aliased CT, unknown column at
		// exec is still a parse-accept).
		if _, err := parseVerdict(t, "SELECT/*!50000 SELE*/CT"); err != nil {
			t.Fatalf("want accept, got: %v", err)
		}
	})
	t.Run("marker replacement does not join operators", func(t *testing.T) {
		// TiDB: 1 - (-2) = 3; the two dashes must not become -- (comment).
		if _, err := parseVerdict(t, "SELECT 1 -/*!50000- 2 */"); err != nil {
			t.Fatalf("want accept, got: %v", err)
		}
	})
}
