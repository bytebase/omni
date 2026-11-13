package catalog

import "testing"

// The loader compatibility corpus (loader_compat_test.go) held to LoadSQL.
// Its verdicts were measured on PostgreSQL 17 when this fork shared pg's
// oracle; Redshift runs no PostgreSQL 17, so they are the loader's own now,
// not a claim about what Redshift accepts.

func TestLoaderCompatAcceptsCorpus(t *testing.T) {
	for _, tc := range loaderCompatAcceptCases() {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := LoadSQL(tc.sql); err != nil {
				t.Fatalf("LoadSQL rejected a statement of the accept corpus: %v\nSQL:\n%s", err, tc.sql)
			}
		})
	}
}

func TestLoaderCompatRejectsCorpus(t *testing.T) {
	for _, tc := range loaderCompatRejectCases() {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := LoadSQL(tc.sql); err == nil {
				t.Fatalf("LoadSQL accepted a statement of the reject corpus\nSQL:\n%s", tc.sql)
			}
		})
	}
}
