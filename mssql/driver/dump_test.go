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

// TestDumpRestores syncs a corpus, dumps it, empties the schema, applies the
// dump, and requires the resynced database to equal the source.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "dump", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	source, restored, _ := dumpAndRestore(t, d, string(corpus))
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}
}

// TestDumpGaps pins what a dump cannot faithfully restore because the driver's
// sync does not carry the object's option, predicate, included columns or
// comment. Each entry fails once the gap closes, so the object moves into
// TestDumpRestores.
func TestDumpGaps(t *testing.T) {
	const setup = `CREATE SCHEMA [gap_seq];
GO
CREATE SEQUENCE [gap_seq].[s] AS INT START WITH 10 INCREMENT BY 5 MINVALUE 1 MAXVALUE 1000 CACHE 20;
GO
CREATE SCHEMA [gap_filter];
GO
CREATE TABLE [gap_filter].[t] ([id] INT NOT NULL, [v] INT NULL);
GO
CREATE UNIQUE INDEX [UX_t_v] ON [gap_filter].[t] ([v]) WHERE [v] IS NOT NULL;
GO
CREATE SCHEMA [gap_include];
GO
CREATE TABLE [gap_include].[t] ([id] INT NOT NULL, [payload] NVARCHAR(20) NULL);
GO
CREATE INDEX [IX_t_payload] ON [gap_include].[t] ([id]) INCLUDE ([payload]);
GO
CREATE SCHEMA [gap_fn];
GO
CREATE FUNCTION [gap_fn].[f]() RETURNS INT AS BEGIN RETURN 1 END;
GO
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'function note', @level0type = N'SCHEMA', @level0name = N'gap_fn', @level1type = N'FUNCTION', @level1name = N'f';
GO
CREATE SCHEMA [gap_view];
GO
CREATE VIEW [gap_view].[v] AS SELECT 1 AS n;
GO
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'view note', @level0type = N'SCHEMA', @level0name = N'gap_view', @level1type = N'VIEW', @level1name = N'v';
GO
CREATE SCHEMA [gap_seq_comment];
GO
CREATE SEQUENCE [gap_seq_comment].[s] AS INT;
GO
EXEC sys.sp_addextendedproperty @name = N'MS_Description', @value = N'sequence note', @level0type = N'SCHEMA', @level0name = N'gap_seq_comment', @level1type = N'SEQUENCE', @level1name = N's';
GO
`
	gaps := []struct {
		name    string
		missing []string
	}{
		{
			// getSequences reads no options, so the dump emits only the type.
			name:    "sequence options",
			missing: []string{"START WITH 10", "INCREMENT BY 5", "CACHE 20"},
		},
		{
			// getIndexes reads no filter_definition.
			name:    "filtered index",
			missing: []string{"WHERE"},
		},
		{
			// getIndexes drops key_ordinal = 0 columns.
			name:    "index include columns",
			missing: []string{"INCLUDE"},
		},
		{name: "function comment", missing: []string{"function note"}},
		{name: "view comment", missing: []string{"view note"}},
		{name: "sequence comment", missing: []string{"sequence note"}},
	}
	d := openTestDriver(t)
	source, restored, dump := dumpAndRestore(t, d, setup)
	for _, gap := range gaps {
		t.Run(gap.name, func(t *testing.T) {
			for _, want := range gap.missing {
				if strings.Contains(dump, want) {
					t.Fatalf("the dump now carries %q; if the gap closed, remove this entry", want)
				}
			}
		})
	}
	if diff := restoreDiff(source, restored); diff == "" {
		t.Fatalf("the restore reproduced the source; if every gap closed, remove this test")
	}
}

func dumpAndRestore(t *testing.T, d *Driver, setup string) (*metadatapb.DatabaseSchemaMetadata, *metadatapb.DatabaseSchemaMetadata, string) {
	t.Helper()
	ctx := context.Background()
	execCorpus(t, d.db, setup)
	source, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dump strings.Builder
	if err := d.Dump(ctx, &dump, source); err != nil {
		t.Fatal(err)
	}
	execCorpus(t, d.db, dropScript(source))
	execCorpus(t, d.db, dump.String())
	restored, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return source, restored, dump.String()
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
