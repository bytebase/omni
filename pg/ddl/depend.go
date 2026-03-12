package ddl

// DepType classifies a dependency.
type DepType byte

const (
	DepNormal       DepType = 'n' // requires CASCADE to drop referent
	DepAuto         DepType = 'a' // auto-dropped with referent
	DepInternal     DepType = 'i' // inseparable from referent
	DepPartitionPri DepType = 'P' // partition object → parent object
	DepPartitionSec DepType = 'S' // partition object → partition table
)

// DepEntry records a dependency between two catalog objects.
//
// pg: src/include/catalog/pg_depend.h
type DepEntry struct {
	ObjType  byte    // 'r'=relation, 'i'=index, 'c'=constraint, 't'=type
	ObjOID   uint32
	ObjSubID int32   // pg: pg_depend.objsubid (attnum, or 0 for whole object)
	RefType  byte
	RefOID   uint32
	RefSubID int32   // pg: pg_depend.refobjsubid (attnum, or 0 for whole object)
	DepType  DepType
}

// recordDependency adds a dependency entry.
//
// pg: src/backend/catalog/pg_depend.c — recordDependencyOn
func (c *Catalog) recordDependency(objType byte, objOID uint32, objSubID int32,
	refType byte, refOID uint32, refSubID int32, depType DepType) {
	c.deps = append(c.deps, DepEntry{
		ObjType:  objType,
		ObjOID:   objOID,
		ObjSubID: objSubID,
		RefType:  refType,
		RefOID:   refOID,
		RefSubID: refSubID,
		DepType:  depType,
	})
}

// findNormalDependents returns all DepNormal entries that reference the given object.
func (c *Catalog) findNormalDependents(refType byte, refOID uint32) []DepEntry {
	var result []DepEntry
	for _, d := range c.deps {
		if d.RefType == refType && d.RefOID == refOID && d.DepType == DepNormal {
			result = append(result, d)
		}
	}
	return result
}

// removeDepsOf removes all dependency entries where the given object is the dependent (Obj).
func (c *Catalog) removeDepsOf(objType byte, objOID uint32) {
	n := 0
	for _, d := range c.deps {
		if d.ObjType == objType && d.ObjOID == objOID {
			continue
		}
		c.deps[n] = d
		n++
	}
	c.deps = c.deps[:n]
}

// recordDependencyOnExpr walks an analyzed expression tree and records
// dependencies on functions, operators, and types it references.
//
// pg: src/backend/catalog/dependency.c — recordDependencyOnExpr
func (c *Catalog) recordDependencyOnExpr(objType byte, objOID uint32, expr AnalyzedExpr, depType DepType) {
	refs := c.findExprReferences(expr)
	for _, ref := range refs {
		c.recordDependency(objType, objOID, 0, ref.refType, ref.refOID, ref.refSubID, depType)
	}
}

// recordDependencyOnSingleRelExpr records dependencies from an expression
// that references a single relation (varno=1, varlevelsup=0). Column
// references to that relation are recorded with selfBehavior; all other
// references (functions, operators, types) are recorded with behavior.
//
// pg: src/backend/catalog/dependency.c — recordDependencyOnSingleRelExpr
func (c *Catalog) recordDependencyOnSingleRelExpr(
	objType byte, objOID uint32,
	expr AnalyzedExpr, relOID uint32,
	behavior, selfBehavior DepType,
) {
	if expr == nil {
		return
	}

	// Walk expression tree for all references (functions, operators, types, columns).
	// pg: find_expr_references_walker with a bogus rangetable for Var resolution.
	var refs []depRef
	seen := make(map[depRef]bool)
	c.walkExprRefsWithRel(expr, relOID, &refs, seen)

	// Separate self-dependencies (column refs to relOID) from external deps.
	// pg: recordDependencyOnSingleRelExpr lines 1635-1653
	for _, ref := range refs {
		if ref.refType == 'r' && ref.refOID == relOID {
			c.recordDependency(objType, objOID, 0, ref.refType, ref.refOID, ref.refSubID, selfBehavior)
		} else {
			c.recordDependency(objType, objOID, 0, ref.refType, ref.refOID, ref.refSubID, behavior)
		}
	}
}

