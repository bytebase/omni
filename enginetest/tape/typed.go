package tape

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// A Type is a named Go type a driver's answers hold that a tape does not know
// on its own: a client's own decimal or UUID type, or a named slice such as
// net.IP. A tape knows Go's booleans, integers, floats and strings,
// time.Time, time.Duration, and the pointers, slices, arrays, maps and
// interfaces built from known types, which is what a client whose rows are
// typed (ClickHouse's) answers for an Array, a Tuple or a Map.
//
// A tape writes a value of a known type by its type's spelling (reflect's
// Type.String) and a literal that the type reads back to the same value: see
// typedValue. A Type that is a struct, whose fields a tape cannot reach, needs
// Text and Value; any other kind may leave both nil and is written by its kind.
type Type struct {
	Type reflect.Type
	// Text writes a value of Type, and Value reads it back. Value must return
	// a value of Type that Text writes as the same text.
	Text  func(v any) (string, error)
	Value func(text string) (any, error)
}

// types is what a session's tape can write: Go's own types, those above, and
// the Server's.
type types struct {
	named map[string]reflect.Type // by spelling
	text  map[reflect.Type]Type   // the types written by Text
	// composite is set for a Server that states Types: its scan types are
	// spelled like its values' types. A Server that states none names only
	// the scan types knownScanTypes holds, as every tape before Types did.
	composite bool
}

var timeType = Type{
	Type: reflect.TypeOf(time.Time{}),
	// The location's name is written beside the instant: RFC 3339 keeps only
	// its offset, and a client states its column's zone by name.
	Text: func(v any) (string, error) {
		t := v.(time.Time)
		return t.Format(time.RFC3339Nano) + " " + t.Location().String(), nil
	},
	Value: func(text string) (any, error) {
		instant, name, ok := strings.Cut(text, " ")
		if !ok {
			return nil, fmt.Errorf("a time without its location: %q", text)
		}
		t, err := time.Parse(time.RFC3339Nano, instant)
		if err != nil {
			return nil, err
		}
		_, offset := t.Zone()
		if name == "UTC" && offset == 0 {
			return t.In(time.UTC), nil
		}
		// A zone this machine knows by name, at the same offset, is the one
		// the client stated; any other is rebuilt as the name and offset the
		// recording saw.
		if name != "" {
			if loc, err := time.LoadLocation(name); err == nil {
				if _, at := t.In(loc).Zone(); at == offset {
					return t.In(loc), nil
				}
			}
		}
		return t.In(time.FixedZone(name, offset)), nil
	},
}

func newTypes(extra []Type) (*types, error) {
	ts := &types{named: map[string]reflect.Type{}, text: map[reflect.Type]Type{}, composite: len(extra) > 0}
	for _, t := range []reflect.Type{
		reflect.TypeOf(false), reflect.TypeOf(""),
		reflect.TypeOf(int(0)), reflect.TypeOf(int8(0)), reflect.TypeOf(int16(0)), reflect.TypeOf(int32(0)), reflect.TypeOf(int64(0)),
		reflect.TypeOf(uint(0)), reflect.TypeOf(uint8(0)), reflect.TypeOf(uint16(0)), reflect.TypeOf(uint32(0)), reflect.TypeOf(uint64(0)),
		reflect.TypeOf(float32(0)), reflect.TypeOf(float64(0)),
		reflect.TypeOf(time.Duration(0)),
	} {
		ts.named[t.String()] = t
	}
	for _, t := range append([]Type{timeType}, extra...) {
		if t.Type == nil {
			return nil, fmt.Errorf("tape: a Type with no Type")
		}
		name := t.Type.String()
		if old, ok := ts.named[name]; ok && old != t.Type {
			return nil, fmt.Errorf("tape: two types are spelled %s", name)
		}
		if (t.Text == nil) != (t.Value == nil) {
			return nil, fmt.Errorf("tape: type %s states one of Text and Value", name)
		}
		if t.Text == nil && t.Type.Kind() == reflect.Struct {
			return nil, fmt.Errorf("tape: type %s is a struct and needs Text and Value", name)
		}
		ts.named[name] = t.Type
		if t.Text != nil {
			ts.text[t.Type] = t
		}
	}
	return ts, nil
}

