package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/mariadb/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of MariaDB's sql_yacc.yy and its near-misses, with omni's
// verdict and MariaDB's. It reads the file and starts no server. The grammar
// is under the GPL, so the repo does not hold it: the test downloads it once
// (../testdata/sources.tsv). The file holds the sql_mode=ORACLE grammar too,
// which omni does not parse: the rows are those of %ifdef MARIADB.
// OMNI_RECORD=1 judges the rows that have no verdict yet on the pinned
// MariaDB.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: enginetest.Fetched(t, filepath.Join("..", "testdata", "sources.tsv"), "sql_yacc.yy"),
		Defined: []string{"MARIADB"},
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

// compatPrepareTimeout bounds one PREPARE, so a statement the server hangs on
// gets no verdict and the recording goes on.
const compatPrepareTimeout = time.Minute

// compatJudge asks the pinned MariaDB, through PREPARE, whether it accepts
// each statement's syntax. testdata/README.md says what the answer can and
// cannot tell.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	ctx := context.Background()
	db := testcontainer.Pinned(t, pin)
	// A prepared statement and USE live on one connection.
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("connect to MariaDB: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	// Select database a so generated rows naming a.a against an unqualified a
	// (such as DELETE a.a FROM a;) resolve to the same schema rather than
	// failing with 1109 (ER_UNKNOWN_TABLE). Create a real key cache omni_cache
	// so structured system-variable stand-ins parse past the lookup.
	for _, setup := range []string{
		"CREATE DATABASE IF NOT EXISTS a",
		"USE a",
		"SET GLOBAL omni_cache.key_buffer_size = 1048576",
	} {
		if _, err := conn.ExecContext(ctx, setup); err != nil {
			t.Fatalf("MariaDB judge setup %q: %v", setup, err)
		}
	}
	j := &compatJudgeMariaDB{ctx: ctx, conn: conn, answers: map[uint16]int{}}
	// The verdict rests on these answers, so a release that gives others must
	// fail here and not be misread.
	for _, probe := range []struct {
		stmt string
		code uint16
	}{
		{"SELECT FROM", 1064},
		{"EXECUTE a", 1295},
		{"SELECT 1" + compat.ProbeTail, 1064},
		{"DELETE a.a FROM a;", 1146},
		{"DO@@omni_cache.key_buffer_size;", 0},
	} {
		err := j.prepare(probe.stmt)
		if probe.code == 0 {
			if err != nil {
				t.Fatalf("MariaDB no longer prepares %q cleanly: got %v", probe.stmt, err)
			}
			continue
		}
		if code, ok := mariaDBErrorCode(err); !ok || code != probe.code {
			t.Fatalf("MariaDB no longer answers PREPARE of %q with %d: got %v", probe.stmt, probe.code, err)
		}
	}
	t.Cleanup(func() { t.Logf("%s", j.report()) })
	return j.verdict
}

type compatJudgeMariaDB struct {
	ctx  context.Context
	conn *sql.Conn
	// answers counts the statements by the error PREPARE gave first, 0 for
	// none; standIns counts statements re-asked after a parse-time lookup
	// error; lookups lists any statement still answered blindly by a lookup
	// error without a stand-in.
	answers  map[uint16]int
	standIns int
	lookups  []string
}

// prepare asks the server to prepare text, without running it.
func (j *compatJudgeMariaDB) prepare(text string) error {
	ctx, cancel := context.WithTimeout(j.ctx, compatPrepareTimeout)
	defer cancel()
	_, err := j.conn.ExecContext(ctx, "PREPARE omni_compat_stmt FROM ?", text)
	if err == nil {
		_, _ = j.conn.ExecContext(ctx, "DEALLOCATE PREPARE omni_compat_stmt")
	}
	return err
}

// maxLookupStandIns bounds how many parse-time lookup errors verdict will
// rewrite in one statement (e.g. a multi-assignment SET STATEMENT).
const maxLookupStandIns = 8

