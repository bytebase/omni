package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	tidbast "github.com/bytebase/omni/tidb/ast"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// CountAffectedRows returns TiDB's planner estimate of the rows statement
// modifies.
func (d *Driver) CountAffectedRows(ctx context.Context, statement string) (int64, error) {
	stmt, err := dmlStatement(statement)
	if err != nil {
		return 0, err
	}
	rows, err := d.db.QueryContext(ctx, "EXPLAIN "+statement)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	idIndex := slices.Index(columns, "id")
	if idIndex < 0 {
		return 0, fmt.Errorf("EXPLAIN returned no id column: %v", columns)
	}
	estRowsIndex := slices.Index(columns, "estRows")
	if estRowsIndex < 0 {
		return 0, fmt.Errorf("EXPLAIN returned no estRows column: %v", columns)
	}

	var plan []planRow
	for rows.Next() {
		scanArgs := make([]any, len(columns))
		for i := range scanArgs {
			scanArgs[i] = new(sql.NullString)
		}
		var row planRow
		scanArgs[idIndex] = &row.id
		scanArgs[estRowsIndex] = &row.estRows
		if err := rows.Scan(scanArgs...); err != nil {
			return 0, err
		}
		plan = append(plan, row)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return getAffectedRows(plan, stmt)
}

// dmlStatement returns statement's one INSERT, UPDATE or DELETE. Only a
// statement that changes rows has an estimate, so anything else is refused
// before it reaches the server.
func dmlStatement(statement string) (tidbast.Node, error) {
	list, err := tidbparser.Parse(statement)
	if err != nil {
		return nil, fmt.Errorf("failed to parse statement: %w", err)
	}
	if list == nil || len(list.Items) != 1 {
		return nil, errors.New("expected exactly one statement")
	}
	switch list.Items[0].(type) {
	case *tidbast.InsertStmt, *tidbast.UpdateStmt, *tidbast.DeleteStmt:
		return list.Items[0], nil
	default:
		return nil, fmt.Errorf("expected an INSERT, UPDATE or DELETE statement, got %T", list.Items[0])
	}
}

// Explain dry-runs statement with EXPLAIN inside a rolled-back transaction.
func (d *Driver) Explain(ctx context.Context, statement string) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, "EXPLAIN "+statement)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
	}
	return rows.Err()
}

// DryRunAutoCommit executes statement in auto-commit mode and drains its rows,
// used for TiDB's BATCH ... DRY RUN statements which cannot run in a
// transaction.
func (d *Driver) DryRunAutoCommit(ctx context.Context, statement string) error {
	rows, err := d.db.QueryContext(ctx, statement)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
	}
	return rows.Err()
}

type planRow struct {
	id      string
	estRows string
}

func getAffectedRows(plan []planRow, stmt tidbast.Node) (int64, error) {
	estimate, err := getAffectedRowsFromPlan(plan)
	if err != nil {
		return 0, err
	}
	return roundRows(estimate * float64(countDMLTargets(stmt))), nil
}

// countDMLTargets returns how many tables an UPDATE or DELETE can change per
// row read.
func countDMLTargets(stmt tidbast.Node) int {
	switch n := stmt.(type) {
	case *tidbast.DeleteStmt:
		return max(len(n.Tables), 1)
	case *tidbast.UpdateStmt:
		var tables []string
		unqualified := 0
		for _, assignment := range n.SetList {
			switch {
			case assignment.Column == nil || assignment.Column.Table == "":
				unqualified++
			default:
				table := strings.ToLower(assignment.Column.Table)
				if !slices.Contains(tables, table) {
					tables = append(tables, table)
				}
			}
		}
		return max(min(len(tables)+unqualified, countTableSources(n.Tables)), 1)
	default:
		return 1
	}
}

func countTableSources(exprs []tidbast.TableExpr) int {
	total := 0
	for _, expr := range exprs {
		total += countTableExpr(expr)
	}
	return total
}

func countTableExpr(expr tidbast.TableExpr) int {
	switch n := expr.(type) {
	case *tidbast.JoinClause:
		count := countTableExpr(n.Left)
		if n.Right != nil {
			count += countTableExpr(n.Right)
		}
		return count
	case *tidbast.TableRef, *tidbast.SubqueryExpr:
		return 1
	default:
		return 0
	}
}

func getAffectedRowsFromPlan(plan []planRow) (float64, error) {
	if len(plan) == 0 {
		return 0, errors.New("EXPLAIN returned no plan")
	}
	if operator, _, _ := strings.Cut(plan[0].id, "_"); operator != "Insert" && operator != "Update" && operator != "Delete" {
		return 0, fmt.Errorf("EXPLAIN plan root %q is not an INSERT, UPDATE, or DELETE", plan[0].id)
	}
	if len(plan) == 1 || (!strings.HasPrefix(plan[1].id, "├─") && !strings.HasPrefix(plan[1].id, "└─")) {
		return 0, nil
	}
	estRows, err := strconv.ParseFloat(plan[1].estRows, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse estRows of %q: %w", plan[1].id, err)
	}
	return estRows, nil
}

func roundRows(rows float64) int64 {
	if rows >= math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(math.Round(rows))
}
