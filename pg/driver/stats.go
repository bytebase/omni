package driver

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	pganalysis "github.com/bytebase/omni/pg/analysis"
)

// CountAffectedRows returns the planner's estimate of the rows statement
// modifies. Any search path prefix from analysis.WithSearchPath is applied in
// a rolled-back transaction before running EXPLAIN.
func (d *Driver) CountAffectedRows(ctx context.Context, statement string) (int64, error) {
	setup, statement := pganalysis.SplitSearchPath(statement)
	queryRow := d.db.QueryRowContext
	if setup != "" {
		tx, err := d.db.BeginTx(ctx, nil)
		if err != nil {
			return 0, err
		}
		defer func() { _ = tx.Rollback() }()
		if _, err := tx.ExecContext(ctx, setup); err != nil {
			return 0, err
		}
		queryRow = tx.QueryRowContext
	}
	var plan string
	if err := queryRow(ctx, "EXPLAIN (FORMAT JSON) "+statement).Scan(&plan); err != nil {
		return 0, err
	}
	return pganalysis.EstimatedAffectedRowsFromExplainJSON(plan)
}

// Explain dry-runs statement with EXPLAIN inside a rolled-back transaction
// after replaying settings.
func (d *Driver) Explain(ctx context.Context, settings []string, statement string) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if err := d.prepareExplainTx(ctx, tx, settings); err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, "EXPLAIN "+statement)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
	}
	return rows.Err()
}

// ExplainJSON runs EXPLAIN (FORMAT JSON) on statement inside a rolled-back
// transaction after replaying settings and returns the JSON plan text.
func (d *Driver) ExplainJSON(ctx context.Context, settings []string, statement string) (string, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()
	if err := d.prepareExplainTx(ctx, tx, settings); err != nil {
		return "", err
	}
	rows, err := tx.QueryContext(ctx, "EXPLAIN (FORMAT JSON) "+statement)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var plans []string
	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			return "", err
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if len(plans) != 1 {
		return "", fmt.Errorf("expected one plan row but got %d", len(plans))
	}
	return plans[0], nil
}

func (d *Driver) prepareExplainTx(ctx context.Context, tx *sql.Tx, settings []string) error {
	if d.tenantMode {
		const ownerQuery = `
		SELECT
			u.rolname
		FROM
			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
		WHERE
			d.datname = current_database();
		`
		var owner string
		if err := tx.QueryRowContext(ctx, ownerQuery).Scan(&owner); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("SET SESSION ROLE %s;", quoteLiteral(owner))); err != nil {
			return err
		}
	}
	for _, preExec := range settings {
		if preExec != "" {
			if _, err := tx.ExecContext(ctx, preExec); err != nil {
				return fmt.Errorf("failed to execute pre-execution: %s: %w", preExec, err)
			}
		}
	}
	return nil
}

// Activity is a session matching an application_name in pg_stat_activity
// together with the sessions it blocks and the sessions blocking it.
type Activity struct {
	Session          *BackendSession
	BlockedSessions  []*BackendSession
	BlockingSessions []*BackendSession
}

// BackendSession is one row read from pg_catalog.pg_stat_activity.
type BackendSession struct {
	PID             string
	BlockedByPIDs   []string
	Query           string
	State           *string
	WaitEventType   *string
	WaitEvent       *string
	Database        *string
	User            *string
	ApplicationName string
	ClientAddr      *string
	ClientPort      *string
	BackendStart    time.Time
	XactStart       *time.Time
	QueryStart      *time.Time
}

// Activity reads the pg_stat_activity session for applicationName and the
// sessions blocking or blocked by it.
func (d *Driver) Activity(ctx context.Context, applicationName string) (*Activity, error) {
	const query = `
			WITH target_session AS (
				SELECT pid FROM pg_catalog.pg_stat_activity WHERE application_name = $1 LIMIT 1
			)
			SELECT
				a.pid,
				pg_blocking_pids(a.pid) AS blocked_by_pids,
				a.query,
				a.state,
				a.wait_event_type,
				a.wait_event,
				a.datname,
				a.usename,
				a.application_name,
				a.client_addr,
				a.client_port,
				a.backend_start,
				a.xact_start,
				a.query_start
			FROM
				pg_catalog.pg_stat_activity a
			WHERE a.application_name = $1
			   OR (SELECT pid FROM target_session) = ANY(pg_blocking_pids(a.pid))
			   OR a.pid = ANY(pg_blocking_pids((SELECT pid FROM target_session)))
			ORDER BY a.pid
		`
	rows, err := d.db.QueryContext(ctx, query, applicationName)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}
	defer rows.Close()

	types := pgtype.NewMap()
	out := &Activity{}
	for rows.Next() {
		var s BackendSession
		var blockedBy []int32
		var backendStart sql.NullTime
		if err := rows.Scan(
			&s.PID,
			types.SQLScanner(&blockedBy),
			&s.Query,
			&s.State,
			&s.WaitEventType,
			&s.WaitEvent,
			&s.Database,
			&s.User,
			&s.ApplicationName,
			&s.ClientAddr,
			&s.ClientPort,
			&backendStart,
			&s.XactStart,
			&s.QueryStart,
		); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		if backendStart.Valid {
			s.BackendStart = backendStart.Time
		}
		if blockedBy != nil {
			s.BlockedByPIDs = make([]string, len(blockedBy))
			for i, pid := range blockedBy {
				s.BlockedByPIDs[i] = strconv.FormatInt(int64(pid), 10)
			}
		}

		if s.ApplicationName == applicationName {
			out.Session = &s
		} else if out.Session != nil {
			if slices.Contains(s.BlockedByPIDs, out.Session.PID) {
				out.BlockedSessions = append(out.BlockedSessions, &s)
			} else if slices.Contains(out.Session.BlockedByPIDs, s.PID) {
				out.BlockingSessions = append(out.BlockingSessions, &s)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}
	return out, nil
}
