package catalog

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

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
	if !strings.Contains(sql, `DEFAULT ON NULL 1`) || strings.Contains(sql, `DROP IDENTITY`) {
		t.Fatalf("DEFAULT ON NULL change wrong:\n%s", sql)
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

func TestDefinitionRendersIdentity(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{
		Name: "T",
		Columns: []*metadatapb.ColumnMetadata{{
			Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true,
			IdentityGeneration: metadatapb.ColumnMetadata_ALWAYS, IdentitySeed: 5, IdentityIncrement: 2,
		}},
	}}}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2)`) {
		t.Fatalf("identity clause missing:\n%s", definition)
	}
	created := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{}), &Catalog{meta: meta})
	if !strings.Contains(created, `GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2)`) {
		t.Fatalf("identity clause missing from migration:\n%s", created)
	}
}

func TestDefinitionDefersCyclicForeignKeys(t *testing.T) {
	a := &metadatapb.TableMetadata{Name: "A", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "B_ID", Type: "NUMBER", Position: 2}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "A_FK", Columns: []string{"B_ID"}, ReferencedTable: "B", ReferencedColumns: []string{"ID"}}}}
	b := &metadatapb.TableMetadata{Name: "B", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "A_ID", Type: "NUMBER", Position: 2}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "B_FK", Columns: []string{"A_ID"}, ReferencedTable: "A", ReferencedColumns: []string{"ID"}}}}
	definition, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{a, b}}}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(definition, `CONSTRAINT "A_FK"`) && strings.Contains(definition[:strings.Index(definition, `CREATE TABLE "B"`)], `CONSTRAINT "A_FK"`) {
		t.Fatalf("cyclic foreign key is embedded in CREATE TABLE:\n%s", definition)
	}
	if !strings.Contains(definition, `ALTER TABLE "A" ADD CONSTRAINT "A_FK"`) || !strings.Contains(definition, `ALTER TABLE "B" ADD CONSTRAINT "B_FK"`) {
		t.Fatalf("cyclic foreign keys were not deferred:\n%s", definition)
	}
}

func TestMigrationSkipsUnnamedSchema(t *testing.T) {
	target, err := LoadSQL(`CREATE TABLE T (ID NUMBER)`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, New(), target)
	if strings.Contains(sql, `CREATE USER`) {
		t.Fatalf("unnamed schema became a user:\n%s", sql)
	}
	if !strings.Contains(sql, `CREATE TABLE "T"`) {
		t.Fatalf("table missing:\n%s", sql)
	}
}

func TestMigrationRestoresMaterializedViewIndexes(t *testing.T) {
	oldView := &metadatapb.MaterializedViewMetadata{Name: "MV", Definition: "SELECT 1 AS X FROM DUAL"}
	newView := &metadatapb.MaterializedViewMetadata{Name: "MV", Definition: "SELECT 2 AS X FROM DUAL", Indexes: []*metadatapb.IndexMetadata{{Name: "MV_IDX", Expressions: []string{"X"}, Visible: true}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{MaterializedViews: []*metadatapb.MaterializedViewMetadata{oldView}}), catalogFromSchema(&metadatapb.SchemaMetadata{MaterializedViews: []*metadatapb.MaterializedViewMetadata{newView}}))
	if !strings.Contains(sql, `CREATE INDEX "MV_IDX" ON "MV"`) {
		t.Fatalf("materialized-view index missing:\n%s", sql)
	}
}

func TestMigrationDefersForeignKeyUntilKeyAddition(t *testing.T) {
	oldParent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}
	newParent := &metadatapb.TableMetadata{Name: "P", Columns: oldParent.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "P_PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	child := &metadatapb.TableMetadata{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldParent}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newParent, child}}))
	key := strings.Index(sql, `ADD CONSTRAINT "P_PK"`)
	fk := strings.Index(sql, `ADD CONSTRAINT "C_FK"`)
	if key < 0 || fk < 0 || key > fk {
		t.Fatalf("child foreign key was not deferred until after the key:\n%s", sql)
	}
}

func TestMigrationRendersInvisibleIndex(t *testing.T) {
	table := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "T_IDX", Expressions: []string{"C"}}}}
	definition, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{table}}}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `CREATE INDEX "T_IDX" ON "T" ("C" ASC) INVISIBLE;`) {
		t.Fatalf("invisible index missing from definition:\n%s", definition)
	}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{table}}))
	if !strings.Contains(sql, `CREATE INDEX "T_IDX" ON "T" ("C" ASC) INVISIBLE;`) {
		t.Fatalf("invisible index missing from migration:\n%s", sql)
	}
}

func TestMigrationAlterProcedureUsesOrReplace(t *testing.T) {
	oldProc := &metadatapb.ProcedureMetadata{Name: "P", Definition: "CREATE PROCEDURE P IS BEGIN NULL; END;"}
	newProc := &metadatapb.ProcedureMetadata{Name: "P", Definition: "CREATE PROCEDURE P IS BEGIN DBMS_OUTPUT.PUT_LINE('x'); END;"}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Procedures: []*metadatapb.ProcedureMetadata{oldProc}}), catalogFromSchema(&metadatapb.SchemaMetadata{Procedures: []*metadatapb.ProcedureMetadata{newProc}}))
	if !strings.Contains(sql, "CREATE OR REPLACE PROCEDURE P IS") {
		t.Fatalf("altered procedure is not CREATE OR REPLACE:\n%s", sql)
	}
}

func TestRenameIdentifierKeepsLiterals(t *testing.T) {
	if got := renameIdentifier(`CODE <> 'CODE' AND "CODE" > 0`, "CODE", "NEW_CODE"); got != `NEW_CODE <> 'CODE' AND "NEW_CODE" > 0` {
		t.Fatalf("renamed expression = %q", got)
	}
}

func TestMigrationViewCommentUsesTable(t *testing.T) {
	oldView := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT 1 AS X FROM DUAL", Comment: "old"}
	newView := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT 2 AS X FROM DUAL", Comment: "new"}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Views: []*metadatapb.ViewMetadata{oldView}}), catalogFromSchema(&metadatapb.SchemaMetadata{Views: []*metadatapb.ViewMetadata{newView}}))
	if !strings.Contains(sql, `COMMENT ON TABLE "V" IS 'new';`) || strings.Contains(sql, `COMMENT ON VIEW`) {
		t.Fatalf("view comment is not COMMENT ON TABLE:\n%s", sql)
	}
}

func TestIdentityOptionsFromSQL(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (ID NUMBER GENERATED BY DEFAULT AS IDENTITY (START WITH 42 INCREMENT BY 7))`)
	if err != nil {
		t.Fatal(err)
	}
	column := meta.Schemas[0].Tables[0].Columns[0]
	if !column.IsIdentity || column.IdentityGeneration != metadatapb.ColumnMetadata_BY_DEFAULT || column.IdentitySeed != 42 || column.IdentityIncrement != 7 {
		t.Fatalf("identity = %+v", column)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `GENERATED BY DEFAULT AS IDENTITY (START WITH 42 INCREMENT BY 7)`) {
		t.Fatalf("identity options missing:\n%s", definition)
	}
}

func TestMigrationAltersIdentity(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_ALWAYS, IdentitySeed: 1, IdentityIncrement: 1}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, IdentitySeed: 42, IdentityIncrement: 7}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, `GENERATED BY DEFAULT AS IDENTITY (START WITH 42 INCREMENT BY 7)`) {
		t.Fatalf("identity alteration missing:\n%s", sql)
	}
}

func TestViewCommentOnlyDoesNotRecreate(t *testing.T) {
	oldView := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT 1 AS X FROM DUAL", Comment: "old"}
	newView := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT 1 AS X FROM DUAL", Comment: "new"}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Views: []*metadatapb.ViewMetadata{oldView}}), catalogFromSchema(&metadatapb.SchemaMetadata{Views: []*metadatapb.ViewMetadata{newView}}))
	if !strings.Contains(sql, `COMMENT ON TABLE "V" IS 'new';`) || strings.Contains(sql, "CREATE VIEW") || strings.Contains(sql, "DROP VIEW") {
		t.Fatalf("comment-only view change recreated the view:\n%s", sql)
	}
}

func TestExtractRejectsCrossSchemaMerge(t *testing.T) {
	if _, err := metadataFromSQL(`CREATE TABLE A.T (ID NUMBER); CREATE TABLE B.T (VAL NUMBER)`); err == nil {
		t.Fatal("two owners were merged into one schema")
	}
}

func TestDefinitionKeepsForeignKeyOnDelete(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}, OnDelete: "CASCADE"}}}}}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `REFERENCES "P" ("ID") ON DELETE CASCADE`) {
		t.Fatalf("foreign-key delete action missing:\n%s", definition)
	}
}

func TestMigrationDropsObjectsFromUnnamedSchema(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE T (ID NUMBER)`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, from, New())
	if !strings.Contains(sql, `DROP TABLE "T"`) || strings.Contains(sql, `DROP USER`) {
		t.Fatalf("unnamed schema objects were not dropped:\n%s", sql)
	}
}

func TestDefinitionOrdersViewAfterDependency(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE Z (ID NUMBER); CREATE VIEW A AS SELECT ID FROM Z`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(definition, `CREATE TABLE "Z"`) > strings.Index(definition, `CREATE VIEW "A"`) {
		t.Fatalf("view precedes its dependency:\n%s", definition)
	}
}

func TestFunctionBasedIndexKeepsQuotedColumn(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T ("CamelCase" NUMBER); CREATE INDEX I ON T ("CamelCase" DESC)`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `("CamelCase" DESC)`) {
		t.Fatalf("quoted index column lost:\n%s", definition)
	}
}

func TestMigrationDefersAlteredTableForeignKeyUntilKeyAddition(t *testing.T) {
	oldParent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}
	newParent := &metadatapb.TableMetadata{Name: "P", Columns: oldParent.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "P_PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	oldChild := &metadatapb.TableMetadata{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}}
	newChild := &metadatapb.TableMetadata{Name: "C", Columns: oldChild.Columns, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldParent, oldChild}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newParent, newChild}}))
	key := strings.Index(sql, `ADD CONSTRAINT "P_PK"`)
	fk := strings.Index(sql, `ADD CONSTRAINT "C_FK"`)
	if key < 0 || fk < 0 || key > fk {
		t.Fatalf("altered table's foreign key was not deferred:\n%s", sql)
	}
}

func TestSchemaDropDropsForeignKeysFirst(t *testing.T) {
	parent := &metadatapb.TableMetadata{Name: "A", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}
	child := &metadatapb.TableMetadata{Name: "Z", Columns: []*metadatapb.ColumnMetadata{{Name: "A_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "Z_FK", Columns: []string{"A_ID"}, ReferencedTable: "A", ReferencedColumns: []string{"ID"}}}}
	from := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "S", Tables: []*metadatapb.TableMetadata{parent, child}})
	to := New()
	sql := migrationFor(t, from, to)
	fk := strings.Index(sql, `DROP CONSTRAINT "Z_FK"`)
	table := strings.Index(sql, `DROP TABLE "S"."A"`)
	if fk < 0 || table < 0 || fk > table {
		t.Fatalf("schema drop did not detach the foreign key first:\n%s", sql)
	}
}

func TestRenameColumnUpdatesForeignKeys(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{
		{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}},
		{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}},
	}})
	if _, err := c.Exec(`ALTER TABLE "C" RENAME COLUMN "P_ID" TO "PARENT_ID"`, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE "P" RENAME COLUMN "ID" TO "PID"`, nil); err != nil {
		t.Fatal(err)
	}
	child := tableByName(c.Metadata().Schemas[0])["C"]
	fk := child.ForeignKeys[0]
	if !slices.Equal(fk.Columns, []string{"PARENT_ID"}) || !slices.Equal(fk.ReferencedColumns, []string{"PID"}) {
		t.Fatalf("foreign key columns = %v -> %v", fk.Columns, fk.ReferencedColumns)
	}
}

func TestExecPreservesUnhandledSchemaMetadata(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Synonyms: []*metadatapb.SynonymMetadata{{Name: "S", TargetSchema: "X", TargetName: "T"}},
		Tables:   []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}},
	}}}
	c := New()
	if _, err := c.LoadMetadata(context.Background(), meta); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE "T" ADD "C" NUMBER`, nil); err != nil {
		t.Fatal(err)
	}
	if len(c.Metadata().Schemas[0].Synonyms) != 1 {
		t.Fatalf("synonyms lost: %+v", c.Metadata().Schemas[0])
	}
}

func TestExecDropsColumnDependents(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{
		Name:             "T",
		Columns:          []*metadatapb.ColumnMetadata{{Name: "A", Type: "NUMBER", Position: 1}, {Name: "B", Type: "NUMBER", Position: 2}},
		Indexes:          []*metadatapb.IndexMetadata{{Name: "T_A", Expressions: []string{"A"}, Visible: true}},
		CheckConstraints: []*metadatapb.CheckConstraintMetadata{{Name: "T_A_CK", Expression: `"A" > 0`}},
	}}})
	if _, err := c.Exec(`ALTER TABLE "T" DROP COLUMN "A"`, nil); err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if len(table.Indexes) != 0 || len(table.CheckConstraints) != 0 {
		t.Fatalf("dependent metadata survived: %+v", table)
	}
}

func TestRenameTableUpdatesForeignKeys(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{
		{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}},
		{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}},
	}})
	if _, err := c.Exec(`ALTER TABLE "P" RENAME TO "P_NEW"`, nil); err != nil {
		t.Fatal(err)
	}
	child := tableByName(c.Metadata().Schemas[0])["C"]
	if child.ForeignKeys[0].ReferencedTable != "P_NEW" {
		t.Fatalf("foreign key still references %q", child.ForeignKeys[0].ReferencedTable)
	}
}

func TestDefinitionRendersSequenceBounds(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Sequences: []*metadatapb.SequenceMetadata{{Name: "S", MinValue: "1", MaxValue: "9999", CacheSize: "20", Cycle: true}}}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `MINVALUE 1`) || !strings.Contains(definition, `CACHE 20`) {
		t.Fatalf("sequence bounds missing:\n%s", definition)
	}
}

func TestQuotedUserDefinedTypeKeepsSpelling(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (A "Address")`)
	if err != nil {
		t.Fatal(err)
	}
	if got := meta.Schemas[0].Tables[0].Columns[0].Type; got != `"Address"` {
		t.Fatalf("type = %q", got)
	}
}

func TestMigrationRestartsSequenceStart(t *testing.T) {
	oldSeq := &metadatapb.SequenceMetadata{Name: "S", Start: "1"}
	newSeq := &metadatapb.SequenceMetadata{Name: "S", Start: "10"}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{oldSeq}}), catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{newSeq}}))
	if !strings.Contains(sql, `RESTART START WITH 10`) {
		t.Fatalf("sequence restart missing:\n%s", sql)
	}
}

func TestMigrationResetsSequenceIncrementToDefault(t *testing.T) {
	oldSeq := &metadatapb.SequenceMetadata{Name: "S", Increment: "7"}
	newSeq := &metadatapb.SequenceMetadata{Name: "S"}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{oldSeq}}), catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{newSeq}}))
	if !strings.Contains(sql, `INCREMENT BY 1`) {
		t.Fatalf("default increment not restored:\n%s", sql)
	}
}

func TestDropOrderResolvesSameSchemaForeignKeys(t *testing.T) {
	parent := &metadatapb.TableMetadata{Name: "A", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}
	child := &metadatapb.TableMetadata{Name: "Z", Columns: []*metadatapb.ColumnMetadata{{Name: "A_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "Z_FK", Columns: []string{"A_ID"}, ReferencedTable: "A", ReferencedColumns: []string{"ID"}}}}
	from := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{parent, child}})
	to := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP"})
	sql := migrationFor(t, from, to)
	childDrop := strings.Index(sql, `DROP TABLE "APP"."Z"`)
	parentDrop := strings.Index(sql, `DROP TABLE "APP"."A"`)
	if childDrop < 0 || parentDrop < 0 || childDrop > parentDrop {
		t.Fatalf("parent dropped before child:\n%s", sql)
	}
}

func TestUnresolvedExternalForeignKeyIsKept(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE C (P_ID NUMBER REFERENCES OTHER.P)`)
	if err != nil {
		t.Fatal(err)
	}
	fk := meta.Schemas[0].Tables[0].ForeignKeys[0]
	if fk.ReferencedSchema != "OTHER" || fk.ReferencedTable != "P" || len(fk.ReferencedColumns) != 0 {
		t.Fatalf("foreign key = %+v", fk)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `REFERENCES "OTHER"."P"`) || strings.Contains(definition, `"OTHER"."P" ()`) {
		t.Fatalf("implicit external reference not rendered:\n%s", definition)
	}
}

func TestVirtualColumnWithoutTypeStaysUntyped(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C VARCHAR2(10), V AS (UPPER(C)))`)
	if err != nil {
		t.Fatal(err)
	}
	column := meta.Schemas[0].Tables[0].Columns[1]
	if column.Type != "" {
		t.Fatalf("virtual column type = %q, want empty", column.Type)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(definition, `"V" NUMBER`) {
		t.Fatalf("virtual column gained a NUMBER type:\n%s", definition)
	}
}

func TestExecAppliesAlterSequence(t *testing.T) {
	c, err := LoadSQL(`CREATE SEQUENCE S INCREMENT BY 1; ALTER SEQUENCE S INCREMENT BY 7`)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Metadata().Schemas[0].Sequences[0].Increment; got != "7" {
		t.Fatalf("increment = %q", got)
	}
}

func TestDefinitionConnectsOwnerQualifiedDependencies(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Name: "APP",
		Tables: []*metadatapb.TableMetadata{
			{Name: "Z", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}},
			{Name: "Y", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "Y_FK", Columns: []string{"ID"}, ReferencedTable: "Z", ReferencedColumns: []string{"ID"}}}},
		},
		Views: []*metadatapb.ViewMetadata{{Name: "A", Definition: "SELECT ID FROM Y", DependencyColumns: []*metadatapb.DependencyColumn{{Schema: "APP", Table: "Y", Column: "*"}}}},
	}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(definition, `CREATE TABLE "Y"`) > strings.Index(definition, `CREATE VIEW "A"`) {
		t.Fatalf("view precedes its owner-qualified dependency:\n%s", definition)
	}
}

