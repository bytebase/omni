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
	require.Len(t, modSM.Views[0].DependencyColumns, 2)

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
	_, err = catBefore.Exec(`
		CREATE FOREIGN DATA WRAPPER dummy;
		CREATE SERVER fsrv FOREIGN DATA WRAPPER dummy;
		CREATE TRIGGER tr_t BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION public.fn_a();
	`, nil)
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
		`CREATE TYPE public.float8_range AS RANGE (subtype = float8, subtype_opclass = float8_ops);`,
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
	require.Equal(t, "float8_ops", arch.RangeTypes[0].SubtypeOpclass)

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

func TestApplyDiffToMetadataRound7Fixes(t *testing.T) {
	// 1. Same-parent detach+reattach with bound change updates Value.
	// 2. Partition rename with colliding name in another schema mutates only matching effective schema.
	// 3. Dotted quoted schema/table identifier comment is applied accurately.
	catBefore := New()
	for _, stmt := range []string{
		`CREATE SCHEMA archive;`,
		`CREATE SCHEMA "app.data";`,
		`CREATE TABLE public.events (id integer) PARTITION BY RANGE (id);`,
		`CREATE TABLE public.logs (id integer) PARTITION BY RANGE (id);`,
		`CREATE TABLE public.dup PARTITION OF public.events FOR VALUES FROM (1) TO (10);`,
		`CREATE TABLE archive.dup PARTITION OF public.logs FOR VALUES FROM (1) TO (10);`,
		`CREATE TABLE "app.data"."events.v1" (id integer, CONSTRAINT "chk.v1" CHECK (id > 0));`,
		`CREATE FUNCTION public.et_r7() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN END; $$;`,
		`CREATE EVENT TRIGGER "et.dot" ON ddl_command_end EXECUTE FUNCTION public.et_r7();`,
	} {
		_, err := catBefore.Exec(stmt, nil)
		require.NoError(t, err, stmt)
	}
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "events",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Partitions: []*metadatapb.TablePartitionMetadata{{
							Name:  "dup",
							Type:  metadatapb.TablePartitionMetadata_RANGE,
							Value: "FOR VALUES FROM (1) TO (10)",
						}},
					},
					{
						Name:    "logs",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Partitions: []*metadatapb.TablePartitionMetadata{{
							Name:   "dup",
							Schema: "archive",
							Type:   metadatapb.TablePartitionMetadata_RANGE,
							Value:  "FOR VALUES FROM (1) TO (10)",
						}},
					},
				},
			},
			{Name: "archive"},
			{
				Name: "app.data",
				Tables: []*metadatapb.TableMetadata{{
					Name:             "events.v1",
					Columns:          []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
					CheckConstraints: []*metadatapb.CheckConstraintMetadata{{Name: "chk.v1", Expression: "id > 0"}},
				}},
			},
		},
		EventTriggers: []*metadatapb.EventTriggerMetadata{{
			Name:           "et.dot",
			Event:          "ddl_command_end",
			FunctionSchema: "public",
			FunctionName:   "et_r7",
			Enabled:        true,
		}},
	}

	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		ALTER TABLE public.events DETACH PARTITION public.dup;
		ALTER TABLE public.events ATTACH PARTITION public.dup FOR VALUES FROM (100) TO (200);
		ALTER TABLE archive.dup RENAME TO dup_renamed;
		COMMENT ON TABLE "app.data"."events.v1" IS 'dotted doc';
		COMMENT ON CONSTRAINT "chk.v1" ON "app.data"."events.v1" IS 'chk doc';
		COMMENT ON EVENT TRIGGER "et.dot" IS 'et doc';
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
	pubDup := findPartitionInSchema(got, "public", "dup")
	require.NotNil(t, pubDup)
	require.Equal(t, "FOR VALUES FROM ('100') TO ('200')", pubDup.Value)

	require.Nil(t, findPartitionInSchema(got, "archive", "dup"))
	archRen := findPartitionInSchema(got, "archive", "dup_renamed")
	require.NotNil(t, archRen)

	dottedTbl := findTable(findSchema(got, "app.data"), "events.v1")
	require.NotNil(t, dottedTbl)
	require.Equal(t, "dotted doc", dottedTbl.Comment)
	require.Equal(t, "chk doc", dottedTbl.CheckConstraints[0].Comment)
	require.Equal(t, "et doc", findEventTrigger(got, "et.dot").Comment)
}

func TestApplyDiffToMetadataRound8Fixes(t *testing.T) {
	catBefore := New()
	for _, stmt := range []string{
		`CREATE SCHEMA archive;`,
		`CREATE FOREIGN DATA WRAPPER dummy;`,
		`CREATE SERVER srv FOREIGN DATA WRAPPER dummy;`,
		`CREATE TABLE public.events (id integer) PARTITION BY RANGE (id);`,
		`CREATE INDEX events_id_idx ON public.events (id);`,
		`CREATE TABLE public.events_p1 PARTITION OF public.events FOR VALUES FROM (1) TO (10);`,
		`CREATE TABLE public.items (id serial, qty integer);`,
		`CREATE VIEW public.v_items AS SELECT id, qty FROM public.items;`,
		`CREATE VIEW public.v_part AS SELECT id FROM public.events_p1;`,
		`CREATE FUNCTION public.fn_part() RETURNS bigint LANGUAGE sql AS $$ SELECT count(*) FROM public.events_p1 $$;`,
		`CREATE SEQUENCE public.custom_seq START 10;`,
		`CREATE FOREIGN TABLE public.ft_orders (id integer, CONSTRAINT chk_ft CHECK (id > 0)) SERVER srv;`,
	} {
		_, err := catBefore.Exec(stmt, nil)
		require.NoError(t, err, stmt)
	}

	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "events",
						Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Indexes: []*metadatapb.IndexMetadata{{Name: "events_id_idx", Type: "btree", Expressions: []string{"id"}}},
						Partitions: []*metadatapb.TablePartitionMetadata{{
							Name:  "events_p1",
							Type:  metadatapb.TablePartitionMetadata_RANGE,
							Value: "FOR VALUES FROM (1) TO (10)",
							Indexes: []*metadatapb.IndexMetadata{{
								Name:              "events_p1_events_id_idx_idx",
								Type:              "btree",
								Expressions:       []string{"id"},
								ParentIndexSchema: "public",
								ParentIndexName:   "events_id_idx",
								Comment:           "part idx doc",
							}},
						}},
					},
					{
						Name: "items",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer", Position: 1, Default: "nextval('public.items_id_seq'::regclass)"},
							{Name: "qty", Type: "integer", Position: 2},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:              "v_items",
						Definition:        "SELECT id, qty FROM public.items",
						DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "items", Column: "qty"}, {Schema: "public", Table: "items", Column: "id"}},
					},
					{
						Name:              "v_part",
						Definition:        "SELECT id FROM public.events_p1",
						DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "events_p1", Column: "id"}},
					},
				},
				Functions: []*metadatapb.FunctionMetadata{{
					Name:             "fn_part",
					Signature:        "fn_part()",
					Definition:       "CREATE FUNCTION public.fn_part() RETURNS bigint LANGUAGE sql AS $$ SELECT count(*) FROM public.events_p1 $$",
					DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "events_p1"}},
				}},
				Sequences: []*metadatapb.SequenceMetadata{
					{Name: "items_id_seq", DataType: "integer", Start: "1", MinValue: "1", MaxValue: "2147483647", Increment: "1", CacheSize: "1", OwnerTable: "items", OwnerColumn: "id"},
					{Name: "custom_seq", DataType: "bigint", Start: "10", MinValue: "1", MaxValue: "9223372036854775807", Increment: "1", CacheSize: "1", LastValue: "42", SkipDump: true, Comment: "seq doc"},
				},
				ExternalTables: []*metadatapb.ExternalTableMetadata{{
					Name:               "ft_orders",
					ExternalServerName: "srv",
					Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
					CheckConstraints:   []*metadatapb.CheckConstraintMetadata{{Name: "chk_ft", Expression: "id > 0", Comment: "ft chk doc"}},
				}},
			},
			{Name: "archive"},
		},
	}

	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		ALTER TABLE public.events_p1 RENAME TO events_p1_ren;
		ALTER TABLE public.events SET SCHEMA archive;
		ALTER TABLE public.items RENAME COLUMN id TO item_id;
		ALTER TABLE public.items RENAME COLUMN qty TO quantity;
		ALTER SEQUENCE public.custom_seq RENAME TO custom_seq_ren;
		ALTER FOREIGN TABLE public.ft_orders RENAME TO ft_orders_ren;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))

	// 1 & 2 & 6: partition rename rewrites deps, preserves index comment, and parent SET SCHEMA updates ParentIndexSchema & preserves child schema.
	arch := findSchema(got, "archive")
	require.NotNil(t, arch)
	movedEvents := findTable(arch, "events")
	require.NotNil(t, movedEvents)
	part := findPartitionInSchema(got, "public", "events_p1_ren")
	require.NotNil(t, part)
	require.Equal(t, "public", part.Schema)
	require.Len(t, part.Indexes, 1)
	require.Equal(t, "part idx doc", part.Indexes[0].Comment)
	require.Equal(t, "archive", part.Indexes[0].ParentIndexSchema)
	require.Equal(t, "events_id_idx", part.Indexes[0].ParentIndexName)

	pub := findSchema(got, "public")
	require.Equal(t, "events_p1_ren", pub.Views[1].DependencyColumns[0].Table)
	require.Equal(t, "events_p1_ren", pub.Functions[0].DependencyTables[0].Table)

	// 3: Sequence rename preserves LastValue, SkipDump, Comment.
	var seq *metadatapb.SequenceMetadata
	for _, s := range pub.Sequences {
		if s.Name == "custom_seq_ren" {
			seq = s
		}
	}
	require.NotNil(t, seq)
	require.Equal(t, "42", seq.LastValue)
	require.True(t, seq.SkipDump)
	require.Equal(t, "seq doc", seq.Comment)

	// 4: Foreign table rename preserves CHECK constraint comment.
	ft := findExternalTable(pub, "ft_orders_ren")
	require.NotNil(t, ft)
	require.Len(t, ft.CheckConstraints, 1)
	require.Equal(t, "ft chk doc", ft.CheckConstraints[0].Comment)

	// 5 & 7: Column renames update view DependencyColumns and serial sequence OwnerColumn.
	vItems := pub.Views[0]
	require.Equal(t, "quantity", vItems.DependencyColumns[0].Column)
	require.Equal(t, "item_id", vItems.DependencyColumns[1].Column)
	var ownedSeq *metadatapb.SequenceMetadata
	for _, s := range pub.Sequences {
		if s.Name == "items_id_seq" {
			ownedSeq = s
		}
	}
	require.NotNil(t, ownedSeq)
	require.Equal(t, "item_id", ownedSeq.OwnerColumn)
}

