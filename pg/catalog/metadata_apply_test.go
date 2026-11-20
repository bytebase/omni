package catalog

import (
	"context"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

// TestClone_WalkThroughIntegration tests the full walk-through flow using Clone:
// load catalog → clone → exec user DDL on clone → diff → apply.
func TestClone_WalkThroughIntegration(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "users",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer", Position: 1},
							{Name: "name", Type: "text", Position: 2, Nullable: true},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{Name: "users_pkey", Expressions: []string{"id"}, Unique: true, Primary: true},
						},
					},
				},
			},
		},
	}

	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	catAfter := catBefore.Clone()

	userSQL := `
		ALTER TABLE public.users ADD COLUMN email text NOT NULL DEFAULT 'unknown';
		CREATE INDEX users_email_idx ON public.users (email);
		CREATE TABLE public.posts (id serial PRIMARY KEY, user_id int REFERENCES public.users(id), title text);
	`
	results, err := catAfter.Exec(userSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "DDL should succeed: %s", r.SQL)
	}

	// Diff
	diff := Diff(catBefore, catAfter)
	require.False(t, diff.IsEmpty())

	// Apply diff to original metadata
	newProto := ApplyDiffToMetadata(meta, catBefore, catAfter, diff)

	// Verify: users table should have 3 columns
	usersSchema := findSchema(newProto, "public")
	require.NotNil(t, usersSchema)
	usersTbl := findTable(usersSchema, "users")
	require.NotNil(t, usersTbl)
	require.Equal(t, 3, len(usersTbl.Columns), "users should have id, name, email")

	// Verify: users should have 2 indexes (pkey + email)
	require.GreaterOrEqual(t, len(usersTbl.Indexes), 2)

	// Verify: posts table should exist
	postsTbl := findTable(usersSchema, "posts")
	require.NotNil(t, postsTbl, "posts table should exist")

	// Verify: original metadata should NOT have posts or email column
	origSchema := findSchema(meta, "public")
	origUsers := findTable(origSchema, "users")
	require.Equal(t, 2, len(origUsers.Columns), "original should still have 2 columns")
	require.Nil(t, findTable(origSchema, "posts"), "original should not have posts")
}

func TestWalkThroughKeepsPartitionsUnderTheirTable(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "orders",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				Partitions: []*metadatapb.TablePartitionMetadata{
					{Name: "orders_2023"},
					{Name: "orders_2024", Subpartitions: []*metadatapb.TablePartitionMetadata{{Name: "orders_2024_q1"}, {Name: "orders_2024_q2"}}},
				},
			}},
		}},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		DROP TABLE public.orders_2023;
		DROP TABLE public.orders_2024_q2;
		CREATE INDEX orders_2024_q1_id_idx ON public.orders_2024_q1 (id);
		ALTER TABLE public.orders_2024 RENAME TO orders_2024_archived;
		CREATE INDEX orders_2024_archived_id_idx ON public.orders_2024_archived (id);
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	schema := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter)).Schemas[0]
	require.Len(t, schema.Tables, 1, "a changed partition must stay under its table")
	partitions := schema.Tables[0].Partitions
	require.Len(t, partitions, 1)
	require.Equal(t, "orders_2024_archived", partitions[0].Name)
	require.Len(t, partitions[0].Indexes, 1)
	require.Equal(t, "orders_2024_archived_id_idx", partitions[0].Indexes[0].Name)
	require.Len(t, partitions[0].Subpartitions, 1)
	q1 := partitions[0].Subpartitions[0]
	require.Equal(t, "orders_2024_q1", q1.Name)
	require.Len(t, q1.Indexes, 1)
	require.Equal(t, "orders_2024_q1_id_idx", q1.Indexes[0].Name)
}

// compositeWalkThroughMetadata has an adversarial name pair (aa_nested sorts
// before its dependency zz_base) plus an enum-referencing composite and a
// table using one, so a successful load proves dependency-ordered install.
func compositeWalkThroughMetadata() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				EnumTypes: []*metadatapb.EnumTypeMetadata{
					{Name: "status", Values: []string{"a", "b"}},
				},
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{
					{
						// References zz_base only through an array suffix.
						Name: "aa_array_only",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "items", Type: "public.zz_base[]"},
						},
					},
					{
						Name: "aa_nested",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "home", Type: "public.zz_base"},
							{Name: "s", Type: "public.status"},
						},
					},
					{
						Name:    "zz_base",
						Comment: "base address type",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "street", Type: "text", Collation: `"C"`, Comment: "street line"},
							{Name: "city", Type: "character varying(50)"},
						},
					},
					{
						Name:     "ext_owned",
						SkipDump: true,
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "x", Type: "integer"},
						},
					},
				},
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "users",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer", Position: 1},
							{Name: "home", Type: "public.zz_base", Position: 2, Nullable: true},
						},
					},
				},
			},
		},
	}
}

