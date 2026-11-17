package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
)

func baseConfig() driver.Config {
	return driver.Config{User: "dba", Host: "db.example.com", Port: "3306"}
}

func TestConnConfigDefaults(t *testing.T) {
	mc, err := connConfig(baseConfig())
	if err != nil {
		t.Fatal(err)
	}
	if !mc.MultiStatements || mc.MaxAllowedPacket != 0 {
		t.Errorf("MultiStatements=%v MaxAllowedPacket=%d, want true and 0", mc.MultiStatements, mc.MaxAllowedPacket)
	}
	if mc.Net != "tcp" || mc.Addr != "db.example.com:3306" || mc.User != "dba" {
		t.Errorf("Net=%q Addr=%q User=%q", mc.Net, mc.Addr, mc.User)
	}
}

func TestConnConfigSettingsAndUnixSocket(t *testing.T) {
	cfg := driver.Config{
		User:     "app?allowAllFiles=true",
		Password: "p@ss/word?tls=false&allowAllFiles=true",
		Host:     "/var/run/mysqld/mysqld.sock",
		Database: "sales",
		Params: map[string]string{
			"timeout":                 "10s",
			"readTimeout":             "30s",
			"allowCleartextPasswords": "true",
		},
	}
	mc, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if mc.Net != "unix" || mc.Addr != "/var/run/mysqld/mysqld.sock" {
		t.Errorf("Net=%q Addr=%q", mc.Net, mc.Addr)
	}
	// Credentials are set directly on gomysql.Config and never parsed as DSN syntax.
	if mc.User != cfg.User || mc.Passwd != cfg.Password || mc.AllowAllFiles {
		t.Errorf("User=%q Passwd=%q AllowAllFiles=%v", mc.User, mc.Passwd, mc.AllowAllFiles)
	}
	if mc.DBName != "sales" || mc.Timeout != 10*time.Second || mc.ReadTimeout != 30*time.Second || !mc.AllowCleartextPasswords {
		t.Errorf("DBName=%q Timeout=%v ReadTimeout=%v Cleartext=%v", mc.DBName, mc.Timeout, mc.ReadTimeout, mc.AllowCleartextPasswords)
	}
}

func TestConnConfigRejectsAllowAllFiles(t *testing.T) {
	for _, key := range []string{"allowAllFiles", "allowallfiles", "AllowAllFiles", "  allowAllFiles  "} {
		cfg := baseConfig()
		cfg.Params = map[string]string{key: "true"}
		if _, err := connConfig(cfg); err == nil || !strings.Contains(strings.ToLower(err.Error()), "allowallfiles") {
			t.Errorf("param %q: err = %v, want allowAllFiles rejection", key, err)
		}
	}
}

func TestConnConfigTLSCannotBeDowngraded(t *testing.T) {
	cfg := baseConfig()
	cfg.Params = map[string]string{
		"tls":                      "false",
		"allowFallbackToPlaintext": "true",
	}
	var gotHost string
	cfg.TLS = func(host string) (*tls.Config, error) {
		gotHost = host
		return &tls.Config{ServerName: host}, nil
	}
	mc, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if gotHost != "db.example.com" || mc.TLS == nil || mc.TLS.ServerName != "db.example.com" || mc.AllowFallbackToPlaintext {
		t.Fatalf("gotHost=%q TLS=%v AllowFallbackToPlaintext=%v", gotHost, mc.TLS, mc.AllowFallbackToPlaintext)
	}

	// With cfg.TLS nil, tls=true in Params infers ServerName from cfg.Host.
	paramTLS := baseConfig()
	paramTLS.Params = map[string]string{"tls": "true"}
	mcParam, err := connConfig(paramTLS)
	if err != nil {
		t.Fatal(err)
	}
	if mcParam.TLS == nil || mcParam.TLS.ServerName != "db.example.com" {
		t.Fatalf("Params tls=true ServerName = %+v, want db.example.com", mcParam.TLS)
	}
}

func TestConnConfigDialWithoutHost(t *testing.T) {
	var dialed bool
	cfg := driver.Config{
		User: "iam-user",
		Dial: func(context.Context, string, string) (net.Conn, error) {
			dialed = true
			return nil, net.ErrClosed
		},
	}
	mc, err := connConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if mc.DialFunc == nil {
		t.Fatal("DialFunc was not set")
	}
	_, _ = mc.DialFunc(context.Background(), "tcp", mc.Addr)
	if !dialed {
		t.Error("DialFunc was not invoked")
	}
}

