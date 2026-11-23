package pg

import (
	"errors"
	"flag"
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/pg/parser"
)

// FuzzParse holds the parser to the engine contract on any input: Split,
// Tokenize and Parse never panic or overflow the stack, Split's segments
// lie inside the script, and Parse fails only with a *parser.ParseError at
// a byte offset inside the text it was given. `go test` replays the seeds
// below and testdata/fuzz/FuzzParse; `go test -fuzz` also seeds from every
// statement of parser/testdata/compat.tsv.
func FuzzParse(f *testing.F) {
	for _, s := range []string{"", "SELECT 1;", "SELECT (((1)));", "CREATE TABLE t (a int);"} {
		f.Add(s)
	}
	if fuzzing() {
		file, err := compat.ReadFile(filepath.Join("parser", "testdata", "compat.tsv"))
		if err != nil {
			f.Fatal(err)
		}
		for _, r := range file.Rows {
			f.Add(r.Statement)
		}
	}
	f.Fuzz(func(t *testing.T, sql string) {
		parser.Tokenize(sql)
		for _, seg := range Split(sql) {
			if seg.ByteStart < 0 || seg.ByteStart > seg.ByteEnd || seg.ByteEnd > len(sql) {
				t.Fatalf("segment [%d,%d) outside a %d-byte script", seg.ByteStart, seg.ByteEnd, len(sql))
			}
			if seg.Empty() {
				continue
			}
			_, err := parser.Parse(seg.Text)
			checkParseError(t, err, len(seg.Text))
		}
		_, err := Parse(sql)
		checkParseError(t, err, len(sql))
	})
}

// fuzzing reports whether the test binary runs with -test.fuzz.
func fuzzing() bool {
	fl := flag.Lookup("test.fuzz")
	return fl != nil && fl.Value.String() != ""
}

func checkParseError(t *testing.T, err error, n int) {
	t.Helper()
	if err == nil {
		return
	}
	var pe *parser.ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("Parse failed with %T, not *ParseError: %v", err, err)
	}
	if pe.Position < -1 || pe.Position > n {
		t.Fatalf("ParseError at %d, outside a %d-byte text: %v", pe.Position, n, err)
	}
}
