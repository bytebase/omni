package catalog

import (
	"context"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

const snapshotDB = "testdb"

func loadMySQLSnapshot(t *testing.T, meta *metadata.DatabaseSchemaMetadata) (*Catalog, *LoadMetadataReport) {
	t.Helper()
	c := New()
	report, err := c.LoadMetadata(context.Background(), meta)
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	return c, report
}

// requireReport fails unless exactly the given objects were replaced by
// stand-ins and exactly the given objects are missing.
func requireReport(t *testing.T, report *LoadMetadataReport, degraded, missing []string) {
	t.Helper()
	keys := func(m map[string]error) []string {
		var out []string
		for k := range m {
			out = append(out, k)
		}
		slices.Sort(out)
		return out
	}
	slices.Sort(degraded)
	slices.Sort(missing)
	if got := keys(report.Degraded); !slices.Equal(got, degraded) {
		t.Errorf("degraded = %v, want %v (errors %v)", got, degraded, report.Degraded)
	}
	if got := keys(report.Missing); !slices.Equal(got, missing) {
		t.Errorf("missing = %v, want %v (errors %v)", got, missing, report.Missing)
	}
}

func snapshotDatabase(t *testing.T, c *Catalog) *Database {
	t.Helper()
	db := c.GetDatabase(snapshotDB)
	if db == nil {
		t.Fatalf("database %q not in catalog", snapshotDB)
	}
	return db
}

func requireTable(t *testing.T, c *Catalog, name string) *Table {
	t.Helper()
	tbl := snapshotDatabase(t, c).Tables[strings.ToLower(name)]
	if tbl == nil {
		t.Fatalf("table %q not in catalog", name)
	}
	return tbl
}

func requireView(t *testing.T, c *Catalog, name string) *View {
	t.Helper()
	v := snapshotDatabase(t, c).Views[strings.ToLower(name)]
	if v == nil {
		t.Fatalf("view %q not in catalog", name)
	}
	return v
}

func requireTrigger(t *testing.T, c *Catalog, name string) *Trigger {
	t.Helper()
	trg := snapshotDatabase(t, c).Triggers[strings.ToLower(name)]
	if trg == nil {
		t.Fatalf("trigger %q not in catalog", name)
	}
	return trg
}

func requireEvent(t *testing.T, c *Catalog, name string) *Event {
	t.Helper()
	ev := snapshotDatabase(t, c).Events[strings.ToLower(name)]
	if ev == nil {
		t.Fatalf("event %q not in catalog", name)
	}
	return ev
}

// snapshot wraps objects in the single unnamed schema MySQL sync produces.
func snapshot(s *metadata.SchemaMetadata) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Name: snapshotDB, Schemas: []*metadata.SchemaMetadata{s}}
}

func snapshotWithTables(tables ...*metadata.TableMetadata) *metadata.DatabaseSchemaMetadata {
	return snapshot(&metadata.SchemaMetadata{Tables: tables})
}

func TestLoadMetadataCreatesAndSelectsDatabase(t *testing.T) {
	c := New()
	c.SetForeignKeyChecks(true)
	if _, err := c.Exec("SET sql_generate_invisible_primary_key = ON; SET explicit_defaults_for_timestamp = OFF;", nil); err != nil {
		t.Fatalf("SET: %v", err)
	}
	meta := snapshotWithTables(&metadata.TableMetadata{
		Name:    "t",
		Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "at", Type: "timestamp"}, {Name: "ts", Type: "timestamp", Nullable: true}},
	})
	meta.CharacterSet, meta.Collation = "latin1", "latin1_bin"
	report, err := c.LoadMetadata(context.Background(), meta)
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	requireReport(t, report, nil, nil)
	// Neither a generated primary key nor the legacy TIMESTAMP rules apply.
	tbl := requireTable(t, c, "t")
	if tbl.GetColumn("my_row_id") != nil || !tbl.GetColumn("ts").Nullable || !tbl.GetColumn("ts").NullExplicit {
		t.Errorf("t = %+v, want no generated key and an explicitly nullable ts", tbl.Columns)
	}
	if at := tbl.GetColumn("at"); at.Default != nil || at.OnUpdate != "" {
		t.Errorf("at = %+v, want no implicit default", at)
	}
	if db := snapshotDatabase(t, c); db.Charset != "latin1" || db.Collation != "latin1_bin" {
		t.Errorf("database charset %q, collation %q; want the snapshot's latin1, latin1_bin", db.Charset, db.Collation)
	}
	if c.CurrentDatabase() != snapshotDB {
		t.Errorf("current database = %q, want %q", c.CurrentDatabase(), snapshotDB)
	}
	if !c.ForeignKeyChecks() || !c.generateGIPK || c.session.ExplicitDefaultsForTimestamp {
		t.Error("the caller's session settings must be back after the load")
	}
	// Loading into an existing database adds to it and takes the snapshot's
	// defaults.
	meta = snapshotWithTables(&metadata.TableMetadata{
		Name:    "u",
		Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}},
	})
	meta.CharacterSet, meta.Collation = "utf8mb4", "utf8mb4_bin"
	if _, err := c.LoadMetadata(context.Background(), meta); err != nil {
		t.Fatalf("second LoadMetadata: %v", err)
	}
	requireTable(t, c, "t")
	requireTable(t, c, "u")
	if db := snapshotDatabase(t, c); db.Charset != "utf8mb4" || db.Collation != "utf8mb4_bin" {
		t.Errorf("database charset %q, collation %q; want the second snapshot's utf8mb4, utf8mb4_bin", db.Charset, db.Collation)
	}
}

