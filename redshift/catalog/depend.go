package catalog

import "strings"

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
	ObjType  byte // 'r'=relation, 'i'=index, 'c'=constraint, 't'=type
	ObjOID   uint32
	ObjSubID int32 // pg: pg_depend.objsubid (attnum, or 0 for whole object)
	RefType  byte
	RefOID   uint32
	RefSubID int32 // pg: pg_depend.refobjsubid (attnum, or 0 for whole object)
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

// dependentObject is one object found by relationDependents.
type dependentObject struct {
	ObjType  byte
	ObjOID   uint32
	ObjSubID int32
	// Normal reports that every entry linking the object to the target is
	// DepNormal. RESTRICT refuses to drop such an object; one linked by any
	// AUTO or INTERNAL entry goes with the target.
	Normal bool
}

// relationDependents returns the objects that depend on the relation relOID,
// or on its column attnum when attnum is not 0, one entry per object. A
// column target matches only the entries on that column. A whole-relation
// target matches the entries on the relation and on any of its columns, and
// skips the relation's own entries (a generated column reading another
// column), since its columns go with it.
//
// pg: src/backend/catalog/dependency.c — findDependentObjects scans
// refobjsubid only for a column target; stack_address_present_add_flags
// treats a column of a relation being dropped as part of it; and
// reportDependentObjects reports only objects that no AUTO or INTERNAL
// dependency reached.
func (c *Catalog) relationDependents(relOID uint32, attnum int16) []dependentObject {
	type objKey struct {
		objType  byte
		objOID   uint32
		objSubID int32
	}
	var result []dependentObject
	pos := make(map[objKey]int)
	for _, d := range c.deps {
		if d.RefType != 'r' || d.RefOID != relOID {
			continue
		}
		if attnum != 0 && d.RefSubID != int32(attnum) {
			continue
		}
		if d.ObjType == 'r' && d.ObjOID == relOID && (attnum == 0 || d.ObjSubID == 0 || d.ObjSubID == int32(attnum)) {
			continue
		}
		k := objKey{d.ObjType, d.ObjOID, d.ObjSubID}
		i, ok := pos[k]
		if !ok {
			i = len(result)
			pos[k] = i
			result = append(result, dependentObject{ObjType: d.ObjType, ObjOID: d.ObjOID, ObjSubID: d.ObjSubID, Normal: true})
		}
		if d.DepType != DepNormal {
			result[i].Normal = false
		}
	}
	return result
}

// hasNormalDependent reports whether any of deps needs CASCADE to drop.
func hasNormalDependent(deps []dependentObject) bool {
	for _, d := range deps {
		if d.Normal {
			return true
		}
	}
	return false
}

// exprColumns returns the columns of the relation relOID that expr reads, in
// order of first appearance.
//
// pg: src/backend/catalog/heap.c — StoreRelCheck (pull_var_clause)
func (c *Catalog) exprColumns(expr AnalyzedExpr, relOID uint32) []int16 {
	var refs []depRef
	c.walkExprRefsWithRel(expr, relOID, &refs, make(map[depRef]bool))
	var cols []int16
	for _, r := range refs {
		if r.refType == 'r' && r.refOID == relOID {
			cols = append(cols, int16(r.refSubID))
		}
	}
	return cols
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
	c.recordDependencyOnSingleRelExprForObject(objType, objOID, 0, expr, relOID, behavior, selfBehavior)
}

func (c *Catalog) recordDependencyOnSingleRelExprForObject(
	objType byte, objOID uint32, objSubID int32,
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
			c.recordDependency(objType, objOID, objSubID, ref.refType, ref.refOID, ref.refSubID, selfBehavior)
		} else {
			c.recordDependency(objType, objOID, objSubID, ref.refType, ref.refOID, ref.refSubID, behavior)
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
		if seqOID := c.nextvalSequenceOID(v); seqOID != 0 {
			addRef('s', seqOID, 0)
		}
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
	case *SubscriptingRefExpr:
		c.walkExprRefsWithRel(v.ContainerExpr, relOID, refs, seen)
		for _, e := range v.SubscriptExprs {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
		for _, e := range v.LowerExprs {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
	case *NamedArgExprQ:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *ArrayCoerceExprQ:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
		c.walkExprRefsWithRel(v.ElemExpr, relOID, refs, seen)
	case *CoerceToDomainExpr:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *RowCompareExprQ:
		for _, e := range v.LArgs {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
		for _, e := range v.RArgs {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
	case *ConvertRowtypeExprQ:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
	case *FieldStoreExprQ:
		c.walkExprRefsWithRel(v.Arg, relOID, refs, seen)
		for _, e := range v.NewVals {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
	case *GroupingFuncExpr:
		for _, e := range v.Args {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
	case *XmlExprQ:
		for _, e := range v.NamedArgs {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
		for _, e := range v.Args {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
	case *JsonConstructorExprQ:
		for _, e := range v.Args {
			c.walkExprRefsWithRel(e, relOID, refs, seen)
		}
	case *JsonExprQ:
		c.walkExprRefsWithRel(v.Expr, relOID, refs, seen)
	case *JsonIsPredicateExpr:
		c.walkExprRefsWithRel(v.Expr, relOID, refs, seen)
	case *JsonValueExprQ:
		c.walkExprRefsWithRel(v.Expr, relOID, refs, seen)
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
		if seqOID := c.nextvalSequenceOID(v); seqOID != 0 {
			addRef('s', seqOID)
		}
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

func (c *Catalog) nextvalSequenceOID(fn *FuncCallExpr) uint32 {
	if fn == nil || !strings.EqualFold(fn.FuncName, "nextval") || len(fn.Args) != 1 {
		return 0
	}
	constArg := unwrapConstExpr(fn.Args[0])
	if constArg == nil || constArg.Value == "" || constArg.TypeOID != REGCLASSOID {
		return 0
	}
	schemaName, seqName := splitRegclassName(constArg.Value)
	seq, err := c.findSequence(schemaName, seqName)
	if err != nil {
		return 0
	}
	return seq.OID
}

func unwrapConstExpr(expr AnalyzedExpr) *ConstExpr {
	switch v := expr.(type) {
	case *ConstExpr:
		return v
	case *RelabelExpr:
		return unwrapConstExpr(v.Arg)
	case *CoerceViaIOExpr:
		return unwrapConstExpr(v.Arg)
	default:
		return nil
	}
}

func splitRegclassName(name string) (string, string) {
	if idx := strings.LastIndexByte(name, '.'); idx >= 0 {
		return name[:idx], name[idx+1:]
	}
	return "", name
}