func TestModifyNullRemovesNamedNotNull(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER CONSTRAINT NN NOT NULL)`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T MODIFY C NULL`, nil); err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if len(table.CheckConstraints) != 0 {
		t.Fatalf("named NOT NULL survived: %+v", table.CheckConstraints)
	}
}

func TestFunctionPrefixHandlesWhitespace(t *testing.T) {
	meta, err := metadataFromSQL("CREATE OR REPLACE\nFUNCTION F RETURN NUMBER IS BEGIN RETURN 1; END;")
	if err != nil {
		t.Fatal(err)
	}
	definition := meta.Schemas[0].Functions[0].Definition
	if strings.Contains(strings.ToUpper(definition), "OR REPLACE") || !strings.HasPrefix(strings.ToUpper(definition), "FUNCTION F") {
		t.Fatalf("function definition = %q", definition)
	}
}

func TestNormalizeExpressionKeepsAlternativeQuoteWhitespace(t *testing.T) {
	if normalizeExpression(`C <> q'[it's  ok]'`) == normalizeExpression(`C <> q'[it's ok]'`) {
		t.Fatal("whitespace inside an alternative-quoted literal was normalized")
	}
}

func TestMigrationQualifiesIndexName(t *testing.T) {
	table := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "I", Expressions: []string{"C"}, Visible: true}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP"}), catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{table}}))
	if !strings.Contains(sql, `CREATE INDEX "APP"."I" ON "APP"."T"`) {
		t.Fatalf("index name is not owner-qualified:\n%s", sql)
	}
}

func TestExecDropsIdentity(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER GENERATED ALWAYS AS IDENTITY)`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T MODIFY ID DROP IDENTITY`, nil); err != nil {
		t.Fatal(err)
	}
	column := tableByName(c.Metadata().Schemas[0])["T"].Columns[0]
	if column.IsIdentity || column.IdentityGeneration != metadatapb.ColumnMetadata_IDENTITY_GENERATION_UNSPECIFIED {
		t.Fatalf("identity survived: %+v", column)
	}
}

func TestNoCacheSurvivesExtractionAndRendering(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE SEQUENCE S NOCACHE`)
	if err != nil {
		t.Fatal(err)
	}
	if got := meta.Schemas[0].Sequences[0].CacheSize; got != "0" {
		t.Fatalf("cache size = %q, want NOCACHE sentinel", got)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `NOCACHE`) {
		t.Fatalf("NOCACHE missing from definition:\n%s", definition)
	}
	oldSeq := &metadatapb.SequenceMetadata{Name: "S", CacheSize: "20"}
	newSeq := &metadatapb.SequenceMetadata{Name: "S", CacheSize: "0"}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{oldSeq}}), catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{newSeq}}))
	if !strings.Contains(sql, `NOCACHE`) {
		t.Fatalf("NOCACHE change missing from migration:\n%s", sql)
	}
}

func TestDefinitionRendersMaterializedViewIndexes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "MV", Definition: "SELECT 1 AS X FROM DUAL", Indexes: []*metadatapb.IndexMetadata{{Name: "MV_IDX", Expressions: []string{"X"}, Visible: true}}}}}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `CREATE INDEX "MV_IDX" ON "MV"`) {
		t.Fatalf("materialized-view index missing from definition:\n%s", definition)
	}
}

func TestKeyReplacementDetachesIncomingForeignKey(t *testing.T) {
	oldParent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	newParent := &metadatapb.TableMetadata{Name: "P", Columns: oldParent.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "UK_P", Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	child := &metadatapb.TableMetadata{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldParent, child}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newParent, child}}))
	drop := strings.Index(sql, `DROP CONSTRAINT "C_FK"`)
	dropKey := strings.Index(sql, `DROP CONSTRAINT "PK_P"`)
	createKey := strings.Index(sql, `ADD CONSTRAINT "UK_P"`)
	add := strings.LastIndex(sql, `ADD CONSTRAINT "C_FK"`)
	if drop < 0 || dropKey < 0 || createKey < 0 || add < 0 || !(drop < dropKey && createKey < add) {
		t.Fatalf("incoming foreign key was not detached and restored:\n%s", sql)
	}
}

func TestLoadMetadataCopiesSnapshot(t *testing.T) {
	table := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{table}}}}
	c := New()
	if _, err := c.LoadMetadata(context.Background(), meta); err != nil {
		t.Fatal(err)
	}
	table.Name = "CHANGED"
	if c.Metadata().Schemas[0].Tables[0].Name != "T" {
		t.Fatal("catalog aliases the caller's snapshot")
	}
}

func TestAddingOverlappingUniqueDoesNotRecreateForeignKey(t *testing.T) {
	oldParent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "CODE", Type: "NUMBER", Position: 2}}, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	newParent := &metadatapb.TableMetadata{Name: "P", Columns: oldParent.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}, {Name: "UK_P", Unique: true, IsConstraint: true, Expressions: []string{"ID", "CODE"}}}}
	child := &metadatapb.TableMetadata{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldParent, child}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newParent, child}}))
	if strings.Contains(sql, `ADD CONSTRAINT "C_FK"`) {
		t.Fatalf("unchanged foreign key was recreated:\n%s", sql)
	}
}

func TestRenameColumnUpdatesIndexAndGenerationExpressions(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C VARCHAR2(10), V AS (UPPER(C))); CREATE INDEX I ON T (LOWER(C))`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T RENAME COLUMN C TO D`, nil); err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if got := table.Indexes[0].Expressions[0]; !strings.Contains(got, "D") || strings.Contains(got, "C") {
		t.Fatalf("index expression = %q", got)
	}
	if got := table.Columns[1].Generation.Expression; strings.Contains(got, "C") {
		t.Fatalf("generation expression = %q", got)
	}
}

func TestDropCheckConstraintKeepsIndependentNotNull(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER NOT NULL, CONSTRAINT CK CHECK (C IS NOT NULL))`)
	if err != nil {
		t.Fatal(err)
	}
	c := &Catalog{meta: meta}
	if _, err := c.Exec(`ALTER TABLE T DROP CONSTRAINT CK`, nil); err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if table.Columns[0].Nullable || len(table.CheckConstraints) != 0 {
		t.Fatalf("column/check = %+v / %+v", table.Columns[0], table.CheckConstraints)
	}
}

func TestAlterTableRejectsAnotherOwner(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}}})
	results, err := c.Exec(`ALTER TABLE OTHER.T ADD C NUMBER`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("alter from another owner was applied")
	}
}

func TestRenameTableLeavesExternalForeignKey(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{
		{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}},
		{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedSchema: "OTHER", ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}},
	}})
	if _, err := c.Exec(`ALTER TABLE P RENAME TO P_NEW`, nil); err != nil {
		t.Fatal(err)
	}
	fk := tableByName(c.Metadata().Schemas[0])["C"].ForeignKeys[0]
	if fk.ReferencedTable != "P" || fk.ReferencedSchema != "OTHER" {
		t.Fatalf("external foreign key rewritten: %+v", fk)
	}
}

func TestIdentityRendersOnNull(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, DefaultOnNull: true}}}}}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `GENERATED BY DEFAULT ON NULL AS IDENTITY`) {
		t.Fatalf("ON NULL identity missing:\n%s", definition)
	}
}

func TestDefinitionRendersComments(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "T", Comment: "table comment", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1, Comment: "column comment"}}}}}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `COMMENT ON TABLE "T" IS 'table comment';`) || !strings.Contains(definition, `COMMENT ON COLUMN "T"."C" IS 'column comment';`) {
		t.Fatalf("comments missing:\n%s", definition)
	}
}

func TestConstraintQuotesReservedColumn(t *testing.T) {
	table := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ORDER", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "T_PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ORDER"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{table}}))
	if !strings.Contains(sql, `PRIMARY KEY ("ORDER")`) {
		t.Fatalf("reserved key column not quoted:\n%s", sql)
	}
}

func TestExecDropsPrimaryKeyWithoutName(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY); ALTER TABLE T DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	for _, index := range table.Indexes {
		if index.Primary {
			t.Fatalf("primary key survived: %+v", index)
		}
	}
}

func TestMigrationAltersIdentityOnNull(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, DefaultOnNull: true}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, `GENERATED BY DEFAULT ON NULL AS IDENTITY`) {
		t.Fatalf("identity ON NULL change missing:\n%s", sql)
	}
}

func TestViewDependenciesExcludeCTE(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER); CREATE VIEW Z AS WITH A AS (SELECT C FROM T) SELECT C FROM A; CREATE VIEW A AS SELECT C FROM Z`)
	if err != nil {
		t.Fatal(err)
	}
	var z *metadatapb.ViewMetadata
	for _, view := range meta.Schemas[0].Views {
		if view.Name == "Z" {
			z = view
		}
	}
	if z == nil {
		t.Fatal("view Z missing")
	}
	for _, dependency := range z.DependencyColumns {
		if dependency.Table == "A" {
			t.Fatalf("CTE name recorded as a dependency: %+v", z.DependencyColumns)
		}
	}
}

func TestDropRejectsAnotherOwner(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}}})
	results, err := c.Exec(`DROP TABLE OTHER.T`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("drop from another owner was applied")
	}
	if tableByName(c.Metadata().Schemas[0])["T"] == nil {
		t.Fatal("local table was deleted")
	}
}

func TestViewDependenciesKeepQuotedCTECase(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER); CREATE VIEW Z AS SELECT C FROM T; CREATE VIEW A AS WITH "z" AS (SELECT C FROM T) SELECT Z.C FROM Z JOIN "z" ON Z.C = "z".C`)
	if err != nil {
		t.Fatal(err)
	}
	var a *metadatapb.ViewMetadata
	for _, view := range meta.Schemas[0].Views {
		if view.Name == "A" {
			a = view
		}
	}
	if a == nil {
		t.Fatal("view A missing")
	}
	found := false
	for _, dependency := range a.DependencyColumns {
		if dependency.Table == "Z" {
			found = true
		}
		if dependency.Table == "z" {
			t.Fatalf("quoted CTE recorded as dependency: %+v", a.DependencyColumns)
		}
	}
	if !found {
		t.Fatalf("view Z dependency missing: %+v", a.DependencyColumns)
	}
}

func TestCreateTableIfNotExistsKeepsExistingColumns(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (ID NUMBER); CREATE TABLE IF NOT EXISTS T (ID NUMBER)`)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(meta.Schemas[0].Tables[0].Columns); got != 1 {
		t.Fatalf("columns = %d, want 1", got)
	}
}

func TestDefinitionDefersFunctionBasedIndexUntilFunction(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Tables:    []*metadatapb.TableMetadata{{Name: "A", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "I", Type: "FUNCTION-BASED NORMAL", Expressions: []string{"Z(ID)"}, Visible: true}}}},
		Functions: []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z(P NUMBER) RETURN NUMBER IS BEGIN RETURN P; END;"}},
	}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, "CREATE OR REPLACE FUNCTION Z")
	index := strings.Index(definition, `CREATE INDEX "I"`)
	if function < 0 || index < 0 || function > index {
		t.Fatalf("function-based index precedes its function:\n%s", definition)
	}
}

func TestRenameColumnUpdatesSelfForeignKey(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY, PARENT_ID NUMBER REFERENCES T(ID))`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T RENAME COLUMN ID TO NEW_ID`, nil); err != nil {
		t.Fatal(err)
	}
	fk := tableByName(c.Metadata().Schemas[0])["T"].ForeignKeys[0]
	if !slices.Equal(fk.ReferencedColumns, []string{"NEW_ID"}) {
		t.Fatalf("referenced columns = %v", fk.ReferencedColumns)
	}
}

func TestRenameColumnQuotesCaseSensitiveName(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER CHECK (C > 0))`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T RENAME COLUMN C TO "d"`, nil); err != nil {
		t.Fatal(err)
	}
	check := tableByName(c.Metadata().Schemas[0])["T"].CheckConstraints[0]
	if !strings.Contains(check.Expression, `"d"`) {
		t.Fatalf("check expression = %q", check.Expression)
	}
}

func TestDropUniqueKeepsOtherUniqueKeys(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (A NUMBER UNIQUE, B NUMBER UNIQUE); ALTER TABLE T DROP UNIQUE (A)`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	var unique [][]string
	for _, index := range table.Indexes {
		if index.Unique && !index.Primary {
			unique = append(unique, index.Expressions)
		}
	}
	if len(unique) != 1 || !slices.Equal(unique[0], []string{"B"}) {
		t.Fatalf("unique keys = %v", unique)
	}
}

func TestViewKeepsCheckOption(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER); CREATE VIEW V AS SELECT C FROM T WHERE C > 0 WITH CHECK OPTION`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, "WITH CHECK OPTION") {
		t.Fatalf("check option missing:\n%s", definition)
	}
}

func TestSequenceKeepsExplicitZero(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE SEQUENCE S MINVALUE 0 START WITH 0 MAXVALUE 100 NOCACHE`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"START WITH 0", "MINVALUE 0", "MAXVALUE 100", "NOCACHE"} {
		if !strings.Contains(definition, want) {
			t.Fatalf("%s missing:\n%s", want, definition)
		}
	}
}

func TestIdentityDefaultIsNotAltered(t *testing.T) {
	oldColumn := &metadatapb.ColumnMetadata{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, Default: `"APP"."ISEQ$$_1".nextval`}
	newColumn := cloneColumn(oldColumn)
	newColumn.DefaultOnNull = true
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{oldColumn}}}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{newColumn}}}}))
	if strings.Contains(sql, "DEFAULT NULL") || strings.Contains(sql, ".nextval") {
		t.Fatalf("identity default altered:\n%s", sql)
	}
}

func TestRenameColumnLeavesExternalForeignKey(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{
		{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}},
		{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedSchema: "OTHER", ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}},
	}})
	if _, err := c.Exec(`ALTER TABLE P RENAME COLUMN ID TO NEW_ID`, nil); err != nil {
		t.Fatal(err)
	}
	fk := tableByName(c.Metadata().Schemas[0])["C"].ForeignKeys[0]
	if !slices.Equal(fk.ReferencedColumns, []string{"ID"}) {
		t.Fatalf("external reference rewritten: %v", fk.ReferencedColumns)
	}
}

func TestDropParenthesizedColumns(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (A NUMBER, B NUMBER, C NUMBER); ALTER TABLE T DROP (A, B)`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if len(table.Columns) != 1 || table.Columns[0].Name != "C" {
		t.Fatalf("columns = %+v", table.Columns)
	}
}

func TestCreateTableAsSelectIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T AS SELECT 1 AS ID FROM DUAL`); err == nil {
		t.Fatal("CTAS installed an empty table")
	}
}

func TestRenameTableRewritesTriggerTarget(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, Triggers: []*metadatapb.TriggerMetadata{{Name: "TRG", Body: "CREATE OR REPLACE TRIGGER TRG BEFORE INSERT ON T FOR EACH ROW BEGIN NULL; END;"}}}}})
	if _, err := c.Exec(`ALTER TABLE T RENAME TO U`, nil); err != nil {
		t.Fatal(err)
	}
	body := tableByName(c.Metadata().Schemas[0])["U"].Triggers[0].Body
	if !strings.Contains(body, " ON U ") || strings.Contains(body, " ON T ") {
		t.Fatalf("trigger body = %q", body)
	}
}

func TestKeyReplacementWithColumnRebuildDropsForeignKeyOnce(t *testing.T) {
	oldColumn := &metadatapb.ColumnMetadata{Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}
	newColumn := cloneColumn(oldColumn)
	newColumn.Generation.Expression = "C+2"
	oldParent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, oldColumn}, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"V"}}}}
	newParent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, newColumn}, Indexes: []*metadatapb.IndexMetadata{{Name: "UK_P", Unique: true, IsConstraint: true, Expressions: []string{"V"}}}}
	child := &metadatapb.TableMetadata{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_V", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_V"}, ReferencedTable: "P", ReferencedColumns: []string{"V"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldParent, child}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newParent, child}}))
	if strings.Count(sql, `DROP CONSTRAINT "C_FK"`) != 1 || strings.Count(sql, `ADD CONSTRAINT "C_FK"`) != 1 {
		t.Fatalf("foreign key detached/restored more than once:\n%s", sql)
	}
}

func TestIdentityOnNullFromSQL(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (ID NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY)`)
	if err != nil {
		t.Fatal(err)
	}
	column := meta.Schemas[0].Tables[0].Columns[0]
	if !column.DefaultOnNull {
		t.Fatalf("identity ON NULL not extracted: %+v", column)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `GENERATED BY DEFAULT ON NULL AS IDENTITY`) {
		t.Fatalf("identity ON NULL missing:\n%s", definition)
	}
}

func TestDropColumnLeavesExternalForeignKey(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{
		{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "B", Type: "NUMBER", Position: 2}}},
		{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedSchema: "OTHER", ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}},
	}})
	if _, err := c.Exec(`ALTER TABLE P DROP COLUMN ID`, nil); err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["C"].ForeignKeys) != 1 {
		t.Fatal("external foreign key removed with a local column")
	}
}

func TestViewDependenciesKeepOuterReferenceBesideNestedCTE(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER); CREATE VIEW Z AS SELECT C FROM T; CREATE VIEW A AS SELECT C FROM Z WHERE EXISTS (WITH Z AS (SELECT C FROM T) SELECT C FROM Z)`)
	if err != nil {
		t.Fatal(err)
	}
	var a *metadatapb.ViewMetadata
	for _, view := range meta.Schemas[0].Views {
		if view.Name == "A" {
			a = view
		}
	}
	found := false
	for _, dependency := range a.DependencyColumns {
		if dependency.Table == "Z" {
			found = true
		}
	}
	if !found {
		t.Fatalf("outer view Z dependency missing: %+v", a.DependencyColumns)
	}
}

func TestTriggerCommentLineIsSignificant(t *testing.T) {
	a := &metadatapb.TriggerMetadata{Name: "TRG", Body: "-- note\n:NEW.C := 1;"}
	b := &metadatapb.TriggerMetadata{Name: "TRG", Body: "-- note :NEW.C := 1;"}
	if triggerEqual(a, b) {
		t.Fatal("a line comment swallowed a statement")
	}
}

