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
	if len(report.Degraded) != 0 || len(report.Missing) != 0 {
		t.Fatalf("want every object installed as defined, got degraded %v, missing %v", report.Degraded, report.Missing)
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
		},
		EnumTypes: []*metadata.EnumTypeMetadata{{Name: "zz_status", Values: []string{"a"}}},
		Views: []*metadata.ViewMetadata{{
			Name:              "aa_view",
			Definition:        "SELECT id, public.zz_code(id) AS code FROM zz_orders",
			DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "zz_orders", Column: "id"}},
		}},
		Functions: []*metadata.FunctionMetadata{
			{Name: "aa_fn", Signature: "aa_fn(public.zz_base)"},
			{Name: "aa_fn", Signature: "aa_fn(public.zz_orders)"},
			{
				Name:       "aa_rows",
				Signature:  "aa_rows()",
				Definition: "CREATE FUNCTION public.aa_rows() RETURNS SETOF public.zz_orders LANGUAGE sql AS $$ SELECT * FROM zz_orders $$;",
			},
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
		Views: []*metadata.ViewMetadata{
			{
				Name:              "over_stand_in",
				Definition:        "SELECT id, geom FROM unknown_type",
				DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "unknown_type", Column: "id"}},
			},
			{
				Name:              "bad_body",
				Definition:        "SELEC nonsense",
				Columns:           []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}},
				DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "unknown_type", Column: "id"}},
			},
		},
	})

	for _, key := range []string{"type:public.dup_status", "rel:public.unknown_type", "rel:public.bad_body"} {
		if report.Degraded[key] == nil {
			t.Errorf("%s: want degraded to a stand-in, report %v", key, report.Degraded)
		}
	}
	if len(report.Missing) != 0 {
		t.Errorf("stand-ins must install, missing %v", report.Missing)
	}
	for _, key := range []string{"rel:public.uses_enum", "rel:public.over_stand_in"} {
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
	requireRelation(t, c, "public", "over_stand_in")
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
	c, report := loadSnapshot(t, false, &metadata.SchemaMetadata{
		Name: "public",
		Views: []*metadata.ViewMetadata{
			{
				Name:              "v_beta",
				Definition:        "SELECT id FROM v_alpha",
				DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v_alpha", Column: "id"}},
			},
			{
				Name:              "v_alpha",
				Definition:        "SELECT id FROM v_beta",
				DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v_beta", Column: "id"}},
			},
		},
	})
	if report.Degraded["rel:public.v_alpha"] == nil {
		t.Errorf("the cycle's lexically first member must install first and fall back, report %v", report.Degraded)
	}
	if report.Degraded["rel:public.v_beta"] != nil {
		t.Errorf("v_beta must install as defined against v_alpha's stand-in: %v", report.Degraded["rel:public.v_beta"])
	}
	requireRelation(t, c, "public", "v_alpha")
	requireRelation(t, c, "public", "v_beta")
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
			// No definition: installs from the signature.
			{Name: "fn", Signature: "fn(integer, integer)"},
		},
	})
	requireCleanReport(t, report)
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
			},
			{
				Name: "rooms",
				Columns: []*metadata.ColumnMetadata{
					{Name: "room", Type: "integer"},
					{Name: "during", Type: "tsrange"},
				},
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
		},
		Sequences: []*metadata.SequenceMetadata{
			{Name: "counter_seq", DataType: "integer", Start: "0", MinValue: "0", MaxValue: "100", Increment: "1", CacheSize: "1"},
			{Name: "users_seq_no_seq", DataType: "bigint", Increment: "5", Start: "10", Cycle: true},
			// Owned by the identity column, which creates it.
			{Name: "users_id_seq", DataType: "bigint", OwnerTable: "users", OwnerColumn: "id"},
		},
		Procedures: []*metadata.ProcedureMetadata{
			{
				Name:       "touch",
				Signature:  "touch()",
				Definition: "CREATE PROCEDURE public.touch() LANGUAGE sql AS $$ SELECT 1 $$;",
			},
			// No definition: installs from the signature, still as a procedure.
			{Name: "touch_by", Signature: "touch_by(integer)"},
		},
	}

	c, report := loadSnapshot(t, true, schema)
	requireCleanReport(t, report)

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
	rooms := requireRelation(t, c, "public", "rooms")
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
	}
	slices.Sort(sequences)
	if want := []string{"counter_seq", "users_id_seq", "users_seq_no_seq"}; !slices.Equal(sequences, want) {
		t.Errorf("sequences = %v, want the explicit ones and the identity column's", sequences)
	}
	for _, name := range []string{"touch", "touch_by"} {
		if procs := c.LookupProcByName(name); len(procs) != 1 || procs[0].Kind != 'p' {
			t.Errorf("%s: want one procedure, got %+v", name, procs)
		}
	}

	// Without Full, only the shapes install.
	c, report = loadSnapshot(t, false, schema)
	requireCleanReport(t, report)
	users = requireRelation(t, c, "public", "users")
	if len(c.IndexesOf(users.OID)) != 0 || len(c.ConstraintsOf(users.OID)) != 0 {
		t.Error("indexes and constraints must wait for Full")
	}
	if users.Columns[0].Identity != 0 || users.Columns[2].HasDefault {
		t.Error("identity and defaults must wait for Full")
	}
	if len(c.SequencesOf("public")) != 0 || len(c.LookupProcByName("touch")) != 0 {
		t.Error("sequences and procedures must wait for Full")
	}
}

func TestLoadMetadataStopsWhenContextIsDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := New().LoadMetadata(ctx, &metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{{Name: "public"}},
	}, LoadMetadataOptions{})
	if err != context.Canceled {
		t.Errorf("err = %v, want context.Canceled", err)
	}
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
			{Name: "v1", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v2"}}},
			{Name: "v2", DependencyColumns: []*metadata.DependencyColumn{{Schema: "public", Table: "v1"}}},
		},
	}
	keys := func(objects []*metadataObject) []string {
		var out []string
		for _, o := range objects {
			out = append(out, o.key())
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
