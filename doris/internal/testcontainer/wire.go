package testcontainer

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// wire is a client of the MySQL protocol Doris speaks that reads what
// go-sql-driver drops: the full column definition of a result set, with the
// org_table and org_name Doris fills for a plain column reference. It logs in
// as root, whose password the image leaves empty, so the answer to any
// authentication method is empty. Spans uses it in a recording only.
type wire struct {
	conn net.Conn
	r    *bufio.Reader
	seq  byte
}

// Capability flags of the handshake.
const (
	clientLongPassword     = 0x00000001
	clientConnectWithDB    = 0x00000008
	clientProtocol41       = 0x00000200
	clientTransactions     = 0x00002000
	clientSecureConnection = 0x00008000
	clientMultiResults     = 0x00020000
	clientPluginAuth       = 0x00080000
)

// wireError is the server's ERR packet.
type wireError struct {
	Number   uint16
	SQLState string
	Message  string
}

func (e *wireError) Error() string { return fmt.Sprintf("%d %s %s", e.Number, e.SQLState, e.Message) }

// wireColumn is one column definition (Protocol::ColumnDefinition41).
type wireColumn struct {
	Schema   string
	Table    string
	OrgTable string
	Name     string
	OrgName  string
}

// wireTimeout bounds one exchange with the live server.
const wireTimeout = 2 * time.Minute

func dialWire(addr, user string) (*wire, error) {
	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		return nil, err
	}
	w := &wire{conn: conn, r: bufio.NewReader(conn)}
	if err := w.login(user); err != nil {
		conn.Close()
		return nil, fmt.Errorf("log in as %s: %w", user, err)
	}
	return w, nil
}

func (w *wire) Close() error { return w.conn.Close() }

func (w *wire) read() ([]byte, error) {
	_ = w.conn.SetDeadline(time.Now().Add(wireTimeout))
	var head [4]byte
	if _, err := io.ReadFull(w.r, head[:]); err != nil {
		return nil, err
	}
	n := int(head[0]) | int(head[1])<<8 | int(head[2])<<16
	w.seq = head[3] + 1
	b := make([]byte, n)
	if _, err := io.ReadFull(w.r, b); err != nil {
		return nil, err
	}
	if n == 0xffffff {
		return nil, errors.New("a packet of 16 MiB or more")
	}
	return b, nil
}

func (w *wire) write(b []byte) error {
	_ = w.conn.SetDeadline(time.Now().Add(wireTimeout))
	head := []byte{byte(len(b)), byte(len(b) >> 8), byte(len(b) >> 16), w.seq}
	w.seq++
	_, err := w.conn.Write(append(head, b...))
	return err
}

// command starts a command: its sequence restarts at 0.
func (w *wire) command(code byte, body string) error {
	w.seq = 0
	return w.write(append([]byte{code}, body...))
}

func (w *wire) login(user string) error {
	hello, err := w.read()
	if err != nil {
		return err
	}
	if len(hello) > 0 && hello[0] == 0xff {
		return parseErr(hello)
	}
	caps := uint32(clientLongPassword | clientProtocol41 | clientTransactions | clientSecureConnection | clientMultiResults | clientPluginAuth)
	var b []byte
	b = binary.LittleEndian.AppendUint32(b, caps)
	b = binary.LittleEndian.AppendUint32(b, 1<<24)
	b = append(b, 45) // utf8mb4_general_ci
	b = append(b, make([]byte, 23)...)
	b = append(b, user...)
	b = append(b, 0)
	b = append(b, 0) // an empty password's answer
	b = append(b, "mysql_native_password"...)
	b = append(b, 0)
	if err := w.write(b); err != nil {
		return err
	}
	for {
		p, err := w.read()
		if err != nil {
			return err
		}
		switch {
		case len(p) == 0:
			return errors.New("an empty packet while logging in")
		case p[0] == 0x00:
			return nil
		case p[0] == 0xff:
			return parseErr(p)
		case p[0] == 0xfe:
			// A switch to another method; an empty password's answer is
			// empty in every one.
			if err := w.write(nil); err != nil {
				return err
			}
		default:
			return fmt.Errorf("an unexpected packet 0x%02x while logging in", p[0])
		}
	}
}

