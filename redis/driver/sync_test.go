package driver

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/redis/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata")

// noconfigUser and noconfigPassword are the ACL account the keyspace test signs
// in with: every command but CONFIG.
const (
	noconfigUser     = "noconfig"
	noconfigPassword = "noconfig-pw"
)

// TestSyncDatabase syncs testdata/corpus.redis as written and holds the synced
// database to testdata/sync/database.json. A Redis database holds keys and no
// schema, so the synced database is its number.
func TestSyncDatabase(t *testing.T) {
	ctx := context.Background()
	d := openDriver(t, testcontainer.Login{}, "0")
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.redis"))
	require.NoError(t, err)
	_, err = d.Execute(ctx, string(corpus), driver.ExecuteOptions{})
	require.NoError(t, err)

	got, err := d.SyncDatabase(ctx)
	require.NoError(t, err)
	writeOrCompareGolden(t, filepath.Join("testdata", "sync", "database.json"), got)
}

// TestSyncInstance holds the pinned server's version and databases to
// testdata/sync/instance.json.
func TestSyncInstance(t *testing.T) {
	d := openDriver(t, testcontainer.Login{}, "0")
	instance, err := d.SyncInstance(context.Background())
	require.NoError(t, err)
	databases := make([]string, 0, len(instance.Databases))
	for _, database := range instance.Databases {
		databases = append(databases, database.Name)
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":   instance.Version,
		"databases": databases,
		"roles":     instance.Roles,
	}, "", "  ")
	require.NoError(t, err)
	writeOrCompareBytes(t, filepath.Join("testdata", "sync", "instance.json"), append(b, '\n'))
}

// TestSyncInstanceFromKeyspace covers an account that may not run CONFIG: the
// databases come from the keyspace section, so only a database holding a key is
// reported, and the fail-open database 0 stands when the section is empty.
func TestSyncInstanceFromKeyspace(t *testing.T) {
	ctx := context.Background()
	admin := openDriver(t, testcontainer.Login{}, "0")
	_, err := admin.Execute(ctx, "ACL SETUSER "+noconfigUser+" on >"+noconfigPassword+" ~* +@all -config", driver.ExecuteOptions{})
	require.NoError(t, err)

	limited := openDriver(t, testcontainer.Login{User: noconfigUser, Password: noconfigPassword}, "0")
	_, err = limited.Execute(ctx, "SET keyspace:key value", driver.ExecuteOptions{})
	require.NoError(t, err)
	_, err = limited.Execute(ctx, "CONFIG GET databases", driver.ExecuteOptions{})
	require.ErrorContains(t, err, "NOPERM")

	instance, err := limited.SyncInstance(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, instance.Version)
	require.Len(t, instance.Databases, 1)
	require.Equal(t, "0", instance.Databases[0].Name)
}

// TestSyncInstanceFailsOpenToDatabaseZero covers an account that may not run
// CONFIG and a keyspace that holds nothing: neither source lists a database, so
// every Redis's database 0 stands.
func TestSyncInstanceFailsOpenToDatabaseZero(t *testing.T) {
	ctx := context.Background()
	admin := openDriver(t, testcontainer.Login{}, "0")
	_, err := admin.Execute(ctx, "ACL SETUSER "+noconfigUser+" on >"+noconfigPassword+" ~* +@all -config", driver.ExecuteOptions{})
	require.NoError(t, err)

	limited := openDriver(t, testcontainer.Login{User: noconfigUser, Password: noconfigPassword}, "0")
	instance, err := limited.SyncInstance(ctx)
	require.NoError(t, err)
	require.Len(t, instance.Databases, 1)
	require.Equal(t, "0", instance.Databases[0].Name)
}

// writeOrCompareGolden rewrites the golden under -update and otherwise compares
// the proto it holds to got. The bytes are not compared: protojson inserts an
// extra space after a key at random per build, so a golden recorded by one build
// does not match another's bytes.
func writeOrCompareGolden(t *testing.T, path string, got *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	if *update {
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		require.NoError(t, err)
		writeOrCompareBytes(t, path, append(b, '\n'))
		return
	}
	want := readGolden(t, path)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("%s (-want +got):\n%s", path, diff)
	}
}

func writeOrCompareBytes(t *testing.T, path string, got []byte) {
	t.Helper()
	if *update {
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, got, 0o644))
		return
	}
	want, err := os.ReadFile(path)
	require.NoError(t, err)
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("%s differs (-want +got):\n%s", path, diff)
	}
}

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	database := &metadatapb.DatabaseSchemaMetadata{}
	require.NoError(t, protojson.Unmarshal(b, database))
	return database
}

// refusedError stands in for an error reply from the server, which go-redis
// marks with the RedisError method.
type refusedError string

func (e refusedError) Error() string { return string(e) }
func (refusedError) RedisError()     {}

func TestDatabasesFromConfigReply(t *testing.T) {
	tests := []struct {
		name         string
		val          map[string]string
		err          error
		want         []int
		wantFailopen bool
		wantErr      bool
	}{
		{
			name: "configured count",
			val:  map[string]string{"databases": "3"},
			want: []int{0, 1, 2},
		},
		{
			name:         "reply without the databases setting",
			val:          map[string]string{},
			wantFailopen: true,
		},
		{
			name:         "command removed by the vendor",
			err:          refusedError("ERR unknown command 'CONFIG', with args beginning with: 'GET' 'databases' "),
			wantFailopen: true,
		},
		{
			name:         "user without permission to run CONFIG",
			err:          refusedError("NOPERM User limited has no permissions to run the 'config|get' command"),
			wantFailopen: true,
		},
		{
			name:    "connection error",
			err:     errors.New("dial tcp 127.0.0.1:6379: connect: connection refused"),
			wantErr: true,
		},
		{
			name:    "count that is not a number",
			val:     map[string]string{"databases": "many"},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, failopen, err := databasesFromConfigReply(tc.val, tc.err)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantFailopen, failopen)
		})
	}
}
