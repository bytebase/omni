package testcontainer

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// everyKind is one value of every BSON kind a document can hold, the deprecated
// ones included, with the values a test's own data reaches.
func everyKind(t *testing.T) []any {
	t.Helper()
	decimal, err := bson.ParseDecimal128("1.5E+100")
	require.NoError(t, err)
	objectID, err := bson.ObjectIDFromHex("6afaf7c2c0c20ddc7f2271b8")
	require.NoError(t, err)
	return []any{
		nil,
		true,
		int32(-7),
		int64(9223372036854775807),
		float64(1779696815227),
		math.Copysign(0, -1),
		math.Inf(1),
		math.NaN(),
		"a string with \"quotes\", a \\ backslash and é",
		bson.D{{Key: "a", Value: int32(1)}, {Key: "b", Value: "x"}},
		bson.M{"sorted": true, "any": bson.A{int32(1), "two"}},
		bson.A{int32(1), "two", bson.D{{Key: "k", Value: true}}},
		objectID,
		bson.DateTime(1700000000123),
		bson.Timestamp{T: 1, I: 2},
		bson.Binary{Subtype: 0x02, Data: []byte{0, 1, 254, 255}},
		bson.Binary{Subtype: 0x04, Data: make([]byte, 16)},
		bson.Regex{Pattern: "a.*", Options: "i"},
		bson.MinKey{},
		bson.MaxKey{},
		bson.Symbol("symbol"),
		bson.DBPointer{DB: "db", Pointer: objectID},
		bson.Undefined{},
		decimal,
		bson.JavaScript("x = 1"),
		bson.CodeWithScope{Code: "x = 1", Scope: bson.D{{Key: "a", Value: int32(1)}}},
	}
}

// TestTheTapeHoldsEveryBSONKind holds the tape's encoding to what a document
// can hold: an answer keeps the bytes of every kind, documents included, so a
// replay hands the driver what the recording's client did.
func TestTheTapeHoldsEveryBSONKind(t *testing.T) {
	for _, value := range everyKind(t) {
		t.Run(kindName(value), func(t *testing.T) {
			doc := bson.M{"name": "value", "value": value}

			encoded, err := encodeDocs([]bson.M{doc})
			require.NoError(t, err)
			decoded, err := decodeDocs(encoded)
			require.NoError(t, err)
			require.Len(t, decoded, 1)
			want, err := encodeDoc(doc)
			require.NoError(t, err)
			got, err := encodeDoc(decoded[0])
			require.NoError(t, err)
			require.Equal(t, string(want), string(got))
		})
	}
}

// TestTheTapeHoldsEveryKindOfExecuteValue holds the same for the values a
// mongosh execution answers with, which the tape carries as an array.
func TestTheTapeHoldsEveryKindOfExecuteValue(t *testing.T) {
	for _, value := range everyKind(t) {
		t.Run(kindName(value), func(t *testing.T) {
			encoded, err := encodeDoc(bson.D{{Key: "values", Value: bson.A{value}}})
			require.NoError(t, err)
			var payload struct {
				Values []json.RawMessage `json:"values"`
			}
			require.NoError(t, unmarshalJSON(encoded, &payload))
			require.Len(t, payload.Values, 1)
			decoded, err := decodeValue(payload.Values[0])
			require.NoError(t, err)
			reencoded, err := encodeDoc(bson.D{{Key: "values", Value: bson.A{decoded}}})
			require.NoError(t, err)
			require.Equal(t, string(encoded), string(reencoded))
		})
	}
}

func kindName(value any) string {
	switch value := value.(type) {
	case nil:
		return "null"
	case bson.Binary:
		return "binary-" + string(rune('0'+value.Subtype))
	case bson.D:
		return "document"
	case bson.M:
		return "map"
	case bson.A:
		return "array"
	default:
		return typeName(value)
	}
}

func typeName(value any) string {
	switch value.(type) {
	case bool:
		return "bool"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case float64:
		return "double"
	case string:
		return "string"
	case bson.ObjectID:
		return "objectid"
	case bson.DateTime:
		return "datetime"
	case bson.Timestamp:
		return "timestamp"
	case bson.Regex:
		return "regex"
	case bson.MinKey:
		return "minkey"
	case bson.MaxKey:
		return "maxkey"
	case bson.Symbol:
		return "symbol"
	case bson.DBPointer:
		return "dbpointer"
	case bson.Undefined:
		return "undefined"
	case bson.Decimal128:
		return "decimal128"
	case bson.JavaScript:
		return "javascript"
	case bson.CodeWithScope:
		return "codewithscope"
	}
	return "unknown"
}
