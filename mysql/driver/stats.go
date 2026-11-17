package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	gomysql "github.com/go-sql-driver/mysql"

	mysqlanalysis "github.com/bytebase/omni/mysql/analysis"
	"github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

const erUnknownSystemVariable = 1193

// CountAffectedRows returns the planner's estimate of the rows statement
// modifies.
func (d *Driver) CountAffectedRows(ctx context.Context, statement string) (int64, error) {
	list, err := mysqlparser.Parse(statement)
	if err != nil {
		return 0, fmt.Errorf("failed to parse statement: %w", err)
	}
	if list == nil || len(list.Items) != 1 {
		return 0, errors.New("expected exactly one statement")
	}
	// Only a statement that changes rows has an estimate. The tabular fallback
	// below stands in for a plan that carries none, so without this check a
	// SELECT of a table would report the rows it reads as rows it modifies.
	switch list.Items[0].(type) {
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
	default:
		return 0, fmt.Errorf("expected an INSERT, UPDATE or DELETE statement, got %T", list.Items[0])
	}
	query := mysqlanalysis.AffectedRowsQuery(list.Items[0], statement)
	plan, err := d.ExplainJSON(ctx, query)
	if err != nil {
		return 0, err
	}
	return d.EstimateAffectedRows(ctx, list.Items[0], query, plan)
}

// EstimateAffectedRows returns the rows stmt modifies according to plan,
// falling back to tabular EXPLAIN when the JSON plan carries no estimate.
func (d *Driver) EstimateAffectedRows(ctx context.Context, stmt ast.Node, query, plan string) (int64, error) {
	count, err := mysqlanalysis.EstimateAffectedRowsFromExplainJSON(stmt, plan)
	if err == nil {
		return count, nil
	}
	rows, tables, ok := d.ExplainTabularEstimate(ctx, query)
	if !ok {
		return 0, err
	}
	return mysqlanalysis.CapAffectedRowsByLimit(stmt, rows*float64(mysqlanalysis.DMLTargetCount(stmt, tables))), nil
}

// ExplainTabularEstimate reads the row product and joined-table count of the
// first query block from EXPLAIN FORMAT=TRADITIONAL.
func (d *Driver) ExplainTabularEstimate(ctx context.Context, query string) (float64, int, bool) {
	// The connection sends every statement in one request (see connConfig),
	// so the probe plans exactly the one statement it was given.
	if err := requireOneStatement(query); err != nil {
		return 0, 0, false
	}
	rows, err := d.db.QueryContext(ctx, "EXPLAIN FORMAT=TRADITIONAL "+query)
	if err != nil {
		return 0, 0, false
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return 0, 0, false
	}
	column := func(name string) int {
		return slices.IndexFunc(columns, func(c string) bool { return strings.EqualFold(c, name) })
	}
	idIndex, rowsIndex, filteredIndex := column("id"), column("rows"), column("filtered")
	if idIndex < 0 || rowsIndex < 0 {
		return 0, 0, false
	}
	values := make([]sql.NullString, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	var blockID string
	estimate, tables := 1.0, 0
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return 0, 0, false
		}
		tableRows, err := strconv.ParseFloat(values[rowsIndex].String, 64)
		if !values[rowsIndex].Valid || err != nil || (tables > 0 && values[idIndex].String != blockID) {
			continue
		}
		blockID = values[idIndex].String
		tables++
		estimate *= tableRows
		if filteredIndex >= 0 && values[filteredIndex].Valid {
			if filtered, err := strconv.ParseFloat(values[filteredIndex].String, 64); err == nil {
				estimate *= filtered / 100
			}
		}
	}
	return estimate, tables, tables > 0 && rows.Err() == nil
}

// ExplainJSON returns the EXPLAIN FORMAT=JSON output for statement in the
// format version mysqlanalysis.EstimateAffectedRowsFromExplainJSON reads.
func (d *Driver) ExplainJSON(ctx context.Context, statement string) (string, error) {
	if err := requireOneStatement(statement); err != nil {
		return "", err
	}
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// The session variable must be set on the connection that runs the
	// EXPLAIN. MySQL before 8.3 lacks it and only produces version 1.
	if _, err := conn.ExecContext(ctx, "SET SESSION explain_json_format_version = 1"); err != nil {
		var mysqlErr *gomysql.MySQLError
		if !errors.As(err, &mysqlErr) || mysqlErr.Number != erUnknownSystemVariable {
			return "", fmt.Errorf("failed to set explain_json_format_version: %w", err)
		}
	} else {
		// The connection returns to the pool, so later statements on it get
		// the server's version back.
		defer func() {
			_, _ = conn.ExecContext(ctx, "SET SESSION explain_json_format_version = DEFAULT")
		}()
	}

	rows, err := conn.QueryContext(ctx, "EXPLAIN FORMAT=JSON "+statement)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	return readExplainJSON(rows)
}

// Explain dry-runs statement with EXPLAIN inside a rolled-back transaction.
func (d *Driver) Explain(ctx context.Context, statement string) error {
	if err := requireOneStatement(statement); err != nil {
		return err
	}
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

// requireOneStatement refuses text that holds anything but one statement: these
// probes prefix EXPLAIN to what they are given, and the connection sends every
// statement of a request, so a second one would run.
func requireOneStatement(statement string) error {
	list, err := mysqlparser.Parse(statement)
	if err != nil {
		return fmt.Errorf("failed to parse statement: %w", err)
	}
	if list == nil || len(list.Items) != 1 {
		return errors.New("expected exactly one statement")
	}
	return nil
}
