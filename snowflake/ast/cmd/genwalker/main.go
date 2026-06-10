// Command genwalker generates walk_generated.go from parsenodes.go and node.go.
//
// It scans all struct types in those files, identifies fields that carry child
// nodes, and generates the walkChildren function that enumerates child nodes
// for Walk.
//
// A field carries child nodes when its type is one of:
//
//   - Node, []Node, [][]Node — the Node interface and slices of it
//   - *T, []*T, T, []T       — where T is a struct implementing Node
//     (has a Tag() method)
//   - *H, []*H, H, []H       — where H is a NON-Node helper struct
//     (e.g. WhenClause, OrderItem, WindowSpec) that
//     transitively holds node-carrying fields
//
// For each node-bearing helper struct a dedicated walk<H> function is
// generated that walks the helper's own node-bearing fields (recursing into
// nested helpers), so expressions stored inside helpers — CASE WHEN
// conditions, OVER (PARTITION BY ... ORDER BY ...) expressions, CTE bodies,
// SELECT-list targets, and the like — are reachable from Walk / Inspect.
//
// If a field's type mentions Node, a Node struct, or a node-bearing helper in
// a shape the generator does not understand (e.g. a map or [][]*T), generation
// fails loudly rather than silently skipping the field.
//
// Usage:
//
//	go run ./snowflake/ast/cmd/genwalker
//
// or via go generate:
//
//	go generate ./snowflake/ast/...
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

// shape classifies how a field's type carries child nodes.
type shape int

const (
	shapeNone             shape = iota // does not carry nodes
	shapeNodeIface                     // Node
	shapeNodeIfaceSlice                // []Node
	shapeNodeIfaceRows                 // [][]Node
	shapeNodePtr                       // *T, T a Node struct
	shapeNodePtrSlice                  // []*T, T a Node struct
	shapeNodeValue                     // T, T a Node struct (by value)
	shapeNodeValueSlice                // []T, T a Node struct (by value)
	shapeHelperPtr                     // *H, H a node-bearing helper struct
	shapeHelperPtrSlice                // []*H, H a node-bearing helper struct
	shapeHelperValue                   // H, H a node-bearing helper struct (by value)
	shapeHelperValueSlice              // []H, H a node-bearing helper struct (by value)
)

type field struct {
	Name string
	Type string // "Node", "[]Node", "*SelectStmt", "[]*WhenClause", etc.
}

type structInfo struct {
	Name   string
	Fields []field
}

