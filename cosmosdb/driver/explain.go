package driver

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
	"github.com/bytebase/omni/cosmosdb/driver/queryengine/gonative"
)

// planCapture is the query engine that answers a plan request. The client
// fetches a statement's plan from the gateway and hands it to the engine a
// query runs through, so this engine keeps the plan and hands back a pipeline
// with nothing to do: the plan is in hand and no partition is read.
type planCapture struct {
	query string
	plan  string
}

var _ queryengine.QueryEngine = (*planCapture)(nil)

// CreateQueryPipeline keeps the gateway's plan and starts no work.
func (c *planCapture) CreateQueryPipeline(query, plan, _ string) (queryengine.QueryPipeline, error) {
	c.query, c.plan = query, plan
	return plannedPipeline{query: query}, nil
}

// CreateReadManyPipeline is unreachable: a plan request is a query, not a list
// of items to read by id.
func (*planCapture) CreateReadManyPipeline([]queryengine.ItemIdentity, string, string, uint8, []string) (queryengine.QueryPipeline, error) {
	return nil, errors.New("a plan request reads no items")
}

// SupportedFeatures is what the query engine this driver runs queries through
// claims, so the gateway states the plan that engine would consume.
func (*planCapture) SupportedFeatures() string { return gonative.Default().SupportedFeatures() }

// plannedPipeline is a pipeline whose work is done before it starts: the plan
// is in hand, so it asks for no data and is complete at once.
type plannedPipeline struct {
	query string
}

var _ queryengine.QueryPipeline = plannedPipeline{}

// Query is the query the gateway's plan says the ranges are read with. No
// range is read, so the statement's own text stands for it.
func (p plannedPipeline) Query() string { return p.query }

func (plannedPipeline) IsComplete() bool { return true }

func (plannedPipeline) Run() (*queryengine.PipelineResult, error) {
	return &queryengine.PipelineResult{IsCompleted: true}, nil
}

func (plannedPipeline) ProvideData([]queryengine.QueryResult) error { return nil }

func (plannedPipeline) Close() {}
