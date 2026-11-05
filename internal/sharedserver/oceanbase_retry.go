package sharedserver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
)

// OceanBaseTransient reports an error OceanBase answers while it is between
// leaders or starved for CPU, which the same statement clears once the server
// settles. On a busy box (load average 380 on 20 cores, 2026-11-05) the root
// service is revoked and taken back within a minute, and a log stream changes
// leader.
//
// A verdict must never read one of these as an answer: engineAccepts-style
// helpers take any code but 1064 to mean the statement parsed.
func OceanBaseTransient(err error) bool {
	var me *gomysql.MySQLError
	if errors.As(err, &me) {
		switch me.Number {
		// 4012 is a timeout, 4550 "rootservice is shutdown", 4638 "The
		// RootServer is not the master", 4767 the tenant's log stream between
		// leaders; the rest are the walk's list from 2026-10-28.
		case 4012, 4023, 4038, 4109, 4550, 4638, 4721, 4767, 6005:
			return true
		}
	}
	return false
}

// retryWindow is how long a statement keeps being retried, retryPause the wait
// between tries. Tests shorten them.
var (
	retryWindow = 3 * time.Minute
	retryPause  = 2 * time.Second
)

// OpenOceanBase opens dsn through a connector that runs each statement again
// while OceanBase answers a transient error, for up to three minutes, so no test
// has to wrap its own statements. A statement still failing after that returns
// the last error as a plain error.
func OpenOceanBase(dsn string) (*sql.DB, error) {
	cfg, err := gomysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	connector, err := gomysql.NewConnector(cfg)
	if err != nil {
		return nil, err
	}
	return sql.OpenDB(retryConnector{connector}), nil
}

type retryConnector struct{ driver.Connector }

func (c retryConnector) Connect(ctx context.Context) (driver.Conn, error) {
	var conn driver.Conn
	err := retry(ctx, func() error {
		var err error
		conn, err = c.Connector.Connect(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &retryConn{conn}, nil
}

// retry runs f until it returns something other than a transient error. A
// transient error that outlasts the window comes back as a plain error, not a
// *mysql.MySQLError: a verdict helper reading codes then finds no server answer
// and fails, where 4550 would have read as "parsed".
func retry(ctx context.Context, f func() error) error {
	deadline := time.Now().Add(retryWindow)
	for {
		err := f()
		if !OceanBaseTransient(err) {
			return err
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("OceanBase stayed unavailable for %v: %v", retryWindow, err)
		}
		select {
		case <-ctx.Done():
			return err
		case <-time.After(retryPause):
		}
	}
}

// retryConn retries Exec and Query and passes the rest of the driver's
// interfaces through, so database/sql treats it as the mysql connection.
type retryConn struct{ driver.Conn }

func (c *retryConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	var res driver.Result
	err := retry(ctx, func() error {
		var err error
		res, err = c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
		return err
	})
	return res, err
}

func (c *retryConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	var rows driver.Rows
	err := retry(ctx, func() error {
		var err error
		rows, err = c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
		return err
	})
	return rows, err
}

func (c *retryConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return c.Conn.(driver.ConnPrepareContext).PrepareContext(ctx, query)
}

func (c *retryConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
}

func (c *retryConn) Ping(ctx context.Context) error {
	return retry(ctx, func() error { return c.Conn.(driver.Pinger).Ping(ctx) })
}

func (c *retryConn) ResetSession(ctx context.Context) error {
	return c.Conn.(driver.SessionResetter).ResetSession(ctx)
}

func (c *retryConn) IsValid() bool { return c.Conn.(driver.Validator).IsValid() }

func (c *retryConn) CheckNamedValue(nv *driver.NamedValue) error {
	return c.Conn.(driver.NamedValueChecker).CheckNamedValue(nv)
}
