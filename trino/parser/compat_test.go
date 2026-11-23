package parser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	enginerules "github.com/bytebase/omni/enginetest/rules"
	"github.com/bytebase/omni/trino/internal/testcontainer"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of Trino's SqlBase.g4 and its near-misses, with omni's verdict
// and Trino's (coarse accept/reject at parser pin judge revision 2; exact
// ParseError.Position byte offsets are pinned in TestErrorPositions). The test
// reads the file and starts no server. OMNI_RECORD=1 judges the rows that
// have no verdict yet on the pinned Trino.
//
// testdata/README.md says what Trino's judge needs.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		PinRole: "parser",
		Grammar: filepath.Join(td, "SqlBase.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
	enginerules.Check(t, enginerules.Suite{
		Compat: filepath.Join("testdata", "compat.tsv"),
		Rules:  filepath.Join(td, "rules.tsv"),
		Root:   filepath.Join("..", ".."),
	})
}

// corpusSuites are the hand-written corpora that are suite: rows of
// compat.tsv, in the order the rows were added. Each held omni's Parse to
// Trino's verdict on a live server; the rows hold it to the recorded one.
var corpusSuites = []struct {
	name    string
	corpora [][]string
}{
	{"cte", [][]string{cteAcceptCorpus, cteRejectCorpus}},
	{"relation", [][]string{relationAcceptCorpus, relationRejectCorpus}},
	{"select", [][]string{selectAcceptCorpus, selectRejectCorpus}},
	{"setop", [][]string{setopAcceptCorpus, setopRejectCorpus}},
	{"window", [][]string{windowAcceptCorpus, windowRejectCorpus}},
	{"ddl", [][]string{ddlOracleCorpus}},
	{"dml", [][]string{dmlOracleCorpus}},
	{"dcl-tcl", [][]string{dclTclOracleCorpus}},
	{"match-recognize", [][]string{matchRecognizeOracleCorpus}},
	{"routines", [][]string{routinesOracleCorpus}},
}

// TestCompatibilitySuitesAreTheCorpora fails while the suite: rows of
// compat.tsv are not the corpora of corpusSuites: each statement once, under
// the first suite that holds it. A statement added to a corpus is added as a
// row with ? in both verdict cells, and judged by a recording.
func TestCompatibilitySuitesAreTheCorpora(t *testing.T) {
	f, err := compat.ReadFile(filepath.Join("testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	var want, got []string
	seen := map[string]bool{}
	for _, s := range corpusSuites {
		for _, corpus := range s.corpora {
			for _, sql := range corpus {
				if !seen[sql] {
					seen[sql] = true
					want = append(want, "suite:"+s.name+"\t"+sql)
				}
			}
		}
	}
	for _, r := range f.Rows {
		if strings.HasPrefix(r.Rule, "suite:") {
			got = append(got, r.Rule+"\t"+r.Statement)
		}
	}
	if slices.Equal(got, want) {
		return
	}
	for i := 0; i < max(len(got), len(want)); i++ {
		g, w := "(none)", "(none)"
		if i < len(got) {
			g = got[i]
		}
		if i < len(want) {
			w = want[i]
		}
		if g != w {
			t.Fatalf("suite row %d of testdata/compat.tsv is %q, and the corpora want %q. A new statement is a row with ? in both verdict cells,\n  ?\t?\t-\tsuite:<name>\t<statement>\njudged by a recording:\n  OMNI_RECORD=1 go test ./trino/parser/ -run '^TestCompatibility$'", i+1, g, w)
		}
	}
}

// omniVerdict reads the row as Trino's statement endpoint does: one
// statement, with nothing after it, returning a coarse accept/reject verdict
// for judge revision 2. Parse reads a script, so a ';' ends the statement
// where Trino's grammar has no place for it, and a text with no statement is
// refused where Trino finds none. Exact ParseError.Position byte offsets are
// pinned separately in TestErrorPositions.
func omniVerdict(sql string) string {
	var segs []Segment
	for _, seg := range Split(sql) {
		if !seg.Empty() {
			segs = append(segs, seg)
		}
	}
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	if len(segs) == 0 || segs[0].ByteEnd < len(sql) {
		return compat.Reject
	}
	return compat.Accept
}

// verifySharedServerPin verifies that a dedicated pin role (such as "parser")
// pins the exact same Trino image, digest, and tag as the default role before
// reusing testcontainer.Judge's default container.
func verifySharedServerPin(pinsPath string, pin enginetest.Pin) error {
	if pin.Role == testcontainer.Default {
		return nil
	}
	serverPin, err := enginetest.LoadPin(pinsPath, testcontainer.Default)
	if err != nil {
		return err
	}
	if serverPin.Image != pin.Image || serverPin.Digest != pin.Digest || serverPin.Tag != pin.Tag {
		return fmt.Errorf("%s pin (%s %s %s) differs from %s server pin (%s %s %s)",
			pin.Role, pin.Image, pin.Digest, pin.Tag,
			testcontainer.Default, serverPin.Image, serverPin.Digest, serverPin.Tag)
	}
	return nil
}

// trinoVerdictOnly reads a Trino server's answer as a coarse accept/reject
// verdict for judge revision 2. It delegates validation to compat.TrinoVerdict
// so non-server errors, empty errorName, and out-of-range error locations
// still fail the recording, and strips any @<byte> offset from rejections.
func trinoVerdictOnly(statement string, err error, read func(error) (compat.TrinoError, bool)) (string, error) {
	v, verr := compat.TrinoVerdict(statement, err, read)
	if verr != nil {
		return "", verr
	}
	if strings.HasPrefix(v, compat.Reject+"@") {
		return compat.Reject, nil
	}
	return v, nil
}

// compatJudge asks the pinned Trino whether its parser accepts each
// statement, by running it and reading the answer with trinoVerdictOnly. A
// statement that parses may run, on the memory catalog the session defaults
// to; nothing it leaves changes whether a later one parses. It runs only in a
// recording, outside any tape.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	if err := verifySharedServerPin(testcontainer.Pins, pin); err != nil {
		t.Fatalf("%v", err)
	}
	o := testcontainer.Judge(t)
	return func(sql string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return trinoVerdictOnly(sql, o.Run(ctx, sql), trinoError)
	}
}

