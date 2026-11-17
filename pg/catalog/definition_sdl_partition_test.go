package catalog

import (
	"context"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
)

// partitionedMetadata is the synced snapshot of a schema whose tables are
// partitioned every way PostgreSQL offers: a range partition sub-partitioned
// by list with a default partition, a hash partition, a primary key that the
// partition key completes, an index on each partitioned table, a partition
// whose bounds are multi-column with MINVALUE and MAXVALUE, and a partition
// with a CHECK and an index of its own. Names are quoted and mixed case, which
// is what a dump that writes them inside double quotes must keep, and the
// expressions carry the quoting the server answered with.
func partitionedMetadata() *metadatapb.DatabaseSchemaMetadata {
	parentIndex := func(name string, columns ...string) *metadatapb.IndexMetadata {
		return &metadatapb.IndexMetadata{
			Name:        name,
			Expressions: columns,
			Type:        "btree",
		}
	}
	childIndex := func(name, parentName string, columns ...string) *metadatapb.IndexMetadata {
		return &metadatapb.IndexMetadata{
			Name:              name,
			Expressions:       columns,
			Type:              "btree",
			ParentIndexName:   parentName,
			ParentIndexSchema: "app",
		}
	}
	primaryKey := func(name string) *metadatapb.IndexMetadata {
		return &metadatapb.IndexMetadata{
			Name:         name,
			Expressions:  []string{`"City"`, `"Day"`},
			Type:         "btree",
			Primary:      true,
			Unique:       true,
			IsConstraint: true,
			OpclassNames: []string{"int4_ops", "date_ops"},
		}
	}
	check := func(name, expression string) *metadatapb.CheckConstraintMetadata {
		return &metadatapb.CheckConstraintMetadata{Name: name, Expression: expression}
	}

	return &metadatapb.DatabaseSchemaMetadata{
		Name: "shop",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "app",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "Measure",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "City", Type: "integer", Nullable: false},
							{Name: "Day", Type: "date", Nullable: false},
							{Name: "Region", Type: "text", Nullable: false},
						},
						CheckConstraints: []*metadatapb.CheckConstraintMetadata{
							check("Measure_city_check", `("City" > 0)`),
						},
						Indexes: []*metadatapb.IndexMetadata{
							primaryKey("Measure_pkey"),
							parentIndex("measure_city", `"City"`, `"Day"`),
						},
						Partitions: []*metadatapb.TablePartitionMetadata{
							{
								Name:       "Measure 2025",
								Type:       metadatapb.TablePartitionMetadata_RANGE,
								Expression: `RANGE ("Day")`,
								Value:      "FOR VALUES FROM ('2025-01-01') TO ('2026-01-01')",
								CheckConstraints: []*metadatapb.CheckConstraintMetadata{
									// The copy of the parent's CHECK the server clones.
									check("Measure_city_check", `("City" > 0)`),
									// The partition's own.
									check("Measure 2025_city_check", `("City" < 1000)`),
								},
								Indexes: []*metadatapb.IndexMetadata{
									primaryKey("Measure 2025_pkey"),
									childIndex("Measure 2025_city_day_idx", "measure_city", `"City"`, `"Day"`),
								},
								Subpartitions: []*metadatapb.TablePartitionMetadata{
									{
										Name:       "Measure 2025 EU",
										Type:       metadatapb.TablePartitionMetadata_LIST,
										Expression: `LIST ("Region")`,
										Value:      "FOR VALUES IN ('eu')",
										CheckConstraints: []*metadatapb.CheckConstraintMetadata{
											check("Measure_city_check", `("City" > 0)`),
											check("Measure 2025_city_check", `("City" < 1000)`),
										},
										Indexes: []*metadatapb.IndexMetadata{
											primaryKey("Measure 2025 EU_pkey"),
											childIndex("Measure 2025 EU_city_day_idx", "Measure 2025_city_day_idx", `"City"`, `"Day"`),
											// A partition's own index.
											parentIndex("Measure 2025 EU_region_idx", `"Region"`),
										},
									},
									{
										Name:       "Measure 2025 Other",
										Type:       metadatapb.TablePartitionMetadata_LIST,
										Expression: `LIST ("Region")`,
										Value:      "DEFAULT",
										CheckConstraints: []*metadatapb.CheckConstraintMetadata{
											check("Measure_city_check", `("City" > 0)`),
											check("Measure 2025_city_check", `("City" < 1000)`),
										},
										Indexes: []*metadatapb.IndexMetadata{
											primaryKey("Measure 2025 Other_pkey"),
											childIndex("Measure 2025 Other_city_day_idx", "Measure 2025_city_day_idx", `"City"`, `"Day"`),
										},
									},
								},
							},
							{
								Name:       "Measure 2026",
								Type:       metadatapb.TablePartitionMetadata_RANGE,
								Expression: `RANGE ("Day")`,
								Value:      "FOR VALUES FROM ('2026-01-01') TO ('2027-01-01')",
								CheckConstraints: []*metadatapb.CheckConstraintMetadata{
									check("Measure_city_check", `("City" > 0)`),
								},
								Indexes: []*metadatapb.IndexMetadata{
									primaryKey("Measure 2026_pkey"),
									childIndex("Measure 2026_city_day_idx", "measure_city", `"City"`, `"Day"`),
								},
							},
							{
								Name:       "Measure 2027",
								Type:       metadatapb.TablePartitionMetadata_RANGE,
								Expression: `RANGE ("City", "Day")`,
								Value:      "FOR VALUES FROM (MINVALUE, '2027-01-01') TO (MAXVALUE, MAXVALUE)",
								CheckConstraints: []*metadatapb.CheckConstraintMetadata{
									check("Measure_city_check", `("City" > 0)`),
								},
							},
						},
					},
					{
						Name: "Shard",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer", Nullable: true},
						},
						Partitions: []*metadatapb.TablePartitionMetadata{
							{
								Name:       "Shard_0",
								Type:       metadatapb.TablePartitionMetadata_HASH,
								Expression: "HASH (id)",
								Value:      "FOR VALUES WITH (modulus 2, remainder 0)",
							},
							{
								Name:       "Shard_1",
								Type:       metadatapb.TablePartitionMetadata_HASH,
								Expression: "HASH (id)",
								Value:      "FOR VALUES WITH (modulus 2, remainder 1)",
							},
						},
					},
				},
			},
		},
	}
}

