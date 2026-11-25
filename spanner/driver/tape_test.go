package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// openDriver returns a driver over database on instance, whose client is t's
// tape: a replay of its tape, or, in a recording, the pinned emulator with every
// call written down. The database itself is created by the test's own fixture.
func openDriver(t testing.TB, instance, database string) *Driver {
	t.Helper()
	return &Driver{
		client:       OpenTape(t, instance, database),
		databaseName: database,
		instancePath: testcontainer.InstancePath(instance),
	}
}

// OpenTape returns t's client to database on instance: a replay of its tape or,
// in a recording, the emulator itself with every call written down.
func OpenTape(t testing.TB, instance, database string) *tapeClient {
	t.Helper()
	raw := tape.OpenRaw(t, testcontainer.Server, testcontainer.Default, func(request string) (string, error) {
		// A recording reaches the emulator on its first call, so a test that
		// replays starts no container.
		return live(t, instance, database, request)
	})
	return &tapeClient{tape: raw}
}

// tapeClient is the driver's client, answered by t's tape. It is what a driver
// test hands the driver in place of a live client; its methods are the calls the
// driver makes, and the recordings hold what a live client answered.
type tapeClient struct {
	tape *tape.Raw
}

func (c *tapeClient) Ping(context.Context) error {
	_, err := c.call(request{Op: opPing})
	return err
}

func (c *tapeClient) Close() error {
	_, err := c.call(request{Op: opClose})
	return err
}

func (c *tapeClient) ListDatabases(context.Context) ([]Database, error) {
	a, err := c.call(request{Op: opListDatabases})
	if err != nil {
		return nil, err
	}
	return a.Databases, nil
}

func (c *tapeClient) GetDatabase(_ context.Context, database string) (Database, error) {
	a, err := c.call(request{Op: opGetDatabase, Database: database})
	if err != nil {
		return Database{}, err
	}
	return Database{Name: database, Dialect: a.Dialect}, nil
}

func (c *tapeClient) UpdateDatabaseDdl(_ context.Context, database string, statements []string) error {
	_, err := c.call(request{Op: opUpdateDatabaseDdl, Database: database, Statements: statements})
	return err
}

func (c *tapeClient) CreateDatabase(_ context.Context, instancePath, createStatement string, extra []string) error {
	_, err := c.call(request{Op: opCreateDatabase, InstancePath: instancePath, CreateStatement: createStatement, ExtraStatements: extra})
	return err
}

func (c *tapeClient) GetDatabaseDdl(_ context.Context, database string) ([]string, error) {
	a, err := c.call(request{Op: opGetDatabaseDdl, Database: database})
	if err != nil {
		return nil, err
	}
	return a.Ddl, nil
}

func (c *tapeClient) Read(_ context.Context, statement string) (RowIterator, error) {
	a, err := c.call(request{Op: opRead, Statement: statement})
	if err != nil {
		return nil, err
	}
	return a.iterator(), nil
}

func (c *tapeClient) Snapshot(context.Context) (Reader, error) {
	return &tapeReader{client: c}, nil
}

func (c *tapeClient) Update(_ context.Context, statement string) (int64, error) {
	a, err := c.call(request{Op: opUpdate, Statement: statement})
	if err != nil {
		return 0, err
	}
	return a.Count, nil
}

func (c *tapeClient) BatchUpdate(_ context.Context, statements []string) ([]int64, error) {
	a, err := c.call(request{Op: opBatchUpdate, Statements: statements})
	if err != nil {
		return nil, err
	}
	return a.Counts, nil
}

func (c *tapeClient) PartitionedUpdate(_ context.Context, statement string) (int64, error) {
	a, err := c.call(request{Op: opPartitionedUpdate, Statement: statement})
	if err != nil {
		return 0, err
	}
	return a.Count, nil
}

func (c *tapeClient) AnalyzeQuery(_ context.Context, statement string, readWrite bool) (*sppb.QueryPlan, error) {
	a, err := c.call(request{Op: opAnalyzeQuery, Statement: statement, ReadWrite: readWrite})
	if err != nil {
		return nil, err
	}
	if a.Plan == nil {
		return nil, nil
	}
	return a.Plan.plan, nil
}

// call sends one call and returns its answer.
func (c *tapeClient) call(r request) (answer, error) {
	text, err := r.text()
	if err != nil {
		return answer{}, err
	}
	answered, err := c.tape.Do(text)
	if err != nil {
		return answer{}, err
	}
	var a answer
	if err := unmarshalJSON([]byte(answered), &a); err != nil {
		return answer{}, fmt.Errorf("test: unreadable answer to %s: %w", text, err)
	}
	return a, a.Error.err()
}

// tapeReader is a snapshot whose queries are t's tape.
type tapeReader struct {
	client *tapeClient
}

func (r *tapeReader) Query(_ context.Context, statement string) (RowIterator, error) {
	a, err := r.client.call(request{Op: opQuery, Statement: statement})
	if err != nil {
		return nil, err
	}
	return a.iterator(), nil
}

func (*tapeReader) Close() {}

