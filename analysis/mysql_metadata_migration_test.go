package analysis

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// fakeMigrationCatalog records what a MySQLMetadataMigration loaded.
type fakeMigrationCatalog struct {
	meta     *metadata.DatabaseSchemaMetadata
	database string
}

func fakeMySQLMetadataMigration(loadErr error, problems map[string]error, sql string) MySQLMetadataMigration[*fakeMigrationCatalog] {
	return MySQLMetadataMigration[*fakeMigrationCatalog]{
		Load: func(_ context.Context, meta *metadata.DatabaseSchemaMetadata, database string) (*fakeMigrationCatalog, []map[string]error, error) {
			if loadErr != nil {
				return nil, nil, loadErr
			}
			return &fakeMigrationCatalog{meta: meta, database: database}, []map[string]error{problems}, nil
		},
		Migration: func(from, to *fakeMigrationCatalog) string {
			if from.database != "bbcatalog" || to.database != "bbcatalog" || from.meta == nil || to.meta == nil {
				return "unexpected load"
			}
			return sql
		},
	}
}

func oneTable(name string, mutate ...func(*metadata.TableMetadata)) *Database {
	table := &metadata.TableMetadata{Name: name, Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}
	for _, m := range mutate {
		m(table)
	}
	return &Database{Name: "d", Metadata: &metadata.DatabaseSchemaMetadata{Name: "d", Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{table}}}}}
}

func TestMySQLMetadataMigrationWithAFakeCatalog(t *testing.T) {
	ctx := context.Background()
	run := func(m MySQLMetadataMigration[*fakeMigrationCatalog], from, to *Database) (string, string) {
		sql, err := m.MetadataMigration(ctx, MetadataMigrationRequest{From: from, To: to})
		if err != nil {
			return sql, err.Error()
		}
		return sql, ""
	}
	plain := fakeMySQLMetadataMigration(nil, nil, "ALTER TABLE `bbcatalog`.`t` ADD COLUMN `b` int")

	if sql, err := run(plain, nil, &Database{}); sql != "ALTER TABLE `t` ADD COLUMN `b` int" || err != "" {
		t.Fatalf("nil and empty sides = %q, %q", sql, err)
	}
	if _, err := run(fakeMySQLMetadataMigration(errors.New("boom"), nil, ""), nil, nil); err != "failed to load source schema: boom" {
		t.Fatalf("load error = %q", err)
	}

	caseSensitive := func(d *Database) *Database { d.ObjectCaseSensitive = true; return d }
	for _, tc := range []struct {
		name     string
		from, to *Database
		want     string
	}{
		{"rename", caseSensitive(oneTable("Foo")), oneTable("foo"), `case-only rename "Foo" to "foo" is not supported by the catalog diff`},
		{"collision", oneTable("Foo"), caseSensitive(&Database{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{Name: "Foo"}},
			Views:  []*metadata.ViewMetadata{{Name: "foo"}},
		}}}}), `case-sensitive schema has "Foo" and "foo", which the catalog cannot distinguish`},
		{"insensitive", oneTable("Foo"), oneTable("foo"), ""},
	} {
		if _, err := run(plain, tc.from, tc.to); err != tc.want {
			t.Errorf("%s: error = %q, want %q", tc.name, err, tc.want)
		}
	}

	broken := fakeMySQLMetadataMigration(nil, map[string]error{"table:t": errors.New("bad default")}, "")
	if _, err := run(broken, oneTable("t"), oneTable("t", func(t *metadata.TableMetadata) { t.Comment = "c" })); err != "cannot generate a migration for objects whose definitions could not be loaded: table:t: bad default" {
		t.Fatalf("changed unloadable table = %q", err)
	}
	statistics := func(t *metadata.TableMetadata) { t.RowCount, t.DataSize, t.IndexSize, t.DataFree = 1, 2, 3, 4 }
	if _, err := run(broken, oneTable("t"), oneTable("t", statistics)); err != "" {
		t.Fatalf("statistics of an unloadable table = %q", err)
	}
	sequence := func(lastValue string) *Database {
		return &Database{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
			Sequences: []*metadata.SequenceMetadata{{Name: "s", Start: "1", LastValue: lastValue}},
		}}}}
	}
	brokenSequence := fakeMySQLMetadataMigration(nil, map[string]error{"sequence:s": errors.New("bad")}, "")
	if _, err := run(brokenSequence, sequence("1"), sequence("42")); err != "" {
		t.Fatalf("last value of an unloadable sequence = %q", err)
	}
}

func TestMySQLMetadataMigrationKeepsAnAliasNamedLikeTheDatabase(t *testing.T) {
	var loaded []string
	m := MySQLMetadataMigration[string]{
		Load: func(_ context.Context, _ *metadata.DatabaseSchemaMetadata, database string) (string, []map[string]error, error) {
			loaded = append(loaded, database)
			return database, nil, nil
		},
		Migration: func(_, to string) string {
			return "CREATE VIEW `" + to + "`.`v` AS select `bbcatalog`.`id` AS `id` from `" + to + "`.`t` `bbcatalog`"
		},
	}
	snapshot := func(definition string) *Database {
		return &Database{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
			Views: []*metadata.ViewMetadata{{Name: "v", Definition: definition}},
		}}}}
	}
	for _, tc := range []struct{ fromDefinition, database string }{
		{"select 1", "bbcatalog_1"},
		{"select 1 from `BBCATALOG_1`", "bbcatalog_2"},
	} {
		loaded = nil
		sql, err := m.MetadataMigration(context.Background(), MetadataMigrationRequest{
			From: snapshot(tc.fromDefinition),
			To:   snapshot("select `bbcatalog`.`id` AS `id` from `t` `bbcatalog`"),
		})
		if want := "CREATE VIEW `v` AS select `bbcatalog`.`id` AS `id` from `t` `bbcatalog`"; sql != want || err != nil {
			t.Errorf("migration = %q, %v, want %q", sql, err, want)
		}
		if want := []string{tc.database, tc.database}; !reflect.DeepEqual(loaded, want) {
			t.Errorf("loaded under %q, want %q", loaded, want)
		}
	}
}

