package testcontainer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
)

// The query span suite of cosmosdb/analysis asks the emulator through these.
// A fixture is a container of its own in SpanDatabase, named after the
// fixture and partitioned on /pk, holding the fixture's documents; a query
// runs in the partition of its first document.
//
// A request is a SpanRequest and an answer a SpanAnswer, both as JSON: the
// tape of the suite's test holds them.
const (
	// SpanDatabase holds the span suite's fixture containers.
	SpanDatabase = "cosmosspan"
	// SpanPartitionKey is the partition key path of a span fixture.
	SpanPartitionKey = "/pk"
)

// SystemProperties are the properties the emulator adds to a stored
// document. Their values name the store that answered (a resource id, an
// etag, a write time), so a SpanAnswer holds SystemValue in their place.
var SystemProperties = []string{"_rid", "_self", "_etag", "_attachments", "_ts"}

// SystemValue stands for the value of a system property in a SpanAnswer.
const SystemValue = "(system)"

// A SpanRequest is one call of the span suite: load a fixture's documents,
// or run a statement on a loaded fixture.
type SpanRequest struct {
	Op        string `json:"op"` // load or query
	Fixture   string `json:"fixture"`
	Documents string `json:"documents,omitempty"` // load: one JSON document per line
	Statement string `json:"statement,omitempty"` // query
}

// A SpanAnswer is the emulator's answer to a SpanRequest. For a query it is
// the items the statement returned, as the gateway sent them, and whether the
// gateway's query plan says the statement selects a bare VALUE; or the error
// code the emulator refused the statement with.
type SpanAnswer struct {
	Error       string            `json:"error,omitempty"`
	SelectValue bool              `json:"selectValue,omitempty"`
	Items       []json.RawMessage `json:"items,omitempty"`
}

// SpanLive answers the span suite's requests from the emulator, starting it
// on first use. It is the live function of the suite's tape.
func SpanLive(t testing.TB) func(string) (string, error) {
	loaded := map[string]string{}
	return func(text string) (string, error) {
		var req SpanRequest
		if err := json.Unmarshal([]byte(text), &req); err != nil {
			return "", fmt.Errorf("testcontainer: unreadable span request %s: %w", text, err)
		}
		shared.Lock()
		err := start()
		endpoint := shared.endpoint
		shared.Unlock()
		if err != nil {
			return "", err
		}
		c, err := client(endpoint)
		if err != nil {
			return "", err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		var a SpanAnswer
		switch req.Op {
		case "load":
			if err := loadSpanFixture(ctx, c, req.Fixture, req.Documents); err != nil {
				return "", err
			}
			loaded[req.Fixture] = req.Documents
		case "query":
			docs, ok := loaded[req.Fixture]
			if !ok {
				return "", fmt.Errorf("testcontainer: span fixture %s queried before it was loaded", req.Fixture)
			}
			a, err = spanQuery(ctx, c, req.Fixture, docs, req.Statement)
			if err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("testcontainer: unknown span request %q", req.Op)
		}
		b, err := json.Marshal(a)
		return string(b), err
	}
}

func client(endpoint string) (*azcosmos.Client, error) {
	credential, err := azcosmos.NewKeyCredential(EmulatorKey)
	if err != nil {
		return nil, fmt.Errorf("emulator key credential: %w", err)
	}
	insecure, err := InsecureTLS("")
	if err != nil {
		return nil, err
	}
	return azcosmos.NewClientWithKey(endpoint, credential, &azcosmos.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &http.Client{Transport: &http.Transport{TLSClientConfig: insecure}},
		},
	})
}

