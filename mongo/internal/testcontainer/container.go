// Package testcontainer is how mongo's tests reach MongoDB: through a tape
// (enginetest/tape) that replays a pinned server's answers, or, in a recording
// (OMNI_RECORD=1), through the pinned server itself. The server is one
// testcontainer per test binary, started by the first test that records and
// removed by Terminate, which each package's TestMain calls.
//
// MongoDB has no database/sql driver, so the tape cannot sit at that seam.
// Open returns the driver's own client instead: the requests are the calls the
// driver makes of the mongo-go-driver client (listdatabases, listcollections,
// collstats, listindexes, a command, a mongosh statement through gomongo), and
// each answer is what the client returned. An answer's BSON is carried as
// canonical Extended JSON, so the tape holds every BSON kind, documents
// included; container_test.go holds that to a corpus of all of them.
//
// A tape records the client library's answers, not the server's wire traffic.
// A behavior that lives between the two, such as how the library negotiates a
// session or decodes a cursor, is not on a tape, and the live parity harness
// covers it. An answer is narrowed to what the driver asks of it, because some
// carry counters and per-run identifiers the metadata and its dump must not:
// collStats keeps its two sizes, and the collection UUIDs a listCollections
// answer reports are recorded as stand-ins.
package testcontainer

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/bytebase/gomongo"
	"github.com/bytebase/gomongo/types"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is mongo/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin role of the one MongoDB the tests ask.
const Default = "default"

var server = &tape.Server{Pins: Pins, Codec: codec{}}

// codec is how a Mongo answer is read for the tape. An answer carries its own
// error, so there is no error spelling to encode; what a tape needs is the
// collection UUIDs, which a server makes up when it creates a collection, held
// as stand-ins so two recordings of one corpus answer the same bytes.
type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return errors.New(text) }

// madeUp matches the values a server makes up that a tape holds as stand-ins:
// the UUID a listCollections answer reports for a collection, the ObjectID a
// document's _id holds when the insert did not name one, and the host name a
// plan reports, which is the container's own name. Each appears as canonical
// Extended JSON writes it.
var madeUp = regexp.MustCompile(`"(?:uuid|userId)":\{"\$binary":\{"base64":"[A-Za-z0-9+/=]+","subType":"04"}}|"\$oid":"[0-9a-f]{24}"|"host":"[0-9a-f]{12}"|"(?:executionTimeMillis|executionTimeMillisEstimate|optimizationTimeMillis)":\{"\$numberInt":"[0-9]+"}`)

func (codec) MadeUp() *regexp.Regexp { return madeUp }

// planMilliseconds matches the timings a plan reports, which are measurements
// of the run and not of the plan.
var planMilliseconds = regexp.MustCompile(`^"(?:executionTimeMillis|executionTimeMillisEstimate|optimizationTimeMillis)":`)

// StandIn is a made-up value's stand-in, counting n from 1, so a tape holds the
// same bytes in every recording. A plan's timing stands for 0, the value it
// reports for a corpus this size: which millisecond a run took is not an answer
// the tape should pin.
func (codec) StandIn(name string, n int) string {
	switch {
	case strings.HasPrefix(name, `"$oid"`):
		return `"$oid":"` + fmt.Sprintf("%024x", n) + `"`
	case strings.HasPrefix(name, `"host"`):
		return `"host":"` + fmt.Sprintf("%012x", n) + `"`
	case planMilliseconds.MatchString(name):
		field, _, _ := strings.Cut(name, ":")
		return field + `:{"$numberInt":"0"}`
	}
	field, _, _ := strings.Cut(name, ":")
	var raw [16]byte
	binary.BigEndian.PutUint64(raw[8:], uint64(n))
	return field + `:{"$binary":{"base64":"` + base64.StdEncoding.EncodeToString(raw[:]) + `","subType":"04"}}`
}

// Client is the driver's MongoDB client, answered by t's tape. It is what a
// test hands the driver in place of the live client; its methods are the calls
// the driver makes, and the recordings hold what the live client answered.
type Client struct {
	tape *tape.Raw
}

// Open returns t's client to the pinned server in pin row role: a replay of
// its tape, or, in a recording, the pinned server with every call written down.
func Open(t testing.TB) *Client {
	t.Helper()
	return &Client{tape: tape.OpenRaw(t, server, Default, func(request string) (string, error) {
		// A recording reaches the server on its first call, so a test that
		// replays never starts a container.
		return liveDo(t, request)
	})}
}

