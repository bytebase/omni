package analysis

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/mongo/ast"
)

// extractPredicateFields returns a sorted, deduplicated slice of field paths
// found in the first argument (the query filter document) of a find/findOne call.
func extractPredicateFields(args []ast.Node) []string {
	if len(args) == 0 {
		return nil
	}
	doc, ok := args[0].(*ast.Document)
	if !ok {
		return nil
	}
	fields := make(map[string]struct{})
	collectFromDocument(doc, "", fields)
	if len(fields) == 0 {
		return nil
	}
	return sortedKeys(fields)
}

// collectFromDocument walks a document and collects field paths into fields.
// For nested documents like {contact: {phone: "123"}}, both "contact" and
// "contact.phone" are emitted as predicate fields. This intentional
// conservative over-reporting ensures all potentially accessed sub-paths are
// included, which is safer for masking decisions.
func collectFromDocument(doc *ast.Document, prefix string, fields map[string]struct{}) {
	for _, kv := range doc.Pairs {
		key := kv.Key
		if strings.HasPrefix(key, "$") {
			if isLogicalOperator(key) {
				collectFromLogicalOp(kv.Value, prefix, fields)
			}
			// other operators (like $expr) are skipped
			continue
		}
		path := joinPath(prefix, key)
		fields[path] = struct{}{}
		collectFromValue(kv.Value, path, fields)
	}
}

// collectFromValue recurses into a value node to collect nested field paths.
func collectFromValue(node ast.Node, prefix string, fields map[string]struct{}) {
	switch v := node.(type) {
	case *ast.Document:
		collectFromDocument(v, prefix, fields)
	case *ast.Array:
		for _, elem := range v.Elements {
			if d, ok := elem.(*ast.Document); ok {
				collectFromDocument(d, prefix, fields)
			}
		}
	}
}

// collectFromLogicalOp handles $and/$or/$nor operator values.
func collectFromLogicalOp(node ast.Node, prefix string, fields map[string]struct{}) {
	collectFromValue(node, prefix, fields)
}

// isLogicalOperator returns true for $and, $or, $nor.
func isLogicalOperator(key string) bool {
	return key == "$and" || key == "$or" || key == "$nor"
}

// joinPath joins prefix and key with a dot, omitting the dot when prefix is empty.
func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// sortedKeys returns sorted keys from a string set map.
func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
