// Package testcontainer is how redis's tests reach Redis: through a tape
// (enginetest/tape) that replays a pinned server's answers, or, in a recording
// (OMNI_RECORD=1), through the pinned server itself. The server is one
// testcontainer per test binary, started by the first test that records and
// removed by Terminate, which each package's TestMain calls.
//
// The pinned server runs cluster-enabled with every hash slot on its one node,
// because the driver's all-nodes query path needs a cluster; a standalone
// client reaches the same server for every other test.
//
// The tape does not sit at database/sql, since go-redis is not a database/sql
// driver, so Open returns the driver's own client: the requests are the calls
// the driver makes of a Redis client (a ping, one INFO section, one CONFIG GET
// parameter, a pipeline of commands, a pipeline to every cluster node), and
// each answer is what go-redis returned — a value, or the error the server
// answered with. A value is written as the kind go-redis decoded it to, so an
// array keeps its elements and a hash keeps its fields.
//
// What a tape cannot hold is what a Redis client does between the calls:
// server selection, cluster topology, connection pooling and the handshake
// each connection runs. It cannot hold an INFO section whole either, because
// the counters in it move between runs, so an answer is narrowed to the lines
// the driver reads (redis_version, and the keyspace's dbN lines). A command's
// answers are held per command, in the order the pipeline sent them, so a
// reply the driver never reads is still in the tape; a pipeline that could not
// be sent at all is not an engine answer and fails a recording.
package testcontainer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is redis/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Driver is the pin row of the server the driver's tapes were recorded on.
const Driver = "driver"

var server = &tape.Server{Pins: Pins, Codec: codec{}}

// codec is how a Redis answer is read for the tape. A server's refusal is one of
// a command's answers, so the tape needs no error spelling from the codec; a
// call-level error is never the engine's answer and fails a recording.
type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return errors.New(text) }

// Login is the account a client talks as. The empty login is Redis's own default
// user, which needs no password until an ACL gives it one.
type Login struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

// Client is the driver's Redis client, answered by t's tape. It is what a test
// hands the driver in place of the live client; its methods are the calls the
// driver makes, and the recordings hold what go-redis answered.
type Client struct {
	tape  *tape.Raw
	login Login
	// cluster is the client the recording asks: a cluster client can run a
	// pipeline on every node, which a standalone one cannot.
	cluster bool
}

// Open returns t's client to the pinned server, as login on one server: a replay
// of its tape, or, in a recording, the pinned server with every call written
// down.
func Open(t testing.TB, login Login) *Client {
	t.Helper()
	return open(t, login, false)
}

// OpenCluster is Open as a client of the pinned server's one-node cluster, for
// the calls that ask every node.
func OpenCluster(t testing.TB, login Login) *Client {
	t.Helper()
	return open(t, login, true)
}

func open(t testing.TB, login Login, cluster bool) *Client {
	t.Helper()
	c := &Client{login: login, cluster: cluster}
	c.tape = tape.OpenRaw(t, server, Driver, func(request string) (string, error) {
		// A recording reaches the server on its first call, so a test that
		// replays never starts a container.
		return liveDo(t, request)
	})
	return c
}

// Ping sends a ping.
func (c *Client) Ping(ctx context.Context) error {
	answer, err := c.do(ctx, request{Op: opPing})
	if err != nil {
		return err
	}
	return answer.err()
}

// Info reads one INFO section, which the tape holds narrowed to the lines the
// driver reads.
func (c *Client) Info(ctx context.Context, section string) (string, error) {
	answer, err := c.do(ctx, request{Op: opInfo, Section: section})
	if err != nil {
		return "", err
	}
	if err := answer.err(); err != nil {
		return "", err
	}
	if answer.Info == nil {
		return "", errors.New("testcontainer: an INFO answer holding no section")
	}
	return *answer.Info, nil
}

// ConfigGet reads one CONFIG GET parameter.
func (c *Client) ConfigGet(ctx context.Context, parameter string) (map[string]string, error) {
	answer, err := c.do(ctx, request{Op: opConfigGet, Parameter: parameter})
	if err != nil {
		return nil, err
	}
	if err := answer.err(); err != nil {
		return nil, err
	}
	return answer.Config, nil
}

