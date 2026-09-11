package catalog

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func loadSnapshot(t *testing.T, full bool, schemas ...*metadata.SchemaMetadata) (*Catalog, *LoadMetadataReport) {
	t.Helper()
	c := New()
	c.SetSearchPath([]string{"public"})
	report, err := c.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{Schemas: schemas}, LoadMetadataOptions{Full: full})
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	return c, report
}

func requireCleanReport(t *testing.T, report *LoadMetadataReport) {
	t.Helper()
	requireDegraded(t, report)
}

// requireDegraded fails unless exactly the given objects were replaced by
// stand-ins and nothing is missing.
func requireDegraded(t *testing.T, report *LoadMetadataReport, keys ...string) {
	t.Helper()
	var got []string
	for key := range report.Degraded {
		got = append(got, key)
	}
	slices.Sort(got)
	slices.Sort(keys)
	if !slices.Equal(got, keys) || len(report.Missing) != 0 {
		t.Fatalf("want degraded %v and nothing missing, got degraded %v, missing %v", keys, report.Degraded, report.Missing)
	}
}

func requireRelation(t *testing.T, c *Catalog, schema, name string) *Relation {
	t.Helper()
	rel := c.GetRelation(schema, name)
	if rel == nil {
		t.Fatalf("relation %s.%s not installed", schema, name)
	}
	return rel
}

func columnTypes(c *Catalog, rel *Relation) map[string]string {
	types := make(map[string]string, len(rel.Columns))
	for _, col := range rel.Columns {
		types[col.Name] = c.FormatType(col.TypeOID, col.TypeMod)
	}
	return types
}

func TestLoadMetadataInstallsSnapshot(t *testing.T) {
	c, report := loadSnapshot(t, false, &metadata.SchemaMetadata{
		Name:      "public",
		EnumTypes: []*metadata.EnumTypeMetadata{{Name: "task_status", Values: []string{"pending", "running"}}},
		Tables: []*metadata.TableMetadata{{
			Name: "tasks",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "integer"},
				{Name: "status", Type: "public.task_status", Nullable: true},
				{Name: "amount", Type: "numeric(10,2)", Nullable: true},
			},
			Partitions: []*metadata.TablePartitionMetadata{{
				Name:          "tasks_2024",
				Subpartitions: []*metadata.TablePartitionMetadata{{Name: "tasks_2024_q1"}},
			}},
		}},
		ExternalTables: []*metadata.ExternalTableMetadata{{
			Name:    "remote_orders",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer"}},
		}},
		Views: []*metadata.ViewMetadata{{
			Name:              "open_tasks",
			Definition:        " SELECT tasks.id, tasks.status FROM tasks;",
			DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "tasks", Column: "id"}},
		}},
	})
	requireCleanReport(t, report)
	requireRelation(t, c, "public", "remote_orders")
	for _, name := range []string{"tasks_2024", "tasks_2024_q1"} {
		if got := columnTypes(c, requireRelation(t, c, "public", name)); len(got) != 3 || got["status"] != "task_status" {
			t.Errorf("partition %s columns = %v, want the parent's", name, got)
		}
	}

	tasks := requireRelation(t, c, "public", "tasks")
	if !tasks.Columns[0].NotNull || tasks.Columns[1].NotNull {
		t.Errorf("nullability not carried over: id NotNull=%v, status NotNull=%v", tasks.Columns[0].NotNull, tasks.Columns[1].NotNull)
	}
	if got := columnTypes(c, tasks)["amount"]; got != "numeric(10,2)" {
		t.Errorf("amount type = %q, want numeric(10,2)", got)
	}
	view := requireRelation(t, c, "public", "open_tasks")
	if got := columnTypes(c, view)["status"]; got != "task_status" {
		t.Errorf("view column status type = %q, want the enum", got)
	}
}

