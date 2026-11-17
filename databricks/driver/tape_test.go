package driver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/databricks/internal/testcontainer"
	metadatapb "github.com/bytebase/omni/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// openTape returns a driver whose client answers from t's tape. It is how the
// recorded tests reach the pinned workspace's answers with no environment and
// no service.
func openTape(t *testing.T) *Driver {
	t.Helper()
	return &Driver{
		client:      testcontainer.Open(t),
		warehouseID: testcontainer.WarehouseID,
		catalog:     testcontainer.Catalog,
	}
}

// The recorded tests. Their names are the tape names: one top-level test per
// file in testdata/tapes.
func TestSyncInstance(t *testing.T) {
	checkSyncInstance(t, openTape(t))
}

func TestSyncDatabase(t *testing.T) {
	d := openTape(t)
	got := checkSyncDatabase(t, d)
	body, err := os.ReadFile(filepath.Join("testdata", "sync", d.catalog+".json"))
	if err != nil {
		t.Fatalf("read sync golden: %v", err)
	}
	want := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(body, want); err != nil {
		t.Fatalf("read sync golden: %v", err)
	}
	if !proto.Equal(got, want) {
		t.Errorf("SyncDatabase = %v, want testdata/sync/%s.json", got, d.catalog)
	}
}

func TestQuery(t *testing.T) {
	checkQuery(t, openTape(t))
}

func TestExecute(t *testing.T) {
	checkExecute(t, openTape(t), "omni_driver_probe")
}

func TestDump(t *testing.T) {
	checkDump(t, openTape(t))
}

func TestPing(t *testing.T) {
	checkPing(t, openTape(t))
}