// Do sends commands as one pipeline and returns each command's answer.
func (c *Client) Do(ctx context.Context, commands [][]any) ([]*redis.Cmd, error) {
	answer, err := c.do(ctx, request{Op: opDo, Commands: commands})
	if err != nil {
		return nil, err
	}
	return answer.commands(ctx, commands)
}

// DoEachShard sends commands as one pipeline to every node of the cluster and
// returns each node's answers.
func (c *Client) DoEachShard(ctx context.Context, commands [][]any) ([][]*redis.Cmd, error) {
	answer, err := c.do(ctx, request{Op: opDoEachShard, Commands: commands})
	if err != nil {
		return nil, err
	}
	shards := make([][]*redis.Cmd, len(answer.Shards))
	for i, shard := range answer.Shards {
		cmds, err := answer.commands(ctx, commands, shard)
		if err != nil {
			return nil, err
		}
		shards[i] = cmds
	}
	return shards, nil
}

// Close answers the close the driver makes; the tape holds nothing for it.
func (*Client) Close() error { return nil }

// do sends one call and returns its answer.
func (c *Client) do(ctx context.Context, r request) (answer, error) {
	r.Login, r.Cluster = c.login, c.cluster
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
		return answer{}, fmt.Errorf("testcontainer: unreadable answer to %s: %w", text, err)
	}
	return a, nil
}

// The calls the driver makes of its client, as a tape names them.
const (
	opPing        = "Ping"
	opInfo        = "Info"
	opConfigGet   = "ConfigGet"
	opDo          = "Do"
	opDoEachShard = "DoEachShard"
)

// request is one call as a tape holds it: the client it came from, the
// operation, and its arguments.
type request struct {
	Login     Login   `json:"login"`
	Cluster   bool    `json:"cluster,omitempty"`
	Op        string  `json:"op"`
	Section   string  `json:"section,omitempty"`
	Parameter string  `json:"parameter,omitempty"`
	Commands  [][]any `json:"commands,omitempty"`
}

func (r request) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// answer is one call's answer as a tape holds it: one of the four shapes below,
// or the error the server returned instead.
type answer struct {
	Info     *string           `json:"info,omitempty"`
	Config   map[string]string `json:"config,omitempty"`
	Commands []command         `json:"commands,omitempty"`
	Shards   [][]command       `json:"shards,omitempty"`
	Error    *string           `json:"error,omitempty"`
}

// err returns the error the server answered the call with, as go-redis hands it
// over, or nil.
func (a answer) err() error {
	if a.Error == nil {
		return nil
	}
	return decodeError(*a.Error)
}

// commands returns a call's answers as the commands go-redis hands them to the
// driver, built with the request's arguments.
func (a answer) commands(ctx context.Context, args [][]any, encoded ...[]command) ([]*redis.Cmd, error) {
	answers := a.Commands
	if len(encoded) > 0 {
		answers = encoded[0]
	}
	cmds := make([]*redis.Cmd, len(answers))
	for i, answered := range answers {
		var command []any
		if i < len(args) {
			command = args[i]
		}
		cmd := redis.NewCmd(ctx, command...)
		if answered.Value != nil {
			value, err := answered.Value.decode()
			if err != nil {
				return nil, err
			}
			cmd.SetVal(value)
		}
		if answered.Error != nil {
			cmd.SetErr(decodeError(*answered.Error))
		}
		cmds[i] = cmd
	}
	return cmds, nil
}

// command is one command's answer: its value, or the error the server answered
// it with.
type command struct {
	Value *value  `json:"value,omitempty"`
	Error *string `json:"error,omitempty"`
}

// value is one Redis value as a tape holds it: the kind go-redis decoded and the
// field that kind uses. A nil answer is a kind too, so a missing key and a key
// holding an empty string stay apart.
type value struct {
	Kind   string  `json:"kind"`
	String string  `json:"string,omitempty"`
	Int    int64   `json:"int,omitempty"`
	Float  float64 `json:"float,omitempty"`
	Bool   bool    `json:"bool,omitempty"`
	Array  []value `json:"array,omitempty"`
	Map    []entry `json:"map,omitempty"`
}

// entry is one field of a map answer, kept in key order so the tape's bytes do
// not change between runs.
type entry struct {
	Key   string `json:"key"`
	Value value  `json:"value"`
}

