package catalog

import (
	"context"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func pgTable(name string, rows int64, columns ...*metadatapb.ColumnMetadata) *metadatapb.TableMetadata {
	return &metadatapb.TableMetadata{Name: name, Columns: columns, RowCount: rows}
}

func newPGTestMetadata(tables []*metadatapb.TableMetadata, views []*metadatapb.ViewMetadata) *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:   "public",
				Tables: tables,
				Views:  views,
			},
		},
	}
}

func testMetadataMigration(oldMeta, newMeta *metadatapb.DatabaseSchemaMetadata) (string, error) {
	var sourcePath, targetPath []string
	if oldMeta != nil {
		sourcePath = parseTestSearchPath(oldMeta.GetSearchPath())
	}
	if newMeta != nil {
		targetPath = parseTestSearchPath(newMeta.GetSearchPath())
	}
	return MetadataMigration(context.Background(), oldMeta, newMeta, MetadataMigrationOptions{
		SourceSearchPath: sourcePath,
		TargetSearchPath: targetPath,
		SkipSchemas:      []string{"bbdataarchive"},
	})
}

func parseTestSearchPath(searchPath string) []string {
	if searchPath == "" {
		return nil
	}
	var parts []string
	start := 0
	var quote byte
	for i := 0; i < len(searchPath); i++ {
		switch {
		case quote != 0:
			if searchPath[i] != quote {
				continue
			}
			if i+1 < len(searchPath) && searchPath[i+1] == quote {
				i++
				continue
			}
			quote = 0
		case searchPath[i] == '"' || searchPath[i] == '\'':
			quote = searchPath[i]
		case searchPath[i] == ',':
			parts = append(parts, searchPath[start:i])
			start = i + 1
		}
	}
	parts = append(parts, searchPath[start:])

	var out []string
	for _, token := range parts {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		schema := strings.ToLower(token)
		switch {
		case len(token) >= 2 && token[0] == '"' && token[len(token)-1] == '"':
			schema = strings.ReplaceAll(token[1:len(token)-1], `""`, `"`)
		case len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\'':
			schema = strings.ReplaceAll(token[1:len(token)-1], `''`, `'`)
		}
		if schema != "" {
			out = append(out, schema)
		}
	}
	return out
}

// TestChangedDegradedProcedureIsRefused covers a procedure the target changes while
// the catalog cannot parse either stored definition: the two stand-ins look alike,
// so the comparison has to refuse rather than let the catalog skip the change.
func TestChangedDegradedProcedureIsRefused(t *testing.T) {
	snapshot := func(definition string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "p",
				Signature:  "p()",
				Definition: definition,
			}},
		}}}
	}
	_, err := testMetadataMigration(snapshot("CREATE PROCEDURE public.p() AS @@@old@@@"),
		snapshot("CREATE PROCEDURE public.p() AS @@@new@@@"))
	require.ErrorContains(t, err, "func:public.p|p()")
}

// TestUnchangedUnloadableForeignKeyDoesNotRefuseUnrelatedTableChange covers a
// foreign key the catalog cannot install because its referenced key is not
// unique: when the key itself is unchanged, an unrelated change to its table must
// still migrate. The report key must compare the constraint, not the whole table.
func TestUnchangedUnloadableForeignKeyDoesNotRefuseUnrelatedTableChange(t *testing.T) {
	snapshot := func(withNote bool) *metadatapb.DatabaseSchemaMetadata {
		columns := []*metadatapb.ColumnMetadata{
			{Name: "id", Type: "integer", Nullable: false},
			{Name: "parent_id", Type: "integer", Nullable: false},
		}
		if withNote {
			columns = append(columns, &metadatapb.ColumnMetadata{Name: "note", Type: "text", Nullable: true})
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{
				{Name: "parent", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Nullable: false}}},
				{
					Name:    "child",
					Columns: columns,
					ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
						Name:              "child_parent_id_fkey",
						Columns:           []string{"parent_id"},
						ReferencedSchema:  "public",
						ReferencedTable:   "parent",
						ReferencedColumns: []string{"id"},
					}},
				},
			},
		}}}
	}

	sql, err := testMetadataMigration(snapshot(false), snapshot(true))
	require.NoError(t, err)
	require.Contains(t, sql, "ADD COLUMN note text")
}

// TestUnchangedUnloadableCheckDoesNotRefuseUnrelatedTableChange is the CHECK
// counterpart: an unsupported check expression must not make every unrelated
// change to its table unmigratable.
func TestUnchangedUnloadableCheckDoesNotRefuseUnrelatedTableChange(t *testing.T) {
	snapshot := func(withNote bool) *metadatapb.DatabaseSchemaMetadata {
		columns := []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Nullable: false}}
		if withNote {
			columns = append(columns, &metadatapb.ColumnMetadata{Name: "note", Type: "text", Nullable: true})
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:             "t",
				Columns:          columns,
				CheckConstraints: []*metadatapb.CheckConstraintMetadata{{Name: "t_missing_check", Expression: "missing <> ''"}},
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot(false), snapshot(true))
	require.NoError(t, err)
	require.Contains(t, sql, "ADD COLUMN note text")
}

// TestRestoreIndexDefinitionsFindsPartitionIndexes covers an index owned by a table
// partition: its stored definition has to replace the generated one.
func TestRestoreIndexDefinitionsFindsPartitionIndexes(t *testing.T) {
	definition := `CREATE INDEX "idx_p1" ON "public"."t_p1" USING btree ("id") WITH (fillfactor=70)`
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{{
			Name: "t",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:    "t_p1",
				Indexes: []*metadatapb.IndexMetadata{{Name: "idx_p1", Definition: definition}},
			}},
		}},
	}}}
	require.NotNil(t, indexByName(meta, "public", "public", "t_p1", "idx_p1"), "a partition owns the index")
	require.NotNil(t, indexByName(meta, "public", "", "", "idx_p1"), "the fallback finds a partition index")

	plan := &MigrationPlan{Ops: []MigrationOp{{
		Type:       OpCreateIndex,
		SchemaName: "public",
		ObjectName: "idx_p1",
		SQL:        `CREATE INDEX "idx_p1" ON "public"."t_p1" USING btree ("id")`,
	}}}
	restoreIndexDefinitions(plan, meta)
	require.Equal(t, definition, plan.Ops[0].SQL)
}

// TestPartitionIndexInAnotherSchema covers a partition whose schema differs from
// its parent table's: the metadata keeps it under the parent, so the owner's own
// schema does not hold it.
func TestPartitionIndexInAnotherSchema(t *testing.T) {
	definition := `CREATE INDEX "idx_p" ON "parted"."t_p" USING btree ("id") WITH (fillfactor=70)`
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "main",
		Tables: []*metadatapb.TableMetadata{{
			Name: "t",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:    "t_p",
				Schema:  "parted",
				Indexes: []*metadatapb.IndexMetadata{{Name: "idx_p", Definition: definition}},
			}},
		}},
	}, {Name: "parted"}}}
	require.NotNil(t, indexByName(meta, "parted", "parted", "t_p", "idx_p"),
		"the partition is held under its parent's schema")
	require.NotNil(t, indexByName(meta, "parted", "", "", "idx_p"))

	plan := &MigrationPlan{Ops: []MigrationOp{{
		Type:       OpCreateIndex,
		SchemaName: "parted",
		ObjectName: "idx_p",
		SQL:        `CREATE INDEX "idx_p" ON "parted"."t_p" USING btree ("id")`,
	}}}
	restoreIndexDefinitions(plan, meta)
	require.Equal(t, definition, plan.Ops[0].SQL)
}

// TestPartitionsAreTypeObjectsAndIndexOwners covers a partition as a relation of its
// own: a type name can resolve to it, and its index can own a comment.
func TestPartitionsAreTypeObjectsAndIndexOwners(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "main",
		Tables: []*metadatapb.TableMetadata{{
			Name: "t",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:    "t_p",
				Schema:  "parted",
				Indexes: []*metadatapb.IndexMetadata{{Name: "idx_p"}},
			}},
		}},
	}, {Name: "parted"}}}
	require.True(t, typeObjects(meta)[objectKey("parted", "t_p")], "a partition provides a row type")

	selected := map[string]bool{objectKey("parted", "t_p"): true}
	require.Equal(t, objectKey("parted", "t_p"), indexNameOwner(meta, "parted", "idx_p", selected),
		"the partition owns the index its comment names")
}

// TestSameNamedPartitionsKeepTheirSchema covers same-named partitions held under
// parents in other schemas: the global lookup must keep the owner clause's schema.
func TestSameNamedPartitionsKeepTheirSchema(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "a", Tables: []*metadatapb.TableMetadata{{
			Name: "parent_a",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:    "p",
				Schema:  "a",
				Indexes: []*metadatapb.IndexMetadata{{Name: "idx", Definition: "a definition"}},
			}},
		}}},
		{Name: "b", Tables: []*metadatapb.TableMetadata{{
			Name: "parent_b",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:    "p",
				Schema:  "b",
				Indexes: []*metadatapb.IndexMetadata{{Name: "idx", Definition: "b definition"}},
			}},
		}}},
	}}
	require.Equal(t, "b definition", indexByName(meta, "b", "b", "p", "idx").GetDefinition())
	require.Equal(t, "a definition", indexByName(meta, "a", "a", "p", "idx").GetDefinition())
}

// TestRoutineDropArgumentsKeepsEmptyList covers an explicitly empty argument list:
// PostgreSQL lets the list go only when the name is unique, so an overload keeps
// its parentheses.
func TestRoutineDropArgumentsKeepsEmptyList(t *testing.T) {
	require.Equal(t, "()", routineDropArguments("f()", ""))
	require.Equal(t, "(integer)", routineDropArguments("f(integer)", ""))
	require.Equal(t, "()", routineDropArguments("", "CREATE FUNCTION public.f() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$"))
	require.Equal(t, "", routineDropArguments("", ""))
}

// TestIndexOwnerSkipsQuotedNames covers an index name that contains the owner
// clause's text: the owner is read after the quoted identifier.
func TestIndexOwnerSkipsQuotedNames(t *testing.T) {
	tests := []struct{ sql, schema, object string }{
		{`CREATE INDEX "x ON y" ON "public"."t" USING btree ("id")`, "public", "t"},
		{`CREATE UNIQUE INDEX "a"" ON b" ON ONLY "public"."t" ("id")`, "public", "t"},
		{`CREATE INDEX "idx" ON "public"."mv" (id)`, "public", "mv"},
	}
	for _, tc := range tests {
		schemaName, objectName := indexOwnerFromSQL(tc.sql)
		require.Equal(t, tc.schema, schemaName, tc.sql)
		require.Equal(t, tc.object, objectName, tc.sql)
	}
}

// TestRelationSetIncludesPartitions covers a removed partition as a relation of its
// own: a composite attribute holding its row type has to be released before it.
func TestRelationSetIncludesPartitions(t *testing.T) {
	withPartition := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main", Tables: []*metadatapb.TableMetadata{{
			Name: "t",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:   "t_p",
				Schema: "parted",
			}},
		}}},
		{Name: "parted"},
	}}
	withoutPartition := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main", Tables: []*metadatapb.TableMetadata{{Name: "t"}}},
		{Name: "parted"},
	}}
	require.True(t, relationSet(withPartition, withoutPartition)[objectKey("parted", "t_p")], "the removed partition is a relation of its own")
	require.False(t, relationSet(withPartition, withPartition)[objectKey("parted", "t_p")], "an unchanged partition is not")
}

// TestRelationSetNestedPartitionKeepsImmediateParentSchema covers a subpartition
// without its own schema under a partition in another schema: the effective
// schema is the immediate parent's, not the top-level table's.
func TestRelationSetNestedPartitionKeepsImmediateParentSchema(t *testing.T) {
	withNested := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main", Tables: []*metadatapb.TableMetadata{{
			Name: "t",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:   "p",
				Schema: "parted",
				Subpartitions: []*metadatapb.TablePartitionMetadata{{
					Name: "sp",
				}},
			}},
		}}},
		{Name: "parted"},
	}}
	withoutNested := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main", Tables: []*metadatapb.TableMetadata{{Name: "t"}}},
		{Name: "parted"},
	}}

	set := relationSet(withNested, withoutNested)
	require.True(t, set[objectKey("parted", "sp")], "the subpartition inherits parted from its immediate parent")
	require.False(t, set[objectKey("main", "sp")], "the subpartition does not inherit the top-level table's schema")
}

// TestSourceOnlyUnloadableObjectsAreDropped covers removing an object whose stored
// definition omni cannot parse: a drop needs no definition, so the plan emits it,
// while adding such an object is still refused.
func TestSourceOnlyUnloadableObjectsAreDropped(t *testing.T) {
	unloadable := func() *metadatapb.ViewMetadata {
		return &metadatapb.ViewMetadata{Name: "v", Definition: "SELECT @@@unsupported@@@"}
	}
	withView := newPGTestMetadata(nil, []*metadatapb.ViewMetadata{unloadable()})
	withoutView := newPGTestMetadata(nil, nil)

	sql, err := testMetadataMigration(withView, withoutView)
	require.NoError(t, err)
	require.Equal(t, `DROP VIEW "public"."v"`, sql)

	_, err = testMetadataMigration(withoutView, withView)
	require.ErrorContains(t, err, "could not be loaded", "creating an unparseable object is still refused")
}

// TestSourceOnlyMissingObjectsAreDropped covers an object the catalog cannot
// install even as a stand-in: a plain index whose expression names no column.
// Its removal still needs a DROP, which comes from the metadata because the
// catalog never had the index to diff.
func TestSourceOnlyMissingObjectsAreDropped(t *testing.T) {
	snapshot := func(indexes []*metadatapb.IndexMetadata) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				Indexes: indexes,
			}},
		}}}
	}
	withIndex := snapshot([]*metadatapb.IndexMetadata{{Name: "t_missing_idx", Expressions: []string{"no_such_column"}, Type: "btree"}})
	withoutIndex := snapshot(nil)

	sql, err := testMetadataMigration(withIndex, withoutIndex)
	require.NoError(t, err)
	require.Equal(t, `DROP INDEX "public"."t_missing_idx"`, sql)

	_, err = testMetadataMigration(withoutIndex, withIndex)
	require.ErrorContains(t, err, "could not be loaded", "creating an object without a stand-in is still refused")
}

// TestSourceOnlyDependentDropsPrecedeTheirProviders pins the order a missing
// routine or constraint needs: its DROP comes before the catalog's drops of the
// type, table or referenced table it depends on, while an index DROP stays before
// the column it names.
func TestSourceOnlyDependentDropsPrecedeTheirProviders(t *testing.T) {
	ops := []MigrationOp{
		{Type: OpDropType, SQL: `DROP TYPE "public"."w"`, Phase: PhasePre},
		{Type: OpDropTable, SQL: `DROP TABLE "public"."t" CASCADE`, Phase: PhasePre},
		{Type: OpCreateTable, SQL: `CREATE TABLE "public"."next" ("id" integer)`, Phase: PhaseMain},
	}
	drops := map[string]MigrationOp{
		"con:public.child.fk_parent": {Type: OpDropConstraint, SQL: `ALTER TABLE "public"."child" DROP CONSTRAINT "fk_parent"`, Phase: PhasePre},
		"func:public.f|f(public.w)":  {Type: OpDropFunction, SQL: `DROP FUNCTION "public"."f"(public.w)`, Phase: PhasePre},
		"idx:public.t.t_idx":         {Type: OpDropIndex, SQL: `DROP INDEX "public"."t_idx"`, Phase: PhasePre},
	}
	var got []string
	for _, op := range insertSourceOnlyDrops(ops, drops) {
		got = append(got, op.SQL)
	}
	require.Equal(t, []string{
		`ALTER TABLE "public"."child" DROP CONSTRAINT "fk_parent"`,
		`DROP FUNCTION "public"."f"(public.w)`,
		`DROP TYPE "public"."w"`,
		`DROP TABLE "public"."t" CASCADE`,
		`DROP INDEX "public"."t_idx"`,
		`CREATE TABLE "public"."next" ("id" integer)`,
	}, got)
}

// TestCockroachMetadataMigrationNoOps pins the cases the old generator covered
// that produce no migration at all. They have no legacy differential, so they
// stay out of the cutover fixture.
func TestCockroachMetadataMigrationNoOps(t *testing.T) {
	unparseableView := &metadatapb.ViewMetadata{Name: "v", Definition: "SELECT @@@omni_unsupported@@@"}
	composites := []*metadatapb.CompositeTypeMetadata{{
		Name:    "addr",
		Comment: "c",
		Attributes: []*metadatapb.CompositeTypeAttribute{
			{Name: "street", Type: "text", Collation: `"C"`, Comment: "s"},
		},
	}}

	tests := []struct {
		name   string
		source *metadatapb.DatabaseSchemaMetadata
		target *metadatapb.DatabaseSchemaMetadata
	}{
		{
			name:   "identical raw view does not load",
			source: newPGTestMetadata(nil, []*metadatapb.ViewMetadata{unparseableView}),
			target: newPGTestMetadata(nil, []*metadatapb.ViewMetadata{unparseableView}),
		},
		{
			name:   "identical composite types",
			source: newPGTestMetadataWithCompositeTypes(composites),
			target: newPGTestMetadataWithCompositeTypes(composites),
		},
		{
			name: "skip-dump composite types",
			source: newPGTestMetadataWithCompositeTypes([]*metadatapb.CompositeTypeMetadata{{
				Name:     "ext_owned",
				SkipDump: true,
				Attributes: []*metadatapb.CompositeTypeAttribute{
					{Name: "x", Type: "integer"},
				},
			}}),
			target: newPGTestMetadataWithCompositeTypes(nil),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sql, err := testMetadataMigration(tc.source, tc.target)
			require.NoError(t, err)
			require.Empty(t, sql)
		})
	}
}

func newPGTestMetadataWithCompositeTypes(compositeTypes []*metadatapb.CompositeTypeMetadata) *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:           "public",
				CompositeTypes: compositeTypes,
			},
		},
	}
}

// TestCockroachMetadataMigrationIgnoresLegacyFunctionComparer records that the
// schema.FunctionComparer registry no longer governs CockroachDB: the catalog
// compares the stored definitions itself, so a comparer that declares two
// different bodies equal cannot suppress the replace.
func TestOmniDiffMetadataMigrationDiffsAroundAnUnchangedUnloadableView(t *testing.T) {
	id := &metadatapb.ColumnMetadata{Name: "id", Type: "integer"}
	note := &metadatapb.ColumnMetadata{Name: "note", Type: "text", Nullable: true}
	view := &metadatapb.ViewMetadata{Name: "v", Definition: "SELECT @@@omni_unsupported@@@"}

	sql, err := testMetadataMigration(newPGTestMetadata([]*metadatapb.TableMetadata{pgTable("t", 0, id)}, []*metadatapb.ViewMetadata{view}),
		newPGTestMetadata([]*metadatapb.TableMetadata{pgTable("t", 0, id, note)}, []*metadatapb.ViewMetadata{view}),
	)
	require.NoError(t, err)
	require.Contains(t, sql, `ADD COLUMN note text`)
	require.NotContains(t, sql, "VIEW")
}

func TestOmniDiffMetadataMigrationRefusesAChangedUnloadableObject(t *testing.T) {
	_, err := testMetadataMigration(newPGTestMetadata(nil, []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT @@@old_definition@@@"}}),
		newPGTestMetadata(nil, []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT @@@new_definition@@@"}}),
	)
	require.ErrorContains(t, err, "rel:public.v")
}

func TestOmniDiffMetadataMigrationIgnoresStatisticsOfAnUnloadableTable(t *testing.T) {
	id := &metadatapb.ColumnMetadata{Name: "id", Type: "integer"}
	unparsed := func(rows int64) *metadatapb.TableMetadata {
		return pgTable("unparsed", rows, &metadatapb.ColumnMetadata{Name: "a", Type: "no_such_type(((", Nullable: true})
	}
	note := &metadatapb.ColumnMetadata{Name: "note", Type: "text", Nullable: true}

	sql, err := testMetadataMigration(newPGTestMetadata([]*metadatapb.TableMetadata{pgTable("t", 0, id), unparsed(10)}, nil),
		newPGTestMetadata([]*metadatapb.TableMetadata{pgTable("t", 0, id, note), unparsed(20)}, nil),
	)
	require.NoError(t, err)
	require.Contains(t, sql, `ADD COLUMN note text`)
}

// TestIdentityGenerationBecomesIdentityColumn covers CockroachDB sync, which
// records identityGeneration without isIdentity. The catalog only installs an
// identity column from the flag, so the shim sets it before loading.
func TestIdentityGenerationBecomesIdentityColumn(t *testing.T) {
	id := &metadatapb.ColumnMetadata{
		Name:               "id",
		Type:               "bigint",
		IdentityGeneration: metadatapb.ColumnMetadata_ALWAYS,
	}
	source := newPGTestMetadata(nil, nil)
	target := newPGTestMetadata([]*metadatapb.TableMetadata{pgTable("t", 0, id)}, nil)

	sql, err := testMetadataMigration(source, target)
	require.NoError(t, err)
	require.Contains(t, sql, "GENERATED ALWAYS AS IDENTITY")
}

// TestRecreatedViewReappliesUnchangedComment covers a dependency-driven view
// recreation: the comment did not change, so the catalog's comment diff is
// empty, but the recreated view would otherwise lose it.
func TestRecreatedViewReappliesUnchangedComment(t *testing.T) {
	view := func(definition string) *metadatapb.ViewMetadata {
		return &metadatapb.ViewMetadata{
			Name:       "v",
			Definition: definition,
			Comment:    "ids",
			DependencyColumns: []*metadatapb.DependencyColumn{{
				Schema: "public", Table: "t", Column: "id",
			}},
		}
	}
	source := newPGTestMetadata(
		[]*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
		[]*metadatapb.ViewMetadata{view("SELECT id FROM public.t")},
	)
	target := newPGTestMetadata(
		[]*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "bigint"})},
		[]*metadatapb.ViewMetadata{view("SELECT id::bigint AS id FROM public.t")},
	)

	sql, err := testMetadataMigration(source, target)
	require.NoError(t, err)
	require.Contains(t, sql, "DROP VIEW")
	require.Contains(t, sql, `COMMENT ON VIEW "public"."v" IS 'ids'`)
	require.Equal(t, 1, strings.Count(sql, "COMMENT ON VIEW"))
}

