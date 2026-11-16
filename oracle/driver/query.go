package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	omnidriver "github.com/bytebase/omni/driver"
	omnioracle "github.com/bytebase/omni/oracle"
	oracleanalysis "github.com/bytebase/omni/oracle/analysis"
	"github.com/bytebase/omni/oracle/ast"
)

// Query runs each non-empty statement of script in order on the session's
// connection. Errors wrapping server errors unwrap to the client's own error
// type. Masking is not applied: an Oracle query result carries no masking
// request of its own, and the caller masks.
func (s *session) Query(ctx context.Context, script string, opts omnidriver.QueryOptions) ([]*omnidriver.Result, error) {
	parsed, err := omnioracle.Parse(script)
	if err != nil {
		return nil, err
	}
	type preparedStmt struct {
		text        string
		returnsData bool
	}
	var prepared []preparedStmt
	for _, statement := range parsed {
		if statement.Empty() {
			continue
		}
		text := statement.Text
		if !opts.Explain && opts.Limit > 0 {
			text = limitStatement(text, opts.Limit, s.driver.engineVersion)
		}
		// The statement's shape decides how it runs, not the read-only
		// validator: a SELECT that declares PL/SQL inline is refused read-only,
		// and still returns its rows.
		returnsData := false
		switch statement.AST.(type) {
		case *ast.SelectStmt, *ast.ExplainPlanStmt:
			returnsData = true
		}
		prepared = append(prepared, preparedStmt{text: text, returnsData: returnsData})
	}
	if len(prepared) == 0 {
		return nil, nil
	}

	var results []*omnidriver.Result
	for _, stmt := range prepared {
		// The session is the caller's, so the schema is set on every statement
		// rather than once: a script that switches it switches for the next.
		if opts.Schema != "" {
			if _, err := s.conn.ExecContext(ctx, fmt.Sprintf("ALTER SESSION SET CURRENT_SCHEMA = %s", quoteIdent(opts.Schema))); err != nil {
				return nil, err
			}
			s.schemaApplied = true
		}
		startTime := time.Now()
		var res *omnidriver.Result
		var err error
		statement := stmt.text
		if opts.Explain {
			res, statement, err = s.explain(ctx, stmt.text, opts)
		} else {
			res, err = s.runStatement(ctx, stmt.text, stmt.returnsData, false, opts.MaxResultBytes)
		}
		stop := false
		if err != nil {
			res = &omnidriver.Result{Err: err}
			stop = true
		}
		res.Statement = statement
		res.Latency = time.Since(startTime)
		results = append(results, res)
		if stop {
			break
		}
	}
	return results, nil
}

// limitStatement bounds the rows statement returns. The rewrite reads the
// statement's first keyword, so a leading comment is set aside and put back
// around the bounded statement; a statement neither form can bound keeps its
// text, as it does in Bytebase.
func limitStatement(statement string, limit int, engineVersion string) string {
	start := skipSpaceAndComments(statement)
	prefix, body := statement[:start], statement[start:]
	if rewritten, err := oracleanalysis.RewriteResultLimit(body, limit, engineVersion); err == nil && rewritten != body {
		return prefix + rewritten
	}
	return prefix + oracleanalysis.RowNumResultLimit(body, limit)
}

// explain plans statement and returns the plan, in the requested format, with
// the text that read it.
func (s *session) explain(ctx context.Context, statement string, opts omnidriver.QueryOptions) (*omnidriver.Result, string, error) {
	statementID := s.driver.nextStatementID()
	if _, err := s.conn.ExecContext(ctx, fmt.Sprintf("EXPLAIN PLAN SET STATEMENT_ID = '%s' FOR %s", statementID, statement)); err != nil {
		return nil, statement, err
	}
	asJSON := isJSONExplainFormat(opts.ExplainFormat)
	planQuery := fmt.Sprintf(textPlanQuery, statementID, statementID)
	if asJSON {
		planQuery = fmt.Sprintf(planStepsQuery, statementID)
	}
	res, err := s.runStatement(ctx, planQuery, true, asJSON, opts.MaxResultBytes)
	if err != nil {
		return nil, planQuery, err
	}
	// Oracle plans a statement without running it.
	format := "TEXT"
	if asJSON {
		format = "JSON"
	}
	res.Plan = &omnidriver.Plan{Format: format}
	return res, planQuery, nil
}