func TestLoadMetadataInstallsDependenciesFirst(t *testing.T) {
	// Every dependency is listed after, and sorts after, the object using it.
	c, report := loadSnapshot(t, true, &metadata.SchemaMetadata{
		Name: "public",
		CompositeTypes: []*metadata.CompositeTypeMetadata{
			{Name: "aa_array_only", Attributes: []*metadata.CompositeTypeAttribute{{Name: "items", Type: "public.zz_base[]"}}},
			{Name: "aa_nested", Attributes: []*metadata.CompositeTypeAttribute{
				{Name: "home", Type: "public.zz_base"},
				{Name: "s", Type: "public.zz_status"},
			}},
			{Name: "zz_base", Attributes: []*metadata.CompositeTypeAttribute{
				{Name: "street", Type: "text", Collation: `"C"`},
				{Name: "city", Type: "character varying(50)"},
			}},
			// A table's row type.
			{Name: "aa_row_holder", Attributes: []*metadata.CompositeTypeAttribute{{Name: "o", Type: "public.zz_orders"}}},
			{Name: "aa_wraps_view", Attributes: []*metadata.CompositeTypeAttribute{{Name: "v", Type: "public.aa_calls_overload"}}},
		},
		EnumTypes: []*metadata.EnumTypeMetadata{{Name: "zz_status", Values: []string{"a"}}},
		Views: []*metadata.ViewMetadata{
			{Name: "aa_view", Definition: "SELECT id, public.zz_code(id) AS code FROM public.zz_orders"},
			{Name: "aa_count", Definition: "SELECT count(*) AS n FROM public.zz_mview"},
			{Name: "aa_uses_fn", Definition: "SELECT public.aa_view_row() AS row"},
			{Name: "aa_from_fn", Definition: "SELECT id FROM public.aa_view_row()"},
			{Name: "aa_calls_g", Definition: "SELECT public.g(1) AS r"},
			{Name: "aa_cast", Definition: "SELECT NULL::public.zz_view AS r"},
			// The CTE shares its name with the qualified view it reads.
			{Name: "aa_cte_count", Definition: "WITH zz_view AS (SELECT 1 AS id) SELECT count(*) AS n FROM public.zz_view"},
			// Calls f(integer); the other overloads take this view's row type,
			// directly or through aa_wraps_view.
			{Name: "aa_calls_overload", Definition: "SELECT public.f(1) AS x"},
			// zz_calls_h calls h(integer); the other overload takes the row type
			// of aa_reads_caller, which reads zz_calls_h.
			{Name: "aa_reads_caller", Definition: "SELECT x FROM public.zz_calls_h"},
			{Name: "zz_calls_h", Definition: "SELECT public.h(1) AS x"},
			// aa_calls_k calls k, whose only overload returns zz_calls_m's rows;
			// zz_calls_m calls m(integer), and the other m takes aa_calls_k's row type.
			{Name: "aa_calls_k", Definition: "SELECT public.k(1) AS r"},
			{Name: "zz_calls_m", Definition: "SELECT public.m(1) AS x"},
			{Name: "zz_view", Definition: "SELECT id FROM public.zz_orders"},
		},
		MaterializedViews: []*metadata.MaterializedViewMetadata{
			{Name: "zz_mview", Definition: "SELECT id FROM public.zz_orders"},
		},
		Functions: []*metadata.FunctionMetadata{
			{Name: "aa_fn", Signature: "aa_fn(public.zz_base)", Definition: "CREATE FUNCTION public.aa_fn(b public.zz_base) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;"},
			{Name: "aa_fn", Signature: "aa_fn(public.zz_orders)", Definition: "CREATE FUNCTION public.aa_fn(o public.zz_orders) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;"},
			{
				Name:       "aa_rows",
				Signature:  "aa_rows()",
				Definition: "CREATE FUNCTION public.aa_rows() RETURNS SETOF public.zz_orders LANGUAGE sql AS $$ SELECT * FROM zz_orders $$;",
			},
			{
				Name:       "aa_view_row",
				Signature:  "aa_view_row()",
				Definition: "CREATE FUNCTION public.aa_view_row() RETURNS public.zz_view LANGUAGE sql AS $$ SELECT * FROM zz_view $$;",
			},
			{Name: "f", Signature: "f(integer)", Definition: "CREATE FUNCTION public.f(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;"},
			// Both overloads return a later view's rows.
			{Name: "g", Signature: "g(integer)", Definition: "CREATE FUNCTION public.g(x integer) RETURNS public.zz_view LANGUAGE sql AS $$ SELECT * FROM zz_view $$;"},
			{Name: "g", Signature: "g(text)", Definition: "CREATE FUNCTION public.g(x text) RETURNS public.zz_view LANGUAGE sql AS $$ SELECT * FROM zz_view $$;"},
			{Name: "f", Signature: "f(public.aa_calls_overload)", Definition: "CREATE FUNCTION public.f(v public.aa_calls_overload) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;"},
			{Name: "h", Signature: "h(integer)", Definition: "CREATE FUNCTION public.h(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;"},
			{Name: "h", Signature: "h(public.aa_reads_caller)", Definition: "CREATE FUNCTION public.h(r public.aa_reads_caller) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;"},
			{Name: "k", Signature: "k(integer)", Definition: "CREATE FUNCTION public.k(x integer) RETURNS public.zz_calls_m LANGUAGE sql AS $$ SELECT 1 $$;"},
			{Name: "m", Signature: "m(integer)", Definition: "CREATE FUNCTION public.m(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT x $$;"},
			{Name: "m", Signature: "m(public.aa_calls_k)", Definition: "CREATE FUNCTION public.m(r public.aa_calls_k) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;"},
			{Name: "f", Signature: "f(public.aa_wraps_view)", Definition: "CREATE FUNCTION public.f(w public.aa_wraps_view) RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$;"},
			// Called by a default, an index expression, and a view body.
			{Name: "zz_code", Signature: "zz_code(integer)", Definition: "CREATE FUNCTION public.zz_code(x integer) RETURNS integer IMMUTABLE LANGUAGE sql AS $$ SELECT x $$;"},
		},
		Tables: []*metadata.TableMetadata{
			{
				Name: "aa_items",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "integer", Default: "nextval('public.zz_seq'::regclass)"},
					{Name: "order_id", Type: "integer"},
					{Name: "order_snapshot", Type: "public.zz_orders", Nullable: true},
					{Name: "code", Type: "integer", Default: "public.zz_code(0)"},
				},
				Indexes: []*metadata.IndexMetadata{{
					Name:       "aa_items_code_idx",
					Definition: "CREATE INDEX aa_items_code_idx ON public.aa_items USING btree (zz_code(order_id))",
				}},
				ForeignKeys: []*metadata.ForeignKeyMetadata{
					{
						Name:              "aa_items_order_fk",
						Columns:           []string{"order_id"},
						ReferencedSchema:  "public",
						ReferencedTable:   "zz_orders",
						ReferencedColumns: []string{"id"},
					},
					{
						Name:              "aa_items_ref_fk",
						Columns:           []string{"order_id"},
						ReferencedTable:   "zz_ref",
						ReferencedColumns: []string{"code"},
					},
				},
			},
			{
				Name:    "zz_ref",
				Columns: []*metadata.ColumnMetadata{{Name: "code", Type: "integer"}},
				// A standalone unique index, not a constraint, can back an FK.
				Indexes: []*metadata.IndexMetadata{{
					Name:        "zz_ref_code_key",
					Unique:      true,
					Expressions: []string{"code"},
					Definition:  "CREATE UNIQUE INDEX zz_ref_code_key ON public.zz_ref USING btree (code)",
				}},
			},
			{
				Name:    "zz_orders",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer"}},
				Indexes: []*metadata.IndexMetadata{{Name: "zz_orders_pkey", Expressions: []string{"id"}, Primary: true, Unique: true, IsConstraint: true}},
			},
		},
		Sequences: []*metadata.SequenceMetadata{{Name: "zz_seq", DataType: "bigint"}},
		Procedures: []*metadata.ProcedureMetadata{{
			Name:       "aa_proc",
			Signature:  "aa_proc()",
			Definition: "CREATE PROCEDURE public.aa_proc(OUT o public.zz_orders) LANGUAGE sql AS $$ SELECT * FROM zz_orders $$;",
		}},
	})
	requireCleanReport(t, report)

	if got := len(requireRelation(t, c, "public", "aa_array_only").Columns); got != 1 {
		t.Errorf("aa_array_only has %d attributes, want its one real attribute", got)
	}
	if got := columnTypes(c, requireRelation(t, c, "public", "aa_nested"))["home"]; got != "zz_base" {
		t.Errorf("aa_nested.home type = %q, want zz_base", got)
	}
}

