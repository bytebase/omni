package driver

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	dbsql "github.com/databricks/databricks-sdk-go/service/sql"

	dbxanalysis "github.com/bytebase/omni/databricks/analysis"
	"github.com/bytebase/omni/driver"
)

var (
	_ driver.Querier = (*Driver)(nil)
	_ driver.Session = (*session)(nil)
)

// session is the driver: the SQL Statement Execution API has no session state a
// statement can set, so each statement carries its own catalog and schema.
type session struct {
	driver *Driver
}

// Session returns a session over the driver's workspace and warehouse.
func (d *Driver) Session(context.Context) (driver.Session, error) {
	return &session{driver: d}, nil
}

// Close releases nothing.
func (*session) Close() error {
	return nil
}

// Query runs each statement of script in order and returns one Result per
// statement run. A statement that fails ends the script: its Result carries the
// error, and the error returned is nil. An error returned means the script
// could not be split.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	commands, err := splitCommands(script)
	if err != nil {
		return nil, fmt.Errorf("failed to split statements: %w", err)
	}
	if len(commands) == 0 {
		return nil, nil
	}

	results := make([]*driver.Result, 0, len(commands))
	for _, command := range commands {
		statement := dbxanalysis.TrimStatement(command.Text)
		if opts.Explain {
			statement = dbxanalysis.ExplainStatement(statement, opts.ExplainFormat)
		}

		start := time.Now()
		result, err := s.driver.queryStatement(ctx, statement, opts)
		if err != nil {
			results = append(results, &driver.Result{Statement: statement, Err: err, Latency: time.Since(start)})
			break
		}
		if format, executed, ok := dbxanalysis.DescribeExplain(statement); ok {
			result.Plan = &driver.Plan{Format: format, Executed: executed}
		}
		result.Statement = statement
		result.Latency = time.Since(start)
		results = append(results, result)
	}
	return results, nil
}

func splitCommands(script string) ([]driver.Command, error) {
	statements, err := dbxanalysis.Split(script)
	if err != nil {
		return nil, err
	}
	commands := make([]driver.Command, 0, len(statements))
	for _, stmt := range statements {
		if stmt.Empty {
			continue
		}
		commands = append(commands, driver.Command{Text: stmt.Text, ByteStart: stmt.Start, ByteEnd: stmt.End})
	}
	return commands, nil
}

// queryStatement runs one statement through the warehouse and decodes its
// answer. The SDK's ExecuteAndWait polls a long-running statement and returns
// the completed response, or the text of a FAILED status.
func (d *Driver) queryStatement(ctx context.Context, statement string, opts driver.QueryOptions) (*driver.Result, error) {
	request := dbsql.ExecuteStatementRequest{
		Statement:   statement,
		WarehouseId: d.warehouseID,
		Catalog:     d.catalog,
		Schema:      opts.Schema,
	}
	// The user-selected row limit is the request's own row limit, as Bytebase
	// sent it, rather than a rewrite of the statement.
	if opts.Limit > 0 {
		request.RowLimit = int64(opts.Limit)
	}
	response, err := d.client.StatementExecution.ExecuteAndWait(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	statementID := response.StatementId
	result, err := resultMetadata(response, statement, dbxanalysis.ClassifyStatement(statement))
	if err != nil {
		return nil, err
	}
	totalBytes := driver.HeaderByteSize(result.Columns)
	var rows [][]string
	if response.Result != nil {
		rows = response.Result.DataArray
	}
	for {
		capReached, err := appendResultRows(result, rows, opts.MaxResultBytes, &totalBytes)
		if err != nil {
			return nil, err
		}
		if response.Result == nil || response.Result.NextChunkIndex == 0 || capReached {
			return result, nil
		}
		next := response.Result.NextChunkIndex
		chunk, err := d.client.StatementExecution.GetStatementResultChunkNByStatementIdAndChunkIndex(ctx, statementID, next)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch result chunk %d: %w", next, err)
		}
		if chunk == nil {
			return nil, fmt.Errorf("no result chunk %d", next)
		}
		rows = chunk.DataArray
		response.Result.NextChunkIndex = chunk.NextChunkIndex
	}
}

// resultMetadata reads the statement's columns and the server's truncation, and
// leaves its rows for appendResultRows so each result chunk is decoded once.
func resultMetadata(response *dbsql.StatementResponse, statement string, returnsRows bool) (*driver.Result, error) {
	if response == nil {
		return nil, fmt.Errorf("no response")
	}

	result := &driver.Result{Statement: statement}
	if response.Manifest != nil && response.Manifest.Schema != nil {
		result.Columns = toColumns(response.Manifest.Schema.Columns)
	}
	if returnsRows && len(result.Columns) == 0 {
		return nil, fmt.Errorf("missing column info")
	}
	if response.Manifest != nil {
		result.Truncated = response.Manifest.Truncated
	}
	return result, nil
}

