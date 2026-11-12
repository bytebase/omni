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
// It runs on the package's own OceanBase (see testcontainer.OceanBase):
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
// different schema behind for each later case. The statement travels as the
// PREPARE's argument, which the connection writes into it as a string
// literal, and a PREPARE of the same name replaces the last one: the question
// is one request, which keeps a sweep's tape small.
//
// PREPARE either carries the statement, refuses to carry it, or rejects its
// syntax with 1064, and only the last of those is the parser's "no". Anything
// that is not a server error -- a dropped connection, a timeout -- is no
// verdict at all and fails the test rather than counting as acceptance.
func prepareAcceptsOceanBase(t *testing.T, ctx context.Context, db *sql.DB, stmt string) bool {
	t.Helper()
	_, err := db.ExecContext(ctx, "PREPARE omni_alter_clause_stmt FROM ?", stmt)
	if err == nil {
		return true
	}
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		t.Fatalf("no verdict for %q: %v", stmt, err)
	}
	return myErr.Number != 1064
}
