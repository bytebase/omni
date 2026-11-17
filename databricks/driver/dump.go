package driver

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	catalogapi "github.com/databricks/databricks-sdk-go/service/catalog"

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
// SHOW CREATE TABLE states a managed or external table and a view, and DESC
// FORMATTED supplies a materialized view's table properties. The system catalog
// and each catalog's information_schema are left out, as Bytebase left them out.
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
		if _, err := fmt.Fprintf(writer, "CREATE CATALOG %s\n", quoteBacktick(catalogName)); err != nil {
			return err
		}

		schemaMap := catalogMap[catalogName]
		schemaNames := make([]string, 0, len(schemaMap))
		for schemaName := range schemaMap {
			schemaNames = append(schemaNames, schemaName)
		}
		sort.Strings(schemaNames)

		for _, schemaName := range schemaNames {
			if schemaName == infoSchema {
				continue
			}
			if schemaName != dftSchema {
				if _, err := fmt.Fprintf(writer, "CREATE SCHEMA %s.%s\n", quoteBacktick(catalogName), quoteBacktick(schemaName)); err != nil {
					return err
				}
			}

			viewDDL := strings.Builder{}
			mtViewDDL := strings.Builder{}
			extTblDDL := strings.Builder{}
			mngTblDDL := strings.Builder{}
			tables := append([]*tableUnion(nil), schemaMap[schemaName]...)
			sort.Slice(tables, func(i, j int) bool { return tables[i].name < tables[j].name })
			for _, tblUnion := range tables {
				qualifiedName, err := getQualifiedTblName(catalogName, schemaName, tblUnion.name)
				if err != nil {
					return err
				}

				switch tblUnion.typeName {
				case catalogapi.TableTypeView:
					ddl, err := d.showCreateTable(ctx, qualifiedName)
					if err != nil {
						return err
					}
					formatDDL := fmt.Sprintf(schemaStmtFmt, "VIEW", qualifiedName, ddl)
					if _, err := viewDDL.WriteString(formatDDL); err != nil {
						return err
					}
				case catalogapi.TableTypeMaterializedView:
					tblAddInfo, err := d.descTbl(ctx, qualifiedName)
					if err != nil {
						return err
					}
					tblProperties := formatTblProperties(tblAddInfo)
					formatDDL := genMaterializedViewDDL(qualifiedName, tblUnion.materialView, tblProperties)
					if _, err := mtViewDDL.WriteString(formatDDL); err != nil {
						return err
					}
				case catalogapi.TableTypeExternal:
					ddl, err := d.showCreateTable(ctx, qualifiedName)
					if err != nil {
						return err
					}
					formatDDL := fmt.Sprintf(schemaStmtFmt, "EXTERNAL TABLE", qualifiedName, ddl)
					if _, err := extTblDDL.WriteString(formatDDL); err != nil {
						return err
					}
				case catalogapi.TableTypeManaged:
					ddl, err := d.showCreateTable(ctx, qualifiedName)
					if err != nil {
						return err
					}
					formatDDL := fmt.Sprintf(schemaStmtFmt, "MANAGED TABLE", qualifiedName, ddl)
					if _, err := mngTblDDL.WriteString(formatDDL); err != nil {
						return err
					}
				default:
					// A streaming table, and any type added later, is not synced.
					continue
				}
			}
			if _, err := writer.Write([]byte(mngTblDDL.String())); err != nil {
				return err
			}
			if _, err := writer.Write([]byte(extTblDDL.String())); err != nil {
				return err
			}
			if _, err := writer.Write([]byte(viewDDL.String())); err != nil {
				return err
			}
			if _, err := writer.Write([]byte(mtViewDDL.String())); err != nil {
				return err
			}
		}
	}

	return nil
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

// genMaterializedViewDDL writes a materialized view from the metadata the sync
// read. It is Bytebase's shape, including its unquoted comment and property
// list; VERIFY.md marks the shape as unconfirmed until a workspace states one.
func genMaterializedViewDDL(qualifiedName string, mtViewMeta *metadatapb.MaterializedViewMetadata, tblProperties string) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "CREATE MATERIALIZED VIEW %s\n", qualifiedName)
	fmt.Fprintf(&builder, "COMMENT %s\n", mtViewMeta.Comment)
	fmt.Fprintf(&builder, "TBLPROPERTIES %s\n", tblProperties)
	return builder.String()
}

type tableAdditionalInfo struct {
	tblProps map[string]string
}

// descTbl reads the table properties DESC FORMATTED states for a materialized
// view.
func (d *Driver) descTbl(ctx context.Context, tblName string) (*tableAdditionalInfo, error) {
	tblAddInfo := &tableAdditionalInfo{
		tblProps: make(map[string]string),
	}
	result, err := d.queryStatement(ctx, fmt.Sprintf("DESC FORMATTED %s", tblName), driver.QueryOptions{})
	if err != nil {
		return nil, err
	}

	for _, row := range result.Rows {
		if len(row) < 2 {
			continue
		}
		if row[0].String != "Table Properties" {
			continue
		}
		for key, value := range parseTableProperties(row[1].String) {
			tblAddInfo.tblProps[key] = value
		}
	}
	return tblAddInfo, nil
}

// tablePropertyPattern is the one Bytebase read: a property key and value of
// letters, digits, hyphens and dots.
var tablePropertyPattern = regexp.MustCompile(`([a-zA-Z0-9-.]*)=([a-zA-Z0-9-.]*)`)

// parseTableProperties reads the key=value pairs out of DESC FORMATTED's table
// properties cell.
func parseTableProperties(value string) map[string]string {
	properties := make(map[string]string)
	for _, match := range tablePropertyPattern.FindAllStringSubmatch(value, -1) {
		properties[match[1]] = match[2]
	}
	return properties
}

// formatTblProperties writes a materialized view's table properties in key
// order, so two dumps of one schema are byte-identical.
func formatTblProperties(tblAddInfo *tableAdditionalInfo) string {
	if tblAddInfo == nil || tblAddInfo.tblProps == nil {
		return ""
	}

	keys := make([]string, 0, len(tblAddInfo.tblProps))
	for key := range tblAddInfo.tblProps {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	builder := strings.Builder{}
	for i, key := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&builder, "%s = %s%s\n", key, tblAddInfo.tblProps[key], comma)
	}
	return builder.String()
}