// TestUnloadableUnchangedViewRecreatedWhenDependencyChanges covers a view whose
// definition omni cannot parse. Both sides hold the same stand-in, so the
// catalog has no dependency edge; DependencyColumns must still force the drop
// and recreation around the column change.
func TestUnloadableUnchangedViewRecreatedWhenDependencyChanges(t *testing.T) {
	view := func() *metadatapb.ViewMetadata {
		return &metadatapb.ViewMetadata{
			Name:       "v",
			Definition: "SELECT FROM",
			DependencyColumns: []*metadatapb.DependencyColumn{{
				Schema: "public", Table: "t", Column: "id",
			}},
		}
	}
	source := newPGTestMetadata(
		[]*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
		[]*metadatapb.ViewMetadata{view()},
	)
	target := newPGTestMetadata(
		[]*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "bigint"})},
		[]*metadatapb.ViewMetadata{view()},
	)

	sql, err := testMetadataMigration(source, target)
	require.NoError(t, err)
	drop := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	alter := strings.Index(sql, "ALTER TABLE public.t ALTER COLUMN id TYPE bigint")
	create := strings.Index(sql, `CREATE VIEW "public"."v" AS`)
	require.GreaterOrEqual(t, drop, 0, sql)
	require.GreaterOrEqual(t, alter, 0, sql)
	require.GreaterOrEqual(t, create, 0, sql)
	require.Less(t, drop, alter, "the view is dropped before the column change")
	require.Less(t, alter, create, "the view is recreated after the column change")
}

// TestUnloadableUnchangedMaterializedViewRecreatedWithIndexes covers the same
// fallback for a materialized view: the recreation keeps WITH NO DATA and its
// indexes, which a bare CREATE ... AS would drop.
func TestUnloadableUnchangedMaterializedViewRecreatedWithIndexes(t *testing.T) {
	view := func() *metadatapb.MaterializedViewMetadata {
		return &metadatapb.MaterializedViewMetadata{
			Name:       "mv",
			Definition: "SELECT FROM",
			Indexes: []*metadatapb.IndexMetadata{{
				Name: "idx_mv", Expressions: []string{"id"}, Type: "btree",
			}},
			DependencyColumns: []*metadatapb.DependencyColumn{{
				Schema: "public", Table: "t", Column: "id",
			}},
		}
	}
	source := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:              "public",
		Tables:            []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{view()},
	}}}
	target := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:              "public",
		Tables:            []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "bigint"})},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{view()},
	}}}

	sql, err := testMetadataMigration(source, target)
	require.NoError(t, err)
	require.Contains(t, sql, `DROP MATERIALIZED VIEW IF EXISTS "public"."mv"`)
	require.Contains(t, sql, `CREATE MATERIALIZED VIEW "public"."mv"`)
	require.Contains(t, sql, "WITH NO DATA")
	require.Contains(t, sql, `CREATE INDEX "idx_mv" ON "public"."mv"`)
}

// TestUnloadableDependentViewsRestoreInDependencyOrder pins the order among the
// restored views: drops run dependents first, creates dependencies first.
func TestUnloadableDependentViewsRestoreInDependencyOrder(t *testing.T) {
	view := func(name string, deps ...*metadatapb.DependencyColumn) *metadatapb.ViewMetadata {
		return &metadatapb.ViewMetadata{Name: name, Definition: "SELECT FROM", DependencyColumns: deps}
	}
	tableDep := func() *metadatapb.DependencyColumn {
		return &metadatapb.DependencyColumn{Schema: "public", Table: "t", Column: "id"}
	}
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{
				view("a", tableDep(), &metadatapb.DependencyColumn{Schema: "public", Table: "b"}),
				view("b", tableDep()),
			},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.NoError(t, err)
	dropA := strings.Index(sql, `DROP VIEW IF EXISTS "public"."a"`)
	dropB := strings.Index(sql, `DROP VIEW IF EXISTS "public"."b"`)
	createA := strings.Index(sql, `CREATE VIEW "public"."a"`)
	createB := strings.Index(sql, `CREATE VIEW "public"."b"`)
	require.GreaterOrEqual(t, dropA, 0, sql)
	require.GreaterOrEqual(t, dropB, 0, sql)
	require.GreaterOrEqual(t, createA, 0, sql)
	require.GreaterOrEqual(t, createB, 0, sql)
	require.Less(t, dropA, dropB, "the dependent view is dropped first")
	require.Less(t, createB, createA, "the dependency is created first")
}

// TestBackupSchemaDependentViewIsIgnored covers backup-only drift: the plan
// filters bbdataarchive, so a backup view that depends on a changed backup table
// is not restored and emits no SQL.
func TestBackupSchemaDependentViewIsIgnored(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "bbdataarchive",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "bbdataarchive", Table: "t", Column: "id",
				}},
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.NoError(t, err)
	require.Empty(t, sql)
}

// TestSplitCompositeRetypeAddsAttributeBeforeItsComment covers a split release
// whose attribute also gains a comment: the ADD must land after the new relation
// exists but before the COMMENT that targets the attribute.
func TestSplitCompositeRetypeAddsAttributeBeforeItsComment(t *testing.T) {
	snapshot := func(table, attributeType, comment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable(table, 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name: "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{{
					Name: "r", Type: attributeType, Comment: comment,
				}},
			}},
		}}}
	}
	sql, err := testMetadataMigration(snapshot("old_t", "public.old_t", ""), snapshot("new_t", "public.new_t", "row comment"))
	require.NoError(t, err)
	add := strings.Index(sql, `ADD ATTRIBUTE "r" public.new_t`)
	comment := strings.Index(sql, "COMMENT ON COLUMN")
	require.GreaterOrEqual(t, add, 0, sql)
	require.GreaterOrEqual(t, comment, 0, sql)
	require.Less(t, add, comment, sql)
}

// TestSplitCompositeRetypeReappliesUnchangedComment covers the split path when
// the attribute keeps its comment: DROP ATTRIBUTE deletes it, the catalog sees
// no comment change, so the add must be followed by the comment.
func TestSplitCompositeRetypeReappliesUnchangedComment(t *testing.T) {
	snapshot := func(table, attributeType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable(table, 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name: "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{{
					Name: "r", Type: attributeType, Comment: "row comment",
				}},
			}},
		}}}
	}
	sql, err := testMetadataMigration(snapshot("old_t", "public.old_t"), snapshot("new_t", "public.new_t"))
	require.NoError(t, err)
	add := strings.Index(sql, `ADD ATTRIBUTE "r" public.new_t`)
	comment := strings.Index(sql, `COMMENT ON COLUMN "public"."c"."r" IS 'row comment'`)
	require.GreaterOrEqual(t, add, 0, sql)
	require.GreaterOrEqual(t, comment, 0, sql)
	require.Less(t, add, comment, sql)
	require.Equal(t, 1, strings.Count(sql, "COMMENT ON COLUMN"), sql)
}

// TestSplitCompositeRetypeKeepsTargetAttributeOrder pins the order of attributes
// split from one dropped row type onto the same created type: composite fields are
// positional, so both ADDs must follow the target's attribute order.
func TestSplitCompositeRetypeKeepsTargetAttributeOrder(t *testing.T) {
	snapshot := func(table string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable(table, 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name: "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{
					{Name: "a", Type: "public." + table},
					{Name: "b", Type: "public." + table},
				},
			}},
		}}}
	}
	sql, err := testMetadataMigration(snapshot("old_t"), snapshot("new_t"))
	require.NoError(t, err)
	dropTable := strings.Index(sql, `DROP TABLE "public"."old_t"`)
	dropA := strings.Index(sql, `DROP ATTRIBUTE "a"`)
	dropB := strings.Index(sql, `DROP ATTRIBUTE "b"`)
	addA := strings.Index(sql, `ADD ATTRIBUTE "a" public.new_t`)
	addB := strings.Index(sql, `ADD ATTRIBUTE "b" public.new_t`)
	for name, index := range map[string]int{"drop table": dropTable, "drop a": dropA, "drop b": dropB, "add a": addA, "add b": addB} {
		require.GreaterOrEqual(t, index, 0, "%s missing from:\n%s", name, sql)
	}
	require.Less(t, dropA, dropTable, "attribute a is released before the old row type is dropped")
	require.Less(t, dropB, dropTable, "attribute b is released before the old row type is dropped")
	require.Less(t, addA, addB, "split attributes keep the target's attribute order")
}

// TestSplitCompositeRetypeReaddsUnchangedSuffix covers a split attribute that
// precedes an unchanged one. ADD ATTRIBUTE appends, so the unchanged attribute
// must be released and re-added too, or the composite ends up in the wrong order.
func TestSplitCompositeRetypeReaddsUnchangedSuffix(t *testing.T) {
	snapshot := func(table string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable(table, 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name: "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{
					{Name: "a", Type: "public." + table},
					{Name: "z", Type: "text", Comment: "keep"},
				},
			}},
		}}}
	}
	sql, err := testMetadataMigration(snapshot("old_t"), snapshot("new_t"))
	require.NoError(t, err)
	dropTable := strings.Index(sql, `DROP TABLE "public"."old_t"`)
	dropA := strings.Index(sql, `DROP ATTRIBUTE "a"`)
	dropZ := strings.Index(sql, `DROP ATTRIBUTE "z"`)
	addA := strings.Index(sql, `ADD ATTRIBUTE "a" public.new_t`)
	addZ := strings.Index(sql, `ADD ATTRIBUTE "z" text`)
	for name, index := range map[string]int{"drop table": dropTable, "drop a": dropA, "drop z": dropZ, "add a": addA, "add z": addZ} {
		require.GreaterOrEqual(t, index, 0, "%s missing from:\n%s", name, sql)
	}
	require.Less(t, dropA, dropTable, sql)
	require.Less(t, dropZ, dropTable, "the unchanged attribute is released before the row type drop")
	require.Less(t, addA, addZ, "the re-added suffix keeps the target's attribute order")
	require.Contains(t, sql, `COMMENT ON COLUMN "public"."c"."z" IS 'keep'`, "the re-added attribute keeps its comment")
}

// TestSplitCompositeReleasesTheWholeSuffixBeforeEveryDrop covers a suffix whose
// later attribute is released from a relation that the triggering attribute never
// named: every suffix release must precede every relation it is released from.
func TestSplitCompositeReleasesTheWholeSuffixBeforeEveryDrop(t *testing.T) {
	source := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{
			pgTable("z_old", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"}),
			pgTable("a_old", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"}),
		},
		CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
			Name: "c",
			Attributes: []*metadatapb.CompositeTypeAttribute{
				{Name: "a", Type: "public.z_old"},
				{Name: "b", Type: "public.a_old"},
			},
		}},
	}}}
	target := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:   "public",
		Tables: []*metadatapb.TableMetadata{pgTable("new_a", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
		CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
			Name: "c",
			Attributes: []*metadatapb.CompositeTypeAttribute{
				{Name: "a", Type: "public.new_a"},
				{Name: "b", Type: "text"},
			},
		}},
	}}}

	sql, err := testMetadataMigration(source, target)
	require.NoError(t, err)
	dropA := strings.Index(sql, `DROP ATTRIBUTE "a"`)
	dropB := strings.Index(sql, `DROP ATTRIBUTE "b"`)
	dropAOld := strings.Index(sql, `DROP TABLE "public"."a_old"`)
	dropZOld := strings.Index(sql, `DROP TABLE "public"."z_old"`)
	for name, index := range map[string]int{"drop a": dropA, "drop b": dropB, "drop a_old": dropAOld, "drop z_old": dropZOld} {
		require.GreaterOrEqual(t, index, 0, "%s missing from:\n%s", name, sql)
	}
	require.Less(t, dropA, dropZOld, "the triggering release precedes its relation drop")
	require.Less(t, dropB, dropAOld, "a later suffix release precedes the relation it holds")
	require.Less(t, dropB, dropZOld, "a later suffix release precedes every relation drop in the plan")
}

// TestRestoredUnloadableViewReleasesRowTypeDependents covers an unchanged
// unloadable view that must be recreated for a column change while an unchanged
// composite attribute holds its row type: the attribute is dropped before the
// view and re-added after it, or the view drop finds a dependent object.
func TestRestoredUnloadableViewReleasesRowTypeDependents(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name:       "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "r", Type: "public.v"}},
			}},
		}}}
	}
	sql, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.NoError(t, err)
	dropAttribute := strings.Index(sql, `ALTER TYPE "public"."c" DROP ATTRIBUTE "r"`)
	dropView := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	createView := strings.Index(sql, `CREATE VIEW "public"."v"`)
	addAttribute := strings.Index(sql, `ALTER TYPE "public"."c" ADD ATTRIBUTE "r" public.v`)
	for name, index := range map[string]int{"drop attribute": dropAttribute, "drop view": dropView, "create view": createView, "add attribute": addAttribute} {
		require.GreaterOrEqual(t, index, 0, "%s missing from:\n%s", name, sql)
	}
	require.Less(t, dropAttribute, dropView, "the row-type dependent is released before the view is dropped")
	require.Less(t, createView, addAttribute, "the attribute is re-added after the view is recreated")
}

// TestRestoredUnloadableViewReleasesRetypedRowTypeAttribute covers an existing
// attribute retyped to a restored view's row type: the catalog cannot order its
// ALTER TYPE against the drop, so it is released and re-added like a split.
func TestRestoredUnloadableViewReleasesRetypedRowTypeAttribute(t *testing.T) {
	snapshot := func(columnType, attributeType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name: "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{
					{Name: "r", Type: attributeType},
					{Name: "z", Type: "text", Comment: "keep"},
				},
			}},
		}}}
	}
	sql, err := testMetadataMigration(snapshot("integer", "text"), snapshot("bigint", "public.v"))
	require.NoError(t, err)
	dropAttribute := strings.Index(sql, `ALTER TYPE "public"."c" DROP ATTRIBUTE "r"`)
	dropView := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	createView := strings.Index(sql, `CREATE VIEW "public"."v"`)
	addAttribute := strings.Index(sql, `ALTER TYPE "public"."c" ADD ATTRIBUTE "r" public.v`)
	require.GreaterOrEqual(t, dropAttribute, 0, sql)
	require.GreaterOrEqual(t, dropView, 0, sql)
	require.GreaterOrEqual(t, createView, 0, sql)
	require.GreaterOrEqual(t, addAttribute, 0, sql)
	require.Less(t, dropAttribute, dropView, "the retyped attribute is released before the view is dropped")
	require.Less(t, createView, addAttribute, "the attribute is re-added after the view is recreated")
	require.NotContains(t, sql, `ADD ATTRIBUTE "r" text`, "the catalog's own alter is replaced")
	require.Contains(t, sql, `COMMENT ON COLUMN "public"."c"."z" IS 'keep'`, "the re-added unchanged attribute keeps its comment")
}

// TestCollationChangeRestoresUnloadableView pins the ordering a collation-only
// column change needs: the view is dropped before the column is altered and
// recreated after it, even though the column's type does not change.
func TestCollationChangeRestoresUnloadableView(t *testing.T) {
	snapshot := func(collation string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "text", Collation: collation})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("C"), snapshot("en_US"))
	require.NoError(t, err)
	dropView := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	alterColumn := strings.Index(sql, `ALTER TABLE public.t ALTER COLUMN id TYPE text COLLATE "en_US"`)
	createView := strings.Index(sql, `CREATE VIEW "public"."v"`)
	for name, index := range map[string]int{"drop view": dropView, "alter column": alterColumn, "create view": createView} {
		require.GreaterOrEqual(t, index, 0, "%s missing from:\n%s", name, sql)
	}
	require.Less(t, dropView, alterColumn, "the view is dropped before the collation change")
	require.Less(t, alterColumn, createView, "the view is recreated after the collation change")
}

// TestUnqualifiedTypeFollowsTheSearchPath covers a name two schemas share: only
// the first visible relation in the search path is the one a type refers to.
func TestUnqualifiedTypeFollowsTheSearchPath(t *testing.T) {
	name := "shared"
	relations := map[string]bool{
		objectKey("first", name):  true,
		objectKey("second", name): true,
		objectKey("public", name): true,
	}
	require.Equal(t, objectKey("second", name), func() string {
		key, ok := referencedRelation(name, relations, []string{"second", "first"})
		require.True(t, ok)
		return key
	}())
	key, ok := referencedRelation(name, relations, []string{"first"})
	require.True(t, ok)
	require.Equal(t, objectKey("first", name), key)
	key, ok = referencedRelation(name, relations, nil)
	require.True(t, ok, "an empty path resolves through the default public schema")
	require.Equal(t, objectKey("public", name), key)
	_, ok = referencedRelation(name, map[string]bool{objectKey("first", name): true}, nil)
	require.False(t, ok, "without public the default path does not choose a schema by key order")
	key, ok = referencedRelation("first."+name, relations, []string{"second"})
	require.True(t, ok, "a qualified name ignores the search path")
	require.Equal(t, objectKey("first", name), key)
}

// TestPgCatalogTypesNameNoRelation covers an unqualified built-in type: pg_catalog
// is searched first, so `text` is not a reference to a relation named text.
func TestPgCatalogTypesNameNoRelation(t *testing.T) {
	relations := map[string]bool{objectKey("public", "text"): true}
	_, ok := referencedRelation("text", relations, []string{"public"})
	require.False(t, ok, "the built-in is not public.text")
	_, ok = referencedRelation(`"text"`, relations, []string{"public"})
	require.False(t, ok, "quoting an exact built-in name does not make it the user relation")
	_, ok = referencedRelation(`"Text"`, map[string]bool{objectKey("public", "Text"): true}, []string{"public"})
	require.True(t, ok, `"Text" is not the built-in lowercase name`)
	require.True(t, builtinMultiwordType("double precision[]"), "an array marker does not hide the built-in")
	require.False(t, builtinMultiwordType("public.precision"))
}

// TestAmbiguityCheckSkipsPgCatalogTypes covers a pg_catalog type that a user schema
// also names: pg_catalog is searched first, so the reference is the built-in.
func TestAmbiguityCheckSkipsPgCatalogTypes(t *testing.T) {
	snapshot := func(withHolder bool) *metadatapb.DatabaseSchemaMetadata {
		schemas := []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
					Name:       "text",
					Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "value", Type: "integer"}},
				}},
			},
		}
		if withHolder {
			schemas = append(schemas, &metadatapb.SchemaMetadata{
				Name:   "app",
				Tables: []*metadatapb.TableMetadata{pgTable("holder", 0, &metadatapb.ColumnMetadata{Name: "value", Type: "text"})},
			})
		}
		return &metadatapb.DatabaseSchemaMetadata{
			SearchPath: `"$user", public`,
			Schemas:    schemas,
		}
	}
	_, err := testMetadataMigration(snapshot(false), snapshot(true))
	require.NoError(t, err, "text is a pg_catalog type, not a reference to public.text")
}

// TestAmbiguityCheckSkipsMultiwordBuiltins covers a built-in spelled with several
// words: its last word cannot be the type a $user shadow would decide.
func TestAmbiguityCheckSkipsMultiwordBuiltins(t *testing.T) {
	snapshot := func(withHolder bool) *metadatapb.DatabaseSchemaMetadata {
		schemas := []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
					Name:       "precision",
					Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "value", Type: "text"}},
				}},
			},
		}
		if withHolder {
			schemas = append(schemas, &metadatapb.SchemaMetadata{
				Name:   "app",
				Tables: []*metadatapb.TableMetadata{pgTable("holder", 0, &metadatapb.ColumnMetadata{Name: "value", Type: "double precision"})},
			})
		}
		return &metadatapb.DatabaseSchemaMetadata{
			SearchPath: `"$user", public`,
			Schemas:    schemas,
		}
	}
	_, err := testMetadataMigration(snapshot(false), snapshot(true))
	require.NoError(t, err, "double precision is a built-in, not a reference to public.precision")
}

// TestUnresolvedSearchPathUserRejectsAmbiguousType covers an unqualified type that
// needs the unresolved $user entry: public has the type, but the login's schema
// might shadow it, so assigning the public dependency would be a guess.
func TestUnresolvedSearchPathUserRejectsAmbiguousType(t *testing.T) {
	snapshot := func(withHolder bool) *metadatapb.DatabaseSchemaMetadata {
		schemas := []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
					Name:       "shared",
					Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "value", Type: "text"}},
				}},
			},
		}
		if withHolder {
			schemas = append(schemas, &metadatapb.SchemaMetadata{
				Name:   "app",
				Tables: []*metadatapb.TableMetadata{pgTable("holder", 0, &metadatapb.ColumnMetadata{Name: "row", Type: "shared"})},
			})
		}
		return &metadatapb.DatabaseSchemaMetadata{
			SearchPath: `"$user", public`,
			Schemas:    schemas,
		}
	}

	_, err := testMetadataMigration(snapshot(false), snapshot(true))
	require.ErrorContains(t, err, "whose $user search-path resolution is ambiguous")

	owner := func(withHolder bool) *metadatapb.DatabaseSchemaMetadata {
		metadata := snapshot(withHolder)
		metadata.Owner = "alice"
		return metadata
	}
	_, err = testMetadataMigration(owner(false), owner(true))
	require.ErrorContains(t, err, "whose $user search-path resolution is ambiguous",
		"the session login decides $user, not the database owner, so this layer refuses")
}