// The kinds a Redis value can have.
const (
	kindNil    = "nil"
	kindString = "string"
	kindInt    = "int"
	kindFloat  = "float"
	kindBool   = "bool"
	kindArray  = "array"
	kindMap    = "map"
)

// encodeValue writes one value as a tape holds it. A kind the tape has no form
// for fails the recording rather than reaching a tape a replay cannot read.
func encodeValue(v any) (value, error) {
	switch x := v.(type) {
	case nil:
		return value{Kind: kindNil}, nil
	case string:
		return value{Kind: kindString, String: x}, nil
	case int64:
		return value{Kind: kindInt, Int: x}, nil
	case float64:
		// JSON has no form for an infinity or a NaN, which RESP3 doubles can
		// carry, so those keep their text.
		if math.IsInf(x, 0) || math.IsNaN(x) {
			return value{Kind: kindFloat, String: strconv.FormatFloat(x, 'g', -1, 64)}, nil
		}
		return value{Kind: kindFloat, Float: x}, nil
	case bool:
		return value{Kind: kindBool, Bool: x}, nil
	case []any:
		list := make([]value, len(x))
		for i, item := range x {
			encoded, err := encodeValue(item)
			if err != nil {
				return value{}, err
			}
			list[i] = encoded
		}
		return value{Kind: kindArray, Array: list}, nil
	case map[any]any:
		entries := make([]entry, 0, len(x))
		for key, item := range x {
			name, ok := key.(string)
			if !ok {
				return value{}, fmt.Errorf("testcontainer: a tape holds no map key of type %T", key)
			}
			encoded, err := encodeValue(item)
			if err != nil {
				return value{}, err
			}
			entries = append(entries, entry{Key: name, Value: encoded})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
		return value{Kind: kindMap, Map: entries}, nil
	default:
		return value{}, fmt.Errorf("testcontainer: a tape holds no %T answer", v)
	}
}

// decode reads the value encodeValue wrote.
func (v value) decode() (any, error) {
	switch v.Kind {
	case kindNil:
		return nil, nil
	case kindString:
		return v.String, nil
	case kindInt:
		return v.Int, nil
	case kindFloat:
		if v.String != "" {
			parsed, err := strconv.ParseFloat(v.String, 64)
			if err != nil {
				return nil, fmt.Errorf("testcontainer: unreadable double %q: %w", v.String, err)
			}
			return parsed, nil
		}
		return v.Float, nil
	case kindBool:
		return v.Bool, nil
	case kindArray:
		list := make([]any, len(v.Array))
		for i, item := range v.Array {
			decoded, err := item.decode()
			if err != nil {
				return nil, err
			}
			list[i] = decoded
		}
		return list, nil
	case kindMap:
		m := make(map[any]any, len(v.Map))
		for _, item := range v.Map {
			decoded, err := item.Value.decode()
			if err != nil {
				return nil, err
			}
			m[item.Key] = decoded
		}
		return m, nil
	default:
		return nil, fmt.Errorf("testcontainer: unknown answer kind %q", v.Kind)
	}
}

// serverError is a Redis error the server answered with: go-redis's own
// spelling for one, so HasErrorPrefix still reads it.
type serverError string

func (e serverError) Error() string { return string(e) }
func (serverError) RedisError()     {}

// decodeError rebuilds the error go-redis hands the driver from its text.
// redis.Nil is the one error a driver compares by identity, so it comes back as
// the sentinel rather than as a new value with its text.
func decodeError(text string) error {
	if text == redis.Nil.Error() {
		return redis.Nil
	}
	return serverError(text)
}

// encodeError writes one error as a tape holds it.
func encodeError(err error) string {
	return err.Error()
}

// isServerError reports whether err is the server's answer rather than a failure
// to reach it.
func isServerError(err error) bool {
	var answered redis.Error
	return errors.As(err, &answered)
}

// unmarshalJSON decodes JSON strictly enough that a tape's typo fails the test
// rather than answering something else.
func unmarshalJSON(b []byte, v any) error {
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// infoLines narrows one INFO section to the lines the driver reads. The counters
// INFO reports move with the server's clock and connections, so a tape that held
// them would change at every recording. A kept line keeps the carriage return
// the server ended it with, so a driver that stops trimming it a replay catches.
func infoLines(section, text string) string {
	var kept []string
	for _, line := range strings.Split(text, "\n") {
		switch section {
		case "server":
			if strings.HasPrefix(line, "redis_version:") {
				kept = append(kept, line)
			}
		case "keyspace":
			if strings.HasPrefix(line, "db") && strings.Contains(line, ":") {
				kept = append(kept, line)
			}
		}
	}
	return strings.Join(kept, "\n")
}

// --- the live server ---

const (
	// clusterSlots is the number of hash slots Redis partitions keys into; the
	// single node this package starts owns all of them.
	clusterSlots = 16384
	// startupTimeout is how long the node takes to answer: a few seconds on an
	// idle host.
	startupTimeout = 3 * time.Minute
)

var liveServer struct {
	once      sync.Once
	container testcontainers.Container
	host      string
	port      string
	err       error

	mu      sync.Mutex
	clients map[liveKey]live
}

// liveKey names one client of the pinned server: the login and whether it is a
// cluster client.
type liveKey struct {
	login   Login
	cluster bool
}

// live is the pinned server as one client asks it.
type live interface {
	Ping(ctx context.Context) error
	Info(ctx context.Context, section string) (string, error)
	ConfigGet(ctx context.Context, parameter string) (map[string]string, error)
	Do(ctx context.Context, commands [][]any) ([]*redis.Cmd, error)
	DoEachShard(ctx context.Context, commands [][]any) ([][]*redis.Cmd, error)
	Close() error
}

// liveClient is the pinned server reached through go-redis, as the driver
// reaches a server.
type liveClient struct {
	client redis.UniversalClient
}

func (c *liveClient) Ping(ctx context.Context) error { return c.client.Ping(ctx).Err() }

func (c *liveClient) Close() error { return c.client.Close() }

func (c *liveClient) Info(ctx context.Context, section string) (string, error) {
	text, err := c.client.Info(ctx, section).Result()
	if err != nil {
		return "", err
	}
	return infoLines(section, text), nil
}

func (c *liveClient) ConfigGet(ctx context.Context, parameter string) (map[string]string, error) {
	return c.client.ConfigGet(ctx, parameter).Result()
}

func (c *liveClient) Do(ctx context.Context, commands [][]any) ([]*redis.Cmd, error) {
	cmds := make([]*redis.Cmd, 0, len(commands))
	if _, err := c.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, command := range commands {
			cmds = append(cmds, p.Do(ctx, command...))
		}
		return nil
	}); err != nil && !isServerError(err) {
		return nil, err
	}
	return cmds, nil
}