func TestConstraintIndexRepresentationsCompareEqual(t *testing.T) {
	a := &metadatapb.IndexMetadata{Name: "PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}
	b := &metadatapb.IndexMetadata{Name: "PK", Primary: true, Unique: true, IsConstraint: true, Type: "NORMAL", Visible: true, Expressions: []string{"ID"}}
	if !indexEqual(a, b) {
		t.Fatal("synced and extracted constraint indexes differ")
	}
}

func TestKeyReplacementDetachesSelfForeignKey(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "PARENT_ID", Type: "NUMBER", Position: 2}}, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_T", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "FK_T", Columns: []string{"PARENT_ID"}, ReferencedTable: "T", ReferencedColumns: []string{"ID"}}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: oldTable.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "UK_T", Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}, ForeignKeys: oldTable.ForeignKeys}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, `DROP CONSTRAINT "FK_T"`) || !strings.Contains(sql, `ADD CONSTRAINT "FK_T"`) {
		t.Fatalf("self foreign key not detached:\n%s", sql)
	}
}

func TestRenameConstraintColumnKeepsRawName(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER PRIMARY KEY)`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T RENAME COLUMN C TO "d"`, nil); err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(c.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `PRIMARY KEY ("d")`) || strings.Contains(definition, `PRIMARY KEY (""d"")`) {
		t.Fatalf("constraint column wrong:\n%s", definition)
	}
}

func TestRenameQuotedColumnLeavesUnquotedReference(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER, "c" NUMBER, CONSTRAINT CK CHECK (C > 0))`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T RENAME COLUMN "c" TO D`, nil); err != nil {
		t.Fatal(err)
	}
	check := tableByName(c.Metadata().Schemas[0])["T"].CheckConstraints[0]
	if check.Expression != "C > 0" {
		t.Fatalf("check expression = %q", check.Expression)
	}
}

func TestTemporaryTableIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE GLOBAL TEMPORARY TABLE SESSION_DATA (ID NUMBER) ON COMMIT DELETE ROWS`); err == nil {
		t.Fatal("temporary table installed as permanent")
	}
}

// TestCreateOrReplaceTriggerOnAnotherTable holds CREATE OR REPLACE TRIGGER
// to Oracle, which does not move a trigger to another table (ORA-04095;
// testdata/programs/simulate.tsv records it).
func TestCreateOrReplaceTriggerOnAnotherTable(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE A (ID NUMBER); CREATE TABLE Z (ID NUMBER); CREATE TRIGGER TR BEFORE INSERT ON Z FOR EACH ROW BEGIN NULL; END;`)
	if err != nil {
		t.Fatal(err)
	}
	results, err := c.Exec(`CREATE OR REPLACE TRIGGER TR BEFORE INSERT ON A FOR EACH ROW BEGIN NULL; END;`, nil)
	if err != nil {
		t.Fatal(err)
	}
	var refusal *Error
	if len(results) != 1 || !errors.As(results[0].Error, &refusal) || refusal.Code != "04095" {
		t.Fatalf("results = %+v, want ORA-04095", results)
	}
	schema := c.Metadata().Schemas[0]
	if len(tableByName(schema)["Z"].Triggers) != 1 || len(tableByName(schema)["A"].Triggers) != 0 {
		t.Fatalf("triggers A=%d Z=%d", len(tableByName(schema)["A"].Triggers), len(tableByName(schema)["Z"].Triggers))
	}
}

func TestDropTableCascadeRemovesIncomingForeignKey(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE P (ID NUMBER PRIMARY KEY); CREATE TABLE C (P_ID NUMBER, CONSTRAINT C_FK FOREIGN KEY (P_ID) REFERENCES P(ID)); DROP TABLE P CASCADE CONSTRAINTS`)
	if err != nil {
		t.Fatal(err)
	}
	child := tableByName(c.Metadata().Schemas[0])["C"]
	if len(child.ForeignKeys) != 0 {
		t.Fatalf("incoming foreign key survived: %+v", child.ForeignKeys)
	}
}

func TestAlterSequenceRejectsAnotherOwner(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Sequences: []*metadatapb.SequenceMetadata{{Name: "S", Increment: "1"}}})
	results, err := c.Exec(`ALTER SEQUENCE OTHER.S INCREMENT BY 10`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("sequence alteration from another owner was applied")
	}
	if got := c.Metadata().Schemas[0].Sequences[0].Increment; got != "1" {
		t.Fatalf("increment = %q", got)
	}
}

func TestFunctionPrefixHandlesOrReplaceOnNewLine(t *testing.T) {
	meta, err := metadataFromSQL("CREATE\nOR REPLACE FUNCTION F RETURN NUMBER IS BEGIN RETURN 1; END;")
	if err != nil {
		t.Fatal(err)
	}
	definition := meta.Schemas[0].Functions[0].Definition
	if strings.Contains(strings.ToUpper(definition), "OR REPLACE") {
		t.Fatalf("definition = %q", definition)
	}
}

func TestTriggerAttributesMatchSyncedForm(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (ID NUMBER); CREATE TRIGGER TRG BEFORE INSERT ON T FOR EACH ROW BEGIN NULL; END;`)
	if err != nil {
		t.Fatal(err)
	}
	trigger := meta.Schemas[0].Tables[0].Triggers[0]
	if trigger.Event != "INSERT" || trigger.Timing != "BEFORE EACH ROW" {
		t.Fatalf("trigger = %+v", trigger)
	}
	synced := clone(trigger)
	synced.Body = "CREATE OR REPLACE TRIGGER APP.TRG BEFORE INSERT ON APP.T FOR EACH ROW BEGIN NULL; END;"
	if !triggerEqual(trigger, synced) {
		t.Fatal("synced and extracted triggers differ")
	}
}

func TestDefinitionOrdersViewsWhenTablesCycle(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE P (ID NUMBER, C_ID NUMBER REFERENCES C(ID)); CREATE TABLE C (ID NUMBER, P_ID NUMBER REFERENCES P(ID)); CREATE VIEW Z AS SELECT ID FROM P; CREATE VIEW A AS SELECT ID FROM Z;`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	z := strings.Index(definition, `CREATE VIEW "Z"`)
	a := strings.Index(definition, `CREATE VIEW "A"`)
	if z < 0 || a < 0 || z > a {
		t.Fatalf("views out of dependency order:\n%s", definition)
	}
}

func TestMigrationDefersFunctionBasedIndexUntilFunction(t *testing.T) {
	table := &metadatapb.TableMetadata{Name: "A", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "I", Type: "FUNCTION-BASED NORMAL", Expressions: []string{"Z(ID)"}, Visible: true}}}
	to := catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{table}, Functions: []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z(P NUMBER) RETURN NUMBER IS BEGIN RETURN P; END;"}}})
	sql := migrationFor(t, New(), to)
	function := strings.Index(sql, "CREATE OR REPLACE FUNCTION Z")
	index := strings.Index(sql, `CREATE INDEX "I"`)
	if function < 0 || index < 0 || function > index {
		t.Fatalf("function-based index precedes its function:\n%s", sql)
	}
}

func TestIdentityToOrdinaryDefault(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, Default: `"APP"."ISEQ$$_1".nextval`}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, Default: "42"}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, `DROP IDENTITY`) || !strings.Contains(sql, `DEFAULT 42`) {
		t.Fatalf("identity-to-default migration:\n%s", sql)
	}
}

func TestAlterTriggerDisable(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER); CREATE TRIGGER TRG BEFORE INSERT ON T FOR EACH ROW BEGIN NULL; END; ALTER TRIGGER TRG DISABLE`)
	if err != nil {
		t.Fatal(err)
	}
	trigger := tableByName(c.Metadata().Schemas[0])["T"].Triggers[0]
	if !strings.Contains(strings.ToUpper(trigger.Body), "ALTER TRIGGER") || !strings.Contains(strings.ToUpper(trigger.Body), "DISABLE") {
		t.Fatalf("trigger body = %q", trigger.Body)
	}
}

func TestDropVirtualColumnBeforeSource(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{
		{Name: "A", Type: "NUMBER", Position: 1},
		{Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "A+1"}},
		{Name: "B", Type: "NUMBER", Position: 3},
	}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "B", Type: "NUMBER", Position: 3}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	v := strings.Index(sql, `DROP COLUMN "V"`)
	a := strings.Index(sql, `DROP COLUMN "A"`)
	if v < 0 || a < 0 || v > a {
		t.Fatalf("virtual column dropped after its source:\n%s", sql)
	}
}

func TestRemovedSelfForeignKeyIsNotRestored(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "PARENT_ID", Type: "NUMBER", Position: 2}}, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_T", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "FK_T", Columns: []string{"PARENT_ID"}, ReferencedTable: "T", ReferencedColumns: []string{"ID"}}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: oldTable.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "UK_T", Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if strings.Count(sql, `DROP CONSTRAINT "FK_T"`) != 1 || strings.Contains(sql, `ADD CONSTRAINT "FK_T"`) {
		t.Fatalf("removed self foreign key handled wrong:\n%s", sql)
	}
}

func TestFallbackOrdersMaterializedViewDependencies(t *testing.T) {
	to, err := LoadSQL(`CREATE TABLE P (ID NUMBER, C_ID NUMBER REFERENCES C(ID)); CREATE TABLE C (ID NUMBER, P_ID NUMBER REFERENCES P(ID)); CREATE MATERIALIZED VIEW Z AS SELECT ID FROM P; CREATE MATERIALIZED VIEW A AS SELECT ID FROM Z;`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, New(), to)
	z := strings.Index(sql, `CREATE MATERIALIZED VIEW "Z"`)
	a := strings.Index(sql, `CREATE MATERIALIZED VIEW "A"`)
	if z < 0 || a < 0 || z > a {
		t.Fatalf("materialized views out of order:\n%s", sql)
	}
}

func TestFallbackSkipsCommentOnlyView(t *testing.T) {
	parent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "C_ID", Type: "NUMBER", Position: 2}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "P_FK", Columns: []string{"C_ID"}, ReferencedTable: "C", ReferencedColumns: []string{"ID"}}}}
	child := &metadatapb.TableMetadata{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "P_ID", Type: "NUMBER", Position: 2}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}}
	oldView := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT ID FROM P", Comment: "old", DependencyColumns: []*metadatapb.DependencyColumn{{Table: "P", Column: "*"}}}
	newView := cloneView(oldView)
	newView.Comment = "new"
	from := catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{parent, child}, Views: []*metadatapb.ViewMetadata{oldView}})
	to := catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{parent, child}, Views: []*metadatapb.ViewMetadata{newView}})
	sql := migrationFor(t, from, to)
	if strings.Contains(sql, "CREATE VIEW") || !strings.Contains(sql, `COMMENT ON TABLE "V" IS 'new';`) {
		t.Fatalf("comment-only view recreated:\n%s", sql)
	}
}

func TestTriggerLiteralOwnerIsSignificant(t *testing.T) {
	a := &metadatapb.TriggerMetadata{Name: "TRG", Body: "CREATE TRIGGER TRG BEFORE INSERT ON T FOR EACH ROW BEGIN :NEW.C := ' ON APP.X'; END;"}
	b := &metadatapb.TriggerMetadata{Name: "TRG", Body: "CREATE TRIGGER TRG BEFORE INSERT ON T FOR EACH ROW BEGIN :NEW.C := ' ON OTHER.X'; END;"}
	if triggerEqual(a, b) {
		t.Fatal("trigger literal owner was normalized away")
	}
}

func TestMaterializedViewKeepsRefreshPolicy(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (ID NUMBER); CREATE MATERIALIZED VIEW MV BUILD IMMEDIATE REFRESH COMPLETE ON COMMIT AS SELECT ID FROM T`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `CREATE MATERIALIZED VIEW "MV" BUILD IMMEDIATE REFRESH COMPLETE ON COMMIT AS SELECT`) {
		t.Fatalf("refresh policy missing:\n%s", definition)
	}
}

func TestRetainedVirtualColumnAltersBeforeSourceDrop(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{
		{Name: "A", Type: "NUMBER", Position: 1},
		{Name: "B", Type: "NUMBER", Position: 2},
		{Name: "V", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "A+1"}},
	}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{
		{Name: "B", Type: "NUMBER", Position: 2},
		{Name: "V", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "B+1"}},
	}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	alter := strings.Index(sql, `MODIFY ("V"`)
	drop := strings.Index(sql, `DROP COLUMN "A"`)
	if alter < 0 || drop < 0 || alter > drop {
		t.Fatalf("virtual column altered after source drop:\n%s", sql)
	}
}

func TestCreateIndexRejectsAnotherTableOwner(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}}}})
	results, err := c.Exec(`CREATE INDEX APP.I ON OTHER.T(C)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("index on another owner was applied")
	}
}

func TestInlineDisableThenEnable(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER); CREATE TRIGGER TRG BEFORE INSERT ON T DISABLE BEGIN NULL; END; ALTER TRIGGER TRG ENABLE`)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.ToUpper(tableByName(c.Metadata().Schemas[0])["T"].Triggers[0].Body)
	if strings.Contains(body, "DISABLE") {
		t.Fatalf("trigger still disabled: %s", body)
	}
}

func TestNonTableTriggerIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TRIGGER AUDIT_TR AFTER LOGON ON DATABASE BEGIN NULL; END;`); err == nil {
		t.Fatal("database trigger installed a phantom table")
	}
}

func TestRenameQuotedColumnKeepsConstraintOnUnquoted(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER, "c" NUMBER, CONSTRAINT PK PRIMARY KEY(C))`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T RENAME COLUMN "c" TO D`, nil); err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	for _, index := range table.Indexes {
		if index.Primary && !slices.Equal(index.Expressions, []string{"C"}) {
			t.Fatalf("primary key moved: %v", index.Expressions)
		}
	}
}

func TestDropQuotedColumnRemovesItsConstraint(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T ("c" NUMBER PRIMARY KEY, D NUMBER); ALTER TABLE T DROP COLUMN "c"`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	for _, index := range table.Indexes {
		if index.Primary {
			t.Fatalf("primary key survived: %+v", index)
		}
	}
}

func TestModifyNullKeepsIndependentCheck(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER NOT NULL, CONSTRAINT CK CHECK ("C" IS NOT NULL)); ALTER TABLE T MODIFY C NULL`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if !table.Columns[0].Nullable || len(table.CheckConstraints) != 1 {
		t.Fatalf("column/check = %+v / %+v", table.Columns[0], table.CheckConstraints)
	}
}

func TestDiffNormalizesUnnamedSyncedSchema(t *testing.T) {
	oldMeta := &metadatapb.DatabaseSchemaMetadata{Name: "APP", Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}}}}}}
	newCatalog, err := LoadSQL(`CREATE TABLE APP.T (C NUMBER NOT NULL)`)
	if err != nil {
		t.Fatal(err)
	}
	diff := Diff(&Catalog{meta: oldMeta}, newCatalog)
	if !diff.IsEmpty() {
		t.Fatalf("equivalent schemas diff: schemas=%d tables=%d", len(diff.SchemaChanges), len(diff.TableChanges))
	}
}

func TestVirtualColumnFunctionOrdersBeforeTable(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Tables:    []*metadatapb.TableMetadata{{Name: "A", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "Z(ID)"}}}}},
		Functions: []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z(P NUMBER) RETURN NUMBER IS BEGIN RETURN P; END;"}},
	}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, "CREATE OR REPLACE FUNCTION Z")
	table := strings.Index(definition, `CREATE TABLE "A"`)
	if function < 0 || table < 0 || function > table {
		t.Fatalf("table precedes its virtual-column function:\n%s", definition)
	}
	sql := migrationFor(t, New(), &Catalog{meta: meta})
	function = strings.Index(sql, "CREATE OR REPLACE FUNCTION Z")
	table = strings.Index(sql, `CREATE TABLE "A"`)
	if function < 0 || table < 0 || function > table {
		t.Fatalf("migration table precedes its virtual-column function:\n%s", sql)
	}
}

func TestTriggerDisableIgnoresBodyText(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER); CREATE TRIGGER TRG BEFORE INSERT ON T FOR EACH ROW BEGIN :NEW.ID := 1; END; ALTER TRIGGER TRG DISABLE`)
	if err != nil {
		t.Fatal(err)
	}
	body := tableByName(c.Metadata().Schemas[0])["T"].Triggers[0].Body
	if !strings.Contains(strings.ToUpper(body), `ALTER TRIGGER "TRG" DISABLE;`) {
		t.Fatalf("disable not recorded: %q", body)
	}
}

func TestRoutineRejectsAnotherOwner(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP"})
	results, err := c.Exec(`CREATE FUNCTION OTHER.F RETURN NUMBER IS BEGIN RETURN 1; END;`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("routine from another owner was installed")
	}
}

func TestConstraintNameSwapFallsBackToDropCreate(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "A", Type: "NUMBER", Position: 1}, {Name: "B", Type: "NUMBER", Position: 2}}, Indexes: []*metadatapb.IndexMetadata{{Name: "U_A", Unique: true, IsConstraint: true, Expressions: []string{"A"}}, {Name: "U_B", Unique: true, IsConstraint: true, Expressions: []string{"B"}}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: oldTable.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "U_A", Unique: true, IsConstraint: true, Expressions: []string{"B"}}, {Name: "U_B", Unique: true, IsConstraint: true, Expressions: []string{"A"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if strings.Contains(sql, "RENAME CONSTRAINT") {
		t.Fatalf("colliding rename emitted:\n%s", sql)
	}
}

func TestDeferrableForeignKeyIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE P (ID NUMBER PRIMARY KEY); CREATE TABLE C (P_ID NUMBER, CONSTRAINT C_FK FOREIGN KEY (P_ID) REFERENCES P(ID) DEFERRABLE INITIALLY DEFERRED)`); err == nil {
		t.Fatal("deferrable foreign key was accepted")
	}
}

