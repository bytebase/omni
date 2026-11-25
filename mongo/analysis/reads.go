package analysis

import (
	"encoding/json"
	"fmt"

	"github.com/bytebase/omni/mongo/ast"
)

// A Namespace names a collection or a view. An empty Database is the database
// the statement runs in.
type Namespace struct {
	Database   string
	Collection string
}

// readMethods are the collection methods that return what they read. Their
// target is read, and an aggregate's pipeline may read more.
var readMethods = map[string]bool{
	"find":                   true,
	"findOne":                true,
	"aggregate":              true,
	"count":                  true,
	"countDocuments":         true,
	"estimatedDocumentCount": true,
	"distinct":               true,
}

// collectionReads returns what a read method on database.collection reads: the
// collection, and every collection its pipeline names. A write method reads
// nothing a caller is shown.
func collectionReads(database, collection, method string, args []ast.Node) []Namespace {
	if !readMethods[method] {
		return nil
	}
	reads := []Namespace{{Database: database, Collection: collection}}
	if method == "aggregate" && len(args) > 0 {
		reads = append(reads, pipelineReads(args[0], database)...)
	}
	return reads
}

// siblingReads returns what db.getSiblingDB("other").<collection>.<method>(...)
// reads, in the other database; cursor methods may follow. Any other chain
// reads nothing it can name.
func siblingReads(s *ast.DatabaseStatement) []Namespace {
	if s.Method != "getSiblingDB" || len(s.Args) != 1 || len(s.ChainedMethods) < 2 {
		return nil
	}
	database, ok := s.Args[0].(*ast.StringLiteral)
	if !ok {
		return nil
	}
	target, method := s.ChainedMethods[0], s.ChainedMethods[1]
	collection := target.Name
	if target.Name == "getCollection" {
		if len(target.Args) != 1 {
			return nil
		}
		name, ok := target.Args[0].(*ast.StringLiteral)
		if !ok {
			return nil
		}
		collection = name.Value
	} else if target.Loc.End-target.Loc.Start != len(target.Name) {
		// A property access ends at its name; a call such as getMongo() is
		// a method of the database, not a collection.
		return nil
	}
	return collectionReads(database.Value, collection, method.Name, method.Args)
}

// pipelineReads returns the collections the stages of pipeline read, in
// database: the from of $lookup and $graphLookup, the coll of $unionWith, and
// whatever the pipelines of $lookup, $unionWith and $facet read in turn. The
// target of $out or $merge is written, not read.
func pipelineReads(pipeline ast.Node, database string) []Namespace {
	type pending struct {
		pipeline ast.Node
		database string
	}
	var reads []Namespace
	stack := []pending{{pipeline, database}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		arr, ok := p.pipeline.(*ast.Array)
		if !ok {
			continue
		}
		for _, elem := range arr.Elements {
			stage, ok := elem.(*ast.Document)
			if !ok || len(stage.Pairs) != 1 {
				continue
			}
			spec := stage.Pairs[0].Value
			switch stage.Pairs[0].Key {
			case "$lookup", "$graphLookup":
				doc, ok := spec.(*ast.Document)
				if !ok {
					continue
				}
				// A $lookup's pipeline runs on the collection it joins, so
				// it names collections of that collection's database.
				inner := p.database
				for _, kv := range doc.Pairs {
					if kv.Key != "from" {
						continue
					}
					if ns, ok := fromNamespace(kv.Value, p.database); ok {
						reads = append(reads, ns)
						inner = ns.Database
					}
				}
				if stage.Pairs[0].Key == "$lookup" {
					if sub := field(doc, "pipeline"); sub != nil {
						stack = append(stack, pending{sub, inner})
					}
				}
			case "$unionWith":
				switch v := spec.(type) {
				case *ast.StringLiteral:
					reads = append(reads, Namespace{Database: p.database, Collection: v.Value})
				case *ast.Document:
					if coll, ok := field(v, "coll").(*ast.StringLiteral); ok {
						reads = append(reads, Namespace{Database: p.database, Collection: coll.Value})
					}
					if sub := field(v, "pipeline"); sub != nil {
						stack = append(stack, pending{sub, p.database})
					}
				}
			case "$facet":
				doc, ok := spec.(*ast.Document)
				if !ok {
					continue
				}
				for _, kv := range doc.Pairs {
					stack = append(stack, pending{kv.Value, p.database})
				}
			}
		}
	}
	return reads
}

