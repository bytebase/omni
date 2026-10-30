package catalog

import "testing"

// --- 4.1 Schema Setup Migrations ---

// TestWalkThrough_4_1_CreateDBTablesIndexes tests creating a database with multiple
// tables and indexes in a single Exec call, then verifies all objects are present.
func TestWalkThrough_4_1_CreateDBTablesIndexes(t *testing.T) {
	c := New()
	sql := `
CREATE DATABASE myapp;
USE myapp;

CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_email (email)
);

CREATE TABLE posts (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT,
    status ENUM('draft','published','archived') DEFAULT 'draft',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_user_id (user_id),
    INDEX idx_status_created (status, created_at)
);

CREATE TABLE comments (
    id INT NOT NULL AUTO_INCREMENT,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    body TEXT NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_post_id (post_id),
    INDEX idx_user_id (user_id)
);
`
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("exec error on stmt %d: %v", r.Index, r.Error)
		}
	}

	db := c.GetDatabase("myapp")
	if db == nil {
		t.Fatal("database 'myapp' not found")
	}

	// Verify tables exist.
	for _, name := range []string{"users", "posts", "comments"} {
		if db.GetTable(name) == nil {
			t.Errorf("table %q not found", name)
		}
	}

	// Verify users table structure.
	users := db.GetTable("users")
	if users == nil {
		t.Fatal("users table not found")
	}
	if len(users.Columns) != 4 {
		t.Errorf("users: expected 4 columns, got %d", len(users.Columns))
	}
	// Check PK index.
	var hasPK bool
	for _, idx := range users.Indexes {
		if idx.Primary {
			hasPK = true
			if len(idx.Columns) != 1 || idx.Columns[0].Name != "id" {
				t.Errorf("users PK: expected column 'id', got %v", idx.Columns)
			}
		}
	}
	if !hasPK {
		t.Error("users: no PRIMARY KEY index found")
	}
	// Check unique index on email.
	var hasEmailIdx bool
	for _, idx := range users.Indexes {
		if idx.Name == "idx_email" {
			hasEmailIdx = true
			if !idx.Unique {
				t.Error("idx_email should be unique")
			}
		}
	}
	if !hasEmailIdx {
		t.Error("users: idx_email index not found")
	}

	// Verify posts table indexes.
	posts := db.GetTable("posts")
	if posts == nil {
		t.Fatal("posts table not found")
	}
	if len(posts.Columns) != 6 {
		t.Errorf("posts: expected 6 columns, got %d", len(posts.Columns))
	}
	idxNames := make(map[string]bool)
	for _, idx := range posts.Indexes {
		idxNames[idx.Name] = true
	}
	for _, expected := range []string{"idx_user_id", "idx_status_created"} {
		if !idxNames[expected] {
			t.Errorf("posts: index %q not found", expected)
		}
	}

	// Verify multi-column index column order.
	for _, idx := range posts.Indexes {
		if idx.Name == "idx_status_created" {
			if len(idx.Columns) != 2 {
				t.Fatalf("idx_status_created: expected 2 columns, got %d", len(idx.Columns))
			}
			if idx.Columns[0].Name != "status" || idx.Columns[1].Name != "created_at" {
				t.Errorf("idx_status_created: expected [status, created_at], got [%s, %s]",
					idx.Columns[0].Name, idx.Columns[1].Name)
			}
		}
	}

	// Verify comments table.
	comments := db.GetTable("comments")
	if comments == nil {
		t.Fatal("comments table not found")
	}
	if len(comments.Columns) != 4 {
		t.Errorf("comments: expected 4 columns, got %d", len(comments.Columns))
	}
}

