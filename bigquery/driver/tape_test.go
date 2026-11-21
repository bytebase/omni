package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"google.golang.org/api/googleapi"

	"github.com/bytebase/omni/bigquery/internal/testcontainer"
	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// openTape returns a driver on the pinned dataset whose client answers from t's
// tape, and verifies connectivity with Ping.
func openTape(t *testing.T) *Driver {
	t.Helper()
	d := openTapeDataset(t, testcontainer.Dataset)
	if err := d.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	return d
}

// openTapeDataset returns a driver on dataset whose client answers from t's
// tape without calling Ping. Replay uses testcontainer.Project; recording uses
// the live project verified against bigquery/testdata/pins.tsv.
func openTapeDataset(t *testing.T, dataset string) *Driver {
	t.Helper()
	rec := &liveRecorder{
		t:       t,
		dataset: dataset,
		jobByID: map[string]*liveJob{},
		idByJob: map[string]string{},
	}
	raw := tape.OpenRaw(t, testcontainer.Server, testcontainer.Driver, rec.do)
	d := &Driver{
		client:       &tapeClient{tape: raw},
		project:      testcontainer.Project,
		databaseName: dataset,
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

// tapeClient implements the driver's client seam over a tape.Raw stream.
type tapeClient struct {
	tape *tape.Raw
}

func (c *tapeClient) Ping(context.Context) error {
	_, err := c.call(tapeRequest{Op: opPing})
	return err
}

func (c *tapeClient) Close() error {
	return nil
}

func (c *tapeClient) ListDatasets(context.Context) ([]string, error) {
	a, err := c.call(tapeRequest{Op: opListDatasets})
	if err != nil {
		return nil, err
	}
	return a.Datasets, nil
}

func (c *tapeClient) DatasetMetadata(_ context.Context, dataset string) (*bigquery.DatasetMetadata, error) {
	a, err := c.call(tapeRequest{Op: opDatasetMetadata, Dataset: dataset})
	if err != nil {
		return nil, err
	}
	if a.Dataset == nil {
		return &bigquery.DatasetMetadata{}, nil
	}
	return &bigquery.DatasetMetadata{
		FullID:   a.Dataset.FullID,
		Location: a.Dataset.Location,
	}, nil
}

func (c *tapeClient) ListTables(_ context.Context, dataset string) ([]string, error) {
	a, err := c.call(tapeRequest{Op: opListTables, Dataset: dataset})
	if err != nil {
		return nil, err
	}
	return a.Tables, nil
}

func (c *tapeClient) TableMetadata(_ context.Context, dataset, table string) (*bigquery.TableMetadata, error) {
	a, err := c.call(tapeRequest{Op: opTableMetadata, Dataset: dataset, Table: table})
	if err != nil {
		return nil, err
	}
	if a.Table == nil {
		return &bigquery.TableMetadata{}, nil
	}
	return &bigquery.TableMetadata{
		Description: a.Table.Description,
		NumRows:     a.Table.NumRows,
		NumBytes:    a.Table.NumBytes,
	}, nil
}

func (c *tapeClient) Read(_ context.Context, dataset, statement string) (rowIterator, error) {
	a, err := c.call(tapeRequest{Op: opRead, Dataset: dataset, Statement: statement})
	if err != nil {
		return nil, err
	}
	return a.iterator(c)
}

func (c *tapeClient) Run(_ context.Context, dataset, statement string) (job, error) {
	a, err := c.call(tapeRequest{Op: opRun, Dataset: dataset, Statement: statement})
	if err != nil {
		return nil, err
	}
	return &tapeJob{client: c, id: a.JobID}, nil
}

func (c *tapeClient) DryRun(_ context.Context, dataset, statement string) (*jobStatus, error) {
	a, err := c.call(tapeRequest{Op: opDryRun, Dataset: dataset, Statement: statement})
	if err != nil {
		return nil, err
	}
	return a.Status.toJobStatus(), nil
}

func (c *tapeClient) call(r tapeRequest) (tapeAnswer, error) {
	text, err := r.text()
	if err != nil {
		return tapeAnswer{}, err
	}
	answered, err := c.tape.Do(text)
	if err != nil {
		return tapeAnswer{}, err
	}
	var a tapeAnswer
	if err := unmarshalJSON([]byte(answered), &a); err != nil {
		return tapeAnswer{}, fmt.Errorf("test: unreadable answer to %s: %w", text, err)
	}
	return a, a.Error.toError()
}

type tapeJob struct {
	client *tapeClient
	id     string
}

func (j *tapeJob) ID() string { return j.id }

func (j *tapeJob) Wait(context.Context) (*jobStatus, error) {
	a, err := j.client.call(tapeRequest{Op: opJobWait, JobID: j.id})
	if err != nil {
		return nil, err
	}
	return a.Status.toJobStatus(), nil
}

func (j *tapeJob) Read(context.Context) (rowIterator, error) {
	a, err := j.client.call(tapeRequest{Op: opJobRead, JobID: j.id})
	if err != nil {
		return nil, err
	}
	return a.iterator(j.client)
}

func (j *tapeJob) Status(context.Context) (*jobStatus, error) {
	a, err := j.client.call(tapeRequest{Op: opJobStatus, JobID: j.id})
	if err != nil {
		return nil, err
	}
	return a.Status.toJobStatus(), nil
}

type tapeRowIterator struct {
	schema    bigquery.Schema
	rows      [][]bigquery.Value
	err       error
	sourceJob job
	next      int
}

func (r *tapeRowIterator) Schema() bigquery.Schema { return r.schema }

func (r *tapeRowIterator) Next() ([]bigquery.Value, error) {
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

func (r *tapeRowIterator) SourceJob() job { return r.sourceJob }

// Operations recorded on a tape.
const (
	opPing            = "ping"
	opListDatasets    = "listDatasets"
	opDatasetMetadata = "datasetMetadata"
	opListTables      = "listTables"
	opTableMetadata   = "tableMetadata"
	opRead            = "read"
	opRun             = "run"
	opDryRun          = "dryRun"
	opJobWait         = "jobWait"
	opJobRead         = "jobRead"
	opJobStatus       = "jobStatus"
)

type tapeRequest struct {
	Op        string `json:"op"`
	Dataset   string `json:"dataset,omitempty"`
	Table     string `json:"table,omitempty"`
	Statement string `json:"statement,omitempty"`
	JobID     string `json:"jobID,omitempty"`
}

func (r tapeRequest) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type tapeAnswer struct {
	Datasets []string         `json:"datasets,omitempty"`
	Dataset  *tapeDatasetMeta `json:"dataset,omitempty"`
	Tables   []string         `json:"tables,omitempty"`
	Table    *tapeTableMeta   `json:"table,omitempty"`
	JobID    string           `json:"jobID,omitempty"`
	Schema   []tapeField      `json:"schema,omitempty"`
	Rows     [][]tapeValue    `json:"rows,omitempty"`
	RowError *tapeFailure     `json:"rowError,omitempty"`
	Status   *tapeJobStatus   `json:"status,omitempty"`
	Error    *tapeFailure     `json:"error,omitempty"`
}

func (a tapeAnswer) iterator(c *tapeClient) (rowIterator, error) {
	schema := make(bigquery.Schema, 0, len(a.Schema))
	for _, f := range a.Schema {
		schema = append(schema, &bigquery.FieldSchema{
			Name: f.Name,
			Type: bigquery.FieldType(f.Type),
		})
	}
	rows := make([][]bigquery.Value, 0, len(a.Rows))
	for _, row := range a.Rows {
		decoded := make([]bigquery.Value, 0, len(row))
		for _, v := range row {
			val, err := v.toBigQueryValue()
			if err != nil {
				return nil, err
			}
			decoded = append(decoded, val)
		}
		rows = append(rows, decoded)
	}
	var src job
	if a.JobID != "" {
		src = &tapeJob{client: c, id: a.JobID}
	}
	return &tapeRowIterator{
		schema:    schema,
		rows:      rows,
		err:       a.RowError.toError(),
		sourceJob: src,
	}, nil
}

type tapeDatasetMeta struct {
	FullID   string `json:"fullID,omitempty"`
	Location string `json:"location,omitempty"`
}

type tapeTableMeta struct {
	Description string `json:"description,omitempty"`
	NumRows     uint64 `json:"numRows,omitempty"`
	NumBytes    int64  `json:"numBytes,omitempty"`
}

type tapeField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type tapeJobStatus struct {
	Error *tapeFailure    `json:"error,omitempty"`
	Stats *tapeQueryStats `json:"stats,omitempty"`
}

func (s *tapeJobStatus) toJobStatus() *jobStatus {
	if s == nil {
		return nil
	}
	out := &jobStatus{Err: s.Error.toError()}
	if s.Stats != nil {
		out.Stats = s.Stats.toQueryStatistics()
	}
	return out
}

type tapeQueryStats struct {
	StatementType       string           `json:"statementType,omitempty"`
	NumDMLAffectedRows  int64            `json:"numDmlAffectedRows,omitempty"`
	TotalBytesProcessed int64            `json:"totalBytesProcessed,omitempty"`
	QueryPlan           []tapeQueryStage `json:"queryPlan,omitempty"`
}

func (s *tapeQueryStats) toQueryStatistics() *bigquery.QueryStatistics {
	stats := &bigquery.QueryStatistics{
		StatementType:       s.StatementType,
		NumDMLAffectedRows:  s.NumDMLAffectedRows,
		TotalBytesProcessed: s.TotalBytesProcessed,
	}
	if len(s.QueryPlan) > 0 {
		stats.QueryPlan = make([]*bigquery.ExplainQueryStage, 0, len(s.QueryPlan))
		for _, st := range s.QueryPlan {
			steps := make([]*bigquery.ExplainQueryStep, 0, len(st.Steps))
			for _, step := range st.Steps {
				steps = append(steps, &bigquery.ExplainQueryStep{
					Kind:     step.Kind,
					Substeps: step.Substeps,
				})
			}
			stats.QueryPlan = append(stats.QueryPlan, &bigquery.ExplainQueryStage{
				ID:                        st.ID,
				Name:                      st.Name,
				Status:                    st.Status,
				InputStages:               st.InputStages,
				Steps:                     steps,
				RecordsRead:               st.RecordsRead,
				RecordsWritten:            st.RecordsWritten,
				ParallelInputs:            st.ParallelInputs,
				CompletedParallelInputs:   st.CompletedParallelInputs,
				ShuffleOutputBytes:        st.ShuffleOutputBytes,
				ShuffleOutputBytesSpilled: st.ShuffleOutputBytesSpilled,
			})
		}
	}
	return stats
}

type tapeQueryStage struct {
	ID                        int64           `json:"id"`
	Name                      string          `json:"name"`
	Status                    string          `json:"status,omitempty"`
	InputStages               []int64         `json:"inputStages,omitempty"`
	Steps                     []tapeQueryStep `json:"steps,omitempty"`
	RecordsRead               int64           `json:"recordsRead,omitempty"`
	RecordsWritten            int64           `json:"recordsWritten,omitempty"`
	ParallelInputs            int64           `json:"parallelInputs,omitempty"`
	CompletedParallelInputs   int64           `json:"completedParallelInputs,omitempty"`
	ShuffleOutputBytes        int64           `json:"shuffleOutputBytes,omitempty"`
	ShuffleOutputBytesSpilled int64           `json:"shuffleOutputBytesSpilled,omitempty"`
}

type tapeQueryStep struct {
	Kind     string   `json:"kind"`
	Substeps []string `json:"substeps,omitempty"`
}

// tapeValue encodes one bigquery.Value losslessly so replay passes the exact Go
// type the BigQuery client returned into convertValue and decodeColumnRow.
type tapeValue struct {
	Kind     string        `json:"k"`
	String   string        `json:"s,omitempty"`
	Int      int64         `json:"i,omitempty"`
	Float    float64       `json:"f,omitempty"`
	Bool     bool          `json:"b,omitempty"`
	Interval *tapeInterval `json:"iv,omitempty"`
	Range    *tapeRange    `json:"rg,omitempty"`
	List     []tapeValue   `json:"l,omitempty"`
}

type tapeInterval struct {
	Months         int32 `json:"m,omitempty"`
	Days           int32 `json:"d,omitempty"`
	SubSecondNanos int32 `json:"n,omitempty"`
}

type tapeRange struct {
	Start *tapeValue `json:"start,omitempty"`
	End   *tapeValue `json:"end,omitempty"`
}

func encodeTapeValue(v bigquery.Value) (tapeValue, error) {
	switch val := v.(type) {
	case nil:
		return tapeValue{Kind: "null"}, nil
	case string:
		return tapeValue{Kind: "string", String: val}, nil
	case []byte:
		return tapeValue{Kind: "bytes", String: base64.StdEncoding.EncodeToString(val)}, nil
	case int64:
		return tapeValue{Kind: "int64", Int: val}, nil
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return tapeValue{Kind: "float64Text", String: fmt.Sprintf("%v", val)}, nil
		}
		return tapeValue{Kind: "float64", Float: val}, nil
	case bool:
		return tapeValue{Kind: "bool", Bool: val}, nil
	case *big.Rat:
		return tapeValue{Kind: "rat", String: val.RatString()}, nil
	case time.Time:
		return tapeValue{Kind: "time", String: val.UTC().Format(time.RFC3339Nano)}, nil
	case civil.Date:
		return tapeValue{Kind: "date", String: val.String()}, nil
	case civil.Time:
		return tapeValue{Kind: "civilTime", String: val.String()}, nil
	case civil.DateTime:
		return tapeValue{Kind: "dateTime", String: val.String()}, nil
	case *bigquery.IntervalValue:
		return tapeValue{
			Kind: "interval",
			Interval: &tapeInterval{
				Months:         val.Months,
				Days:           val.Days,
				SubSecondNanos: val.SubSecondNanos,
			},
		}, nil
	case *bigquery.RangeValue:
		rg := &tapeRange{}
		if val.Start != nil {
			s, err := encodeTapeValue(val.Start)
			if err != nil {
				return tapeValue{}, err
			}
			rg.Start = &s
		}
		if val.End != nil {
			e, err := encodeTapeValue(val.End)
			if err != nil {
				return tapeValue{}, err
			}
			rg.End = &e
		}
		return tapeValue{Kind: "range", Range: rg}, nil
	case []bigquery.Value:
		list := make([]tapeValue, 0, len(val))
		for _, elem := range val {
			tv, err := encodeTapeValue(elem)
			if err != nil {
				return tapeValue{}, err
			}
			list = append(list, tv)
		}
		return tapeValue{Kind: "list", List: list}, nil
	default:
		return tapeValue{}, fmt.Errorf("unsupported bigquery.Value type %T", v)
	}
}

func (v tapeValue) toBigQueryValue() (bigquery.Value, error) {
	switch v.Kind {
	case "null":
		return nil, nil
	case "string":
		return v.String, nil
	case "bytes":
		return base64.StdEncoding.DecodeString(v.String)
	case "int64":
		return v.Int, nil
	case "float64":
		return v.Float, nil
	case "float64Text":
		switch v.String {
		case "NaN":
			return math.NaN(), nil
		case "+Inf":
			return math.Inf(1), nil
		case "-Inf":
			return math.Inf(-1), nil
		default:
			return nil, fmt.Errorf("invalid float64Text %q", v.String)
		}
	case "bool":
		return v.Bool, nil
	case "rat":
		rat, ok := new(big.Rat).SetString(v.String)
		if !ok {
			return nil, fmt.Errorf("invalid rat %q", v.String)
		}
		return rat, nil
	case "time":
		return time.Parse(time.RFC3339Nano, v.String)
	case "date":
		return civil.ParseDate(v.String)
	case "civilTime":
		return civil.ParseTime(v.String)
	case "dateTime":
		return civil.ParseDateTime(v.String)
	case "interval":
		if v.Interval == nil {
			return &bigquery.IntervalValue{}, nil
		}
		return &bigquery.IntervalValue{
			Months:         v.Interval.Months,
			Days:           v.Interval.Days,
			SubSecondNanos: v.Interval.SubSecondNanos,
		}, nil
	case "range":
		out := &bigquery.RangeValue{}
		if v.Range != nil {
			if v.Range.Start != nil {
				s, err := v.Range.Start.toBigQueryValue()
				if err != nil {
					return nil, err
				}
				out.Start = s
			}
			if v.Range.End != nil {
				e, err := v.Range.End.toBigQueryValue()
				if err != nil {
					return nil, err
				}
				out.End = e
			}
		}
		return out, nil
	case "list":
		out := make([]bigquery.Value, 0, len(v.List))
		for _, elem := range v.List {
			val, err := elem.toBigQueryValue()
			if err != nil {
				return nil, err
			}
			out = append(out, val)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unknown tapeValue kind %q", v.Kind)
	}
}

// tapeFailure preserves the exact error type (*googleapi.Error or
// *bigquery.Error) so errors.As, notFoundError and Error() behave identically on
// replay.
type tapeFailure struct {
	Kind     string                `json:"kind,omitempty"`
	Code     int                   `json:"code,omitempty"`
	Reason   string                `json:"reason,omitempty"`
	Location string                `json:"location,omitempty"`
	Message  string                `json:"message,omitempty"`
	Errors   []googleapi.ErrorItem `json:"errors,omitempty"`
	Text     string                `json:"text,omitempty"`
}

func failureOf(err error) *tapeFailure {
	if err == nil {
		return nil
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		f := &tapeFailure{
			Kind:    "googleapi",
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Errors:  apiErr.Errors,
		}
		rebuilt := &googleapi.Error{Code: f.Code, Message: f.Message, Errors: f.Errors}
		if err.Error() != rebuilt.Error() {
			f.Text = err.Error()
		}
		return f
	}
	var bqErr *bigquery.Error
	if errors.As(err, &bqErr) {
		f := &tapeFailure{
			Kind:     "bigquery",
			Reason:   bqErr.Reason,
			Location: bqErr.Location,
			Message:  bqErr.Message,
		}
		rebuilt := &bigquery.Error{Reason: f.Reason, Location: f.Location, Message: f.Message}
		if err.Error() != rebuilt.Error() {
			f.Text = err.Error()
		}
		return f
	}
	return &tapeFailure{Text: err.Error()}
}

func (f *tapeFailure) toError() error {
	if f == nil {
		return nil
	}
	var inner error
	switch f.Kind {
	case "googleapi":
		inner = &googleapi.Error{
			Code:    f.Code,
			Message: f.Message,
			Errors:  f.Errors,
		}
	case "bigquery":
		inner = &bigquery.Error{
			Reason:   f.Reason,
			Location: f.Location,
			Message:  f.Message,
		}
	default:
		return errors.New(f.Text)
	}
	if f.Text != "" && f.Text != inner.Error() {
		return &wrappedTapeError{text: f.Text, cause: inner}
	}
	return inner
}

type wrappedTapeError struct {
	text  string
	cause error
}

func (e *wrappedTapeError) Error() string { return e.text }
func (e *wrappedTapeError) Unwrap() error { return e.cause }

// liveRecorder executes tape requests against the pinned live BigQuery project
// during OMNI_RECORD=1 and normalizes nondeterministic fields (job IDs, listing
// orders, dataset/table timestamps, plan stage timings) before the tape writes
// them.
type liveRecorder struct {
	t       testing.TB
	dataset string

	mu      sync.Mutex
	live    *liveClient
	nextJob int
	jobByID map[string]*liveJob
	idByJob map[string]string
}

func (r *liveRecorder) ensureLive(ctx context.Context) (*liveClient, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.live != nil {
		return r.live, nil
	}
	pin, err := enginetest.LoadPin(testcontainer.Pins, testcontainer.Driver)
	if err != nil {
		return nil, err
	}
	cfg := testcontainer.Live(r.t)
	if cfg.Dataset != testcontainer.Dataset {
		r.t.Fatalf("testcontainer: %s is %q, but recorded tests expect %q", testcontainer.DatasetVariable, cfg.Dataset, testcontainer.Dataset)
	}
	d, err := Open(ctx, driver.Config{
		Database: r.dataset,
		Password: cfg.Password,
		Params:   map[string]string{ParamProject: cfg.Project},
	})
	if err != nil {
		return nil, err
	}
	lc := d.client.(*liveClient)
	lc.disableQueryCache = true
	if err := testcontainer.VerifyPin(ctx, lc.bq, pin, cfg.Project, cfg.Dataset); err != nil {
		_ = lc.Close()
		r.t.Fatalf("testcontainer: %v", err)
	}
	r.live = lc
	r.t.Cleanup(func() { _ = lc.Close() })
	return lc, nil
}

func (r *liveRecorder) registerJob(j job) string {
	if j == nil {
		return ""
	}
	lj, ok := j.(*liveJob)
	if !ok {
		return j.ID()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if standIn, ok := r.idByJob[lj.ID()]; ok {
		return standIn
	}
	r.nextJob++
	standIn := fmt.Sprintf("job-%d", r.nextJob)
	r.idByJob[lj.ID()] = standIn
	r.jobByID[standIn] = lj
	return standIn
}

func (r *liveRecorder) lookupJob(standIn string) (*liveJob, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lj, ok := r.jobByID[standIn]
	if !ok {
		return nil, fmt.Errorf("test: unknown job %q", standIn)
	}
	return lj, nil
}

func (r *liveRecorder) do(text string) (string, error) {
	var req tapeRequest
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("test: unreadable request %s: %w", text, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	lc, err := r.ensureLive(ctx)
	if err != nil {
		return "", err
	}

	var a tapeAnswer
	switch req.Op {
	case opPing:
		err = lc.Ping(ctx)
	case opListDatasets:
		a.Datasets, err = lc.ListDatasets(ctx)
		if err == nil {
			sort.Strings(a.Datasets)
		}
	case opDatasetMetadata:
		var meta *bigquery.DatasetMetadata
		meta, err = lc.DatasetMetadata(ctx, req.Dataset)
		if err == nil && meta != nil {
			a.Dataset = &tapeDatasetMeta{
				FullID:   meta.FullID,
				Location: meta.Location,
			}
		}
	case opListTables:
		a.Tables, err = lc.ListTables(ctx, req.Dataset)
		if err == nil {
			sort.Strings(a.Tables)
		}
	case opTableMetadata:
		var meta *bigquery.TableMetadata
		meta, err = lc.TableMetadata(ctx, req.Dataset, req.Table)
		if err == nil && meta != nil {
			a.Table = &tapeTableMeta{
				Description: meta.Description,
				NumRows:     meta.NumRows,
				NumBytes:    meta.NumBytes,
			}
		}
	case opRead:
		var it rowIterator
		it, err = lc.Read(ctx, req.Dataset, req.Statement)
		if err == nil {
			var encErr error
			a, encErr = r.drainIterator(it)
			if encErr != nil {
				return "", encErr
			}
		}
	case opRun:
		var j job
		j, err = lc.Run(ctx, req.Dataset, req.Statement)
		if err == nil {
			a.JobID = r.registerJob(j)
		}
	case opDryRun:
		var status *jobStatus
		status, err = lc.DryRun(ctx, req.Dataset, req.Statement)
		if err == nil {
			a.Status = encodeJobStatus(status)
		}
	case opJobWait:
		lj, lookupErr := r.lookupJob(req.JobID)
		if lookupErr != nil {
			return "", lookupErr
		}
		var status *jobStatus
		status, err = lj.Wait(ctx)
		if err == nil {
			a.Status = encodeJobStatus(status)
		}
	case opJobRead:
		lj, lookupErr := r.lookupJob(req.JobID)
		if lookupErr != nil {
			return "", lookupErr
		}
		var it rowIterator
		it, err = lj.Read(ctx)
		if err == nil {
			var encErr error
			a, encErr = r.drainIterator(it)
			if encErr != nil {
				return "", encErr
			}
		}
	case opJobStatus:
		lj, lookupErr := r.lookupJob(req.JobID)
		if lookupErr != nil {
			return "", lookupErr
		}
		var status *jobStatus
		status, err = lj.Status(ctx)
		if err == nil {
			a.Status = encodeJobStatus(status)
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

func (r *liveRecorder) drainIterator(it rowIterator) (tapeAnswer, error) {
	var a tapeAnswer
	for {
		values, err := it.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			a.RowError = failureOf(err)
			break
		}
		encoded := make([]tapeValue, 0, len(values))
		for _, v := range values {
			tv, encErr := encodeTapeValue(v)
			if encErr != nil {
				return tapeAnswer{}, encErr
			}
			encoded = append(encoded, tv)
		}
		a.Rows = append(a.Rows, encoded)
	}
	for _, field := range it.Schema() {
		if field == nil {
			continue
		}
		a.Schema = append(a.Schema, tapeField{
			Name: field.Name,
			Type: string(field.Type),
		})
	}
	a.JobID = r.registerJob(it.SourceJob())
	return a, nil
}

func encodeJobStatus(status *jobStatus) *tapeJobStatus {
	if status == nil {
		return nil
	}
	out := &tapeJobStatus{Error: failureOf(status.Err)}
	if status.Stats != nil {
		qs := &tapeQueryStats{
			StatementType:       status.Stats.StatementType,
			NumDMLAffectedRows:  status.Stats.NumDMLAffectedRows,
			TotalBytesProcessed: status.Stats.TotalBytesProcessed,
		}
		// Timing durations (Wait/Read/Compute/Write Avg/Max) fluctuate between
		// runs, so only the deterministic structural fields of QueryPlan are
		// recorded.
		for _, stage := range status.Stats.QueryPlan {
			if stage == nil {
				continue
			}
			steps := make([]tapeQueryStep, 0, len(stage.Steps))
			for _, step := range stage.Steps {
				if step == nil {
					continue
				}
				steps = append(steps, tapeQueryStep{
					Kind:     step.Kind,
					Substeps: step.Substeps,
				})
			}
			qs.QueryPlan = append(qs.QueryPlan, tapeQueryStage{
				ID:                        stage.ID,
				Name:                      stage.Name,
				Status:                    stage.Status,
				InputStages:               stage.InputStages,
				Steps:                     steps,
				RecordsRead:               stage.RecordsRead,
				RecordsWritten:            stage.RecordsWritten,
				ParallelInputs:            stage.ParallelInputs,
				CompletedParallelInputs:   stage.CompletedParallelInputs,
				ShuffleOutputBytes:        stage.ShuffleOutputBytes,
				ShuffleOutputBytesSpilled: stage.ShuffleOutputBytesSpilled,
			})
		}
		out.Stats = qs
	}
	return out
}

func unmarshalJSON(b []byte, v any) error {
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