// runStatement runs one statement and reads its rows, or its affected count.
func (s *session) runStatement(ctx context.Context, statement string, returnsData, planDocument bool, maxResultBytes int64) (*omnidriver.Result, error) {
	if !returnsData {
		sqlResult, err := s.conn.ExecContext(ctx, statement)
		if err != nil {
			return nil, err
		}
		// A DDL statement reports no affected rows.
		affected, err := sqlResult.RowsAffected()
		if err != nil {
			affected = 0
		}
		return &omnidriver.Result{AffectedRows: &affected}, nil
	}
	rows, err := s.conn.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	if planDocument {
		res, err := readPlanSteps(rows)
		if err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return res, nil
	}
	res, localZones, err := readRows(rows, maxResultBytes)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// A TIMESTAMP WITH LOCAL TIME ZONE comes back as the instant it names, and
	// SQL*Plus prints that instant in the session's zone. Reading the zone needs
	// the statement's rows closed, so the conversion happens here.
	if len(localZones) > 0 {
		// The connection holds one result set at a time.
		if err := rows.Close(); err != nil {
			return nil, err
		}
		location, err := s.localTimeZone(ctx)
		if err != nil {
			return nil, err
		}
		for _, cell := range localZones {
			inZone := res.Rows[cell.row][cell.column].Time.In(location)
			res.Rows[cell.row][cell.column].Time = time.Date(
				inZone.Year(), inZone.Month(), inZone.Day(),
				inZone.Hour(), inZone.Minute(), inZone.Second(), inZone.Nanosecond(),
				time.UTC,
			)
		}
	}
	return res, nil
}

// localTimeZone returns the zone the session prints a TIMESTAMP WITH LOCAL TIME
// ZONE in, as SESSIONTIMEZONE states it: an offset such as +02:00, or a region
// name such as Europe/Berlin.
func (s *session) localTimeZone(ctx context.Context) (*time.Location, error) {
	var zone string
	if err := s.conn.QueryRowContext(ctx, "SELECT SESSIONTIMEZONE FROM DUAL").Scan(&zone); err != nil {
		return nil, fmt.Errorf("failed to read the session time zone: %w", err)
	}
	return oracleLocation(zone), nil
}

var oracleLocations sync.Map

// oracleLocation reads SESSIONTIMEZONE, which names a zone as an offset or as a
// region such as Europe/Berlin, and leaves a name the system cannot resolve in
// UTC.
func oracleLocation(zone string) *time.Location {
	if zone == "" || zone == "UTC" {
		return time.UTC
	}
	if cached, ok := oracleLocations.Load(zone); ok {
		if location, ok := cached.(*time.Location); ok {
			return location
		}
	}
	var location *time.Location
	if offset, ok := zoneOffset(zone); ok {
		location = time.FixedZone(zone, offset)
	} else if loaded, err := time.LoadLocation(zone); err == nil {
		location = loaded
	} else {
		location = time.UTC
	}
	oracleLocations.Store(zone, location)
	return location
}

// zoneOffset reads SESSIONTIMEZONE's offset spelling, hours and minutes east of
// UTC, with an optional fractional part on the minutes.
func zoneOffset(zone string) (int, bool) {
	if len(zone) < 3 || (zone[0] != '+' && zone[0] != '-') {
		return 0, false
	}
	var hours, minutes int
	if _, err := fmt.Sscanf(zone[1:], "%d:%d", &hours, &minutes); err != nil {
		return 0, false
	}
	offset := hours*3600 + minutes*60
	if zone[0] == '-' {
		offset = -offset
	}
	return offset, true
}

// isJSONExplainFormat reports whether format asks for the plan as a document.
// Oracle states its plan as rows, as SQL*Plus prints it, so the default and
// TEXT keep those rows and only JSON renders them as one document.
func isJSONExplainFormat(format string) bool {
	return strings.EqualFold(format, "JSON")
}