func TestStripMySQLDatabaseQualifier(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"ALTER TABLE `t` ADD COLUMN `b` int", "ALTER TABLE `t` ADD COLUMN `b` int"},
		{"CREATE VIEW `bbcatalog`.`v` AS select `bbcatalog`.`t`.`a` from `bbcatalog`.`t`", "CREATE VIEW `v` AS select `t`.`a` from `t`"},
		{"select '`bbcatalog`.`t`', \"`bbcatalog`.`u`\" from `bbcatalog`.`t`", "select '`bbcatalog`.`t`', \"`bbcatalog`.`u`\" from `t`"},
		{"select 'it''s \\' `bbcatalog`.`t`' from `bbcatalog`.`t`", "select 'it''s \\' `bbcatalog`.`t`' from `t`"},
		{"CREATE TABLE `bbcatalog`.`bbcatalog` (`bbcatalog` int)", "CREATE TABLE `bbcatalog` (`bbcatalog` int)"},
		{"BEGIN\n-- don't `bbcatalog`.`a`\n# it's `bbcatalog`.`b`\n/* it's `bbcatalog`.`c` */ SELECT 1 FROM `bbcatalog`.`t`; END", "BEGIN\n-- don't `bbcatalog`.`a`\n# it's `bbcatalog`.`b`\n/* it's `bbcatalog`.`c` */ SELECT 1 FROM `t`; END"},
		{"SELECT 1 /*!50000 FROM `bbcatalog`.`t` */", "SELECT 1 /*!50000 FROM `t` */"},
		{"SELECT 1--`bbcatalog`.`t`", "SELECT 1--`t`"},
		{"SELECT `bbcatalog`.x, `bbcatalog`.`t", "SELECT `bbcatalog`.x, `t"},
	} {
		if got := StripMySQLDatabaseQualifier(tc.in, "bbcatalog"); got != tc.want {
			t.Errorf("StripMySQLDatabaseQualifier(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

type fakeScriptCatalog struct {
	scripts *[]string
	results []CatalogResult
	runErr  error
	current string
}

func (c fakeScriptCatalog) Run(script string) ([]CatalogResult, error) {
	*c.scripts = append(*c.scripts, script)
	if len(*c.scripts) == 1 {
		return nil, nil
	}
	return c.results, c.runErr
}

func (c fakeScriptCatalog) CurrentDatabase() string { return c.current }

func (c fakeScriptCatalog) DatabaseMetadata(database string) *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Name: database}
}

func TestMySQLMetadataExtractionWithAFakeCatalog(t *testing.T) {
	ctx := context.Background()
	var scripts []string
	extract := func(c fakeScriptCatalog, export func(MySQLScriptCatalog) *metadata.DatabaseSchemaMetadata, script string) (*metadata.DatabaseSchemaMetadata, string) {
		scripts = nil
		c.scripts = &scripts
		x := MySQLMetadataExtraction{Database: "tmp", NewCatalog: func() MySQLScriptCatalog { return c }, Export: export}
		meta, err := x.ExtractMetadata(ctx, script)
		if err != nil {
			return meta, err.Error()
		}
		return meta, ""
	}

	if meta, err := extract(fakeScriptCatalog{}, nil, " \n\t"); err != "" || meta == nil || meta.Name != "" || scripts != nil {
		t.Fatalf("empty script = %v, %q, ran %q", meta, err, scripts)
	}
	meta, err := extract(fakeScriptCatalog{results: []CatalogResult{{}}}, nil, "CREATE TABLE t (id int);")
	if err != "" || meta.GetName() != "tmp" {
		t.Fatalf("script = %v, %q", meta, err)
	}
	if want := []string{"SET foreign_key_checks = 0;\nCREATE DATABASE IF NOT EXISTS `tmp`;\nUSE `tmp`;", "CREATE TABLE t (id int);"}; !reflect.DeepEqual(scripts, want) {
		t.Fatalf("scripts = %q, want %q", scripts, want)
	}
	current := func(c MySQLScriptCatalog) *metadata.DatabaseSchemaMetadata {
		return c.DatabaseMetadata(c.CurrentDatabase())
	}
	if meta, _ := extract(fakeScriptCatalog{current: "app"}, current, "USE app;"); meta.GetName() != "app" {
		t.Fatalf("export = %v", meta)
	}
	if _, err := extract(fakeScriptCatalog{runErr: errors.New("syntax")}, nil, "x"); err != "failed to parse schema DDL: syntax" {
		t.Fatalf("parse error = %q", err)
	}
	failed := []CatalogResult{{}, {Index: 1, Error: errors.New("first")}, {Index: 2, Error: errors.New("second")}}
	if _, err := extract(fakeScriptCatalog{results: failed}, nil, "x"); err != "failed to execute schema DDL: first" {
		t.Fatalf("statement error = %q", err)
	}
}
