package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/driver"
	omnimssql "github.com/bytebase/omni/mssql"
	mssqlanalysis "github.com/bytebase/omni/mssql/analysis"
	"github.com/bytebase/omni/mssql/ast"
)

// stmtType says what a statement of a batch returns.
type stmtType = int

const (
	stmtTypeUnknown stmtType = 1 << iota
	// stmtTypeResultSetGenerating returns exactly one result set.
	stmtTypeResultSetGenerating
	stmtTypeRowCountGenerating
	// stmtTypeProcedure returns any number of result sets and row counts.
	stmtTypeProcedure
	// stmtTypeOutput is a data modification whose OUTPUT clause returns rows.
	stmtTypeOutput
)

// Session opens one server connection held until Close, so what a statement
// sets (a session option, a temporary table) holds for the next.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &session{driver: d, conn: conn}, nil
}

type session struct {
	driver *Driver
	conn   *sql.Conn
}

func (s *session) Close() error { return s.conn.Close() }

// Query runs the script's batches in order and returns one result per
// statement. A statement that fails ends the script: its Result carries the
// error, and the error returned is nil. Errors wrapping SQL Server's own
// unwrap to gomssqldb.Error, whose text the driver unpacks to every message
// the server sent.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	var results []*driver.Result
	batch := mssqlanalysis.NewBatcher(script)
	for {
		command, err := batch.Next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return results, fmt.Errorf("failed to get next batch for statement: %s: %w", batch.Batch().Text, err)
			}
			if v := batch.Batch(); v != nil && v.Text != "" {
				r, stop, err := s.queryBatch(ctx, v.Text, opts)
				results = append(results, r...)
				if err != nil {
					return results, err
				}
				if stop {
					return results, nil
				}
			}
			break
		}
		if command == nil {
			continue
		}
		goCommand, ok := command.(*mssqlanalysis.GoCommand)
		if !ok {
			return results, fmt.Errorf("unsupported command type: %T", command)
		}
		v := batch.Batch()
		for i := uint(0); i < goCommand.Count; i++ {
			r, stop, err := s.queryBatch(ctx, v.Text, opts)
			results = append(results, r...)
			if err != nil {
				return results, err
			}
			if stop {
				return results, nil
			}
		}
		batch.Reset(nil)
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}

// queryBatch runs one batch and returns one result per statement. stop reports
// that a statement failed, which ends the script; err reports that the batch
// could not be split or classified at all.
func (s *session) queryBatch(ctx context.Context, batch string, opts driver.QueryOptions) ([]*driver.Result, bool, error) {
	statements, err := omnimssql.Parse(batch)
	if err != nil {
		return nil, true, err
	}
	var nonEmpty []omnimssql.Statement
	for _, statement := range statements {
		if !statement.Empty() {
			nonEmpty = append(nonEmpty, statement)
		}
	}
	if len(nonEmpty) == 0 {
		return nil, false, nil
	}

	if opts.Explain {
		return s.explainBatch(ctx, nonEmpty, opts)
	}
	return s.runBatch(ctx, nonEmpty, opts)
}

// explainBatch runs each statement under the session's plan output and returns
// the plan as the statement's result.
func (s *session) explainBatch(ctx context.Context, statements []omnimssql.Statement, opts driver.QueryOptions) ([]*driver.Result, bool, error) {
	showplan := showplanStatistic(opts.ExplainFormat)
	if _, err := s.conn.ExecContext(ctx, "SET "+showplan+" ON"); err != nil {
		return nil, true, fmt.Errorf("failed to enable explain mode: %w", err)
	}
	defer func() {
		// A cancelled context cannot send the option; the session is dropped with it.
		_, _ = s.conn.ExecContext(context.WithoutCancel(ctx), "SET "+showplan+" OFF")
	}()

	results := make([]*driver.Result, 0, len(statements))
	for _, statement := range statements {
		start := time.Now()
		rows, err := s.conn.QueryContext(ctx, statement.Text)
		var result *driver.Result
		if err != nil {
			result = &driver.Result{Err: err}
		} else {
			result, err = readRows(rows, opts.MaxResultBytes)
			_ = rows.Close()
			if err != nil {
				result = &driver.Result{Err: err}
			}
		}
		result.Statement = statement.Text
		result.Latency = time.Since(start)
		if result.Err != nil {
			results = append(results, result)
			return results, true, nil
		}
		results = append(results, result)
	}
	return results, false, nil
}

