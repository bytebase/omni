package driver

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	catalogapi "github.com/databricks/databricks-sdk-go/service/catalog"

	dbxanalysis "github.com/bytebase/omni/databricks/analysis"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

const (
	schemaStmtFmt = "" +
		"--\n" +
		"-- %s structure for %s\n" +
		"--\n" +
		"%s;\n"
)

var (
	sysCatalog = "system"
	infoSchema = "information_schema"
	dftSchema  = "default"
)

// Dump writes every catalog's, schema's and object's DDL as the server states
// it. The metadata is not rendered and does not narrow the dump: omni has no
// databricks catalog that turns a synced schema back into Databricks DDL, so
// SHOW CREATE TABLE states a managed or external table, a view and a
// materialized view. The system catalog and each catalog's information_schema
// are left out, as Bytebase left them out.
func (d *Driver) Dump(ctx context.Context, writer io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	catalogMap, err := d.listCatalogTables(ctx, "")
	if err != nil {
		return err
	}

	catalogNames := make([]string, 0, len(catalogMap))
	for catalogName := range catalogMap {
		catalogNames = append(catalogNames, catalogName)
	}
	sort.Strings(catalogNames)

	for _, catalogName := range catalogNames {
		if catalogName == sysCatalog {
			continue
		}
		if _, err := fmt.Fprintf(writer, "CREATE CATALOG %s;\n", quoteBacktick(catalogName)); err != nil {
			return err
		}
	}

	for _, catalogName := range catalogNames {
		if catalogName == sysCatalog {
			continue
		}
		for _, schemaName := range sortedSchemas(catalogMap[catalogName]) {
			if schemaName == infoSchema || schemaName == dftSchema {
				continue
			}
			if _, err := fmt.Fprintf(writer, "CREATE SCHEMA %s.%s;\n", quoteBacktick(catalogName), quoteBacktick(schemaName)); err != nil {
				return err
			}
		}
	}

	// Base tables first, across every schema, so a view may name one in a
	// schema the dump has not reached yet.
	for _, catalogName := range catalogNames {
		if catalogName == sysCatalog {
			continue
		}
		schemaMap := catalogMap[catalogName]
		for _, schemaName := range sortedSchemas(schemaMap) {
			for _, tblUnion := range sortedTables(schemaMap[schemaName]) {
				if tblUnion.typeName != catalogapi.TableTypeManaged && tblUnion.typeName != catalogapi.TableTypeExternal {
					continue
				}
				qualifiedName, err := getQualifiedTblName(catalogName, schemaName, tblUnion.name)
				if err != nil {
					return err
				}
				ddl, err := d.showCreateTable(ctx, qualifiedName)
				if err != nil {
					return err
				}
				label := "MANAGED TABLE"
				if tblUnion.typeName == catalogapi.TableTypeExternal {
					label = "EXTERNAL TABLE"
				}
				if _, err := fmt.Fprintf(writer, schemaStmtFmt, label, qualifiedName, ddl); err != nil {
					return err
				}
			}
		}
	}

	// Views and materialized views after the objects they may read, in
	// dependency order across schemas.
	var objectViews []dumpObject
	dependencies := map[string][]string{}
	for _, catalogName := range catalogNames {
		if catalogName == sysCatalog {
			continue
		}
		schemaMap := catalogMap[catalogName]
		for _, schemaName := range sortedSchemas(schemaMap) {
			if schemaName == infoSchema {
				continue
			}
			for _, tblUnion := range sortedTables(schemaMap[schemaName]) {
				if tblUnion.typeName != catalogapi.TableTypeView && tblUnion.typeName != catalogapi.TableTypeMaterializedView {
					continue
				}
				object := dumpObject{catalog: catalogName, schema: schemaName, name: tblUnion.name, kind: tblUnion.typeName}
				objectViews = append(objectViews, object)
				// ListTables does not state a view's dependencies, so read
				// the object itself before ordering it.
				info, err := d.client.Tables.GetByFullName(ctx, object.key())
				if err != nil {
					return fmt.Errorf("failed to read dependencies of %q: %w", object.key(), err)
				}
				if info.ViewDependencies != nil {
					for _, dependency := range info.ViewDependencies.Dependencies {
						if dependency.Table != nil {
							dependencies[object.key()] = append(dependencies[object.key()], dependency.Table.TableFullName)
						}
					}
				}
				if len(dependencies[object.key()]) == 0 {
					var viewDef string
					if tblUnion.view != nil {
						viewDef = tblUnion.view.Definition
					} else if tblUnion.materialView != nil {
						viewDef = tblUnion.materialView.Definition
					}
					if viewDef != "" {
						if extracted, err := dbxanalysis.ExtractViewDependencies(viewDef, catalogName, schemaName); err == nil {
							for _, dep := range extracted {
								dependencies[object.key()] = append(dependencies[object.key()], dep.FullName())
							}
						}
					}
				}
			}
		}
	}

	for _, object := range orderViewObjects(objectViews, dependencies) {
		qualifiedName, err := getQualifiedTblName(object.catalog, object.schema, object.name)
		if err != nil {
			return err
		}
		ddl, err := d.showCreateTable(ctx, qualifiedName)
		if err != nil {
			return err
		}
		label := "MATERIALIZED VIEW"
		if object.kind == catalogapi.TableTypeView {
			label = "VIEW"
			ddl = qualifyViewDDL(ddl, qualifiedName)
		}
		if _, err := fmt.Fprintf(writer, schemaStmtFmt, label, qualifiedName, ddl); err != nil {
			return err
		}
	}

	return nil
}