// planStepsQuery reads the plan Oracle recorded for one EXPLAIN PLAN: a row per
// operator, each naming its parent, which is how Oracle states a plan. The
// columns are what the query-plan viewer reads — the tree's shape, what each
// step does, the object it reads, and what the optimizer expects it to cost.
const planStepsQuery = `SELECT id, parent_id, position, operation, options, object_owner, object_name, object_type, optimizer, cost, cardinality, bytes, partition_start, partition_stop, access_predicates, filter_predicates FROM plan_table WHERE statement_id = '%s' ORDER BY id`

// textPlanQuery is the plan as Oracle prints it: one indented row per operator.
const textPlanQuery = `SELECT LPAD(' ', LEVEL-1) || OPERATION || ' (' || OPTIONS || ')' "Operation", OBJECT_NAME "Object", OPTIMIZER "Optimizer", COST "Cost", CARDINALITY "Cardinality", BYTES "Bytes", PARTITION_START "Partition Start", PARTITION_ID "Partition ID", ACCESS_PREDICATES "Access Predicates" FROM PLAN_TABLE START WITH ID = 0 AND statement_id = '%s' CONNECT BY PRIOR ID=PARENT_ID AND statement_id = '%s' ORDER BY id`

// readPlanSteps returns the plan rows as one JSON array, which is the form the
// query-plan viewer parses; the viewer rebuilds the tree from the parent links.
//
// Oracle renders a plan as rows rather than as a document, and no EXPLAIN option
// renders it as one, so this is where the rendering happens. It is built here
// rather than with Oracle's own JSON_ARRAYAGG so that it works on releases older
// than 12.2, which do not have those functions. A tab-separated grid would not
// do: a predicate can hold a tab or a newline, and nothing downstream could tell
// that apart from the next column or the next row.
func readPlanSteps(rows *sql.Rows) (*omnidriver.Result, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var steps [][]any
	for rows.Next() {
		values := make([]any, len(columns))
		scan := make([]any, len(columns))
		for i := range values {
			scan[i] = &values[i]
		}
		if err := rows.Scan(scan...); err != nil {
			return nil, err
		}
		steps = append(steps, values)
	}
	payload, err := planJSON(columns, steps)
	if err != nil {
		return nil, err
	}
	return &omnidriver.Result{
		Columns: []omnidriver.Column{{Name: "PLAN", TypeName: "CLOB"}},
		Rows:    [][]omnidriver.Value{{{Kind: omnidriver.String, String: payload}}},
	}, nil
}

// localZoneCell names one value a TIMESTAMP WITH LOCAL TIME ZONE column
// returned, which the session's zone renders.
type localZoneCell struct {
	row    int
	column int
}

// readRows reads one result set, stopping at maxResultBytes, and reports the
// values the session's time zone renders.
func readRows(rows *sql.Rows, maxResultBytes int64) (*omnidriver.Result, []localZoneCell, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}
	columns := make([]omnidriver.Column, len(columnNames))
	for i, columnType := range columnTypes {
		columns[i] = omnidriver.Column{
			Name:     columnNames[i],
			TypeName: strings.ToUpper(columnType.DatabaseTypeName()),
		}
	}
	res := &omnidriver.Result{Columns: columns}
	if len(columns) == 0 {
		return res, nil, nil
	}
	totalBytes := omnidriver.HeaderByteSize(columns)
	var localZones []localZoneCell
	for rows.Next() {
		dest := make([]any, len(columns))
		for i := range columns {
			dest[i] = makeValueByTypeName(columns[i].TypeName)
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, nil, err
		}
		row := make([]omnidriver.Value, len(columns))
		for i := range columns {
			row[i] = convertValue(columns[i].TypeName, columnTypes[i], dest[i])
			if columns[i].TypeName == "TIMESTAMPLTZ_DTY" && row[i].Kind == omnidriver.Timestamp {
				localZones = append(localZones, localZoneCell{row: len(res.Rows), column: i})
			}
		}
		res.Rows = append(res.Rows, row)
		if maxResultBytes > 0 {
			totalBytes += omnidriver.RowByteSize(row)
			if totalBytes > maxResultBytes {
				res.Truncated = true
				break
			}
		}
	}
	return res, localZones, nil
}

