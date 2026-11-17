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
// with a CHECK, a key, a unique index and an index of its own. Names are quoted
// and mixed case, which
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
	// primaryKey is a partitioned table's primary key: the parent's index, or
	// the copy a partition inherits, which is attached to the parent's.
	primaryKey := func(name, parentName string) *metadatapb.IndexMetadata {
		index := &metadatapb.IndexMetadata{
			Name:         name,
			Expressions:  []string{`"City"`, `"Day"`},
			Type:         "btree",
			Primary:      true,
			Unique:       true,
			IsConstraint: true,
			OpclassNames: []string{"int4_ops", "date_ops"},
		}
		if parentName != "" {
			index.ParentIndexName = parentName
			index.ParentIndexSchema = "app"
		}
		return index
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
							primaryKey("Measure_pkey", ""),
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
									primaryKey("Measure 2025_pkey", "Measure_pkey"),
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
											primaryKey("Measure 2025 EU_pkey", "Measure 2025_pkey"),
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
											primaryKey("Measure 2025 Other_pkey", "Measure 2025_pkey"),
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
									primaryKey("Measure 2026_pkey", "Measure_pkey"),
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
							{Name: "id", Type: "integer", Nullable: false},
							{Name: "v", Type: "text", Nullable: false},
						},
						Partitions: []*metadatapb.TablePartitionMetadata{
							{
								Name:       "Shard_0",
								Type:       metadatapb.TablePartitionMetadata_HASH,
								Expression: "HASH (id)",
								Value:      "FOR VALUES WITH (modulus 2, remainder 0)",
								Indexes: []*metadatapb.IndexMetadata{
									// A key the partition holds alone.
									{
										Name:         "Shard_0_pkey",
										Expressions:  []string{"id"},
										Type:         "btree",
										Primary:      true,
										Unique:       true,
										IsConstraint: true,
									},
									{
										Name:         "Shard_0_v_key",
										Expressions:  []string{"v"},
										Type:         "btree",
										Unique:       true,
										IsConstraint: true,
									},
								},
							},
							{
								Name:       "Shard_1",
								Type:       metadatapb.TablePartitionMetadata_HASH,
								Expression: "HASH (id)",
								Value:      "FOR VALUES WITH (modulus 2, remainder 1)",
								Indexes: []*metadatapb.IndexMetadata{
									// A unique index the partition holds alone, which
									// is not a constraint.
									{
										Name:        "Shard_1_v_idx",
										Expressions: []string{"v"},
										Type:        "btree",
										Unique:      true,
									},
								},
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
	// The diff takes a partition's unattached index for the parent's, so the
	// attachment is checked here: plain PostgreSQL keeps a partitioned index
	// invalid until every partition's index is attached to it.
	assertPartitionIndexesAttached(t, got, partitionedMetadata())
}

// assertPartitionIndexesAttached requires every partition index the snapshot
// attaches to a parent index to be attached in c too.
func assertPartitionIndexesAttached(t *testing.T, c *Catalog, meta *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	attached := 0
	for _, schema := range meta.Schemas {
		for _, table := range schema.Tables {
			var walk func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata)
			walk = func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata) {
				for _, partition := range partitions {
					partitionSchema := partition.GetSchema()
					if partitionSchema == "" {
						partitionSchema = parentSchema
					}
					rel := c.GetRelation(partitionSchema, partition.Name)
					if rel == nil {
						t.Errorf("the SDL leaves out partition %q", partition.Name)
						continue
					}
					indexes := make(map[string]*Index)
					for _, index := range c.IndexesOf(rel.OID) {
						indexes[index.Name] = index
					}
					for _, index := range partition.Indexes {
						if index.ParentIndexName == "" {
							continue
						}
						attached++
						got := indexes[index.Name]
						if got == nil {
							t.Errorf("the SDL leaves out index %q of partition %q", index.Name, partition.Name)
							continue
						}
						if got.ParentIndexOID == 0 {
							t.Errorf("the SDL leaves index %q of partition %q unattached to %q", index.Name, partition.Name, index.ParentIndexName)
						}
					}
					walk(partitionSchema, partition.Subpartitions)
				}
			}
			walk(schema.Name, table.Partitions)
		}
	}
	if attached == 0 {
		t.Fatal("the snapshot holds no attached partition index, so nothing checks the attachment")
	}
}

// TestDatabaseDefinitionSDLPartitionQuotedName pins the escaping of the names a
// partition's statement writes: a partition named with a double quote comes out
// with the quote doubled, so the statement names the identifier it stands for.
func TestDatabaseDefinitionSDLPartitionQuotedName(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "app",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "p",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Nullable: true}},
						Partitions: []*metadatapb.TablePartitionMetadata{
							{
								Name:       `part"1`,
								Comment:    "note",
								Expression: "RANGE (id)",
								Value:      "FOR VALUES FROM (1) TO (2)",
							},
						},
					},
				},
			},
		},
	}
	sdl, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: true})
	if err != nil {
		t.Fatalf("DatabaseDefinition: %v", err)
	}
	want := `CREATE TABLE "app"."part""1" PARTITION OF "app"."p" FOR VALUES FROM (1) TO (2);`
	if !strings.Contains(sdl, want) {
		t.Errorf("the SDL does not hold %s, which names the partition it stands for:\n%s", want, sdl)
	}
	wantComment := `COMMENT ON TABLE "app"."part""1" IS 'note';`
	if !strings.Contains(sdl, wantComment) {
		t.Errorf("the SDL does not hold %s, which names the partition its comment is on:\n%s", wantComment, sdl)
	}
}
