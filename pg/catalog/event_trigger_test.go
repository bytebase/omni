package catalog

import (
	"strings"
	"testing"
)

const eventTriggerFuncSQL = `
CREATE FUNCTION evt_fn() RETURNS event_trigger
LANGUAGE plpgsql AS $$ BEGIN END $$;
`

func TestEventTriggerCatalog(t *testing.T) {
	t.Run("create alter comment and drop", func(t *testing.T) {
		c, err := LoadSQL(eventTriggerFuncSQL + `
	CREATE EVENT TRIGGER evt_start ON ddl_command_start
		WHEN TAG IN ('create table', 'DROP TABLE')
	EXECUTE FUNCTION evt_fn();
ALTER EVENT TRIGGER evt_start ENABLE ALWAYS;
COMMENT ON EVENT TRIGGER evt_start IS 'event trigger comment';
`)
		if err != nil {
			t.Fatal(err)
		}

		evt := c.eventTriggerByName["evt_start"]
		if evt == nil {
			t.Fatal("expected event trigger evt_start to be tracked")
		}
		if evt.EventName != "ddl_command_start" {
			t.Fatalf("expected ddl_command_start, got %q", evt.EventName)
		}
		if evt.Enabled != 'A' {
			t.Fatalf("expected enabled state A, got %q", evt.Enabled)
		}
		if len(evt.Tags) != 2 || evt.Tags[0] != "CREATE TABLE" || evt.Tags[1] != "DROP TABLE" {
			t.Fatalf("expected canonical tags, got %#v", evt.Tags)
		}
		if got := c.comments[commentKey{ObjType: 'E', ObjOID: evt.OID}]; got != "event trigger comment" {
			t.Fatalf("expected event trigger comment, got %q", got)
		}

		if _, err := c.Exec("DROP EVENT TRIGGER evt_start", nil); err != nil {
			t.Fatal(err)
		}
		if c.eventTriggerByName["evt_start"] != nil {
			t.Fatal("expected event trigger to be removed after DROP")
		}
	})

	t.Run("accepts postgres event trigger command tags", func(t *testing.T) {
		c, err := LoadSQL(eventTriggerFuncSQL + `
	CREATE EVENT TRIGGER evt_start ON ddl_command_start
		WHEN TAG IN (
			'REFRESH MATERIALIZED VIEW',
			'REINDEX',
			'SECURITY LABEL',
			'REVOKE',
			'SELECT INTO',
			'DROP CONSTRAINT',
			'LOGIN'
		)
		EXECUTE FUNCTION evt_fn();
	`)
		if err != nil {
			t.Fatal(err)
		}
		evt := c.eventTriggerByName["evt_start"]
		if evt == nil {
			t.Fatal("expected event trigger evt_start to be tracked")
		}
		want := []string{
			"REFRESH MATERIALIZED VIEW",
			"REINDEX",
			"SECURITY LABEL",
			"REVOKE",
			"SELECT INTO",
			"DROP CONSTRAINT",
			"LOGIN",
		}
		if strings.Join(evt.Tags, ",") != strings.Join(want, ",") {
			t.Fatalf("expected postgres command tags %#v, got %#v", want, evt.Tags)
		}
	})

	t.Run("accepts postgres table rewrite command tags", func(t *testing.T) {
		c, err := LoadSQL(eventTriggerFuncSQL + `
	CREATE EVENT TRIGGER evt_rewrite ON table_rewrite
		WHEN TAG IN ('ALTER MATERIALIZED VIEW', 'ALTER TABLE', 'ALTER TYPE')
		EXECUTE FUNCTION evt_fn();
	`)
		if err != nil {
			t.Fatal(err)
		}
		evt := c.eventTriggerByName["evt_rewrite"]
		if evt == nil {
			t.Fatal("expected event trigger evt_rewrite to be tracked")
		}
		if len(evt.Tags) != 3 || evt.Tags[0] != "ALTER MATERIALIZED VIEW" || evt.Tags[1] != "ALTER TABLE" || evt.Tags[2] != "ALTER TYPE" {
			t.Fatalf("expected table rewrite tags, got %#v", evt.Tags)
		}
	})

	t.Run("rejects postgres incompatible definitions", func(t *testing.T) {
		tests := []struct {
			name    string
			sql     string
			wantErr string
		}{
			{
				name:    "invalid event name",
				sql:     eventTriggerFuncSQL + `CREATE EVENT TRIGGER bad ON elephant_bootstrap EXECUTE FUNCTION evt_fn();`,
				wantErr: `unrecognized event name "elephant_bootstrap"`,
			},
			{
				name: "trigger function return type",
				sql: `
CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END $$;
CREATE EVENT TRIGGER bad ON ddl_command_start EXECUTE FUNCTION trg_fn();
`,
				wantErr: `function trg_fn must return type event_trigger`,
			},
			{
				name:    "duplicate tag filter",
				sql:     eventTriggerFuncSQL + `CREATE EVENT TRIGGER bad ON ddl_command_start WHEN TAG IN ('CREATE TABLE') AND TAG IN ('DROP TABLE') EXECUTE FUNCTION evt_fn();`,
				wantErr: `filter variable "tag" specified more than once`,
			},
			{
				name:    "unknown filter variable",
				sql:     eventTriggerFuncSQL + `CREATE EVENT TRIGGER bad ON ddl_command_start WHEN foo IN ('CREATE TABLE') EXECUTE FUNCTION evt_fn();`,
				wantErr: `unrecognized filter variable "foo"`,
			},
			{
				name:    "unsupported ddl tag",
				sql:     eventTriggerFuncSQL + `CREATE EVENT TRIGGER bad ON ddl_command_start WHEN TAG IN ('DROP EVENT TRIGGER') EXECUTE FUNCTION evt_fn();`,
				wantErr: `event triggers are not supported for DROP EVENT TRIGGER`,
			},
			{
				name:    "unknown ddl tag",
				sql:     eventTriggerFuncSQL + `CREATE EVENT TRIGGER bad ON ddl_command_start WHEN TAG IN ('sandwich') EXECUTE FUNCTION evt_fn();`,
				wantErr: `filter value "SANDWICH" not recognized for filter variable "tag"`,
			},
			{
				name:    "login tag filter",
				sql:     eventTriggerFuncSQL + `CREATE EVENT TRIGGER bad ON login WHEN TAG IN ('CREATE TABLE') EXECUTE FUNCTION evt_fn();`,
				wantErr: `tag filtering is not supported for login event triggers`,
			},
			{
				name:    "table rewrite rejects create table",
				sql:     eventTriggerFuncSQL + `CREATE EVENT TRIGGER bad ON table_rewrite WHEN TAG IN ('CREATE TABLE') EXECUTE FUNCTION evt_fn();`,
				wantErr: `event triggers are not supported for CREATE TABLE`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := LoadSQL(tt.sql)
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
			})
		}
	})
}