func TestWalkThroughCompositeFallbackPreservesAttributeNames(t *testing.T) {
	// LoadMetadata does not install domains, so a domain-typed attribute makes
	// the composite stand in — which must keep attribute names so later DDL
	// targeting them still resolves.
	meta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{
					{
						Name: "with_domain",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "p", Type: "public.pos_int"},
							{Name: "note", Type: "text"},
						},
					},
				},
			},
		},
	}

	cat := New()
	_, err := cat.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	rel := cat.GetRelation("public", "with_domain")
	require.NotNil(t, rel, "composite must fall back, not vanish")
	require.Len(t, rel.Columns, 2, "fallback must preserve attribute names")
	require.Equal(t, "p", rel.Columns[0].Name)
	require.Equal(t, "note", rel.Columns[1].Name)

	// DDL targeting a preserved attribute must succeed against the fallback.
	catAfter := cat.Clone()
	results, err := catAfter.Exec(`ALTER TYPE public.with_domain DROP ATTRIBUTE note;`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "DDL against fallback composite should succeed: %s", r.SQL)
	}

	// A rename keeps the attribute number, so the renamed attribute of a
	// stood-in composite must also keep its real previous type.
	renameResults, err := catAfter.Exec(`ALTER TYPE public.with_domain RENAME ATTRIBUTE p TO q;`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range renameResults {
		require.NoError(t, r.Error, "rename against fallback composite should succeed: %s", r.SQL)
	}

	// Applying the diff must not rewrite the untouched attribute to the
	// fallback's text type — unchanged attributes keep previous metadata.
	diff := Diff(cat, catAfter)
	require.False(t, diff.IsEmpty())
	newProto := ApplyDiffToMetadata(meta, cat, catAfter, diff)
	var applied *metadatapb.CompositeTypeMetadata
	for _, composite := range newProto.Schemas[0].CompositeTypes {
		if composite.Name == "with_domain" {
			applied = composite
		}
	}
	require.NotNil(t, applied)
	require.Len(t, applied.Attributes, 1, "dropped attribute must be removed")
	require.Equal(t, "q", applied.Attributes[0].Name, "rename must be reflected")
	require.Equal(t, "public.pos_int", applied.Attributes[0].Type,
		"renamed attribute must keep its real previous type, not the fallback text")
}

