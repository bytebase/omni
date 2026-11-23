package testcontainer

import (
	"fmt"
	"math/big"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/shopspring/decimal"

	"github.com/bytebase/omni/enginetest/tape"
)

// Types are the client's own Go types that its rows hold, for a tape to write.
// clickhouse-go answers a Decimal and a UUID as their text (both implement
// driver.Valuer), but names them as a column's scan type and holds them inside
// an Array, a Tuple or a Map; an Int128 or wider is a *big.Int, an IPv4 or IPv6
// a net.IP, and a geometry one of orb's. A Variant, a Dynamic and a JSON value
// are the client's structs with unexported fields, which a tape does not
// hold: a recording that reads one fails.
var Types = []tape.Type{
	{
		Type: reflect.TypeOf(decimal.Decimal{}),
		// The coefficient and exponent the client built it from, which a
		// Decimal's own text trims.
		Text: func(v any) (string, error) {
			d := v.(decimal.Decimal)
			return d.Coefficient().String() + "e" + strconv.FormatInt(int64(d.Exponent()), 10), nil
		},
		Value: func(text string) (any, error) {
			coefficient, exponent, ok := strings.Cut(text, "e")
			n, okN := new(big.Int).SetString(coefficient, 10)
			e, err := strconv.ParseInt(exponent, 10, 32)
			if !ok || !okN || err != nil {
				return nil, fmt.Errorf("a decimal %q", text)
			}
			return decimal.NewFromBigInt(n, int32(e)), nil
		},
	},
	{
		Type: reflect.TypeOf(uuid.UUID{}),
		Text: func(v any) (string, error) { return v.(uuid.UUID).String(), nil },
		Value: func(text string) (any, error) {
			return uuid.Parse(text)
		},
	},
	{
		Type: reflect.TypeOf(big.Int{}),
		Text: func(v any) (string, error) {
			n := v.(big.Int)
			return n.String(), nil
		},
		Value: func(text string) (any, error) {
			n, ok := new(big.Int).SetString(text, 10)
			if !ok {
				return nil, fmt.Errorf("an integer %q", text)
			}
			return *n, nil
		},
	},
	// The rest are written by their kind: an IPv4's four bytes stay four.
	{Type: reflect.TypeOf(net.IP(nil))},
	{Type: reflect.TypeOf(orb.Point{})},
	{Type: reflect.TypeOf(orb.LineString(nil))},
	{Type: reflect.TypeOf(orb.MultiLineString(nil))},
	{Type: reflect.TypeOf(orb.Ring(nil))},
	{Type: reflect.TypeOf(orb.Polygon(nil))},
	{Type: reflect.TypeOf(orb.MultiPolygon(nil))},
}