func TestConnConfigValidation(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*driver.Config)
		want string
	}{
		{"blank user", func(c *driver.Config) { c.User = "" }, "user must be set"},
		{"blank host", func(c *driver.Config) { c.Host = "  " }, "host must be set"},
		{"blank port", func(c *driver.Config) { c.Port = "  " }, "port must be set"},
	} {
		cfg := baseConfig()
		tc.edit(&cfg)
		if _, err := connConfig(cfg); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err = %v, want %q", tc.name, err, tc.want)
		}
	}
}

func TestOpenDoesNotConnect(t *testing.T) {
	cfg := baseConfig()
	cfg.Host, cfg.Port = "127.0.0.1", "1"
	d, err := Open(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Ping(context.Background()); err == nil {
		t.Error("Ping reached a server on port 1")
	}
	if err := d.Close(); err != nil {
		t.Error(err)
	}
}

func TestParseVersion(t *testing.T) {
	for _, tc := range []struct {
		in, wantVer, wantRest string
	}{
		{"8.0.27", "8.0.27", ""},
		{"5.7.22-log", "5.7.22", "-log"},
		{"5.6.29_ddm_3.0.1.7", "5.6.29", "_ddm_3.0.1.7"},
		{"10.4.7-MariaDB", "10.4.7", "-MariaDB"},
	} {
		ver, rest, err := parseVersion(tc.in)
		if err != nil || ver != tc.wantVer || rest != tc.wantRest {
			t.Errorf("parseVersion(%q) = (%q, %q, %v)", tc.in, ver, rest, err)
		}
	}
}

// TestOpenSyncDatabaseTranslatesHandshake1049ToErrNotFound exercises the public
// Open path when Config.Database names a database the server rejects during
// handshake with MySQL error 1049 (ER_BAD_DB_ERROR).
func TestOpenSyncDatabaseTranslatesHandshake1049ToErrNotFound(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		writePacket := func(seq byte, payload []byte) error {
			hdr := []byte{byte(len(payload)), byte(len(payload) >> 8), byte(len(payload) >> 16), seq}
			if _, err := conn.Write(hdr); err != nil {
				return err
			}
			_, err := conn.Write(payload)
			return err
		}

		// Minimal MySQL v10 initial handshake packet (clientProtocol41 |
		// clientSecureConn | clientConnectWithDB | clientPluginAuth).
		var greeting []byte
		greeting = append(greeting, 0x0a)
		greeting = append(greeting, "8.0.46\x00"...)
		greeting = append(greeting, 1, 0, 0, 0)
		greeting = append(greeting, "12345678"...)
		greeting = append(greeting, 0x00)
		// Capability flags low: CLIENT_LONG_PASSWORD(1) | CLIENT_CONNECT_WITH_DB(8) |
		// CLIENT_PROTOCOL_41(0x200) | CLIENT_SECURE_CONNECTION(0x8000) = 0x8209.
		greeting = append(greeting, 0x09, 0x82)
		greeting = append(greeting, 0x21)
		greeting = append(greeting, 0x02, 0x00)
		// Capability flags high: CLIENT_PLUGIN_AUTH(0x00080000 >> 16 = 0x0008).
		greeting = append(greeting, 0x08, 0x00)
		greeting = append(greeting, 21)
		greeting = append(greeting, make([]byte, 10)...)
		greeting = append(greeting, "123456789012\x00"...)
		greeting = append(greeting, "mysql_native_password\x00"...)
		if err := writePacket(0, greeting); err != nil {
			return
		}

		// Read client HandshakeResponse41 header + body.
		var hdr [4]byte
		if _, err := io.ReadFull(conn, hdr[:]); err != nil {
			return
		}
		bodyLen := int(hdr[0]) | int(hdr[1])<<8 | int(hdr[2])<<16
		if _, err := io.CopyN(io.Discard, conn, int64(bodyLen)); err != nil {
			return
		}

		// Reply with ERR_Packet: 0xff, errno 1049 (0x0419 LE), '#42000', message.
		var errPkt []byte
		errPkt = append(errPkt, 0xff, 0x19, 0x04)
		errPkt = append(errPkt, "#42000Unknown database 'missing_db'"...)
		_ = writePacket(2, errPkt)
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	d, err := Open(context.Background(), driver.Config{
		Host:     host,
		Port:     port,
		Database: "missing_db",
		User:     "root",
		Password: "secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	_, err = d.SyncDatabase(context.Background())
	if !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase err = %v, want driver.ErrNotFound", err)
	}
	var mysqlErr *gomysql.MySQLError
	if !errors.As(err, &mysqlErr) || mysqlErr.Number != 1049 {
		t.Fatalf("SyncDatabase err = %v, want wrapped MySQLError 1049", err)
	}
}

// TestOpenSessionKillQueryOnCancel is the kill-on-cancel path over the public
// Open: a canceled statement is stopped on the server with KILL QUERY of the
// connection id the driver read when it opened the session.
func TestOpenSessionKillQueryOnCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	sleepStarted := make(chan struct{})
	killed := make(chan string, 2)

	appendPacket := func(dst []byte, seq byte, payload []byte) []byte {
		dst = append(dst, byte(len(payload)), byte(len(payload)>>8), byte(len(payload)>>16), seq)
		return append(dst, payload...)
	}
	writePacket := func(conn net.Conn, seq byte, payload []byte) error {
		_, err := conn.Write(appendPacket(nil, seq, payload))
		return err
	}
	writeSingleCell := func(c net.Conn, val string) {
		var out []byte
		out = appendPacket(out, 1, []byte{1})
		col := []byte{3, 'd', 'e', 'f', 0, 0, 0, 2, 'i', 'd', 0, 0x0c, 0x21, 0x00, 10, 0, 0, 0, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00}
		out = appendPacket(out, 2, col)
		out = appendPacket(out, 3, []byte{0xfe, 0, 0, 0x02, 0x00})
		row := append([]byte{byte(len(val))}, val...)
		out = appendPacket(out, 4, row)
		out = appendPacket(out, 5, []byte{0xfe, 0, 0, 0x02, 0x00})
		_, _ = c.Write(out)
	}
	readPacket := func(conn net.Conn) ([]byte, error) {
		var hdr [4]byte
		if _, err := io.ReadFull(conn, hdr[:]); err != nil {
			return nil, err
		}
		n := int(hdr[0]) | int(hdr[1])<<8 | int(hdr[2])<<16
		buf := make([]byte, n)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
		return buf, nil
	}
	greeting := func() []byte {
		var g []byte
		g = append(g, 0x0a)
		g = append(g, "8.0.46\x00"...)
		g = append(g, 1, 0, 0, 0)
		g = append(g, "12345678"...)
		g = append(g, 0x00, 0x09, 0x82, 0x21, 0x02, 0x00, 0x08, 0x00, 21)
		g = append(g, make([]byte, 10)...)
		g = append(g, "123456789012\x00"...)
		g = append(g, "mysql_native_password\x00"...)
		return g
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(5 * time.Second))
				if err := writePacket(c, 0, greeting()); err != nil {
					return
				}
				if _, err := readPacket(c); err != nil {
					return
				}
				if err := writePacket(c, 2, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}); err != nil {
					return
				}
				for {
					pkt, err := readPacket(c)
					if err != nil || len(pkt) == 0 || pkt[0] == 0x01 {
						return
					}
					if pkt[0] != 0x03 {
						_ = writePacket(c, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})
						continue
					}
					sql := string(pkt[1:])
					switch {
					case strings.Contains(sql, "max_allowed_packet"):
						writeSingleCell(c, "67108864")
					case strings.Contains(sql, "CONNECTION_ID()"):
						writeSingleCell(c, "42")
					case strings.Contains(sql, "SLEEP"):
						close(sleepStarted)
						<-connClosed(c)
						return
					case strings.HasPrefix(sql, "KILL QUERY "):
						killed <- sql
						_ = writePacket(c, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})
					default:
						_ = writePacket(c, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})
					}
				}
			}(conn)
		}
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	d, err := Open(context.Background(), driver.Config{
		Host: host, Port: port, User: "root", ApplicationName: "bytebase",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()

	sess, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sess.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan []*driver.Result, 1)
	go func() {
		res, _ := sess.Query(ctx, "SELECT SLEEP(30);", driver.QueryOptions{})
		done <- res
	}()
	<-sleepStarted
	cancel()

	select {
	case got := <-killed:
		if got != "KILL QUERY 42" {
			t.Fatalf("kill SQL = %q, want KILL QUERY 42", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("never sent KILL QUERY on context cancel")
	}
	res := <-done
	if len(res) != 1 || res[0].Err == nil {
		t.Fatalf("canceled result = %+v", res)
	}
}

func connClosed(c net.Conn) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		var b [1]byte
		_, _ = c.Read(b[:])
	}()
	return ch
}