func main() {
	fset := token.NewFileSet()

	// Parse both source files. The generator runs from the snowflake/ast
	// directory thanks to //go:generate go run ./cmd/genwalker in walk.go,
	// so the relative paths point at sibling files.
	sources := []string{"parsenodes.go", "node.go"}
	var files []*ast.File
	for _, src := range sources {
		f, err := parser.ParseFile(fset, src, nil, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", src, err)
			os.Exit(1)
		}
		files = append(files, f)
	}

	// nodeStructs tracks which structs implement the Node interface (have a
	// Tag() method). Only these structs get their own case in the walkChildren
	// switch.
	nodeStructs := map[string]bool{}
	// structFields records every struct's full field list (all fields, not
	// just node-bearing ones) so helper-struct reachability can be computed.
	structFields := map[string][]field{}

	for _, f := range files {
		// Scan function declarations for Tag() methods to identify Node types.
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || fd.Name.Name != "Tag" {
				continue
			}
			// Check it returns a single result and has no parameters.
			if fd.Type.Params != nil && len(fd.Type.Params.List) > 0 {
				continue
			}
			if fd.Type.Results == nil || len(fd.Type.Results.List) != 1 {
				continue
			}
			// Extract the receiver type name (handles both *T and T receivers).
			for _, recv := range fd.Recv.List {
				recvName := ""
				switch rt := recv.Type.(type) {
				case *ast.StarExpr:
					if id, ok := rt.X.(*ast.Ident); ok {
						recvName = id.Name
					}
				case *ast.Ident:
					recvName = rt.Name
				}
				if recvName != "" {
					nodeStructs[recvName] = true
				}
			}
		}
		// Collect fields for each struct.
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts := spec.(*ast.TypeSpec)
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				var fields []field
				for _, fld := range st.Fields.List {
					if len(fld.Names) == 0 {
						continue // embedded
					}
					typStr := typeString(fld.Type)
					for _, name := range fld.Names {
						fields = append(fields, field{Name: name.Name, Type: typStr})
					}
				}
				structFields[ts.Name.Name] = fields
			}
		}
	}

	// Compute the set of node-bearing helper structs: non-Node structs that
	// transitively hold node-carrying fields. Fixpoint iteration handles
	// helper-inside-helper nesting (e.g. WindowSpec → WindowFrame →
	// WindowBound).
	bearingHelpers := map[string]bool{}
	for changed := true; changed; {
		changed = false
		for name, fields := range structFields {
			if nodeStructs[name] || bearingHelpers[name] || name == "Loc" {
				continue
			}
			for _, f := range fields {
				if classify(f.Type, nodeStructs, bearingHelpers) != shapeNone {
					bearingHelpers[name] = true
					changed = true
					break
				}
			}
		}
	}

	// Guard: any field whose type mentions Node, a Node struct, or a
	// node-bearing helper must classify to a known shape — otherwise the
	// walker would silently skip child nodes. Fail generation instead.
	var unsupported []string
	for name, fields := range structFields {
		if !nodeStructs[name] && !bearingHelpers[name] {
			continue
		}
		for _, f := range fields {
			if classify(f.Type, nodeStructs, bearingHelpers) != shapeNone {
				continue
			}
			b := baseType(f.Type)
			if b == "Node" || nodeStructs[b] || bearingHelpers[b] {
				unsupported = append(unsupported, fmt.Sprintf("%s.%s %s", name, f.Name, f.Type))
			}
		}
	}
	if len(unsupported) > 0 {
		sort.Strings(unsupported)
		fmt.Fprintf(os.Stderr, "genwalker: node-bearing fields with unsupported type shapes:\n")
		for _, u := range unsupported {
			fmt.Fprintf(os.Stderr, "  %s\n", u)
		}
		os.Exit(1)
	}

	// Sort node structs by name for deterministic output.
	var names []string
	for n := range structFields {
		names = append(names, n)
	}
	sort.Strings(names)

	// Generate code.
	var buf bytes.Buffer
	buf.WriteString("// Code generated by genwalker; DO NOT EDIT.\n\n")
	buf.WriteString("package ast\n\n")
	buf.WriteString("// walkChildren walks the child nodes of node, calling Walk(v, child)\n")
	buf.WriteString("// for each child. This function is generated from parsenodes.go and node.go.\n")
	buf.WriteString("func walkChildren(v Visitor, node Node) {\n")
	buf.WriteString("\tswitch n := node.(type) {\n")

	cases := 0
	caseFields := 0
	for _, s := range names {
		// Only emit cases for structs that implement Node (have a Tag() method).
		if !nodeStructs[s] {
			continue
		}
		var walkable []field
		for _, f := range structFields[s] {
			if classify(f.Type, nodeStructs, bearingHelpers) != shapeNone {
				walkable = append(walkable, f)
			}
		}
		if len(walkable) == 0 {
			continue
		}
		cases++
		caseFields += len(walkable)
		fmt.Fprintf(&buf, "\tcase *%s:\n", s)
		for _, f := range walkable {
			emitFieldWalk(&buf, "\t\t", "n."+f.Name, f.Type, nodeStructs, bearingHelpers)
		}
	}

	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	// Emit one walk function per node-bearing helper struct, sorted by name.
	helperCount := 0
	for _, s := range names {
		if !bearingHelpers[s] {
			continue
		}
		helperCount++
		fmt.Fprintf(&buf, "\n// walk%s walks the node-bearing fields of the non-Node helper struct\n", s)
		fmt.Fprintf(&buf, "// %s. The helper itself is not visited (it is not a Node); its child\n", s)
		buf.WriteString("// nodes are walked in field order.\n")
		fmt.Fprintf(&buf, "func walk%s(v Visitor, n *%s) {\n", s, s)
		buf.WriteString("\tif n == nil {\n\t\treturn\n\t}\n")
		for _, f := range structFields[s] {
			if classify(f.Type, nodeStructs, bearingHelpers) == shapeNone {
				continue
			}
			emitFieldWalk(&buf, "\t", "n."+f.Name, f.Type, nodeStructs, bearingHelpers)
		}
		buf.WriteString("}\n")
	}

	// Format with gofmt rules.
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "format: %v\n", err)
		// Write unformatted output for debugging.
		_ = os.WriteFile("walk_generated.go", buf.Bytes(), 0644)
		os.Exit(1)
	}

	if err := os.WriteFile("walk_generated.go", formatted, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated walk_generated.go: %d cases, %d child fields, %d helper walkers\n",
		cases, caseFields, helperCount)
}

