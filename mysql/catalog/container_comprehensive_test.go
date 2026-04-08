package catalog

import "testing"

func TestDDLWorkflow_Container(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	steps := []struct {
		name  string
		sql   string
		check string // table to SHOW CREATE TABLE after this step
	}{
		{"create_basic", "CREATE TABLE users (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(100) NOT NULL, email VARCHAR(255), PRIMARY KEY (id), UNIQUE KEY idx_email (email)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4", "users"},
		{"add_column", "ALTER TABLE users ADD COLUMN age INT DEFAULT 0", "users"},
		{"add_index", "CREATE INDEX idx_name ON users (name)", "users"},
		{"modify_column", "ALTER TABLE users MODIFY COLUMN name VARCHAR(200) NOT NULL", "users"},
		{"drop_index", "DROP INDEX idx_name ON users", "users"},
		{"create_orders", "CREATE TABLE orders (id INT NOT NULL AUTO_INCREMENT, user_id INT NOT NULL, amount DECIMAL(10,2), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id), KEY idx_user (user_id)) ENGINE=InnoDB", "orders"},
		{"rename_column", "ALTER TABLE users CHANGE COLUMN email email_address VARCHAR(255)", "users"},
		{"drop_column", "ALTER TABLE users DROP COLUMN age", "users"},
	}

	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")

	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			if err := ctr.execSQL(step.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			results, err := c.Exec(step.sql, nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			if results[0].Error != nil {
				t.Fatalf("omni exec error: %v", results[0].Error)
			}

			if step.check == "" {
				return
			}

			ctrDDL, err := ctr.showCreateTable(step.check)
			if err != nil {
				t.Fatalf("container show create: %v", err)
			}
			omniDDL := c.ShowCreateTable("test", step.check)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("SHOW CREATE TABLE mismatch:\n--- container ---\n%s\n--- omni ---\n%s",
					ctrDDL, omniDDL)
			}
		})
	}
}

func TestShowCreateTable_ContainerComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	ctr, cleanup := startContainer(t)
	defer cleanup()

	cases := []struct {
		name  string
		sql   string
		table string
	}{
		{"basic_types", "CREATE TABLE t_types (a INT, b VARCHAR(100), c TEXT, d DECIMAL(10,2), e DATETIME)", "t_types"},
		{"not_null_default", "CREATE TABLE t_defaults (id INT NOT NULL, name VARCHAR(50) DEFAULT 'test', active TINYINT(1) DEFAULT 1)", "t_defaults"},
		{"auto_increment_pk", "CREATE TABLE t_auto (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id))", "t_auto"},
		{"multi_col_pk", "CREATE TABLE t_multi_pk (a INT NOT NULL, b INT NOT NULL, c VARCHAR(10), PRIMARY KEY (a, b))", "t_multi_pk"},
		{"unique_index", "CREATE TABLE t_unique (id INT, email VARCHAR(255), UNIQUE KEY idx_email (email))", "t_unique"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctr.execSQL("DROP TABLE IF EXISTS " + tc.table)
			if err := ctr.execSQL(tc.sql); err != nil {
				t.Fatalf("container exec: %v", err)
			}
			ctrDDL, _ := ctr.showCreateTable(tc.table)

			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.sql, nil)
			omniDDL := c.ShowCreateTable("test", tc.table)

			if normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
				t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL, omniDDL)
			}
		})
	}
}
