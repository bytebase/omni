package driver

import (
	"context"
	"errors"
	"net"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/bytebase/omni/driver"
)

type wireReply struct {
	notices []pgproto3.NoticeResponse
	params  map[string]string
	fields  []pgproto3.FieldDescription
	rows    [][][]byte
	tags    []string
	err     *pgproto3.ErrorResponse
}

type wireServer struct {
	ln        net.Listener
	cancels   chan pgproto3.CancelRequest
	nextPID   atomic.Uint32
	handle    func(connID uint32, sql string, args [][]byte) wireReply
	closeOnce sync.Once
	done      chan struct{}
}

func startWireServer(t *testing.T, initTimeZone string, handle func(connID uint32, sql string, args [][]byte) wireReply) *wireServer {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ws := &wireServer{
		ln:      ln,
		cancels: make(chan pgproto3.CancelRequest, 8),
		handle:  handle,
		done:    make(chan struct{}),
	}
	t.Cleanup(func() {
		ws.closeOnce.Do(func() {
			close(ws.done)
			_ = ln.Close()
		})
	})

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go ws.serveConn(c, initTimeZone)
		}
	}()
	return ws
}

func (ws *wireServer) config() driver.Config {
	host, port, _ := net.SplitHostPort(ws.ln.Addr().String())
	return driver.Config{
		Host:   host,
		Port:   port,
		User:   "dba",
		Params: map[string]string{"sslmode": "disable"},
	}
}

func (ws *wireServer) serveConn(c net.Conn, initTimeZone string) {
	defer c.Close()
	go func() {
		<-ws.done
		_ = c.Close()
	}()

	b := pgproto3.NewBackend(c, c)
	for {
		startup, err := b.ReceiveStartupMessage()
		if err != nil {
			return
		}
		switch m := startup.(type) {
		case *pgproto3.SSLRequest:
			if _, err := c.Write([]byte{'N'}); err != nil {
				return
			}
			continue
		case *pgproto3.CancelRequest:
			select {
			case ws.cancels <- *m:
			default:
			}
			return
		case *pgproto3.StartupMessage:
			pid := ws.nextPID.Add(1)
			b.Send(&pgproto3.AuthenticationOk{})
			b.Send(&pgproto3.BackendKeyData{ProcessID: pid, SecretKey: []byte{1, 2, 3, 4}})
			if initTimeZone != "" {
				b.Send(&pgproto3.ParameterStatus{Name: "TimeZone", Value: initTimeZone})
			}
			b.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
			if err := b.Flush(); err != nil {
				return
			}
			ws.loop(b, pid)
			return
		default:
			return
		}
	}
}

