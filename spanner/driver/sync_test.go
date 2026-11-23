package driver

import (
	"context"
	"errors"
	"flag"
	"net"
	"os"
	"path/filepath"
	"testing"

	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// TestSyncDatabase syncs the fixture database and holds the result to
// testdata/sync/fixture.json: the tables with their columns, indexes and foreign
// keys, and the view with its definition.
func TestSyncDatabase(t *testing.T) {
	a := require.New(t)
	instance, database := syncFixture(t)
	d := openDriver(t, instance, database)
	defer d.Close()

	got, err := d.SyncDatabase(context.Background())
	a.NoError(err)
	syncGolden(t, got, filepath.Join("testdata", "sync", "fixture.json"))
}

// TestSyncDatabasePostgres syncs a database of the PostgreSQL dialect and holds
// the result to testdata/sync/postgres.json: the PostgreSQL catalog's own
// schema name, type names, default expression and match option.
func TestSyncDatabasePostgres(t *testing.T) {
	a := require.New(t)
	instance, database := postgresSyncFixture(t)
	d := openDriverWithDialect(t, instance, database, postgresSQL)
	defer d.Close()

	got, err := d.SyncDatabase(context.Background())
	a.NoError(err)
	syncGolden(t, got, filepath.Join("testdata", "sync", "postgres.json"))
}

// syncGolden holds a synced schema to a golden file, or writes it with -update.
func syncGolden(t *testing.T, got *metadatapb.DatabaseSchemaMetadata, golden string) {
	t.Helper()
	a := require.New(t)
	if *update {
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		a.NoError(err)
		a.NoError(os.MkdirAll(filepath.Dir(golden), 0o755))
		a.NoError(os.WriteFile(golden, append(b, '\n'), 0o644))
		return
	}
	want := readGolden(t, golden)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("SyncDatabase differs from %s (-want +got):\n%s", golden, diff)
	}
}

// TestSyncDatabaseMissing holds a database the server does not have to
// driver.ErrNotFound, which is what a caller reports as a missing database
// rather than a failed sync.
func TestSyncDatabaseMissing(t *testing.T) {
	a := require.New(t)
	instance, _ := syncFixture(t)
	d := openDriver(t, instance, "missing")
	defer d.Close()

	_, err := d.SyncDatabase(context.Background())
	a.ErrorIs(err, driver.ErrNotFound)
}

// TestSyncInstance reads the instance's databases with their names and no
// schemas, both dialects included: a database of the PostgreSQL dialect is a
// database, and SyncDatabase reads its catalog.
func TestSyncInstance(t *testing.T) {
	a := require.New(t)
	instance := "syncinstance"
	testcontainer.Create(t, instance, "fixture", fixtureDDL, nil)
	testcontainer.CreatePostgres(t, instance, "pgdialect", nil, nil)
	d := openDriver(t, instance, "fixture")
	defer d.Close()

	got, err := d.SyncInstance(context.Background())
	a.NoError(err)
	var names []string
	for _, database := range got.Databases {
		names = append(names, database.Name)
	}
	a.ElementsMatch([]string{"fixture", "pgdialect"}, names)
	a.Empty(got.Version)
	a.Empty(got.Roles)
}

// TestSyncDatabaseEmpty reads a database that holds no schema: it has no table
// and no view, and so reports no schema.
func TestSyncDatabaseEmpty(t *testing.T) {
	a := require.New(t)
	testcontainer.Create(t, "empty", "empty", nil, nil)
	d := openDriver(t, "empty", "empty")
	defer d.Close()

	got, err := d.SyncDatabase(context.Background())
	a.NoError(err)
	a.Equal("empty", got.Name)
	a.Empty(got.Schemas)
}

// readGolden reads a metadata golden file.
func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	metadata := &metadatapb.DatabaseSchemaMetadata{}
	require.NoError(t, protojson.Unmarshal(b, metadata))
	return metadata
}

type stubListDatabasesClient struct {
	client
	databases []Database
	err       error
}

func (c *stubListDatabasesClient) ListDatabases(context.Context) ([]Database, error) {
	return c.databases, c.err
}

