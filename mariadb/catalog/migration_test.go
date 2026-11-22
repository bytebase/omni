package catalog

import (
	"context"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func planFor(t *testing.T, from, to string) *MigrationPlan {
	t.Helper()
	a, b := loadCat(t, dbDDL+from), loadCat(t, dbDDL+to)
	return GenerateMigration(a, b, Diff(a, b))
}

func TestMigration_EmptyDiffEmptyPlan(t *testing.T) {
	if sql := planFor(t, "CREATE TABLE t (a INT);", "CREATE TABLE t (a INT(11));").SQL(); sql != "" {
		t.Errorf("plan = %q, want empty", sql)
	}
	if plan := GenerateMigration(New(), New(), nil); len(plan.Ops) != 0 {
		t.Errorf("nil diff: %d ops", len(plan.Ops))
	}
}

func TestMigration_Deterministic(t *testing.T) {
	from := "CREATE TABLE a (id INT); CREATE TABLE b (id INT); CREATE SEQUENCE s;"
	to := "CREATE TABLE a (id INT, x INT, KEY k (x)); CREATE TABLE c (id INT); CREATE SEQUENCE s INCREMENT BY 2; CREATE SEQUENCE u;"
	first := planFor(t, from, to).SQL()
	for i := 0; i < 20; i++ {
		if got := planFor(t, from, to).SQL(); got != first {
			t.Fatalf("plan changed between runs:\n%s\n---\n%s", first, got)
		}
	}
}

func TestMigration_SequenceLimitsLeavingOutTheStartWarn(t *testing.T) {
	plan := planFor(t, "CREATE SEQUENCE s;", "CREATE SEQUENCE s START WITH 100 MINVALUE 10;")
	if strings.Contains(plan.SQL(), "RESTART") {
		t.Errorf("the plan restarts the sequence:\n%s", plan.SQL())
	}
	if !plan.HasWarnings() {
		t.Errorf("no warning that the next value may lie outside the new limits:\n%s", plan.SQL())
	}
}

func TestMigration_ViewBodyQualifiesItsRelations(t *testing.T) {
	got := planFor(t, "", "CREATE TABLE t (a INT); CREATE VIEW v AS WITH c AS (SELECT a FROM t) SELECT c.a FROM c JOIN t ON c.a = t.a;").SQL()
	if strings.Count(got, "`app`.`t`") < 2 || strings.Contains(got, "`app`.`c`") {
		t.Errorf("want t qualified with app and the CTE c bare:\n%s", got)
	}
}

func TestMigration_AddedObjectsKeepTheirSessionContext(t *testing.T) {
	a := loadCat(t, dbDDL+"CREATE TABLE t (a VARCHAR(10));")
	b := loadCat(t, dbDDL+"CREATE TABLE t (a VARCHAR(10)); CREATE FUNCTION f() RETURNS VARCHAR(10) DETERMINISTIC RETURN 'a' || 'b'; CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 'x' || 'y'; CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 'a' || 'b';")
	ctx := SessionContext{SQLMode: "PIPES_AS_CONCAT"}
	b.ApplySessionContext(SessionContextMap{
		Functions: map[string]SessionContext{"f": ctx},
		Triggers:  map[string]SessionContext{"tr": ctx},
		Events:    map[string]SessionContext{"e": ctx},
	})
	for _, op := range GenerateMigration(a, b, Diff(a, b)).Ops {
		switch op.Type {
		case OpCreateFunction, OpCreateTrigger, OpCreateEvent:
			if !strings.Contains(op.SQL, "SET sql_mode = 'PIPES_AS_CONCAT'") {
				t.Errorf("%s created without its sql_mode:\n%s", op.ObjectName, op.SQL)
			}
		}
	}
}

func TestMigration_MariaDBStatements(t *testing.T) {
	for _, tc := range []struct{ name, from, to, want string }{
		{"sequence create", "", "CREATE SEQUENCE s START WITH 5;", "CREATE SEQUENCE `app`.`s` start with 5 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB"},
		{"sequence drop", "CREATE SEQUENCE s;", "", "DROP SEQUENCE `app`.`s`"},
		{"sequence alter", "CREATE SEQUENCE s;", "CREATE SEQUENCE s INCREMENT BY 2 COMMENT 'c';", "ALTER SEQUENCE `app`.`s` INCREMENT BY 2;\nALTER TABLE `app`.`s` COMMENT='c'"},
		{"check drop", "CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 0));", "CREATE TABLE t (a INT);", "ALTER TABLE `app`.`t` DROP CONSTRAINT `c`"},
		{"column check", "CREATE TABLE t (a INT);", "CREATE TABLE t (a INT CHECK (a > 0));", "MODIFY COLUMN `a` int(11) DEFAULT NULL CHECK (`a` > 0)"},
		{"versioning", "CREATE TABLE t (a INT);", "CREATE TABLE t (a INT) WITH SYSTEM VERSIONING;", "ALTER TABLE `app`.`t` ADD SYSTEM VERSIONING"},
		{"keep history", "CREATE TABLE t (a INT) WITH SYSTEM VERSIONING;", "CREATE TABLE t (a INT, b INT) WITH SYSTEM VERSIONING;", "SET @@system_versioning_alter_history = KEEP"},
		{"period", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL);", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR p (s, e));", "ALTER TABLE `app`.`t` ADD PERIOD FOR `p` (`s`, `e`)"},
		{"option reset", "CREATE TABLE t (a INT) STATS_PERSISTENT=1 COMMENT 'x';", "CREATE TABLE t (a INT);", "ALTER TABLE `app`.`t` STATS_PERSISTENT=DEFAULT COMMENT=''"},
		{"event qualified", "", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1;", "CREATE EVENT `app`.`e`"},
		{"ignored key", "CREATE TABLE t (a INT, KEY k (a));", "CREATE TABLE t (a INT, KEY k (a) IGNORED);", "ADD KEY `k` (`a`) IGNORED"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := planFor(t, tc.from, tc.to).SQL(); !strings.Contains(got, tc.want) {
				t.Errorf("plan:\n%s\nwant it to contain:\n%s", got, tc.want)
			}
		})
	}
}

