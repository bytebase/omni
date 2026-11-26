package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/dynamodb"
)

// Session opens one session. DynamoDB has no session state this driver sets,
// so a session is the driver itself.
func (d *Driver) Session(_ context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

type session struct {
	driver *Driver
}

func (*session) Close() error { return nil }

// Query runs each statement of script and returns one Result per statement
// run. A statement that fails ends the script: its Result carries the error and
// nothing after it runs.
//
// A statement's rows are the items its answer held, one row per item, and its
// columns are the attributes those items named, sorted by name. An item that
// leaves an attribute out, or a later item that names an attribute the earlier
// ones did not, is a nil cell for that row: DynamoDB's items are schemaless, and
// the columns of one answer are the union over it.
//
// QueryOptions.Limit caps the rows the result holds and QueryOptions.MaxResultBytes
// the bytes it takes, both while the answer is read: the driver asks no further
// page once it has the rows, and stops adding them once their encoded size
// passes the byte cap.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	if opts.Explain {
		return nil, errors.New("DynamoDB does not support EXPLAIN")
	}

	singleSQLs := splitStatements(script)
	if len(singleSQLs) == 0 {
		return nil, nil
	}

	var results []*driver.Result
	for _, singleSQL := range singleSQLs {
		startTime := time.Now()
		result, err := s.driver.querySinglePartiQL(ctx, singleSQL.Text, opts)
		stop := false
		if err != nil {
			result = &driver.Result{
				Err: err,
			}
			stop = true
		}
		result.Latency = time.Since(startTime)
		result.Statement = singleSQL.Text
		results = append(results, result)
		if stop {
			break
		}
	}
	return results, nil
}

// splitStatements hands a script to dynamodb.Split, so statement boundaries are
// the PartiQL grammar's judgment, and drops the segments that hold no
// statement.
func splitStatements(script string) []dynamodb.Segment {
	var statements []dynamodb.Segment
	for _, segment := range dynamodb.Split(script) {
		if segment.Empty() {
			continue
		}
		statements = append(statements, segment)
	}
	return statements
}

type dynamodbQueryResultMeta struct {
	columnType string
	value      *driver.Value
}

