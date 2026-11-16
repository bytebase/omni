package catalog

import (
	"context"
	"slices"
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