func (ws *wireServer) loop(b *pgproto3.Backend, pid uint32) {
	var curSQL string
	var curArgs [][]byte
	txStatus := byte('I')

	sendReply := func(rep wireReply) {
		for i := range rep.notices {
			n := rep.notices[i]
			b.Send(&n)
		}
		for k, v := range rep.params {
			b.Send(&pgproto3.ParameterStatus{Name: k, Value: v})
		}
		if rep.err != nil {
			if txStatus == 'T' {
				txStatus = 'E'
			}
			b.Send(rep.err)
			return
		}
		if len(rep.fields) > 0 {
			b.Send(&pgproto3.RowDescription{Fields: rep.fields})
			for _, row := range rep.rows {
				b.Send(&pgproto3.DataRow{Values: row})
			}
		}
		tags := rep.tags
		if len(tags) == 0 {
			tags = []string{"OK"}
		}
		for _, tag := range tags {
			b.Send(&pgproto3.CommandComplete{CommandTag: []byte(tag)})
		}
	}

	updateTx := func(sql string, failed bool) {
		upper := strings.ToUpper(strings.TrimSpace(sql))
		switch {
		case strings.HasPrefix(upper, "BEGIN"):
			if !failed {
				txStatus = 'T'
			}
		case strings.HasPrefix(upper, "COMMIT"), strings.HasPrefix(upper, "ROLLBACK"):
			txStatus = 'I'
		}
	}

	for {
		msg, err := b.Receive()
		if err != nil {
			return
		}
		switch m := msg.(type) {
		case *pgproto3.Terminate:
			return
		case *pgproto3.Query:
			rep := ws.handle(pid, m.String, nil)
			sendReply(rep)
			updateTx(m.String, rep.err != nil)
			b.Send(&pgproto3.ReadyForQuery{TxStatus: txStatus})
			if err := b.Flush(); err != nil {
				return
			}
		case *pgproto3.Parse:
			curSQL = m.Query
			b.Send(&pgproto3.ParseComplete{})
		case *pgproto3.Describe:
			b.Send(&pgproto3.ParameterDescription{})
			b.Send(&pgproto3.NoData{})
		case *pgproto3.Bind:
			curArgs = m.Parameters
			b.Send(&pgproto3.BindComplete{})
		case *pgproto3.Execute:
			rep := ws.handle(pid, curSQL, curArgs)
			sendReply(rep)
			updateTx(curSQL, rep.err != nil)
		case *pgproto3.Sync:
			b.Send(&pgproto3.ReadyForQuery{TxStatus: txStatus})
			if err := b.Flush(); err != nil {
				return
			}
		}
	}
}

func TestOpenTenantModeFailureClosesPool(t *testing.T) {
	cfg := baseConfig()
	cfg.Host, cfg.Port = "127.0.0.1", "1"
	if _, err := Open(context.Background(), cfg, Options{TenantMode: true}); err == nil {
		t.Fatal("Open in TenantMode against closed port succeeded")
	}

	ws := startWireServer(t, "UTC", func(_ uint32, sql string, _ [][]byte) wireReply {
		if strings.Contains(sql, "pg_roles") {
			return wireReply{
				fields: []pgproto3.FieldDescription{{Name: []byte("rolname"), DataTypeOID: pgtype.TextOID, TypeModifier: -1}},
				rows:   [][][]byte{{[]byte("forbidden_owner")}},
				tags:   []string{"SELECT 1"},
			}
		}
		return wireReply{err: &pgproto3.ErrorResponse{Severity: "ERROR", Code: "42501", Message: "permission denied to set role"}}
	})
	if _, err := Open(context.Background(), ws.config(), Options{TenantMode: true}); err == nil || !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("Open when SET ROLE fails = %v, want permission denied", err)
	}
}

func TestOpenTenantModeSetsOwnerRole(t *testing.T) {
	byConn := map[uint32][]string{}
	var mu sync.Mutex
	ws := startWireServer(t, "UTC", func(connID uint32, sql string, _ [][]byte) wireReply {
		mu.Lock()
		byConn[connID] = append(byConn[connID], strings.TrimSpace(sql))
		mu.Unlock()
		if strings.Contains(sql, "pg_roles") {
			return wireReply{
				fields: []pgproto3.FieldDescription{{Name: []byte("rolname"), DataTypeOID: pgtype.TextOID, TypeModifier: -1}},
				rows:   [][][]byte{{[]byte(`app"owner`)}},
				tags:   []string{"SELECT 1"},
			}
		}
		return wireReply{tags: []string{"SET"}}
	})

	d, err := Open(context.Background(), ws.config(), Options{TenantMode: true})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	// Holding two concurrent sessions opens a second physical connection;
	// OptionAfterConnect must initialize SET ROLE on both.
	s1, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer s1.Close()
	s2, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()

	mu.Lock()
	defer mu.Unlock()
	if len(byConn) < 2 {
		t.Fatalf("opened %d physical connections, want at least 2", len(byConn))
	}
	for connID, queries := range byConn {
		if !slices.Contains(queries, `SET ROLE "app""owner";`) {
			t.Fatalf("conn %d queries = %q, missing SET ROLE", connID, queries)
		}
	}
}