// TestUnresolvedSearchPathUserRejectsAmbiguousViewReference covers a view body
// whose unqualified relation is actually in the login's schema while public has a
// same-named relation: loading with the literal $user would bind the wrong one.
func TestUnresolvedSearchPathUserRejectsAmbiguousViewReference(t *testing.T) {
	snapshot := func(definition string) *metadatapb.DatabaseSchemaMetadata {
		schemas := []*metadatapb.SchemaMetadata{
			{Name: "alice", Tables: []*metadatapb.TableMetadata{pgTable("v", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})}},
			{Name: "public", Tables: []*metadatapb.TableMetadata{pgTable("v", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})}},
		}
		if definition != "" {
			schemas = append(schemas, &metadatapb.SchemaMetadata{
				Name: "app",
				Views: []*metadatapb.ViewMetadata{{
					Name:       "w",
					Definition: definition,
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "alice", Table: "v",
					}},
				}},
			})
		}
		return &metadatapb.DatabaseSchemaMetadata{
			SearchPath: `"$user", public`,
			Schemas:    schemas,
		}
	}

	_, err := testMetadataMigration(snapshot(""), snapshot("SELECT * FROM v"))
	require.ErrorContains(t, err, `view "app"."w" has an unqualified reference to "alice"."v"`)

	_, err = testMetadataMigration(snapshot(""), snapshot("SELECT * FROM alice.v"))
	require.NoError(t, err, "a qualified reference does not depend on the user schema")
}

// TestExplicitPgCatalogPositionDecidesBuiltin covers PostgreSQL's rule that
// pg_catalog is implicit first only when the path does not list it; an explicit
// user schema before pg_catalog can shadow a built-in name.
func TestExplicitPgCatalogPositionDecidesBuiltin(t *testing.T) {
	relations := map[string]bool{objectKey("public", "text"): true}
	key, ok := referencedRelation("text", relations, []string{"public", "pg_catalog"})
	require.True(t, ok, "an explicit pg_catalog after public lets public.text shadow text")
	require.Equal(t, objectKey("public", "text"), key)
	_, ok = referencedRelation("text", relations, []string{"pg_catalog", "public"})
	require.False(t, ok, "pg_catalog before public keeps the built-in")
	_, ok = referencedRelation("text", relations, []string{"public"})
	require.False(t, ok, "an implicit pg_catalog stays first")

	require.Equal(t, []string{"public", "pg_catalog"}, parseTestSearchPath("public, pg_catalog"), "an explicit pg_catalog keeps its position")
}

// TestTypeNameMatchingKeepsIdentifierBoundaries pins quoted identifier parsing:
// a dot inside quotes belongs to that identifier, so `"a.b"."c"` and `"a"."b.c"`
// name different relations.
func TestTypeNameMatchingKeepsIdentifierBoundaries(t *testing.T) {
	dottedSchema := objectKey("a.b", "c")
	dottedObject := objectKey("a", "b.c")
	require.True(t, referencesType(`"a.b"."c"`, dottedSchema))
	require.True(t, referencesType(`"a"."b.c"`, dottedObject))
	require.False(t, referencesType(`"a"."b.c"`, dottedSchema), `"a"."b.c" is not "a.b"."c"`)
	require.False(t, referencesType(`"a.b"."c"`, dottedObject))
	require.True(t, referencesType("c", dottedSchema), "an unqualified name matches the object part")
	require.False(t, referencesType("c", dottedObject))

	for _, tc := range []struct{ typeName, key string }{
		{`"a"."b.c"`, dottedSchema},
		{`"a.b"."c"`, dottedObject},
	} {
		_, ok := referencedRelation(tc.typeName, map[string]bool{tc.key: true}, nil)
		require.False(t, ok, "%s does not name %s", tc.typeName, tc.key)
	}
	for _, tc := range []struct{ typeName, key string }{
		{`"a.b"."c"`, dottedSchema},
		{`"a"."b.c"`, dottedObject},
		{`"a.b"."c"[]`, dottedSchema},
	} {
		key, ok := referencedRelation(tc.typeName, map[string]bool{tc.key: true}, nil)
		require.True(t, ok, "%s names %s", tc.typeName, tc.key)
		require.Equal(t, tc.key, key)
	}

	upper := objectKey("public", "Foo")
	lower := objectKey("public", "foo")
	require.True(t, referencesType(`"public"."Foo"`, upper))
	require.False(t, referencesType(`"public"."Foo"`, lower), "a quoted identifier matches its case exactly")
	require.True(t, referencesType("public.foo", lower), "an unquoted identifier folds case")
	require.True(t, referencesType("public.Foo", lower))
	_, ok := referencedRelation(`"public"."Foo"`, map[string]bool{lower: true}, nil)
	require.False(t, ok, `"public"."Foo" is not "public"."foo"`)
	mixedUpper := objectKey("Foo", "bar")
	mixedLower := objectKey("foo", "bar")
	require.True(t, referencesType(`"Foo".bar`, mixedUpper))
	require.False(t, referencesType(`"Foo".bar`, mixedLower), "a quoted schema keeps its case while the object folds")
	require.False(t, referencesType("public.Foo", upper), "an unquoted reference folds to lower case before resolving")
	require.True(t, referencesType("public.Foo", lower))
	_, ok = referencedRelation("public.FOO", map[string]bool{upper: true}, nil)
	require.False(t, ok, "an unquoted reference does not name a quoted upper-case relation")
	require.True(t, referencesType(`foo."Bar"`, objectKey("foo", "Bar")))
	require.False(t, referencesType(`foo."Bar"`, objectKey("foo", "bar")), "a quoted object keeps its case while the schema folds")
	key, ok := referencedRelation("public.Foo", map[string]bool{lower: true}, nil)
	require.True(t, ok)
	require.Equal(t, lower, key)
}

// TestReleaseMovesTheAttributeNamedLikeItsType covers an attribute whose name
// repeats its composite type's name: the release must find the operation for that
// attribute, not the first one that mentions the type.
func TestReleaseMovesTheAttributeNamedLikeItsType(t *testing.T) {
	snapshot := func(table string) *metadatapb.DatabaseSchemaMetadata {
		attributeType := "text"
		schemaMeta := &metadatapb.SchemaMetadata{
			Name: "public",
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name: "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{
					{Name: "x", Type: attributeType},
					{Name: "c", Type: attributeType},
				},
			}},
		}
		if table != "" {
			schemaMeta.Tables = []*metadatapb.TableMetadata{pgTable(table, 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})}
			schemaMeta.CompositeTypes[0].Attributes[0].Type = "public." + table
			schemaMeta.CompositeTypes[0].Attributes[1].Type = "public." + table
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{schemaMeta}}
	}
	sql, err := testMetadataMigration(snapshot("old_t"), snapshot(""))
	require.NoError(t, err)
	dropTable := strings.Index(sql, `DROP TABLE "public"."old_t"`)
	attributeC := strings.Index(sql, `ALTER ATTRIBUTE "c" TYPE text`)
	require.GreaterOrEqual(t, dropTable, 0, sql)
	require.GreaterOrEqual(t, attributeC, 0, sql)
	require.Less(t, attributeC, dropTable, "the attribute named like its type is released before the drop")
}

// TestOpTargetsSelectedKeepsDistinctCommentOwners keeps a comment on `"a"."b.c"`
// with its owner when `"a.b"."c"` is the selected view: their flattened
// descriptions are equal, but the quoted targets are not.
func TestOpTargetsSelectedKeepsDistinctCommentOwners(t *testing.T) {
	selected := map[string]bool{objectKey("a.b", "c"): true}
	other := MigrationOp{Type: OpComment, ObjectName: "a.b.c", SQL: `COMMENT ON VIEW "a"."b.c" IS 'new'`}
	require.False(t, opTargetsSelected(other, selected, nil, nil, nil), "the other view's comment is not the selected view's")
	own := MigrationOp{Type: OpComment, ObjectName: "a.b.c", SQL: `COMMENT ON VIEW "a.b"."c" IS 'new'`}
	require.True(t, opTargetsSelected(own, selected, nil, nil, nil))
	column := MigrationOp{Type: OpComment, ObjectName: "a.b.c", SQL: `COMMENT ON COLUMN "a.b"."c"."id" IS 'id'`}
	require.True(t, opTargetsSelected(column, selected, nil, nil, nil), "an accessory of the selected view belongs to it")
}

// TestMultiwordTypeIsNotARelationReference covers a created table named like the
// first word of a built-in type: `timestamp without time zone` must not depend on
// a table called `timestamp`, or the false edge can hide the real order.
func TestMultiwordTypeIsNotARelationReference(t *testing.T) {
	require.False(t, referencesType("timestamp without time zone", objectKey("public", "timestamp")))
	_, referenced := referencedRelation("timestamp without time zone", map[string]bool{objectKey("public", "timestamp"): true}, nil)
	require.False(t, referenced, "a multiword type does not name a relation")
	require.False(t, referencesType("character varying(10)", objectKey("public", "character")))
	require.True(t, referencesType("timestamp", objectKey("public", "timestamp")))
	require.True(t, referencesType("public.t[]", objectKey("public", "t")))
	require.True(t, referencesType(`numeric(10,2)`, objectKey("public", "numeric")), "a modifier belongs to the type")

	source := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	target := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{
			pgTable("timestamp", 0, &metadatapb.ColumnMetadata{Name: "row", Type: "public.z_base"}),
			pgTable("holder", 0, &metadatapb.ColumnMetadata{Name: "at", Type: "timestamp without time zone"}),
			pgTable("z_base", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"}),
		},
	}}}
	sql, err := testMetadataMigration(source, target)
	require.NoError(t, err)
	provider := strings.Index(sql, `CREATE TABLE "public"."z_base"`)
	user := strings.Index(sql, `CREATE TABLE "public"."timestamp"`)
	require.GreaterOrEqual(t, provider, 0, sql)
	require.GreaterOrEqual(t, user, 0, sql)
	require.Less(t, provider, user, sql)
}

// TestCreatedTableRowTypeOrdersItsProviderFirst covers two tables created in one
// migration where one column uses the other's row type: the provider's create
// must come first, or the type does not exist yet.
func TestCreatedTableRowTypeOrdersItsProviderFirst(t *testing.T) {
	source := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	target := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{
			pgTable("a_holder", 0, &metadatapb.ColumnMetadata{Name: "row", Type: "public.z_base"}),
			pgTable("z_base", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"}),
		},
	}}}

	sql, err := testMetadataMigration(source, target)
	require.NoError(t, err)
	provider := strings.Index(sql, `CREATE TABLE "public"."z_base"`)
	user := strings.Index(sql, `CREATE TABLE "public"."a_holder"`)
	require.GreaterOrEqual(t, provider, 0, sql)
	require.GreaterOrEqual(t, user, 0, sql)
	require.Less(t, provider, user, "the table providing the row type is created first")
}

// TestAccessoryCommentsGoWithTheRestoredView covers an index comment on a
// restored materialized view and a trigger comment on a restored view: both are
// carried by the synthesized create, so the plan must not also emit them earlier,
// while their index and trigger do not exist.
// TestTargetSearchPathOrdersCreatedTypes covers a type name two schemas share
// while the source and target search paths differ: target create ordering must
// resolve it through the target path. The false provider depends on the consumer,
// so resolving through the source path would form a cycle and leave the real
// provider after its consumer.
func TestTargetSearchPathOrdersCreatedTypes(t *testing.T) {
	composite := func(name, attributeType string) *metadatapb.CompositeTypeMetadata {
		return &metadatapb.CompositeTypeMetadata{
			Name:       name,
			Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "v", Type: attributeType}},
		}
	}
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "afirst", CompositeTypes: []*metadatapb.CompositeTypeMetadata{composite("shared", "app.holder")}},
		{Name: "app", CompositeTypes: []*metadatapb.CompositeTypeMetadata{composite("holder", "shared")}},
		{Name: "zprovider", CompositeTypes: []*metadatapb.CompositeTypeMetadata{composite("shared", "integer")}},
	}}
	plan := &MigrationPlan{Ops: []MigrationOp{
		{Type: OpCreateType, SchemaName: "app", ObjectName: "holder"},
		{Type: OpCreateType, SchemaName: "afirst", ObjectName: "shared"},
		{Type: OpCreateType, SchemaName: "zprovider", ObjectName: "shared"},
	}}

	orderTypeAndTableCreates(plan, meta, []string{"zprovider", "afirst"})

	positions := make(map[string]int, len(plan.Ops))
	for i, op := range plan.Ops {
		positions[objectKey(op.SchemaName, op.ObjectName)] = i
	}
	provider := positions[objectKey("zprovider", "shared")]
	consumer := positions[objectKey("app", "holder")]
	require.Less(t, provider, consumer, "the target search path makes zprovider.shared the provider")
}

func TestAccessoryCommentsGoWithTheRestoredView(t *testing.T) {
	snapshot := func(columnType, indexComment, triggerComment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "lower",
					Definition: "SELECT FROM",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "t", Column: "id",
					}},
				},
				{
					Name:       "upper",
					Definition: "SELECT id FROM public.lower",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "lower",
					}},
					Triggers: []*metadatapb.TriggerMetadata{{
						Name: "trig", Comment: triggerComment,
						Body: `CREATE TRIGGER "trig" BEFORE INSERT ON "public"."upper" FOR EACH ROW EXECUTE FUNCTION public.touch()`,
					}},
				},
			},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv",
				Definition: "SELECT id FROM public.lower",
				Indexes:    []*metadatapb.IndexMetadata{{Name: "idx_mv", Expressions: []string{"id"}, Comment: indexComment}},
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "lower",
				}},
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer", "old", "old"), snapshot("bigint", "new", "new"))
	require.NoError(t, err)
	indexComment := strings.Index(sql, `COMMENT ON INDEX "public"."idx_mv" IS 'new'`)
	triggerComment := strings.Index(sql, `COMMENT ON TRIGGER "trig" ON "public"."upper" IS 'new'`)
	require.GreaterOrEqual(t, indexComment, 0, sql)
	require.GreaterOrEqual(t, triggerComment, 0, sql)
	require.Equal(t, 1, strings.Count(sql, `COMMENT ON INDEX "public"."idx_mv"`), "the index comment is emitted once, with the recreation")
	require.Equal(t, 1, strings.Count(sql, `COMMENT ON TRIGGER "trig" ON "public"."upper"`), "the trigger comment is emitted once, with the recreation")
	require.Less(t, strings.Index(sql, `CREATE MATERIALIZED VIEW "public"."mv"`), indexComment, sql)
	require.Less(t, strings.Index(sql, `CREATE VIEW "public"."upper"`), triggerComment, sql)
}

// TestTargetOnlyIndexCommentGoesWithItsMatView covers an index comment on a
// materialized view the target adds that depends on a restored view: the index
// exists only in the target metadata, and the comment must go with the recreation
// rather than run before it.
func TestTargetOnlyIndexCommentGoesWithItsMatView(t *testing.T) {
	snapshot := func(columnType string, addMaterializedView bool) *metadatapb.DatabaseSchemaMetadata {
		schemaMeta := &metadatapb.SchemaMetadata{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "lower",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
		}
		if addMaterializedView {
			schemaMeta.MaterializedViews = []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv",
				Definition: "SELECT id FROM public.lower",
				Indexes:    []*metadatapb.IndexMetadata{{Name: "idx_mv", Expressions: []string{"id"}, Comment: "the index"}},
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "lower",
				}},
			}}
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{schemaMeta}}
	}

	sql, err := testMetadataMigration(snapshot("integer", false), snapshot("bigint", true))
	require.NoError(t, err)
	comment := strings.Index(sql, `COMMENT ON INDEX "public"."idx_mv" IS 'the index'`)
	require.GreaterOrEqual(t, comment, 0, sql)
	require.Equal(t, 1, strings.Count(sql, `COMMENT ON INDEX "public"."idx_mv"`), "the index comment is emitted once, with the recreation")
	require.Less(t, strings.Index(sql, `CREATE MATERIALIZED VIEW "public"."mv"`), comment, sql)
}

// TestRemoveCompositeAttributeOpsKeepsUnrelatedComments keeps a comment whose
// value mentions the rebuilt attribute with its own object.
func TestRemoveCompositeAttributeOpsKeepsUnrelatedComments(t *testing.T) {
	ops := []MigrationOp{
		{Type: OpComment, SchemaName: "public", ObjectName: "other", SQL: `COMMENT ON TABLE "public"."other" IS 'mentions "public"."c"."a"'`},
		{Type: OpComment, SchemaName: "public", ObjectName: "c", SQL: `COMMENT ON COLUMN "public"."c"."a" IS 'rebuilt'`},
		{Type: OpComment, SchemaName: "public", ObjectName: "c", SQL: `COMMENT ON COLUMN "public"."c"."ab" IS 'sibling'`},
	}
	got := removeCompositeAttributeOps(ops, "public", "c", "a")
	require.Len(t, got, 2, "only the rebuilt attribute's comment is removed")
	require.Contains(t, got[0].SQL, `"public"."other"`)
	require.Contains(t, got[1].SQL, `"public"."c"."ab"`)
}

// TestRefuseAmbiguousIndexComments rejects a comment whose index name several of a
// schema's relations own while one of them is selected: nothing can attribute the
// comment, so the migration is refused rather than run against the wrong index.
func TestRefuseAmbiguousIndexComments(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{{
			Name:    "t",
			Indexes: []*metadatapb.IndexMetadata{{Name: "idx"}},
		}},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
			Name:    "mv",
			Indexes: []*metadatapb.IndexMetadata{{Name: "idx"}},
		}},
	}}}
	ops := []MigrationOp{{Type: OpComment, SQL: `COMMENT ON INDEX "public"."idx" IS 'z'`}}
	err := refuseAmbiguousIndexComments(ops, map[string]bool{objectKey("public", "mv"): true}, meta, meta)
	require.ErrorContains(t, err, "used by more than one relation")
	require.NoError(t, refuseAmbiguousIndexComments(ops, map[string]bool{objectKey("public", "other"): true}, meta, meta), "an unrelated selection keeps the comment")
}

// TestOpTargetsSelectedParsesCommentTargets covers COMMENT ON statements whose
// text contains a selected view's name without targeting it: an escaped quote can
// spell another owner, and a function of that name shares the schema.
func TestOpTargetsSelectedParsesCommentTargets(t *testing.T) {
	selected := map[string]bool{objectKey("a", "b"): true}
	tests := []struct {
		name   string
		sql    string
		claims bool
	}{
		{name: "view", sql: `COMMENT ON VIEW "a"."b" IS 'z'`, claims: true},
		{name: "materialized view", sql: `COMMENT ON MATERIALIZED VIEW "a"."b" IS 'z'`, claims: true},
		{name: "view column", sql: `COMMENT ON COLUMN "a"."b"."c" IS 'z'`, claims: true},
		{name: "escaped quote", sql: `COMMENT ON VIEW "x""a"."b" IS 'z'`},
		{name: "function", sql: `COMMENT ON FUNCTION "a"."b"(integer) IS 'z'`},
		{name: "other schema", sql: `COMMENT ON VIEW "x"."b" IS 'z'`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			op := MigrationOp{Type: OpComment, ObjectName: "a.b", SQL: tc.sql}
			require.Equal(t, tc.claims, opTargetsSelected(op, selected, nil, nil, nil), tc.sql)
		})
	}
}

// TestIndexNameOwnerIsUnique covers an index name that several of a schema's
// relations own: CockroachDB allows that, so the name alone claims nothing.
func TestIndexNameOwnerIsUnique(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{{
			Name:    "t",
			Indexes: []*metadatapb.IndexMetadata{{Name: "idx"}},
		}},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
			Name:    "mv",
			Indexes: []*metadatapb.IndexMetadata{{Name: "idx"}},
		}},
	}}}
	require.Empty(t, indexNameOwner(meta, "public", "idx", map[string]bool{objectKey("public", "mv"): true}), "an index name owned twice claims neither")
	require.Empty(t, indexNameOwner(meta, "public", "idx", map[string]bool{objectKey("public", "t"): true}), "an unselected relation keeps its index drop")
	require.Empty(t, indexNameOwner(meta, "public", "missing", map[string]bool{objectKey("public", "mv"): true}))
}

// TestOpTargetsSelectedKeepsATriggerNamedLikeTheView keeps a trigger whose name
// equals a selected view's name with the table it belongs to: a trigger operation
// carries the trigger's name, not the relation's.
func TestOpTargetsSelectedKeepsATriggerNamedLikeTheView(t *testing.T) {
	selected := map[string]bool{objectKey("public", "v"): true}
	trigger := MigrationOp{
		Type:         OpCreateTrigger,
		SchemaName:   "public",
		ObjectName:   "v",
		ParentObject: "t",
		SQL:          `CREATE TRIGGER "v" BEFORE INSERT ON "public"."t" FOR EACH ROW EXECUTE FUNCTION public.touch()`,
	}
	require.False(t, opTargetsSelected(trigger, selected, nil, nil, nil), "a trigger named like the view belongs to its own table")
	onTheView := MigrationOp{
		Type:         OpCreateTrigger,
		SchemaName:   "public",
		ObjectName:   "audit",
		ParentObject: "v",
		SQL:          `CREATE TRIGGER "audit" AFTER INSERT ON "public"."v" FOR EACH ROW EXECUTE FUNCTION public.touch()`,
	}
	require.True(t, opTargetsSelected(onTheView, selected, nil, nil, nil), "a trigger on the selected view belongs to it")
	view := MigrationOp{Type: OpDropView, SchemaName: "public", ObjectName: "v", SQL: `DROP VIEW "public"."v"`}
	require.True(t, opTargetsSelected(view, selected, nil, nil, nil))

	callsTheView := MigrationOp{
		Type:         OpCreateTrigger,
		SchemaName:   "public",
		ObjectName:   "x",
		ParentObject: "t",
		SQL:          `CREATE TRIGGER "x" BEFORE INSERT ON "public"."t" FOR EACH ROW EXECUTE FUNCTION "public"."v"()`,
	}
	require.False(t, opTargetsSelected(callsTheView, selected, nil, nil, nil), "a trigger belongs to its table, not to the function it calls")
}