// Ping sends a ping.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.do(ctx, request{Op: opPing})
	return err
}

// Disconnect closes the client.
func (c *Client) Disconnect(ctx context.Context) error {
	_, err := c.do(ctx, request{Op: opDisconnect})
	return err
}

// ListDatabaseNames returns every database's name.
func (c *Client) ListDatabaseNames(ctx context.Context) ([]string, error) {
	doc, err := c.do(ctx, request{Op: opListDatabaseNames})
	if err != nil {
		return nil, err
	}
	var names struct {
		Names []string `json:"names"`
	}
	if err := unmarshalJSON(doc, &names); err != nil {
		return nil, err
	}
	return names.Names, nil
}

// ListCollections returns the collection and view descriptions of database, as
// listCollections reports them.
func (c *Client) ListCollections(ctx context.Context, database string) ([]bson.M, error) {
	doc, err := c.do(ctx, request{Op: opListCollections, Database: database})
	if err != nil {
		return nil, err
	}
	return decodeDocs(doc)
}

// EstimatedDocumentCount returns a collection's estimated document count.
func (c *Client) EstimatedDocumentCount(ctx context.Context, database, collection string) (int64, error) {
	doc, err := c.do(ctx, request{Op: opEstimatedDocumentCount, Database: database, Collection: collection})
	if err != nil {
		return 0, err
	}
	var answer struct {
		Count *int64 `json:"count"`
	}
	if err := unmarshalJSON(doc, &answer); err != nil {
		return 0, err
	}
	if answer.Count == nil {
		return 0, errors.New("testcontainer: an estimated count answer without a count")
	}
	return *answer.Count, nil
}

// CollStats returns the collStats command's answer for one collection.
func (c *Client) CollStats(ctx context.Context, database, collection string) (bson.M, error) {
	doc, err := c.do(ctx, request{Op: opCollStats, Database: database, Collection: collection})
	if err != nil {
		return nil, err
	}
	return decodeDoc(doc)
}

// ListIndexes returns the collection's index descriptions, as listIndexes
// reports them.
func (c *Client) ListIndexes(ctx context.Context, database, collection string) ([]bson.M, error) {
	doc, err := c.do(ctx, request{Op: opListIndexes, Database: database, Collection: collection})
	if err != nil {
		return nil, err
	}
	return decodeDocs(doc)
}

// RunCommand runs command on database and returns its answer.
func (c *Client) RunCommand(ctx context.Context, database string, command bson.D) (bson.M, error) {
	encoded, err := bson.MarshalExtJSON(command, true, false)
	if err != nil {
		return nil, err
	}
	doc, err := c.do(ctx, request{Op: opRunCommand, Database: database, Command: encoded})
	if err != nil {
		return nil, err
	}
	return decodeDoc(doc)
}

// Execute runs one mongosh statement, as gomongo does, reading at most maxRows
// rows when that is positive.
func (c *Client) Execute(ctx context.Context, database, statement string, maxRows int64) (*gomongo.Result, error) {
	doc, err := c.do(ctx, request{Op: opExecute, Database: database, Statement: statement, MaxRows: maxRows})
	if err != nil {
		return nil, err
	}
	var payload struct {
		Operation operationType     `json:"operation"`
		Values    []json.RawMessage `json:"values"`
	}
	if err := unmarshalJSON(doc, &payload); err != nil {
		return nil, err
	}
	result := &gomongo.Result{Operation: types.OperationType(payload.Operation)}
	for _, raw := range payload.Values {
		value, err := decodeValue(raw)
		if err != nil {
			return nil, err
		}
		result.Value = append(result.Value, value)
	}
	return result, nil
}

// do sends one call and returns its answer as canonical Extended JSON.
func (c *Client) do(ctx context.Context, r request) (json.RawMessage, error) {
	text, err := r.text()
	if err != nil {
		return nil, err
	}
	answered, err := c.tape.Do(text)
	if err != nil {
		return nil, err
	}
	var a answer
	if err := unmarshalJSON([]byte(answered), &a); err != nil {
		return nil, fmt.Errorf("testcontainer: unreadable answer to %s: %w", text, err)
	}
	if a.Error != nil {
		return nil, a.Error.err()
	}
	if len(a.Value) == 0 {
		return json.RawMessage("{}"), nil
	}
	return a.Value, nil
}