// verdict reads PREPARE's answer: 1064 rejects, a parse-time lookup error is
// re-asked with a known-valid stand-in so MariaDB parses past the lookup and
// judges the rest of the statement (falling back to acceptance only if no
// stand-in applies), and compat.PrepareVerdict's probe decides every other
// error.
func (j *compatJudgeMariaDB) verdict(stmt string) (string, error) {
	cur := stmt
	var firstCode uint16
	for step := 0; ; step++ {
		first := j.prepare(cur)
		code, sent := mariaDBErrorCode(first)
		if first != nil && !sent {
			return "", first
		}
		if step == 0 {
			firstCode = code
		}
		if mariaDBLookupErrors[code] {
			if step < maxLookupStandIns {
				var myErr *gomysql.MySQLError
				if errors.As(first, &myErr) {
					if next, ok := mariaDBStandIn(cur, myErr); ok && next != cur {
						if step == 0 {
							j.standIns++
						}
						cur = next
						continue
					}
				}
			}
			j.lookups = append(j.lookups, fmt.Sprintf("%d %s", code, stmt))
			j.answers[firstCode]++
			return compat.Accept, nil
		}
		accepts, err := compat.PrepareVerdict(first, mariaDBErrorCode, func() error {
			return j.prepare(cur + compat.ProbeTail)
		})
		if err != nil {
			return "", err
		}
		j.answers[firstCode]++
		if accepts {
			return compat.Accept, nil
		}
		return compat.Reject, nil
	}
}

// mariaDBLookupErrors are errors MariaDB raises inside a grammar action while
// parsing, for a name the server looks up and does not have (or a feature the
// build disables): a character set (1115), a plugin or loadable function
// (1128), a system variable (1193, 1272, 4081), a collation (1253, 1273), a
// storage engine (1286), and SHOW ... CODE without debug (1289). Because they
// abort the parse before later tokens are read, verdict rewrites the reported
// lookup with a known-valid stand-in (mariaDBStandIn) and asks MariaDB again.
var mariaDBLookupErrors = map[uint16]bool{
	1115: true, 1128: true, 1193: true, 1253: true, 1272: true,
	1273: true, 1286: true, 1289: true, 4081: true,
}

var (
	reUnknownEngine      = regexp.MustCompile(`Unknown storage engine '([^']*)'`)
	reUnknownCharset     = regexp.MustCompile(`Unknown character set: '([^']*)'`)
	reUnknownCollation   = regexp.MustCompile(`Unknown collation: '([^']*)'`)
	reCollationMismatch  = regexp.MustCompile(`COLLATION '([^']*)' is not valid for CHARACTER SET '([^']*)'`)
	reFunctionNotDefined = regexp.MustCompile(`Function '([^']*)' is not defined`)
	reUnknownSysvar      = regexp.MustCompile(`Unknown system variable '([^']*)'`)
	reNotVarComponent    = regexp.MustCompile(`Variable '([^']*)' is not a variable component`)
	reUnknownStructVar   = regexp.MustCompile(`Unknown structured system variable or ROW routine variable '([^']*)'`)
	reShowCode           = regexp.MustCompile(`(?i)\bSHOW\s+(PROCEDURE|FUNCTION|PACKAGE\s+BODY)\s+CODE\b`)
	reSetStatement       = regexp.MustCompile(`(?i)^\s*SET\s+STATEMENT\s+[^:=]`)
	reASCIIBeforeCollate = regexp.MustCompile(`(?i)\bASCII(\s+COLLATE\b)`)
)