func TestFallbackOrdersVirtualColumnFunctionBeforeTable(t *testing.T) {
	to, err := LoadSQL(`CREATE TABLE P (ID NUMBER, C_ID NUMBER REFERENCES C(ID)); CREATE TABLE C (ID NUMBER, P_ID NUMBER REFERENCES P(ID)); CREATE FUNCTION Z(N NUMBER) RETURN NUMBER IS BEGIN RETURN N; END; CREATE TABLE A (ID NUMBER, V AS (Z(ID)))`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(to.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, "CREATE OR REPLACE FUNCTION Z")
	table := strings.Index(definition, `CREATE TABLE "A"`)
	if function < 0 || table < 0 || function > table {
		t.Fatalf("definition table precedes function:\n%s", definition)
	}
	sql := migrationFor(t, New(), to)
	function = strings.Index(sql, "CREATE OR REPLACE FUNCTION Z")
	table = strings.Index(sql, `CREATE TABLE "A"`)
	if function < 0 || table < 0 || function > table {
		t.Fatalf("migration table precedes function:\n%s", sql)
	}
}

func TestVirtualAlterAddsNewSourceColumnFirst(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "A", Type: "NUMBER", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "A+1"}}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "B", Type: "NUMBER", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "B+1"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	add := strings.Index(sql, `ADD "B"`)
	alter := strings.Index(sql, `MODIFY ("V"`)
	drop := strings.Index(sql, `DROP COLUMN "A"`)
	if add < 0 || alter < 0 || drop < 0 || add > alter || alter > drop {
		t.Fatalf("new source column ordering wrong:\n%s", sql)
	}
}

func TestExecRejectsQualifiedDDLOnUnnamedSyncedSchema(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Name: "APP", Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}}}}}
	c := &Catalog{meta: meta, namedNotNull: map[string]bool{}}
	results, err := c.Exec(`ALTER TABLE OTHER.T ADD X NUMBER`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("qualified DDL from another owner was applied")
	}
}

func TestTriggerRejectsAnotherOwner(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}}})
	results, err := c.Exec(`CREATE OR REPLACE TRIGGER OTHER.TR BEFORE INSERT ON APP.T BEGIN NULL; END;`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("trigger from another owner was installed")
	}
}

func TestDropPrimaryKeyRestoresNullability(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY); ALTER TABLE T DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if !table.Columns[0].Nullable {
		t.Fatalf("primary-key nullability not restored: %+v", table.Columns[0])
	}
}

func TestStandaloneRenameTable(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER); RENAME T TO U`)
	if err != nil {
		t.Fatal(err)
	}
	schema := c.Metadata().Schemas[0]
	if tableByName(schema)["U"] == nil || tableByName(schema)["T"] != nil {
		t.Fatalf("standalone rename not applied: %+v", schema.Tables)
	}
}

func TestDropPrimaryKeyKeepsExplicitNotNull(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER NOT NULL, CONSTRAINT PK PRIMARY KEY(ID)); ALTER TABLE T DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if table.Columns[0].Nullable {
		t.Fatalf("explicit NOT NULL cleared: %+v", table.Columns[0])
	}
}

func TestDropUniqueMatchesQuotedColumnCase(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER, "c" NUMBER, CONSTRAINT U1 UNIQUE(C), CONSTRAINT U2 UNIQUE("c")); ALTER TABLE T DROP UNIQUE ("c")`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	var unique [][]string
	for _, index := range table.Indexes {
		if index.Unique && !index.Primary {
			unique = append(unique, index.Expressions)
		}
	}
	if len(unique) != 1 || !slices.Equal(unique[0], []string{"C"}) {
		t.Fatalf("unique keys = %v", unique)
	}
}

func TestRenameTableRewritesQualifiedTriggerTarget(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE APP.T (C NUMBER); CREATE TRIGGER APP.TR BEFORE INSERT ON APP.T BEGIN NULL; END; ALTER TABLE APP.T RENAME TO U`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["U"]
	if table == nil || len(table.Triggers) != 1 {
		t.Fatalf("renamed table/triggers: %+v", c.Metadata().Schemas[0].Tables)
	}
	if !strings.Contains(table.Triggers[0].Body, " ON APP.U ") {
		t.Fatalf("trigger body = %q", table.Triggers[0].Body)
	}
}

func TestCreateSchemaEmitsValidAuthentication(t *testing.T) {
	to, err := LoadSQL(`CREATE TABLE APP.T (C NUMBER)`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, New(), to)
	if !strings.Contains(sql, `CREATE USER "APP" NO AUTHENTICATION;`) {
		t.Fatalf("schema creation lacks authentication:\n%s", sql)
	}
}

func TestDropInlinePrimaryKeyKeepsExplicitNotNull(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER NOT NULL PRIMARY KEY); ALTER TABLE T DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	if tableByName(c.Metadata().Schemas[0])["T"].Columns[0].Nullable {
		t.Fatal("explicit NOT NULL cleared")
	}
}

func TestCommentRejectsAnotherOwner(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}}}})
	results, err := c.Exec(`COMMENT ON TABLE OTHER.T IS 'other'`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatal("comment from another owner was applied")
	}
}

func TestAlterIndexRenameAndVisibility(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE INDEX I ON T(C); ALTER INDEX I RENAME TO J; ALTER INDEX J INVISIBLE`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	found := false
	for _, index := range table.Indexes {
		if index.Name == "J" {
			found = true
			if index.Visible {
				t.Fatalf("index still visible: %+v", index)
			}
		}
	}
	if !found {
		t.Fatalf("renamed index missing: %+v", table.Indexes)
	}
}

func TestConstraintDisableIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER, CONSTRAINT CK CHECK (C > 0)); ALTER TABLE T DISABLE CONSTRAINT CK`); err == nil {
		t.Fatal("constraint disable was silently accepted")
	}
}

func TestViewFunctionOrdersBeforeView(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Views:     []*metadatapb.ViewMetadata{{Name: "A", Definition: "SELECT Z() AS C FROM DUAL"}},
		Functions: []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z RETURN NUMBER IS BEGIN RETURN 1; END;"}},
	}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, "CREATE OR REPLACE FUNCTION Z")
	view := strings.Index(definition, `CREATE VIEW "A"`)
	if function < 0 || view < 0 || function > view {
		t.Fatalf("view precedes its function:\n%s", definition)
	}
	sql := migrationFor(t, New(), &Catalog{meta: meta})
	function = strings.Index(sql, "CREATE OR REPLACE FUNCTION Z")
	view = strings.Index(sql, `CREATE VIEW "A"`)
	if function < 0 || view < 0 || function > view {
		t.Fatalf("migration view precedes its function:\n%s", sql)
	}
}

func TestRenameTableQuotesTriggerTarget(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE TRIGGER TR BEFORE INSERT ON T BEGIN NULL; END; ALTER TABLE T RENAME TO "u"`)
	if err != nil {
		t.Fatal(err)
	}
	body := tableByName(c.Metadata().Schemas[0])["u"].Triggers[0].Body
	if !strings.Contains(body, ` ON "u" `) {
		t.Fatalf("trigger target = %q", body)
	}
}

func TestDropColumnKeepsExternalSameNameForeignKey(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, {Name: "D", Type: "NUMBER", Position: 2}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "FK", Columns: []string{"D"}, ReferencedSchema: "OTHER", ReferencedTable: "T", ReferencedColumns: []string{"C"}}}}}})
	if _, err := c.Exec(`ALTER TABLE T DROP COLUMN C`, nil); err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["T"].ForeignKeys) != 1 {
		t.Fatal("external same-name foreign key removed")
	}
}

func TestRenameColumnSkipsAlternativeQuotedLiteral(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER CHECK (C <> 0)); ALTER TABLE T RENAME COLUMN C TO D`)
	if err != nil {
		t.Fatal(err)
	}
	_ = c
	expr := renameIdentifier(`C <> q'[don't C]'`, "C", "D")
	if expr != `D <> q'[don't C]'` {
		t.Fatalf("renamed expression = %q", expr)
	}
}

func TestDeferrableUniqueIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER, CONSTRAINT U UNIQUE(C) DEFERRABLE INITIALLY DEFERRED)`); err == nil {
		t.Fatal("deferrable unique constraint accepted")
	}
}

func TestDropDefaultOnNullBeforeRelaxingNotNull(t *testing.T) {
	oldColumn := &metadatapb.ColumnMetadata{Name: "C", Type: "NUMBER", Position: 1, Default: "1", DefaultOnNull: true, Nullable: false}
	newColumn := &metadatapb.ColumnMetadata{Name: "C", Type: "NUMBER", Position: 1, Default: "1", Nullable: true}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{oldColumn}}}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{newColumn}}}}))
	defaultDrop := strings.Index(sql, "DEFAULT NULL")
	null := strings.Index(sql, `MODIFY "C" NULL`)
	if defaultDrop < 0 || null < 0 || defaultDrop > null {
		t.Fatalf("DEFAULT ON NULL not removed first:\n%s", sql)
	}
}

func TestOrdinaryToIdentityIsReportedUnsupported(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, "not supported") || strings.Contains(sql, `MODIFY "ID" GENERATED`) {
		t.Fatalf("ordinary-to-identity conversion looks executable:\n%s", sql)
	}
}

func TestConstraintColumnWhitespaceIsSignificant(t *testing.T) {
	a := &metadatapb.IndexMetadata{Name: "U", Unique: true, IsConstraint: true, Expressions: []string{"A B"}}
	b := &metadatapb.IndexMetadata{Name: "U", Unique: true, IsConstraint: true, Expressions: []string{"A  B"}}
	if indexEqual(a, b) {
		t.Fatal("constraint column whitespace normalized away")
	}
}

func TestDefaultOnNullImpliesNotNull(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER DEFAULT ON NULL 1)`)
	if err != nil {
		t.Fatal(err)
	}
	column := tableByName(c.Metadata().Schemas[0])["T"].Columns[0]
	if column.Nullable || !column.DefaultOnNull {
		t.Fatalf("column = %+v", column)
	}
}

func TestDropColumnKeepsCheckInAlternativeQuote(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER, D VARCHAR2(30), CONSTRAINT CK CHECK (D <> q'[don't C]')); ALTER TABLE T DROP COLUMN C`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if len(table.CheckConstraints) != 1 {
		t.Fatalf("check removed: %+v", table.CheckConstraints)
	}
}

func TestDropConstraintKeepsStandaloneSameNameIndex(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER, D NUMBER, CONSTRAINT CK CHECK(C > 0)); CREATE UNIQUE INDEX CK ON T(D); ALTER TABLE T DROP CONSTRAINT CK`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	found := false
	for _, index := range table.Indexes {
		if index.Name == "CK" {
			found = true
		}
	}
	if !found || len(table.CheckConstraints) != 0 {
		t.Fatalf("index/check = %+v / %+v", table.Indexes, table.CheckConstraints)
	}
}

func TestRenamePrimaryKeyColumnMovesProvenance(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY); ALTER TABLE T RENAME COLUMN ID TO X; ALTER TABLE T DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if !table.Columns[0].Nullable {
		t.Fatalf("renamed primary-key column kept NOT NULL: %+v", table.Columns[0])
	}
}

func TestRenameConstraintLeavesStandaloneIndex(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER, D NUMBER, CONSTRAINT CK CHECK(C > 0)); CREATE UNIQUE INDEX CK ON T(D); ALTER TABLE T RENAME CONSTRAINT CK TO CK_NEW`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	found := false
	for _, index := range table.Indexes {
		if index.Name == "CK" {
			found = true
		}
	}
	if !found {
		t.Fatalf("standalone index renamed: %+v", table.Indexes)
	}
}

func TestDropColumnKeepsCheckFunctionName(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T ("ABS" NUMBER, D NUMBER, CONSTRAINT CK CHECK(ABS(D) > 0)); ALTER TABLE T DROP COLUMN "ABS"`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["T"].CheckConstraints) != 1 {
		t.Fatal("check depending on function name was removed")
	}
}

func TestIncomingForeignKeyFindsUnnamedSnapshotSchema(t *testing.T) {
	oldParent := &metadatapb.TableMetadata{Name: "P", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	newParent := &metadatapb.TableMetadata{Name: "P", Columns: oldParent.Columns, Indexes: []*metadatapb.IndexMetadata{{Name: "UK_P", Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}
	child := &metadatapb.TableMetadata{Name: "C", Columns: []*metadatapb.ColumnMetadata{{Name: "P_ID", Type: "NUMBER", Position: 1}}, ForeignKeys: []*metadatapb.ForeignKeyMetadata{{Name: "C_FK", Columns: []string{"P_ID"}, ReferencedTable: "P", ReferencedColumns: []string{"ID"}}}}
	from := &Catalog{meta: &metadatapb.DatabaseSchemaMetadata{Name: "APP", Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{oldParent, child}}}}}
	to := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{newParent, child}})
	sql := migrationFor(t, from, to)
	if !strings.Contains(sql, `DROP CONSTRAINT "C_FK"`) || !strings.Contains(sql, `ADD CONSTRAINT "C_FK"`) {
		t.Fatalf("incoming foreign key not detached:\n%s", sql)
	}
}

func TestDefinitionKeepsLiteralSequenceLikeDefault(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "VARCHAR2(30 BYTE)", Default: "'ISEQ$$_tag'", Position: 1}}}}}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `DEFAULT 'ISEQ$$_tag'`) {
		t.Fatalf("literal default dropped:\n%s", definition)
	}
}

func TestRenameColumnKeepsFunctionName(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T ("ABS" NUMBER, D NUMBER, CONSTRAINT CK CHECK(ABS(D) > 0)); ALTER TABLE T RENAME COLUMN "ABS" TO X`)
	if err != nil {
		t.Fatal(err)
	}
	check := tableByName(c.Metadata().Schemas[0])["T"].CheckConstraints[0]
	if !strings.Contains(check.Expression, "ABS(D)") {
		t.Fatalf("function call renamed: %q", check.Expression)
	}
}

func TestRenameTableMovesPrimaryKeyProvenance(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY); ALTER TABLE T RENAME TO U; ALTER TABLE U DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["U"]
	if !table.Columns[0].Nullable {
		t.Fatalf("renamed table's primary-key column kept NOT NULL: %+v", table.Columns[0])
	}
}

func TestFunctionReferenceAllowsWhitespace(t *testing.T) {
	if !expressionReferencesFunction("Z (C)", "Z") {
		t.Fatal("function call with whitespace not detected")
	}
}

func TestDropTableClearsNullabilityProvenance(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY); DROP TABLE T; CREATE TABLE T (ID NUMBER NOT NULL PRIMARY KEY); ALTER TABLE T DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	if tableByName(c.Metadata().Schemas[0])["T"].Columns[0].Nullable {
		t.Fatal("stale provenance cleared an explicit NOT NULL")
	}
}

func TestAlterTriggerRename(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER); CREATE TRIGGER TR BEFORE INSERT ON T BEGIN NULL; END; ALTER TRIGGER TR RENAME TO TR2; DROP TRIGGER TR2`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["T"].Triggers) != 0 {
		t.Fatal("renamed trigger survived its drop")
	}
}

func TestQualifiedStatementEstablishesOwner(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE TABLE APP.U (C NUMBER); CREATE TABLE OTHER.V (C NUMBER)`); err == nil {
		t.Fatal("a second owner was accepted after APP")
	}
}

func TestSetUnusedIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER, D NUMBER); ALTER TABLE T SET UNUSED COLUMN C`); err == nil {
		t.Fatal("SET UNUSED silently accepted")
	}
}

func TestFunctionReferenceContinuesAfterAlias(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Views:     []*metadatapb.ViewMetadata{{Name: "A", Definition: "SELECT 1 AS Z, Z() AS C FROM DUAL"}},
		Functions: []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z RETURN NUMBER IS BEGIN RETURN 1; END;"}},
	}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, "CREATE OR REPLACE FUNCTION Z")
	view := strings.Index(definition, `CREATE VIEW "A"`)
	if function < 0 || view < 0 || function > view {
		t.Fatalf("view precedes function after alias:\n%s", definition)
	}
}

func TestDropColumnSkipsBlockComment(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER, D NUMBER, CONSTRAINT CK CHECK(D /* C */ > 0)); ALTER TABLE T DROP COLUMN C`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["T"].CheckConstraints) != 1 {
		t.Fatal("check removed because of a comment")
	}
}

func TestPreAlterVirtualColumnDetachesIndex(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, {Name: "D", Type: "NUMBER", Position: 2}, {Name: "V", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}}, Indexes: []*metadatapb.IndexMetadata{{Name: "I", Expressions: []string{"V"}, Visible: true}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "D", Type: "NUMBER", Position: 2}, {Name: "V", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "D+1"}}}, Indexes: oldTable.Indexes}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	dropIndex := strings.Index(sql, `DROP INDEX "I"`)
	alter := strings.Index(sql, `MODIFY ("V"`)
	createIndex := strings.Index(sql, `CREATE INDEX "I"`)
	if dropIndex < 0 || alter < 0 || createIndex < 0 || !(dropIndex < alter && alter < createIndex) {
		t.Fatalf("index not detached around virtual alteration:\n%s", sql)
	}
}

func TestIdentityDropBeforeRelaxingNull(t *testing.T) {
	oldColumn := &metadatapb.ColumnMetadata{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, Nullable: false}
	newColumn := &metadatapb.ColumnMetadata{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{oldColumn}}}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{newColumn}}}}))
	drop := strings.Index(sql, "DROP IDENTITY")
	null := strings.Index(sql, `MODIFY "ID" NULL`)
	if drop < 0 || null < 0 || drop > null {
		t.Fatalf("identity not dropped first:\n%s", sql)
	}
}

func TestLoadMetadataResetsProvenance(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY)`)
	if err != nil {
		t.Fatal(err)
	}
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, Indexes: []*metadatapb.IndexMetadata{{Name: "PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}}}}}}}}
	if _, err := c.LoadMetadata(context.Background(), meta); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec(`ALTER TABLE T DROP PRIMARY KEY`, nil); err != nil {
		t.Fatal(err)
	}
	if tableByName(c.Metadata().Schemas[0])["T"].Columns[0].Nullable {
		t.Fatal("stale provenance cleared a replacement snapshot")
	}
}

func TestDropPrimaryKeyCascadeRemovesForeignKey(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE P (ID NUMBER PRIMARY KEY); CREATE TABLE C (P_ID NUMBER REFERENCES P(ID)); ALTER TABLE P DROP PRIMARY KEY CASCADE`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["C"].ForeignKeys) != 0 {
		t.Fatal("cascading key drop left the foreign key")
	}
}

