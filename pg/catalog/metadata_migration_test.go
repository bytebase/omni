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
	require.Equal(t, `DROP INDEX IF EXISTS "public"."t_missing_idx"`, sql)

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

// TestSourceOnlyProcedureListedAsFunction covers a procedure sync stored among
// functions: its synthesized drop must be DROP PROCEDURE.
func TestSourceOnlyProcedureListedAsFunction(t *testing.T) {
	def := "CREATE OR REPLACE PROCEDURE public.p(a missing_type)\n LANGUAGE sql\nAS $procedure$ SELECT 1 $procedure$"
	sig := "p(a missing_type)"
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name:      "public",
		Functions: []*metadatapb.FunctionMetadata{{Name: "p", Signature: sig, Definition: def}},
	}}}
	newMeta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	key := "func:public.p|" + routineIdentity(sig, def)
	drops := sourceOnlyMissingDrops(oldMeta, newMeta, &LoadMetadataReport{Missing: map[string]error{key: nil}})
	drop, ok := drops[key]
	require.True(t, ok, drops)
	require.True(t, strings.HasPrefix(drop.SQL, "DROP PROCEDURE "), drop.SQL)
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
	sql, err := testMetadataMigration(withOwner(""), withOwner("postgres"))
	require.NoError(t, err)
	require.Empty(t, sql)

	sql, err = testMetadataMigration(withOwner("alice"), withOwner("bob"))
	require.NoError(t, err)
	require.Contains(t, sql, `ALTER TABLE "public"."t" OWNER TO bob`)
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
	require.Contains(t, trigSQL, `EXECUTE FUNCTION "public"."audit_fn"('tenant', 'ops')`)
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
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "sup_fn",
				Signature:  "sup_fn(internal)",
				Definition: `CREATE FUNCTION public.sup_fn(internal) RETURNS internal LANGUAGE internal STRICT AS 'sup_fn'`,
			}},
		}},
	}
	newFnMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "sup_fn",
					Signature:  "sup_fn(internal)",
					Definition: `CREATE FUNCTION public.sup_fn(internal) RETURNS internal LANGUAGE internal STRICT AS 'sup_fn'`,
				},
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

func TestRound17Fixes(t *testing.T) {
	// 1. Routine parameter name, mode, and multiple OUT param type changes emit DROP + CREATE.
	paramFnMeta := func(def string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "pfn",
					Signature:  "pfn(integer)",
					Definition: def,
				}},
			}},
		}
	}
	renameParamSQL, err := testMetadataMigration(
		paramFnMeta(`CREATE FUNCTION public.pfn(a integer, OUT x integer, OUT y text) RETURNS record LANGUAGE sql AS $$ SELECT a, 'x' $$`),
		paramFnMeta(`CREATE FUNCTION public.pfn(b integer, OUT x integer, OUT y text) RETURNS record LANGUAGE sql AS $$ SELECT b, 'x' $$`),
	)
	require.NoError(t, err)
	require.Contains(t, renameParamSQL, `DROP FUNCTION "public"."pfn"(integer)`)
	require.Contains(t, renameParamSQL, `CREATE FUNCTION "public"."pfn"("b" integer, OUT "x" integer, OUT "y" text)`)

	outTypeSQL, err := testMetadataMigration(
		paramFnMeta(`CREATE FUNCTION public.pfn(a integer, OUT x integer, OUT y text) RETURNS record LANGUAGE sql AS $$ SELECT a, 'x' $$`),
		paramFnMeta(`CREATE FUNCTION public.pfn(a integer, OUT x integer, OUT y bigint) RETURNS record LANGUAGE sql AS $$ SELECT a, 1::bigint $$`),
	)
	require.NoError(t, err)
	require.Contains(t, outTypeSQL, `DROP FUNCTION "public"."pfn"(integer)`)
	require.Contains(t, outTypeSQL, `OUT "y" bigint`)

	// 2. Trigger calling cross-schema function preserves qualified function schema.
	crossTrigMeta := func(withTrig bool) *metadatapb.DatabaseSchemaMetadata {
		tbl := pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})
		if withTrig {
			tbl.Triggers = []*metadatapb.TriggerMetadata{{
				Name:   "trg_app",
				Timing: "BEFORE",
				Event:  "INSERT",
				Body:   `CREATE TRIGGER trg_app BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION app.audit()`,
			}}
		}
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{Name: "public", Tables: []*metadatapb.TableMetadata{tbl}},
				{
					Name: "app",
					Functions: []*metadatapb.FunctionMetadata{{
						Name:       "audit",
						Signature:  "audit()",
						Definition: `CREATE FUNCTION app.audit() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$`,
					}},
				},
			},
		}
	}
	crossTrigSQL, err := MetadataMigration(
		context.Background(),
		crossTrigMeta(false),
		crossTrigMeta(true),
		MetadataMigrationOptions{
			SourceSearchPath: []string{"public"},
			TargetSearchPath: []string{"public"},
		},
	)
	require.NoError(t, err)
	require.Contains(t, crossTrigSQL, `EXECUTE FUNCTION "app"."audit"()`)

	// 3. Table ownership migrations: owner on added table and owner-only change on existing table.
	ownerMeta := func(owner string, include bool) *metadatapb.DatabaseSchemaMetadata {
		sm := &metadatapb.SchemaMetadata{Name: "public"}
		if include {
			tbl := pgTable("owned_tbl", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})
			tbl.Owner = owner
			sm.Tables = []*metadatapb.TableMetadata{tbl}
		}
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{sm}}
	}
	addOwnerSQL, err := testMetadataMigration(ownerMeta("", false), ownerMeta("alice", true))
	require.NoError(t, err)
	require.Contains(t, addOwnerSQL, `CREATE TABLE "public"."owned_tbl"`)
	require.Contains(t, addOwnerSQL, `ALTER TABLE "public"."owned_tbl" OWNER TO alice`)

	chgOwnerSQL, err := testMetadataMigration(ownerMeta("alice", true), ownerMeta("bob", true))
	require.NoError(t, err)
	require.Contains(t, chgOwnerSQL, `ALTER TABLE "public"."owned_tbl" OWNER TO bob`)
}

func TestRound18Fixes(t *testing.T) {
	// 1. Replacing a void function with a procedure drops the function and creates the procedure.
	fnMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "do_work",
				Signature:  "do_work(integer)",
				Definition: "CREATE FUNCTION public.do_work(x integer) RETURNS void LANGUAGE plpgsql AS $$ BEGIN END; $$",
			}},
		}},
	}
	procMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "do_work",
				Signature:  "do_work(integer)",
				Definition: "CREATE PROCEDURE public.do_work(x integer) LANGUAGE plpgsql AS $$ BEGIN END; $$",
			}},
		}},
	}
	kindSQL, err := testMetadataMigration(fnMeta, procMeta)
	require.NoError(t, err)
	require.Contains(t, kindSQL, `DROP FUNCTION "public"."do_work"(integer)`)
	require.Contains(t, kindSQL, `CREATE PROCEDURE "public"."do_work"("x" integer)`)

	// 2. Database ownership migrations when nonempty database names match.
	dbOwnerOnlySQL, err := testMetadataMigration(
		&metadatapb.DatabaseSchemaMetadata{Name: "appdb", Owner: "alice"},
		&metadatapb.DatabaseSchemaMetadata{Name: "appdb", Owner: "bob"},
	)
	require.NoError(t, err)
	require.Equal(t, `ALTER DATABASE appdb OWNER TO bob`, dbOwnerOnlySQL)

	dbOwnerWithTableSQL, err := testMetadataMigration(
		&metadatapb.DatabaseSchemaMetadata{
			Name:    "appdb",
			Owner:   "alice",
			Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
		},
		&metadatapb.DatabaseSchemaMetadata{
			Name:  "appdb",
			Owner: "bob",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name:   "public",
				Tables: []*metadatapb.TableMetadata{pgTable("t", 0, &metadatapb.ColumnMetadata{Name: "id", Type: "integer"})},
			}},
		},
	)
	require.NoError(t, err)
	require.Contains(t, dbOwnerWithTableSQL, `CREATE TABLE "public"."t"`)
	require.Contains(t, dbOwnerWithTableSQL, `ALTER DATABASE appdb OWNER TO bob`)

	noDBOwnerSQL, err := testMetadataMigration(
		&metadatapb.DatabaseSchemaMetadata{Name: "db_a", Owner: "alice"},
		&metadatapb.DatabaseSchemaMetadata{Name: "db_b", Owner: "bob"},
	)
	require.NoError(t, err)
	require.Empty(t, noDBOwnerSQL)
}

func TestRound19Fixes(t *testing.T) {
	// 1. Standalone table with matching index attached as partition suppresses duplicate ALTER INDEX ATTACH PARTITION.
	p0 := &metadatapb.TablePartitionMetadata{
		Name:       "events_p0",
		Expression: "RANGE (id)",
		Value:      "FOR VALUES FROM (-100) TO (0)",
	}
	standaloneMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{
				{
					Name:       "events",
					Columns:    []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					Indexes:    []*metadatapb.IndexMetadata{{Name: "events_id_idx", Expressions: []string{"id"}}},
					Partitions: []*metadatapb.TablePartitionMetadata{p0},
				},
				{
					Name:    "events_p1",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					Indexes: []*metadatapb.IndexMetadata{{Name: "events_p1_id_idx", Expressions: []string{"id"}}},
				},
			},
		}},
	}
	attachedMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "events",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				Indexes: []*metadatapb.IndexMetadata{{Name: "events_id_idx", Expressions: []string{"id"}}},
				Partitions: []*metadatapb.TablePartitionMetadata{
					p0,
					{
						Name:       "events_p1",
						Expression: "RANGE (id)",
						Value:      "FOR VALUES FROM (0) TO (100)",
						Indexes:    []*metadatapb.IndexMetadata{{Name: "events_p1_id_idx", Expressions: []string{"id"}}},
					},
				},
			}},
		}},
	}
	attachSQL, err := testMetadataMigration(standaloneMeta, attachedMeta)
	require.NoError(t, err)
	require.Contains(t, attachSQL, `ALTER TABLE "public"."events" ATTACH PARTITION "public"."events_p1" FOR VALUES FROM ('0') TO ('100')`)
	require.NotContains(t, attachSQL, `ATTACH PARTITION "public"."events_p1_id_idx"`)

	// 2. Quoted 'NULL' partition bound value is preserved as string literal rather than SQL NULL.
	pSeed := &metadatapb.TablePartitionMetadata{
		Name:       "items_seed",
		Expression: "LIST (tag)",
		Value:      "FOR VALUES IN ('seed')",
	}
	listParentMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:       "items",
				Columns:    []*metadatapb.ColumnMetadata{{Name: "tag", Type: "text"}},
				Partitions: []*metadatapb.TablePartitionMetadata{pSeed},
			}},
		}},
	}
	listChildMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "items",
				Columns: []*metadatapb.ColumnMetadata{{Name: "tag", Type: "text"}},
				Partitions: []*metadatapb.TablePartitionMetadata{
					pSeed,
					{
						Name:       "items_str_null",
						Expression: "LIST (tag)",
						Value:      "FOR VALUES IN ('NULL')",
					},
					{
						Name:       "items_sql_null",
						Expression: "LIST (tag)",
						Value:      "FOR VALUES IN (NULL)",
					},
				},
			}},
		}},
	}
	listPartSQL, err := testMetadataMigration(listParentMeta, listChildMeta)
	require.NoError(t, err)
	require.Contains(t, listPartSQL, `CREATE TABLE "public"."items_str_null" PARTITION OF "public"."items" FOR VALUES IN ('NULL')`)
	require.Contains(t, listPartSQL, `CREATE TABLE "public"."items_sql_null" PARTITION OF "public"."items" FOR VALUES IN (NULL)`)
}

func TestRound20Fixes(t *testing.T) {
	// 1. Supported function added alongside its alphabetically-later support routine orders support first.
	supportFrom := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}
	supportTo := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "aaa_caller",
					Signature:  "aaa_caller(integer)",
					Definition: `CREATE FUNCTION public.aaa_caller(integer) RETURNS integer LANGUAGE sql SUPPORT public.zzz_sup AS 'SELECT $1'`,
				},
				{
					Name:       "zzz_sup",
					Signature:  "zzz_sup(internal)",
					Definition: `CREATE FUNCTION public.zzz_sup(internal) RETURNS internal LANGUAGE internal AS 'sup'`,
				},
			},
		}},
	}
	supSQL, err := testMetadataMigration(supportFrom, supportTo)
	require.NoError(t, err)
	supIdx := strings.Index(supSQL, `CREATE FUNCTION "public"."zzz_sup"`)
	callerIdx := strings.Index(supSQL, `CREATE FUNCTION public.aaa_caller`)
	require.NotEqual(t, -1, supIdx, supSQL)
	require.NotEqual(t, -1, callerIdx, supSQL)
	require.Less(t, supIdx, callerIdx, "support function must be created before caller")

	// 2. Switching between ordinary function and WINDOW function drops and recreates.
	winFrom := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "wf",
				Signature:  "wf(integer)",
				Definition: `CREATE FUNCTION public.wf(integer) RETURNS integer LANGUAGE internal AS 'wf'`,
			}},
		}},
	}
	winTo := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "wf",
				Signature:  "wf(integer)",
				Definition: `CREATE FUNCTION public.wf(integer) RETURNS integer LANGUAGE internal WINDOW AS 'wf'`,
			}},
		}},
	}
	winSQL, err := testMetadataMigration(winFrom, winTo)
	require.NoError(t, err)
	require.Contains(t, winSQL, `DROP FUNCTION "public"."wf"(integer)`)
	require.Contains(t, winSQL, `WINDOW`)
	require.NotContains(t, winSQL, `CREATE OR REPLACE`)

	// 3. Routine comments quote mixed-case, reserved, and hyphenated identifiers.
	commentFrom := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}
	commentTo := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "Audit",
					Signature:  "Audit(integer)",
					Definition: `CREATE FUNCTION public."Audit"(integer) RETURNS integer LANGUAGE sql AS 'SELECT $1'`,
					Comment:    "mixed",
				},
				{
					Name:       "select",
					Signature:  "select(integer)",
					Definition: `CREATE FUNCTION public."select"(integer) RETURNS integer LANGUAGE sql AS 'SELECT $1'`,
					Comment:    "keyword",
				},
			},
			Procedures: []*metadatapb.ProcedureMetadata{{
				Name:       "my-func",
				Signature:  "my-func(integer)",
				Definition: `CREATE PROCEDURE public."my-func"(integer) LANGUAGE plpgsql AS $$BEGIN END$$`,
				Comment:    "hyphen",
			}},
		}},
	}
	commentSQL, err := testMetadataMigration(commentFrom, commentTo)
	require.NoError(t, err)
	require.Contains(t, commentSQL, `COMMENT ON FUNCTION public."Audit"(integer) IS 'mixed'`)
	require.Contains(t, commentSQL, `COMMENT ON FUNCTION public."select"(integer) IS 'keyword'`)
	require.Contains(t, commentSQL, `COMMENT ON PROCEDURE public."my-func"(integer) IS 'hyphen'`)
}

func TestRound21Fixes(t *testing.T) {
	// 1. Explicit postgres ownership on new table and partition child emits OWNER TO postgres.
	partFrom := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}
	partTo := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "events",
				Owner:   "postgres",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:       "events_1",
					Expression: "RANGE (id)",
					Value:      "FOR VALUES FROM (0) TO (10)",
				}},
			}},
		}},
	}
	partSQL, err := testMetadataMigration(partFrom, partTo)
	require.NoError(t, err)
	require.Contains(t, partSQL, `ALTER TABLE "public"."events" OWNER TO postgres`)
	require.Contains(t, partSQL, `ALTER TABLE "public"."events_1" OWNER TO postgres`)

	// 2. Rebuilding a table on relkind change ('r' -> 'p') restores unchanged owner, outbound/inbound FKs,
	// triggers, RLS, policies, grants, comments (including constraint-backed indexes), and managed sequences.
	setupRebuild := func(partitioned bool) *Catalog {
		c := New()
		tableClause := ""
		if partitioned {
			tableClause = " PARTITION BY RANGE (id)"
		}
		sql := strings.ReplaceAll(`
			CREATE ROLE alice;
			CREATE FUNCTION public.trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$;
			CREATE TABLE public.ref_target (id integer PRIMARY KEY);
			CREATE TABLE public.child (
				id integer PRIMARY KEY,
				ref_id integer CONSTRAINT child_ref_fk REFERENCES public.ref_target(id)
			)__PART__;
			ALTER TABLE public.child OWNER TO alice;
			ALTER TABLE public.child ENABLE ROW LEVEL SECURITY;
			ALTER TABLE public.child FORCE ROW LEVEL SECURITY;
			CREATE POLICY child_pol ON public.child USING (true);
			CREATE TRIGGER child_trg BEFORE INSERT ON public.child FOR EACH ROW EXECUTE FUNCTION public.trg_fn();
			GRANT SELECT ON public.child TO PUBLIC;
			COMMENT ON TABLE public.child IS 'kept comment';
			CREATE TABLE public.inbound (
				child_id integer CONSTRAINT inbound_child_fk REFERENCES public.child(id)
			);
		`, "__PART__", tableClause)
		res, execErr := c.Exec(sql, nil)
		require.NoError(t, execErr)
		for _, r := range res {
			require.NoError(t, r.Error, r.SQL)
		}
		return c
	}
	inhFrom := setupRebuild(false)
	inhTo := setupRebuild(true)

	rebuildSQL := GenerateMigration(inhFrom, inhTo, Diff(inhFrom, inhTo)).SQL()
	require.Contains(t, rebuildSQL, `DROP TABLE "public"."child" CASCADE`)
	createChildIdx := strings.Index(rebuildSQL, `CREATE TABLE "public"."child"`)
	rlsIdx := strings.Index(rebuildSQL, `ENABLE ROW LEVEL SECURITY`)
	require.NotEqual(t, -1, createChildIdx)
	require.NotEqual(t, -1, rlsIdx)
	require.Less(t, createChildIdx, rlsIdx, "CREATE TABLE must precede ENABLE ROW LEVEL SECURITY:\n%s", rebuildSQL)
	require.Contains(t, rebuildSQL, `ALTER TABLE "public"."child" OWNER TO alice`)
	require.Contains(t, rebuildSQL, `ADD CONSTRAINT child_ref_fk FOREIGN KEY`)
	require.Contains(t, rebuildSQL, `ADD CONSTRAINT inbound_child_fk FOREIGN KEY`)
	require.Contains(t, rebuildSQL, `CREATE TRIGGER "child_trg"`)
	require.Contains(t, rebuildSQL, `FORCE ROW LEVEL SECURITY`)
	require.Contains(t, rebuildSQL, `CREATE POLICY "child_pol" ON "public"."child"`)
	require.Contains(t, rebuildSQL, `GRANT SELECT ON TABLE "public"."child" TO PUBLIC`)
	require.Contains(t, rebuildSQL, `COMMENT ON TABLE "public"."child" IS 'kept comment'`)

	// 3. Routine identifier containing '(' keeps full routine name quoted in comments.
	parenFrom := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	parenTo := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "audit(v2",
				Signature:  "audit(v2(integer)",
				Definition: `CREATE FUNCTION public."audit(v2"(integer) RETURNS integer LANGUAGE sql AS 'SELECT $1'`,
				Comment:    "has paren",
			}},
		}},
	}
	parenSQL, err := testMetadataMigration(parenFrom, parenTo)
	require.NoError(t, err)
	require.Contains(t, parenSQL, `COMMENT ON FUNCTION public."audit(v2"(integer) IS 'has paren'`)
}

