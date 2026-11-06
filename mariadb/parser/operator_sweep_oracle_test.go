package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
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
		switch myErr.Number {
		case 1064:
			return false, nil
		case 1115, 1128, 1193, 1253, 1273, 1286, 1289, 4081:
			return true, nil
		}
		// A code the parser's actions raise comes back with a syntax error after the
		// statement; one raised after parsing turns into that syntax error.
		if _, err := conn.ExecContext(ctx, "SET @omni_sweep = ?", stmt+"\n; SELECT FROM"); err != nil {
			return false, err
		}
		_, err := conn.ExecContext(ctx, "PREPARE omni_sweep_stmt FROM @omni_sweep")
		var again *gomysql.MySQLError
		if !errors.As(err, &again) {
			return false, err
		}
		return again.Number != myErr.Number, nil
	}
	return false, err
}

// TestOperatorFormsMatchEngine holds operator_forms_test.go's verdicts to the
// server, so an image upgrade that changes an answer fails here rather than leaving
// a stale expectation behind.
func TestOperatorFormsMatchEngine(t *testing.T) {
	o := startMariaDB(t)
	accepts := func(sql string) (bool, error) { return operatorVerdict(o.ctx, o.db, sql) }
	for _, c := range operatorFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := accepts(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("MariaDB 11.8.8 accepts %q = %v, this file records %v", c.sql, got, c.accept)
			}
		})
	}
}

// TestOperatorTreesMatchEngine evaluates each pinned shape on the server against
// the text as written, over every assignment of a, b and c from 0, 1, 2 and NULL. A
// shape that groups differently from the engine shows as a row where the two differ.
func TestOperatorTreesMatchEngine(t *testing.T) {
	o := startMariaDB(t)
	db, ctx := o.db, o.ctx
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
// OPERATOR_SWEEP_N and OPERATOR_SWEEP_SEED size and seed it; a failure names the seed.
func TestOperatorSweepMatchesEngine(t *testing.T) {
	o := startMariaDB(t)
	db := o.db
	ctx := o.ctx
	n, seed := 1000, int64(1)
	if v, err := strconv.Atoi(os.Getenv("OPERATOR_SWEEP_N")); err == nil {
		n = v
	}
	if v, err := strconv.ParseInt(os.Getenv("OPERATOR_SWEEP_SEED"), 10, 64); err == nil {
		seed = v
	}
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