func TestSequenceRestartRequiresKnownOldStart(t *testing.T) {
	oldSeq := &metadatapb.SequenceMetadata{Name: "S"}
	newSeq := &metadatapb.SequenceMetadata{Name: "S", Start: "1"}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{oldSeq}}), catalogFromSchema(&metadatapb.SchemaMetadata{Sequences: []*metadatapb.SequenceMetadata{newSeq}}))
	if strings.Contains(sql, "RESTART") {
		t.Fatalf("restart inferred from unknown start:\n%s", sql)
	}
}

func TestColumnAlterationBeforeDependentView(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "DATE", Position: 1}}}
	view := &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT EXTRACT(YEAR FROM C) AS Y FROM T", DependencyColumns: []*metadatapb.DependencyColumn{{Table: "T", Column: "C"}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}, Views: []*metadatapb.ViewMetadata{view}}))
	alter := strings.Index(sql, `MODIFY "C" DATE`)
	create := strings.Index(sql, `CREATE VIEW "V"`)
	if alter < 0 || create < 0 || alter > create {
		t.Fatalf("view created before column alteration:\n%s", sql)
	}
}

func TestEarlyVirtualAlterExcludesDroppedIndex(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, {Name: "D", Type: "NUMBER", Position: 2}, {Name: "V", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}}, Indexes: []*metadatapb.IndexMetadata{{Name: "I", Expressions: []string{"V"}, Visible: true}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "D", Type: "NUMBER", Position: 2}, {Name: "V", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "D+1"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if strings.Count(sql, `DROP INDEX "I"`) != 1 || strings.Contains(sql, `CREATE INDEX "I"`) {
		t.Fatalf("already-dropped index handled twice:\n%s", sql)
	}
}

func TestMaterializedViewKeepsRefreshSchedule(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER); CREATE MATERIALIZED VIEW MV REFRESH COMPLETE START WITH SYSDATE NEXT SYSDATE + 1 AS SELECT C FROM T`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, "START WITH SYSDATE NEXT SYSDATE + 1") {
		t.Fatalf("refresh schedule missing:\n%s", definition)
	}
}

func TestCascadeKeyDropOnlyRemovesThatKey(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE P (ID NUMBER, CODE NUMBER, CONSTRAINT PK PRIMARY KEY(ID), CONSTRAINT UK UNIQUE(CODE)); CREATE TABLE C (P_ID NUMBER REFERENCES P(ID)); ALTER TABLE P DROP CONSTRAINT UK CASCADE`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["C"].ForeignKeys) != 1 {
		t.Fatal("unrelated foreign key removed by cascading unique drop")
	}
}

func TestDropColumnClearsReusedProvenance(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER PRIMARY KEY, D NUMBER); ALTER TABLE T DROP COLUMN ID; ALTER TABLE T ADD ID NUMBER NOT NULL PRIMARY KEY; ALTER TABLE T DROP PRIMARY KEY`)
	if err != nil {
		t.Fatal(err)
	}
	for _, column := range tableByName(c.Metadata().Schemas[0])["T"].Columns {
		if column.Name == "ID" && column.Nullable {
			t.Fatal("reused column inherited stale provenance")
		}
	}
}

func TestRenameTableHandlesTriggerOnNewline(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE TRIGGER TR BEFORE INSERT ON
T BEGIN NULL; END; ALTER TABLE T RENAME TO U`)
	if err != nil {
		t.Fatal(err)
	}
	body := tableByName(c.Metadata().Schemas[0])["U"].Triggers[0].Body
	if strings.Contains(body, "ON\nT") || !strings.Contains(body, "ON\nU") {
		t.Fatalf("trigger target = %q", body)
	}
}

func TestDisabledConstraintCreationIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER, CONSTRAINT CK CHECK(C > 0) DISABLE)`); err == nil {
		t.Fatal("disabled constraint installed as enabled")
	}
}

func TestDisabledLiteralDoesNotRejectConstraint(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (STATUS VARCHAR2(20), CONSTRAINT CK CHECK (STATUS <> 'DISABLE'))`); err != nil {
		t.Fatal(err)
	}
}

func TestCascadeKeyDropRemovesSelfForeignKey(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER, PARENT_ID NUMBER, CONSTRAINT PK PRIMARY KEY(ID), CONSTRAINT FK FOREIGN KEY(PARENT_ID) REFERENCES T(ID)); ALTER TABLE T DROP CONSTRAINT PK CASCADE`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["T"].ForeignKeys) != 0 {
		t.Fatal("self foreign key survived cascading key drop")
	}
}

func TestRenameDisabledTriggerUpdatesMarker(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (ID NUMBER); CREATE TRIGGER TR BEFORE INSERT ON T BEGIN NULL; END; ALTER TRIGGER TR DISABLE; ALTER TRIGGER TR RENAME TO TR2`)
	if err != nil {
		t.Fatal(err)
	}
	body := tableByName(c.Metadata().Schemas[0])["T"].Triggers[0].Body
	if !strings.Contains(body, `ALTER TRIGGER "TR2" DISABLE;`) {
		t.Fatalf("disable marker = %q", body)
	}
}

func TestDeferrableCheckIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER, CONSTRAINT CK CHECK(C > 0) DEFERRABLE INITIALLY DEFERRED)`); err == nil {
		t.Fatal("deferrable check accepted")
	}
}

func TestTableWideTriggerDisableIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (ID NUMBER); CREATE TRIGGER TR BEFORE INSERT ON T BEGIN NULL; END; ALTER TABLE T DISABLE ALL TRIGGERS`); err == nil {
		t.Fatal("table-wide trigger disable accepted")
	}
}

func TestViewDependencyUsesEffectiveSchema(t *testing.T) {
	to, err := LoadSQL(`CREATE TABLE Z (C NUMBER); CREATE VIEW A AS SELECT C FROM Z; CREATE TABLE APP.U (C NUMBER)`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, New(), to)
	table := strings.Index(sql, `CREATE TABLE "APP"."Z"`)
	view := strings.Index(sql, `CREATE VIEW "APP"."A"`)
	if table < 0 || view < 0 || table > view {
		t.Fatalf("view precedes its table:\n%s", sql)
	}
}

func TestMaterializedViewSkipsSystemGeneratedIndex(t *testing.T) {
	// Oracle names a materialized view's container indexes
	// SYS_C_SNAP$_<object id>, and the id differs between two snapshots of the
	// same view. A definition must not write one (the view creates it again),
	// and a diff must not treat it as a change.
	view := func(objectID string) *metadatapb.MaterializedViewMetadata {
		return &metadatapb.MaterializedViewMetadata{
			Name:       "MV",
			Definition: "SELECT ID, NAME FROM T\n",
			Indexes: []*metadatapb.IndexMetadata{
				{Name: "SYS_C_SNAP$_" + objectID + "MV", Expressions: []string{"ID"}, Unique: true},
				{Name: "MV_NAME_IDX", Expressions: []string{"NAME"}},
			},
		}
	}
	from := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", MaterializedViews: []*metadatapb.MaterializedViewMetadata{view("9000013")}})
	to := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", MaterializedViews: []*metadatapb.MaterializedViewMetadata{view("9000099")}})

	definition, err := DatabaseDefinition(to.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(definition, "SYS_C_SNAP$") {
		t.Fatalf("a definition wrote the view's system-generated index:\n%s", definition)
	}
	if !strings.Contains(definition, `"MV_NAME_IDX"`) {
		t.Fatalf("a definition dropped the view's own index:\n%s", definition)
	}

	if !Diff(from, to).IsEmpty() {
		t.Fatalf("a generated index's object id changed the diff:\n%s", migrationFor(t, from, to))
	}

	sql := migrationFor(t, New(), to)
	if strings.Contains(sql, "SYS_C_SNAP$") {
		t.Fatalf("a migration wrote the view's system-generated index:\n%s", sql)
	}
	if !strings.Contains(sql, `"MV_NAME_IDX"`) {
		t.Fatalf("a migration dropped the view's own index:\n%s", sql)
	}
}

func TestCascadeUnnamedPrimaryKeyDropKeepsForeignKeyOnOtherKey(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE P (ID NUMBER PRIMARY KEY, CODE NUMBER CONSTRAINT P_CODE_UK UNIQUE);
CREATE TABLE C (CODE NUMBER, CONSTRAINT C_FK FOREIGN KEY (CODE) REFERENCES P (CODE));
ALTER TABLE P DROP PRIMARY KEY CASCADE`)
	if err != nil {
		t.Fatal(err)
	}
	schema := c.Metadata().Schemas[0]
	if got := tableByName(schema)["C"].ForeignKeys; len(got) != 1 {
		t.Fatalf("the cascade removed a foreign key on another key: %v", got)
	}
	p := tableByName(schema)["P"]
	if len(p.Indexes) != 1 || !p.Indexes[0].Unique {
		t.Fatalf("the unique key on CODE was removed with the primary key: %v", p.Indexes)
	}
}

func TestDefinitionOrdersViewsThroughEffectiveSchema(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE Z (C NUMBER);
CREATE VIEW B AS SELECT C FROM Z;
CREATE VIEW A AS SELECT C FROM B;
CREATE TABLE APP.U (C NUMBER)`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(c.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	inner := strings.Index(definition, `CREATE VIEW "B"`)
	outer := strings.Index(definition, `CREATE VIEW "A"`)
	if inner < 0 || outer < 0 || inner > outer {
		t.Fatalf("a view precedes the view it reads:\n%s", definition)
	}
}

func TestCreateSchemaAuthorizationMustMatch(t *testing.T) {
	_, err := LoadSQL(`CREATE TABLE APP.T (C NUMBER); CREATE SCHEMA AUTHORIZATION OTHER CREATE TABLE U (C NUMBER)`)
	if err == nil || !strings.Contains(err.Error(), `statement names schema "OTHER"`) {
		t.Fatalf("CREATE SCHEMA AUTHORIZATION OTHER: %v, want the owner mismatch", err)
	}
}

func TestStandaloneRenameMovesSequencesAndViews(t *testing.T) {
	c, err := LoadSQL(`CREATE SEQUENCE S;
RENAME S TO S2;
DROP SEQUENCE S2;
CREATE VIEW V AS SELECT 1 AS X FROM DUAL;
RENAME V TO V2`)
	if err != nil {
		t.Fatal(err)
	}
	schema := c.Metadata().Schemas[0]
	if len(schema.Sequences) != 0 {
		t.Fatalf("a dropped, renamed sequence stayed: %v", schema.Sequences)
	}
	if len(schema.Views) != 1 || schema.Views[0].Name != "V2" {
		t.Fatalf("views = %v, want V2", schema.Views)
	}
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE INDEX I ON T (C); RENAME I TO J`); err == nil {
		t.Fatal("renaming an index was silently ignored")
	}
}

func TestQuotedFunctionReferenceOrdersBeforeTable(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{
		Tables:    []*metadatapb.TableMetadata{{Name: "A", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: `"Z"(C)`}}}}},
		Functions: []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z(P NUMBER) RETURN NUMBER IS BEGIN RETURN P; END;"}},
	}}}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, "CREATE OR REPLACE FUNCTION Z")
	table := strings.Index(definition, `CREATE TABLE "A"`)
	if function < 0 || table < 0 || function > table {
		t.Fatalf("quoted function dependency missing:\n%s", definition)
	}
}

func TestMaterializedViewKeepsScheduleWithoutMethod(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER); CREATE MATERIALIZED VIEW MV REFRESH START WITH SYSDATE NEXT SYSDATE + 1 AS SELECT C FROM T`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, "START WITH SYSDATE NEXT SYSDATE + 1") {
		t.Fatalf("schedule missing:\n%s", definition)
	}
}

func TestRenameViewRewritesTriggerTarget(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE VIEW V AS SELECT C FROM T; CREATE TRIGGER TR INSTEAD OF INSERT ON V BEGIN NULL; END; RENAME V TO U`)
	if err != nil {
		t.Fatal(err)
	}
	var view *metadatapb.ViewMetadata
	for _, candidate := range c.Metadata().Schemas[0].Views {
		if candidate.Name == "U" {
			view = candidate
		}
	}
	if view == nil || len(view.Triggers) != 1 || !strings.Contains(view.Triggers[0].Body, " ON U ") {
		t.Fatalf("view trigger = %+v", view)
	}
}

func TestUsingExistingIndexConstraintIsRejected(t *testing.T) {
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE UNIQUE INDEX I ON T(C); ALTER TABLE T ADD CONSTRAINT UK UNIQUE(C) USING INDEX I`); err == nil {
		t.Fatal("USING INDEX constraint accepted")
	}
}

func TestQuotedFunctionDependencyOrder(t *testing.T) {
	c, err := LoadSQL(`CREATE FUNCTION "Z" (C NUMBER) RETURN NUMBER AS BEGIN RETURN C; END;
CREATE TABLE A (C NUMBER, V GENERATED ALWAYS AS ("Z"(C)) VIRTUAL)`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(c.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, `FUNCTION "Z"`)
	table := strings.Index(definition, `CREATE TABLE "A"`)
	if function < 0 || table < 0 || function > table {
		t.Fatalf("a table precedes the function its virtual column calls:\n%s", definition)
	}
	sql := migrationFor(t, New(), c)
	function, table = strings.Index(sql, `FUNCTION "Z"`), strings.Index(sql, `CREATE TABLE "A"`)
	if function < 0 || table < 0 || function > table {
		t.Fatalf("a migration creates the table before its function:\n%s", sql)
	}
}

func TestMaterializedViewKeepsScheduleWithoutRefreshMode(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER);
CREATE MATERIALIZED VIEW MV REFRESH START WITH SYSDATE NEXT SYSDATE + 1 AS SELECT C FROM T`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(c.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, "REFRESH START WITH SYSDATE NEXT SYSDATE + 1") {
		t.Fatalf("the definition dropped the refresh schedule:\n%s", definition)
	}

	without := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name:              "APP",
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "MV", Definition: "SELECT C FROM T\n"}},
	})
	if Diff(without, c).IsEmpty() {
		t.Fatal("a refresh schedule added to a materialized view produced an empty diff")
	}
	if sql := migrationFor(t, without, c); !strings.Contains(sql, "REFRESH START WITH") {
		t.Fatalf("the migration dropped the refresh schedule:\n%s", sql)
	}
}

func TestRenameMovesAttachedTriggers(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER);
CREATE VIEW V AS SELECT C FROM T;
CREATE TRIGGER TR INSTEAD OF INSERT ON V BEGIN NULL; END;
RENAME V TO U`)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(c.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `CREATE VIEW "U"`) {
		t.Fatalf("the view did not move:\n%s", definition)
	}
	if strings.Contains(definition, " ON V ") {
		t.Fatalf("the view's trigger still targets the old name:\n%s", definition)
	}
	if !strings.Contains(definition, ` ON U `) {
		t.Fatalf("the view's trigger did not follow the rename:\n%s", definition)
	}

	extractor := &oracleMetadataExtractor{
		materializedViews: map[string]*metadatapb.MaterializedViewMetadata{
			"MV": {Name: "MV", Triggers: []*metadatapb.TriggerMetadata{{Body: `CREATE OR REPLACE TRIGGER "TR" BEFORE INSERT ON MV FOR EACH ROW BEGIN NULL; END;`}}},
		},
		currentSchema: "APP",
	}
	if !extractor.renameMaterializedView("MV", "MV2") {
		t.Fatal("the materialized view was not renamed")
	}
	if body := extractor.materializedViews["MV2"].Triggers[0].Body; !strings.Contains(body, " ON MV2 ") {
		t.Fatalf("the materialized view's trigger did not follow the rename: %s", body)
	}
}

func TestConstraintUsingExistingIndexIsRejected(t *testing.T) {
	_, err := LoadSQL(`CREATE TABLE T (C NUMBER);
CREATE UNIQUE INDEX I ON T (C);
ALTER TABLE T ADD CONSTRAINT UK UNIQUE (C) USING INDEX I`)
	if err == nil || !strings.Contains(err.Error(), `backed by the existing index "I"`) {
		t.Fatalf("a constraint backed by an existing index: %v, want the unsupported error", err)
	}
	// The physical-properties form names no index and stays a documented gap.
	if _, err := LoadSQL(`CREATE TABLE T (C NUMBER, CONSTRAINT UK UNIQUE (C) USING INDEX TABLESPACE USERS)`); err != nil {
		t.Fatalf("USING INDEX TABLESPACE: %v", err)
	}
}