func TestMigration_RecreateInheritsSourceCommentEscaping(t *testing.T) {
	// The target has no captured context; the recreate inherits the source's
	// NO_BACKSLASH_ESCAPES, so the comment's backslash must render literally
	// (single backslash), matching the mode the statement is framed with.
	from := loadCat(t, dbDDL+`CREATE FUNCTION f() RETURNS INT DETERMINISTIC COMMENT 'a\\b' RETURN 1;`)
	from.ApplySessionContext(SessionContextMap{Functions: map[string]SessionContext{"f": {SQLMode: "NO_BACKSLASH_ESCAPES"}}})
	to := loadCat(t, dbDDL+`CREATE FUNCTION f() RETURNS INT DETERMINISTIC COMMENT 'c\\d' RETURN 2;`)
	var create string
	for _, op := range GenerateMigration(from, to, Diff(from, to)).Ops {
		if op.Type == OpCreateFunction {
			create = op.SQL
		}
	}
	if !strings.Contains(create, "SET sql_mode = 'NO_BACKSLASH_ESCAPES'") {
		t.Fatalf("recreate not framed with the source mode:\n%s", create)
	}
	if !strings.Contains(create, `COMMENT 'c\d'`) || strings.Contains(create, `COMMENT 'c\\d'`) {
		t.Errorf("comment not escaped for NO_BACKSLASH_ESCAPES:\n%s", create)
	}
}