func TestLoadMetadataRecordsSequenceDependencies(t *testing.T) {
	// The table sorts first; the sequence must still install before it, or the
	// default records no dependency on it.
	c, report := loadSnapshot(t, true, &metadata.SchemaMetadata{
		Name: "public",
		Tables: []*metadata.TableMetadata{{
			Name:    "aa_items",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer", Default: "nextval('zz_seq'::regclass)"}},
		}},
		Sequences: []*metadata.SequenceMetadata{{Name: "zz_seq", DataType: "integer"}},
	})
	requireCleanReport(t, report)
	results, err := c.Exec("DROP SEQUENCE public.zz_seq;", &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if results[0].Error == nil {
		t.Error("dropping a sequence a column default draws from must fail")
	}
}

func TestLoadMetadataStandsInForFailedObjects(t *testing.T) {
	c, report := loadSnapshot(t, false, &metadata.SchemaMetadata{
		Name: "public",
		// Enum labels must be unique, so this enum fails.
		EnumTypes: []*metadata.EnumTypeMetadata{{Name: "dup_status", Values: []string{"a", "a"}}},
		Tables: []*metadata.TableMetadata{
			{
				Name: "unknown_type",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "integer"},
					{Name: "geom", Type: "public.not_in_snapshot"},
				},
			},
			{
				Name:    "uses_enum",
				Columns: []*metadata.ColumnMetadata{{Name: "s", Type: "public.dup_status"}},
			},
		},
		Functions: []*metadata.FunctionMetadata{
			// Parses, but omni does not know the language.
			{Name: "js_fn", Signature: "js_fn(integer)", Definition: "CREATE FUNCTION public.js_fn(x integer) RETURNS integer LANGUAGE plv8 AS $$ return x $$;"},
			// No definition; the signature names a table that sorts after it.
			{Name: "a_sig", Signature: "a_sig(public.unknown_type)"},
		},
		Views: []*metadata.ViewMetadata{
			{Name: "over_stand_in", Definition: "SELECT id, geom FROM unknown_type"},
			{Name: "bad_body", Definition: "SELEC nonsense", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}}},
			{Name: "over_mv", Definition: "SELECT exposed FROM public.mv_alias"},
		},
		MaterializedViews: []*metadata.MaterializedViewMetadata{{
			Name:              "mv_alias",
			Definition:        "SELECT id AS exposed FROM public.unknown_type WHERE public.missing_fn(id)",
			DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "unknown_type", Column: "id"}},
		}},
	})

	for _, key := range []string{"type:public.dup_status", "rel:public.unknown_type", "rel:public.bad_body", "rel:public.mv_alias", "func:public.js_fn|js_fn(integer)", "func:public.a_sig|a_sig(public.unknown_type)"} {
		if report.Degraded[key] == nil {
			t.Errorf("%s: want degraded to a stand-in, report %v", key, report.Degraded)
		}
	}
	if len(report.Missing) != 0 {
		t.Errorf("stand-ins must install, missing %v", report.Missing)
	}
	for _, key := range []string{"rel:public.uses_enum", "rel:public.over_stand_in", "rel:public.over_mv"} {
		if report.Degraded[key] != nil {
			t.Errorf("%s depends only on stand-ins and must install as defined: %v", key, report.Degraded[key])
		}
	}
	if got := columnTypes(c, requireRelation(t, c, "public", "unknown_type")); got["id"] != "text" || got["geom"] != "text" {
		t.Errorf("table stand-in columns = %v, want the snapshot's names typed text", got)
	}
	if got := columnTypes(c, requireRelation(t, c, "public", "bad_body")); len(got) != 2 || got["a"] != "text" || got["b"] != "text" {
		t.Errorf("view stand-in columns = %v, want the view's own column names", got)
	}
	if got := columnTypes(c, requireRelation(t, c, "public", "mv_alias")); len(got) != 1 || got["exposed"] != "text" {
		t.Errorf("materialized view stand-in columns = %v, want the names its body outputs", got)
	}
	requireRelation(t, c, "public", "over_stand_in")
	for _, fn := range []string{"js_fn", "a_sig"} {
		if len(c.LookupProcByName(fn)) != 1 {
			t.Errorf("%s: want its signature stand-in installed", fn)
		}
	}
}