// TestDatabaseDefinitionSDLPartitionRoundTrip holds the SDL rendering of a
// partitioned schema to the snapshot it was rendered from: every partition, its
// bound, and its index attachment the renderer writes must load back into the
// catalog LoadMetadata installs from the same snapshot, so a schema written as
// SDL and the database it was read from differ in nothing.
func TestDatabaseDefinitionSDLPartitionRoundTrip(t *testing.T) {
	meta := partitionedMetadata()
	want := New()
	if _, err := want.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true}); err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}

	sdl, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: true})
	if err != nil {
		t.Fatalf("DatabaseDefinition: %v", err)
	}
	t.Run("single file", func(t *testing.T) {
		assertSDLReloads(t, want, sdl)
	})

	files, err := MultiFileDefinition(meta)
	if err != nil {
		t.Fatalf("MultiFileDefinition: %v", err)
	}
	t.Run("multi file", func(t *testing.T) {
		var joined strings.Builder
		for _, file := range files {
			joined.WriteString(file.Content)
			joined.WriteString("\n")
		}
		assertSDLReloads(t, want, joined.String())
	})
}

// assertSDLReloads loads sdl and holds it to want: the migration that takes one
// to the other must be empty.
func assertSDLReloads(t *testing.T, want *Catalog, sdl string) {
	t.Helper()
	got, err := LoadSDL(sdl)
	if err != nil {
		t.Fatalf("LoadSDL: %v\nSDL:\n%s", err, sdl)
	}
	if sql := GenerateMigration(want, got, Diff(want, got)).SQL(); sql != "" {
		t.Errorf("the SDL does not reload the snapshot, migrating it would:\n%s\nSDL:\n%s", sql, sdl)
	}
}
