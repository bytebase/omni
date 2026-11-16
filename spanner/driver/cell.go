package driver

import (
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/protobuf/types/known/structpb"
)

// nullString is a STRING that may be null: what COLUMN_DEFAULT is read as.
type nullString struct {
	StringVal string
	Valid     bool
}

// scan decodes cells into the destinations, the way spanner.Row.Columns did for
// the sync's queries: the row came as the client's JSON-shaped values, so each
// cell is decoded from what the type says it is.
func scan(cells []Cell, dests ...any) error {
	if len(cells) != len(dests) {
		return fmt.Errorf("row holds %d values, want %d", len(cells), len(dests))
	}
	for i, dest := range dests {
		if err := decodeCell(cells[i], dest); err != nil {
			return fmt.Errorf("column %d: %w", i, err)
		}
	}
	return nil
}

func decodeCell(cell Cell, dest any) error {
	value := cell.Value
	if _, null := value.GetKind().(*structpb.Value_NullValue); value == nil || null || value.GetKind() == nil {
		if nullable, ok := dest.(*nullString); ok {
			nullable.Valid = false
			return nil
		}
		return errors.New("unexpected NULL")
	}
	switch d := dest.(type) {
	case *string:
		*d = value.GetStringValue()
		return nil
	case *bool:
		*d = value.GetBoolValue()
		return nil
	case *nullString:
		d.StringVal = value.GetStringValue()
		d.Valid = true
		return nil
	case *int64:
		// Spanner encodes INT64 as a string to preserve precision, and a value
		// the client read as a number arrives as one.
		if kind := value.GetKind(); kind != nil {
			if _, ok := kind.(*structpb.Value_NumberValue); ok {
				*d = int64(value.GetNumberValue())
				return nil
			}
		}
		parsed, err := strconv.ParseInt(value.GetStringValue(), 10, 64)
		if err != nil {
			return err
		}
		*d = parsed
		return nil
	case *[]string:
		list := value.GetListValue()
		values := make([]string, 0, len(list.GetValues()))
		for _, item := range list.GetValues() {
			values = append(values, item.GetStringValue())
		}
		*d = values
		return nil
	default:
		return fmt.Errorf("cannot decode a %s into %T", cell.Type.GetCode().String(), dest)
	}
}