func TestRound22Fixes(t *testing.T) {
	// 1. Regular table inheritance switch uses ALTER TABLE ... NO INHERIT / INHERIT instead of DROP + CREATE.
	inhFrom := New()
	res, err := inhFrom.Exec(`
		CREATE TABLE public.p1 (id integer);
		CREATE TABLE public.p2 (id integer, extra text);
		CREATE TABLE public.child () INHERITS (public.p1);
	`, nil)
	require.NoError(t, err)
	for _, r := range res {
		require.NoError(t, r.Error, r.SQL)
	}
	inhTo := New()
	res, err = inhTo.Exec(`
		CREATE TABLE public.p2 (id integer, extra text);
		CREATE TABLE public.child () INHERITS (public.p2);
	`, nil)
	require.NoError(t, err)
	for _, r := range res {
		require.NoError(t, r.Error, r.SQL)
	}
	inhSQL := GenerateMigration(inhFrom, inhTo, Diff(inhFrom, inhTo)).SQL()
	require.NotContains(t, inhSQL, `DROP TABLE "public"."child"`)
	require.NotContains(t, inhSQL, `CREATE TABLE "public"."child"`)
	noInheritIdx := strings.Index(inhSQL, `ALTER TABLE "public"."child" NO INHERIT "public"."p1"`)
	dropParentIdx := strings.Index(inhSQL, `DROP TABLE "public"."p1" CASCADE`)
	addColIdx := strings.Index(inhSQL, `ALTER TABLE public.child ADD COLUMN extra text`)
	inheritIdx := strings.Index(inhSQL, `ALTER TABLE "public"."child" INHERIT "public"."p2"`)
	require.NotEqual(t, -1, noInheritIdx, inhSQL)
	require.NotEqual(t, -1, dropParentIdx, inhSQL)
	require.NotEqual(t, -1, addColIdx, inhSQL)
	require.NotEqual(t, -1, inheritIdx, inhSQL)
	require.Less(t, noInheritIdx, dropParentIdx, inhSQL)
	require.Less(t, addColIdx, inheritIdx, inhSQL)

	appliedCat := inhFrom.Clone()
	applyRes, err := appliedCat.Exec(inhSQL, nil)
	require.NoError(t, err)
	for _, r := range applyRes {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(appliedCat, inhTo).Relations)

	// 2. Rebuilding a table restores managed sequence options, grants, comments, and constraint-backed index comments.
	setupSeqRebuild := func(partitioned bool) *Catalog {
		c := New()
		partClause := ""
		if partitioned {
			partClause = " PARTITION BY RANGE (id)"
		}
		sql := `
			CREATE TABLE public.items (
				id integer GENERATED BY DEFAULT AS IDENTITY (INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 START WITH 20 CACHE 3 CYCLE),
				code text CONSTRAINT items_code_key UNIQUE,
				PRIMARY KEY (id)
			)` + partClause + `;
			GRANT USAGE ON SEQUENCE public.items_id_seq TO PUBLIC;
			COMMENT ON SEQUENCE public.items_id_seq IS 'managed seq comment';
			COMMENT ON INDEX public.items_code_key IS 'backing idx comment';
		`
		execRes, execErr := c.Exec(sql, nil)
		require.NoError(t, execErr)
		for _, r := range execRes {
			require.NoError(t, r.Error, r.SQL)
		}
		return c
	}
	seqFrom := setupSeqRebuild(false)
	seqTo := setupSeqRebuild(true)
	seqRebuildSQL := GenerateMigration(seqFrom, seqTo, Diff(seqFrom, seqTo)).SQL()
	require.Contains(t, seqRebuildSQL, `DROP TABLE "public"."items" CASCADE`)
	require.Contains(t, seqRebuildSQL, `ALTER SEQUENCE "public"."items_id_seq"`)
	require.Contains(t, seqRebuildSQL, `INCREMENT BY 5`)
	require.Contains(t, seqRebuildSQL, `MINVALUE 10`)
	require.Contains(t, seqRebuildSQL, `MAXVALUE 1000`)
	require.Contains(t, seqRebuildSQL, `START WITH 20`)
	require.Contains(t, seqRebuildSQL, `CACHE 3`)
	require.Contains(t, seqRebuildSQL, `CYCLE`)
	require.Contains(t, seqRebuildSQL, `GRANT USAGE ON SEQUENCE "public"."items_id_seq" TO PUBLIC`)
	require.Contains(t, seqRebuildSQL, `COMMENT ON SEQUENCE "public"."items_id_seq" IS 'managed seq comment'`)
	require.Contains(t, seqRebuildSQL, `COMMENT ON INDEX "public"."items_code_key" IS 'backing idx comment'`)

	// Verify round-trip convergence in catalog execution.
	seqApplied := seqFrom.Clone()
	seqExecRes, err := seqApplied.Exec(seqRebuildSQL, nil)
	require.NoError(t, err)
	for _, r := range seqExecRes {
		require.NoError(t, r.Error, r.SQL)
	}
	afterRebuildDiff := Diff(seqApplied, seqTo)
	require.Empty(t, afterRebuildDiff.Relations)
	require.Empty(t, afterRebuildDiff.Sequences)
	require.Empty(t, afterRebuildDiff.Grants)
	require.Empty(t, afterRebuildDiff.Comments)
}

func TestRound23Fixes(t *testing.T) {
	// 1. MetadataMigration retains unchanged indexes and index comments when rebuilding a materialized view
	// (including when rebuilt via restoreUnloadableDependentViews).
	mvFrom := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "src",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:       "mv_items",
				Definition: "SELECT id FROM public.src",
				DependencyColumns: []*metadatapb.DependencyColumn{{
					Schema: "public", Table: "src", Column: "id",
				}},
				Indexes: []*metadatapb.IndexMetadata{{
					Name:        "mv_items_id_idx",
					Definition:  `CREATE INDEX "mv_items_id_idx" ON "public"."mv_items" USING btree (id) WITH (fillfactor = 70)`,
					Expressions: []string{"id"},
					Type:        "btree",
					Comment:     "mv index doc",
				}},
			}},
		}},
	}
	mvTo := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "src",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "bigint"}},
			}},
			MaterializedViews: mvFrom.Schemas[0].MaterializedViews,
		}},
	}
	mvSQL, err := testMetadataMigration(mvFrom, mvTo)
	require.NoError(t, err)
	require.Contains(t, mvSQL, `DROP MATERIALIZED VIEW`)
	require.Contains(t, mvSQL, `CREATE MATERIALIZED VIEW`)
	require.Contains(t, mvSQL, `WITH (fillfactor = 70)`)
	require.Contains(t, mvSQL, `COMMENT ON INDEX "public"."mv_items_id_idx" IS 'mv index doc'`)
	createMVPos := strings.Index(mvSQL, `CREATE MATERIALIZED VIEW`)
	createIdxPos := strings.Index(mvSQL, `CREATE INDEX`)
	require.NotEqual(t, -1, createMVPos)
	require.NotEqual(t, -1, createIdxPos)
	require.Less(t, createMVPos, createIdxPos, mvSQL)

	// 2. Direct GenerateMigration generates DDL for Diff.ForeignTables (add, alter, server switch, drop)
	// and converges when executed.
	setupFT := func(sql string) *Catalog {
		c := New()
		res, execErr := c.Exec(`
			CREATE FOREIGN DATA WRAPPER fdw_a;
			CREATE FOREIGN DATA WRAPPER fdw_b;
			CREATE SERVER srv_a FOREIGN DATA WRAPPER fdw_a;
			CREATE SERVER srv_b FOREIGN DATA WRAPPER fdw_b;
		`+sql, &ExecOptions{ContinueOnError: true})
		require.NoError(t, execErr)
		for _, r := range res {
			require.NoError(t, r.Error, r.SQL)
		}
		return c
	}
	ftFrom := setupFT(`
		CREATE FOREIGN TABLE public.ft_mod (id integer OPTIONS (column_name 'id'), drop_me text) SERVER srv_a OPTIONS (schema_name 's1', old_opt 'v');
		CREATE FOREIGN TABLE public.ft_srv (id integer) SERVER srv_a;
		COMMENT ON FOREIGN TABLE public.ft_srv IS 'kept ft comment';
		CREATE FOREIGN TABLE public.ft_drop (id integer) SERVER srv_a;
	`)
	ftTo := setupFT(`
		CREATE FOREIGN TABLE public.ft_mod (id bigint OPTIONS (column_name 'id_new', extra '1') NOT NULL, note text) SERVER srv_a OPTIONS (schema_name 's2', new_opt 'x');
		CREATE FOREIGN TABLE public.ft_srv (id integer) SERVER srv_b;
		COMMENT ON FOREIGN TABLE public.ft_srv IS 'kept ft comment';
		CREATE FOREIGN TABLE public.ft_add (id integer) SERVER srv_a OPTIONS (table_name 't1');
	`)
	ftDiff := Diff(ftFrom, ftTo)
	require.NotEmpty(t, ftDiff.ForeignTables)
	ftSQL := GenerateMigration(ftFrom, ftTo, ftDiff).SQL()
	require.Contains(t, ftSQL, `DROP FOREIGN TABLE "public"."ft_drop" CASCADE`)
	require.Contains(t, ftSQL, `CREATE FOREIGN TABLE "public"."ft_add"`)
	require.Contains(t, ftSQL, `ALTER FOREIGN TABLE "public"."ft_mod" OPTIONS (DROP old_opt, SET schema_name 's2', ADD new_opt 'x')`)
	require.Contains(t, ftSQL, `ALTER FOREIGN TABLE "public"."ft_mod" ALTER COLUMN id OPTIONS (SET column_name 'id_new', ADD extra '1')`)
	require.Contains(t, ftSQL, `COMMENT ON FOREIGN TABLE "public"."ft_srv" IS 'kept ft comment'`)

	ftApplied := ftFrom.Clone()
	ftExecRes, err := ftApplied.Exec(ftSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range ftExecRes {
		require.NoError(t, r.Error, r.SQL)
	}
	remFT := Diff(ftApplied, ftTo)
	require.Empty(t, remFT.ForeignTables)
	require.Empty(t, remFT.Comments)
}

