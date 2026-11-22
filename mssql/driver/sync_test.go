package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	gomssqldb "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// openTestDriver returns a driver on the database of the top-level test, with
// the corpus applied.
func openTestDriver(t *testing.T) *Driver {
	t.Helper()
	db := testcontainer.Open(t)
	t.Cleanup(func() { _ = db.Close() })
	return &Driver{db: db, databaseName: testcontainer.Database(t)}
}

// execCorpus applies a script of GO-separated batches.
func execCorpus(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	var batch strings.Builder
	flush := func() {
		text := strings.TrimSpace(batch.String())
		batch.Reset()
		if text == "" {
			return
		}
		if _, err := db.ExecContext(context.Background(), text); err != nil {
			t.Fatalf("%s:\n%v", text, err)
		}
	}
	for _, line := range strings.Split(script, "\n") {
		if strings.EqualFold(strings.TrimSpace(line), "GO") {
			flush()
			continue
		}
		batch.WriteString(line)
		batch.WriteString("\n")
	}
	flush()
}

func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	execCorpus(t, d.db, string(corpus))

	got, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	golden := filepath.Join("testdata", "sync", "database.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want := readGolden(t, golden)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("SyncDatabase differs from %s (-want +got):\n%s", golden, diff)
	}

	var sequences, indexes int
	for _, schema := range got.Schemas {
		if schema.Name != "app" {
			continue
		}
		for _, sequence := range schema.Sequences {
			switch sequence.Name {
			case "seq_orders":
				sequences++
				if sequence.Start != "10" || sequence.Increment != "2" || sequence.MinValue != "1" ||
					sequence.MaxValue != "10000" || !sequence.Cycle || sequence.CacheSize != "50" {
					t.Errorf("seq_orders synced without its options: %+v", sequence)
				}
			case "seq_decimal":
				sequences++
				if sequence.DataType != "decimal(38,0)" {
					t.Errorf("seq_decimal synced without its precision: %q", sequence.DataType)
				}
			}
		}
		for _, table := range schema.Tables {
			if table.Name != "accounts" {
				continue
			}
			for _, index := range table.Indexes {
				if index.Name != "idx_accounts_status" {
					continue
				}
				indexes++
				if len(index.Expressions) != 1 || index.Expressions[0] != "status" {
					t.Errorf("idx_accounts_status synced its INCLUDE column as a key: %v", index.Expressions)
				}
				if !strings.Contains(index.Definition, "INCLUDE:email") || !strings.Contains(index.Definition, "WHERE (") {
					t.Errorf("idx_accounts_status synced without its INCLUDE or predicate: %q", index.Definition)
				}
			}
		}
	}
	if sequences != 2 || indexes != 1 {
		t.Fatalf("the corpus objects are missing from the sync: %d sequences, %d indexes", sequences, indexes)
	}
}

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, m); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return m
}

func TestSyncInstance(t *testing.T) {
	d := openTestDriver(t)
	if _, err := d.db.ExecContext(context.Background(), "CREATE DATABASE "+quoteIdent("omni_instance_extra")); err != nil {
		t.Fatal(err)
	}
	instance, err := d.SyncInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// The server's own databases and whatever an earlier test left differ
	// between runs, so the golden holds the version and this test's database.
	var names []string
	for _, database := range instance.Databases {
		if database.Name == "omni_instance_extra" || database.Name == testcontainer.Database(t) {
			names = append(names, fmt.Sprintf("%s %s", database.Name, database.Collation))
		}
	}
	out, err := json.MarshalIndent(map[string]any{"version": instance.Version, "databases": names}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got := string(out) + "\n"

	golden := filepath.Join("testdata", "sync", "instance.json")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), got); diff != "" {
		t.Errorf("SyncInstance differs from %s (-want +got):\n%s", golden, diff)
	}
}