// The calls the driver makes of its client, as a tape names them.
const (
	opPing                   = "Ping"
	opDisconnect             = "Disconnect"
	opListDatabaseNames      = "ListDatabaseNames"
	opListCollections        = "ListCollections"
	opEstimatedDocumentCount = "EstimatedDocumentCount"
	opCollStats              = "CollStats"
	opListIndexes            = "ListIndexes"
	opRunCommand             = "RunCommand"
	opExecute                = "Execute"
)

// request is one call as a tape holds it: the operation, its arguments, and a
// command as canonical Extended JSON.
type request struct {
	Op         string          `json:"op"`
	Database   string          `json:"database,omitempty"`
	Collection string          `json:"collection,omitempty"`
	Statement  string          `json:"statement,omitempty"`
	MaxRows    int64           `json:"maxRows,omitempty"`
	Command    json.RawMessage `json:"command,omitempty"`
}

func (r request) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// answer is one call's answer as a tape holds it: the client's value as
// canonical Extended JSON, or the error it returned instead.
type answer struct {
	Value json.RawMessage `json:"value,omitempty"`
	Error *clientError    `json:"error,omitempty"`
}

// clientError is an error the client returned, as a tape holds it. A MongoDB
// command error keeps the fields the driver reads, so a replay's type
// assertion finds the same type a recording's error had; any other error is
// its text alone.
type clientError struct {
	Code    int32    `json:"code,omitempty"`
	Name    string   `json:"name,omitempty"`
	Labels  []string `json:"labels,omitempty"`
	Message string   `json:"message"`
}

func (e *clientError) err() error {
	if e.Name == "" && e.Code == 0 && len(e.Labels) == 0 {
		return errors.New(e.Message)
	}
	return mongo.CommandError{Code: e.Code, Name: e.Name, Labels: e.Labels, Message: e.Message}
}

// encodeError renders err as an answer's error half.
func encodeError(err error) *clientError {
	var commandError mongo.CommandError
	if errors.As(err, &commandError) {
		return &clientError{
			Code:    commandError.Code,
			Name:    commandError.Name,
			Labels:  commandError.Labels,
			Message: commandError.Message,
		}
	}
	return &clientError{Message: err.Error()}
}

// operationType names a gomongo operation, so a tape reads as the operation it
// recorded rather than as its number.
type operationType types.OperationType

var operationNames = map[types.OperationType]string{
	types.OpUnknown:                "OpUnknown",
	types.OpNoOp:                   "OpNoOp",
	types.OpFind:                   "OpFind",
	types.OpFindOne:                "OpFindOne",
	types.OpAggregate:              "OpAggregate",
	types.OpShowDatabases:          "OpShowDatabases",
	types.OpShowCollections:        "OpShowCollections",
	types.OpGetCollectionNames:     "OpGetCollectionNames",
	types.OpGetCollectionInfos:     "OpGetCollectionInfos",
	types.OpGetIndexes:             "OpGetIndexes",
	types.OpCountDocuments:         "OpCountDocuments",
	types.OpEstimatedDocumentCount: "OpEstimatedDocumentCount",
	types.OpDistinct:               "OpDistinct",
	types.OpInsertOne:              "OpInsertOne",
	types.OpInsertMany:             "OpInsertMany",
	types.OpUpdateOne:              "OpUpdateOne",
	types.OpUpdateMany:             "OpUpdateMany",
	types.OpReplaceOne:             "OpReplaceOne",
	types.OpDeleteOne:              "OpDeleteOne",
	types.OpDeleteMany:             "OpDeleteMany",
	types.OpFindOneAndUpdate:       "OpFindOneAndUpdate",
	types.OpFindOneAndReplace:      "OpFindOneAndReplace",
	types.OpFindOneAndDelete:       "OpFindOneAndDelete",
	types.OpCreateIndex:            "OpCreateIndex",
	types.OpCreateIndexes:          "OpCreateIndexes",
	types.OpDropIndex:              "OpDropIndex",
	types.OpDropIndexes:            "OpDropIndexes",
	types.OpDrop:                   "OpDrop",
	types.OpCreateCollection:       "OpCreateCollection",
	types.OpDropDatabase:           "OpDropDatabase",
	types.OpRenameCollection:       "OpRenameCollection",
	types.OpDbStats:                "OpDbStats",
	types.OpCollectionStats:        "OpCollectionStats",
	types.OpServerStatus:           "OpServerStatus",
	types.OpServerBuildInfo:        "OpServerBuildInfo",
	types.OpDbVersion:              "OpDbVersion",
	types.OpHostInfo:               "OpHostInfo",
	types.OpListCommands:           "OpListCommands",
	types.OpRunCommand:             "OpRunCommand",
	types.OpDataSize:               "OpDataSize",
	types.OpStorageSize:            "OpStorageSize",
	types.OpTotalIndexSize:         "OpTotalIndexSize",
	types.OpTotalSize:              "OpTotalSize",
	types.OpIsCapped:               "OpIsCapped",
	types.OpValidate:               "OpValidate",
	types.OpLatencyStats:           "OpLatencyStats",
}