func TestLoadMetadataCompositeStandInKeepsAttributeNames(t *testing.T) {
	c, report := loadSnapshot(t, false, &metadata.SchemaMetadata{
		Name: "public",
		CompositeTypes: []*metadata.CompositeTypeMetadata{
			{Name: "with_domain", Attributes: []*metadata.CompositeTypeAttribute{
				{Name: "p", Type: "public.pos_int"},
				{Name: "note", Type: "text"},
			}},
		},
	})
	if report.Degraded["type:public.with_domain"] == nil {
		t.Fatalf("want with_domain degraded, report %v", report.Degraded)
	}
	rel := requireRelation(t, c, "public", "with_domain")
	if len(rel.Columns) != 2 || rel.Columns[0].Name != "p" || rel.Columns[1].Name != "note" {
		t.Errorf("stand-in attributes = %v, want p and note", columnTypes(c, rel))
	}
}

func TestLoadMetadataBreaksCycles(t *testing.T) {
	// a reads b, b reads c, and c reads a. Only a, the first member, needs a
	// stand-in: c installs against it, then b against c. v_f calls fv, which
	// returns v_f's rows; the stand-in goes to the view, which needs nothing.
	view := func(name, def string) *metadata.ViewMetadata {
		return &metadata.ViewMetadata{Name: name, Definition: def, Columns: []*metadata.ColumnMetadata{{Name: "id"}}}
	}
	c, report := loadSnapshot(t, false, &metadata.SchemaMetadata{
		Name: "public",
		Views: []*metadata.ViewMetadata{
			view("v_b", "SELECT id FROM public.v_c"),
			view("v_c", "SELECT id FROM public.v_a"),
			view("v_a", "SELECT id FROM public.v_b"),
			view("v_f", "SELECT id FROM public.fv()"),
		},
		Functions: []*metadata.FunctionMetadata{{
			Name:       "fv",
			Signature:  "fv()",
			Definition: "CREATE FUNCTION public.fv() RETURNS SETOF public.v_f LANGUAGE sql AS $$ SELECT 1 $$;",
		}},
	})
	requireDegraded(t, report, "rel:public.v_a", "rel:public.v_f")
	for _, name := range []string{"v_a", "v_b", "v_c", "v_f"} {
		requireRelation(t, c, "public", name)
	}
}

