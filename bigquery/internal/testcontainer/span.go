package testcontainer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Span is the pin role bigquery/analysis's span.tsv and its tape were recorded
// at: the project and location of the dataset a recording builds.
const Span = "span"

// SpanDataset is the dataset a span recording builds each fixture in and deletes
// when it is done. The fixtures and queries name it, so it is fixed.
const SpanDataset = "omni_span"

// spanLabel marks the dataset as the span recording's own: a recording deletes
// a dataset of that name it finds only when it carries the label.
const spanLabel = "omni"

// spanRequest is one request on the span tape.
type spanRequest struct {
	Op  string `json:"op"` // "fixture" or "query"
	SQL string `json:"sql"`
}

// A SpanAnswer is BigQuery's answer to a span request.
type SpanAnswer struct {
	// Error is the reason BigQuery states for refusing a query
	// (invalidQuery, notFound, ...). Empty for an answer.
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`

	// Relations is a fixture's stored schema: every table and view the dataset
	// lists, in name order, with its columns as the table's schema states them.
	Relations []SpanRelation `json:"relations,omitempty"`

	// Columns is a dry-run query's output field names, as its schema states.
	Columns []string `json:"columns,omitempty"`
	// Tables is the tables and views the dry run reports the query references
	// (referencedTables and referencedLogicalViews), as dataset.table, sorted.
	Tables []string `json:"tables,omitempty"`
}

// A SpanRelation is one table or view of a fixture, with its top-level
// columns in order. Their types are not kept: the analysis reads none.
type SpanRelation struct {
	Name    string   `json:"name"`
	View    bool     `json:"view,omitempty"`
	Columns []string `json:"columns"`
}

// SpanOracle answers span requests: from t's tape in a replay, and from BigQuery
// itself in a recording.
type SpanOracle struct {
	raw *tape.Raw
}

// OpenSpan returns the span oracle on t's tape. A recording resolves
// BIGQUERY_CREDENTIALS and BIGQUERY_PROJECT, checks the project and location
// against the span pin, and deletes SpanDataset when t ends.
func OpenSpan(t testing.TB) *SpanOracle {
	t.Helper()
	var live func(string) (string, error)
	if enginetest.Recording(t) {
		live = openSpanLive(t).do
	}
	return &SpanOracle{raw: tape.OpenRaw(t, Server, Span, live)}
}

// Fixture builds a fixture from its DDL in a fresh SpanDataset and returns the
// relations it holds.
func (o *SpanOracle) Fixture(ddl string) (SpanAnswer, error) {
	a, err := o.ask(spanRequest{Op: "fixture", SQL: ddl})
	if err == nil && a.Error != "" {
		err = fmt.Errorf("fixture: %s: %s", a.Error, a.Message)
	}
	return a, err
}

// Query dry-runs query on the last fixture built, with SpanDataset as the
// default dataset. A query BigQuery refuses is an answer with Error set.
func (o *SpanOracle) Query(query string) (SpanAnswer, error) {
	return o.ask(spanRequest{Op: "query", SQL: query})
}

func (o *SpanOracle) ask(req spanRequest) (SpanAnswer, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return SpanAnswer{}, err
	}
	text, err := o.raw.Do(string(b))
	if err != nil {
		return SpanAnswer{}, err
	}
	var a SpanAnswer
	if err := json.Unmarshal([]byte(text), &a); err != nil {
		return SpanAnswer{}, fmt.Errorf("decode the span answer %q: %w", text, err)
	}
	return a, nil
}

// spanLive is a recording's connection to the pinned project.
type spanLive struct {
	project string
	client  *bigquery.Client
	http    *http.Client
}

func openSpanLive(t testing.TB) *spanLive {
	t.Helper()
	cfg := DryRun(t)
	pin, err := enginetest.LoadPin(Pins, Span)
	if err != nil {
		t.Fatalf("%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	client, err := DryRunClient(ctx, cfg)
	if err != nil {
		t.Fatalf("bigquery span: create client: %v", err)
	}
	opts := []option.ClientOption{option.WithScopes(bigquery.Scope)}
	if cfg.Password != "" {
		cred, err := credentialOption([]byte(cfg.Password))
		if err != nil {
			t.Fatalf("bigquery span: %v", err)
		}
		opts = append(opts, cred)
	}
	hc, _, err := htransport.NewClient(context.Background(), opts...)
	if err != nil {
		t.Fatalf("bigquery span: create HTTP client: %v", err)
	}
	if err := VerifyDryRunPin(ctx, client, pin, cfg.Project); err != nil {
		t.Fatalf("bigquery span: %v", err)
	}
	l := &spanLive{project: cfg.Project, client: client, http: hc}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := l.dropDataset(ctx); err != nil {
			t.Errorf("bigquery span: delete dataset %s: %v", SpanDataset, err)
		}
		_ = client.Close()
	})
	return l
}

func (l *spanLive) do(text string) (string, error) {
	var req spanRequest
	if err := json.Unmarshal([]byte(text), &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var (
		a   SpanAnswer
		err error
	)
	switch req.Op {
	case "fixture":
		a, err = l.fixture(ctx, req.SQL)
	case "query":
		a, err = l.dryRun(ctx, req.SQL)
	default:
		err = fmt.Errorf("unknown span request %q", req.Op)
	}
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(a)
	return string(b), err
}

// dropDataset deletes SpanDataset and everything in it, if it is there. A
// dataset of that name without the span label is not this recording's, and is
// left alone.
func (l *spanLive) dropDataset(ctx context.Context) error {
	ds := l.client.Dataset(SpanDataset)
	meta, err := ds.Metadata(ctx)
	if isNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if meta.Labels[spanLabel] != Span {
		return fmt.Errorf("dataset %s exists without label %s=%s, so it is not the span recording's: delete it by hand", SpanDataset, spanLabel, Span)
	}
	return ds.DeleteWithContents(ctx)
}

// fixture builds ddl in a fresh SpanDataset and reads back its relations as
// bigquery/driver's sync lists them: every table and view of the dataset, in
// name order, with its top-level columns in order.
func (l *spanLive) fixture(ctx context.Context, ddl string) (SpanAnswer, error) {
	if err := l.dropDataset(ctx); err != nil {
		return SpanAnswer{}, err
	}
	ds := l.client.Dataset(SpanDataset)
	if err := ds.Create(ctx, &bigquery.DatasetMetadata{
		Location:               "US",
		Labels:                 map[string]string{spanLabel: Span},
		DefaultTableExpiration: 24 * time.Hour,
	}); err != nil {
		return SpanAnswer{}, fmt.Errorf("create dataset %s: %w", SpanDataset, err)
	}
	q := l.client.Query(ddl)
	q.Location = "US"
	job, err := q.Run(ctx)
	if err != nil {
		return SpanAnswer{Error: reason(err), Message: err.Error()}, nil
	}
	st, err := job.Wait(ctx)
	if err == nil {
		err = st.Err()
	}
	if err != nil {
		return SpanAnswer{Error: reason(err), Message: err.Error()}, nil
	}
	var names []string
	it := ds.Tables(ctx)
	for {
		tb, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return SpanAnswer{}, err
		}
		names = append(names, tb.TableID)
	}
	slices.Sort(names)
	var a SpanAnswer
	for _, name := range names {
		meta, err := ds.Table(name).Metadata(ctx)
		if err != nil {
			return SpanAnswer{}, err
		}
		rel := SpanRelation{Name: name, View: meta.Type == bigquery.ViewTable}
		for _, f := range meta.Schema {
			rel.Columns = append(rel.Columns, f.Name)
		}
		a.Relations = append(a.Relations, rel)
	}
	return a, nil
}

// dryRun submits query as a dry-run job, which reads no data and costs
// nothing, and returns the output fields and referenced tables it states. The
// REST answer is read whole because the Go client drops
// referencedLogicalViews, the views a query reads.
func (l *spanLive) dryRun(ctx context.Context, query string) (SpanAnswer, error) {
	body, err := json.Marshal(map[string]any{
		"configuration": map[string]any{
			"dryRun": true,
			"query": map[string]any{
				"query":          query,
				"useLegacySql":   false,
				"defaultDataset": map[string]string{"projectId": l.project, "datasetId": SpanDataset},
			},
		},
		"jobReference": map[string]string{"location": "US"},
	})
	if err != nil {
		return SpanAnswer{}, err
	}
	url := "https://bigquery.googleapis.com/bigquery/v2/projects/" + l.project + "/jobs"
	backoff := time.Second
	for attempt := 1; ; attempt++ {
		a, retry, err := l.dryRunOnce(ctx, url, body)
		if retry && attempt < 6 {
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
		return a, err
	}
}

type tableRef struct {
	ProjectID string `json:"projectId"`
	DatasetID string `json:"datasetId"`
	TableID   string `json:"tableId"`
}

func (l *spanLive) dryRunOnce(ctx context.Context, url string, body []byte) (SpanAnswer, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return SpanAnswer{}, false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := l.http.Do(req)
	if err != nil {
		return SpanAnswer{}, true, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpanAnswer{}, true, err
	}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return SpanAnswer{}, true, fmt.Errorf("dry run: HTTP %d: %s", resp.StatusCode, raw)
	}
	var job struct {
		Error *struct {
			Message string `json:"message"`
			Errors  []struct {
				Reason string `json:"reason"`
			} `json:"errors"`
		} `json:"error"`
		Statistics struct {
			Query struct {
				Schema struct {
					Fields []struct {
						Name string `json:"name"`
					} `json:"fields"`
				} `json:"schema"`
				ReferencedTables       []tableRef `json:"referencedTables"`
				ReferencedLogicalViews []tableRef `json:"referencedLogicalViews"`
			} `json:"query"`
		} `json:"statistics"`
	}
	if err := json.Unmarshal(raw, &job); err != nil {
		return SpanAnswer{}, false, fmt.Errorf("dry run: decode %q: %w", raw, err)
	}
	if job.Error != nil {
		if resp.StatusCode < 400 || len(job.Error.Errors) == 0 {
			return SpanAnswer{}, false, fmt.Errorf("dry run: HTTP %d: %s", resp.StatusCode, raw)
		}
		return SpanAnswer{Error: job.Error.Errors[0].Reason, Message: job.Error.Message}, false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return SpanAnswer{}, false, fmt.Errorf("dry run: HTTP %d: %s", resp.StatusCode, raw)
	}
	a := SpanAnswer{Columns: []string{}, Tables: []string{}}
	for _, f := range job.Statistics.Query.Schema.Fields {
		a.Columns = append(a.Columns, f.Name)
	}
	q := job.Statistics.Query
	for _, r := range append(q.ReferencedTables, q.ReferencedLogicalViews...) {
		// A table read at a point in time (FOR SYSTEM_TIME AS OF) is named
		// with the snapshot it read, table@<ms>, which moves every run.
		table, _, _ := strings.Cut(r.TableID, "@")
		a.Tables = append(a.Tables, r.DatasetID+"."+table)
	}
	slices.Sort(a.Tables)
	a.Tables = slices.Compact(a.Tables)
	return a, false, nil
}

func isNotFound(err error) bool {
	var gerr *googleapi.Error
	return errors.As(err, &gerr) && gerr.Code == http.StatusNotFound
}

// reason is the first reason a BigQuery error states, or "error".
func reason(err error) string {
	var gerr *googleapi.Error
	if errors.As(err, &gerr) && len(gerr.Errors) > 0 && gerr.Errors[0].Reason != "" {
		return gerr.Errors[0].Reason
	}
	var berr *bigquery.Error
	if errors.As(err, &berr) && berr.Reason != "" {
		return berr.Reason
	}
	return "error"
}
