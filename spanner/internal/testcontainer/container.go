// Package testcontainer is how spanner's driver tests reach Spanner: through a
// tape (enginetest/tape) that replays a pinned server's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned server itself. The server is the
// Cloud Spanner emulator, one testcontainer per test binary, started by the
// first test that records and removed by Terminate, which the package's TestMain
// calls. A replay starts nothing and needs no Docker.
//
// Spanner's client is not a database/sql driver, so the tape cannot sit at that
// seam. The driver's own tests answer the client seam (spanner/driver's client)
// instead: a request is one call the driver makes of its clients — the admin
// API's listing, get, DDL update, creation and DDL read, a single-use read, a
// snapshot's query, an update, a batch update, a partitioned update, a plan ask,
// a ping and a close — and an answer is the databases it listed, the DDL it
// read, the rows a read returned or the client's error. Those tests own the
// exchange's spelling because they own the seam's types.
//
// The emulator answers the admin API, DDL, DML, queries and INFORMATION_SCHEMA,
// so sync, query, execute, dump and ping are all recorded from it. Three facts a
// real instance has and an emulator does not are never recorded, because it
// cannot answer them: IAM and the credential chain (it accepts any caller), a
// Private Service Connect endpoint and multi-region routing (it is one process
// on one address), and query plans (it has no optimizer). The driver's package
// doc states each.
//
// Create creates a GoogleSQL database from a backtick-quoted CREATE DATABASE.
// CreatePostgres creates a PostgreSQL-dialect one, which the driver skips and
// refuses; it is there so a test records the database beside a GoogleSQL one.
//
// A database is created per fixture, through the admin API, before the first
// request of a recording: a replay answers as if it were there, and the fixture
// is not part of what the driver asked. The database's id is the caller's, so
// two tests hold their own state and two recordings of one test write the same
// resource names.
package testcontainer

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	spannerdb "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is spanner/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin role of the one emulator the tests ask.
const Default = "default"

// Port is the port the emulator serves gRPC on; its REST gateway is 9020 and
// answers no error code the driver reads, so nothing here uses it.
const Port = "9010"

// Project is the Google Cloud project every fixture lives in. The emulator
// accepts any project, and creating an instance is part of a fixture rather than
// of what the driver asked.
const Project = "omni"

// Server is the tape server the driver's recordings are written to. An answer
// carries its own error, so there is no error spelling to encode; what a tape
// needs is the request ID a client puts in an error's text, which names one
// request.
var Server = &tape.Server{Pins: Pins, Codec: codec{}}

type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return fmt.Errorf("%s", text) } //nolint:revive // an answer's text, not a cause to wrap

// madeUp matches the request ID a client puts in an error's text, which differs
// between two runs of one request. The text reaches the tape inside a JSON
// string, so its quotes are escaped.
var madeUp = regexp.MustCompile(`requestID = \\"[^"\\]*\\"`)

func (codec) MadeUp() *regexp.Regexp { return madeUp }

// StandIn is a made-up value's stand-in, counting n from 1: a request ID that
// says it is one, so two recordings of one answer write the same bytes.
func (codec) StandIn(_ string, n int) string {
	return fmt.Sprintf(`requestID = \"omni-request-%d\"`, n)
}

type server struct {
	sync.Mutex
	container testcontainers.Container
	host      string
	instances map[string]bool
	databases map[string]bool
	err       error
}

var shared = &server{instances: map[string]bool{}, databases: map[string]bool{}}

// Host starts the pinned emulator on first use and returns its gRPC address. It
// reaches the server, so only a recording calls it.
func Host(t testing.TB) string {
	t.Helper()
	s := shared
	s.Lock()
	defer s.Unlock()
	if s.host == "" && s.err == nil {
		s.host, s.container, s.err = startContainer(t)
	}
	if s.err != nil {
		t.Fatalf("start the Spanner emulator: %v", s.err)
	}
	return s.host
}

