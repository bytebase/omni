package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// collectExprDeps walks an expression AST node and collects references to
// functions, relations (tables/sequences/views), and types. Only references
// that exist in the declared set are considered dependencies by the caller.
func collectExprDeps(expr nodes.Node) (funcRefs, relRefs, typeRefs []string) {
	if expr == nil {
		return
	}

	// First pass: collect CTE names to exclude from relRefs.
	cteNames := map[string]bool{}
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if cte, ok := n.(*nodes.CommonTableExpr); ok {
			cteNames[cte.Ctename] = true
		}
		return true
	})

	// Second pass: collect all references.
	nodes.Inspect(expr, func(n nodes.Node) bool {
		switch x := n.(type) {
		case *nodes.FuncCall:
			name := extractFuncCallName(x.Funcname)
			if name != "" {
				funcRefs = append(funcRefs, name)
			}
			// Special: nextval/currval/setval -> extract sequence name from string arg.
			if isSeqFunc(name) {
				if seqName := extractStringArg(x.Args); seqName != "" {
					relRefs = append(relRefs, normalizeQualifiedName(seqName))
				}
			}
		case *nodes.RangeVar:
			if x.Relname != "" && !cteNames[x.Relname] {
				relRefs = append(relRefs, qualifiedRangeVar(x))
			}
		case *nodes.TypeCast:
			if x.TypeName != nil {
				name := typeNameQualified(x.TypeName)
				if name != "" {
					typeRefs = append(typeRefs, name)
				}
			}
		}
		return true
	})
	return
}

// extractFuncCallName extracts the qualified function name from a FuncCall's Funcname list.
// Returns "schema.name" or "public.name" for unqualified names.
// For built-in functions (pg_catalog schema), returns the bare name with pg_catalog prefix.
func extractFuncCallName(funcname *nodes.List) string {
	if funcname == nil || len(funcname.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(funcname.Items))
	for _, item := range funcname.Items {
		if s, ok := item.(*nodes.String); ok {
			parts = append(parts, s.Str)
		}
	}
	if len(parts) == 1 {
		return "public." + parts[0]
	}
	return strings.Join(parts, ".")
}

// isSeqFunc returns true if the function name (qualified) is a sequence function.
func isSeqFunc(qualifiedName string) bool {
	// Strip schema prefix to get bare name.
	name := qualifiedName
	if idx := strings.LastIndex(qualifiedName, "."); idx >= 0 {
		name = qualifiedName[idx+1:]
	}
	switch name {
	case "nextval", "currval", "setval":
		return true
	}
	return false
}

// extractStringArg extracts the first string constant argument from a function's Args list.
func extractStringArg(args *nodes.List) string {
	if args == nil || len(args.Items) == 0 {
		return ""
	}
	// The first argument should be a string constant (A_Const with String val)
	// or a TypeCast wrapping a string constant.
	first := args.Items[0]
	return extractStringFromNode(first)
}

// extractStringFromNode extracts a string value from a node that may be
// an A_Const with String val, or a TypeCast wrapping one.
func extractStringFromNode(n nodes.Node) string {
	switch x := n.(type) {
	case *nodes.A_Const:
		if s, ok := x.Val.(*nodes.String); ok {
			return s.Str
		}
	case *nodes.TypeCast:
		return extractStringFromNode(x.Arg)
	}
	return ""
}

// normalizeQualifiedName takes a possibly schema-qualified name string like
// "public.my_seq" or "my_seq" and returns "schema.name".
func normalizeQualifiedName(name string) string {
	if strings.Contains(name, ".") {
		return name
	}
	return "public." + name
}