// spell is t's spelling, or an error when a tape cannot read t back from it.
func (ts *types) spell(t reflect.Type) (string, error) {
	name := t.String()
	back, err := ts.parse(name)
	if err != nil || back != t {
		return "", fmt.Errorf("a %s value, which a tape does not hold", name)
	}
	return name, nil
}

// parse reads a type's spelling.
func (ts *types) parse(spelling string) (reflect.Type, error) {
	t, rest, err := ts.parseType(spelling)
	if err != nil {
		return nil, err
	}
	if rest != "" {
		return nil, fmt.Errorf("type %q: %q after the type", spelling, rest)
	}
	return t, nil
}

func (ts *types) parseType(s string) (reflect.Type, string, error) {
	switch {
	case strings.HasPrefix(s, "*"):
		elem, rest, err := ts.parseType(s[1:])
		if err != nil {
			return nil, "", err
		}
		return reflect.PointerTo(elem), rest, nil
	case strings.HasPrefix(s, "[]"):
		elem, rest, err := ts.parseType(s[2:])
		if err != nil {
			return nil, "", err
		}
		return reflect.SliceOf(elem), rest, nil
	case strings.HasPrefix(s, "["):
		n, rest, ok := strings.Cut(s[1:], "]")
		length, err := strconv.Atoi(n)
		if !ok || err != nil || length < 0 {
			return nil, "", fmt.Errorf("an array type with length %q", n)
		}
		elem, rest, err := ts.parseType(rest)
		if err != nil {
			return nil, "", err
		}
		return reflect.ArrayOf(length, elem), rest, nil
	case strings.HasPrefix(s, "map["):
		key, rest, err := ts.parseType(s[len("map["):])
		if err != nil {
			return nil, "", err
		}
		rest, ok := strings.CutPrefix(rest, "]")
		if !ok || !key.Comparable() {
			return nil, "", fmt.Errorf("a map type keyed by %s", key)
		}
		elem, rest, err := ts.parseType(rest)
		if err != nil {
			return nil, "", err
		}
		return reflect.MapOf(key, elem), rest, nil
	case strings.HasPrefix(s, "interface {}"):
		return anyType, s[len("interface {}"):], nil
	}
	end := strings.IndexFunc(s, func(r rune) bool {
		return !(r == '_' || r == '.' || '0' <= r && r <= '9' || 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z')
	})
	if end < 0 {
		end = len(s)
	}
	if t, ok := ts.named[s[:end]]; ok {
		return t, s[end:], nil
	}
	return nil, "", fmt.Errorf("type %q, which a tape does not know", s[:end])
}

// scanType rebuilds a recorded scan type, or nil.
func (ts *types) scanType(name string) reflect.Type {
	if !ts.composite {
		return namedScanType(name)
	}
	if name == "" {
		return nil
	}
	t, err := ts.parse(name)
	if err != nil {
		return nil
	}
	return t
}

// knowsScanType reports whether a recording can name scan type t.
func (ts *types) knowsScanType(t reflect.Type) bool {
	if t == nil {
		return false
	}
	if !ts.composite {
		return knownScanType(t)
	}
	_, err := ts.spell(t)
	return err == nil
}

// fromDriver is the package's fromDriver, and a typed value (kind "go") for a
// value of any other type the tape knows.
func (ts *types) fromDriver(v driver.Value) (value, error) {
	out, err := fromDriver(v)
	if err == nil {
		return out, nil
	}
	spelling, err := ts.spell(reflect.TypeOf(v))
	if err != nil {
		return value{}, err
	}
	var b strings.Builder
	if err := ts.write(&b, reflect.ValueOf(v)); err != nil {
		return value{}, err
	}
	return value{kind: "go", typ: spelling, text: b.String()}, nil
}

// driverValue is v as the driver answered it.
func (ts *types) driverValue(v value) (driver.Value, error) {
	if v.kind != "go" {
		return v.driver()
	}
	t, err := ts.parse(v.typ)
	if err != nil {
		return nil, err
	}
	r := &literal{ts: ts, s: v.text}
	out, err := r.read(t)
	if err != nil {
		return nil, fmt.Errorf("a %s value: %w", v.typ, err)
	}
	r.space()
	if r.i != len(r.s) {
		return nil, fmt.Errorf("a %s value: %q after the value", v.typ, r.s[r.i:])
	}
	return out.Interface(), nil
}

// write writes v's literal. The type decides how its literal reads, so only
// an interface's value names its type:
//
//	a boolean, an integer or a float   as strconv formats it (NaN, +Inf, -0)
//	a string, or a type with Text      a Go quoted string
//	a pointer                          nil, or & and its element
//	a slice                            nil, or [a, b], [] when empty
//	an array                           [a, b]
//	a map                              nil, or {k: v, ...} in its keys' literal order, {} when empty
//	an interface                       nil, or its value's type and (literal)
func (ts *types) write(b *strings.Builder, v reflect.Value) error {
	t := v.Type()
	if tt, ok := ts.text[t]; ok {
		text, err := tt.Text(v.Interface())
		if err != nil {
			return fmt.Errorf("a %s value: %w", t, err)
		}
		b.WriteString(strconv.Quote(text))
		return nil
	}
	switch t.Kind() {
	case reflect.Bool:
		b.WriteString(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.WriteString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		b.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		b.WriteString(strconv.FormatFloat(v.Float(), 'g', -1, t.Bits()))
	case reflect.String:
		b.WriteString(strconv.Quote(v.String()))
	case reflect.Pointer:
		if v.IsNil() {
			b.WriteString("nil")
			return nil
		}
		b.WriteString("&")
		return ts.write(b, v.Elem())
	case reflect.Interface:
		if v.IsNil() {
			b.WriteString("nil")
			return nil
		}
		spelling, err := ts.spell(v.Elem().Type())
		if err != nil {
			return err
		}
		b.WriteString(spelling + "(")
		if err := ts.write(b, v.Elem()); err != nil {
			return err
		}
		b.WriteString(")")
	case reflect.Slice, reflect.Array:
		if t.Kind() == reflect.Slice && v.IsNil() {
			b.WriteString("nil")
			return nil
		}
		b.WriteString("[")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			if err := ts.write(b, v.Index(i)); err != nil {
				return err
			}
		}
		b.WriteString("]")
	case reflect.Map:
		if v.IsNil() {
			b.WriteString("nil")
			return nil
		}
		type entry struct{ key, value string }
		entries := make([]entry, 0, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			var k, e strings.Builder
			if err := ts.write(&k, iter.Key()); err != nil {
				return err
			}
			if err := ts.write(&e, iter.Value()); err != nil {
				return err
			}
			entries = append(entries, entry{k.String(), e.String()})
		}
		// A map has no order, and a tape must write the same text twice.
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].key != entries[j].key {
				return entries[i].key < entries[j].key
			}
			return entries[i].value < entries[j].value
		})
		b.WriteString("{")
		for i, e := range entries {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(e.key + ": " + e.value)
		}
		b.WriteString("}")
	default:
		return fmt.Errorf("a %s value, which a tape does not hold", t)
	}
	return nil
}

