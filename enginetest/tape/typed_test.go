package tape

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest"
)

// testDecimal stands in for a client's struct type, which a tape writes
// through its Text.
type testDecimal struct {
	coef *big.Int
	exp  int32
}

// testIP stands in for a named slice such as net.IP, written by its kind.
type testIP []byte

var testTypes = []Type{
	{
		Type: reflect.TypeOf(testDecimal{}),
		Text: func(v any) (string, error) {
			d := v.(testDecimal)
			return fmt.Sprintf("%se%d", d.coef, d.exp), nil
		},
		Value: func(text string) (any, error) {
			coef, exp, ok := strings.Cut(text, "e")
			n, okN := new(big.Int).SetString(coef, 10)
			var e int32
			_, err := fmt.Sscan(exp, &e)
			if !ok || !okN || err != nil {
				return nil, fmt.Errorf("a decimal %q", text)
			}
			return testDecimal{n, e}, nil
		},
	},
	{Type: reflect.TypeOf(testIP(nil))},
}

func ptr[T any](v T) *T { return &v }

// typedCell is one value a typed client answers, the column's scan type, and
// the literal the tape holds for it.
type typedCell struct {
	scan    reflect.Type
	value   driver.Value
	literal string
}

func typedCells() []typedCell {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	at := time.Date(2026, 3, 4, 5, 6, 7, 891000000, time.UTC)
	return []typedCell{
		{reflect.TypeOf([]int32(nil)), []int32{1, -2, 3}, "[1, -2, 3]"},
		{reflect.TypeOf([]int32(nil)), []int32{}, "[]"},
		{reflect.TypeOf([]int32(nil)), []int32(nil), "nil"},
		{reflect.TypeOf([]*int32(nil)), []*int32{ptr(int32(1)), nil}, "[&1, nil]"},
		{reflect.TypeOf([][]string(nil)), [][]string{{"a", "b,c"}, {}}, `[["a", "b,c"], []]`},
		{reflect.TypeOf([][]uint8(nil)), [][]uint8{{0, 255}}, "[[0, 255]]"},
		{reflect.TypeOf([]float64(nil)), []float64{1.5, math.Inf(1), math.Inf(-1), math.Copysign(0, -1)}, "[1.5, +Inf, -Inf, -0]"},
		{reflect.TypeOf(map[string]uint64(nil)), map[string]uint64{"b": 2, "a": 1}, `{"a": 1, "b": 2}`},
		{reflect.TypeOf(map[float64]string(nil)), map[float64]string{2.5: "x"}, `{2.5: "x"}`},
		{reflect.TypeOf(map[string]*int64(nil)), map[string]*int64{"n": nil}, `{"n": nil}`},
		{reflect.TypeOf([]any(nil)), []any{int32(1), "a", nil, (*int8)(nil), []string{"x"}}, `[int32(1), string("a"), nil, *int8(nil), []string(["x"])]`},
		{reflect.TypeOf(map[string]any(nil)), map[string]any{"x": float64(1.5), "y": map[string]any{}}, `{"x": float64(1.5), "y": map[string]interface {}({})}`},
		{reflect.TypeOf([2]float64{}), [2]float64{1, 2}, "[1, 2]"},
		{reflect.TypeOf(testDecimal{}), testDecimal{big.NewInt(150), -2}, `"150e-2"`},
		{reflect.TypeOf([]testDecimal(nil)), []testDecimal{{big.NewInt(-5), 0}}, `["-5e0"]`},
		{reflect.TypeOf(testIP(nil)), testIP{127, 0, 0, 1}, "[127, 0, 0, 1]"},
		{reflect.TypeOf([]time.Time(nil)), []time.Time{at, at.In(tokyo), at.In(time.FixedZone("", 3600))}, `["2026-03-04T05:06:07.891Z UTC", "2026-03-04T14:06:07.891+09:00 Asia/Tokyo", "2026-03-04T06:06:07.891+01:00 "]`},
		{reflect.TypeOf([]time.Duration(nil)), []time.Duration{90 * time.Second}, "[90000000000]"},
	}
}

// typedRows answers typedCells as one row, with each column's scan type.
type typedRows struct {
	cells []typedCell
	sent  bool
}

func (r *typedRows) Columns() []string {
	names := make([]string, len(r.cells))
	for i := range names {
		names[i] = fmt.Sprintf("c%d", i)
	}
	return names
}
func (r *typedRows) ColumnTypeDatabaseTypeName(int) string { return "T" }
func (r *typedRows) ColumnTypeScanType(i int) reflect.Type { return r.cells[i].scan }
func (r *typedRows) Close() error                          { return nil }
func (r *typedRows) Next(dest []driver.Value) error {
	if r.sent {
		return io.EOF
	}
	for i, c := range r.cells {
		dest[i] = c.value
	}
	r.sent = true
	return nil
}

type typedConnector struct{ cells []typedCell }