func (c *liveClient) DoEachShard(ctx context.Context, commands [][]any) ([][]*redis.Cmd, error) {
	cluster, ok := c.client.(*redis.ClusterClient)
	if !ok {
		return nil, fmt.Errorf("expect *redis.ClusterClient, but get %T", c.client)
	}
	var mu sync.Mutex
	var shards []shardCmds
	err := cluster.ForEachShard(ctx, func(ctx context.Context, client *redis.Client) error {
		cmds := make([]*redis.Cmd, 0, len(commands))
		_, _ = client.Pipelined(ctx, func(p redis.Pipeliner) error {
			for _, command := range commands {
				cmds = append(cmds, p.Do(ctx, command...))
			}
			return nil
		})
		mu.Lock()
		defer mu.Unlock()
		shards = append(shards, shardCmds{addr: client.Options().Addr, cmds: cmds})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(shards, func(i, j int) bool { return shards[i].addr < shards[j].addr })
	answers := make([][]*redis.Cmd, len(shards))
	for i, shard := range shards {
		answers[i] = shard.cmds
	}
	return answers, nil
}

// shardCmds is one node's answers, kept with its address so the tape's order is
// the addresses' rather than the run's.
type shardCmds struct {
	addr string
	cmds []*redis.Cmd
}

// liveDo answers one request from the pinned server, in the encoding a tape
// holds.
func liveDo(t testing.TB, text string) (string, error) {
	t.Helper()
	var req request
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("testcontainer: unreadable request %q: %w", text, err)
	}
	client := liveFor(t, req)
	ctx := context.Background()
	var a answer
	switch req.Op {
	case opPing:
		if err := client.Ping(ctx); err != nil {
			if !isServerError(err) {
				return "", err
			}
			message := encodeError(err)
			a.Error = &message
		}
	case opInfo:
		text, err := client.Info(ctx, req.Section)
		if err != nil {
			if !isServerError(err) {
				return "", err
			}
			message := encodeError(err)
			a.Error = &message
			break
		}
		a.Info = &text
	case opConfigGet:
		config, err := client.ConfigGet(ctx, req.Parameter)
		if err != nil {
			if !isServerError(err) {
				return "", err
			}
			message := encodeError(err)
			a.Error = &message
			break
		}
		a.Config = config
	case opDo:
		cmds, err := client.Do(ctx, req.Commands)
		if err != nil {
			return "", err
		}
		if a.Commands, err = encodeCommands(cmds); err != nil {
			return "", err
		}
	case opDoEachShard:
		shards, err := client.DoEachShard(ctx, req.Commands)
		if err != nil {
			return "", err
		}
		for _, cmds := range shards {
			encoded, err := encodeCommands(cmds)
			if err != nil {
				return "", err
			}
			a.Shards = append(a.Shards, encoded)
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

// encodeCommands writes the answers go-redis gave one pipeline.
func encodeCommands(cmds []*redis.Cmd) ([]command, error) {
	answers := make([]command, len(cmds))
	for i, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			message := encodeError(err)
			answers[i].Error = &message
			continue
		}
		encoded, err := encodeValue(cmd.Val())
		if err != nil {
			return nil, err
		}
		answers[i].Value = &encoded
	}
	return answers, nil
}

// liveFor returns the request's client of the pinned server, starting the
// container on first use and emptying it at the start of each top-level test.
// Only a recording calls it.
func liveFor(t testing.TB, req request) live {
	liveServer.once.Do(func() {
		ctx := context.Background()
		liveServer.container, liveServer.host, liveServer.port, liveServer.err = start(ctx)
	})
	if liveServer.err != nil {
		t.Fatalf("start Redis: %v", liveServer.err)
	}
	liveServer.mu.Lock()
	if liveServer.clients == nil {
		liveServer.clients = map[liveKey]live{}
	}
	liveServer.mu.Unlock()
	client := clientFor(t, req.Login, req.Cluster)
	emptyOnce(t)
	return client
}

// clientFor returns the client one login and topology asks the pinned server
// through, building it on first use.
func clientFor(t testing.TB, login Login, cluster bool) live {
	key := liveKey{login: login, cluster: cluster}
	liveServer.mu.Lock()
	defer liveServer.mu.Unlock()
	if client, ok := liveServer.clients[key]; ok {
		return client
	}
	address := net.JoinHostPort(liveServer.host, liveServer.port)
	var client live
	if cluster {
		client = &liveClient{client: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{address},
			Username: login.User,
			Password: login.Password,
		})}
	} else {
		client = &liveClient{client: redis.NewClient(&redis.Options{
			Addr:     address,
			Username: login.User,
			Password: login.Password,
		})}
	}
	liveServer.clients[key] = client
	return client
}

