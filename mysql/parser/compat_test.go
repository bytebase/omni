package parser

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of MySQL's sql_yacc.yy and its near-misses, with omni's verdict
// and MySQL's. It reads the file and starts no server. The grammar is under
// the GPL, so the repo does not hold it: the test downloads it once
// (../testdata/sources.tsv). OMNI_RECORD=1 judges the rows that have no
// verdict yet on the pinned MySQL.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		PinRole: "parser",
		Grammar: enginetest.Fetched(t, filepath.Join("..", "testdata", "sources.tsv"), "sql_yacc.yy"),
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict is what omni reads syntactically, and it reads a name the server
// looks up while it parses the way the judge does (enginetest/compat's
// lookupErrorCodes): `ALTER DATABASE COLLATE utf8mb4` is text MySQL reads and
// refuses on the name, so the row stays accept/accept, and the parser's Code
// says which name it refused. Without that, asking omni the name would count as
// a rejection and the rows would read as a regression.
func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		var pe *ParseError
		if errors.As(err, &pe) && pe.Code != 0 {
			return compat.Accept
		}
		return compat.Reject
	}
	return compat.Accept
}

// compatPrepareTimeout bounds one PREPARE, so a statement the server hangs on
// gets no verdict and the recording goes on.
const compatPrepareTimeout = time.Minute

// compatJudge asks the pinned MySQL, through PREPARE, whether it accepts
// each statement's syntax. The server checks its digest and release when it
// starts.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	// testcontainer registers the shared MySQL 8.0.46 server under the default
	// role; when a dedicated pin role (such as parser) shares that release,
	// verify it pins the exact same server before reusing it.
	serverPin := pin
	if pin.Role != testcontainer.Default {
		var err error
		serverPin, err = enginetest.LoadPin(testcontainer.Pins, testcontainer.Default)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if serverPin.Image != pin.Image || serverPin.Digest != pin.Digest || serverPin.Tag != pin.Tag {
			t.Fatalf("%s pin (%s %s %s) differs from default server pin (%s %s %s)",
				pin.Role, pin.Image, pin.Digest, pin.Tag, serverPin.Image, serverPin.Digest, serverPin.Tag)
		}
	}
	ctx := context.Background()
	db := testcontainer.Judge(t, serverPin, "compat")
	prepare := func(text string) error {
		ctx, cancel := context.WithTimeout(ctx, compatPrepareTimeout)
		defer cancel()
		_, err := db.ExecContext(ctx, "PREPARE omni_compat FROM ?", text)
		return err
	}
	// The verdict reads 1064 as the parser's rejection, so a release that
	// answers a syntax error otherwise must fail here and not be misread.
	if code, ok := mysqlErrorCode(prepare("SELECT FROM")); !ok || code != 1064 {
		t.Fatalf("MySQL no longer answers PREPARE of a syntax error with 1064: %v", prepare("SELECT FROM"))
	}
	// Judge revision 2 re-asks WITH PARSER's 1128 plugin lookup with MySQL's
	// built-in ngram full-text parser plugin; verify ngram is loaded and
	// trailing syntax errors still answer 1064.
	if v, err := compatJudgeVerdict("ALTER TABLE t ADD FULLTEXT(c) WITH PARSER ngram;", prepare); err != nil || v != compat.Accept {
		t.Fatalf("MySQL no longer accepts built-in WITH PARSER ngram: verdict=%q err=%v", v, err)
	}
	if code, ok := mysqlErrorCode(prepare("ALTER TABLE t ADD FULLTEXT(c) WITH PARSER ngram zzz;")); !ok || code != 1064 {
		t.Fatalf("MySQL no longer answers WITH PARSER ngram trailing syntax error with 1064: %v", prepare("ALTER TABLE t ADD FULLTEXT(c) WITH PARSER ngram zzz;"))
	}
	return func(stmt string) (string, error) {
		return compatJudgeVerdict(stmt, prepare)
	}
}

// builtInFulltextParser is a full-text parser plugin shipped ready by default
// in MySQL 8.0, used to re-ask WITH PARSER statements that abort with 1128
// (ER_FUNCTION_NOT_DEFINED) on an unloaded plugin name.
const builtInFulltextParser = "ngram"

