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
// refusedPairs lists, by "<release>/<pair>", the diffPairs a release refuses
// to create, and why. The oracles require the server to refuse the pair's
// second schema there and skip it; an entry the server takes fails.
var refusedPairs = map[string]string{
	"10.7.8/checks reading through to a generated column":     refusesCheckOnGenerated107,
	"10.7.8/kept column checking an added generated one":      refusesCheckOnGenerated107,
	"10.7.8/auto key over a column checking a generated one":  refusesCheckOnGenerated107,
	"10.7.8/kept column checking a replaced one":              refusesCheckOnGenerated107,
	"10.7.8/every column replaced, checks reading each other": refusesCheckOnGenerated107,
	"10.7.8/only column replaced by a group":                  refusesCheckOnGenerated107,
	"10.7.8/restore a reader check after a merged add":        refusesCheckOnGenerated107,
}

const refusesCheckOnGenerated107 = "10.7 refuses a column CHECK that reads a generated column, wherever it stands (4029)"

// pairRefused reports whether srv refuses p, after requiring the server to
// refuse p's second schema.
func pairRefused(t *testing.T, srv oracleServer, conn *sql.Conn, database string, p diffPair, seen map[string]bool) bool {
	t.Helper()
	key := srv.name + "/" + p.name
	why, ok := refusedPairs[key]
	if !ok {
		return false
	}
	seen[key] = true
	resetDatabase(t, conn, database)
	for _, stmt := range splitLikeMySQLClient(p.b + ";") {
		if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
			return true
		}
	}
	t.Errorf("%s: listed as refused (%s) but the server takes it: remove the entry", key, why)
	return true
}

func gapFor(gaps map[string]string, srv oracleServer, key string) (string, bool) {
	if gap, ok := gaps[srv.name+"/"+key]; ok {
		return gap, true
	}
	gap, ok := gaps[key]
	return gap, ok
}
