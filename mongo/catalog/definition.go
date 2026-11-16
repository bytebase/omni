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
// A collection is created with the options its metadata holds, which a synced
// time-series collection and a capped or validated collection need. A view is
// created from its definition, the options listCollections reports. An index
// is created from its key plus its definition, the options listIndexes reports
// beyond the key and the name and uniqueness the metadata also carries.
func DatabaseDefinition(meta *metadata.DatabaseSchemaMetadata) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "// MongoDB schema dump of database %s\n", strconv.Quote(meta.GetName()))
	for _, schema := range meta.GetSchemas() {
		for _, table := range schema.GetTables() {
			if err := writeCollection(&b, table); err != nil {
				return "", err
			}
			for _, index := range table.GetIndexes() {
				if err := writeIndex(&b, table.GetName(), index); err != nil {
					return "", err
				}
			}
		}
		for _, view := range schema.GetViews() {
			if err := writeView(&b, view); err != nil {
				return "", err
			}
		}
	}
	return b.String(), nil
}

// writeCollection writes the statement that creates table with its options.
func writeCollection(b *strings.Builder, table *metadata.TableMetadata) error {
	options := table.GetCreateOptions()
	if options == "" {
		fmt.Fprintf(b, "db.createCollection(%s);\n", strconv.Quote(table.GetName()))
		return nil
	}
	literal, err := renderExtendedJSON(options)
	if err != nil {
		return fmt.Errorf("collection %q options: %w", table.GetName(), err)
	}
	if literal == "{}" {
		fmt.Fprintf(b, "db.createCollection(%s);\n", strconv.Quote(table.GetName()))
		return nil
	}
	body, err := objectBody(literal)
	if err != nil {
		return fmt.Errorf("collection %q options: %w", table.GetName(), err)
	}
	// The create command takes the collection's name beside the options; the
	// server's answer holds the options alone.
	fmt.Fprintf(b, "db.runCommand({\"create\": %s, %s});\n", strconv.Quote(table.GetName()), body)
	return nil
}

// writeView writes the statement that creates view. Metadata synced before a
// view carried its definition still has none, and a view without one cannot be
// recreated, so the dump names it in a comment instead of inventing a
// definition.
func writeView(b *strings.Builder, view *metadata.ViewMetadata) error {
	definition := view.GetDefinition()
	if definition == "" {
		fmt.Fprintf(b, "// View %s is not restorable: the synced metadata carries no definition.\n", strconv.Quote(view.GetName()))
		return nil
	}
	literal, err := renderExtendedJSON(definition)
	if err != nil {
		return fmt.Errorf("view %q definition: %w", view.GetName(), err)
	}
	if literal == "{}" {
		fmt.Fprintf(b, "// View %s is not restorable: the synced metadata carries no definition.\n", strconv.Quote(view.GetName()))
		return nil
	}
	body, err := objectBody(literal)
	if err != nil {
		return fmt.Errorf("view %q definition: %w", view.GetName(), err)
	}
	fmt.Fprintf(b, "db.runCommand({\"create\": %s, %s});\n", strconv.Quote(view.GetName()), body)
	return nil
}

// writeIndex writes the statement that creates index on collection. The
// createIndexes command carries every option listIndexes reports; mongosh's
// createIndex wrapper, and the subset of it gomongo executes, does not.
func writeIndex(b *strings.Builder, collection string, index *metadata.IndexMetadata) error {
	if len(index.GetExpressions()) != 1 {
		return fmt.Errorf("index %q on collection %q has %d expressions; an index's key is one document",
			index.GetName(), collection, len(index.GetExpressions()))
	}
	key := index.GetExpressions()[0]
	if !json.Valid([]byte(key)) {
		return fmt.Errorf("index %q on collection %q has a key that is not a document: %s", index.GetName(), collection, key)
	}
	fields := []string{`"key": ` + key, `"name": ` + strconv.Quote(index.GetName())}
	if index.GetUnique() {
		fields = append(fields, `"unique": true`)
	}
	if definition := index.GetDefinition(); definition != "" {
		literal, err := renderExtendedJSON(definition)
		if err != nil {
			return fmt.Errorf("index %q on collection %q: %w", index.GetName(), collection, err)
		}
		body, err := objectBody(literal)
		if err != nil {
			return fmt.Errorf("index %q on collection %q: %w", index.GetName(), collection, err)
		}
		if body != "" {
			fields = append(fields, body)
		}
	}
	fmt.Fprintf(b, "db.runCommand({\"createIndexes\": %s, \"indexes\": [{%s}]});\n",
		strconv.Quote(collection), strings.Join(fields, ", "))
	return nil
}

// objectBody returns the fields of a rendered object literal, without its
// braces, so a caller can put them in another object.
func objectBody(literal string) (string, error) {
	if len(literal) < 2 || literal[0] != '{' || literal[len(literal)-1] != '}' {
		return "", fmt.Errorf("the options are not a document: %s", literal)
	}
	return literal[1 : len(literal)-1], nil
}