// fromNamespace reads a $lookup or $graphLookup from: a collection of database,
// or {db: ..., coll: ...}.
func fromNamespace(from ast.Node, database string) (Namespace, bool) {
	switch v := from.(type) {
	case *ast.StringLiteral:
		return Namespace{Database: database, Collection: v.Value}, true
	case *ast.Document:
		db, dbOK := field(v, "db").(*ast.StringLiteral)
		coll, collOK := field(v, "coll").(*ast.StringLiteral)
		if dbOK && collOK {
			return Namespace{Database: db.Value, Collection: coll.Value}, true
		}
	}
	return Namespace{}, false
}

// field returns the value of doc's last key, which is the one MongoDB reads.
func field(doc *ast.Document, key string) ast.Node {
	var value ast.Node
	for _, kv := range doc.Pairs {
		if kv.Key == key {
			value = kv.Value
		}
	}
	return value
}

// ExpandViews returns reads followed by what each view among them reads, and
// what each view among those reads, until only collections are left. Each
// namespace appears once, the views included, so a policy on a view still
// applies to a query over it.
//
// definition returns a view's options as Bytebase's sync stores them
// (metadata.ViewMetadata.Definition): the listCollections options, in
// Extended JSON, {"viewOn": ..., "pipeline": [...]}. It reports false for a
// collection. A view reads the collection or view it is on, in its own
// database, and whatever its pipeline reads.
func ExpandViews(reads []Namespace, definition func(Namespace) (string, bool)) ([]Namespace, error) {
	seen := map[Namespace]bool{}
	var out []Namespace
	queue := append([]Namespace(nil), reads...)
	for len(queue) > 0 {
		ns := queue[0]
		queue = queue[1:]
		if seen[ns] {
			continue
		}
		seen[ns] = true
		out = append(out, ns)
		text, ok := definition(ns)
		if !ok {
			continue
		}
		var view struct {
			ViewOn   string `json:"viewOn"`
			Pipeline any    `json:"pipeline"`
		}
		if err := json.Unmarshal([]byte(text), &view); err != nil {
			return nil, fmt.Errorf("view %s.%s: unreadable definition: %w", ns.Database, ns.Collection, err)
		}
		if view.ViewOn == "" {
			return nil, fmt.Errorf("view %s.%s: a definition without viewOn", ns.Database, ns.Collection)
		}
		queue = append(queue, Namespace{Database: ns.Database, Collection: view.ViewOn})
		queue = append(queue, pipelineReads(jsonNode(view.Pipeline), ns.Database)...)
	}
	return out, nil
}

// jsonNode returns the shape of a decoded JSON value that pipelineReads reads:
// documents, arrays and strings. Any other value is nil, which no stage reads
// a collection name from.
func jsonNode(v any) ast.Node {
	type pending struct {
		value any
		set   func(ast.Node)
	}
	var root ast.Node
	stack := []pending{{v, func(n ast.Node) { root = n }}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch x := p.value.(type) {
		case string:
			p.set(&ast.StringLiteral{Value: x})
		case []any:
			arr := &ast.Array{Elements: make([]ast.Node, len(x))}
			p.set(arr)
			for i, e := range x {
				stack = append(stack, pending{e, func(n ast.Node) { arr.Elements[i] = n }})
			}
		case map[string]any:
			doc := &ast.Document{Pairs: make([]ast.KeyValue, 0, len(x))}
			p.set(doc)
			for key, e := range x {
				doc.Pairs = append(doc.Pairs, ast.KeyValue{Key: key})
				i := len(doc.Pairs) - 1
				stack = append(stack, pending{e, func(n ast.Node) { doc.Pairs[i].Value = n }})
			}
		}
	}
	return root
}
