package driver

import (
	"bytes"
	"compress/zlib"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	crrawparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	crrawparsertree "github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/driver"
	omnipg "github.com/bytebase/omni/pg"
	pganalysis "github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/ast"
)

// Session opens one server connection held until Close, so session state
// such as search_path and TimeZone persists across Query calls.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &session{conn: conn}, nil
}

type session struct {
	conn *sql.Conn
}

func (s *session) Close() error {
	return s.conn.Close()
}

// Query runs each non-empty statement of script in order on the session's
// connection. Errors wrapping CockroachDB server errors unwrap to
// *pgconn.PgError. A script the parser cannot read is the call's error and no
// statement of it runs: every statement is parsed before the first is sent.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	split, err := splitStatements(script)
	if err != nil {
		return nil, err
	}
	if len(split) == 0 {
		return nil, nil
	}

	type preparedStmt struct {
		text        string
		returnsData bool
		planAsJSON  bool
		plan        *driver.Plan
	}
	prepared := make([]preparedStmt, 0, len(split))
	for _, single := range split {
		statement := single.Text
		planAsJSON := false
		if opts.Explain {
			explained, err := pganalysis.CockroachExplainStatement(statement, strings.EqualFold(opts.ExplainFormat, "JSON"))
			if err != nil {
				return nil, err
			}
			planAsJSON = strings.EqualFold(opts.ExplainFormat, "JSON")
			statement = explained
		} else if opts.Limit > 0 {
			statement = statementWithResultLimit(statement, opts.Limit)
		}

		_, returnsData, err := validateQuery(statement)
		if err != nil {
			return nil, err
		}
		var plan *driver.Plan
		if opts.Explain {
			format := "TEXT"
			if planAsJSON {
				format = "JSON"
			}
			plan = &driver.Plan{Format: format}
		} else if format, executed, ok := describeExplain(statement); ok {
			plan = &driver.Plan{Format: format, Executed: executed}
		}
		prepared = append(prepared, preparedStmt{text: statement, returnsData: returnsData, planAsJSON: planAsJSON, plan: plan})
	}

	var results []*driver.Result
	for _, stmt := range prepared {
		// Sanitize the schema name by escaping any quotes.
		if opts.Schema != "" {
			safeSchemaName := strings.ReplaceAll(opts.Schema, `"`, `""`)
			if err := crdb.Execute(func() error {
				_, err := s.conn.ExecContext(ctx, fmt.Sprintf(`SET search_path TO "%s";`, safeSchemaName)) // NOSONAR(go:S2077) safeSchemaName is sanitized by escaping double quotes above
				return err
			}); err != nil {
				return nil, err
			}
		}

		startTime := time.Now()
		queryResult, err := func() (*driver.Result, error) {
			if stmt.returnsData {
				location := pgSessionLocation(ctx, s.conn)
				var r *driver.Result
				if err := crdb.Execute(func() error {
					rows, err := s.conn.QueryContext(ctx, stmt.text)
					if err != nil {
						return err
					}
					defer rows.Close()
					r, err = readRows(rows, location, opts.MaxResultBytes)
					if err != nil {
						return err
					}
					return rows.Err()
				}); err != nil {
					return nil, err
				}
				if stmt.planAsJSON {
					return planDiagram(r)
				}
				return r, nil
			}

			var sqlResult sql.Result
			if err := crdb.Execute(func() error {
				var err error
				sqlResult, err = s.conn.ExecContext(ctx, stmt.text)
				return err
			}); err != nil {
				return nil, err
			}
			res := &driver.Result{}
			if affectedRows, err := sqlResult.RowsAffected(); err == nil {
				res.AffectedRows = &affectedRows
			}
			return res, nil
		}()
		stop := false
		if err != nil {
			queryResult = &driver.Result{Err: err}
			stop = true
		}
		queryResult.Statement = stmt.text
		if queryResult.Err == nil {
			queryResult.Plan = stmt.plan
		}
		queryResult.Latency = time.Since(startTime)
		results = append(results, queryResult)
		if stop {
			break
		}
	}

	return results, nil
}

// validateQuery reports whether the statement returns rows. CockroachDB's
// own parser decides, so its native syntax (SHOW, a cast like 1:::INT8) and a
// SELECT over a mutating CTE are read the same way the server reads them.
// readOnly is Bytebase's editor question and comes from the PostgreSQL
// grammar when that grammar accepts the statement.
func validateQuery(statement string) (readOnly, returnsData bool, err error) {
	parsed, err := crrawparser.Parse(statement)
	if err != nil {
		return false, false, err
	}
	for _, stmt := range parsed {
		node, ok := stmt.AST.(crrawparsertree.Statement)
		if !ok {
			continue
		}
		if node.StatementReturnType() == crrawparsertree.Rows {
			returnsData = true
		}
	}
	if stmts, perr := omnipg.Parse(statement); perr == nil {
		nodes := make([]ast.Node, 0, len(stmts))
		for _, stmt := range stmts {
			nodes = append(nodes, stmt.AST)
		}
		readOnly, _ = pganalysis.ValidateQuery(nodes)
	}
	return readOnly, returnsData, nil
}

