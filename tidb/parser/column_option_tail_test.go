package parser

import (
	"errors"
	"strings"
	"testing"
)

// TestColumnOptionTailRequired pins the tail of a column option/constraint to
// what TiDB v8.5.5 requires, closing four §3.6 negatives-corpus over-acceptances:
//
//   - COLUMN_FORMAT needs a format token (DEFAULT|FIXED|DYNAMIC); a bare
//     COLUMN_FORMAT is a 1064. (ColumnOption#18/drop-2, ColumnFormat rule.)
//   - GENERATED needs ALWAYS before AS; `GENERATED AS (expr)` is a 1064.
//     (GeneratedAlways#2/drop-2.)
//   - A trailing NOT after CHECK(expr) must begin NOT ENFORCED (or NOT NULL,
//     which the outer loop reads); a dangling NOT is a 1064.
//     (EnforcedOrNotOrNotNullOpt#1/drop-2.)
//   - DEFAULT does not take CURRENT_USER or CURRENT_ROLE: both are reserved
//     Username/Rolename keywords absent from DefaultValueExpr.
//     (NowSymOptionFraction#1/swap, NowSymOptionFraction#2/swap, CurdateSym#2/swap.)
//
// Each reject row is measured reject on pingcap/tidb v8.5.5. The accept rows are
// the legitimate neighbours that must keep parsing — the completeness guard
// against narrowing a real column option or DEFAULT value.
//
// Ref: ColumnOption, ColumnFormat, GeneratedAlways, EnforcedOrNotOrNotNullOpt,
// DefaultValueExpr in TiDB v8.5.5 pkg/parser/parser.y.
func TestColumnOptionTailRequired(t *testing.T) {
	reject := []struct {
		sql string
		at  string // error points at the last occurrence; "" means end of input
	}{
		{"ALTER TABLE ident ADD ident JSON COLUMN_FORMAT", ""},
		{"ALTER TABLE ident ADD ident JSON GENERATED AS(ident)", "AS"},
		{"ALTER TABLE ident ADD ident JSON CHECK(ident) NOT", ""},
		{"ALTER TABLE ident ADD ident JSON DEFAULT CURRENT_USER", "CURRENT_USER"},
		{"ALTER TABLE ident ADD ident JSON DEFAULT CURRENT_USER()", "CURRENT_USER"},
		{"ALTER TABLE ident ADD ident JSON DEFAULT CURRENT_ROLE()", "CURRENT_ROLE"},
	}
	for _, tc := range reject {
		t.Run("reject/"+tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want a *ParseError", tc.sql, err)
			}
			want := len(tc.sql)
			if tc.at != "" {
				want = strings.LastIndex(tc.sql, tc.at)
			}
			if pe.Position != want {
				t.Errorf("error Position = %d, want %d (%q): %v", pe.Position, want, tc.at, err)
			}
		})
	}

	// The legitimate forms must still parse: every COLUMN_FORMAT value, GENERATED
	// ALWAYS and the AS shorthand, CHECK with both ENFORCED tails and NOT NULL,
	// and the full span of DEFAULT value forms the swap mutations came from.
	accept := []string{
		"ALTER TABLE ident ADD ident JSON COLUMN_FORMAT DEFAULT",
		"ALTER TABLE ident ADD ident JSON COLUMN_FORMAT FIXED",
		"ALTER TABLE ident ADD ident JSON COLUMN_FORMAT DYNAMIC",
		"ALTER TABLE ident ADD ident INT GENERATED ALWAYS AS (ident)",
		"ALTER TABLE ident ADD ident INT GENERATED ALWAYS AS (ident) STORED",
		"ALTER TABLE ident ADD ident INT AS (ident)",
		"ALTER TABLE ident ADD ident INT CHECK (ident) ENFORCED",
		"ALTER TABLE ident ADD ident INT CHECK (ident) NOT ENFORCED",
		"ALTER TABLE ident ADD ident INT CHECK (ident) NOT NULL",
		"ALTER TABLE ident ADD ident INT DEFAULT 0",
		"ALTER TABLE ident ADD ident INT DEFAULT -1",
		"ALTER TABLE ident ADD ident JSON DEFAULT NULL",
		"ALTER TABLE ident ADD ident TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		"ALTER TABLE ident ADD ident TIMESTAMP DEFAULT CURRENT_TIMESTAMP()",
		"ALTER TABLE ident ADD ident TIMESTAMP DEFAULT CURRENT_TIMESTAMP(6)",
		"ALTER TABLE ident ADD ident TIMESTAMP DEFAULT LOCALTIME",
		"ALTER TABLE ident ADD ident TIMESTAMP DEFAULT NOW()",
		"ALTER TABLE ident ADD ident DATE DEFAULT CURRENT_DATE",
		"ALTER TABLE ident ADD ident DATE DEFAULT CURDATE()",
		"ALTER TABLE ident ADD ident TIMESTAMP DEFAULT (NOW())",
		"ALTER TABLE ident ADD ident INT DEFAULT NEXT VALUE FOR ident",
		"ALTER TABLE ident ADD ident INT DEFAULT NEXTVAL(ident)",
	}
	for _, sql := range accept {
		t.Run("accept/"+sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q) error = %v, want nil", sql, err)
			}
		})
	}
}