// trinoError reads a Trino answer for compat.TrinoVerdict.
func trinoError(err error) (compat.TrinoError, bool) {
	var qe *testcontainer.QueryError
	if !errors.As(err, &qe) {
		return compat.TrinoError{}, false
	}
	return compat.TrinoError{Name: qe.Name, Line: qe.Line, Column: qe.Column}, true
}

func TestTrinoVerdictOnly(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
		err  error
		want string
	}{
		{"ran", "SELECT 1", nil, compat.Accept},
		{"semantic_error", "SELECT * FROM missing_tbl", &testcontainer.QueryError{Name: "TABLE_NOT_FOUND", Line: 1, Column: 15}, compat.Accept},
		{"incompatible_client", "START TRANSACTION", &testcontainer.QueryError{Name: "INCOMPATIBLE_CLIENT"}, compat.Accept},
		{"syntax_error_with_location", "SELECT 1 FROM FROM", &testcontainer.QueryError{Name: "SYNTAX_ERROR", Line: 1, Column: 15}, compat.Reject},
		{"syntax_error_multiline_utf8", "SELECT '😀'\nFROM FROM", &testcontainer.QueryError{Name: "SYNTAX_ERROR", Line: 2, Column: 6}, compat.Reject},
		{"syntax_error_no_location", "DROP TABLE a.b.c.d", &testcontainer.QueryError{Name: "SYNTAX_ERROR"}, compat.Reject},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := trinoVerdictOnly(tc.sql, tc.err, trinoError)
			if err != nil {
				t.Fatalf("trinoVerdictOnly(%q) unexpected error: %v", tc.sql, err)
			}
			if got != tc.want {
				t.Errorf("trinoVerdictOnly(%q) = %q, want %q", tc.sql, got, tc.want)
			}
		})
	}

	for _, tc := range []struct {
		name    string
		sql     string
		err     error
		wantErr string
	}{
		{"non_server_error", "SELECT 1", errors.New("connection refused"), "connection refused"},
		{"empty_error_name", "SELECT 1", &testcontainer.QueryError{Line: 1, Column: 1}, "no errorName"},
		{"out_of_range_location", "SELECT (1", &testcontainer.QueryError{Name: "SYNTAX_ERROR", Line: 1, Column: 11}, "outside the statement"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := trinoVerdictOnly(tc.sql, tc.err, trinoError)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("trinoVerdictOnly(%q) = (%q, %v), want error containing %q", tc.sql, got, err, tc.wantErr)
			}
		})
	}
}

func TestOmniVerdict(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
		want string
	}{
		{"valid_single_stmt", "SELECT 1.5", compat.Accept},
		{"syntax_error", "SELECT 1.5 IN(1.5 zzz)", compat.Reject},
		{"too_many_dots", "DROP TABLE a.b.c.d", compat.Reject},
		{"empty_input", "", compat.Reject},
		{"comment_only", "-- just a comment", compat.Reject},
		{"trailing_semicolon", "SELECT 1;", compat.Reject},
		{"two_statements", "SELECT 1; SELECT 2", compat.Reject},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := omniVerdict(tc.sql); got != tc.want {
				t.Errorf("omniVerdict(%q) = %q, want %q", tc.sql, got, tc.want)
			}
		})
	}
}

func TestVerifySharedServerPin(t *testing.T) {
	pinsPath := filepath.Join(t.TempDir(), "pins.tsv")
	content := "default\ttrinodb/trino:482\tsha256:aaa\t482\nparser\ttrinodb/trino:482\tsha256:aaa\t482\t2\n"
	if err := os.WriteFile(pinsPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := verifySharedServerPin(pinsPath, enginetest.Pin{Role: "default", Image: "other", Digest: "sha256:zzz", Tag: "999"}); err != nil {
		t.Errorf("default role should skip check, got %v", err)
	}
	if err := verifySharedServerPin(pinsPath, enginetest.Pin{Role: "parser", Image: "trinodb/trino:482", Digest: "sha256:aaa", Tag: "482", Judge: "2"}); err != nil {
		t.Errorf("matching parser pin returned error: %v", err)
	}
	for _, bad := range []struct {
		name string
		pin  enginetest.Pin
	}{
		{"image_mismatch", enginetest.Pin{Role: "parser", Image: "trinodb/trino:483", Digest: "sha256:aaa", Tag: "482", Judge: "2"}},
		{"digest_mismatch", enginetest.Pin{Role: "parser", Image: "trinodb/trino:482", Digest: "sha256:bbb", Tag: "482", Judge: "2"}},
		{"tag_mismatch", enginetest.Pin{Role: "parser", Image: "trinodb/trino:482", Digest: "sha256:aaa", Tag: "483", Judge: "2"}},
	} {
		t.Run(bad.name, func(t *testing.T) {
			if err := verifySharedServerPin(pinsPath, bad.pin); err == nil {
				t.Errorf("verifySharedServerPin(%+v) = nil, want mismatch error", bad.pin)
			}
		})
	}
}