// describeExplain reports a directly submitted EXPLAIN's format and whether
// it executed the statement it plans.
func describeExplain(statement string) (format string, executed bool, ok bool) {
	parsed, err := crrawparser.Parse(statement)
	if err != nil || len(parsed) != 1 {
		return "", false, false
	}
	switch node := parsed[0].AST.(type) {
	case *crrawparsertree.Explain:
		return explainFormat(node.ExplainOptions), false, true
	case *crrawparsertree.ExplainAnalyze:
		return explainFormat(node.ExplainOptions), true, true
	}
	return "", false, false
}

func explainFormat(opts crrawparsertree.ExplainOptions) string {
	if opts.Flags[crrawparsertree.ExplainFlagJSON] {
		return "JSON"
	}
	return "TEXT"
}

func statementWithResultLimit(statement string, limit int) string {
	if rewritten, err := pganalysis.RewriteResultLimit(statement, limit); err == nil {
		return rewritten
	}
	return fmt.Sprintf("WITH result AS (\n%s\n) SELECT * FROM result LIMIT %d;", trimStatement(statement), limit)
}

func readRows(rows *sql.Rows, location *time.Location, maxResultBytes int64) (*driver.Result, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columns := make([]driver.Column, len(columnNames))
	for i, ct := range columnTypes {
		columns[i] = driver.Column{
			Name:     columnNames[i],
			TypeName: strings.ToUpper(ct.DatabaseTypeName()),
		}
	}
	res := &driver.Result{Columns: columns}
	totalBytes := driver.HeaderByteSize(columns)

	if len(columns) > 0 {
		for rows.Next() {
			dest := make([]any, len(columns))
			for i, col := range columns {
				dest[i] = makeValueByTypeName(col.TypeName, columnTypes[i])
			}
			if err := rows.Scan(dest...); err != nil {
				return nil, err
			}
			row := make([]driver.Value, len(columns))
			for i, col := range columns {
				row[i] = convertValue(col.TypeName, columnTypes[i], dest[i], location)
			}
			res.Rows = append(res.Rows, row)
			if maxResultBytes > 0 {
				totalBytes += driver.RowByteSize(row)
				if totalBytes > maxResultBytes {
					res.Truncated = true
					break
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func makeValueByTypeName(typeName string, _ *sql.ColumnType) any {
	switch typeName {
	case "VARCHAR", "TEXT", "UUID":
		return new(sql.NullString)
	case "BOOL":
		return new(sql.NullBool)
	case "INT", "INTEGER", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "INT2", "INT4", "INT8":
		return new(sql.NullInt64)
	case "FLOAT", "DOUBLE", "FLOAT4", "FLOAT8":
		return new(sql.NullFloat64)
	case "DATE":
		return new(pgtype.Date)
	case "TIMESTAMP", "TIMESTAMPTZ":
		return new(pgtype.Timestamptz)
	case "BIT", "VARBIT", "BYTEA":
		return new([]byte)
	default:
		return new(sql.NullString)
	}
}

var timeTzOID = fmt.Sprintf("%d", pgtype.TimetzOID)

func convertValue(typeName string, columnType *sql.ColumnType, value any, location *time.Location) driver.Value {
	switch raw := value.(type) {
	case *sql.NullString:
		if raw.Valid {
			// TODO: Fix DatabaseTypeName for 1266, Object ID for TIME WITHOUT TIME ZONE
			if columnType.DatabaseTypeName() == "TIME" || columnType.DatabaseTypeName() == timeTzOID || columnType.DatabaseTypeName() == "INTERVAL" {
				return driver.Value{Kind: driver.String, String: padZeroes(raw.String, 6)}
			}
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
			return driver.Value{Kind: driver.Bool, Bool: raw.Bool}
		}
	case *sql.NullFloat64:
		if raw.Valid {
			return driver.Value{Kind: driver.Float64, Float64: raw.Float64}
		}
	case *pgtype.Date:
		if raw.Valid {
			if raw.InfinityModifier != pgtype.Finite {
				return driver.Value{Kind: driver.String, String: raw.InfinityModifier.String()}
			}
			if typeName == "DATE" {
				return driver.Value{Kind: driver.String, String: raw.Time.Format(time.DateOnly)}
			}
		}
	case *pgtype.Timestamptz:
		if raw.Valid {
			if raw.InfinityModifier != pgtype.Finite {
				return driver.Value{Kind: driver.String, String: raw.InfinityModifier.String()}
			}
			// pgx v5.13+ reports a temporal column's fractional-seconds precision as its
			// scale, an unspecified precision already normalized to PostgreSQL's default
			// of 6 (jackc/pgx#2704).
			_, scale, ok := columnType.DecimalSize()
			if !ok {
				scale = 6
			}
			if typeName == "TIMESTAMP" {
				return driver.Value{Kind: driver.Timestamp, Time: raw.Time, Precision: int32(scale)}
			}
			return driver.Value{Kind: driver.TimestampTZ, Time: raw.Time.In(location), Precision: int32(scale)}
		}
	}
	return driver.Value{Kind: driver.Null}
}

// Padding 0's to nanosecond precision to make sure it's always 6 digits.
// Since the data cannot be formatted into a time.Time, we need to pad it here.
// The value arrives as trimmed text (PostgreSQL drops trailing fractional zeros), so the
// accuracy to pad to is supplied by the caller rather than derived from the value.
func padZeroes(rawStr string, acc int) string {
	dotIndex := strings.Index(rawStr, ".")
	if dotIndex < 0 {
		return rawStr
	}
	// End index is used to cut off the time zone information.
	endIndex := len(rawStr)
	if plusIndex := strings.Index(rawStr[dotIndex+1:], "+"); plusIndex >= 0 {
		endIndex = dotIndex + 1 + plusIndex
	} else if minusIndex := strings.Index(rawStr[dotIndex+1:], "-"); minusIndex >= 0 {
		endIndex = dotIndex + 1 + minusIndex
	}
	decimalPart := rawStr[dotIndex+1 : endIndex]
	if len(decimalPart) < acc {
		rawStr = rawStr[:endIndex] + strings.Repeat("0", acc-len(decimalPart)) + rawStr[endIndex:]
	}
	return rawStr
}

// diagramLink is the line EXPLAIN (DISTSQL) prints after the text plan. What
// follows the fragment marker is the plan itself, deflated and base64'd.
var diagramLink = regexp.MustCompile(`https?://\S*#([A-Za-z0-9+/_=-]+)`)

// planDiagram returns the physical plan EXPLAIN (DISTSQL) encodes in the link
// it prints, as the JSON document the query-plan viewer parses.
//
// CockroachDB publishes no plan as a document: every EXPLAIN option of its own
// returns text, and the one machine-readable form it has is the payload behind
// this link, which its own diagram page decodes in the browser. Decoding it
// here is what lets the viewer draw the plan without the browser fetching
// anything, and it is the same move the Oracle driver makes in rendering the
// plan table.
func planDiagram(result *driver.Result) (*driver.Result, error) {
	var payload string
	for _, row := range result.Rows {
		for _, value := range row {
			if value.Kind == driver.String {
				if match := diagramLink.FindStringSubmatch(value.String); match != nil {
					payload = match[1]
				}
			}
		}
	}
	if payload == "" {
		return nil, fmt.Errorf("the plan carried no diagram to read")
	}
	// CockroachDB encodes with the standard alphabet, but the payload rides in
	// a URL fragment, so accept either spelling of the two bytes that differ.
	normalized := strings.NewReplacer("-", "+", "_", "/").Replace(strings.TrimRight(payload, "="))
	deflated, err := base64.RawStdEncoding.DecodeString(normalized)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the query plan: %w", err)
	}
	reader, err := zlib.NewReader(bytes.NewReader(deflated))
	if err != nil {
		return nil, fmt.Errorf("failed to read the query plan: %w", err)
	}
	defer reader.Close()
	// A plan is small; the bound keeps a corrupt payload from inflating
	// without end.
	plan, err := io.ReadAll(io.LimitReader(reader, maxPlanBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to read the query plan: %w", err)
	}

	return &driver.Result{
		Columns: []driver.Column{{Name: "plan", TypeName: "TEXT"}},
		Rows:    [][]driver.Value{{{Kind: driver.String, String: string(plan)}}},
	}, nil
}

// maxPlanBytes bounds the inflated plan. CockroachDB's diagrams run to a few
// tens of kilobytes; this is far above that and far below anything that would
// trouble the editor.
const maxPlanBytes = 8 << 20

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}), unicode.IsSpace)
}

var pgLocations sync.Map

// pgSessionLocation reads the session's TimeZone parameter immediately before
// a statement runs so every cell of its result set renders in one location;
// POSIX specifications Go's zoneinfo rejects fall back to UTC. A connector
// that does not expose pgx, such as the tape, is asked with SHOW timezone.
func pgSessionLocation(ctx context.Context, conn *sql.Conn) *time.Location {
	location := time.UTC
	known := false
	_ = conn.Raw(func(driverConn any) error {
		pgxConn, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return nil
		}
		location = pgLocation(pgxConn.Conn().PgConn().ParameterStatus("TimeZone"))
		known = true
		return nil
	})
	if known {
		return location
	}
	var name string
	if err := conn.QueryRowContext(ctx, "SHOW timezone").Scan(&name); err == nil {
		return pgLocation(name)
	}
	return location
}

func pgLocation(name string) *time.Location {
	if name == "" || name == "UTC" {
		return time.UTC
	}
	if cached, ok := pgLocations.Load(name); ok {
		if location, ok := cached.(*time.Location); ok {
			return location
		}
	}
	location, err := time.LoadLocation(name)
	if err != nil {
		location = time.UTC
	}
	pgLocations.Store(name, location)
	return location
}
