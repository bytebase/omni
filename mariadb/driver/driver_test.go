package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
)

func baseConfig() driver.Config {
	return driver.Config{User: "dba", Host: "db.example.com", Port: "3306"}
}

func TestConnConfigAndOpen(t *testing.T) {
	mc, err := connConfig(baseConfig())
	if err != nil {
		t.Fatal(err)
	}
	if !mc.MultiStatements || mc.MaxAllowedPacket != 0 || mc.Net != "tcp" || mc.Addr != "db.example.com:3306" {
		t.Fatalf("unexpected defaults: %+v", mc)
	}

	// Special characters in credentials and Unix socket path never inject DSN parameters.
	unixCfg := driver.Config{
		User:     "app?allowAllFiles=true",
		Password: "p@ss/word?tls=false&allowAllFiles=true",
		Host:     "/var/run/mysqld/mysqld.sock",
		Database: "sales",
		Params:   map[string]string{"timeout": "10s", "allowCleartextPasswords": "true"},
	}
	mcUnix, err := connConfig(unixCfg)
	if err != nil {
		t.Fatal(err)
	}
	if mcUnix.Net != "unix" || mcUnix.User != unixCfg.User || mcUnix.Passwd != unixCfg.Password || mcUnix.AllowAllFiles {
		t.Fatalf("unexpected unix config: %+v", mcUnix)
	}

	for _, key := range []string{"allowAllFiles", "allowallfiles", "AllowAllFiles", " allowAllFiles "} {
		bad := baseConfig()
		bad.Params = map[string]string{key: "true"}
		if _, err := connConfig(bad); err == nil {
			t.Errorf("accepted %q", key)
		}
	}

	// TLS cannot be downgraded by Params; nil cfg.TLS with tls=true infers host ServerName.
	tlsCfg := baseConfig()
	tlsCfg.Params = map[string]string{"tls": "false", "allowFallbackToPlaintext": "true"}
	tlsCfg.TLS = func(host string) (*tls.Config, error) { return &tls.Config{ServerName: host}, nil }
	mcTLS, err := connConfig(tlsCfg)
	if err != nil || mcTLS.TLS == nil || mcTLS.TLS.ServerName != "db.example.com" || mcTLS.AllowFallbackToPlaintext {
		t.Fatalf("TLS config = %+v, err = %v", mcTLS, err)
	}

	paramTLS := baseConfig()
	paramTLS.Params = map[string]string{"tls": "true"}
	mcParam, err := connConfig(paramTLS)
	if err != nil || mcParam.TLS == nil || mcParam.TLS.ServerName != "db.example.com" {
		t.Fatalf("Params tls=true = %+v, err = %v", mcParam, err)
	}

	// Open does not reach the server.
	dead := baseConfig()
	dead.Host, dead.Port = "127.0.0.1", "1"
	d, err := Open(context.Background(), dead)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Ping(context.Background()); err == nil {
		t.Error("Ping succeeded on closed port")
	}
	_ = d.Close()
}

func TestBuildExecuteCommandsDelimiterAndBounds(t *testing.T) {
	stmt := "DELIMITER //\nCREATE PROCEDURE p()\nBEGIN\n  SELECT 1;\nEND//\nDELIMITER ;\n"
	cmds := buildExecuteCommands(stmt)
	if len(cmds) != 1 || strings.Contains(cmds[0].Text, "DELIMITER") || !strings.Contains(cmds[0].Text, "CREATE PROCEDURE p()") {
		t.Fatalf("normalized delimiter cmds = %+v", cmds)
	}

	var many strings.Builder
	for range maxSplitCommands + 5 {
		many.WriteString("SELECT 1;\n")
	}
	if got := buildExecuteCommands(many.String()); len(got) != 1 || got[0].Text != many.String() {
		t.Fatalf("too-many-commands fallback = %d cmds", len(got))
	}

	large := "SELECT 1;\n" + strings.Repeat(" ", maxSheetCheckBytes)
	if got := buildExecuteCommands(large); len(got) != 1 || got[0].Text != large {
		t.Fatalf("large-sheet fallback = %d cmds", len(got))
	}
}

func TestOpenSyncDatabaseHandshake1049AndKillOnCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	var connNum atomic.Int32
	killed := make(chan string, 2)
	sleepStarted := make(chan struct{})

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
	readPacket := func(conn net.Conn) (byte, []byte, error) {
		var hdr [4]byte
		if _, err := io.ReadFull(conn, hdr[:]); err != nil {
			return 0, nil, err
		}
		n := int(hdr[0]) | int(hdr[1])<<8 | int(hdr[2])<<16
		buf := make([]byte, n)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return 0, nil, err
		}
		return hdr[3], buf, nil
	}
	greeting := func() []byte {
		var g []byte
		g = append(g, 0x0a)
		g = append(g, "11.8.8-MariaDB\x00"...)
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
			n := connNum.Add(1)
			go func(c net.Conn, idx int32) {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(5 * time.Second))
				if err := writePacket(c, 0, greeting()); err != nil {
					return
				}
				_, body, err := readPacket(c)
				if err != nil {
					return
				}
				if strings.Contains(string(body), "missing_db") {
					var errPkt []byte
					errPkt = append(errPkt, 0xff, 0x19, 0x04)
					errPkt = append(errPkt, "#42000Unknown database 'missing_db'"...)
					_ = writePacket(c, 2, errPkt)
					return
				}
				if err := writePacket(c, 2, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}); err != nil {
					return
				}
				for {
					_, pkt, err := readPacket(c)
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
						<-ctxDone(c)
						return
					case strings.HasPrefix(sql, "KILL QUERY "):
						killed <- sql
						_ = writePacket(c, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})
					default:
						_ = writePacket(c, 1, []byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00})
					}
				}
			}(conn, n)
		}
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	dMissing, err := Open(context.Background(), driver.Config{
		Host: host, Port: port, Database: "missing_db", User: "root",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer dMissing.Close()
	_, err = dMissing.SyncDatabase(context.Background())
	var mysqlErr *gomysql.MySQLError
	if !errors.Is(err, driver.ErrNotFound) || !errors.As(err, &mysqlErr) || mysqlErr.Number != 1049 {
		t.Fatalf("missing db SyncDatabase err = %v", err)
	}

	dLive, err := Open(context.Background(), driver.Config{
		Host: host, Port: port, User: "root", ApplicationName: "bytebase",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer dLive.Close()

	sess, err := dLive.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

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

func ctxDone(c net.Conn) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		var b [1]byte
		_, _ = c.Read(b[:])
	}()
	return ch
}
