package driver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/driver"
	omnipg "github.com/bytebase/omni/pg"
	pganalysis "github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/ast"
)

// Session opens one server connection held until Close, so session state such
// as search_path, TimeZone and temporary tables persists across Query calls.
func (d *Driver) Session(ctx context.Context) (driver.Session, error) {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	s := &session{driver: d, conn: conn}
	_ = conn.Raw(func(driverConn any) error {
		sc, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return nil
		}
		s.pgConn = sc.Conn().PgConn()
		return nil
	})
	if s.pgConn != nil {
		d.noticesMu.Lock()
		if d.notices == nil {
			d.notices = make(map[*pgconn.PgConn]*[]driver.Message)
		}
		d.notices[s.pgConn] = &s.messages
		d.noticesMu.Unlock()
	}
	return s, nil
}

type session struct {
	driver   *Driver
	conn     *sql.Conn
	pgConn   *pgconn.PgConn
	messages []driver.Message
}

func (s *session) Close() error {
	if s.pgConn != nil {
		s.driver.noticesMu.Lock()
		delete(s.driver.notices, s.pgConn)
		s.driver.noticesMu.Unlock()
	}
	return s.conn.Close()
}

func (s *session) takeMessages() []driver.Message {
	if s.pgConn == nil {
		return nil
	}
	s.driver.noticesMu.Lock()
	defer s.driver.noticesMu.Unlock()
	msgs := s.messages
	s.messages = nil
	return msgs
}

// Query runs each non-empty statement of script in order on the session's
// connection. Errors wrapping PostgreSQL server errors unwrap to
// *pgconn.PgError. A script the parser cannot read is the call's error and no
// statement of it runs: every statement is parsed before the first is sent.
func (s *session) Query(ctx context.Context, script string, opts driver.QueryOptions) ([]*driver.Result, error) {
	segments := omnipg.Split(script)
	var statements []string
	for _, seg := range segments {
		if !seg.Empty() {
			statements = append(statements, seg.Text)
		}
	}
	if len(statements) == 0 {
		return nil, nil
	}

	type preparedStmt struct {
		text        string
		returnsData bool
	}
	prepared := make([]preparedStmt, 0, len(statements))
	for _, statement := range statements {
		if opts.Explain {
			explained, err := pganalysis.ExplainStatement(statement, strings.ToLower(opts.ExplainFormat))
			if err != nil {
				return nil, err
			}
			statement = explained
		} else if opts.Limit > 0 {
			if rewritten, err := pganalysis.RewriteResultLimit(statement, opts.Limit); err == nil {
				statement = rewritten
			} else {
				statement = fmt.Sprintf("WITH result AS (\n%s\n) SELECT * FROM result LIMIT %d;", trimStatement(statement), opts.Limit)
			}
		}

		parsed, err := omnipg.Parse(statement)
		if err != nil {
			return nil, err
		}
		nodes := make([]ast.Node, 0, len(parsed))
		for _, stmt := range parsed {
			nodes = append(nodes, stmt.AST)
		}
		_, returnsData := pganalysis.ValidateQuery(nodes)
		prepared = append(prepared, preparedStmt{text: statement, returnsData: returnsData})
	}

	var results []*driver.Result
	for _, stmt := range prepared {
		// Keep public after the selected schema because PostgreSQL extensions
		// commonly install helper functions there and call them unqualified.
		if opts.Schema != "" {
			if _, err := s.conn.ExecContext(ctx, fmt.Sprintf("SET search_path TO %s, public;", quoteIdent(opts.Schema))); err != nil {
				return nil, err
			}
		}

		_ = s.takeMessages()
		startTime := time.Now()
		res, err := s.runStatement(ctx, stmt.text, stmt.returnsData, opts.MaxResultBytes)
		stop := false
		if err != nil {
			_ = s.takeMessages()
			res = &driver.Result{Err: err}
			stop = true
		}
		res.Statement = stmt.text
		res.Latency = time.Since(startTime)
		results = append(results, res)
		if stop {
			break
		}
	}
	return results, nil
}

