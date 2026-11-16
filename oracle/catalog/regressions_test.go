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

func TestDropNamedNotNullConstraintMakesColumnNullable(t *testing.T) {
	c := catalogFromSchema(&metadatapb.SchemaMetadata{Tables: []*metadatapb.TableMetadata{{
		Name:             "T",
		Columns:          []*metadatapb.ColumnMetadata{{Name: "C", Type: "NUMBER", Position: 1}},
		CheckConstraints: []*metadatapb.CheckConstraintMetadata{{Name: "NN", Expression: `"C" IS NOT NULL`}},
	}}})
	if _, err := c.Exec(`ALTER TABLE T DROP CONSTRAINT NN`, nil); err != nil {
		t.Fatal(err)
	}
	table := tableByName(c.Metadata().Schemas[0])["T"]
	if !table.Columns[0].Nullable || len(table.CheckConstraints) != 0 {
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

func TestCreateOrReplaceTriggerMovesTarget(t *testing.T) {
	c, err := LoadSQL(`CREATE TABLE A (ID NUMBER); CREATE TABLE Z (ID NUMBER); CREATE TRIGGER TR BEFORE INSERT ON Z FOR EACH ROW BEGIN NULL; END; CREATE OR REPLACE TRIGGER TR BEFORE INSERT ON A FOR EACH ROW BEGIN NULL; END;`)
	if err != nil {
		t.Fatal(err)
	}
	schema := c.Metadata().Schemas[0]
	if len(tableByName(schema)["Z"].Triggers) != 0 || len(tableByName(schema)["A"].Triggers) != 1 {
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
