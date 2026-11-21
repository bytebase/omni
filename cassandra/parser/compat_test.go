package parser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	enginerules "github.com/bytebase/omni/enginetest/rules"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/Cql.g4 and its near-misses, with omni's verdict
// and Cassandra's. Cql.g4 is Apache Cassandra's own CQL grammar, ANTLR 3,
// converted to ANTLR 4 by cmd/antlr3to4 (TestCqlG4IsTheConversion holds it to
// that). The test reads the file and starts no server; OMNI_RECORD=1 judges
// every unjudged row on the pinned Cassandra via CQL PREPARE, and
// OMNI_RECORD=verify re-judges every row against the pinned server without
// writing.
//
// testdata/README.md documents the judge's verdict rule.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "Cql.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
	enginerules.Check(t, enginerules.Suite{
		Compat: filepath.Join("testdata", "compat.tsv"),
		Rules:  filepath.Join(td, "rules.tsv"),
		Root:   filepath.Join("..", ".."),
	})
}

// TestGrammarVendorDigest holds the vendored ANTLR 3 grammar files (Cql.g,
// Parser.g, Lexer.g) to the grammar row of ../testdata/pins.tsv, so none is
// edited without updating its pin.
func TestGrammarVendorDigest(t *testing.T) {
	td := filepath.Join("..", "testdata")
	pin, err := enginetest.LoadPin(filepath.Join(td, "pins.tsv"), "grammar")
	if err != nil {
		t.Fatal(err)
	}
	vendor := []string{
		filepath.Join(td, "Cql.g"),
		filepath.Join(td, "Parser.g"),
		filepath.Join(td, "Lexer.g"),
	}
	got, err := enginetest.FilesDigest(vendor...)
	if err != nil {
		t.Fatal(err)
	}
	if got != pin.Digest {
		t.Fatalf("pins.tsv grammar row pins %s, and %s hash to %s", pin.Digest, strings.Join(vendor, ", "), got)
	}
}

// omniVerdict accepts a row that parses as exactly one statement.
func omniVerdict(sql string) string {
	list, err := Parse(sql)
	if err != nil || list == nil || len(list.Items) != 1 {
		return compat.Reject
	}
	return compat.Accept
}

const (
	// cassandraStartupTimeout bounds Cassandra's JVM boot, system keyspace
	// initialization, and gossip settlement.
	cassandraStartupTimeout = 5 * time.Minute
	cassandraConnectTimeout = 30 * time.Second
	cassandraPrepareTimeout = 30 * time.Second

	// cassandraProbeTail is appended to a statement whose PREPARE returned a
	// non-syntax RequestError (0x2200 Invalid, 0x2100 Unauthorized, 0x2300
	// Config, 0x2400 AlreadyExists) to distinguish parse-time rejections from
	// post-parse checks.
	//
	// In Cql.g, query is `cqlStatement (';')* EOF`. The leading newline closes
	// any trailing single-line comment (-- or //), ';' ends cqlStatement, and
	// ')' cannot follow ';' before EOF:
	//   - If stmt parsed cleanly to EOF and only failed later in prepare() or
	//     validate(), preparing stmt + cassandraProbeTail fails at ')' with
	//     ErrCodeSyntax (0x2000) -> Accept.
	//   - If a Parser.g inline Java action threw InvalidRequestException
	//     (0x2200) before EOF (e.g. CreateTableStatement.Raw.addColumn,
	//     extendClusteringOrder, FunctionResource.functionFromCql), preparing
	//     stmt + cassandraProbeTail aborts in the same action with
	//     ErrCodeInvalid (0x2200) -> Reject.
	cassandraProbeTail = "\n; )"
)

// compatJudge starts the pinned Cassandra server and judges each statement by
// sending a CQL native-protocol PREPARE frame via gocql's StatementMetadata.
// PREPARE parses and prepares without executing, so no row mutates server
// state.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	session := startPinnedCassandra(t, pin)
	prepare := func(text string) error {
		ctx, cancel := context.WithTimeout(context.Background(), cassandraPrepareTimeout)
		defer cancel()
		_, err := session.StatementMetadata(ctx, text, "")
		return err
	}

	// Sanity-check the error-code contract before judging any row so a server
	// release that changes how PREPARE reports syntax or Parser.g action
	// errors fails immediately rather than silently misclassifying rows.
	assertProbe := func(stmt, want string) {
		got, err := cassandraVerdict(prepare(stmt), func() error {
			return prepare(stmt + cassandraProbeTail)
		})
		if err != nil || got != want {
			t.Fatalf("Cassandra sanity check on %q = (%q, %v), want (%q, nil)", stmt, got, err, want)
		}
	}
	assertProbe("SELECT FROM", compat.Reject)
	assertProbe("CREATE TABLE a(a ASCII, a ASCII)", compat.Reject)
	assertProbe("SELECT a FROM a", compat.Accept)

	return func(stmt string) (string, error) {
		return cassandraVerdict(prepare(stmt), func() error {
			return prepare(stmt + cassandraProbeTail)
		})
	}
}