func TestLoadMetadataInstallsTable(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshotWithTables(&metadata.TableMetadata{
		Name:          "users",
		Engine:        "InnoDB",
		Charset:       "utf8mb4",
		Collation:     "utf8mb4_0900_ai_ci",
		Comment:       "users table",
		CreateOptions: "row_format=COMPRESSED stats_persistent=0",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id", Type: "bigint unsigned", Default: autoIncrementDefault},
			{Name: "email", Type: "varchar(255)", CharacterSet: "latin1"},
			{Name: "created_at", Type: "timestamp", Default: "CURRENT_TIMESTAMP", OnUpdate: "CURRENT_TIMESTAMP"},
			{Name: "bio", Type: "text", Nullable: true, Comment: "free form text"},
			{Name: "status", Type: "varchar(16)", Default: "'active'"},
			{Name: "score", Type: "int", Default: "0"},
			{Name: "seed", Type: "double", Default: "(rand())"},
		},
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}},
			{Name: "uniq_email", Type: "BTREE", Unique: true, Expressions: []string{"email"}},
		},
	}))
	requireReport(t, report, nil, nil)

	tbl := requireTable(t, c, "users")
	if len(tbl.Columns) != 7 || tbl.Engine != "InnoDB" || tbl.Charset != "utf8mb4" || tbl.Comment != "users table" {
		t.Errorf("table = %d columns, engine %q, charset %q, comment %q", len(tbl.Columns), tbl.Engine, tbl.Charset, tbl.Comment)
	}
	id := tbl.GetColumn("id")
	if !id.AutoIncrement || id.Nullable || !strings.EqualFold(id.DataType, "bigint") {
		t.Errorf("id = auto increment %v, nullable %v, type %q", id.AutoIncrement, id.Nullable, id.DataType)
	}
	if tbl.RowFormat != "COMPRESSED" {
		t.Errorf("row format = %q, want the declared COMPRESSED", tbl.RowFormat)
	}
	if email := tbl.GetColumn("email"); email.Charset != "latin1" {
		t.Errorf("email charset = %q, want the column's latin1", email.Charset)
	}
	// The catalog may print CURRENT_TIMESTAMP as now(), so only presence counts.
	if created := tbl.GetColumn("created_at"); created.Default == nil || *created.Default == "" || created.OnUpdate == "" {
		t.Errorf("created_at lost its default or ON UPDATE: %+v", created)
	}
	if bio := tbl.GetColumn("bio"); bio.Comment != "free form text" {
		t.Errorf("bio comment = %q", bio.Comment)
	}
	// A literal default stays a constant; a parenthesized one stays an expression.
	for name, constant := range map[string]bool{"status": true, "score": true, "seed": false} {
		if col := tbl.GetColumn(name); (col.DefaultKind == ColumnDefaultConstant) != constant {
			t.Errorf("%s default kind = %v, want constant %v", name, col.DefaultKind, constant)
		}
	}
	var hasPK, hasUnique bool
	for _, con := range tbl.Constraints {
		hasPK = hasPK || con.Type == ConPrimaryKey
		hasUnique = hasUnique || (con.Type == ConUniqueKey && con.Name == "uniq_email")
	}
	if !hasPK || !hasUnique {
		t.Errorf("primary key %v, uniq_email %v; want both", hasPK, hasUnique)
	}
}

