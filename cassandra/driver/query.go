package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"reflect"
	"strings"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"gopkg.in/inf.v0"

	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/driver"
)

// Session returns one session. Cassandra has no per-session state this driver
// sets, so a session is the driver itself.
func (d *Driver) Session(_ context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

type session struct {
	driver *Driver
}

func (*session) Close() error { return nil }

// Query runs each statement of script in order and returns one Result per
// statement run. A statement's error is gocql's, as the package doc describes.
//
// A statement that fails ends the script with its error on its own Result;
// the error returned is nil, since the script ran as far as it could.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	var results []*driver.Result
	for _, segment := range cassandra.Split(script) {
		// An empty segment, only whitespace and comments, is no statement:
		// sending a trailing "-- done" would fail a script whose earlier
		// statements already ran.
		if segment.Empty {
			continue
		}
		statement := trimStatement(segment.Text)
		if statement == "" {
			continue
		}
		startTime := time.Now()
		result, err := s.run(ctx, statement, opts)
		if err != nil {
			result = &driver.Result{Statement: statement, Err: fmt.Errorf("failed to fetch row data: %w", err)}
		}
		result.Latency = time.Since(startTime)
		results = append(results, result)
		if result.Err != nil {
			break
		}
	}
	return results, nil
}

// run runs one statement page by page: Limit is the page size and the row
// count a statement stops at, and MaxResultBytes stops the reading and sets
// Truncated.
func (s *session) run(ctx context.Context, statement string, opts driver.QueryOptions) (*driver.Result, error) {
	c, err := s.driver.conn()
	if err != nil {
		return nil, err
	}
	result := &driver.Result{Statement: statement}
	pageSize := 0
	switch {
	case opts.Limit > 0:
		pageSize = opts.Limit
	case opts.MaxResultBytes > 0:
		// Without a page size Cassandra answers a whole table in one reply,
		// and the cap could not stop the reading: the page bounds what one
		// answer holds, the cap bounds how many rows are kept.
		pageSize = pageSizeForACappedRead
	}
	var (
		totalBytes int64
		pageState  []byte
	)
	for {
		columns, rows, nextPageState, err := c.Query(ctx, statement, pageSize, pageState)
		if err != nil {
			return nil, err
		}
		if len(result.Columns) == 0 {
			for _, column := range columns {
				result.Columns = append(result.Columns, driver.Column{
					Name:     column.Name,
					TypeName: typeToString(column.TypeInfo.Type()),
				})
			}
			totalBytes = driver.HeaderByteSize(result.Columns)
		}
		groups := valueGroups(columns, rows)
		for _, rowData := range rows {
			row := make([]driver.Value, 0, len(groups))
			next := 0
			for _, size := range groups {
				if next+size > len(rowData.Values) {
					// A row shorter than its columns is read as it came, so a
					// client that answers fewer values never panics here.
					row = append(row, convertRowValue(rowData.Values[next]))
					next++
					continue
				}
				if size == 1 {
					row = append(row, convertRowValue(rowData.Values[next]))
				} else {
					row = append(row, convertTuple(rowData.Values[next:next+size]))
				}
				next += size
			}
			result.Rows = append(result.Rows, row)
			totalBytes += driver.RowByteSize(row)
			n := len(result.Rows)
			if opts.MaxResultBytes > 0 && totalBytes > opts.MaxResultBytes {
				result.Truncated = true
				return result, nil
			}
			if opts.Limit > 0 && opts.Limit == n {
				return result, nil
			}
		}
		if len(nextPageState) == 0 {
			break
		}
		pageState = nextPageState
	}
	return result, nil
}

// valueGroups returns how many of a row's values each column contributes.
// gocql expands a tuple column into one value per element and names those row
// columns "<name>[i]"; every other column is one value, even one whose alias
// looks like a tuple element ("x[0]"). A grouping that does not consume exactly
// the row's values is not one this driver can align, and every column then takes
// one value.
func valueGroups(columns []gocql.ColumnInfo, rows []gocql.RowData) []int {
	groups := make([]int, len(columns))
	for i := range groups {
		groups[i] = 1
	}
	if len(rows) == 0 {
		return groups
	}
	names := rows[0].Columns
	if len(names) == 0 || len(names) == len(columns) {
		return groups
	}
	at := 0
	for i, column := range columns {
		if at >= len(names) {
			break
		}
		size := 1
		if arity := tupleArity(column.TypeInfo); arity > 0 && names[at] == column.Name+"[0]" {
			size = min(arity, len(names)-at)
		}
		groups[i] = size
		at += size
	}
	if at != len(names) {
		for i := range groups {
			groups[i] = 1
		}
	}
	return groups
}

// tupleArity is how many values a tuple column expands into, and 0 for a column
// that is not a tuple. A live column's type is gocql's own tuple type; a tape's
// carries the element count it recorded.
func tupleArity(info gocql.TypeInfo) int {
	if info == nil || info.Type() != gocql.TypeTuple {
		return 0
	}
	if tuple, ok := info.(gocql.TupleTypeInfo); ok {
		return len(tuple.Elems)
	}
	if counted, ok := info.(interface{ Elements() int }); ok {
		return counted.Elements()
	}
	return 0
}

// convertTuple reads a tuple column's elements as the JSON spelling the column
// shows, the way a list or a user-defined type is read. gocql expands a tuple
// into one value per element and leaves each nil for a NULL, so a tuple whose
// every element is NULL is read as NULL: gocql cannot say whether the column
// itself was NULL, and a tuple of NULLs is the one case the contract loses.
func convertTuple(values []any) driver.Value {
	allNull := true
	for _, value := range values {
		if !isNull(value) {
			allNull = false
			break
		}
	}
	if allNull {
		return driver.Value{Kind: driver.Null}
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return driver.Value{Kind: driver.String, String: fmt.Sprintf("%v", values)}
	}
	return driver.Value{Kind: driver.String, String: string(encoded)}
}