func (d *Driver) querySinglePartiQL(ctx context.Context, statement string, opts driver.QueryOptions) (*driver.Result, error) {
	result := &driver.Result{}
	request := StatementRequest{
		Statement: statement,
	}
	if opts.Limit > 0 {
		request.Limit = int32(opts.Limit)
	}

	// DynamoDB's items are schemaless, so the answer's columns are the union of
	// the attributes its items named, and the type of a column is the type of
	// the last item that set it.
	rowMap := make(map[string][]*driver.Value)
	columnTypeMap := make(map[string]string)
	totalRowCount := 0
	// What the caller may receive is measured while the answer is read, so a
	// read the byte cap stops does not ask for the pages after it: the cells an
	// item carried are added as it arrives, and no page is asked for once the
	// cap is passed. The exact cut of the rows is made below, where the columns
	// and their framing are known, and that is where Result.Truncated is set:
	// this measure is a lower bound of the same size, so it passes the cap only
	// where the rows do.
	ingestedBytes := int64(0)
read:
	for {
		answer, err := d.client.ExecuteStatement(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, item := range answer.Items {
			totalRowCount++
			meta := convertAttributeValueMapToRow(item)
			allKeySet := make(map[string]bool)
			for key := range rowMap {
				allKeySet[key] = true
			}
			curKeySet := make(map[string]*dynamodbQueryResultMeta, len(meta))
			for key, value := range meta {
				curKeySet[key] = value
				allKeySet[key] = true
			}
			for key := range allKeySet {
				_, inRowMap := rowMap[key]
				_, inCurKeySet := curKeySet[key]
				// 1. The key appears in the rowMap and in the current row, so
				// the value is appended to the column.
				if inRowMap && inCurKeySet {
					rowMap[key] = append(rowMap[key], curKeySet[key].value)
					columnTypeMap[key] = curKeySet[key].columnType
				}
				// 2. The key appears in the rowMap but not in the current row,
				// so a nil is appended.
				if inRowMap && !inCurKeySet {
					rowMap[key] = append(rowMap[key], nil)
				}
				// 3. The key appears in the current row but not in the rowMap,
				// so it is a new column and the earlier rows are backfilled with
				// nils.
				if !inRowMap && inCurKeySet {
					for range totalRowCount - 1 {
						rowMap[key] = append(rowMap[key], nil)
					}
					rowMap[key] = append(rowMap[key], curKeySet[key].value)
					columnTypeMap[key] = curKeySet[key].columnType
				}
			}
			for _, value := range curKeySet {
				ingestedBytes += value.value.ByteSize()
			}
			if opts.MaxResultBytes > 0 && ingestedBytes > opts.MaxResultBytes {
				break read
			}
		}
		if answer.NextToken == "" {
			break
		}
		if opts.Limit > 0 && totalRowCount >= opts.Limit {
			// The rows the caller asked for are read: the rest of the answer is
			// not, whatever the endpoint's next page holds.
			break
		}
		request.NextToken = answer.NextToken
		if opts.Limit > 0 {
			request.Limit = int32(opts.Limit - totalRowCount)
		}
	}
	sortedColumnNames := make([]string, 0, len(rowMap))
	for key := range rowMap {
		sortedColumnNames = append(sortedColumnNames, key)
	}
	slices.Sort(sortedColumnNames)
	for _, name := range sortedColumnNames {
		result.Columns = append(result.Columns, driver.Column{
			Name:     name,
			TypeName: columnTypeMap[name],
		})
	}
	// Flatten the row map to rows.
	totalBytes := driver.HeaderByteSize(result.Columns)
	if len(rowMap) > 0 {
		for i := range totalRowCount {
			row := make([]driver.Value, 0, len(sortedColumnNames))
			for _, columnName := range sortedColumnNames {
				value := rowMap[columnName][i]
				if value == nil {
					row = append(row, driver.Value{Kind: driver.Null})
					continue
				}
				row = append(row, *value)
			}
			result.Rows = append(result.Rows, row)
			if opts.MaxResultBytes > 0 {
				totalBytes += driver.RowByteSize(row)
				if totalBytes > opts.MaxResultBytes {
					result.Truncated = true
					break
				}
			}
		}
	}
	return result, nil
}

func convertAttributeValueMapToRow(items map[string]types.AttributeValue) map[string]*dynamodbQueryResultMeta {
	result := make(map[string]*dynamodbQueryResultMeta, len(items))
	for item, attributeValue := range items {
		var columnType string
		switch attributeValue.(type) {
		case *types.AttributeValueMemberB:
			columnType = "Binary"
		case *types.AttributeValueMemberBOOL:
			columnType = "Boolean"
		case *types.AttributeValueMemberBS:
			columnType = "BinarySet"
		case *types.AttributeValueMemberL:
			columnType = "List"
		case *types.AttributeValueMemberM:
			columnType = "Map"
		case *types.AttributeValueMemberN:
			columnType = "Number"
		case *types.AttributeValueMemberNS:
			columnType = "NumberSet"
		case *types.AttributeValueMemberNULL:
			columnType = "Null"
		case *types.AttributeValueMemberS:
			columnType = "String"
		case *types.AttributeValueMemberSS:
			columnType = "StringSet"
		default:
			columnType = "Unknown"
		}
		value, err := convertAttributeValueToValue(attributeValue)
		if err != nil {
			value = driver.Value{Kind: driver.String, String: err.Error()}
		}
		result[item] = &dynamodbQueryResultMeta{
			columnType: columnType,
			value:      &value,
		}
	}
	return result
}

func convertAttributeValueToValue(attributeValue types.AttributeValue) (driver.Value, error) {
	switch attributeValue := attributeValue.(type) {
	case *types.AttributeValueMemberB:
		// A binary is bytes, not text: a proto string holding one that is not
		// valid UTF-8 cannot be marshalled, so it crosses as the contract's
		// Bytes kind and reaches the caller unmangled.
		return driver.Value{Kind: driver.Bytes, Bytes: attributeValue.Value}, nil
	case *types.AttributeValueMemberBOOL:
		return driver.Value{Kind: driver.Bool, Bool: attributeValue.Value}, nil
	case *types.AttributeValueMemberBS, *types.AttributeValueMemberL, *types.AttributeValueMemberM,
		*types.AttributeValueMemberNS, *types.AttributeValueMemberSS:
		// A set, a list and a map are documents rather than scalars, so they
		// cross as their JSON text. A binary inside a document is base64, as
		// json.Marshal would render any other invalid UTF-8 as a replacement
		// character and lose it.
		document := convertAttributeValueToGoPrimitives(attributeValue)
		b, err := json.Marshal(document)
		if err != nil {
			return driver.Value{}, fmt.Errorf("failed to marshal attribute value: %w", err)
		}
		return driver.Value{Kind: driver.String, String: string(b)}, nil
	case *types.AttributeValueMemberN:
		return driver.Value{Kind: driver.String, String: attributeValue.Value}, nil
	case *types.AttributeValueMemberNULL:
		return driver.Value{Kind: driver.Null}, nil
	case *types.AttributeValueMemberS:
		return driver.Value{Kind: driver.String, String: attributeValue.Value}, nil
	}
	return driver.Value{}, fmt.Errorf("unsupported attribute value type: %T", attributeValue)
}

func convertAttributeValueToGoPrimitives(av types.AttributeValue) any {
	switch av := av.(type) {
	case *types.AttributeValueMemberB:
		return base64.StdEncoding.EncodeToString(av.Value)
	case *types.AttributeValueMemberBOOL:
		return av.Value
	case *types.AttributeValueMemberBS:
		ss := make([]string, 0, len(av.Value))
		for _, b := range av.Value {
			ss = append(ss, base64.StdEncoding.EncodeToString(b))
		}
		return ss
	case *types.AttributeValueMemberL:
		ss := make([]any, 0, len(av.Value))
		for _, v := range av.Value {
			ss = append(ss, convertAttributeValueToGoPrimitives(v))
		}
		return ss
	case *types.AttributeValueMemberM:
		m := make(map[string]any, len(av.Value))
		for k, v := range av.Value {
			m[k] = convertAttributeValueToGoPrimitives(v)
		}
		return m
	case *types.AttributeValueMemberN:
		return av.Value
	case *types.AttributeValueMemberNS:
		return av.Value
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return av.Value
	case *types.AttributeValueMemberSS:
		return av.Value
	}
	return nil
}
