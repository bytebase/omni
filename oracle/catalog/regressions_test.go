package catalog

import (
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
)

func catalogFromSchema(schema *metadatapb.SchemaMetadata) *Catalog {
	return &Catalog{meta: &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{schema}}}
}

func migrationFor(t *testing.T, from, to *Catalog) string {
	t.Helper()
	return GenerateMigration(from, to, nil).SQL()
}

func TestNormalizeExpressionKeepsLiteralWhitespace(t *testing.T) {
	if normalizeExpression("STATUS = 'A  B'") == normalizeExpression("STATUS = 'A B'") {
		t.Fatal("whitespace inside a string literal was normalized away")
	}
	if normalizeExpression("A   =\n B") != normalizeExpression("A = B") {
		t.Fatal("whitespace outside a literal was not normalized")
	}
}

func TestDiffAltersSequence(t *testing.T) {
	oldMeta := &metadatapb.SequenceMetadata{Name: "S", Increment: "1", CacheSize: "20"}
	newMeta := &metadatapb.SequenceMetadata{Name: "S", Increment: "10", CacheSize: "40"}
	from := catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{oldMeta}})
	to := catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{newMeta}})
	if Diff(from, to).IsEmpty() {
		t.Fatal("a changed sequence produced an empty diff")
	}
	sql := migrationFor(t, from, to)
	if !strings.Contains(sql, `ALTER SEQUENCE "S"`) || !strings.Contains(sql, `INCREMENT BY 10`) || !strings.Contains(sql, `CACHE 40`) {
		t.Fatalf("sequence alteration missing:\n%s", sql)
	}
}

func TestMigrationKeepsDefaultOnNull(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Default: "1", Position: 1}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Default: "1", DefaultOnNull: true, Position: 1}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, `DEFAULT ON NULL 1`) {
		t.Fatalf("DEFAULT ON NULL change missing:\n%s", sql)
	}

	created := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(created, `DEFAULT ON NULL 1`) {
		t.Fatalf("DEFAULT ON NULL missing from CREATE TABLE:\n%s", created)
	}
}

func TestExtractCrossSchemaForeignKey(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE C (ID NUMBER, P_ID NUMBER, CONSTRAINT C_FK FOREIGN KEY (P_ID) REFERENCES OTHER.P(ID))`)
	if err != nil {
		t.Fatal(err)
	}
	fk := meta.Schemas[0].Tables[0].ForeignKeys[0]
	if fk.ReferencedSchema != "OTHER" {
		t.Fatalf("referenced schema = %q, want OTHER", fk.ReferencedSchema)
	}
}

func TestNationalCharacterLengthsAreNotDoubled(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (A NVARCHAR2(100), B NCHAR(100))`)
	if err != nil {
		t.Fatal(err)
	}
	columns := meta.Schemas[0].Tables[0].Columns
	if columns[0].Type != "NVARCHAR2(100)" || columns[1].Type != "NCHAR(100)" {
		t.Fatalf("types = %q, %q; want NVARCHAR2(100), NCHAR(100)", columns[0].Type, columns[1].Type)
	}
}

func TestViewTriggerSurvivesMigration(t *testing.T) {
	trigger := &metadatapb.TriggerMetadata{Name: "V_TRG", Body: "CREATE OR REPLACE TRIGGER V_TRG INSTEAD OF INSERT ON V FOR EACH ROW BEGIN NULL; END;"}
	oldView := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT 1 AS X FROM DUAL"}
	newView := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT 2 AS X FROM DUAL", Triggers: []*metadatapb.TriggerMetadata{trigger}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Views: []*metadatapb.ViewMetadata{oldView}}), catalogFromSchema(&metadatapb.SchemaMetadata{Views: []*metadatapb.ViewMetadata{newView}}))
	if !strings.Contains(sql, "CREATE OR REPLACE TRIGGER V_TRG") {
		t.Fatalf("view trigger missing from migration:\n%s", sql)
	}
}

func TestBitmapIndexMigration(t *testing.T) {
	table := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "T_BITMAP", Type: "BITMAP", Expressions: []string{"C"}, Visible: true}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{table}}))
	if !strings.Contains(sql, `CREATE BITMAP INDEX "T_BITMAP"`) {
		t.Fatalf("bitmap modifier missing from migration:\n%s", sql)
	}
}

func TestConstraintColumnQuotesCaseSensitiveName(t *testing.T) {
	table := &metadatapb.TableMetadata{
		Name:    "T",
		Columns: []*metadatapb.ColumnMetadata{{Name: "CamelCase", Type: "NUMBER", Position: 1}},
		Indexes: []*metadatapb.IndexMetadata{{Name: "T_PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"CamelCase"}}},
	}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{table}}))
	if !strings.Contains(sql, `PRIMARY KEY ("CamelCase")`) {
		t.Fatalf("case-sensitive key column not quoted:\n%s", sql)
	}
	if !strings.Contains(sql, `CREATE TABLE "T"`) {
		t.Fatalf("table missing:\n%s", sql)
	}
}

func TestConstraintRenameKeepsIncomingForeignKeys(t *testing.T) {
	columns := []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}
	oldTable := &metadatapb.TableMetadata{Name: "P", Columns: columns, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_OLD", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	newTable := &metadatapb.TableMetadata{Name: "P", Columns: columns, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_NEW", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, `RENAME CONSTRAINT "PK_OLD" TO "PK_NEW"`) {
		t.Fatalf("constraint rename missing:\n%s", sql)
	}
	if strings.Contains(sql, `DROP CONSTRAINT "PK_OLD"`) {
		t.Fatalf("constraint dropped instead of renamed:\n%s", sql)
	}
}