func TestLoadMetadataInstallsIndexKinds(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshotWithTables(&metadata.TableMetadata{
		Name: "t",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id", Type: "int", Default: autoIncrementDefault},
			{Name: "body", Type: "text", Nullable: true},
			{Name: "location", Type: "geometry"},
			{Name: "tag", Type: "varchar(64)", Nullable: true},
			{Name: "name", Type: "varchar(255)", Nullable: true},
		},
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}},
			{Name: "ft_body", Type: "FULLTEXT", Expressions: []string{"body"}},
			{Name: "sp_location", Type: "SPATIAL", Expressions: []string{"location"}},
			// Functional key parts, with and without outer parentheses, alone
			// and next to a plain column.
			{Name: "idx_abs", Type: "BTREE", Expressions: []string{"(abs(id))"}},
			{Name: "idx_lower_name", Type: "BTREE", Expressions: []string{"lower(`name`)"}},
			{Name: "idx_mixed", Type: "BTREE", Expressions: []string{"tag", "(lower(`name`))"}},
			{Name: "idx_prefix", Type: "BTREE", Expressions: []string{"name"}, KeyLength: []int64{10}, Descending: []bool{true}},
			{Name: "idx_shown", Type: "BTREE", Expressions: []string{"tag"}, Visible: true},
			{Name: "idx_hidden", Type: "BTREE", Expressions: []string{"tag"}, Comment: "rarely used"},
		},
	}))
	requireReport(t, report, nil, nil)

	indexes := make(map[string]*Index)
	for _, idx := range requireTable(t, c, "t").Indexes {
		indexes[idx.Name] = idx
	}
	for _, name := range []string{"ft_body", "sp_location", "idx_abs", "idx_lower_name", "idx_mixed", "idx_prefix"} {
		if indexes[name] == nil {
			t.Errorf("index %s not installed", name)
		}
	}
	if idx := indexes["ft_body"]; idx != nil && !idx.Fulltext {
		t.Error("ft_body is not FULLTEXT")
	}
	if idx := indexes["sp_location"]; idx != nil && !idx.Spatial {
		t.Error("sp_location is not SPATIAL")
	}
	if idx := indexes["idx_mixed"]; idx != nil && len(idx.Columns) != 2 {
		t.Errorf("idx_mixed has %d key parts, want 2", len(idx.Columns))
	}
	if idx := indexes["idx_abs"]; idx != nil && (len(idx.Columns) != 1 || !strings.Contains(strings.ToLower(idx.Columns[0].Expr), "abs(")) {
		t.Errorf("idx_abs key parts = %+v, want the abs() expression", idx.Columns)
	}
	if idx := indexes["idx_prefix"]; idx != nil && (len(idx.Columns) != 1 || idx.Columns[0].Length != 10 || !idx.Columns[0].Descending) {
		t.Errorf("idx_prefix key parts = %+v, want name(10) DESC", idx.Columns)
	}
	// The primary key is visible though the snapshot does not say so.
	for name, visible := range map[string]bool{"PRIMARY": true, "idx_shown": true, "idx_hidden": false} {
		if idx := indexes[name]; idx == nil || idx.Visible != visible {
			t.Errorf("%s = %+v, want visible %v", name, idx, visible)
		}
	}
	if idx := indexes["idx_hidden"]; idx == nil || idx.Comment != "rarely used" {
		t.Errorf("idx_hidden = %+v, want its comment", idx)
	}
}

func TestLoadMetadataLeavesDefaultIndexTypesUnwritten(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshotWithTables(
		&metadata.TableMetadata{Name: "disk", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}}, Indexes: []*metadata.IndexMetadata{
			{Name: "k", Type: "BTREE", Expressions: []string{"a"}},
		}},
		&metadata.TableMetadata{Name: "mem", Engine: "MEMORY", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}}, Indexes: []*metadata.IndexMetadata{
			{Name: "k_hash", Type: "HASH", Expressions: []string{"a"}},
			{Name: "k_btree", Type: "BTREE", Expressions: []string{"a"}},
		}},
	))
	requireReport(t, report, nil, nil)
	types := make(map[string]string)
	for _, name := range []string{"disk", "mem"} {
		for _, idx := range requireTable(t, c, name).Indexes {
			types[name+"."+idx.Name] = idx.IndexType
		}
	}
	if want := map[string]string{"disk.k": "", "mem.k_hash": "", "mem.k_btree": "BTREE"}; !maps.Equal(types, want) {
		t.Errorf("index types = %v, want %v", types, want)
	}
}