// showplanStatistic picks the session option that turns on the requested plan
// output. SQL Server has no JSON showplan, so anything but XML takes the
// tabular plan.
func showplanStatistic(format string) string {
	if strings.EqualFold(format, "XML") {
		return "SHOWPLAN_XML"
	}
	return "SHOWPLAN_ALL"
}

// runBatch runs the batch's statements one at a time and returns one result
// for each: a statement that returns rows yields its result sets, and one that
// only reports a count yields that count. A statement that fails ends the
// script, its Result carrying the error.
//
// SQL Server reports the affected rows of a batch's statements only through a
// collector on the request, so a statement that returns rows reports no count
// here. Bytebase's driver reads that collector; the driver contract's Result
// holds one count or one set of rows, not both.
func (s *session) runBatch(ctx context.Context, statements []omnimssql.Statement, opts driver.QueryOptions) ([]*driver.Result, bool, error) {
	results := make([]*driver.Result, 0, len(statements))
	for _, statement := range statements {
		t, err := getStmtType(statement.AST)
		if err != nil {
			return results, true, err
		}
		text := statement.Text
		if opts.Limit > 0 {
			if limited, err := omnimssql.StatementWithResultLimit(text, opts.Limit); err == nil {
				text = limited
			}
		}

		start := time.Now()
		var batch []*driver.Result
		if t&stmtTypeResultSetGenerating != 0 || t&stmtTypeProcedure != 0 {
			batch, err = s.queryStatement(ctx, text, opts)
		} else {
			batch, err = s.execStatement(ctx, text)
		}
		if err != nil {
			return append(results, &driver.Result{Statement: text, Err: unpackMSSQLError(err), Latency: time.Since(start)}), true, nil
		}
		for _, result := range batch {
			result.Statement = text
			result.Latency = time.Since(start)
			if opts.MaskingEnabled && t&(stmtTypeProcedure|stmtTypeOutput) != 0 && len(result.Columns) > 0 {
				withholdRows(result)
			}
			results = append(results, result)
		}
	}
	return results, false, nil
}

// queryStatement runs one statement that returns rows, one Result per result
// set the server sends.
func (s *session) queryStatement(ctx context.Context, statement string, opts driver.QueryOptions) ([]*driver.Result, error) {
	rows, err := s.conn.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*driver.Result
	for {
		result, err := readRows(rows, opts.MaxResultBytes)
		if err != nil {
			return results, err
		}
		results = append(results, result)
		if !rows.NextResultSet() {
			break
		}
	}
	if err := rows.Err(); err != nil {
		return results, err
	}
	return results, nil
}

// execStatement runs one statement that reports no rows, as its affected
// count.
func (s *session) execStatement(ctx context.Context, statement string) ([]*driver.Result, error) {
	sqlResult, err := s.conn.ExecContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	affected, err := sqlResult.RowsAffected()
	if err != nil {
		// A DDL statement reports no count, which the client reports as an error.
		affected = 0
	}
	return []*driver.Result{{AffectedRows: &affected}}, nil
}