func TestApplyDiffToMetadataRound9Fixes(t *testing.T) {
	// 1. DROP EXTENSION removes extension-owned members from result.
	// 2. COMMENT ON TRIGGER updates ExternalTableMetadata.Triggers.
	extName := "test_r9_ext"
	RegisterExtensionSQL(extName, `
		CREATE TABLE public.ext_tbl (id integer);
		CREATE SEQUENCE public.ext_seq;
		CREATE TYPE public.ext_enum AS ENUM ('a', 'b');
		CREATE FUNCTION public.ext_fn() RETURNS integer LANGUAGE sql AS 'SELECT 1';
	`)
	t.Cleanup(func() { delete(extensionScripts, extName) })

	catBefore := New()
	for _, stmt := range []string{
		`CREATE EXTENSION test_r9_ext;`,
		`CREATE FOREIGN DATA WRAPPER dummy;`,
		`CREATE SERVER srv FOREIGN DATA WRAPPER dummy;`,
		`CREATE FOREIGN TABLE public.ft_events (id integer) SERVER srv;`,
		`CREATE FUNCTION public.ft_trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;`,
		`CREATE TRIGGER ft_trig BEFORE INSERT ON public.ft_events FOR EACH ROW EXECUTE FUNCTION public.ft_trig_fn();`,
	} {
		_, err := catBefore.Exec(stmt, nil)
		require.NoError(t, err, stmt)
	}

	meta := &metadatapb.DatabaseSchemaMetadata{
		Name:       "postgres",
		Extensions: []*metadatapb.ExtensionMetadata{{Name: extName, Schema: "public", Version: "1.0"}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:     "ext_tbl",
				SkipDump: true,
				Columns:  []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
			}},
			Sequences: []*metadatapb.SequenceMetadata{{
				Name:      "ext_seq",
				DataType:  "bigint",
				SkipDump:  true,
				Start:     "1",
				MinValue:  "1",
				MaxValue:  "9223372036854775807",
				Increment: "1",
				CacheSize: "1",
			}},
			EnumTypes: []*metadatapb.EnumTypeMetadata{{
				Name:     "ext_enum",
				Values:   []string{"a", "b"},
				SkipDump: true,
			}},
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "ext_fn", Signature: "ext_fn()", Definition: "CREATE FUNCTION public.ext_fn() RETURNS integer LANGUAGE sql AS 'SELECT 1'", SkipDump: true},
				{Name: "ft_trig_fn", Signature: "ft_trig_fn()", Definition: "CREATE FUNCTION public.ft_trig_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$"},
			},
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "ft_events",
				ExternalServerName: "srv",
				Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				Triggers:           []*metadatapb.TriggerMetadata{{Name: "ft_trig", Body: "CREATE TRIGGER ft_trig BEFORE INSERT ON public.ft_events FOR EACH ROW EXECUTE FUNCTION public.ft_trig_fn()"}},
			}},
		}},
	}

	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		DROP EXTENSION test_r9_ext;
		COMMENT ON TRIGGER ft_trig ON public.ft_events IS 'ft trigger doc';
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
	require.Empty(t, got.Extensions)
	pub := findSchema(got, "public")
	require.NotNil(t, pub)
	require.Nil(t, findTable(pub, "ext_tbl"))
	require.Empty(t, pub.Sequences)
	require.Empty(t, pub.EnumTypes)
	require.Len(t, pub.Functions, 1)
	require.Equal(t, "ft_trig_fn", pub.Functions[0].Name)

	ft := findExternalTable(pub, "ft_events")
	require.NotNil(t, ft)
	require.Len(t, ft.Triggers, 1)
	require.Equal(t, "ft trigger doc", ft.Triggers[0].Comment)
}