func TestSyncInstanceDialectFilter(t *testing.T) {
	ctx := context.Background()
	mixed := []Database{
		{Name: "projects/p/instances/i/databases/gsql_a", Dialect: DialectGoogleSQL},
		{Name: "projects/p/instances/i/databases/pg_a", Dialect: DialectPostgresSQL},
		{Name: "projects/p/instances/i/databases/default_db", Dialect: ""},
		{Name: "projects/p/instances/i/databases/unspecified_db", Dialect: databasepb.DatabaseDialect_DATABASE_DIALECT_UNSPECIFIED.String()},
		{Name: "projects/p/instances/i/databases/pg_b", Dialect: DialectPostgresSQL},
		{Name: "projects/p/instances/i/databases/gsql_b", Dialect: DialectGoogleSQL},
	}

	t.Run("default keeps all dialects in order", func(t *testing.T) {
		d := &Driver{client: &stubListDatabasesClient{databases: mixed}}
		got, err := d.SyncInstance(ctx)
		require.NoError(t, err)
		require.Equal(t, []string{"gsql_a", "pg_a", "default_db", "unspecified_db", "pg_b", "gsql_b"}, databaseNames(got))
	})

	t.Run("GoogleSQLOnly skips PostgreSQL databases and preserves order", func(t *testing.T) {
		d := &Driver{
			client:        &stubListDatabasesClient{databases: mixed},
			googleSQLOnly: true,
		}
		got, err := d.SyncInstance(ctx)
		require.NoError(t, err)
		require.Equal(t, []string{"gsql_a", "default_db", "unspecified_db", "gsql_b"}, databaseNames(got))
	})

	t.Run("GoogleSQLOnly skips malformed DSN on filtered PostgreSQL database", func(t *testing.T) {
		d := &Driver{
			client: &stubListDatabasesClient{databases: []Database{
				{Name: "not-a-dsn", Dialect: DialectPostgresSQL},
				{Name: "projects/p/instances/i/databases/gsql_a", Dialect: DialectGoogleSQL},
			}},
			googleSQLOnly: true,
		}
		got, err := d.SyncInstance(ctx)
		require.NoError(t, err)
		require.Equal(t, []string{"gsql_a"}, databaseNames(got))
	})

	t.Run("malformed DSN on kept database returns error", func(t *testing.T) {
		d := &Driver{
			client: &stubListDatabasesClient{databases: []Database{
				{Name: "not-a-dsn", Dialect: DialectGoogleSQL},
			}},
			googleSQLOnly: true,
		}
		_, err := d.SyncInstance(ctx)
		require.ErrorContains(t, err, "failed to get database name from not-a-dsn")
	})

	t.Run("ListDatabases error propagates", func(t *testing.T) {
		wantErr := errors.New("permission denied")
		for _, googleSQLOnly := range []bool{false, true} {
			d := &Driver{
				client:        &stubListDatabasesClient{err: wantErr},
				googleSQLOnly: googleSQLOnly,
			}
			_, err := d.SyncInstance(ctx)
			require.ErrorIs(t, err, wantErr)
		}
	})
}

type fakeDatabaseAdminServer struct {
	databasepb.UnimplementedDatabaseAdminServer
	databases []*databasepb.Database
}

func (s *fakeDatabaseAdminServer) ListDatabases(_ context.Context, _ *databasepb.ListDatabasesRequest) (*databasepb.ListDatabasesResponse, error) {
	return &databasepb.ListDatabasesResponse{Databases: s.databases}, nil
}

// TestOpenGoogleSQLOnly exercises Options.GoogleSQLOnly through Open and the
// live gRPC admin client against an in-process DatabaseAdminServer.
func TestOpenGoogleSQLOnly(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv := grpc.NewServer()
	databasepb.RegisterDatabaseAdminServer(srv, &fakeDatabaseAdminServer{
		databases: []*databasepb.Database{
			{Name: "projects/p/instances/i/databases/gsql_db", DatabaseDialect: databasepb.DatabaseDialect_GOOGLE_STANDARD_SQL},
			{Name: "projects/p/instances/i/databases/pg_db", DatabaseDialect: databasepb.DatabaseDialect_POSTGRESQL},
			{Name: "projects/p/instances/i/databases/default_db", DatabaseDialect: databasepb.DatabaseDialect_DATABASE_DIALECT_UNSPECIFIED},
		},
	})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	host, port, err := net.SplitHostPort(lis.Addr().String())
	require.NoError(t, err)
	cfg := driver.Config{
		Host: host,
		Port: port,
		Params: map[string]string{
			ParamProject:  "p",
			ParamInstance: "i",
			ParamEmulator: "1",
		},
	}

	ctx := context.Background()
	allDriver, err := Open(ctx, cfg)
	require.NoError(t, err)
	defer allDriver.Close()
	allInst, err := allDriver.SyncInstance(ctx)
	require.NoError(t, err)
	require.Equal(t, []string{"gsql_db", "pg_db", "default_db"}, databaseNames(allInst))

	filteredDriver, err := Open(ctx, cfg, Options{GoogleSQLOnly: true})
	require.NoError(t, err)
	defer filteredDriver.Close()
	filteredInst, err := filteredDriver.SyncInstance(ctx)
	require.NoError(t, err)
	require.Equal(t, []string{"gsql_db", "default_db"}, databaseNames(filteredInst))
}

func databaseNames(inst *driver.Instance) []string {
	names := make([]string, 0, len(inst.Databases))
	for _, db := range inst.Databases {
		names = append(names, db.Name)
	}
	return names
}