func TestLoadMetadataReportsMissingObjects(t *testing.T) {
	_, report := loadSnapshot(t, true, &metadata.SchemaMetadata{
		Name: "public",
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer"}},
			Indexes: []*metadata.IndexMetadata{{Name: "t_missing_idx", Expressions: []string{"no_such_column"}, Type: "btree"}},
		}},
		Functions: []*metadata.FunctionMetadata{{Name: "broken", Signature: "broken(not a type;;)", Definition: "not sql"}},
	})
	for _, key := range []string{"func:public.broken|broken(not a type;;)", "idx:public.t.t_missing_idx"} {
		if report.Missing[key] == nil {
			t.Errorf("%s: want missing, report %v", key, report.Missing)
		}
	}
	if len(report.Degraded) != 0 {
		t.Errorf("functions and indexes have no stand-in, degraded %v", report.Degraded)
	}
}

func TestLoadMetadataInstallsFunctionOverloads(t *testing.T) {
	c, report := loadSnapshot(t, false, &metadata.SchemaMetadata{
		Name: "public",
		Functions: []*metadata.FunctionMetadata{
			{Name: "fn", Signature: "fn(integer)", Definition: "CREATE OR REPLACE FUNCTION public.fn(x integer) RETURNS integer LANGUAGE sql AS $$ SELECT $1 $$;"},
			{Name: "fn", Signature: "fn(text)", Definition: "CREATE OR REPLACE FUNCTION public.fn(x text) RETURNS text LANGUAGE sql AS $$ SELECT $1 $$;"},
			// No definition: its signature stand-in takes its place.
			{Name: "fn", Signature: "fn(integer, integer)"},
		},
	})
	requireDegraded(t, report, "func:public.fn|fn(integer, integer)")
	if got := len(c.LookupProcByName("fn")); got != 3 {
		t.Errorf("installed %d overloads of fn, want 3", got)
	}
}