// fulltextParserPluginRE matches `WITH PARSER <IDENT_sys>` in fulltext_index_option.
var fulltextParserPluginRE = regexp.MustCompile("(?i)\\b(WITH\\s+PARSER\\s+)(`(?:``|[^`])*`|[A-Za-z0-9_$]+)")

// replaceFulltextParserPlugin substitutes plugin for any other parser name in
// a WITH PARSER clause. It reports whether the statement changed.
func replaceFulltextParserPlugin(stmt, plugin string) (string, bool) {
	changed := false
	rewritten := fulltextParserPluginRE.ReplaceAllStringFunc(stmt, func(match string) string {
		parts := fulltextParserPluginRE.FindStringSubmatch(match)
		if len(parts) != 3 || strings.EqualFold(parts[2], plugin) {
			return match
		}
		changed = true
		return parts[1] + plugin
	})
	return rewritten, changed
}

// compatJudgeVerdict reads PREPARE's answer for stmt. MySQL's bison action for
// fulltext_index_option (`WITH PARSER IDENT_sys`) checks plugin_is_ready while
// parsing and aborts with 1128 (ER_FUNCTION_NOT_DEFINED) before reading past
// the clause; re-asking with builtInFulltextParser lets MySQL parse the full
// statement so valid clauses accept and trailing near-misses reject with 1064.
func compatJudgeVerdict(stmt string, prepare func(string) error) (string, error) {
	first := prepare(stmt)
	if code, ok := mysqlErrorCode(first); ok && code == 1128 {
		if rewritten, replaced := replaceFulltextParserPlugin(stmt, builtInFulltextParser); replaced {
			stmt = rewritten
			first = prepare(stmt)
		}
	}
	accepts, err := compat.PrepareVerdict(first, mysqlErrorCode, func() error {
		return prepare(stmt + compat.ProbeTail)
	})
	switch {
	case err != nil:
		return "", err
	case accepts:
		return compat.Accept, nil
	default:
		return compat.Reject, nil
	}
}

// mysqlErrorCode returns the server's error number, and false for an error
// the server did not send.
func mysqlErrorCode(err error) (uint16, bool) {
	var me *gomysql.MySQLError
	if !errors.As(err, &me) {
		return 0, false
	}
	return me.Number, true
}

func TestCompatJudgeFulltextParserRewrite(t *testing.T) {
	prepare := func(text string) error {
		switch text {
		case "ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) WITH PARSER utf8mb4;",
			"ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) WITH PARSER utf8mb4 zzz;":
			return &gomysql.MySQLError{Number: 1128, Message: "Function 'utf8mb4' is not defined"}
		case "ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) WITH PARSER ngram;":
			return &gomysql.MySQLError{Number: 1146, Message: "Table 'test.utf8mb4' doesn't exist"}
		case "ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) WITH PARSER ngram;" + compat.ProbeTail,
			"ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) WITH PARSER ngram zzz;":
			return &gomysql.MySQLError{Number: 1064, Message: "syntax error"}
		default:
			t.Fatalf("unexpected PREPARE text: %q", text)
			return nil
		}
	}

	got, err := compatJudgeVerdict("ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) WITH PARSER utf8mb4;", prepare)
	if err != nil || got != compat.Accept {
		t.Fatalf("valid WITH PARSER: got (%q, %v), want (%q, nil)", got, err, compat.Accept)
	}
	got, err = compatJudgeVerdict("ALTER TABLE utf8mb4 ADD FULLTEXT(utf8mb4) WITH PARSER utf8mb4 zzz;", prepare)
	if err != nil || got != compat.Reject {
		t.Fatalf("trailing near-miss WITH PARSER: got (%q, %v), want (%q, nil)", got, err, compat.Reject)
	}
	if _, replaced := replaceFulltextParserPlugin("ALTER TABLE t ADD FULLTEXT(c) WITH PARSER ngram;", builtInFulltextParser); replaced {
		t.Fatal("already-ngram WITH PARSER should not report replaced")
	}
	if _, replaced := replaceFulltextParserPlugin("SELECT nosuch_udf();", builtInFulltextParser); replaced {
		t.Fatal("statement without WITH PARSER should not report replaced")
	}
}