// emptied records the top-level tests whose first live call has emptied the
// server, so a test reads the same answers recorded alone as in a run of its
// package.
var emptied struct {
	mu   sync.Mutex
	done map[string]bool
}

// emptyOnce drops every key on t's top-level test's first live call, so what an
// earlier test left behind cannot reach the tape.
func emptyOnce(t testing.TB) {
	t.Helper()
	top := enginetest.TopName(t.Name())
	emptied.mu.Lock()
	defer emptied.mu.Unlock()
	if emptied.done[top] {
		return
	}
	if emptied.done == nil {
		emptied.done = map[string]bool{}
	}
	emptied.done[top] = true
	client := clientFor(t, Login{}, false)
	cmds, err := client.Do(context.Background(), [][]any{{"FLUSHALL"}})
	if err != nil {
		t.Fatalf("empty Redis: %v", err)
	}
	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			t.Fatalf("empty Redis: %v", err)
		}
	}
}

// start starts the pinned server and bootstraps its cluster.
func start(ctx context.Context) (testcontainers.Container, string, string, error) {
	var err error
	pin, err := enginetest.LoadPin(Pins, Driver)
	if err != nil {
		return nil, "", "", err
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
		testcontainers.WithCmd("redis-server",
			"--cluster-enabled", "yes",
			"--cluster-config-file", "/tmp/nodes.conf",
			"--cluster-node-timeout", "5000",
		),
		testcontainers.WithExposedPorts("6379/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").WithStartupTimeout(startupTimeout)),
	)
	if err != nil {
		return nil, "", "", err
	}
	if err := checkDigest(ctx, container, pin.Digest); err != nil {
		return nil, "", "", err
	}
	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, "", "", err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", "", err
	}
	// Cluster clients follow the addresses CLUSTER SLOTS reports, so the node
	// is told to announce the host-mapped address rather than its address
	// inside Docker, which a test outside Docker cannot reach.
	announce, err := announceIP(ctx, host)
	if err != nil {
		return nil, "", "", err
	}
	if err := bootstrap(ctx, announce, port.Port(), net.JoinHostPort(host, port.Port())); err != nil {
		return nil, "", "", err
	}
	if err := checkVersion(ctx, net.JoinHostPort(host, port.Port()), pin.Tag); err != nil {
		return nil, "", "", err
	}
	return container, host, port.Port(), nil
}