// mariaDBStandIn substitutes a known-valid MariaDB 11.8.8 name or equivalent
// production for the lookup error myErr reported, returning the rewritten
// statement and true if a substitution was made.
func mariaDBStandIn(stmt string, myErr *gomysql.MySQLError) (string, bool) {
	switch myErr.Number {
	case 1286:
		m := reUnknownEngine.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		re := regexp.MustCompile(`(?i)(\bENGINE\s+)` + sqlNamePattern(m[1]))
		return replaceOnce(re, stmt, "${1}InnoDB")

	case 1115:
		m := reUnknownCharset.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		re := regexp.MustCompile(`(?i)(\b(?:NAMES|CHARSET|CHARACTER\s+SET)\s+)` + sqlNamePattern(m[1]))
		return replaceOnce(re, stmt, "${1}utf8mb4")

	case 1273:
		m := reUnknownCollation.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		re := regexp.MustCompile(`(?i)(\bCOLLATE(?:\s*=\s*|\s+))` + sqlNamePattern(m[1]))
		return replaceOnce(re, stmt, "${1}utf8mb4_bin")

	case 1253:
		m := reCollationMismatch.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		coll, cs := m[1], m[2]
		targetCS := "utf8mb4"
		if strings.EqualFold(coll, "binary") {
			targetCS = "binary"
		} else if prefix, _, ok := strings.Cut(coll, "_"); ok && prefix != "" {
			targetCS = prefix
		}
		if strings.EqualFold(cs, "latin1") {
			if next, ok := replaceOnce(reASCIIBeforeCollate, stmt, "CHARSET "+targetCS+"${1}"); ok {
				return next, true
			}
		}
		re := regexp.MustCompile(`(?i)(\b(?:CHARSET|CHARACTER\s+SET)\s+)` + sqlNamePattern(cs) + `(\s+COLLATE\b)`)
		return replaceOnce(re, stmt, "${1}"+targetCS+"${2}")

	case 1128:
		m := reFunctionNotDefined.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		if schema, fn, ok := strings.Cut(m[1], "."); ok {
			re := regexp.MustCompile(`(?i)` + sqlNamePattern(schema) + `(\s*\.\s*` + regexp.QuoteMeta(fn) + `\s*\()`)
			return replaceOnce(re, stmt, "mariadb_schema${1}")
		}
		// Fulltext parser plugin (fulltext_key_opt: WITH PARSER IDENT_sys).
		// MariaDB 11.8.8 ships no ftparser plugin; COMMENT 's' is another
		// fulltext_key_opt alternative with the same follow set.
		re := regexp.MustCompile(`(?i)\bWITH\s+PARSER\s+` + sqlNamePattern(m[1]))
		return replaceOnce(re, stmt, "COMMENT 's'")

	case 1289:
		// SHOW PROCEDURE|FUNCTION|PACKAGE BODY CODE sp_name is disabled in
		// non-debug builds; SHOW CREATE PROCEDURE|FUNCTION|PACKAGE BODY sp_name
		// is in the same show_param rule with the same sp_name tail.
		return replaceOnce(reShowCode, stmt, "SHOW CREATE ${1}")

	case 1272:
		m := reNotVarComponent.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		re := regexp.MustCompile(`(?i)(@@\s*(?:(?:GLOBAL|LOCAL|SESSION)\s*\.\s*)?)(?:` + "`[^`]+`" + `|[A-Za-z0-9_$]+)\s*\.\s*` + sqlNamePattern(m[1]))
		return replaceOnce(re, stmt, "${1}omni_cache.key_buffer_size")

	case 4081:
		m := reUnknownStructVar.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		re := regexp.MustCompile(`(?i)(\bSET\s+(?:(?:STATEMENT|GLOBAL|LOCAL|SESSION)\s+)?|,\s*(?:(?:GLOBAL|LOCAL|SESSION)\s+)?|@@\s*(?:(?:GLOBAL|LOCAL|SESSION)\s*\.\s*)?)` + sqlNamePattern(m[1]) + `\s*\.\s*(?:` + "`[^`]+`" + `|[A-Za-z0-9_$]+)`)
		return replaceOnce(re, stmt, "${1}omni_cache.key_buffer_size")

	case 1193:
		m := reUnknownSysvar.FindStringSubmatch(myErr.Message)
		if m == nil {
			return "", false
		}
		namePat := sqlNamePattern(m[1])
		// 1. Structured @@<struct>.<name> where <struct> is not a scope prefix.
		reStruct := regexp.MustCompile(`(?i)(@@\s*(?:(?:GLOBAL|LOCAL|SESSION)\s*\.\s*)?)(` + "`[^`]+`" + `|[A-Za-z0-9_$]+)\s*\.\s*` + namePat)
		for _, loc := range reStruct.FindAllStringSubmatchIndex(stmt, -1) {
			structName := strings.Trim(stmt[loc[4]:loc[5]], "`")
			if !isVarScope(structName) {
				return stmt[:loc[4]] + "omni_cache.key_buffer_size" + stmt[loc[1]:], true
			}
		}
		// 2. Plain @@[GLOBAL.|LOCAL.|SESSION.]<name>.
		reAtAt := regexp.MustCompile(`(?i)(@@\s*(?:(?:GLOBAL|LOCAL|SESSION)\s*\.\s*)?)` + namePat)
		if next, ok := replaceOnce(reAtAt, stmt, "${1}autocommit"); ok {
			return next, true
		}
		// 3. SET STATEMENT <name> = ... (uses a SET STATEMENT-capable variable).
		if reSetStatement.MatchString(stmt) {
			reStmtVar := regexp.MustCompile(`(?i)(\bSTATEMENT\s+|,\s*)` + namePat + `(\s*:?=)`)
			if next, ok := replaceOnce(reStmtVar, stmt, "${1}max_statement_time${2}"); ok {
				return next, true
			}
		}
		// 4. SET [GLOBAL|LOCAL|SESSION] <name> = ...
		reSetVar := regexp.MustCompile(`(?i)(\bSET\s+(?:(?:GLOBAL|LOCAL|SESSION)\s+)?|,\s*(?:(?:GLOBAL|LOCAL|SESSION)\s+)?)` + namePat + `(\s*:?=)`)
		return replaceOnce(reSetVar, stmt, "${1}autocommit${2}")
	}
	return "", false
}