func (c typedConnector) Connect(context.Context) (driver.Conn, error) { return typedConn(c), nil }
func (c typedConnector) Driver() driver.Driver                        { return nil }

type typedConn typedConnector

func (c typedConn) Prepare(string) (driver.Stmt, error) { return nil, errPrepared }
func (c typedConn) Close() error                        { return nil }
func (c typedConn) Begin() (driver.Tx, error)           { return fakeTx{}, nil }
func (c typedConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &typedRows{cells: c.cells}, nil
}

func typedServer(t *testing.T, cells []typedCell, types []Type) *Server {
	srv, _ := harness(t)
	srv.Live = func(testing.TB, enginetest.Pin) driver.Connector { return typedConnector{cells} }
	srv.Types = types
	return srv
}

// readTyped reads the one row typedRows answers, each cell into its column's
// scan type, and returns the scan types and values' types. The values are
// appended to got.
func readTyped(srv *Server, got *[]any) func(t testing.TB) []string {
	return func(t testing.TB) []string {
		rows, err := Open(t, srv, "default").Query("SELECT typed")
		if err != nil {
			return []string{err.Error()}
		}
		defer rows.Close()
		types, err := rows.ColumnTypes()
		if err != nil {
			return []string{err.Error()}
		}
		var out []string
		for rows.Next() {
			dest := make([]any, len(types))
			for i, c := range types {
				dest[i] = reflect.New(c.ScanType()).Interface()
			}
			if err := rows.Scan(dest...); err != nil {
				return []string{err.Error()}
			}
			for i, d := range dest {
				v := reflect.ValueOf(d).Elem().Interface()
				*got = append(*got, v)
				out = append(out, fmt.Sprintf("%s %T", types[i].ScanType(), v))
			}
		}
		if err := rows.Err(); err != nil {
			return []string{err.Error()}
		}
		return out
	}
}

// A client whose rows are typed answers slices, maps and its own types. The
// tape writes each by its type, and a replay hands database/sql the same type
// and value, and the same scan type for the column.
func TestTypedValuesReplay(t *testing.T) {
	cells := typedCells()
	srv := typedServer(t, cells, testTypes)

	var live, replayedValues []any
	record(t)
	recorded, errs := run(t, "TestFake", readTyped(srv, &live))
	if len(errs) > 0 || len(recorded) != len(cells) {
		t.Fatalf("a recording read %q (%v)", recorded, errs)
	}
	b, err := os.ReadFile(filepath.Join(dir, "TestFake.tape"))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cells {
		line := fmt.Sprintf(": go %s\n| %s\n", reflect.TypeOf(c.value), c.literal)
		if !strings.Contains(string(b), line) {
			t.Errorf("the tape lacks %q:\n%s", line, b)
		}
		column := fmt.Sprintf(" %q\n", c.scan.String())
		if !strings.Contains(string(b), column) {
			t.Errorf("the tape lacks the scan type %q", c.scan)
		}
	}

	replay(t)
	replayed, errs := run(t, "TestFake", readTyped(srv, &replayedValues))
	if len(errs) > 0 {
		t.Fatalf("a replay failed: %v", errs)
	}
	if !reflect.DeepEqual(replayed, recorded) {
		t.Errorf("a replay read\n%s\nwant\n%s", strings.Join(replayed, "\n"), strings.Join(recorded, "\n"))
	}
	for i, c := range cells {
		if !reflect.DeepEqual(live[i], c.value) || !reflect.DeepEqual(replayedValues[i], c.value) {
			t.Errorf("column %d: recorded %#v and replayed %#v, want %#v", i, live[i], replayedValues[i], c.value)
		}
	}
}

// Each value comes back as the same Go value: the same type, the same nil and
// empty slices and maps, the same pointers, and a time in the same location.
func TestTypedValuesRoundTrip(t *testing.T) {
	ts, err := newTypes(testTypes)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range typedCells() {
		v, err := ts.fromDriver(c.value)
		if err != nil {
			t.Errorf("%T: %v", c.value, err)
			continue
		}
		back, err := ts.driverValue(v)
		if err != nil {
			t.Errorf("%T %s: %v", c.value, v.text, err)
			continue
		}
		if !reflect.DeepEqual(back, c.value) {
			t.Errorf("%T: %#v came back as %#v", c.value, c.value, back)
		}
		if times, ok := c.value.([]time.Time); ok {
			for i, tm := range back.([]time.Time) {
				if got, want := tm.Location().String(), times[i].Location().String(); got != want {
					t.Errorf("time %d is in %q, want %q", i, got, want)
				}
			}
		}
	}

	// NaN is never equal to itself, so DeepEqual cannot compare it.
	v, err := ts.fromDriver([]float32{float32(math.NaN())})
	if err != nil || v.text != "[NaN]" {
		t.Fatalf("NaN: %+v, %v", v, err)
	}
	back, err := ts.driverValue(v)
	if f, ok := back.([]float32); err != nil || !ok || len(f) != 1 || !math.IsNaN(float64(f[0])) {
		t.Errorf("NaN came back as %#v, %v", back, err)
	}
}