func TestWalkThroughDropReaddAttributeReadsCatalogType(t *testing.T) {
	// Dropping and re-adding an attribute with the same name assigns a new
	// attnum; the rebuilt metadata must take the catalog's new type, not
	// carry the stale previous metadata.
	meta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{
					{
						Name: "with_domain",
						Attributes: []*metadatapb.CompositeTypeAttribute{
							{Name: "p", Type: "public.pos_int", Comment: "old comment"},
							{Name: "note", Type: "text"},
						},
					},
				},
			},
		},
	}

	cat := New()
	_, err := cat.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	catAfter := cat.Clone()
	results, err := catAfter.Exec(`
		ALTER TYPE public.with_domain DROP ATTRIBUTE p;
		ALTER TYPE public.with_domain ADD ATTRIBUTE p text;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "DDL should succeed: %s", r.SQL)
	}

	diff := Diff(cat, catAfter)
	require.False(t, diff.IsEmpty())
	newProto := ApplyDiffToMetadata(meta, cat, catAfter, diff)
	var applied *metadatapb.CompositeTypeMetadata
	for _, composite := range newProto.Schemas[0].CompositeTypes {
		if composite.Name == "with_domain" {
			applied = composite
		}
	}
	require.NotNil(t, applied)
	require.Len(t, applied.Attributes, 2)
	var p *metadatapb.CompositeTypeAttribute
	for _, attribute := range applied.Attributes {
		if attribute.Name == "p" {
			p = attribute
		}
	}
	require.NotNil(t, p)
	require.Equal(t, "text", p.Type, "re-added attribute must take the catalog's new type: %v", applied)
	require.Empty(t, p.Comment, "re-added attribute must not inherit the dropped attribute's comment")
}

func TestWalkThroughAppliesCompositeTypeChanges(t *testing.T) {
	meta := compositeWalkThroughMetadata()

	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	catAfter := catBefore.Clone()
	userSQL := `
		CREATE TYPE public.geo AS (lat numeric(9,6), lng numeric(9,6));
		CREATE TYPE public.geo_wrap AS (g geo, gs geo[]);
		CREATE SCHEMA "select";
		CREATE TYPE "select".kw AS (x int);
		CREATE TYPE public.kw_wrap AS (k "select".kw);
		CREATE TABLE public.places (id int, location public.geo);
		DROP TYPE public.aa_nested;
		ALTER TYPE public.zz_base ADD ATTRIBUTE zip text;
		ALTER TYPE public.ext_owned ADD ATTRIBUTE y integer;
	`
	results, err := catAfter.Exec(userSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "DDL should succeed: %s", r.SQL)
	}

	diff := Diff(catBefore, catAfter)
	require.False(t, diff.IsEmpty())

	newProto := ApplyDiffToMetadata(meta, catBefore, catAfter, diff)

	publicSchema := findSchema(newProto, "public")
	require.NotNil(t, publicSchema)

	composites := make(map[string]*metadatapb.CompositeTypeMetadata)
	for _, composite := range publicSchema.CompositeTypes {
		composites[composite.Name] = composite
	}
	require.Contains(t, composites, "geo", "created composite must be applied to metadata")
	require.Len(t, composites["geo"].Attributes, 2)
	require.Equal(t, "lat", composites["geo"].Attributes[0].Name)
	require.Equal(t, "numeric(9,6)", composites["geo"].Attributes[0].Type,
		"built-in attribute types must stay unqualified")

	// User-defined attribute types must come back schema-qualified even when
	// the DDL referenced them unqualified via the search path.
	require.Contains(t, composites, "geo_wrap")
	require.Equal(t, "public.geo", composites["geo_wrap"].Attributes[0].Type)
	require.Equal(t, "public.geo[]", composites["geo_wrap"].Attributes[1].Type)

	// A schema named like a reserved keyword must come back quoted.
	require.Contains(t, composites, "kw_wrap")
	require.Equal(t, `"select".kw`, composites["kw_wrap"].Attributes[0].Type)
	require.NotContains(t, composites, "aa_nested", "dropped composite must be removed from metadata")
	require.Contains(t, composites, "zz_base", "altered composite must remain")

	// The modified composite gains the new attribute while preserving the
	// type comment and surviving attributes' comments/collations.
	base := composites["zz_base"]
	require.Equal(t, "base address type", base.Comment, "type comment must survive modification")
	require.Len(t, base.Attributes, 3)
	require.Equal(t, "street", base.Attributes[0].Name)
	require.Equal(t, `"C"`, base.Attributes[0].Collation, "attribute collation must survive modification")
	require.Equal(t, "street line", base.Attributes[0].Comment, "attribute comment must survive modification")
	require.Equal(t, "zip", base.Attributes[2].Name, "added attribute must appear")

	// skip_dump survives modification (extension-owned types must stay
	// excluded from dumps).
	require.Contains(t, composites, "ext_owned")
	require.True(t, composites["ext_owned"].SkipDump, "skip_dump must survive modification")

	// Original metadata is not mutated.
	require.Len(t, meta.Schemas[0].CompositeTypes, 4)
}

func TestApplyDiffToMetadataEmptyDiffClonesOriginal(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadatapb.TableMetadata{{Name: "t"}},
		}},
	}
	cloned := ApplyDiffToMetadata(meta, nil, nil, nil)
	require.NotSame(t, meta, cloned)
	cloned.Schemas[0].Tables[0].Name = "mutated"
	require.Equal(t, "t", meta.Schemas[0].Tables[0].Name)
}

func TestApplyDiffToMetadataRenameParentTablePreservesPartitions(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "orders",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer"}},
				Partitions: []*metadatapb.TablePartitionMetadata{
					{Name: "orders_2024"},
				},
			}},
		}},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`ALTER TABLE public.orders RENAME TO orders_new;`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
	sm := findSchema(got, "public")
	require.NotNil(t, sm)
	require.Nil(t, findTable(sm, "orders"))
	renamed := findTable(sm, "orders_new")
	require.NotNil(t, renamed)
	require.Len(t, renamed.Partitions, 1)
	require.Equal(t, "orders_2024", renamed.Partitions[0].Name)
}

func TestApplyDiffToMetadataModifiedColumnRetainsUnchangedFields(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t",
				Columns: []*metadatapb.ColumnMetadata{{
					Name:            "note",
					Type:            "text",
					Position:        1,
					Nullable:        true,
					Collation:       `"C"`,
					CollationSchema: "pg_catalog",
					Comment:         "keep me",
					Options:         []string{"column_name=note_remote"},
				}},
			}},
		}},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`ALTER TABLE public.t ALTER COLUMN note SET DEFAULT 'x';`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
	col := findTable(findSchema(got, "public"), "t").Columns[0]
	require.Equal(t, "'x'::text", col.Default)
	require.Equal(t, "keep me", col.Comment)
	require.Equal(t, `"C"`, col.Collation)
	require.Equal(t, "pg_catalog", col.CollationSchema)
	require.Equal(t, []string{"column_name=note_remote"}, col.Options)
}

func TestApplyDiffToMetadataRangesTriggersCommentsIndexesSequencesAndTableFunctions(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name: "t",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Type: "integer", Position: 1},
					{Name: "val", Type: "text", Position: 2, Nullable: true},
				},
			}},
		}},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		CREATE TYPE public.float8_range AS RANGE (SUBTYPE = float8);
		CREATE FUNCTION public.trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
		CREATE TRIGGER t_audit AFTER INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION public.trig_fn();
		COMMENT ON TABLE public.t IS 'table doc';
		COMMENT ON COLUMN public.t.val IS 'col doc';
		CREATE INDEX t_partial_idx ON public.t (id) INCLUDE (val) WHERE id > 0;
		CREATE SEQUENCE public.t_seq OWNED BY public.t.id;
		CREATE FUNCTION public.tab_fn(p integer) RETURNS TABLE(out_id integer, out_val text) LANGUAGE sql AS $$ SELECT p, 'a' $$;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
	sm := findSchema(got, "public")
	require.NotNil(t, sm)
	require.Len(t, sm.RangeTypes, 1)
	require.Equal(t, "float8_range", sm.RangeTypes[0].Name)

	tbl := findTable(sm, "t")
	require.Equal(t, "table doc", tbl.Comment)
	require.Equal(t, "col doc", tbl.Columns[1].Comment)
	require.Len(t, tbl.Triggers, 1)
	require.Equal(t, "t_audit", tbl.Triggers[0].Name)
	require.Len(t, tbl.Indexes, 1)
	require.Equal(t, []string{"id"}, tbl.Indexes[0].Expressions)
	require.Contains(t, tbl.Indexes[0].Definition, "INCLUDE")
	require.Contains(t, tbl.Indexes[0].Definition, "WHERE")

	require.Len(t, sm.Sequences, 1)
	require.Equal(t, "t", sm.Sequences[0].OwnerTable)
	require.Equal(t, "id", sm.Sequences[0].OwnerColumn)

	var tabFn *metadatapb.FunctionMetadata
	for _, fn := range sm.Functions {
		if fn.Name == "tab_fn" {
			tabFn = fn
		}
	}
	require.NotNil(t, tabFn)
	require.Contains(t, tabFn.Definition, `tab_fn("p" integer)`)
	require.Contains(t, tabFn.Definition, "RETURNS TABLE(")
}

func TestApplyDiffToMetadataPreservesMetadataOnRenamePartitionViewSequenceAndOwner(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name:  "s1",
			Owner: "alice",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "events",
				Owner:   "alice",
				Comment: "events table",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "id", Type: "integer", Position: 1},
					{Name: "region", Type: "text", Position: 2},
				},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:       "events_us",
					Type:       metadatapb.TablePartitionMetadata_LIST,
					Expression: "LIST (region)",
					Value:      "FOR VALUES IN ('us')",
					Comment:    "us partition",
				}},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:              "v_events",
				Definition:        "SELECT id FROM s1.events;",
				Comment:           "view doc",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "s1", Table: "events", Column: "id"}},
			}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
				Name:              "mv_events",
				Definition:        "SELECT id FROM s1.events;",
				Comment:           "mv doc",
				DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "s1", Table: "events", Column: "id"}},
			}},
			Sequences: []*metadatapb.SequenceMetadata{{
				Name:      "seq_id",
				DataType:  "bigint",
				Start:     "1",
				MinValue:  "1",
				MaxValue:  "9223372036854775807",
				Increment: "1",
				CacheSize: "1",
				LastValue: "42",
				Comment:   "seq doc",
			}},
		}},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	// Case 1: new partition, view update preserving comment/deps, sequence alter, owner change.
	catMod := catBefore.Clone()
	results, err := catMod.Exec(`
		CREATE TABLE s1.events_eu PARTITION OF s1.events FOR VALUES IN ('eu');
		CREATE OR REPLACE VIEW s1.v_events AS SELECT id, region FROM s1.events;
		ALTER SEQUENCE s1.seq_id INCREMENT BY 2;
		ALTER TABLE s1.events OWNER TO bob;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}
	modGot := ApplyDiffToMetadata(meta, catBefore, catMod, Diff(catBefore, catMod))
	modSM := findSchema(modGot, "s1")
	require.NotNil(t, modSM)
	tbl := findTable(modSM, "events")
	require.NotNil(t, tbl)
	require.Equal(t, "bob", tbl.Owner)
	require.Nil(t, findTable(modSM, "events_eu"), "new partition must nest under parent table")
	require.Len(t, tbl.Partitions, 2)
	require.Equal(t, "events_eu", tbl.Partitions[1].Name)
	require.Equal(t, "FOR VALUES IN ('eu')", tbl.Partitions[1].Value)
	require.Equal(t, metadatapb.TablePartitionMetadata_LIST, tbl.Partitions[1].Type)

	require.Len(t, modSM.Views, 1)
	require.Equal(t, "view doc", modSM.Views[0].Comment)
	require.Len(t, modSM.Views[0].DependencyColumns, 1)

	require.Len(t, modSM.Sequences, 1)
	require.Equal(t, "2", modSM.Sequences[0].Increment)
	require.Equal(t, "seq doc", modSM.Sequences[0].Comment)
	require.Equal(t, "42", modSM.Sequences[0].LastValue)

	// Case 2: schema rename preserves tables, partitions, views, sequences, and owner.
	catRenamed := catBefore.Clone()
	results, err = catRenamed.Exec(`ALTER SCHEMA s1 RENAME TO s2;`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}
	renGot := ApplyDiffToMetadata(meta, catBefore, catRenamed, Diff(catBefore, catRenamed))
	require.Nil(t, findSchema(renGot, "s1"))
	renSM := findSchema(renGot, "s2")
	require.NotNil(t, renSM)
	require.Equal(t, "alice", renSM.Owner)
	renTbl := findTable(renSM, "events")
	require.NotNil(t, renTbl)
	require.Equal(t, "alice", renTbl.Owner)
	require.Equal(t, "events table", renTbl.Comment)
	require.Len(t, renTbl.Partitions, 1)
	require.Equal(t, "us partition", renTbl.Partitions[0].Comment)
	require.Len(t, renSM.Sequences, 1)
	require.Equal(t, "seq doc", renSM.Sequences[0].Comment)
	require.Len(t, renSM.Views, 1)
	require.Equal(t, "view doc", renSM.Views[0].Comment)
	require.Len(t, renSM.MaterializedViews, 1)
	require.Equal(t, "mv doc", renSM.MaterializedViews[0].Comment)
}