func TestApplyDiffToMetadataRound10Fixes(t *testing.T) {
	// 1. Direct alterations to SkipDump extension members.
	t.Run("extension_member_direct_alteration", func(t *testing.T) {
		catBefore := New()
		results, err := catBefore.Exec(`
			CREATE EXTENSION ext_direct;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		catBefore.extensionMembers = map[uint32]string{}
		results, err = catBefore.Exec(`
			CREATE TABLE public.ext_tbl (id integer);
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		extRel := catBefore.GetRelation("public", "ext_tbl")
		require.NotNil(t, extRel)
		catBefore.extensionMembers[extRel.OID] = "ext_direct"

		meta := &metadatapb.DatabaseSchemaMetadata{
			Extensions: []*metadatapb.ExtensionMetadata{{Name: "ext_direct", Schema: "public", Version: "1.0"}},
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{{
					Name:     "ext_tbl",
					SkipDump: true,
					Columns:  []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				}},
			}},
		}

		catAfter := catBefore.Clone()
		results, err = catAfter.Exec(`ALTER TABLE public.ext_tbl ADD COLUMN extra text;`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		diff := Diff(catBefore, catAfter)
		require.True(t, diff.IsEmpty(), "extension member changes are hidden from Diff")

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, diff)
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		tbl := findTable(pub, "ext_tbl")
		require.NotNil(t, tbl)
		require.True(t, tbl.SkipDump)
		require.Len(t, tbl.Columns, 2)
		require.Equal(t, "extra", tbl.Columns[1].Name)
	})

	// 2, 3, 4, 7: View columns, column comments (view + composite), DependencyColumns, Rules on table & view.
	t.Run("views_composites_rules_and_comments", func(t *testing.T) {
		catBefore := New()
		results, err := catBefore.Exec(`
			CREATE TYPE public.addr AS (city text);
			CREATE TABLE public.src (a integer, b text);
			CREATE TABLE public.src_log (a integer);
			CREATE VIEW public.v AS SELECT a FROM public.src;
			CREATE MATERIALIZED VIEW public.mv AS SELECT a FROM public.src;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{{
					Name:       "addr",
					Attributes: []*metadatapb.CompositeTypeAttribute{{Name: "city", Type: "text"}},
				}},
				Tables: []*metadatapb.TableMetadata{
					{Name: "src", Columns: []*metadatapb.ColumnMetadata{
						{Name: "a", Type: "integer", Position: 1},
						{Name: "b", Type: "text", Position: 2},
					}},
					{Name: "src_log", Columns: []*metadatapb.ColumnMetadata{{Name: "a", Type: "integer", Position: 1}}},
				},
				Views: []*metadatapb.ViewMetadata{{
					Name:              "v",
					Definition:        "SELECT a FROM public.src",
					Columns:           []*metadatapb.ColumnMetadata{{Name: "a", Type: "integer", Position: 1}},
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "src", Column: "a"}},
				}},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
					Name:              "mv",
					Definition:        "SELECT a FROM public.src",
					DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "public", Table: "src", Column: "a"}},
				}},
			}},
		}

		catAfter := catBefore.Clone()
		results, err = catAfter.Exec(`
			COMMENT ON COLUMN public.v.a IS 'view col doc';
			COMMENT ON COLUMN public.addr.city IS 'comp attr doc';
			CREATE OR REPLACE VIEW public.v AS SELECT a, b FROM public.src;
			CREATE RULE v_ins AS ON INSERT TO public.v DO INSTEAD INSERT INTO public.src (a) VALUES (NEW.a);
			CREATE RULE src_ins AS ON INSERT TO public.src DO ALSO INSERT INTO public.src_log (a) VALUES (NEW.a);
			ALTER TABLE public.src DISABLE RULE src_ins;
			DROP MATERIALIZED VIEW public.mv;
			CREATE MATERIALIZED VIEW public.mv AS SELECT a, b FROM public.src;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Equal(t, "comp attr doc", pub.CompositeTypes[0].Attributes[0].Comment)

		v := findView(pub, "v")
		require.NotNil(t, v)
		require.Len(t, v.Columns, 2)
		require.Equal(t, "view col doc", v.Columns[0].Comment)
		require.Equal(t, "b", v.Columns[1].Name)
		require.Len(t, v.DependencyColumns, 2)
		require.Equal(t, "a", v.DependencyColumns[0].Column)
		require.Equal(t, "b", v.DependencyColumns[1].Column)
		require.Len(t, v.Rules, 1)
		require.Equal(t, "v_ins", v.Rules[0].Name)
		require.Equal(t, "INSERT", v.Rules[0].Event)
		require.True(t, v.Rules[0].IsInstead)
		require.True(t, v.Rules[0].IsEnabled)

		srcTbl := findTable(pub, "src")
		require.NotNil(t, srcTbl)
		require.Len(t, srcTbl.Rules, 1)
		require.Equal(t, "src_ins", srcTbl.Rules[0].Name)
		require.False(t, srcTbl.Rules[0].IsEnabled)

		require.Len(t, pub.MaterializedViews, 1)
		require.Len(t, pub.MaterializedViews[0].DependencyColumns, 2)
	})

	// 5 & 9: Routine COST/ROWS/WINDOW/SET and Collations.
	t.Run("routine_options_and_collations", func(t *testing.T) {
		catBefore := New()
		results, err := catBefore.Exec(`
			CREATE SCHEMA other;
			CREATE FUNCTION public.fn_cfg(x integer) RETURNS SETOF integer LANGUAGE plpgsql COST 100 ROWS 10 SET search_path = public AS $$ BEGIN RETURN NEXT x; END; $$;
			CREATE COLLATION public.my_coll (provider = icu, locale = 'und-u-ks-level2', deterministic = false);
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "public",
					Functions: []*metadatapb.FunctionMetadata{{
						Name:       "fn_cfg",
						Signature:  "public.fn_cfg(integer)",
						Definition: "CREATE FUNCTION public.fn_cfg(x integer) RETURNS SETOF integer LANGUAGE plpgsql COST 100 ROWS 10 SET search_path = public AS $$ BEGIN RETURN NEXT x; END; $$",
					}},
					Collations: []*metadatapb.CollationMetadata{{
						Name:     "my_coll",
						Provider: "i",
						Locale:   "und-u-ks-level2",
						Comment:  "coll doc",
						SkipDump: true,
					}},
				},
				{Name: "other"},
			},
		}

		catAfter := catBefore.Clone()
		results, err = catAfter.Exec(`
			ALTER FUNCTION public.fn_cfg(integer) COST 500 ROWS 25 SET search_path = other;
			CREATE FUNCTION public.fn_win(x integer) RETURNS integer LANGUAGE sql WINDOW COST 50 AS $$ SELECT x $$;
			ALTER COLLATION public.my_coll RENAME TO my_coll2;
			ALTER COLLATION public.my_coll2 SET SCHEMA other;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Empty(t, pub.Collations)

		var fnCfg, fnWin *metadatapb.FunctionMetadata
		for _, fn := range pub.Functions {
			switch fn.Name {
			case "fn_cfg":
				fnCfg = fn
			case "fn_win":
				fnWin = fn
			}
		}
		require.NotNil(t, fnCfg)
		require.Contains(t, fnCfg.Definition, "COST 500")
		require.Contains(t, fnCfg.Definition, "ROWS 25")
		require.Contains(t, fnCfg.Definition, `SET "search_path" = 'other'`)
		require.NotContains(t, fnCfg.Definition, `'public'`)

		require.NotNil(t, fnWin)
		require.Contains(t, fnWin.Definition, "WINDOW")
		require.Contains(t, fnWin.Definition, "COST 50")

		other := findSchema(got, "other")
		require.NotNil(t, other)
		require.Len(t, other.Collations, 1)
		coll := other.Collations[0]
		require.Equal(t, "my_coll2", coll.Name)
		require.Equal(t, "i", coll.Provider)
		require.Equal(t, "und-u-ks-level2", coll.Locale)
		require.NotNil(t, coll.Deterministic)
		require.False(t, *coll.Deterministic)
		require.Equal(t, "coll doc", coll.Comment)
		require.True(t, coll.SkipDump)
	})
}

func TestApplyDiffToMetadataRound11Fixes(t *testing.T) {
	// 1. SET ... FROM CURRENT resolves session GUC values instead of replaying FROM CURRENT.
	t.Run("routine_set_from_current", func(t *testing.T) {
		catBefore := New()
		catBefore.SetSearchPath([]string{"app", "public"})
		results, err := catBefore.Exec(`
			CREATE SCHEMA app;
			SET timezone = 'UTC';
			CREATE FUNCTION public.fn_cur() RETURNS integer LANGUAGE sql SET search_path FROM CURRENT SET timezone FROM CURRENT AS $$ SELECT 1 $$;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}}}
		catAfter := catBefore.Clone()
		catAfter.SetSearchPath([]string{"archive", "public"})
		results, err = catAfter.Exec(`
			SET timezone = 'Asia/Tokyo';
			CREATE FUNCTION public.fn_new() RETURNS integer LANGUAGE sql SET search_path FROM CURRENT SET timezone FROM CURRENT AS $$ SELECT 2 $$;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, New(), catBefore, Diff(New(), catBefore))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Len(t, pub.Functions, 1)
		require.Contains(t, pub.Functions[0].Definition, `SET "search_path" = 'app', 'public'`)
		require.Contains(t, pub.Functions[0].Definition, `SET "timezone" = 'UTC'`)
		require.NotContains(t, pub.Functions[0].Definition, `FROM CURRENT`)

		// Also verify GenerateMigration emits COST/ROWS/WINDOW/SET on direct diff.
		catCostBefore := New()
		_, err = catCostBefore.Exec(`CREATE FUNCTION public.fn_c() RETURNS integer LANGUAGE sql COST 10 AS $$ SELECT 1 $$;`, nil)
		require.NoError(t, err)
		catCostAfter := catCostBefore.Clone()
		_, err = catCostAfter.Exec(`ALTER FUNCTION public.fn_c() COST 5 SET search_path = public;`, nil)
		require.NoError(t, err)
		plan := GenerateMigration(catCostBefore, catCostAfter, Diff(catCostBefore, catCostAfter))
		require.Contains(t, plan.SQL(), "COST 5")
		require.Contains(t, plan.SQL(), `SET "search_path" = 'public'`)
	})

	// 2. ALTER FOREIGN TABLE ... OPTIONS and CREATE FOREIGN TABLE IF NOT EXISTS target schema lookup.
	t.Run("foreign_table_options_and_if_not_exists", func(t *testing.T) {
		catBefore := New()
		results, err := catBefore.Exec(`
			CREATE SCHEMA app;
			CREATE SCHEMA fallback;
			SET search_path = app, fallback, public;
			CREATE FOREIGN DATA WRAPPER dummy;
			CREATE SERVER srv_fallback FOREIGN DATA WRAPPER dummy;
			CREATE SERVER srv_app FOREIGN DATA WRAPPER dummy;
			CREATE FOREIGN TABLE fallback.ft (id integer) SERVER srv_fallback OPTIONS (table_name 'fb');
			CREATE FOREIGN TABLE app.ft_loaded (id integer) SERVER srv_app OPTIONS (schema_name 's0', table_name 't0');
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "app",
					ExternalTables: []*metadatapb.ExternalTableMetadata{{
						Name:               "ft_loaded",
						ExternalServerName: "srv_app",
						Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Options:            []string{"schema_name=s0", "table_name=t0"},
					}},
				},
				{
					Name: "fallback",
					ExternalTables: []*metadatapb.ExternalTableMetadata{{
						Name:               "ft",
						ExternalServerName: "srv_fallback",
						Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
						Options:            []string{"table_name=fb"},
					}},
				},
			},
		}

		catAfter := catBefore.Clone()
		results, err = catAfter.Exec(`
			ALTER FOREIGN TABLE app.ft_loaded OPTIONS (SET schema_name 's1', DROP table_name, ADD batch_size '100');
			CREATE FOREIGN TABLE IF NOT EXISTS ft (id integer) SERVER srv_app OPTIONS (table_name 'app_ft');
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		appSM := findSchema(got, "app")
		require.NotNil(t, appSM)
		ftLoaded := findExternalTable(appSM, "ft_loaded")
		require.NotNil(t, ftLoaded)
		require.Equal(t, []string{"schema_name=s1", "batch_size=100"}, ftLoaded.Options)

		ftNew := findExternalTable(appSM, "ft")
		require.NotNil(t, ftNew, "IF NOT EXISTS must check target schema rather than search_path")
		require.Equal(t, "srv_app", ftNew.ExternalServerName)
		require.Equal(t, []string{"table_name=app_ft"}, ftNew.Options)
	})
}

func TestApplyDiffToMetadataRound12Fixes(t *testing.T) {
	// 1. Newly created range with explicit collation, copied collation FROM,
	// column collation & schema on new columns, and unmodeled index key options.
	t.Run("range_collation_copied_collation_column_collation_and_index_options", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{Name: "app"}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)
		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			CREATE COLLATION app.c_copy FROM pg_catalog."C";
			CREATE COLLATION app.c_user (PROVIDER = icu, LOCALE = 'en-US');
			CREATE TYPE app.text_range AS RANGE (SUBTYPE = text, COLLATION = "C");
			CREATE TYPE app.user_range AS RANGE (SUBTYPE = text, COLLATION = app.c_user);
			CREATE TABLE app.items (
				a text COLLATE "C",
				b text COLLATE app.c_user
			);
			CREATE INDEX ix_coll ON app.items (a COLLATE "C");
			CREATE INDEX ix_nulls ON app.items (a NULLS FIRST);
			CREATE INDEX ix_desc_nulls_last ON app.items (a DESC NULLS LAST);
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		appSM := findSchema(got, "app")
		require.NotNil(t, appSM)

		copied := findCollationInSchema(appSM, "c_copy")
		require.NotNil(t, copied)
		require.Equal(t, "c", copied.Provider)
		require.Equal(t, "C", copied.LcCollate)
		require.Equal(t, "C", copied.LcCtype)
		require.NotNil(t, copied.Deterministic)
		require.True(t, *copied.Deterministic)

		tr := findRangeType(appSM, "text_range")
		require.NotNil(t, tr)
		require.Equal(t, `"C"`, tr.Collation)

		ur := findRangeType(appSM, "user_range")
		require.NotNil(t, ur)
		require.Equal(t, `"app"."c_user"`, ur.Collation)

		tbl := findTable(appSM, "items")
		require.NotNil(t, tbl)
		require.Len(t, tbl.Columns, 2)
		require.Equal(t, "C", tbl.Columns[0].Collation)
		require.Equal(t, "", tbl.Columns[0].CollationSchema)
		require.Equal(t, "c_user", tbl.Columns[1].Collation)
		require.Equal(t, "app", tbl.Columns[1].CollationSchema)

		ixColl := findIndex(tbl.Indexes, "ix_coll")
		require.NotNil(t, ixColl)
		require.Contains(t, ixColl.Definition, `COLLATE "C"`)

		ixNulls := findIndex(tbl.Indexes, "ix_nulls")
		require.NotNil(t, ixNulls)
		require.Contains(t, ixNulls.Definition, `NULLS FIRST`)

		ixDescNullsLast := findIndex(tbl.Indexes, "ix_desc_nulls_last")
		require.NotNil(t, ixDescNullsLast)
		require.Contains(t, ixDescNullsLast.Definition, `DESC NULLS LAST`)
	})

	// 2. Schema ownership changes from RENAME ROLE / REASSIGN OWNED.
	t.Run("schema_owner_from_role_operations", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{Name: "s_rename", Owner: "old_role"},
				{Name: "s_reassign", Owner: "donor_role"},
			},
		}
		catBefore := New()
		results, err := catBefore.Exec(`CREATE ROLE old_role; CREATE ROLE donor_role; CREATE ROLE target_role;`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}
		_, err = catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err = catAfter.Exec(`
			ALTER ROLE old_role RENAME TO new_role;
			REASSIGN OWNED BY donor_role TO target_role;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		diff := Diff(catBefore, catAfter)
		got := ApplyDiffToMetadata(meta, catBefore, catAfter, diff)
		require.Equal(t, "new_role", findSchema(got, "s_rename").Owner)
		require.Equal(t, "target_role", findSchema(got, "s_reassign").Owner)
	})
}

func TestApplyDiffToMetadataRound13Fixes(t *testing.T) {
	// 1 & 3. SQL-standard routine bodies and safe dollar-quoting tags.
	t.Run("sql_standard_bodies_and_dollar_quoting", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			CREATE FUNCTION public.f_ret(x int) RETURNS int LANGUAGE sql RETURN x + 1;
			CREATE FUNCTION public.f_atomic(x int) RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN x * 2; END;
			CREATE FUNCTION public.f_dollar() RETURNS text LANGUAGE sql AS $tag$ SELECT '$function$'; $tag$;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)

		var retDef, atomicDef, dollarDef string
		for _, fn := range pub.Functions {
			switch fn.Name {
			case "f_ret":
				retDef = fn.Definition
			case "f_atomic":
				atomicDef = fn.Definition
			case "f_dollar":
				dollarDef = fn.Definition
			}
		}
		require.Contains(t, retDef, "RETURN (x + 1)")
		require.NotContains(t, retDef, "$function$$function$")
		require.Contains(t, atomicDef, "BEGIN ATOMIC")
		require.Contains(t, atomicDef, "RETURN (x * 2);")
		require.NotContains(t, atomicDef, "$function$$function$")
		require.Contains(t, dollarDef, "AS $fn$ SELECT '$function$'; $fn$")

		// Re-loading the applied snapshot must succeed cleanly without syntax errors.
		reloaded := New()
		report, err := reloaded.LoadMetadata(context.Background(), got, LoadMetadataOptions{Full: true})
		require.NoError(t, err)
		require.Empty(t, report.Degraded)
		require.Empty(t, report.Missing)
	})

	// 2. Replacing a function body updates DependencyTables; metadata-only updates preserve them.
	t.Run("refresh_function_dependency_tables_on_body_replacement", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{Name: "a", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}}},
					{Name: "b", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}}},
				},
				Functions: []*metadatapb.FunctionMetadata{
					{
						Name:             "read_rel",
						Signature:        "read_rel()",
						Definition:       "CREATE FUNCTION public.read_rel() RETURNS bigint LANGUAGE sql AS $$ SELECT count(*) FROM public.a $$",
						DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "a"}},
					},
					{
						Name:             "keep_rel",
						Signature:        "keep_rel()",
						Definition:       "CREATE FUNCTION public.keep_rel() RETURNS bigint LANGUAGE sql AS $$ SELECT count(*) FROM public.a $$",
						DependencyTables: []*metadatapb.DependencyTable{{Schema: "public", Table: "a"}},
					},
				},
			}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			CREATE OR REPLACE FUNCTION public.read_rel() RETURNS bigint LANGUAGE sql AS $$ SELECT count(*) FROM public.b $$;
			ALTER FUNCTION public.keep_rel() COST 50;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		for _, fn := range pub.Functions {
			switch fn.Name {
			case "read_rel":
				require.Len(t, fn.DependencyTables, 1)
				require.Equal(t, "b", fn.DependencyTables[0].Table)
			case "keep_rel":
				require.Len(t, fn.DependencyTables, 1)
				require.Equal(t, "a", fn.DependencyTables[0].Table)
			}
		}
	})

	// 4. Temp relations/views/sequences must not enter persistent applied metadata.
	t.Run("exclude_temporary_relations_from_metadata", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			CREATE TEMP TABLE tmp_items (id serial);
			CREATE TEMP VIEW tmp_v AS SELECT 1 AS v;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		for _, sm := range got.Schemas {
			require.False(t, isTempSchemaName(sm.Name), "unexpected temp schema %q in metadata", sm.Name)
		}
		require.Len(t, got.Schemas, 1)
		require.Equal(t, "public", got.Schemas[0].Name)
	})

	// 5. Dropping routines with schema-local signatures removes them cleanly.
	t.Run("drop_routines_with_schema_local_signature", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "f",
					Signature:  "f(integer)",
					Definition: "CREATE FUNCTION public.f(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$",
				}},
				Procedures: []*metadatapb.ProcedureMetadata{{
					Name:       "p",
					Signature:  "p(integer)",
					Definition: "CREATE PROCEDURE public.p(x integer) LANGUAGE sql AS $$ SELECT 1 $$",
				}},
			}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			DROP FUNCTION public.f(integer);
			DROP PROCEDURE public.p(integer);
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Empty(t, pub.Functions)
		require.Empty(t, pub.Procedures)
	})
}

func TestApplyDiffToMetadataRound14Fixes(t *testing.T) {
	// 1. Composite attributes preserve schema-qualified collations outside search path.
	t.Run("composite_attribute_qualified_collation", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{Name: "public"},
				{
					Name:       "app",
					Collations: []*metadatapb.CollationMetadata{{Name: "c_user", Provider: "c"}},
				},
			},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			CREATE COLLATION IF NOT EXISTS app.c_user (LOCALE = 'C');
			CREATE TYPE public.addr AS (street text COLLATE app.c_user);
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Len(t, pub.CompositeTypes, 1)
		require.Equal(t, `"app"."c_user"`, pub.CompositeTypes[0].Attributes[0].Collation)
	})

	// 2. Newly created range types retain subtype_opclass, canonical, and subtype_diff.
	t.Run("range_creation_support_options", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		_, err = catAfter.Exec(`
			CREATE TYPE public.f8_range;
			CREATE FUNCTION public.canon_fn(public.f8_range) RETURNS public.f8_range LANGUAGE C IMMUTABLE AS 'internal', 'canon_fn';
			CREATE FUNCTION public.diff_fn(float8, float8) RETURNS float8 LANGUAGE sql IMMUTABLE AS 'SELECT $1 - $2';
			CREATE TYPE public.f8_range AS RANGE (SUBTYPE = float8, SUBTYPE_OPCLASS = float8_ops, CANONICAL = public.canon_fn, SUBTYPE_DIFF = public.diff_fn);
		`, nil)
		require.NoError(t, err)

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Len(t, pub.RangeTypes, 1)
		rt := pub.RangeTypes[0]
		require.Equal(t, "float8_ops", rt.SubtypeOpclass)
		require.Equal(t, "public.canon_fn", rt.Canonical)
		require.Equal(t, "public.diff_fn", rt.SubtypeDiff)
	})

	// 3. Database ownership changes update DatabaseSchemaMetadata.Owner even as the sole change.
	t.Run("database_owner_change", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Name:    "appdb",
			Owner:   "alice",
			Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			CREATE ROLE bob;
			ALTER DATABASE appdb OWNER TO bob;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		require.Equal(t, "bob", got.Owner)
	})

	// 4. Foreign data wrapper and foreign server DDL updates metadata and foreign tables.
	t.Run("fdw_and_foreign_server_ddl", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{Name: "public"},
				{
					Name: "app",
					Functions: []*metadatapb.FunctionMetadata{
						{Name: "h", Signature: "h()", Definition: "CREATE FUNCTION app.h() RETURNS fdw_handler LANGUAGE internal AS 'h'"},
						{Name: "v", Signature: "v(text[], oid)", Definition: "CREATE FUNCTION app.v(text[], oid) RETURNS void LANGUAGE internal AS 'v'"},
					},
				},
			},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			CREATE FOREIGN DATA WRAPPER custom_fdw HANDLER app.h OPTIONS (debug 'true');
			ALTER FOREIGN DATA WRAPPER custom_fdw VALIDATOR app.v OPTIONS (SET debug 'false');
			CREATE SERVER srv_one TYPE 'pg' VERSION '16' FOREIGN DATA WRAPPER custom_fdw OPTIONS (host 'db1');
			ALTER SERVER srv_one VERSION '17' OPTIONS (SET host 'db2');
			ALTER SERVER srv_one RENAME TO srv_two;
			CREATE FOREIGN TABLE public.remote_tbl (id integer) SERVER srv_two;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		require.Len(t, got.ForeignDataWrappers, 1)
		fdw := got.ForeignDataWrappers[0]
		require.Equal(t, "custom_fdw", fdw.Name)
		require.Equal(t, "app.h", fdw.Handler)
		require.Equal(t, "app.v", fdw.Validator)
		require.Equal(t, []string{"debug=false"}, fdw.Options)

		require.Len(t, got.ForeignServers, 1)
		srv := got.ForeignServers[0]
		require.Equal(t, "srv_two", srv.Name)
		require.Equal(t, "custom_fdw", srv.Fdw)
		require.NotNil(t, srv.Type)
		require.Equal(t, "pg", *srv.Type)
		require.NotNil(t, srv.Version)
		require.Equal(t, "17", *srv.Version)
		require.Equal(t, []string{"host=db2"}, srv.Options)

		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Len(t, pub.ExternalTables, 1)
		require.Equal(t, "srv_two", pub.ExternalTables[0].ExternalServerName)
	})

	// 5. Explicit empty server VERSION ''/TYPE '', cascading DROP SERVER, and foreign column OPTIONS.
	t.Run("server_empty_version_cascade_and_foreign_column_options", func(t *testing.T) {
		meta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{{
					Name:       "et_fn",
					Signature:  "et_fn()",
					Definition: "CREATE FUNCTION public.et_fn() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN END; $$",
				}},
			}},
			EventTriggers: []*metadatapb.EventTriggerMetadata{{
				Name:           "et_keep_def",
				Event:          "ddl_command_end",
				FunctionSchema: "public",
				FunctionName:   "et_fn",
				Enabled:        true,
				Definition:     "CREATE EVENT TRIGGER et_keep_def ON ddl_command_end EXECUTE FUNCTION public.et_fn()",
			}},
		}
		catBefore := New()
		_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
		require.NoError(t, err)
		_, err = catBefore.Exec(`
			CREATE EVENT TRIGGER et_keep_def ON ddl_command_end EXECUTE FUNCTION public.et_fn();
			CREATE FOREIGN DATA WRAPPER pg_fdw;
			CREATE SERVER s_drop FOREIGN DATA WRAPPER pg_fdw;
			CREATE FOREIGN TABLE public.ft_drop (id integer OPTIONS (column_name 'id')) SERVER s_drop;
			CREATE VIEW public.v_drop AS SELECT id FROM public.ft_drop;
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		meta.Schemas[0].ExternalTables = []*metadatapb.ExternalTableMetadata{{
			Name:               "ft_drop",
			ExternalServerName: "s_drop",
			Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1, Options: []string{"column_name=id"}}},
		}}
		meta.Schemas[0].Views = []*metadatapb.ViewMetadata{{
			Name:       "v_drop",
			Definition: "SELECT id FROM public.ft_drop",
		}}
		meta.ForeignServers = []*metadatapb.ForeignServerMetadata{{Name: "s_drop", Fdw: "pg_fdw"}}

		catAfter := catBefore.Clone()
		results, err := catAfter.Exec(`
			ALTER EVENT TRIGGER et_keep_def DISABLE;
			CREATE EVENT TRIGGER et_added ON sql_drop EXECUTE FUNCTION public.et_fn();
			CREATE SERVER s_empty TYPE '' VERSION '' FOREIGN DATA WRAPPER pg_fdw;
			DROP SERVER s_drop CASCADE;
			CREATE FOREIGN TABLE public.ft_keep (id integer OPTIONS (column_name 'orig')) SERVER s_empty;
			ALTER FOREIGN TABLE public.ft_keep ALTER COLUMN id OPTIONS (SET column_name 'updated', ADD schema_name 'ext');
		`, &ExecOptions{ContinueOnError: true})
		require.NoError(t, err)
		for _, r := range results {
			require.NoError(t, r.Error, r.SQL)
		}

		got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
		var keepEt, addedEt *metadatapb.EventTriggerMetadata
		for _, et := range got.EventTriggers {
			switch et.Name {
			case "et_keep_def":
				keepEt = et
			case "et_added":
				addedEt = et
			}
		}
		require.NotNil(t, keepEt)
		require.False(t, keepEt.Enabled)
		require.Equal(t, "CREATE EVENT TRIGGER et_keep_def ON ddl_command_end EXECUTE FUNCTION public.et_fn()", keepEt.Definition)
		require.NotNil(t, addedEt)
		require.Contains(t, addedEt.Definition, `CREATE EVENT TRIGGER "et_added" ON "sql_drop"`)

		require.Len(t, got.ForeignServers, 1)
		srvEmpty := got.ForeignServers[0]
		require.Equal(t, "s_empty", srvEmpty.Name)
		require.NotNil(t, srvEmpty.Type)
		require.Equal(t, "", *srvEmpty.Type)
		require.NotNil(t, srvEmpty.Version)
		require.Equal(t, "", *srvEmpty.Version)

		pub := findSchema(got, "public")
		require.NotNil(t, pub)
		require.Empty(t, pub.Views, "view depending on dropped server's foreign table must be cascaded")
		require.Len(t, pub.ExternalTables, 1)
		require.Equal(t, "ft_keep", pub.ExternalTables[0].Name)
		require.Equal(t, []string{"column_name=updated", "schema_name=ext"}, pub.ExternalTables[0].Columns[0].Options)
	})
}

