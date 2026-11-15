package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// diffSyncLosses lists the schema sides whose snapshot diffs non-empty against
// the schema text it was synced from, because sync leaves out what the text
// declares (see syncLosses), or whose text omni does not load. An entry that
// loads and diffs empty fails, so a fix removes it.
var diffSyncLosses = map[string]string{
	"Alter table with column additions, type changes, and constraints/old":                  lossOnUpdate,
	"Reverse: Alter table columns/new":                                                      lossOnUpdate,
	"Multiple triggers and audit logging/old":                                               lossOnUpdate,
	"Multiple triggers and audit logging/new":                                               lossOnUpdate,
	"Reverse: Drop triggers and views/old":                                                  lossOnUpdate,
	"Reverse: Drop triggers and views/new":                                                  lossOnUpdate,
	"MariaDB column types and defaults/new":                                                 lossOnUpdate,
	"MariaDB timestamps/new":                                                                lossOnUpdate,
	"Operations on partitioned tables/new":                                                  lossPartitions,
	"Re-adding a primary key and an overlapping secondary index on a partitioned table/new": lossPartitions,
	"Recreating a partitioned table from metadata/new":                                      lossPartitions,
	"MariaDB partitioned tables/old":                                                        lossPartitions,
	"MariaDB partitioned tables/new":                                                        lossPartitions,
	"MariaDB system-time partitioning/new":                                                  lossPartitions,
	"MariaDB invisible columns/new":                                                         lossInvisible,
	"MariaDB system-versioned tables/new":                                                   lossVersioning,
	"MariaDB application-time periods/new":                                                  lossPeriod,
	"MariaDB table options/new":                                                             lossHash,
	"MariaDB index forms/new":                                                               lossIgnored,
	"MariaDB views/new":                                                                     lossView,
	"MariaDB vector index/new":                                                              lossVector,
	"MariaDB compressed columns/new":                                                        gapCompressed,
	"MariaDB routines/new":                                                                  gapForRange,
}

// TestOracle_DiffSnapshotAgainstSchemaText loads each snapshot synced from the
// pinned server and the schema text the server held when it was taken, and
// requires an empty diff in both directions: the server's answer is that the
// two are one schema. It also requires the snapshot's SDL to diff empty
// against the snapshot.
func TestOracle_DiffSnapshotAgainstSchemaText(t *testing.T) {
	for _, tc := range readSnapshotCases(t, pinnedSnapshots) {
		for _, side := range []struct {
			name   string
			schema string
			meta   map[string]any
		}{{"old", tc.OldSchema, tc.OldMetadata}, {"new", tc.NewSchema, tc.NewMetadata}} {
			key := tc.Description + "/" + side.name
			t.Run(key, func(t *testing.T) {
				meta := metadataFixture(t, side.meta)
				snapshot := New()
				if _, err := snapshot.LoadMetadataInto(context.Background(), meta, "bbcatalog"); err != nil {
					t.Fatal(err)
				}
				loss, listed := diffSyncLosses[key]
				text, err := schemaTextCatalog(side.schema)
				if err != nil {
					if !listed {
						t.Fatal(err)
					}
					return
				}
				forward, backward := Diff(text, snapshot), Diff(snapshot, text)
				switch empty := forward.IsEmpty() && backward.IsEmpty(); {
				case !empty && !listed:
					t.Errorf("snapshot and schema text differ:\n%s", describeDiff(forward))
				case empty && listed:
					t.Errorf("listed as a sync loss (%s) but diffs empty: remove the entry", loss)
				}

				sdl, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: true})
				if err != nil {
					t.Fatal(err)
				}
				fromSDL, err := LoadSDL("CREATE DATABASE `bbcatalog`;\nUSE `bbcatalog`;\n" + sdl)
				if err != nil {
					t.Fatalf("LoadSDL of the snapshot's SDL: %v\n%s", err, sdl)
				}
				if d := Diff(snapshot, fromSDL); !d.IsEmpty() {
					t.Errorf("snapshot and its SDL differ:\n%s", describeDiff(d))
				}
			})
		}
	}
}

// describeDiff renders a diff compactly for a failure message, with each
// changed column's and key's two sides.
func describeDiff(d *SchemaDiff) string {
	var b strings.Builder
	n := NormalizerFor(MariaDB118)
	for _, te := range d.Tables {
		fmt.Fprintf(&b, "table %s %s\n", te.Name, te.Action)
		if te.Action == DiffModify && n.CanonicalTableOptions(te.From) != n.CanonicalTableOptions(te.To) {
			fmt.Fprintf(&b, "  options: %s\n        -> %s\n", n.CanonicalTableOptions(te.From), n.CanonicalTableOptions(te.To))
		}
		for _, ce := range te.Columns {
			fmt.Fprintf(&b, "  column %s %s", ce.Name, ce.Action)
			if ce.Action == DiffModify {
				fmt.Fprintf(&b, ": %s\n        -> %s", n.CanonicalColumn(te.From, ce.From), n.CanonicalColumn(te.To, ce.To))
			}
			b.WriteString("\n")
		}
		for _, ie := range te.Indexes {
			fmt.Fprintf(&b, "  index %s %s", ie.Name, ie.Action)
			if ie.From != nil && ie.To != nil {
				fmt.Fprintf(&b, ": %s -> %s", showIndex(ie.From), showIndex(ie.To))
			}
			b.WriteString("\n")
		}
		for _, fe := range te.ForeignKeys {
			fmt.Fprintf(&b, "  foreign key %s %s\n", fe.Name, fe.Action)
		}
		for _, ce := range te.Checks {
			fmt.Fprintf(&b, "  check %s %s\n", ce.Name, ce.Action)
		}
		if te.PartitionChanged {
			b.WriteString("  partitioning\n")
		}
	}
	for _, e := range d.Views {
		fmt.Fprintf(&b, "view %s %s\n", e.Name, e.Action)
	}
	for _, e := range d.Functions {
		fmt.Fprintf(&b, "function %s %s\n", e.Name, e.Action)
	}
	for _, e := range d.Procedures {
		fmt.Fprintf(&b, "procedure %s %s\n", e.Name, e.Action)
	}
	for _, e := range d.Triggers {
		fmt.Fprintf(&b, "trigger %s %s\n", e.Name, e.Action)
	}
	for _, e := range d.Events {
		fmt.Fprintf(&b, "event %s %s\n", e.Name, e.Action)
	}
	for _, e := range d.Sequences {
		fmt.Fprintf(&b, "sequence %s %s\n", e.Name, e.Action)
	}
	return b.String()
}

// schemaTextCatalog loads a fixture's schema text into the database bbcatalog:
// as SDL, or, for a script SDL refuses (one that alters what it created),
// statement by statement with foreign key checks off.
func schemaTextCatalog(schema string) (*Catalog, error) {
	const prefix = "CREATE DATABASE `bbcatalog`;\nUSE `bbcatalog`;\n"
	c, err := LoadSDL(prefix + schema)
	if err == nil || !strings.Contains(err.Error(), "SDL does not allow") {
		return c, err
	}
	c = New()
	c.SetForeignKeyChecks(false)
	results, err := c.Exec(prefix+schema, &ExecOptions{ContinueOnError: true})
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.Error != nil {
			return nil, fmt.Errorf("line %d: %w", r.Line, r.Error)
		}
	}
	return c, nil
}
