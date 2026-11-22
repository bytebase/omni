package parser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of TiDB's parser.y and its near-misses, with omni's verdict and
// TiDB's. It reads the file and starts no server. OMNI_RECORD=1 judges the
// rows that have no verdict yet on the pinned TiDB.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		PinRole: "parser",
		Grammar: filepath.Join("..", "testdata", "parser.y"),
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the pinned TiDB, through PREPARE, whether it accepts each
// statement's syntax. It runs only in a recording, on a connection outside any
// tape. The statement is bound to a user variable and prepared from it, so it
// reaches the server byte for byte: rows hold invalid UTF-8, trailing spaces
// and a CR inside a literal. The variable lives on one connection, so the
// judge holds one.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	// testcontainer.Judge starts the shared TiDB server under the default role;
	// when a dedicated pin role (such as parser) shares that release, verify it
	// pins the exact same server before reusing it.
	if pin.Role != "default" {
		serverPin, err := enginetest.LoadPin(testcontainer.Pins, "default")
		if err != nil {
			t.Fatalf("%v", err)
		}
		if serverPin.Image != pin.Image || serverPin.Digest != pin.Digest || serverPin.Tag != pin.Tag {
			t.Fatalf("%s pin (%s %s %s) differs from default server pin (%s %s %s)",
				pin.Role, pin.Image, pin.Digest, pin.Tag, serverPin.Image, serverPin.Digest, serverPin.Tag)
		}
	}
	ctx := context.Background()
	conn, err := testcontainer.Judge(t).Conn(ctx)
	if err != nil {
		t.Fatalf("connect to TiDB: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	prepare := func(text string) error {
		if _, err := conn.ExecContext(ctx, "SET @omni_compat = ?", text); err != nil {
			// %v, not %w: a failed bind is no verdict, whatever its code.
			return fmt.Errorf("binding the statement: %v", err)
		}
		_, err := conn.ExecContext(ctx, "PREPARE omni_compat_stmt FROM @omni_compat")
		return err
	}
	// The verdict reads 1064 as the parser's rejection, so a release that
	// answers a syntax error otherwise must fail here and not be misread.
	if code, ok := mysqlErrorCode(prepare("SELECT FROM")); !ok || code != 1064 {
		t.Fatalf("TiDB no longer answers PREPARE of a syntax error with 1064: %v", prepare("SELECT FROM"))
	}
	// Judge revision 2 re-asks CollationName's 1273 (Unknown collation) parse-time
	// lookup with TiDB's built-in utf8mb4_bin stand-in; verify both the built-in
	// stand-in and the unknown-collation rewrite accept valid syntax and reject
	// trailing syntax errors.
	if v, err := compatJudgeVerdict("ALTER DATABASE COLLATE 'utf8mb4_bin'", prepare); err != nil || v != compat.Accept {
		t.Fatalf("TiDB no longer accepts built-in COLLATE 'utf8mb4_bin': verdict=%q err=%v", v, err)
	}
	if code, ok := mysqlErrorCode(prepare("ALTER DATABASE COLLATE 'utf8mb4_bin' zzz")); !ok || code != 1064 {
		t.Fatalf("TiDB no longer answers COLLATE 'utf8mb4_bin' trailing syntax error with 1064: %v", prepare("ALTER DATABASE COLLATE 'utf8mb4_bin' zzz"))
	}
	if v, err := compatJudgeVerdict("ALTER DATABASE COLLATE 's'", prepare); err != nil || v != compat.Accept {
		t.Fatalf("compatJudge no longer accepts valid unknown-collation statement: verdict=%q err=%v", v, err)
	}
	if v, err := compatJudgeVerdict("ALTER DATABASE COLLATE 's' zzz", prepare); err != nil || v != compat.Reject {
		t.Fatalf("compatJudge no longer rejects unknown-collation trailing syntax error: verdict=%q err=%v", v, err)
	}
	return func(sql string) (string, error) {
		return compatJudgeVerdict(sql, prepare)
	}
}

// builtInCollation is a collation built into TiDB, used to re-ask statements
// that abort during parsing with 1273 (Unknown collation) on an unknown name.
const builtInCollation = "utf8mb4_bin"