func TestApplyDiffToMetadataRound16Fixes(t *testing.T) {
	cat := New()
	_, err := cat.Exec(`CREATE FOREIGN DATA WRAPPER pg_fdw;`, nil)
	require.NoError(t, err)

	// Missing FDW or missing server must be rejected rather than succeeding silently.
	res, err := cat.Exec(`CREATE SERVER bad_srv FOREIGN DATA WRAPPER missing_fdw;`, nil)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.ErrorContains(t, res[0].Error, `foreign-data wrapper "missing_fdw" does not exist`)
	require.Nil(t, cat.foreignServers["bad_srv"])
	res, err = cat.Exec(`CREATE FOREIGN TABLE public.bad_ft (id integer) SERVER missing_srv;`, nil)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.ErrorContains(t, res[0].Error, `server "missing_srv" does not exist`)
	require.Nil(t, cat.GetRelation("public", "bad_ft"))

	// VERSION NULL clears Version to nil (unlike VERSION '' which keeps non-nil ""), and SUPPORT + two-part AS survive DDL.
	v1 := "1.0"
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name:                "postgres",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{Name: "pg_fdw"}},
		ForeignServers:      []*metadatapb.ForeignServerMetadata{{Name: "srv_null", Fdw: "pg_fdw", Version: &v1}},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{
					Name:       "sup_old",
					Signature:  "sup_old(internal)",
					Definition: "CREATE FUNCTION public.sup_old(internal) RETURNS internal LANGUAGE internal STRICT AS 'sup_old'",
				},
				{
					Name:       "sup_new",
					Signature:  "sup_new(internal)",
					Definition: "CREATE FUNCTION public.sup_new(internal) RETURNS internal LANGUAGE internal STRICT AS 'sup_new'",
				},
				{
					Name:       "my_c_fn",
					Signature:  "my_c_fn()",
					Definition: "CREATE FUNCTION public.my_c_fn() RETURNS integer LANGUAGE c SUPPORT sup_old AS '$libdir/foo', 'foo_entry'",
				},
			},
		}},
	}
	catBefore := New()
	_, err = catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	catAfter := catBefore.Clone()
	results, err := catAfter.Exec(`
		CREATE SERVER srv_create_null VERSION NULL FOREIGN DATA WRAPPER pg_fdw;
		ALTER SERVER srv_null VERSION NULL;
		ALTER FUNCTION public.my_c_fn() SUPPORT public.sup_new;
		CREATE FUNCTION public.c_new() RETURNS integer LANGUAGE c AS '$libdir/bar', 'bar_sym';
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
	bySrv := make(map[string]*metadatapb.ForeignServerMetadata)
	for _, srv := range got.ForeignServers {
		bySrv[srv.Name] = srv
	}
	require.Contains(t, bySrv, "srv_null")
	require.Nil(t, bySrv["srv_null"].Version)
	require.Contains(t, bySrv, "srv_create_null")
	require.Nil(t, bySrv["srv_create_null"].Version)

	pub := findSchema(got, "public")
	require.NotNil(t, pub)
	byFn := make(map[string]*metadatapb.FunctionMetadata)
	for _, fn := range pub.Functions {
		byFn[fn.Name] = fn
	}
	require.Contains(t, byFn, "my_c_fn")
	require.Contains(t, byFn["my_c_fn"].Definition, `SUPPORT public.sup_new`)
	require.Contains(t, byFn["my_c_fn"].Definition, `AS '$libdir/foo', 'foo_entry'`)
	require.Contains(t, byFn, "c_new")
	require.Contains(t, byFn["c_new"].Definition, `AS '$libdir/bar', 'bar_sym'`)
}

func TestApplyDiffToMetadataRound17Fixes(t *testing.T) {
	// 1. Loaded snapshot catalogs enforce foreign-option action preconditions, FDW routine resolution/atomicity, and multi-server DROP atomicity.
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{{
			Name:    "pg_fdw",
			Options: []string{"debug=true"},
		}},
		ForeignServers: []*metadatapb.ForeignServerMetadata{
			{Name: "s_free", Fdw: "pg_fdw", Options: []string{"host=h1"}},
			{Name: "s_used", Fdw: "pg_fdw", Options: []string{"host=h2"}},
		},
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Functions: []*metadatapb.FunctionMetadata{
				{Name: "good_h", Signature: "good_h()", Definition: "CREATE FUNCTION public.good_h() RETURNS fdw_handler LANGUAGE internal AS 'h'"},
				{Name: "bad_h", Signature: "bad_h()", Definition: "CREATE FUNCTION public.bad_h() RETURNS integer LANGUAGE sql AS 'SELECT 1'"},
			},
			ExternalTables: []*metadatapb.ExternalTableMetadata{{
				Name:               "ft_used",
				ExternalServerName: "s_used",
				Columns:            []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
			}},
			Tables: []*metadatapb.TableMetadata{{
				Name:    "orders",
				Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:             "orders_p1",
					Expression:       "RANGE (id)",
					Value:            "FOR VALUES FROM (1) TO (10)",
					CheckConstraints: []*metadatapb.CheckConstraintMetadata{{Name: "chk_p1", Expression: "(id > 0)"}},
				}},
			}},
		}},
	}
	cat := New()
	_, err := cat.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	// Snapshot must reject SET/DROP of missing options and ADD of existing options.
	for _, badSQL := range []string{
		`ALTER SERVER s_free OPTIONS (SET missing 'x');`,
		`ALTER SERVER s_free OPTIONS (DROP missing);`,
		`ALTER SERVER s_free OPTIONS (ADD host 'dup');`,
		`ALTER FOREIGN DATA WRAPPER pg_fdw OPTIONS (SET missing 'x');`,
	} {
		res, execErr := cat.Exec(badSQL, nil)
		require.NoError(t, execErr)
		require.Len(t, res, 1)
		require.Error(t, res[0].Error, badSQL)
	}
	require.Equal(t, []string{"host=h1"}, cat.foreignServers["s_free"].Options)

	// FDW handler/validator resolution & atomic option check on ALTER FDW.
	for _, badFdw := range []string{
		`CREATE FOREIGN DATA WRAPPER missing_h_fdw HANDLER no_such_fn;`,
		`CREATE FOREIGN DATA WRAPPER wrong_ret_fdw HANDLER public.bad_h;`,
		`CREATE FOREIGN DATA WRAPPER missing_v_fdw VALIDATOR no_such_val;`,
		`ALTER FOREIGN DATA WRAPPER pg_fdw HANDLER public.good_h OPTIONS (DROP missing);`,
	} {
		res, execErr := cat.Exec(badFdw, nil)
		require.NoError(t, execErr)
		require.Len(t, res, 1)
		require.Error(t, res[0].Error, badFdw)
	}
	require.Empty(t, cat.fdws["pg_fdw"].Handler, "rejected ALTER FDW must not mutate Handler")

	// Multi-target DROP SERVER must not delete s_free when s_used has dependents.
	res, err := cat.Exec(`DROP SERVER s_free, s_used;`, nil)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Error(t, res[0].Error)
	require.NotNil(t, cat.foreignServers["s_free"], "failed multi-target DROP SERVER must leave s_free intact")

	// 2. Partition constraint comments apply to TablePartitionMetadata.
	catAfter := cat.Clone()
	res, err = catAfter.Exec(`COMMENT ON CONSTRAINT chk_p1 ON public.orders_p1 IS 'part chk doc';`, nil)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)

	got := ApplyDiffToMetadata(meta, cat, catAfter, Diff(cat, catAfter))
	pub := findSchema(got, "public")
	require.NotNil(t, pub)
	orders := findTable(pub, "orders")
	require.NotNil(t, orders)
	require.Len(t, orders.Partitions, 1)
	require.Equal(t, "part chk doc", orders.Partitions[0].CheckConstraints[0].Comment)
}

func TestApplyDiffToMetadataRound18Fixes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		ForeignDataWrappers: []*metadatapb.ForeignDataWrapperMetadata{
			{Name: "pg_fdw", Handler: ""},
		},
		ForeignServers: []*metadatapb.ForeignServerMetadata{
			{Name: "srv", Fdw: "pg_fdw"},
		},
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "good_h", Signature: "good_h()", Definition: "CREATE FUNCTION public.good_h() RETURNS fdw_handler LANGUAGE C STRICT AS 'pg_fdw', 'good_h'"},
				},
			},
			{
				Name: "hidden",
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "hidden_h", Signature: "hidden_h()", Definition: "CREATE FUNCTION hidden.hidden_h() RETURNS fdw_handler LANGUAGE C STRICT AS 'pg_fdw', 'hidden_h'"},
					{Name: "hidden_v", Signature: "hidden_v(text[], oid)", Definition: "CREATE FUNCTION hidden.hidden_v(text[], oid) RETURNS void LANGUAGE internal STRICT AS 'pg_fdw_validator'"},
				},
			},
		},
	}
	cat := New()
	cat.SetSearchPath([]string{"public"})
	_, err := cat.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	// 1. Duplicate FDW HANDLER / VALIDATOR clauses must fail without mutating existing FDW.
	// 2. Unqualified FDW handlers/validators outside search_path must fail.
	// 3. Duplicate foreign-table options in CREATE FOREIGN TABLE must fail without registering the table.
	for _, badSQL := range []string{
		`CREATE FOREIGN DATA WRAPPER dup_h_fdw HANDLER public.good_h HANDLER public.good_h;`,
		`CREATE FOREIGN DATA WRAPPER dup_v_fdw VALIDATOR hidden.hidden_v VALIDATOR hidden.hidden_v;`,
		`ALTER FOREIGN DATA WRAPPER pg_fdw HANDLER public.good_h HANDLER public.good_h;`,
		`CREATE FOREIGN DATA WRAPPER unqual_hidden_h HANDLER hidden_h;`,
		`CREATE FOREIGN DATA WRAPPER unqual_hidden_v VALIDATOR hidden_v;`,
		`CREATE FOREIGN TABLE public.bad_ft (id integer) SERVER srv OPTIONS (fmt 'csv', fmt 'json');`,
	} {
		res, execErr := cat.Exec(badSQL, nil)
		require.NoError(t, execErr)
		require.Len(t, res, 1)
		require.Error(t, res[0].Error, badSQL)
	}
	require.Empty(t, cat.fdws["pg_fdw"].Handler)
	require.Nil(t, cat.fdws["unqual_hidden_h"])
	require.Nil(t, cat.fdws["unqual_hidden_v"])
	require.Nil(t, cat.GetRelation("public", "bad_ft"))

	// Qualified reference to a routine outside search_path still succeeds.
	res, err := cat.Exec(`CREATE FOREIGN DATA WRAPPER qual_fdw HANDLER hidden.hidden_h VALIDATOR hidden.hidden_v;`, nil)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.NoError(t, res[0].Error)
}

func TestApplyDiffToMetadataRound19Fixes(t *testing.T) {
	// 1. Unqualified FDW handler stops at first matching search_path schema, FDW dependencies block RESTRICT / cascade or rename,
	// SUPPORT validates routine existence & (internal)->internal signature, and function<->procedure conversion + quoted LANGUAGE survive.
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name:       "postgres",
		SearchPath: "first_sp, second_sp, public",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "first_sp",
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "h_ambig", Signature: "h_ambig()", Definition: "CREATE FUNCTION first_sp.h_ambig() RETURNS integer LANGUAGE sql AS 'SELECT 1'"},
				},
			},
			{
				Name: "second_sp",
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "h_ambig", Signature: "h_ambig()", Definition: "CREATE FUNCTION second_sp.h_ambig() RETURNS fdw_handler LANGUAGE internal STRICT AS 'h'"},
				},
			},
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "my_h", Signature: "my_h()", Definition: "CREATE FUNCTION public.my_h() RETURNS fdw_handler LANGUAGE internal STRICT AS 'my_h'"},
					{Name: "bad_sup", Signature: "bad_sup(integer)", Definition: "CREATE FUNCTION public.bad_sup(integer) RETURNS integer LANGUAGE sql AS 'SELECT 1'"},
					{Name: "sup_ok", Signature: "sup_ok(internal)", Definition: "CREATE FUNCTION public.sup_ok(internal) RETURNS internal LANGUAGE internal STRICT AS 'sup_ok'"},
					{Name: "swap_me", Signature: "swap_me(integer)", Definition: "CREATE FUNCTION public.swap_me(x integer) RETURNS void LANGUAGE plpgsql AS 'SELECT 1'"},
				},
			},
		},
	}
	cat := New()
	cat.SetSearchPath([]string{"first_sp", "second_sp", "public"})
	_, err := cat.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	// Unqualified HANDLER h_ambig must fail on first_sp.h_ambig wrong return type instead of falling through to second_sp.h_ambig.
	// Missing/wrong-signature SUPPORT must fail without mutating existing functions.
	for _, badSQL := range []string{
		`CREATE FOREIGN DATA WRAPPER shadowed_fdw HANDLER h_ambig;`,
		`CREATE FUNCTION public.f_missing_sup(integer) RETURNS integer LANGUAGE c SUPPORT public.no_such_sup AS 'foo';`,
		`CREATE FUNCTION public.f_bad_sup(integer) RETURNS integer LANGUAGE c SUPPORT public.bad_sup AS 'foo';`,
		`ALTER FUNCTION public.swap_me(integer) COST 42 SUPPORT public.bad_sup;`,
	} {
		res, execErr := cat.Exec(badSQL, nil)
		require.NoError(t, execErr)
		require.Len(t, res, 1)
		require.Error(t, res[0].Error, badSQL)
	}
	require.Equal(t, "swap_me_", cat.findExactProc(cat.schemaByName["public"], "swap_me", []uint32{INT4OID}).Name+"_"+cat.userProcs[cat.findExactProc(cat.schemaByName["public"], "swap_me", []uint32{INT4OID}).OID].Cost)

	catAfter := cat.Clone()
	// Verify valid SUPPORT, FDW handler dependency (RESTRICT error, RENAME propagation, CASCADE drop), and function -> procedure kind conversion.
	for _, badDrop := range []string{
		`CREATE FOREIGN DATA WRAPPER f_dep HANDLER public.my_h; DROP FUNCTION public.my_h() RESTRICT;`,
	} {
		res, execErr := catAfter.Exec(badDrop, &ExecOptions{ContinueOnError: false})
		require.NoError(t, execErr)
		require.Error(t, res[len(res)-1].Error)
	}

	catAfter = cat.Clone()
	results, err := catAfter.Exec(`
		CREATE FOREIGN DATA WRAPPER f_dep HANDLER public.my_h;
		ALTER FUNCTION public.my_h() RENAME TO my_h2;
		DROP FUNCTION public.swap_me(integer);
		CREATE PROCEDURE public.swap_me(x integer) LANGUAGE plpgsql AS 'SELECT 1';
		CREATE FUNCTION public.c_sup(integer) RETURNS integer LANGUAGE c SUPPORT public.sup_ok AS '$libdir/x', 'x';
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, cat, catAfter, Diff(cat, catAfter))
	var depFdw *metadatapb.ForeignDataWrapperMetadata
	for _, f := range got.ForeignDataWrappers {
		if f.Name == "f_dep" {
			depFdw = f
		}
	}
	require.NotNil(t, depFdw)
	require.Equal(t, "public.my_h2", depFdw.Handler)

	pub := findSchema(got, "public")
	require.NotNil(t, pub)
	for _, fn := range pub.Functions {
		require.NotEqual(t, "swap_me", fn.Name, "old function must be removed when converted to procedure")
	}
	var swappedProc *metadatapb.ProcedureMetadata
	for _, p := range pub.Procedures {
		if p.Name == "swap_me" {
			swappedProc = p
		}
	}
	require.NotNil(t, swappedProc)
	require.Contains(t, swappedProc.Definition, `LANGUAGE plpgsql`)

	customLangDDL := buildUserProcDDL(cat, &UserProc{
		Schema:   cat.schemaByName["public"],
		Name:     "custom_fn",
		Kind:     'f',
		RetType:  INT4OID,
		Language: "my-lang",
		Body:     "body",
	})
	require.Contains(t, customLangDDL, `LANGUAGE "my-lang"`)

	cascadeCat := cat.Clone()
	cascadeRes, err := cascadeCat.Exec(`
		CREATE FOREIGN DATA WRAPPER f_dep HANDLER public.my_h;
		DROP FUNCTION public.my_h() CASCADE;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range cascadeRes {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Nil(t, cascadeCat.fdws["f_dep"])
}

func TestApplyDiffToMetadataRound20Fixes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "h", Signature: "h()", Definition: "CREATE FUNCTION public.h() RETURNS fdw_handler LANGUAGE internal STRICT AS 'h'"},
					{Name: "v_int", Signature: "v_int(text[], oid)", Definition: "CREATE FUNCTION public.v_int(text[], oid) RETURNS integer LANGUAGE internal STRICT AS 'v_int'"},
					{Name: "v_void", Signature: "v_void(text[], oid)", Definition: "CREATE FUNCTION public.v_void(text[], oid) RETURNS void LANGUAGE internal STRICT AS 'v_void'"},
					{Name: "sup1", Signature: "sup1(internal)", Definition: "CREATE FUNCTION public.sup1(internal) RETURNS internal LANGUAGE internal STRICT AS 'sup1'"},
					{Name: "sup2", Signature: "sup2(internal)", Definition: "CREATE FUNCTION public.sup2(internal) RETURNS internal LANGUAGE internal STRICT AS 'sup2'"},
					{Name: "caller", Signature: "caller(integer)", Definition: "CREATE FUNCTION public.caller(integer) RETURNS integer LANGUAGE sql SUPPORT public.sup1 AS 'SELECT $1'"},
				},
				Procedures: []*metadatapb.ProcedureMetadata{
					{Name: "v_proc", Signature: "v_proc(text[], oid)", Definition: "CREATE PROCEDURE public.v_proc(text[], oid) LANGUAGE plpgsql AS $$BEGIN END$$"},
				},
			},
			{Name: "other"},
		},
	}

	cat := New()
	_, err := cat.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	// 1. FDW VALIDATOR rejects non-void functions and procedures.
	for _, badSQL := range []string{
		`CREATE FOREIGN DATA WRAPPER bad_ret_fdw VALIDATOR public.v_int;`,
		`CREATE FOREIGN DATA WRAPPER bad_proc_fdw VALIDATOR public.v_proc;`,
	} {
		res, execErr := cat.Exec(badSQL, nil)
		require.NoError(t, execErr)
		require.Len(t, res, 1)
		require.Error(t, res[0].Error, badSQL)
	}

	// 2. SUPPORT dependency blocks RESTRICT drop and updates on ALTER FUNCTION ... SUPPORT + schema move + rename.
	res, err := cat.Exec(`DROP FUNCTION public.sup1(internal) RESTRICT;`, nil)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Error(t, res[0].Error)

	after := cat.Clone()
	diffSQL := `
		CREATE FOREIGN DATA WRAPPER moved_fdw HANDLER public.h VALIDATOR public.v_void;
		ALTER FUNCTION public.h() SET SCHEMA other;
		ALTER FUNCTION public.v_void(text[], oid) RENAME TO v_renamed;
		ALTER FUNCTION public.v_renamed(text[], oid) SET SCHEMA other;
		ALTER FUNCTION public.sup1(internal) SET SCHEMA other;
		ALTER FUNCTION other.sup1(internal) RENAME TO sup1_renamed;
	`
	execRes, err := after.Exec(diffSQL, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range execRes {
		require.NoError(t, r.Error, r.SQL)
	}
	applied := ApplyDiffToMetadata(meta, cat, after, Diff(cat, after))
	var movedFdw *metadatapb.ForeignDataWrapperMetadata
	for _, fdw := range applied.ForeignDataWrappers {
		if fdw.Name == "moved_fdw" {
			movedFdw = fdw
		}
	}
	require.NotNil(t, movedFdw)
	require.Equal(t, "other.h", movedFdw.Handler)
	require.Equal(t, "other.v_renamed", movedFdw.Validator)

	pub := findSchema(applied, "public")
	require.NotNil(t, pub)
	var callerFn *metadatapb.FunctionMetadata
	for _, fn := range pub.Functions {
		if fn.Name == "caller" {
			callerFn = fn
		}
	}
	require.NotNil(t, callerFn)
	require.Contains(t, callerFn.Definition, "SUPPORT other.sup1_renamed")

	// Switching SUPPORT via ALTER FUNCTION unblocks dropping old support function and drops cleanly on CASCADE.
	retargetCat := cat.Clone()
	retRes, err := retargetCat.Exec(`
		ALTER FUNCTION public.caller(integer) SUPPORT public.sup2;
		DROP FUNCTION public.sup1(internal) RESTRICT;
		DROP FUNCTION public.sup2(internal) CASCADE;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range retRes {
		require.NoError(t, r.Error, r.SQL)
	}
	require.Empty(t, retargetCat.findUserProcsByName(retargetCat.schemaByName["public"], "caller"))

	// 3. ALTER FUNCTION rejects invalid combinations atomically before updating SUPPORT or dependencies.
	atomicCat := cat.Clone()
	atomicRes, err := atomicCat.Exec(`ALTER FUNCTION public.sup2(internal) SUPPORT public.sup1 ROWS 1;`, nil)
	require.NoError(t, err)
	require.Len(t, atomicRes, 1)
	require.Error(t, atomicRes[0].Error)
	sup2Proc := atomicCat.findExactProc(atomicCat.schemaByName["public"], "sup2", []uint32{INTERNALOID})
	require.NotNil(t, sup2Proc)
	require.Empty(t, atomicCat.userProcs[sup2Proc.OID].Support)
	require.Zero(t, atomicCat.userProcs[sup2Proc.OID].SupportOID)
	dropSup1Res, err := atomicCat.Exec(`DROP FUNCTION public.sup1(internal) RESTRICT;`, nil)
	require.NoError(t, err)
	require.Len(t, dropSup1Res, 1)
	require.Error(t, dropSup1Res[0].Error) // still referenced by caller, NOT by sup2!

	// Drop caller first; sup1 must now drop cleanly if atomicCat did not leak a sup2 -> sup1 dep.
	dropCleanRes, err := atomicCat.Exec(`
		DROP FUNCTION public.caller(integer);
		DROP FUNCTION public.sup1(internal) RESTRICT;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range dropCleanRes {
		require.NoError(t, r.Error, r.SQL)
	}
}

func TestApplyDiffToMetadataImplicitPostgresOwnerOmitted(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name:    "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public", Tables: []*metadatapb.TableMetadata{{Name: "existing", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "integer", Position: 1}}}}}},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	catAfter := catBefore.Clone()
	execRes, err := catAfter.Exec(`
		CREATE TABLE public.implicit_owner (id integer);
		ALTER TABLE public.existing ADD COLUMN note text;
	`, nil)
	require.NoError(t, err)
	for _, r := range execRes {
		require.NoError(t, r.Error, r.SQL)
	}

	got := ApplyDiffToMetadata(meta, catBefore, catAfter, Diff(catBefore, catAfter))
	pub := findSchema(got, "public")
	require.NotNil(t, pub)
	require.Empty(t, findTable(pub, "implicit_owner").Owner)
	require.Empty(t, findTable(pub, "existing").Owner)
}

func TestApplyDiffToMetadataRound23Fixes(t *testing.T) {
	// 1. Detaching a partitioned partition preserves its Subpartitions under the promoted table.
	partMeta := &metadatapb.DatabaseSchemaMetadata{
		Name: "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{
			Name: "public",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "parent",
				Columns: []*metadatapb.ColumnMetadata{{Name: "a", Type: "integer", Position: 1}, {Name: "b", Type: "integer", Position: 2}},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:       "child",
					Expression: "RANGE (a)",
					Value:      "FOR VALUES FROM (0) TO (100)",
					Type:       metadatapb.TablePartitionMetadata_RANGE,
					Subpartitions: []*metadatapb.TablePartitionMetadata{{
						Name:       "leaf",
						Expression: "LIST (b)",
						Value:      "FOR VALUES IN (1)",
						Type:       metadatapb.TablePartitionMetadata_LIST,
					}},
				}},
			}},
		}},
	}
	catBefore := New()
	_, err := catBefore.LoadMetadata(context.Background(), partMeta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)
	catAfter := catBefore.Clone()
	detRes, err := catAfter.Exec(`ALTER TABLE public.parent DETACH PARTITION public.child;`, nil)
	require.NoError(t, err)
	for _, r := range detRes {
		require.NoError(t, r.Error, r.SQL)
	}
	detGot := ApplyDiffToMetadata(partMeta, catBefore, catAfter, Diff(catBefore, catAfter))
	pub := findSchema(detGot, "public")
	require.NotNil(t, pub)
	childTbl := findTable(pub, "child")
	require.NotNil(t, childTbl)
	require.Len(t, childTbl.Partitions, 1)
	require.Equal(t, "leaf", childTbl.Partitions[0].Name)

	// Reattaching preserves Subpartitions back onto the partition entry.
	catReattach := catAfter.Clone()
	attRes, err := catReattach.Exec(`ALTER TABLE public.parent ATTACH PARTITION public.child FOR VALUES FROM (0) TO (100);`, nil)
	require.NoError(t, err)
	for _, r := range attRes {
		require.NoError(t, r.Error, r.SQL)
	}
	attGot := ApplyDiffToMetadata(detGot, catAfter, catReattach, Diff(catAfter, catReattach))
	attPub := findSchema(attGot, "public")
	require.Nil(t, findTable(attPub, "child"))
	parentTbl := findTable(attPub, "parent")
	require.Len(t, parentTbl.Partitions, 1)
	require.Len(t, parentTbl.Partitions[0].Subpartitions, 1)
	require.Equal(t, "leaf", parentTbl.Partitions[0].Subpartitions[0].Name)

	// 2. Untracked SET ... FROM CURRENT is rejected instead of preserving symbolic FROM CURRENT.
	gucCat := New()
	gucRes, err := gucCat.Exec(`CREATE FUNCTION public.fn_untracked() RETURNS integer LANGUAGE sql SET work_mem FROM CURRENT AS $$ SELECT 1 $$;`, nil)
	require.NoError(t, err)
	require.Len(t, gucRes, 1)
	require.Error(t, gucRes[0].Error)
}
