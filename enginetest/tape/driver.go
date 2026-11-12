package tape

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"testing"
)

// errPrepared is what a test gets for a driver-level prepared statement. A
// tape holds each request with its arguments; the live connector inlines
// them (interpolateParams in go-sql-driver), so a test never needs one.
var errPrepared = errors.New("tape: prepared statements are not recorded; send the statement with its arguments, which the live connector inlines")

// A session is one stream's state while its test runs: the answers it replays,
// or the exchanges it records.
type session struct {
	t     testing.TB
	name  string
	codec Codec
	rec   bool

	// poison ends the recording of this stream's file. Recording only.
	poison func() error

	mu     sync.Mutex
	want   []exchange // replay: the tape; recording: the committed stream, if any
	next   int
	broken error
	got    []exchange // recording
}

// replay answers req from the tape, or fails the test when the tape holds
// something else next.
func (s *session) replay(req request) (answer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.broken != nil {
		return answer{}, s.broken
	}
	if s.next >= len(s.want) || s.want[s.next].req.render() != req.render() {
		msg := fmt.Sprintf("no recorded answer for %s (request %d of stream %q)", req.describe(), s.next+1, s.name)
		if s.next < len(s.want) {
			msg += fmt.Sprintf("; the tape holds %s there", s.want[s.next].req.describe())
		}
		s.broken = fmt.Errorf("%s. Record it:\n  %s", msg, recordCommand(s.t))
		s.t.Errorf("%v", s.broken)
		return answer{}, s.broken
	}
	a := s.want[s.next].ans
	s.next++
	return a, nil
}

func (s *session) record(req request, a answer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.got = append(s.got, exchange{req, a})
}

// finish fails a replay that left answers no request asked for.
func (s *session) finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.broken == nil && s.next < len(s.want) {
		s.t.Errorf("stream %q holds %d recorded requests the test never sent, starting with %s. Record it:\n  %s",
			s.name, len(s.want)-s.next, s.want[s.next].req.describe(), recordCommand(s.t))
	}
}

// reproduces fails a recording whose answers differ from the committed
// tape's for the same requests: such an answer can't be replayed.
func (s *session) reproduces() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.got) && i < len(s.want); i++ {
		if s.got[i].req.render() != s.want[i].req.render() {
			return
		}
		if was, now := s.want[i].ans.render(), s.got[i].ans.render(); was != now {
			s.t.Errorf("stream %q: the engine answered request %d (%s) differently than the committed tape. Make the test's requests deterministic, or pin what varies.\nwas:\n%snow:\n%s",
				s.name, i+1, s.got[i].req.describe(), was, now)
			return
		}
	}
}

// errorAnswer turns an error the live server returned into an answer. An
// error the engine did not answer (a timeout, a driver fault) fails the
// test and poisons the recording: it must never reach the tape, where a
// replay would hand it to the test as the engine's verdict. The test gets
// that error as the driver returned it.
func (s *session) errorAnswer(err error) answer {
	if errors.Is(err, driver.ErrBadConn) {
		return answer{kind: "badconn"}
	}
	text, ok := s.codec.EncodeError(err)
	if !ok {
		s.t.Errorf("stream %q: %v is not an answer from the engine, so it can't be recorded", s.name, err)
		if perr := s.poison(); perr != nil {
			s.t.Errorf("tape: %v", perr)
		}
		return answer{kind: "error", err: err.Error(), raw: err}
	}
	return answer{kind: "error", err: text}
}

// result turns an answer back into what the driver returns, the same way
// in a recording as in a replay.
func (s *session) result(a answer) (answer, error) {
	switch a.kind {
	case "error":
		if a.raw != nil {
			return a, a.raw
		}
		return a, s.codec.DecodeError(a.err)
	case "badconn":
		return a, driver.ErrBadConn
	}
	return a, nil
}

type connector struct {
	s    *session
	live driver.Connector // nil in a replay
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	if c.live == nil {
		return &conn{s: c.s}, nil
	}
	lc, err := c.live.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &conn{s: c.s, live: lc}, nil
}