func TestEventTriggerSDL(t *testing.T) {
	c, err := LoadSDL(`
COMMENT ON EVENT TRIGGER evt_start IS 'event trigger comment';
ALTER EVENT TRIGGER evt_start DISABLE;
CREATE EVENT TRIGGER evt_start ON ddl_command_start
	WHEN TAG IN ('CREATE TABLE')
	EXECUTE FUNCTION evt_fn();
CREATE FUNCTION evt_fn() RETURNS event_trigger
LANGUAGE plpgsql AS $$ BEGIN END $$;
`)
	if err != nil {
		t.Fatal(err)
	}
	evt := c.eventTriggerByName["evt_start"]
	if evt == nil {
		t.Fatal("expected event trigger from SDL")
	}
	if evt.Enabled != 'D' {
		t.Fatalf("expected disabled event trigger, got %q", evt.Enabled)
	}
	if got := c.comments[commentKey{ObjType: 'E', ObjOID: evt.OID}]; got != "event trigger comment" {
		t.Fatalf("expected event trigger comment, got %q", got)
	}
}

func TestEventTriggerDiffAndMigration(t *testing.T) {
	before := eventTriggerFuncSQL + `
CREATE EVENT TRIGGER evt_start ON ddl_command_start EXECUTE FUNCTION evt_fn();
`
	after := eventTriggerFuncSQL + `
CREATE EVENT TRIGGER evt_start ON ddl_command_start
	WHEN TAG IN ('CREATE TABLE')
	EXECUTE FUNCTION evt_fn();
ALTER EVENT TRIGGER evt_start ENABLE REPLICA;
COMMENT ON EVENT TRIGGER evt_start IS 'event trigger comment';
`

	from, err := LoadSDL(before)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSDL(after)
	if err != nil {
		t.Fatal(err)
	}
	diff := Diff(from, to)
	if len(diff.EventTriggers) != 1 {
		t.Fatalf("expected 1 event trigger diff, got %d", len(diff.EventTriggers))
	}
	if diff.EventTriggers[0].Action != DiffModify {
		t.Fatalf("expected event trigger modify diff, got %d", diff.EventTriggers[0].Action)
	}
	if len(diff.Comments) != 1 {
		t.Fatalf("expected 1 comment diff, got %d", len(diff.Comments))
	}

	plan := GenerateMigration(from, to, diff)
	sql := plan.SQL()
	if !strings.Contains(sql, "DROP EVENT TRIGGER") || !strings.Contains(sql, "CREATE EVENT TRIGGER") {
		t.Fatalf("expected drop/create event trigger migration, got:\n%s", sql)
	}
	if !strings.Contains(sql, "ALTER EVENT TRIGGER") || !strings.Contains(sql, "ENABLE REPLICA") {
		t.Fatalf("expected enable replica migration, got:\n%s", sql)
	}
	if !strings.Contains(sql, "COMMENT ON EVENT TRIGGER") {
		t.Fatalf("expected event trigger comment migration, got:\n%s", sql)
	}

	migrated, err := LoadSQL(before + ";\n" + sql)
	if err != nil {
		t.Fatalf("migration failed: %v\nSQL:\n%s", err, sql)
	}
	if diff2 := Diff(migrated, to); !diff2.IsEmpty() {
		t.Fatalf("expected empty roundtrip diff, got %#v", diff2)
	}
}