// deferredViewSnapshot builds the pair used by the deferred-create tests: a target
// that adds public.upper above a restored public.lower.
func deferredViewSnapshot(columnType string, addUpper bool) *metadatapb.DatabaseSchemaMetadata {
	schemaMeta := &metadatapb.SchemaMetadata{
		Name:   "public",
		Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
		Views: []*metadatapb.ViewMetadata{{
			Name:       "lower",
			Definition: "SELECT FROM",
			DependencyColumns: []*metadatapb.DependencyColumn{{
				Schema: "public", Table: "t", Column: "id",
			}},
		}},
	}
	if addUpper {
		schemaMeta.Views = append(schemaMeta.Views, &metadatapb.ViewMetadata{
			Name:       "upper",
			Definition: "SELECT id FROM public.lower",
			DependencyColumns: []*metadatapb.DependencyColumn{{
				Schema: "public", Table: "lower",
			}},
		})
	}
	return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{schemaMeta}}
}

// TestDeferredCreatedViewRefusesARowTypeColumn covers a view the target adds on a
// restored view: its create is deferred with the rest of the recreation, so a
// table column of its row type cannot be ordered and the migration is refused.
func TestDeferredCreatedViewRefusesARowTypeColumn(t *testing.T) {
	source := deferredViewSnapshot("integer", false)
	target := deferredViewSnapshot("bigint", true)
	target.Schemas[0].Tables = append(target.Schemas[0].Tables,
		pgTable("holds", 0, &metadatapb.ColumnMetadata{Name: "row", Type: "public.upper"}))
	_, err := testMetadataMigration(source, target)
	require.ErrorContains(t, err, `column "row" of table "public"."holds" uses the row type of view "public"."upper"`)
}

// TestDeferredCreatedViewOrdersItsAttributeAfterwards covers an existing composite
// that gains an attribute of a view the target adds: the attribute is added after
// the view's deferred create.
func TestDeferredCreatedViewOrdersItsAttributeAfterwards(t *testing.T) {
	composite := func(columnType string, addAttribute bool) *metadatapb.DatabaseSchemaMetadata {
		meta := deferredViewSnapshot(columnType, addAttribute)
		attributes := []*metadatapb.CompositeTypeAttribute{{Name: "x", Type: "integer"}}
		if addAttribute {
			attributes = append(attributes, &metadatapb.CompositeTypeAttribute{Name: "r", Type: "public.upper"})
		}
		meta.Schemas[0].CompositeTypes = []*metadatapb.CompositeTypeMetadata{{Name: "c", Attributes: attributes}}
		return meta
	}
	sql, err := testMetadataMigration(composite("integer", false), composite("bigint", true))
	require.NoError(t, err)
	createView := strings.Index(sql, `CREATE VIEW "public"."upper"`)
	addAttribute := strings.Index(sql, `ADD ATTRIBUTE "r" public.upper`)
	require.GreaterOrEqual(t, createView, 0, sql)
	require.GreaterOrEqual(t, addAttribute, 0, sql)
	require.Less(t, createView, addAttribute, "the attribute is added after the view it needs is created")
}

// TestDeferredCreatedViewReleasesARetypedAttribute covers an existing attribute
// retyped to a view the target adds: the view has no drop to precede, but the
// release still has to happen before its deferred create, and the attribute is
// re-added after it.
func TestDeferredCreatedViewReleasesARetypedAttribute(t *testing.T) {
	snapshot := func(columnType, attributeType string) *metadatapb.DatabaseSchemaMetadata {
		schemaMeta := &metadatapb.SchemaMetadata{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "lower",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
			CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name:       "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "r", Type: attributeType}},
			}},
		}
		if attributeType != "text" {
			schemaMeta.Views = append(schemaMeta.Views, &metadatapb.ViewMetadata{
				Name:       "upper",
				Definition: "SELECT id FROM public.lower",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "lower",
				}},
			})
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{schemaMeta}}
	}
	sql, err := testMetadataMigration(snapshot("integer", "text"), snapshot("bigint", "public.upper"))
	require.NoError(t, err)
	release := strings.Index(sql, `ALTER TYPE "public"."c" DROP ATTRIBUTE "r"`)
	createView := strings.Index(sql, `CREATE VIEW "public"."upper"`)
	add := strings.Index(sql, `ALTER TYPE "public"."c" ADD ATTRIBUTE "r" public.upper`)
	for name, index := range map[string]int{"release": release, "create view": createView, "add": add} {
		require.GreaterOrEqual(t, index, 0, "%s missing from:\n%s", name, sql)
	}
	require.Less(t, release, createView, "the release happens before the deferred create")
	require.Less(t, createView, add, "the attribute is re-added after the create")
}

// TestRestoredViewRefusesRoutinesHoldingItsRowType covers an unchanged routine whose
// parameter or return type is a restored view's row type: the view drop fails while
// the routine holds it. Procedures are routines too, and a routine that only
// mentions the view in its body is not a dependency.
func TestRestoredViewRefusesRoutinesHoldingItsRowType(t *testing.T) {
	type world struct {
		viewName  string
		function  *metadatapb.FunctionMetadata
		procedure *metadatapb.ProcedureMetadata
	}
	snapshot := func(columnType string, w world) *metadatapb.DatabaseSchemaMetadata {
		viewName := w.viewName
		if viewName == "" {
			viewName = "v"
		}
		schemaMeta := &metadatapb.SchemaMetadata{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       viewName,
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
		}
		if w.function != nil {
			schemaMeta.Functions = []*metadatapb.FunctionMetadata{w.function}
		}
		if w.procedure != nil {
			schemaMeta.Procedures = []*metadatapb.ProcedureMetadata{w.procedure}
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{schemaMeta}}
	}
	tests := []struct {
		name    string
		world   world
		refused bool
	}{
		{
			name: "parameter",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "f",
				Signature:  "f(public.v)",
				Definition: "CREATE FUNCTION public.f(public.v) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$",
			}},
			refused: true,
		},
		{
			name: "return type",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "g",
				Signature:  "g()",
				Definition: "CREATE FUNCTION public.g() RETURNS public.v LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "return type beginning with table",
			world: world{viewName: "table_v", function: &metadatapb.FunctionMetadata{
				Name:       "gt",
				Signature:  "gt()",
				Definition: "CREATE FUNCTION public.gt() RETURNS public.table_v LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "procedure",
			world: world{procedure: &metadatapb.ProcedureMetadata{
				Name:       "p",
				Signature:  "p(public.v)",
				Definition: "CREATE PROCEDURE public.p(public.v) LANGUAGE sql AS $$ SELECT 1 $$",
			}},
			refused: true,
		},
		{
			name: "partly quoted",
			world: world{viewName: "V", function: &metadatapb.FunctionMetadata{
				Name:       "f",
				Signature:  `f(public."V")`,
				Definition: `CREATE FUNCTION public.f(public."V") RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$`,
			}},
			refused: true,
		},
		{
			name: "out parameter",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "o",
				Signature:  "o()",
				Definition: "CREATE FUNCTION public.o(OUT public.v) RETURNS record LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "out parameter name contains default",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "od",
				Signature:  "od()",
				Definition: "CREATE FUNCTION public.od(OUT mydefault public.v) RETURNS record LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "returns followed by a newline",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rn",
				Signature:  "rn()",
				Definition: "CREATE FUNCTION public.rn() RETURNS\npublic.v LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "returns followed by a tab",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rt",
				Signature:  "rt()",
				Definition: "CREATE FUNCTION public.rt() RETURNS public.v\tLANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "returns followed by another clause",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rp",
				Signature:  "rp()",
				Definition: "CREATE FUNCTION public.rp() RETURNS public.v PARALLEL SAFE LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "returns followed by called on null input",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rc",
				Signature:  "rc()",
				Definition: "CREATE FUNCTION public.rc() RETURNS public.v CALLED ON NULL INPUT LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "setof type on the next line",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rs",
				Signature:  "rs()",
				Definition: "CREATE FUNCTION public.rs() RETURNS SETOF\npublic.v LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "returns followed by not leakproof",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rnl",
				Signature:  "rnl()",
				Definition: "CREATE FUNCTION public.rnl() RETURNS public.v NOT LEAKPROOF LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "returns followed by external security",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "re",
				Signature:  "re()",
				Definition: "CREATE FUNCTION public.re() RETURNS public.v EXTERNAL SECURITY DEFINER LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "comment after return type",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rcm",
				Signature:  "rcm()",
				Definition: "CREATE FUNCTION public.rcm() RETURNS public.v /* note */ LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "comment after out parameter type",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rcp",
				Signature:  "rcp()",
				Definition: "CREATE FUNCTION public.rcp(OUT x public.v /* note */) RETURNS record LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "comment before out parameter type",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rcq",
				Signature:  "rcq()",
				Definition: "CREATE FUNCTION public.rcq(OUT x /* = */ public.v) RETURNS record LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "comment with parenthesis in out parameter",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rcr",
				Signature:  "rcr()",
				Definition: "CREATE FUNCTION public.rcr(OUT x public.v /* ( */) RETURNS record LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "comment between setof and type",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "rcs",
				Signature:  "rcs()",
				Definition: "CREATE FUNCTION public.rcs() RETURNS SETOF/* note */public.v LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "quoted parameter before an out parameter",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "q",
				Signature:  "q(text)",
				Definition: `CREATE FUNCTION public.q(")" text, OUT result public.v) RETURNS record LANGUAGE sql AS $$ SELECT $$`,
			}},
			refused: true,
		},
		{
			name: "quoted routine name",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       `f(`,
				Signature:  "f(text)",
				Definition: `CREATE FUNCTION public."f("(OUT result public.v) RETURNS record LANGUAGE sql AS $$ SELECT $$`,
			}},
			refused: true,
		},
		{
			name: "default mentioning the view",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "d",
				Signature:  "d(text)",
				Definition: `CREATE FUNCTION public.d(x text DEFAULT ' public.v') RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$`,
			}},
		},
		{
			name: "dollar-quoted default with a closing parenthesis",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "dd",
				Signature:  "dd(text)",
				Definition: `CREATE FUNCTION public.dd(x text DEFAULT $$)$$, OUT result public.v) RETURNS record LANGUAGE sql AS $$ SELECT $$`,
			}},
			refused: true,
		},
		{
			name: "unqualified parameter type",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "u",
				Signature:  "u(v)",
				Definition: "CREATE FUNCTION public.u(x v) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$",
			}},
			refused: true,
		},
		{
			name: "parameter named like the view",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "n",
				Signature:  "n(text)",
				Definition: "CREATE FUNCTION public.n(v text) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$",
			}},
		},
		{
			name: "returns table",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "t",
				Signature:  "t()",
				Definition: "CREATE FUNCTION public.t() RETURNS TABLE(result public.v) LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "returns table with a keyword column name",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "tk",
				Signature:  "tk()",
				Definition: "CREATE FUNCTION public.tk() RETURNS TABLE (x public.v, cost integer) LANGUAGE sql AS $$ SELECT $$",
			}},
			refused: true,
		},
		{
			name: "return type with a clause word in its name",
			world: world{viewName: "v LANGUAGE suffix", function: &metadatapb.FunctionMetadata{
				Name:       "w",
				Signature:  "w()",
				Definition: `CREATE FUNCTION public.w() RETURNS public."v LANGUAGE suffix" LANGUAGE sql AS $$ SELECT $$`,
			}},
			refused: true,
		},
		{
			name: "routine named like the view",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "v",
				Signature:  "v(integer)",
				Definition: "CREATE FUNCTION public.v(integer) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$",
			}},
		},
		{
			name: "returns inside a quoted routine name",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "RETURNS decoy",
				Signature:  "decoy()",
				Definition: `CREATE FUNCTION public."RETURNS decoy"() RETURNS public.v LANGUAGE sql AS $$ SELECT $$`,
			}},
			refused: true,
		},
		{
			name: "e-string default with an escaped quote",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "es",
				Signature:  "es(text)",
				Definition: `CREATE FUNCTION public.es(x text DEFAULT E'\')', OUT result public.v) RETURNS record LANGUAGE sql AS $$ SELECT $$`,
			}},
			refused: true,
		},
		{
			name: "multiword builtin ending in the view's name",
			world: world{viewName: "precision", function: &metadatapb.FunctionMetadata{
				Name:       "p",
				Signature:  "p(double precision)",
				Definition: "CREATE FUNCTION public.p(x double precision) RETURNS double precision LANGUAGE sql AS $$ SELECT 1 $$",
			}},
		},
		{
			name: "body reference",
			world: world{function: &metadatapb.FunctionMetadata{
				Name:       "h",
				Signature:  "h()",
				Definition: "CREATE FUNCTION public.h() RETURNS integer LANGUAGE sql AS $$ SELECT 1 FROM public.v $$",
			}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := testMetadataMigration(snapshot("integer", tc.world), snapshot("bigint", tc.world))
			if tc.refused {
				require.ErrorContains(t, err, "uses the row type of view")
				return
			}
			require.NoError(t, err, "a routine that only mentions the view in its body does not block the drop")
		})
	}
}

// TestUnqualifiedRoutineTypeFollowsTheSearchPath resolves a routine's unqualified
// parameter or result type the way PostgreSQL does: an earlier same-named view
// shadows the restored later view, while a search path that reaches the restored
// view first still refuses the migration.
func TestUnqualifiedRoutineTypeFollowsTheSearchPath(t *testing.T) {
	type routine struct {
		name       string
		signature  string
		definition string
	}
	tests := []routine{
		{
			name:       "parameter type",
			signature:  "f(v)",
			definition: "CREATE FUNCTION later.f(x v) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$",
		},
		{
			name:       "return type",
			signature:  "f()",
			definition: "CREATE FUNCTION later.f() RETURNS v LANGUAGE sql AS $$ SELECT 1 $$",
		},
	}
	snapshot := func(columnType, searchPath string, r routine) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			SearchPath: searchPath,
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name:  "earlier",
					Views: []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT 1"}},
				},
				{
					Name:   "later",
					Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
					Views: []*metadatapb.ViewMetadata{{
						Name:       "v",
						Definition: "SELECT FROM",
						DependencyColumns: []*metadatapb.DependencyColumn{{
							Schema: "later", Table: "t", Column: "id",
						}},
					}},
					Functions: []*metadatapb.FunctionMetadata{{
						Name:       "f",
						Signature:  r.signature,
						Definition: r.definition,
					}},
				},
			},
		}
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sql, err := testMetadataMigration(snapshot("integer", "earlier, later", tc),
				snapshot("bigint", "earlier, later", tc))
			require.NoError(t, err, "the routine depends on earlier.v, not the restored later.v")
			require.Contains(t, sql, `DROP VIEW IF EXISTS "later"."v"`, "the later view is still restored")

			_, err = testMetadataMigration(snapshot("integer", "later, earlier", tc),
				snapshot("bigint", "later, earlier", tc))
			require.ErrorContains(t, err, `routine "later"."f" uses the row type of view "later"."v"`, "the search path reaches the restored view first")
		})
	}

	t.Run("target search path", func(t *testing.T) {
		tc := tests[0]
		_, err := testMetadataMigration(snapshot("integer", "earlier, later", tc),
			snapshot("bigint", "later, earlier", tc))
		require.ErrorContains(t, err, `routine "later"."f" uses the row type of view "later"."v"`, "the target routine resolves through the target path")

		_, err = testMetadataMigration(snapshot("integer", "later, earlier", tc),
			snapshot("bigint", "earlier, later", tc))
		require.ErrorContains(t, err, `routine "later"."f" uses the row type of view "later"."v"`, "the source routine still resolves through the source path")
	})
}

// TestRestoredUnloadableViewRefusesRowTypeDependents keeps a table column, or an
// attribute the target adds, out of a restored view's drop: releasing a column
// would discard its data, and a new attribute would need the view before its
// ALTER TYPE, so the migration is refused instead.
func TestRestoredUnloadableViewRefusesRowTypeDependents(t *testing.T) {
	type world struct {
		tables         []*metadatapb.TableMetadata
		compositeTypes []*metadatapb.CompositeTypeMetadata
	}
	snapshot := func(columnType string, w world) *metadatapb.DatabaseSchemaMetadata {
		tables := append([]*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})}, w.tables...)
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: tables,
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
			CompositeTypes: w.compositeTypes,
		}}}
	}
	rowTypeColumn := func() *metadatapb.TableMetadata {
		return pgTable("holds", 0, &metadatapb.ColumnMetadata{Name: "r", Type: "public.v"})
	}
	tests := []struct {
		name   string
		source world
		target world
	}{
		{
			name:   "unchanged column",
			source: world{tables: []*metadatapb.TableMetadata{rowTypeColumn()}},
			target: world{tables: []*metadatapb.TableMetadata{rowTypeColumn()}},
		},
		{
			name:   "removed column",
			source: world{tables: []*metadatapb.TableMetadata{rowTypeColumn()}},
			target: world{tables: []*metadatapb.TableMetadata{pgTable("holds", 0)}},
		},
		{
			name:   "retyped column",
			source: world{tables: []*metadatapb.TableMetadata{rowTypeColumn()}},
			target: world{tables: []*metadatapb.TableMetadata{pgTable("holds", 0, &metadatapb.ColumnMetadata{Name: "r", Type: "text"})}},
		},
		{
			name:   "added column",
			source: world{tables: []*metadatapb.TableMetadata{pgTable("holds", 0)}},
			target: world{tables: []*metadatapb.TableMetadata{rowTypeColumn()}},
		},
		{
			name:   "added attribute",
			source: world{},
			target: world{compositeTypes: []*metadatapb.CompositeTypeMetadata{{
				Name:       "c",
				Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "r", Type: "public.v"}},
			}}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := testMetadataMigration(snapshot("integer", tc.source), snapshot("bigint", tc.target))
			require.ErrorContains(t, err, "uses the row type of view")
		})
	}
}

// TestLoadableDependentViewFollowsUnloadableView covers a loadable unchanged
// view that depends on an unloadable one: dropping the lower view fails while
// the upper still references it, so the closure must restore both.
func TestLoadableDependentViewFollowsUnloadableView(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "upper",
					Definition: "SELECT id FROM public.lower",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "lower",
					}},
				},
				{
					Name:       "lower",
					Definition: "SELECT FROM",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "t", Column: "id",
					}},
				},
			},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.NoError(t, err)
	dropUpper := strings.Index(sql, `DROP VIEW IF EXISTS "public"."upper"`)
	dropLower := strings.Index(sql, `DROP VIEW IF EXISTS "public"."lower"`)
	createLower := strings.Index(sql, `CREATE VIEW "public"."lower"`)
	createUpper := strings.Index(sql, `CREATE VIEW "public"."upper"`)
	require.GreaterOrEqual(t, dropUpper, 0, sql)
	require.GreaterOrEqual(t, dropLower, 0, sql)
	require.GreaterOrEqual(t, createLower, 0, sql)
	require.GreaterOrEqual(t, createUpper, 0, sql)
	require.Less(t, dropUpper, dropLower, sql)
	require.Less(t, createLower, createUpper, sql)
}

// TestChangedDependentViewIsDroppedBeforeUnloadableView covers a changed upper
// view: it is not recreated by the fallback, but it must still be dropped before
// the lower view it depends on.
func TestChangedDependentViewIsDroppedBeforeUnloadableView(t *testing.T) {
	snapshot := func(columnType, upperDefinition string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{
				{
					Name:       "upper",
					Definition: upperDefinition,
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "lower",
					}},
				},
				{
					Name:       "lower",
					Definition: "SELECT FROM",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "t", Column: "id",
					}},
				},
			},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer", "SELECT id FROM public.lower"),
		snapshot("bigint", "SELECT id, 1 AS x FROM public.lower"))
	require.NoError(t, err)
	dropUpper := strings.Index(sql, `DROP VIEW IF EXISTS "public"."upper"`)
	dropLower := strings.Index(sql, `DROP VIEW IF EXISTS "public"."lower"`)
	recreateUpper := strings.Index(sql, `CREATE VIEW "public"."upper"`)
	require.GreaterOrEqual(t, dropUpper, 0, sql)
	require.GreaterOrEqual(t, dropLower, 0, sql)
	require.GreaterOrEqual(t, recreateUpper, 0, sql)
	require.Less(t, dropUpper, dropLower, sql)
	require.Less(t, dropLower, recreateUpper, sql)
}

