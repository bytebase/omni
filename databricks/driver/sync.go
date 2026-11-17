package driver

import (
	"context"
	"fmt"
	"sort"
	"strings"

	catalogapi "github.com/databricks/databricks-sdk-go/service/catalog"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// tableUnion is one object a sync found, with the metadata kind it becomes.
type tableUnion struct {
	externalTable *metadatapb.ExternalTableMetadata
	table         *metadatapb.TableMetadata
	view          *metadatapb.ViewMetadata
	materialView  *metadatapb.MaterializedViewMetadata
	typeName      catalogapi.TableType
	name          string
}

type databricksSchemaMap = map[string][]*tableUnion
type databricksCatalogMap = map[string]databricksSchemaMap

// SyncInstance reads the SQL version and the workspace's catalog names. A
// catalog is the unit a caller syncs, so the databases are names only:
// SyncDatabase reads one catalog's schemas and tables. Roles stay empty because
// Unity Catalog exposes no role a sync reads, which is where Bytebase's TODO
// left it.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	result, err := d.queryStatement(ctx, "SELECT VERSION()", driver.QueryOptions{})
	if err != nil {
		return nil, wrapError("failed to read the instance version", err)
	}
	version, err := versionFromResult(result)
	if err != nil {
		return nil, err
	}
	catalogNames, err := d.listCatalogs(ctx)
	if err != nil {
		return nil, wrapError("failed to list catalogs", err)
	}

	instance := &driver.Instance{Version: version}
	for _, name := range catalogNames {
		instance.Databases = append(instance.Databases, &metadatapb.DatabaseSchemaMetadata{Name: name})
	}
	return instance, nil
}

// SyncDatabase reads the catalog the driver was opened on. An empty catalog
// answers nothing, as Bytebase's sync did.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	if d.catalog == "" {
		return nil, nil
	}
	catalogMap, err := d.listCatalogTables(ctx, d.catalog)
	if err != nil {
		return nil, wrapError("failed to list catalog tables", err)
	}
	schemaMap, ok := catalogMap[d.catalog]
	if !ok {
		return nil, fmt.Errorf("%w: catalog %q", driver.ErrNotFound, d.catalog)
	}
	return &metadatapb.DatabaseSchemaMetadata{
		Name:    d.catalog,
		Schemas: convertToStorepbSchemas(schemaMap),
	}, nil
}

// listCatalogs returns the workspace's catalog names in name order.
func (d *Driver) listCatalogs(ctx context.Context) ([]string, error) {
	catalogs, err := d.client.Catalogs.ListAll(ctx, catalogapi.ListCatalogsRequest{})
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(catalogs))
	for _, catalogInfo := range catalogs {
		names = append(names, catalogInfo.Name)
	}
	sort.Strings(names)
	return names, nil
}

// listCatalogTables reads the schemas and tables of targetCatalogName, or of
// every catalog when it is empty. Each catalog map holds every schema the
// catalog lists, including an empty one, so a schema with no tables still
// reaches the synced metadata.
func (d *Driver) listCatalogTables(ctx context.Context, targetCatalogName string) (databricksCatalogMap, error) {
	catalogMap := make(databricksCatalogMap)

	catalogsInfo, err := d.client.Catalogs.ListAll(ctx, catalogapi.ListCatalogsRequest{})
	if err != nil {
		return nil, err
	}
	sort.Slice(catalogsInfo, func(i, j int) bool { return catalogsInfo[i].Name < catalogsInfo[j].Name })

	for _, catalogInfo := range catalogsInfo {
		if targetCatalogName != "" && targetCatalogName != catalogInfo.Name {
			continue
		}
		schemasInfo, err := d.client.Schemas.ListAll(ctx, catalogapi.ListSchemasRequest{
			CatalogName: catalogInfo.Name,
		})
		if err != nil {
			return nil, err
		}
		sort.Slice(schemasInfo, func(i, j int) bool { return schemasInfo[i].Name < schemasInfo[j].Name })

		catalogMap[catalogInfo.Name] = make(databricksSchemaMap)
		for _, schemaInfo := range schemasInfo {
			catalogMap[catalogInfo.Name][schemaInfo.Name] = []*tableUnion{}

			tablesInfo, err := d.client.Tables.ListAll(ctx, catalogapi.ListTablesRequest{
				CatalogName: catalogInfo.Name,
				SchemaName:  schemaInfo.Name,
			})
			if err != nil {
				return nil, err
			}
			sort.Slice(tablesInfo, func(i, j int) bool { return tablesInfo[i].Name < tablesInfo[j].Name })
			if err := appendSchemaTables(catalogMap, tablesInfo); err != nil {
				return nil, err
			}
		}
	}

	return catalogMap, nil
}