// loadSpanFixture writes a fixture's documents into its container. An
// upsert makes a second load of the same fixture write the same documents.
func loadSpanFixture(ctx context.Context, c *azcosmos.Client, fixture, documents string) error {
	if _, err := c.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: SpanDatabase}, nil); err != nil && !alreadyExists(err) {
		return fmt.Errorf("create the span database: %w", err)
	}
	database, err := c.NewDatabase(SpanDatabase)
	if err != nil {
		return err
	}
	if err := createContainer(ctx, database, fixture, SpanPartitionKey, 400); err != nil {
		return err
	}
	container, err := database.NewContainer(fixture)
	if err != nil {
		return err
	}
	for _, doc := range spanDocuments(documents) {
		pk, err := partitionKey(doc)
		if err != nil {
			return fmt.Errorf("span fixture %s: %w", fixture, err)
		}
		if _, err := container.UpsertItem(ctx, azcosmos.NewPartitionKeyString(pk), []byte(doc), nil); err != nil {
			return fmt.Errorf("span fixture %s: write %s: %w", fixture, doc, err)
		}
	}
	return nil
}

func spanDocuments(documents string) []string {
	var out []string
	for _, line := range strings.Split(documents, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}

func partitionKey(doc string) (string, error) {
	var d struct {
		PK *string `json:"pk"`
	}
	if err := json.Unmarshal([]byte(doc), &d); err != nil || d.PK == nil {
		return "", fmt.Errorf("a document needs a string pk: %s", doc)
	}
	return *d.PK, nil
}

// scCode is the code the query compiler gives a refusal (SC1001 for a syntax
// error), which the gateway's error message carries in a JSON document.
var scCode = regexp.MustCompile(`\\?"code\\?"\s*:\s*\\?"(SC\d+)\\?"`)

// spanQuery runs a statement in the partition of the fixture's first
// document and asks the gateway for its plan.
func spanQuery(ctx context.Context, c *azcosmos.Client, fixture, documents, statement string) (SpanAnswer, error) {
	docs := spanDocuments(documents)
	if len(docs) == 0 {
		return SpanAnswer{}, fmt.Errorf("span fixture %s holds no document", fixture)
	}
	pk, err := partitionKey(docs[0])
	if err != nil {
		return SpanAnswer{}, err
	}
	container, err := c.NewContainer(SpanDatabase, fixture)
	if err != nil {
		return SpanAnswer{}, err
	}
	var a SpanAnswer
	pager := container.NewQueryItemsPager(statement, azcosmos.NewPartitionKeyString(pk), nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return refusal(err)
		}
		for _, item := range page.Items {
			narrowed, err := narrowSystem(item)
			if err != nil {
				return SpanAnswer{}, err
			}
			a.Items = append(a.Items, narrowed)
		}
	}

	// The plan is what the gateway states about the statement before any
	// partition is read; the client fetches it for a query its query engine
	// runs, and the capture keeps it and reads nothing.
	capture := &planCapture{}
	plan := container.NewQueryItemsPager(statement, azcosmos.NewPartitionKey(), &azcosmos.QueryOptions{QueryEngine: capture})
	if _, err := plan.NextPage(ctx); err != nil {
		return SpanAnswer{}, fmt.Errorf("the gateway's plan of %s: %w", statement, err)
	}
	var p struct {
		QueryInfo struct {
			HasSelectValue bool `json:"hasSelectValue"`
		} `json:"queryInfo"`
	}
	if err := json.Unmarshal([]byte(capture.plan), &p); err != nil {
		return SpanAnswer{}, fmt.Errorf("the gateway's plan of %s: %w: %s", statement, err, capture.plan)
	}
	a.SelectValue = p.QueryInfo.HasSelectValue
	return a, nil
}

// refusal is the answer for a query the emulator refused: the compiler's
// code, or the status the gateway answered with when it names none.
func refusal(err error) (SpanAnswer, error) {
	var responseErr *azcore.ResponseError
	if !errors.As(err, &responseErr) {
		return SpanAnswer{}, err
	}
	if m := scCode.FindStringSubmatch(err.Error()); m != nil {
		return SpanAnswer{Error: m[1]}, nil
	}
	return SpanAnswer{Error: fmt.Sprintf("%d %s", responseErr.StatusCode, responseErr.ErrorCode)}, nil
}