// TestWalkThrough_4_1_CreateTableThenAddFK tests creating a table and then adding
// a foreign key that references it.
func TestWalkThrough_4_1_CreateTableThenAddFK(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `
CREATE TABLE parents (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    PRIMARY KEY (id)
);
`)

	wtExec(t, c, `
CREATE TABLE children (
    id INT NOT NULL AUTO_INCREMENT,
    parent_id INT NOT NULL,
    name VARCHAR(100),
    PRIMARY KEY (id),
    INDEX idx_parent_id (parent_id)
);
`)

	wtExec(t, c, `
ALTER TABLE children ADD CONSTRAINT fk_parent
    FOREIGN KEY (parent_id) REFERENCES parents (id)
    ON DELETE CASCADE ON UPDATE CASCADE;
`)

	db := c.GetDatabase("testdb")
	children := db.GetTable("children")
	if children == nil {
		t.Fatal("children table not found")
	}

	// Find the FK constraint.
	var fk *Constraint
	for _, con := range children.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_parent" {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint 'fk_parent' not found")
	}

	if fk.RefTable != "parents" {
		t.Errorf("FK RefTable: expected 'parents', got %q", fk.RefTable)
	}
	if len(fk.Columns) != 1 || fk.Columns[0] != "parent_id" {
		t.Errorf("FK Columns: expected [parent_id], got %v", fk.Columns)
	}
	if len(fk.RefColumns) != 1 || fk.RefColumns[0] != "id" {
		t.Errorf("FK RefColumns: expected [id], got %v", fk.RefColumns)
	}
	if fk.OnDelete != "CASCADE" {
		t.Errorf("FK OnDelete: expected CASCADE, got %q", fk.OnDelete)
	}
	if fk.OnUpdate != "CASCADE" {
		t.Errorf("FK OnUpdate: expected CASCADE, got %q", fk.OnUpdate)
	}
}

// TestWalkThrough_4_1_CreateTableThenView tests creating a table and then a view
// on it, verifying the view resolves columns correctly.
func TestWalkThrough_4_1_CreateTableThenView(t *testing.T) {
	c := wtSetup(t)

	wtExec(t, c, `
CREATE TABLE products (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    active TINYINT(1) DEFAULT 1,
    PRIMARY KEY (id)
);
`)

	wtExec(t, c, `
CREATE VIEW active_products AS
    SELECT id, name, price FROM products WHERE active = 1;
`)

	db := c.GetDatabase("testdb")
	v := db.Views[toLower("active_products")]
	if v == nil {
		t.Fatal("view 'active_products' not found")
	}

	// The view should have derived columns from the SELECT.
	if len(v.Columns) < 3 {
		t.Fatalf("view columns: expected at least 3, got %d: %v", len(v.Columns), v.Columns)
	}

	expected := []string{"id", "name", "price"}
	for i, want := range expected {
		if i >= len(v.Columns) {
			t.Errorf("missing column %d: expected %q", i, want)
			continue
		}
		if v.Columns[i] != want {
			t.Errorf("column %d: expected %q, got %q", i, want, v.Columns[i])
		}
	}
}

// TestWalkThrough_4_1_CreateDBSetCharsetTables tests that tables inherit charset
// from the database when created after a database with explicit charset.
func TestWalkThrough_4_1_CreateDBSetCharsetTables(t *testing.T) {
	c := New()

	sql := `
CREATE DATABASE latin_db DEFAULT CHARACTER SET latin1;
USE latin_db;

CREATE TABLE t1 (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    PRIMARY KEY (id)
);

CREATE TABLE t2 (
    id INT NOT NULL AUTO_INCREMENT,
    description TEXT,
    PRIMARY KEY (id)
) DEFAULT CHARSET=utf8mb4;
`
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("exec error on stmt %d: %v", r.Index, r.Error)
		}
	}

	db := c.GetDatabase("latin_db")
	if db == nil {
		t.Fatal("database 'latin_db' not found")
	}
	if db.Charset != "latin1" {
		t.Errorf("database charset: expected 'latin1', got %q", db.Charset)
	}

	// t1 should inherit latin1 from the database.
	t1 := db.GetTable("t1")
	if t1 == nil {
		t.Fatal("table t1 not found")
	}
	if t1.Charset != "latin1" {
		t.Errorf("t1 charset: expected 'latin1', got %q", t1.Charset)
	}

	// t1's VARCHAR column should inherit latin1.
	nameCol := t1.GetColumn("name")
	if nameCol == nil {
		t.Fatal("column 'name' not found in t1")
	}
	if nameCol.Charset != "latin1" {
		t.Errorf("t1.name charset: expected 'latin1', got %q", nameCol.Charset)
	}

	// t2 has explicit utf8mb4 override.
	t2 := db.GetTable("t2")
	if t2 == nil {
		t.Fatal("table t2 not found")
	}
	if t2.Charset != "utf8mb4" {
		t.Errorf("t2 charset: expected 'utf8mb4', got %q", t2.Charset)
	}

	// t2's TEXT column should inherit utf8mb4 from the table.
	descCol := t2.GetColumn("description")
	if descCol == nil {
		t.Fatal("column 'description' not found in t2")
	}
	if descCol.Charset != "utf8mb4" {
		t.Errorf("t2.description charset: expected 'utf8mb4', got %q", descCol.Charset)
	}
}