var operationValues = func() map[string]types.OperationType {
	values := make(map[string]types.OperationType, len(operationNames))
	for op, name := range operationNames {
		values[name] = op
	}
	return values
}()

func (o operationType) MarshalJSON() ([]byte, error) {
	name, ok := operationNames[types.OperationType(o)]
	if !ok {
		return nil, fmt.Errorf("testcontainer: unknown gomongo operation %d", int(o))
	}
	return json.Marshal(name)
}

func (o *operationType) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}
	value, ok := operationValues[name]
	if !ok {
		return fmt.Errorf("testcontainer: unknown gomongo operation %q", name)
	}
	*o = operationType(value)
	return nil
}

// unmarshalJSON decodes JSON strictly enough that a tape's typo fails the test
// rather than answering something else.
func unmarshalJSON(b []byte, v any) error {
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// encodeDocs renders documents as a tape holds them: a document whose "docs"
// field is the array.
func encodeDocs(docs []bson.M) (json.RawMessage, error) {
	values := make(bson.A, len(docs))
	for i, doc := range docs {
		values[i] = doc
	}
	return encodeDoc(bson.D{{Key: "docs", Value: values}})
}

// encodeDoc renders a document as a tape holds it: canonical Extended JSON,
// which keeps every BSON kind, with the keys of every bson.M in a fixed order,
// since a Go map iterates in a random one.
func encodeDoc(doc any) (json.RawMessage, error) {
	b, err := bson.MarshalExtJSON(ordered(doc), true, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ordered returns v with every bson.M under it as a bson.D of sorted keys. A
// bson.D keeps its order, which is the server's.
func ordered(v any) any {
	switch x := v.(type) {
	case bson.M:
		keys := make([]string, 0, len(x))
		for key := range x {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		doc := make(bson.D, len(keys))
		for i, key := range keys {
			doc[i] = bson.E{Key: key, Value: ordered(x[key])}
		}
		return doc
	case map[string]any:
		return ordered(bson.M(x))
	case bson.D:
		doc := make(bson.D, len(x))
		for i, e := range x {
			doc[i] = bson.E{Key: e.Key, Value: ordered(e.Value)}
		}
		return doc
	case bson.A:
		values := make(bson.A, len(x))
		for i, e := range x {
			values[i] = ordered(e)
		}
		return values
	case []any:
		values := make(bson.A, len(x))
		for i, e := range x {
			values[i] = ordered(e)
		}
		return values
	}
	return v
}

// decodeDocs reads the documents encodeDocs wrote.
func decodeDocs(doc json.RawMessage) ([]bson.M, error) {
	var payload struct {
		Docs []json.RawMessage `json:"docs"`
	}
	if err := unmarshalJSON(doc, &payload); err != nil {
		return nil, err
	}
	docs := make([]bson.M, len(payload.Docs))
	for i, raw := range payload.Docs {
		d, err := decodeDoc(raw)
		if err != nil {
			return nil, err
		}
		docs[i] = d
	}
	return docs, nil
}

// decodeDoc reads one document as the client decodes it: a bson.M at the top,
// with nested documents as bson.D, which is what decoding into a bson.M gives.
func decodeDoc(doc json.RawMessage) (bson.M, error) {
	var m bson.M
	if err := bson.UnmarshalExtJSON(doc, true, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// decodeValue reads one value of an execute answer: a document's field holds
// it, since Extended JSON has no top-level scalar.
func decodeValue(raw json.RawMessage) (any, error) {
	var d bson.D
	if err := bson.UnmarshalExtJSON(append(append([]byte(`{"v":`), raw...), '}'), true, &d); err != nil {
		return nil, err
	}
	if len(d) != 1 {
		return nil, fmt.Errorf("testcontainer: want one value, got %d", len(d))
	}
	return d[0].Value, nil
}

// --- the live server ---

const (
	user     = "testuser"
	password = "testpass"
	// connectTimeout keeps a listener that never answers from hanging a test.
	connectTimeout = 10 * time.Second
)

var liveServer struct {
	once   sync.Once
	client *mongo.Client
	err    error
}

// live returns the pinned server's client, starting the container on the first
// call. Only a recording calls it.
func live() *mongo.Client {
	liveServer.once.Do(func() {
		liveServer.client, liveServer.err = start(context.Background())
	})
	return liveServer.client
}

func start(ctx context.Context) (*mongo.Client, error) {
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		return nil, err
	}
	// testcontainers-go panics from MustExtractDockerHost when the Docker
	// socket is reachable but access is denied, and a test must report that
	// rather than crash the binary.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	container, err := testcontainers.Run(ctx, pin.Image,
		testcontainers.WithEnv(map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": user,
			"MONGO_INITDB_ROOT_PASSWORD": password,
		}),
		testcontainers.WithExposedPorts("27017/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").WithStartupTimeout(3*time.Minute)),
	)
	if err != nil {
		return nil, err
	}
	if err := checkDigest(ctx, container, pin.Digest); err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return nil, err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	client, err := connect(ctx, net.JoinHostPort(host, port.Port()))
	if err != nil {
		return nil, err
	}
	var build bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "buildInfo", Value: 1}}).Decode(&build); err != nil {
		return nil, fmt.Errorf("MongoDB does not answer: %w", err)
	}
	version, _ := build["version"].(string)
	if want := strings.TrimPrefix(pin.Tag, "r"); version != want {
		return nil, fmt.Errorf("the server reports %q, but pins.tsv pins MongoDB %s (%s)", version, want, pin.Tag)
	}
	return client, nil
}