func TestBuildExecuteCommandsDelimiterAndBounds(t *testing.T) {
	const procedure = "DELIMITER //\nCREATE PROCEDURE p()\nBEGIN\n  SELECT 1;\nEND//\nDELIMITER ;\n"
	commands := buildExecuteCommands(procedure)
	if len(commands) != 1 || strings.Contains(commands[0].Text, "DELIMITER") || !strings.Contains(commands[0].Text, "CREATE PROCEDURE p()") {
		t.Fatalf("delimiter script commands = %+v", commands)
	}

	var many strings.Builder
	for range maxSplitCommands + 5 {
		many.WriteString("SELECT 1;\n")
	}
	if got := buildExecuteCommands(many.String()); len(got) != 1 || got[0].Text != many.String() {
		t.Fatalf("too-many-commands fallback = %d commands", len(got))
	}

	large := "SELECT 1;\n" + strings.Repeat(" ", maxSheetCheckBytes)
	if got := buildExecuteCommands(large); len(got) != 1 || got[0].Text != large {
		t.Fatalf("large-sheet fallback = %d commands", len(got))
	}

	// A script with directives stays split past both limits: the directives
	// cannot reach the server.
	var manyProcedures strings.Builder
	manyProcedures.WriteString("DELIMITER //\n")
	for range maxSplitCommands + 5 {
		manyProcedures.WriteString("SELECT 1//\n")
	}
	manyProcedures.WriteString("DELIMITER ;\n")
	got := buildExecuteCommands(manyProcedures.String())
	if len(got) != maxSplitCommands+5 {
		t.Fatalf("delimiter script past the command limit = %d commands, want %d", len(got), maxSplitCommands+5)
	}
	for _, command := range got {
		if strings.Contains(strings.ToUpper(command.Text), "DELIMITER") {
			t.Fatalf("command holds a directive: %q", command.Text)
		}
	}

	if got := buildExecuteCommands("  ;\n-- nothing\n"); len(got) != 0 {
		t.Fatalf("empty script = %+v", got)
	}
}

