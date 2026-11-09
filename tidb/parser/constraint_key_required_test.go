package parser

import (
	"errors"
	"strings"
	"testing"
)

// TestTableConstraintKeyRequired pins that KEY is a required terminal after
// PRIMARY and FOREIGN in a table-level constraint. The table constraint is
// `PRIMARY KEY (key_part,...)` and `FOREIGN KEY [index_name] (col,...)
// reference_definition` (ConstraintKeywordOpt/Constraint in TiDB v8.5.5
// pkg/parser/parser.y), so `PRIMARY (col)` and `FOREIGN (col) REFERENCES ...`
// without KEY are rejected. These are the §3.6 negatives-corpus rows
// ConstraintElem#1/drop-2 and ConstraintElem#5/drop-2, each measured reject on
// pingcap/tidb v8.5.5. UNIQUE keeps its optional INDEX/KEY and stays accepted.
func TestTableConstraintKeyRequired(t *testing.T) {
	reject := []struct {
		sql string
		at  string // error points at the last occurrence of this text
	}{
		{"ALTER TABLE t ADD PRIMARY (a)", "(a)"},
		{"ALTER TABLE t ADD FOREIGN (a) REFERENCES p (id)", "(a)"},
		{"CREATE TABLE t (a INT, PRIMARY (a))", "(a))"},
		{"CREATE TABLE t (a INT, FOREIGN (a) REFERENCES p (id))", "(a)"},
	}
	for _, tc := range reject {
		t.Run("reject/"+tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want a *ParseError", tc.sql, err)
			}
			want := strings.LastIndex(tc.sql, tc.at)
			if pe.Position != want {
				t.Errorf("error Position = %d, want %d (%q): %v", pe.Position, want, tc.at, err)
			}
		})
	}

	// Valid forms must still parse: KEY present, CONSTRAINT symbols, and the
	// UNIQUE case whose INDEX/KEY is genuinely optional.
	accept := []string{
		"CREATE TABLE t (a INT, PRIMARY KEY (a))",
		"CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a))",
		"ALTER TABLE t ADD PRIMARY KEY (a)",
		"ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id)",
		"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES p (id)",
		"ALTER TABLE t ADD UNIQUE (a)",
		"ALTER TABLE t ADD UNIQUE KEY (a)",
	}
	for _, sql := range accept {
		t.Run("accept/"+sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q) error = %v, want nil", sql, err)
			}
		})
	}
}