// TestChangedViewKeepsOldDependencyEdge covers a target definition that stops
// depending on the unloadable view: the currently installed view still does, so
// its drop must be ordered before the lower view's drop.
func TestChangedViewKeepsOldDependencyEdge(t *testing.T) {
	snapshot := func(columnType, upperDefinition string, upperDeps []*metadatapb.DependencyColumn) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{
				{Name: "upper", Definition: upperDefinition, DependencyColumns: upperDeps},
				{
					Name:       "lower",
					Definition: "SELECT FROM",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "t", Column: "id",
					}},
				},
			},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer", "SELECT id FROM public.lower", []*metadatapb.DependencyColumn{{Schema: "public", Table: "lower"}}),
		snapshot("bigint", "SELECT 1 AS one", nil))
	require.NoError(t, err)
	dropUpper := strings.Index(sql, `DROP VIEW IF EXISTS "public"."upper"`)
	dropLower := strings.Index(sql, `DROP VIEW IF EXISTS "public"."lower"`)
	require.GreaterOrEqual(t, dropUpper, 0, sql)
	require.GreaterOrEqual(t, dropLower, 0, sql)
	require.Less(t, dropUpper, dropLower, sql)
}

// TestCommentOnlyViewChangeIsRecreated covers a selected dependent whose only
// change is its comment: the catalog emits no replacement, so the fallback must
// synthesize one after the drop.
func TestCommentOnlyViewChangeIsRecreated(t *testing.T) {
	snapshot := func(columnType, comment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{
				{
					Name: "upper", Definition: "SELECT id FROM public.lower", Comment: comment,
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "lower"}},
				},
				{
					Name: "lower", Definition: "SELECT FROM",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
				},
			},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer", "old comment"), snapshot("bigint", "new comment"))
	require.NoError(t, err)
	dropUpper := strings.Index(sql, `DROP VIEW IF EXISTS "public"."upper"`)
	createUpper := strings.Index(sql, `CREATE VIEW "public"."upper"`)
	require.GreaterOrEqual(t, dropUpper, 0, sql)
	require.GreaterOrEqual(t, createUpper, 0, sql)
	require.Less(t, dropUpper, createUpper, sql)
	comment := strings.Index(sql, "COMMENT ON VIEW")
	require.GreaterOrEqual(t, comment, 0, sql)
	require.Less(t, createUpper, comment, sql)
}

// TestIndexOnlyMatViewChangeRecreatesBeforeItsIndex covers a selected dependent
// materialized view whose only change is an index: the synthesized definition
// carries the index, so the catalog index ops must not run before it.
func TestIndexOnlyMatViewChangeRecreatesBeforeItsIndex(t *testing.T) {
	snapshot := func(columnType string, unique bool) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name: "lower", Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv",
				Definition: "SELECT id FROM public.lower",
				Indexes:    []*metadatapb.IndexMetadata{{Name: "idx_mv", Expressions: []string{"id"}, Unique: unique}},
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "lower",
				}},
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer", false), snapshot("bigint", true))
	require.NoError(t, err)
	createView := strings.Index(sql, `CREATE MATERIALIZED VIEW "public"."mv"`)
	createIndex := strings.Index(sql, `CREATE UNIQUE INDEX "idx_mv"`)
	require.GreaterOrEqual(t, createView, 0, sql)
	require.GreaterOrEqual(t, createIndex, 0, sql)
	require.Less(t, createView, createIndex, sql)
	require.NotContains(t, sql, "DROP INDEX")
}

// TestIndexByNameUsesOwner covers snapshots where two relations carry an index
// of the same name: the lookup must use the relation named in the create SQL.
func TestIndexByNameUsesOwner(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{
			{Name: "a", Indexes: []*metadatapb.IndexMetadata{{Name: "idx", Definition: "definition a"}}},
			{Name: "b", Indexes: []*metadatapb.IndexMetadata{{Name: "idx", Definition: "definition b"}}},
		},
	}}}
	require.Equal(t, "definition a", indexByName(meta, "public", "public", "a", "idx").GetDefinition())
	require.Equal(t, "definition b", indexByName(meta, "public", "public", "b", "idx").GetDefinition())
}

// TestRemovedDependentViewDropsBeforeUnloadableView covers an old-only upper
// view: it is not recreated, but its drop must precede the lower view it
// depends on.
func TestRemovedDependentViewDropsBeforeUnloadableView(t *testing.T) {
	lower := func() *metadatapb.ViewMetadata {
		return &metadatapb.ViewMetadata{
			Name: "lower", Definition: "SELECT FROM",
			DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
		}
	}
	upper := &metadatapb.ViewMetadata{
		Name: "upper", Definition: "SELECT id FROM public.lower",
		DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "lower"}},
	}
	snapshot := func(columnType string, views ...*metadatapb.ViewMetadata) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views:  views,
		}}}
	}
	sql, err := testMetadataMigration(snapshot("integer", upper, lower()), snapshot("bigint", lower()))
	require.NoError(t, err)
	dropUpper := strings.Index(sql, `DROP VIEW IF EXISTS "public"."upper"`)
	dropLower := strings.Index(sql, `DROP VIEW IF EXISTS "public"."lower"`)
	require.GreaterOrEqual(t, dropUpper, 0, sql)
	require.GreaterOrEqual(t, dropLower, 0, sql)
	require.Less(t, dropUpper, dropLower, sql)
}

// TestSplitRetypeToCreatedType covers a composite attribute moving to a type
// this migration creates: the attribute is dropped before the old table and
// added after the new type exists.
func TestSplitRetypeToCreatedType(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:           "public",
		Tables:         []*metadatapb.TableMetadata{pgTable("old_t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
		CompositeTypes: []*metadatapb.CompositeTypeMetadata{{Name: "c", Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "r", Type: "public.old_t"}}}},
	}}}
	newMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:           "public",
		Tables:         []*metadatapb.TableMetadata{pgTable("new_t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
		EnumTypes:      []*metadatapb.EnumTypeMetadata{{Name: "e", Values: []string{"a", "b"}}},
		CompositeTypes: []*metadatapb.CompositeTypeMetadata{{Name: "c", Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "r", Type: "public.e"}}}},
	}}}
	sql, err := testMetadataMigration(oldMeta, newMeta)
	require.NoError(t, err)
	dropAttribute := strings.Index(sql, `DROP ATTRIBUTE "r"`)
	createType := strings.Index(sql, `CREATE TYPE "public"."e"`)
	addAttribute := strings.Index(sql, `ADD ATTRIBUTE "r" public.e`)
	require.GreaterOrEqual(t, dropAttribute, 0, sql)
	require.GreaterOrEqual(t, createType, 0, sql)
	require.GreaterOrEqual(t, addAttribute, 0, sql)
	require.Less(t, dropAttribute, createType, sql)
	require.Less(t, createType, addAttribute, sql)
}

// TestReversedDependentViewsUseSeparateDropAndCreateOrders covers dependents
// whose relationship reverses: drops follow the old graph, creates the target.
func TestReversedDependentViewsUseSeparateDropAndCreateOrders(t *testing.T) {
	view := func(name, definition string, deps ...string) *metadatapb.ViewMetadata {
		v := &metadatapb.ViewMetadata{Name: name, Definition: definition}
		for _, dep := range deps {
			v.DependencyColumns = append(v.DependencyColumns, &metadatapb.DependencyColumn{Schema: "public", Table: dep})
		}
		return v
	}
	snapshot := func(columnType string, reversed bool) *metadatapb.DatabaseSchemaMetadata {
		lower := view("lower", "SELECT FROM", "t")
		lower.Columns = []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}
		views := []*metadatapb.ViewMetadata{lower}
		if reversed {
			views = append(views, view("a", "SELECT id FROM public.lower", "lower"), view("b", "SELECT id FROM public.a", "a"))
		} else {
			views = append(views, view("b", "SELECT id FROM public.lower", "lower"), view("a", "SELECT id FROM public.b", "b"))
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views:  views,
		}}}
	}

	sql, err := testMetadataMigration(snapshot("integer", false), snapshot("bigint", true))
	require.NoError(t, err)
	dropB := strings.Index(sql, `DROP VIEW IF EXISTS "public"."b"`)
	dropA := strings.Index(sql, `DROP VIEW IF EXISTS "public"."a"`)
	dropLower := strings.Index(sql, `DROP VIEW IF EXISTS "public"."lower"`)
	createLower := strings.Index(sql, `CREATE VIEW "public"."lower"`)
	createA := strings.Index(sql, `CREATE VIEW "public"."a"`)
	createB := strings.Index(sql, `CREATE VIEW "public"."b"`)
	require.GreaterOrEqual(t, dropB, 0, sql)
	require.GreaterOrEqual(t, dropA, 0, sql)
	require.GreaterOrEqual(t, dropLower, 0, sql)
	require.GreaterOrEqual(t, createLower, 0, sql)
	require.GreaterOrEqual(t, createA, 0, sql)
	require.GreaterOrEqual(t, createB, 0, sql)
	require.Less(t, dropA, dropB, sql)
	require.Less(t, dropB, dropLower, sql)
	require.Less(t, createLower, createA, sql)
	require.Less(t, createA, createB, sql)
}

// TestAddedDependentViewFollowsRestoredView covers a target-only view above a
// restored view: it has no drop, but its create must follow the lower create.
func TestAddedDependentViewFollowsRestoredView(t *testing.T) {
	snapshot := func(columnType string, addUpper bool) *metadatapb.DatabaseSchemaMetadata {
		views := []*metadatapb.ViewMetadata{{
			Name: "lower", Definition: "SELECT FROM",
			DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
		}}
		if addUpper {
			views = append(views, &metadatapb.ViewMetadata{
				Name: "upper", Definition: "SELECT id FROM public.lower",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "lower"}},
			})
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views:  views,
		}}}
	}
	sql, err := testMetadataMigration(snapshot("integer", false), snapshot("bigint", true))
	require.NoError(t, err)
	createLower := strings.Index(sql, `CREATE VIEW "public"."lower"`)
	createUpper := strings.Index(sql, `CREATE VIEW "public"."upper"`)
	require.GreaterOrEqual(t, createLower, 0, sql)
	require.GreaterOrEqual(t, createUpper, 0, sql)
	require.Less(t, createLower, createUpper, sql)
	require.NotContains(t, sql, `DROP VIEW IF EXISTS "public"."upper"`)
}

// TestIndirectUnloadableDependentIsRestored covers an unloadable view whose
// dependency on the retyped column is through a loadable lower view.
func TestIndirectUnloadableDependentIsRestored(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{
				{
					Name: "lower", Definition: "SELECT id FROM public.t",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
				},
				{
					Name: "upper", Definition: "SELECT FROM",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "lower"}},
				},
			},
		}}}
	}
	sql, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.NoError(t, err)
	dropUpper := strings.Index(sql, `DROP VIEW IF EXISTS "public"."upper"`)
	createUpper := strings.Index(sql, `CREATE VIEW "public"."upper"`)
	require.GreaterOrEqual(t, dropUpper, 0, sql)
	require.GreaterOrEqual(t, createUpper, 0, sql)
	require.Less(t, dropUpper, createUpper, sql)
}

// TestAddedSchemaDependentViewFollowsRestoredView covers a view in a target-only
// schema: it is a create-only candidate and must follow the restored view it
// references.
func TestAddedSchemaDependentViewFollowsRestoredView(t *testing.T) {
	snapshot := func(columnType string, addSchema bool) *metadatapb.DatabaseSchemaMetadata {
		schemas := []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name: "lower", Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
			}},
		}}
		if addSchema {
			schemas = append(schemas, &metadatapb.SchemaMetadata{
				Name: "extra",
				Views: []*metadatapb.ViewMetadata{{
					Name: "upper", Definition: "SELECT id FROM public.lower",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "lower"}},
				}},
			})
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: schemas}
	}
	sql, err := testMetadataMigration(snapshot("integer", false), snapshot("bigint", true))
	require.NoError(t, err)
	createLower := strings.Index(sql, `CREATE VIEW "public"."lower"`)
	createUpper := strings.Index(sql, `CREATE VIEW "extra"."upper"`)
	require.GreaterOrEqual(t, createLower, 0, sql)
	require.GreaterOrEqual(t, createUpper, 0, sql)
	require.Less(t, createLower, createUpper, sql)
}

// TestObjectKeyIsUnambiguous covers identifiers containing dots: the key must
// not collide for schema "a.b"/view "c" and schema "a"/view "b.c".
func TestObjectKeyIsUnambiguous(t *testing.T) {
	require.NotEqual(t, objectKey("a.b", "c"), objectKey("a", "b.c"))
	schemaName, objectName := splitObjectKey(objectKey("a.b", "c"))
	require.Equal(t, "a.b", schemaName)
	require.Equal(t, "c", objectName)
}

// TestViewDependencyQualificationIgnoresLiterals covers a view whose text merely
// contains "alice.v" in a string literal while the actual unqualified reference
// resolves to public.v. The literal is not a qualified reference, so the
// unresolved $user path is still ambiguous.
func TestViewDependencyQualificationIgnoresLiterals(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "alice", Views: []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT 1"}}},
			{Name: "public", Views: []*metadatapb.ViewMetadata{
				{Name: "v", Definition: "SELECT 1"},
				{Name: "w", Definition: "SELECT 'alice.v', * FROM v", DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "alice", Table: "v"}}},
			}},
		},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "alice"."v"`)

	meta.Schemas[1].Views[1].Definition = "SELECT * FROM alice.v"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), "a parsed qualified reference is safe")
}

// TestSequenceDoesNotShadowARoutineType covers a sequence with the same name as a
// later relation: a sequence has no named row type, so the routine type resolves
// to the restored view and must refuse the view drop.
func TestSequenceDoesNotShadowARoutineType(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			SearchPath: "earlier, later",
			Schemas: []*metadatapb.SchemaMetadata{
				{Name: "earlier", Sequences: []*metadatapb.SequenceMetadata{{Name: "v", DataType: "bigint", Start: "1", Increment: "1"}}},
				{
					Name:   "later",
					Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
					Views: []*metadatapb.ViewMetadata{{
						Name:       "v",
						Definition: "SELECT FROM",
						DependencyColumns: []*metadatapb.DependencyColumn{{
							Schema: "later", Table: "t", Column: "id",
						}},
					}},
					Functions: []*metadatapb.FunctionMetadata{{
						Name:       "f",
						Signature:  "f(v)",
						Definition: "CREATE FUNCTION later.f(x v) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$",
					}},
				},
			},
		}
	}

	_, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.ErrorContains(t, err, `routine "later"."f" uses the row type of view "later"."v"`)
}

// TestRoutineChangeRestoresUnloadableView covers a routine replacement with no
// column change: the catalog has no edge from the unloadable view to the routine,
// so the view must still be restored around the routine operation.
func TestRoutineChangeRestoresUnloadableView(t *testing.T) {
	snapshot := func(body string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:  "public",
			Views: []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT FROM"}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "f",
				Signature:  "f()",
				Definition: "CREATE FUNCTION public.f() RETURNS integer LANGUAGE sql AS $$ SELECT " + body + " $$",
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("1"), snapshot("2"))
	require.NoError(t, err)
	require.Contains(t, sql, `DROP VIEW IF EXISTS "public"."v"`)
	require.Contains(t, sql, `CREATE VIEW "public"."v"`)
}

// TestPartitionColumnChangeRestoresUnloadableView covers a view that names a
// partition directly: the parent's column change is inherited, and the partition
// key must be in the affected set so the view is restored around the ALTER.
func TestPartitionColumnChangeRestoresUnloadableView(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
			{Name: "main", Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: columnType}},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name: "p", Schema: "parted",
				}},
			}}},
			{Name: "parted"},
			{Name: "public", Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "parted", Table: "p", Column: "id",
				}},
			}}},
		}}
	}

	sql, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.NoError(t, err)
	drop := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	alter := strings.Index(sql, `ALTER TABLE main.t ALTER COLUMN id TYPE bigint`)
	create := strings.Index(sql, `CREATE VIEW "public"."v"`)
	require.GreaterOrEqual(t, drop, 0, sql)
	require.GreaterOrEqual(t, alter, 0, sql)
	require.GreaterOrEqual(t, create, 0, sql)
	require.Less(t, drop, alter, sql)
	require.Less(t, alter, create, sql)
}

// TestAffectedTableColumnsIncludesRemovedTable covers a table removed in the
// target: every column, and every partition's effective key, is affected so a
// retained unloadable view is restored before the drop.
func TestAffectedTableColumnsIncludesRemovedTable(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main", Tables: []*metadatapb.TableMetadata{{
			Name:       "t",
			Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
			Partitions: []*metadatapb.TablePartitionMetadata{{Name: "p", Schema: "parted"}},
		}}},
		{Name: "parted"},
	}}
	newMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main"},
		{Name: "parted"},
	}}

	affected := affectedTableColumns(oldMeta, newMeta, map[string]bool{"bbdataarchive": true})
	require.True(t, affected[objectKey("main", "t")]["id"])
	require.True(t, affected[objectKey("parted", "p")]["id"])
}

// TestObsoleteRoutineDropsBeforeRestoredView covers a source routine that holds a
// restored view's row type but is removed from the target: its drop is moved
// before the synthesized view drop, so the migration is valid and need not be
// refused.
func TestObsoleteRoutineDropsBeforeRestoredView(t *testing.T) {
	snapshot := func(columnType string, withFunction bool) *metadatapb.DatabaseSchemaMetadata {
		schemas := []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "t", Column: "id",
				}},
			}},
		}}
		if withFunction {
			schemas[0].Functions = []*metadatapb.FunctionMetadata{{
				Name:       "f",
				Signature:  "f(v)",
				Definition: "CREATE FUNCTION public.f(x v) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$",
			}}
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: schemas}
	}

	sql, err := testMetadataMigration(snapshot("integer", true), snapshot("bigint", false))
	require.NoError(t, err)
	dropFunction := strings.Index(sql, `DROP FUNCTION`)
	dropView := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	createView := strings.Index(sql, `CREATE VIEW "public"."v"`)
	require.GreaterOrEqual(t, dropFunction, 0, sql)
	require.GreaterOrEqual(t, dropView, 0, sql)
	require.GreaterOrEqual(t, createView, 0, sql)
	require.Less(t, dropFunction, dropView, sql)
	require.Less(t, dropView, createView, sql)
}

// TestUnresolvedUserRejectsUnqualifiedRelation covers a dependency in a schema
// after the unresolved $user entry: skipping $user to find public.t does not
// make the unqualified reference safe, because the login's schema can shadow it.
func TestUnresolvedUserRejectsUnqualifiedRelation(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestUnresolvedUserRejectsOneUnqualifiedOccurrence covers a view that uses the
// dependency both qualified and unqualified: the second occurrence can bind
// through $user, so the whole dependency is ambiguous.
func TestUnresolvedUserRejectsOneUnqualifiedOccurrence(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t a JOIN t b ON a.id = b.id",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestUnresolvedUserAllowsAliasedQualifiedRelation covers an alias equal to the
// relation name: only the FROM reference is a relation, so a qualified relation
// with `AS t` is not ambiguous.
func TestUnresolvedUserAllowsAliasedQualifiedRelation(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t AS t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))

	meta.Schemas[0].Views[0].Definition = "SELECT * FROM public.t t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), "a bare alias is not a relation reference")
}

// TestViewDependencySeesSequences covers a view selecting from a sequence: the
// sequence is absent from the type universe but still participates in the
// unresolved $user ambiguity check.
func TestViewDependencySeesSequences(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			Sequences: []*metadatapb.SequenceMetadata{{Name: "s", DataType: "bigint", Start: "1", Increment: "1"}},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT last_value FROM s",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "s"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."s"`)
}

// TestViewDependencyQualificationSkipsCTEs covers a CTE named like the
// dependency: only the CTE body's qualified relation is a relation reference,
// so the final FROM t is safe.
func TestViewDependencyQualificationSkipsCTEs(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "WITH t AS (SELECT * FROM public.t) SELECT * FROM t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationScopesCTEs covers a nested CTE: its name is not
// visible after the subquery closes, so the final unqualified reference is still
// ambiguous even though an earlier occurrence was qualified.
func TestViewDependencyQualificationScopesCTEs(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t WHERE EXISTS (WITH t AS (SELECT * FROM public.t) SELECT * FROM t) UNION ALL SELECT * FROM t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "WITH RECURSIVE t AS (SELECT * FROM t) SELECT * FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), "a recursive CTE may reference itself")

	meta.Schemas[0].Views[0].Definition = "WITH RECURSIVE x AS (SELECT * FROM t), t AS (SELECT 1) SELECT * FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), "a recursive CTE may reference a later sibling")
}

// TestSnapshotObjectsIncludesPartitionConstraints covers a partition CHECK
// constraint reported under con:<schema>.<partition>.<name>.
func TestSnapshotObjectsIncludesPartitionConstraints(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main", Tables: []*metadatapb.TableMetadata{{
			Name: "t",
			Partitions: []*metadatapb.TablePartitionMetadata{{
				Name:             "p",
				Schema:           "parted",
				CheckConstraints: []*metadatapb.CheckConstraintMetadata{{Name: "c", Expression: "missing <> ''"}},
			}},
		}}},
		{Name: "parted"},
	}}
	objects := snapshotObjects(meta)
	require.NotNil(t, objects["con:parted.p.c"])
}