// TestSessionCloseDiscardsTheConnection pins that a session that selected
// another database does not hand it to the pool: a later Execute on the same
// driver still writes the driver's database.
func TestSessionCloseDiscardsTheConnection(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	var (
		mu         sync.Mutex
		handshakes int
		statements []string
	)
	appendPacket := func(dst []byte, seq byte, payload []byte) []byte {
		dst = append(dst, byte(len(payload)), byte(len(payload)>>8), byte(len(payload)>>16), seq)
		return append(dst, payload...)
	}
	readPacket := func(conn net.Conn) ([]byte, error) {
		var hdr [4]byte
		if _, err := io.ReadFull(conn, hdr[:]); err != nil {
			return nil, err
		}
		n := int(hdr[0]) | int(hdr[1])<<8 | int(hdr[2])<<16
		buf := make([]byte, n)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
		return buf, nil
	}
	writeSingleCell := func(c net.Conn, val string) {
		var out []byte
		out = appendPacket(out, 1, []byte{1})
		col := []byte{3, 'd', 'e', 'f', 0, 0, 0, 2, 'i', 'd', 0, 0x0c, 0x21, 0x00, 10, 0, 0, 0, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00}
		out = appendPacket(out, 2, col)
		out = appendPacket(out, 3, []byte{0xfe, 0, 0, 0x02, 0x00})
		row := append([]byte{byte(len(val))}, val...)
		out = appendPacket(out, 4, row)
		out = appendPacket(out, 5, []byte{0xfe, 0, 0, 0x02, 0x00})
		_, _ = c.Write(out)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(5 * time.Second))
				mu.Lock()
				handshakes++
				mu.Unlock()

				var greeting []byte
				greeting = append(greeting, 0x0a)
				greeting = append(greeting, "8.0.46\x00"...)
				greeting = append(greeting, 1, 0, 0, 0)
				greeting = append(greeting, "12345678"...)
				greeting = append(greeting, 0x00, 0x09, 0x82, 0x21, 0x02, 0x00, 0x08, 0x00, 21)
				greeting = append(greeting, make([]byte, 10)...)
				greeting = append(greeting, "123456789012\x00"...)
				greeting = append(greeting, "mysql_native_password\x00"...)
				if _, err := c.Write(appendPacket(nil, 0, greeting)); err != nil {
					return
				}
				if _, err := readPacket(c); err != nil {
					return
				}
				if _, err := c.Write(appendPacket(nil, 2, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})); err != nil {
					return
				}
				for {
					pkt, err := readPacket(c)
					if err != nil || len(pkt) == 0 || pkt[0] == 0x01 {
						return
					}
					if pkt[0] != 0x03 {
						_, _ = c.Write(appendPacket(nil, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}))
						continue
					}
					sql := string(pkt[1:])
					mu.Lock()
					statements = append(statements, sql)
					mu.Unlock()
					if strings.Contains(sql, "max_allowed_packet") {
						writeSingleCell(c, "67108864")
						continue
					}
					if strings.Contains(sql, "CONNECTION_ID()") {
						writeSingleCell(c, "7")
						continue
					}
					_, _ = c.Write(appendPacket(nil, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}))
				}
			}(conn)
		}
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	d, err := Open(context.Background(), driver.Config{Host: host, Port: port, User: "root", Database: "app"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()

	sess, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sess.Query(context.Background(), "USE other;", driver.QueryOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := sess.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Execute(context.Background(), "UPDATE t SET a = 1;", driver.ExecuteOptions{Transaction: driver.AutoCommit}); err != nil {
		t.Fatal(err)
	}
	// A script that changes the session is not pooled either.
	if _, err := d.Execute(context.Background(), "UPDATE t SET a = 2;", driver.ExecuteOptions{Transaction: driver.AutoCommit}); err != nil {
		t.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	if handshakes < 3 {
		t.Errorf("handshakes = %d, want each execute and the session to open a connection of its own", handshakes)
	}
	var sawUse bool
	for _, statement := range statements {
		if strings.Contains(statement, "USE other") {
			sawUse = true
		}
	}
	if !sawUse {
		t.Errorf("statements = %v", statements)
	}
	if last := statements[len(statements)-1]; !strings.Contains(last, "UPDATE t SET a = 2") {
		t.Errorf("last statement = %q", last)
	}
}

// TestExecuteRetriesLockWaitTimeout pins the retry loop: a script that fails on
// MySQL's lock wait timeout runs again, each attempt from the beginning, and
// the logger sees the retry.
// lockServer is a fake MySQL server that counts the statements it receives and
// answers a lock wait timeout for the ones the test names. It answers the two
// probes a connection makes (max_allowed_packet, CONNECTION_ID) and OK for
// everything else.
type lockServer struct {
	mu sync.Mutex
	// ran counts the statements by marker.
	ran map[string]int
	// failAlways times out every matching statement; failFirst times out its
	// first one only.
	failAlways []string
	failFirst  []string
	seenFirst  map[string]bool
}

func startLockServer(t *testing.T) (*lockServer, string, string) {
	t.Helper()
	s := &lockServer{ran: map[string]int{}, seenFirst: map[string]bool{}}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	appendPacket := func(dst []byte, seq byte, payload []byte) []byte {
		dst = append(dst, byte(len(payload)), byte(len(payload)>>8), byte(len(payload)>>16), seq)
		return append(dst, payload...)
	}
	readPacket := func(conn net.Conn) ([]byte, error) {
		var hdr [4]byte
		if _, err := io.ReadFull(conn, hdr[:]); err != nil {
			return nil, err
		}
		n := int(hdr[0]) | int(hdr[1])<<8 | int(hdr[2])<<16
		buf := make([]byte, n)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
		return buf, nil
	}
	writeSingleCell := func(c net.Conn, val string) {
		var out []byte
		out = appendPacket(out, 1, []byte{1})
		col := []byte{3, 'd', 'e', 'f', 0, 0, 0, 2, 'i', 'd', 0, 0x0c, 0x21, 0x00, 10, 0, 0, 0, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00}
		out = appendPacket(out, 2, col)
		out = appendPacket(out, 3, []byte{0xfe, 0, 0, 0x02, 0x00})
		row := append([]byte{byte(len(val))}, val...)
		out = appendPacket(out, 4, row)
		out = appendPacket(out, 5, []byte{0xfe, 0, 0, 0x02, 0x00})
		_, _ = c.Write(out)
	}
	ok := func(c net.Conn) {
		_, _ = c.Write(appendPacket(nil, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}))
	}
	lockTimeout := func(c net.Conn) {
		var errPkt []byte
		errPkt = append(errPkt, 0xff, 0xb5, 0x04)
		errPkt = append(errPkt, "#HY000Lock wait timeout exceeded; try restarting transaction"...)
		_, _ = c.Write(appendPacket(nil, 1, errPkt))
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(5 * time.Second))
				var greeting []byte
				greeting = append(greeting, 0x0a)
				greeting = append(greeting, "8.0.46\x00"...)
				greeting = append(greeting, 1, 0, 0, 0)
				greeting = append(greeting, "12345678"...)
				greeting = append(greeting, 0x00, 0x09, 0x82, 0x21, 0x02, 0x00, 0x08, 0x00, 21)
				greeting = append(greeting, make([]byte, 10)...)
				greeting = append(greeting, "123456789012\x00"...)
				greeting = append(greeting, "mysql_native_password\x00"...)
				if _, err := c.Write(appendPacket(nil, 0, greeting)); err != nil {
					return
				}
				if _, err := readPacket(c); err != nil {
					return
				}
				if _, err := c.Write(appendPacket(nil, 2, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})); err != nil {
					return
				}
				for {
					pkt, err := readPacket(c)
					if err != nil || len(pkt) == 0 || pkt[0] == 0x01 {
						return
					}
					if pkt[0] != 0x03 {
						ok(c)
						continue
					}
					sql := string(pkt[1:])
					switch {
					case strings.Contains(sql, "max_allowed_packet"):
						writeSingleCell(c, "67108864")
					case strings.Contains(sql, "CONNECTION_ID()"):
						writeSingleCell(c, "7")
					default:
						key := statementKey(sql)
						s.mu.Lock()
						s.ran[key]++
						count := s.ran[key]
						timedOut := slices.Contains(s.failAlways, key) || (slices.Contains(s.failFirst, key) && count == 1)
						s.mu.Unlock()
						if timedOut {
							lockTimeout(c)
							continue
						}
						ok(c)
					}
				}
			}(conn)
		}
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	return s, host, port
}