func TestRound24Fixes(t *testing.T) {
	// 1. FDW -> Server -> ForeignTable create/alter/drop ordering and convergence.
	fdwFrom := New()
	_, err := fdwFrom.Exec(`
		CREATE FUNCTION public.fdw_h1() RETURNS fdw_handler AS '' LANGUAGE C;
		CREATE FUNCTION public.fdw_v1(text[], oid) RETURNS void AS '' LANGUAGE C;
		CREATE FUNCTION public.fdw_h_drop() RETURNS fdw_handler AS '' LANGUAGE C;
		CREATE FOREIGN DATA WRAPPER w_keep HANDLER public.fdw_h1 VALIDATOR public.fdw_v1 OPTIONS (foo 'bar');
		CREATE FOREIGN DATA WRAPPER w_drop HANDLER public.fdw_h_drop;
		CREATE SERVER srv_keep TYPE 'pg' VERSION '16' FOREIGN DATA WRAPPER w_keep OPTIONS (host 'h1');
		CREATE SERVER srv_rebuild TYPE 'old' FOREIGN DATA WRAPPER w_keep;
		CREATE SERVER srv_drop FOREIGN DATA WRAPPER w_drop;
		CREATE FOREIGN TABLE public.ft_on_rebuild (id integer) SERVER srv_rebuild;
		CREATE FOREIGN TABLE public.ft_drop (id integer) SERVER srv_drop;
		COMMENT ON FOREIGN TABLE public.ft_drop IS 'drop me cleanly';
	`, nil)
	require.NoError(t, err)

	fdwTo := New()
	_, err = fdwTo.Exec(`
		CREATE FUNCTION public.fdw_h1() RETURNS fdw_handler AS '' LANGUAGE C;
		CREATE FUNCTION public.fdw_h2() RETURNS fdw_handler AS '' LANGUAGE C;
		CREATE FOREIGN DATA WRAPPER w_keep HANDLER public.fdw_h1 NO VALIDATOR OPTIONS (foo 'baz', extra '1');
		CREATE FOREIGN DATA WRAPPER w_new HANDLER public.fdw_h2;
		CREATE SERVER srv_keep TYPE 'pg' VERSION NULL FOREIGN DATA WRAPPER w_keep OPTIONS (host 'h2');
		CREATE SERVER srv_rebuild TYPE 'new' VERSION '17' FOREIGN DATA WRAPPER w_keep;
		CREATE SERVER srv_new FOREIGN DATA WRAPPER w_new OPTIONS (dbname 'd1');
		CREATE FOREIGN TABLE public.ft_on_rebuild (id integer) SERVER srv_rebuild;
		CREATE FOREIGN TABLE public.ft_new (id integer) SERVER srv_new;
	`, nil)
	require.NoError(t, err)

	fdwDiff := Diff(fdwFrom, fdwTo)
	require.NotEmpty(t, fdwDiff.ForeignDataWrappers)
	require.NotEmpty(t, fdwDiff.ForeignServers)
	fdwPlan := GenerateMigration(fdwFrom, fdwTo, fdwDiff)
	fdwSQL := fdwPlan.SQL()
	require.NotContains(t, fdwSQL, `COMMENT ON FOREIGN TABLE "public"."ft_drop" IS NULL`)
	require.NotContains(t, fdwSQL, `DROP SERVER "srv_rebuild" CASCADE`)
	require.NotEmpty(t, fdwPlan.Warnings())

	fdwApplied := fdwFrom.Clone()
	execRes, err := fdwApplied.Exec(fdwSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range execRes {
		require.NoError(t, r.Error, r.SQL)
	}
	remFdw := Diff(fdwApplied, fdwTo)
	require.Empty(t, remFdw.ForeignDataWrappers)
	require.Empty(t, remFdw.ForeignTables)
	require.Empty(t, remFdw.Functions)
	// srv_keep and srv_new converged; only unalterable srv_rebuild remains in diff.
	require.Len(t, remFdw.ForeignServers, 1)
	require.Equal(t, "srv_rebuild", remFdw.ForeignServers[0].Name)

	// 2. Foreign table server switch preserves dependent views and single retained trigger;
	// column type change + modified CHECK constraint avoids duplicate drop; NOT VALID is preserved.
	setupFTComplex := func(ddl string) *Catalog {
		c := New()
		_, err := c.Exec(`
			CREATE FOREIGN DATA WRAPPER fdw_a;
			CREATE SERVER srv_a FOREIGN DATA WRAPPER fdw_a;
			CREATE SERVER srv_b FOREIGN DATA WRAPPER fdw_a;
			CREATE FUNCTION public.ft_trig_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
		`+ddl, nil)
		require.NoError(t, err)
		return c
	}

	ftFrom := setupFTComplex(`
		CREATE FOREIGN TABLE public.ft_switch (id integer, val integer) SERVER srv_a;
		CREATE TRIGGER tr_keep BEFORE INSERT ON public.ft_switch FOR EACH ROW EXECUTE FUNCTION public.ft_trig_fn();
		CREATE VIEW public.v_ft AS SELECT id, val FROM public.ft_switch;
		CREATE FOREIGN TABLE public.ft_chk (
			id integer,
			val integer,
			CONSTRAINT chk_mod CHECK (val > 0),
			CONSTRAINT chk_drop CHECK (val < 100),
			CONSTRAINT chk_same CHECK (val <> 0)
		) SERVER srv_a;
	`)
	ftTo := setupFTComplex(`
		CREATE FOREIGN TABLE public.ft_switch (id integer, val integer) SERVER srv_b;
		CREATE TRIGGER tr_keep BEFORE INSERT ON public.ft_switch FOR EACH ROW EXECUTE FUNCTION public.ft_trig_fn();
		CREATE VIEW public.v_ft AS SELECT id, val FROM public.ft_switch;
		CREATE FOREIGN TABLE public.ft_chk (
			id integer,
			val bigint,
			CONSTRAINT chk_mod CHECK (val > 10),
			CONSTRAINT chk_same CHECK (val <> 0),
			CONSTRAINT chk_nv CHECK (id > 0) NOT VALID
		) SERVER srv_a;
	`)

	ftDiff := Diff(ftFrom, ftTo)
	ftSQL := GenerateMigration(ftFrom, ftTo, ftDiff).SQL()
	require.Equal(t, 1, strings.Count(ftSQL, `CREATE TRIGGER "tr_keep"`))
	require.Equal(t, 1, strings.Count(ftSQL, `DROP CONSTRAINT chk_mod`))
	require.Contains(t, ftSQL, `CREATE VIEW "public"."v_ft"`)
	require.Contains(t, ftSQL, "NOT VALID")

	ftApplied := ftFrom.Clone()
	ftExecRes, err := ftApplied.Exec(ftSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range ftExecRes {
		require.NoError(t, r.Error, r.SQL)
	}
	remFT2 := Diff(ftApplied, ftTo)
	require.Empty(t, remFT2.ForeignTables)
	require.Empty(t, remFT2.Relations)
	require.Empty(t, remFT2.Comments)
}

func TestRound25Fixes(t *testing.T) {
	// 1. Foreign table partitions (add partition, detach, bound change) + grant revokes suppressed on dropped/altered foreign tables
	// + FDW/Server ownership changes + range type support functions.
	setupR25 := func(ddl string) *Catalog {
		c := New()
		_, err := c.Exec(`
			CREATE ROLE alice;
			CREATE FOREIGN DATA WRAPPER w_own;
			CREATE SERVER s_own FOREIGN DATA WRAPPER w_own;
			CREATE TABLE public.part_parent (id integer) PARTITION BY RANGE (id);
			CREATE TABLE public.part_other (id integer) PARTITION BY RANGE (id);
		`+ddl, nil)
		require.NoError(t, err)
		return c
	}

	rFrom := setupR25(`
		CREATE FOREIGN TABLE public.ft_detach PARTITION OF public.part_parent FOR VALUES FROM (0) TO (10) SERVER s_own;
		CREATE FOREIGN TABLE public.ft_rebind PARTITION OF public.part_parent FOR VALUES FROM (10) TO (20) SERVER s_own;
		CREATE FOREIGN TABLE public.ft_drop_grant (id integer, secret text) SERVER s_own;
		CREATE FOREIGN TABLE public.ft_col_grant (id integer, drop_secret text) SERVER s_own;
		GRANT SELECT ON public.ft_drop_grant TO alice;
		GRANT SELECT (id, drop_secret) ON public.ft_drop_grant TO alice;
		GRANT SELECT (id, drop_secret) ON public.ft_col_grant TO alice;
		CREATE FUNCTION public.r_canon(int4range) RETURNS int4range LANGUAGE sql AS 'SELECT $1';
		CREATE FUNCTION public.r_canon2(int4range) RETURNS int4range LANGUAGE sql AS 'SELECT $1';
	`)
	rTo := setupR25(`
		ALTER FOREIGN DATA WRAPPER w_own OWNER TO alice;
		ALTER SERVER s_own OWNER TO alice;
		CREATE FOREIGN TABLE public.ft_detach (id integer) SERVER s_own;
		CREATE FOREIGN TABLE public.ft_rebind PARTITION OF public.part_other FOR VALUES FROM (100) TO (200) SERVER s_own;
		CREATE FOREIGN TABLE public.ft_add_part PARTITION OF public.part_parent FOR VALUES FROM (20) TO (30) SERVER s_own;
		CREATE FOREIGN TABLE public.ft_col_grant (id integer) SERVER s_own;
		GRANT SELECT (id) ON public.ft_col_grant TO alice;
		CREATE FUNCTION public.r_canon(int4range) RETURNS int4range LANGUAGE sql IMMUTABLE AS 'SELECT $1';
		CREATE FUNCTION public.r_canon2(int4range) RETURNS int4range LANGUAGE sql IMMUTABLE AS 'SELECT $1';
		CREATE FUNCTION public.f_diff(float8, float8) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT $1 - $2';
		CREATE TYPE public.f8range AS RANGE (SUBTYPE = float8, SUBTYPE_DIFF = public.f_diff);
	`)
	// Range modification case via standalone subtype_diff change.
	rDiffFrom := New()
	_, err := rDiffFrom.Exec(`
		CREATE FUNCTION public.d1(float8, float8) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT $1 - $2';
		CREATE FUNCTION public.d2(float8, float8) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT $1 - $2';
		CREATE TYPE public.myfloatrange AS RANGE (SUBTYPE = float8, SUBTYPE_DIFF = public.d1);
	`, nil)
	require.NoError(t, err)
	rDiffTo := New()
	_, err = rDiffTo.Exec(`
		CREATE FUNCTION public.d1(float8, float8) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT $1 - $2';
		CREATE FUNCTION public.d2(float8, float8) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT $1 - $2';
		CREATE TYPE public.myfloatrange AS RANGE (SUBTYPE = float8, SUBTYPE_OPCLASS = float8_ops, SUBTYPE_DIFF = public.d2);
	`, nil)
	require.NoError(t, err)
	rDiff := Diff(rDiffFrom, rDiffTo)
	require.Len(t, rDiff.Ranges, 1)
	rSQL := GenerateMigration(rDiffFrom, rDiffTo, rDiff).SQL()
	require.Contains(t, rSQL, "SUBTYPE_OPCLASS = float8_ops")
	require.Contains(t, rSQL, `SUBTYPE_DIFF = public.d2`)
	rApplied := rDiffFrom.Clone()
	execRange, err := rApplied.Exec(rSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range execRange {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(rApplied, rDiffTo).Ranges)

	diff25 := Diff(rFrom, rTo)
	require.NotEmpty(t, diff25.ForeignDataWrappers)
	require.NotEmpty(t, diff25.ForeignServers)
	sql25 := GenerateMigration(rFrom, rTo, diff25).SQL()
	require.Contains(t, sql25, `ALTER FOREIGN DATA WRAPPER w_own OWNER TO alice`)
	require.Contains(t, sql25, `ALTER SERVER s_own OWNER TO alice`)
	require.Contains(t, sql25, `DETACH PARTITION "public"."ft_detach"`)
	require.Contains(t, sql25, `ATTACH PARTITION "public"."ft_add_part" FOR VALUES FROM ('20') TO ('30')`)
	require.Contains(t, sql25, `SUBTYPE_DIFF = public.f_diff`)
	require.NotContains(t, sql25, `ON TABLE "public"."ft_drop_grant"`)
	require.NotContains(t, sql25, `REVOKE SELECT ("drop_secret") ON TABLE "public"."ft_col_grant"`)

	applied25 := rFrom.Clone()
	exec25, err := applied25.Exec(sql25, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range exec25 {
		require.NoError(t, r.Error, r.SQL)
	}
	rem25 := Diff(applied25, rTo)
	require.Empty(t, rem25.ForeignDataWrappers)
	require.Empty(t, rem25.ForeignServers)
	require.Empty(t, rem25.ForeignTables)
	require.Empty(t, rem25.Ranges)
	require.Empty(t, rem25.Grants)

	// 2. External table check constraints & triggers in LoadMetadata + MetadataMigration.
	extFrom := &metadatapb.DatabaseSchemaMetadata{
		Name:                "postgres",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{Name: "w1"}},
		ForeignServers:      []*metadatapb.ForeignServerMetadata{{Name: "srv1", Fdw: "w1"}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{{
				Name:       "ext_trig_fn",
				Signature:  "ext_trig_fn()",
				Definition: "CREATE FUNCTION public.ext_trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$",
			}},
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "ext_t",
				ExternalServerName: "srv1",
				Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				CheckConstraints:   []*metadatapb.CheckConstraintMetadata{{Name: "chk_id", Expression: "(id > 0)"}},
				Triggers: []*metadatapb.TriggerMetadata{{
					Name:   "tr_ext",
					Timing: "BEFORE",
					Event:  "INSERT",
					Body:   "CREATE TRIGGER tr_ext BEFORE INSERT ON public.ext_t FOR EACH ROW EXECUTE FUNCTION public.ext_trig_fn()",
				}},
			}},
		}},
	}
	catLoaded := New()
	_, err = catLoaded.LoadMetadata(context.Background(), extFrom, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	catLoadedExtRel := catLoaded.GetRelation("public", "ext_t")
	require.NotNil(t, catLoadedExtRel)
	require.Len(t, catLoaded.TriggersOf(catLoadedExtRel.OID), 1)
	require.Equal(t, "tr_ext", catLoaded.TriggersOf(catLoadedExtRel.OID)[0].Name)
	catAfterCheck := catLoaded.Clone()
	_, err = catAfterCheck.Exec(`
		ALTER FOREIGN TABLE public.ext_t DROP CONSTRAINT chk_id;
		DROP TRIGGER tr_ext ON public.ext_t;
	`, nil)
	require.NoError(t, err)
	appliedExt := ApplyDiffToMetadata(extFrom, catLoaded, catAfterCheck, Diff(catLoaded, catAfterCheck))
	require.Empty(t, findSchema(appliedExt, "public").ExternalTables[0].CheckConstraints)
	require.Empty(t, findSchema(appliedExt, "public").ExternalTables[0].Triggers)

	extLoaded, extReport, err := loadMetadataForMigration(context.Background(), extFrom, nil)
	require.NoError(t, err)
	require.Empty(t, extReport.Missing)
	extRel := extLoaded.GetRelation("public", "ext_t")
	require.NotNil(t, extRel)
	require.Len(t, extLoaded.ConstraintsOf(extRel.OID), 1)
	require.Equal(t, "chk_id", extLoaded.ConstraintsOf(extRel.OID)[0].Name)
	require.Len(t, extLoaded.TriggersOf(extRel.OID), 1)
	require.Equal(t, "tr_ext", extLoaded.TriggersOf(extRel.OID)[0].Name)
}

func TestRound26Fixes(t *testing.T) {
	// 1. Canonical range uses shell-type -> function -> full CREATE TYPE AS RANGE sequence.
	canonFrom := New()
	canonTo := New()
	canonRes, err := canonTo.Exec(`
		CREATE TYPE public.c_range;
		CREATE FUNCTION public.c_range_canon(public.c_range) RETURNS public.c_range LANGUAGE C IMMUTABLE AS '';
		CREATE TYPE public.c_range AS RANGE (SUBTYPE = integer, CANONICAL = public.c_range_canon);
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range canonRes {
		require.NoError(t, r.Error, r.SQL)
	}
	canonDiff := Diff(canonFrom, canonTo)
	canonSQL := GenerateMigration(canonFrom, canonTo, canonDiff).SQL()
	require.Contains(t, canonSQL, `CREATE TYPE "public"."c_range";`)
	require.Contains(t, canonSQL, `CANONICAL = public.c_range_canon`)
	canonApplied := canonFrom.Clone()
	execCanon, err := canonApplied.Exec(canonSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range execCanon {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(canonApplied, canonTo).Ranges)
	require.Empty(t, Diff(canonApplied, canonTo).Functions)

	// 2. New FDW and Server restore default postgres owner when migration runs under another role.
	ownFrom := New()
	_, err = ownFrom.Exec(`CREATE ROLE alice;`, nil)
	require.NoError(t, err)
	ownTo := New()
	_, err = ownTo.Exec(`
		CREATE ROLE alice;
		CREATE FOREIGN DATA WRAPPER w_def;
		CREATE SERVER s_def FOREIGN DATA WRAPPER w_def;
	`, nil)
	require.NoError(t, err)
	ownSQL := GenerateMigration(ownFrom, ownTo, Diff(ownFrom, ownTo)).SQL()
	require.Contains(t, ownSQL, `ALTER FOREIGN DATA WRAPPER w_def OWNER TO postgres`)
	require.Contains(t, ownSQL, `ALTER SERVER s_def OWNER TO postgres`)
	ownApplied := ownFrom.Clone()
	_, err = ownApplied.Exec(`SET ROLE alice;`, nil)
	require.NoError(t, err)
	execOwn, err := ownApplied.Exec(ownSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range execOwn {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(ownApplied, ownTo).ForeignDataWrappers)
	require.Empty(t, Diff(ownApplied, ownTo).ForeignServers)

	// 3. Recreating a range without previous options clears those options in ApplyDiffToMetadata.
	rangeMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			RangeTypes: []*metadatapb.RangeTypeMetadata{{
				Name:           "r_opt",
				Subtype:        "float8",
				SubtypeOpclass: "float8_ops",
				SubtypeDiff:    "public.r_diff",
				Comment:        "keep comment",
			}},
		}},
	}
	rBefore := New()
	_, err = rBefore.Exec(`
		CREATE FUNCTION public.r_diff(float8, float8) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT $1 - $2';
		CREATE TYPE public.r_opt AS RANGE (SUBTYPE = float8, SUBTYPE_OPCLASS = float8_ops, SUBTYPE_DIFF = public.r_diff);
	`, nil)
	require.NoError(t, err)
	rAfter := rBefore.Clone()
	_, err = rAfter.Exec(`
		DROP TYPE public.r_opt;
		CREATE TYPE public.r_opt AS RANGE (SUBTYPE = float8);
	`, nil)
	require.NoError(t, err)
	rAppliedMeta := ApplyDiffToMetadata(rangeMeta, rBefore, rAfter, Diff(rBefore, rAfter))
	rGot := findSchema(rAppliedMeta, "public").RangeTypes[0]
	require.Empty(t, rGot.SubtypeOpclass)
	require.Empty(t, rGot.SubtypeDiff)
	require.Equal(t, "keep comment", rGot.Comment)

	// 4. Table rebuild recreates modified RLS policy instead of emitting ALTER POLICY on dropped table.
	polFrom := New()
	_, err = polFrom.Exec(`
		CREATE TABLE public.t_pol (id integer);
		ALTER TABLE public.t_pol ENABLE ROW LEVEL SECURITY;
		CREATE POLICY p_mod ON public.t_pol FOR SELECT USING (id > 0);
	`, nil)
	require.NoError(t, err)
	polTo := New()
	_, err = polTo.Exec(`
		CREATE TABLE public.t_pol (id integer) PARTITION BY RANGE (id);
		ALTER TABLE public.t_pol ENABLE ROW LEVEL SECURITY;
		CREATE POLICY p_mod ON public.t_pol FOR SELECT USING (id > 10);
	`, nil)
	require.NoError(t, err)
	polSQL := GenerateMigration(polFrom, polTo, Diff(polFrom, polTo)).SQL()
	require.NotContains(t, polSQL, `ALTER POLICY`)
	require.Contains(t, polSQL, `CREATE POLICY "p_mod"`)
	polApplied := polFrom.Clone()
	execPol, err := polApplied.Exec(polSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range execPol {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(polApplied, polTo).Relations)

	// 5. Dotted schema/function names in trigger identities are disambiguated.
	trgFrom := New()
	_, err = trgFrom.Exec(`
		CREATE SCHEMA "a.b";
		CREATE SCHEMA a;
		CREATE FUNCTION "a.b".c() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE FUNCTION a."b.c"() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE TABLE public.t_trig (id integer);
		CREATE TRIGGER tr_dot BEFORE INSERT ON public.t_trig FOR EACH ROW EXECUTE FUNCTION "a.b".c();
	`, nil)
	require.NoError(t, err)
	trgTo := New()
	_, err = trgTo.Exec(`
		CREATE SCHEMA "a.b";
		CREATE SCHEMA a;
		CREATE FUNCTION "a.b".c() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE FUNCTION a."b.c"() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE TABLE public.t_trig (id integer);
		CREATE TRIGGER tr_dot BEFORE INSERT ON public.t_trig FOR EACH ROW EXECUTE FUNCTION a."b.c"();
	`, nil)
	require.NoError(t, err)
	trgDiff := Diff(trgFrom, trgTo)
	require.Len(t, trgDiff.Relations, 1)
	require.Len(t, trgDiff.Relations[0].Triggers, 1)
	trgSQL := GenerateMigration(trgFrom, trgTo, trgDiff).SQL()
	require.Contains(t, trgSQL, `"b.c"()`)
}

func TestRound27Fixes(t *testing.T) {
	// 1. DefineRange validates SUBTYPE_DIFF and CANONICAL signatures and records dependencies.
	c := New()
	_, err := c.Exec(`CREATE FUNCTION public.bad_diff(float8, float8) RETURNS integer LANGUAGE sql AS 'SELECT 1';`, nil)
	require.NoError(t, err)
	res, err := c.Exec(`CREATE TYPE public.missing_diff_r AS RANGE (SUBTYPE = float8, SUBTYPE_DIFF = public.no_such_diff);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Error(t, res[0].Error)
	var catErr *Error
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeUndefinedFunction, catErr.Code)
	require.Nil(t, c.typeByName[typeKey{ns: PublicNamespace, name: "missing_diff_r"}])

	res, err = c.Exec(`CREATE TYPE public.bad_ret_diff_r AS RANGE (SUBTYPE = float8, SUBTYPE_DIFF = public.bad_diff);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Error(t, res[0].Error)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeInvalidObjectDefinition, catErr.Code)
	require.Nil(t, c.typeByName[typeKey{ns: PublicNamespace, name: "bad_ret_diff_r"}])

	res, err = c.Exec(`CREATE TYPE public.missing_canon_r AS RANGE (SUBTYPE = int4, CANONICAL = public.no_canon);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Error(t, res[0].Error)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeUndefinedFunction, catErr.Code)
	require.Nil(t, c.typeByName[typeKey{ns: PublicNamespace, name: "missing_canon_r"}])

	_, err = c.Exec(`
		CREATE TYPE public.bad_canon_r;
		CREATE FUNCTION public.bad_canon_fn(public.bad_canon_r) RETURNS integer LANGUAGE C IMMUTABLE AS 'internal', 'bad_canon_fn';
	`, nil)
	require.NoError(t, err)
	res, err = c.Exec(`CREATE TYPE public.bad_canon_r AS RANGE (SUBTYPE = int4, CANONICAL = public.bad_canon_fn);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Error(t, res[0].Error)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeInvalidObjectDefinition, catErr.Code)
	require.False(t, c.typeByName[typeKey{ns: PublicNamespace, name: "bad_canon_r"}].IsDefined)

	_, err = c.Exec(`
		CREATE TYPE public.good_r;
		CREATE FUNCTION public.good_canon(public.good_r) RETURNS public.good_r LANGUAGE C IMMUTABLE AS 'internal', 'good_canon';
		CREATE FUNCTION public.good_diff(int4, int4) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT ($1 - $2)::float8';
		CREATE TYPE public.good_r AS RANGE (SUBTYPE = int4, CANONICAL = public.good_canon, SUBTYPE_DIFF = public.good_diff);
	`, nil)
	require.NoError(t, err)
	res, err = c.Exec(`DROP FUNCTION public.good_canon(public.good_r);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Error(t, res[0].Error)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeDependentObjects, catErr.Code)
	res, err = c.Exec(`DROP FUNCTION public.good_diff(int4, int4);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Error(t, res[0].Error)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeDependentObjects, catErr.Code)

	crFrom := New()
	crTo := New()
	_, err = crTo.Exec(`
		CREATE TYPE public.good_r;
		CREATE FUNCTION public.good_canon(public.good_r) RETURNS public.good_r LANGUAGE C IMMUTABLE AS 'internal', 'good_canon';
		CREATE FUNCTION public.good_diff(int4, int4) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT ($1 - $2)::float8';
		CREATE TYPE public.good_r AS RANGE (SUBTYPE = int4, CANONICAL = public.good_canon, SUBTYPE_DIFF = public.good_diff);
		CREATE TABLE public.t_good_r (r public.good_r);
	`, nil)
	require.NoError(t, err)
	crSQL := GenerateMigration(crFrom, crTo, Diff(crFrom, crTo)).SQL()
	crApplied := crFrom.Clone()
	crRes, err := crApplied.Exec(crSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range crRes {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(crApplied, crTo).Ranges)
	require.Empty(t, Diff(crApplied, crTo).Relations)

	// 2. Reordered foreign options do not trigger false DiffModify entries.
	optFrom := New()
	_, err = optFrom.Exec(`
		CREATE FOREIGN DATA WRAPPER w_ord OPTIONS (a '1', b '2');
		CREATE SERVER s_ord FOREIGN DATA WRAPPER w_ord OPTIONS (host 'h1', port '5432');
		CREATE FOREIGN TABLE public.ft_ord (id integer OPTIONS (col_a 'x', col_b 'y')) SERVER s_ord OPTIONS (schema_name 's', table_name 't');
	`, nil)
	require.NoError(t, err)
	optTo := New()
	_, err = optTo.Exec(`
		CREATE FOREIGN DATA WRAPPER w_ord OPTIONS (b '2', a '1');
		CREATE SERVER s_ord FOREIGN DATA WRAPPER w_ord OPTIONS (port '5432', host 'h1');
		CREATE FOREIGN TABLE public.ft_ord (id integer OPTIONS (col_b 'y', col_a 'x')) SERVER s_ord OPTIONS (table_name 't', schema_name 's');
	`, nil)
	require.NoError(t, err)
	optDiff := Diff(optFrom, optTo)
	require.Empty(t, optDiff.ForeignDataWrappers)
	require.Empty(t, optDiff.ForeignServers)
	require.Empty(t, optDiff.ForeignTables)
	require.Empty(t, GenerateMigration(optFrom, optTo, optDiff).Ops)
}

func TestRound28Fixes(t *testing.T) {
	// 1. Validate SUBTYPE_OPCLASS and track dependency on user operator classes.
	c := New()
	var catErr *Error
	res, err := c.Exec(`CREATE TYPE public.bad_opc_r AS RANGE (SUBTYPE = int4, SUBTYPE_OPCLASS = no_such_ops);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeUndefinedObject, catErr.Code)

	res, err = c.Exec(`CREATE TYPE public.mismatch_opc_r AS RANGE (SUBTYPE = int4, SUBTYPE_OPCLASS = text_ops);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeDatatypeMismatch, catErr.Code)

	_, err = c.Exec(`
		CREATE OPERATOR CLASS public.custom_int4_ops FOR TYPE int4 USING btree AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int4, int4);
		CREATE TYPE public.custom_opc_r AS RANGE (SUBTYPE = int4, SUBTYPE_OPCLASS = public.custom_int4_ops);
	`, nil)
	require.NoError(t, err)
	res, err = c.Exec(`DROP OPERATOR CLASS public.custom_int4_ops USING btree;`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeDependentObjects, catErr.Code)

	// 2. ALTER FUNCTION RENAME TO and SET SCHEMA update range support routine names.
	rCat := New()
	_, err = rCat.Exec(`
		CREATE SCHEMA s1;
		CREATE TYPE public.ren_r;
		CREATE FUNCTION public.ren_canon(public.ren_r) RETURNS public.ren_r LANGUAGE C IMMUTABLE AS 'internal', 'ren_canon';
		CREATE FUNCTION public.ren_diff(int4, int4) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT ($1 - $2)::float8';
		CREATE TYPE public.ren_r AS RANGE (SUBTYPE = int4, CANONICAL = public.ren_canon, SUBTYPE_DIFF = public.ren_diff);
		ALTER FUNCTION public.ren_canon(public.ren_r) RENAME TO ren_canon2;
		ALTER FUNCTION public.ren_diff(int4, int4) SET SCHEMA s1;
	`, nil)
	require.NoError(t, err)
	rtOID, _, err := rCat.ResolveType(TypeName{Schema: "public", Name: "ren_r", TypeMod: -1})
	require.NoError(t, err)
	rtInfo := rCat.RangeInfo(rtOID)
	require.NotNil(t, rtInfo)
	require.Equal(t, "public.ren_canon2", rtInfo.Canonical)
	require.Equal(t, "s1.ren_diff", rtInfo.SubtypeDiff)

	empty := New()
	renSQL := GenerateMigration(empty, rCat, Diff(empty, rCat)).SQL()
	require.Contains(t, renSQL, `CANONICAL = public.ren_canon2`)
	require.Contains(t, renSQL, `SUBTYPE_DIFF = s1.ren_diff`)
	renApplied := empty.Clone()
	renExec, err := renApplied.Exec(renSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range renExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(renApplied, rCat).Ranges)

	// 3. Explicit MULTIRANGE_TYPE_NAME is created, serialized by ApplyDiffToMetadata, and migrated.
	mrBefore := New()
	_, err = mrBefore.Exec(`CREATE SCHEMA app;`, nil)
	require.NoError(t, err)
	mrAfter := mrBefore.Clone()
	_, err = mrAfter.Exec(`
		CREATE TYPE app.r_custom AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = app.custom_mr);
		CREATE TABLE app.use_mr (m app.custom_mr);
	`, nil)
	require.NoError(t, err)
	mrMeta := ApplyDiffToMetadata(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "app"}},
	}, mrBefore, mrAfter, Diff(mrBefore, mrAfter))
	appMeta := findSchema(mrMeta, "app")
	require.NotNil(t, appMeta)
	require.Len(t, appMeta.RangeTypes, 1)
	require.Equal(t, "app.custom_mr", appMeta.RangeTypes[0].MultirangeName)

	mrSQL := GenerateMigration(mrBefore, mrAfter, Diff(mrBefore, mrAfter)).SQL()
	require.Contains(t, mrSQL, `MULTIRANGE_TYPE_NAME = "app"."custom_mr"`)
	mrApplied := mrBefore.Clone()
	mrExec, err := mrApplied.Exec(mrSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range mrExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(mrApplied, mrAfter).Ranges)
	require.Empty(t, Diff(mrApplied, mrAfter).Relations)
}

func TestRound29Fixes(t *testing.T) {
	// 1. Subtypes without a default btree operator class (e.g. json) are rejected unless supplied.
	c := New()
	var catErr *Error
	res, err := c.Exec(`CREATE TYPE public.json_r AS RANGE (SUBTYPE = json);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeUndefinedObject, catErr.Code)

	// Custom default btree operator class allows range creation and records dependency.
	_, err = c.Exec(`
		CREATE TYPE public.custom_base;
		CREATE FUNCTION public.cb_in(cstring) RETURNS public.custom_base LANGUAGE C STRICT IMMUTABLE AS 'internal', 'cb_in';
		CREATE FUNCTION public.cb_out(public.custom_base) RETURNS cstring LANGUAGE C STRICT IMMUTABLE AS 'internal', 'cb_out';
		CREATE TYPE public.custom_base (INPUT = public.cb_in, OUTPUT = public.cb_out);
		CREATE FUNCTION public.cb_cmp(public.custom_base, public.custom_base) RETURNS integer LANGUAGE sql IMMUTABLE AS 'SELECT 0';
		CREATE OPERATOR CLASS public.cb_default_ops DEFAULT FOR TYPE public.custom_base USING btree AS FUNCTION 1 public.cb_cmp(public.custom_base, public.custom_base);
		CREATE TYPE public.cb_r AS RANGE (SUBTYPE = public.custom_base);
	`, nil)
	require.NoError(t, err)
	res, err = c.Exec(`DROP OPERATOR CLASS public.cb_default_ops USING btree;`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.ErrorAs(t, res[0].Error, &catErr)
	require.Equal(t, CodeDependentObjects, catErr.Code)

	// 2. Non-immutable CANONICAL and SUBTYPE_DIFF functions are rejected.
	vCat := New()
	_, err = vCat.Exec(`
		CREATE FUNCTION public.mut_diff(float8, float8) RETURNS float8 LANGUAGE sql AS 'SELECT $1 - $2';
		CREATE FUNCTION public.st_diff(float8, float8) RETURNS float8 LANGUAGE sql STABLE AS 'SELECT $1 - $2';
		CREATE TYPE public.mut_canon_r;
		CREATE FUNCTION public.mut_canon(public.mut_canon_r) RETURNS public.mut_canon_r LANGUAGE C AS 'internal', 'mut_canon';
	`, nil)
	require.NoError(t, err)
	for _, ddl := range []string{
		`CREATE TYPE public.r_mut_d AS RANGE (SUBTYPE = float8, SUBTYPE_DIFF = public.mut_diff);`,
		`CREATE TYPE public.r_st_d AS RANGE (SUBTYPE = float8, SUBTYPE_DIFF = public.st_diff);`,
		`CREATE TYPE public.mut_canon_r AS RANGE (SUBTYPE = int4, CANONICAL = public.mut_canon);`,
	} {
		vRes, vErr := vCat.Exec(ddl, &ExecOptions{ContinueOnError: true})
		require.NoError(t, vErr)
		require.ErrorAs(t, vRes[0].Error, &catErr, ddl)
		require.Equal(t, CodeInvalidObjectDefinition, catErr.Code, ddl)
	}

	// 3. Generated multirange name follows PostgreSQL makeMultirangeTypeName rules.
	mBefore := New()
	mAfter := mBefore.Clone()
	_, err = mAfter.Exec(`
		CREATE TYPE public.price_range AS RANGE (SUBTYPE = int4);
		CREATE TABLE public.price_tbl (p public.price_multirange);
	`, nil)
	require.NoError(t, err)
	mMeta := ApplyDiffToMetadata(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}, mBefore, mAfter, Diff(mBefore, mAfter))
	pubMeta := findSchema(mMeta, "public")
	require.Len(t, pubMeta.RangeTypes, 1)
	require.Equal(t, "public.price_multirange", pubMeta.RangeTypes[0].MultirangeName)
	mSQL := GenerateMigration(mBefore, mAfter, Diff(mBefore, mAfter)).SQL()
	require.NotContains(t, mSQL, `MULTIRANGE_TYPE_NAME`)
	mApplied := mBefore.Clone()
	mExec, err := mApplied.Exec(mSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range mExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(mApplied, mAfter).Ranges)

	// 4. Cross-schema MULTIRANGE_TYPE_NAME is created in target schema, serialized, and migrated.
	xBefore := New()
	_, err = xBefore.Exec(`CREATE SCHEMA a; CREATE SCHEMA b;`, nil)
	require.NoError(t, err)
	xAfter := xBefore.Clone()
	_, err = xAfter.Exec(`
		CREATE TYPE a.r_cross AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = b.mr_cross);
		CREATE TABLE b.t_cross (c b.mr_cross);
	`, nil)
	require.NoError(t, err)
	xMeta := ApplyDiffToMetadata(&metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{{Name: "a"}, {Name: "b"}},
	}, xBefore, xAfter, Diff(xBefore, xAfter))
	aMeta := findSchema(xMeta, "a")
	require.Len(t, aMeta.RangeTypes, 1)
	require.Equal(t, "b.mr_cross", aMeta.RangeTypes[0].MultirangeName)
	xSQL := GenerateMigration(xBefore, xAfter, Diff(xBefore, xAfter)).SQL()
	require.Contains(t, xSQL, `MULTIRANGE_TYPE_NAME = "b"."mr_cross"`)
	xApplied := xBefore.Clone()
	xExec, err := xApplied.Exec(xSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range xExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(xApplied, xAfter).Ranges)
	require.Empty(t, Diff(xApplied, xAfter).Relations)
}

