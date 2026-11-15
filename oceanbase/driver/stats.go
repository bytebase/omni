package driver

import (
	"context"
	"database/sql"
	"strings"

	oceanbaseanalysis "github.com/bytebase/omni/oceanbase/analysis"
)

// CountAffectedRows returns OceanBase's planner estimate of the rows statement
// modifies from EXPLAIN FORMAT=JSON.
func (d *Driver) CountAffectedRows(ctx context.Context, statement string) (int64, error) {
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