// dumpObject is one view or materialized view in the order the dump writes it.
type dumpObject struct {
	catalog, schema, name string
	kind                  catalogapi.TableType
}

func (o dumpObject) key() string { return o.catalog + "." + o.schema + "." + o.name }

func sortedSchemas(schemaMap databricksSchemaMap) []string {
	names := make([]string, 0, len(schemaMap))
	for name := range schemaMap {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedTables(tables []*tableUnion) []*tableUnion {
	out := append([]*tableUnion(nil), tables...)
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

// orderViewObjects puts a view after every view it reads. The service states a
// view's dependencies by full table name; a dependency on a base table is
// already created and does not affect the order. Lookups normalize identifiers
// to lowercase to match Databricks' case-insensitive identifier semantics.
func orderViewObjects(objects []dumpObject, dependencies map[string][]string) []dumpObject {
	byKey := make(map[string]dumpObject, len(objects))
	indegree := make(map[string]int, len(objects))
	children := map[string][]string{}
	for _, object := range objects {
		key := strings.ToLower(object.key())
		byKey[key] = object
		indegree[key] = 0
	}
	for rawKey, list := range dependencies {
		key := strings.ToLower(rawKey)
		if _, ok := byKey[key]; !ok {
			continue
		}
		seen := map[string]bool{}
		for _, rawDep := range list {
			dependency := strings.ToLower(rawDep)
			if dependency == key || seen[dependency] {
				continue
			}
			if _, ok := byKey[dependency]; !ok {
				continue
			}
			seen[dependency] = true
			children[dependency] = append(children[dependency], key)
			indegree[key]++
		}
	}
	for key := range children {
		sort.Strings(children[key])
	}

	ready := make([]string, 0, len(objects))
	for key, degree := range indegree {
		if degree == 0 {
			ready = append(ready, key)
		}
	}
	sort.Strings(ready)

	ordered := make([]dumpObject, 0, len(objects))
	for len(ready) > 0 {
		key := ready[0]
		ready = ready[1:]
		ordered = append(ordered, byKey[key])
		for _, child := range children[key] {
			indegree[child]--
			if indegree[child] == 0 {
				ready = append(ready, child)
			}
		}
		sort.Strings(ready)
	}
	if len(ordered) != len(objects) {
		// The service should not state a cycle. Keep the rest in name order
		// rather than dropping it.
		remaining := make([]string, 0, len(objects)-len(ordered))
		for key := range byKey {
			if indegree[key] > 0 {
				remaining = append(remaining, key)
			}
		}
		sort.Strings(remaining)
		for _, key := range remaining {
			ordered = append(ordered, byKey[key])
		}
	}
	return ordered
}

// qualifyViewDDL replaces a view statement's target with the catalog-qualified
// name. The workspace states a view as CREATE VIEW schema.view, which would
// create it in another catalog when the dump is replayed elsewhere.
func qualifyViewDDL(ddl, qualifiedName string) string {
	return dbxanalysis.QualifyViewDDL(ddl, qualifiedName)
}

// showCreateTable asks the server for one object's own DDL.
func (d *Driver) showCreateTable(ctx context.Context, qualifiedTblName string) (string, error) {
	result, err := d.queryStatement(ctx, fmt.Sprintf("SHOW CREATE TABLE %s", qualifiedTblName), driver.QueryOptions{})
	if err != nil {
		return "", err
	}
	if result == nil || len(result.Rows) != 1 || len(result.Rows[0]) != 1 {
		return "", fmt.Errorf("invalid ddl format")
	}
	return result.Rows[0][0].String, nil
}