// The calls the driver makes of its client, as a tape names them.
const (
	opPing              = "ping"
	opClose             = "close"
	opListDatabases     = "listDatabases"
	opGetDatabase       = "getDatabase"
	opUpdateDatabaseDdl = "updateDatabaseDdl"
	opCreateDatabase    = "createDatabase"
	opGetDatabaseDdl    = "getDatabaseDdl"
	opRead              = "read"
	opQuery             = "query"
	opUpdate            = "update"
	opBatchUpdate       = "batchUpdate"
	opPartitionedUpdate = "partitionedUpdate"
	opAnalyzeQuery      = "analyzeQuery"
)

// request is one call as a tape holds it.
type request struct {
	Op              string   `json:"op"`
	Database        string   `json:"database,omitempty"`
	InstancePath    string   `json:"instancePath,omitempty"`
	Statements      []string `json:"statements,omitempty"`
	CreateStatement string   `json:"createStatement,omitempty"`
	ExtraStatements []string `json:"extraStatements,omitempty"`
	Statement       string   `json:"statement,omitempty"`
	ReadWrite       bool     `json:"readWrite,omitempty"`
}

func (r request) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// answer is one call's answer as a tape holds it.
type answer struct {
	Databases []Database `json:"databases,omitempty"`
	// Dialect is what a database read stated.
	Dialect string   `json:"dialect,omitempty"`
	Ddl     []string `json:"ddl,omitempty"`
	Count   int64    `json:"count,omitempty"`
	Counts  []int64  `json:"counts,omitempty"`
	Columns []column `json:"columns,omitempty"`
	Rows    [][]cell `json:"rows,omitempty"`
	Plan    *plan    `json:"plan,omitempty"`
	Error   *failure `json:"error,omitempty"`
}

// iterator is an answer to a read as the seam's row iterator: the columns it
// declares, the rows it read, and the error the read ended with, after them.
func (a answer) iterator() RowIterator {
	columns := make([]Column, 0, len(a.Columns))
	for _, c := range a.Columns {
		columns = append(columns, Column{Name: c.Name, Type: c.Type.spannerType})
	}
	rows := make([][]Cell, 0, len(a.Rows))
	for _, row := range a.Rows {
		cells := make([]Cell, 0, len(row))
		for _, c := range row {
			cells = append(cells, Cell{Type: c.Type.spannerType, Value: c.Value.spannerValue})
		}
		rows = append(rows, cells)
	}
	return &tapeRowIterator{columns: columns, rows: rows, err: a.Error.err()}
}

// tapeRowIterator is a read's answer, read on demand.
type tapeRowIterator struct {
	columns []Column
	rows    [][]Cell
	// err is what the read failed with after the rows it returned, if it did.
	err  error
	next int
}

func (r *tapeRowIterator) Columns() []Column { return r.columns }

func (r *tapeRowIterator) Next() ([]Cell, error) {
	if r.next < len(r.rows) {
		row := r.rows[r.next]
		r.next++
		return row, nil
	}
	if r.err != nil {
		return nil, r.err
	}
	return nil, io.EOF
}

func (r *tapeRowIterator) Stop() {}

// column is one result column as a tape holds it.
type column struct {
	Name string      `json:"name"`
	Type spannerType `json:"type"`
}

// cell is one value as a tape holds it.
type cell struct {
	Type  spannerType  `json:"type"`
	Value spannerValue `json:"value"`
}

// spannerType is a Spanner type as a tape holds it: the proto's own JSON form,
// so a recording holds the type the client stated rather than what a decoder
// made of it.
type spannerType struct {
	spannerType *sppb.Type
}

func (t spannerType) MarshalJSON() ([]byte, error) {
	if t.spannerType == nil {
		return []byte("null"), nil
	}
	return protojson.Marshal(t.spannerType)
}

func (t *spannerType) UnmarshalJSON(b []byte) error {
	t.spannerType = &sppb.Type{}
	return protojson.Unmarshal(b, t.spannerType)
}

// spannerValue is a cell's value as a tape holds it: the proto's own JSON form.
type spannerValue struct {
	spannerValue *structpb.Value
}

func (v spannerValue) MarshalJSON() ([]byte, error) {
	if v.spannerValue == nil {
		return []byte("null"), nil
	}
	return protojson.Marshal(v.spannerValue)
}

func (v *spannerValue) UnmarshalJSON(b []byte) error {
	v.spannerValue = &structpb.Value{}
	return protojson.Unmarshal(b, v.spannerValue)
}

// plan is a query plan as a tape holds it: the proto's own JSON form.
type plan struct {
	plan *sppb.QueryPlan
}

func (p plan) MarshalJSON() ([]byte, error) {
	if p.plan == nil {
		return []byte("null"), nil
	}
	return protojson.Marshal(p.plan)
}

func (p *plan) UnmarshalJSON(b []byte) error {
	p.plan = &sppb.QueryPlan{}
	return protojson.Unmarshal(b, p.plan)
}

