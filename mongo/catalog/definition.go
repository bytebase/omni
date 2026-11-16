package catalog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// DatabaseDefinition renders meta as the mongosh script that recreates its
// collections and their indexes. It renders the metadata it is given and
// nothing else.
//
// A MongoDB view is not restorable: a synced view's metadata carries its name
// and not the pipeline it runs, so the script names it in a comment. Index
// options beyond a name and uniqueness are not restorable either, for the same
// reason: listIndexes reports them and the metadata does not keep them.
func DatabaseDefinition(meta *metadata.DatabaseSchemaMetadata) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "// MongoDB schema dump of database %s\n", strconv.Quote(meta.GetName()))
	for _, schema := range meta.GetSchemas() {
		for _, table := range schema.GetTables() {
			fmt.Fprintf(&b, "db.createCollection(%s);\n", strconv.Quote(table.GetName()))
			for _, index := range table.GetIndexes() {
				if len(index.GetExpressions()) != 1 {
					return "", fmt.Errorf("index %q on collection %q has %d expressions; an index's key is one document",
						index.GetName(), table.GetName(), len(index.GetExpressions()))
				}
				key := index.GetExpressions()[0]
				if !json.Valid([]byte(key)) {
					return "", fmt.Errorf("index %q on collection %q has a key that is not a document: %s", index.GetName(), table.GetName(), key)
				}
				fmt.Fprintf(&b, "db.getCollection(%s).createIndex(%s, %s);\n",
					strconv.Quote(table.GetName()), key, indexOptions(index))
			}
		}
		for _, view := range schema.GetViews() {
			fmt.Fprintf(&b, "// View %s is not restorable: the synced metadata carries no pipeline.\n", strconv.Quote(view.GetName()))
		}
	}
	return b.String(), nil
}

// indexOptions renders the options createIndex takes for index: its name, and
// its uniqueness when it is unique.
func indexOptions(index *metadata.IndexMetadata) string {
	options := `{"name": ` + strconv.Quote(index.GetName())
	if index.GetUnique() {
		options += `, "unique": true`
	}
	return options + "}"
}