// collectExprDepsFromCreateStmt extracts expression-level dependencies from a CreateStmt.
// It walks column defaults, column-level constraints, and table-level constraints.
func collectExprDepsFromCreateStmt(s *nodes.CreateStmt) (funcRefs, relRefs, typeRefs []string) {
	if s.TableElts == nil {
		return
	}
	for _, elt := range s.TableElts.Items {
		switch n := elt.(type) {
		case *nodes.ColumnDef:
			// Walk RawDefault (column DEFAULT expression).
			if n.RawDefault != nil {
				fr, rr, tr := collectExprDeps(n.RawDefault)
				funcRefs = append(funcRefs, fr...)
				relRefs = append(relRefs, rr...)
				typeRefs = append(typeRefs, tr...)
			}
			// Walk column-level constraints.
			if n.Constraints != nil {
				for _, c := range n.Constraints.Items {
					if cons, ok := c.(*nodes.Constraint); ok && cons.RawExpr != nil {
						fr, rr, tr := collectExprDeps(cons.RawExpr)
						funcRefs = append(funcRefs, fr...)
						relRefs = append(relRefs, rr...)
						typeRefs = append(typeRefs, tr...)
					}
				}
			}
		case *nodes.Constraint:
			// Standalone table-level constraint in TableElts (e.g., CHECK).
			if n.RawExpr != nil {
				fr, rr, tr := collectExprDeps(n.RawExpr)
				funcRefs = append(funcRefs, fr...)
				relRefs = append(relRefs, rr...)
				typeRefs = append(typeRefs, tr...)
			}
		}
	}
	// Walk table-level constraints.
	if s.Constraints != nil {
		for _, item := range s.Constraints.Items {
			if cons, ok := item.(*nodes.Constraint); ok && cons.RawExpr != nil {
				fr, rr, tr := collectExprDeps(cons.RawExpr)
				funcRefs = append(funcRefs, fr...)
				relRefs = append(relRefs, rr...)
				typeRefs = append(typeRefs, tr...)
			}
		}
	}
	return
}

// collectExprDepsFromViewStmt extracts expression-level dependencies from a ViewStmt.
func collectExprDepsFromViewStmt(s *nodes.ViewStmt) (funcRefs, relRefs, typeRefs []string) {
	if s.Query == nil {
		return
	}
	return collectExprDeps(s.Query)
}

// collectExprDepsFromIndexStmt extracts expression-level dependencies from an IndexStmt.
func collectExprDepsFromIndexStmt(s *nodes.IndexStmt) (funcRefs, relRefs, typeRefs []string) {
	// Walk WhereClause (partial index).
	if s.WhereClause != nil {
		fr, rr, tr := collectExprDeps(s.WhereClause)
		funcRefs = append(funcRefs, fr...)
		relRefs = append(relRefs, rr...)
		typeRefs = append(typeRefs, tr...)
	}
	// Walk IndexParams for expression indexes.
	if s.IndexParams != nil {
		for _, item := range s.IndexParams.Items {
			if ie, ok := item.(*nodes.IndexElem); ok && ie.Expr != nil {
				fr, rr, tr := collectExprDeps(ie.Expr)
				funcRefs = append(funcRefs, fr...)
				relRefs = append(relRefs, rr...)
				typeRefs = append(typeRefs, tr...)
			}
		}
	}
	return
}

// collectExprDepsFromTriggerStmt extracts expression-level dependencies from a CreateTrigStmt.
func collectExprDepsFromTriggerStmt(s *nodes.CreateTrigStmt) (funcRefs, relRefs, typeRefs []string) {
	// Walk WhenClause.
	if s.WhenClause != nil {
		fr, rr, tr := collectExprDeps(s.WhenClause)
		funcRefs = append(funcRefs, fr...)
		relRefs = append(relRefs, rr...)
		typeRefs = append(typeRefs, tr...)
	}
	return
}

// collectExprDepsFromPolicyStmt extracts expression-level dependencies from a CreatePolicyStmt.
func collectExprDepsFromPolicyStmt(s *nodes.CreatePolicyStmt) (funcRefs, relRefs, typeRefs []string) {
	// Walk Qual (USING expression).
	if s.Qual != nil {
		fr, rr, tr := collectExprDeps(s.Qual)
		funcRefs = append(funcRefs, fr...)
		relRefs = append(relRefs, rr...)
		typeRefs = append(typeRefs, tr...)
	}
	// Walk WithCheck expression.
	if s.WithCheck != nil {
		fr, rr, tr := collectExprDeps(s.WithCheck)
		funcRefs = append(funcRefs, fr...)
		relRefs = append(relRefs, rr...)
		typeRefs = append(typeRefs, tr...)
	}
	return
}