func TestLoadMetadataFull(t *testing.T) {
	schema := &metadata.SchemaMetadata{
		Name: "public",
		Tables: []*metadata.TableMetadata{
			{
				Name: "users",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "bigint", IsIdentity: true, IdentityGeneration: metadata.ColumnMetadata_ALWAYS},
					{Name: "email", Type: "text"},
					{Name: "seq_no", Type: "bigint", Default: "nextval('users_seq_no_seq'::regclass)"},
					{Name: "email_lower", Type: "text", Nullable: true, Generation: &metadata.GenerationMetadata{
						Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "lower(email)",
					}},
					{Name: "age", Type: "integer", Nullable: true},
					{Name: "code", Type: "text", Nullable: true, Collation: "C"},
					{Name: "bad_gen", Type: "text", Nullable: true, Generation: &metadata.GenerationMetadata{
						Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "lower(((",
					}},
				},
				Indexes: []*metadata.IndexMetadata{
					{Name: "users_pkey", Expressions: []string{"id"}, Primary: true, Unique: true, IsConstraint: true, Type: "btree"},
					{Name: "users_email_key", Expressions: []string{"email"}, Unique: true, IsConstraint: true, Type: "btree"},
					{Name: "users_lower_email_idx", Expressions: []string{"lower(email)", "age"}, Descending: []bool{false, true}, Type: "btree"},
					{Name: "users_adult_email_idx", Unique: true, Definition: "CREATE UNIQUE INDEX users_adult_email_idx ON public.users USING btree (email) WHERE (age > 18)"},
				},
				CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "users_age_check", Expression: "(age > 0)"}},
				// Shares the identity column, so it adds no sequence of its own.
				Partitions: []*metadata.TablePartitionMetadata{{
					Name:             "users_2024",
					Indexes:          []*metadata.IndexMetadata{{Name: "users_2024_pkey", Expressions: []string{"id"}, Primary: true, Unique: true, IsConstraint: true, Type: "btree"}},
					CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "users_age_check", Expression: "(age > 0)"}},
				}},
			},
			{
				Name: "rooms",
				Columns: []*metadata.ColumnMetadata{
					{Name: "room", Type: "integer"},
					{Name: "during", Type: "tsrange"},
				},
				// Sync lists the EXCLUDE constraint's backing index too.
				Indexes: []*metadata.IndexMetadata{{
					Name:         "rooms_no_overlap",
					IsConstraint: true,
					Type:         "gist",
					Definition:   "CREATE INDEX rooms_no_overlap ON public.rooms USING gist (room, during)",
				}},
				ExcludeConstraints: []*metadata.ExcludeConstraintMetadata{{
					Name:       "rooms_no_overlap",
					Expression: "EXCLUDE USING gist (room WITH =, during WITH &&)",
				}},
			},
			{
				Name:    "posts",
				Columns: []*metadata.ColumnMetadata{{Name: "user_id", Type: "bigint", Nullable: true}},
				ForeignKeys: []*metadata.ForeignKeyMetadata{{
					Name:              "posts_user_fk",
					Columns:           []string{"user_id"},
					ReferencedTable:   "users",
					ReferencedColumns: []string{"id"},
					OnDelete:          "CASCADE",
				}},
			},
			{
				// Degrades to a stand-in, which still owns its serial sequence.
				Name: "broken",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "integer", Default: "nextval('public.broken_id_seq'::regclass)"},
					{Name: "geom", Type: "public.not_in_snapshot"},
				},
			},
		},
		Sequences: []*metadata.SequenceMetadata{
			{Name: "broken_id_seq", DataType: "integer", OwnerTable: "broken", OwnerColumn: "id"},
			{Name: "counter_seq", DataType: "integer", Start: "0", MinValue: "0", MaxValue: "100", Increment: "1", CacheSize: "1"},
			// Owned by a plain column, as SERIAL creates it.
			{Name: "users_seq_no_seq", DataType: "bigint", Increment: "5", Start: "10", Cycle: true, OwnerTable: "users", OwnerColumn: "seq_no"},
			// Behind the identity column, renamed and tuned.
			{Name: "users_id_custom", DataType: "bigint", Start: "100", Increment: "10", OwnerTable: "users", OwnerColumn: "id"},
		},
		Procedures: []*metadata.ProcedureMetadata{
			{
				Name:       "touch",
				Signature:  "touch()",
				Definition: "CREATE PROCEDURE public.touch() LANGUAGE sql AS $$ SELECT 1 $$;",
			},
			// No definition: its signature stand-in is still a procedure.
			{Name: "touch_by", Signature: "touch_by(integer)"},
		},
		// PostgreSQL sync lists procedures among the functions.
		Functions: []*metadata.FunctionMetadata{
			{Name: "sweep", Signature: "sweep()", Definition: "CREATE PROCEDURE public.sweep() LANGUAGE sql AS $$ SELECT 1 $$;"},
			// Parses, but omni does not know the language.
			{Name: "js_proc", Signature: "js_proc(integer)", Definition: "CREATE PROCEDURE public.js_proc(x integer) LANGUAGE plv8 AS $$ return $$;"},
			// Does not parse, but still reads as a procedure.
			{Name: "cut_proc", Signature: "cut_proc(integer)", Definition: "CREATE OR REPLACE PROCEDURE public.cut_proc(x integer)\n LANGUAGE sql\nAS"},
		},
	}

	c, report := loadSnapshot(t, true, schema)
	requireDegraded(t, report, "func:public.touch_by|touch_by(integer)", "func:public.js_proc|js_proc(integer)", "func:public.cut_proc|cut_proc(integer)", "rel:public.broken")

	users := requireRelation(t, c, "public", "users")
	byName := make(map[string]*Column, len(users.Columns))
	for _, col := range users.Columns {
		byName[col.Name] = col
	}
	if byName["id"].Identity != 'a' {
		t.Errorf("id identity = %q, want ALWAYS", byName["id"].Identity)
	}
	if !byName["seq_no"].HasDefault {
		t.Error("seq_no default not installed")
	}
	if byName["email_lower"].Generated != 's' {
		t.Errorf("email_lower generated = %q, want stored", byName["email_lower"].Generated)
	}
	if byName["bad_gen"].Generated != 0 {
		t.Error("a generation expression that does not parse must be dropped, not leave the column marked generated")
	}
	if byName["code"].CollationName != "C" {
		t.Errorf("code collation = %q, want C", byName["code"].CollationName)
	}
	var indexes []string
	for _, idx := range c.IndexesOf(users.OID) {
		indexes = append(indexes, idx.Name)
		if idx.Name == "users_adult_email_idx" && idx.WhereClause == "" {
			t.Error("the partial index lost its predicate")
		}
	}
	slices.Sort(indexes)
	if want := []string{"users_adult_email_idx", "users_email_key", "users_lower_email_idx", "users_pkey"}; !slices.Equal(indexes, want) {
		t.Errorf("indexes = %v, want %v", indexes, want)
	}
	constraints := make(map[string]ConstraintType)
	for _, con := range c.ConstraintsOf(users.OID) {
		constraints[con.Name] = con.Type
	}
	if constraints["users_pkey"] != ConstraintPK || constraints["users_email_key"] != ConstraintUnique || constraints["users_age_check"] != ConstraintCheck {
		t.Errorf("users constraints = %v", constraints)
	}
	partition := requireRelation(t, c, "public", "users_2024")
	partitionConstraints := make(map[string]ConstraintType)
	for _, con := range c.ConstraintsOf(partition.OID) {
		partitionConstraints[con.Name] = con.Type
	}
	if len(partitionConstraints) != 2 || partitionConstraints["users_2024_pkey"] != ConstraintPK || partitionConstraints["users_age_check"] != ConstraintCheck {
		t.Errorf("users_2024 constraints = %v, want its own primary key and check", partitionConstraints)
	}
	rooms := requireRelation(t, c, "public", "rooms")
	if idx := c.IndexesOf(rooms.OID); len(idx) != 1 {
		t.Errorf("rooms has %d indexes, want only the EXCLUDE constraint's", len(idx))
	}
	if excl := c.ConstraintsOf(rooms.OID); len(excl) != 1 || excl[0].Type != ConstraintExclude || !slices.Equal(excl[0].ExclOps, []string{"=", "&&"}) {
		t.Errorf("rooms constraints = %+v, want the EXCLUDE with its operators", excl)
	}
	posts := requireRelation(t, c, "public", "posts")
	fks := c.ConstraintsOf(posts.OID)
	if len(fks) != 1 || fks[0].Type != ConstraintFK || fks[0].FKDelAction != 'c' {
		t.Errorf("posts constraints = %+v, want one FK with ON DELETE CASCADE", fks)
	}
	var sequences []string
	for _, seq := range c.SequencesOf("public") {
		sequences = append(sequences, seq.Name)
		if seq.Name == "counter_seq" && (seq.TypeOID != INT4OID || seq.Start != 0 || seq.MinValue != 0) {
			t.Errorf("counter_seq = %+v, want an integer sequence starting at 0", seq)
		}
		if seq.Name == "users_id_custom" && (seq.Start != 100 || seq.Increment != 10 || seq.OwnerRelOID != users.OID) {
			t.Errorf("users_id_custom = %+v, want the identity column's sequence starting at 100 by 10", seq)
		}
		if seq.Name == "users_seq_no_seq" && seq.OwnerRelOID != users.OID {
			t.Errorf("users_seq_no_seq = %+v, want it owned by users.seq_no", seq)
		}
		if seq.Name == "broken_id_seq" && seq.OwnerRelOID != requireRelation(t, c, "public", "broken").OID {
			t.Errorf("broken_id_seq = %+v, want it owned by the stand-in for broken", seq)
		}
	}
	slices.Sort(sequences)
	if want := []string{"broken_id_seq", "counter_seq", "users_id_custom", "users_seq_no_seq"}; !slices.Equal(sequences, want) {
		t.Errorf("sequences = %v, want the explicit ones and the identity column's", sequences)
	}
	for _, name := range []string{"touch", "touch_by", "sweep", "js_proc", "cut_proc"} {
		if procs := c.LookupProcByName(name); len(procs) != 1 || procs[0].Kind != 'p' {
			t.Errorf("%s: want one procedure, got %+v", name, procs)
		}
	}

	// Without Full, only the shapes install.
	c, report = loadSnapshot(t, false, schema)
	requireDegraded(t, report, "rel:public.broken")
	users = requireRelation(t, c, "public", "users")
	if len(c.IndexesOf(users.OID)) != 0 || len(c.ConstraintsOf(users.OID)) != 0 {
		t.Error("indexes and constraints must wait for Full")
	}
	if users.Columns[0].Identity != 0 || users.Columns[2].HasDefault {
		t.Error("identity and defaults must wait for Full")
	}
	if len(c.SequencesOf("public")) != 0 || len(c.LookupProcByName("touch")) != 0 || len(c.LookupProcByName("sweep")) != 0 || len(c.LookupProcByName("cut_proc")) != 0 {
		t.Error("sequences and procedures must wait for Full")
	}
}

