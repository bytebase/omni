package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/partiql"
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
		result, err := s.driver.querySinglePartiQL(ctx, singleSQL.Text, opts.Limit)
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

// splitStatements hands a script to partiql.Split, so statement boundaries are
// the PartiQL grammar's judgment, and drops the segments that hold no
// statement.
func splitStatements(script string) []partiql.Segment {
	var statements []partiql.Segment
	for _, segment := range partiql.Split(script) {
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

func (d *Driver) querySinglePartiQL(ctx context.Context, statement string, limit int) (*driver.Result, error) {
	result := &driver.Result{}
	request := StatementRequest{
		Statement: statement,
	}
	if limit > 0 {
		request.Limit = int32(limit)
	}

	// DynamoDB's items are schemaless, so the answer's columns are the union of
	// the attributes its items named, and the type of a column is the type of
	// the last item that set it.
	rowMap := make(map[string][]*driver.Value)
	columnTypeMap := make(map[string]string)
	totalRowCount := 0
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
		}
		if answer.NextToken == "" {
			break
		}
		request.NextToken = answer.NextToken
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
		return driver.Value{Kind: driver.String, String: string(attributeValue.Value)}, nil
	case *types.AttributeValueMemberBOOL:
		return driver.Value{Kind: driver.Bool, Bool: attributeValue.Value}, nil
	case *types.AttributeValueMemberBS, *types.AttributeValueMemberL, *types.AttributeValueMemberM,
		*types.AttributeValueMemberNS, *types.AttributeValueMemberSS:
		// A set, a list and a map are documents rather than scalars, so they
		// cross as their JSON text, as Bytebase's driver answered them.
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
		return string(av.Value)
	case *types.AttributeValueMemberBOOL:
		return av.Value
	case *types.AttributeValueMemberBS:
		ss := make([]string, 0, len(av.Value))
		for _, b := range av.Value {
			ss = append(ss, string(b))
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