// maxCollationStandIns bounds how many 1273 (Unknown collation) rewrites
// compatJudgeVerdict performs on a single statement.
const maxCollationStandIns = 8

var reUnknownCollation = regexp.MustCompile(`Unknown collation: '(.*)'$`)

// replaceUnknownCollation substitutes 'utf8mb4_bin' for the first SQL COLLATE
// target matching name in stmt, scanning SQL tokens so string literals, quoted
// identifiers, and comments are skipped. It reports whether stmt changed.
func replaceUnknownCollation(stmt, name string) (string, bool) {
	if strings.EqualFold(name, builtInCollation) {
		return stmt, false
	}
	toks := scanCollationTokens(stmt)
	for i := 0; i < len(toks); i++ {
		if toks[i].kind != sqlTokWord || !strings.EqualFold(toks[i].val, "COLLATE") {
			continue
		}
		if i > 0 && toks[i-1].kind == sqlTokPunct && toks[i-1].val == "." {
			continue
		}
		target := i + 1
		if target < len(toks) && toks[target].kind == sqlTokPunct && toks[target].val == "=" {
			target++
		}
		if target >= len(toks) {
			break
		}
		tok := toks[target]
		switch tok.kind {
		case sqlTokString, sqlTokQuotedIdent, sqlTokWord:
			if strings.EqualFold(tok.val, name) {
				next := stmt[:tok.start] + "'" + builtInCollation + "'" + stmt[tok.end:]
				return next, next != stmt
			}
		}
	}
	return stmt, false
}

type sqlTokKind int

const (
	sqlTokWord sqlTokKind = iota
	sqlTokString
	sqlTokQuotedIdent
	sqlTokPunct
)

type sqlTok struct {
	kind       sqlTokKind
	start, end int
	val        string
}

// scanCollationTokens lexes stmt into coarse SQL tokens while skipping
// whitespace, line comments (-- and #), and block comments (stepping into
// MySQL/TiDB executable comments /*! ... */).
func scanCollationTokens(s string) []sqlTok {
	var toks []sqlTok
	n := len(s)
	i := 0
	for i < n {
		c := s[i]
		// Whitespace.
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\f' || c == '\v' {
			i++
			continue
		}
		// Hash line comment.
		if c == '#' {
			for i < n && s[i] != '\n' {
				i++
			}
			continue
		}
		// Dash line comment: -- followed by whitespace, control char, or EOF.
		if c == '-' && i+1 < n && s[i+1] == '-' && (i+2 == n || s[i+2] <= ' ') {
			i += 2
			for i < n && s[i] != '\n' {
				i++
			}
			continue
		}
		// Block comment or MySQL/TiDB executable comment.
		if c == '/' && i+1 < n && s[i+1] == '*' {
			if i+2 < n && s[i+2] == '!' {
				i += 3
				if i+1 < n && (s[i] == 'T' || s[i] == 't') && s[i+1] == '!' {
					i += 2
				}
				for i < n && s[i] >= '0' && s[i] <= '9' {
					i++
				}
				continue
			}
			i += 2
			for i+1 < n && !(s[i] == '*' && s[i+1] == '/') {
				i++
			}
			if i+1 < n {
				i += 2
			} else {
				i = n
			}
			continue
		}
		// Closing delimiter of an executable comment.
		if c == '*' && i+1 < n && s[i+1] == '/' {
			i += 2
			continue
		}
		// Single- or double-quoted string literal.
		if c == '\'' || c == '"' {
			start := i
			quote := c
			i++
			var b strings.Builder
			for i < n {
				ch := s[i]
				if ch == '\\' && i+1 < n {
					i++
					b.WriteByte(unescapeSQLByte(s[i]))
					i++
					continue
				}
				if ch == quote {
					if i+1 < n && s[i+1] == quote {
						b.WriteByte(quote)
						i += 2
						continue
					}
					i++
					break
				}
				b.WriteByte(ch)
				i++
			}
			toks = append(toks, sqlTok{kind: sqlTokString, start: start, end: i, val: b.String()})
			continue
		}
		// Backtick-quoted identifier.
		if c == '`' {
			start := i
			i++
			var b strings.Builder
			for i < n {
				if s[i] == '`' {
					if i+1 < n && s[i+1] == '`' {
						b.WriteByte('`')
						i += 2
						continue
					}
					i++
					break
				}
				b.WriteByte(s[i])
				i++
			}
			toks = append(toks, sqlTok{kind: sqlTokQuotedIdent, start: start, end: i, val: b.String()})
			continue
		}
		// Unquoted keyword or identifier.
		if isSQLIdentChar(c) {
			start := i
			for i < n && isSQLIdentChar(s[i]) {
				i++
			}
			toks = append(toks, sqlTok{kind: sqlTokWord, start: start, end: i, val: s[start:i]})
			continue
		}
		toks = append(toks, sqlTok{kind: sqlTokPunct, start: i, end: i + 1, val: s[i : i+1]})
		i++
	}
	return toks
}

func isSQLIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '$' || c >= 0x80
}

func unescapeSQLByte(c byte) byte {
	switch c {
	case '0':
		return 0
	case 'b':
		return '\b'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case 'Z':
		return 0x1a
	default:
		return c
	}
}

// compatJudgeVerdict reads PREPARE's answer for stmt. TiDB's CollationName
// grammar action calls charset.GetCollationByName while parsing and aborts
// with 1273 (Unknown collation) before reading what follows; re-asking with
// builtInCollation lets TiDB parse the full statement so valid clauses accept
// and trailing syntax errors reject with 1064.
func compatJudgeVerdict(stmt string, prepare func(string) error) (string, error) {
	cur := stmt
	var first error
	for step := 0; ; step++ {
		first = prepare(cur)
		var myErr *gomysql.MySQLError
		if step < maxCollationStandIns && errors.As(first, &myErr) && myErr.Number == 1273 {
			if m := reUnknownCollation.FindStringSubmatch(myErr.Message); m != nil {
				if next, ok := replaceUnknownCollation(cur, m[1]); ok {
					cur = next
					continue
				}
			}
		}
		break
	}
	accepts, err := prepareVerdict(first, func() error { return prepare(cur + compat.ProbeTail) })
	if err != nil {
		return "", err
	}
	if accepts {
		return compat.Accept, nil
	}
	return compat.Reject, nil
}