func TestMigration_EventInheritsSourceCommentEscaping(t *testing.T) {
	from := loadCat(t, dbDDL+`CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' COMMENT 'a\\b' DO SELECT 1;`)
	from.ApplySessionContext(SessionContextMap{Events: map[string]SessionContext{"e": {SQLMode: "NO_BACKSLASH_ESCAPES"}}})
	to := loadCat(t, dbDDL+`CREATE EVENT e ON SCHEDULE EVERY 2 DAY STARTS '2030-01-01 00:00:00' COMMENT 'c\\d' DO SELECT 1;`)
	var stmt string
	for _, op := range GenerateMigration(from, to, Diff(from, to)).Ops {
		if op.Type == OpCreateEvent {
			stmt = op.SQL
		}
	}
	if !strings.Contains(stmt, "SET sql_mode = 'NO_BACKSLASH_ESCAPES'") {
		t.Fatalf("event statement not framed with the source mode:\n%s", stmt)
	}
	if !strings.Contains(stmt, `COMMENT 'c\d'`) || strings.Contains(stmt, `COMMENT 'c\\d'`) {
		t.Errorf("event comment not escaped for NO_BACKSLASH_ESCAPES:\n%s", stmt)
	}
}

func TestMigration_MetadataSystemVersioning(t *testing.T) {
	implicit := &metadata.TableMetadata{
		Name:          "history",
		Engine:        "InnoDB",
		CreateOptions: "WITH SYSTEM VERSIONING",
		Columns:       []*metadata.ColumnMetadata{{Name: "id", Type: "int", Nullable: false}},
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Expressions: []string{"id"}, Primary: true, Unique: true, Visible: true, Type: "BTREE"},
		},
	}
	explicit := &metadata.TableMetadata{
		Name:          "explicit_period",
		Engine:        "InnoDB",
		CreateOptions: "WITH SYSTEM VERSIONING",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id", Type: "int", Nullable: true},
			{Name: "valid_from", Type: "timestamp(6)", Generation: &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "ROW START"}},
			{Name: "valid_to", Type: "timestamp(6)", Generation: &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "ROW END"}},
		},
	}
	load := func(tables ...*metadata.TableMetadata) *Catalog {
		t.Helper()
		c := New()
		meta := &metadata.DatabaseSchemaMetadata{
			Name:    "db",
			Schemas: []*metadata.SchemaMetadata{{Name: "", Tables: tables}},
		}
		report, err := c.LoadMetadataInto(context.Background(), meta, "bbcatalog")
		if err != nil {
			t.Fatalf("LoadMetadataInto: %v", err)
		}
		requireReport(t, report, nil, nil)
		return c
	}

	empty := load()
	withSV := load(implicit, explicit)
	got := GenerateMigration(empty, withSV, Diff(empty, withSV)).SQL()
	for _, want := range []string{
		"CREATE TABLE `bbcatalog`.`explicit_period`",
		"`valid_from` timestamp(6) GENERATED ALWAYS AS ROW START",
		"`valid_to` timestamp(6) GENERATED ALWAYS AS ROW END",
		"PERIOD FOR SYSTEM_TIME (`valid_from`, `valid_to`)",
		") ENGINE=InnoDB",
		"WITH SYSTEM VERSIONING",
		"CREATE TABLE `bbcatalog`.`history`",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("migration missing %q:\n%s", want, got)
		}
	}

	if selfDiff := GenerateMigration(withSV, load(implicit, explicit), Diff(withSV, load(implicit, explicit))).SQL(); selfDiff != "" {
		t.Errorf("self-diff of system-versioned tables not empty:\n%s", selfDiff)
	}
}