// TestSyncSchemaWithApostrophe pins that a schema name carrying an apostrophe
// reaches the catalog queries as one literal.
func TestSyncSchemaWithApostrophe(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE SCHEMA [owner's];
GO
CREATE TABLE [owner's].items (id INT NOT NULL PRIMARY KEY, note NVARCHAR(20));
`)
	metadata, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatalf("a schema named with an apostrophe failed the sync: %v", err)
	}
	var found bool
	for _, schema := range metadata.Schemas {
		if schema.Name != "owner's" {
			continue
		}
		found = true
		if len(schema.Tables) != 1 || schema.Tables[0].Name != "items" || len(schema.Tables[0].Columns) != 2 {
			t.Errorf("schema %q synced %+v", schema.Name, schema.Tables)
		}
	}
	if !found {
		t.Fatalf("the schema is missing from %+v", metadata.Schemas)
	}
}

// TestSyncDatabaseMissing holds that a database SQL Server says the login may
// not open is driver.ErrNotFound, on the connection the pool cannot make: a
// database that is not there is refused when the connection opens, before any
// statement reaches the server.
func TestSyncDatabaseMissing(t *testing.T) {
	db := testcontainer.OpenMissingDatabase(t, missingDatabase)
	t.Cleanup(func() { _ = db.Close() })

	// What the server answered, with the client's own error type, as a replay
	// holds it. 4063 is the warning a login that may fall back to its default
	// database gets, and go-mssqldb reports it as the connect error.
	err := db.PingContext(context.Background())
	var engineError gomssqldb.Error
	if !errors.As(err, &engineError) || engineError.Number != 4063 {
		t.Fatalf("Ping = %v, want a SQL Server error 4063", err)
	}

	d := &Driver{db: db, databaseName: missingDatabase}
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase = %v, want driver.ErrNotFound", err)
	}
}

// missingDatabase is a database no test creates, so a recording's server does
// not have it.
const missingDatabase = "omni_no_such_database"

// TestWrapSyncError holds the numbers SQL Server answers for a database the
// login may not open, and that a failure of another kind keeps its own error.
func TestWrapSyncError(t *testing.T) {
	for _, number := range []int32{911, 4060, 4063} {
		err := wrapSyncError(missingDatabase, gomssqldb.Error{Number: number, Message: "Cannot open database"})
		if !errors.Is(err, driver.ErrNotFound) {
			t.Errorf("wrapSyncError(%d) = %v, want driver.ErrNotFound", number, err)
		}
	}
	other := gomssqldb.Error{Number: 229, Message: "The SELECT permission was denied"}
	if err := wrapSyncError(missingDatabase, other); errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapSyncError(229) = %v, want the engine's own error", err)
	}
}

func requireSyncedTable(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata, schemaName, tableName string) *metadatapb.TableMetadata {
	t.Helper()
	for _, schema := range metadata.Schemas {
		if schema.Name != schemaName {
			continue
		}
		for _, table := range schema.Tables {
			if table.Name == tableName {
				return table
			}
		}
	}
	t.Fatalf("table %s.%s not found in synced metadata", schemaName, tableName)
	return nil
}

// TestSyncUniqueIndex pins which synced indexes carry Unique: standalone
// CREATE UNIQUE INDEX (is_unique) vs UNIQUE constraint (is_unique_constraint)
// vs PRIMARY KEY vs plain CREATE INDEX.
func TestSyncUniqueIndex(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE TABLE dbo.index_kinds (
	id INT IDENTITY(1,1) PRIMARY KEY,
	code VARCHAR(50) UNIQUE,
	name NVARCHAR(100) NOT NULL,
	status INT,
	tag NVARCHAR(50)
);
GO
CREATE UNIQUE INDEX idx_unique_name ON dbo.index_kinds(name);
CREATE UNIQUE INDEX idx_unique_filtered ON dbo.index_kinds(status) WHERE status IS NOT NULL;
CREATE INDEX idx_plain_tag ON dbo.index_kinds(tag);
CREATE UNIQUE INDEX idx_unique_desc ON dbo.index_kinds(name DESC, status ASC);
GO
`)

	metadata, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}

	table := requireSyncedTable(t, metadata, "dbo", "index_kinds")
	byName := make(map[string]*metadatapb.IndexMetadata, len(table.Indexes))
	for _, index := range table.Indexes {
		byName[index.Name] = index
	}

	for _, name := range []string{"idx_unique_name", "idx_unique_filtered", "idx_unique_desc"} {
		index := byName[name]
		if index == nil {
			t.Fatalf("index %s not synced", name)
		}
		if !index.Unique || index.Primary || index.IsConstraint {
			t.Errorf("%s = (unique=%v primary=%v isConstraint=%v), want (true, false, false)", name, index.Unique, index.Primary, index.IsConstraint)
		}
	}

	plain := byName["idx_plain_tag"]
	if plain == nil || plain.Unique {
		t.Errorf("idx_plain_tag = %+v, want non-unique index", plain)
	}

	var columnConstraint *metadatapb.IndexMetadata
	for name, index := range byName {
		if index.IsConstraint && !index.Primary {
			if columnConstraint != nil {
				t.Fatalf("expected one unique constraint, also got %s", name)
			}
			columnConstraint = index
		}
	}
	if columnConstraint == nil || !columnConstraint.Unique || !cmp.Equal([]string{"code"}, columnConstraint.Expressions) {
		t.Fatalf("unique constraint = %+v, want unique on [code]", columnConstraint)
	}

	var primary *metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if index.Primary {
			primary = index
		}
	}
	if primary == nil || !primary.Unique || !cmp.Equal([]string{"id"}, primary.Expressions) {
		t.Fatalf("primary key = %+v, want unique primary on [id]", primary)
	}

	if diff := cmp.Diff([]string{"name", "status"}, byName["idx_unique_desc"].Expressions); diff != "" {
		t.Errorf("idx_unique_desc expressions (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]bool{true, false}, byName["idx_unique_desc"].Descending); diff != "" {
		t.Errorf("idx_unique_desc descending (-want +got):\n%s", diff)
	}
}