func TestLoadMetadataInstallsForeignKeysInAnyOrder(t *testing.T) {
	// The child sorts first and references the parent.
	c, report := loadMySQLSnapshot(t, snapshotWithTables(
		&metadata.TableMetadata{
			Name: "a_child",
			Columns: []*metadata.ColumnMetadata{
				{Name: "id", Type: "bigint unsigned", Default: autoIncrementDefault},
				{Name: "parent_id", Type: "bigint unsigned"},
			},
			Indexes: []*metadata.IndexMetadata{
				{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}},
				{Name: "idx_parent_id", Type: "BTREE", Expressions: []string{"parent_id"}},
			},
			ForeignKeys: []*metadata.ForeignKeyMetadata{{
				Name:              "fk_child_parent",
				Columns:           []string{"parent_id"},
				ReferencedTable:   "z_parent",
				ReferencedColumns: []string{"id"},
				OnDelete:          "CASCADE",
				OnUpdate:          "RESTRICT",
			}},
		},
		&metadata.TableMetadata{
			Name:    "z_parent",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "bigint unsigned", Default: autoIncrementDefault}},
			Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}}},
		},
	))
	requireReport(t, report, nil, nil)

	var fk *Constraint
	for _, con := range requireTable(t, c, "a_child").Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_child_parent" {
			fk = con
		}
	}
	if fk == nil {
		t.Fatal("fk_child_parent not installed")
	}
	if !strings.EqualFold(fk.OnDelete, "CASCADE") || !strings.EqualFold(fk.OnUpdate, "RESTRICT") {
		t.Errorf("fk actions = ON DELETE %q ON UPDATE %q", fk.OnDelete, fk.OnUpdate)
	}
}

func TestLoadMetadataInstallsColumnValues(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshotWithTables(&metadata.TableMetadata{
		Name: "t",
		Columns: []*metadata.ColumnMetadata{
			{Name: "a", Type: "int", Nullable: true},
			{Name: "b", Type: "int", Nullable: true},
			{Name: "c", Type: "int", Nullable: true, Generation: &metadata.GenerationMetadata{
				Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "a + b",
			}},
			// Sync reports DEFAULT NULL for every nullable column, including
			// types that show no default.
			{Name: "payload_json", Type: "json", Nullable: true, Default: "NULL"},
			{Name: "blob_body", Type: "blob", Nullable: true, Default: "NULL"},
			{Name: "location", Type: "geometry", Nullable: true, Default: "NULL"},
			{Name: "notes", Type: "text", Nullable: true, Default: "NULL"},
			{Name: "spot", Type: "point", Nullable: true, Default: "NULL"},
			{Name: "name", Type: "varchar(255)", Nullable: true, Default: "NULL"},
		},
	}))
	requireReport(t, report, nil, nil)

	tbl := requireTable(t, c, "t")
	if gen := tbl.GetColumn("c").Generated; gen == nil || !gen.Stored || !strings.Contains(gen.Expr, "a") {
		t.Errorf("generated column = %+v, want stored a + b", gen)
	}
	for _, name := range []string{"payload_json", "blob_body", "location", "notes", "spot"} {
		if col := tbl.GetColumn(name); col.Default != nil {
			t.Errorf("%s: default = %q, want none", name, *col.Default)
		}
	}
	if col := tbl.GetColumn("name"); col.Default == nil {
		t.Error("name lost its DEFAULT NULL")
	}
}

func TestLoadMetadataInstallsViewsInAnyOrder(t *testing.T) {
	// Each view reads the next, which installs after it, and m_view's columns
	// come only from analyzing its body.
	c, report := loadMySQLSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{
			Name:    "users",
			Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}, {Name: "name", Type: "varchar(255)", Nullable: true}},
		}},
		Views: []*metadata.ViewMetadata{
			{Name: "a_view", Definition: "SELECT id FROM m_view"},
			{Name: "m_view", Definition: "SELECT * FROM z_view"},
			{Name: "z_view", Definition: "SELECT id, name FROM users"},
		},
	}))
	requireReport(t, report, nil, nil)
	// The columns come from analyzing the parsed body.
	if v := requireView(t, c, "z_view"); !strings.Contains(strings.ToUpper(v.Definition), "SELECT") || !slices.Equal(v.Columns, []string{"id", "name"}) {
		t.Errorf("z_view = definition %q, columns %v; want its body and id, name", v.Definition, v.Columns)
	}
	for _, name := range []string{"a_view", "m_view"} {
		if v := requireView(t, c, name); v.AnalyzedQuery == nil {
			t.Errorf("%s was not analyzed once the view it reads installed", name)
		}
	}
}