func (c *connector) Driver() driver.Driver { return tapeDriver{} }

type tapeDriver struct{}

func (tapeDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("tape: open a connection with tape.Open")
}

type conn struct {
	s    *session
	live driver.Conn // nil in a replay
}

var (
	_ driver.ExecerContext      = (*conn)(nil)
	_ driver.QueryerContext     = (*conn)(nil)
	_ driver.ConnBeginTx        = (*conn)(nil)
	_ driver.NamedValueChecker  = (*conn)(nil)
	_ driver.Pinger             = (*conn)(nil)
	_ driver.RowsNextResultSet  = (*rows)(nil)
	_ driver.SessionResetter    = (*conn)(nil)
	_ driver.ConnPrepareContext = (*conn)(nil)
)

func (c *conn) Prepare(string) (driver.Stmt, error) { return nil, errPrepared }

func (c *conn) PrepareContext(context.Context, string) (driver.Stmt, error) {
	return nil, errPrepared
}

func (c *conn) Close() error {
	if c.live != nil {
		return c.live.Close()
	}
	return nil
}

func (c *conn) CheckNamedValue(nv *driver.NamedValue) error {
	v, err := driver.DefaultParameterConverter.ConvertValue(nv.Value)
	nv.Value = v
	return err
}

func (c *conn) Ping(ctx context.Context) error {
	if p, ok := c.live.(driver.Pinger); ok {
		return p.Ping(ctx)
	}
	return nil
}

func (c *conn) ResetSession(ctx context.Context) error {
	if r, ok := c.live.(driver.SessionResetter); ok {
		return r.ResetSession(ctx)
	}
	return nil
}

func newRequest(kind, text string, args []driver.NamedValue) (request, error) {
	r := request{kind: kind, text: text}
	for _, a := range args {
		if a.Name != "" {
			return r, fmt.Errorf("tape: named argument %q; use positional arguments", a.Name)
		}
		v, err := fromDriver(a.Value)
		if err != nil {
			return r, fmt.Errorf("tape: argument %d: %w", a.Ordinal, err)
		}
		r.args = append(r.args, v)
	}
	return r, nil
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	req, err := newRequest("exec", query, args)
	if err != nil {
		return nil, err
	}
	var a answer
	if c.live == nil {
		if a, err = c.s.replay(req); err != nil {
			return nil, err
		}
	} else {
		ex, ok := c.live.(driver.ExecerContext)
		if !ok {
			return nil, errors.New("tape: the live driver has no ExecContext")
		}
		res, err := ex.ExecContext(ctx, query, args)
		switch {
		case errors.Is(err, driver.ErrSkip):
			return nil, errPrepared
		case err != nil:
			a = c.s.errorAnswer(err)
		default:
			a = answer{kind: "ok", affected: "-", lastID: "-"}
			if n, err := res.RowsAffected(); err == nil {
				a.affected = strconv.FormatInt(n, 10)
			}
			if n, err := res.LastInsertId(); err == nil {
				a.lastID = strconv.FormatInt(n, 10)
			}
		}
		c.s.record(req, a)
	}
	if a, err = c.s.result(a); err != nil {
		return nil, err
	}
	if a.kind != "ok" {
		return nil, fmt.Errorf("tape: stream %q answers exec with %q", c.s.name, a.kind)
	}
	return result{a}, nil
}

type result struct{ a answer }

func (r result) RowsAffected() (int64, error) { return parseCount(r.a.affected) }
func (r result) LastInsertId() (int64, error) { return parseCount(r.a.lastID) }