// connect returns a client for the pinned server at addr, signing in as the
// container's root user.
func connect(ctx context.Context, addr string) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin&directConnection=true&connectTimeoutMS=%d",
		user, password, addr, connectTimeout.Milliseconds())
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	return client, nil
}

func checkDigest(ctx context.Context, container testcontainers.Container, digest string) error {
	info, err := container.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the MongoDB container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the MongoDB image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the MongoDB container runs an image with digests %v, but pins.tsv pins %s", image.RepoDigests, digest)
}

// emptied records the top-level tests whose first live call has emptied the
// server, so a test reads the same answers recorded alone as in a run of its
// package.
var emptied struct {
	mu   sync.Mutex
	done map[string]bool
}

// emptyOnce drops every database on t's top-level test's first live call but
// MongoDB's own, so what an earlier test left behind cannot reach the tape.
func emptyOnce(t testing.TB, client *mongo.Client) {
	t.Helper()
	top := enginetest.TopName(t.Name())
	emptied.mu.Lock()
	defer emptied.mu.Unlock()
	if emptied.done[top] {
		return
	}
	ctx := context.Background()
	names, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		t.Fatalf("empty MongoDB: %v", err)
	}
	for _, name := range names {
		if name == "admin" || name == "config" || name == "local" {
			continue
		}
		if err := client.Database(name).Drop(ctx); err != nil {
			t.Fatalf("empty MongoDB: drop %s: %v", name, err)
		}
	}
	if emptied.done == nil {
		emptied.done = map[string]bool{}
	}
	emptied.done[top] = true
}

// Terminate removes the server this run started.
func Terminate() {
	if liveServer.client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := liveServer.client.Disconnect(ctx); err != nil {
		fmt.Printf("testcontainer: disconnect MongoDB: %v\n", err)
	}
}

// --- the live client's answers ---

