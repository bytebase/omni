package parser

import "testing"

// TestInsertFormsMatchEngine holds insert_forms_test.go's verdicts to the server, so
// an image upgrade that changes an answer fails here rather than leaving a stale
// expectation behind.
func TestInsertFormsMatchEngine(t *testing.T) {
	o := startParserOracle(t)
	accepts := func(sql string) (bool, error) { return operatorVerdict(o.ctx, o.db, sql) }
	for _, c := range insertFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := accepts(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("MySQL 8.0 accepts %q = %v, this file records %v", c.sql, got, c.accept)
			}
		})
	}
}