// TestBuiltinMultiwordTypeIsCaseInsensitive covers SQL keywords in metadata
// spelled in upper or mixed case; only quoted names remain case-sensitive.
func TestBuiltinMultiwordTypeIsCaseInsensitive(t *testing.T) {
	require.True(t, builtinMultiwordType("DOUBLE PRECISION"))
	require.True(t, builtinMultiwordType("char varying"))
	require.True(t, builtinMultiwordType("char varying(16)"))
	require.True(t, builtinMultiwordType("national character"))
	require.True(t, builtinMultiwordType("national char"))
	require.True(t, builtinMultiwordType("national character varying(32)"))
	require.True(t, builtinMultiwordType("national char varying"))
	require.True(t, builtinMultiwordType("nchar varying"))
	require.True(t, builtinMultiwordType("character\n\tvarying [ ]"))
	require.True(t, builtinMultiwordType("character varying array[4]"))
	require.True(t, builtinMultiwordType("timestamp WITHOUT TIME ZONE"))
	require.True(t, builtinMultiwordType("timestamp(3) with time zone"))
	require.True(t, builtinMultiwordType("time(6) without time zone"))
	require.True(t, builtinMultiwordType("interval(6) day to second(3)"))
	require.True(t, builtinMultiwordType("interval day to second"))
	require.True(t, builtinMultiwordType("interval year to month"))
	require.False(t, builtinMultiwordType("public.PRECISION"))
	require.False(t, builtinMultiwordType(`character "varying"`))

	objects := map[string]bool{
		objectKey("public", "varying"):   true,
		objectKey("public", "character"): true,
	}
	searchPath := []string{"$user", "public"}
	require.False(t, ambiguousRoutineTypes("f()", "CREATE FUNCTION app.f() RETURNS char varying LANGUAGE sql AS $$SELECT 'x'$$", objects, searchPath))
	require.False(t, ambiguousRoutineTypes("g()", "CREATE FUNCTION app.g() RETURNS national character varying LANGUAGE sql AS $$SELECT 'x'$$", objects, searchPath))
}

// TestReleasedSourceOnlyRoutineDropsAreSorted covers the map iteration order of
// source-only routine drops that depend on a restored view.
func TestReleasedSourceOnlyRoutineDropsAreSorted(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:  "public",
		Views: []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT 1"}},
		Functions: []*metadatapb.FunctionMetadata{
			{Name: "b", Signature: "b(v)"},
			{Name: "a", Signature: "a(v)"},
		},
	}}}
	sourceDrops := map[string]MigrationOp{
		"func:public.b|b(v)": {Type: OpDropFunction, SchemaName: "public", ObjectName: "b", SQL: `DROP FUNCTION "public"."b"(v)`},
		"func:public.a|a(v)": {Type: OpDropFunction, SchemaName: "public", ObjectName: "a", SQL: `DROP FUNCTION "public"."a"(v)`},
	}
	ops, moved := releaseRoutineOps(&MigrationPlan{}, oldMeta, map[string]bool{objectKey("public", "v"): true}, []string{"public"}, sourceDrops)
	require.Len(t, ops, 2)
	require.Equal(t, "a", ops[0].ObjectName)
	require.Equal(t, "b", ops[1].ObjectName)
	require.True(t, moved[routineKey("public", "a", "v")])
	require.Empty(t, sourceDrops)
}

// TestSourceOnlyRoutineDropForcesViewRestoration covers a routine that exists only
// in source and is reported missing: its synthesized drop must force the
// unloadable view restoration even when no table column changed.
func TestSourceOnlyRoutineDropForcesViewRestoration(t *testing.T) {
	snapshot := func(withFunction bool) *metadatapb.DatabaseSchemaMetadata {
		schemas := []*metadatapb.SchemaMetadata{{Name: "public", Views: []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT FROM"}}}}
		if withFunction {
			schemas[0].Functions = []*metadatapb.FunctionMetadata{{
				Name:       "f",
				Signature:  "f(v)",
				Definition: "CREATE FUNCTION public.f(x v) RETURNS integer LANGUAGE sql AS @@invalid@@",
			}}
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: schemas}
	}

	sql, err := testMetadataMigration(snapshot(true), snapshot(false))
	require.NoError(t, err)
	dropFunction := strings.Index(sql, `DROP FUNCTION`)
	dropView := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	createView := strings.Index(sql, `CREATE VIEW "public"."v"`)
	require.GreaterOrEqual(t, dropFunction, 0, sql)
	require.GreaterOrEqual(t, dropView, 0, sql)
	require.GreaterOrEqual(t, createView, 0, sql)
	require.Less(t, dropFunction, dropView, sql)
}

// TestCTEVisibilityFollowsDeclarationOrder covers a later CTE not being visible
// to an earlier CTE body: the first `t` still names the dependency relation.
func TestCTEVisibilityFollowsDeclarationOrder(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "WITH x AS (SELECT * FROM t), t AS (SELECT 1) SELECT * FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestViewDependencyQualificationSeesTableQuery covers PostgreSQL's TABLE
// relation query form: the unqualified `TABLE t` is still a relation reference.
func TestViewDependencyQualificationSeesTableQuery(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t UNION ALL TABLE t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestProviderRelationDropRestoresUnloadableView covers an unloadable view that
// depends on a materialized view whose definition changes: the catalog drops the
// materialized view, so the dependent view must be restored around that drop even
// though no column or routine changed.
func TestProviderRelationDropRestoresUnloadableView(t *testing.T) {
	snapshot := func(materializedDefinition string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv",
				Definition: materializedDefinition,
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v",
				Definition: "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "mv",
				}},
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("SELECT 1"), snapshot("SELECT 2"))
	require.NoError(t, err)
	dropView := strings.Index(sql, `DROP VIEW IF EXISTS "public"."v"`)
	dropProvider := strings.Index(sql, `DROP MATERIALIZED VIEW`)
	createView := strings.Index(sql, `CREATE VIEW "public"."v"`)
	require.GreaterOrEqual(t, dropView, 0, sql)
	require.GreaterOrEqual(t, dropProvider, 0, sql)
	require.GreaterOrEqual(t, createView, 0, sql)
	require.Less(t, dropView, dropProvider, sql)
	require.Less(t, dropProvider, createView, sql)
}

// TestViewDependencyQualificationKeepsOuterRelationList covers a derived table:
// the parentheses must not lose the outer FROM-list state, so the final `t` is
// still an unqualified relation reference.
func TestViewDependencyQualificationKeepsOuterRelationList(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t, (SELECT 1) x, t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestCTENameMatchesEquivalentQuotedSpelling covers a lowercase quoted CTE name:
// `"t"` is the same identifier as unquoted `t`, but `"T"` is not.
func TestCTENameMatchesEquivalentQuotedSpelling(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `WITH "t" AS (SELECT * FROM public.t) SELECT * FROM public.t UNION ALL SELECT * FROM t`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))

	meta.Schemas[0].Views[0].Definition = `WITH "T" AS (SELECT * FROM public.t) SELECT * FROM public.t UNION ALL SELECT * FROM t`
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestViewDependencyQualificationSeesParenthesizedJoin covers a parenthesized
// joined-table group: the first relation inside it is still a relation reference.
func TestViewDependencyQualificationSeesParenthesizedJoin(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t a, (t JOIN public.u u ON true) AS j",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestProviderRelationDropLeavesUnrelatedViews covers a materialized-view drop:
// only the unloadable view whose dependency graph reaches the dropped relation is
// restored, so an unrelated populated materialized view is not rebuilt WITHOUT
// DATA.
func TestProviderRelationDropLeavesUnrelatedViews(t *testing.T) {
	snapshot := func(materializedDefinition string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{
				{Name: "changed", Definition: materializedDefinition},
				{Name: "unrelated", Definition: "SELECT FROM"},
			},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "v",
				Definition:        "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "changed"}},
			}},
		}}}
	}

	sql, err := testMetadataMigration(snapshot("SELECT 1"), snapshot("SELECT 2"))
	require.NoError(t, err)
	require.Contains(t, sql, `DROP VIEW IF EXISTS "public"."v"`, sql)
	require.NotContains(t, sql, `"unrelated"`, sql)
}

// TestViewDependencyQualificationSkipsAliasColumnList covers `FROM public.t AS
// x(t)`: the parenthesized alias column list is not a joined relation.
func TestViewDependencyQualificationSkipsAliasColumnList(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t AS x(t)",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationHandlesNonASCIIIdentifiers covers PostgreSQL's
// non-ASCII unquoted identifiers.
func TestViewDependencyQualificationHandlesNonASCIIIdentifiers(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("名", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT * FROM public."名" JOIN 名 ON true`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "名"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."名"`)
}

// TestMarkSchemaRelationsAffected covers a removed schema: its tables, views,
// materialized views and sequences become affected dependencies.
func TestMarkSchemaRelationsAffected(t *testing.T) {
	schema := &metadatapb.SchemaMetadata{
		Name:              "gone",
		Tables:            []*metadatapb.TableMetadata{{Name: "t"}},
		Views:             []*metadatapb.ViewMetadata{{Name: "v"}},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "mv"}},
		Sequences:         []*metadatapb.SequenceMetadata{{Name: "s"}},
	}
	affected := make(affectedColumnSet)
	markSchemaRelationsAffected(affected, schema)
	for _, key := range []string{"gone.t", "gone.v", "gone.mv", "gone.s"} {
		require.Contains(t, affected, objectKey("gone", strings.TrimPrefix(key, "gone.")), key)
	}
}

// TestViewDependencyQualificationSeesRegclass covers a qualified regclass
// literal: it names the sequence explicitly, while an unqualified regclass still
// follows the ambiguous $user path.
func TestViewDependencyQualificationSeesRegclass(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			Sequences: []*metadatapb.SequenceMetadata{{Name: "s", DataType: "bigint", Start: "1", Increment: "1"}},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT nextval('public.s'::regclass)",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "s"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))

	meta.Schemas[0].Views[0].Definition = "SELECT nextval('s'::regclass)"
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."s"`)
}

// TestViewDependencyQualificationSkipsDerivedTableAlias covers a bare alias on a
// derived table: the alias `t` is not a second relation reference.
func TestViewDependencyQualificationSkipsDerivedTableAlias(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t base, (SELECT 1) t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSkipsTablesample covers TABLESAMPLE after a
// qualified relation: the sampling method is not a relation reference.
func TestViewDependencyQualificationSkipsTablesample(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("system", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT * FROM public."system" TABLESAMPLE SYSTEM (10)`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "system"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestSourceOnlyConstraintIndexUsesConstraintKey covers a constraint-backed index
// reported missing under its constraint key: the synthesized drop must use con:.
func TestSourceOnlyConstraintIndexUsesConstraintKey(t *testing.T) {
	withIndex := func() *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t",
				Indexes: []*metadatapb.IndexMetadata{{
					Name:         "t_pkey_idx",
					IsConstraint: true,
					Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "t_pkey"}},
				}},
			}},
		}}}
	}
	withoutIndex := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public", Tables: []*metadatapb.TableMetadata{{Name: "t"}}}}}
	drops := sourceOnlyMissingDrops(withIndex(), withoutIndex, &LoadMetadataReport{Missing: map[string]error{"con:public.t.t_pkey": nil}})
	drop, ok := drops["con:public.t.t_pkey"]
	require.True(t, ok)
	require.Equal(t, OpDropConstraint, drop.Type)
	require.Equal(t, "t_pkey", drop.ObjectName)
}

// TestInsertSourceOnlyDropsSkipsSchemaCoveredChildren covers a source-only child
// whose schema is dropped: the schema CASCADE already removes it.
func TestInsertSourceOnlyDropsSkipsSchemaCoveredChildren(t *testing.T) {
	ops := []MigrationOp{{Type: OpDropSchema, SchemaName: "gone", Phase: PhasePre}}
	drops := map[string]MigrationOp{
		"type:gone.t": {Type: OpDropType, SchemaName: "gone", ObjectName: "t"},
		"seq:gone.s":  {Type: OpDropSequence, SchemaName: "gone", ObjectName: "s"},
	}
	result := insertSourceOnlyDrops(ops, drops)
	require.Len(t, result, 1)
	require.Equal(t, OpDropSchema, result[0].Type)
}

// TestSnapshotObjectsUsesConstraintName covers an index whose constraint has a
// different name from the index: the report key uses the constraint name.
func TestSnapshotObjectsUsesConstraintName(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "public",
		Tables: []*metadatapb.TableMetadata{{
			Name: "t",
			Indexes: []*metadatapb.IndexMetadata{{
				Name:         "t_pkey_idx",
				IsConstraint: true,
				Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "t_pkey"}},
			}},
		}},
	}}}
	objects := snapshotObjects(meta)
	require.NotNil(t, objects["con:public.t.t_pkey"])
}

// TestMarkDroppedSchemaAffectedIncludesCrossSchemaPartition covers a partition in
// the dropped schema whose parent table lives elsewhere.
func TestMarkDroppedSchemaAffectedIncludesCrossSchemaPartition(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
		{Name: "main", Tables: []*metadatapb.TableMetadata{{
			Name:       "t",
			Partitions: []*metadatapb.TablePartitionMetadata{{Name: "p", Schema: "parted"}},
		}}},
		{Name: "parted"},
	}}
	affected := make(affectedColumnSet)
	markDroppedSchemaAffected(affected, meta, "parted")
	require.Contains(t, affected, objectKey("parted", "p"))
}

// TestCTESearchAndCycleClausesDoNotHideSiblings covers SEARCH and CYCLE clauses
// between a CTE body and the next sibling declaration.
func TestCTESearchAndCycleClausesDoNotHideSiblings(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "WITH RECURSIVE x AS (SELECT 1) SEARCH DEPTH FIRST BY id SET ord, t AS (SELECT 1) SELECT * FROM public.t UNION ALL SELECT * FROM t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))

	meta.Schemas[0].Views[0].Definition = "WITH RECURSIVE x AS (SELECT 1) CYCLE id SET is_cycle USING path, t AS (SELECT 1) SELECT * FROM public.t UNION ALL SELECT * FROM t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestCompositeReorderInBackupSchemaIsIgnored covers backup-only metadata drift:
// the schema is filtered from the generated plan, so its reorder must not fail
// validation before the filter.
func TestCompositeReorderInBackupSchemaIsIgnored(t *testing.T) {
	snapshot := func(reordered bool) *metadatapb.DatabaseSchemaMetadata {
		attributes := []*metadatapb.CompositeTypeAttribute{{Name: "a", Type: "text"}, {Name: "b", Type: "integer"}}
		if reordered {
			attributes[0], attributes[1] = attributes[1], attributes[0]
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public"},
			{Name: "bbdataarchive", CompositeTypes: []*metadatapb.CompositeTypeMetadata{{Name: "addr", Attributes: attributes}}},
		}}
	}
	sql, err := testMetadataMigration(snapshot(false), snapshot(true))
	require.NoError(t, err)
	require.Empty(t, sql)
}

// TestCTECollectsNestedWith covers a WITH declared inside another CTE body.
func TestCTECollectsNestedWith(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "WITH x AS (WITH t AS (SELECT 1) SELECT * FROM t) SELECT * FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesUnqualifiedCast covers a row-type cast
// outside the FROM list: `NULL::t` still depends on the table's row type.
func TestViewDependencyQualificationSeesUnqualifiedCast(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT NULL::t AS r FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT NULL::public.t AS r FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestUnloadableBackupChangeIsIgnored covers changed unparseable metadata in
// bbdataarchive: the plan filters it, so preflight must not refuse it.
func TestUnloadableBackupChangeIsIgnored(t *testing.T) {
	snapshot := func(definition string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public"},
			{Name: "bbdataarchive", Views: []*metadatapb.ViewMetadata{{Name: "v", Definition: definition}}},
		}}
	}
	sql, err := testMetadataMigration(snapshot("SELECT @@@old@@@"), snapshot("SELECT @@@new@@@"))
	require.NoError(t, err)
	require.Empty(t, sql)
}

// TestViewDependencyQualificationSeesSQLStandardCast covers CAST(NULL AS t):
// the type position is a row-type dependency like a :: cast.
func TestViewDependencyQualificationSeesSQLStandardCast(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT CAST(NULL AS t) FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT CAST(NULL AS public.t) FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestBackupRoutineChangeDoesNotRestoreManagedMatview covers a backup-only
// routine change: it must not force a managed unloadable materialized view to be
// rebuilt WITHOUT DATA.
func TestBackupRoutineChangeDoesNotRestoreManagedMatview(t *testing.T) {
	snapshot := func(body string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "mv", Definition: "SELECT FROM"}}},
			{Name: "bbdataarchive", Functions: []*metadatapb.FunctionMetadata{{
				Name:       "f",
				Signature:  "f()",
				Definition: "CREATE FUNCTION bbdataarchive.f() RETURNS integer LANGUAGE sql AS $$ SELECT " + body + " $$",
			}}},
		}}
	}
	sql, err := testMetadataMigration(snapshot("1"), snapshot("2"))
	require.NoError(t, err)
	require.Empty(t, sql)
}

// TestRoutineReturnsClauseIsCaseInsensitive covers lowercase SQL in a routine
// definition.
func TestRoutineReturnsClauseIsCaseInsensitive(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "v",
				Definition:        "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "f",
				Signature:  "f()",
				Definition: "create function public.f() returns public.v language sql as $$ select 1 $$",
			}},
		}}}
	}
	_, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.ErrorContains(t, err, `routine "public"."f" uses the row type of view "public"."v"`)
}

// TestViewDependencyQualificationSeesCastRegclass covers CAST('s' AS regclass)
// alongside a qualified regclass literal.
func TestViewDependencyQualificationSeesCastRegclass(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			Sequences: []*metadatapb.SequenceMetadata{{Name: "s", DataType: "bigint", Start: "1", Increment: "1"}},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT nextval(CAST('s' AS regclass)), nextval('public.s'::regclass)",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "s"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."s"`)

	meta.Schemas[0].Views[0].Definition = "SELECT nextval(CAST('public.s' AS regclass)), nextval('public.s'::regclass)"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesThreePartRelation covers CockroachDB's
// database.schema.relation form.
func TestViewDependencyQualificationSeesThreePartRelation(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM db.public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesIsOf covers IS OF type lists.
func TestViewDependencyQualificationSeesIsOf(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT NULL::public.t IS OF (t) FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT NULL::public.t IS OF (public.t) FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesTypedConstant covers `type 'literal'`
// typed-constant expressions outside relation lists.
func TestViewDependencyQualificationSeesTypedConstant(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT NULL::public.t, t '(1)' FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT NULL::public.t, public.t '(1)' FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSkipsBuiltinCast covers a built-in type whose
// name also matches the dependency relation.
func TestViewDependencyQualificationSkipsBuiltinCast(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("text", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT NULL::text FROM public.text",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "text"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSkipsFromFunctionCall covers a set-returning
// function in FROM whose name matches the dependency relation.
func TestViewDependencyQualificationSkipsFromFunctionCall(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("generate_series", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.generate_series, generate_series(1, 2)",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "generate_series"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationScansFunctionArguments covers a subquery inside a
// FROM function argument list.
func TestViewDependencyQualificationScansFunctionArguments(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t base, unnest(ARRAY(SELECT id FROM t)) x",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestViewDependencyQualificationSkipsStringPrefix covers an escape-string
// prefix that is also the dependency name.
func TestViewDependencyQualificationSkipsStringPrefix(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("e", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT E'\n' FROM public.e`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "e"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestUnloadableBackupColumnChangeIsIgnored covers a backup-only column change
// that a managed unloadable view references.
func TestUnloadableBackupColumnChangeIsIgnored(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
					Name:       "mv",
					Definition: "SELECT FROM",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "bbdataarchive", Table: "t", Column: "id",
					}},
				}},
			},
			{
				Name:   "bbdataarchive",
				Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			},
		}}
	}
	sql, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.NoError(t, err)
	require.Empty(t, sql)
}

// TestViewDependencyQualificationSkipsUnicodeStringPrefix covers U&'...'.
func TestViewDependencyQualificationSkipsUnicodeStringPrefix(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("u", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT U&'d\0061t\+000061' FROM public.u`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "u"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesRegtype covers a regtype literal like a
// regclass literal.
func TestViewDependencyQualificationSeesRegtype(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT 't'::regtype, NULL::public.t FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT 'public.t'::regtype, NULL::public.t FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestReleasedRoutineKeysIncludeSignature covers two overloads with one name:
// only the overload that uses the view row type is released.
func TestReleasedRoutineKeysIncludeSignature(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:  "public",
		Views: []*metadatapb.ViewMetadata{{Name: "v", Definition: "SELECT 1"}},
		Functions: []*metadatapb.FunctionMetadata{
			{Name: "f", Signature: "f(v)"},
			{Name: "f", Signature: "f(integer)"},
		},
	}}}
	sourceDrops := map[string]MigrationOp{
		"func:public.f|f(integer)": {Type: OpDropFunction, SchemaName: "public", ObjectName: "f", SQL: `DROP FUNCTION "public"."f"(integer)`},
	}
	ops, moved := releaseRoutineOps(&MigrationPlan{}, oldMeta, map[string]bool{objectKey("public", "v"): true}, []string{"public"}, sourceDrops)
	require.Empty(t, ops)
	require.Empty(t, moved)
	require.Contains(t, sourceDrops, "func:public.f|f(integer)", "the unrelated overload is not released")
}

// TestViewDependencyQualificationSkipsRowsFrom covers ROWS FROM (...), whose
// ROWS keyword is not a relation named rows.
func TestViewDependencyQualificationSkipsRowsFrom(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("rows", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.rows r CROSS JOIN ROWS FROM (generate_series(1, 2)) f(n)",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "rows"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSkipsNestedBlockComment covers PostgreSQL's
// nested block comments.
func TestViewDependencyQualificationSkipsNestedBlockComment(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t /* outer /* inner */ FROM t */",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationKeepsQuotedTypedConstant covers a quoted type
// name before a literal: it is a typed constant, not a string prefix.
func TestViewDependencyQualificationKeepsQuotedTypedConstant(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("e", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT "e" '(1)', NULL::public.e FROM public.e`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "e"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."e"`)
}

// TestBackupViewIsExcludedFromAmbiguityValidation covers an unchanged backup view
// whose unqualified reference would otherwise block a managed migration.
func TestBackupViewIsExcludedFromAmbiguityValidation(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name:   "public",
				Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			},
			{
				Name: "bbdataarchive",
				Views: []*metadatapb.ViewMetadata{{
					Name:              "v",
					Definition:        "SELECT * FROM t",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
				}},
			},
		},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesDollarRegclass covers a dollar-quoted
// regclass literal.
func TestViewDependencyQualificationSeesDollarRegclass(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT $$t$$::regclass, 'public.t'::regclass FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT $$public.t$$::regclass, 'public.t'::regclass FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestBackupTypeAmbiguityIsIgnored covers a backup column with an unqualified
// user-defined type under an unresolved $user path.
func TestBackupTypeAmbiguityIsIgnored(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
					Name:       "shared",
					Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "x", Type: "text"}},
				}},
			},
			{
				Name:   "bbdataarchive",
				Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "c", Type: "shared"})},
			},
		},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestRoutineReturnTypeNamedLikeClauseKeyword covers `RETURNS public.cost COST 1`.
func TestRoutineReturnTypeNamedLikeClauseKeyword(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "cost",
				Definition:        "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "f",
				Signature:  "f()",
				Definition: "CREATE FUNCTION public.f() RETURNS public.cost COST 1 LANGUAGE sql AS $$ SELECT $$",
			}},
		}}}
	}
	_, err := testMetadataMigration(snapshot("integer"), snapshot("bigint"))
	require.ErrorContains(t, err, `routine "public"."f" uses the row type of view "public"."cost"`)
}