func isVarScope(s string) bool {
	return strings.EqualFold(s, "GLOBAL") || strings.EqualFold(s, "LOCAL") || strings.EqualFold(s, "SESSION")
}

func sqlNamePattern(name string) string {
	q := regexp.QuoteMeta(name)
	parts := []string{"`" + q + "`", "'" + q + "'"}
	if len(name) > 0 && isWordChar(name[0]) && isWordChar(name[len(name)-1]) {
		parts = append(parts, `\b`+q+`\b`)
	} else {
		parts = append(parts, q)
	}
	return "(?:" + strings.Join(parts, "|") + ")"
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

func replaceOnce(re *regexp.Regexp, s, repl string) (string, bool) {
	loc := re.FindStringSubmatchIndex(s)
	if loc == nil {
		return "", false
	}
	expanded := re.ExpandString(nil, repl, s, loc)
	next := s[:loc[0]] + string(expanded) + s[loc[1]:]
	return next, next != s
}

// report says how the server answered the statements this run judged.
func (j *compatJudgeMariaDB) report() string {
	var b strings.Builder
	codes := make([]uint16, 0, len(j.answers))
	for c := range j.answers {
		codes = append(codes, c)
	}
	slices.SortFunc(codes, func(a, b uint16) int { return j.answers[b] - j.answers[a] })
	b.WriteString("PREPARE's first answer, by error code (0 is none):")
	for _, c := range codes {
		fmt.Fprintf(&b, " %d×%d", c, j.answers[c])
	}
	fmt.Fprintf(&b, "\n%d statements were re-asked with a stand-in after a parse-time lookup error", j.standIns)
	if len(j.lookups) == 0 {
		b.WriteString("\n0 statements were answered for a name the server does not have without a stand-in")
	} else {
		fmt.Fprintf(&b, "\n%d statements were answered for a name the server does not have, and read as accepted:\n  %s",
			len(j.lookups), strings.Join(j.lookups, "\n  "))
	}
	return b.String()
}

// mariaDBErrorCode returns the server's error number, and false for an error
// the server did not send.
func mariaDBErrorCode(err error) (uint16, bool) {
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		return 0, false
	}
	return myErr.Number, true
}