func TestLoadMetadataKeepsViewColumnListsOnlyWhereNeeded(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{Name: "users", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}},
		Views: []*metadata.ViewMetadata{
			// Created as CREATE VIEW renamed (user_id) AS SELECT id FROM users.
			{Name: "renamed", Definition: "SELECT id FROM users", Columns: []*metadata.ColumnMetadata{{Name: "user_id"}}},
			{Name: "same", Definition: "SELECT id FROM users", Columns: []*metadata.ColumnMetadata{{Name: "ID"}}},
			{Name: "star", Definition: "SELECT * FROM users", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
			{Name: "star_qualified", Definition: "SELECT u.* FROM users u", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
			// Created as CREATE VIEW star_renamed (user_id) AS SELECT * FROM users.
			{Name: "star_renamed", Definition: "SELECT * FROM users", Columns: []*metadata.ColumnMetadata{{Name: "user_id"}}},
			// A body that does not resolve names nothing of its own. A qualified
			// star leaves one unnamed column behind rather than none.
			{Name: "unresolved", Definition: "SELECT * FROM missing", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}}},
			{Name: "unresolved_qualified", Definition: "SELECT m.* FROM missing m", Columns: []*metadata.ColumnMetadata{{Name: "a"}}},
		},
	}))
	requireReport(t, report, nil, nil)
	for _, name := range []string{"renamed", "star_renamed"} {
		if v := requireView(t, c, name); !v.ExplicitColumns || !slices.Equal(v.Columns, []string{"user_id"}) {
			t.Errorf("%s = explicit %v, columns %v; want the snapshot's user_id", name, v.ExplicitColumns, v.Columns)
		}
	}
	// A body naming its own columns keeps them; one that does not resolve takes
	// the snapshot's, inferred.
	for name, want := range map[string][]string{"same": {"id"}, "star": {"id"}, "star_qualified": {"id"}, "unresolved": {"a", "b"}, "unresolved_qualified": {"a"}} {
		if v := requireView(t, c, name); v.ExplicitColumns || !slices.Equal(v.Columns, want) {
			t.Errorf("%s = explicit %v, columns %v; want %v inferred", name, v.ExplicitColumns, v.Columns, want)
		}
	}
	// Whatever a view ends up advertising, its inferred column metadata answers
	// to the same names.
	for _, name := range []string{"renamed", "same", "star", "star_qualified", "star_renamed", "unresolved", "unresolved_qualified"} {
		v := requireView(t, c, name)
		var inferred []string
		for _, col := range v.ColumnMetadata {
			inferred = append(inferred, col.Name)
		}
		if !slices.Equal(inferred, v.Columns) {
			t.Errorf("%s = columns %v, column metadata %v; want them named alike", name, v.Columns, inferred)
		}
	}
}

func TestLoadMetadataDropsViewColumnMetadataAnUnresolvedStarShifts(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{Name: "users", Columns: []*metadata.ColumnMetadata{{Name: "name", Type: "varchar(8)", Nullable: true}}}},
		Views: []*metadata.ViewMetadata{
			{Name: "shifted", Definition: "SELECT m.*, u.name FROM missing m, users u", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "name"}}},
			{Name: "aligned", Definition: "SELECT u.name FROM users u", Columns: []*metadata.ColumnMetadata{{Name: "label"}}},
		},
	}))
	requireReport(t, report, nil, nil)
	// A star the catalog cannot expand shifts every column after it, so what the
	// body inferred belongs to none of the names the snapshot gives.
	for i, col := range requireView(t, c, "shifted").ColumnMetadata {
		if col.Collation != "" || col.Nullable {
			t.Errorf("shifted column %d (%s) = nullable %v, collation %q; want nothing carried over", i, col.Name, col.Nullable, col.Collation)
		}
	}
	// A body with a column for each of the snapshot's keeps what it inferred.
	if col := requireView(t, c, "aligned").ColumnMetadata[0]; col.Name != "label" || !col.Nullable || col.Collation == "" {
		t.Errorf("aligned column = %+v; want users.name's inferred nullability and collation under label", col)
	}
}