// TestViewDependencyQualificationSeesUnicodeEscapedIdentifier covers U&"t".
func TestViewDependencyQualificationSeesUnicodeEscapedIdentifier(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT * FROM public.t p, U&"t" q`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestPrefixedLiteralRegclass covers E't'::regclass.
func TestPrefixedLiteralRegclass(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT E't'::regclass, 'public.t'::regclass FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestUnicodeEscapedIdentifierDecodes covers U&"\0074" and UESCAPE.
func TestUnicodeEscapedIdentifierDecodes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT * FROM public.t p, U&"\0074" q`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = `SELECT * FROM public.t p, U&"!0074" UESCAPE '!' q`
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestEscapedRegclassLiteralDecodes covers U&'...' and E'...' regclass literals.
func TestEscapedRegclassLiteralDecodes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT U&'\\0074'::regclass, 'public.t'::regclass FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT E'\\164'::regclass, 'public.t'::regclass FROM public.t"
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestViewDependencyQualificationStopsAtFetch covers FETCH FIRST ROW ONLY.
func TestViewDependencyQualificationStopsAtFetch(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("first", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.first FETCH FIRST ROW ONLY",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "first"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationKeepsOperatorAdjacency covers `t || 'x'`: the
// operator prevents treating the column as a typed constant.
func TestViewDependencyQualificationKeepsOperatorAdjacency(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "t", Type: "text"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT t || 'x' FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesDollarTypedConstant covers `t $$(1)$$`.
func TestViewDependencyQualificationSeesDollarTypedConstant(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT t $$(1)$$, NULL::public.t FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestViewDependencyQualificationScansTablesampleArguments covers a subquery in
// a TABLESAMPLE argument.
func TestViewDependencyQualificationScansTablesampleArguments(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t TABLESAMPLE system ((SELECT 1 FROM t LIMIT 1))",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestViewDependencyQualificationSkipsWithOrdinality covers WITH ORDINALITY.
func TestViewDependencyQualificationSkipsWithOrdinality(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("ordinality", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.ordinality o CROSS JOIN unnest(ARRAY[1]) WITH ORDINALITY AS u",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "ordinality"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSeesPostfixInheritance covers `public.t*`.
func TestViewDependencyQualificationSeesPostfixInheritance(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.t* AS p, t AS q",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestRoutineReturnTypeAllowsWhitespaceAroundDot covers `RETURNS other . v`.
func TestRoutineReturnTypeAllowsWhitespaceAroundDot(t *testing.T) {
	snapshot := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: columnType})}},
			{Name: "other", Views: []*metadatapb.ViewMetadata{{
				Name:              "v",
				Definition:        "SELECT FROM",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
			}}},
		}}
	}
	snapshotWithFunction := func(columnType string) *metadatapb.DatabaseSchemaMetadata {
		db := snapshot(columnType)
		db.Schemas[0].Functions = []*metadatapb.FunctionMetadata{{
			Name:       "f",
			Signature:  "f()",
			Definition: "CREATE FUNCTION public.f() RETURNS other . v LANGUAGE sql AS $$ SELECT $$",
		}}
		return db
	}
	_, err := testMetadataMigration(snapshotWithFunction("integer"), snapshotWithFunction("bigint"))
	require.ErrorContains(t, err, `routine "public"."f" uses the row type of view "other"."v"`)
}

// TestReleasedRoutineKeysKeepQuotedCase covers f(public."V") vs f(public."v").
func TestReleasedRoutineKeysKeepQuotedCase(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:  "public",
		Views: []*metadatapb.ViewMetadata{{Name: "V", Definition: "SELECT 1"}},
		Functions: []*metadatapb.FunctionMetadata{
			{Name: "f", Signature: `f(public."V")`},
			{Name: "f", Signature: `f(public."v")`},
		},
	}}}
	sourceDrops := map[string]MigrationOp{
		`func:public.f|f(public."V")`: {Type: OpDropFunction, SchemaName: "public", ObjectName: "f", SQL: `DROP FUNCTION "public"."f"(public."V")`},
		`func:public.f|f(public."v")`: {Type: OpDropFunction, SchemaName: "public", ObjectName: "f", SQL: `DROP FUNCTION "public"."f"(public."v")`},
	}
	ops, _ := releaseRoutineOps(&MigrationPlan{}, oldMeta, map[string]bool{objectKey("public", "V"): true}, []string{"public"}, sourceDrops)
	require.Len(t, ops, 1)
	require.Equal(t, `DROP FUNCTION "public"."f"(public."V")`, ops[0].SQL)
	require.Contains(t, sourceDrops, `func:public.f|f(public."v")`)
}

// TestReleasedRoutineKeysKeepMultiwordBoundaries covers f(double precision) vs
// f(doubleprecision).
func TestReleasedRoutineKeysKeepMultiwordBoundaries(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:  "public",
		Views: []*metadatapb.ViewMetadata{{Name: "doubleprecision", Definition: "SELECT 1"}},
		Functions: []*metadatapb.FunctionMetadata{
			{Name: "f", Signature: "f(double precision)"},
			{Name: "f", Signature: "f(doubleprecision)"},
		},
	}}}
	sourceDrops := map[string]MigrationOp{
		`func:public.f|f(double precision)`: {Type: OpDropFunction, SchemaName: "public", ObjectName: "f", SQL: `DROP FUNCTION "public"."f"(double precision)`},
	}
	ops, moved := releaseRoutineOps(&MigrationPlan{}, oldMeta, map[string]bool{objectKey("public", "doubleprecision"): true}, []string{"public"}, sourceDrops)
	require.Empty(t, ops)
	require.Empty(t, moved)
	require.Contains(t, sourceDrops, `func:public.f|f(double precision)`)
}

// TestViewDependencyQualificationSkipsFunctionAlias covers a bare alias after a
// table function.
func TestViewDependencyQualificationSkipsFunctionAlias(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("g", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM public.g base CROSS JOIN unnest(ARRAY[1]) g",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "g"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSkipsJoinedTableAlias covers an alias after a
// parenthesized joined-table group.
func TestViewDependencyQualificationSkipsJoinedTableAlias(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT * FROM (public.t JOIN public.u ON true) t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationIgnoresJoinModifiers covers a dependency named
// cross beside CROSS JOIN.
func TestViewDependencyQualificationIgnoresJoinModifiers(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("cross", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        `SELECT * FROM public."cross" c CROSS JOIN public.u u ON true`,
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "cross"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationScansFunctionColumnTypes covers
// record_func() AS x(c t), where t is a column type.
func TestViewDependencyQualificationScansFunctionColumnTypes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT x.c, NULL::public.t FROM record_func() AS x(c t)",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

// TestViewDependencyQualificationSeesRegclassTypedConstant covers
// `regclass 't'`.
func TestViewDependencyQualificationSeesRegclassTypedConstant(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT regclass 't', 'public.t'::regclass FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT regclass 'public.t', 'public.t'::regclass FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationSpacesQualifiedRegclass covers whitespace
// around the dot in `::pg_catalog . regclass`.
func TestViewDependencyQualificationSpacesQualifiedRegclass(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT 't'::pg_catalog . regclass, 'public.t'::regclass FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.ErrorContains(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}), `view "public"."w" has an unqualified reference to "public"."t"`)

	meta.Schemas[0].Views[0].Definition = "SELECT 'public.t'::pg_catalog . regclass, 'public.t'::regclass FROM public.t"
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationIgnoresQualifiedRegclassType covers a
// user-defined type whose final component is named regclass.
func TestViewDependencyQualificationIgnoresQualifiedRegclassType(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT 't'::public.regclass FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestUnqualifiedRegclassRespectsSearchPath covers a user-defined regclass type
// before pg_catalog: the unqualified cast is not the built-in relation type.
func TestUnqualifiedRegclassRespectsSearchPath(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `types, "$user", public, pg_catalog`,
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "types",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
					Name:       "regclass",
					Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "x", Type: "text"}},
				}},
			},
			{
				Name:   "public",
				Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
				Views: []*metadatapb.ViewMetadata{{
					Name:              "w",
					Definition:        "SELECT 't'::regclass FROM public.t",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
				}},
			},
		},
	}
	paths := []string{"types", "$user", "public", "pg_catalog"}
	require.NoError(t, rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true}))
}