// appendResultRows decodes one chunk's rows onto result. It reports whether
// the local byte cap stopped it; the server's own truncation is separate.
func appendResultRows(result *driver.Result, rows [][]string, maxResultBytes int64, totalBytes *int64) (bool, error) {
	for _, rowData := range rows {
		if len(rowData) != len(result.Columns) {
			return false, fmt.Errorf("row has %d cells but the statement states %d columns", len(rowData), len(result.Columns))
		}
		row := make([]driver.Value, len(result.Columns))
		for i, value := range rowData {
			converted, err := toValue(result.Columns[i].TypeName, value)
			if err != nil {
				return false, fmt.Errorf("column %q: %w", result.Columns[i].Name, err)
			}
			row[i] = converted
		}
		result.Rows = append(result.Rows, row)
		if maxResultBytes > 0 {
			*totalBytes += driver.RowByteSize(row)
			if *totalBytes > maxResultBytes {
				result.Truncated = true
				return true, nil
			}
		}
	}
	return false, nil
}

// resultFromResponse turns the statement execution API's answer into the
// contract's result. A query is expected to state its columns even when it
// returned no rows; a command states neither columns nor an affected count.
func resultFromResponse(response *dbsql.StatementResponse, statement string, returnsRows bool, maxResultBytes int64) (*driver.Result, error) {
	result, err := resultMetadata(response, statement, returnsRows)
	if err != nil {
		return nil, err
	}
	totalBytes := driver.HeaderByteSize(result.Columns)
	var rows [][]string
	if response.Result != nil {
		rows = response.Result.DataArray
	}
	if _, err := appendResultRows(result, rows, maxResultBytes, &totalBytes); err != nil {
		return nil, err
	}
	return result, nil
}

// toColumns names the columns the statement returned. ColumnInfo.TypeName is
// the base type, which is what Bytebase reported; a complex type's full
// spelling is TypeText and is not used here.
func toColumns(columnInfo []dbsql.ColumnInfo) []driver.Column {
	columns := make([]driver.Column, len(columnInfo))
	for i, col := range columnInfo {
		columns[i] = driver.Column{
			Name:     col.Name,
			TypeName: strings.ToUpper(string(col.TypeName)),
		}
	}
	return columns
}

// toValue maps one JSON_ARRAY cell onto the contract's kinds, as Bytebase's
// conversion did: the typed primitives become typed values, and every other
// type stays the text the API stated. A JSON null is the empty string in this
// format, so an empty value in a non-STRING column is NULL; an empty STRING
// stays the empty string.
func toValue(typeName string, value string) (driver.Value, error) {
	switch strings.ToUpper(typeName) {
	case string(dbsql.ColumnInfoTypeNameBoolean):
		if value == "" {
			return driver.Value{Kind: driver.Null}, nil
		}
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return driver.Value{}, err
		}
		return driver.Value{Kind: driver.Bool, Bool: parsed}, nil
	case string(dbsql.ColumnInfoTypeNameBinary):
		if value == "" {
			return driver.Value{Kind: driver.Null}, nil
		}
		return driver.Value{Kind: driver.Bytes, Bytes: []byte(value)}, nil
	case string(dbsql.ColumnInfoTypeNameShort), string(dbsql.ColumnInfoTypeNameInt):
		if value == "" {
			return driver.Value{Kind: driver.Null}, nil
		}
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return driver.Value{}, err
		}
		return driver.Value{Kind: driver.Int32, Int64: parsed}, nil
	case string(dbsql.ColumnInfoTypeNameLong):
		if value == "" {
			return driver.Value{Kind: driver.Null}, nil
		}
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return driver.Value{}, err
		}
		return driver.Value{Kind: driver.Int64, Int64: parsed}, nil
	case string(dbsql.ColumnInfoTypeNameFloat):
		if value == "" {
			return driver.Value{Kind: driver.Null}, nil
		}
		parsed, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return driver.Value{}, err
		}
		return driver.Value{Kind: driver.Float32, Float64: parsed}, nil
	case string(dbsql.ColumnInfoTypeNameDouble):
		if value == "" {
			return driver.Value{Kind: driver.Null}, nil
		}
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return driver.Value{}, err
		}
		return driver.Value{Kind: driver.Float64, Float64: parsed}, nil
	default:
		if value == "" && !strings.EqualFold(typeName, string(dbsql.ColumnInfoTypeNameString)) {
			return driver.Value{Kind: driver.Null}, nil
		}
		return driver.Value{Kind: driver.String, String: value}, nil
	}
}