func (s *session) runStatement(ctx context.Context, statement string, returnsData bool, maxResultBytes int64) (*driver.Result, error) {
	if returnsData {
		location := pgSessionLocation(s.conn)
		rows, err := s.conn.QueryContext(ctx, statement)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		res, err := readRows(rows, location, maxResultBytes)
		if err != nil {
			return nil, err
		}
		res.Messages = s.takeMessages()
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if format, executed, ok := pganalysis.DescribeExplain(statement); ok {
			res.Plan = &driver.Plan{
				Format:   strings.ToUpper(format),
				Executed: executed,
			}
		}
		return res, nil
	}

	// database/sql's Result states a count for every statement, so it cannot
	// tell a statement that reported none (CREATE TABLE, SET) from one that
	// reported 0 (UPDATE 0). Read PostgreSQL's command tag through pgx, as the
	// executor does; a tape's connection has no pgx underneath, so fall back to
	// database/sql there.
	var result *driver.Result
	usedPgx := false
	err := s.conn.Raw(func(driverConn any) error {
		sc, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return nil
		}
		usedPgx = true
		mrr := sc.Conn().PgConn().Exec(ctx, statement)
		result = &driver.Result{}
		var affectedRows int64
		var reported bool
		for mrr.NextResult() {
			// Close drains any rows without retaining them: a DML statement with
			// RETURNING, or a SELECT with a data-modifying CTE, is classified as
			// not returning data, but can still answer with a large result set.
			tag, err := mrr.ResultReader().Close()
			if err != nil {
				_ = mrr.Close()
				return err
			}
			if n, ok := commandTagRowsAffected(tag); ok {
				affectedRows += n
				reported = true
			}
		}
		if err := mrr.Close(); err != nil {
			return err
		}
		if reported {
			result.AffectedRows = &affectedRows
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if usedPgx {
		result.Messages = s.takeMessages()
		return result, nil
	}

	sqlResult, err := s.conn.ExecContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	res := &driver.Result{Messages: s.takeMessages()}
	if affectedRows, err := sqlResult.RowsAffected(); err == nil {
		res.AffectedRows = &affectedRows
	}
	return res, nil
}

// commandTagRowsAffected returns the row count a command tag carries, and false
// when it carries none. PostgreSQL ends a counted tag (INSERT 0 3, UPDATE 3,
// SELECT 3, COPY 3) with the count and a tag with no count (CREATE TABLE,
// ALTER TABLE, TRUNCATE TABLE, SET) with a word.
func commandTagRowsAffected(tag pgconn.CommandTag) (int64, bool) {
	s := tag.String()
	if s == "" || s[len(s)-1] < '0' || s[len(s)-1] > '9' {
		return 0, false
	}
	return tag.RowsAffected(), true
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
				dest[i] = makeValueByTypeName(col.TypeName)
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
	case "DATE":
		return new(pgtype.Date)
	case "TIMESTAMP", "TIMESTAMPTZ":
		return new(pgtype.Timestamptz)
	case "BIT", "VARBIT", "BYTEA", "GEOMETRY", "GEOGRAPHY":
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
			dbType := columnType.DatabaseTypeName()
			if dbType == "TIME" || dbType == timeTzOID || dbType == "INTERVAL" {
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

// FormatTimestamp renders t in PostgreSQL's text representation: a space
// separator, no zone for timestamp without time zone, and a " BC" suffix for
// years before 1 AD (Go year 0 is 1 BC).
func FormatTimestamp(t time.Time, withZone bool) string {
	year := t.Year()
	suffix := ""
	if year <= 0 {
		year = 1 - year
		suffix = " BC"
	}
	s := fmt.Sprintf("%04d-%s", year, t.Format("01-02 15:04:05.999999999"))
	if withZone {
		s += t.Format("-07:00")
	}
	return s + suffix
}

// padZeroes pads the fractional seconds of a TIME, TIMETZ or INTERVAL text
// value to acc digits, since PostgreSQL omits trailing fractional zeros.
func padZeroes(rawStr string, acc int) string {
	dotIndex := strings.Index(rawStr, ".")
	if dotIndex < 0 {
		return rawStr
	}
	endIndex := len(rawStr)
	if plusIndex := strings.Index(rawStr, "+"); plusIndex >= 0 {
		endIndex = plusIndex
	} else if minusIndex := strings.Index(rawStr, "-"); minusIndex >= 0 {
		if minusIndex > dotIndex {
			endIndex = minusIndex
		}
	}
	if endIndex <= dotIndex {
		return rawStr
	}
	decimalPart := rawStr[dotIndex+1 : endIndex]
	if len(decimalPart) < acc {
		rawStr = rawStr[:endIndex] + strings.Repeat("0", acc-len(decimalPart)) + rawStr[endIndex:]
	}
	return rawStr
}

var pgLocations sync.Map

// pgSessionLocation reads the session's TimeZone parameter immediately before
// a statement runs so every cell of its result set renders in one location;
// POSIX specifications Go's zoneinfo rejects fall back to UTC.
func pgSessionLocation(conn *sql.Conn) *time.Location {
	location := time.UTC
	_ = conn.Raw(func(driverConn any) error {
		pgxConn, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return nil
		}
		location = pgLocation(pgxConn.Conn().PgConn().ParameterStatus("TimeZone"))
		return nil
	})
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

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}), unicode.IsSpace)
}