// walkExprRefsWithRel extends walkExprRefs to also resolve VarExpr nodes.
// VarExpr references to the given relOID are added with per-column SubIDs.
//
// pg: find_expr_references_walker — Var case (lines 1710-1751)
func (c *Catalog) walkExprRefsWithRel(expr AnalyzedExpr, relOID uint32, refs *[]depRef, seen map[depRef]bool) {
	if expr == nil {
		return
	}
	addRef := func(t byte, oid uint32, subID int32) {
		if oid == 0 || oid < FirstNormalObjectId {
			return
		}
		r := depRef{refType: t, refOID: oid, refSubID: subID}
		if !seen[r] {
			seen[r] = true
			*refs = append(*refs, r)
		}
	}

	switch v := expr.(type) {
	case *VarExpr:
		// pg: find_expr_references_walker — Var case
		// A whole-row Var (attno=0) adds no column dependency.
		// Otherwise add a per-column reference to the relation.
		if v.AttNum > 0 && v.LevelsUp == 0 {
			addRef('r', relOID, int32(v.AttNum))
		}
	case *FuncCallExpr:
		addRef('f', v.FuncOID, 0)
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
	case *AggExpr:
		addRef('f', v.AggFuncOID, 0)
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
	case *OpExpr:
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code, 0)
		}
		c.walkExprRefsWithRel(v.Left, relOID, refs, seen)
		c.walkExprRefsWithRel(v.Right, relOID, refs, seen)
	case *DistinctExprQ:
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code, 0)
		}
		c.walkExprRefsWithRel(v.Left, relOID, refs, seen)
		c.walkExprRefsWithRel(v.Right, relOID, refs, seen)
	case *ScalarArrayOpExpr:
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code, 0)
		}
		c.walkExprRefsWithRel(v.Left, relOID, refs, seen)
		c.walkExprRefsWithRel(v.Right, relOID, refs, seen)
	case *NullIfExprQ:
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code, 0)
		}
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
	case *CoerceViaIOExpr:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *RelabelExpr:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *BoolExprQ:
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
	case *CaseExprQ:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
		for _, w := range v.When {
			c.walkExprRefsWithRel(w.Condition, relOID, refs, seen)
			c.walkExprRefsWithRel(w.Result, relOID, refs, seen)
		}
		c.walkExprRefsWithRel(v.Default, relOID, refs, seen)
	case *CoalesceExprQ:
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
	case *MinMaxExprQ:
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
	case *NullTestExpr:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *BooleanTestExpr:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *SubLinkExpr:
		c.walkExprRefsWithRel(v.TestExpr, relOID, refs, seen)
		if v.SubQuery != nil {
			for _, te := range v.SubQuery.TargetList {
				c.walkExprRefsWithRel(te.Expr, relOID, refs, seen)
			}
		}
	case *ArrayExprQ:
		for _, elem := range v.Elements {
			c.walkExprRefsWithRel(elem, relOID, refs, seen)
		}
	case *RowExprQ:
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
	case *CollateExprQ:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *FieldSelectExprQ:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *WindowFuncExpr:
		addRef('f', v.FuncOID, 0)
		for _, arg := range v.Args {
			c.walkExprRefsWithRel(arg, relOID, refs, seen)
		}
		c.walkExprRefsWithRel(v.AggFilter, relOID, refs, seen)
	// ConstExpr, CoerceToDomainValueExpr, SQLValueFuncExpr — no references
	}
}

// depRef is an object reference found in an expression tree.
type depRef struct {
	refType  byte
	refOID   uint32
	refSubID int32
}