// emitFieldWalk writes the walk statement(s) for one node-bearing field.
// expr is the field access expression (e.g. "n.Whens"), indent is the leading
// whitespace for the first line.
func emitFieldWalk(buf *bytes.Buffer, indent, expr, typStr string, nodeStructs, bearingHelpers map[string]bool) {
	switch classify(typStr, nodeStructs, bearingHelpers) {
	case shapeNodeIface:
		fmt.Fprintf(buf, "%sWalk(v, %s)\n", indent, expr)
	case shapeNodeIfaceSlice:
		fmt.Fprintf(buf, "%swalkNodes(v, %s)\n", indent, expr)
	case shapeNodeIfaceRows:
		fmt.Fprintf(buf, "%swalkNodeRows(v, %s)\n", indent, expr)
	case shapeNodePtr:
		fmt.Fprintf(buf, "%sif %s != nil {\n", indent, expr)
		fmt.Fprintf(buf, "%s\tWalk(v, %s)\n", indent, expr)
		fmt.Fprintf(buf, "%s}\n", indent)
	case shapeNodePtrSlice:
		fmt.Fprintf(buf, "%sfor _, item := range %s {\n", indent, expr)
		fmt.Fprintf(buf, "%s\tif item != nil {\n", indent)
		fmt.Fprintf(buf, "%s\t\tWalk(v, item)\n", indent)
		fmt.Fprintf(buf, "%s\t}\n", indent)
		fmt.Fprintf(buf, "%s}\n", indent)
	case shapeNodeValue:
		fmt.Fprintf(buf, "%sWalk(v, &%s)\n", indent, expr)
	case shapeNodeValueSlice:
		fmt.Fprintf(buf, "%sfor i := range %s {\n", indent, expr)
		fmt.Fprintf(buf, "%s\tWalk(v, &%s[i])\n", indent, expr)
		fmt.Fprintf(buf, "%s}\n", indent)
	case shapeHelperPtr:
		fmt.Fprintf(buf, "%swalk%s(v, %s)\n", indent, baseType(typStr), expr)
	case shapeHelperPtrSlice:
		fmt.Fprintf(buf, "%sfor _, item := range %s {\n", indent, expr)
		fmt.Fprintf(buf, "%s\twalk%s(v, item)\n", indent, baseType(typStr))
		fmt.Fprintf(buf, "%s}\n", indent)
	case shapeHelperValue:
		fmt.Fprintf(buf, "%swalk%s(v, &%s)\n", indent, baseType(typStr), expr)
	case shapeHelperValueSlice:
		fmt.Fprintf(buf, "%sfor i := range %s {\n", indent, expr)
		fmt.Fprintf(buf, "%s\twalk%s(v, &%s[i])\n", indent, baseType(typStr), expr)
		fmt.Fprintf(buf, "%s}\n", indent)
	}
}

// classify maps a field type string to the shape of child-node traversal it
// needs. Returns shapeNone for types that carry no nodes (enums, Ident, Loc,
// scalars, inert helper structs).
func classify(typStr string, nodeStructs, bearingHelpers map[string]bool) shape {
	switch typStr {
	case "Node":
		return shapeNodeIface
	case "[]Node":
		return shapeNodeIfaceSlice
	case "[][]Node":
		return shapeNodeIfaceRows
	}
	name := typStr
	slice := false
	if strings.HasPrefix(name, "[]") {
		slice = true
		name = name[2:]
	}
	ptr := false
	if strings.HasPrefix(name, "*") {
		ptr = true
		name = name[1:]
	}
	// Reject any remaining wrapper (e.g. [][]*T, **T, maps) — the guard in
	// main reports node-bearing fields that end up here.
	if strings.ContainsAny(name, "[]*. ") {
		return shapeNone
	}
	if name == "Loc" {
		return shapeNone
	}
	switch {
	case nodeStructs[name] && ptr && slice:
		return shapeNodePtrSlice
	case nodeStructs[name] && ptr:
		return shapeNodePtr
	case nodeStructs[name] && slice:
		return shapeNodeValueSlice
	case nodeStructs[name]:
		return shapeNodeValue
	case bearingHelpers[name] && ptr && slice:
		return shapeHelperPtrSlice
	case bearingHelpers[name] && ptr:
		return shapeHelperPtr
	case bearingHelpers[name] && slice:
		return shapeHelperValueSlice
	case bearingHelpers[name]:
		return shapeHelperValue
	}
	return shapeNone
}

// typeString returns the string representation of a Go type expression.
// Supports identifiers, pointers, selector expressions, and slice types.
func typeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.SelectorExpr:
		return typeString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeString(t.Elt)
	default:
		return ""
	}
}

// baseType strips slice and pointer wrappers: "[]*WhenClause" → "WhenClause".
func baseType(t string) string {
	for {
		switch {
		case strings.HasPrefix(t, "[]"):
			t = t[2:]
		case strings.HasPrefix(t, "*"):
			t = t[1:]
		default:
			return t
		}
	}
}