func TestCascadeCompositeKeyMatchesCompleteColumns(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE P (A NUMBER, B NUMBER, C NUMBER, CONSTRAINT U_AB UNIQUE(A,B), CONSTRAINT U_AC UNIQUE(A,C)); CREATE TABLE CH (A NUMBER, C NUMBER, CONSTRAINT CH_FK FOREIGN KEY(A,C) REFERENCES P(A,C)); ALTER TABLE P DROP CONSTRAINT U_AB CASCADE`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tableByName(c.Metadata().Schemas[0])["CH"].ForeignKeys) != 1 {
		t.Fatal("foreign key for the retained composite key was removed")
	}
}

func TestRenameQualifiedTrigger(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1}}, Triggers: []*metadatapb.TriggerMetadata{{Name: "TR", Body: "CREATE OR REPLACE TRIGGER APP.TR BEFORE INSERT ON APP.T BEGIN NULL; END;"}}}}})
	if _, err := c.Exec(`ALTER TRIGGER APP.TR RENAME TO TR2`, nil); err != nil {
		t.Fatal(err)
	}
	body := tableByName(c.Metadata().Schemas[0])["T"].Triggers[0].Body
	if !strings.Contains(body, "TRIGGER APP.TR2") {
		t.Fatalf("trigger body = %q", body)
	}
}

func TestMaterializedViewFunctionIndexAfterFunction(t *testing.T) {
	old := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}}}})
	new := catalogFromSchema(&metadatapb.SchemaMetadata{Name: "APP",
		Tables:            []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}}},
		Functions:         []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z(P NUMBER) RETURN NUMBER IS BEGIN RETURN P; END;"}},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "MV", Definition: "SELECT C FROM T", Indexes: []*metadatapb.IndexMetadata{{Name: "I", Type: "FUNCTION-BASED NORMAL", Expressions: []string{"Z(C)"}, Visible: true}}}},
	})
	sql := migrationFor(t, old, new)
	function := strings.Index(sql, "CREATE OR REPLACE FUNCTION Z")
	index := strings.Index(sql, `CREATE INDEX "APP"."I"`)
	if function < 0 || index < 0 || function > index {
		t.Fatalf("materialized-view index precedes its function:\n%s", sql)
	}
}

func TestIdentityModeChangeDoesNotRestart(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_ALWAYS, IdentitySeed: 1, IdentityIncrement: 1}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, IsIdentity: true, IdentityGeneration: metadatapb.ColumnMetadata_BY_DEFAULT, IdentitySeed: 1, IdentityIncrement: 1}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if strings.Contains(sql, "START WITH") {
		t.Fatalf("identity mode change restarts the sequence:\n%s", sql)
	}
}

func TestVirtualColumnDetachedBeforeSourceTypeChange(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "BINARY_DOUBLE", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	drop := strings.Index(sql, `DROP COLUMN "V"`)
	modify := strings.Index(sql, `MODIFY "C" BINARY_DOUBLE`)
	add := strings.LastIndex(sql, `ADD "V"`)
	if drop < 0 || modify < 0 || add < 0 || !(drop < modify && modify < add) {
		t.Fatalf("virtual column not detached around source type change:\n%s", sql)
	}
}

func TestDefinitionRendersViewComments(t *testing.T) {
	to := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name: "APP",
		Views: []*metadatapb.ViewMetadata{{
			Name:       "V",
			Definition: "SELECT C FROM T\n",
			Comment:    "it's a view",
		}},
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
			Name:       "MV",
			Definition: "SELECT C FROM T\n",
			Comment:    "a materialized view",
		}},
	})
	definition, err := DatabaseDefinition(to.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `COMMENT ON TABLE "V" IS 'it''s a view';`) {
		t.Fatalf("the view's comment is missing or unescaped:\n%s", definition)
	}
	if !strings.Contains(definition, `COMMENT ON MATERIALIZED VIEW "MV" IS 'a materialized view';`) {
		t.Fatalf("the materialized view's comment is missing:\n%s", definition)
	}
}

func TestSharedVirtualDependentsDropOnce(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, {Name: "V1", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}, {Name: "V2", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+2"}}}, Indexes: []*metadatapb.IndexMetadata{{Name: "I", Expressions: []string{"V1", "V2"}, Visible: true}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "BINARY_DOUBLE", Position: 1}, {Name: "V1", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}, {Name: "V2", Type: "NUMBER", Position: 3, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+2"}}}, Indexes: oldTable.Indexes}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if strings.Count(sql, `DROP INDEX "I"`) != 1 || strings.Count(sql, `CREATE INDEX "I"`) != 1 {
		t.Fatalf("shared index handled more than once:\n%s", sql)
	}
	if strings.Index(sql, `ADD "V2"`) > strings.Index(sql, `CREATE INDEX "I"`) {
		t.Fatalf("index recreated before all virtual columns:\n%s", sql)
	}
}

func TestVirtualCommentChangeStillDetachesForTypeChange(t *testing.T) {
	oldTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}}}
	newTable := &metadatapb.TableMetadata{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "BINARY_DOUBLE", Position: 1}, {Name: "V", Type: "NUMBER", Position: 2, Comment: "changed", Generation: &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_VIRTUAL, Expression: "C+1"}}}}
	sql := migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{oldTable}}), catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{newTable}}))
	if !strings.Contains(sql, `DROP COLUMN "V"`) || !strings.Contains(sql, `MODIFY "C" BINARY_DOUBLE`) {
		t.Fatalf("comment-only virtual column not detached:\n%s", sql)
	}
}

func TestRecreatedMaterializedViewKeepsFunctionIndex(t *testing.T) {
	oldView := &metadatapb.MaterializedViewMetadata{Name: "MV", Definition: "SELECT C FROM T"}
	newView := &metadatapb.MaterializedViewMetadata{Name: "MV", Definition: "SELECT C+1 AS C FROM T", Indexes: []*metadatapb.IndexMetadata{{Name: "I", Type: "FUNCTION-BASED NORMAL", Expressions: []string{"Z(C)"}, Visible: true}}}
	meta := &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{Name: "APP", Tables: []*metadatapb.TableMetadata{{Name: "T", Columns: []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}}}, Functions: []*metadatapb.FunctionMetadata{{Name: "Z", Definition: "FUNCTION Z(P NUMBER) RETURN NUMBER IS BEGIN RETURN P; END;"}}, MaterializedViews: []*metadatapb.MaterializedViewMetadata{oldView}}}}
	to := clone(meta)
	to.Schemas[0].MaterializedViews = []*metadatapb.MaterializedViewMetadata{newView}
	sql := migrationFor(t, &Catalog{meta: meta}, &Catalog{meta: to})
	if !strings.Contains(sql, `CREATE INDEX "APP"."I"`) {
		t.Fatalf("recreated materialized view lost its index:\n%s", sql)
	}
}

func TestDropUniqueKeepIndex(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE T (C NUMBER CONSTRAINT UK UNIQUE); ALTER TABLE T DROP UNIQUE(C) KEEP INDEX`)
	if err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	found := false
	for _, index := range table.Indexes {
		if index.Name == "UK" && index.Unique && !index.IsConstraint {
			found = true
		}
	}
	if !found {
		t.Fatalf("kept index missing: %+v", table.Indexes)
	}
}

func TestVirtualColumnDependentKeyDetachesIncomingForeignKey(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE T (C NUMBER, V AS (CAST(C AS NUMBER)), CONSTRAINT T_V_UK UNIQUE (V));
CREATE TABLE CH (V NUMBER, CONSTRAINT CH_FK FOREIGN KEY (V) REFERENCES T (V))`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE T (C BINARY_DOUBLE, V AS (CAST(C AS NUMBER)), CONSTRAINT T_V_UK UNIQUE (V));
CREATE TABLE CH (V NUMBER, CONSTRAINT CH_FK FOREIGN KEY (V) REFERENCES T (V))`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, from, to)
	drop := strings.Index(sql, `ALTER TABLE "CH" DROP CONSTRAINT "CH_FK"`)
	key := strings.Index(sql, `DROP CONSTRAINT "T_V_UK"`)
	add := strings.Index(sql, `ADD CONSTRAINT "CH_FK"`)
	if drop < 0 || key < 0 || add < 0 || !(drop < key && key < add) {
		t.Fatalf("the child foreign key does not detach around the key's unique constraint (drop=%d key=%d add=%d):\n%s", drop, key, add, sql)
	}
}

func TestDefinitionDefersMaterializedViewFunctionIndex(t *testing.T) {
	to := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name: "APP",
		MaterializedViews: []*metadatapb.MaterializedViewMetadata{{
			Name:       "A",
			Definition: "SELECT C FROM T\n",
			Indexes: []*metadatapb.IndexMetadata{{
				Name:        "I",
				Type:        "FUNCTION-BASED NORMAL",
				Expressions: []string{"Z(C)"},
			}},
		}},
		Functions: []*metadatapb.FunctionMetadata{{
			Name:       "Z",
			Definition: "CREATE OR REPLACE FUNCTION \"Z\" (C NUMBER) RETURN NUMBER AS BEGIN RETURN C; END;",
		}},
	})
	definition, err := DatabaseDefinition(to.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	function := strings.Index(definition, `FUNCTION "Z"`)
	index := strings.Index(definition, `CREATE INDEX "I"`)
	if function < 0 || index < 0 || function > index {
		t.Fatalf("the view's function-based index precedes its function (function=%d index=%d):\n%s", function, index, definition)
	}
}

func TestTriggerDisableClauseSurvivesWhitespace(t *testing.T) {
	c, err := LoadSQL("CREATE TABLE T (C NUMBER);\nCREATE TRIGGER TR BEFORE INSERT ON T\nDISABLE\nBEGIN NULL; END;")
	if err != nil {
		t.Fatal(err)
	}
	definition, err := DatabaseDefinition(c.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, "DISABLE") {
		t.Fatalf("a disabled trigger's declaration lost its clause:\n%s", definition)
	}
	if _, err := c.Exec("ALTER TRIGGER TR ENABLE", nil); err != nil {
		t.Fatal(err)
	}
	definition, err = DatabaseDefinition(c.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(definition, "DISABLE") {
		t.Fatalf("an enabled trigger kept the declaration's DISABLE clause:\n%s", definition)
	}
}

func TestVirtualColumnExpressionChangeDetaches(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE T (C NUMBER, V AS (C + 1))`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE T (C BINARY_DOUBLE, V AS (C + 2))`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, from, to)
	drop := strings.Index(sql, `DROP COLUMN "V"`)
	modify := strings.Index(sql, `MODIFY "C" BINARY_DOUBLE`)
	add := strings.Index(sql, `GENERATED ALWAYS AS (C + 2) VIRTUAL`)
	if drop < 0 || modify < 0 || add < 0 || !(drop < modify && modify < add) {
		t.Fatalf("a changed virtual column does not detach around its source (drop=%d modify=%d add=%d):\n%s", drop, modify, add, sql)
	}
	if mod := strings.Index(sql, `MODIFY "V"`); mod >= 0 {
		t.Fatalf("the virtual column was modified instead of recreated:\n%s", sql)
	}
}

func TestNewForeignKeyWaitsForDetachedVirtualKey(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE P (C NUMBER, V AS (CAST(C AS NUMBER)), CONSTRAINT P_V_UK UNIQUE (V))`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE P (C BINARY_DOUBLE, V AS (CAST(C AS NUMBER)), CONSTRAINT P_V_UK UNIQUE (V));
CREATE TABLE CH (V NUMBER, CONSTRAINT CH_FK FOREIGN KEY (V) REFERENCES P (V))`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, from, to)
	key := strings.Index(sql, `ADD CONSTRAINT "P_V_UK" UNIQUE`)
	fk := strings.Index(sql, `ADD CONSTRAINT "CH_FK"`)
	if key < 0 || fk < 0 || key > fk {
		t.Fatalf("the new foreign key precedes the key it references (key=%d fk=%d):\n%s", key, fk, sql)
	}
}

func TestDetachedVirtualColumnKeepsComment(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE T (C NUMBER, V AS (C + 1));
COMMENT ON COLUMN T.V IS 'v'`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE T (C BINARY_DOUBLE, V AS (C + 1));
COMMENT ON COLUMN T.V IS 'v'`)
	if err != nil {
		t.Fatal(err)
	}
	if sql := migrationFor(t, from, to); !strings.Contains(sql, `COMMENT ON COLUMN "T"."V" IS 'v'`) {
		t.Fatalf("a recreated virtual column lost its comment:\n%s", sql)
	}
}

func TestFunctionReplacementRebuildsItsIndexes(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE T (C NUMBER);
CREATE FUNCTION F (X NUMBER) RETURN NUMBER AS BEGIN RETURN X + 1; END;
CREATE INDEX I ON T (F(C))`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE T (C NUMBER);
CREATE FUNCTION F (X NUMBER) RETURN NUMBER AS BEGIN RETURN X + 2; END;
CREATE INDEX I ON T (F(C))`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, from, to)
	function := strings.Index(sql, `RETURN X + 2`)
	rebuild := strings.Index(sql, `ALTER INDEX "I" REBUILD`)
	if function < 0 || rebuild < 0 || function > rebuild {
		t.Fatalf("a replaced function's index is not rebuilt after it (function=%d rebuild=%d):\n%s", function, rebuild, sql)
	}
}

func TestProcedureCreatePrefixAcrossWhitespace(t *testing.T) {
	from := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name:       "APP",
		Procedures: []*metadatapb.ProcedureMetadata{{Name: "P", Definition: "CREATE\nPROCEDURE P IS BEGIN NULL; END;"}},
	})
	to := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name:       "APP",
		Procedures: []*metadatapb.ProcedureMetadata{{Name: "P", Definition: "CREATE\nPROCEDURE P IS BEGIN NULL; NULL; END;"}},
	})
	sql := migrationFor(t, from, to)
	if strings.Count(strings.ToUpper(sql), "CREATE OR REPLACE") != 1 || strings.Contains(strings.ToUpper(sql), "CREATE\nPROCEDURE") {
		t.Fatalf("a procedure whose declaration puts a newline after CREATE:\n%s", sql)
	}
}

func TestFunctionIndexDetachesBeforeOrdinaryTypeChange(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE T (C NUMBER); CREATE INDEX I ON T (ABS(C))`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE T (C BINARY_DOUBLE); CREATE INDEX I ON T (ABS(C))`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, from, to)
	drop := strings.Index(sql, `DROP INDEX "I"`)
	modify := strings.Index(sql, `MODIFY "C" BINARY_DOUBLE`)
	create := strings.Index(sql, `CREATE INDEX "I"`)
	if drop < 0 || modify < 0 || create < 0 || !(drop < modify && modify < create) {
		t.Fatalf("a function-based index does not detach around a datatype change (drop=%d modify=%d create=%d):\n%s", drop, modify, create, sql)
	}
}

func TestModifyConstraintStateIsRejected(t *testing.T) {
	_, err := LoadSQL(`CREATE TABLE T (C NUMBER, CONSTRAINT UK UNIQUE (C)); ALTER TABLE T MODIFY CONSTRAINT UK DISABLE`)
	if err == nil || !strings.Contains(err.Error(), "constraint enable/disable is not supported") {
		t.Fatalf("MODIFY CONSTRAINT UK DISABLE: %v, want the unsupported-state error", err)
	}
}

func TestVirtualColumnDetachesWhenSourceLeavesTheExpression(t *testing.T) {
	from, err := LoadSQL(`CREATE TABLE T (C NUMBER, D NUMBER, V AS (C + 1))`)
	if err != nil {
		t.Fatal(err)
	}
	to, err := LoadSQL(`CREATE TABLE T (C BINARY_DOUBLE, D NUMBER, V AS (D + 1))`)
	if err != nil {
		t.Fatal(err)
	}
	sql := migrationFor(t, from, to)
	drop := strings.Index(sql, `DROP COLUMN "V"`)
	modify := strings.Index(sql, `MODIFY "C" BINARY_DOUBLE`)
	if drop < 0 || modify < 0 || drop > modify {
		t.Fatalf("the old virtual expression's source change precedes its detach (drop=%d modify=%d):\n%s", drop, modify, sql)
	}
}

func TestCreateSequenceIfNotExistsKeepsTheSequence(t *testing.T) {
	c, err := LoadSQL(`CREATE SEQUENCE S START WITH 100 INCREMENT BY 10; CREATE SEQUENCE IF NOT EXISTS S`)
	if err != nil {
		t.Fatal(err)
	}
	sequences := c.Metadata().Schemas[0].Sequences
	if len(sequences) != 1 || sequences[0].Start != "100" || sequences[0].Increment != "10" {
		t.Fatalf("CREATE SEQUENCE IF NOT EXISTS replaced the sequence: %v", sequences)
	}
}

func TestCycleFallbackWritesFunctionBeforeView(t *testing.T) {
	table := func(name, other string) *metadatapb.TableMetadata {
		return &metadatapb.TableMetadata{
			Name:    name,
			Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true}},
			ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
				Name:              name + "_FK",
				Columns:           []string{"ID"},
				ReferencedTable:   other,
				ReferencedColumns: []string{"ID"},
			}},
		}
	}
	schema := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name:   "APP",
		Tables: []*metadatapb.TableMetadata{table("A", "B"), table("B", "A")},
		Functions: []*metadatapb.FunctionMetadata{{
			Name:       "F",
			Definition: "CREATE OR REPLACE FUNCTION \"F\" (X NUMBER) RETURN NUMBER AS BEGIN RETURN X; END;",
		}},
		Views: []*metadatapb.ViewMetadata{{
			Name:       "V",
			Definition: "SELECT F(ID) FROM A",
		}},
	})
	definition, err := DatabaseDefinition(schema.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	function, view := strings.Index(definition, `FUNCTION "F"`), strings.Index(definition, `CREATE VIEW "V"`)
	if function < 0 || view < 0 || function > view {
		t.Fatalf("the cycle fallback writes the view before the function it calls (function=%d view=%d):\n%s", function, view, definition)
	}
	sql := migrationFor(t, New(), schema)
	function, view = strings.Index(sql, `FUNCTION "F"`), strings.Index(sql, `CREATE VIEW`)
	if function < 0 || view < 0 || function > view {
		t.Fatalf("the migration's cycle fallback writes the view before the function it calls (function=%d view=%d):\n%s", function, view, sql)
	}
}

func TestSyncedQuotedKeysMatchExtractedOnes(t *testing.T) {
	// The driver's sync quotes a key column to keep a case-sensitive name; an
	// extracted snapshot states it raw. The same column is not a change, and it
	// is rendered with one pair of quotes either way.
	quoted := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name: "APP",
		Tables: []*metadatapb.TableMetadata{{
			Name: "T",
			Columns: []*metadatapb.ColumnMetadata{
				{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true},
				{Name: "REF", Type: "NUMBER", Position: 2, Nullable: true},
			},
			Indexes: []*metadatapb.IndexMetadata{
				{Name: "T_PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{`"ID"`}},
				{Name: "T_REF_IDX", Type: "NORMAL", Expressions: []string{`"REF"`}, Descending: []bool{false}},
			},
			ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
				Name:              "T_FK",
				Columns:           []string{`"REF"`},
				ReferencedTable:   "T",
				ReferencedColumns: []string{`"ID"`},
			}},
		}},
	})
	definition, err := DatabaseDefinition(quoted.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`PRIMARY KEY ("ID")`, `CREATE INDEX "T_REF_IDX" ON "T" ("REF" ASC)`, `FOREIGN KEY ("REF") REFERENCES "T" ("ID")`} {
		if !strings.Contains(definition, want) {
			t.Fatalf("the definition does not state %s:\n%s", want, definition)
		}
	}
	if strings.Contains(definition, `""`) {
		t.Fatalf("a quoted key column was quoted twice:\n%s", definition)
	}

	raw := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name: "APP",
		Tables: []*metadatapb.TableMetadata{{
			Name: "T",
			Columns: []*metadatapb.ColumnMetadata{
				{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true},
				{Name: "REF", Type: "NUMBER", Position: 2, Nullable: true},
			},
			Indexes: []*metadatapb.IndexMetadata{
				{Name: "T_PK", Primary: true, Unique: true, IsConstraint: true, Expressions: []string{"ID"}},
				{Name: "T_REF_IDX", Type: "NORMAL", Expressions: []string{"REF"}, Descending: []bool{false}},
			},
			ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
				Name:              "T_FK",
				Columns:           []string{"REF"},
				ReferencedTable:   "T",
				ReferencedColumns: []string{"ID"},
			}},
		}},
	})
	if !Diff(raw, quoted).IsEmpty() {
		t.Fatalf("a quoted key column read as a change:\n%s", migrationFor(t, raw, quoted))
	}
}

