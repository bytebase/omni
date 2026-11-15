package driver

import "math/bits"

// HeaderByteSize returns the encoded byte size of cols for MaxResultBytes
// accounting, matching the wire framing of repeated column name and type
// fields.
func HeaderByteSize(cols []Column) int64 {
	var n int64
	for _, c := range cols {
		n += delimitedSize(1, int64(len(c.Name)))
		n += delimitedSize(1, int64(len(c.TypeName)))
	}
	return n
}

// RowByteSize returns the encoded byte size of one result row for
// MaxResultBytes accounting.
func RowByteSize(row []Value) int64 {
	var body int64
	for _, v := range row {
		body += v.ByteSize()
	}
	return delimitedSize(1, body)
}

// ByteSize returns the encoded byte size of one cell value inside a row.
func (v Value) ByteSize() int64 {
	var inner int64
	switch v.Kind {
	case Null, Bool:
		inner = 2
	case Bytes:
		inner = delimitedSize(1, int64(len(v.Bytes)))
	case String:
		inner = delimitedSize(1, int64(len(v.String)))
	case Float32:
		inner = 5
	case Float64:
		inner = 9
	case Int32, Int64:
		inner = 1 + varintSize(uint64(v.Int64))
	case Uint32, Uint64:
		inner = 1 + varintSize(v.Uint64)
	case Timestamp:
		inner = delimitedSize(1, timestampBodySize(v))
	case TimestampTZ:
		zone, offset := v.Time.Zone()
		body := timestampBodySize(v)
		if zone != "" {
			body += delimitedSize(1, int64(len(zone)))
		}
		if offset != 0 {
			body += 1 + varintSize(uint64(int64(offset)))
		}
		inner = delimitedSize(1, body)
	case Structured:
		inner = delimitedSize(1, structuredSize(v.Structured))
	default:
		inner = 2
	}
	return delimitedSize(1, inner)
}

func timestampBodySize(v Value) int64 {
	var ts int64
	if sec := v.Time.Unix(); sec != 0 {
		ts += 1 + varintSize(uint64(sec))
	}
	if nsec := v.Time.Nanosecond(); nsec != 0 {
		ts += 1 + varintSize(uint64(nsec))
	}
	body := delimitedSize(1, ts)
	if v.Precision != 0 {
		body += 1 + varintSize(uint64(v.Precision))
	}
	return body
}

func structuredSize(v any) int64 {
	switch x := v.(type) {
	case nil, bool:
		return 2
	case float64:
		return 9
	case string:
		return delimitedSize(1, int64(len(x)))
	case []any:
		var list int64
		for _, item := range x {
			list += delimitedSize(1, structuredSize(item))
		}
		return delimitedSize(1, list)
	case map[string]any:
		var fields int64
		for k, item := range x {
			entry := delimitedSize(1, int64(len(k))) + delimitedSize(1, structuredSize(item))
			fields += delimitedSize(1, entry)
		}
		return delimitedSize(1, fields)
	default:
		return 2
	}
}

func delimitedSize(tagLen, body int64) int64 {
	return tagLen + varintSize(uint64(body)) + body
}

func varintSize(x uint64) int64 {
	return int64((bits.Len64(x|1) + 6) / 7)
}