// literal reads what write wrote.
type literal struct {
	ts *types
	s  string
	i  int
}

func (r *literal) space() {
	for r.i < len(r.s) && r.s[r.i] == ' ' {
		r.i++
	}
}

func (r *literal) errorf(format string, args ...any) error {
	return fmt.Errorf("at byte %d: %s", r.i, fmt.Sprintf(format, args...))
}

// take consumes s when it comes next.
func (r *literal) take(s string) bool {
	r.space()
	if strings.HasPrefix(r.s[r.i:], s) {
		r.i += len(s)
		return true
	}
	return false
}

func (r *literal) want(s string) error {
	if !r.take(s) {
		return r.errorf("want %q", s)
	}
	return nil
}

// token is a number or a boolean: everything up to the next delimiter.
func (r *literal) token() string {
	r.space()
	start := r.i
	for r.i < len(r.s) && !strings.ContainsRune(" ,:)]}", rune(r.s[r.i])) {
		r.i++
	}
	return r.s[start:r.i]
}

func (r *literal) quoted() (string, error) {
	r.space()
	q, err := strconv.QuotedPrefix(r.s[r.i:])
	if err != nil {
		return "", r.errorf("want a quoted string")
	}
	r.i += len(q)
	return strconv.Unquote(q)
}

// list reads the elements of [a, b] or {k: v} up to close, calling each for
// every one.
func (r *literal) list(close string, each func() error) error {
	if r.take(close) {
		return nil
	}
	for {
		if err := each(); err != nil {
			return err
		}
		if r.take(close) {
			return nil
		}
		if err := r.want(","); err != nil {
			return err
		}
	}
}