func TestLoadMetadataInstallsTriggers(t *testing.T) {
	account := func(triggers ...*metadata.TriggerMetadata) *metadata.TableMetadata {
		return &metadata.TableMetadata{
			Name:     "account",
			Columns:  []*metadata.ColumnMetadata{{Name: "acct_num", Type: "int", Nullable: true}, {Name: "amount", Type: "decimal(10,2)", Nullable: true}},
			Triggers: triggers,
		}
	}
	c, report := loadMySQLSnapshot(t, snapshotWithTables(account(
		&metadata.TriggerMetadata{Name: "ins_sum", Timing: "BEFORE", Event: "INSERT", Body: "SET @sum = @sum + NEW.amount", SqlMode: "ANSI_QUOTES"},
		&metadata.TriggerMetadata{Name: "upd_check", Timing: "BEFORE", Event: "UPDATE", Body: `BEGIN
			IF NEW.amount < 0 THEN
				SET NEW.amount = 0;
			ELSEIF NEW.amount > 100 THEN
				SET NEW.amount = 100;
			END IF;
		END`},
		&metadata.TriggerMetadata{Name: "validate_amount", Timing: "AFTER", Event: "UPDATE", Body: `BEGIN
			IF NEW.amount < 0 THEN
				SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Negative amount not allowed';
			END IF;
		END`},
		// A body that does not parse is kept as text.
		&metadata.TriggerMetadata{Name: "bad_body", Timing: "BEFORE", Event: "INSERT", Body: "this is not sql $$$"},
	)))
	requireReport(t, report, nil, nil)

	for _, tt := range []struct{ name, timing, event, body string }{
		{"ins_sum", "BEFORE", "INSERT", "@SUM"},
		{"upd_check", "BEFORE", "UPDATE", "IF NEW.AMOUNT"},
		{"validate_amount", "AFTER", "UPDATE", "SIGNAL"},
		{"bad_body", "BEFORE", "INSERT", "NOT SQL"},
	} {
		trg := requireTrigger(t, c, tt.name)
		if !strings.EqualFold(trg.Timing, tt.timing) || !strings.EqualFold(trg.Event, tt.event) || !strings.EqualFold(trg.Table, "account") || !strings.Contains(strings.ToUpper(trg.Body), tt.body) {
			t.Errorf("%s = %s %s ON %s, body %q", tt.name, trg.Timing, trg.Event, trg.Table, trg.Body)
		}
	}
	if trg := requireTrigger(t, c, "ins_sum"); !trg.HasSessionContext || trg.SQLMode != "ANSI_QUOTES" {
		t.Errorf("ins_sum sql_mode = %q (context %v), want the snapshot's ANSI_QUOTES", trg.SQLMode, trg.HasSessionContext)
	}
}

func TestLoadMetadataInstallsEvents(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshot(&metadata.SchemaMetadata{Events: []*metadata.EventMetadata{
		{Name: "myevent", Definition: "CREATE EVENT `myevent` ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 HOUR DO UPDATE mytable SET mycol = mycol + 1"},
		{Name: "e_hourly", Definition: "CREATE EVENT `e_hourly` ON SCHEDULE EVERY 1 HOUR COMMENT 'Clears out sessions table each hour.' DO DELETE FROM site_activity.sessions", TimeZone: "+08:00"},
		{Name: "e_daily", Definition: "CREATE EVENT `e_daily` ON SCHEDULE EVERY 1 DAY STARTS CURRENT_TIMESTAMP + INTERVAL 5 HOUR COMMENT 'Saves total then clears each day' " +
			"DO BEGIN INSERT INTO totals (time, total) SELECT NOW(), COUNT(*) FROM sessions; DELETE FROM sessions; END"},
	}}))
	requireReport(t, report, nil, nil)

	for name, schedule := range map[string]string{"myevent": "INTERVAL", "e_hourly": "EVERY", "e_daily": "STARTS"} {
		if ev := requireEvent(t, c, name); !strings.Contains(strings.ToUpper(ev.Schedule), schedule) {
			t.Errorf("%s schedule = %q, want %s in it", name, ev.Schedule, schedule)
		}
	}
	if ev := requireEvent(t, c, "e_hourly"); !strings.Contains(ev.Comment, "sessions") {
		t.Errorf("e_hourly comment = %q", ev.Comment)
	}
	if ev := requireEvent(t, c, "e_hourly"); !ev.HasSessionContext || ev.TimeZone != "+08:00" {
		t.Errorf("e_hourly time zone = %q (context %v), want the snapshot's +08:00", ev.TimeZone, ev.HasSessionContext)
	}
}

func TestLoadMetadataInstallsRoutines(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Functions: []*metadata.FunctionMetadata{
			{Name: "add_one", Definition: "CREATE FUNCTION `add_one`(x INT) RETURNS int DETERMINISTIC RETURN x + 1", SqlMode: "PIPES_AS_CONCAT", CharacterSetClient: "utf8mb4", CollationConnection: "utf8mb4_0900_ai_ci"},
			{Name: "broken_fn", Definition: "not a function"},
		},
		Procedures: []*metadata.ProcedureMetadata{
			{Name: "touch", Definition: "CREATE PROCEDURE `touch`() BEGIN SELECT 1; END"},
		},
	}))
	// A function or procedure that fails has no stand-in.
	requireReport(t, report, nil, []string{"function:broken_fn"})
	db := snapshotDatabase(t, c)
	if db.Functions[strings.ToLower("add_one")] == nil || db.Procedures[strings.ToLower("touch")] == nil {
		t.Errorf("functions %v, procedures %v; want add_one and touch", db.Functions, db.Procedures)
	}
	// A routine takes the context the snapshot records; one with none stays bare.
	if fn := db.Functions["add_one"]; fn == nil || !fn.HasSessionContext || fn.SQLMode != "PIPES_AS_CONCAT" || fn.CollationConnection != "utf8mb4_0900_ai_ci" {
		t.Errorf("add_one = %+v, want the snapshot's session context", fn)
	}
	if p := db.Procedures["touch"]; p == nil || p.HasSessionContext {
		t.Errorf("touch = %+v, want no session context", p)
	}
}

