package tape

import (
	"database/sql/driver"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// A tape is a text file of streams, each the requests one connection sent
// and the engine's answers, in order:
//
//	-- pin: default pingcap/tidb:v8.5.5 sha256:... v8.5.5
//
//	== TestX default#1
//	> exec
//	| CREATE TABLE t (a INT)
//	< ok 0 0
//	> query
//	| SHOW CREATE TABLE t
//	< result
//	< column "Table" "VARCHAR"
//	< column "Create Table" "VARCHAR"
//	< row
//	: bytes
//	| t
//	: bytes
//	| CREATE TABLE `t` (
//	|   `a` int DEFAULT NULL
//	| )
//
// Every line starts with a sigil, so SQL and values never need escaping:
// "| " carries one line of text, and "|q " a whole text as a Go quoted
// string where the raw form would not survive an editor (a carriage return,
// a control character, invalid UTF-8). A request is "> exec" or "> query"
// with its text and then one ": " value per argument, or "> begin",
// "> commit", "> rollback". An answer is "< ok", "< ok <rows affected>
// <last insert id>", "< error" with the codec's text, "< badconn", or one
// "< result" per result set with its columns and rows.

type exchange struct {
	req request
	ans answer
}

type request struct {
	kind string // exec, query, begin, commit, rollback
	text string
	args []value
}

type answer struct {
	kind     string // ok, error, badconn, rows
	affected string // exec only: decimal, or - when the driver could not say
	lastID   string
	err      string // the codec's text
	sets     []resultSet

	// raw is an error that is not the engine's answer, which only a
	// recording sees. It is never written: it poisons the recording.
	raw error
}

type resultSet struct {
	cols []column
	rows [][]value
}

type column struct{ name, typ string }

// A value is a driver.Value as text: null, int64, float64, bool, time,
// bytes or string.
type value struct {
	kind, text string
}

type stream struct {
	name      string
	exchanges []exchange
}

type tapeFile struct {
	pins    map[string]string // role -> enginetest.Pin.String()
	streams []*stream
}

func (f *tapeFile) stream(name string) *stream {
	for _, s := range f.streams {
		if s.name == name {
			return s
		}
	}
	return nil
}

const (
	pinPrefix    = "-- pin: "
	streamPrefix = "== "
)

func fromDriver(v driver.Value) (value, error) {
	switch v := v.(type) {
	case nil:
		return value{kind: "null"}, nil
	case int64:
		return value{"int64", strconv.FormatInt(v, 10)}, nil
	case float64:
		return value{"float64", strconv.FormatFloat(v, 'g', -1, 64)}, nil
	case bool:
		return value{"bool", strconv.FormatBool(v)}, nil
	case time.Time:
		return value{"time", v.Format(time.RFC3339Nano)}, nil
	case []byte:
		return value{"bytes", string(v)}, nil
	case string:
		return value{"string", v}, nil
	}
	return value{}, fmt.Errorf("a %T value, which a tape does not hold", v)
}

func (v value) driver() (driver.Value, error) {
	switch v.kind {
	case "null":
		return nil, nil
	case "int64":
		return strconv.ParseInt(v.text, 10, 64)
	case "float64":
		return strconv.ParseFloat(v.text, 64)
	case "bool":
		return strconv.ParseBool(v.text)
	case "time":
		return time.Parse(time.RFC3339Nano, v.text)
	case "bytes":
		return []byte(v.text), nil
	case "string":
		return v.text, nil
	}
	return nil, fmt.Errorf("a value of kind %q", v.kind)
}

// encode renders the file. Pins come first in role order, then the streams
// in the order given.
func (f *tapeFile) encode() string {
	var b strings.Builder
	roles := make([]string, 0, len(f.pins))
	for r := range f.pins {
		roles = append(roles, r)
	}
	sort.Strings(roles)
	for _, r := range roles {
		b.WriteString(pinPrefix + f.pins[r] + "\n")
	}
	for _, s := range f.streams {
		b.WriteString("\n" + streamPrefix + s.name + "\n")
		for _, e := range s.exchanges {
			writeRequest(&b, e.req)
			writeAnswer(&b, e.ans)
		}
	}
	return b.String()
}

func writeRequest(b *strings.Builder, r request) {
	b.WriteString("> " + r.kind + "\n")
	if r.kind == "exec" || r.kind == "query" {
		writeText(b, r.text)
		for _, a := range r.args {
			writeValue(b, a)
		}
	}
}

func writeAnswer(b *strings.Builder, a answer) {
	switch a.kind {
	case "ok":
		if a.affected == "" {
			b.WriteString("< ok\n")
		} else {
			b.WriteString("< ok " + a.affected + " " + a.lastID + "\n")
		}
	case "error":
		b.WriteString("< error\n")
		writeText(b, a.err)
	case "badconn":
		b.WriteString("< badconn\n")
	case "rows":
		for _, s := range a.sets {
			b.WriteString("< result\n")
			for _, c := range s.cols {
				b.WriteString("< column " + strconv.Quote(c.name) + " " + strconv.Quote(c.typ) + "\n")
			}
			for _, row := range s.rows {
				b.WriteString("< row\n")
				for _, v := range row {
					writeValue(b, v)
				}
			}
		}
	}
}

func writeValue(b *strings.Builder, v value) {
	switch v.kind {
	case "bytes", "string":
		b.WriteString(": " + v.kind + "\n")
		writeText(b, v.text)
	case "null":
		b.WriteString(": null\n")
	default:
		b.WriteString(": " + v.kind + " " + v.text + "\n")
	}
}

func writeText(b *strings.Builder, s string) {
	if !plain(s) {
		b.WriteString("|q " + strconv.Quote(s) + "\n")
		return
	}
	for _, line := range strings.Split(s, "\n") {
		if line == "" {
			b.WriteString("|\n")
		} else {
			b.WriteString("| " + line + "\n")
		}
	}
}

// plain reports whether s can be written as raw lines and read back
// unchanged: valid UTF-8 with no control character but tab and newline.
func plain(s string) bool {
	if !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if r < 0x20 && r != '\t' && r != '\n' || r == 0x7f {
			return false
		}
	}
	return true
}