// statementKey names a statement for the counters, from the marker a test puts
// in it.
func statementKey(sql string) string {
	for _, marker := range []string{"UPDATE t", "UPDATE locked", "UPDATE stable", "CREATE TABLE kept", "CREATE TABLE t", "ALTER TABLE t"} {
		if strings.Contains(sql, marker) {
			return marker
		}
	}
	return "other"
}

func (s *lockServer) count(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ran[key]
}

// TestExecuteRetriesLockWaitTimeout pins when Execute retries: only an attempt
// that committed nothing may run again, and the logger sees each retry.
func TestExecuteRetriesLockWaitTimeout(t *testing.T) {
	server, host, port := startLockServer(t)
	server.failFirst = []string{"UPDATE t", "ALTER TABLE t"}
	server.failAlways = []string{"UPDATE locked"}

	d, err := Open(context.Background(), driver.Config{Host: host, Port: port, User: "root"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()

	log := &executeStepLogger{}
	affected, err := d.Execute(context.Background(), "UPDATE t SET a = 1;", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		MaxRetries:  2,
		Log:         log,
	})
	if err != nil {
		t.Fatalf("Execute = (%d, %v)", affected, err)
	}
	if got := server.count("UPDATE t"); got != 2 {
		t.Fatalf("UPDATE ran %d times, want 2", got)
	}
	want := []string{"tx:BEGIN", "exec[0:19]", "resp:err", "tx:ROLLBACK", "retry:1/2", "tx:BEGIN", "exec[0:19]", "resp:0:[0]", "tx:COMMIT"}
	if !slices.Equal(log.steps, want) {
		t.Fatalf("steps = %v, want %v", log.steps, want)
	}

	// Auto-commit commits the first statement, so a lock timeout on the second
	// must not replay the script.
	autoLog := &executeStepLogger{}
	if _, err := d.Execute(context.Background(), "UPDATE stable SET a = 1;\nUPDATE locked SET a = 1;", driver.ExecuteOptions{
		Transaction: driver.AutoCommit,
		MaxRetries:  2,
		Log:         autoLog,
	}); err == nil {
		t.Fatal("expected the lock timeout")
	}
	if got := server.count("UPDATE stable"); got != 1 {
		t.Errorf("the committed statement ran %d times, want 1", got)
	}
	if got := server.count("UPDATE locked"); got != 1 {
		t.Errorf("the failing statement ran %d times, want 1", got)
	}
	for _, step := range autoLog.steps {
		if strings.HasPrefix(step, "retry:") {
			t.Fatalf("auto-commit retried after a committed statement: %v", autoLog.steps)
		}
	}

	// A DDL statement that fails on a lock timeout commits the DML before it
	// before it runs, so the failure cannot be retried even though nothing
	// succeeded before it.
	alterLog := &executeStepLogger{}
	if _, err := d.Execute(context.Background(), "UPDATE t SET a = 2;\nALTER TABLE t ADD COLUMN b INT;", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		MaxRetries:  2,
		Log:         alterLog,
	}); err == nil {
		t.Fatal("expected the lock timeout")
	}
	if got := server.count("ALTER TABLE t"); got != 1 {
		t.Errorf("the failing DDL ran %d times, want 1", got)
	}
	for _, step := range alterLog.steps {
		if strings.HasPrefix(step, "retry:") {
			t.Fatalf("transactional retried after a failing DDL: %v", alterLog.steps)
		}
	}

	// A DDL statement inside a transaction commits implicitly, so the failure
	// after it is not repeatable either.
	ddlLog := &executeStepLogger{}
	if _, err := d.Execute(context.Background(), "CREATE TABLE kept (a INT);\nUPDATE locked SET a = 1;", driver.ExecuteOptions{
		Transaction: driver.Transactional,
		MaxRetries:  2,
		Log:         ddlLog,
	}); err == nil {
		t.Fatal("expected the lock timeout")
	}
	if got := server.count("CREATE TABLE kept"); got != 1 {
		t.Errorf("the implicitly committed statement ran %d times, want 1", got)
	}
	for _, step := range ddlLog.steps {
		if strings.HasPrefix(step, "retry:") {
			t.Fatalf("transactional retried after an implicit commit: %v", ddlLog.steps)
		}
	}
}

