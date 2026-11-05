package parser

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestExecFallback_OceanBase proves execFallbackCases and execVersionCases
// against a live OceanBase 4.3.5 (DSN-gated, see oceanBaseOracle): the verdict
// through PREPARE, and for an accepted query, the rows it returns, which are
// the rows of the text the case says it reads.
func TestExecFallback_OceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, cases := range [][]execFallbackCase{execFallbackCases, execVersionCases} {
		for _, c := range cases {
			t.Run(c.sql, func(t *testing.T) {
				if got := prepareAcceptsOceanBase(t, ctx, db, c.sql); got != c.accept {
					t.Fatalf("OceanBase accepts=%v, the case says %v", got, c.accept)
				}
				if _, err := Parse(c.sql); (err == nil) != c.accept {
					t.Errorf("DIVERGENCE omni accepts=%v (%v), OceanBase accepts=%v", err == nil, err, c.accept)
				}
				if !c.accept || (c.reads != "" && !strings.HasPrefix(c.reads, "select")) {
					return
				}
				got := oceanBaseRows(t, ctx, db, c.sql)
				want := "[]"
				if c.reads != "" {
					want = oceanBaseRows(t, ctx, db, c.reads)
				}
				if got != want {
					t.Errorf("OceanBase returns %s, and %s for %q", got, want, c.reads)
				}
			})
		}
	}
}

// TestExecFallbackScript_OceanBase proves TestExecFallbackPerStatement's
// premise against a live OceanBase 4.3.5: in a script, each statement is read
// again on its own.
func TestExecFallbackScript_OceanBase(t *testing.T) {
	oceanBaseOracle(t)
	db, err := sql.Open("mysql", os.Getenv(oceanBaseDSNEnv)+"?multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	for _, c := range []struct{ sql, want string }{
		{"select 1 /*! + 1 */; select 1 /*! from */", "[[2]] [[1]]"},
		{"select 1 /*! from */; select 1 /*! + 1 */", "[[1]] [[2]]"},
		{"select 1; /*! select from */; select 2", "[[1]] [[2]]"},
		{"select 1 /*! from ; select 2 */", "[[1]]"},
	} {
		if got := oceanBaseRows(t, ctx, db, c.sql); got != c.want {
			t.Errorf("OceanBase returns %s for %q, want %s", got, c.sql, c.want)
		}
	}
}

// oceanBaseRows runs query and returns the rows of every result set that has
// columns, as "[[a b] [c d]] [[e]]".
func oceanBaseRows(t *testing.T, ctx context.Context, db *sql.DB, query string) string {
	t.Helper()
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		t.Fatalf("OceanBase refused %q: %v", query, err)
	}
	defer rows.Close()
	var sets []string
	for {
		cols, err := rows.Columns()
		if err != nil {
			t.Fatalf("%q: %v", query, err)
		}
		var set []string
		for rows.Next() {
			raw := make([]sql.RawBytes, len(cols))
			ptrs := make([]any, len(cols))
			for i := range raw {
				ptrs[i] = &raw[i]
			}
			if err := rows.Scan(ptrs...); err != nil {
				t.Fatalf("%q: %v", query, err)
			}
			cells := make([]string, len(raw))
			for i, r := range raw {
				cells[i] = string(r)
			}
			set = append(set, fmt.Sprintf("%v", cells))
		}
		if err := rows.Err(); err != nil {
			t.Fatalf("%q: %v", query, err)
		}
		if len(cols) > 0 {
			sets = append(sets, "["+strings.Join(set, " ")+"]")
		}
		if !rows.NextResultSet() {
			break
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("%q: %v", query, err)
	}
	if len(sets) == 0 {
		return "[]"
	}
	return strings.Join(sets, " ")
}
