package driver

import (
	"context"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// client is what the driver asks of BigQuery: connectivity, dataset and table
// listings and metadata, query reads, command and script jobs, and dry runs.
// It is the seam a test replaces with a tape so the driver's own code runs
// during both recording and offline replay; liveClient below is the thin
// adapter to *bigquery.Client.
type client interface {
	Ping(ctx context.Context) error
	Close() error
	ListDatasets(ctx context.Context) ([]string, error)
	DatasetMetadata(ctx context.Context, dataset string) (*bigquery.DatasetMetadata, error)
	ListTables(ctx context.Context, dataset string) ([]string, error)
	TableMetadata(ctx context.Context, dataset, table string) (*bigquery.TableMetadata, error)
	Read(ctx context.Context, dataset, statement string) (rowIterator, error)
	Run(ctx context.Context, dataset, statement string) (job, error)
	DryRun(ctx context.Context, dataset, statement string) (*jobStatus, error)
}

// rowIterator is a result set returned by Read or job.Read.
type rowIterator interface {
	Schema() bigquery.Schema
	// Next returns the next row's values, or io.EOF at the end.
	Next() ([]bigquery.Value, error)
	// SourceJob returns the job that produced these rows, or nil if none.
	SourceJob() job
}

// job is a submitted BigQuery query job.
type job interface {
	ID() string
	Wait(ctx context.Context) (*jobStatus, error)
	Read(ctx context.Context) (rowIterator, error)
	Status(ctx context.Context) (*jobStatus, error)
}

// jobStatus is the finished status of a BigQuery job, carrying the job error
// (whose field on *bigquery.JobStatus is unexported) and its query statistics.
type jobStatus struct {
	Err   error
	Stats *bigquery.QueryStatistics
}

func jobStatusOf(status *bigquery.JobStatus) *jobStatus {
	if status == nil {
		return nil
	}
	stats, _ := statisticsOf(status)
	return &jobStatus{
		Err:   status.Err(),
		Stats: stats,
	}
}

// liveClient wraps *bigquery.Client to implement the client seam.
type liveClient struct {
	bq *bigquery.Client
	// disableQueryCache is set only by the tape recorder so back-to-back
	// recordings do not flip between fresh runs (which emit query plans) and
	// BigQuery query-cache hits (which omit them). Production Open leaves it
	// false.
	disableQueryCache bool
}

func (c *liveClient) Ping(ctx context.Context) error {
	q := c.bq.Query("SELECT 1")
	q.DisableQueryCache = c.disableQueryCache
	_, err := q.Read(ctx)
	return err
}

func (c *liveClient) Close() error {
	return c.bq.Close()
}

func (c *liveClient) ListDatasets(ctx context.Context) ([]string, error) {
	var ids []string
	it := c.bq.Datasets(ctx)
	for {
		ds, err := it.Next()
		if err == iterator.Done {
			return ids, nil
		}
		if err != nil {
			return nil, err
		}
		ids = append(ids, ds.DatasetID)
	}
}

func (c *liveClient) DatasetMetadata(ctx context.Context, dataset string) (*bigquery.DatasetMetadata, error) {
	return c.bq.Dataset(dataset).Metadata(ctx)
}

func (c *liveClient) ListTables(ctx context.Context, dataset string) ([]string, error) {
	var ids []string
	it := c.bq.Dataset(dataset).Tables(ctx)
	for {
		tbl, err := it.Next()
		if err == iterator.Done {
			return ids, nil
		}
		if err != nil {
			return nil, err
		}
		ids = append(ids, tbl.TableID)
	}
}

func (c *liveClient) TableMetadata(ctx context.Context, dataset, table string) (*bigquery.TableMetadata, error) {
	return c.bq.Dataset(dataset).Table(table).Metadata(ctx)
}

func (c *liveClient) query(dataset, statement string) *bigquery.Query {
	q := c.bq.Query(statement)
	q.DefaultDatasetID = dataset
	q.DisableQueryCache = c.disableQueryCache
	return q
}

func (c *liveClient) Read(ctx context.Context, dataset, statement string) (rowIterator, error) {
	it, err := c.query(dataset, statement).Read(ctx)
	if err != nil {
		return nil, err
	}
	return &liveRowIterator{iter: it}, nil
}

func (c *liveClient) Run(ctx context.Context, dataset, statement string) (job, error) {
	j, err := c.query(dataset, statement).Run(ctx)
	if err != nil {
		return nil, err
	}
	return &liveJob{job: j}, nil
}

func (c *liveClient) DryRun(ctx context.Context, dataset, statement string) (*jobStatus, error) {
	q := c.query(dataset, statement)
	q.DryRun = true
	j, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	return jobStatusOf(j.LastStatus()), nil
}

type liveRowIterator struct {
	iter *bigquery.RowIterator
}

func (r *liveRowIterator) Schema() bigquery.Schema {
	return r.iter.Schema
}

func (r *liveRowIterator) Next() ([]bigquery.Value, error) {
	var values []bigquery.Value
	err := r.iter.Next(&values)
	if err == iterator.Done {
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (r *liveRowIterator) SourceJob() job {
	if j := r.iter.SourceJob(); j != nil {
		return &liveJob{job: j}
	}
	return nil
}

type liveJob struct {
	job *bigquery.Job
}

func (j *liveJob) ID() string {
	return j.job.ID()
}

func (j *liveJob) Wait(ctx context.Context) (*jobStatus, error) {
	status, err := j.job.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return jobStatusOf(status), nil
}

func (j *liveJob) Read(ctx context.Context) (rowIterator, error) {
	it, err := j.job.Read(ctx)
	if err != nil {
		return nil, err
	}
	return &liveRowIterator{iter: it}, nil
}

func (j *liveJob) Status(ctx context.Context) (*jobStatus, error) {
	status, err := j.job.Status(ctx)
	if err != nil {
		return nil, err
	}
	return jobStatusOf(status), nil
}
