package testcontainer

import (
	"math/big"
	"testing"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/inf.v0"
)

// TestTheTapeHoldsEveryCellKind holds a cell's encoding to every Go type the
// driver's row conversion reads: what a recording writes, a replay hands back
// as the same value.
func TestTheTapeHoldsEveryCellKind(t *testing.T) {
	uuid, err := gocql.ParseUUID("6afaf7c2-c0c2-0ddc-7f22-71b8f0a0b0c0")
	require.NoError(t, err)
	decimal, ok := new(inf.Dec).SetString("-1.2500")
	require.True(t, ok)
	varint, ok := new(big.Int).SetString("123456789012345678901234567890", 10)
	require.True(t, ok)
	timestamp := time.Date(2026, 1, 2, 3, 4, 5, 678000000, time.UTC)
	blob := []byte{0, 1, 254, 255}
	empty := ""

	values := []any{
		&empty,
		new("a string with \"quotes\" and é"),
		(*string)(nil),
		new(int(7)),
		new(int64(-9223372036854775808)),
		new(int32(300)),
		new(int16(-300)),
		new(int8(-7)),
		new(true),
		new(float32(1.25)),
		new(float64(-1.5)),
		&blob,
		&timestamp,
		new(2*time.Hour + 3*time.Second),
		new(gocql.Duration{Months: 1, Days: 2, Nanoseconds: 3}),
		&uuid,
		&decimal,
		new(*inf.Dec),
		&varint,
		new(*big.Int),
		&[]int{1, 2, 3},
		&[]int{9223372036854775807},
		&map[string]int{"b": 2, "a": 1},
		&map[string]any{},
		&[]string{},
		&[]any{int64(1), "two"},
	}
	for _, value := range values {
		encoded, err := encodeCell(value)
		require.NoError(t, err, "%T", value)
		decoded, err := decodeCell(encoded)
		require.NoError(t, err, "%T", value)
		reencoded, err := encodeCell(decoded)
		require.NoError(t, err, "%T", value)
		require.Equal(t, encoded, reencoded, "%T did not survive the tape", value)
	}
}

// TestTheTapeRefusesAValueItCannotHold fails a recording rather than writing an
// answer a replay would read as something else.
func TestTheTapeRefusesAValueItCannotHold(t *testing.T) {
	_, err := encodeCell(make(chan int))
	require.Error(t, err)
}
