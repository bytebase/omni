package parser

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestAlterTableClauseFormsOceanBaseOracle asserts omni agrees with a live
// OceanBase on alterTableClauseCases.
//
// It runs on the package's own OceanBase (see sharedserver.PrivateOceanBase):
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
//
// The expectations in alterTableClauseCases are what
// oceanbase/oceanbase-ce:4.3.5-lts (4.3.5.6) answered on 2026-11-04, so an
// image upgrade that changes a verdict fails here instead of passing silently.
func TestAlterTableClauseFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, c := range alterTableClauseCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts := prepareAcceptsOceanBase(t, ctx, db, c.sql)
			if accepts != c.accept {
				t.Fatalf("OceanBase accepts=%v, want %v for %q", accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v",
					omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}

// prepareAcceptsOceanBase asks the server whether sql parses, without running
// it: the cases create and alter tables, and executing them would leave a
// different schema behind for each later case. The statement is bound to a
// user variable, so its own quoting stays the server's business, and the two
// statements share one connection because a user variable lives on the
// connection that set it.
//
// PREPARE either carries the statement, refuses to carry it, or rejects its
// syntax with 1064, and only the last of those is the parser's "no". Anything
// that is not a server error -- a dropped connection, a timeout -- is no
// verdict at all and fails the test rather than counting as acceptance.
func prepareAcceptsOceanBase(t *testing.T, ctx context.Context, db *sql.DB, stmt string) bool {
	t.Helper()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("taking a connection for %q: %v", stmt, err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "SET @omni_alter_clause = ?", stmt); err != nil {
		t.Fatalf("binding %q: %v", stmt, err)
	}
	_, err = conn.ExecContext(ctx, "PREPARE omni_alter_clause_stmt FROM @omni_alter_clause")
	if err == nil {
		_, _ = conn.ExecContext(ctx, "DEALLOCATE PREPARE omni_alter_clause_stmt")
		return true
	}
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		t.Fatalf("no verdict for %q: %v", stmt, err)
	}
	return myErr.Number != 1064
}