func TestRound30Fixes(t *testing.T) {
	// 1. Canonical range rebuilt via modified collation (no retained column) emits shell CREATE TYPE -> CREATE FUNCTION -> full CREATE TYPE AS RANGE.
	canonColFrom := New()
	_, err := canonColFrom.Exec(`
		CREATE COLLATION public.r_coll (LOCALE = 'en_US.utf8');
		CREATE TYPE public.canon_coll_r;
		CREATE FUNCTION public.canon_coll_fn(public.canon_coll_r) RETURNS public.canon_coll_r LANGUAGE C IMMUTABLE AS 'internal', 'canon_coll_fn';
		CREATE TYPE public.canon_coll_r AS RANGE (SUBTYPE = text, COLLATION = public.r_coll, CANONICAL = public.canon_coll_fn);
	`, nil)
	require.NoError(t, err)
	canonColTo := New()
	_, err = canonColTo.Exec(`
		CREATE COLLATION public.r_coll (LOCALE = 'C.utf8');
		CREATE TYPE public.canon_coll_r;
		CREATE FUNCTION public.canon_coll_fn(public.canon_coll_r) RETURNS public.canon_coll_r LANGUAGE C IMMUTABLE AS 'internal', 'canon_coll_fn';
		CREATE TYPE public.canon_coll_r AS RANGE (SUBTYPE = text, COLLATION = public.r_coll, CANONICAL = public.canon_coll_fn);
	`, nil)
	require.NoError(t, err)
	canonColSQL := GenerateMigration(canonColFrom, canonColTo, Diff(canonColFrom, canonColTo)).SQL()
	require.Contains(t, canonColSQL, `CREATE TYPE "public"."canon_coll_r";`)
	canonColApplied := canonColFrom.Clone()
	ccRes, err := canonColApplied.Exec(canonColSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range ccRes {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(canonColApplied, canonColTo).Ranges)
	require.Empty(t, Diff(canonColApplied, canonColTo).Functions)

	// 2. Quoted identifiers with dots in multirange schema/name are compared structurally in diffRanges.
	dotA := New()
	_, err = dotA.Exec(`
		CREATE SCHEMA "a.b";
		CREATE SCHEMA "a";
		CREATE TYPE "a".r_dot AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = "a.b"."c");
	`, nil)
	require.NoError(t, err)
	dotB := New()
	_, err = dotB.Exec(`
		CREATE SCHEMA "a.b";
		CREATE SCHEMA "a";
		CREATE TYPE "a".r_dot AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = "a"."b.c");
	`, nil)
	require.NoError(t, err)
	dotDiff := Diff(dotA, dotB)
	require.Len(t, dotDiff.Ranges, 1)
	require.Equal(t, DiffModify, dotDiff.Ranges[0].Action)

	// 3. Unqualified MULTIRANGE_TYPE_NAME with qualified range resolves via search_path, not range's schema.
	spBefore := New()
	_, err = spBefore.Exec(`CREATE SCHEMA range_ns; CREATE SCHEMA mr_ns; SET search_path = mr_ns, public;`, nil)
	require.NoError(t, err)
	spAfter := spBefore.Clone()
	_, err = spAfter.Exec(`
		CREATE TYPE range_ns.r_sp AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = mr_unqual);
	`, nil)
	require.NoError(t, err)
	require.Zero(t, resolveTypeOIDByName(spAfter, "range_ns", "mr_unqual"))
	require.NotZero(t, resolveTypeOIDByName(spAfter, "mr_ns", "mr_unqual"))
	spSQL := GenerateMigration(spBefore, spAfter, Diff(spBefore, spAfter)).SQL()
	require.Contains(t, spSQL, `MULTIRANGE_TYPE_NAME = "mr_ns"."mr_unqual"`)
	spApplied := spBefore.Clone()
	spRes, err := spApplied.Exec(spSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range spRes {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(spApplied, spAfter).Ranges)

	// 4. Auto-generated multirange names respect PostgreSQL's 63-byte identifier limit.
	longNoRange := strings.Repeat("x", 63)
	mrLong := makeMultirangeTypeName(longNoRange)
	require.Equal(t, 63, len(mrLong))
	require.Equal(t, strings.Repeat("x", 52)+"_multirange", mrLong)

	longWithRange := strings.Repeat("y", 58) + "range"
	mrWithRange := makeMultirangeTypeName(longWithRange)
	require.Equal(t, 63, len(mrWithRange))
	require.Equal(t, strings.Repeat("y", 58)+"multi", mrWithRange)

	mbName := strings.Repeat("\u00e9", 26) + "x" // 53 bytes
	mrMB := makeMultirangeTypeName(mbName)
	require.LessOrEqual(t, len(mrMB), 63)
	require.True(t, strings.HasSuffix(mrMB, "_multirange"))

	// 5. Foreign server FDW change refuses DROP SERVER CASCADE and preserves foreign tables.
	srvFdwFrom := New()
	_, err = srvFdwFrom.Exec(`
		CREATE FOREIGN DATA WRAPPER w1;
		CREATE FOREIGN DATA WRAPPER w2;
		CREATE SERVER srv_fdw FOREIGN DATA WRAPPER w1;
		CREATE FOREIGN TABLE public.ft_keep (id integer) SERVER srv_fdw;
	`, nil)
	require.NoError(t, err)
	srvFdwTo := New()
	_, err = srvFdwTo.Exec(`
		CREATE FOREIGN DATA WRAPPER w1;
		CREATE FOREIGN DATA WRAPPER w2;
		CREATE SERVER srv_fdw FOREIGN DATA WRAPPER w2;
		CREATE FOREIGN TABLE public.ft_keep (id integer) SERVER srv_fdw;
	`, nil)
	require.NoError(t, err)
	srvPlan := GenerateMigration(srvFdwFrom, srvFdwTo, Diff(srvFdwFrom, srvFdwTo))
	require.NotEmpty(t, srvPlan.Warnings())
	require.NotContains(t, srvPlan.SQL(), "DROP SERVER")
	require.NotContains(t, srvPlan.SQL(), "CREATE FOREIGN TABLE")
}

func TestRound31Fixes(t *testing.T) {
	// 1. Grants on retained functions rebuilt with ranges are restored cleanly.
	grFrom := New()
	_, err := grFrom.Exec(`
		CREATE ROLE app_user;
		CREATE TYPE public.r_grant AS RANGE (SUBTYPE = int4);
		CREATE FUNCTION public.use_r(x public.r_grant) RETURNS int4 LANGUAGE sql IMMUTABLE AS 'SELECT 1';
		GRANT EXECUTE ON FUNCTION public.use_r(public.r_grant) TO app_user;
		GRANT USAGE ON TYPE public.r_grant TO app_user;
	`, nil)
	require.NoError(t, err)
	grTo := New()
	_, err = grTo.Exec(`
		CREATE ROLE app_user;
		CREATE TYPE public.r_grant AS RANGE (SUBTYPE = int4, SUBTYPE_OPCLASS = int4_ops);
		CREATE FUNCTION public.use_r(x public.r_grant) RETURNS int4 LANGUAGE sql IMMUTABLE AS 'SELECT 1';
		GRANT EXECUTE ON FUNCTION public.use_r(public.r_grant) TO app_user;
		GRANT USAGE ON TYPE public.r_grant TO app_user;
	`, nil)
	require.NoError(t, err)
	grDiff := Diff(grFrom, grTo)
	require.NotEmpty(t, grDiff.Ranges)
	require.Empty(t, grDiff.Grants)
	grSQL := GenerateMigration(grFrom, grTo, grDiff).SQL()
	require.Contains(t, grSQL, `GRANT USAGE ON TYPE "public"."r_grant" TO "app_user"`)
	require.Contains(t, grSQL, `GRANT EXECUTE ON FUNCTION public.use_r(r_grant) TO "app_user"`)
	grApplied := grFrom.Clone()
	grExec, err := grApplied.Exec(grSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range grExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(grApplied, grTo).Grants)
	require.Empty(t, Diff(grApplied, grTo).Functions)
	require.Empty(t, Diff(grApplied, grTo).Ranges)

	// 2. Quoted identifiers named "none" are preserved instead of being treated as NONE keyword.
	noneBefore := New()
	_, err = noneBefore.Exec(`
		CREATE FUNCTION public."none"() RETURNS fdw_handler LANGUAGE internal AS 'textin';
		CREATE FUNCTION public.none_canon(public.r_qnone) RETURNS public.r_qnone LANGUAGE C IMMUTABLE AS 'internal', 'nc';
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	noneAfter := New()
	_, err = noneAfter.Exec(`
		CREATE FUNCTION public."none"() RETURNS fdw_handler LANGUAGE internal AS 'textin';
		CREATE FOREIGN DATA WRAPPER fdw_qnone HANDLER public."none";
		CREATE TYPE public.r_qnone AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = "none");
		CREATE TYPE public.r_unqnone AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = none);
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.NotNil(t, noneAfter.fdws["fdw_qnone"])
	require.Contains(t, noneAfter.fdws["fdw_qnone"].Handler, "none")
	require.NotZero(t, resolveTypeOIDByName(noneAfter, "public", "none"))
	require.Zero(t, resolveTypeOIDByName(noneAfter, "public", "r_unqnone_none"))
	noneSQL := GenerateMigration(noneBefore, noneAfter, Diff(noneBefore, noneAfter)).SQL()
	require.Contains(t, noneSQL, `MULTIRANGE_TYPE_NAME = "public"."none"`)
	noneApplied := noneBefore.Clone()
	noneExec, err := noneApplied.Exec(noneSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range noneExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(noneApplied, noneAfter).ForeignDataWrappers)
	require.Empty(t, Diff(noneApplied, noneAfter).Ranges)

	// 3. Dropping a table column with column-level grants suppresses REVOKE on the dropped column.
	colFrom := New()
	_, err = colFrom.Exec(`
		CREATE ROLE reader;
		CREATE TABLE public.t_cgrant (id int, a int, b text);
		GRANT SELECT (a) ON public.t_cgrant TO reader;
		GRANT SELECT (a, b) ON public.t_cgrant TO reader;
		GRANT UPDATE (b) ON public.t_cgrant TO reader;
	`, nil)
	require.NoError(t, err)
	colTo := colFrom.Clone()
	_, err = colTo.Exec(`ALTER TABLE public.t_cgrant DROP COLUMN a;`, nil)
	require.NoError(t, err)
	colSQL := GenerateMigration(colFrom, colTo, Diff(colFrom, colTo)).SQL()
	require.NotContains(t, colSQL, `REVOKE SELECT ("a")`)
	require.NotContains(t, colSQL, `REVOKE SELECT ("a", "b")`)
	colApplied := colFrom.Clone()
	colExec, err := colApplied.Exec(colSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range colExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(colApplied, colTo).Relations)
	require.Empty(t, Diff(colApplied, colTo).Grants)

	// 4. Renaming a foreign table preserves dependent views and triggers.
	ftRenFrom := New()
	_, err = ftRenFrom.Exec(`
		CREATE ROLE ft_reader;
		CREATE FOREIGN DATA WRAPPER w_ft;
		CREATE SERVER srv_ft FOREIGN DATA WRAPPER w_ft;
		CREATE FUNCTION public.ft_ren_trig() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE FOREIGN TABLE public.ft_old (id int, note text) SERVER srv_ft;
		CREATE TRIGGER tr_ft_ren BEFORE INSERT ON public.ft_old FOR EACH ROW EXECUTE FUNCTION public.ft_ren_trig();
		CREATE VIEW public.v_ft_const AS SELECT 1 AS ok FROM public.ft_old;
		CREATE VIEW public.v_ft_col AS SELECT id, note FROM public.ft_old;
		GRANT SELECT ON public.ft_old TO ft_reader;
	`, nil)
	require.NoError(t, err)
	ftRenTo := ftRenFrom.Clone()
	_, err = ftRenTo.Exec(`ALTER FOREIGN TABLE public.ft_old RENAME TO ft_new;`, nil)
	require.NoError(t, err)
	ftRenSQL := GenerateMigration(ftRenFrom, ftRenTo, Diff(ftRenFrom, ftRenTo)).SQL()
	ftRenApplied := ftRenFrom.Clone()
	ftRenExec, err := ftRenApplied.Exec(ftRenSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range ftRenExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.NotNil(t, ftRenApplied.GetRelation("public", "v_ft_const"))
	require.NotNil(t, ftRenApplied.GetRelation("public", "v_ft_col"))
	require.Empty(t, Diff(ftRenApplied, ftRenTo).ForeignTables)
	require.Empty(t, Diff(ftRenApplied, ftRenTo).Relations)
	require.Empty(t, Diff(ftRenApplied, ftRenTo).Grants)
}

func TestRound32Fixes(t *testing.T) {
	// 1. FDW rename emits ALTER FOREIGN DATA WRAPPER RENAME TO and preserves dependent servers and foreign tables.
	fdwFrom := New()
	_, err := fdwFrom.Exec(`
		CREATE FOREIGN DATA WRAPPER fdw_old OPTIONS (debug 'false');
		CREATE SERVER srv_keep FOREIGN DATA WRAPPER fdw_old;
		CREATE FOREIGN TABLE public.ft_keep (id int) SERVER srv_keep;
	`, nil)
	require.NoError(t, err)
	fdwTo := fdwFrom.Clone()
	_, err = fdwTo.Exec(`
		ALTER FOREIGN DATA WRAPPER fdw_old RENAME TO fdw_new;
		ALTER FOREIGN DATA WRAPPER fdw_new OPTIONS (SET debug 'true');
	`, nil)
	require.NoError(t, err)
	fdwSQL := GenerateMigration(fdwFrom, fdwTo, Diff(fdwFrom, fdwTo)).SQL()
	require.Contains(t, fdwSQL, `ALTER FOREIGN DATA WRAPPER fdw_old RENAME TO fdw_new`)
	require.NotContains(t, fdwSQL, `DROP FOREIGN DATA WRAPPER`)
	fdwApplied := fdwFrom.Clone()
	fdwExec, err := fdwApplied.Exec(fdwSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range fdwExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(fdwApplied, fdwTo).ForeignDataWrappers)
	require.Empty(t, Diff(fdwApplied, fdwTo).ForeignServers)
	require.Empty(t, Diff(fdwApplied, fdwTo).ForeignTables)

	// 2. Dropping an old SUPPORT routine after switching SUPPORT on a retained function defers DROP FUNCTION until after CREATE OR REPLACE.
	supFrom := New()
	_, err = supFrom.Exec(`
		CREATE FUNCTION public.sup_old(internal) RETURNS internal LANGUAGE C STRICT AS 'internal', 's1';
		CREATE FUNCTION public.sup_new(internal) RETURNS internal LANGUAGE C STRICT AS 'internal', 's2';
		CREATE FUNCTION public.fn_with_sup(x int) RETURNS int LANGUAGE sql IMMUTABLE SUPPORT public.sup_old AS 'SELECT x';
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	supTo := supFrom.Clone()
	_, err = supTo.Exec(`
		CREATE OR REPLACE FUNCTION public.fn_with_sup(x int) RETURNS int LANGUAGE sql IMMUTABLE SUPPORT public.sup_new AS 'SELECT x + 1';
		DROP FUNCTION public.sup_old(internal);
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	supSQL := GenerateMigration(supFrom, supTo, Diff(supFrom, supTo)).SQL()
	supReplaceIdx := strings.Index(supSQL, `CREATE OR REPLACE FUNCTION "public"."fn_with_sup"`)
	supDropIdx := strings.Index(supSQL, `DROP FUNCTION "public"."sup_old"`)
	require.NotEqual(t, -1, supReplaceIdx)
	require.NotEqual(t, -1, supDropIdx)
	require.Less(t, supReplaceIdx, supDropIdx)
	supApplied := supFrom.Clone()
	supExec, err := supApplied.Exec(supSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range supExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(supApplied, supTo).Functions)

	// 3. Extension-membership (SkipDump) transitions on triggers do not emit spurious CREATE/DROP TRIGGER.
	makeTrgMeta := func(skipDump bool) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "trg_fn",
					Signature:  "trg_fn()",
					Definition: `CREATE FUNCTION public.trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$`,
				}},
				Tables: []*metadatapb.TableMetadata{{
					Name:    "t_ext_trg",
					Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
					Triggers: []*metadatapb.TriggerMetadata{{
						Name:     "tr_ext",
						Body:     `CREATE TRIGGER tr_ext BEFORE INSERT ON public.t_ext_trg FOR EACH ROW EXECUTE FUNCTION public.trg_fn()`,
						SkipDump: skipDump,
					}},
				}},
			}},
		}
	}
	sqlDetach, err := MetadataMigration(context.Background(), makeTrgMeta(true), makeTrgMeta(false), MetadataMigrationOptions{})
	require.NoError(t, err)
	require.NotContains(t, sqlDetach, `CREATE TRIGGER`)
	require.NotContains(t, sqlDetach, `DROP TRIGGER`)
	sqlAttach, err := MetadataMigration(context.Background(), makeTrgMeta(false), makeTrgMeta(true), MetadataMigrationOptions{})
	require.NoError(t, err)
	require.NotContains(t, sqlAttach, `CREATE TRIGGER`)
	require.NotContains(t, sqlAttach, `DROP TRIGGER`)

	// 4. Removing a parameter default from a retained function rebuilds dependent unchanged views.
	defFrom := New()
	_, err = defFrom.Exec(`
		CREATE FUNCTION public.fn_def(a int DEFAULT 0) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT a + 1';
		CREATE VIEW public.v_fn_def AS SELECT public.fn_def(1) AS val;
	`, nil)
	require.NoError(t, err)
	defTo := New()
	_, err = defTo.Exec(`
		CREATE FUNCTION public.fn_def(a int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT a + 1';
		CREATE VIEW public.v_fn_def AS SELECT public.fn_def(1) AS val;
	`, nil)
	require.NoError(t, err)
	defSQL := GenerateMigration(defFrom, defTo, Diff(defFrom, defTo)).SQL()
	defApplied := defFrom.Clone()
	defExec, err := defApplied.Exec(defSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range defExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(defApplied, defTo).Functions)
	require.Empty(t, Diff(defApplied, defTo).Relations)
}