type executeStepLogger struct{ steps []string }

func (l *executeStepLogger) CommandExecute(cmd driver.Command) {
	l.steps = append(l.steps, fmt.Sprintf("exec[%d:%d]", cmd.ByteStart, cmd.ByteEnd))
}

func (l *executeStepLogger) CommandResponse(affected int64, perResult []int64, err error) {
	if err != nil {
		l.steps = append(l.steps, "resp:err")
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("resp:%d:%v", affected, perResult))
}

func (l *executeStepLogger) TransactionControl(op driver.TransactionOp, err error) {
	if err != nil {
		l.steps = append(l.steps, fmt.Sprintf("tx:%s:err", op))
		return
	}
	l.steps = append(l.steps, fmt.Sprintf("tx:%s", op))
}

func (l *executeStepLogger) Retry(_ error, attempt, max int) {
	l.steps = append(l.steps, fmt.Sprintf("retry:%d/%d", attempt, max))
}

func TestSanitizeCommentToken(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"bytebase", "bytebase"},
		{"a*/b", "a* /b"},
		{"**//", "** //"},
		{"*/ SELECT 1", "* / SELECT 1"},
	} {
		if got := sanitizeCommentToken(tc.in); got != tc.want {
			t.Errorf("sanitizeCommentToken(%q) = %q, want %q", tc.in, got, tc.want)
		}
		if strings.Contains(sanitizeCommentToken(tc.in), "*/") {
			t.Errorf("sanitizeCommentToken(%q) still holds a comment terminator", tc.in)
		}
	}
}