// makeValueByTypeName creates appropriate Go types for Oracle column types.
//
// DATE: date. TIMESTAMPDTY: timestamp. TIMESTAMPTZ_DTY: timestamp with time
// zone. TIMESTAMPLTZ_DTY: timezone with local time zone.
func makeValueByTypeName(typeName string) any {
	switch typeName {
	case "VARCHAR", "TEXT", "UUID":
		return new(sql.NullString)
	case "BOOL":
		return new(sql.NullBool)
	case "INT", "INTEGER", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "INT2", "INT4", "INT8":
		return new(sql.NullInt64)
	case "FLOAT", "DOUBLE", "FLOAT4", "FLOAT8":
		return new(sql.NullFloat64)
	case "BIT", "VARBIT":
		return new([]byte)
	case "DATE", "TIMESTAMPDTY", "TIMESTAMPLTZ_DTY", "TIMESTAMPTZ_DTY":
		return new(sql.NullTime)
	default:
		return new(sql.NullString)
	}
}

// convertValue turns one scanned value into the contract's value.
func convertValue(typeName string, columnType *sql.ColumnType, value any) omnidriver.Value {
	switch raw := value.(type) {
	case *sql.NullString:
		if raw.Valid {
			return omnidriver.Value{Kind: omnidriver.String, String: raw.String}
		}
	case *sql.NullInt64:
		if raw.Valid {
			return omnidriver.Value{Kind: omnidriver.Int64, Int64: raw.Int64}
		}
	case *[]byte:
		if len(*raw) > 0 {
			return omnidriver.Value{Kind: omnidriver.Bytes, Bytes: *raw}
		}
	case *sql.NullBool:
		if raw.Valid {
			return omnidriver.Value{Kind: omnidriver.Bool, Bool: raw.Bool}
		}
	case *sql.NullFloat64:
		if raw.Valid {
			return omnidriver.Value{Kind: omnidriver.Float64, Float64: raw.Float64}
		}
	case *sql.NullTime:
		if raw.Valid {
			return convertTimestamp(typeName, columnType, raw.Time)
		}
	}
	return omnidriver.Value{Kind: omnidriver.Null}
}

// convertTimestamp handles Oracle timestamp type conversions. The client
// retrieves the database timezone from the wire protocol and appends it to the
// timestamp, so DATE and TIMESTAMP keep their wall clock and only TIMESTAMP WITH
// TIME ZONE carries a zone.
// https://github.com/sijms/go-ora/blob/2962e725e7a756a667a546fb360ef09afd4c8bd0/v2/parameter.go#L616
func convertTimestamp(typeName string, columnType *sql.ColumnType, t time.Time) omnidriver.Value {
	_, scale, _ := columnType.DecimalSize()
	switch typeName {
	case "DATE", "TIMESTAMPDTY":
		stripped := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
		return omnidriver.Value{Kind: omnidriver.Timestamp, Time: stripped, Precision: int32(scale)}
	case "TIMESTAMPLTZ_DTY":
		// The client states the value as the instant it names; the session's zone
		// renders it, in the pass that runs once the statement's rows are closed.
		return omnidriver.Value{Kind: omnidriver.Timestamp, Time: t.UTC(), Precision: int32(scale)}
	default:
		// The client returns the value in the database's zone, which is the
		// zone the session renders it in.
		return omnidriver.Value{Kind: omnidriver.TimestampTZ, Time: t, Precision: int32(scale)}
	}
}

// planJSON renders plan rows as a JSON array, one object per operator keyed by
// the lower-cased column names, leaving out what the optimizer had nothing to
// say about.
func planJSON(columns []string, rows [][]any) (string, error) {
	steps := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		step := make(map[string]any, len(columns))
		for i, column := range columns {
			if i >= len(row) {
				break
			}
			if value := jsonValue(row[i]); value != nil {
				step[strings.ToLower(column)] = value
			}
		}
		steps = append(steps, step)
	}
	payload, err := json.Marshal(steps)
	if err != nil {
		return "", fmt.Errorf("failed to render the query plan: %w", err)
	}
	return string(payload), nil
}

// jsonValue renders one scanned column, or nil where the row has no value for
// it. The driver hands back text as bytes, which JSON would otherwise encode as
// base64.
func jsonValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case []byte:
		return string(typed)
	default:
		return typed
	}
}