func TestRound33Fixes(t *testing.T) {
	// 1. Renaming a foreign server emits ALTER SERVER ... RENAME TO and preserves foreign tables and views.
	srvFrom := New()
	_, err := srvFrom.Exec(`
		CREATE FOREIGN DATA WRAPPER w_srv;
		CREATE SERVER srv_old FOREIGN DATA WRAPPER w_srv OPTIONS (host 'h1');
		CREATE FOREIGN TABLE public.ft_srv (id int) SERVER srv_old;
		CREATE VIEW public.v_ft_srv AS SELECT id FROM public.ft_srv;
	`, nil)
	require.NoError(t, err)
	srvTo := srvFrom.Clone()
	_, err = srvTo.Exec(`
		ALTER SERVER srv_old RENAME TO srv_new;
		ALTER SERVER srv_new OPTIONS (SET host 'h2');
	`, nil)
	require.NoError(t, err)
	srvSQL := GenerateMigration(srvFrom, srvTo, Diff(srvFrom, srvTo)).SQL()
	require.Contains(t, srvSQL, `ALTER SERVER srv_old RENAME TO srv_new`)
	require.NotContains(t, srvSQL, `DROP SERVER`)
	require.NotContains(t, srvSQL, `DROP FOREIGN TABLE`)
	srvApplied := srvFrom.Clone()
	srvExec, err := srvApplied.Exec(srvSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range srvExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(srvApplied, srvTo).ForeignServers)
	require.Empty(t, Diff(srvApplied, srvTo).ForeignTables)
	require.Empty(t, Diff(srvApplied, srvTo).Relations)

	// 2. Removing a function parameter default rebuilds all retained dependents (defaults, generated columns, CHECK, indexes, policies, triggers, views).
	fnFrom := New()
	_, err = fnFrom.Exec(`
		CREATE FUNCTION public.fn_all_dep(a int DEFAULT 0) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT a + 1';
		CREATE FUNCTION public.fn_all_dep_trig() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE TABLE public.t_fn_dep (
			id int,
			x int DEFAULT public.fn_all_dep(1),
			g int GENERATED ALWAYS AS (public.fn_all_dep(id)) STORED,
			CONSTRAINT chk_fn_dep CHECK (public.fn_all_dep(id) > 0)
		);
		CREATE INDEX idx_fn_expr ON public.t_fn_dep ((public.fn_all_dep(id)));
		CREATE INDEX idx_fn_part ON public.t_fn_dep (id) WHERE public.fn_all_dep(id) > 0;
		CREATE POLICY pol_fn_dep ON public.t_fn_dep USING (public.fn_all_dep(id) > 0);
		CREATE TRIGGER tr_fn_dep BEFORE INSERT ON public.t_fn_dep FOR EACH ROW WHEN (public.fn_all_dep(NEW.id) > 0) EXECUTE FUNCTION public.fn_all_dep_trig();
		CREATE VIEW public.v_fn_dep AS SELECT public.fn_all_dep(id) AS val FROM public.t_fn_dep;
	`, nil)
	require.NoError(t, err)
	fnTo := New()
	_, err = fnTo.Exec(`
		CREATE FUNCTION public.fn_all_dep(a int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT a + 1';
		CREATE FUNCTION public.fn_all_dep_trig() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE TABLE public.t_fn_dep (
			id int,
			x int DEFAULT public.fn_all_dep(1),
			g int GENERATED ALWAYS AS (public.fn_all_dep(id)) STORED,
			CONSTRAINT chk_fn_dep CHECK (public.fn_all_dep(id) > 0)
		);
		CREATE INDEX idx_fn_expr ON public.t_fn_dep ((public.fn_all_dep(id)));
		CREATE INDEX idx_fn_part ON public.t_fn_dep (id) WHERE public.fn_all_dep(id) > 0;
		CREATE POLICY pol_fn_dep ON public.t_fn_dep USING (public.fn_all_dep(id) > 0);
		CREATE TRIGGER tr_fn_dep BEFORE INSERT ON public.t_fn_dep FOR EACH ROW WHEN (public.fn_all_dep(NEW.id) > 0) EXECUTE FUNCTION public.fn_all_dep_trig();
		CREATE VIEW public.v_fn_dep AS SELECT public.fn_all_dep(id) AS val FROM public.t_fn_dep;
	`, nil)
	require.NoError(t, err)
	fnSQL := GenerateMigration(fnFrom, fnTo, Diff(fnFrom, fnTo)).SQL()
	fnApplied := fnFrom.Clone()
	fnExec, err := fnApplied.Exec(fnSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range fnExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(fnApplied, fnTo).Functions)
	require.Empty(t, Diff(fnApplied, fnTo).Relations)

	// 3. Foreign tables reject PK, UNIQUE, FK, and EXCLUDE constraints while permitting CHECK and NOT NULL.
	ftCat := New()
	_, err = ftCat.Exec(`
		CREATE FOREIGN DATA WRAPPER w_con;
		CREATE SERVER srv_con FOREIGN DATA WRAPPER w_con;
		CREATE TABLE public.ref_tbl (id int PRIMARY KEY);
	`, nil)
	require.NoError(t, err)
	for _, badSQL := range []string{
		`CREATE FOREIGN TABLE public.ft_pk (id int PRIMARY KEY) SERVER srv_con;`,
		`CREATE FOREIGN TABLE public.ft_uniq (id int, UNIQUE (id)) SERVER srv_con;`,
		`CREATE FOREIGN TABLE public.ft_fk (id int REFERENCES public.ref_tbl(id)) SERVER srv_con;`,
		`CREATE FOREIGN TABLE public.ft_excl (id int, EXCLUDE USING btree (id WITH =)) SERVER srv_con;`,
	} {
		res, execErr := ftCat.Exec(badSQL, nil)
		if execErr == nil {
			require.NotEmpty(t, res)
			execErr = res[0].Error
		}
		require.Error(t, execErr, badSQL)
	}
	resOK, err := ftCat.Exec(`
		CREATE FOREIGN TABLE public.ft_ok (id int NOT NULL, CONSTRAINT chk_ft CHECK (id > 0)) SERVER srv_con;
		ALTER FOREIGN TABLE public.ft_ok ADD CONSTRAINT ft_ok_pk PRIMARY KEY (id);
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Len(t, resOK, 2)
	require.NoError(t, resOK[0].Error)
	require.Error(t, resOK[1].Error)
}

func TestRound34Fixes(t *testing.T) {
	// 1. Companion/multirange array collision aborts DefineRange before registering the range type.
	mrCat := New()
	_, err := mrCat.Exec(`
		CREATE SCHEMA a;
		CREATE SCHEMA b;
		CREATE TYPE b._mr AS ENUM ('x');
	`, nil)
	require.NoError(t, err)
	mrRes, err := mrCat.Exec(`CREATE TYPE a.r AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = b.mr);`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Len(t, mrRes, 1)
	require.Error(t, mrRes[0].Error)
	require.Zero(t, resolveTypeOIDByName(mrCat, "a", "r"))
	require.Zero(t, resolveTypeOIDByName(mrCat, "a", "_r"))

	// 2. Synced metadata with unqualified SUPPORT resolves through configured search_path.
	supMeta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "app",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "f",
					Signature:  "f(integer)",
					Definition: `CREATE FUNCTION app.f(x integer) RETURNS integer LANGUAGE sql IMMUTABLE SUPPORT helper AS 'SELECT x'`,
				}},
			},
			{
				Name: "custom",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "helper",
					Signature:  "helper(internal)",
					Definition: `CREATE FUNCTION custom.helper(internal) RETURNS internal LANGUAGE internal STRICT AS 'int4pl'`,
				}},
			},
		},
	}
	supCat := New()
	supCat.SetSearchPath([]string{"custom", "public"})
	supReport, err := supCat.LoadMetadata(context.Background(), supMeta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, supReport.Degraded)
	supMigSQL, err := MetadataMigration(context.Background(), &metadatapb.DatabaseSchemaMetadata{}, supMeta, MetadataMigrationOptions{
		TargetSearchPath: []string{"custom", "public"},
	})
	require.NoError(t, err)
	require.Contains(t, supMigSQL, `CREATE FUNCTION app.f`)
	require.Contains(t, supMigSQL, `CREATE FUNCTION "custom"."helper"`)
	require.Less(t, strings.Index(supMigSQL, `"custom"."helper"`), strings.Index(supMigSQL, `CREATE FUNCTION app.f`))

	// 3. SQL-standard function bodies (RETURN / BEGIN ATOMIC) record relation dependencies directly from AST.
	sqlBodyFrom := New()
	sqlBodyTo := New()
	_, err = sqlBodyTo.Exec(`
		CREATE TABLE public.z_src (id int);
		CREATE FUNCTION public.a_ret() RETURNS int LANGUAGE sql RETURN (SELECT id FROM public.z_src LIMIT 1);
		CREATE FUNCTION public.a_blk() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN (SELECT id FROM public.z_src LIMIT 1); END;
	`, nil)
	require.NoError(t, err)
	sqlBodyPlan := GenerateMigration(sqlBodyFrom, sqlBodyTo, Diff(sqlBodyFrom, sqlBodyTo))
	sqlBodySQL := sqlBodyPlan.SQL()
	tblIdx := strings.Index(sqlBodySQL, `CREATE TABLE "public"."z_src"`)
	retIdx := strings.Index(sqlBodySQL, `CREATE FUNCTION "public"."a_ret"`)
	blkIdx := strings.Index(sqlBodySQL, `CREATE FUNCTION "public"."a_blk"`)
	require.NotEqual(t, -1, tblIdx)
	require.Less(t, tblIdx, retIdx)
	require.Less(t, tblIdx, blkIdx)
	sqlBodyApplied := sqlBodyFrom.Clone()
	sqlBodyExec, err := sqlBodyApplied.Exec(sqlBodySQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range sqlBodyExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(sqlBodyApplied, sqlBodyTo).Functions)
	require.Empty(t, Diff(sqlBodyApplied, sqlBodyTo).Relations)

	// 4. Renaming / moving schema of a range SUBTYPE_OPCLASS refreshes stored SubtypeOpclass references.
	opcFrom := New()
	_, err = opcFrom.Exec(`
		CREATE SCHEMA op_dst;
		CREATE OPERATOR FAMILY public.my_fam USING btree;
		CREATE OPERATOR CLASS public.my_opc FOR TYPE int4 USING btree FAMILY public.my_fam AS
			OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >, FUNCTION 1 btint4cmp(int4, int4);
		CREATE TYPE public.r_opc AS RANGE (SUBTYPE = int4, SUBTYPE_OPCLASS = public.my_opc);
	`, nil)
	require.NoError(t, err)
	opcTo := opcFrom.Clone()
	_, err = opcTo.Exec(`
		ALTER OPERATOR CLASS public.my_opc USING btree RENAME TO my_opc2;
		ALTER OPERATOR CLASS public.my_opc2 USING btree SET SCHEMA op_dst;
	`, nil)
	require.NoError(t, err)
	rtOID := resolveTypeOIDByName(opcTo, "public", "r_opc")
	require.NotZero(t, rtOID)
	require.Equal(t, `op_dst.my_opc2`, opcTo.RangeInfo(rtOID).SubtypeOpclass)
}

func TestRound35Fixes(t *testing.T) {
	// 1. CREATE OR REPLACE FUNCTION refreshes default-expression dependencies.
	defCat := New()
	_, err := defCat.Exec(`
		CREATE FUNCTION public.old_helper() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;
		CREATE FUNCTION public.new_helper() RETURNS integer LANGUAGE sql AS $$ SELECT 2 $$;
		CREATE FUNCTION public.f_def(x integer DEFAULT public.old_helper()) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;
		CREATE OR REPLACE FUNCTION public.f_def(x integer DEFAULT public.new_helper()) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;
	`, nil)
	require.NoError(t, err)

	// Dropping old_helper must now succeed without CASCADE, while new_helper is blocked.
	dropNewRes, err := defCat.Exec(`DROP FUNCTION public.new_helper()`, nil)
	require.NoError(t, err)
	require.Len(t, dropNewRes, 1)
	require.Error(t, dropNewRes[0].Error)

	dropOldRes, err := defCat.Exec(`DROP FUNCTION public.old_helper()`, nil)
	require.NoError(t, err)
	require.Len(t, dropOldRes, 1)
	require.NoError(t, dropOldRes[0].Error)

	// Migration from old_helper default to new_helper default executes cleanly (replaces f_def before dropping old_helper).
	migDefFrom := New()
	_, err = migDefFrom.Exec(`
		CREATE FUNCTION public.old_helper() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;
		CREATE FUNCTION public.f_def(x integer DEFAULT public.old_helper()) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;
	`, nil)
	require.NoError(t, err)
	migDefTo := New()
	_, err = migDefTo.Exec(`
		CREATE FUNCTION public.new_helper() RETURNS integer LANGUAGE sql AS $$ SELECT 2 $$;
		CREATE FUNCTION public.f_def(x integer DEFAULT public.new_helper()) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;
	`, nil)
	require.NoError(t, err)
	defSQL := GenerateMigration(migDefFrom, migDefTo, Diff(migDefFrom, migDefTo)).SQL()
	defApplied := migDefFrom.Clone()
	execRes, err := defApplied.Exec(defSQL, nil)
	require.NoError(t, err)
	for _, r := range execRes {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(defApplied, migDefTo).Functions)

	// 2. WINDOW attribute sets routine kind 'w' on both UserProc and BuiltinProc and enforces kind on replace.
	winCat := New()
	_, err = winCat.Exec(`
		CREATE FUNCTION public.my_win() RETURNS integer LANGUAGE internal WINDOW AS 'window_row_number';
	`, nil)
	require.NoError(t, err)
	var winProc *UserProc
	for _, up := range winCat.userProcs {
		if up.Name == "my_win" {
			winProc = up
			break
		}
	}
	require.NotNil(t, winProc)
	require.Equal(t, byte('w'), winProc.Kind)
	require.Equal(t, byte('w'), winCat.procByOID[winProc.OID].Kind)
	var foundRow bool
	for _, r := range winCat.QueryPgProc("public") {
		if r.ProName == "my_win" {
			foundRow = true
			require.Equal(t, byte('w'), r.ProKind)
		}
	}
	require.True(t, foundRow)

	okView, err := winCat.Exec(`CREATE VIEW public.v_win_ok AS SELECT public.my_win() OVER () AS w`, nil)
	require.NoError(t, err)
	require.NoError(t, okView[0].Error)

	badView, err := winCat.Exec(`CREATE VIEW public.v_win_bad AS SELECT public.my_win() AS w`, nil)
	require.NoError(t, err)
	require.Error(t, badView[0].Error)

	badReplace, err := winCat.Exec(`CREATE OR REPLACE FUNCTION public.my_win() RETURNS integer LANGUAGE internal AS 'window_row_number'`, nil)
	require.NoError(t, err)
	require.Error(t, badReplace[0].Error)

	okReplace, err := winCat.Exec(`CREATE OR REPLACE FUNCTION public.my_win() RETURNS integer LANGUAGE internal WINDOW AS 'window_row_number'`, nil)
	require.NoError(t, err)
	require.NoError(t, okReplace[0].Error)

	// 3. Multiple inheritance parent replacement preserves retained later parents without NO INHERIT/INHERIT churn.
	inhFrom := New()
	_, err = inhFrom.Exec(`
		CREATE TABLE public.p1 (a int);
		CREATE TABLE public.p2 (a int);
		CREATE TABLE public.kept (k int);
		CREATE TABLE public.child () INHERITS (public.p1, public.kept);
	`, nil)
	require.NoError(t, err)
	inhTo := New()
	_, err = inhTo.Exec(`
		CREATE TABLE public.p1 (a int);
		CREATE TABLE public.p2 (a int);
		CREATE TABLE public.kept (k int);
		CREATE TABLE public.child () INHERITS (public.p2, public.kept);
	`, nil)
	require.NoError(t, err)
	inhSQL := GenerateMigration(inhFrom, inhTo, Diff(inhFrom, inhTo)).SQL()
	require.Contains(t, inhSQL, `NO INHERIT`)
	require.Contains(t, inhSQL, `INHERIT`)
	require.NotContains(t, inhSQL, `kept`)

	inhApplied := inhFrom.Clone()
	inhExec, err := inhApplied.Exec(inhSQL, nil)
	require.NoError(t, err)
	for _, r := range inhExec {
		require.NoError(t, r.Error, r.SQL)
	}
	// Columns from kept remain non-local with InhCount=1, so dropping kept.k removes child.k too.
	childRel := inhApplied.GetRelation("public", "child")
	require.NotNil(t, childRel)
	kCol := childRel.Columns[childRel.colByName["k"]]
	require.False(t, kCol.IsLocal)
	require.Equal(t, 1, kCol.InhCount)
	dropK, err := inhApplied.Exec(`ALTER TABLE public.kept DROP COLUMN k`, nil)
	require.NoError(t, err)
	require.NoError(t, dropK[0].Error)
	for _, col := range childRel.Columns {
		require.NotEqual(t, "k", col.Name)
	}
}

func TestRound36Fixes(t *testing.T) {
	// 1. SQL-standard body calling another newly added function (alphabetically later) orders callee first.
	callFrom := New()
	callTo := New()
	_, err := callTo.Exec(`
		CREATE FUNCTION public.z_callee() RETURNS integer LANGUAGE sql RETURN 42;
		CREATE FUNCTION public.a_caller() RETURNS integer LANGUAGE sql RETURN public.z_callee();
	`, nil)
	require.NoError(t, err)
	callSQL := GenerateMigration(callFrom, callTo, Diff(callFrom, callTo)).SQL()
	require.Contains(t, callSQL, `CREATE FUNCTION "public"."z_callee"`)
	require.Contains(t, callSQL, `CREATE FUNCTION "public"."a_caller"`)
	require.Less(t, strings.Index(callSQL, `"public"."z_callee"`), strings.Index(callSQL, `"public"."a_caller"`))

	callApplied := callFrom.Clone()
	callExec, err := callApplied.Exec(callSQL, nil)
	require.NoError(t, err)
	for _, r := range callExec {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, Diff(callApplied, callTo).Functions)

	// Unqualified call in a custom search_path schema also loads cleanly via LoadMetadata.
	unqMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "custom",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "a_caller",
					Signature:  "custom.a_caller()",
					Definition: `CREATE FUNCTION custom.a_caller() RETURNS integer LANGUAGE sql RETURN z_callee()`,
				},
				{
					Name:       "z_callee",
					Signature:  "custom.z_callee()",
					Definition: `CREATE FUNCTION custom.z_callee() RETURNS integer LANGUAGE sql RETURN 42`,
				},
			},
		}},
	}
	unqCat := New()
	unqReport, err := unqCat.LoadMetadata(context.Background(), unqMeta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, unqReport.Degraded)

	// 2. ApplyDiffToMetadata serializes EXCLUDE constraint expressions when adding EXCLUDE constraints.
	exclBefore := New()
	_, err = exclBefore.Exec(`CREATE TABLE public.booking (id int, r int4range);`, nil)
	require.NoError(t, err)
	exclOrigMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "booking",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Position: 1, Type: "integer", Nullable: true},
					{Name: "r", Position: 2, Type: "int4range", Nullable: true},
				},
			}},
		}},
	}
	exclAfter := exclBefore.Clone()
	_, err = exclAfter.Exec(`ALTER TABLE public.booking ADD CONSTRAINT booking_r_excl EXCLUDE USING gist (r WITH &&);`, nil)
	require.NoError(t, err)
	exclUpdated := ApplyDiffToMetadata(exclOrigMeta, exclBefore, exclAfter, Diff(exclBefore, exclAfter))
	require.Len(t, exclUpdated.Schemas[0].Tables[0].ExcludeConstraints, 1)
	exclProto := exclUpdated.Schemas[0].Tables[0].ExcludeConstraints[0]
	require.Equal(t, "booking_r_excl", exclProto.Name)
	require.Equal(t, `EXCLUDE USING gist (r WITH &&)`, exclProto.Expression)

	// Round-trip loading the updated metadata recreates the EXCLUDE constraint identically.
	reloaded := New()
	reloadReport, err := reloaded.LoadMetadata(context.Background(), exclUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, reloadReport.Degraded)
	require.Empty(t, Diff(reloaded, exclAfter).Relations)
}

func TestRound37Fixes(t *testing.T) {
	// 1. A SQL-standard body calling a quoted mixed-case routine orders the callee first.
	quotedFrom := New()
	quotedTo := New()
	_, err := quotedTo.Exec(`
		CREATE FUNCTION public."Z"() RETURNS integer LANGUAGE sql RETURN 42;
		CREATE FUNCTION public."A"() RETURNS integer LANGUAGE sql RETURN public."Z"();
	`, nil)
	require.NoError(t, err)
	quotedSQL := GenerateMigration(quotedFrom, quotedTo, Diff(quotedFrom, quotedTo)).SQL()
	require.Less(t, strings.Index(quotedSQL, `"public"."Z"`), strings.Index(quotedSQL, `"public"."A"`), quotedSQL)
	require.Error(t, round37ExecErr(quotedTo, `DROP FUNCTION public."Z"();`))

	// 2. A cast in a SQL-standard body depends on the type it names, as PostgreSQL records.
	castCat := New()
	_, err = castCat.Exec(`
		CREATE SCHEMA app;
		CREATE TYPE app.my_type AS (a int);
		CREATE FUNCTION public.uses_cast() RETURNS boolean LANGUAGE sql RETURN NULL::app.my_type IS NULL;
		CREATE FUNCTION public.uses_array_cast() RETURNS boolean LANGUAGE sql RETURN NULL::app.my_type[] IS NULL;
	`, nil)
	require.NoError(t, err)
	require.Error(t, round37ExecErr(castCat, `DROP TYPE app.my_type;`))
	require.NoError(t, round37ExecErr(castCat, `DROP FUNCTION public.uses_cast();`))
	require.Error(t, round37ExecErr(castCat, `DROP TYPE app.my_type;`))

	// 3. ApplyDiffToMetadata keeps INCLUDE columns of added PRIMARY KEY and UNIQUE constraints.
	inclBefore := New()
	_, err = inclBefore.Exec(`CREATE TABLE public.cov (id int, n int, m int);`, nil)
	require.NoError(t, err)
	inclOrigMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "cov",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Position: 1, Type: "integer", Nullable: true},
					{Name: "n", Position: 2, Type: "integer", Nullable: true},
					{Name: "m", Position: 3, Type: "integer", Nullable: true},
				},
			}},
		}},
	}
	inclAfter := inclBefore.Clone()
	_, err = inclAfter.Exec(`
		ALTER TABLE public.cov ADD CONSTRAINT cov_pk PRIMARY KEY (id) INCLUDE (n);
		ALTER TABLE public.cov ADD CONSTRAINT cov_uq UNIQUE (m) INCLUDE (n, id);
	`, nil)
	require.NoError(t, err)
	inclUpdated := ApplyDiffToMetadata(inclOrigMeta, inclBefore, inclAfter, Diff(inclBefore, inclAfter))
	defs := map[string]string{}
	for _, idx := range inclUpdated.Schemas[0].Tables[0].Indexes {
		defs[idx.Name] = idx.Definition
	}
	require.Contains(t, defs["cov_pk"], "INCLUDE (n)")
	require.Contains(t, defs["cov_uq"], "INCLUDE (n, id)")
	inclReloaded := New()
	inclReport, err := inclReloaded.LoadMetadata(context.Background(), inclUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, inclReport.Degraded)
	require.True(t, Diff(inclReloaded, inclAfter).IsEmpty())
	covered := 0
	for _, idx := range inclReloaded.IndexesOf(inclReloaded.GetRelation("public", "cov").OID) {
		require.Less(t, idx.NKeyColumns, len(idx.Columns), idx.Name)
		covered++
	}
	require.Equal(t, 2, covered)

	// 4. A foreign key added with MATCH FULL keeps its match mode in the migration.
	fkFrom := New()
	_, err = fkFrom.Exec(`
		CREATE TABLE public.p (a int, b int, PRIMARY KEY (a, b));
		CREATE TABLE public.c (a int, b int);
	`, nil)
	require.NoError(t, err)
	fkTo := fkFrom.Clone()
	_, err = fkTo.Exec(`ALTER TABLE public.c ADD CONSTRAINT c_fk FOREIGN KEY (a, b) REFERENCES public.p (a, b) MATCH FULL ON DELETE CASCADE;`, nil)
	require.NoError(t, err)
	fkSQL := GenerateMigration(fkFrom, fkTo, Diff(fkFrom, fkTo)).SQL()
	require.Contains(t, fkSQL, `REFERENCES public.p (a, b) MATCH FULL ON DELETE CASCADE`)
	fkApplied := fkFrom.Clone()
	require.NoError(t, round37ExecErr(fkApplied, fkSQL))
	require.True(t, Diff(fkApplied, fkTo).IsEmpty())

	// 5. ApplyDiffToMetadata leaves a new aggregate out of the functions, as sync does,
	// rather than writing a CREATE FUNCTION that would recreate something else.
	aggBefore := New()
	aggOrigMeta := &metadatapb.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	aggAfter := aggBefore.Clone()
	_, err = aggAfter.Exec(`
		CREATE FUNCTION public.add_int(int, int) RETURNS int LANGUAGE sql IMMUTABLE RETURN $1 + $2;
		CREATE AGGREGATE public.my_sum(int) (SFUNC = public.add_int, STYPE = int);
	`, nil)
	require.NoError(t, err)
	aggUpdated := ApplyDiffToMetadata(aggOrigMeta, aggBefore, aggAfter, Diff(aggBefore, aggAfter))
	var aggNames []string
	for _, fn := range aggUpdated.Schemas[0].Functions {
		aggNames = append(aggNames, fn.Name)
		require.NotContains(t, fn.Definition, "aggregate_dummy")
	}
	require.Equal(t, []string{"add_int"}, aggNames)

	// 6. A constraint trigger keeps its FROM table in migrations and in applied metadata.
	trigFrom := New()
	_, err = trigFrom.Exec(`
		CREATE TABLE public.ref (id int);
		CREATE TABLE public.src (id int);
		CREATE FUNCTION public.trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NULL; END$$;
	`, nil)
	require.NoError(t, err)
	trigTo := trigFrom.Clone()
	_, err = trigTo.Exec(`CREATE CONSTRAINT TRIGGER src_ct AFTER INSERT ON public.src FROM public.ref DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION public.trg_fn();`, nil)
	require.NoError(t, err)
	trigSQL := GenerateMigration(trigFrom, trigTo, Diff(trigFrom, trigTo)).SQL()
	require.Contains(t, trigSQL, `ON "public"."src" FROM "public"."ref" DEFERRABLE INITIALLY DEFERRED`)
	trigApplied := trigFrom.Clone()
	require.NoError(t, round37ExecErr(trigApplied, trigSQL))
	require.True(t, Diff(trigApplied, trigTo).IsEmpty())
	// Dropping the referenced table removes the trigger, as PostgreSQL does.
	require.NoError(t, round37ExecErr(trigApplied, `DROP TABLE public.ref CASCADE;`))
	require.Empty(t, trigApplied.triggersByRel[trigApplied.GetRelation("public", "src").OID])

	// 7. A SQL-standard body call depends only on the overload it resolves to.
	ovCat := New()
	_, err = ovCat.Exec(`
		CREATE FUNCTION public.z(integer) RETURNS integer LANGUAGE sql RETURN 1;
		CREATE FUNCTION public.z(text) RETURNS integer LANGUAGE sql RETURN 2;
		CREATE FUNCTION public.calls_int() RETURNS integer LANGUAGE sql RETURN public.z(1);
	`, nil)
	require.NoError(t, err)
	require.Error(t, round37ExecErr(ovCat.Clone(), `DROP FUNCTION public.z(integer);`))
	require.NoError(t, round37ExecErr(ovCat, `DROP FUNCTION public.z(text);`))

	// 8. ApplyDiffToMetadata quotes column keys, so a column named like an expression reloads as that column.
	qBefore := New()
	_, err = qBefore.Exec(`CREATE TABLE public.qt (a int, b int, "a + b" int, "Mixed" int);`, nil)
	require.NoError(t, err)
	qOrigMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "qt",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "a", Position: 1, Type: "integer", Nullable: true},
					{Name: "b", Position: 2, Type: "integer", Nullable: true},
					{Name: "a + b", Position: 3, Type: "integer", Nullable: true},
					{Name: "Mixed", Position: 4, Type: "integer", Nullable: true},
				},
			}},
		}},
	}
	qAfter := qBefore.Clone()
	_, err = qAfter.Exec(`CREATE INDEX qt_sum_col ON public.qt ("a + b", "Mixed", a);`, nil)
	require.NoError(t, err)
	qUpdated := ApplyDiffToMetadata(qOrigMeta, qBefore, qAfter, Diff(qBefore, qAfter))
	require.Len(t, qUpdated.Schemas[0].Tables[0].Indexes, 1)
	require.Equal(t, []string{`"a + b"`, `"Mixed"`, `a`}, qUpdated.Schemas[0].Tables[0].Indexes[0].Expressions)
	qReloaded := New()
	qReport, err := qReloaded.LoadMetadata(context.Background(), qUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, qReport.Degraded)
	require.True(t, Diff(qReloaded, qAfter).IsEmpty())

	// 9. A SQL-standard body records calls nested in array constructors and row expressions.
	arrCat := New()
	_, err = arrCat.Exec(`
		CREATE FUNCTION public.z() RETURNS integer LANGUAGE sql RETURN 1;
		CREATE FUNCTION public.y() RETURNS integer LANGUAGE sql RETURN 2;
		CREATE FUNCTION public.arr() RETURNS integer[] LANGUAGE sql RETURN ARRAY[public.z()];
		CREATE FUNCTION public.rowf() RETURNS boolean LANGUAGE sql RETURN ROW(public.y()) IS NOT NULL;
	`, nil)
	require.NoError(t, err)
	require.Error(t, round37ExecErr(arrCat.Clone(), `DROP FUNCTION public.z();`))
	require.Error(t, round37ExecErr(arrCat.Clone(), `DROP FUNCTION public.y();`))

	// 10. UNIQUE NULLS NOT DISTINCT survives migrations and ApplyDiffToMetadata, for constraints and plain indexes.
	nndBefore := New()
	_, err = nndBefore.Exec(`CREATE TABLE public.nnd (a int, b int);`, nil)
	require.NoError(t, err)
	nndAfter := nndBefore.Clone()
	_, err = nndAfter.Exec(`
		ALTER TABLE public.nnd ADD CONSTRAINT nnd_a_key UNIQUE NULLS NOT DISTINCT (a);
		CREATE UNIQUE INDEX nnd_b_idx ON public.nnd (b) NULLS NOT DISTINCT;
	`, nil)
	require.NoError(t, err)
	nndSQL := GenerateMigration(nndBefore, nndAfter, Diff(nndBefore, nndAfter)).SQL()
	require.Contains(t, nndSQL, `UNIQUE NULLS NOT DISTINCT (a)`)
	require.Contains(t, nndSQL, `(b) NULLS NOT DISTINCT`)
	nndApplied := nndBefore.Clone()
	require.NoError(t, round37ExecErr(nndApplied, nndSQL))
	require.True(t, Diff(nndApplied, nndAfter).IsEmpty())
	nndMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "nnd",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "a", Position: 1, Type: "integer", Nullable: true},
					{Name: "b", Position: 2, Type: "integer", Nullable: true},
				},
			}},
		}},
	}
	nndUpdated := ApplyDiffToMetadata(nndMeta, nndBefore, nndAfter, Diff(nndBefore, nndAfter))
	nndReloaded := New()
	nndReport, err := nndReloaded.LoadMetadata(context.Background(), nndUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, nndReport.Degraded)
	require.True(t, Diff(nndReloaded, nndAfter).IsEmpty())
	// Diff sees a constraint that loses NULLS NOT DISTINCT.
	nndPlain := nndBefore.Clone()
	_, err = nndPlain.Exec(`ALTER TABLE public.nnd ADD CONSTRAINT nnd_a_key UNIQUE (a); CREATE UNIQUE INDEX nnd_b_idx ON public.nnd (b) NULLS NOT DISTINCT;`, nil)
	require.NoError(t, err)
	require.False(t, Diff(nndPlain, nndAfter).IsEmpty())

	// 11. Detaching and renaming a partition in one script makes it a top-level table.
	partBase := New()
	_, err = partBase.Exec(`
		CREATE TABLE public.parent (id int, k int) PARTITION BY LIST (k);
		CREATE TABLE public.child PARTITION OF public.parent FOR VALUES IN (1);
		COMMENT ON TABLE public.child IS 'kept';
	`, nil)
	require.NoError(t, err)
	partCols := []*metadatapb.ColumnMetadata{
		{Name: "id", Position: 1, Type: "integer", Nullable: true},
		{Name: "k", Position: 2, Type: "integer", Nullable: true},
	}
	partMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "parent",
				Columns: partCols,
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:    "child",
					Type:    metadatapb.TablePartitionMetadata_LIST,
					Value:   "FOR VALUES IN (1)",
					Comment: "kept",
				}},
			}},
		}},
	}
	partAfter := partBase.Clone()
	_, err = partAfter.Exec(`
		ALTER TABLE public.parent DETACH PARTITION public.child;
		ALTER TABLE public.child RENAME TO archived;
	`, nil)
	require.NoError(t, err)
	partUpdated := ApplyDiffToMetadata(partMeta, partBase, partAfter, Diff(partBase, partAfter))
	var partSchema *metadatapb.SchemaMetadata
	for _, sm := range partUpdated.Schemas {
		if sm.Name == "public" {
			partSchema = sm
		}
	}
	require.NotNil(t, partSchema)
	tables := map[string]*metadatapb.TableMetadata{}
	for _, tbl := range partSchema.Tables {
		tables[tbl.Name] = tbl
	}
	require.Contains(t, tables, "archived")
	require.Equal(t, "kept", tables["archived"].Comment)
	require.Empty(t, tables["parent"].Partitions)
	partReloaded := New()
	partReport, err := partReloaded.LoadMetadata(context.Background(), partUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, partReport.Degraded)
	require.Zero(t, partReloaded.GetRelation("public", "archived").PartitionOf)

	// 12. Dropping one column of a mixed column grant still revokes the surviving column.
	grFrom := New()
	_, err = grFrom.Exec(`
		CREATE ROLE reader;
		CREATE TABLE public.gt (a int, b int, c int);
		GRANT SELECT (a, b) ON public.gt TO reader;
	`, nil)
	require.NoError(t, err)
	grTo := grFrom.Clone()
	_, err = grTo.Exec(`ALTER TABLE public.gt DROP COLUMN a; REVOKE SELECT (b) ON public.gt FROM reader;`, nil)
	require.NoError(t, err)
	grSQL := GenerateMigration(grFrom, grTo, Diff(grFrom, grTo)).SQL()
	require.Contains(t, grSQL, `REVOKE SELECT ("b")`)
	require.NotContains(t, grSQL, `"a", "b"`)
	grApplied := grFrom.Clone()
	require.NoError(t, round37ExecErr(grApplied, grSQL))
	require.True(t, Diff(grApplied, grTo).IsEmpty())
	// Dropping every granted column needs no revoke.
	grAll := grFrom.Clone()
	_, err = grAll.Exec(`ALTER TABLE public.gt DROP COLUMN a, DROP COLUMN b;`, nil)
	require.NoError(t, err)
	grAllSQL := GenerateMigration(grFrom, grAll, Diff(grFrom, grAll)).SQL()
	require.NotContains(t, grAllSQL, "REVOKE")
	grAllApplied := grFrom.Clone()
	require.NoError(t, round37ExecErr(grAllApplied, grAllSQL))
	require.True(t, Diff(grAllApplied, grAll).IsEmpty())

	// 13. A partition created in the same script as its parent nests under it, at any depth.
	nestFrom := New()
	nestTo := New()
	_, err = nestTo.Exec(`
		CREATE TABLE public.zparent (id int, k int) PARTITION BY LIST (k);
		CREATE TABLE public.mid PARTITION OF public.zparent FOR VALUES IN (1) PARTITION BY LIST (id);
		CREATE TABLE public.aleaf PARTITION OF public.mid FOR VALUES IN (7);
	`, nil)
	require.NoError(t, err)
	nestMeta := ApplyDiffToMetadata(&metadatapb.DatabaseSchemaMetadata{Name: "db"}, nestFrom, nestTo, Diff(nestFrom, nestTo))
	require.Len(t, nestMeta.Schemas, 1)
	require.Len(t, nestMeta.Schemas[0].Tables, 1)
	require.Equal(t, "zparent", nestMeta.Schemas[0].Tables[0].Name)
	require.Len(t, nestMeta.Schemas[0].Tables[0].Partitions, 1)
	require.Equal(t, "mid", nestMeta.Schemas[0].Tables[0].Partitions[0].Name)
	require.Len(t, nestMeta.Schemas[0].Tables[0].Partitions[0].Subpartitions, 1)
	require.Equal(t, "aleaf", nestMeta.Schemas[0].Tables[0].Partitions[0].Subpartitions[0].Name)

	// 14. Renaming a server keeps its attached foreign-table partition attached, and still detaches one the target detaches.
	fpFrom := New()
	_, err = fpFrom.Exec(`
		CREATE FOREIGN DATA WRAPPER w_fp;
		CREATE SERVER fp_old FOREIGN DATA WRAPPER w_fp;
		CREATE TABLE public.fp_parent (id int, k int) PARTITION BY LIST (k);
		CREATE FOREIGN TABLE public.fp_kept PARTITION OF public.fp_parent FOR VALUES IN (1) SERVER fp_old;
		CREATE FOREIGN TABLE public.fp_gone PARTITION OF public.fp_parent FOR VALUES IN (2) SERVER fp_old;
	`, nil)
	require.NoError(t, err)
	fpTo := fpFrom.Clone()
	_, err = fpTo.Exec(`
		ALTER SERVER fp_old RENAME TO fp_new;
		ALTER TABLE public.fp_parent DETACH PARTITION public.fp_gone;
	`, nil)
	require.NoError(t, err)
	fpSQL := GenerateMigration(fpFrom, fpTo, Diff(fpFrom, fpTo)).SQL()
	require.Contains(t, fpSQL, `ALTER SERVER fp_old RENAME TO fp_new`)
	require.NotContains(t, fpSQL, `ATTACH PARTITION`)
	require.Contains(t, fpSQL, `DETACH PARTITION`)
	fpApplied := fpFrom.Clone()
	require.NoError(t, round37ExecErr(fpApplied, fpSQL))
	require.True(t, Diff(fpApplied, fpTo).IsEmpty())

	// 15. Expression keys in applied index metadata carry the parentheses an index element needs.
	exBefore := New()
	_, err = exBefore.Exec(`CREATE TABLE public.ex (a int, b int, s text);`, nil)
	require.NoError(t, err)
	exMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "ex",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "a", Position: 1, Type: "integer", Nullable: true},
					{Name: "b", Position: 2, Type: "integer", Nullable: true},
					{Name: "s", Position: 3, Type: "text", Nullable: true},
				},
			}},
		}},
	}
	exAfter := exBefore.Clone()
	_, err = exAfter.Exec(`CREATE INDEX ex_sum ON public.ex ((a + b), lower(s));`, nil)
	require.NoError(t, err)
	exUpdated := ApplyDiffToMetadata(exMeta, exBefore, exAfter, Diff(exBefore, exAfter))
	require.Len(t, exUpdated.Schemas[0].Tables[0].Indexes, 1)
	require.Equal(t, []string{`(a + b)`, `lower(s)`}, exUpdated.Schemas[0].Tables[0].Indexes[0].Expressions)
	exReloaded := New()
	exReport, err := exReloaded.LoadMetadata(context.Background(), exUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, exReport.Degraded)
	require.True(t, Diff(exReloaded, exAfter).IsEmpty())

	// 16. Applied column metadata schema-qualifies user types visible on the search path.
	utBefore := New()
	_, err = utBefore.Exec(`CREATE TYPE public.status AS ENUM ('on', 'off'); CREATE TABLE public.ut (id int);`, nil)
	require.NoError(t, err)
	utMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			EnumTypes: []*metadatapb.EnumTypeMetadata{{Name: "status", Values: []string{"on", "off"}}},
			Tables: []*metadatapb.TableMetadata{{
				Name:    "ut",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Position: 1, Type: "integer", Nullable: true}},
			}},
		}},
	}
	utAfter := utBefore.Clone()
	_, err = utAfter.Exec(`ALTER TABLE public.ut ADD COLUMN s status, ADD COLUMN ss status[];`, nil)
	require.NoError(t, err)
	utUpdated := ApplyDiffToMetadata(utMeta, utBefore, utAfter, Diff(utBefore, utAfter))
	utTypes := map[string]string{}
	for _, col := range utUpdated.Schemas[0].Tables[0].Columns {
		utTypes[col.Name] = col.Type
	}
	require.Equal(t, "integer", utTypes["id"])
	require.Equal(t, "public.status", utTypes["s"])
	require.Equal(t, "public.status[]", utTypes["ss"])
	utReloaded := New()
	utReport, err := utReloaded.LoadMetadata(context.Background(), utUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, utReport.Degraded)
	require.True(t, Diff(utReloaded, utAfter).IsEmpty())

	// 17. Applied check metadata keeps NOT VALID and NO INHERIT, as pg_get_constraintdef writes them.
	nvBefore := New()
	_, err = nvBefore.Exec(`CREATE TABLE public.nv (a int);`, nil)
	require.NoError(t, err)
	nvMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "nv",
				Columns: []*metadatapb.ColumnMetadata{{Name: "a", Position: 1, Type: "integer", Nullable: true}},
			}},
		}},
	}
	nvAfter := nvBefore.Clone()
	_, err = nvAfter.Exec(`
		ALTER TABLE public.nv ADD CONSTRAINT nv_pos CHECK (a > 0) NOT VALID;
		ALTER TABLE public.nv ADD CONSTRAINT nv_small CHECK (a < 100) NO INHERIT;
		CREATE TABLE public.nv2 (b int, CONSTRAINT nv2_b CHECK (b > 0) NO INHERIT);
	`, nil)
	require.NoError(t, err)
	nvUpdated := ApplyDiffToMetadata(nvMeta, nvBefore, nvAfter, Diff(nvBefore, nvAfter))
	nvExprs := map[string]string{}
	for _, tbl := range nvUpdated.Schemas[0].Tables {
		for _, ck := range tbl.CheckConstraints {
			nvExprs[ck.Name] = ck.Expression
		}
	}
	require.True(t, strings.HasSuffix(nvExprs["nv_pos"], " NOT VALID"), nvExprs["nv_pos"])
	require.True(t, strings.HasSuffix(nvExprs["nv_small"], " NO INHERIT"), nvExprs["nv_small"])
	require.True(t, strings.HasSuffix(nvExprs["nv2_b"], " NO INHERIT"), nvExprs["nv2_b"])
	nvReloaded := New()
	nvReport, err := nvReloaded.LoadMetadata(context.Background(), nvUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, nvReport.Degraded)
	require.True(t, Diff(nvReloaded, nvAfter).IsEmpty())

	// 18. A new table keeps CHECK NO INHERIT and UNIQUE NULLS NOT DISTINCT in its CREATE TABLE.
	ctFrom := New()
	ctTo := New()
	_, err = ctTo.Exec(`CREATE TABLE public.ct (a int, b int, CONSTRAINT ct_a CHECK (a > 0) NO INHERIT, CONSTRAINT ct_b UNIQUE NULLS NOT DISTINCT (b));`, nil)
	require.NoError(t, err)
	ctSQL := GenerateMigration(ctFrom, ctTo, Diff(ctFrom, ctTo)).SQL()
	require.Contains(t, ctSQL, `CHECK ((a > 0)) NO INHERIT`)
	require.Contains(t, ctSQL, `UNIQUE NULLS NOT DISTINCT ("b")`)
	ctApplied := ctFrom.Clone()
	require.NoError(t, round37ExecErr(ctApplied, ctSQL))
	require.True(t, Diff(ctApplied, ctTo).IsEmpty())

	// 19. Columns added out of name order are appended in physical order.
	coBefore := New()
	_, err = coBefore.Exec(`CREATE TABLE public.co (id int);`, nil)
	require.NoError(t, err)
	coMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "co",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Position: 1, Type: "integer", Nullable: true}},
			}},
		}},
	}
	coAfter := coBefore.Clone()
	_, err = coAfter.Exec(`ALTER TABLE public.co ADD COLUMN z integer, ADD COLUMN a text, ADD COLUMN m date;`, nil)
	require.NoError(t, err)
	coUpdated := ApplyDiffToMetadata(coMeta, coBefore, coAfter, Diff(coBefore, coAfter))
	var coNames []string
	for _, col := range coUpdated.Schemas[0].Tables[0].Columns {
		coNames = append(coNames, col.Name)
	}
	require.Equal(t, []string{"id", "z", "a", "m"}, coNames)

	// 20. Attaching a table as a partition and renaming it in one script nests it under the parent.
	atBase := New()
	_, err = atBase.Exec(`
		CREATE TABLE public.at_parent (id int, k int) PARTITION BY LIST (k);
		CREATE TABLE public.at_t (id int, k int);
	`, nil)
	require.NoError(t, err)
	atCols := []*metadatapb.ColumnMetadata{
		{Name: "id", Position: 1, Type: "integer", Nullable: true},
		{Name: "k", Position: 2, Type: "integer", Nullable: true},
	}
	atMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{
				{Name: "at_parent", Columns: atCols},
				{Name: "at_t", Columns: atCols, Comment: "kept"},
			},
		}},
	}
	atAfter := atBase.Clone()
	_, err = atAfter.Exec(`
		ALTER TABLE public.at_parent ATTACH PARTITION public.at_t FOR VALUES IN (1);
		ALTER TABLE public.at_t RENAME TO at_renamed;
	`, nil)
	require.NoError(t, err)
	atUpdated := ApplyDiffToMetadata(atMeta, atBase, atAfter, Diff(atBase, atAfter))
	require.Len(t, atUpdated.Schemas[0].Tables, 1)
	require.Equal(t, "at_parent", atUpdated.Schemas[0].Tables[0].Name)
	require.Len(t, atUpdated.Schemas[0].Tables[0].Partitions, 1)
	require.Equal(t, "at_renamed", atUpdated.Schemas[0].Tables[0].Partitions[0].Name)
	require.Equal(t, "kept", atUpdated.Schemas[0].Tables[0].Partitions[0].Comment)

	// 21. ApplyDiffToMetadata keeps serial and identity sequences, as sync stores them, and drops them with their column.
	sqBefore := New()
	sqMeta := &metadatapb.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
	sqAfter := sqBefore.Clone()
	_, err = sqAfter.Exec(`
		CREATE TABLE public.sq (id serial, n int GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2), keep int);
	`, nil)
	require.NoError(t, err)
	sqUpdated := ApplyDiffToMetadata(sqMeta, sqBefore, sqAfter, Diff(sqBefore, sqAfter))
	sqByName := map[string]*metadatapb.SequenceMetadata{}
	for _, seq := range sqUpdated.Schemas[0].Sequences {
		sqByName[seq.Name] = seq
	}
	require.Contains(t, sqByName, "sq_id_seq")
	require.Equal(t, "sq", sqByName["sq_id_seq"].OwnerTable)
	require.Equal(t, "id", sqByName["sq_id_seq"].OwnerColumn)
	require.Contains(t, sqByName, "sq_n_seq")
	require.Equal(t, "5", sqByName["sq_n_seq"].Start)
	require.Equal(t, "2", sqByName["sq_n_seq"].Increment)
	sqReloaded := New()
	sqReport, err := sqReloaded.LoadMetadata(context.Background(), sqUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, sqReport.Degraded)
	require.True(t, Diff(sqReloaded, sqAfter).IsEmpty())
	// Dropping the serial column drops its sequence from the metadata.
	sqDropped := sqAfter.Clone()
	_, err = sqDropped.Exec(`ALTER TABLE public.sq DROP COLUMN id;`, nil)
	require.NoError(t, err)
	sqDroppedMeta := ApplyDiffToMetadata(sqUpdated, sqAfter, sqDropped, Diff(sqAfter, sqDropped))
	var sqLeft []string
	for _, seq := range sqDroppedMeta.Schemas[0].Sequences {
		sqLeft = append(sqLeft, seq.Name)
	}
	require.Equal(t, []string{"sq_n_seq"}, sqLeft)

	// 22. A change only to a serial sequence still reaches the metadata.
	ssBefore := New()
	_, err = ssBefore.Exec(`CREATE TABLE public.ss (id serial);`, nil)
	require.NoError(t, err)
	emptyBefore := New()
	ssMeta := ApplyDiffToMetadata(&metadatapb.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}, emptyBefore, ssBefore, Diff(emptyBefore, ssBefore))
	ssAfter := ssBefore.Clone()
	_, err = ssAfter.Exec(`ALTER SEQUENCE public.ss_id_seq CACHE 20;`, nil)
	require.NoError(t, err)
	require.True(t, Diff(ssBefore, ssAfter).IsEmpty())
	ssUpdated := ApplyDiffToMetadata(ssMeta, ssBefore, ssAfter, Diff(ssBefore, ssAfter))
	require.Len(t, ssUpdated.Schemas[0].Sequences, 1)
	require.Equal(t, "20", ssUpdated.Schemas[0].Sequences[0].CacheSize)

	// 23. Dropping or replacing a procedure that sync listed among functions removes that entry.
	prBefore := New()
	_, err = prBefore.Exec(`
		CREATE PROCEDURE public.p(a integer) LANGUAGE sql AS 'SELECT 1';
		CREATE PROCEDURE public.q(a integer) LANGUAGE sql AS 'SELECT 1';
	`, nil)
	require.NoError(t, err)
	prMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "p", Signature: "p(a integer)", Definition: "CREATE OR REPLACE PROCEDURE public.p(a integer)\n LANGUAGE sql\nAS $procedure$SELECT 1$procedure$\n"},
				{Name: "q", Signature: "q(a integer)", Definition: "CREATE OR REPLACE PROCEDURE public.q(a integer)\n LANGUAGE sql\nAS $procedure$SELECT 1$procedure$\n", Comment: "kept"},
			},
		}},
	}
	prAfter := prBefore.Clone()
	_, err = prAfter.Exec(`
		DROP PROCEDURE public.p(integer);
		CREATE OR REPLACE PROCEDURE public.q(a integer) LANGUAGE sql AS 'SELECT 2';
	`, nil)
	require.NoError(t, err)
	prUpdated := ApplyDiffToMetadata(prMeta, prBefore, prAfter, Diff(prBefore, prAfter))
	require.Empty(t, prUpdated.Schemas[0].Functions)
	require.Len(t, prUpdated.Schemas[0].Procedures, 1)
	require.Equal(t, "q", prUpdated.Schemas[0].Procedures[0].Name)
	require.Equal(t, "kept", prUpdated.Schemas[0].Procedures[0].Comment)

	// 24. A string-bodied SQL function records no dependency on the routines and types it calls, as PostgreSQL.
	strCat := New()
	_, err = strCat.Exec(`
		CREATE TYPE public.st AS (a int);
		CREATE FUNCTION public.zz() RETURNS integer LANGUAGE sql RETURN 1;
		CREATE FUNCTION public.str_body() RETURNS integer LANGUAGE sql AS $$ SELECT public.zz() + (NULL::public.st).a $$;
	`, nil)
	require.NoError(t, err)
	require.NoError(t, round37ExecErr(strCat, `DROP FUNCTION public.zz(); DROP TYPE public.st;`))
	require.NotNil(t, strCat.procByName["str_body"])

	// 25. A NO INHERIT check rebuilt around a column type change stays NO INHERIT.
	niFrom := New()
	_, err = niFrom.Exec(`CREATE TABLE public.ni (a int, CONSTRAINT ni_a CHECK (a > 0) NO INHERIT);`, nil)
	require.NoError(t, err)
	niTo := niFrom.Clone()
	_, err = niTo.Exec(`ALTER TABLE public.ni ALTER COLUMN a TYPE bigint;`, nil)
	require.NoError(t, err)
	niSQL := GenerateMigration(niFrom, niTo, Diff(niFrom, niTo)).SQL()
	require.Contains(t, niSQL, "NO INHERIT")
	niApplied := niFrom.Clone()
	require.NoError(t, round37ExecErr(niApplied, niSQL))
	require.True(t, Diff(niApplied, niTo).IsEmpty())

	// 26. Applied routine definitions schema-qualify user types in parameters, returns, and RETURNS TABLE.
	rtBefore := New()
	_, err = rtBefore.Exec(`CREATE TYPE public.status AS ENUM ('on', 'off');`, nil)
	require.NoError(t, err)
	rtMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:      "public",
			EnumTypes: []*metadatapb.EnumTypeMetadata{{Name: "status", Values: []string{"on", "off"}}},
		}},
	}
	rtAfter := rtBefore.Clone()
	_, err = rtAfter.Exec(`
		CREATE FUNCTION public.flip(s status) RETURNS status LANGUAGE sql AS 'SELECT s';
		CREATE FUNCTION public.rows_of(s status) RETURNS TABLE(v status) LANGUAGE sql AS 'SELECT s';
	`, nil)
	require.NoError(t, err)
	rtUpdated := ApplyDiffToMetadata(rtMeta, rtBefore, rtAfter, Diff(rtBefore, rtAfter))
	rtDefs := map[string]string{}
	for _, fn := range rtUpdated.Schemas[0].Functions {
		rtDefs[fn.Name] = fn.Definition
	}
	require.Contains(t, rtDefs["flip"], `"s" public.status`)
	require.Contains(t, rtDefs["flip"], `RETURNS public.status`)
	require.Contains(t, rtDefs["rows_of"], `TABLE("v" public.status)`)
	rtReloaded := New()
	rtReport, err := rtReloaded.LoadMetadata(context.Background(), rtUpdated, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, rtReport.Degraded)
	require.True(t, Diff(rtReloaded, rtAfter).IsEmpty())

	// 27. A serial default in applied metadata names its sequence schema-qualified, as sync's pg_get_expr does.
	ndBefore := New()
	ndAfter := ndBefore.Clone()
	_, err = ndAfter.Exec(`CREATE TABLE public.nd (id serial); CREATE TABLE public."Nd2" (id serial);`, nil)
	require.NoError(t, err)
	ndMeta := ApplyDiffToMetadata(&metadatapb.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}, ndBefore, ndAfter, Diff(ndBefore, ndAfter))
	ndDefaults := map[string]string{}
	for _, tbl := range ndMeta.Schemas[0].Tables {
		ndDefaults[tbl.Name] = tbl.Columns[0].Default
	}
	require.Equal(t, `nextval('public.nd_id_seq'::regclass)`, ndDefaults["nd"])
	require.Equal(t, `nextval('public."Nd2_id_seq"'::regclass)`, ndDefaults["Nd2"])
	ndReloaded := New()
	ndReport, err := ndReloaded.LoadMetadata(context.Background(), ndMeta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, ndReport.Degraded)

	// 28. An unqualified body call resolves its overload across search_path schemas.
	spCat := New()
	_, err = spCat.Exec(`
		CREATE SCHEMA app;
		SET search_path = public, app;
		CREATE FUNCTION public.zq(text) RETURNS integer LANGUAGE sql RETURN 1;
		CREATE FUNCTION app.zq(integer) RETURNS integer LANGUAGE sql RETURN 2;
		CREATE FUNCTION public.calls_zq() RETURNS integer LANGUAGE sql RETURN zq(1);
	`, nil)
	require.NoError(t, err)
	require.Error(t, round37ExecErr(spCat.Clone(), `DROP FUNCTION app.zq(integer);`))
	require.NoError(t, round37ExecErr(spCat, `DROP FUNCTION public.zq(text);`))

	// 29. A foreign partition is listed under its parent as well as among external tables, and leaves both when dropped.
	fpaBase := New()
	_, err = fpaBase.Exec(`
		CREATE FOREIGN DATA WRAPPER w_fpa;
		CREATE SERVER fpa_srv FOREIGN DATA WRAPPER w_fpa;
		CREATE TABLE public.fpa_parent (id int, k int) PARTITION BY RANGE (k);
	`, nil)
	require.NoError(t, err)
	fpaMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "fpa_parent",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Position: 1, Type: "integer", Nullable: true},
					{Name: "k", Position: 2, Type: "integer", Nullable: true},
				},
			}},
		}},
	}
	fpaAfter := fpaBase.Clone()
	_, err = fpaAfter.Exec(`CREATE FOREIGN TABLE public.fpa_remote PARTITION OF public.fpa_parent FOR VALUES FROM (10) TO (20) SERVER fpa_srv;`, nil)
	require.NoError(t, err)
	fpaUpdated := ApplyDiffToMetadata(fpaMeta, fpaBase, fpaAfter, Diff(fpaBase, fpaAfter))
	require.Len(t, fpaUpdated.Schemas[0].ExternalTables, 1)
	require.Len(t, fpaUpdated.Schemas[0].Tables[0].Partitions, 1)
	require.Equal(t, "fpa_remote", fpaUpdated.Schemas[0].Tables[0].Partitions[0].Name)
	require.Contains(t, fpaUpdated.Schemas[0].Tables[0].Partitions[0].Value, "FOR VALUES FROM")
	fpaDropped := fpaAfter.Clone()
	_, err = fpaDropped.Exec(`DROP FOREIGN TABLE public.fpa_remote;`, nil)
	require.NoError(t, err)
	fpaDroppedMeta := ApplyDiffToMetadata(fpaUpdated, fpaAfter, fpaDropped, Diff(fpaAfter, fpaDropped))
	require.Empty(t, fpaDroppedMeta.Schemas[0].ExternalTables)
	require.Empty(t, fpaDroppedMeta.Schemas[0].Tables[0].Partitions)

	// 30. Reattaching a partition under another bound and renaming it refreshes its bound.
	rbBase := New()
	_, err = rbBase.Exec(`
		CREATE TABLE public.rb_parent (id int, k int) PARTITION BY RANGE (k);
		CREATE TABLE public.rb_child PARTITION OF public.rb_parent FOR VALUES FROM (0) TO (10);
	`, nil)
	require.NoError(t, err)
	rbMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "rb_parent",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Position: 1, Type: "integer", Nullable: true},
					{Name: "k", Position: 2, Type: "integer", Nullable: true},
				},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:  "rb_child",
					Type:  metadatapb.TablePartitionMetadata_RANGE,
					Value: "FOR VALUES FROM (0) TO (10)",
				}},
			}},
		}},
	}
	rbAfter := rbBase.Clone()
	_, err = rbAfter.Exec(`
		ALTER TABLE public.rb_parent DETACH PARTITION public.rb_child;
		ALTER TABLE public.rb_parent ATTACH PARTITION public.rb_child FOR VALUES FROM (10) TO (20);
		ALTER TABLE public.rb_child RENAME TO rb_moved;
	`, nil)
	require.NoError(t, err)
	rbUpdated := ApplyDiffToMetadata(rbMeta, rbBase, rbAfter, Diff(rbBase, rbAfter))
	require.Len(t, rbUpdated.Schemas[0].Tables[0].Partitions, 1)
	require.Equal(t, "rb_moved", rbUpdated.Schemas[0].Tables[0].Partitions[0].Name)
	require.Contains(t, rbUpdated.Schemas[0].Tables[0].Partitions[0].Value, "20")

	// 31. A body call that resolves to a built-in records no dependency on a same-named user overload.
	biCat := New()
	_, err = biCat.Exec(`
		CREATE FUNCTION public.abs(text) RETURNS integer LANGUAGE sql RETURN 1;
		CREATE FUNCTION public.uses_abs() RETURNS integer LANGUAGE sql RETURN abs(1);
	`, nil)
	require.NoError(t, err)
	require.NoError(t, round37ExecErr(biCat, `DROP FUNCTION public.abs(text);`))

	// 32. A change only to an index's storage options migrates with ALTER INDEX SET / RESET.
	ixMeta := func(def string) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Name: "db",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name:    "ix_t",
					Columns: []*metadatapb.ColumnMetadata{{Name: "a", Position: 1, Type: "integer", Nullable: true}},
					Indexes: []*metadatapb.IndexMetadata{{Name: "ix_a", Type: "btree", Expressions: []string{"a"}, Definition: def}},
				}},
			}},
		}
	}
	ixSQL, err := testMetadataMigration(
		ixMeta(`CREATE INDEX ix_a ON public.ix_t USING btree (a) WITH (fillfactor='70', deduplicate_items='off')`),
		ixMeta(`CREATE INDEX ix_a ON public.ix_t USING btree (a) WITH (fillfactor='90')`),
	)
	require.NoError(t, err)
	require.Contains(t, ixSQL, `ALTER INDEX "public"."ix_a" RESET (deduplicate_items)`)
	require.Contains(t, ixSQL, `ALTER INDEX "public"."ix_a" SET (fillfactor='90')`)

	// 33. A source-only index dropped with its column is dropped IF EXISTS, after the column drop removed it.
	dcMeta := func(withCol bool) *metadatapb.DatabaseSchemaMetadata {
		cols := []*metadatapb.ColumnMetadata{{Name: "a", Position: 1, Type: "integer", Nullable: true}}
		var idx []*metadatapb.IndexMetadata
		if withCol {
			cols = append(cols, &metadatapb.ColumnMetadata{Name: "c", Position: 2, Type: "integer", Nullable: true})
			idx = []*metadatapb.IndexMetadata{{Name: "dc_idx", Type: "btree", Expressions: []string{"(c::no_such_domain)"},
				Definition: `CREATE INDEX dc_idx ON public.dc USING btree (((c)::no_such_domain))`}}
		}
		return &metadatapb.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{{Name: "dc", Columns: cols, Indexes: idx}},
		}}}
	}
	dcSQL, err := testMetadataMigration(dcMeta(true), dcMeta(false))
	require.NoError(t, err)
	require.Contains(t, dcSQL, `DROP COLUMN`)
	require.Contains(t, dcSQL, `DROP INDEX IF EXISTS "public"."dc_idx"`)
}

// round37ExecErr runs sql and returns its first error, statement errors included.
func round37ExecErr(c *Catalog, sql string) error {
	res, err := c.Exec(sql, nil)
	if err != nil {
		return err
	}
	for _, r := range res {
		if r.Error != nil {
			return r.Error
		}
	}
	return nil
}
