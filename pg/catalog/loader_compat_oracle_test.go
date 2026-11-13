package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/pg/internal/testcontainer"
)

type loaderCompatOracle struct {
	db  *sql.DB
	ctx context.Context
}

// startLoaderCompatOracle returns t's connection to PostgreSQL 17: a replay
// of its tape, or the pinned server in a recording.
func startLoaderCompatOracle(t *testing.T) *loaderCompatOracle {
	t.Helper()
	return &loaderCompatOracle{db: testcontainer.Open(t), ctx: context.Background()}
}

func (o *loaderCompatOracle) execIsolated(t *testing.T, ddl string) {
	t.Helper()
	if err := o.execIsolatedErr(t, ddl); err != nil {
		t.Fatalf("PG17 rejected SQL that the compatibility corpus expects to be accepted: %v\nSQL:\n%s", err, ddl)
	}
}

func (o *loaderCompatOracle) execIsolatedErr(t *testing.T, ddl string) error {
	t.Helper()
	schema := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "_"))
	schema = strings.ReplaceAll(schema, "-", "_")
	if _, err := o.db.ExecContext(o.ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schema)); err != nil {
		return fmt.Errorf("PG cleanup failed: %w", err)
	}
	if _, err := o.db.ExecContext(o.ctx, fmt.Sprintf(`CREATE SCHEMA %q`, schema)); err != nil {
		return fmt.Errorf("PG schema setup failed: %w", err)
	}
	t.Cleanup(func() {
		_, _ = o.db.ExecContext(o.ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schema))
	})

	sql := fmt.Sprintf("SET search_path TO %q, public;\n%s", schema, ddl)
	_, err := o.db.ExecContext(o.ctx, sql)
	return err
}

func TestLoaderCompatPG17OracleAcceptsCorpus(t *testing.T) {
	oracle := startLoaderCompatOracle(t)
	for _, tc := range loaderCompatAcceptCases() {
		t.Run(tc.name, func(t *testing.T) {
			oracle.execIsolated(t, tc.sql)
			if _, err := LoadSQL(tc.sql); err != nil {
				t.Fatalf("LoadSQL rejected PG17-accepted SQL: %v\nSQL:\n%s", err, tc.sql)
			}
		})
	}
}

func TestLoaderCompatPG17OracleRejectsCorpus(t *testing.T) {
	oracle := startLoaderCompatOracle(t)
	for _, tc := range loaderCompatRejectCases() {
		t.Run(tc.name, func(t *testing.T) {
			if oracle.execIsolatedErr(t, tc.sql) == nil {
				t.Fatalf("PG17 accepted SQL that the reject corpus expects to be rejected\nSQL:\n%s", tc.sql)
			}
			if _, err := LoadSQL(tc.sql); err == nil {
				t.Fatalf("LoadSQL accepted SQL that the reject corpus expects to be rejected\nSQL:\n%s", tc.sql)
			}
		})
	}
}