// TestSyncSpatialIndex pins what SyncDatabase reports for SQL Server spatial
// indexes: tessellation scheme and grid levels, bounding box, storage options,
// dimensional descriptor, mixed index types, and computed columns.
func TestSyncSpatialIndex(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name     string
		setupSQL string
		validate func(*testing.T, *metadatapb.DatabaseSchemaMetadata)
	}{
		{
			name: "geometry_and_geography_indexes",
			setupSQL: `
CREATE SCHEMA geo;
GO
CREATE TABLE geo.locations (
	id INT IDENTITY(1,1) PRIMARY KEY,
	name NVARCHAR(200) NOT NULL,
	location_point GEOMETRY NOT NULL,
	boundary_polygon GEOMETRY,
	geo_location GEOGRAPHY NOT NULL,
	route_line GEOGRAPHY
);
GO
CREATE SPATIAL INDEX idx_location_point ON geo.locations(location_point)
USING GEOMETRY_GRID
WITH (
	BOUNDING_BOX = (-180, -90, 180, 90),
	GRIDS = (LEVEL_1 = MEDIUM, LEVEL_2 = HIGH, LEVEL_3 = MEDIUM, LEVEL_4 = LOW),
	CELLS_PER_OBJECT = 32,
	PAD_INDEX = ON,
	FILLFACTOR = 90,
	ALLOW_ROW_LOCKS = ON,
	ALLOW_PAGE_LOCKS = ON
);
GO
CREATE SPATIAL INDEX idx_boundary_polygon ON geo.locations(boundary_polygon)
USING GEOMETRY_GRID
WITH (
	BOUNDING_BOX = (0, 0, 1000, 1000),
	GRIDS = (LEVEL_1 = LOW, LEVEL_2 = LOW, LEVEL_3 = HIGH, LEVEL_4 = HIGH),
	CELLS_PER_OBJECT = 64
);
GO
CREATE SPATIAL INDEX idx_geo_location ON geo.locations(geo_location)
USING GEOGRAPHY_GRID
WITH (
	GRIDS = (LEVEL_1 = LOW, LEVEL_2 = MEDIUM, LEVEL_3 = HIGH, LEVEL_4 = MEDIUM),
	CELLS_PER_OBJECT = 16,
	FILLFACTOR = 85,
	PAD_INDEX = OFF,
	ALLOW_ROW_LOCKS = OFF,
	ALLOW_PAGE_LOCKS = OFF
);
GO
CREATE SPATIAL INDEX idx_route_line ON geo.locations(route_line)
USING GEOGRAPHY_GRID
WITH (
	GRIDS = (LEVEL_1 = HIGH, LEVEL_2 = HIGH, LEVEL_3 = HIGH, LEVEL_4 = HIGH),
	CELLS_PER_OBJECT = 128
);
GO
`,
			validate: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				table := requireSyncedTable(t, metadata, "geo", "locations")
				if len(table.Columns) != 6 {
					t.Fatalf("columns = %d, want 6", len(table.Columns))
				}
				columns := make(map[string]*metadatapb.ColumnMetadata, len(table.Columns))
				for _, column := range table.Columns {
					columns[column.Name] = column
				}
				for name, wantType := range map[string]string{
					"location_point":   "geometry",
					"boundary_polygon": "geometry",
					"geo_location":     "geography",
					"route_line":       "geography",
				} {
					if got := strings.ToLower(columns[name].Type); got != wantType {
						t.Errorf("column %s type = %q, want %q", name, got, wantType)
					}
				}
				if got := len(spatialIndexes(table)); got != 4 {
					t.Fatalf("spatial indexes = %d, want 4", got)
				}

				lp := requireSpatialIndex(t, table, "idx_location_point", "location_point", "GEOMETRY")
				assertSpatialCells(t, lp, 32)
				assertBoundingBox(t, lp, -180, -90, 180, 90)
				assertGridLevels(t, lp, "MEDIUM", "HIGH", "MEDIUM", "LOW")
				if s := lp.SpatialConfig.Storage; s.Fillfactor != 90 || !s.PadIndex || !s.AllowRowLocks || !s.AllowPageLocks {
					t.Errorf("idx_location_point storage = %+v", s)
				}

				bp := requireSpatialIndex(t, table, "idx_boundary_polygon", "boundary_polygon", "GEOMETRY")
				assertSpatialCells(t, bp, 64)
				assertBoundingBox(t, bp, 0, 0, 1000, 1000)
				assertGridLevels(t, bp, "LOW", "LOW", "HIGH", "HIGH")

				gl := requireSpatialIndex(t, table, "idx_geo_location", "geo_location", "GEOGRAPHY")
				assertSpatialCells(t, gl, 16)
				assertGridLevels(t, gl, "LOW", "MEDIUM", "HIGH", "MEDIUM")
				if s := gl.SpatialConfig.Storage; s.Fillfactor != 85 || s.PadIndex || s.AllowRowLocks || s.AllowPageLocks {
					t.Errorf("idx_geo_location storage = %+v", s)
				}

				rl := requireSpatialIndex(t, table, "idx_route_line", "route_line", "GEOGRAPHY")
				assertSpatialCells(t, rl, 128)
				assertGridLevels(t, rl, "HIGH", "HIGH", "HIGH", "HIGH")
			},
		},
		{
			name: "configuration_bounds",
			setupSQL: `
CREATE TABLE dbo.spatial_edge_cases (
	id INT IDENTITY(1,1) PRIMARY KEY,
	min_config_geom GEOMETRY NOT NULL,
	max_config_geom GEOMETRY NOT NULL,
	min_config_geog GEOGRAPHY NOT NULL,
	max_config_geog GEOGRAPHY NOT NULL
);
GO
CREATE SPATIAL INDEX idx_min_geom ON dbo.spatial_edge_cases(min_config_geom)
USING GEOMETRY_GRID
WITH (
	BOUNDING_BOX = (0, 0, 1, 1),
	GRIDS = (LEVEL_1 = LOW, LEVEL_2 = LOW, LEVEL_3 = LOW, LEVEL_4 = LOW),
	CELLS_PER_OBJECT = 1
);
GO
CREATE SPATIAL INDEX idx_max_geom ON dbo.spatial_edge_cases(max_config_geom)
USING GEOMETRY_GRID
WITH (
	BOUNDING_BOX = (-1000000, -1000000, 1000000, 1000000),
	GRIDS = (LEVEL_1 = HIGH, LEVEL_2 = HIGH, LEVEL_3 = HIGH, LEVEL_4 = HIGH),
	CELLS_PER_OBJECT = 8192,
	PAD_INDEX = ON,
	FILLFACTOR = 100,
	ALLOW_ROW_LOCKS = ON,
	ALLOW_PAGE_LOCKS = ON,
	MAXDOP = 4,
	DATA_COMPRESSION = PAGE
);
GO
CREATE SPATIAL INDEX idx_min_geog ON dbo.spatial_edge_cases(min_config_geog)
USING GEOGRAPHY_GRID
WITH (
	GRIDS = (LEVEL_1 = LOW, LEVEL_2 = LOW, LEVEL_3 = LOW, LEVEL_4 = LOW),
	CELLS_PER_OBJECT = 1
);
GO
CREATE SPATIAL INDEX idx_max_geog ON dbo.spatial_edge_cases(max_config_geog)
USING GEOGRAPHY_GRID
WITH (
	GRIDS = (LEVEL_1 = HIGH, LEVEL_2 = HIGH, LEVEL_3 = HIGH, LEVEL_4 = HIGH),
	CELLS_PER_OBJECT = 8192,
	FILLFACTOR = 50
);
GO
`,
			validate: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				table := requireSyncedTable(t, metadata, "dbo", "spatial_edge_cases")
				if got := len(spatialIndexes(table)); got != 4 {
					t.Fatalf("spatial indexes = %d, want 4", got)
				}

				minGeom := requireSpatialIndex(t, table, "idx_min_geom", "min_config_geom", "GEOMETRY")
				assertSpatialCells(t, minGeom, 1)
				assertBoundingBox(t, minGeom, 0, 0, 1, 1)
				assertGridLevels(t, minGeom, "LOW", "LOW", "LOW", "LOW")

				maxGeom := requireSpatialIndex(t, table, "idx_max_geom", "max_config_geom", "GEOMETRY")
				assertSpatialCells(t, maxGeom, 8192)
				assertBoundingBox(t, maxGeom, -1000000, -1000000, 1000000, 1000000)
				assertGridLevels(t, maxGeom, "HIGH", "HIGH", "HIGH", "HIGH")
				if s := maxGeom.SpatialConfig.Storage; s.Fillfactor != 100 || !s.PadIndex || !s.AllowRowLocks || !s.AllowPageLocks {
					t.Errorf("idx_max_geom storage = %+v", s)
				}

				minGeog := requireSpatialIndex(t, table, "idx_min_geog", "min_config_geog", "GEOGRAPHY")
				assertSpatialCells(t, minGeog, 1)
				assertGridLevels(t, minGeog, "LOW", "LOW", "LOW", "LOW")

				maxGeog := requireSpatialIndex(t, table, "idx_max_geog", "max_config_geog", "GEOGRAPHY")
				assertSpatialCells(t, maxGeog, 8192)
				assertGridLevels(t, maxGeog, "HIGH", "HIGH", "HIGH", "HIGH")
				if got := maxGeog.SpatialConfig.Storage.Fillfactor; got != 50 {
					t.Errorf("idx_max_geog fillfactor = %d, want 50", got)
				}
			},
		},
		{
			name: "mixed_index_types",
			setupSQL: `
CREATE TABLE dbo.mixed_indexes (
	id INT IDENTITY(1,1) PRIMARY KEY,
	name NVARCHAR(100) NOT NULL,
	code VARCHAR(50) UNIQUE,
	location GEOMETRY NOT NULL,
	area GEOGRAPHY,
	status INT,
	created_at DATETIME2 DEFAULT GETDATE()
);
GO
CREATE INDEX idx_name ON dbo.mixed_indexes(name);
CREATE INDEX idx_status_created ON dbo.mixed_indexes(status, created_at);
CREATE UNIQUE INDEX idx_name_status ON dbo.mixed_indexes(name, status) WHERE status IS NOT NULL;
GO
CREATE SPATIAL INDEX idx_location ON dbo.mixed_indexes(location)
USING GEOMETRY_GRID
WITH (
	BOUNDING_BOX = (-100, -100, 100, 100),
	GRIDS = (LEVEL_1 = MEDIUM, LEVEL_2 = MEDIUM, LEVEL_3 = MEDIUM, LEVEL_4 = MEDIUM),
	CELLS_PER_OBJECT = 16
);
GO
CREATE SPATIAL INDEX idx_area ON dbo.mixed_indexes(area)
USING GEOGRAPHY_GRID
WITH (
	GRIDS = (LEVEL_1 = HIGH, LEVEL_2 = LOW, LEVEL_3 = HIGH, LEVEL_4 = LOW),
	CELLS_PER_OBJECT = 32
);
GO
CREATE NONCLUSTERED COLUMNSTORE INDEX idx_columnstore ON dbo.mixed_indexes(id, name, status);
GO
`,
			validate: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				table := requireSyncedTable(t, metadata, "dbo", "mixed_indexes")
				counts := make(map[string]int)
				for _, index := range table.Indexes {
					switch {
					case index.Primary:
						counts["PRIMARY"]++
					case index.Type == "SPATIAL":
						counts["SPATIAL"]++
					case index.Type == "NONCLUSTERED COLUMNSTORE":
						counts["COLUMNSTORE"]++
					case index.Unique:
						counts["UNIQUE"]++
					default:
						counts["REGULAR"]++
					}
				}
				want := map[string]int{
					"PRIMARY":     1,
					"UNIQUE":      2,
					"REGULAR":     2,
					"SPATIAL":     2,
					"COLUMNSTORE": 1,
				}
				if diff := cmp.Diff(want, counts); diff != "" {
					t.Errorf("index kind counts (-want +got):\n%s", diff)
				}

				location := requireSpatialIndex(t, table, "idx_location", "location", "GEOMETRY")
				assertSpatialCells(t, location, 16)
				assertBoundingBox(t, location, -100, -100, 100, 100)
				assertGridLevels(t, location, "MEDIUM", "MEDIUM", "MEDIUM", "MEDIUM")

				area := requireSpatialIndex(t, table, "idx_area", "area", "GEOGRAPHY")
				assertSpatialCells(t, area, 32)
				assertGridLevels(t, area, "HIGH", "LOW", "HIGH", "LOW")
			},
		},
		{
			name: "computed_columns",
			setupSQL: `
CREATE TABLE dbo.spatial_computed (
	id INT IDENTITY(1,1) PRIMARY KEY,
	shape GEOMETRY NOT NULL,
	area AS shape.STArea() PERSISTED,
	perimeter AS shape.STLength() PERSISTED,
	center_x AS shape.STCentroid().STX PERSISTED,
	center_y AS shape.STCentroid().STY PERSISTED,
	as_text AS shape.STAsText() PERSISTED
);
GO
CREATE SPATIAL INDEX idx_shape ON dbo.spatial_computed(shape)
USING GEOMETRY_GRID
WITH (
	BOUNDING_BOX = (-500, -500, 500, 500),
	GRIDS = (LEVEL_1 = LOW, LEVEL_2 = MEDIUM, LEVEL_3 = HIGH, LEVEL_4 = MEDIUM),
	CELLS_PER_OBJECT = 24,
	FILLFACTOR = 75
);
GO
CREATE INDEX idx_area ON dbo.spatial_computed(area);
CREATE INDEX idx_center ON dbo.spatial_computed(center_x, center_y);
GO
`,
			validate: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				table := requireSyncedTable(t, metadata, "dbo", "spatial_computed")
				columns := make(map[string]*metadatapb.ColumnMetadata, len(table.Columns))
				for _, column := range table.Columns {
					columns[column.Name] = column
				}
				for _, name := range []string{"area", "perimeter", "center_x", "center_y", "as_text"} {
					if _, ok := columns[name]; !ok {
						t.Errorf("computed column %s missing", name)
					}
				}

				shape := requireSpatialIndex(t, table, "idx_shape", "shape", "GEOMETRY")
				assertSpatialCells(t, shape, 24)
				assertBoundingBox(t, shape, -500, -500, 500, 500)
				assertGridLevels(t, shape, "LOW", "MEDIUM", "HIGH", "MEDIUM")
				if got := shape.SpatialConfig.Storage.Fillfactor; got != 75 {
					t.Errorf("idx_shape fillfactor = %d, want 75", got)
				}

				regular := 0
				for _, index := range table.Indexes {
					if index.Type != "SPATIAL" && !index.Primary && !index.Unique {
						regular++
					}
				}
				if regular != 2 {
					t.Errorf("regular indexes over computed columns = %d, want 2", regular)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := openTestDriver(t)
			execCorpus(t, d.db, tc.setupSQL)
			metadata, err := d.SyncDatabase(ctx)
			if err != nil {
				t.Fatal(err)
			}
			tc.validate(t, metadata)
		})
	}
}

