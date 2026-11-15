package catalog

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bytebase/omni/mariadb/internal/testcontainer"
)

// oracleServer is a release the oracle tests hold the catalog to: its server
// role in pins.tsv, its Version, and the snapshots Bytebase's sync took on it.
type oracleServer struct {
	name    string
	role    string
	version Version
	fixture string
}

var oracleServers = []oracleServer{
	{"11.8.8", "", MariaDB118, "testdata/metadata_migration/11.8.8.yaml"},
	{"11.4.13", testcontainer.MariaDB114, MariaDB114, "testdata/metadata_migration/11.4.13.yaml"},
	{"10.11.19", testcontainer.MariaDB1011, MariaDB1011, "testdata/metadata_migration/10.11.19.yaml"},
	{"10.7.8", testcontainer.MariaDB107, MariaDB107, "testdata/metadata_migration/10.7.8.yaml"},
}

// forEachServer runs fn on each release's server, on a connection with a
// database of the test's own selected and empty.
func forEachServer(t *testing.T, fn func(t *testing.T, srv oracleServer, conn *sql.Conn, database string)) {
	for _, srv := range oracleServers {
		t.Run(srv.name, func(t *testing.T) {
			conn, err := testcontainer.OpenDatabaseOn(t, srv.role).Conn(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			fn(t, srv, conn, testcontainer.DatabaseName(t))
		})
	}
}

// newCatalogAt is New at release v.
func newCatalogAt(v Version) *Catalog {
	c := New()
	c.SetVersion(v)
	return c
}

// gapFor looks a gap up for a release: an entry under "<release>/<key>" first,
// then one under the key alone, which holds on every release.
func gapFor(gaps map[string]string, srv oracleServer, key string) (string, bool) {
	if gap, ok := gaps[srv.name+"/"+key]; ok {
		return gap, true
	}
	gap, ok := gaps[key]
	return gap, ok
}