// findExprReferences walks an expression tree to collect referenced objects.
//
// pg: src/backend/catalog/dependency.c — find_expr_references_walker
func (c *Catalog) findExprReferences(expr AnalyzedExpr) []depRef {
	if expr == nil {
		return nil
	}
	var refs []depRef
	seen := make(map[depRef]bool)
	c.walkExprRefs(expr, &refs, seen)
	return refs
}

func (c *Catalog) walkExprRefs(expr AnalyzedExpr, refs *[]depRef, seen map[depRef]bool) {
	if expr == nil {
		return
	}
	addRef := func(t byte, oid uint32) {
		if oid == 0 || oid < FirstNormalObjectId {
			return
		}
		r := depRef{refType: t, refOID: oid}
		if !seen[r] {
			seen[r] = true
			*refs = append(*refs, r)
		}
	}

	switch v := expr.(type) {
	case *FuncCallExpr:
		addRef('f', v.FuncOID)
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
	case *AggExpr:
		addRef('f', v.AggFuncOID)
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
	case *OpExpr:
		// pg: OpExpr — record dependency on the operator's function
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code)
		}
		c.walkExprRefs(v.Left, refs, seen)
		c.walkExprRefs(v.Right, refs, seen)
	case *DistinctExprQ:
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code)
		}
		c.walkExprRefs(v.Left, refs, seen)
		c.walkExprRefs(v.Right, refs, seen)
	case *ScalarArrayOpExpr:
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code)
		}
		c.walkExprRefs(v.Left, refs, seen)
		c.walkExprRefs(v.Right, refs, seen)
	case *NullIfExprQ:
		if op := c.operByOID[v.OpOID]; op != nil {
			addRef('f', op.Code)
		}
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
	case *CoerceViaIOExpr:
		c.walkExprRefs(v.Arg, refs, seen)
	case *RelabelExpr:
		c.walkExprRefs(v.Arg, refs, seen)
	case *BoolExprQ:
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
	case *CaseExprQ:
		c.walkExprRefs(v.Arg, refs, seen)
		for _, w := range v.When {
			c.walkExprRefs(w.Condition, refs, seen)
			c.walkExprRefs(w.Result, refs, seen)
		}
		c.walkExprRefs(v.Default, refs, seen)
	case *CoalesceExprQ:
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
	case *MinMaxExprQ:
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
	case *NullTestExpr:
		c.walkExprRefs(v.Arg, refs, seen)
	case *BooleanTestExpr:
		c.walkExprRefs(v.Arg, refs, seen)
	case *SubLinkExpr:
		c.walkExprRefs(v.TestExpr, refs, seen)
		if v.SubQuery != nil {
			for _, te := range v.SubQuery.TargetList {
				c.walkExprRefs(te.Expr, refs, seen)
			}
		}
	case *ArrayExprQ:
		for _, elem := range v.Elements {
			c.walkExprRefs(elem, refs, seen)
		}
	case *RowExprQ:
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
	case *CollateExprQ:
		c.walkExprRefs(v.Arg, refs, seen)
	case *FieldSelectExprQ:
		c.walkExprRefs(v.Arg, refs, seen)
	case *WindowFuncExpr:
		addRef('f', v.FuncOID)
		for _, arg := range v.Args {
			c.walkExprRefs(arg, refs, seen)
		}
		c.walkExprRefs(v.AggFilter, refs, seen)
	// VarExpr, ConstExpr, CoerceToDomainValueExpr, SQLValueFuncExpr — no references
	}
}

// removeDepsOn removes all dependency entries where the given object is the referent (Ref).
func (c *Catalog) removeDepsOn(refType byte, refOID uint32) {
	n := 0
	for _, d := range c.deps {
		if d.RefType == refType && d.RefOID == refOID {
			continue
		}
		c.deps[n] = d
		n++
	}
	c.deps = c.deps[:n]
}