func TestApplyDiffToMetadataExtensionsExternalTablesRenamesTriggersAndDefaults(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Extensions: []*metadatapb.ExtensionMetadata{{
			Name:        "old_ext",
			Schema:      "public",
			Version:     "1.0",
			Description: "preexisting",
		}},
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name: "t",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "id", Type: "integer", Position: 1},
						{Name: "old_col", Type: "text", Position: 2, Nullable: true, Comment: "col comment"},
					},
					Triggers: []*metadatapb.TriggerMetadata{{
						Name:    "tr_t",
						Timing:  "BEFORE",
						Event:   "INSERT",
						Body:    "CREATE TRIGGER tr_t BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION public.fn_a()",
						Comment: "trig doc",
					}},
				}},
				Views: []*metadatapb.ViewMetadata{{
					Name:              "v_t",
					Definition:        "SELECT id FROM public.t;",
					Comment:           "v doc",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
					Rules:             []*metadatapb.RuleMetadata{{Name: "r1"}},
				}},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
					Name:              "mv_t",
					Definition:        "SELECT id FROM public.t;",
					Comment:           "mv doc",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "t", Column: "id"}},
				}},
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "fn_a", Signature: "fn_a()", Definition: "CREATE FUNCTION public.fn_a() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$"},
					{Name: "fn_b", Signature: "fn_b()", Definition: "CREATE FUNCTION public.fn_b() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$"},
				},
			},
			{Name: "other"},
		},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	_, err = catBefore.Exec(`CREATE TRIGGER tr_t BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION public.fn_a();`, nil)
	require.NoError(t, err)

	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		DROP EXTENSION old_ext;
		CREATE EXTENSION pg_trgm WITH SCHEMA other;
		CREATE FOREIGN TABLE public.ft (id integer NOT NULL, note text) SERVER fsrv OPTIONS (schema_name 'remote', table_name 'ft');
		ALTER TABLE public.t RENAME COLUMN old_col TO new_col;
		CREATE OR REPLACE TRIGGER tr_t BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION public.fn_b();
		ALTER VIEW public.v_t RENAME TO v_t2;
		ALTER MATERIALIZED VIEW public.mv_t RENAME TO mv_t2;
		CREATE FUNCTION public.greet(who text DEFAULT 'world', times integer DEFAULT 2) RETURNS text LANGUAGE sql AS $$ SELECT who $$;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))

	// Extensions applied.
	require.Len(t, got.Extensions, 1)
	require.Equal(t, "pg_trgm", got.Extensions[0].Name)
	require.Equal(t, "other", got.Extensions[0].Schema)

	pub := findSchema(got, "public")
	require.NotNil(t, pub)

	// Foreign table stored under ExternalTables, not Tables.
	require.Len(t, pub.Tables, 1)
	require.Len(t, pub.ExternalTables, 1)
	require.Equal(t, "ft", pub.ExternalTables[0].Name)
	require.Equal(t, "fsrv", pub.ExternalTables[0].ExternalServerName)
	require.Equal(t, []string{"schema_name=remote", "table_name=ft"}, pub.ExternalTables[0].Options)

	// Renamed column kept its Comment; replaced trigger kept its Comment.
	tbl := findTable(pub, "t")
	require.NotNil(t, tbl)
	require.Len(t, tbl.Columns, 2)
	require.Equal(t, "new_col", tbl.Columns[1].Name)
	require.Equal(t, "col comment", tbl.Columns[1].Comment)
	require.Len(t, tbl.Triggers, 1)
	require.Equal(t, "trig doc", tbl.Triggers[0].Comment)
	require.Contains(t, tbl.Triggers[0].Body, "fn_b")

	// Renamed view and materialized view kept metadata.
	require.Len(t, pub.Views, 1)
	require.Equal(t, "v_t2", pub.Views[0].Name)
	require.Equal(t, "v doc", pub.Views[0].Comment)
	require.Len(t, pub.Views[0].DependencyColumns, 1)
	require.Len(t, pub.Views[0].Rules, 1)
	require.Len(t, pub.MaterializedViews, 1)
	require.Equal(t, "mv_t2", pub.MaterializedViews[0].Name)
	require.Equal(t, "mv doc", pub.MaterializedViews[0].Comment)
	require.Len(t, pub.MaterializedViews[0].DependencyColumns, 1)

	// Function parameter DEFAULTs preserved in generated definition.
	var greetFn *metadatapb.FunctionMetadata
	for _, fn := range pub.Functions {
		if fn.Name == "greet" {
			greetFn = fn
			break
		}
	}
	require.NotNil(t, greetFn)
	require.Contains(t, greetFn.Definition, "DEFAULT 'world'")
	require.Contains(t, greetFn.Definition, "DEFAULT 2")
}