func (r *literal) read(t reflect.Type) (reflect.Value, error) {
	if tt, ok := r.ts.text[t]; ok {
		text, err := r.quoted()
		if err != nil {
			return reflect.Value{}, err
		}
		v, err := tt.Value(text)
		if err != nil {
			return reflect.Value{}, r.errorf("%v", err)
		}
		out := reflect.ValueOf(v)
		if !out.IsValid() || out.Type() != t {
			return reflect.Value{}, r.errorf("%s's Value returned a %T", t, v)
		}
		return out, nil
	}
	out := reflect.New(t).Elem()
	switch t.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(r.token())
		if err != nil {
			return out, r.errorf("%v", err)
		}
		out.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(r.token(), 10, t.Bits())
		if err != nil {
			return out, r.errorf("%v", err)
		}
		out.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(r.token(), 10, t.Bits())
		if err != nil {
			return out, r.errorf("%v", err)
		}
		out.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(r.token(), t.Bits())
		if err != nil {
			return out, r.errorf("%v", err)
		}
		out.SetFloat(f)
	case reflect.String:
		s, err := r.quoted()
		if err != nil {
			return out, err
		}
		out.SetString(s)
	case reflect.Pointer:
		if r.take("nil") {
			return out, nil
		}
		if err := r.want("&"); err != nil {
			return out, err
		}
		elem, err := r.read(t.Elem())
		if err != nil {
			return out, err
		}
		out.Set(reflect.New(t.Elem()))
		out.Elem().Set(elem)
	case reflect.Interface:
		if r.take("nil") {
			return out, nil
		}
		r.space()
		end := strings.IndexByte(r.s[r.i:], '(')
		if end < 0 {
			return out, r.errorf("want a type and (")
		}
		dynamic, err := r.ts.parse(r.s[r.i : r.i+end])
		if err != nil {
			return out, r.errorf("%v", err)
		}
		if !dynamic.AssignableTo(t) {
			return out, r.errorf("a %s in a %s", dynamic, t)
		}
		r.i += end + 1
		elem, err := r.read(dynamic)
		if err != nil {
			return out, err
		}
		if err := r.want(")"); err != nil {
			return out, err
		}
		out.Set(elem)
	case reflect.Slice:
		if r.take("nil") {
			return out, nil
		}
		if err := r.want("["); err != nil {
			return out, err
		}
		out.Set(reflect.MakeSlice(t, 0, 0))
		err := r.list("]", func() error {
			elem, err := r.read(t.Elem())
			if err == nil {
				out.Set(reflect.Append(out, elem))
			}
			return err
		})
		if err != nil {
			return out, err
		}
	case reflect.Array:
		if err := r.want("["); err != nil {
			return out, err
		}
		n := 0
		err := r.list("]", func() error {
			if n >= t.Len() {
				return r.errorf("more than the %d elements of a %s", t.Len(), t)
			}
			elem, err := r.read(t.Elem())
			if err == nil {
				out.Index(n).Set(elem)
				n++
			}
			return err
		})
		if err != nil {
			return out, err
		}
		if n != t.Len() {
			return out, r.errorf("%d elements of a %s", n, t)
		}
	case reflect.Map:
		if r.take("nil") {
			return out, nil
		}
		if err := r.want("{"); err != nil {
			return out, err
		}
		out.Set(reflect.MakeMap(t))
		err := r.list("}", func() error {
			key, err := r.read(t.Key())
			if err != nil {
				return err
			}
			if err := r.want(":"); err != nil {
				return err
			}
			elem, err := r.read(t.Elem())
			if err != nil {
				return err
			}
			out.SetMapIndex(key, elem)
			return nil
		})
		if err != nil {
			return out, err
		}
	default:
		return out, r.errorf("a %s value, which a tape does not hold", t)
	}
	return out, nil
}