func TestLoadMetadataStampsOnlyObjectsItInstalls(t *testing.T) {
	c := New()
	// Same-named functions in another database, and already in the destination.
	if _, err := c.Exec("CREATE DATABASE other; USE other; CREATE FUNCTION add_one(x INT) RETURNS int DETERMINISTIC RETURN x + 1; CREATE DATABASE testdb; USE testdb; CREATE FUNCTION dup(x INT) RETURNS int DETERMINISTIC RETURN x;", nil); err != nil {
		t.Fatalf("setup: %v", err)
	}
	report, err := c.LoadMetadata(context.Background(), snapshot(&metadata.SchemaMetadata{Functions: []*metadata.FunctionMetadata{
		{Name: "add_one", Definition: "CREATE FUNCTION `add_one`(x INT) RETURNS int DETERMINISTIC RETURN x + 1", SqlMode: "PIPES_AS_CONCAT"},
		{Name: "dup", Definition: "CREATE FUNCTION `dup`(x INT) RETURNS int DETERMINISTIC RETURN x", SqlMode: "PIPES_AS_CONCAT"},
	}}))
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	requireReport(t, report, nil, []string{"function:dup"})
	db := snapshotDatabase(t, c)
	if fn := db.Functions["add_one"]; fn == nil || fn.SQLMode != "PIPES_AS_CONCAT" {
		t.Errorf("testdb add_one = %+v, want the snapshot's context", fn)
	}
	for name, fn := range map[string]*Routine{"other add_one": c.GetDatabase("other").Functions["add_one"], "testdb dup": db.Functions["dup"]} {
		if fn == nil || fn.HasSessionContext {
			t.Errorf("%s = %+v, want its own context untouched", name, fn)
		}
	}
}

func TestLoadMetadataStandsInForFailedObjects(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{
			{
				Name: "broken",
				Columns: []*metadata.ColumnMetadata{
					{Name: "id", Type: "int"},
					{Name: "bad", Type: "varchar((("},
				},
				Triggers: []*metadata.TriggerMetadata{
					// Attaches to the stand-in.
					{Name: "ins_stub", Timing: "BEFORE", Event: "INSERT", Body: "SET @x = NEW.id"},
					// No timing or event: the stand-in takes defaults.
					{Name: "trg_needs_defaults"},
				},
			},
		},
		Views: []*metadata.ViewMetadata{
			{Name: "bad_view", Definition: "SELEC nonsense", Columns: []*metadata.ColumnMetadata{{Name: "x"}}},
		},
		Events: []*metadata.EventMetadata{
			{Name: "e_empty"},
			{Name: "e_broken", Definition: "this is definitely not a CREATE EVENT statement $$"},
		},
	}))
	requireReport(t, report, []string{"table:broken", "trigger:broken.trg_needs_defaults", "view:bad_view", "event:e_empty", "event:e_broken"}, nil)

	tbl := requireTable(t, c, "broken")
	if len(tbl.Columns) != 2 {
		t.Fatalf("stand-in has %d columns, want 2", len(tbl.Columns))
	}
	for _, col := range tbl.Columns {
		if !strings.EqualFold(col.DataType, "text") {
			t.Errorf("stand-in column %s is %s, want text", col.Name, col.DataType)
		}
	}
	requireTrigger(t, c, "ins_stub")
	if trg := requireTrigger(t, c, "trg_needs_defaults"); trg.Timing == "" || trg.Event == "" || !strings.EqualFold(trg.Table, "broken") {
		t.Errorf("trigger stand-in = %s %s ON %s", trg.Timing, trg.Event, trg.Table)
	}
	if cols := requireView(t, c, "bad_view").Columns; !slices.Equal(cols, []string{"x"}) {
		t.Errorf("bad_view stand-in columns = %v, want x", cols)
	}
	requireEvent(t, c, "e_empty")
	requireEvent(t, c, "e_broken")
}

