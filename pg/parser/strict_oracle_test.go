package parser

import "testing"

// syntaxCase is one statement and whether PostgreSQL's grammar accepts it.
type syntaxCase struct {
	sql    string
	accept bool
}

// assertSyntaxMatchesPG checks every case against PostgreSQL 17 and then
// against Parse. PG accepts a statement when it reports anything other than
// SQLSTATE 42601: the statements run without fixtures, so an accepted one
// usually fails name resolution instead. Checking PG first keeps the
// expectation tied to the engine rather than to memory.
func assertSyntaxMatchesPG(t *testing.T, cases []syntaxCase) {
	t.Helper()
	o := startFirstSetOracle(t)
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			if pgAccept := o.probe(t, c.sql) == probeAccept; pgAccept != c.accept {
				t.Fatalf("postgres accept = %v, want %v", pgAccept, c.accept)
			}
			_, err := Parse(c.sql)
			if omniAccept := err == nil; omniAccept != c.accept {
				t.Fatalf("omni accept = %v, want %v (err = %v)", omniAccept, c.accept, err)
			}
		})
	}
}