func TestApplyDiffToMetadataTableAndFunctionRenamesAndOwnerOnly(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name:  "orders",
					Owner: "alice",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "id", Type: "integer", Position: 1},
					},
					Indexes: []*metadatapb.IndexMetadata{{
						Name:         "orders_pkey",
						Expressions:  []string{"id"},
						Primary:      true,
						Unique:       true,
						IsConstraint: true,
						Comment:      "pk doc",
						Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "orders_pkey", Primary: true}},
					}},
				}},
				Functions: []*metadatapb.FunctionMetadata{{
					Name:             "fn_old",
					Signature:        "fn_old(integer)",
					Definition:       "CREATE FUNCTION public.fn_old(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$",
					Comment:          "fn doc",
					SkipDump:         true,
					DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "orders"}},
				}},
			},
			{Name: "other"},
		},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	// Ensure fn_old is installed as non-skip-dump for rename test.
	_, err = catBefore.Exec(`CREATE FUNCTION public.fn_ren(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;`, nil)
	require.NoError(t, err)
	meta.Schemas[0].Functions = append(meta.Schemas[0].Functions, &metadatapb.FunctionMetadata{
		Name:             "fn_ren",
		Signature:        "public.fn_ren(integer)",
		Definition:       "CREATE FUNCTION public.fn_ren(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$",
		Comment:          "ren doc",
		DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "orders"}},
	})

	// Case 1: owner-only change updates TableMetadata.Owner even with empty diff.
	catOwner := catBefore.Clone()
	_, err = catOwner.Exec(`ALTER TABLE public.orders OWNER TO bob;`, nil)
	require.NoError(t, err)
	ownerGot := ApplyDiffToMetadata(meta, catBefore, catOwner, Diff(catBefore, catOwner))
	require.Equal(t, "bob", findTable(findSchema(ownerGot, "public"), "orders").Owner)

	// Case 2: table rename preserves index comment & Constraints; function rename preserves DependencyTables & Comment.
	catRen := catBefore.Clone()
	results, err := catRen.Exec(`
		ALTER TABLE public.orders RENAME TO orders_new;
		ALTER FUNCTION public.fn_ren(integer) RENAME TO fn_new;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}
	renGot := ApplyDiffToMetadata(meta, catBefore, catRen, Diff(catBefore, catRen))
	pub := findSchema(renGot, "public")
	tbl := findTable(pub, "orders_new")
	require.NotNil(t, tbl)
	require.Len(t, tbl.Indexes, 1)
	require.Equal(t, "pk doc", tbl.Indexes[0].Comment)
	require.Len(t, tbl.Indexes[0].Constraints, 1)

	var fnNew *metadatapb.FunctionMetadata
	for _, fn := range pub.Functions {
		if fn.Name == "fn_new" {
			fnNew = fn
			break
		}
	}
	require.NotNil(t, fnNew)
	require.Equal(t, "ren doc", fnNew.Comment)
	require.Len(t, fnNew.DependencyTables, 1)
	require.Equal(t, "orders_new", fnNew.DependencyTables[0].Table)
}

func TestApplyDiffToMetadataRound5Fixes(t *testing.T) {
	ctx := context.Background()

	// 1. Range rename & SET SCHEMA preserve metadata.
	// 2. EventTrigger add/modify(disable)/rename/drop preserve Comment & SkipDump.
	// 3. Constraint DiffModify preserves Comment.
	// 4. Cross-schema partition drop removes partition from parent schema table.
	// 5. Table rename rewrites ViewMetadata.DependencyColumns as well as FunctionMetadata.DependencyTables.
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "parents",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Indexes: []*metadatapb.IndexMetadata{{
							Name:         "parents_pkey",
							Expressions:  []string{"id"},
							Primary:      true,
							Unique:       true,
							IsConstraint: true,
							Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "parents_pkey", Primary: true}},
						}},
					},
					{
						Name:    "orders",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}, {Name: "parent_id", Type: "integer", Position: 2}},
					},
					{
						Name:    "items",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}, {Name: "parent_id", Type: "integer", Position: 2}},
						ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
							Name:              "items_parent_fk",
							Columns:           []string{"parent_id"},
							ReferencedSchema:  "public",
							ReferencedTable:   "parents",
							ReferencedColumns: []string{"id"},
							OnDelete:          "NO ACTION",
							Comment:           "fk comment",
						}},
					},
					{
						Name:    "logs",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Partitions: []*metadatapb.TablePartitionMetadata{{
							Name:   "logs_2026",
							Schema: "archive",
							Type:   metadatapb.TablePartitionMetadata_RANGE,
							Value:  "FOR VALUES FROM (2026) TO (2027)",
						}},
					},
				},
				Views: []*metadatapb.ViewMetadata{{
					Name:              "orders_v",
					Definition:        "SELECT id FROM public.orders",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "orders", Column: "id"}},
				}},
				RangeTypes: []*metadatapb.RangeTypeMetadata{{
					Name:           "float8_range",
					Subtype:        "float8",
					SubtypeOpclass: "float8_ops",
					Canonical:      "custom_canon",
					Comment:        "range comment",
					SkipDump:       true,
				}},
			},
			{Name: "archive"},
		},
		EventTriggers: []*metadatapb.EventTriggerMetadata{{
			Name:           "et_existing",
			Event:          "ddl_command_end",
			FunctionSchema: "public",
			FunctionName:   "et_fn",
			Enabled:        true,
			Comment:        "et doc",
			SkipDump:       true,
		}},
	}

	catBefore := New()
	_, err := catBefore.LoadMetadata(ctx, meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	// Install objects that LoadMetadata does not install natively (range type, partitioned parent + cross-schema child, event trigger).
	for _, stmt := range []string{
		`CREATE FUNCTION public.et_fn() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN END; $$;`,
		`CREATE EVENT TRIGGER et_existing ON ddl_command_end EXECUTE FUNCTION public.et_fn();`,
		`CREATE TYPE public.float8_range AS RANGE (subtype = float8);`,
		`CREATE TABLE public.logs (id integer) PARTITION BY RANGE (id);`,
		`CREATE TABLE archive.logs_2026 PARTITION OF public.logs FOR VALUES FROM (2026) TO (2027);`,
	} {
		_, err = catBefore.Exec(stmt, nil)
		require.NoError(t, err, stmt)
	}

	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		ALTER TYPE public.float8_range RENAME TO f8_range;
		ALTER TYPE public.f8_range SET SCHEMA archive;
		ALTER EVENT TRIGGER et_existing DISABLE;
		ALTER EVENT TRIGGER et_existing RENAME TO et_renamed;
		CREATE EVENT TRIGGER et_new ON sql_drop WHEN TAG IN ('DROP TABLE') EXECUTE FUNCTION public.et_fn();
		ALTER TABLE public.items DROP CONSTRAINT items_parent_fk, ADD CONSTRAINT items_parent_fk FOREIGN KEY (parent_id) REFERENCES public.parents(id) ON DELETE CASCADE;
		DROP TABLE archive.logs_2026;
		ALTER TABLE public.orders RENAME TO orders_renamed;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))

	// 1. Range metadata preserved after RENAME + SET SCHEMA.
	pub := findSchema(got, "public")
	arch := findSchema(got, "archive")
	require.Empty(t, pub.RangeTypes)
	require.Len(t, arch.RangeTypes, 1)
	require.Equal(t, "f8_range", arch.RangeTypes[0].Name)
	require.Equal(t, "range comment", arch.RangeTypes[0].Comment)
	require.True(t, arch.RangeTypes[0].SkipDump)
	require.Equal(t, "custom_canon", arch.RangeTypes[0].Canonical)

	// 2. Event triggers updated, rename preserved Comment & SkipDump, Disable reflected.
	require.Len(t, got.EventTriggers, 2)
	var etRen, etNew *metadatapb.EventTriggerMetadata
	for _, et := range got.EventTriggers {
		switch et.Name {
		case "et_renamed":
			etRen = et
		case "et_new":
			etNew = et
		}
	}
	require.NotNil(t, etRen)
	require.False(t, etRen.Enabled)
	require.Equal(t, "et doc", etRen.Comment)
	require.True(t, etRen.SkipDump)
	require.NotNil(t, etNew)
	require.True(t, etNew.Enabled)
	require.Equal(t, "sql_drop", etNew.Event)
	require.Equal(t, []string{"DROP TABLE"}, etNew.Tags)
	require.Equal(t, "public", etNew.FunctionSchema)
	require.Equal(t, "et_fn", etNew.FunctionName)

	// 3. Constraint DiffModify preserved Comment.
	itemsTbl := findTable(pub, "items")
	require.NotNil(t, itemsTbl)
	require.Len(t, itemsTbl.ForeignKeys, 1)
	require.Equal(t, "CASCADE", itemsTbl.ForeignKeys[0].OnDelete)
	require.Equal(t, "fk comment", itemsTbl.ForeignKeys[0].Comment)

	// 4. Cross-schema partition removed from public.logs.
	logsTbl := findTable(pub, "logs")
	require.NotNil(t, logsTbl)
	require.Empty(t, logsTbl.Partitions)

	// 5. View DependencyColumns rewritten on table rename.
	require.NotNil(t, findTable(pub, "orders_renamed"))
	require.Len(t, pub.Views, 1)
	require.Len(t, pub.Views[0].DependencyColumns, 1)
	require.Equal(t, "orders_renamed", pub.Views[0].DependencyColumns[0].Table)
}

func TestApplyDiffToMetadataRound6Fixes(t *testing.T) {
	ctx := context.Background()

	// 1. Cross-schema partition comment & index comment.
	// 2. Partition index ParentIndexSchema & ParentIndexName on create & ATTACH.
	// 3. ATTACH PARTITION moves standalone table under parent Partitions, DETACH moves it back.
	// 4. Quoted routine identifier ("select") in reconstructed definition.
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "events",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Indexes: []*metadatapb.IndexMetadata{{
							Name:        "events_id_idx",
							Type:        "btree",
							Expressions: []string{"id"},
						}},
						Partitions: []*metadatapb.TablePartitionMetadata{{
							Name:   "events_p1",
							Schema: "archive",
							Type:   metadatapb.TablePartitionMetadata_RANGE,
							Value:  "FOR VALUES FROM (1) TO (10)",
							Indexes: []*metadatapb.IndexMetadata{{
								Name:              "events_p1_id_idx",
								Type:              "btree",
								Expressions:       []string{"id"},
								ParentIndexSchema: "public",
								ParentIndexName:   "events_id_idx",
							}},
						}},
					},
					{
						Name:    "events_p2",
						Comment: "p2 doc",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Indexes: []*metadatapb.IndexMetadata{{
							Name:        "events_p2_idx",
							Type:        "btree",
							Expressions: []string{"id"},
							Comment:     "p2 idx doc",
						}},
					},
				},
			},
			{Name: "archive"},
		},
	}

	catBefore := New()
	for _, stmt := range []string{
		`CREATE SCHEMA archive;`,
		`CREATE TABLE public.events (id integer) PARTITION BY RANGE (id);`,
		`CREATE INDEX events_id_idx ON ONLY public.events (id);`,
		`CREATE TABLE archive.events_p1 PARTITION OF public.events FOR VALUES FROM (1) TO (10);`,
		`CREATE INDEX events_p1_id_idx ON archive.events_p1 (id);`,
		`ALTER INDEX public.events_id_idx ATTACH PARTITION archive.events_p1_id_idx;`,
		`CREATE TABLE public.events_p2 (id integer);`,
		`CREATE INDEX events_p2_idx ON public.events_p2 (id);`,
	} {
		_, err := catBefore.Exec(stmt, nil)
		require.NoError(t, err, stmt)
	}

	// Step A: comment on cross-schema partition & its index, ATTACH events_p2 & index, create partition events_p3 with cloned index, create quoted routine "select".
	catAttach := catBefore.Clone()
	results, err := catAttach.Exec(`
		COMMENT ON TABLE archive.events_p1 IS 'p1 table doc';
		COMMENT ON INDEX archive.events_p1_id_idx IS 'p1 idx doc';
		ALTER TABLE public.events ATTACH PARTITION public.events_p2 FOR VALUES FROM (10) TO (20);
		CREATE TABLE archive.events_p3 PARTITION OF public.events FOR VALUES FROM (20) TO (30);
		CREATE FUNCTION public."select"(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	attachedMeta := ApplyDiffToMetadata(meta, catBefore, catAttach, Diff(catBefore, catAttach))
	pub := findSchema(attachedMeta, "public")
	require.Nil(t, findTable(pub, "events_p2"))
	eventsTbl := findTable(pub, "events")
	require.NotNil(t, eventsTbl)
	require.Len(t, eventsTbl.Partitions, 3)

	p1 := findPartitionInSchema(attachedMeta, "archive", "events_p1")
	require.NotNil(t, p1)
	require.Equal(t, "p1 table doc", p1.Comment)
	require.Len(t, p1.Indexes, 1)
	require.Equal(t, "p1 idx doc", p1.Indexes[0].Comment)
	require.Equal(t, "public", p1.Indexes[0].ParentIndexSchema)
	require.Equal(t, "events_id_idx", p1.Indexes[0].ParentIndexName)

	p2 := findPartitionInSchema(attachedMeta, "public", "events_p2")
	require.NotNil(t, p2)
	require.Equal(t, "p2 doc", p2.Comment)
	require.Len(t, p2.Indexes, 2)
	require.Equal(t, "p2 idx doc", p2.Indexes[0].Comment)
	require.Equal(t, "public", p2.Indexes[1].ParentIndexSchema)
	require.Equal(t, "events_id_idx", p2.Indexes[1].ParentIndexName)

	p3 := findPartitionInSchema(attachedMeta, "archive", "events_p3")
	require.NotNil(t, p3)
	require.Len(t, p3.Indexes, 1)
	require.Equal(t, "public", p3.Indexes[0].ParentIndexSchema)
	require.Equal(t, "events_id_idx", p3.Indexes[0].ParentIndexName)

	// Verify quoted routine definition reloads cleanly into a fresh catalog.
	var selectFn *metadatapb.FunctionMetadata
	for _, fn := range pub.Functions {
		if fn.Name == "select" {
			selectFn = fn
			break
		}
	}
	require.NotNil(t, selectFn)
	require.Contains(t, selectFn.Definition, `public."select"`)
	reloadCat := New()
	report, err := reloadCat.LoadMetadata(ctx, attachedMeta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	require.Empty(t, report.Degraded)
	require.Len(t, reloadCat.findUserProcsByName(reloadCat.GetSchema("public"), "select"), 1)

	// Step B: DETACH archive.events_p1 moves it from public.events.Partitions into archive.Tables without duplication.
	catDetach := catAttach.Clone()
	_, err = catDetach.Exec(`ALTER TABLE public.events DETACH PARTITION archive.events_p1;`, nil)
	require.NoError(t, err)
	detachedMeta := ApplyDiffToMetadata(attachedMeta, catAttach, catDetach, Diff(catAttach, catDetach))
	require.Nil(t, findPartitionInSchema(detachedMeta, "archive", "events_p1"))
	arch := findSchema(detachedMeta, "archive")
	detachedTbl := findTable(arch, "events_p1")
	require.NotNil(t, detachedTbl)
	require.Equal(t, "p1 table doc", detachedTbl.Comment)
	var detachedIdx *metadatapb.IndexMetadata
	for _, idx := range detachedTbl.Indexes {
		require.Empty(t, idx.ParentIndexName)
		if idx.Name == "events_p1_id_idx" {
			detachedIdx = idx
		}
	}
	require.NotNil(t, detachedIdx)
	require.Equal(t, "p1 idx doc", detachedIdx.Comment)
}