func TestEventTriggerFunctionDependency(t *testing.T) {
	c, err := LoadSQL(eventTriggerFuncSQL + `
CREATE EVENT TRIGGER evt_start ON ddl_command_start EXECUTE FUNCTION evt_fn();
`)
	if err != nil {
		t.Fatal(err)
	}
	results, err := c.Exec("DROP FUNCTION evt_fn()", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Error == nil {
		t.Fatalf("expected DROP FUNCTION to be blocked by event trigger dependency, got %#v", results)
	}
	results, err = c.Exec("DROP FUNCTION evt_fn() CASCADE", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Error != nil {
		t.Fatalf("expected DROP FUNCTION CASCADE to succeed, got %#v", results)
	}
	if len(c.eventTriggers) != 0 {
		t.Fatalf("expected cascade to drop event trigger, got %d", len(c.eventTriggers))
	}
}

func TestEventTriggerFunctionDefinitionValidation(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
	}{
		{
			name:    "sql functions cannot return event_trigger",
			sql:     `CREATE FUNCTION evt_sql_fn() RETURNS event_trigger LANGUAGE sql AS $$ SELECT 1 $$;`,
			wantErr: `SQL functions cannot return type event_trigger`,
		},
		{
			name:    "event trigger functions cannot have arguments",
			sql:     `CREATE FUNCTION evt_arg_fn(name text) RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN END $$;`,
			wantErr: `event trigger functions cannot have declared arguments`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadSQL(tt.sql)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}