// InstancePath is the resource name of the instance a fixture lives on.
func InstancePath(instance string) string {
	return fmt.Sprintf("projects/%s/instances/%s", Project, instance)
}

// DatabasePath is the resource name of database on instance.
func DatabasePath(instance, database string) string {
	return fmt.Sprintf("%s/databases/%s", InstancePath(instance), database)
}

// Create returns the resource name of database on instance, creating the
// instance on first use and the database with statements as its DDL and rows as
// the DML it then runs, in a recording. A replay starts no server: the database
// is not there, and the tape answers as if it were, so the name alone is what
// the test needs.
func Create(t testing.TB, instance, database string, statements []string, rows []string) string {
	t.Helper()
	if !enginetest.Recording(t) {
		return DatabasePath(instance, database)
	}
	host := Host(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := ensureDatabase(ctx, host, instance, database, statements, rows); err != nil {
		t.Fatalf("create database %s: %v", database, err)
	}
	return DatabasePath(instance, database)
}

// ensureDatabase creates database on instance, once per run, with statements as
// its DDL and rows as the DML it then runs, and a PostgreSQL dialect when
// postgres asks for one.
func ensureDatabase(ctx context.Context, host, instance, database string, statements, rows []string) error {
	shared.Lock()
	known := shared.databases[DatabasePath(instance, database)]
	shared.Unlock()
	if known {
		return nil
	}
	if err := ensureInstance(ctx, host, instance); err != nil {
		return err
	}
	admin, err := databaseAdmin(ctx, host)
	if err != nil {
		return err
	}
	defer func() { _ = admin.Close() }()
	op, err := admin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          InstancePath(instance),
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", database),
		ExtraStatements: statements,
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}

	if len(rows) > 0 {
		if err := writeRows(ctx, host, DatabasePath(instance, database), rows); err != nil {
			return err
		}
	}
	shared.Lock()
	shared.databases[DatabasePath(instance, database)] = true
	shared.Unlock()
	return nil
}

// writeRows runs rows as one read-write transaction's DML, so a fixture's data
// is there before the first recorded request reads it.
func writeRows(ctx context.Context, host, databasePath string, rows []string) error {
	client, err := spanner.NewClient(ctx, databasePath, emulatorOptions(host)...)
	if err != nil {
		return err
	}
	defer client.Close()
	if _, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, rwt *spanner.ReadWriteTransaction) error {
		for _, row := range rows {
			if _, err := rwt.Update(ctx, spanner.NewStatement(row)); err != nil {
				return fmt.Errorf("write a fixture row: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("write the fixture rows of %s: %w", databasePath, err)
	}
	return nil
}

// ensureInstance creates instance, once per run: the emulator starts with no
// instance at all.
func ensureInstance(ctx context.Context, host, instanceName string) error {
	shared.Lock()
	known := shared.instances[instanceName]
	shared.Unlock()
	if known {
		return nil
	}
	admin, err := instanceAdmin(ctx, host)
	if err != nil {
		return err
	}
	defer func() { _ = admin.Close() }()
	op, err := admin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", Project),
		InstanceId: instanceName,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/emulator-config", Project),
			DisplayName: instanceName,
			NodeCount:   1,
		},
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil
		}
		return err
	}
	if _, err := op.Wait(ctx); err != nil && status.Code(err) != codes.AlreadyExists {
		return err
	}
	shared.Lock()
	shared.instances[instanceName] = true
	shared.Unlock()
	return nil
}

func instanceAdmin(ctx context.Context, host string) (*instance.InstanceAdminClient, error) {
	return instance.NewInstanceAdminClient(ctx, emulatorOptions(host)...)
}

func databaseAdmin(ctx context.Context, host string) (*spannerdb.DatabaseAdminClient, error) {
	return spannerdb.NewDatabaseAdminClient(ctx, emulatorOptions(host)...)
}

