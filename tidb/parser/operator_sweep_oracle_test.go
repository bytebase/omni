package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// operatorVerdict asks the engine to parse stmt without running it. A 1064 is the
// parser's refusal; any other server error means the text parsed. Anything that is
// not a server error is no verdict, and comes back as an error rather than passing
// a case that expects either answer.
//
// The statement travels as a bound parameter and through one held connection,
// because the variable that carries it lives on the connection that set it.
func operatorVerdict(ctx context.Context, db *sql.DB, stmt string) (bool, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.ExecContext(ctx, "SET @omni_sweep = ?", stmt); err != nil {
		return false, err
	}
	_, err = conn.ExecContext(ctx, "PREPARE omni_sweep_stmt FROM @omni_sweep")
	if err == nil {
		_, _ = conn.ExecContext(ctx, "DEALLOCATE PREPARE omni_sweep_stmt")
		return true, nil
	}
	var myErr *gomysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number != 1064, nil
	}
	return false, err
}

// TestOperatorFormsMatchEngine holds operator_forms_test.go's verdicts to the
// server, so an image upgrade that changes an answer fails here rather than leaving
// a stale expectation behind.
func TestOperatorFormsMatchEngine(t *testing.T) {
	tc := startTiDB(t)
	accepts := tc.prepareSyntax
	for _, c := range operatorFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := accepts(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("TiDB v8.5.5 accepts %q = %v, this file records %v", c.sql, got, c.accept)
			}
		})
	}
}

// TestOperatorTreesMatchEngine evaluates each pinned shape on the server against
// the text as written, over every assignment of a, b and c from 0, 1, 2 and NULL. A
// shape that groups differently from the engine shows as a row where the two differ.
func TestOperatorTreesMatchEngine(t *testing.T) {
	tc := startTiDB(t)
	db, ctx := tc.db, tc.ctx
	values := func(name string) string {
		return "(SELECT 0 AS " + name + " UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT NULL)"
	}
	from := values("a") + " ta CROSS JOIN " + values("b") + " tb CROSS JOIN " + values("c") + " tc"
	for _, c := range operatorTreeCases {
		t.Run(c.expr, func(t *testing.T) {
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE NOT ((%s) <=> (%s))", from, c.expr, c.shape)
			var differ int
			if err := db.QueryRowContext(ctx, query).Scan(&differ); err != nil {
				t.Fatalf("%s: %v", query, err)
			}
			if differ != 0 {
				t.Errorf("%d rows differ between %s and %s", differ, c.expr, c.shape)
			}
		})
	}
}

// sweepExpr writes a random expression over small constants, with no parentheses
// beyond the few it picks, so the grouping is left entirely to the grammar.
const (
	operatorSweepN          = 1000
	operatorSweepSeed int64 = 1
)

func sweepExpr(r *rand.Rand, depth int) string {
	atoms := []string{"0", "1", "2", "3", "NULL"}
	if depth <= 0 {
		return atoms[r.Intn(len(atoms))]
	}
	binary := []string{
		"OR", "||", "XOR", "AND", "&&", "=", "<>", "!=", "<", ">", "<=", ">=", "<=>",
		"|", "&", "<<", ">>", "+", "-", "*", "/", "DIV", "%", "MOD", "^",
	}
	prefix := []string{"NOT ", "! ", "~ ", "- ", "+ "}
	postfix := []string{" IS NULL", " IS NOT NULL", " IS TRUE", " IS NOT FALSE", " IS UNKNOWN"}
	maybeNot := func() string {
		if r.Intn(2) == 0 {
			return "NOT "
		}
		return ""
	}
	switch x := r.Intn(100); {
	case x < 50:
		return sweepExpr(r, depth-1) + " " + binary[r.Intn(len(binary))] + " " + sweepExpr(r, depth-1)
	case x < 65:
		return prefix[r.Intn(len(prefix))] + sweepExpr(r, depth-1)
	case x < 75:
		return sweepExpr(r, depth-1) + postfix[r.Intn(len(postfix))]
	case x < 83:
		return sweepExpr(r, depth-1) + " " + maybeNot() + "BETWEEN " + sweepExpr(r, depth-1) + " AND " + sweepExpr(r, depth-1)
	case x < 90:
		return sweepExpr(r, depth-1) + " " + maybeNot() + "IN (" + sweepExpr(r, depth-1) + ", " + sweepExpr(r, depth-1) + ")"
	case x < 93:
		return sweepExpr(r, depth-1) + " " + maybeNot() + "LIKE " + sweepExpr(r, depth-1)
	case x < 95:
		return sweepExpr(r, depth-1) + " " + maybeNot() + "REGEXP " + sweepExpr(r, depth-1)
	case x < 96:
		return "@v := " + sweepExpr(r, depth-1)
	case x < 98:
		return "(" + sweepExpr(r, depth-1) + ")"
	}
	return atoms[r.Intn(len(atoms))]
}

// TestOperatorSweepMatchesEngine asks the engine two questions about random
// expressions. Does it parse them, compared with Parse? And for those it does, does
// `(expr) <=> (exprShape(tree))` hold? The second is what a verdict corpus cannot do:
// the engine evaluates the text as written and the grouping this parser gave it, so a
// wrong tree shows as a wrong value with no evaluator or expected tree written here.
//
// operatorSweepN and operatorSweepSeed size and seed it. Its tape holds the stream
// they make, so changing either needs a recording. The second question is built
// from omni's tree, so a change to the tree of a swept expression needs one too.
func TestOperatorSweepMatchesEngine(t *testing.T) {
	tc := startTiDB(t)
	db := tc.db
	ctx := tc.ctx
	if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS omni_sweep"); err != nil {
		t.Fatal(err)
	}
	n, seed := operatorSweepN, operatorSweepSeed
	r := rand.New(rand.NewSource(seed))
	var compared int
	for i := 0; i < n; i++ {
		expr := sweepExpr(r, 1+r.Intn(3))
		stmt := "SELECT " + expr
		engine, err := operatorVerdict(ctx, db, stmt)
		if err != nil {
			t.Fatalf("seed %d: the engine gave no verdict for %q: %v", seed, stmt, err)
		}
		_, parseErr := Parse(stmt)
		if (parseErr == nil) != engine {
			t.Errorf("seed %d: Parse(%q) err = %v, engine accepts = %v", seed, stmt, parseErr, engine)
			continue
		}
		if parseErr != nil {
			continue
		}
		shape := exprShape(t, selectExpr(t, expr))
		var same sql.NullInt64
		compare := fmt.Sprintf("SELECT (%s) <=> (%s)", expr, shape)
		if err := db.QueryRowContext(ctx, compare).Scan(&same); err != nil {
			var myErr *gomysql.MySQLError
			if errors.As(err, &myErr) {
				continue // a value out of range, say: the constants can overflow
			}
			t.Fatalf("seed %d: %s: %v", seed, compare, err)
		}
		compared++
		if !same.Valid || same.Int64 != 1 {
			t.Errorf("seed %d: the engine does not read %q as %q", seed, expr, shape)
		}
	}
	if compared < n/2 {
		t.Errorf("only %d of %d expressions were compared: the sweep has stopped saying anything", compared, n)
	}
}
