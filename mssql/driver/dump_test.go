package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDumpRestores syncs a corpus, dumps it, checks the renderings a restore
// cannot judge, empties the schema, applies the dump, and requires the resynced
// database to equal the source.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "dump", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	ctx := context.Background()
	execCorpus(t, d.db, string(corpus))
	source, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var script strings.Builder
	if err := d.Dump(ctx, &script, source); err != nil {
		t.Fatal(err)
	}
	dump := script.String()
	assertRendering(t, dump)
	execCorpus(t, d.db, dropScript(source))
	execCorpus(t, d.db, dump)
	restored, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}

	// The fixture's partitioned rowstore index reports its partitioning column
	// at key_ordinal 0, which the sync must put in neither the key list nor
	// INCLUDE. Restoring without partitioning then agrees, because the index
	// carries only the key it declares.
	var found bool
	for _, schema := range source.Schemas {
		if schema.Name != "dump" {
			continue
		}
		for _, table := range schema.Tables {
			if table.Name != "events" {
				continue
			}
			for _, index := range table.Indexes {
				if index.Name != "IX_events_id" {
					continue
				}
				found = true
				if diff := cmp.Diff([]string{"id"}, index.Expressions); diff != "" {
					t.Errorf("the partitioned index kept the implicit partition column as a key (-want +got):\n%s", diff)
				}
			}
		}
	}
	if !found {
		t.Fatal("the partitioned fixture's index is missing from the sync; if [dump].[events] left the corpus, remove this check")
	}
}

// TestDumpKeepsAnEncryptedModuleCommentOut pins the placeholder a sync writes
// where the server withholds an encrypted module's body. The placeholder is a
// comment, so restoring it creates nothing, and mssql/catalog's
// hasRestorableDefinition recognizes exactly its text to keep the module's
// comment from being attached to an object the dump never creates. Changing
// the text on either side fails here, which is the point: the dump would
// otherwise emit sp_addextendedproperty for a module that does not exist.
func TestDumpKeepsAnEncryptedModuleCommentOut(t *testing.T) {
	const setup = `CREATE SCHEMA [enc];
GO
CREATE PROCEDURE [enc].[p] WITH ENCRYPTION AS SELECT 1 AS n;
GO
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'procedure note', @level0type = N'SCHEMA', @level0name = N'enc', @level1type = N'PROCEDURE', @level1name = N'p';
GO
CREATE VIEW [enc].[v] WITH ENCRYPTION AS SELECT 1 AS n;
GO
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'view note', @level0type = N'SCHEMA', @level0name = N'enc', @level1type = N'VIEW', @level1name = N'v';
GO
CREATE FUNCTION [enc].[f]() RETURNS INT WITH ENCRYPTION AS BEGIN RETURN 1 END;
GO
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'function note', @level0type = N'SCHEMA', @level0name = N'enc', @level1type = N'FUNCTION', @level1name = N'f';
GO
`
	d := openTestDriver(t)
	execCorpus(t, d.db, setup)
	ctx := context.Background()
	source, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dump strings.Builder
	if err := d.Dump(ctx, &dump, source); err != nil {
		t.Fatal(err)
	}
	for _, note := range []string{"procedure note", "view note", "function note"} {
		if strings.Contains(dump.String(), note) {
			t.Errorf("the dump attaches the comment %q to a module whose body the server withholds; the placeholder no longer matches mssql/catalog's hasRestorableDefinition", note)
		}
	}
	for _, placeholder := range []string{
		encryptedDefinition("procedure", "enc.p"),
		encryptedDefinition("view", "enc.v"),
		encryptedDefinition("function", "enc.f"),
	} {
		if !strings.Contains(dump.String(), placeholder) {
			t.Errorf("the dump does not carry the placeholder %q a sync writes for an encrypted module", placeholder)
		}
	}
}