// announceIP returns the address the node announces. Redis refuses a name
// there, and a host outside Docker may be one.
func announceIP(ctx context.Context, host string) (string, error) {
	if ip := net.ParseIP(host); ip != nil {
		return host, nil
	}
	addresses, err := net.DefaultResolver.LookupHost(ctx, host)
	if err != nil {
		return "", fmt.Errorf("resolve the Redis host %s: %w", host, err)
	}
	for _, address := range addresses {
		if ip := net.ParseIP(address); ip != nil && ip.To4() != nil {
			return address, nil
		}
	}
	if len(addresses) > 0 {
		return addresses[0], nil
	}
	return "", fmt.Errorf("the Redis host %s resolves to no address", host)
}

// bootstrap gives the node every hash slot and the address it announces.
func bootstrap(ctx context.Context, host, port, address string) error {
	client := redis.NewClient(&redis.Options{Addr: address, Protocol: 2})
	defer func() { _ = client.Close() }()
	if err := client.ConfigSet(ctx, "cluster-announce-ip", host).Err(); err != nil {
		return fmt.Errorf("announce the cluster's address: %w", err)
	}
	if err := client.ConfigSet(ctx, "cluster-announce-port", port).Err(); err != nil {
		return fmt.Errorf("announce the cluster's port: %w", err)
	}
	if err := client.ClusterAddSlotsRange(ctx, 0, clusterSlots-1).Err(); err != nil {
		return fmt.Errorf("assign the cluster's slots: %w", err)
	}
	// Owning every slot is not enough on its own: the node reports state fail
	// until its cluster cron has taken the new slots into account.
	deadline := time.Now().Add(2 * time.Minute)
	for {
		info, err := client.ClusterInfo(ctx).Result()
		if err != nil {
			return fmt.Errorf("read the cluster's state: %w", err)
		}
		if strings.Contains(info, "cluster_state:ok") {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("the Redis cluster did not reach state ok: %s", info)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// checkVersion holds the pinned image to the release pins.tsv names.
func checkVersion(ctx context.Context, address, version string) error {
	client := redis.NewClient(&redis.Options{Addr: address, Protocol: 2})
	defer func() { _ = client.Close() }()
	info, err := client.Info(ctx, "server").Result()
	if err != nil {
		return fmt.Errorf("Redis does not answer: %w", err)
	}
	reported := ""
	for _, line := range strings.Split(info, "\n") {
		if after, ok := strings.CutPrefix(line, "redis_version:"); ok {
			reported = strings.TrimSpace(after)
		}
	}
	if reported != version {
		return fmt.Errorf("the server reports %q, but pins.tsv pins Redis %s", reported, version)
	}
	return nil
}

func checkDigest(ctx context.Context, container testcontainers.Container, digest string) error {
	info, err := container.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Redis container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Redis image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Redis container runs an image with digests %v, but pins.tsv pins %s", image.RepoDigests, digest)
}

// Terminate removes the server this run started, so a recording leaves no
// container behind even when the reaper is disabled.
func Terminate() {
	if liveServer.container == nil {
		return
	}
	liveServer.mu.Lock()
	for _, client := range liveServer.clients {
		_ = client.Close()
	}
	liveServer.mu.Unlock()
	_ = testcontainers.TerminateContainer(liveServer.container)
}