// A server that states no Types names only the scan types a tape always has,
// so a recording of another engine's tapes is unchanged, and its values of
// composite types are still written.
func TestScanTypesWithoutTypes(t *testing.T) {
	cells := []typedCell{
		{reflect.TypeOf(int64(0)), int64(1), ""},
		{reflect.TypeOf([]int32(nil)), []int32{1}, "[1]"},
	}
	srv := typedServer(t, cells, nil)
	record(t)
	if _, errs := run(t, "TestFake", readTyped(srv, new([]any))); len(errs) > 0 {
		t.Fatalf("a recording failed: %v", errs)
	}
	b, err := os.ReadFile(filepath.Join(dir, "TestFake.tape"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range []string{`< column "c0" "T" "int64"`, `< column "c1" "T"` + "\n", ": go []int32\n| [1]\n"} {
		if !strings.Contains(string(b), line) {
			t.Errorf("the tape lacks %q:\n%s", line, b)
		}
	}
}

type unknownStruct struct{ A int }

// A value of a type the tape cannot spell fails the recording, which keeps
// the committed tape, rather than writing what a replay would read otherwise.
func TestAnUnknownTypeFailsTheRecording(t *testing.T) {
	for _, v := range []driver.Value{
		unknownStruct{1},
		[]unknownStruct{{1}},
		[]any{unknownStruct{1}},
		map[string]chan int{"c": nil},
		[]complex128{1},
	} {
		srv := typedServer(t, []typedCell{{reflect.TypeOf(v), v, ""}}, testTypes)
		record(t)
		if _, errs := run(t, "TestFake", readTyped(srv, new([]any))); len(errs) == 0 || !strings.Contains(strings.Join(errs, "\n"), "which a tape does not hold") {
			t.Errorf("%T recorded: %v", v, errs)
		}
		if _, err := os.Stat(filepath.Join(dir, "TestFake.tape")); !os.IsNotExist(err) {
			t.Errorf("%T: a failed recording wrote a tape (%v)", v, err)
		}
	}
}

func TestTypesRejects(t *testing.T) {
	// A type declared in a function is spelled as one declared outside it.
	type testIP []byte
	text := func(any) (string, error) { return "", nil }
	value := func(string) (any, error) { return nil, nil }
	for _, c := range []struct {
		name  string
		types []Type
		want  string
	}{
		{"no type", []Type{{}}, "with no Type"},
		{"a struct without Text", []Type{{Type: reflect.TypeOf(unknownStruct{})}}, "needs Text and Value"},
		{"Text without Value", []Type{{Type: reflect.TypeOf(testIP(nil)), Text: text}}, "one of Text and Value"},
		{"Value without Text", []Type{{Type: reflect.TypeOf(testIP(nil)), Value: value}}, "one of Text and Value"},
		{"a spelling twice", []Type{testTypes[1], {Type: reflect.TypeOf(testIP(nil))}}, "two types are spelled"},
	} {
		if _, err := newTypes(c.types); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want %q", c.name, err, c.want)
		}
	}
}

func TestTypedValueRejects(t *testing.T) {
	ts, err := newTypes(testTypes)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ typ, text, want string }{
		{"[]int8", "[1, 300]", "out of range"},
		{"[]int8", "[1 2]", `want ","`},
		{"[]int8", "[1, 2] x", "after the value"},
		{"[2]int8", "[1]", "1 elements of a [2]int8"},
		{"[1]int8", "[1, 2]", "more than the 1 elements"},
		{"map[string]int8", `{"a" 1}`, `want ":"`},
		{"[]interface {}", "[int8 1]", "want a type and ("},
		{"[]interface {}", "[chan(1)]", "does not know"},
		{"*int8", "1", `want "&"`},
		{"[]string", "[a]", "want a quoted string"},
		{"tape.testDecimal", `"x"`, "a decimal"},
		{"time.Time", `"2026-03-04T05:06:07Z"`, "without its location"},
		{"map[[]int8]int8", "{}", "a map type keyed by"},
		{"[x]int8", "[]", "an array type with length"},
		{"[]int8 x", "[]", "after the type"},
	} {
		_, err := ts.driverValue(value{kind: "go", typ: c.typ, text: c.text})
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s %s: got %v, want %q", c.typ, c.text, err, c.want)
		}
	}
	if _, err := parse("x.tape", "== s\n> query\n| x\n: go []int8\n| [1]\n< ok\n"); err == nil || !strings.Contains(err.Error(), "only answers hold one") {
		t.Errorf("a typed argument parsed: %v", err)
	}
	if _, err := parse("x.tape", "== s\n> query\n| x\n< result\n< column \"a\" \"T\"\n< row\n: go\n| [1]\n"); err == nil || !strings.Contains(err.Error(), "with no type") {
		t.Errorf("a typed value with no type parsed: %v", err)
	}
}