func parseCount(s string) (int64, error) {
	if s == "-" {
		return 0, errors.New("tape: the driver gave no count for this statement")
	}
	return strconv.ParseInt(s, 10, 64)
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	req, err := newRequest("query", query, args)
	if err != nil {
		return nil, err
	}
	var a answer
	if c.live == nil {
		if a, err = c.s.replay(req); err != nil {
			return nil, err
		}
	} else {
		q, ok := c.live.(driver.QueryerContext)
		if !ok {
			return nil, errors.New("tape: the live driver has no QueryContext")
		}
		lr, err := q.QueryContext(ctx, query, args)
		switch {
		case errors.Is(err, driver.ErrSkip):
			return nil, errPrepared
		case err != nil:
			a = c.s.errorAnswer(err)
		default:
			a, err = readAll(lr)
			if err != nil {
				a = c.s.errorAnswer(err)
			}
		}
		c.s.record(req, a)
	}
	if a, err = c.s.result(a); err != nil {
		return nil, err
	}
	if a.kind != "rows" {
		return nil, fmt.Errorf("tape: stream %q answers query with %q", c.s.name, a.kind)
	}
	return &rows{sets: a.sets}, nil
}

// readAll reads every result set of a live answer, so the test then reads
// the recorded copy exactly as a replay would.
func readAll(lr driver.Rows) (answer, error) {
	defer lr.Close()
	a := answer{kind: "rows"}
	for {
		var s resultSet
		typed, _ := lr.(driver.RowsColumnTypeDatabaseTypeName)
		for i, name := range lr.Columns() {
			c := column{name: name}
			if typed != nil {
				c.typ = typed.ColumnTypeDatabaseTypeName(i)
			}
			s.cols = append(s.cols, c)
		}
		dest := make([]driver.Value, len(s.cols))
		for {
			err := lr.Next(dest)
			if err == io.EOF {
				break
			}
			if err != nil {
				return a, err
			}
			row := make([]value, len(dest))
			for i, d := range dest {
				v, err := fromDriver(d) // copies []byte, which the driver reuses
				if err != nil {
					return a, fmt.Errorf("tape: column %q: %w", s.cols[i].name, err)
				}
				row[i] = v
			}
			s.rows = append(s.rows, row)
		}
		a.sets = append(a.sets, s)
		more, ok := lr.(driver.RowsNextResultSet)
		if !ok || !more.HasNextResultSet() {
			return a, nil
		}
		if err := more.NextResultSet(); err != nil {
			if err == io.EOF {
				return a, nil
			}
			return a, err
		}
	}
}

type rows struct {
	sets []resultSet
	set  int
	row  int
}

func (r *rows) Columns() []string {
	cols := r.sets[r.set].cols
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.name
	}
	return names
}

func (r *rows) ColumnTypeDatabaseTypeName(i int) string { return r.sets[r.set].cols[i].typ }

func (r *rows) Close() error { return nil }

func (r *rows) Next(dest []driver.Value) error {
	s := r.sets[r.set]
	if r.row >= len(s.rows) {
		return io.EOF
	}
	for i, v := range s.rows[r.row] {
		d, err := v.driver()
		if err != nil {
			return err
		}
		dest[i] = d
	}
	r.row++
	return nil
}

func (r *rows) HasNextResultSet() bool { return r.set+1 < len(r.sets) }

func (r *rows) NextResultSet() error {
	if !r.HasNextResultSet() {
		return io.EOF
	}
	r.set++
	r.row = 0
	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	var ltx driver.Tx
	err := c.step(request{kind: "begin"}, func() error {
		var err error
		if b, ok := c.live.(driver.ConnBeginTx); ok {
			ltx, err = b.BeginTx(ctx, opts)
		} else {
			ltx, err = c.live.Begin() //nolint:staticcheck // the drivers that lack BeginTx
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &tx{c: c, live: ltx}, nil
}

type tx struct {
	c    *conn
	live driver.Tx
}

func (t *tx) Commit() error {
	return t.c.step(request{kind: "commit"}, func() error { return t.live.Commit() })
}

func (t *tx) Rollback() error {
	return t.c.step(request{kind: "rollback"}, func() error { return t.live.Rollback() })
}

// step runs a request whose answer is only ok or an error.
func (c *conn) step(req request, run func() error) error {
	var a answer
	if c.live == nil {
		var err error
		if a, err = c.s.replay(req); err != nil {
			return err
		}
	} else {
		a = answer{kind: "ok"}
		if err := run(); err != nil {
			a = c.s.errorAnswer(err)
		}
		c.s.record(req, a)
	}
	_, err := c.s.result(a)
	return err
}
