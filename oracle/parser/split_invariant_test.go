package parser

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestSplitInvariantsForCuratedScripts(t *testing.T) {
	for _, sql := range []string{
		"",
		";;;",
		"SELECT 1 FROM dual;",
		"SELECT 'a;b' FROM dual; SELECT q'[c;d]' FROM dual;",
		"SELECT 1 /* ; */ FROM dual; -- ;\nSELECT 2 FROM dual;",
		"BEGIN NULL; END;\n/\nSELECT 1 FROM dual;",
		"CREATE PACKAGE BODY pkg IS\nPROCEDURE p IS BEGIN NULL; END;\nEND pkg;\n/\n",
		"SET DEFINE OFF\nPROMPT before ;\nSELECT 1 FROM dual;\nSPOOL out.log\nSELECT 2 FROM dual;",
		"/\n/\nSELECT 1 FROM dual\n/\n/",
		"SELECT 'unterminated; SELECT 2 FROM dual;",
		"CREATE PROCEDURE p IS BEGIN NULL;",
		string(allBytes()),
		"SELECT 1 FROM dual;\x00\x01SELECT 2 FROM dual;",
	} {
		t.Run(fmt.Sprintf("%q", sql), func(t *testing.T) {
			assertSplitInvariants(t, sql)
		})
	}
}

func TestSplitInvariantsForGeneratedScripts(t *testing.T) {
	commands := []string{
		"",
		"SET DEFINE OFF\n",
		"PROMPT preparing ; script\n",
		"SPOOL install.log\n",
		"WHENEVER SQLERROR EXIT SQL.SQLCODE ROLLBACK\n",
		"REM ignored ; semicolon\n",
		"@preflight.sql\n",
	}
	literals := []string{
		"'plain'",
		"'semi;colon'",
		"q'[semi;colon]'",
		`"quoted;identifier"`,
		"/* block ; comment */ 1",
	}
	plsql := []string{
		"BEGIN NULL; END;",
		"BEGIN IF 1 = 1 THEN NULL; END IF; END;",
		"DECLARE v NUMBER := 1; BEGIN v := v + 1; END;",
		"CREATE PROCEDURE p IS BEGIN NULL; NULL; END;",
		"CREATE TRIGGER trg BEFORE INSERT ON t BEGIN NULL; END;",
	}

	for i, command := range commands {
		for j, literal := range literals {
			for k, block := range plsql {
				sql := command +
					fmt.Sprintf("SELECT %s FROM dual;\n", literal) +
					block + "\n/\n" +
					commands[(i+j+k)%len(commands)] +
					"SELECT 10 / 2 FROM dual;"
				t.Run(fmt.Sprintf("command_%d_literal_%d_plsql_%d", i, j, k), func(t *testing.T) {
					assertSplitInvariants(t, sql)
				})
			}
		}
	}
}

func TestSplitSoftFailInvariantsForGeneratedInvalidInputs(t *testing.T) {
	fragments := []string{
		"SELECT 1 FROM dual;",
		"BEGIN NULL;",
		"CREATE PROCEDURE p IS BEGIN",
		"SELECT 'unterminated",
		"SELECT q'[unterminated",
		"SELECT \"unterminated",
		"/* unterminated",
		"\x00\x01\x02",
		";\n/\n",
	}

	rng := rand.New(rand.NewSource(20260509))
	for i := 0; i < 128; i++ {
		var b strings.Builder
		for j := 0; j < 1+rng.Intn(8); j++ {
			b.WriteString(fragments[rng.Intn(len(fragments))])
			switch rng.Intn(4) {
			case 0:
				b.WriteByte('\n')
			case 1:
				b.WriteByte(';')
			case 2:
				b.WriteString("\n/\n")
			}
		}
		sql := b.String()
		t.Run(fmt.Sprintf("generated_%03d", i), func(t *testing.T) {
			assertSplitInvariants(t, sql)
		})
	}
}

func FuzzSplitInvariants(f *testing.F) {
	for _, seed := range []string{
		"",
		"SELECT 1 FROM dual;",
		"SELECT 'a;b' FROM dual; SELECT 2 FROM dual;",
		"BEGIN NULL; END;\n/\nSELECT 1 FROM dual;",
		"SET DEFINE OFF\nSELECT 1 FROM dual;",
		"SELECT 'unterminated; SELECT 2 FROM dual;",
		string(allBytes()),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, sql string) {
		assertSplitInvariants(t, sql)
	})
}

func assertSplitInvariants(t *testing.T, sql string) {
	t.Helper()

	var segments []Segment
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Split panicked for input %q: %v", sql, r)
			}
		}()
		segments = Split(sql)
	}()

	prevEnd := 0
	for i, seg := range segments {
		if seg.ByteStart < prevEnd {
			t.Fatalf("segment[%d] starts at %d before previous end %d", i, seg.ByteStart, prevEnd)
		}
		if seg.ByteStart < 0 || seg.ByteEnd < seg.ByteStart || seg.ByteEnd > len(sql) {
			t.Fatalf("segment[%d] invalid range [%d,%d] for input length %d", i, seg.ByteStart, seg.ByteEnd, len(sql))
		}
		if seg.Text != sql[seg.ByteStart:seg.ByteEnd] {
			t.Fatalf("segment[%d] text %q does not match input range %q", i, seg.Text, sql[seg.ByteStart:seg.ByteEnd])
		}
		if seg.Empty() {
			t.Fatalf("segment[%d] is empty: %q", i, seg.Text)
		}

		resplit := Split(seg.Text)
		if len(resplit) != 1 {
			t.Fatalf("segment[%d] resplit into %d segments: %#v", i, len(resplit), resplit)
		}
		if resplit[0].Text != seg.Text {
			t.Fatalf("segment[%d] resplit text = %q, want %q", i, resplit[0].Text, seg.Text)
		}
		prevEnd = seg.ByteEnd
	}
}

func allBytes() []byte {
	b := make([]byte, 256)
	for i := range b {
		b[i] = byte(i)
	}
	return b
}