func TestLoadMetadataStopsWhenContextIsDone(t *testing.T) {
	for _, tt := range []struct {
		name   string
		checks int
		meta   *metadata.DatabaseSchemaMetadata
	}{
		{name: "before planning", checks: 0, meta: &metadata.DatabaseSchemaMetadata{}},
		{name: "before an install", checks: 1, meta: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{Name: "public"}}}},
	} {
		if _, err := New().LoadMetadata(&doneAfter{Context: context.Background(), checks: tt.checks}, tt.meta, LoadMetadataOptions{}); err != context.Canceled {
			t.Errorf("%s: err = %v, want context.Canceled", tt.name, err)
		}
	}
}

// doneAfter is a context that is done once Err has been called checks times.
type doneAfter struct {
	context.Context
	checks int
}

func (c *doneAfter) Err() error {
	if c.checks == 0 {
		return context.Canceled
	}
	c.checks--
	return nil
}

// BenchmarkLoadMetadata loads 2,000 tables, each with a primary key, and 200
// views over them.
func BenchmarkLoadMetadata(b *testing.B) {
	schema := &metadata.SchemaMetadata{Name: "public"}
	for i := range 2000 {
		name := fmt.Sprintf("t%04d", i)
		schema.Tables = append(schema.Tables, &metadata.TableMetadata{
			Name: name,
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "bigint"},
				{Name: "name", Type: "character varying(255)", Nullable: true},
				{Name: "created_at", Type: "timestamp with time zone", Default: "now()"},
			},
			Indexes: []*metadata.IndexMetadata{{Name: name + "_pkey", Expressions: []string{"id"}, Primary: true, Unique: true, IsConstraint: true}},
		})
		if i%10 == 0 {
			schema.Views = append(schema.Views, &metadata.ViewMetadata{
				Name:              "v_" + name,
				Definition:        fmt.Sprintf(" SELECT %s.id, %s.name FROM %s;", name, name, name),
				DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: name, Column: "id"}},
			})
		}
	}
	meta := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{schema}}
	for _, full := range []bool{false, true} {
		b.Run(fmt.Sprintf("full=%v", full), func(b *testing.B) {
			for b.Loop() {
				c := New()
				c.SetSearchPath([]string{"public"})
				report, err := c.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: full})
				if err != nil || len(report.Degraded)+len(report.Missing) > 0 {
					b.Fatalf("LoadMetadata: %v, %+v", err, report)
				}
			}
		})
	}
}