// ok reads the answer of a command that answers OK or ERR.
func (w *wire) ok() error {
	p, err := w.read()
	if err != nil {
		return err
	}
	switch {
	case len(p) > 0 && p[0] == 0xff:
		return parseErr(p)
	case len(p) > 0 && p[0] == 0x00:
		return nil
	}
	return fmt.Errorf("an unexpected packet 0x%02x", p[0])
}

// Use selects database (COM_INIT_DB).
func (w *wire) Use(database string) error {
	if err := w.command(0x02, database); err != nil {
		return err
	}
	return w.ok()
}

// Exec runs a statement that returns no rows (COM_QUERY).
func (w *wire) Exec(query string) error {
	if err := w.command(0x03, query); err != nil {
		return err
	}
	return w.ok()
}

// Query runs query (COM_QUERY) and returns its result set's column
// definitions and its rows as text; NULL is "".
func (w *wire) Query(query string) ([]wireColumn, [][]string, error) {
	if err := w.command(0x03, query); err != nil {
		return nil, nil, err
	}
	p, err := w.read()
	if err != nil {
		return nil, nil, err
	}
	if p[0] == 0xff {
		return nil, nil, parseErr(p)
	}
	if p[0] == 0x00 {
		return nil, nil, nil
	}
	n, _ := lenencInt(p)
	cols, err := w.columns(int(n))
	if err != nil {
		return nil, nil, err
	}
	var rows [][]string
	for {
		p, err := w.read()
		if err != nil {
			return nil, nil, err
		}
		if p[0] == 0xff {
			return nil, nil, parseErr(p)
		}
		if p[0] == 0xfe && len(p) < 9 {
			if binary.LittleEndian.Uint16(p[3:5])&0x08 != 0 {
				return nil, nil, errors.New("more than one result set")
			}
			return cols, rows, nil
		}
		row := make([]string, n)
		rest := p
		for i := range row {
			if rest[0] == 0xfb {
				rest = rest[1:]
				continue
			}
			row[i], rest = lenencString(rest)
		}
		rows = append(rows, row)
	}
}

// columns reads n column definitions and the EOF after them, if n > 0.
func (w *wire) columns(n int) ([]wireColumn, error) {
	if n == 0 {
		return nil, nil
	}
	out := make([]wireColumn, n)
	for i := range out {
		p, err := w.read()
		if err != nil {
			return nil, err
		}
		var f [6]string
		rest := p
		for k := range f {
			f[k], rest = lenencString(rest)
		}
		// catalog, schema, table, org_table, name, org_name
		out[i] = wireColumn{Schema: f[1], Table: f[2], OrgTable: f[3], Name: f[4], OrgName: f[5]}
	}
	p, err := w.read()
	if err != nil {
		return nil, err
	}
	if len(p) == 0 || p[0] != 0xfe {
		return nil, errors.New("no EOF after the column definitions")
	}
	return out, nil
}

func parseErr(p []byte) error {
	e := &wireError{Number: binary.LittleEndian.Uint16(p[1:3]), SQLState: "-"}
	rest := p[3:]
	if len(rest) >= 6 && rest[0] == '#' {
		e.SQLState, rest = string(rest[1:6]), rest[6:]
	}
	e.Message = string(rest)
	return e
}

func lenencInt(b []byte) (uint64, []byte) {
	switch b[0] {
	case 0xfc:
		return uint64(binary.LittleEndian.Uint16(b[1:3])), b[3:]
	case 0xfd:
		return uint64(b[1]) | uint64(b[2])<<8 | uint64(b[3])<<16, b[4:]
	case 0xfe:
		return binary.LittleEndian.Uint64(b[1:9]), b[9:]
	}
	return uint64(b[0]), b[1:]
}

func lenencString(b []byte) (string, []byte) {
	n, rest := lenencInt(b)
	return string(rest[:n]), rest[n:]
}