func TestParseTimestampOrStringValue(t *testing.T) {
	// The default text form, the form database/sql renders a parseTime
	// time.Time as in the default location and in one east of it, and a value
	// that is not a timestamp at all.
	for _, tc := range []struct {
		raw  string
		kind driver.Kind
		want time.Time
	}{
		{"2026-03-04 05:06:07.123456", driver.Timestamp, time.Date(2026, 3, 4, 5, 6, 7, 123456000, time.UTC)},
		{"2026-03-04 05:06:07", driver.Timestamp, time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)},
		{"2026-03-04T05:06:07.123456Z", driver.Timestamp, time.Date(2026, 3, 4, 5, 6, 7, 123456000, time.UTC)},
		// parseTime with loc renders this form, whose instant is the previous
		// day in UTC; the wall clock is the one the column holds.
		{"2026-03-04T05:06:07.123456+08:00", driver.Timestamp, time.Date(2026, 3, 4, 5, 6, 7, 123456000, time.UTC)},
		{"0000-00-00 00:00:00", driver.String, time.Time{}},
		{"not a time", driver.String, time.Time{}},
	} {
		got := parseTimestampOrStringValue(tc.raw, 6)
		if got.Kind != tc.kind {
			t.Errorf("parseTimestampOrStringValue(%q) kind = %v, want %v", tc.raw, got.Kind, tc.kind)
			continue
		}
		if tc.kind == driver.Timestamp {
			if !got.Time.Equal(tc.want) || got.Time.Location() != time.UTC || got.Precision != 6 {
				t.Errorf("parseTimestampOrStringValue(%q) = %+v, want %v at UTC with precision 6", tc.raw, got, tc.want)
			}
		}
	}
}