// collectExprDepsFromDomainStmt extracts expression-level dependencies from a CreateDomainStmt.
func collectExprDepsFromDomainStmt(s *nodes.CreateDomainStmt) (funcRefs, relRefs, typeRefs []string) {
	// Walk constraint expressions.
	if s.Constraints != nil {
		for _, item := range s.Constraints.Items {
			if cons, ok := item.(*nodes.Constraint); ok && cons.RawExpr != nil {
				fr, rr, tr := collectExprDeps(cons.RawExpr)
				funcRefs = append(funcRefs, fr...)
				relRefs = append(relRefs, rr...)
				typeRefs = append(typeRefs, tr...)
			}
		}
	}
	// Base type dependency.
	if s.Typname != nil {
		name := typeNameQualified(s.Typname)
		if name != "" {
			typeRefs = append(typeRefs, name)
		}
	}
	return
}

// collectStructuralDepsFromFunctionStmt extracts type dependencies from function signatures.
func collectStructuralDepsFromFunctionStmt(s *nodes.CreateFunctionStmt) (funcRefs, relRefs, typeRefs []string) {
	// Return type.
	if s.ReturnType != nil {
		name := typeNameQualified(s.ReturnType)
		if name != "" {
			typeRefs = append(typeRefs, name)
		}
		// RETURNS SETOF table — the TypeName will have Setof=true.
		// The name is already captured above.
		// For SETOF <table>, it appears as a type ref but the declared set
		// might have it as a relation. We add it as both a type and rel ref.
		if s.ReturnType.Setof {
			relRefs = append(relRefs, name)
		}
	}
	// Parameter types and default expressions.
	if s.Parameters != nil {
		for _, item := range s.Parameters.Items {
			fp, ok := item.(*nodes.FunctionParameter)
			if !ok {
				continue
			}
			if fp.ArgType != nil {
				name := typeNameQualified(fp.ArgType)
				if name != "" {
					typeRefs = append(typeRefs, name)
				}
			}
			// Parameter default expression.
			if fp.Defexpr != nil {
				fr, rr, tr := collectExprDeps(fp.Defexpr)
				funcRefs = append(funcRefs, fr...)
				relRefs = append(relRefs, rr...)
				typeRefs = append(typeRefs, tr...)
			}
		}
	}
	return
}

// collectStructuralDepsFromCompositeType extracts type dependencies from composite type columns.
func collectStructuralDepsFromCompositeType(s *nodes.CompositeTypeStmt) (funcRefs, relRefs, typeRefs []string) {
	if s.Coldeflist == nil {
		return
	}
	for _, item := range s.Coldeflist.Items {
		col, ok := item.(*nodes.ColumnDef)
		if !ok {
			continue
		}
		if col.TypeName != nil {
			name := typeNameQualified(col.TypeName)
			if name != "" {
				typeRefs = append(typeRefs, name)
			}
		}
	}
	return
}

// collectStructuralDepsFromRangeType extracts type dependencies from range type params.
func collectStructuralDepsFromRangeType(s *nodes.CreateRangeStmt) (funcRefs, relRefs, typeRefs []string) {
	if s.Params == nil {
		return
	}
	for _, item := range s.Params.Items {
		de, ok := item.(*nodes.DefElem)
		if !ok {
			continue
		}
		if de.Defname == "subtype" {
			if tn, ok := de.Arg.(*nodes.TypeName); ok {
				name := typeNameQualified(tn)
				if name != "" {
					typeRefs = append(typeRefs, name)
				}
			}
		}
	}
	return
}