func TestOrderMetadataObjectsIsDeterministic(t *testing.T) {
	schema := &metadata.SchemaMetadata{
		Name:      "public",
		EnumTypes: []*metadata.EnumTypeMetadata{{Name: "e"}},
		Tables: []*metadata.TableMetadata{
			{Name: "b", Columns: []*metadata.ColumnMetadata{{Name: "x", Type: "public.e"}}},
			{Name: "a"},
			{Name: "c"},
		},
		Views: []*metadata.ViewMetadata{
			{Name: "v1", Definition: "SELECT * FROM public.v2"},
			{Name: "v2", Definition: "SELECT * FROM public.v1"},
		},
	}
	keys := func(groups [][]*metadataObject) []string {
		var out []string
		for _, group := range groups {
			for _, o := range group {
				out = append(out, o.key())
			}
		}
		return out
	}
	objects := collectMetadataObjects(&metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{schema}}, false)
	want := keys(orderMetadataObjects(objects))
	if !slices.Equal(want, []string{"schema:public", "type:public.e", "rel:public.a", "rel:public.b", "rel:public.c", "rel:public.v1", "rel:public.v2"}) {
		t.Fatalf("order = %v", want)
	}
	r := rand.New(rand.NewSource(1))
	for range 20 {
		shuffled := slices.Clone(objects)
		r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		if got := keys(orderMetadataObjects(shuffled)); !slices.Equal(got, want) {
			t.Fatalf("order depends on input order: %v, want %v", got, want)
		}
	}
}