// TestSyncSpatialIndexIsRepeatable pins that a second sync of an unchanged
// database reports the same spatial indexes as the first.
func TestSyncSpatialIndexIsRepeatable(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE SCHEMA spatial_test;
GO
CREATE TABLE spatial_test.mixed_spatial (
	id INT IDENTITY(1,1) PRIMARY KEY,
	name NVARCHAR(100) NOT NULL,
	geo_point GEOGRAPHY NOT NULL,
	geo_line GEOGRAPHY NOT NULL,
	created_at DATETIME2 DEFAULT GETDATE()
);
GO
CREATE SPATIAL INDEX idx_geo_point ON spatial_test.mixed_spatial(geo_point)
USING GEOGRAPHY_GRID
WITH (
	GRIDS = (LEVEL_1 = MEDIUM, LEVEL_2 = HIGH, LEVEL_3 = MEDIUM, LEVEL_4 = LOW),
	CELLS_PER_OBJECT = 16
);
GO
CREATE SPATIAL INDEX idx_geo_line ON spatial_test.mixed_spatial(geo_line)
USING GEOGRAPHY_GRID
WITH (
	GRIDS = (LEVEL_1 = LOW, LEVEL_2 = LOW, LEVEL_3 = MEDIUM, LEVEL_4 = HIGH),
	CELLS_PER_OBJECT = 8
);
GO
`)

	first, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	table := requireSyncedTable(t, first, "spatial_test", "mixed_spatial")
	if got := len(spatialIndexes(table)); got != 2 {
		t.Fatalf("spatial indexes = %d, want 2", got)
	}
	geoPoint := requireSpatialIndex(t, table, "idx_geo_point", "geo_point", "GEOGRAPHY")
	assertSpatialCells(t, geoPoint, 16)
	assertGridLevels(t, geoPoint, "MEDIUM", "HIGH", "MEDIUM", "LOW")

	geoLine := requireSpatialIndex(t, table, "idx_geo_line", "geo_line", "GEOGRAPHY")
	assertSpatialCells(t, geoLine, 8)
	assertGridLevels(t, geoLine, "LOW", "LOW", "MEDIUM", "HIGH")

	second, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	byName := func(metadata *metadatapb.DatabaseSchemaMetadata) map[string]*metadatapb.IndexMetadata {
		indexes := make(map[string]*metadatapb.IndexMetadata)
		for _, index := range spatialIndexes(requireSyncedTable(t, metadata, "spatial_test", "mixed_spatial")) {
			indexes[index.Name] = index
		}
		return indexes
	}
	if diff := cmp.Diff(byName(first), byName(second), protocmp.Transform()); diff != "" {
		t.Errorf("second sync differs from first (-first +second):\n%s", diff)
	}
}

func spatialIndexes(table *metadatapb.TableMetadata) []*metadatapb.IndexMetadata {
	var indexes []*metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if index.Type == "SPATIAL" {
			indexes = append(indexes, index)
		}
	}
	return indexes
}

func requireSpatialIndex(t *testing.T, table *metadatapb.TableMetadata, name, column, dataType string) *metadatapb.IndexMetadata {
	t.Helper()
	var index *metadatapb.IndexMetadata
	for _, candidate := range table.Indexes {
		if candidate.Name == name {
			index = candidate
			break
		}
	}
	if index == nil {
		t.Fatalf("spatial index %s not synced", name)
	}
	if index.Type != "SPATIAL" || index.Unique || index.Primary {
		t.Errorf("%s = (type=%q unique=%v primary=%v), want (SPATIAL, false, false)", name, index.Type, index.Unique, index.Primary)
	}
	if !cmp.Equal([]string{column}, index.Expressions) || !cmp.Equal([]bool{false}, index.Descending) {
		t.Errorf("%s keys = %v desc=%v, want [%s] [false]", name, index.Expressions, index.Descending, column)
	}
	sc := index.SpatialConfig
	if sc == nil || sc.Method != "SPATIAL" || sc.Tessellation == nil || sc.Storage == nil || sc.Dimensional == nil {
		t.Fatalf("%s incomplete SpatialConfig: %+v", name, sc)
	}
	if sc.Tessellation.Scheme != dataType+"_GRID" {
		t.Errorf("%s scheme = %q, want %q", name, sc.Tessellation.Scheme, dataType+"_GRID")
	}
	if sc.Dimensional.DataType != dataType || sc.Dimensional.Dimensions != 2 {
		t.Errorf("%s dimensional = %+v, want (%s, 2)", name, sc.Dimensional, dataType)
	}
	if dataType == "GEOGRAPHY" && sc.Tessellation.BoundingBox != nil {
		t.Errorf("%s GEOGRAPHY index should not carry a bounding box", name)
	}
	return index
}

func assertSpatialCells(t *testing.T, index *metadatapb.IndexMetadata, want int32) {
	t.Helper()
	if got := index.SpatialConfig.Tessellation.CellsPerObject; got != want {
		t.Errorf("%s cells_per_object = %d, want %d", index.Name, got, want)
	}
}

func assertBoundingBox(t *testing.T, index *metadatapb.IndexMetadata, xmin, ymin, xmax, ymax float64) {
	t.Helper()
	bbox := index.SpatialConfig.Tessellation.BoundingBox
	if bbox == nil {
		t.Fatalf("%s missing bounding box", index.Name)
	}
	if bbox.Xmin != xmin || bbox.Ymin != ymin || bbox.Xmax != xmax || bbox.Ymax != ymax {
		t.Errorf("%s bbox = (%v,%v,%v,%v), want (%v,%v,%v,%v)", index.Name, bbox.Xmin, bbox.Ymin, bbox.Xmax, bbox.Ymax, xmin, ymin, xmax, ymax)
	}
}

func assertGridLevels(t *testing.T, index *metadatapb.IndexMetadata, level1, level2, level3, level4 string) {
	t.Helper()
	densities := make(map[int32]string)
	for _, level := range index.SpatialConfig.Tessellation.GridLevels {
		densities[level.Level] = level.Density
	}
	want := map[int32]string{1: level1, 2: level2, 3: level3, 4: level4}
	if diff := cmp.Diff(want, densities); diff != "" {
		t.Errorf("%s grid levels (-want +got):\n%s", index.Name, diff)
	}
}

// TestSyncDefaultSchemaAndSynonyms pins multi-part synonym target parsing and
// that DefaultSchema follows the syncing user's default schema rather than
// hardcoding dbo.
func TestSyncDefaultSchemaAndSynonyms(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t)
	execCorpus(t, d.db, `
CREATE SCHEMA app;
GO
CREATE SYNONYM dbo.s1 FOR obj3;
CREATE SYNONYM dbo.s4 FOR srv1.db1.sch1.obj1;
CREATE SYNONYM app.s2 FOR sch2.obj2;
GO
`)

	metadata, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if metadata.DefaultSchema != "dbo" {
		t.Fatalf("DefaultSchema = %q, want dbo", metadata.DefaultSchema)
	}
	synonyms := make(map[string][]*metadatapb.SynonymMetadata, len(metadata.Schemas))
	for _, schema := range metadata.Schemas {
		synonyms[schema.Name] = schema.Synonyms
	}
	want := map[string][]*metadatapb.SynonymMetadata{
		"dbo": {
			{Name: "s1", TargetName: "obj3"},
			{Name: "s4", TargetServer: "srv1", TargetDatabase: "db1", TargetSchema: "sch1", TargetName: "obj1"},
		},
		"app": {{Name: "s2", TargetSchema: "sch2", TargetName: "obj2"}},
	}
	if diff := cmp.Diff(want, synonyms, protocmp.Transform(), cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("synonyms (-want +got):\n%s", diff)
	}

	// Impersonate a non-sysadmin user whose default schema is app to prove
	// SyncDatabase reads the active user's default schema via SCHEMA_NAME().
	execCorpus(t, d.db, `
CREATE USER app_user WITHOUT LOGIN WITH DEFAULT_SCHEMA = app;
ALTER ROLE db_owner ADD MEMBER app_user;
EXECUTE AS USER = 'app_user';
`)
	defer func() {
		_, _ = d.db.ExecContext(ctx, "REVERT;")
	}()

	metadata, err = d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if metadata.DefaultSchema != "app" {
		t.Fatalf("DefaultSchema as app_user = %q, want app", metadata.DefaultSchema)
	}
}