func TestMariaDBStandIn(t *testing.T) {
	tests := []struct {
		stmt string
		code uint16
		msg  string
		want string
	}{
		{
			stmt: "SHOW ENGINE STATUS;",
			code: 1286,
			msg:  "Unknown storage engine 'STATUS'",
			want: "SHOW ENGINE InnoDB;",
		},
		{
			stmt: "SHOW ENGINE a zzz STATUS;",
			code: 1286,
			msg:  "Unknown storage engine 'a'",
			want: "SHOW ENGINE InnoDB zzz STATUS;",
		},
		{
			stmt: "SET NAMES a zzz;",
			code: 1115,
			msg:  "Unknown character set: 'a'",
			want: "SET NAMES utf8mb4 zzz;",
		},
		{
			stmt: "SET CHARSET BINLOG;",
			code: 1115,
			msg:  "Unknown character set: 'BINLOG'",
			want: "SET CHARSET utf8mb4;",
		},
		{
			stmt: "ALTER DATABASE COLLATE zzz utf8mb4_bin;",
			code: 1273,
			msg:  "Unknown collation: 'zzz'",
			want: "ALTER DATABASE COLLATE utf8mb4_bin utf8mb4_bin;",
		},
		{
			stmt: "ALTER DATABASE COLLATE = zzz utf8mb4_bin;",
			code: 1273,
			msg:  "Unknown collation: 'zzz'",
			want: "ALTER DATABASE COLLATE = utf8mb4_bin utf8mb4_bin;",
		},
		{
			stmt: "ALTER TABLE a ADD a CHAR ASCII COLLATE utf8mb4_bin zzz;",
			code: 1253,
			msg:  "COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'",
			want: "ALTER TABLE a ADD a CHAR CHARSET utf8mb4 COLLATE utf8mb4_bin zzz;",
		},
		{
			stmt: "ALTER TABLE a ADD a CHAR ASCII COLLATE BINARY zzz;",
			code: 1253,
			msg:  "COLLATION 'binary' is not valid for CHARACTER SET 'latin1'",
			want: "ALTER TABLE a ADD a CHAR CHARSET binary COLLATE BINARY zzz;",
		},
		{
			stmt: "DO a.SUBSTRING(a FROM a) zzz;",
			code: 1128,
			msg:  "Function 'a.SUBSTRING' is not defined",
			want: "DO mariadb_schema.SUBSTRING(a FROM a) zzz;",
		},
		{
			stmt: "ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a zzz;",
			code: 1128,
			msg:  "Function 'a' is not defined",
			want: "ALTER TABLE a ADD FULLTEXT(a) COMMENT 's' zzz;",
		},
		{
			stmt: "SHOW PACKAGE BODY CODE a zzz;",
			code: 1289,
			msg:  "The 'SHOW PROCEDURE|FUNCTION CODE' feature is disabled; you need MariaDB built with '--with-debug' to have it working",
			want: "SHOW CREATE PACKAGE BODY a zzz;",
		},
		{
			stmt: "DO@@GOTO.autocommit;",
			code: 1272,
			msg:  "Variable 'autocommit' is not a variable component (can't be used as XXXX.variable_name)",
			want: "DO@@omni_cache.key_buffer_size;",
		},
		{
			stmt: "SET STATEMENT.a = a FOR COMMIT;",
			code: 4081,
			msg:  "Unknown structured system variable or ROW routine variable 'STATEMENT'",
			want: "SET omni_cache.key_buffer_size = a FOR COMMIT;",
		},
		{
			stmt: "SET STATEMENT a.a = a zzz FOR COMMIT;",
			code: 4081,
			msg:  "Unknown structured system variable or ROW routine variable 'a'",
			want: "SET STATEMENT omni_cache.key_buffer_size = a zzz FOR COMMIT;",
		},
		{
			stmt: "SET@@autocommit.a = a zzz;",
			code: 4081,
			msg:  "Unknown structured system variable or ROW routine variable 'autocommit'",
			want: "SET@@omni_cache.key_buffer_size = a zzz;",
		},
		{
			stmt: "DO@@autocommit.a zzz;",
			code: 1193,
			msg:  "Unknown system variable 'a'",
			want: "DO@@omni_cache.key_buffer_size zzz;",
		},
		{
			stmt: "SET@@DEFAULT.a = a zzz;",
			code: 1193,
			msg:  "Unknown system variable 'a'",
			want: "SET@@omni_cache.key_buffer_size = a zzz;",
		},
		{
			stmt: "DO@@GLOBAL.zzz autocommit;",
			code: 1193,
			msg:  "Unknown system variable 'zzz'",
			want: "DO@@GLOBAL.autocommit autocommit;",
		},
		{
			stmt: "DO@@BIT zzz;",
			code: 1193,
			msg:  "Unknown system variable 'BIT'",
			want: "DO@@autocommit zzz;",
		},
		{
			stmt: "SET STATEMENT a = a, a = a zzz FOR COMMIT;",
			code: 1193,
			msg:  "Unknown system variable 'a'",
			want: "SET STATEMENT max_statement_time = a, a = a zzz FOR COMMIT;",
		},
		{
			stmt: "SET STATEMENT max_statement_time = a, a = a zzz FOR COMMIT;",
			code: 1193,
			msg:  "Unknown system variable 'a'",
			want: "SET STATEMENT max_statement_time = a, max_statement_time = a zzz FOR COMMIT;",
		},
		{
			stmt: "SET STATEMENT = a FOR COMMIT;",
			code: 1193,
			msg:  "Unknown system variable 'STATEMENT'",
			want: "SET autocommit = a FOR COMMIT;",
		},
		{
			stmt: "SET GLOBAL = a, CHARSET utf8mb4;",
			code: 1193,
			msg:  "Unknown system variable 'GLOBAL'",
			want: "SET autocommit = a, CHARSET utf8mb4;",
		},
		{
			stmt: "SET GLOBAL zzz = a;",
			code: 1193,
			msg:  "Unknown system variable 'zzz'",
			want: "SET GLOBAL autocommit = a;",
		},
		{
			stmt: "SET ROLE = a zzz;",
			code: 1193,
			msg:  "Unknown system variable 'ROLE'",
			want: "SET autocommit = a zzz;",
		},
	}
	for _, tc := range tests {
		got, ok := mariaDBStandIn(tc.stmt, &gomysql.MySQLError{Number: tc.code, Message: tc.msg})
		if !ok {
			t.Errorf("mariaDBStandIn(%q, %d) returned ok=false, want %q", tc.stmt, tc.code, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("mariaDBStandIn(%q, %d) =\n  %q\nwant:\n  %q", tc.stmt, tc.code, got, tc.want)
		}
	}
}