// TestQueryReportsAnErrorAfterAResultSet pins that a procedure's error, which
// the server states after the result set it already sent, reaches the caller:
// only the first set is the result, but the rest must be drained to see it.
func TestQueryReportsAnErrorAfterAResultSet(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	appendPacket := func(dst []byte, seq byte, payload []byte) []byte {
		dst = append(dst, byte(len(payload)), byte(len(payload)>>8), byte(len(payload)>>16), seq)
		return append(dst, payload...)
	}
	readPacket := func(conn net.Conn) ([]byte, error) {
		var hdr [4]byte
		if _, err := io.ReadFull(conn, hdr[:]); err != nil {
			return nil, err
		}
		n := int(hdr[0]) | int(hdr[1])<<8 | int(hdr[2])<<16
		buf := make([]byte, n)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
		return buf, nil
	}
	writeSingleCell := func(c net.Conn, val string) {
		var out []byte
		out = appendPacket(out, 1, []byte{1})
		col := []byte{3, 'd', 'e', 'f', 0, 0, 0, 2, 'i', 'd', 0, 0x0c, 0x21, 0x00, 10, 0, 0, 0, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00}
		out = appendPacket(out, 2, col)
		out = appendPacket(out, 3, []byte{0xfe, 0, 0, 0x02, 0x00})
		row := append([]byte{byte(len(val))}, val...)
		out = appendPacket(out, 4, row)
		out = appendPacket(out, 5, []byte{0xfe, 0, 0, 0x02, 0x00})
		_, _ = c.Write(out)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(5 * time.Second))
				var greeting []byte
				greeting = append(greeting, 0x0a)
				greeting = append(greeting, "8.0.46\x00"...)
				greeting = append(greeting, 1, 0, 0, 0)
				greeting = append(greeting, "12345678"...)
				greeting = append(greeting, 0x00, 0x09, 0x82, 0x21, 0x02, 0x00, 0x08, 0x00, 21)
				greeting = append(greeting, make([]byte, 10)...)
				greeting = append(greeting, "123456789012\x00"...)
				greeting = append(greeting, "mysql_native_password\x00"...)
				if _, err := c.Write(appendPacket(nil, 0, greeting)); err != nil {
					return
				}
				if _, err := readPacket(c); err != nil {
					return
				}
				if _, err := c.Write(appendPacket(nil, 2, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})); err != nil {
					return
				}
				for {
					pkt, err := readPacket(c)
					if err != nil || len(pkt) == 0 || pkt[0] == 0x01 {
						return
					}
					if pkt[0] != 0x03 {
						_, _ = c.Write(appendPacket(nil, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}))
						continue
					}
					sql := string(pkt[1:])
					switch {
					case strings.Contains(sql, "max_allowed_packet"):
						writeSingleCell(c, "67108864")
					case strings.Contains(sql, "CONNECTION_ID()"):
						writeSingleCell(c, "7")
					case strings.Contains(sql, "CALL"):
						// One result set, then the error the procedure raised.
						var out []byte
						out = appendPacket(out, 1, []byte{1})
						col := []byte{3, 'd', 'e', 'f', 0, 0, 0, 2, 'i', 'd', 0, 0x0c, 0x21, 0x00, 10, 0, 0, 0, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00}
						out = appendPacket(out, 2, col)
						out = appendPacket(out, 3, []byte{0xfe, 0, 0, 0x02, 0x00})
						out = appendPacket(out, 4, []byte{1, '1'})
						// The EOF states SERVER_MORE_RESULTS_EXISTS (0x0008) beside
						// autocommit, which is what makes the client ask for the
						// next one.
						out = appendPacket(out, 5, []byte{0xfe, 0, 0, 0x0a, 0x00})
						errPkt := []byte{0xff, 0x6c, 0x06}
						errPkt = append(errPkt, "#45000the procedure failed"...)
						out = appendPacket(out, 6, errPkt)
						_, _ = c.Write(out)
					default:
						_, _ = c.Write(appendPacket(nil, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}))
					}
				}
			}(conn)
		}
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	d, err := Open(context.Background(), driver.Config{Host: host, Port: port, User: "root"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()

	sess, err := d.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sess.Close() }()

	res, err := sess.Query(context.Background(), "CALL failing();", driver.QueryOptions{})
	if err != nil || len(res) != 1 {
		t.Fatalf("CALL = %+v, %v", res, err)
	}
	if res[0].Err == nil || !strings.Contains(res[0].Err.Error(), "the procedure failed") {
		t.Fatalf("CALL err = %v, want the procedure's error", res[0].Err)
	}
}