// appendSchemaTables maps catalog.TableInfo answers into the metadata kinds a
// sync keeps. A streaming table, and any type Unity Catalog adds later, is
// skipped as Bytebase skipped it.
func appendSchemaTables(catalogMap databricksCatalogMap, tablesInfo []catalogapi.TableInfo) error {
	for _, tableInfo := range tablesInfo {
		table := tableUnion{
			typeName: tableInfo.TableType,
			name:     tableInfo.Name,
		}

		switch tableInfo.TableType {
		case catalogapi.TableTypeView:
			table.view = &metadatapb.ViewMetadata{
				Name:              table.name,
				Definition:        tableInfo.ViewDefinition,
				Comment:           tableInfo.Comment,
				DependencyColumns: convertToDependencyColumns(tableInfo.SchemaName, tableInfo.Name, tableInfo.Columns),
			}
		case catalogapi.TableTypeMaterializedView:
			table.materialView = &metadatapb.MaterializedViewMetadata{
				Name:       table.name,
				Definition: tableInfo.ViewDefinition,
				Comment:    tableInfo.Comment,
			}
		case catalogapi.TableTypeExternal:
			table.externalTable = &metadatapb.ExternalTableMetadata{
				Name: table.name,
			}
		case catalogapi.TableTypeManaged:
			table.table = &metadatapb.TableMetadata{
				Name:    table.name,
				Columns: convertToColumnMetadata(tableInfo.Columns),
				Comment: tableInfo.Comment,
			}
		default:
			// A streaming table, and any type added later, is not synced.
			continue
		}

		schemaMap, ok := catalogMap[tableInfo.CatalogName]
		if !ok {
			return fmt.Errorf("catalog %q was not initialized", tableInfo.CatalogName)
		}
		tables, ok := schemaMap[tableInfo.SchemaName]
		if !ok {
			return fmt.Errorf("schema %q.%q was not initialized", tableInfo.CatalogName, tableInfo.SchemaName)
		}
		catalogMap[tableInfo.CatalogName][tableInfo.SchemaName] = append(tables, &table)
	}
	return nil
}

// convertToColumnMetadata reads a table's columns in position order.
func convertToColumnMetadata(columnInfo []catalogapi.ColumnInfo) []*metadatapb.ColumnMetadata {
	columns := make([]*metadatapb.ColumnMetadata, 0, len(columnInfo))
	for _, col := range columnInfo {
		columns = append(columns, &metadatapb.ColumnMetadata{
			Name:     col.Name,
			Position: int32(col.Position),
			Nullable: col.Nullable,
			Type:     string(col.TypeName),
			Comment:  col.Comment,
		})
	}
	sort.SliceStable(columns, func(i, j int) bool {
		if columns[i].Position != columns[j].Position {
			return columns[i].Position < columns[j].Position
		}
		return columns[i].Name < columns[j].Name
	})
	return columns
}

// convertToDependencyColumns reads a view's columns as its dependencies, in
// name order so the synced metadata does not depend on the listing's order.
func convertToDependencyColumns(schema, table string, columnInfo []catalogapi.ColumnInfo) []*metadatapb.DependencyColumn {
	columns := make([]*metadatapb.DependencyColumn, 0, len(columnInfo))
	for _, col := range columnInfo {
		columns = append(columns, &metadatapb.DependencyColumn{
			Schema: schema,
			Table:  table,
			Column: col.Name,
		})
	}
	sort.Slice(columns, func(i, j int) bool { return columns[i].Column < columns[j].Column })
	return columns
}

// convertToStorepbSchemas turns one catalog's map into metadata, with schemas
// and each schema's tables in name order.
func convertToStorepbSchemas(schemaMap databricksSchemaMap) []*metadatapb.SchemaMetadata {
	schemaNames := make([]string, 0, len(schemaMap))
	for schemaName := range schemaMap {
		schemaNames = append(schemaNames, schemaName)
	}
	sort.Strings(schemaNames)

	schemas := make([]*metadatapb.SchemaMetadata, 0, len(schemaNames))
	for _, schemaName := range schemaNames {
		schemaMetadata := &metadatapb.SchemaMetadata{Name: schemaName}
		tables := append([]*tableUnion(nil), schemaMap[schemaName]...)
		sort.Slice(tables, func(i, j int) bool { return tables[i].name < tables[j].name })

		for _, table := range tables {
			switch table.typeName {
			case catalogapi.TableTypeView:
				schemaMetadata.Views = append(schemaMetadata.Views, table.view)
			case catalogapi.TableTypeMaterializedView:
				schemaMetadata.MaterializedViews = append(schemaMetadata.MaterializedViews, table.materialView)
			case catalogapi.TableTypeExternal:
				schemaMetadata.ExternalTables = append(schemaMetadata.ExternalTables, table.externalTable)
			case catalogapi.TableTypeManaged:
				schemaMetadata.Tables = append(schemaMetadata.Tables, table.table)
			default:
				// A streaming table, and any type added later, is not synced.
			}
		}
		schemas = append(schemas, schemaMetadata)
	}
	return schemas
}

// versionFromResult reads the release out of SELECT VERSION(). Databricks
// answers two words, the SQL release and its build, and Bytebase kept the
// first.
func versionFromResult(result *driver.Result) (string, error) {
	if result == nil || len(result.Rows) != 1 || len(result.Rows[0]) != 1 {
		return "", fmt.Errorf("invalid version result: expected one row with one value")
	}
	value := result.Rows[0][0]
	if value.Kind != driver.String {
		return "", fmt.Errorf("invalid version result: kind %s", value.Kind)
	}
	fields := strings.Fields(value.String)
	if len(fields) == 0 {
		return "", fmt.Errorf("invalid version result %q", value.String)
	}
	return fields[0], nil
}

// getQualifiedTblName quotes a table's parts for SHOW CREATE TABLE. Databricks
// escapes a backtick inside a part by doubling it.
func getQualifiedTblName(catalogName, schemaName, tblName string) (string, error) {
	if tblName == "" || schemaName == "" {
		return "", fmt.Errorf("table name and schema name must be specified")
	}
	qualifiedName := quoteBacktick(schemaName) + "." + quoteBacktick(tblName)
	if catalogName != "" {
		qualifiedName = quoteBacktick(catalogName) + "." + qualifiedName
	}
	return qualifiedName, nil
}

// quoteBacktick quotes one Databricks identifier.
func quoteBacktick(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