// TestMigrationWritesOnlyCommentCategoriesOracleTakes pins the comment
// statements a migration writes: COMMENT ON TABLE, which is how a view is
// commented, and COMMENT ON MATERIALIZED VIEW. Oracle names no other category:
// COMMENT ON SEQUENCE, FUNCTION and INDEX are ORA-32594, "invalid object
// category", so a metadata comment on one of those has no statement to render
// and must not produce one.
func TestMigrationWritesOnlyCommentCategoriesOracleTakes(t *testing.T) {
	columns := []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}}
	table := func(comment string) *metadatapb.TableMetadata {
		return &metadatapb.TableMetadata{
			Name: "T", Columns: columns,
			// Function-based, so the index is created by a statement of its own.
			Indexes: []*metadatapb.IndexMetadata{{Name: "T_FBI", Type: "FUNCTION-BASED NORMAL", Expressions: []string{"UPPER(C)"}, Visible: true, Comment: comment}},
		}
	}
	sequence := func(comment string) *metadatapb.SequenceMetadata {
		return &metadatapb.SequenceMetadata{Name: "S", Increment: "1", CacheSize: "20", Comment: comment}
	}
	function := func(comment string) *metadatapb.FunctionMetadata {
		return &metadatapb.FunctionMetadata{Name: "F", Definition: "FUNCTION F RETURN NUMBER IS BEGIN RETURN 1; END;", Comment: comment}
	}
	view := func(comment string) *metadatapb.ViewMetadata {
		return &metadatapb.ViewMetadata{Name: "V", Definition: "SELECT 1 AS X FROM DUAL", Comment: comment}
	}
	materializedView := func(comment string) *metadatapb.MaterializedViewMetadata {
		return &metadatapb.MaterializedViewMetadata{Name: "MV", Definition: "SELECT 1 AS X FROM DUAL", Comment: comment}
	}
	emitted := []string{"COMMENT ON TABLE", "COMMENT ON MATERIALIZED VIEW"}
	check := func(t *testing.T, sql string) {
		t.Helper()
		for _, statement := range emitted {
			if !strings.Contains(sql, statement) {
				t.Errorf("missing %q:\n%s", statement, sql)
			}
		}
		for _, statement := range []string{"COMMENT ON INDEX", "COMMENT ON SEQUENCE", "COMMENT ON FUNCTION"} {
			if strings.Contains(sql, statement) {
				t.Errorf("wrote %q, which Oracle refuses with ORA-32594:\n%s", statement, sql)
			}
		}
	}

	commented := func(comment string) *Catalog {
		return catalogFromSchema(&metadatapb.SchemaMetadata{
			Tables:            []*metadatapb.TableMetadata{table(comment)},
			Sequences:         []*metadatapb.SequenceMetadata{sequence(comment)},
			Functions:         []*metadatapb.FunctionMetadata{function(comment)},
			Views:             []*metadatapb.ViewMetadata{view(comment)},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{materializedView(comment)},
		})
	}
	// A comment that arrives on an object that exists, and one that arrives
	// with the object itself.
	check(t, migrationFor(t, commented(""), commented("a comment")))
	check(t, migrationFor(t, catalogFromSchema(&metadatapb.SchemaMetadata{}), commented("a comment")))
}

func TestQuotedKeyWithSpacesIsAChange(t *testing.T) {
	index := func(expression string) *Catalog {
		return catalogFromSchema(&metadatapb.SchemaMetadata{
			Name: "APP",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "T",
				Columns: []*metadatapb.ColumnMetadata{{Name: "A", Type: "NUMBER", Position: 1, Nullable: true}},
				Indexes: []*metadatapb.IndexMetadata{{
					Name:        "T_IDX",
					Type:        "NORMAL",
					Expressions: []string{expression},
					Descending:  []bool{false},
				}},
			}},
		})
	}
	if Diff(index(`"A  B"`), index(`"A B"`)).IsEmpty() {
		t.Fatal("a quoted key with repeated spaces compared equal to one with a single space")
	}
	if !Diff(index(`"A"`), index(`A`)).IsEmpty() {
		t.Fatal("a quoted key compared unequal to the same raw key")
	}
}

func TestReusedIndexConstraintRendersUsingIndex(t *testing.T) {
	to := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name: "APP",
		Tables: []*metadatapb.TableMetadata{{
			Name: "T",
			Columns: []*metadatapb.ColumnMetadata{
				{Name: "ID", Type: "NUMBER", Position: 1, Nullable: false},
				{Name: "CODE", Type: "VARCHAR2(20)", Position: 2, Nullable: true},
			},
			Indexes: []*metadatapb.IndexMetadata{
				{
					Name:         "T_CODE_IDX",
					Type:         "NORMAL",
					Expressions:  []string{`"CODE"`},
					Descending:   []bool{false},
					Visible:      true,
					IsConstraint: true,
					Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "T_CODE_UK", Columns: []string{`"CODE"`}}},
				},
				{
					Name:         "T_ID_IDX",
					Type:         "NORMAL",
					Primary:      true,
					Expressions:  []string{`"ID"`},
					Descending:   []bool{false},
					Visible:      true,
					IsConstraint: true,
					Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "T_PK", Columns: []string{`"ID"`}, Primary: true}},
				},
			},
		}},
	})
	definition, err := DatabaseDefinition(to.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(definition, `CONSTRAINT "T_CODE_UK",`) {
		t.Fatalf("a reused-index constraint was rendered inline without its type:\n%s", definition)
	}
	index := strings.Index(definition, `CREATE INDEX "T_CODE_IDX" ON "T" ("CODE" ASC)`)
	constraint := strings.Index(definition, `ALTER TABLE "T" ADD CONSTRAINT "T_CODE_UK" UNIQUE ("CODE") USING INDEX "T_CODE_IDX"`)
	if index < 0 || constraint < 0 || index > constraint {
		t.Fatalf("the constraint does not follow the index it reuses (index=%d constraint=%d):\n%s", index, constraint, definition)
	}
	if !strings.Contains(definition, `ADD CONSTRAINT "T_PK" PRIMARY KEY ("ID") USING INDEX "T_ID_IDX"`) {
		t.Fatalf("the primary key does not name its index:\n%s", definition)
	}

	sql := migrationFor(t, New(), to)
	index = strings.Index(sql, `CREATE INDEX "APP"."T_CODE_IDX"`)
	constraint = strings.Index(sql, `ADD CONSTRAINT "T_CODE_UK" UNIQUE ("CODE") USING INDEX "T_CODE_IDX"`)
	if index < 0 || constraint < 0 || index > constraint {
		t.Fatalf("the migration does not create the index before its constraint (index=%d constraint=%d):\n%s", index, constraint, sql)
	}
}

func TestDefinitionDefersSelfReferenceToReusedKey(t *testing.T) {
	to := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name: "APP",
		Tables: []*metadatapb.TableMetadata{{
			Name: "T",
			Columns: []*metadatapb.ColumnMetadata{
				{Name: "ID", Type: "NUMBER", Position: 1, Nullable: false},
				{Name: "PARENT_ID", Type: "NUMBER", Position: 2, Nullable: true},
			},
			Indexes: []*metadatapb.IndexMetadata{{
				Name:         "T_ID_IDX",
				Type:         "NORMAL",
				Primary:      true,
				Expressions:  []string{`"ID"`},
				Descending:   []bool{false},
				Visible:      true,
				IsConstraint: true,
				Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "T_PK", Columns: []string{`"ID"`}, Primary: true}},
			}},
			ForeignKeys: []*metadatapb.ForeignKeyMetadata{{
				Name:              "T_FK",
				Columns:           []string{`"PARENT_ID"`},
				ReferencedTable:   "T",
				ReferencedColumns: []string{`"ID"`},
			}},
		}},
	})
	definition, err := DatabaseDefinition(to.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	table := definition[:strings.Index(definition, "\n);")]
	if strings.Contains(table, "FOREIGN KEY") {
		t.Fatalf("a self reference to a reused key stayed inside CREATE TABLE:\n%s", definition)
	}
	key := strings.Index(definition, `ADD CONSTRAINT "T_PK"`)
	fk := strings.Index(definition, `ADD CONSTRAINT "T_FK"`)
	if key < 0 || fk < 0 || key > fk {
		t.Fatalf("the foreign key precedes the key it references (key=%d fk=%d):\n%s", key, fk, definition)
	}
}

func TestReplacingReusedConstraintDropsItsIndex(t *testing.T) {
	table := func(primary bool) *metadatapb.SchemaMetadata {
		return &metadatapb.SchemaMetadata{
			Name: "APP",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "T",
				Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, Nullable: false}},
				Indexes: []*metadatapb.IndexMetadata{{
					Name:         "K",
					Type:         "NORMAL",
					Primary:      primary,
					Expressions:  []string{`"ID"`},
					Descending:   []bool{false},
					Visible:      true,
					IsConstraint: true,
					Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "K", Columns: []string{`"ID"`}, Primary: primary}},
				}},
			}},
		}
	}
	sql := migrationFor(t, catalogFromSchema(table(false)), catalogFromSchema(table(true)))
	constraint := strings.Index(sql, `DROP CONSTRAINT "K"`)
	index := strings.Index(sql, `DROP INDEX "APP"."K"`)
	create := strings.Index(sql, `CREATE INDEX "APP"."K"`)
	if constraint < 0 || index < 0 || create < 0 || !(constraint < create && index < create) {
		t.Fatalf("the retained index is not dropped before it is created again (constraint=%d index=%d create=%d):\n%s", constraint, index, create, sql)
	}
}

func TestReusedUniqueIndexKeepsConstraintName(t *testing.T) {
	to := catalogFromSchema(&metadatapb.SchemaMetadata{
		Name: "APP",
		Tables: []*metadatapb.TableMetadata{{
			Name: "T",
			Columns: []*metadatapb.ColumnMetadata{
				{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true},
				{Name: "V", Type: "NUMBER", Position: 2, Nullable: true},
			},
			Indexes: []*metadatapb.IndexMetadata{{
				Name:         "UK_VIRTUAL_IDX",
				Type:         "NORMAL",
				Unique:       true,
				Expressions:  []string{`"V"`},
				Descending:   []bool{false},
				Visible:      true,
				IsConstraint: true,
				Constraints:  []*metadatapb.IndexConstraintMetadata{{Name: "UK_VIRTUAL_UK", Columns: []string{`"V"`}}},
			}},
		}},
	})
	definition, err := DatabaseDefinition(to.Metadata())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(definition, `ALTER TABLE "T" ADD CONSTRAINT "UK_VIRTUAL_UK" UNIQUE ("V") USING INDEX "UK_VIRTUAL_IDX"`) {
		t.Fatalf("a reused unique index lost its constraint name:\n%s", definition)
	}
	if strings.Contains(definition, `CONSTRAINT "UK_VIRTUAL_IDX"`) {
		t.Fatalf("the index name was used as the constraint name:\n%s", definition)
	}
}

// TestViewDependenciesNameTheDerivedAnalyticViewsBase holds a view over an
// inline analytic view and one over a WITH analytic-view element to the
// persistent analytic view each names: the change depends on it, as on any
// other relation its query reads.
func TestViewDependenciesNameTheDerivedAnalyticViewsBase(t *testing.T) {
	meta, err := metadataFromSQL(`CREATE TABLE T (C NUMBER);
CREATE VIEW Z AS SELECT C FROM ANALYTIC VIEW(USING AV);
CREATE VIEW W AS WITH A ANALYTIC VIEW AS (USING AV) SELECT C FROM A`)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Z", "W"} {
		var view *metadatapb.ViewMetadata
		for _, candidate := range meta.Schemas[0].Views {
			if candidate.Name == name {
				view = candidate
			}
		}
		if view == nil {
			t.Fatalf("view %s missing", name)
		}
		found := false
		for _, dependency := range view.DependencyColumns {
			if dependency.Table == "AV" {
				found = true
			}
		}
		if !found {
			t.Fatalf("view %s: dependency on AV missing: %+v", name, view.DependencyColumns)
		}
	}
}

// TestPackageSpecNamingItselfCompiles holds a package spec that qualifies its
// own items, or two specs that name each other's, to compiling once: the
// checker reads a spec it is already reading as declared.
func TestPackageSpecNamingItselfCompiles(t *testing.T) {
	c, err := LoadSQL("CREATE PACKAGE p AS a NUMBER := 1; b NUMBER := p.a; END p;\n/\n" +
		"CREATE PACKAGE q AS x NUMBER := p.b; END q;\n/\n" +
		"CREATE OR REPLACE PACKAGE p AS a NUMBER := q.x; b NUMBER := p.a; END p;\n/\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range c.unitStates() {
		if u.status != "VALID" {
			t.Errorf("%s %s is %s %s", u.kind, u.name, u.status, u.err)
		}
	}
}