// cassandraVerdict classifies Cassandra's answer to CQL PREPARE on one
// statement, asking probe (PREPARE of stmt + cassandraProbeTail) when needed
// to separate Parser.g inline action rejections from post-parse errors.
//
// An answer that cannot be read cleanly (transport/timeout error, server or
// protocol error, unexpected code, or ambiguous probe) returns an error so the
// recording fails instead of guessing.
func cassandraVerdict(err error, probe func() error) (string, error) {
	if err == nil {
		return compat.Accept, nil
	}
	var reqErr gocql.RequestError
	if !errors.As(err, &reqErr) {
		return "", err
	}
	switch reqErr.Code() {
	case gocql.ErrCodeSyntax:
		return compat.Reject, nil

	case gocql.ErrCodeInvalid:
		perr := probe()
		if perr == nil {
			return "", fmt.Errorf("probe unexpectedly succeeded after 0x2200 (%s)", reqErr.Message())
		}
		var preq gocql.RequestError
		if !errors.As(perr, &preq) {
			return "", fmt.Errorf("probe failed: %w", perr)
		}
		switch preq.Code() {
		case gocql.ErrCodeSyntax:
			return compat.Accept, nil
		case gocql.ErrCodeInvalid:
			return compat.Reject, nil
		default:
			return "", fmt.Errorf("probe returned unexpected Cassandra error 0x%04x (%s) after 0x2200 (%s)",
				preq.Code(), preq.Message(), reqErr.Message())
		}

	case gocql.ErrCodeUnauthorized, gocql.ErrCodeConfig, gocql.ErrCodeAlreadyExists:
		// Confirm the statement parsed cleanly to EOF before the post-parse
		// auth/config/schema check ran.
		perr := probe()
		if perr == nil {
			return "", fmt.Errorf("probe unexpectedly succeeded after 0x%04x (%s)", reqErr.Code(), reqErr.Message())
		}
		var preq gocql.RequestError
		if !errors.As(perr, &preq) {
			return "", fmt.Errorf("probe failed: %w", perr)
		}
		if preq.Code() == gocql.ErrCodeSyntax {
			return compat.Accept, nil
		}
		return "", fmt.Errorf("probe returned 0x%04x (%s) instead of syntax error after 0x%04x (%s)",
			preq.Code(), preq.Message(), reqErr.Code(), reqErr.Message())

	default:
		return "", fmt.Errorf("unexpected Cassandra error 0x%04x (%s)", reqErr.Code(), reqErr.Message())
	}
}

