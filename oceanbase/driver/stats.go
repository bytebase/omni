package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	oceanbaseanalysis "github.com/bytebase/omni/oceanbase/analysis"
	oceanbaseast "github.com/bytebase/omni/oceanbase/ast"
	oceanbaseparser "github.com/bytebase/omni/oceanbase/parser"
)

// CountAffectedRows returns OceanBase's planner estimate of the rows statement
// modifies from EXPLAIN FORMAT=JSON.
func (d *Driver) CountAffectedRows(ctx context.Context, statement string) (int64, error) {
	if _, err := dmlStatement(statement); err != nil {
		return 0, err
	}
	rows, err := d.db.QueryContext(ctx, "EXPLAIN FORMAT=JSON "+statement)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	plan, err := readExplainJSON(rows)
	if err != nil {
		return 0, err
	}
	return oceanbaseanalysis.EstimateAffectedRowsFromExplainJSON(plan)
}

// dmlStatement returns statement's one INSERT, UPDATE or DELETE. Only a
// statement that changes rows has an estimate, so anything else is refused
// before it reaches the server.
func dmlStatement(statement string) (oceanbaseast.Node, error) {
	list, err := oceanbaseparser.Parse(statement)
	if err != nil {
		return nil, fmt.Errorf("failed to parse statement: %w", err)
	}
	if list == nil || len(list.Items) != 1 {
		return nil, errors.New("expected exactly one statement")
	}
	switch list.Items[0].(type) {
	case *oceanbaseast.InsertStmt, *oceanbaseast.UpdateStmt, *oceanbaseast.DeleteStmt:
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

// ExplainJSON runs EXPLAIN FORMAT=JSON on statement inside a rolled-back
// transaction and returns the concatenated JSON plan text.
func (d *Driver) ExplainJSON(ctx context.Context, statement string) (string, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, "EXPLAIN FORMAT=JSON "+statement)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	return readExplainJSON(rows)
}

// ColumnExtra reads INFORMATION_SCHEMA.COLUMNS.EXTRA for one column, used by
// OceanBase's offline-DDL review rule.
func (d *Driver) ColumnExtra(ctx context.Context, database, table, column string) (string, error) {
	var extra string
	err := d.db.QueryRowContext(ctx,
		"SELECT EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ?",
		database, table, column,
	).Scan(&extra)
	if err != nil {
		return "", err
	}
	return extra, nil
}

func readExplainJSON(rows *sql.Rows) (string, error) {
	var plan strings.Builder
	for rows.Next() {
		var line sql.NullString
		if err := rows.Scan(&line); err != nil {
			return "", err
		}
		plan.WriteString(line.String)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return plan.String(), nil
}