// render is one answer as tape text, for comparing two and for messages.
func (a answer) render() string {
	var b strings.Builder
	writeAnswer(&b, a)
	return b.String()
}

func (r request) render() string {
	var b strings.Builder
	writeRequest(&b, r)
	return b.String()
}

// describe names a request in one line for a message.
func (r request) describe() string {
	if r.kind != "exec" && r.kind != "query" {
		return r.kind
	}
	text := r.text
	if len(text) > 200 {
		text = text[:200] + "..."
	}
	s := r.kind + " " + strconv.Quote(text)
	if len(r.args) > 0 {
		s += fmt.Sprintf(" with %d arguments", len(r.args))
	}
	return s
}

// parse reads a tape. name is for errors.
func parse(name, data string) (*tapeFile, error) {
	p := &parser{name: name, lines: strings.Split(data, "\n")}
	if n := len(p.lines); n > 0 && p.lines[n-1] == "" {
		p.lines = p.lines[:n-1]
	}
	f := &tapeFile{pins: map[string]string{}}
	var cur *stream
	for p.i < len(p.lines) {
		line := p.lines[p.i]
		switch {
		case line == "":
			p.i++
		case strings.HasPrefix(line, pinPrefix):
			pin := strings.TrimPrefix(line, pinPrefix)
			role, _, _ := strings.Cut(pin, " ")
			if cur != nil {
				return nil, p.errorf("a pin after the first stream")
			}
			f.pins[role] = pin
			p.i++
		case strings.HasPrefix(line, streamPrefix):
			name := strings.TrimPrefix(line, streamPrefix)
			if f.stream(name) != nil {
				return nil, p.errorf("stream %q appears twice", name)
			}
			cur = &stream{name: name}
			f.streams = append(f.streams, cur)
			p.i++
		case strings.HasPrefix(line, "> "):
			if cur == nil {
				return nil, p.errorf("a request before any stream")
			}
			e, err := p.exchange()
			if err != nil {
				return nil, err
			}
			cur.exchanges = append(cur.exchanges, e)
		default:
			return nil, p.errorf("unexpected line %q", line)
		}
	}
	return f, nil
}

type parser struct {
	name  string
	lines []string
	i     int
}

func (p *parser) errorf(format string, args ...any) error {
	return fmt.Errorf("%s:%d: %s", p.name, p.i+1, fmt.Sprintf(format, args...))
}