// startPinnedCassandra starts the container named by pin, verifies its image
// digest and reported release_version against pin, and registers t.Cleanup to
// close the session and remove the container.
func startPinnedCassandra(t testing.TB, pin enginetest.Pin) *gocql.Session {
	t.Helper()
	ctx := context.Background()
	container, err := runCassandraContainer(ctx, pin)
	if err != nil {
		t.Fatalf("start Cassandra (%s): %v", pin.ImageRef(), err)
	}
	t.Cleanup(func() {
		_ = testcontainers.TerminateContainer(container)
	})
	if err := checkCassandraDigest(ctx, container, pin.Digest); err != nil {
		t.Fatalf("verify Cassandra digest: %v", err)
	}
	port, err := container.MappedPort(ctx, "9042/tcp")
	if err != nil {
		t.Fatalf("read Cassandra port: %v", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("read Cassandra host: %v", err)
	}
	cluster := gocql.NewCluster(host + ":" + port.Port())
	cluster.ConnectTimeout = cassandraConnectTimeout
	session, err := cluster.CreateSession()
	if err != nil {
		t.Fatalf("connect to Cassandra: %v", err)
	}
	t.Cleanup(session.Close)

	var version string
	if err := session.Query("SELECT release_version FROM system.local").ScanContext(ctx, &version); err != nil {
		t.Fatalf("read Cassandra release_version: %v", err)
	}
	if version != pin.Tag {
		t.Fatalf("server reports Cassandra %q, but pins.tsv pins %q", version, pin.Tag)
	}
	return session
}

// runCassandraContainer starts the pinned Cassandra container, converting a
// testcontainers Docker-socket panic into an error.
func runCassandraContainer(ctx context.Context, pin enginetest.Pin) (container testcontainers.Container, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	return testcontainers.Run(ctx, pin.ImageRef(),
		testcontainers.WithEnv(map[string]string{
			"CASSANDRA_CLUSTER_NAME": "omni",
			"MAX_HEAP_SIZE":          "1G",
			"HEAP_NEWSIZE":           "200M",
		}),
		testcontainers.WithExposedPorts("9042/tcp"),
		testcontainers.WithWaitStrategyAndDeadline(cassandraStartupTimeout,
			wait.ForLog("Starting listening for CQL clients").WithStartupTimeout(cassandraStartupTimeout)),
	)
}

// checkCassandraDigest proves the running container uses the pinned digest.
func checkCassandraDigest(ctx context.Context, container testcontainers.Container, digest string) error {
	info, err := container.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("container image digests %v do not include %s", image.RepoDigests, digest)
}

type fakeRequestError struct {
	code int
	msg  string
}

func (e fakeRequestError) Code() int       { return e.code }
func (e fakeRequestError) Message() string { return e.msg }
func (e fakeRequestError) Error() string   { return e.msg }

func TestCassandraVerdict(t *testing.T) {
	syntaxErr := fakeRequestError{code: gocql.ErrCodeSyntax, msg: "line 1:0 no viable alternative"}
	invalidErr := fakeRequestError{code: gocql.ErrCodeInvalid, msg: "No keyspace has been specified"}
	actionInvalidErr := fakeRequestError{code: gocql.ErrCodeInvalid, msg: "Duplicate column 'a' declaration for table 'a'"}
	unauthErr := fakeRequestError{code: gocql.ErrCodeUnauthorized, msg: "You have to be logged in"}
	configErr := fakeRequestError{code: gocql.ErrCodeConfig, msg: "config error"}
	existsErr := fakeRequestError{code: gocql.ErrCodeAlreadyExists, msg: "already exists"}
	serverErr := fakeRequestError{code: gocql.ErrCodeServer, msg: "internal error"}
	overloadedErr := fakeRequestError{code: gocql.ErrCodeOverloaded, msg: "overloaded"}
	unpreparedErr := fakeRequestError{code: gocql.ErrCodeUnprepared, msg: "unprepared"}
	unknownErr := fakeRequestError{code: 0x9999, msg: "unknown"}
	transportErr := context.DeadlineExceeded

	cases := []struct {
		name           string
		err            error
		probeErr       error
		wantVerdict    string
		wantErr        bool
		wantProbeCalls int
	}{
		{
			name:        "prepare succeeds",
			err:         nil,
			wantVerdict: compat.Accept,
		},
		{
			name:    "non-RequestError fails",
			err:     transportErr,
			wantErr: true,
		},
		{
			name:        "syntax error rejects without probing",
			err:         syntaxErr,
			wantVerdict: compat.Reject,
		},
		{
			name:           "post-parse invalid accepts when probe hits syntax error",
			err:            invalidErr,
			probeErr:       syntaxErr,
			wantVerdict:    compat.Accept,
			wantProbeCalls: 1,
		},
		{
			name:           "Parser.g action invalid rejects when probe stays invalid",
			err:            actionInvalidErr,
			probeErr:       actionInvalidErr,
			wantVerdict:    compat.Reject,
			wantProbeCalls: 1,
		},
		{
			name:           "invalid fails when probe unexpectedly succeeds",
			err:            invalidErr,
			probeErr:       nil,
			wantErr:        true,
			wantProbeCalls: 1,
		},
		{
			name:           "invalid fails when probe returns transport error",
			err:            invalidErr,
			probeErr:       transportErr,
			wantErr:        true,
			wantProbeCalls: 1,
		},
		{
			name:           "invalid fails when probe returns unexpected Cassandra code",
			err:            invalidErr,
			probeErr:       serverErr,
			wantErr:        true,
			wantProbeCalls: 1,
		},
		{
			name:           "unauthorized accepts when probe hits syntax error",
			err:            unauthErr,
			probeErr:       syntaxErr,
			wantVerdict:    compat.Accept,
			wantProbeCalls: 1,
		},
		{
			name:           "config accepts when probe hits syntax error",
			err:            configErr,
			probeErr:       syntaxErr,
			wantVerdict:    compat.Accept,
			wantProbeCalls: 1,
		},
		{
			name:           "already exists accepts when probe hits syntax error",
			err:            existsErr,
			probeErr:       syntaxErr,
			wantVerdict:    compat.Accept,
			wantProbeCalls: 1,
		},
		{
			name:           "unauthorized fails if probe also returns unauthorized",
			err:            unauthErr,
			probeErr:       unauthErr,
			wantErr:        true,
			wantProbeCalls: 1,
		},
		{
			name:           "unauthorized fails if probe unexpectedly succeeds",
			err:            unauthErr,
			probeErr:       nil,
			wantErr:        true,
			wantProbeCalls: 1,
		},
		{
			name:           "unauthorized fails if probe returns transport error",
			err:            unauthErr,
			probeErr:       transportErr,
			wantErr:        true,
			wantProbeCalls: 1,
		},
		{
			name:    "server error fails without probing",
			err:     serverErr,
			wantErr: true,
		},
		{
			name:    "overloaded error fails without probing",
			err:     overloadedErr,
			wantErr: true,
		},
		{
			name:    "unprepared error fails without probing",
			err:     unpreparedErr,
			wantErr: true,
		},
		{
			name:    "unknown error code fails without probing",
			err:     unknownErr,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			got, err := cassandraVerdict(tc.err, func() error {
				calls++
				return tc.probeErr
			})
			if calls != tc.wantProbeCalls {
				t.Errorf("probe calls = %d, want %d", calls, tc.wantProbeCalls)
			}
			if tc.wantErr {
				if err == nil {
					t.Fatalf("cassandraVerdict() = (%q, nil), want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("cassandraVerdict() unexpected error: %v", err)
			}
			if got != tc.wantVerdict {
				t.Errorf("cassandraVerdict() = %q, want %q", got, tc.wantVerdict)
			}
		})
	}
}