func TestSessionPreValidatesWholeScriptBeforeExecuting(t *testing.T) {
	var executed []string
	var mu sync.Mutex
	ws := startWireServer(t, "UTC", func(_ uint32, sql string, _ [][]byte) wireReply {
		mu.Lock()
		executed = append(executed, sql)
		mu.Unlock()
		return wireReply{tags: []string{"INSERT 0 1"}}
	})

	d, err := Open(context.Background(), ws.config())
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	sess, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	if _, err := sess.Query(context.Background(), "INSERT INTO t VALUES (1); SELECT FROM WHERE;", driver.QueryOptions{}); err == nil {
		t.Fatal("expected parse error on second statement")
	}
	if _, err := sess.Query(context.Background(), "SELECT 1; EXPLAIN ANALYZE SELECT 1;", driver.QueryOptions{Explain: true}); err == nil {
		t.Fatal("expected explain error on second statement")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(executed) != 0 {
		t.Fatalf("server executed %q before failing script pre-validation", executed)
	}
}

func TestSessionNoticesPerConnection(t *testing.T) {
	ws := startWireServer(t, "UTC", func(connID uint32, sql string, _ [][]byte) wireReply {
		if strings.Contains(sql, "1/0") {
			return wireReply{
				notices: []pgproto3.NoticeResponse{{Severity: "NOTICE", Message: "before error"}},
				err:     &pgproto3.ErrorResponse{Severity: "ERROR", Code: "22012", Message: "division by zero"},
			}
		}
		return wireReply{
			notices: []pgproto3.NoticeResponse{
				{Severity: "NOTICE", Message: "from conn " + string(rune('0'+connID))},
				{Severity: "WARNING", Message: "be careful"},
				{Severity: "CUSTOM", Message: "unknown severity"},
			},
			fields: []pgproto3.FieldDescription{{Name: []byte("v"), DataTypeOID: pgtype.Int4OID, TypeModifier: -1}},
			rows:   [][][]byte{{[]byte("1")}},
			tags:   []string{"SELECT 1"},
		}
	})

	d, err := Open(context.Background(), ws.config())
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	s1, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer s1.Close()

	s2, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()

	// A failed query clears its notices instead of leaking them into the next
	// Query call on the same session.
	failed, err := s1.Query(context.Background(), "SELECT 1/0;", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var pgErr *pgconn.PgError
	if len(failed) != 1 || !errors.As(failed[0].Err, &pgErr) || pgErr.Code != "22012" {
		t.Fatalf("failed result = %+v", failed)
	}
	if len(failed[0].Messages) != 0 {
		t.Fatalf("failed query kept messages %+v", failed[0].Messages)
	}

	res1, err := s1.Query(context.Background(), "SELECT 1;", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	res2, err := s2.Query(context.Background(), "SELECT 1;", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	want1 := []driver.Message{
		{Level: "NOTICE", Content: "from conn 1"},
		{Level: "WARNING", Content: "be careful"},
		{Level: "", Content: "unknown severity"},
	}
	want2 := []driver.Message{
		{Level: "NOTICE", Content: "from conn 2"},
		{Level: "WARNING", Content: "be careful"},
		{Level: "", Content: "unknown severity"},
	}
	if !slices.Equal(res1[0].Messages, want1) || !slices.Equal(res2[0].Messages, want2) {
		t.Fatalf("s1 messages = %+v, s2 messages = %+v", res1[0].Messages, res2[0].Messages)
	}
}

func TestSessionTimestampZonePrecisionAndTimetz(t *testing.T) {
	ws := startWireServer(t, "UTC", func(_ uint32, sql string, _ [][]byte) wireReply {
		if strings.HasPrefix(sql, "SET TimeZone") {
			return wireReply{
				params: map[string]string{"TimeZone": "Asia/Shanghai"},
				tags:   []string{"SET"},
			}
		}
		return wireReply{
			fields: []pgproto3.FieldDescription{
				{Name: []byte("tz_default"), DataTypeOID: pgtype.TimestamptzOID, TypeModifier: -1},
				{Name: []byte("tz_scale3"), DataTypeOID: pgtype.TimestamptzOID, TypeModifier: 3},
				{Name: []byte("ts_default"), DataTypeOID: pgtype.TimestampOID, TypeModifier: -1},
				{Name: []byte("tz_far"), DataTypeOID: pgtype.TimestamptzOID, TypeModifier: -1},
				{Name: []byte("ts_bc"), DataTypeOID: pgtype.TimestampOID, TypeModifier: -1},
				{Name: []byte("ttz"), DataTypeOID: pgtype.TimetzOID, TypeModifier: -1},
			},
			rows: [][][]byte{{
				[]byte("2026-03-04 05:06:07.891234+00"),
				[]byte("2026-03-04 05:06:07.891+00"),
				[]byte("2026-03-04 05:06:07.891234"),
				[]byte("10000-01-01 00:00:00+00"),
				[]byte("0001-12-31 15:00:00 BC"),
				[]byte("12:34:56.123+08"),
			}},
			tags: []string{"SELECT 1"},
		}
	})

	d, err := Open(context.Background(), ws.config())
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	sess, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	results, err := sess.Query(context.Background(), "SET TimeZone TO 'Asia/Shanghai'; SELECT 1;", driver.QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[0].AffectedRows == nil || *results[0].AffectedRows != 0 {
		t.Fatalf("SET result = %+v", results)
	}
	row := results[1].Rows[0]

	if row[0].Kind != driver.TimestampTZ || row[0].Precision != 6 {
		t.Fatalf("tz_default = %+v", row[0])
	}
	if zone, offset := row[0].Time.Zone(); zone != "CST" || offset != 8*3600 {
		t.Fatalf("tz_default zone = %q %d, want CST +28800", zone, offset)
	}
	if row[1].Kind != driver.TimestampTZ || row[1].Precision != 3 {
		t.Fatalf("tz_scale3 = %+v, want Precision 3", row[1])
	}
	if row[2].Kind != driver.Timestamp || row[2].Precision != 6 || row[2].Time.Location() != time.UTC {
		t.Fatalf("ts_default = %+v", row[2])
	}
	if got := FormatTimestamp(row[3].Time, true); got != "10000-01-01 08:00:00+08:00" {
		t.Fatalf("tz_far formatted = %q", got)
	}
	if got := FormatTimestamp(row[4].Time, false); got != "0001-12-31 15:00:00 BC" {
		t.Fatalf("ts_bc formatted = %q", got)
	}
	if row[5].Kind != driver.String || row[5].String != "12:34:56.123000+08" {
		t.Fatalf("timetz = %+v, want zero-padded 12:34:56.123000+08", row[5])
	}
}

func TestSessionCancelStopsServerQuery(t *testing.T) {
	started := make(chan struct{})
	unblock := make(chan struct{})
	ws := startWireServer(t, "UTC", func(_ uint32, sql string, _ [][]byte) wireReply {
		if strings.Contains(sql, "pg_sleep") {
			close(started)
			<-unblock
			return wireReply{err: &pgproto3.ErrorResponse{Severity: "ERROR", Code: "57014", Message: "canceling statement due to user request"}}
		}
		return wireReply{tags: []string{"OK"}}
	})

	d, err := Open(context.Background(), ws.config())
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	sess, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan []*driver.Result, 1)
	go func() {
		res, qerr := sess.Query(ctx, "SELECT pg_sleep(30);", driver.QueryOptions{})
		if qerr != nil {
			t.Errorf("Query returned script error %v", qerr)
		}
		done <- res
	}()

	<-started
	cancel()

	select {
	case req := <-ws.cancels:
		if req.ProcessID == 0 {
			t.Fatal("CancelRequest carried zero ProcessID")
		}
		close(unblock)
	case <-time.After(5 * time.Second):
		t.Fatal("server never received PostgreSQL CancelRequest")
	}

	results := <-done
	if len(results) != 1 || results[0].Err == nil {
		t.Fatalf("canceled results = %+v", results)
	}
	if err := d.Ping(context.Background()); err != nil {
		t.Fatalf("Ping after cancel failed: %v", err)
	}
}

type recordedExecLog struct {
	events []string
}

func (l *recordedExecLog) CommandExecute(cmd driver.Command) {
	l.events = append(l.events, "exec:"+cmd.Text)
}

func (l *recordedExecLog) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.events = append(l.events, "resp:err:"+err.Error())
		return
	}
	var parts []string
	for _, n := range perResult {
		parts = append(parts, string(rune('0'+n)))
	}
	l.events = append(l.events, "resp:"+string(rune('0'+affected))+":["+strings.Join(parts, ",")+"]")
}

func (l *recordedExecLog) TransactionControl(op driver.TransactionOp, err error) {
	if err != nil {
		l.events = append(l.events, "tx:"+string(op)+":err")
		return
	}
	l.events = append(l.events, "tx:"+string(op))
}

func (l *recordedExecLog) Retry(_ error, attempt, max int) {
	l.events = append(l.events, "retry:"+string(rune('0'+attempt))+"/"+string(rune('0'+max)))
}

func TestExecutePgxMultiResultAndLockRetry(t *testing.T) {
	var attempts atomic.Int32
	ws := startWireServer(t, "UTC", func(_ uint32, sql string, _ [][]byte) wireReply {
		switch {
		case strings.Contains(sql, "pg_roles"):
			return wireReply{
				fields: []pgproto3.FieldDescription{{Name: []byte("rolname"), DataTypeOID: pgtype.TextOID, TypeModifier: -1}},
				rows:   [][][]byte{{[]byte("app_owner")}},
				tags:   []string{"SELECT 1"},
			}
		case strings.HasPrefix(sql, "SET LOCAL ROLE NONE;"):
			if attempts.Add(1) == 1 {
				return wireReply{
					err: &pgproto3.ErrorResponse{
						Severity: "ERROR",
						Code:     "55P03",
						Message:  "canceling statement due to lock timeout",
					},
				}
			}
			return wireReply{tags: []string{"SET", "CREATE EXTENSION", "SET"}}
		case strings.HasPrefix(strings.TrimSpace(sql), "INSERT"):
			return wireReply{tags: []string{"INSERT 0 2"}}
		default:
			return wireReply{tags: []string{"OK"}}
		}
	})

	d, err := Open(context.Background(), ws.config(), Options{TenantMode: true})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	log := &recordedExecLog{}
	affected, err := d.Execute(context.Background(), "CREATE EXTENSION pg_trgm;\nINSERT INTO t VALUES (1), (2);", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		MaxRetries:  2,
		Log:         log,
	})
	if err != nil {
		t.Fatal(err)
	}
	if affected != 2 {
		t.Fatalf("affected = %d, want 2", affected)
	}
	want := []string{
		"tx:BEGIN",
		"exec:SET LOCAL ROLE NONE;CREATE EXTENSION pg_trgm;SET LOCAL ROLE 'app_owner';",
		"resp:err:ERROR: canceling statement due to lock timeout (SQLSTATE 55P03)",
		"tx:ROLLBACK",
		"retry:1/2",
		"tx:BEGIN",
		"exec:SET LOCAL ROLE NONE;CREATE EXTENSION pg_trgm;SET LOCAL ROLE 'app_owner';",
		"resp:0:[0,0,0]",
		"exec:\nINSERT INTO t VALUES (1), (2);",
		"resp:2:[2]",
		"tx:COMMIT",
	}
	if !slices.Equal(log.events, want) {
		t.Fatalf("events =\n%v\nwant\n%v", log.events, want)
	}
}

func TestActivity(t *testing.T) {
	ws := startWireServer(t, "UTC", func(_ uint32, _ string, _ [][]byte) wireReply {
		return wireReply{
			fields: []pgproto3.FieldDescription{
				{Name: []byte("pid"), DataTypeOID: pgtype.Int4OID, TypeModifier: -1},
				{Name: []byte("blocked_by_pids"), DataTypeOID: pgtype.Int4ArrayOID, TypeModifier: -1},
				{Name: []byte("query"), DataTypeOID: pgtype.TextOID, TypeModifier: -1},
				{Name: []byte("state"), DataTypeOID: pgtype.TextOID, TypeModifier: -1},
				{Name: []byte("wait_event_type"), DataTypeOID: pgtype.TextOID, TypeModifier: -1},
				{Name: []byte("wait_event"), DataTypeOID: pgtype.TextOID, TypeModifier: -1},
				{Name: []byte("datname"), DataTypeOID: pgtype.TextOID, TypeModifier: -1},
				{Name: []byte("usename"), DataTypeOID: pgtype.TextOID, TypeModifier: -1},
				{Name: []byte("application_name"), DataTypeOID: pgtype.TextOID, TypeModifier: -1},
				{Name: []byte("client_addr"), DataTypeOID: pgtype.InetOID, TypeModifier: -1},
				{Name: []byte("client_port"), DataTypeOID: pgtype.Int4OID, TypeModifier: -1},
				{Name: []byte("backend_start"), DataTypeOID: pgtype.TimestamptzOID, TypeModifier: -1},
				{Name: []byte("xact_start"), DataTypeOID: pgtype.TimestamptzOID, TypeModifier: -1},
				{Name: []byte("query_start"), DataTypeOID: pgtype.TimestamptzOID, TypeModifier: -1},
			},
			rows: [][][]byte{
				{
					[]byte("100"), []byte("{200}"), []byte("ALTER TABLE t ADD COLUMN c int"),
					[]byte("active"), []byte("Lock"), []byte("relation"),
					[]byte("sales"), []byte("dba"), []byte("bytebase-taskrun-42"),
					[]byte("10.0.0.1"), []byte("51234"),
					[]byte("2026-03-04 05:00:00+00"), []byte("2026-03-04 05:01:00+00"), []byte("2026-03-04 05:01:02+00"),
				},
				{
					[]byte("200"), []byte("{}"), []byte("SELECT * FROM t FOR UPDATE"),
					[]byte("idle in transaction"), nil, nil,
					[]byte("sales"), []byte("other"), []byte("psql"),
					nil, nil,
					[]byte("2026-03-04 04:55:00+00"), []byte("2026-03-04 04:56:00+00"), nil,
				},
				{
					[]byte("300"), []byte("{100}"), []byte("SELECT * FROM t"),
					[]byte("active"), []byte("Lock"), []byte("relation"),
					[]byte("sales"), []byte("reader"), []byte("app"),
					[]byte("10.0.0.2"), []byte("51235"),
					[]byte("2026-03-04 05:02:00+00"), nil, []byte("2026-03-04 05:02:01+00"),
				},
			},
			tags: []string{"SELECT 3"},
		}
	})

	d, err := Open(context.Background(), ws.config())
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	act, err := d.Activity(context.Background(), "bytebase-taskrun-42")
	if err != nil {
		t.Fatal(err)
	}
	if act.Session == nil || act.Session.PID != "100" || !slices.Equal(act.Session.BlockedByPIDs, []string{"200"}) {
		t.Fatalf("Session = %+v", act.Session)
	}
	if len(act.BlockingSessions) != 1 || act.BlockingSessions[0].PID != "200" || !slices.Equal(act.BlockingSessions[0].BlockedByPIDs, []string{}) {
		t.Fatalf("BlockingSessions = %+v", act.BlockingSessions)
	}
	if len(act.BlockedSessions) != 1 || act.BlockedSessions[0].PID != "300" {
		t.Fatalf("BlockedSessions = %+v", act.BlockedSessions)
	}
}