func (p *parser) peek() string {
	if p.i < len(p.lines) {
		return p.lines[p.i]
	}
	return ""
}

func (p *parser) exchange() (exchange, error) {
	var e exchange
	e.req.kind = strings.TrimPrefix(p.peek(), "> ")
	p.i++
	switch e.req.kind {
	case "exec", "query":
		text, err := p.text()
		if err != nil {
			return e, err
		}
		e.req.text = text
		for strings.HasPrefix(p.peek(), ": ") {
			v, err := p.value()
			if err != nil {
				return e, err
			}
			e.req.args = append(e.req.args, v)
		}
	case "begin", "commit", "rollback":
	default:
		return e, p.errorf("unknown request %q", e.req.kind)
	}
	a, err := p.answer()
	e.ans = a
	return e, err
}

func (p *parser) answer() (answer, error) {
	line := p.peek()
	switch {
	case line == "< ok":
		p.i++
		return answer{kind: "ok"}, nil
	case strings.HasPrefix(line, "< ok "):
		f := strings.Fields(strings.TrimPrefix(line, "< ok "))
		if len(f) != 2 {
			return answer{}, p.errorf("want < ok <rows affected> <last insert id>")
		}
		p.i++
		return answer{kind: "ok", affected: f[0], lastID: f[1]}, nil
	case line == "< error":
		p.i++
		text, err := p.text()
		return answer{kind: "error", err: text}, err
	case line == "< badconn":
		p.i++
		return answer{kind: "badconn"}, nil
	case line == "< result":
		a := answer{kind: "rows"}
		for p.peek() == "< result" {
			p.i++
			var s resultSet
			for strings.HasPrefix(p.peek(), "< column ") {
				c, err := p.column()
				if err != nil {
					return a, err
				}
				s.cols = append(s.cols, c)
			}
			for p.peek() == "< row" {
				p.i++
				row := make([]value, len(s.cols))
				for j := range row {
					v, err := p.value()
					if err != nil {
						return a, err
					}
					row[j] = v
				}
				s.rows = append(s.rows, row)
			}
			a.sets = append(a.sets, s)
		}
		return a, nil
	}
	return answer{}, p.errorf("want an answer (< ok, < error, < badconn, < result), got %q", line)
}

func (p *parser) column() (column, error) {
	rest := strings.TrimPrefix(p.peek(), "< column ")
	name, err := strconv.QuotedPrefix(rest)
	if err != nil {
		return column{}, p.errorf("a column's name: %v", err)
	}
	typ := strings.TrimPrefix(rest[len(name):], " ")
	n, _ := strconv.Unquote(name)
	t, err := strconv.Unquote(typ)
	if err != nil {
		return column{}, p.errorf("a column's type: %v", err)
	}
	p.i++
	return column{n, t}, nil
}

func (p *parser) value() (value, error) {
	line := p.peek()
	if !strings.HasPrefix(line, ": ") {
		return value{}, p.errorf("want a value, got %q", line)
	}
	kind, text, _ := strings.Cut(strings.TrimPrefix(line, ": "), " ")
	p.i++
	v := value{kind: kind, text: text}
	switch kind {
	case "bytes", "string":
		t, err := p.text()
		v.text = t
		return v, err
	case "null":
		return v, nil
	case "int64", "float64", "bool", "time":
		if _, err := v.driver(); err != nil {
			return v, p.errorf("a %s value: %v", kind, err)
		}
		return v, nil
	}
	return v, p.errorf("unknown value kind %q", kind)
}

// text reads one text: a single "|q " line, or one or more "|" lines.
func (p *parser) text() (string, error) {
	line := p.peek()
	if strings.HasPrefix(line, "|q ") {
		p.i++
		s, err := strconv.Unquote(strings.TrimPrefix(line, "|q "))
		if err != nil {
			return "", p.errorf("a quoted text: %v", err)
		}
		return s, nil
	}
	var lines []string
	for {
		line := p.peek()
		if line == "|" {
			lines = append(lines, "")
		} else if strings.HasPrefix(line, "| ") {
			lines = append(lines, line[2:])
		} else {
			break
		}
		p.i++
	}
	if lines == nil {
		return "", p.errorf("want text (| lines), got %q", p.peek())
	}
	return strings.Join(lines, "\n"), nil
}