func TestMigration_MetadataColumnAttributes(t *testing.T) {
	srid := func(v uint32) *uint32 { return &v }
	id := &metadata.ColumnMetadata{Name: "id", Type: "int", Nullable: false}
	diffCols := func(oldCol, newCol *metadata.ColumnMetadata) string {
		t.Helper()
		load := func(col *metadata.ColumnMetadata) *Catalog {
			c := New()
			meta := &metadata.DatabaseSchemaMetadata{
				Name: "d",
				Schemas: []*metadata.SchemaMetadata{{
					Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{id, col}}},
				}},
			}
			report, err := c.LoadMetadataInto(context.Background(), meta, "bbcatalog")
			if err != nil {
				t.Fatalf("LoadMetadataInto: %v", err)
			}
			requireReport(t, report, nil, nil)
			return c
		}
		from, to := load(oldCol), load(newCol)
		return GenerateMigration(from, to, Diff(from, to)).SQL()
	}

	// MariaDB has no SRID column attribute, so an SRID-only change is a no-op.
	if got := diffCols(
		&metadata.ColumnMetadata{Name: "pt", Type: "point", Nullable: false},
		&metadata.ColumnMetadata{Name: "pt", Type: "point", Nullable: false, Srid: srid(4326)},
	); got != "" {
		t.Errorf("SRID-only change should be empty in MariaDB, got:\n%s", got)
	}

	// INVISIBLE uses plain MariaDB syntax rather than MySQL's version-gated comment.
	got := diffCols(
		&metadata.ColumnMetadata{Name: "c", Type: "int", Nullable: true, Default: "NULL"},
		&metadata.ColumnMetadata{Name: "c", Type: "int", Nullable: true, Default: "NULL", IsInvisible: true},
	)
	if want := "ALTER TABLE `bbcatalog`.`t` MODIFY COLUMN `c` int(11) DEFAULT NULL INVISIBLE"; got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestMigration_MetadataAcrossDatabaseNames(t *testing.T) {
	id := &metadata.ColumnMetadata{Name: "id", Type: "int"}
	view := func(db string) *metadata.ViewMetadata {
		return &metadata.ViewMetadata{Name: "v", Definition: "select `" + db + "`.`t`.`id` AS `id` from `" + db + "`.`t`"}
	}
	load := func(db string, views ...*metadata.ViewMetadata) *Catalog {
		t.Helper()
		c := New()
		meta := &metadata.DatabaseSchemaMetadata{
			Name: db,
			Schemas: []*metadata.SchemaMetadata{{
				Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{id}}},
				Views:  views,
			}},
		}
		report, err := c.LoadMetadataInto(context.Background(), meta, "bbcatalog")
		if err != nil {
			t.Fatalf("LoadMetadataInto: %v", err)
		}
		requireReport(t, report, nil, nil)
		return c
	}

	prodWithView := load("prod", view("prod"))
	stagingWithView := load("staging", view("staging"))
	if got := GenerateMigration(prodWithView, stagingWithView, Diff(prodWithView, stagingWithView)).SQL(); got != "" {
		t.Errorf("identical views across database names should diff empty, got:\n%s", got)
	}

	prodNoView := load("prod")
	got := GenerateMigration(prodNoView, stagingWithView, Diff(prodNoView, stagingWithView)).SQL()
	if strings.Contains(got, "staging") {
		t.Errorf("new view leaked source database name:\n%s", got)
	}
	if !strings.Contains(got, "VIEW `bbcatalog`.`v` AS select `t`.`id` AS `id` from `bbcatalog`.`t`") {
		t.Errorf("unexpected view migration SQL:\n%s", got)
	}
}

func TestMigration_JSONColumnsAndFunctionalIndexes(t *testing.T) {
	from := `
CREATE TABLE users (
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(100) NOT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB;`
	to := `
CREATE TABLE users (
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(100) NOT NULL,
	preferences JSON,
	PRIMARY KEY (id)
) ENGINE=InnoDB;
CREATE TABLE products (
	id INT NOT NULL AUTO_INCREMENT,
	details JSON NOT NULL,
	brand VARCHAR(50) AS (json_unquote(json_extract(details, '$.brand'))) VIRTUAL,
	PRIMARY KEY (id),
	INDEX idx_brand (brand),
	CONSTRAINT chk_details_not_empty CHECK (json_type(details) IS NOT NULL)
) ENGINE=InnoDB;`

	got := planFor(t, from, to).SQL()
	for _, want := range []string{
		"CREATE TABLE `app`.`products`",
		"`details` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`details`))",
		"ALTER TABLE `app`.`users` ADD COLUMN `preferences` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`preferences`))",
		"ALTER TABLE `app`.`products` ADD KEY `idx_brand` (`brand`)",
		"ALTER TABLE `app`.`products` ADD CONSTRAINT `chk_details_not_empty` CHECK (json_type(`details`) is not null)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("JSON migration missing %q:\n%s", want, got)
		}
	}
}