func TestCompatJudgeUnknownCollationRewrite(t *testing.T) {
	prepare := func(text string) error {
		switch text {
		case "ALTER DATABASE COLLATE 's'",
			"ALTER DATABASE COLLATE 's' zzz":
			return &gomysql.MySQLError{Number: 1273, Message: "Unknown collation: 's'"}
		case "ALTER DATABASE COLLATE 'utf8mb4_bin'":
			return &gomysql.MySQLError{Number: 1295, Message: "This command is not supported in the prepared statement protocol yet"}
		case "ALTER DATABASE COLLATE 'utf8mb4_bin'" + compat.ProbeTail,
			"ALTER DATABASE COLLATE 'utf8mb4_bin' zzz":
			return &gomysql.MySQLError{Number: 1064, Message: "syntax error"}
		case "CREATE TABLE t (a VARCHAR(10) COLLATE 's1', b VARCHAR(10) COLLATE 's2') zzz":
			return &gomysql.MySQLError{Number: 1273, Message: "Unknown collation: 's1'"}
		case "CREATE TABLE t (a VARCHAR(10) COLLATE 'utf8mb4_bin', b VARCHAR(10) COLLATE 's2') zzz":
			return &gomysql.MySQLError{Number: 1273, Message: "Unknown collation: 's2'"}
		case "CREATE TABLE t (a VARCHAR(10) COLLATE 'utf8mb4_bin', b VARCHAR(10) COLLATE 'utf8mb4_bin') zzz":
			return &gomysql.MySQLError{Number: 1064, Message: "syntax error"}
		case "DO 'COLLATE s' COLLATE s":
			return &gomysql.MySQLError{Number: 1273, Message: "Unknown collation: 's'"}
		case "DO 'COLLATE s' COLLATE 'utf8mb4_bin'":
			return nil
		default:
			t.Fatalf("unexpected PREPARE text: %q", text)
			return nil
		}
	}

	got, err := compatJudgeVerdict("ALTER DATABASE COLLATE 's'", prepare)
	if err != nil || got != compat.Accept {
		t.Fatalf("valid unknown collation: got (%q, %v), want (%q, nil)", got, err, compat.Accept)
	}
	got, err = compatJudgeVerdict("ALTER DATABASE COLLATE 's' zzz", prepare)
	if err != nil || got != compat.Reject {
		t.Fatalf("trailing near-miss unknown collation: got (%q, %v), want (%q, nil)", got, err, compat.Reject)
	}
	got, err = compatJudgeVerdict("CREATE TABLE t (a VARCHAR(10) COLLATE 's1', b VARCHAR(10) COLLATE 's2') zzz", prepare)
	if err != nil || got != compat.Reject {
		t.Fatalf("multiple unknown collations with trailing syntax error: got (%q, %v), want (%q, nil)", got, err, compat.Reject)
	}
	got, err = compatJudgeVerdict("DO 'COLLATE s' COLLATE s", prepare)
	if err != nil || got != compat.Accept {
		t.Fatalf("string literal containing COLLATE text: got (%q, %v), want (%q, nil)", got, err, compat.Accept)
	}

	for _, tc := range []struct {
		stmt string
		name string
		want string
		ok   bool
	}{
		{"ALTER DATABASE COLLATE 's' zzz", "s", "ALTER DATABASE COLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE = 's' zzz", "s", "ALTER DATABASE COLLATE = 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE='s' zzz", "s", "ALTER DATABASE COLLATE='utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE's'zzz", "s", "ALTER DATABASE COLLATE'utf8mb4_bin'zzz", true},
		{`ALTER DATABASE COLLATE "s" zzz`, "s", "ALTER DATABASE COLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE `s` zzz", "s", "ALTER DATABASE COLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE s zzz", "s", "ALTER DATABASE COLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE = s zzz", "s", "ALTER DATABASE COLLATE = 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE 'S' zzz", "s", "ALTER DATABASE COLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE '' zzz", "", "ALTER DATABASE COLLATE 'utf8mb4_bin' zzz", true},
		{"DO 'COLLATE s' COLLATE s", "s", "DO 'COLLATE s' COLLATE 'utf8mb4_bin'", true},
		{"DO 'COLLATE ''s''' COLLATE 's' zzz", "s", "DO 'COLLATE ''s''' COLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE /* COLLATE 's' */ COLLATE 's' zzz", "s", "ALTER DATABASE /* COLLATE 's' */ COLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE -- COLLATE 's'\nCOLLATE 's' zzz", "s", "ALTER DATABASE -- COLLATE 's'\nCOLLATE 'utf8mb4_bin' zzz", true},
		{"ALTER DATABASE COLLATE /* c */ = /* c */ 's' zzz", "s", "ALTER DATABASE COLLATE /* c */ = /* c */ 'utf8mb4_bin' zzz", true},
		{"CREATE TABLE t (a CHAR(1)) /*!40100 COLLATE 's' */ zzz", "s", "CREATE TABLE t (a CHAR(1)) /*!40100 COLLATE 'utf8mb4_bin' */ zzz", true},
		{"ALTER DATABASE `COLLATE` 's'", "s", "ALTER DATABASE `COLLATE` 's'", false},
		{"SELECT t.collate s, a COLLATE s FROM t", "s", "SELECT t.collate s, a COLLATE 'utf8mb4_bin' FROM t", true},
		{"ALTER DATABASE COLLATE 'utf8mb4_bin'", "utf8mb4_bin", "ALTER DATABASE COLLATE 'utf8mb4_bin'", false},
		{"ALTER DATABASE COLLATE super zzz", "s", "ALTER DATABASE COLLATE super zzz", false},
		{"SELECT 1", "s", "SELECT 1", false},
	} {
		got, ok := replaceUnknownCollation(tc.stmt, tc.name)
		if ok != tc.ok || got != tc.want {
			t.Errorf("replaceUnknownCollation(%q, %q) = (%q, %v), want (%q, %v)", tc.stmt, tc.name, got, ok, tc.want, tc.ok)
		}
	}
}