// pageSizeForACappedRead is the page size a read with only a byte cap is
// fetched in.
const pageSizeForACappedRead = 1000

// isNull reports whether a value gocql handed back is a NULL.
func isNull(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Ptr && rv.IsNil()
}

// convertRowValue reads one cell gocql handed back. gocql allocates a pointer
// per column from the CQL type's Go type, and a list, map, set, tuple or UDT
// value is read as its JSON spelling, which is what Bytebase showed. A NULL is
// that pointer left nil: a nil decimal or varint is a nil pointer one level
// down, as its Go type is a pointer itself.
func convertRowValue(v any) driver.Value {
	if isNull(v) {
		return driver.Value{Kind: driver.Null}
	}
	switch v := v.(type) {
	case *string:
		return driver.Value{Kind: driver.String, String: *v}
	case *int64:
		return driver.Value{Kind: driver.Int64, Int64: *v}
	case *time.Duration:
		return driver.Value{Kind: driver.String, String: timeOfDayString(*v)}
	case *time.Time:
		return driver.Value{Kind: driver.Timestamp, Time: *v, Precision: 3}
	case *[]byte:
		return driver.Value{Kind: driver.Bytes, Bytes: *v}
	case *bool:
		return driver.Value{Kind: driver.Bool, Bool: *v}
	case *float32:
		return driver.Value{Kind: driver.Float32, Float64: float64(*v)}
	case *float64:
		return driver.Value{Kind: driver.Float64, Float64: *v}
	case *int:
		return driver.Value{Kind: driver.Int64, Int64: int64(*v)}
	case *int16:
		return driver.Value{Kind: driver.Int32, Int64: int64(*v)}
	case *int8:
		return driver.Value{Kind: driver.Int32, Int64: int64(*v)}
	case *gocql.UUID:
		return driver.Value{Kind: driver.String, String: v.String()}
	case **inf.Dec:
		if *v == nil {
			return driver.Value{Kind: driver.Null}
		}
		return driver.Value{Kind: driver.String, String: (*v).String()}
	case **big.Int:
		if *v == nil {
			return driver.Value{Kind: driver.Null}
		}
		return driver.Value{Kind: driver.String, String: (*v).String()}
	case *gocql.Duration:
		if v == nil {
			return driver.Value{Kind: driver.Null}
		}
		display := time.Duration(v.Nanoseconds).String()
		if v.Days != 0 {
			display = fmt.Sprintf("%dd%s", v.Days, display)
		}
		if v.Months != 0 {
			display = fmt.Sprintf("%dmo%s", v.Months, display)
		}
		return driver.Value{Kind: driver.String, String: display}

	default:
		// A collection or a user-defined type is read as its JSON spelling.
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			slog.Error("failed to marshal Cassandra complex value", "error", err)
			return driver.Value{Kind: driver.String, String: fmt.Sprintf("%v", v)}
		}
		return driver.Value{Kind: driver.String, String: string(jsonBytes)}
	}
}

// timeOfDayString renders a CQL time, which gocql carries as nanoseconds since
// midnight, the way Bytebase showed it: HH:MM:SS and up to six fractional
// digits with trailing zeros trimmed.
func timeOfDayString(d time.Duration) string {
	const day = 24 * time.Hour
	ns := int64(d) % int64(day)
	if ns < 0 {
		ns += int64(day)
	}
	hour := ns / int64(time.Hour)
	minute := ns / int64(time.Minute) % 60
	second := ns / int64(time.Second) % 60
	micros := ns / int64(time.Microsecond) % 1_000_000
	display := fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
	if micros > 0 {
		display += "." + strings.TrimRight(fmt.Sprintf("%06d", micros), "0")
	}
	return display
}

// trimStatement trims a segment the way Bytebase's query path trimmed the
// statement it sent, so a caller reads the text without its padding.
func trimStatement(statement string) string {
	return strings.Trim(statement, " \t\n\r\f;")
}

// typeToString is Cassandra's own name for a CQL type, in the lowercase gocql
// reports.
func typeToString(t gocql.Type) string {
	switch t {
	case gocql.TypeCustom:
		return "custom"
	case gocql.TypeAscii:
		return "ascii"
	case gocql.TypeBigInt:
		return "bigint"
	case gocql.TypeBlob:
		return "blob"
	case gocql.TypeBoolean:
		return "boolean"
	case gocql.TypeCounter:
		return "counter"
	case gocql.TypeDecimal:
		return "decimal"
	case gocql.TypeDouble:
		return "double"
	case gocql.TypeFloat:
		return "float"
	case gocql.TypeInt:
		return "int"
	case gocql.TypeText:
		return "text"
	case gocql.TypeTimestamp:
		return "timestamp"
	case gocql.TypeUUID:
		return "uuid"
	case gocql.TypeVarchar:
		return "varchar"
	case gocql.TypeVarint:
		return "varint"
	case gocql.TypeTimeUUID:
		return "timeuuid"
	case gocql.TypeInet:
		return "inet"
	case gocql.TypeDate:
		return "date"
	case gocql.TypeTime:
		return "time"
	case gocql.TypeSmallInt:
		return "smallint"
	case gocql.TypeTinyInt:
		return "tinyint"
	case gocql.TypeDuration:
		return "duration"
	case gocql.TypeList:
		return "list"
	case gocql.TypeMap:
		return "map"
	case gocql.TypeSet:
		return "set"
	case gocql.TypeUDT:
		return "udt"
	case gocql.TypeTuple:
		return "tuple"
	default:
		return fmt.Sprintf("unknown(%d)", t)
	}
}
