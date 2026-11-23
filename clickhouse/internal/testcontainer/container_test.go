package testcontainer

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/bytebase/omni/enginetest"
)

func TestPin(t *testing.T) {
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		t.Fatal(err)
	}
	if pin.Image != "clickhouse/clickhouse-server:"+pin.Tag {
		t.Errorf("image = %q, tag = %q", pin.Image, pin.Tag)
	}
	if !strings.HasPrefix(pin.Digest, "sha256:") {
		t.Errorf("digest = %q", pin.Digest)
	}
}

func TestDatabase(t *testing.T) {
	for _, tc := range []struct {
		test, want string
	}{
		{"TestSyncDatabase", "omni_testsyncdatabase"},
		{"TestExecuteAutoCommit", "omni_testexecuteautocommit"},
	} {
		got := Database(fakeTB{name: tc.test})
		if got != tc.want {
			t.Errorf("Database(%q) = %q, want %q", tc.test, got, tc.want)
		}
	}
}

type fakeTB struct {
	testing.TB
	name string
}

func (f fakeTB) Name() string { return f.name }

// Each type a tape writes by Text reads back as a value it writes the same
// way, which is what a replay hands the driver.
func TestTypesRoundTrip(t *testing.T) {
	samples := map[reflect.Type][]any{
		reflect.TypeOf(decimal.Decimal{}): {decimal.New(150, -2), decimal.New(-1, -3), decimal.New(0, -2), decimal.RequireFromString("-12345678901234567890.0123456789")},
		reflect.TypeOf(uuid.UUID{}):       {uuid.MustParse("0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f"), uuid.UUID{}},
		reflect.TypeOf(big.Int{}):         {*big.NewInt(0), *new(big.Int).Lsh(big.NewInt(-1), 127)},
	}
	for _, typ := range Types {
		if typ.Text == nil {
			continue
		}
		values := samples[typ.Type]
		if len(values) == 0 {
			t.Errorf("no sample of %s", typ.Type)
		}
		for _, v := range values {
			text, err := typ.Text(v)
			if err != nil {
				t.Fatalf("%s %v: %v", typ.Type, v, err)
			}
			back, err := typ.Value(text)
			if err != nil {
				t.Fatalf("%s %q: %v", typ.Type, text, err)
			}
			again, err := typ.Text(back)
			if err != nil || again != text || reflect.TypeOf(back) != typ.Type {
				t.Errorf("%s %q came back as %T %q (%v)", typ.Type, text, back, again, err)
			}
		}
	}
	if text, _ := Types[0].Text(decimal.New(150, -2)); text != "150e-2" {
		t.Errorf("1.50 is written %q, want its coefficient and exponent", text)
	}
}