func TestStandInStatements(t *testing.T) {
	c := New()
	c.LoadMetadata(context.Background(), snapshot(&metadata.SchemaMetadata{})) //nolint:errcheck // creates the database
	for _, stmt := range []struct {
		name string
		err  error
	}{
		{"pt", c.DefineTable(standInTableStmt(&metadata.TableMetadata{Name: "pt", Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}, {Name: "a"}}}))},
		{"empty", c.DefineTable(standInTableStmt(&metadata.TableMetadata{Name: "empty"}))},
		{"pv", c.DefineView(standInViewStmt(&metadata.ViewMetadata{Name: "pv", Columns: []*metadata.ColumnMetadata{{Name: "x"}, {Name: "y"}}}))},
		{"dv", c.DefineView(standInViewStmt(&metadata.ViewMetadata{Name: "dv", DependencyColumns: []*metadata.DependencyColumn{{Column: "id"}}}))},
	} {
		if stmt.err != nil {
			t.Fatalf("%s: %v", stmt.name, stmt.err)
		}
	}
	if cols := requireTable(t, c, "pt").Columns; len(cols) != 2 || !strings.EqualFold(cols[0].DataType, "text") {
		t.Errorf("pt stand-in columns = %+v, want a and b as text", cols)
	}
	if cols := requireTable(t, c, "empty").Columns; len(cols) != 1 || cols[0].Name != standInPlaceholder {
		t.Errorf("empty stand-in columns = %+v, want the placeholder", cols)
	}
	if cols := requireView(t, c, "pv").Columns; len(cols) != 2 {
		t.Errorf("pv stand-in has %d columns, want 2", len(cols))
	}
	if cols := requireView(t, c, "dv").Columns; len(cols) != 1 || cols[0] != "id" {
		t.Errorf("dv stand-in columns = %v, want the column it reads", cols)
	}
}

func TestLoadMetadataStopsWhenContextIsDone(t *testing.T) {
	for _, tt := range []struct {
		name   string
		checks int
		meta   *metadata.DatabaseSchemaMetadata
	}{
		{name: "before planning", checks: 0, meta: snapshot(&metadata.SchemaMetadata{})},
		{name: "before an install", checks: 1, meta: snapshotWithTables(&metadata.TableMetadata{Name: "t"})},
		{name: "before a view reinstall", checks: 3, meta: snapshot(&metadata.SchemaMetadata{Views: []*metadata.ViewMetadata{
			{Name: "a", Definition: "SELECT id FROM z"},
			{Name: "z", Definition: "SELECT 1 AS id"},
		}})},
	} {
		if _, err := New().LoadMetadata(&doneAfter{Context: context.Background(), checks: tt.checks}, tt.meta); err != context.Canceled {
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

func TestParseTypeName(t *testing.T) {
	for typ, want := range map[string]string{
		"int":                          "INT",
		"int unsigned":                 "INT",
		"bigint(20) unsigned zerofill": "BIGINT",
		"varchar(255)":                 "VARCHAR",
		"char(1)":                      "CHAR",
		"text":                         "TEXT",
		"decimal(10,2)":                "DECIMAL",
		"datetime(6)":                  "DATETIME",
		"timestamp":                    "TIMESTAMP",
		"tinyint(1)":                   "TINYINT",
		"blob":                         "BLOB",
		"json":                         "JSON",
		"enum('a','b','c')":            "ENUM",
		"set('x','y')":                 "SET",
		"geometry":                     "GEOMETRY",
		"bit(1)":                       "BIT",
	} {
		dt, err := parseTypeName(typ)
		if err != nil || !strings.EqualFold(dt.Name, want) {
			t.Errorf("parseTypeName(%q) = %v, %v; want %s", typ, dt, err, want)
		}
	}
	for _, typ := range []string{"", "not a real type (((", "int UNSIGNED ZEROFILL NOT_A_THING"} {
		if _, err := parseTypeName(typ); err == nil {
			t.Errorf("parseTypeName(%q) succeeded, want an error", typ)
		}
	}
}

func TestParseExpr(t *testing.T) {
	for _, expr := range []string{`'hello'`, `CURRENT_TIMESTAMP`, `NULL`, `0`, `a + b`, `CONCAT('x', 'y')`, `DATE_FORMAT(NOW(), '%Y-%m-%d')`} {
		if _, err := parseExpr(expr); err != nil {
			t.Errorf("parseExpr(%q): %v", expr, err)
		}
	}
	for _, expr := range []string{"", "(((", "this is not sql at all ###"} {
		if _, err := parseExpr(expr); err == nil {
			t.Errorf("parseExpr(%q) succeeded, want an error", expr)
		}
	}
}
