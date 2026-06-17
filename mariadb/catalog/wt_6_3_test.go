package catalog

import (
	"strings"
	"testing"
)

func TestWalkThrough_6_3_FKActionsRendering(t *testing.T) {
	// Scenario 1: FK with ON DELETE CASCADE ON UPDATE SET NULL — both actions rendered in SHOW CREATE
	t.Run("fk_with_cascade_and_set_null", func(t *testing.T) {
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id)
				ON DELETE CASCADE ON UPDATE SET NULL
		)`)

		got := c.ShowCreateTable("test", "child")
		if !strings.Contains(got, "ON DELETE CASCADE") {
			t.Errorf("expected ON DELETE CASCADE in output:\n%s", got)
		}
		if !strings.Contains(got, "ON UPDATE SET NULL") {
			t.Errorf("expected ON UPDATE SET NULL in output:\n%s", got)
		}
	})

	// Scenario 2: FK with no action specified — defaults not rendered in SHOW CREATE
	t.Run("fk_with_no_action_defaults_omitted", func(t *testing.T) {
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			parent_id INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES parent (id)
		)`)

		got := c.ShowCreateTable("test", "child")
		// MySQL 8.0 does not show ON DELETE or ON UPDATE when using default (NO ACTION).
		if strings.Contains(got, "ON DELETE") {
			t.Errorf("default FK action should not show ON DELETE:\n%s", got)
		}
		if strings.Contains(got, "ON UPDATE") {
			t.Errorf("default FK action should not show ON UPDATE:\n%s", got)
		}
		// Verify FK is still rendered.
		if !strings.Contains(got, "CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`)") {
			t.Errorf("FK constraint not rendered correctly:\n%s", got)
		}
	})

	// Scenario 3: Multi-column FK with actions — actions on composite FK rendered correctly
	t.Run("multi_column_fk_with_actions", func(t *testing.T) {
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, `CREATE TABLE parent (
			a INT NOT NULL,
			b INT NOT NULL,
			PRIMARY KEY (a, b)
		)`)
		mustExec(t, c, `CREATE TABLE child (
			id INT NOT NULL,
			pa INT,
			pb INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_composite FOREIGN KEY (pa, pb) REFERENCES parent (a, b)
				ON DELETE CASCADE ON UPDATE SET NULL
		)`)

		got := c.ShowCreateTable("test", "child")
		// Verify multi-column FK is rendered with both columns.
		if !strings.Contains(got, "FOREIGN KEY (`pa`, `pb`) REFERENCES `parent` (`a`, `b`)") {
			t.Errorf("multi-column FK not rendered correctly:\n%s", got)
		}
		// Actions rendered once for the whole FK.
		if !strings.Contains(got, "ON DELETE CASCADE") {
			t.Errorf("expected ON DELETE CASCADE on composite FK:\n%s", got)
		}
		if !strings.Contains(got, "ON UPDATE SET NULL") {
			t.Errorf("expected ON UPDATE SET NULL on composite FK:\n%s", got)
		}
		// Verify only one occurrence of each action (not per-column).
		if strings.Count(got, "ON DELETE CASCADE") != 1 {
			t.Errorf("expected exactly one ON DELETE CASCADE, got %d:\n%s",
				strings.Count(got, "ON DELETE CASCADE"), got)
		}
		if strings.Count(got, "ON UPDATE SET NULL") != 1 {
			t.Errorf("expected exactly one ON UPDATE SET NULL, got %d:\n%s",
				strings.Count(got, "ON UPDATE SET NULL"), got)
		}
	})
}