// emulatorOptions is how a request reaches the pinned emulator: its address, no
// credential, and no TLS.
func emulatorOptions(host string) []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint(host),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	}
}

// startContainer runs the pinned image and waits for its gRPC port.
func startContainer(t testing.TB) (host string, container testcontainers.Container, err error) {
	t.Helper()
	// testcontainers panics when Docker is missing; fail the test instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		return "", nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	containerPort := Port + "/tcp"
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.ImageRef(),
			ExposedPorts: []string{containerPort},
			WaitingFor:   wait.ForListeningPort(Port + "/tcp").WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		if c != nil {
			_ = c.Terminate(ctx, testcontainers.StopTimeout(0))
		}
		return "", nil, fmt.Errorf("start the Spanner emulator: %w", err)
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		_ = c.Terminate(ctx, testcontainers.StopTimeout(0))
		return "", nil, err
	}
	address, err := c.Host(ctx)
	if err != nil {
		return "", c, fmt.Errorf("the Spanner emulator's host: %w", err)
	}
	mapped, err := c.MappedPort(ctx, Port+"/tcp")
	if err != nil {
		return "", c, fmt.Errorf("the Spanner emulator's port: %w", err)
	}
	return address + ":" + mapped.Port(), c, nil
}

// checkDigest proves the container runs the pinned image. The pull names
// image@digest, so a moved tag cannot select another image; this check reads
// what the daemon actually started, which is where the recorded answers came
// from.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Spanner emulator container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Spanner emulator image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Spanner emulator runs an image with digests %v, but pins.tsv pins %s", image.RepoDigests, digest)
}

// Terminate removes the emulator this run started, if any. The package's
// TestMain calls it.
func Terminate() {
	shared.Lock()
	container := shared.container
	shared.container = nil
	shared.host = ""
	shared.instances = map[string]bool{}
	shared.databases = map[string]bool{}
	shared.Unlock()
	if container == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_ = container.Terminate(ctx, testcontainers.StopTimeout(0))
	cancel()
}

// CreatePostgres returns the resource name of a database of the PostgreSQL
// dialect named database on instance, creating the instance on first use, the
// database with statements as its DDL and rows as the DML it then runs, in a
// recording. A replay starts no server.
func CreatePostgres(t testing.TB, instance, database string, statements, rows []string) string {
	t.Helper()
	if !enginetest.Recording(t) {
		return DatabasePath(instance, database)
	}
	host := Host(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := ensurePostgresDatabase(ctx, host, instance, database, statements, rows); err != nil {
		t.Fatalf("create the PostgreSQL database %s: %v", database, err)
	}
	return DatabasePath(instance, database)
}

// ensurePostgresDatabase creates a database of the PostgreSQL dialect, with
// statements as its DDL and rows as the DML it then runs.
func ensurePostgresDatabase(ctx context.Context, host, instance, database string, statements, rows []string) error {
	shared.Lock()
	known := shared.databases[DatabasePath(instance, database)]
	shared.Unlock()
	if known {
		return nil
	}
	if err := ensureInstance(ctx, host, instance); err != nil {
		return err
	}
	admin, err := databaseAdmin(ctx, host)
	if err != nil {
		return err
	}
	defer func() { _ = admin.Close() }()
	// A PostgreSQL database's creation statement is written in its own dialect, as
	// is every statement the fixture applies after it.
	op, err := admin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          InstancePath(instance),
		CreateStatement: fmt.Sprintf("CREATE DATABASE \"%s\"", database),
		DatabaseDialect: databasepb.DatabaseDialect_POSTGRESQL,
		ExtraStatements: statements,
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	if len(rows) > 0 {
		if err := writeRows(ctx, host, DatabasePath(instance, database), rows); err != nil {
			return err
		}
	}
	shared.Lock()
	shared.databases[DatabasePath(instance, database)] = true
	shared.Unlock()
	return nil
}
