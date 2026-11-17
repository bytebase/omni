package driver

import (
	"bytes"
	"context"
	"slices"
	"testing"

	"github.com/bytebase/omni/driver"
)

// spatialKindsSetup declares a column of every spatial type MySQL names, each
// holding a value of its declared type.
const spatialKindsSetup = `DROP TABLE IF EXISTS qgeo;
CREATE TABLE qgeo (
  id BIGINT PRIMARY KEY,
  g GEOMETRY,
  p POINT,
  l LINESTRING,
  poly POLYGON,
  mp MULTIPOINT,
  ml MULTILINESTRING,
  mpoly MULTIPOLYGON,
  gc GEOMETRYCOLLECTION
);
INSERT INTO qgeo VALUES (1,
  ST_GeomFromText('POINT(1 2)'),
  ST_GeomFromText('POINT(1 2)'),
  ST_GeomFromText('LINESTRING(1 2,3 4)'),
  ST_GeomFromText('POLYGON((0 0,1 0,1 1,0 0))'),
  ST_GeomFromText('MULTIPOINT((1 2),(3 4))'),
  ST_GeomFromText('MULTILINESTRING((1 2,3 4))'),
  ST_GeomFromText('MULTIPOLYGON(((0 0,1 0,1 1,0 0)))'),
  ST_GeomFromText('GEOMETRYCOLLECTION(POINT(1 2))'));`

// TestQuerySpatialKinds pins that MySQL names every spatial column GEOMETRY,
// whatever the column was declared with: the decode path needs only that name,
// and a narrower one never fires.
func TestQuerySpatialKinds(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "default", "omni_qe_spatial")
	sess, err := d.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sess.Close() }()

	setup, err := sess.Query(ctx, spatialKindsSetup, driver.QueryOptions{})
	if err != nil {
		t.Fatalf("create qgeo: %v", err)
	}
	for _, res := range setup {
		if res.Err != nil {
			t.Fatalf("%s: %v", res.Statement, res.Err)
		}
	}
	geo, err := sess.Query(ctx, "SELECT * FROM qgeo ORDER BY id;", driver.QueryOptions{})
	if err != nil || len(geo) != 1 || geo[0].Err != nil {
		t.Fatalf("qgeo query = %+v, %v", geo, err)
	}
	columns, row := geo[0].Columns, geo[0].Rows[0]
	wantTypes := []string{"BIGINT", "GEOMETRY", "GEOMETRY", "GEOMETRY", "GEOMETRY", "GEOMETRY", "GEOMETRY", "GEOMETRY", "GEOMETRY"}
	var types []string
	for _, c := range columns {
		types = append(types, c.TypeName)
	}
	if !slices.Equal(types, wantTypes) {
		t.Fatalf("qgeo types = %v, want %v", types, wantTypes)
	}
	for i := 1; i < len(columns); i++ {
		if got := row[i]; got.Kind != driver.Bytes || len(got.Bytes) == 0 {
			t.Errorf("%q = %+v, want bytes", columns[i].Name, got)
		}
	}
	if !bytes.Equal(row[2].Bytes, row[1].Bytes) {
		t.Errorf("POINT(1 2) in %q and %q differ: %x and %x", columns[1].Name, columns[2].Name, row[1].Bytes, row[2].Bytes)
	}
}