// liveDo answers one request from the pinned server, in the encoding a tape
// holds.
func liveDo(t testing.TB, text string) (string, error) {
	t.Helper()
	var req request
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %q: %w", text, err)
	}
	ctx := context.Background()
	client := live()
	if client == nil {
		t.Fatalf("start MongoDB: %v", liveServer.err)
	}
	emptyOnce(t, client)
	var a answer
	switch req.Op {
	case opPing:
		if err := client.Ping(ctx, nil); err != nil {
			a.Error = encodeError(err)
		}
	case opDisconnect:
	case opListDatabaseNames:
		names, err := client.ListDatabaseNames(ctx, bson.M{})
		if err != nil {
			a.Error = encodeError(err)
			break
		}
		payload, err := json.Marshal(struct {
			Names []string `json:"names"`
		}{names})
		if err != nil {
			return "", err
		}
		a.Value = payload
	case opEstimatedDocumentCount:
		count, err := client.Database(req.Database).Collection(req.Collection).EstimatedDocumentCount(ctx)
		if err != nil {
			a.Error = encodeError(err)
			break
		}
		payload, err := json.Marshal(struct {
			Count int64 `json:"count"`
		}{count})
		if err != nil {
			return "", err
		}
		a.Value = payload
	case opListCollections:
		docs, err := liveListCollections(ctx, client, req.Database)
		if err != nil {
			a.Error = encodeError(err)
			break
		}
		if a.Value, err = encodeDocs(docs); err != nil {
			return "", err
		}
	case opCollStats:
		doc, err := liveRunCommand(ctx, client, req.Database, bson.D{{Key: "collStats", Value: req.Collection}})
		if err != nil {
			a.Error = encodeError(err)
			break
		}
		doc = collStatsFields(doc)
		if a.Value, err = encodeDoc(doc); err != nil {
			return "", err
		}
	case opListIndexes:
		docs, err := liveListIndexes(ctx, client, req.Database, req.Collection)
		if err != nil {
			a.Error = encodeError(err)
			break
		}
		if a.Value, err = encodeDocs(docs); err != nil {
			return "", err
		}
	case opRunCommand:
		var command bson.D
		if err := bson.UnmarshalExtJSON(req.Command, true, &command); err != nil {
			return "", fmt.Errorf("testcontainer: unreadable command in %s: %w", text, err)
		}
		doc, err := liveRunCommand(ctx, client, req.Database, command)
		if err != nil {
			a.Error = encodeError(err)
			break
		}
		if a.Value, err = encodeDoc(doc); err != nil {
			return "", err
		}
	case opExecute:
		var opts []gomongo.ExecuteOption
		if req.MaxRows > 0 {
			opts = append(opts, gomongo.WithMaxRows(req.MaxRows))
		}
		result, err := gomongo.NewClient(client).Execute(ctx, req.Database, req.Statement, opts...)
		if err != nil {
			a.Error = encodeError(err)
			break
		}
		values := make(bson.A, len(result.Value))
		copy(values, result.Value)
		name, ok := operationNames[result.Operation]
		if !ok {
			return "", fmt.Errorf("testcontainer: gomongo answered an unknown operation %d", int(result.Operation))
		}
		doc := bson.D{{Key: "operation", Value: name}, {Key: "values", Value: values}}
		if a.Value, err = encodeDoc(doc); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("testcontainer: unknown operation %q", req.Op)
	}
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// liveListCollections returns the collection descriptions the client decodes
// from a listCollections cursor, ordered by name as the driver's client orders
// them: MongoDB lists collections in a hash order that changes between runs.
func liveListCollections(ctx context.Context, client *mongo.Client, database string) ([]bson.M, error) {
	cursor, err := client.Database(database).ListCollections(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var docs []bson.M
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	sortDocsByName(docs)
	return docs, cursor.Close(ctx)
}

// liveListIndexes returns the index descriptions the client decodes from a
// listIndexes cursor, ordered by name as the driver's client orders them.
func liveListIndexes(ctx context.Context, client *mongo.Client, database, collection string) ([]bson.M, error) {
	cursor, err := client.Database(database).Collection(collection).Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	var docs []bson.M
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	sortDocsByName(docs)
	return docs, cursor.Close(ctx)
}

// sortDocsByName orders documents by their "name".
func sortDocsByName(docs []bson.M) {
	sort.SliceStable(docs, func(i, j int) bool {
		left, _ := docs[i]["name"].(string)
		right, _ := docs[j]["name"].(string)
		return left < right
	})
}

// collStatsFields keeps what the driver reads of a collStats answer, as
// mongo/driver's client does: the rest carries statistics counters that move
// between runs.
func collStatsFields(doc bson.M) bson.M {
	fields := bson.M{}
	for _, key := range []string{"size", "totalIndexSize"} {
		if value, ok := doc[key]; ok {
			fields[key] = value
		}
	}
	return fields
}

// liveRunCommand runs command and decodes its answer the way the driver does.
func liveRunCommand(ctx context.Context, client *mongo.Client, database string, command bson.D) (bson.M, error) {
	var doc bson.M
	if err := client.Database(database).RunCommand(ctx, command).Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}