// TestViewDependencyQualificationScansRowsFromColumnTypes covers a column
// definition list inside ROWS FROM.
func TestViewDependencyQualificationSeesRegclassArrayElement(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT '{t}'::regclass[], 'public.t'::regclass FROM public.t",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

func TestViewDependencyQualificationScansRowsFromColumnTypes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		SearchPath: `"$user", public`,
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "w",
				Definition:        "SELECT x.c, NULL::public.t FROM ROWS FROM (record_func() AS (c t)) AS x",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
			}},
		}},
	}
	paths := []string{"$user", "public"}
	err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
	require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`)
}

func TestViewDependencyQualificationSeesParenthesizedRegclassCasts(t *testing.T) {
	paths := []string{"$user", "public"}
	for _, definition := range []string{
		"SELECT CAST(('t') AS regclass), 'public.t'::regclass FROM public.t",
		"SELECT (('t'))::regclass, 'public.t'::regclass FROM public.t",
		"SELECT regclass(('t')), 'public.t'::regclass FROM public.t",
		"SELECT ARRAY['t', 'public.t']::regclass[] FROM public.t",
		"SELECT '{t,public.t}'::regclass[] FROM public.t",
		"SELECT '[0:0]={t}'::regclass[], 'public.t'::regclass FROM public.t",
		"SELECT '[1:1][1:1]={{t}}'::regclass[], 'public.t'::regclass FROM public.t",
		`SELECT '{"\"t\""}'::regclass[], 'public.t'::regclass FROM public.t`,
		"SELECT 't' /*c*/ :: /*c*/ regclass, 'public.t'::regclass FROM public.t",
		"SELECT CAST(ARRAY['t'] AS regclass[]), 'public.t'::regclass FROM public.t",
	} {
		meta := &metadatapb.DatabaseSchemaMetadata{
			SearchPath: `"$user", public`,
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:   "public",
				Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
				Views: []*metadatapb.ViewMetadata{{
					Name:              "w",
					Definition:        definition,
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t"}},
				}},
			}},
		}
		err := rejectAmbiguousUserTypes(meta, meta, paths, paths, map[string]bool{"bbdataarchive": true})
		require.ErrorContains(t, err, `view "public"."w" has an unqualified reference to "public"."t"`, definition)
	}
}

func TestSkippedSchemaCommentsAreFiltered(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:    "bbdataarchive",
			Comment: "old comment",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Comment: "old table comment",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Comment: "old col comment"}},
			}},
		}},
	}
	newMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:    "bbdataarchive",
			Comment: "new comment",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t",
				Comment: "new table comment",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Comment: "new col comment"}},
			}},
		}},
	}
	sql, err := testMetadataMigration(oldMeta, newMeta)
	require.NoError(t, err)
	require.Empty(t, sql)
}

func TestChangedCollationIsRefused(t *testing.T) {
	withCollation := func(locale string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Collations: []*metadatapb.CollationMetadata{{
					Name:   "custom_col",
					Locale: locale,
				}},
			}},
		}
	}
	// Identical collations do not block unrelated migrations.
	sql, err := testMetadataMigration(withCollation("en_US.utf8"), withCollation("en_US.utf8"))
	require.NoError(t, err)
	require.Empty(t, sql)

	// Changed collations are refused rather than silently returning no-op.
	_, err = testMetadataMigration(withCollation("en_US.utf8"), withCollation("C.utf8"))
	require.ErrorContains(t, err, "collation")
}

func TestExternalTableChangesAreRefused(t *testing.T) {
	withExternalTable := func(option string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				ExternalTables: []*metadatapb.ExternalTableMetadata{{
					Name:    "ext_users",
					Options: []string{option},
				}},
			}},
		}
	}
	sql, err := testMetadataMigration(withExternalTable("table_name=users"), withExternalTable("table_name=users"))
	require.NoError(t, err)
	require.Empty(t, sql)

	_, err = testMetadataMigration(&metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}, withExternalTable("table_name=users"))
	require.ErrorContains(t, err, "external table")

	_, err = testMetadataMigration(withExternalTable("table_name=users"), withExternalTable("table_name=accounts"))
	require.ErrorContains(t, err, "external table")
}

func TestRoutineParameterDefaultsPreserved(t *testing.T) {
	empty := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	withFunc := func(def string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "greet",
					Signature:  "public.greet(name text DEFAULT 'world', times integer DEFAULT 1)",
					Definition: def,
				}},
			}},
		}
	}

	def1 := "CREATE FUNCTION public.greet(name text DEFAULT 'world', times integer DEFAULT 1) RETURNS text LANGUAGE sql AS $$ SELECT name $$"
	def2 := "CREATE FUNCTION public.greet(name text DEFAULT 'world', times integer DEFAULT 2) RETURNS text LANGUAGE sql AS $$ SELECT name $$"

	sql, err := testMetadataMigration(empty, withFunc(def1))
	require.NoError(t, err)
	require.Contains(t, sql, "DEFAULT 'world'")
	require.Contains(t, sql, "DEFAULT 1")

	sql, err = testMetadataMigration(withFunc(def1), withFunc(def2))
	require.NoError(t, err)
	require.Contains(t, sql, "CREATE OR REPLACE FUNCTION")
	require.Contains(t, sql, "DEFAULT 2")
}

func TestUnloadedSnapshotObjectsAreRefusedWhenChanged(t *testing.T) {
	rangeMeta := func(subtype string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:       "public",
				RangeTypes: []*metadatapb.RangeTypeMetadata{{Name: "num_range", Subtype: subtype}},
			}},
		}
	}
	sql, err := testMetadataMigration(rangeMeta("numeric"), rangeMeta("numeric"))
	require.NoError(t, err)
	require.Empty(t, sql)
	_, err = testMetadataMigration(rangeMeta("numeric"), rangeMeta("float8"))
	require.ErrorContains(t, err, "range types")

	serverMeta := func(host string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			ForeignServers: []*metadatapb.ForeignServerMetadata{{Name: "srv", Fdw: "postgres_fdw", Options: []string{"host=" + host}}},
		}
	}
	sql, err = testMetadataMigration(serverMeta("a"), serverMeta("a"))
	require.NoError(t, err)
	require.Empty(t, sql)
	_, err = testMetadataMigration(serverMeta("a"), serverMeta("b"))
	require.ErrorContains(t, err, "foreign servers")

	fdwMeta := func(handler string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{Name: "my_fdw", Handler: handler}},
		}
	}
	sql, err = testMetadataMigration(fdwMeta("h1"), fdwMeta("h1"))
	require.NoError(t, err)
	require.Empty(t, sql)
	_, err = testMetadataMigration(fdwMeta("h1"), fdwMeta("h2"))
	require.ErrorContains(t, err, "foreign data wrappers")
}

func TestTableOwnerOnlyChangeDoesNotProduceSpuriousDiff(t *testing.T) {
	withOwner := func(owner string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:   "public",
				Tables: []*metadatapb.TableMetadata{{Name: "t", Owner: owner, Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}}}},
			}},
		}
	}
	sql, err := testMetadataMigration(withOwner("alice"), withOwner("bob"))
	require.NoError(t, err)
	require.Empty(t, sql)
}

func TestFunctionSetClauseWithoutDefaultsPreserved(t *testing.T) {
	empty := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	withFunc := func(def string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "do_work",
					Signature:  "public.do_work(x integer)",
					Definition: def,
				}},
			}},
		}
	}

	def1 := "CREATE FUNCTION public.do_work(x integer) RETURNS integer LANGUAGE plpgsql SET search_path = public AS $$ BEGIN RETURN x; END; $$"
	def2 := "CREATE FUNCTION public.do_work(x integer) RETURNS integer LANGUAGE plpgsql SET search_path = archive AS $$ BEGIN RETURN x; END; $$"

	sql, err := testMetadataMigration(empty, withFunc(def1))
	require.NoError(t, err)
	require.Contains(t, sql, "SET search_path")

	sql, err = testMetadataMigration(withFunc(def1), withFunc(def2))
	require.NoError(t, err)
	require.Contains(t, sql, "CREATE OR REPLACE FUNCTION")
	require.Contains(t, sql, "SET search_path")
	require.Contains(t, sql, "archive")
}

func TestPartitionAttachAndDetachMigration(t *testing.T) {
	c0 := &metadatapb.TablePartitionMetadata{
		Name:       "c0",
		Type:       metadatapb.TablePartitionMetadata_RANGE,
		Expression: "RANGE (id)",
		Value:      "FOR VALUES FROM (0) TO (1)",
	}
	c1 := &metadatapb.TablePartitionMetadata{
		Name:       "c1",
		Type:       metadatapb.TablePartitionMetadata_RANGE,
		Expression: "RANGE (id)",
		Value:      "FOR VALUES FROM (1) TO (10)",
	}
	standalone := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{
				{
					Name:       "p",
					Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
					Partitions: []*metadatapb.TablePartitionMetadata{c0},
				},
				{
					Name:    "c1",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				},
			},
		}},
	}
	attached := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:       "p",
				Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				Partitions: []*metadatapb.TablePartitionMetadata{c0, c1},
			}},
		}},
	}

	sqlAttach, err := testMetadataMigration(standalone, attached)
	require.NoError(t, err)
	require.Contains(t, sqlAttach, `ALTER TABLE "public"."p" ATTACH PARTITION "public"."c1" FOR VALUES FROM ('1') TO ('10')`)
	require.NotContains(t, sqlAttach, `DROP TABLE`)

	sqlDetach, err := testMetadataMigration(attached, standalone)
	require.NoError(t, err)
	require.Contains(t, sqlDetach, `ALTER TABLE "public"."p" DETACH PARTITION "public"."c1"`)
	require.NotContains(t, sqlDetach, `DROP TABLE`)

	// Bound change under the same parent detaches and reattaches with the new bound.
	rebound := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "p",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				Partitions: []*metadatapb.TablePartitionMetadata{
					c0,
					{
						Name:       "c1",
						Type:       metadatapb.TablePartitionMetadata_RANGE,
						Expression: "RANGE (id)",
						Value:      "FOR VALUES FROM (10) TO (20)",
					},
				},
			}},
		}},
	}
	sqlRebound, err := testMetadataMigration(attached, rebound)
	require.NoError(t, err)
	require.Contains(t, sqlRebound, `ALTER TABLE "public"."p" DETACH PARTITION "public"."c1"`)
	require.Contains(t, sqlRebound, `ALTER TABLE "public"."p" ATTACH PARTITION "public"."c1" FOR VALUES FROM ('10') TO ('20')`)
	require.NotContains(t, sqlRebound, `DROP TABLE`)

	// Cross-schema partition attach is excluded when the child schema is skipped.
	crossStandalone := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name:       "p",
					Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
					Partitions: []*metadatapb.TablePartitionMetadata{c0},
				}},
			},
			{
				Name: "bbdataarchive",
				Tables: []*metadatapb.TableMetadata{{
					Name:    "c_skip",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				}},
			},
		},
	}
	crossAttached := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name:    "p",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
					Partitions: []*metadatapb.TablePartitionMetadata{
						c0,
						{
							Name:   "c_skip",
							Schema: "bbdataarchive",
							Type:   metadatapb.TablePartitionMetadata_RANGE,
							Value:  "FOR VALUES FROM (1) TO (10)",
						},
					},
				}},
			},
			{Name: "bbdataarchive"},
		},
	}
	sqlSkipped, err := MetadataMigration(context.Background(), crossStandalone, crossAttached, MetadataMigrationOptions{
		SkipSchemas: []string{"bbdataarchive"},
	})
	require.NoError(t, err)
	require.NotContains(t, sqlSkipped, `ATTACH PARTITION`)
	require.NotContains(t, sqlSkipped, `DETACH PARTITION`)
}

func TestFunctionRemovingDefaultDropsBeforeCreate(t *testing.T) {
	withDef := func(def string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "greet",
					Signature:  "public.greet(name text)",
					Definition: def,
				}},
			}},
		}
	}
	oldDef := "CREATE FUNCTION public.greet(name text DEFAULT 'world') RETURNS text LANGUAGE sql AS $$ SELECT name $$"
	newDef := "CREATE FUNCTION public.greet(name text) RETURNS text LANGUAGE sql AS $$ SELECT name $$"

	sql, err := testMetadataMigration(withDef(oldDef), withDef(newDef))
	require.NoError(t, err)
	require.Contains(t, sql, `DROP FUNCTION "public"."greet"(text)`)
	require.Contains(t, sql, `CREATE FUNCTION`)
	require.NotContains(t, sql, `CREATE OR REPLACE FUNCTION`)
}

func TestDottedCommentTargetsMigration(t *testing.T) {
	base := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "app.data",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "events.v1",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				Indexes: []*metadatapb.IndexMetadata{{
					Name:         "pk.v1",
					Type:         "btree",
					Expressions:  []string{"id"},
					Primary:      true,
					Unique:       true,
					IsConstraint: true,
				}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "fn.dot",
				Signature:  `"app.data"."fn.dot"(integer)`,
				Definition: `CREATE FUNCTION "app.data"."fn.dot"(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$`,
			}},
		}},
	}
	commented := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "app.data",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "events.v1",
				Comment: "tbl doc",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1, Comment: "col doc"}},
				Indexes: []*metadatapb.IndexMetadata{{
					Name:         "pk.v1",
					Type:         "btree",
					Expressions:  []string{"id"},
					Primary:      true,
					Unique:       true,
					IsConstraint: true,
					Comment:      "pk doc",
				}},
			}},
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "fn.dot",
				Signature:  `"app.data"."fn.dot"(integer)`,
				Definition: `CREATE FUNCTION "app.data"."fn.dot"(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$`,
				Comment:    "fn doc",
			}},
		}},
	}

	sql, err := testMetadataMigration(base, commented)
	require.NoError(t, err)
	require.Contains(t, sql, `COMMENT ON TABLE "app.data"."events.v1" IS 'tbl doc'`)
	require.Contains(t, sql, `COMMENT ON COLUMN "app.data"."events.v1"."id" IS 'col doc'`)
	require.Contains(t, sql, `COMMENT ON CONSTRAINT "pk.v1" ON "app.data"."events.v1" IS 'pk doc'`)
	require.Contains(t, sql, `COMMENT ON FUNCTION "app.data"."fn.dot"(integer) IS 'fn doc'`)
}

func TestListPartitionNullBoundMigration(t *testing.T) {
	empty := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}
	withNullList := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "events",
				Columns: []*metadatapb.ColumnMetadata{{Name: "code", Type: "integer", Position: 1}},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:       "events_null",
					Type:       metadatapb.TablePartitionMetadata_LIST,
					Expression: "LIST (code)",
					Value:      "FOR VALUES IN (NULL, 1)",
				}},
			}},
		}},
	}

	sql, err := testMetadataMigration(empty, withNullList)
	require.NoError(t, err)
	require.Contains(t, sql, `FOR VALUES IN (NULL, '1')`)
	require.NotContains(t, sql, `'NULL'`)

	c := New()
	_, err = c.Exec(sql, nil)
	require.NoError(t, err)
}

func TestRound10CommentAndRangeMigrations(t *testing.T) {
	// RangeType comment-only change should not be refused and should emit COMMENT ON TYPE.
	rangeSnapshot := func(comment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				RangeTypes: []*metadatapb.RangeTypeMetadata{{
					Name:    "num_range",
					Subtype: "numeric",
					Comment: comment,
				}},
			}},
		}
	}
	sql, err := testMetadataMigration(rangeSnapshot(""), rangeSnapshot("range doc"))
	require.NoError(t, err)
	require.Contains(t, sql, `COMMENT ON TYPE "public"."num_range" IS 'range doc'`)

	sql, err = testMetadataMigration(rangeSnapshot("range doc"), rangeSnapshot(""))
	require.NoError(t, err)
	require.Contains(t, sql, `COMMENT ON TYPE "public"."num_range" IS NULL`)

	// FK, CHECK, and TRIGGER comment-only changes should emit COMMENT statements.
	commentedSnapshot := func(fkComment, chkComment, partChkComment, trigComment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "parent",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer", Position: 1},
						},
						Indexes: []*metadatapb.IndexMetadata{{
							Name:         "parent_pkey",
							Type:         "btree",
							Expressions:  []string{"id"},
							Primary:      true,
							Unique:       true,
							IsConstraint: true,
						}},
					},
					{
						Name: "child",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer", Position: 1},
							{Name: "parent_id", Type: "integer", Position: 2},
						},
						ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
							Name:              "fk_child_parent",
							Columns:           []string{"parent_id"},
							ReferencedSchema:  "public",
							ReferencedTable:   "parent",
							ReferencedColumns: []string{"id"},
							Comment:           fkComment,
						}},
						CheckConstraints: []*metadatapb.CheckConstraintMetadata{{
							Name:       "chk_pos",
							Expression: "(id > 0)",
							Comment:    chkComment,
						}},
						Triggers: []*metadatapb.TriggerMetadata{{
							Name:    "trg_child",
							Body:    `CREATE TRIGGER trg_child BEFORE INSERT ON public.child FOR EACH ROW EXECUTE FUNCTION public.touch()`,
							Comment: trigComment,
						}},
					},
					{
						Name: "logs",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer", Position: 1},
						},
						Partitions: []*metadatapb.TablePartitionMetadata{{
							Name:       "logs_p0",
							Type:       metadatapb.TablePartitionMetadata_RANGE,
							Expression: "RANGE (id)",
							Value:      "FOR VALUES FROM (0) TO (10)",
							CheckConstraints: []*metadatapb.CheckConstraintMetadata{{
								Name:       "chk_part_pos",
								Expression: "(id >= 0)",
								Comment:    partChkComment,
							}},
						}},
					},
				},
			}},
		}
	}

	sql, err = testMetadataMigration(commentedSnapshot("", "", "", ""), commentedSnapshot("fk doc", "chk doc", "pchk doc", "trig doc"))
	require.NoError(t, err)
	require.Contains(t, sql, `COMMENT ON CONSTRAINT "fk_child_parent" ON "public"."child" IS 'fk doc'`)
	require.Contains(t, sql, `COMMENT ON CONSTRAINT "chk_pos" ON "public"."child" IS 'chk doc'`)
	require.Contains(t, sql, `COMMENT ON CONSTRAINT "chk_part_pos" ON "public"."logs_p0" IS 'pchk doc'`)
	require.Contains(t, sql, `COMMENT ON TRIGGER "trg_child" ON "public"."child" IS 'trig doc'`)
}

func TestDottedSkipSchemaIgnoresDegradedObjects(t *testing.T) {
	makeMeta := func(viewDef string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "app",
					Tables: []*metadatapb.TableMetadata{{
						Name:    "t",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
					}},
				},
				{
					Name: "app.data",
					Views: []*metadatapb.ViewMetadata{{
						Name:       "broken",
						Definition: viewDef,
					}},
				},
			},
		}
	}

	// Unparseable view in skipped dotted schema must not cause unloadableChanges failure.
	sql, err := MetadataMigration(context.Background(), makeMeta("SELECT FROM"), makeMeta("SELECT +"), MetadataMigrationOptions{
		SkipSchemas: []string{"app.data"},
	})
	require.NoError(t, err)
	require.Empty(t, sql)
}

func TestRound12CollationExtTriggerAndViewColumnComments(t *testing.T) {
	// 1 & 2. Comment-only changes on user collations and external-table triggers.
	commentOnlyMeta := func(collComment, trigComment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Collations: []*metadatapb.CollationMetadata{{
					Name:     "c_custom",
					Provider: "c",
					Comment:  collComment,
				}},
				ExternalTables: []*metadatapb.ExternalTableMetadata{{
					Name:               "ext_events",
					ExternalServerName: "srv",
					Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
					Triggers: []*metadatapb.TriggerMetadata{{
						Name:    "trg_ext",
						Timing:  "BEFORE",
						Event:   "INSERT",
						Body:    "CREATE TRIGGER trg_ext BEFORE INSERT ON public.ext_events FOR EACH ROW EXECUTE FUNCTION public.fn()",
						Comment: trigComment,
					}},
				}},
			}},
		}
	}

	sql, err := testMetadataMigration(commentOnlyMeta("", ""), commentOnlyMeta("coll doc", "ext trig doc"))
	require.NoError(t, err)
	require.Contains(t, sql, `COMMENT ON COLLATION "public"."c_custom" IS 'coll doc'`)
	require.Contains(t, sql, `COMMENT ON TRIGGER "trg_ext" ON "public"."ext_events" IS 'ext trig doc'`)

	// 8. Recreated view restores column comments as well as view comment.
	viewMeta := func(colType string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: colType}),
				},
				Views: []*metadatapb.ViewMetadata{{
					Name:       "v",
					Definition: "SELECT id FROM public.t",
					Comment:    "view doc",
					DependencyColumns: []*metadatapb.DependencyColumn{{
						Schema: "public", Table: "t", Column: "id",
					}},
					Columns: []*metadatapb.ColumnMetadata{{
						Name:    "id",
						Type:    colType,
						Comment: "view col doc",
					}},
				}},
			}},
		}
	}

	viewSQL, err := testMetadataMigration(viewMeta("integer"), viewMeta("bigint"))
	require.NoError(t, err)
	require.Contains(t, viewSQL, `COMMENT ON VIEW "public"."v" IS 'view doc'`)
	require.Contains(t, viewSQL, `COMMENT ON COLUMN "public"."v"."id" IS 'view col doc'`)
}

func TestRound14Fixes(t *testing.T) {
	// 1. Relation triggers on tables, views, and materialized views migrate structurally.
	trigBaseFn := &metadatapb.FunctionMetadata{
		Name:       "audit_fn",
		Signature:  "audit_fn()",
		Definition: "CREATE FUNCTION public.audit_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$",
	}
	trigMeta := func(hasTableTrig, hasViewTrig bool, matTrigBody string) *metadatapb.DatabaseSchemaMetadata {
		tbl := pgTable("events", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})
		if hasTableTrig {
			tbl.Triggers = []*metadatapb.TriggerMetadata{{
				Name:   "trg_events",
				Timing: "BEFORE",
				Event:  "INSERT",
				Body:   `CREATE TRIGGER trg_events BEFORE INSERT ON public.events FOR EACH ROW EXECUTE FUNCTION public.audit_fn()`,
			}}
		}
		v := &metadatapb.ViewMetadata{
			Name:              "v_events",
			Definition:        "SELECT id FROM public.events",
			DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "events", Column: "id"}},
		}
		if hasViewTrig {
			v.Triggers = []*metadatapb.TriggerMetadata{{
				Name:   "trg_view",
				Timing: "INSTEAD OF",
				Event:  "INSERT",
				Body:   `CREATE TRIGGER trg_view INSTEAD OF INSERT ON public.v_events FOR EACH ROW EXECUTE FUNCTION public.audit_fn('tenant', 'ops')`,
			}}
		}
		mv := &metadatapb.MaterializedViewMetadata{
			Name:              "mv_events",
			Definition:        "SELECT id FROM public.events",
			DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "events", Column: "id"}},
		}
		if matTrigBody != "" {
			mv.Triggers = []*metadatapb.TriggerMetadata{{
				Name:   "trg_mv",
				Timing: "BEFORE",
				Event:  "TRUNCATE",
				Body:   matTrigBody,
			}}
		}
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:              "public",
				Functions:         []*metadatapb.FunctionMetadata{trigBaseFn},
				Tables:            []*metadatapb.TableMetadata{tbl},
				Views:             []*metadatapb.ViewMetadata{v},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{mv},
			}},
		}
	}
	trigSQL, err := testMetadataMigration(
		trigMeta(true, false, `CREATE TRIGGER trg_mv BEFORE TRUNCATE ON public.mv_events FOR EACH STATEMENT EXECUTE FUNCTION public.audit_fn()`),
		trigMeta(false, true, `CREATE TRIGGER trg_mv AFTER TRUNCATE ON public.mv_events FOR EACH STATEMENT EXECUTE FUNCTION public.audit_fn()`),
	)
	require.NoError(t, err)
	require.Contains(t, trigSQL, `DROP TRIGGER "trg_events" ON "public"."events"`)
	require.Contains(t, trigSQL, `CREATE TRIGGER "trg_view" INSTEAD OF INSERT ON "public"."v_events"`)
	require.Contains(t, trigSQL, `EXECUTE FUNCTION "audit_fn"('tenant', 'ops')`)
	require.Contains(t, trigSQL, `DROP TRIGGER "trg_mv" ON "public"."mv_events"`)
	require.Contains(t, trigSQL, `CREATE TRIGGER "trg_mv" AFTER TRUNCATE ON "public"."mv_events"`)

	// Unloadable trigger body changes are refused rather than silently skipped.
	unloadTrigMeta := func(body string) *metadatapb.DatabaseSchemaMetadata {
		t1 := pgTable("events", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})
		t1.Triggers = []*metadatapb.TriggerMetadata{{
			Name:   "trg_missing_fn",
			Timing: "BEFORE",
			Event:  "INSERT",
			Body:   body,
		}}
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:   "public",
				Tables: []*metadatapb.TableMetadata{t1},
			}},
		}
	}
	_, err = testMetadataMigration(
		unloadTrigMeta(`CREATE TRIGGER trg_missing_fn BEFORE INSERT ON public.events FOR EACH ROW EXECUTE FUNCTION public.missing_fn()`),
		unloadTrigMeta(`CREATE TRIGGER trg_missing_fn AFTER INSERT ON public.events FOR EACH ROW EXECUTE FUNCTION public.missing_fn()`),
	)
	require.ErrorContains(t, err, "trg:public.events.trg_missing_fn")

	// 2. Rewrite rules on tables and views migrate cleanly (ignoring view _RETURN rules).
	ruleMeta := func(rules []*metadatapb.RuleMetadata, viewRules []*metadatapb.RuleMetadata) *metadatapb.DatabaseSchemaMetadata {
		tbl := pgTable("items", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})
		tbl.Rules = rules
		v := &metadatapb.ViewMetadata{
			Name:       "v_items",
			Definition: "SELECT id FROM public.items",
			Rules: append([]*metadatapb.RuleMetadata{{
				Name: "_RETURN", Event: "SELECT", IsInstead: true, IsEnabled: true,
				Definition: "CREATE RULE \"_RETURN\" AS ON SELECT TO public.v_items DO INSTEAD SELECT id FROM public.items;",
			}}, viewRules...),
		}
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:   "public",
				Tables: []*metadatapb.TableMetadata{tbl},
				Views:  []*metadatapb.ViewMetadata{v},
			}},
		}
	}
	ruleSQL, err := testMetadataMigration(
		ruleMeta(
			[]*metadatapb.RuleMetadata{{Name: "r_drop", Event: "INSERT", IsInstead: true, IsEnabled: true, Definition: "CREATE RULE r_drop AS ON INSERT TO public.items DO INSTEAD NOTHING;"}},
			nil,
		),
		ruleMeta(
			[]*metadatapb.RuleMetadata{{Name: "r_add", Event: "DELETE", IsInstead: true, IsEnabled: false, Definition: "CREATE RULE r_add AS ON DELETE TO public.items DO INSTEAD NOTHING;"}},
			[]*metadatapb.RuleMetadata{{Name: "vr_ins", Event: "INSERT", IsInstead: true, IsEnabled: true, Definition: "CREATE RULE vr_ins AS ON INSERT TO public.v_items DO INSTEAD NOTHING;"}},
		),
	)
	require.NoError(t, err)
	require.Contains(t, ruleSQL, `DROP RULE "r_drop" ON "public"."items"`)
	require.Contains(t, ruleSQL, `CREATE RULE "r_add" AS ON DELETE TO "public"."items" DO INSTEAD NOTHING`)
	require.Contains(t, ruleSQL, `ALTER TABLE "public"."items" DISABLE RULE "r_add"`)
	require.Contains(t, ruleSQL, `CREATE RULE "vr_ins" AS ON INSERT TO "public"."v_items" DO INSTEAD NOTHING`)

	// 3. Event triggers migrate cleanly.
	etFn := &metadatapb.FunctionMetadata{
		Name:       "et_fn",
		Signature:  "et_fn()",
		Definition: "CREATE FUNCTION public.et_fn() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN END; $$",
	}
	etOld := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			Functions: []*metadatapb.FunctionMetadata{etFn},
		}},
		EventTriggers: []*metadatapb.EventTriggerMetadata{{
			Name:           "et_drop",
			Event:          "ddl_command_start",
			FunctionSchema: "public",
			FunctionName:   "et_fn",
			Enabled:        true,
		}},
	}
	etNew := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			Functions: []*metadatapb.FunctionMetadata{etFn},
		}},
		EventTriggers: []*metadatapb.EventTriggerMetadata{{
			Name:           "et_add",
			Event:          "sql_drop",
			Tags:           []string{"DROP TABLE"},
			FunctionSchema: "public",
			FunctionName:   "et_fn",
			Enabled:        false,
		}},
	}
	etSQL, err := testMetadataMigration(etOld, etNew)
	require.NoError(t, err)
	require.Contains(t, etSQL, `DROP EVENT TRIGGER "et_drop"`)
	require.Contains(t, etSQL, `CREATE EVENT TRIGGER "et_add" ON "sql_drop"`)
	require.Contains(t, etSQL, `ALTER EVENT TRIGGER "et_add" DISABLE`)

	// 4. External-table column comment migrations succeed while structural changes stay refused.
	extMeta := func(colType, colComment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				ExternalTables: []*metadatapb.ExternalTableMetadata{{
					Name:               "ext_events",
					ExternalServerName: "srv",
					Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: colType, Comment: colComment}},
				}},
			}},
		}
	}
	extCommentSQL, err := testMetadataMigration(extMeta("integer", ""), extMeta("integer", "ext col doc"))
	require.NoError(t, err)
	require.Contains(t, extCommentSQL, `COMMENT ON COLUMN "public"."ext_events"."id" IS 'ext col doc'`)
	_, err = testMetadataMigration(extMeta("integer", ""), extMeta("bigint", ""))
	require.ErrorContains(t, err, "external tables")

	// 5. SkipSchemas filters out schema create, drop, and alter owner operations.
	skipSQL, err := MetadataMigration(
		context.Background(),
		&metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{Name: "public"},
				{Name: "backup_drop", Owner: "alice"},
				{Name: "backup_alter", Owner: "alice"},
			},
		},
		&metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{Name: "public"},
				{Name: "backup_add", Owner: "alice"},
				{Name: "backup_alter", Owner: "bob"},
			},
		},
		MetadataMigrationOptions{SkipSchemas: []string{"backup_drop", "backup_alter", "backup_add"}},
	)
	require.NoError(t, err)
	require.Empty(t, skipSQL)

	// 6. Source-only missing drops inside a dropped schema are suppressed behind DROP SCHEMA CASCADE.
	oldDropSchemaMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public"},
			{
				Name: "legacy",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "broken_fn",
					Signature:  "broken_fn(nonexistent_type)",
					Definition: "broken sql here",
				}},
				Tables: []*metadatapb.TableMetadata{{
					Name:    "broken_tbl",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					Indexes: []*metadatapb.IndexMetadata{{
						Name:       "broken_idx",
						Definition: "CREATE INDEX broken_idx ON legacy.broken_tbl USING btree ((id + 1)) WHERE unknown_fn(id)",
					}},
				}},
			},
		},
	}
	newDropSchemaMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}
	dropSchemaSQL, err := testMetadataMigration(oldDropSchemaMeta, newDropSchemaMeta)
	require.NoError(t, err)
	require.Contains(t, dropSchemaSQL, `DROP SCHEMA legacy CASCADE`)
	require.NotContains(t, dropSchemaSQL, `DROP FUNCTION "legacy"."broken_fn"`)
	require.NotContains(t, dropSchemaSQL, `DROP INDEX "legacy"."broken_idx"`)
}

func TestRound16Fixes(t *testing.T) {
	// 1. Foreign-table check constraint comment-only changes emit COMMENT ON CONSTRAINT; structural changes remain refused.
	extChkMeta := func(expr, chkComment string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				ExternalTables: []*metadatapb.ExternalTableMetadata{{
					Name:               "ext_events",
					ExternalServerName: "srv",
					Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					CheckConstraints:   []*metadatapb.CheckConstraintMetadata{{Name: "chk_id", Expression: expr, Comment: chkComment}},
				}},
			}},
		}
	}
	extChkSQL, err := testMetadataMigration(extChkMeta("id > 0", ""), extChkMeta("id > 0", "ext chk doc"))
	require.NoError(t, err)
	require.Contains(t, extChkSQL, `COMMENT ON CONSTRAINT "chk_id" ON "public"."ext_events" IS 'ext chk doc'`)
	_, err = testMetadataMigration(extChkMeta("id > 0", ""), extChkMeta("id > 1", ""))
	require.ErrorContains(t, err, "external tables")

	// 2. Routine SUPPORT and two-part C AS 'object_file', 'link_symbol' round-trip cleanly.
	oldFnMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}
	newFnMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "c_add",
					Signature:  "c_add(integer)",
					Definition: `CREATE FUNCTION public.c_add(integer) RETURNS integer LANGUAGE c SUPPORT public.sup_fn AS '$libdir/mylib', 'c_add_entry'`,
				},
				{
					Name:       "c_entry",
					Signature:  "c_entry(integer)",
					Definition: `CREATE FUNCTION public.c_entry(integer) RETURNS integer LANGUAGE c AS '$libdir/libfoo', 'foo_entry'`,
				},
			},
		}},
	}
	fnSQL, err := testMetadataMigration(oldFnMeta, newFnMeta)
	require.NoError(t, err)
	require.Contains(t, fnSQL, `SUPPORT public.sup_fn`)
	require.Contains(t, fnSQL, `AS '$libdir/mylib', 'c_add_entry'`)
	require.Contains(t, fnSQL, `AS '$libdir/libfoo', 'foo_entry'`)
}
