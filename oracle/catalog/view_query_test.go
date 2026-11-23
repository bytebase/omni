package catalog

import (
	"errors"
	"testing"
)

func TestViewQueryDerivesColumns(t *testing.T) {
	c, err := LoadSQL(`
CREATE TABLE dept (deptno NUMBER(2) PRIMARY KEY, dname VARCHAR2(14) NOT NULL, loc VARCHAR2(13));
CREATE VIEW v AS SELECT d.*, UPPER(loc) up FROM dept d;
`)
	if err != nil {
		t.Fatal(err)
	}
	views := c.Metadata().GetSchemas()[0].GetViews()
	if len(views) != 1 {
		t.Fatalf("views = %d, want 1", len(views))
	}
	v := views[0]
	if want := `SELECT d."DEPTNO",d."DNAME",d."LOC", UPPER(loc) up FROM dept d`; v.GetDefinition() != want {
		t.Errorf("definition = %q, want %q", v.GetDefinition(), want)
	}
	want := []struct {
		name, typ string
		nullable  bool
	}{
		{"DEPTNO", "NUMBER(2)", false},
		{"DNAME", "VARCHAR2(14 BYTE)", false},
		{"LOC", "VARCHAR2(13 BYTE)", true},
		{"UP", "VARCHAR2(13 BYTE)", true},
	}
	if len(v.GetColumns()) != len(want) {
		t.Fatalf("columns = %d, want %d", len(v.GetColumns()), len(want))
	}
	for i, w := range want {
		got := v.GetColumns()[i]
		if got.GetName() != w.name || got.GetType() != w.typ || got.GetNullable() != w.nullable {
			t.Errorf("column %d = %s %s nullable=%v, want %s %s nullable=%v", i, got.GetName(), got.GetType(), got.GetNullable(), w.name, w.typ, w.nullable)
		}
	}
}

func TestViewQueryValidation(t *testing.T) {
	const missing = `CREATE VIEW v AS SELECT * FROM nosuch`

	// By default a view over objects the catalog does not hold loads, as a
	// schema dump that names other schemas must.
	c, err := LoadSQL(missing)
	if err != nil {
		t.Fatalf("default config refused the view: %v", err)
	}
	if v := c.Metadata().GetSchemas()[0].GetViews()[0]; v.GetDefinition() != "SELECT * FROM nosuch" || len(v.GetColumns()) != 0 {
		t.Errorf("view = %q with %d columns, want its text and no columns", v.GetDefinition(), len(v.GetColumns()))
	}

	// The pinned server's configuration refuses it as Oracle does.
	c = New()
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(missing, nil)
	if err != nil {
		t.Fatal(err)
	}
	var refusal *Error
	if len(results) != 1 || !errors.As(results[0].Error, &refusal) || refusal.Code != "00942" {
		t.Fatalf("results = %+v, want ORA-00942", results)
	}

	// FORCE keeps the view either way.
	results, err = c.Exec(`CREATE FORCE VIEW v AS SELECT * FROM nosuch`, nil)
	if err != nil || results[0].Error != nil {
		t.Fatalf("FORCE view refused: %v %+v", err, results)
	}
}