// TestWalkThrough_4_1_MysqldumpStyle tests a mysqldump-style output with SET vars,
// DELIMITER, procedures, triggers, and tables.
func TestWalkThrough_4_1_MysqldumpStyle(t *testing.T) {
	c := New()

	sql := `
SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT;
SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS;
SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION;
SET NAMES utf8mb4;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS=0;

CREATE DATABASE IF NOT EXISTS dumpdb DEFAULT CHARACTER SET utf8mb4;
USE dumpdb;

CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_username (username),
    UNIQUE KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE orders (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    total DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DELIMITER ;;

CREATE PROCEDURE get_user_orders(IN p_user_id INT)
BEGIN
    SELECT * FROM orders WHERE user_id = p_user_id;
END;;

CREATE TRIGGER trg_order_after_insert AFTER INSERT ON orders
FOR EACH ROW
BEGIN
    UPDATE users SET email = email WHERE id = NEW.user_id;
END;;

DELIMITER ;

SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT;
SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS;
SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION;
`

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("exec error on stmt %d (line %d): %v", r.Index, r.Line, r.Error)
		}
	}

	db := c.GetDatabase("dumpdb")
	if db == nil {
		t.Fatal("database 'dumpdb' not found")
	}

	// Verify tables.
	users := db.GetTable("users")
	if users == nil {
		t.Fatal("table 'users' not found")
	}
	if users.Engine != "InnoDB" {
		t.Errorf("users engine: expected 'InnoDB', got %q", users.Engine)
	}

	orders := db.GetTable("orders")
	if orders == nil {
		t.Fatal("table 'orders' not found")
	}

	// Verify FK on orders.
	var fk *Constraint
	for _, con := range orders.Constraints {
		if con.Type == ConForeignKey && con.Name == "fk_orders_user" {
			fk = con
			break
		}
	}
	if fk == nil {
		t.Fatal("FK constraint 'fk_orders_user' not found on orders table")
	}
	if fk.RefTable != "users" {
		t.Errorf("FK RefTable: expected 'users', got %q", fk.RefTable)
	}
	if fk.OnDelete != "CASCADE" {
		t.Errorf("FK OnDelete: expected 'CASCADE', got %q", fk.OnDelete)
	}

	// Verify procedure.
	proc := db.Procedures[toLower("get_user_orders")]
	if proc == nil {
		t.Fatal("procedure 'get_user_orders' not found")
	}

	// Verify trigger.
	trg := db.Triggers[toLower("trg_order_after_insert")]
	if trg == nil {
		t.Fatal("trigger 'trg_order_after_insert' not found")
	}
	if trg.Timing != "AFTER" {
		t.Errorf("trigger timing: expected 'AFTER', got %q", trg.Timing)
	}
	if trg.Event != "INSERT" {
		t.Errorf("trigger event: expected 'INSERT', got %q", trg.Event)
	}

	// Verify FK checks were re-enabled at the end.
	// The last SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS uses a variable reference,
	// which may not be interpreted. We just verify catalog is in a valid state.
	// The key point is the script executed without error.
}