// TestTriggerOnMaterializedView holds a DML trigger on a materialized view,
// which Oracle creates, to being created.
func TestTriggerOnMaterializedView(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE t (id NUMBER PRIMARY KEY); CREATE MATERIALIZED VIEW mv AS SELECT id FROM t; CREATE TRIGGER trg BEFORE INSERT ON mv FOR EACH ROW BEGIN NULL; END;`)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.unitStates()) != 1 {
		t.Fatalf("units = %+v, want the trigger", c.unitStates())
	}
}

// TestPackageSpecChainCompilesOnce holds a chain of specs, each reading the
// one before twice, to reading each spec once rather than twice per level.
func TestPackageSpecChainCompilesOnce(t *testing.T) {
	var b strings.Builder
	b.WriteString("CREATE PACKAGE p0 AS a NUMBER := 1; END p0;\n/\n")
	for i := 1; i <= 40; i++ {
		fmt.Fprintf(&b, "CREATE PACKAGE p%d AS a NUMBER := p%d.a; b NUMBER := p%d.a; END p%d;\n/\n", i, i-1, i-1, i)
	}
	done := make(chan error, 1)
	go func() {
		_, err := LoadSQL(b.String())
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("40 chained package specs took over 30s")
	}
}

// TestTypeNameBlocksReplacedSynonym holds CREATE OR REPLACE SYNONYM to
// Oracle's namespace, which a type shares (ORA-00955).
func TestTypeNameBlocksReplacedSynonym(t *testing.T) {
	c, err := LoadSQL("CREATE TABLE x (id NUMBER);\nCREATE TYPE t AS OBJECT (a NUMBER);\n/\n")
	if err != nil {
		t.Fatal(err)
	}
	results, err := c.Exec("CREATE OR REPLACE SYNONYM t FOR x", nil)
	if err != nil {
		t.Fatal(err)
	}
	var refusal *Error
	if len(results) != 1 || !errors.As(results[0].Error, &refusal) || refusal.Code != "00955" {
		t.Fatalf("results = %+v, want ORA-00955", results)
	}
}

// TestTriggerDroppedAndCreatedAgain holds a trigger dropped and created
// again in one GetDatabaseMetadata script, which reads every statement on
// one extractor, to being created.
func TestTriggerDroppedAndCreatedAgain(t *testing.T) {
	_, err := GetDatabaseMetadata("CREATE TABLE x (id NUMBER);\nCREATE TRIGGER trg BEFORE INSERT ON x FOR EACH ROW BEGIN NULL; END;\n/\nDROP TRIGGER trg;\nCREATE TRIGGER trg BEFORE INSERT ON x FOR EACH ROW BEGIN NULL; END;\n/\n")
	if err != nil {
		t.Fatal(err)
	}
}

// TestUnitChangesReachTheirReaders holds the units book to what Oracle does
// around a unit's readers: a call through a synonym reads its target,
// dropping a table drops its triggers, an editioning view takes a BEFORE
// trigger, and ALTER ... COMPILE compiles a unit again.
func TestUnitChangesReachTheirReaders(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE PROCEDURE p IS BEGIN NULL; END;\n/\nCREATE SYNONYM s FOR p;\nCREATE PROCEDURE caller IS BEGIN s; END;\n/\nDROP PROCEDURE p;\n")
	if err != nil {
		t.Fatal(err)
	}
	if got := status(c, unitProcedure, "CALLER"); got != "INVALID" {
		t.Errorf("caller of a dropped procedure through a synonym: %s, want INVALID", got)
	}

	if _, err := GetDatabaseMetadata("CREATE TABLE t (id NUMBER);\nCREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW BEGIN NULL; END;\n/\nDROP TABLE t;\nCREATE TABLE t (id NUMBER);\nCREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW BEGIN NULL; END;\n/\n"); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadSQL("CREATE TABLE t (id NUMBER);\nCREATE EDITIONING VIEW ev AS SELECT id FROM t;\nCREATE TRIGGER trg BEFORE INSERT ON ev FOR EACH ROW BEGIN NULL; END;\n/\n"); err != nil {
		t.Fatal(err)
	}

	c, err = LoadSQL("CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT id INTO v FROM later; END;\n/\nCREATE TABLE later (id NUMBER);\n")
	if err != nil {
		t.Fatal(err)
	}
	if got := status(c, unitProcedure, "P"); got != "INVALID ORA-00942" {
		t.Fatalf("before COMPILE: %s, want INVALID ORA-00942", got)
	}
	if _, err := c.Exec("ALTER PROCEDURE p COMPILE", nil); err != nil {
		t.Fatal(err)
	}
	if got := status(c, unitProcedure, "P"); got != "VALID" {
		t.Errorf("after COMPILE: %s, want VALID", got)
	}

	c, err = LoadSQL("CREATE TABLE t (id NUMBER);\nCREATE PROCEDURE p IS v NUMBER; BEGIN SELECT id INTO v FROM t; END;\n/\nALTER TABLE t RENAME TO u;\n")
	if err != nil {
		t.Fatal(err)
	}
	if got := status(c, unitProcedure, "P"); got != "INVALID" {
		t.Errorf("reader of a table ALTER TABLE renamed: %s, want INVALID", got)
	}

	c, err = LoadSQL("CREATE TABLE t (id NUMBER);\nCREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW BEGIN nope; END;\n/\nALTER TRIGGER trg RENAME TO trg2;\n")
	if err != nil {
		t.Fatal(err)
	}
	if got := status(c, unitTrigger, "TRG2"); got != "INVALID PLS-00201" {
		t.Errorf("renamed invalid trigger: %s, want INVALID PLS-00201", got)
	}

	if _, err := GetDatabaseMetadata("CREATE TABLE t (id NUMBER PRIMARY KEY);\nCREATE MATERIALIZED VIEW mv AS SELECT id FROM t;\nCREATE TRIGGER trg BEFORE INSERT ON mv FOR EACH ROW BEGIN NULL; END;\n/\nDROP MATERIALIZED VIEW mv;\nCREATE MATERIALIZED VIEW mv AS SELECT id FROM t;\nCREATE TRIGGER trg BEFORE INSERT ON mv FOR EACH ROW BEGIN NULL; END;\n/\n"); err != nil {
		t.Fatal(err)
	}
}

// TestUnitReadsBindsTrimAndRowInserts holds three forms of a unit's body to
// Oracle: TRIM's mode word is no name, a trigger's :NEW field in static
// SQL is the table's column (PLS-00049 otherwise), and an INSERT of a whole
// row through a synonym reads the target table's every column.
func TestUnitReadsBindsTrimAndRowInserts(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE FUNCTION f RETURN VARCHAR2 IS BEGIN RETURN TRIM(BOTH 'x' FROM 'xxabcxx'); END;\n/\n" +
		"CREATE TABLE t (a NUMBER, b NUMBER);\nCREATE TABLE log (id NUMBER);\n" +
		"CREATE TRIGGER ok_trg BEFORE INSERT ON t FOR EACH ROW BEGIN INSERT INTO log (id) VALUES (:NEW.a); END;\n/\n" +
		"CREATE TRIGGER bad_trg BEFORE INSERT ON t FOR EACH ROW BEGIN INSERT INTO log (id) VALUES (:NEW.nope); END;\n/\n" +
		"CREATE SYNONYM s FOR t;\nCREATE PROCEDURE p IS BEGIN INSERT INTO s VALUES (1, 2); END;\n/\n" +
		"ALTER TABLE t DROP COLUMN b;\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitFunction, "F", "VALID"},
		{unitTrigger, "OK_TRG", "VALID"},
		{unitTrigger, "BAD_TRG", "INVALID PLS-00049"},
		{unitProcedure, "P", "INVALID"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}
}

// TestUnitNamesOwnSchemaLabelsAndMembers holds four more of Oracle's
// readings: a call naming the unit's own schema reads the callee, a
// materialized view anchors %ROWTYPE, a GOTO reaches a label in a loop,
// and replacing a spec invalidates only the callers of a member it drops.
func TestUnitNamesOwnSchemaLabelsAndMembers(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE TABLE app.t (id NUMBER PRIMARY KEY);\n" +
		"CREATE PROCEDURE app.p IS BEGIN NULL; END;\n/\n" +
		"CREATE PROCEDURE app.caller IS BEGIN app.p; END;\n/\n" +
		"CREATE MATERIALIZED VIEW app.mv AS SELECT id FROM app.t;\n" +
		"CREATE PROCEDURE app.rowtyped IS r mv%ROWTYPE; BEGIN NULL; END;\n/\n" +
		"CREATE PROCEDURE app.jumps IS BEGIN LOOP GOTO done; <<done>> EXIT; END LOOP; END;\n/\n" +
		"CREATE PACKAGE app.pk AS PROCEDURE run; PROCEDURE keep; END;\n/\n" +
		"CREATE PROCEDURE app.runs IS BEGIN pk.run; END;\n/\n" +
		"CREATE PROCEDURE app.keeps IS BEGIN pk.keep; END;\n/\n" +
		"CREATE OR REPLACE PACKAGE app.pk AS PROCEDURE other; PROCEDURE keep; END;\n/\n" +
		"DROP PROCEDURE app.p;\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitProcedure, "CALLER", "INVALID"},
		{unitProcedure, "ROWTYPED", "VALID"},
		{unitProcedure, "JUMPS", "VALID"},
		{unitProcedure, "RUNS", "INVALID"},
		{unitProcedure, "KEEPS", "VALID"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}
}

// TestUnitReadsPackageCursorsAndQuantifiers holds a cursor FOR loop over a
// package's cursor and a quantified comparison's list to compiling.
func TestUnitReadsPackageCursorsAndQuantifiers(t *testing.T) {
	c, err := LoadSQL("CREATE TABLE t (id NUMBER);\n" +
		"CREATE PACKAGE pk AS CURSOR c IS SELECT id FROM t; END;\n/\n" +
		"CREATE PROCEDURE p IS n NUMBER; BEGIN FOR r IN pk.c LOOP NULL; END LOOP; SELECT COUNT(*) INTO n FROM dual WHERE 1 = ANY (1, 2) AND 2 = SOME (1, 2) AND 3 <> ALL (1, 2); END;\n/\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range c.unitStates() {
		if u.status != "VALID" {
			t.Errorf("%s %s: %s %s, want VALID", u.kind, u.name, u.status, u.err)
		}
	}
}

// TestTypeMethodsCallEachOther holds a type body whose methods call each
// other, and themselves, without SELF to compiling.
func TestTypeMethodsCallEachOther(t *testing.T) {
	c, err := LoadSQL("CREATE TYPE pt AS OBJECT (x NUMBER, MEMBER FUNCTION f RETURN NUMBER, MEMBER FUNCTION g RETURN NUMBER);\n/\n" +
		"CREATE TYPE BODY pt AS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN g() + SELF.g + x; END; MEMBER FUNCTION g RETURN NUMBER IS BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN f; END; END;\n/\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range c.unitStates() {
		if u.status != "VALID" {
			t.Errorf("%s %s: %s %s, want VALID", u.kind, u.name, u.status, u.err)
		}
	}
}

// TestCallsParseCalleeOnce holds many calls to one stored procedure to one
// parse of its text per statement.
func TestCallsParseCalleeOnce(t *testing.T) {
	var body strings.Builder
	for range 2000 {
		body.WriteString("callee(1); ")
	}
	c, err := LoadSQL("CREATE PROCEDURE callee (n NUMBER) IS BEGIN NULL; END;\n/\nCREATE PROCEDURE caller IS BEGIN " + body.String() + "END;\n/\n")
	if err != nil {
		t.Fatal(err)
	}
	e := newMetadataExtractorFromMetadata(c.meta, "")
	e.savedSignature(c.meta.Schemas[0].Procedures[0].Definition, false)
	if len(e.signatures) != 1 {
		t.Fatalf("signatures cached = %d, want 1", len(e.signatures))
	}
}

// TestTriggersRowsAndReturning holds three more of Oracle's readings: an
// INSTEAD OF trigger reads :NEW without FOR EACH ROW, adding a column
// invalidates a whole-row INSERT but not one with a column list, and a
// RETURNING clause reads the table's columns into the unit's variables.
func TestTriggersRowsAndReturning(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE TABLE t (id NUMBER);\nCREATE VIEW v AS SELECT id FROM t;\n" +
		"CREATE TRIGGER trg INSTEAD OF INSERT ON v BEGIN INSERT INTO t (id) VALUES (:NEW.id); END;\n/\n" +
		"CREATE PROCEDURE whole IS BEGIN INSERT INTO t VALUES (1); END;\n/\n" +
		"CREATE PROCEDURE listed IS BEGIN INSERT INTO t (id) VALUES (1); END;\n/\n" +
		"CREATE PROCEDURE ret_ok IS v NUMBER; BEGIN INSERT INTO t (id) VALUES (1) RETURNING id INTO v; END;\n/\n" +
		"CREATE PROCEDURE ret_col IS v NUMBER; BEGIN UPDATE t SET id = 2 RETURNING nope INTO v; END;\n/\n" +
		"CREATE PROCEDURE ret_var IS BEGIN DELETE FROM t RETURNING id INTO missing_var; END;\n/\n" +
		"ALTER TABLE t ADD (b NUMBER);\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitTrigger, "TRG", "VALID"},
		{unitProcedure, "WHOLE", "INVALID"},
		{unitProcedure, "LISTED", "VALID"},
		{unitProcedure, "RET_OK", "VALID"},
		{unitProcedure, "RET_COL", "INVALID ORA-00904"},
		{unitProcedure, "RET_VAR", "INVALID PLS-00201"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}
}

// TestTypesConstructorsPipesAndWhen holds four more readings: a variable
// typed through a synonym reads the type, a constructor returns SELF AS
// RESULT with a bare RETURN, PIPE ROW's value is checked, and a trigger's
// WHEN condition reads its column.
func TestTypesConstructorsPipesAndWhen(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE TYPE ty AS OBJECT (a NUMBER, CONSTRUCTOR FUNCTION ty (x NUMBER) RETURN SELF AS RESULT);\n/\n" +
		"CREATE TYPE BODY ty AS CONSTRUCTOR FUNCTION ty (x NUMBER) RETURN SELF AS RESULT IS BEGIN a := x; RETURN; END; END;\n/\n" +
		"CREATE SYNONYM tys FOR ty;\nCREATE PROCEDURE typed IS v tys; BEGIN NULL; END;\n/\n" +
		"CREATE TYPE nums AS TABLE OF NUMBER;\n/\n" +
		"CREATE FUNCTION piped RETURN nums PIPELINED IS BEGIN PIPE ROW (missing_var); RETURN; END;\n/\n" +
		"CREATE TABLE t (a NUMBER, b NUMBER);\n" +
		"CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW WHEN (NEW.b > 0) BEGIN NULL; END;\n/\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{unitType, unitTypeBody} {
		if got := status(c, kind, "TY"); got != "VALID" {
			t.Errorf("%s TY with a constructor: %s, want VALID", kind, got)
		}
	}
	if _, err := c.Exec("DROP TYPE ty FORCE", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec("ALTER TABLE t DROP COLUMN b", nil); err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitProcedure, "TYPED", "INVALID"},
		{unitFunction, "PIPED", "INVALID PLS-00201"},
		{unitTrigger, "TRG", "INVALID"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}
}

// TestSubtypesIteratorsAndPackageBodies holds four more readings: static
// SQL reads a materialized view, a subtype's body reads inherited
// attributes, a generalized FOR loop declares every iterand, and a body
// that leaves out a parameterless spec procedure is PLS-00323.
func TestSubtypesIteratorsAndPackageBodies(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE TABLE t (id NUMBER PRIMARY KEY);\nCREATE MATERIALIZED VIEW mv AS SELECT id FROM t;\n" +
		"CREATE PROCEDURE reads_mv IS n NUMBER; BEGIN SELECT COUNT(*) INTO n FROM mv; END;\n/\n" +
		"CREATE TYPE base_t AS OBJECT (a NUMBER) NOT FINAL;\n/\n" +
		"CREATE TYPE sub_t UNDER base_t (MEMBER FUNCTION f RETURN NUMBER);\n/\n" +
		"CREATE TYPE BODY sub_t AS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN SELF.a + a; END; END;\n/\n" +
		"CREATE PROCEDURE pairs IS TYPE nt IS TABLE OF NUMBER; xs nt := nt(10, 20); BEGIN FOR i, v IN PAIRS OF xs LOOP DBMS_OUTPUT.PUT_LINE(i + v); END LOOP; END;\n/\n" +
		"CREATE PACKAGE pk AS PROCEDURE run; END;\n/\n" +
		"CREATE PACKAGE BODY pk AS PROCEDURE other IS BEGIN NULL; END; END;\n/\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitProcedure, "READS_MV", "VALID"},
		{unitTypeBody, "SUB_T", "VALID"},
		{unitProcedure, "PAIRS", "VALID"},
		{unitPackageBody, "PK", "INVALID PLS-00323"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}
}

// TestDefaultsRowSetsAndDependencies holds four more readings: DEFAULT in
// VALUES and SET, a row-valued SET, a trigger's UPDATE OF column, and a
// type created after a metadata load.
func TestDefaultsRowSetsAndDependencies(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE TABLE t (a NUMBER, b NUMBER);\n" +
		"CREATE PROCEDURE defaults IS BEGIN UPDATE t SET a = DEFAULT; INSERT INTO t (a) VALUES (DEFAULT); END;\n/\n" +
		"CREATE PROCEDURE row_set IS BEGIN UPDATE t SET (a, b) = (SELECT 1, 2 FROM dual); END;\n/\n" +
		"CREATE PROCEDURE row_set_wide IS BEGIN UPDATE t SET (a, b) = (SELECT 1, 2, 3 FROM dual); END;\n/\n" +
		"CREATE TRIGGER trg AFTER UPDATE OF b ON t FOR EACH ROW BEGIN NULL; END;\n/\n" +
		"ALTER TABLE t DROP COLUMN b;\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitProcedure, "DEFAULTS", "VALID"},
		{unitProcedure, "ROW_SET", "INVALID"},
		{unitProcedure, "ROW_SET_WIDE", "INVALID ORA-00913"},
		{unitTrigger, "TRG", "INVALID"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}

	m := New()
	m.SetMetadata(&metadatapb.DatabaseSchemaMetadata{Name: "APP", Schemas: []*metadatapb.SchemaMetadata{{Name: "APP"}}})
	for _, stmt := range []string{"CREATE TYPE ty AS OBJECT (a NUMBER)", "CREATE PROCEDURE p IS v ty; BEGIN NULL; END;", "DROP TYPE ty FORCE"} {
		if _, err := m.Exec(stmt, nil); err != nil {
			t.Fatal(err)
		}
	}
	if got := status(m, unitProcedure, "P"); got != "INVALID" {
		t.Errorf("reader of a dropped type after a metadata load: %s, want INVALID", got)
	}
}

// TestPatternsRecordsAndSnapshotTypeBodies holds four more readings: a
// compound MATCH_RECOGNIZE measure, a statement mixing a known table with a
// dictionary view, VALUES rec, and a type body after a metadata load.
func TestPatternsRecordsAndSnapshotTypeBodies(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE TABLE t (id NUMBER, v NUMBER);\n" +
		"CREATE PROCEDURE pattern IS n NUMBER; BEGIN SELECT MAX(m) INTO n FROM t MATCH_RECOGNIZE (ORDER BY id MEASURES MATCH_NUMBER() + 1 AS m PATTERN (a b*) DEFINE b AS v > PREV(v)); END;\n/\n" +
		"CREATE PROCEDURE mixed IS n NUMBER; BEGIN SELECT COUNT(*) INTO n FROM t WHERE EXISTS (SELECT 1 FROM user_tables); END;\n/\n" +
		"CREATE PROCEDURE rec IS BEGIN INSERT INTO t VALUES missing_rec; END;\n/\n" +
		"DROP TABLE t;\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitProcedure, "PATTERN", "INVALID"},
		{unitProcedure, "MIXED", "INVALID"},
		{unitProcedure, "REC", "INVALID PLS-00201"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}

	m := New()
	m.SetMetadata(&metadatapb.DatabaseSchemaMetadata{Name: "APP", Schemas: []*metadatapb.SchemaMetadata{{Name: "APP"}}})
	if _, err := m.Exec("CREATE OR REPLACE TYPE BODY synced_t AS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;", nil); err != nil {
		t.Fatal(err)
	}
	if got := status(m, unitTypeBody, "SYNCED_T"); got != "VALID" {
		t.Errorf("type body after a metadata load: %s, want VALID", got)
	}
}

// TestXMLTypeSynonymsAndUsingKeys holds three more readings: a built-in
// type's constructor in static SQL, dropping a public synonym a unit reads
// through, and a JOIN USING key's column.
func TestXMLTypeSynonymsAndUsingKeys(t *testing.T) {
	status := func(c *Catalog, kind, name string) string {
		for _, u := range c.unitStates() {
			if u.kind == kind && u.name == name {
				return strings.TrimSpace(u.status + " " + u.err)
			}
		}
		return "missing"
	}
	c, err := LoadSQL("CREATE TABLE t (id NUMBER, x NUMBER);\nCREATE TABLE u (id NUMBER, y NUMBER);\nCREATE TABLE w (id NUMBER);\nCREATE PUBLIC SYNONYM ws FOR w;\n" +
		"CREATE PROCEDURE xml IS x XMLTYPE; BEGIN SELECT XMLTYPE('<r/>') INTO x FROM dual; x := XMLTYPE('<s/>'); END;\n/\n" +
		"CREATE PROCEDURE pub IS n NUMBER; BEGIN SELECT COUNT(*) INTO n FROM ws; END;\n/\n" +
		"CREATE PROCEDURE joined IS n NUMBER; BEGIN SELECT COUNT(*) INTO n FROM t JOIN u USING (id); END;\n/\n" +
		"DROP PUBLIC SYNONYM ws;\nALTER TABLE t DROP COLUMN id;\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct{ kind, name, status string }{
		{unitProcedure, "XML", "VALID"},
		{unitProcedure, "PUB", "INVALID"},
		{unitProcedure, "JOINED", "INVALID"},
	} {
		if got := status(c, want.kind, want.name); got != want.status {
			t.Errorf("%s %s: %s, want %s", want.kind, want.name, got, want.status)
		}
	}
}

// TestForeignTargetsAndSnapshotSupertypes holds a DML statement writing
// another schema's table, and a subtype after a metadata load, to compiling.
func TestForeignTargetsAndSnapshotSupertypes(t *testing.T) {
	m := New()
	m.SetMetadata(&metadatapb.DatabaseSchemaMetadata{Name: "APP", Schemas: []*metadatapb.SchemaMetadata{{Name: "APP"}}})
	for _, stmt := range []string{
		"CREATE PROCEDURE p IS BEGIN INSERT INTO other.t (id) VALUES (1); UPDATE other.t SET id = 2; DELETE FROM other.t; END;",
		"CREATE TYPE child_t UNDER base_t (b NUMBER)",
	} {
		if _, err := m.Exec(stmt, nil); err != nil {
			t.Fatal(err)
		}
	}
	for _, u := range m.unitStates() {
		if u.status != "VALID" {
			t.Errorf("%s %s: %s %s, want VALID", u.kind, u.name, u.status, u.err)
		}
	}
}
