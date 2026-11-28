package analysis

import (
	"context"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/catalog"
)

var _ analysis.MetadataMigrator = Engine

// MetadataMigration generates the migration with catalog.MetadataMigration,
// each side resolving names on its own synced search_path.
func (engine) MetadataMigration(ctx context.Context, request analysis.MetadataMigrationRequest) (string, error) {
	return catalog.MetadataMigration(ctx, databaseMetadata(request.From), databaseMetadata(request.To), catalog.MetadataMigrationOptions{
		SourceSearchPath: migrationSearchPath(request.From),
		TargetSearchPath: migrationSearchPath(request.To),
		SkipSchemas:      request.SkipSchemas,
	})
}

func databaseMetadata(database *analysis.Database) *metadata.DatabaseSchemaMetadata {
	if database == nil {
		return nil
	}
	return database.Metadata
}

// migrationSearchPath returns the schemas the database's search_path setting
// names, in order: each entry unquoted (a doubled quote undone) or lower-cased,
// with $user and system schemas kept, since pg_catalog's position decides
// whether a user schema shadows a built-in. It returns nil when none remain.
func migrationSearchPath(database *analysis.Database) []string {
	var searchPath []string
	for _, token := range splitSearchPathSetting(databaseMetadata(database).GetSearchPath()) {
		token = strings.TrimSpace(token)
		schema := strings.ToLower(token)
		switch {
		case len(token) >= 2 && token[0] == '"' && token[len(token)-1] == '"':
			schema = strings.ReplaceAll(token[1:len(token)-1], `""`, `"`)
		case len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\'':
			schema = strings.ReplaceAll(token[1:len(token)-1], `''`, `'`)
		default:
		}
		if schema != "" {
			searchPath = append(searchPath, schema)
		}
	}
	return searchPath
}

// splitSearchPathSetting splits a search_path setting at the commas outside
// quotes.
func splitSearchPathSetting(setting string) []string {
	var parts []string
	start := 0
	var quote byte
	for i := 0; i < len(setting); i++ {
		switch {
		case quote != 0:
			if setting[i] != quote {
				continue
			}
			if i+1 < len(setting) && setting[i+1] == quote {
				i++
				continue
			}
			quote = 0
		case setting[i] == '"' || setting[i] == '\'':
			quote = setting[i]
		case setting[i] == ',':
			parts = append(parts, setting[start:i])
			start = i + 1
		default:
		}
	}
	return append(parts, setting[start:])
}