// assertRendering pins what restoring the dump cannot judge. A filtered index
// restored without its predicate, and a view followed by a second semicolon,
// both execute and resync into a database SQL Server treats as the source's,
// so only the script's text holds the shapes Dump renders.
func assertRendering(t *testing.T, dump string) {
	t.Helper()
	// The predicate is SQL Server's text, brackets and term order kept, and it
	// follows the key and INCLUDE lists as part of the index.
	const filteredIndex = ") INCLUDE ([status])\nWHERE ([parent_id]>(0) AND [note] IS NOT NULL);"
	if !strings.Contains(dump, filteredIndex) {
		t.Errorf("the filtered index does not render its predicate after its INCLUDE list; restoring the dump then drops the filter. Want:\n%s", filteredIndex)
	}
	// A synced view definition keeps the trailing semicolon the CREATE VIEW was
	// run with, so the dump trims it and adds one back. Appending
	// unconditionally would grow a semicolon on every restore.
	const view = "CREATE VIEW [dump].[v_parents] AS\n" +
		"\tSELECT [id], [name] FROM [dump].[parents];\n" +
		"\nGO\n\n" +
		"EXEC sp_addextendedproperty 'MS_Description', N'view note', 'SCHEMA', 'dump', 'VIEW', 'v_parents';\n" +
		"GO\n"
	if !strings.Contains(dump, view) {
		t.Errorf("the view does not render its body once, ending in one semicolon, with its comment in the next batch. Want:\n%s", view)
	}
}

func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source, restored, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "data_free", "row_count"),
	)
}

// dropScript empties every schema in snapshot so a dump can create it again.
func dropScript(snapshot *metadatapb.DatabaseSchemaMetadata) string {
	var lines []string
	schemas := slices.Clone(snapshot.Schemas)
	slices.SortFunc(schemas, func(a, b *metadatapb.SchemaMetadata) int { return strings.Compare(a.Name, b.Name) })
	for _, schema := range schemas {
		views := slices.Clone(schema.Views)
		slices.SortFunc(views, func(a, b *metadatapb.ViewMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, view := range views {
			lines = append(lines, fmt.Sprintf("DROP VIEW IF EXISTS [%s].[%s];", schema.Name, view.Name))
		}
		procedures := slices.Clone(schema.Procedures)
		slices.SortFunc(procedures, func(a, b *metadatapb.ProcedureMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, procedure := range procedures {
			lines = append(lines, fmt.Sprintf("DROP PROCEDURE IF EXISTS [%s].[%s];", schema.Name, procedure.Name))
		}
		functions := slices.Clone(schema.Functions)
		slices.SortFunc(functions, func(a, b *metadatapb.FunctionMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, function := range functions {
			lines = append(lines, fmt.Sprintf("DROP FUNCTION IF EXISTS [%s].[%s];", schema.Name, function.Name))
		}
		synonyms := slices.Clone(schema.Synonyms)
		slices.SortFunc(synonyms, func(a, b *metadatapb.SynonymMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, synonym := range synonyms {
			lines = append(lines, fmt.Sprintf("DROP SYNONYM IF EXISTS [%s].[%s];", schema.Name, synonym.Name))
		}
		sequences := slices.Clone(schema.Sequences)
		slices.SortFunc(sequences, func(a, b *metadatapb.SequenceMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, sequence := range sequences {
			lines = append(lines, fmt.Sprintf("DROP SEQUENCE IF EXISTS [%s].[%s];", schema.Name, sequence.Name))
		}
	}
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			fks := slices.Clone(table.ForeignKeys)
			slices.SortFunc(fks, func(a, b *metadatapb.ForeignKeyMetadata) int { return strings.Compare(a.Name, b.Name) })
			for _, fk := range fks {
				lines = append(lines, fmt.Sprintf("ALTER TABLE [%s].[%s] DROP CONSTRAINT IF EXISTS [%s];", schema.Name, table.Name, fk.Name))
			}
		}
	}
	for _, schema := range schemas {
		tables := slices.Clone(schema.Tables)
		slices.SortFunc(tables, func(a, b *metadatapb.TableMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, table := range tables {
			lines = append(lines, fmt.Sprintf("DROP TABLE IF EXISTS [%s].[%s];", schema.Name, table.Name))
		}
	}
	for _, schema := range schemas {
		if !strings.EqualFold(schema.Name, "dbo") {
			lines = append(lines, fmt.Sprintf("DROP SCHEMA IF EXISTS [%s];", schema.Name))
		}
	}
	return strings.Join(lines, "\n")
}
