package v1review

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// columnCase is one column rule test: the findings, as described, that the
// rule reports for sql.
type columnCase struct {
	name string
	sql  string
	want []string
}

// describe renders a finding as "code table.name value @statement: anchor",
// with the text its range anchors, or "-" when the range is empty.
func describe(sql string, f v1review.Finding) string {
	anchor := "-"
	if f.Start != 0 || f.End != 0 {
		anchor = sql[f.Start:f.End]
	}
	return fmt.Sprintf("%s %s.%s %s @%d: %s", f.Code, f.Table, f.Name, f.Value, f.Statement, anchor)
}

func runColumnCases(t *testing.T, rule func([]ast.Node) []v1review.Finding, cases []columnCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, f := range rule(parse(t, tt.sql)) {
				got = append(got, describe(tt.sql, f))
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