// failure is an error the client returned, as a tape holds it: the gRPC code it
// stated, when it stated one, and the text the driver reports.
type failure struct {
	// Code is the gRPC code's number, absent for an error that stated none.
	Code int32  `json:"code,omitempty"`
	Text string `json:"text"`
}

func failureOf(err error) *failure {
	if err == nil {
		return nil
	}
	if s, ok := status.FromError(err); ok {
		return &failure{Code: int32(s.Code()), Text: err.Error()}
	}
	return &failure{Text: err.Error()}
}

func (f *failure) err() error {
	if f == nil {
		return nil
	}
	if f.Code == 0 {
		return errors.New(f.Text)
	}
	return &tapeError{code: codes.Code(f.Code), text: f.Text}
}

// tapeError is a recorded client error: the text the driver reports and the gRPC
// code a caller reads off it. The client's own error type is the SDK's, which a
// tape cannot hold.
type tapeError struct {
	code codes.Code
	text string
}

func (e *tapeError) Error() string { return e.text }

func (e *tapeError) GRPCStatus() *status.Status { return status.New(e.code, e.text) }

// live answers one request from the pinned emulator, through the client the
// driver itself opens. It skips Open, which refuses a PostgreSQL-dialect
// database, so a test can record that refusal.
func live(t testing.TB, instance, database, text string) (string, error) {
	var req request
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("test: unreadable request %s: %w", text, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	c, err := newClient(ctx, driver.Config{
		Host:     testcontainer.Host(t),
		Database: database,
		Params: map[string]string{
			ParamProject:  testcontainer.Project,
			ParamInstance: instance,
			ParamEmulator: "1",
		},
	})
	if err != nil {
		return "", err
	}
	defer func() { _ = c.Close() }()

	a := answer{}
	switch req.Op {
	case opPing:
		err = c.Ping(ctx)
	case opClose:
		err = c.Close()
	case opListDatabases:
		a.Databases, err = c.ListDatabases(ctx)
	case opGetDatabase:
		var database Database
		database, err = c.GetDatabase(ctx, req.Database)
		a.Dialect = database.Dialect
	case opUpdateDatabaseDdl:
		err = c.UpdateDatabaseDdl(ctx, req.Database, req.Statements)
	case opCreateDatabase:
		err = c.CreateDatabase(ctx, req.InstancePath, req.CreateStatement, req.ExtraStatements)
	case opGetDatabaseDdl:
		a.Ddl, err = c.GetDatabaseDdl(ctx, req.Database)
	case opRead:
		var iter RowIterator
		iter, err = c.Read(ctx, req.Statement)
		if err == nil {
			a, err = readAnswer(iter, a)
		}
	case opQuery:
		var reader Reader
		reader, err = c.Snapshot(ctx)
		if err == nil {
			defer reader.Close()
			var iter RowIterator
			iter, err = reader.Query(ctx, req.Statement)
			if err == nil {
				a, err = readAnswer(iter, a)
			}
		}
	case opUpdate:
		a.Count, err = c.Update(ctx, req.Statement)
	case opBatchUpdate:
		a.Counts, err = c.BatchUpdate(ctx, req.Statements)
	case opPartitionedUpdate:
		a.Count, err = c.PartitionedUpdate(ctx, req.Statement)
	case opAnalyzeQuery:
		var queryPlan *sppb.QueryPlan
		queryPlan, err = c.AnalyzeQuery(ctx, req.Statement, req.ReadWrite)
		if queryPlan != nil {
			a.Plan = &plan{plan: queryPlan}
		}
	default:
		return "", fmt.Errorf("test: unknown request %q", req.Op)
	}
	if err != nil {
		a.Error = failureOf(err)
	}
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// readAnswer reads a result into the answer a tape holds: its columns, its rows,
// and the error the read ended with. A read that fails before any row is the
// answer's error alone.
func readAnswer(iter RowIterator, a answer) (answer, error) {
	defer iter.Stop()
	cells, err := iter.Next()
	if errors.Is(err, io.EOF) {
		a.Columns = tapeColumns(iter.Columns())
		return a, nil
	}
	if err != nil {
		a.Error = failureOf(err)
		return a, nil
	}
	a.Columns = tapeColumns(iter.Columns())
	for {
		a.Rows = append(a.Rows, tapeCells(cells))
		cells, err = iter.Next()
		if errors.Is(err, io.EOF) {
			return a, nil
		}
		if err != nil {
			a.Error = failureOf(err)
			return a, nil
		}
	}
}

func tapeColumns(columns []Column) []column {
	taped := make([]column, 0, len(columns))
	for _, c := range columns {
		taped = append(taped, column{Name: c.Name, Type: spannerType{spannerType: c.Type}})
	}
	return taped
}

func tapeCells(cells []Cell) []cell {
	taped := make([]cell, 0, len(cells))
	for _, c := range cells {
		taped = append(taped, cell{Type: spannerType{spannerType: c.Type}, Value: spannerValue{spannerValue: c.Value}})
	}
	return taped
}

// unmarshalJSON decodes JSON strictly enough that a tape's typo fails the test
// rather than answering something else.
func unmarshalJSON(b []byte, v any) error {
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
