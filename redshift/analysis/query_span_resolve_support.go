package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
)

func (q *querySpanResolver) getDatabaseMetadata(database string) (*metadata.DatabaseSchemaMetadata, error) {
	if database == "" {
		database = q.defaultDatabase
	}
	if q.opts.GetDatabaseMetadata == nil {
		return nil, nil
	}
	dbMeta, err := q.opts.GetDatabaseMetadata(database)
	if err != nil {
		return nil, fmt.Errorf(redshiftGetDatabaseMetadataError+": %w", database, err)
	}
	return dbMeta, nil
}

func schemaMetadata(meta *metadata.DatabaseSchemaMetadata, name string) *metadata.SchemaMetadata {
	for _, schema := range meta.GetSchemas() {
		if schema.GetName() == name {
			return schema
		}
	}
	return nil
}

func schemaHasRelation(schema *metadata.SchemaMetadata, name string) bool {
	for _, table := range schema.GetTables() {
		if table.GetName() == name || tableHasPartition(table, name) {
			return true
		}
	}
	for _, view := range schema.GetViews() {
		if view.GetName() == name {
			return true
		}
	}
	for _, view := range schema.GetMaterializedViews() {
		if view.GetName() == name {
			return true
		}
	}
	for _, table := range schema.GetExternalTables() {
		if table.GetName() == name {
			return true
		}
	}
	for _, sequence := range schema.GetSequences() {
		if sequence.GetName() == name {
			return true
		}
	}
	return false
}

func tableHasPartition(table *metadata.TableMetadata, name string) bool {
	for _, partition := range table.GetPartitions() {
		if partition.GetName() == name || partitionHasSubpartition(partition, name) {
			return true
		}
	}
	return false
}

func partitionHasSubpartition(partition *metadata.TablePartitionMetadata, name string) bool {
	for _, sub := range partition.GetSubpartitions() {
		if sub.GetName() == name || partitionHasSubpartition(sub, name) {
			return true
		}
	}
	return false
}

func searchRelation(meta *metadata.DatabaseSchemaMetadata, searchPath []string, name string) (string, string) {
	for _, schemaName := range searchPath {
		schema := schemaMetadata(meta, schemaName)
		if schema == nil {
			continue
		}
		if schemaHasRelation(schema, name) {
			return schema.GetName(), name
		}
	}
	return "", ""
}

func searchFunctions(meta *metadata.DatabaseSchemaMetadata, searchPath []string, name string) ([]string, []*metadata.FunctionMetadata) {
	var schemas []string
	var funcs []*metadata.FunctionMetadata
	for _, schemaName := range searchPath {
		schema := schemaMetadata(meta, schemaName)
		if schema == nil {
			continue
		}
		for _, function := range schema.GetFunctions() {
			if function.GetName() == name {
				schemas = append(schemas, schema.GetName())
				funcs = append(funcs, function)
			}
		}
	}
	return schemas, funcs
}

func relationColumnNames(table *metadata.TableMetadata) []string {
	names := make([]string, 0, len(table.GetColumns()))
	for _, column := range table.GetColumns() {
		names = append(names, column.GetName())
	}
	return names
}

func viewColumnNames(view *metadata.ViewMetadata) []string {
	names := make([]string, 0, len(view.GetColumns()))
	for _, column := range view.GetColumns() {
		names = append(names, column.GetName())
	}
	return names
}

func externalTableColumnNames(table *metadata.ExternalTableMetadata) []string {
	names := make([]string, 0, len(table.GetColumns()))
	for _, column := range table.GetColumns() {
		names = append(names, column.GetName())
	}
	return names
}

func normalizedSchemaName(meta *metadata.DatabaseSchemaMetadata, name string) *metadata.SchemaMetadata {
	schema := schemaMetadata(meta, name)
	if schema != nil {
		return schema
	}
	for _, candidate := range meta.GetSchemas() {
		if strings.EqualFold(candidate.GetName(), name) {
			return candidate
		}
	}
	return nil
}