// narrowSystem replaces the value of every system property of an item with
// SystemValue, keeping the item's members in the order the gateway sent them.
func narrowSystem(item []byte) (json.RawMessage, error) {
	dec := json.NewDecoder(strings.NewReader(string(item)))
	dec.UseNumber()
	var b strings.Builder
	if err := copyNarrowed(dec, &b, false); err != nil {
		return nil, fmt.Errorf("an item the gateway returned: %w: %s", err, item)
	}
	return json.RawMessage(b.String()), nil
}

// copyNarrowed copies one JSON value from dec to b; system says the value is
// a system property's.
func copyNarrowed(dec *json.Decoder, b *strings.Builder, system bool) error {
	// A loop with an explicit stack would be needed for deep values; an item
	// is a stored document, which Cosmos DB caps at 128 levels.
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if system {
		if d, ok := tok.(json.Delim); ok && (d == '{' || d == '[') {
			if err := skipRest(dec); err != nil {
				return err
			}
		}
		v, _ := json.Marshal(SystemValue)
		b.Write(v)
		return nil
	}
	switch d := tok.(type) {
	case json.Delim:
		switch d {
		case '{':
			b.WriteByte('{')
			for i := 0; dec.More(); i++ {
				key, err := dec.Token()
				if err != nil {
					return err
				}
				if i > 0 {
					b.WriteByte(',')
				}
				k, _ := json.Marshal(key)
				b.Write(k)
				b.WriteByte(':')
				name, _ := key.(string)
				if err := copyNarrowed(dec, b, isSystem(name)); err != nil {
					return err
				}
			}
			if _, err := dec.Token(); err != nil {
				return err
			}
			b.WriteByte('}')
		case '[':
			b.WriteByte('[')
			for i := 0; dec.More(); i++ {
				if i > 0 {
					b.WriteByte(',')
				}
				if err := copyNarrowed(dec, b, false); err != nil {
					return err
				}
			}
			if _, err := dec.Token(); err != nil {
				return err
			}
			b.WriteByte(']')
		}
	default:
		v, err := json.Marshal(d)
		if err != nil {
			return err
		}
		b.Write(v)
	}
	return nil
}

// skipRest consumes the rest of an object or array whose opening delimiter
// was read.
func skipRest(dec *json.Decoder) error {
	for depth := 1; depth > 0; {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		if d, ok := tok.(json.Delim); ok {
			if d == '{' || d == '[' {
				depth++
			} else {
				depth--
			}
		}
	}
	return nil
}

func isSystem(name string) bool {
	for _, p := range SystemProperties {
		if p == name {
			return true
		}
	}
	return false
}

// planCapture is the query engine of a plan request: it keeps the plan the
// gateway stated and hands back a pipeline with nothing to do.
type planCapture struct {
	plan string
}

var _ queryengine.QueryEngine = (*planCapture)(nil)

func (c *planCapture) CreateQueryPipeline(query, plan, _ string) (queryengine.QueryPipeline, error) {
	c.plan = plan
	return donePipeline{query: query}, nil
}

func (*planCapture) CreateReadManyPipeline([]queryengine.ItemIdentity, string, string, uint8, []string) (queryengine.QueryPipeline, error) {
	return nil, errors.New("a plan request reads no items")
}

// SupportedFeatures is what cosmosdb/driver's query engine claims, so the
// gateway states the plan the driver would run the statement with.
func (*planCapture) SupportedFeatures() string {
	return "Aggregate,NonValueAggregate,MultipleAggregates,Distinct,DistinctValue,OrderBy,Top,OffsetAndLimit,GroupBy"
}

type donePipeline struct{ query string }

func (p donePipeline) Query() string                             { return p.query }
func (donePipeline) IsComplete() bool                            { return true }
func (donePipeline) ProvideData([]queryengine.QueryResult) error { return nil }
func (donePipeline) Close()                                      {}
func (donePipeline) Run() (*queryengine.PipelineResult, error) {
	return &queryengine.PipelineResult{IsCompleted: true}, nil
}