// readRows reads one result set, stopping at maxResultBytes.
func readRows(rows *sql.Rows, maxResultBytes int64) (*driver.Result, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columns := make([]driver.Column, len(columnNames))
	for i, columnType := range columnTypes {
		columns[i] = driver.Column{
			Name:     columnNames[i],
			TypeName: strings.ToUpper(columnType.DatabaseTypeName()),
		}
	}
	result := &driver.Result{Columns: columns}
	totalBytes := driver.HeaderByteSize(columns)
	if len(columns) == 0 {
		return result, nil
	}
	for rows.Next() {
		dest := make([]any, len(columns))
		for i, columnType := range columnTypes {
			dest[i] = makeValueByTypeName(columnType.DatabaseTypeName())
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		row := make([]driver.Value, len(columns))
		for i, columnType := range columnTypes {
			row[i] = convertValue(columnType.DatabaseTypeName(), columnType, dest[i])
		}
		result.Rows = append(result.Rows, row)
		if maxResultBytes > 0 {
			totalBytes += driver.RowByteSize(row)
			if totalBytes > maxResultBytes {
				result.Truncated = true
				break
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// withholdRows empties a result set of a procedure call or an OUTPUT clause,
// whose rows masking cannot trace to columns.
func withholdRows(result *driver.Result) {
	if len(result.Columns) == 0 {
		return
	}
	result.Rows = nil
	// The sentence is Bytebase's own, shown to the user as it stands.
	result.Err = errors.New("Rows returned by a procedure call or an OUTPUT clause are hidden because data masking can't be applied to them.") //nolint:staticcheck // the text is the message Bytebase shows
}

// getStmtType classifies the statement's AST, which the parser already read.
func getStmtType(node ast.Node) (stmtType, error) {
	return classifyStmtType(node)
}

func classifyStmtType(node ast.Node) (stmtType, error) {
	switch n := node.(type) {
	case *ast.SelectStmt:
		if mssqlanalysis.HasSelectInto(n) {
			// SELECT ... INTO materialises a new table: no result set, only a row count.
			return stmtTypeRowCountGenerating, nil
		}
		if assignsVariables(n) {
			// SELECT @v = ... returns no result set or row count.
			return stmtTypeUnknown, nil
		}
		return stmtTypeResultSetGenerating | stmtTypeRowCountGenerating, nil

	case *ast.InsertStmt:
		return dmlStmtType(n.OutputClause), nil
	case *ast.UpdateStmt:
		return dmlStmtType(n.OutputClause), nil
	case *ast.DeleteStmt:
		return dmlStmtType(n.OutputClause), nil
	case *ast.MergeStmt:
		return dmlStmtType(n.OutputClause), nil

	case *ast.BulkInsertStmt, *ast.InsertBulkStmt, *ast.CopyIntoStmt,
		*ast.ReadtextStmt, *ast.WritetextStmt, *ast.UpdatetextStmt,
		*ast.ReceiveStmt:
		return stmtTypeRowCountGenerating, nil

	case *ast.ExecStmt:
		return stmtTypeProcedure, nil

	case *ast.SetOptionStmt:
		// One SET can list several options, as in SET ANSI_NULLS, NOEXEC ON.
		for _, option := range n.Options {
			if strings.EqualFold(option, "NOEXEC") || strings.EqualFold(option, "PARSEONLY") {
				// These skip statements without an error, so the result sets that do
				// arrive could not be matched to their statements.
				return stmtTypeUnknown, fmt.Errorf("SET %s is not supported", strings.ToUpper(option))
			}
		}
		return stmtTypeUnknown, nil

	case *ast.IfStmt, *ast.WhileStmt, *ast.BeginEndStmt, *ast.TryCatchStmt,
		*ast.ReturnStmt, *ast.BreakStmt, *ast.ContinueStmt, *ast.GotoStmt,
		*ast.LabelStmt, *ast.WaitForStmt:
		return stmtTypeUnknown, errors.New("unsupported control flow statement")

	default:
		return stmtTypeUnknown, nil
	}
}

// dmlStmtType classifies a data modification, which also returns the affected
// rows as a result set when its OUTPUT clause has no INTO.
func dmlStmtType(output *ast.OutputClause) stmtType {
	if output != nil && output.IntoTable == nil {
		return stmtTypeResultSetGenerating | stmtTypeRowCountGenerating | stmtTypeOutput
	}
	return stmtTypeRowCountGenerating
}

func assignsVariables(sel *ast.SelectStmt) bool {
	if sel.TargetList == nil {
		return false
	}
	for _, item := range sel.TargetList.Items {
		if target, ok := item.(*ast.ResTarget); ok {
			if _, ok := target.Val.(*ast.SelectAssign); ok {
				return true
			}
		}
	}
	return false
}

// unpackMSSQLError reports every message the server sent, with the fields
// Bytebase shows for one.
func unpackMSSQLError(err error) error {
	var e gomssqldb.Error
	if !errors.As(err, &e) {
		return err
	}
	if len(e.All) <= 1 {
		return errors.New(formatMSSQLError(e))
	}
	var messages []string
	for _, message := range e.All {
		if message.Message == "" {
			continue
		}
		messages = append(messages, formatMSSQLError(message))
	}
	return errors.New(strings.Join(messages, "\n"))
}

func formatMSSQLError(e gomssqldb.Error) string {
	b := new(strings.Builder)
	if len(e.ProcName) > 0 {
		fmt.Fprintf(b, "Msg %d, Level %d, State %d, Server %s, Procedure %s, Line %d\n", e.Number, e.Class, e.State, e.ServerName, e.ProcName, e.LineNo)
	} else {
		fmt.Fprintf(b, "Msg %d, Level %d, State %d, Server %s, Line %d\n", e.Number, e.Class, e.State, e.ServerName, e.LineNo)
	}
	b.WriteString(e.Message)
	return b.String()
}

// makeValueByTypeName returns the scan target for a column type.
func makeValueByTypeName(typeName string) any {
	switch typeName {
	case "UNIQUEIDENTIFIER":
		return new(gomssqldb.NullUniqueIdentifier)
	case "TINYINT", "SMALLINT", "INT", "BIGINT":
		return new(sql.NullInt64)
	case "VARCHAR", "NVARCHAR", "CHAR", "NCHAR", "TEXT", "NTEXT", "XML":
		return new(sql.NullString)
	case "REAL", "FLOAT":
		return new(sql.NullFloat64)
	case "BIT":
		// Every SQL editor shows 0/1 for BIT, so the value converts to an integer.
		return new(sql.NullBool)
	case "DECIMAL", "SMALLMONEY", "MONEY":
		return new(sql.NullString)
	case "TIME", "DATE", "SMALLDATETIME", "DATETIME", "DATETIME2", "DATETIMEOFFSET":
		return new(sql.NullTime)
	case "VARBINARY", "IMAGE", "BINARY", "SQL_VARIANT", "GEOMETRY", "GEOGRAPHY":
		return new([]byte)
	default:
		// An unknown type is read as text, which most values survive.
		return new(sql.NullString)
	}
}

// convertValue turns one scanned value into the contract's value.
func convertValue(typeName string, columnType *sql.ColumnType, value any) driver.Value {
	switch raw := value.(type) {
	case *sql.NullString:
		if raw.Valid {
			return driver.Value{Kind: driver.String, String: raw.String}
		}
	case *sql.NullInt64:
		if raw.Valid {
			return driver.Value{Kind: driver.Int64, Int64: raw.Int64}
		}
	case *[]byte:
		if *raw != nil {
			return driver.Value{Kind: driver.Bytes, Bytes: *raw}
		}
	case *sql.NullBool:
		if raw.Valid {
			var v int64
			if raw.Bool {
				v = 1
			}
			return driver.Value{Kind: driver.Int64, Int64: v}
		}
	case *sql.NullFloat64:
		if raw.Valid {
			return driver.Value{Kind: driver.Float64, Float64: raw.Float64}
		}
	case *gomssqldb.NullUniqueIdentifier:
		if raw.Valid {
			return driver.Value{Kind: driver.String, String: raw.UUID.String()}
		}
	case *sql.NullTime:
		if raw.Valid {
			switch columnType.DatabaseTypeName() {
			case "TIME":
				return driver.Value{Kind: driver.String, String: raw.Time.Format(time.TimeOnly)}
			case "DATE":
				return driver.Value{Kind: driver.String, String: raw.Time.Format(time.DateOnly)}
			}
			// The client reports a temporal column's fractional-seconds precision as
			// its scale, and nothing when the column declares none.
			_, scale, _ := columnType.DecimalSize()
			if typeName == "DATETIME" || typeName == "DATETIME2" || typeName == "SMALLDATETIME" {
				return driver.Value{Kind: driver.Timestamp, Time: raw.Time, Precision: int32(scale)}
			}
			return driver.Value{Kind: driver.TimestampTZ, Time: raw.Time, Precision: int32(scale)}
		}
	}
	return driver.Value{Kind: driver.Null}
}
