package catalog

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Rule represents a rewrite rule defined via CREATE RULE on a table or view.
type Rule struct {
	OID         uint32
	OwnerRelOID uint32
	Name        string
	Event       nodes.CmdType // CMD_SELECT, CMD_INSERT, CMD_UPDATE, CMD_DELETE
	Instead     bool
	Mode        byte        // 'O'=origin (enabled), 'D'=disabled, 'R'=replica, 'A'=always
	WhereClause nodes.Node  // raw AST qualification (nil if none)
	Actions     *nodes.List // raw AST action statements (nil = DO NOTHING)
	bindings    *ruleBindings
}

type ruleColBinding struct {
	relOID     uint32
	attNum     int16
	prefix     string // "old", "new", explicit table alias, or ""
	useRelName bool   // true when qualified with unaliased relation name
}

type ruleTypeBinding struct {
	typeOID uint32
	isArray bool
}

type ruleColDep struct {
	relOID uint32
	attNum int16
}

type ruleJoinUsingCol struct {
	name string
	deps []ruleColDep
}

type ruleRegclassBinding struct {
	seqOID uint32
	relOID uint32
	idxOID uint32
}

type ruleStarCol struct {
	relOID     uint32
	attNum     int16
	prefix     string
	useRelName bool
	colAlias   string
	origName   string
	typeOID    uint32
	typeMod    int32
}

type ruleBindings struct {
	ownerRelOID    uint32
	isOnSelect     bool
	allowOld       bool
	allowNew       bool
	inRuleWhere    bool
	relByRV        map[*nodes.RangeVar]uint32
	seqByRV        map[*nodes.RangeVar]uint32
	tsmByRTS       map[*nodes.RangeTableSample]uint32
	colByRef       map[*nodes.ColumnRef]ruleColBinding
	colByRT        map[*nodes.ResTarget]ruleColBinding
	targetNameByRT map[*nodes.ResTarget]string
	starExpandByRT map[*nodes.ResTarget][]ruleStarCol
	colByIE        map[*nodes.IndexElem]ruleColBinding
	collByIE       map[*nodes.IndexElem]uint32
	opcByIE        map[*nodes.IndexElem]uint32
	conByInfer     map[*nodes.InferClause]uint32
	typeByTN       map[*nodes.TypeName]ruleTypeBinding
	collByCC       map[*nodes.CollateClause]uint32
	funcByFC       map[*nodes.FuncCall]uint32
	castTypeByFC   map[*nodes.FuncCall]uint32
	operByList     map[*nodes.List]uint32
	regclassByNode map[nodes.Node]ruleRegclassBinding
	exprTypeByNode map[nodes.Node]uint32
	joinUsing      map[*nodes.JoinExpr][]ruleJoinUsingCol
	naturalJoined  map[*nodes.JoinExpr]bool
	stmtCtx        map[*nodes.SelectStmt]*ruleEvalContext
	srfCalls       map[*nodes.FuncCall]bool

	// fn is the SQL function whose body these bindings hold, nil for a
	// rule. A name no column takes is one of its parameters, and a column
	// prints qualified, as print_function_sqlbody prints it.
	fn *sqlFnParams
	// fnQualify reports that the statement being bound has a range table
	// whatever its FROM holds (a DML target), so a parameter in it prints
	// qualified.
	fnQualify   bool
	paramByNode map[nodes.Node]sqlFnParamRef

	relDeps  []uint32
	colDeps  []ruleColDep
	conDeps  []uint32
	opcDeps  []uint32
	seqDeps  []uint32
	typeDeps []uint32
	collDeps []uint32
	funcDeps []uint32
	operDeps []uint32

	seenRel  map[uint32]bool
	seenCol  map[ruleColDep]bool
	seenCon  map[uint32]bool
	seenOpc  map[uint32]bool
	seenSeq  map[uint32]bool
	seenType map[uint32]bool
	seenColl map[uint32]bool
	seenFunc map[uint32]bool
	seenOper map[uint32]bool
}

func cloneRuleBindings(b *ruleBindings) *ruleBindings {
	if b == nil {
		return nil
	}
	cp := &ruleBindings{
		ownerRelOID:    b.ownerRelOID,
		isOnSelect:     b.isOnSelect,
		allowOld:       b.allowOld,
		allowNew:       b.allowNew,
		inRuleWhere:    b.inRuleWhere,
		relByRV:        make(map[*nodes.RangeVar]uint32, len(b.relByRV)),
		seqByRV:        make(map[*nodes.RangeVar]uint32, len(b.seqByRV)),
		tsmByRTS:       make(map[*nodes.RangeTableSample]uint32, len(b.tsmByRTS)),
		colByRef:       make(map[*nodes.ColumnRef]ruleColBinding, len(b.colByRef)),
		colByRT:        make(map[*nodes.ResTarget]ruleColBinding, len(b.colByRT)),
		targetNameByRT: make(map[*nodes.ResTarget]string, len(b.targetNameByRT)),
		starExpandByRT: make(map[*nodes.ResTarget][]ruleStarCol, len(b.starExpandByRT)),
		colByIE:        make(map[*nodes.IndexElem]ruleColBinding, len(b.colByIE)),
		collByIE:       make(map[*nodes.IndexElem]uint32, len(b.collByIE)),
		opcByIE:        make(map[*nodes.IndexElem]uint32, len(b.opcByIE)),
		conByInfer:     make(map[*nodes.InferClause]uint32, len(b.conByInfer)),
		typeByTN:       make(map[*nodes.TypeName]ruleTypeBinding, len(b.typeByTN)),
		collByCC:       make(map[*nodes.CollateClause]uint32, len(b.collByCC)),
		funcByFC:       make(map[*nodes.FuncCall]uint32, len(b.funcByFC)),
		castTypeByFC:   make(map[*nodes.FuncCall]uint32, len(b.castTypeByFC)),
		operByList:     make(map[*nodes.List]uint32, len(b.operByList)),
		regclassByNode: make(map[nodes.Node]ruleRegclassBinding, len(b.regclassByNode)),
		exprTypeByNode: make(map[nodes.Node]uint32, len(b.exprTypeByNode)),
		joinUsing:      make(map[*nodes.JoinExpr][]ruleJoinUsingCol, len(b.joinUsing)),
		naturalJoined:  make(map[*nodes.JoinExpr]bool, len(b.naturalJoined)),
		relDeps:        append([]uint32(nil), b.relDeps...),
		colDeps:        append([]ruleColDep(nil), b.colDeps...),
		conDeps:        append([]uint32(nil), b.conDeps...),
		opcDeps:        append([]uint32(nil), b.opcDeps...),
		seqDeps:        append([]uint32(nil), b.seqDeps...),
		typeDeps:       append([]uint32(nil), b.typeDeps...),
		collDeps:       append([]uint32(nil), b.collDeps...),
		funcDeps:       append([]uint32(nil), b.funcDeps...),
		operDeps:       append([]uint32(nil), b.operDeps...),
	}
	for k, v := range b.relByRV {
		cp.relByRV[k] = v
	}
	for k, v := range b.seqByRV {
		cp.seqByRV[k] = v
	}
	for k, v := range b.tsmByRTS {
		cp.tsmByRTS[k] = v
	}
	for k, v := range b.colByRef {
		cp.colByRef[k] = v
	}
	for k, v := range b.colByRT {
		cp.colByRT[k] = v
	}
	for k, v := range b.targetNameByRT {
		cp.targetNameByRT[k] = v
	}
	for k, v := range b.starExpandByRT {
		cp.starExpandByRT[k] = append([]ruleStarCol(nil), v...)
	}
	for k, v := range b.colByIE {
		cp.colByIE[k] = v
	}
	for k, v := range b.collByIE {
		cp.collByIE[k] = v
	}
	for k, v := range b.opcByIE {
		cp.opcByIE[k] = v
	}
	for k, v := range b.conByInfer {
		cp.conByInfer[k] = v
	}
	for k, v := range b.typeByTN {
		cp.typeByTN[k] = v
	}
	for k, v := range b.collByCC {
		cp.collByCC[k] = v
	}
	for k, v := range b.funcByFC {
		cp.funcByFC[k] = v
	}
	for k, v := range b.castTypeByFC {
		cp.castTypeByFC[k] = v
	}
	for k, v := range b.operByList {
		cp.operByList[k] = v
	}
	for k, v := range b.regclassByNode {
		cp.regclassByNode[k] = v
	}
	for k, v := range b.exprTypeByNode {
		cp.exprTypeByNode[k] = v
	}
	for k, v := range b.joinUsing {
		cpCols := make([]ruleJoinUsingCol, len(v))
		for i, uc := range v {
			cpCols[i] = ruleJoinUsingCol{
				name: uc.name,
				deps: append([]ruleColDep(nil), uc.deps...),
			}
		}
		cp.joinUsing[k] = cpCols
	}
	for k, v := range b.naturalJoined {
		cp.naturalJoined[k] = v
	}
	cp.fn = b.fn
	cp.paramByNode = make(map[nodes.Node]sqlFnParamRef, len(b.paramByNode))
	for k, v := range b.paramByNode {
		cp.paramByNode[k] = v
	}
	return cp
}

func (b *ruleBindings) addRelDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenRel[oid] {
		return
	}
	b.seenRel[oid] = true
	b.relDeps = append(b.relDeps, oid)
}

func (b *ruleBindings) addColDep(relOID uint32, attNum int16) {
	if relOID < FirstNormalObjectId || attNum <= 0 {
		return
	}
	key := ruleColDep{relOID: relOID, attNum: attNum}
	if b.seenCol[key] {
		return
	}
	b.seenCol[key] = true
	b.colDeps = append(b.colDeps, key)
}

func (b *ruleBindings) addConDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenCon[oid] {
		return
	}
	b.seenCon[oid] = true
	b.conDeps = append(b.conDeps, oid)
}

func (b *ruleBindings) addOpcDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenOpc[oid] {
		return
	}
	b.seenOpc[oid] = true
	b.opcDeps = append(b.opcDeps, oid)
}

func (b *ruleBindings) addSeqDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenSeq[oid] {
		return
	}
	b.seenSeq[oid] = true
	b.seqDeps = append(b.seqDeps, oid)
}

func (b *ruleBindings) addTypeDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenType[oid] {
		return
	}
	b.seenType[oid] = true
	b.typeDeps = append(b.typeDeps, oid)
}

func (b *ruleBindings) addCollDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenColl[oid] {
		return
	}
	b.seenColl[oid] = true
	b.collDeps = append(b.collDeps, oid)
}

func (b *ruleBindings) addFuncDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenFunc[oid] {
		return
	}
	b.seenFunc[oid] = true
	b.funcDeps = append(b.funcDeps, oid)
}

func (b *ruleBindings) addOperDep(oid uint32) {
	if oid < FirstNormalObjectId || b.seenOper[oid] {
		return
	}
	b.seenOper[oid] = true
	b.operDeps = append(b.operDeps, oid)
}

// CreateRuleStmt validates and registers a rewrite rule from a RuleStmt AST.
//
// pg: src/backend/rewrite/rewriteDefine.c — DefineQueryRewrite
func (c *Catalog) CreateRuleStmt(s *nodes.RuleStmt) error {
	if s.Relation == nil {
		return nil
	}
	_, rel, err := c.findRelation(s.Relation.Schemaname, s.Relation.Relname)
	if err != nil {
		return err
	}
	switch rel.RelKind {
	case 'r', 'v', 'p':
	case 'm':
		return &Error{Code: CodeFeatureNotSupported, Message: "rules on materialized views are not supported"}
	default:
		return errWrongObjectType(rel.Name, "a table or view")
	}
	if s.Rulename == "" {
		return nil
	}
	if s.Event == nodes.CMD_SELECT {
		return c.createOnSelectRule(rel, s)
	}
	if rel.RelKind == 'v' && s.Rulename == "_RETURN" {
		if !s.Replace {
			return errDuplicateObject("rule", s.Rulename)
		}
		return &Error{
			Code:    CodeInvalidObjectDefinition,
			Message: fmt.Sprintf("cannot change event of rule %q", s.Rulename),
		}
	}
	if rel.Rules != nil && rel.Rules[s.Rulename] != nil {
		if !s.Replace {
			return errDuplicateObject("rule", s.Rulename)
		}
		if rel.Rules[s.Rulename].Event != s.Event {
			return &Error{
				Code:    CodeInvalidObjectDefinition,
				Message: fmt.Sprintf("cannot change event of rule %q", s.Rulename),
			}
		}
	}
	if s.WhereClause != nil && s.Actions != nil {
		for _, item := range s.Actions.Items {
			if _, isNotify := item.(*nodes.NotifyStmt); isNotify {
				return &Error{
					Code:    CodeInvalidObjectDefinition,
					Message: "rules with WHERE conditions can only have SELECT, INSERT, UPDATE, or DELETE actions",
				}
			}
		}
	}
	if err := c.validateRuleReturningClauses(rel, s); err != nil {
		return err
	}

	bindings, err := c.buildRuleBindings(rel, s.Event, s.WhereClause, s.Actions)
	if err != nil {
		return err
	}

	if rel.Rules == nil {
		rel.Rules = make(map[string]*Rule)
	}
	mode := byte('O')
	var oid uint32
	if existing := rel.Rules[s.Rulename]; existing != nil {
		mode = existing.Mode
		oid = existing.OID
		c.removeDepsOf('R', oid)
	} else {
		oid = c.oidGen.Next()
	}

	rule := &Rule{
		OID:         oid,
		OwnerRelOID: rel.OID,
		Name:        s.Rulename,
		Event:       s.Event,
		Instead:     s.Instead,
		Mode:        mode,
		WhereClause: s.WhereClause,
		Actions:     s.Actions,
		bindings:    bindings,
	}
	rel.Rules[s.Rulename] = rule
	c.rules[oid] = rule
	rel.HasRules = true
	c.recordRuleDependencies(rule)
	return nil
}

func (c *Catalog) recordRuleDependencies(rule *Rule) {
	if rule == nil || rule.OID == 0 {
		return
	}
	// Automatic dependency on the owning relation so DROP TABLE/VIEW drops its own rules.
	c.recordDependency('R', rule.OID, 0, 'r', rule.OwnerRelOID, 0, DepAuto)

	b := rule.bindings
	if b == nil {
		return
	}
	for _, relOID := range b.relDeps {
		if relOID != rule.OwnerRelOID {
			c.recordDependency('R', rule.OID, 0, 'r', relOID, 0, DepNormal)
		}
	}
	for _, cd := range b.colDeps {
		c.recordDependency('R', rule.OID, 0, 'r', cd.relOID, int32(cd.attNum), DepNormal)
	}
	for _, conOID := range b.conDeps {
		c.recordDependency('R', rule.OID, 0, 'c', conOID, 0, DepNormal)
	}
	for _, opcOID := range b.opcDeps {
		c.recordDependency('R', rule.OID, 0, 'q', opcOID, 0, DepNormal)
	}
	for _, seqOID := range b.seqDeps {
		c.recordDependency('R', rule.OID, 0, 's', seqOID, 0, DepNormal)
	}
	for _, typeOID := range b.typeDeps {
		c.recordDependency('R', rule.OID, 0, 't', typeOID, 0, DepNormal)
	}
	for _, collOID := range b.collDeps {
		c.recordDependency('R', rule.OID, 0, 'L', collOID, 0, DepNormal)
	}
	for _, funcOID := range b.funcDeps {
		c.recordDependency('R', rule.OID, 0, 'f', funcOID, 0, DepNormal)
	}
	for _, operOID := range b.operDeps {
		c.recordDependency('R', rule.OID, 0, 'O', operOID, 0, DepNormal)
	}
}

func (c *Catalog) removeRule(rel *Relation, rule *Rule) {
	if rule == nil {
		return
	}
	if rel != nil && rel.Rules != nil {
		if rel.Rules[rule.Name] == rule {
			delete(rel.Rules, rule.Name)
		}
		rel.HasRules = len(rel.Rules) > 0
	}
	delete(c.rules, rule.OID)
	delete(c.extensionMembers, rule.OID)
	c.removeComments('R', rule.OID)
	c.removeDepsOf('R', rule.OID)
	c.removeDepsOn('R', rule.OID)
}

func (c *Catalog) removeRulesForRelation(rel *Relation) {
	if rel == nil || len(rel.Rules) == 0 {
		return
	}
	rules := make([]*Rule, 0, len(rel.Rules))
	for _, r := range rel.Rules {
		if r != nil {
			rules = append(rules, r)
		}
	}
	for _, r := range rules {
		c.removeRule(rel, r)
	}
	rel.Rules = nil
	rel.HasRules = false
}

// createOnSelectRule validates PostgreSQL 17's ON SELECT (_RETURN) rule
// constraints and replaces the view's query definition.
//
// pg: src/backend/rewrite/rewriteDefine.c — DefineQueryRewrite (CMD_SELECT)
func (c *Catalog) createOnSelectRule(rel *Relation, s *nodes.RuleStmt) error {
	if rel.RelKind != 'v' {
		return &Error{
			Code:    CodeWrongObjectType,
			Message: fmt.Sprintf("relation %q cannot have ON SELECT rules", rel.Name),
		}
	}
	if !s.Replace {
		return &Error{
			Code:    CodeObjectNotInPrerequisiteState,
			Message: fmt.Sprintf("%q is already a view", rel.Name),
		}
	}
	if s.Rulename != "_RETURN" {
		return &Error{
			Code:    CodeInvalidObjectDefinition,
			Message: fmt.Sprintf("view rule for %q must be named \"_RETURN\"", rel.Name),
		}
	}
	if s.WhereClause != nil {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "event qualifications are not implemented for rules on SELECT",
		}
	}
	if s.Actions != nil && len(s.Actions.Items) > 1 {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "multiple actions for rules on SELECT are not implemented",
		}
	}
	if !s.Instead || s.Actions == nil || len(s.Actions.Items) == 0 {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "rules on SELECT must have action INSTEAD SELECT",
		}
	}
	selStmt, ok := s.Actions.Items[0].(*nodes.SelectStmt)
	if !ok || selStmt == nil {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "rules on SELECT must have action INSTEAD SELECT",
		}
	}
	if _, err := c.buildRuleBindings(rel, s.Event, nil, s.Actions); err != nil {
		return err
	}
	q, err := c.analyzeSelectStmt(selStmt)
	if err != nil {
		if c.snapshot {
			return nil
		}
		return err
	}
	var nonJunk []*TargetEntry
	for _, te := range q.TargetList {
		if !te.ResJunk {
			nonJunk = append(nonJunk, te)
		}
	}
	if !c.snapshot {
		if len(nonJunk) < len(rel.Columns) {
			return &Error{
				Code:    CodeInvalidObjectDefinition,
				Message: "SELECT rule's target list has too few entries",
			}
		}
		if len(nonJunk) > len(rel.Columns) {
			return &Error{
				Code:    CodeInvalidObjectDefinition,
				Message: "SELECT rule's target list has too many entries",
			}
		}
		for i, col := range rel.Columns {
			te := nonJunk[i]
			if te.ResName != col.Name {
				return &Error{
					Code:    CodeInvalidObjectDefinition,
					Message: fmt.Sprintf("SELECT rule's target entry %d has different column name from column %q", i+1, col.Name),
				}
			}
			if te.Expr.exprType() != col.TypeOID || (col.TypeMod >= 0 && te.Expr.exprTypMod() >= 0 && te.Expr.exprTypMod() != col.TypeMod) {
				return &Error{
					Code:    CodeInvalidObjectDefinition,
					Message: fmt.Sprintf("SELECT rule's target entry %d has different type from column %q", i+1, col.Name),
				}
			}
		}
	}
	rel.AnalyzedQuery = q
	c.removeDepsOf('r', rel.OID)
	c.recordViewDepsQ(rel.OID, q)
	return nil
}

func ruleActionsHaveReturning(actions *nodes.List) bool {
	if actions == nil {
		return false
	}
	for _, item := range actions.Items {
		var retList *nodes.List
		switch act := item.(type) {
		case *nodes.InsertStmt:
			retList = act.ReturningList
		case *nodes.UpdateStmt:
			retList = act.ReturningList
		case *nodes.DeleteStmt:
			retList = act.ReturningList
		case *nodes.MergeStmt:
			retList = act.ReturningList
		}
		if retList != nil && len(retList.Items) > 0 {
			return true
		}
	}
	return false
}

// validateRuleReturningClauses enforces PostgreSQL's checkRuleResultList rules
// on RETURNING clauses inside rule actions.
//
// pg: src/backend/rewrite/rewriteDefine.c — DefineQueryRewrite, checkRuleResultList
func (c *Catalog) validateRuleReturningClauses(rel *Relation, s *nodes.RuleStmt) error {
	if s.Actions == nil {
		return nil
	}
	returningCount := 0
	for _, item := range s.Actions.Items {
		var retList *nodes.List
		switch act := item.(type) {
		case *nodes.InsertStmt:
			retList = act.ReturningList
		case *nodes.UpdateStmt:
			retList = act.ReturningList
		case *nodes.DeleteStmt:
			retList = act.ReturningList
		case *nodes.MergeStmt:
			retList = act.ReturningList
		}
		if retList == nil || len(retList.Items) == 0 {
			continue
		}
		returningCount++
		if !s.Instead || s.WhereClause != nil {
			return &Error{
				Code:    CodeFeatureNotSupported,
				Message: "RETURNING lists are not supported in conditional or non-INSTEAD rules",
			}
		}
		if returningCount > 1 {
			return &Error{
				Code:    CodeFeatureNotSupported,
				Message: "cannot have RETURNING lists in multiple rules",
			}
		}
	}
	if returningCount > 0 && rel != nil && rel.Rules != nil {
		for name, existing := range rel.Rules {
			if existing == nil || existing.Event != s.Event {
				continue
			}
			if s.Replace && name == s.Rulename {
				continue
			}
			if ruleActionsHaveReturning(existing.Actions) {
				return &Error{
					Code:    CodeFeatureNotSupported,
					Message: "cannot have RETURNING lists in multiple rules",
				}
			}
		}
	}
	return nil
}

func newRuleBindings(ownerRelOID uint32, event nodes.CmdType) *ruleBindings {
	return &ruleBindings{
		ownerRelOID:    ownerRelOID,
		isOnSelect:     event == nodes.CMD_SELECT,
		allowOld:       event == nodes.CMD_UPDATE || event == nodes.CMD_DELETE,
		allowNew:       event == nodes.CMD_INSERT || event == nodes.CMD_UPDATE,
		relByRV:        make(map[*nodes.RangeVar]uint32),
		seqByRV:        make(map[*nodes.RangeVar]uint32),
		tsmByRTS:       make(map[*nodes.RangeTableSample]uint32),
		colByRef:       make(map[*nodes.ColumnRef]ruleColBinding),
		colByRT:        make(map[*nodes.ResTarget]ruleColBinding),
		targetNameByRT: make(map[*nodes.ResTarget]string),
		starExpandByRT: make(map[*nodes.ResTarget][]ruleStarCol),
		colByIE:        make(map[*nodes.IndexElem]ruleColBinding),
		collByIE:       make(map[*nodes.IndexElem]uint32),
		opcByIE:        make(map[*nodes.IndexElem]uint32),
		conByInfer:     make(map[*nodes.InferClause]uint32),
		typeByTN:       make(map[*nodes.TypeName]ruleTypeBinding),
		collByCC:       make(map[*nodes.CollateClause]uint32),
		funcByFC:       make(map[*nodes.FuncCall]uint32),
		castTypeByFC:   make(map[*nodes.FuncCall]uint32),
		operByList:     make(map[*nodes.List]uint32),
		regclassByNode: make(map[nodes.Node]ruleRegclassBinding),
		exprTypeByNode: make(map[nodes.Node]uint32),
		joinUsing:      make(map[*nodes.JoinExpr][]ruleJoinUsingCol),
		naturalJoined:  make(map[*nodes.JoinExpr]bool),
		stmtCtx:        make(map[*nodes.SelectStmt]*ruleEvalContext),
		srfCalls:       make(map[*nodes.FuncCall]bool),
		paramByNode:    make(map[nodes.Node]sqlFnParamRef),
		seenRel:        make(map[uint32]bool),
		seenCol:        make(map[ruleColDep]bool),
		seenCon:        make(map[uint32]bool),
		seenOpc:        make(map[uint32]bool),
		seenSeq:        make(map[uint32]bool),
		seenType:       make(map[uint32]bool),
		seenColl:       make(map[uint32]bool),
		seenFunc:       make(map[uint32]bool),
		seenOper:       make(map[uint32]bool),
	}
}

func (c *Catalog) validatePreparedDMLStmt(n nodes.Node, paramTypes []uint32, fixedParams ...bool) error {
	b := newRuleBindings(0, 0)
	if err := c.bindRuleStmt(b, n, &ruleEvalContext{}); err != nil {
		return err
	}
	isFixed := len(fixedParams) > 0 && fixedParams[0]
	if err := c.inferPreparedDMLParamTypes(b, n, paramTypes, isFixed); err != nil && !c.snapshot {
		return err
	}
	return nil
}

// renameRule handles ALTER RULE name ON table RENAME TO new_name.
//
// pg: src/backend/rewrite/rewriteDefine.c — RenameRewriteRule
func (c *Catalog) renameRule(stmt *nodes.RenameStmt) error {
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}
	_, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		if stmt.MissingOk && missingObject(err) {
			return nil
		}
		return err
	}
	if rel.RelKind != 'r' && rel.RelKind != 'v' && rel.RelKind != 'p' {
		return errWrongObjectType(rel.Name, "a table or view")
	}
	var rule *Rule
	if rel.Rules != nil {
		rule = rel.Rules[stmt.Subname]
	}
	if rule == nil {
		if c.snapshot {
			return nil
		}
		return &Error{
			Code:    CodeUndefinedObject,
			Message: fmt.Sprintf("rule %q for relation %q does not exist", stmt.Subname, rel.Name),
		}
	}
	if stmt.Newname != stmt.Subname && rel.Rules[stmt.Newname] != nil {
		return &Error{
			Code:    CodeDuplicateObject,
			Message: fmt.Sprintf("rule %q for relation %q already exists", stmt.Newname, rel.Name),
		}
	}
	delete(rel.Rules, stmt.Subname)
	rule.Name = stmt.Newname
	rel.Rules[stmt.Newname] = rule
	return nil
}

type ruleScopeEntry struct {
	relOID               uint32
	hasAlias             bool
	alias                string
	colAliases           []string
	isVirtual            bool // CTE, subquery, sequence, table function, XMLTABLE, JSON_TABLE, or aliased JOIN
	virtualCols          map[string]int
	virtualColTypes      map[string]uint32
	virtualColOrder      []string
	anyCol               bool
	isJoinAlias          bool
	hasJoinColAliases    bool
	hiddenByJoinColAlias bool
	joinRels             []uint32
}

type ruleEvalContext struct {
	scope         map[string]ruleScopeEntry
	scopeOrder    []string
	outer         *ruleEvalContext
	cteNames      map[string]bool
	cteCols       map[string][]string
	cteColTypes   map[string][]uint32
	cteAny        map[string]bool
	usingCols     map[string]int
	targetAliases map[string]bool
	windowNames   map[string]bool
	windowDefs    map[string]*nodes.WindowDef
	hasOpenSource bool
	disallowAgg   string
	disallowWin   string
	// disallowSRF names the construct that refuses set-returning calls
	// (parse_func.c check_srf_call_placement), or "" where they are allowed.
	disallowSRF string
	// insertSource marks the context an INSERT binds its source query in: a
	// single VALUES row there is a target list, which allows set-returning
	// calls (analyze.c transformInsertStmt).
	insertSource bool
}

func (ctx *ruleEvalContext) setScope(name string, se ruleScopeEntry, snapshot bool) error {
	if ctx == nil || name == "" {
		return nil
	}
	if ctx.scope == nil {
		ctx.scope = make(map[string]ruleScopeEntry)
	}
	if _, exists := ctx.scope[name]; exists {
		if !snapshot && !strings.HasPrefix(name, "\x00") {
			return &Error{
				Code:    CodeDuplicateAlias,
				Message: fmt.Sprintf("table name %q specified more than once", name),
			}
		}
	} else {
		ctx.scopeOrder = append(ctx.scopeOrder, name)
	}
	ctx.scope[name] = se
	return nil
}

func (c *Catalog) buildRuleBindings(ownerRel *Relation, event nodes.CmdType, whereClause nodes.Node, actions *nodes.List) (*ruleBindings, error) {
	b := newRuleBindings(ownerRel.OID, event)
	if whereClause != nil {
		b.inRuleWhere = true
		err := c.bindRuleBoolExpr(b, whereClause, &ruleEvalContext{}, "WHERE")
		b.inRuleWhere = false
		if err != nil {
			return nil, err
		}
	}
	if actions != nil {
		for _, item := range actions.Items {
			if err := c.bindRuleStmt(b, item, &ruleEvalContext{}); err != nil {
				return nil, err
			}
		}
	}
	return b, nil
}

func (c *Catalog) bindRuleBoolExpr(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext, construct string) error {
	if n == nil {
		return nil
	}
	if ctx != nil {
		prevAgg, prevWin, prevSRF := ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF
		if ctx.disallowSRF == "" {
			switch construct {
			case "WHERE", "HAVING", "FILTER":
				ctx.disallowSRF = construct
			case "JOIN/ON":
				ctx.disallowSRF = "JOIN conditions"
			}
		}
		switch construct {
		case "WHERE", "JOIN/ON", "WHEN", "FILTER":
			if ctx.disallowAgg == "" {
				ctx.disallowAgg = construct
			}
			if ctx.disallowWin == "" {
				ctx.disallowWin = construct
			}
		case "HAVING":
			if ctx.disallowWin == "" {
				ctx.disallowWin = construct
			}
		}
		defer func() {
			ctx.disallowAgg = prevAgg
			ctx.disallowWin = prevWin
			ctx.disallowSRF = prevSRF
		}()
	}
	if err := c.bindRuleExpr(b, n, ctx); err != nil {
		return err
	}
	if c.snapshot {
		return nil
	}
	t := c.inferRuleExprType(b, n, ctx)
	if t == 0 || t == BOOLOID || t == UNKNOWNOID {
		return nil
	}
	if !c.CanCoerce(t, BOOLOID, 'a') {
		return &Error{
			Code:    CodeDatatypeMismatch,
			Message: fmt.Sprintf("argument of %s must be type boolean, not type %s", construct, c.typeName(t)),
		}
	}
	return nil
}

func extractAliasColnames(alias *nodes.Alias) []string {
	if alias == nil || alias.Colnames == nil || len(alias.Colnames.Items) == 0 {
		return nil
	}
	out := make([]string, 0, len(alias.Colnames.Items))
	for _, it := range alias.Colnames.Items {
		out = append(out, stringVal(it))
	}
	return out
}

func extractNodeListStrings(list *nodes.List) []string {
	if list == nil || len(list.Items) == 0 {
		return nil
	}
	out := make([]string, 0, len(list.Items))
	for _, it := range list.Items {
		out = append(out, stringVal(it))
	}
	return out
}

func applyVirtualAliasCols(baseCols []string, baseTypes []uint32, baseAny bool, aliasCols []string) (map[string]int, map[string]uint32, []string, bool) {
	if len(baseCols) == 0 && len(aliasCols) == 0 {
		return nil, nil, nil, baseAny
	}
	n := len(baseCols)
	if len(aliasCols) > n {
		n = len(aliasCols)
	}
	outMap := make(map[string]int, n)
	outTypes := make(map[string]uint32, n)
	outOrder := make([]string, 0, n)
	for i := 0; i < n; i++ {
		var name string
		if i < len(aliasCols) && aliasCols[i] != "" {
			name = aliasCols[i]
		} else if i < len(baseCols) {
			name = baseCols[i]
		}
		if name != "" {
			outMap[name]++
			if i < len(baseTypes) && baseTypes[i] != 0 && outTypes[name] == 0 {
				outTypes[name] = baseTypes[i]
			}
		}
		outOrder = append(outOrder, name)
	}
	return outMap, outTypes, outOrder, baseAny
}

func (c *Catalog) bindRangeVar(b *ruleBindings, rv *nodes.RangeVar, cteNames map[string]bool, forDMLTarget bool) (uint32, string, error) {
	if rv == nil {
		return 0, "", nil
	}
	if rv.Catalogname != "" {
		if err := c.checkDatabaseQualifier(rv.Catalogname); err != nil {
			return 0, "", err
		}
	}
	if !forDMLTarget && rv.Schemaname == "" && cteNames[rv.Relname] {
		return 0, "cte", nil
	}
	_, rel, err := c.findRelation(rv.Schemaname, rv.Relname)
	if err != nil {
		if !forDMLTarget {
			if seq, seqErr := c.findSequence(rv.Schemaname, rv.Relname); seqErr == nil && seq != nil {
				if b != nil && b.seqByRV != nil {
					b.seqByRV[rv] = seq.OID
				}
				b.addSeqDep(seq.OID)
				return 0, "seq", nil
			}
		}
		if c.snapshot || b != nil && b.fn != nil && !forDMLTarget && systemRelationName(rv) {
			// The catalog holds no system catalog or information_schema
			// view, which a function body may read.
			return 0, "open", nil
		}
		return 0, "", err
	}
	if rel.RelKind == 'c' {
		return 0, "", errWrongObjectType(rel.Name, "a table")
	}
	if forDMLTarget {
		switch rel.RelKind {
		case 'r', 'p', 'v', 'f':
		case 'm':
			return 0, "", &Error{
				Code:    CodeWrongObjectType,
				Message: fmt.Sprintf("cannot change materialized view %q", rel.Name),
			}
		default:
			return 0, "", errWrongObjectType(rel.Name, "a table")
		}
	}
	b.relByRV[rv] = rel.OID
	b.addRelDep(rel.OID)
	return rel.OID, "", nil
}

func addScopeRangeVar(c *Catalog, b *ruleBindings, ctx *ruleEvalContext, rv *nodes.RangeVar, forDMLTarget bool) (uint32, error) {
	relOID, virtualKind, err := c.bindRangeVar(b, rv, ctx.cteNames, forDMLTarget)
	if err != nil {
		return 0, err
	}
	if rv == nil {
		return 0, nil
	}
	name := rv.Relname
	hasAlias := false
	alias := ""
	if rv.Alias != nil && rv.Alias.Aliasname != "" {
		name = rv.Alias.Aliasname
		hasAlias = true
		alias = rv.Alias.Aliasname
	}
	if virtualKind != "" {
		var baseCols []string
		var baseTypes []uint32
		baseAny := false
		switch virtualKind {
		case "cte":
			baseCols = ctx.cteCols[rv.Relname]
			if ctx.cteColTypes != nil {
				baseTypes = ctx.cteColTypes[rv.Relname]
			}
			baseAny = ctx.cteAny[rv.Relname]
		case "seq":
			baseCols = []string{"last_value", "log_cnt", "is_called"}
			baseTypes = []uint32{INT8OID, INT8OID, BOOLOID}
		default:
			baseAny = true
		}
		aliasCols := extractAliasColnames(rv.Alias)
		if !baseAny && len(aliasCols) > len(baseCols) && !c.snapshot {
			return 0, &Error{
				Code:    CodeInvalidColumnReference,
				Message: fmt.Sprintf("table %q has %d columns available but %d columns specified", name, len(baseCols), len(aliasCols)),
			}
		}
		vCols, vTypes, vOrder, vAny := applyVirtualAliasCols(baseCols, baseTypes, baseAny, aliasCols)
		if err := ctx.setScope(name, ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           alias,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColTypes: vTypes,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}, c.snapshot); err != nil {
			return 0, err
		}
		if vAny {
			ctx.hasOpenSource = true
		}
		return 0, nil
	}
	if relOID != 0 {
		aliasCols := extractAliasColnames(rv.Alias)
		if rel := c.relationByOID[relOID]; rel != nil && len(aliasCols) > len(rel.Columns) && !c.snapshot {
			return 0, &Error{
				Code:    CodeInvalidColumnReference,
				Message: fmt.Sprintf("table %q has %d columns available but %d columns specified", name, len(rel.Columns), len(aliasCols)),
			}
		}
		if err := ctx.setScope(name, ruleScopeEntry{
			relOID:     relOID,
			hasAlias:   hasAlias,
			alias:      alias,
			colAliases: aliasCols,
		}, c.snapshot); err != nil {
			return 0, err
		}
	}
	return relOID, nil
}

func (c *Catalog) collectFromScope(b *ruleBindings, list *nodes.List, ctx *ruleEvalContext) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		if err := c.collectTableRefScope(b, item, ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) collectTableRefScope(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext) error {
	switch v := n.(type) {
	case *nodes.RangeVar:
		_, err := addScopeRangeVar(c, b, ctx, v, false)
		return err
	case *nodes.RangeSubselect:
		subCtx := &ruleEvalContext{cteNames: ctx.cteNames, cteCols: ctx.cteCols, cteColTypes: ctx.cteColTypes, cteAny: ctx.cteAny}
		if v.Lateral {
			subCtx.outer = ctx
		}
		if err := c.bindRuleStmt(b, v.Subquery, subCtx); err != nil {
			return err
		}
		baseCols, baseAny := c.extractRuleQueryCols(b, v.Subquery)
		baseTypes, _ := c.extractRuleQueryColTypes(b, v.Subquery)
		aliasCols := extractAliasColnames(v.Alias)
		if !baseAny && len(aliasCols) > len(baseCols) && !c.snapshot {
			aliasName := ""
			if v.Alias != nil {
				aliasName = v.Alias.Aliasname
			}
			return &Error{
				Code:    CodeInvalidColumnReference,
				Message: fmt.Sprintf("table %q has %d columns available but %d columns specified", aliasName, len(baseCols), len(aliasCols)),
			}
		}
		vCols, vTypes, vOrder, vAny := applyVirtualAliasCols(baseCols, baseTypes, baseAny, aliasCols)
		if v.Alias != nil && v.Alias.Aliasname != "" {
			if err := ctx.setScope(v.Alias.Aliasname, ruleScopeEntry{
				hasAlias:        true,
				alias:           v.Alias.Aliasname,
				isVirtual:       true,
				virtualCols:     vCols,
				virtualColTypes: vTypes,
				virtualColOrder: vOrder,
				anyCol:          vAny,
			}, c.snapshot); err != nil {
				return err
			}
		}
		if vAny {
			ctx.hasOpenSource = true
		}
	case *nodes.RangeFunction:
		if v.Functions != nil {
			// A function in FROM sees the items before it with or without
			// LATERAL, for which it is a noise word (parse_clause.c).
			fnCtx := ctx
			for _, fnItem := range v.Functions.Items {
				if pair, ok := fnItem.(*nodes.List); ok && len(pair.Items) > 0 {
					if err := c.bindRuleExpr(b, pair.Items[0], fnCtx); err != nil {
						return err
					}
					if len(pair.Items) > 1 {
						if cdList, ok := pair.Items[1].(*nodes.List); ok {
							if err := c.bindRuleTableFuncColDefs(b, cdList); err != nil {
								return err
							}
						}
					}
				} else if err := c.bindRuleExpr(b, fnItem, fnCtx); err != nil {
					return err
				}
			}
		}
		if err := c.bindRuleTableFuncColDefs(b, v.Coldeflist); err != nil {
			return err
		}
		baseCols, baseAny := c.extractRangeFunctionCols(v)
		aliasCols := extractAliasColnames(v.Alias)
		vCols, vTypes, vOrder, vAny := applyVirtualAliasCols(baseCols, nil, baseAny, aliasCols)
		scopeName := ""
		hasAlias := false
		aliasName := ""
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName = v.Alias.Aliasname
			hasAlias = true
			aliasName = v.Alias.Aliasname
		} else if v.Functions != nil && len(v.Functions.Items) == 1 {
			scopeName = singleRangeFuncName(v.Functions.Items[0])
		}
		if !baseAny && len(aliasCols) > len(baseCols) && !c.snapshot {
			return &Error{
				Code:    CodeInvalidColumnReference,
				Message: fmt.Sprintf("table %q has %d columns available but %d columns specified", scopeName, len(baseCols), len(aliasCols)),
			}
		}
		if scopeName == "" {
			scopeName = fmt.Sprintf("\x00rowsfrom_%d", len(ctx.scope))
		}
		if err := ctx.setScope(scopeName, ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           aliasName,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColTypes: vTypes,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}, c.snapshot); err != nil {
			return err
		}
		if vAny {
			ctx.hasOpenSource = true
		}
	case *nodes.RangeTableFunc:
		fnCtx := ctx
		if !v.Lateral {
			fnCtx = &ruleEvalContext{outer: ctx.outer, cteNames: ctx.cteNames, cteCols: ctx.cteCols, cteColTypes: ctx.cteColTypes, cteAny: ctx.cteAny}
		}
		if v.Namespaces != nil {
			for _, item := range v.Namespaces.Items {
				if rt, ok := item.(*nodes.ResTarget); ok && rt != nil {
					if err := c.bindRuleExpr(b, rt.Val, fnCtx); err != nil {
						return err
					}
				} else if err := c.bindRuleExpr(b, item, fnCtx); err != nil {
					return err
				}
			}
		}
		if err := c.bindRuleExpr(b, v.Docexpr, fnCtx); err != nil {
			return err
		}
		if err := c.bindRuleExpr(b, v.Rowexpr, fnCtx); err != nil {
			return err
		}
		var baseCols []string
		var baseTypes []uint32
		if v.Columns != nil {
			for _, item := range v.Columns.Items {
				fc, ok := item.(*nodes.RangeTableFuncCol)
				if !ok || fc == nil {
					continue
				}
				var colType uint32
				if fc.ForOrdinality {
					colType = INT4OID
				} else if fc.TypeName != nil {
					if err := c.bindRuleTypeName(b, fc.TypeName); err != nil {
						return err
					}
					colType, _, _ = c.resolveTypeName(fc.TypeName)
				}
				if fc.Colname != "" {
					baseCols = append(baseCols, fc.Colname)
					baseTypes = append(baseTypes, colType)
				}
				if err := c.bindRuleExpr(b, fc.Colexpr, fnCtx); err != nil {
					return err
				}
				if err := c.bindRuleExpr(b, fc.Coldefexpr, fnCtx); err != nil {
					return err
				}
			}
		}
		aliasCols := extractAliasColnames(v.Alias)
		scopeName, hasAlias, aliasName := "xmltable", false, ""
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName, hasAlias, aliasName = v.Alias.Aliasname, true, v.Alias.Aliasname
		}
		if len(aliasCols) > len(baseCols) && !c.snapshot {
			return &Error{
				Code:    CodeInvalidColumnReference,
				Message: fmt.Sprintf("table %q has %d columns available but %d columns specified", scopeName, len(baseCols), len(aliasCols)),
			}
		}
		vCols, vTypes, vOrder, vAny := applyVirtualAliasCols(baseCols, baseTypes, false, aliasCols)
		if err := ctx.setScope(scopeName, ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           aliasName,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColTypes: vTypes,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}, c.snapshot); err != nil {
			return err
		}
	case *nodes.JsonTable:
		fnCtx := ctx
		if !v.Lateral {
			fnCtx = &ruleEvalContext{outer: ctx.outer, cteNames: ctx.cteNames, cteCols: ctx.cteCols, cteColTypes: ctx.cteColTypes, cteAny: ctx.cteAny}
		}
		if v.ContextItem != nil {
			if err := c.bindRuleExpr(b, v.ContextItem, fnCtx); err != nil {
				return err
			}
		}
		if v.Pathspec != nil {
			if err := c.bindRuleExpr(b, v.Pathspec.String, fnCtx); err != nil {
				return err
			}
		}
		if v.Passing != nil {
			for _, item := range v.Passing.Items {
				if ja, ok := item.(*nodes.JsonArgument); ok && ja != nil {
					if err := c.bindRuleExpr(b, ja.Val, fnCtx); err != nil {
						return err
					}
				}
			}
		}
		var baseCols []string
		if err := c.bindRuleJsonTableColumns(b, v.Columns, fnCtx, &baseCols); err != nil {
			return err
		}
		if v.OnError != nil {
			if err := c.bindRuleExpr(b, v.OnError.Expr, fnCtx); err != nil {
				return err
			}
		}
		aliasCols := extractAliasColnames(v.Alias)
		scopeName, hasAlias, aliasName := "json_table", false, ""
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName, hasAlias, aliasName = v.Alias.Aliasname, true, v.Alias.Aliasname
		}
		if len(aliasCols) > len(baseCols) && !c.snapshot {
			return &Error{
				Code:    CodeInvalidColumnReference,
				Message: fmt.Sprintf("table %q has %d columns available but %d columns specified", scopeName, len(baseCols), len(aliasCols)),
			}
		}
		vCols, vTypes, vOrder, vAny := applyVirtualAliasCols(baseCols, nil, false, aliasCols)
		if err := ctx.setScope(scopeName, ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           aliasName,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColTypes: vTypes,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}, c.snapshot); err != nil {
			return err
		}
	case *nodes.RangeTableSample:
		if err := c.collectTableRefScope(b, v.Relation, ctx); err != nil {
			return err
		}
		if !c.snapshot {
			rv, isRV := v.Relation.(*nodes.RangeVar)
			if !isRV || rv == nil {
				return &Error{
					Code:    CodeFeatureNotSupported,
					Message: "TABLESAMPLE clause can only be applied to tables and materialized views",
				}
			}
			relOID := uint32(0)
			if b != nil {
				relOID = b.relByRV[rv]
			}
			rel := c.relationByOID[relOID]
			if rel == nil || (rel.RelKind != 'r' && rel.RelKind != 'p' && rel.RelKind != 'm') {
				return &Error{
					Code:    CodeFeatureNotSupported,
					Message: "TABLESAMPLE clause can only be applied to tables and materialized views",
				}
			}
		}
		handlerOID, err := c.resolveTablesampleHandler(v.Method)
		if err != nil {
			if !c.snapshot {
				return err
			}
		} else if handlerOID >= FirstNormalObjectId && b != nil {
			if b.tsmByRTS != nil {
				b.tsmByRTS[v] = handlerOID
			}
			b.addFuncDep(handlerOID)
		}
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		if err := c.bindRuleExpr(b, v.Repeatable, ctx); err != nil {
			return err
		}
		if v.Repeatable != nil && !c.snapshot {
			if t := c.inferRuleExprType(b, v.Repeatable, ctx); t != 0 && t != UNKNOWNOID && !c.CanCoerce(t, FLOAT8OID, 'a') {
				return &Error{
					Code:    CodeDatatypeMismatch,
					Message: fmt.Sprintf("TABLESAMPLE REPEATABLE requires numeric2 or float8, not %s", c.typeName(t)),
				}
			}
		}
	case *nodes.JoinExpr:
		beforeKeys := make(map[string]bool, len(ctx.scope))
		for k := range ctx.scope {
			beforeKeys[k] = true
		}
		if err := c.collectTableRefScope(b, v.Larg, ctx); err != nil {
			return err
		}
		if err := c.collectTableRefScope(b, v.Rarg, ctx); err != nil {
			return err
		}
		lCols, lAny, lRels := c.tableRefColumns(v.Larg, ctx)
		rCols, rAny, rRels := c.tableRefColumns(v.Rarg, ctx)
		aliasCols := extractAliasColnames(v.Alias)
		if v.UsingClause != nil {
			if ctx.usingCols == nil {
				ctx.usingCols = make(map[string]int)
			}
			seenUsing := make(map[string]bool)
			var usingBindings []ruleJoinUsingCol
			for _, u := range v.UsingClause.Items {
				uCol := stringVal(u)
				if uCol == "" {
					continue
				}
				if seenUsing[uCol] && !c.snapshot {
					return &Error{
						Code:    CodeDuplicateColumn,
						Message: fmt.Sprintf("column name %q appears more than once in USING clause", uCol),
					}
				}
				if (!lAny && !stringSliceContains(lCols, uCol)) || (!rAny && !stringSliceContains(rCols, uCol)) {
					if !c.snapshot {
						return errUndefinedColumn(uCol)
					}
				}
				lDeps := c.resolveTableRefColDeps(b, v.Larg, ctx, uCol)
				rDeps := c.resolveTableRefColDeps(b, v.Rarg, ctx, uCol)
				usingBindings = append(usingBindings, ruleJoinUsingCol{
					name: uCol,
					deps: append(lDeps, rDeps...),
				})
				if !seenUsing[uCol] {
					seenUsing[uCol] = true
					if len(aliasCols) == 0 {
						ctx.usingCols[uCol]++
					}
				}
			}
			if b != nil && b.joinUsing != nil {
				b.joinUsing[v] = usingBindings
			}
			if v.JoinUsing != nil && v.JoinUsing.Aliasname != "" {
				juCols := make(map[string]int, len(usingBindings))
				juOrder := make([]string, 0, len(usingBindings))
				for _, ub := range usingBindings {
					if ub.name != "" {
						juCols[ub.name]++
						juOrder = append(juOrder, ub.name)
					}
				}
				if err := ctx.setScope(v.JoinUsing.Aliasname, ruleScopeEntry{
					hasAlias:        true,
					alias:           v.JoinUsing.Aliasname,
					isVirtual:       true,
					virtualCols:     juCols,
					virtualColOrder: juOrder,
					isJoinAlias:     true,
					joinRels:        append(append([]uint32(nil), lRels...), rRels...),
				}, c.snapshot); err != nil {
					return err
				}
			}
		} else if v.IsNatural {
			if ctx.usingCols == nil {
				ctx.usingCols = make(map[string]int)
			}
			seenNat := make(map[string]bool)
			var usingBindings []ruleJoinUsingCol
			for _, lc := range lCols {
				if lc != "" && stringSliceContains(rCols, lc) && !seenNat[lc] {
					seenNat[lc] = true
					lDeps := c.resolveTableRefColDeps(b, v.Larg, ctx, lc)
					rDeps := c.resolveTableRefColDeps(b, v.Rarg, ctx, lc)
					usingBindings = append(usingBindings, ruleJoinUsingCol{
						name: lc,
						deps: append(lDeps, rDeps...),
					})
					if len(aliasCols) == 0 {
						ctx.usingCols[lc]++
					}
				}
			}
			if b != nil && !lAny && !rAny {
				if b.naturalJoined != nil {
					b.naturalJoined[v] = true
				}
				if b.joinUsing != nil {
					b.joinUsing[v] = usingBindings
				}
			}
		}
		if v.Quals != nil {
			// JOIN/ON sees only the join's own inputs and outer query levels,
			// not earlier items of the FROM list (parse_clause.c
			// transformJoinOnClause).
			onCtx := *ctx
			onCtx.scope = make(map[string]ruleScopeEntry, len(ctx.scope))
			onCtx.scopeOrder = nil
			for _, k := range ctx.scopeOrder {
				if !beforeKeys[k] {
					if se, ok := ctx.scope[k]; ok {
						onCtx.scope[k] = se
						onCtx.scopeOrder = append(onCtx.scopeOrder, k)
					}
				}
			}
			if err := c.bindRuleBoolExpr(b, v.Quals, &onCtx, "JOIN/ON"); err != nil {
				return err
			}
		}
		if v.Alias != nil && v.Alias.Aliasname != "" {
			if len(aliasCols) > 0 {
				for k, se := range ctx.scope {
					if !beforeKeys[k] {
						se.hiddenByJoinColAlias = true
						ctx.scope[k] = se
					}
				}
			}
			combinedCols := combineJoinCols(lCols, rCols, v.UsingClause, v.IsNatural)
			if !(lAny || rAny) && len(aliasCols) > len(combinedCols) && !c.snapshot {
				return &Error{
					Code:    CodeInvalidColumnReference,
					Message: fmt.Sprintf("table %q has %d columns available but %d columns specified", v.Alias.Aliasname, len(combinedCols), len(aliasCols)),
				}
			}
			vCols, vTypes, vOrder, vAny := applyVirtualAliasCols(combinedCols, nil, lAny || rAny, aliasCols)
			joinRels := append(append([]uint32(nil), lRels...), rRels...)
			if err := ctx.setScope(v.Alias.Aliasname, ruleScopeEntry{
				hasAlias:          true,
				alias:             v.Alias.Aliasname,
				isVirtual:         true,
				virtualCols:       vCols,
				virtualColTypes:   vTypes,
				virtualColOrder:   vOrder,
				anyCol:            vAny,
				isJoinAlias:       true,
				hasJoinColAliases: len(aliasCols) > 0,
				joinRels:          joinRels,
			}, c.snapshot); err != nil {
				return err
			}
			if vAny {
				ctx.hasOpenSource = true
			}
		}
	}
	return nil
}

func (c *Catalog) bindRuleTypeName(b *ruleBindings, tn *nodes.TypeName) error {
	if tn == nil {
		return nil
	}
	oid, _, err := c.resolveTypeName(tn)
	if err != nil {
		if !c.snapshot {
			return err
		}
		return nil
	}
	if oid != 0 {
		isArr := tn.ArrayBounds != nil && len(tn.ArrayBounds.Items) > 0
		elemOID := oid
		if bt := c.typeByOID[oid]; bt != nil && bt.Elem != 0 {
			elemOID = bt.Elem
		}
		b.typeByTN[tn] = ruleTypeBinding{typeOID: elemOID, isArray: isArr}
		b.addTypeDep(elemOID)
	}
	return nil
}

func (c *Catalog) bindRuleJsonTableColumns(b *ruleBindings, list *nodes.List, ctx *ruleEvalContext, baseCols *[]string) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		jtc, ok := item.(*nodes.JsonTableColumn)
		if !ok || jtc == nil {
			continue
		}
		if jtc.Coltype == nodes.JTC_NESTED {
			if jtc.Pathspec != nil {
				if err := c.bindRuleExpr(b, jtc.Pathspec.String, ctx); err != nil {
					return err
				}
			}
			if err := c.bindRuleJsonTableColumns(b, jtc.Columns, ctx, baseCols); err != nil {
				return err
			}
			continue
		}
		if jtc.Name != "" {
			*baseCols = append(*baseCols, jtc.Name)
		}
		if jtc.TypeName != nil {
			if err := c.bindRuleTypeName(b, jtc.TypeName); err != nil {
				return err
			}
		}
		if jtc.Pathspec != nil {
			if err := c.bindRuleExpr(b, jtc.Pathspec.String, ctx); err != nil {
				return err
			}
		}
		if jtc.OnEmpty != nil {
			if err := c.bindRuleExpr(b, jtc.OnEmpty.Expr, ctx); err != nil {
				return err
			}
		}
		if jtc.OnError != nil {
			if err := c.bindRuleExpr(b, jtc.OnError.Expr, ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func stringSliceContains(xs []string, target string) bool {
	for _, x := range xs {
		if x == target {
			return true
		}
	}
	return false
}

func (c *Catalog) resolveTableRefColDeps(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext, colName string) []ruleColDep {
	if n == nil || ctx == nil || colName == "" {
		return nil
	}
	switch v := n.(type) {
	case *nodes.RangeVar:
		if v == nil {
			return nil
		}
		name := v.Relname
		if v.Alias != nil && v.Alias.Aliasname != "" {
			name = v.Alias.Aliasname
		}
		se, ok := ctx.scope[name]
		if !ok || se.isVirtual {
			return nil
		}
		rel := c.relationByOID[se.relOID]
		if rel == nil {
			return nil
		}
		if attNum, found := resolveScopeColumn(rel, se, colName); found {
			b.addColDep(se.relOID, attNum)
			if len(se.colAliases) == 0 {
				return []ruleColDep{{relOID: se.relOID, attNum: attNum}}
			}
		}
	case *nodes.RangeTableSample:
		return c.resolveTableRefColDeps(b, v.Relation, ctx, colName)
	case *nodes.JoinExpr:
		if v.Alias != nil && v.Alias.Aliasname != "" && len(extractAliasColnames(v.Alias)) > 0 {
			if se, ok := ctx.scope[v.Alias.Aliasname]; ok {
				for _, rOID := range se.joinRels {
					recordAllRelColumnDeps(c, b, rOID)
				}
			}
			return nil
		}
		lDeps := c.resolveTableRefColDeps(b, v.Larg, ctx, colName)
		rDeps := c.resolveTableRefColDeps(b, v.Rarg, ctx, colName)
		return append(lDeps, rDeps...)
	}
	return nil
}

func combineJoinCols(lCols, rCols []string, usingClause *nodes.List, isNatural bool) []string {
	var out []string
	common := make(map[string]bool)
	if usingClause != nil && len(usingClause.Items) > 0 {
		for _, u := range usingClause.Items {
			uCol := stringVal(u)
			if uCol != "" && !common[uCol] {
				common[uCol] = true
				out = append(out, uCol)
			}
		}
	} else if isNatural {
		for _, lc := range lCols {
			if lc != "" && stringSliceContains(rCols, lc) && !common[lc] {
				common[lc] = true
				out = append(out, lc)
			}
		}
	}
	for _, lc := range lCols {
		if !common[lc] {
			out = append(out, lc)
		}
	}
	for _, rc := range rCols {
		if !common[rc] {
			out = append(out, rc)
		}
	}
	return out
}

func (c *Catalog) tableRefColumns(n nodes.Node, ctx *ruleEvalContext) ([]string, bool, []uint32) {
	if n == nil || ctx == nil {
		return nil, false, nil
	}
	switch v := n.(type) {
	case *nodes.RangeVar:
		if v == nil {
			return nil, false, nil
		}
		name := v.Relname
		if v.Alias != nil && v.Alias.Aliasname != "" {
			name = v.Alias.Aliasname
		}
		se, ok := ctx.scope[name]
		if !ok {
			return nil, true, nil
		}
		if se.isVirtual {
			return se.virtualColOrder, se.anyCol, se.joinRels
		}
		rel := c.relationByOID[se.relOID]
		if rel == nil {
			return nil, true, nil
		}
		cols := make([]string, len(rel.Columns))
		for i, col := range rel.Columns {
			if i < len(se.colAliases) && se.colAliases[i] != "" {
				cols[i] = se.colAliases[i]
			} else {
				cols[i] = col.Name
			}
		}
		return cols, false, []uint32{se.relOID}
	case *nodes.RangeSubselect:
		if v.Alias != nil && v.Alias.Aliasname != "" {
			if se, ok := ctx.scope[v.Alias.Aliasname]; ok {
				return se.virtualColOrder, se.anyCol, nil
			}
		}
	case *nodes.RangeFunction:
		scopeName := ""
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName = v.Alias.Aliasname
		} else if v.Functions != nil && len(v.Functions.Items) == 1 {
			scopeName = singleRangeFuncName(v.Functions.Items[0])
		}
		if scopeName != "" {
			if se, ok := ctx.scope[scopeName]; ok {
				return se.virtualColOrder, se.anyCol, nil
			}
		}
	case *nodes.RangeTableFunc:
		scopeName := "xmltable"
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName = v.Alias.Aliasname
		}
		if se, ok := ctx.scope[scopeName]; ok {
			return se.virtualColOrder, se.anyCol, nil
		}
	case *nodes.JsonTable:
		scopeName := "json_table"
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName = v.Alias.Aliasname
		}
		if se, ok := ctx.scope[scopeName]; ok {
			return se.virtualColOrder, se.anyCol, nil
		}
	case *nodes.RangeTableSample:
		return c.tableRefColumns(v.Relation, ctx)
	case *nodes.JoinExpr:
		if v.Alias != nil && v.Alias.Aliasname != "" {
			if se, ok := ctx.scope[v.Alias.Aliasname]; ok {
				return se.virtualColOrder, se.anyCol, se.joinRels
			}
		}
		lCols, lAny, lRels := c.tableRefColumns(v.Larg, ctx)
		rCols, rAny, rRels := c.tableRefColumns(v.Rarg, ctx)
		return combineJoinCols(lCols, rCols, v.UsingClause, v.IsNatural), lAny || rAny, append(append([]uint32(nil), lRels...), rRels...)
	}
	return nil, true, nil
}

func singleRangeFuncName(fnItem nodes.Node) string {
	expr := fnItem
	if pair, ok := fnItem.(*nodes.List); ok && len(pair.Items) > 0 {
		expr = pair.Items[0]
	}
	if fc, ok := expr.(*nodes.FuncCall); ok && fc != nil && fc.Funcname != nil && len(fc.Funcname.Items) > 0 {
		return stringVal(fc.Funcname.Items[len(fc.Funcname.Items)-1])
	}
	if name := figureColName(expr); name != "" && name != "?column?" {
		return name
	}
	return ""
}

func (c *Catalog) extractRangeFunctionCols(v *nodes.RangeFunction) ([]string, bool) {
	if v == nil {
		return nil, false
	}
	var cols []string
	anyCol := false
	if v.Coldeflist != nil && len(v.Coldeflist.Items) > 0 {
		for _, item := range v.Coldeflist.Items {
			if cd, ok := item.(*nodes.ColumnDef); ok && cd != nil {
				cols = append(cols, cd.Colname)
			}
		}
	} else if v.Functions != nil {
		isSingle := !v.IsRowsfrom && len(v.Functions.Items) == 1
		for _, fnItem := range v.Functions.Items {
			var fnExpr nodes.Node
			var cdList *nodes.List
			if pair, ok := fnItem.(*nodes.List); ok && len(pair.Items) > 0 {
				fnExpr = pair.Items[0]
				if len(pair.Items) > 1 {
					cdList, _ = pair.Items[1].(*nodes.List)
				}
			} else {
				fnExpr = fnItem
			}
			if cdList != nil && len(cdList.Items) > 0 {
				for _, item := range cdList.Items {
					if cd, ok := item.(*nodes.ColumnDef); ok && cd != nil {
						cols = append(cols, cd.Colname)
					}
				}
			} else if fc, ok := fnExpr.(*nodes.FuncCall); ok && fc != nil {
				fCols, fAny := c.resolveTableFuncCols(fc, v.Alias, isSingle)
				cols = append(cols, fCols...)
				if fAny {
					anyCol = true
				}
			} else if isSingle && v.Alias != nil && v.Alias.Aliasname != "" {
				cols = append(cols, v.Alias.Aliasname)
			} else if name := figureColName(fnExpr); name != "" && name != "?column?" {
				cols = append(cols, name)
			} else {
				anyCol = true
			}
		}
	}
	if v.Ordinality {
		cols = append(cols, "ordinality")
	}
	return cols, anyCol
}

func (c *Catalog) resolveTableFuncCols(fc *nodes.FuncCall, alias *nodes.Alias, isSingleFunc bool) ([]string, bool) {
	fSchema, fName := qualifiedName(fc.Funcname)
	if fName == "" {
		return nil, true
	}
	var searchNS []uint32
	if fSchema != "" {
		if ns := c.schemaByName[fSchema]; ns != nil {
			searchNS = []uint32{ns.OID}
		}
	} else {
		searchNS = c.routineSearchPathWithCatalog()
	}
	procs := c.procByName[fName]
	var matched []*BuiltinProc
	for _, nsOID := range searchNS {
		var exact []*BuiltinProc
		for _, p := range procs {
			if p.Kind == 'p' {
				continue
			}
			pNS := uint32(PGCatalogNamespace)
			if p.OID >= FirstNormalObjectId {
				up := c.userProcs[p.OID]
				if up == nil || up.Schema == nil {
					continue
				}
				pNS = up.Schema.OID
			}
			if pNS != nsOID {
				continue
			}
			if _, ok := procMatchesRuleCallArity(p, fc); ok {
				exact = append(exact, p)
			}
		}
		if len(exact) > 0 {
			matched = exact
			break
		}
	}
	if len(matched) == 0 {
		return nil, true
	}
	var cols []string
	seen := make(map[string]bool)
	addCol := func(name string) {
		if name != "" && !seen[name] {
			seen[name] = true
			cols = append(cols, name)
		}
	}
	for _, p := range matched {
		var argModes []byte
		var argNames []string
		if up := c.userProcs[p.OID]; up != nil {
			argModes = up.ArgModes
			argNames = up.ArgNames
		} else {
			argModes = p.ArgModes
			argNames = p.ArgNames
		}
		outCount := 0
		for i, m := range argModes {
			if m == 'o' || m == 'b' || m == 't' {
				outCount++
				if i < len(argNames) && argNames[i] != "" {
					addCol(argNames[i])
				} else {
					return nil, true
				}
			}
		}
		if outCount > 0 {
			continue
		}
		if bt := c.typeByOID[p.RetType]; bt != nil && bt.RelID != 0 {
			if rel := c.relationByOID[bt.RelID]; rel != nil {
				for _, col := range rel.Columns {
					addCol(col.Name)
				}
				continue
			}
		}
		if p.RetType == RECORDOID {
			return nil, true
		}
		if isSingleFunc && alias != nil && alias.Aliasname != "" {
			addCol(alias.Aliasname)
		} else {
			addCol(fName)
		}
	}
	return cols, false
}

func (c *Catalog) extractRuleQueryCols(b *ruleBindings, n nodes.Node) ([]string, bool) {
	if n == nil {
		return nil, false
	}
	switch s := n.(type) {
	case *nodes.SelectStmt:
		if s.Op != nodes.SETOP_NONE {
			return c.extractRuleQueryCols(b, s.Larg)
		}
		if s.ValuesLists != nil {
			if len(s.ValuesLists.Items) > 0 {
				if row, ok := s.ValuesLists.Items[0].(*nodes.List); ok {
					cols := make([]string, len(row.Items))
					for i := range row.Items {
						cols[i] = fmt.Sprintf("column%d", i+1)
					}
					return cols, false
				}
			}
			return nil, true
		}
		if s.TargetList == nil {
			return nil, false
		}
		var ctx *ruleEvalContext
		if b != nil {
			ctx = b.stmtCtx[s]
		}
		var cols []string
		anyCol := false
		for _, item := range s.TargetList.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok || rt == nil {
				continue
			}
			if rt.Name != "" {
				cols = append(cols, rt.Name)
				continue
			}
			if b != nil && b.starExpandByRT != nil {
				if starCols, hasStar := b.starExpandByRT[rt]; hasStar && len(starCols) > 0 {
					for _, sc := range starCols {
						cols = append(cols, sc.origName)
					}
					continue
				}
			}
			if cr, ok := rt.Val.(*nodes.ColumnRef); ok && cr.Fields != nil && len(cr.Fields.Items) > 0 {
				fields := cr.Fields.Items
				if _, isStar := fields[len(fields)-1].(*nodes.A_Star); isStar {
					if len(fields) == 1 {
						if ctx == nil || s.FromClause == nil {
							anyCol = true
						} else {
							for _, fromItem := range s.FromClause.Items {
								fCols, fAny, _ := c.tableRefColumns(fromItem, ctx)
								cols = append(cols, fCols...)
								if fAny {
									anyCol = true
								}
							}
						}
					} else {
						qual := stringVal(fields[len(fields)-2])
						if se, ok := lookupScopeEntry(ctx, qual); ok {
							if se.isVirtual {
								cols = append(cols, se.virtualColOrder...)
								if se.anyCol {
									anyCol = true
								}
							} else if rel := c.relationByOID[se.relOID]; rel != nil {
								for i, col := range rel.Columns {
									if i < len(se.colAliases) && se.colAliases[i] != "" {
										cols = append(cols, se.colAliases[i])
									} else {
										cols = append(cols, col.Name)
									}
								}
							} else {
								anyCol = true
							}
						} else {
							anyCol = true
						}
					}
					continue
				}
				cols = append(cols, stringVal(fields[len(fields)-1]))
				continue
			}
			cols = append(cols, figureColName(rt.Val))
		}
		return cols, anyCol
	case *nodes.InsertStmt:
		return c.extractReturningTargetCols(s.ReturningList, s.Relation)
	case *nodes.UpdateStmt:
		return c.extractReturningTargetCols(s.ReturningList, s.Relation)
	case *nodes.DeleteStmt:
		return c.extractReturningTargetCols(s.ReturningList, s.Relation)
	case *nodes.MergeStmt:
		return c.extractReturningTargetCols(s.ReturningList, s.Relation)
	}
	return nil, true
}

func (c *Catalog) extractReturningTargetCols(list *nodes.List, targetRV *nodes.RangeVar) ([]string, bool) {
	if list == nil {
		return nil, false
	}
	var cols []string
	anyCol := false
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok || rt == nil {
			continue
		}
		if rt.Name != "" {
			cols = append(cols, rt.Name)
			continue
		}
		if cr, ok := rt.Val.(*nodes.ColumnRef); ok && cr.Fields != nil && len(cr.Fields.Items) > 0 {
			fields := cr.Fields.Items
			if _, isStar := fields[len(fields)-1].(*nodes.A_Star); isStar {
				if targetRV != nil {
					if _, rel, err := c.findRelation(targetRV.Schemaname, targetRV.Relname); err == nil && rel != nil {
						for _, col := range rel.Columns {
							cols = append(cols, col.Name)
						}
						continue
					}
				}
				anyCol = true
				continue
			}
			cols = append(cols, stringVal(fields[len(fields)-1]))
			continue
		}
		cols = append(cols, figureColName(rt.Val))
	}
	return cols, anyCol
}

func (c *Catalog) bindWithClause(b *ruleBindings, w *nodes.WithClause, parent *ruleEvalContext) (map[string]bool, map[string][]string, map[string][]uint32, map[string]bool, error) {
	cteNames := make(map[string]bool)
	cteCols := make(map[string][]string)
	cteColTypes := make(map[string][]uint32)
	cteAny := make(map[string]bool)
	if parent != nil {
		for k, v := range parent.cteNames {
			cteNames[k] = v
		}
		for k, v := range parent.cteCols {
			cteCols[k] = v
		}
		for k, v := range parent.cteColTypes {
			cteColTypes[k] = v
		}
		for k, v := range parent.cteAny {
			cteAny[k] = v
		}
	}
	if w == nil || w.Ctes == nil {
		return cteNames, cteCols, cteColTypes, cteAny, nil
	}
	if w.Recursive && !c.snapshot {
		if err := checkRuleRecursiveWith(w); err != nil {
			return nil, nil, nil, nil, err
		}
	}
	if w.Recursive {
		for _, item := range w.Ctes.Items {
			if cte, ok := item.(*nodes.CommonTableExpr); ok && cte.Ctename != "" {
				cteNames[cte.Ctename] = true
				if cte.Aliascolnames != nil && len(cte.Aliascolnames.Items) > 0 {
					cteCols[cte.Ctename] = extractNodeListStrings(cte.Aliascolnames)
				} else {
					cteAny[cte.Ctename] = true
				}
			}
		}
	}
	seenCTEName := make(map[string]bool)
	for _, item := range w.Ctes.Items {
		if cte, ok := item.(*nodes.CommonTableExpr); ok {
			if cte.Ctename != "" {
				if seenCTEName[cte.Ctename] && !c.snapshot {
					return nil, nil, nil, nil, &Error{
						Code:    CodeDuplicateAlias,
						Message: fmt.Sprintf("WITH query name %q specified more than once", cte.Ctename),
					}
				}
				seenCTEName[cte.Ctename] = true
			}
			var outer *ruleEvalContext
			if parent != nil {
				outer = parent
			}
			cteCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteColTypes: cteColTypes, cteAny: cteAny}
			if err := c.bindRuleStmt(b, cte.Ctequery, cteCtx); err != nil {
				return nil, nil, nil, nil, err
			}
			if cc, ok := cte.CycleClause.(*nodes.CTECycleClause); ok && cc != nil {
				if err := c.bindRuleExpr(b, cc.CycleMarkValue, cteCtx); err != nil {
					return nil, nil, nil, nil, err
				}
				if err := c.bindRuleExpr(b, cc.CycleMarkDefault, cteCtx); err != nil {
					return nil, nil, nil, nil, err
				}
				if cc.CycleMarkValue != nil && cc.CycleMarkDefault != nil && !c.snapshot {
					vt := c.inferRuleExprType(b, cc.CycleMarkValue, cteCtx)
					dt := c.inferRuleExprType(b, cc.CycleMarkDefault, cteCtx)
					if vt != 0 && dt != 0 && vt != UNKNOWNOID && dt != UNKNOWNOID && vt != dt {
						if _, err := c.selectCommonType([]AnalyzedExpr{
							&ConstExpr{TypeOID: vt, TypeMod: -1},
							&ConstExpr{TypeOID: dt, TypeMod: -1},
						}, "CYCLE"); err != nil {
							return nil, nil, nil, nil, errDatatypeMismatch(err.Error())
						}
					}
				}
			}
			if cte.Ctename != "" {
				cteNames[cte.Ctename] = true
				baseCols, baseAny := c.extractRuleQueryCols(b, cte.Ctequery)
				baseTypes, _ := c.extractRuleQueryColTypes(b, cte.Ctequery)
				aliasCols := extractNodeListStrings(cte.Aliascolnames)
				if !baseAny && len(aliasCols) > len(baseCols) && !c.snapshot {
					return nil, nil, nil, nil, &Error{
						Code:    CodeInvalidColumnReference,
						Message: fmt.Sprintf("WITH query %q has %d columns available but %d columns specified", cte.Ctename, len(baseCols), len(aliasCols)),
					}
				}
				_, _, vOrder, vAny := applyVirtualAliasCols(baseCols, baseTypes, baseAny, aliasCols)
				cteOutputCols := append([]string(nil), vOrder...)
				if sc, ok := cte.SearchClause.(*nodes.CTESearchClause); ok && sc != nil {
					if !w.Recursive && !c.snapshot {
						return nil, nil, nil, nil, &Error{
							Code:    CodeSyntaxError,
							Message: "SEARCH clause can only be specified for recursive WITH queries",
						}
					}
					if sc.SearchColList != nil {
						seenSearch := make(map[string]bool, len(sc.SearchColList.Items))
						for _, colNode := range sc.SearchColList.Items {
							cName := stringVal(colNode)
							if !vAny && !stringSliceContains(cteOutputCols, cName) && !c.snapshot {
								return nil, nil, nil, nil, errUndefinedColumn(cName)
							}
							if seenSearch[cName] && !c.snapshot {
								return nil, nil, nil, nil, errDuplicateColumn(cName)
							}
							seenSearch[cName] = true
						}
					}
					if sc.SearchSeqColumn != "" {
						if !vAny && stringSliceContains(vOrder, sc.SearchSeqColumn) && !c.snapshot {
							return nil, nil, nil, nil, errDuplicateColumn(sc.SearchSeqColumn)
						}
						vOrder = append(vOrder, sc.SearchSeqColumn)
						baseTypes = append(baseTypes, RECORDOID)
					}
				}
				if cc, ok := cte.CycleClause.(*nodes.CTECycleClause); ok && cc != nil {
					if !w.Recursive && !c.snapshot {
						return nil, nil, nil, nil, &Error{
							Code:    CodeSyntaxError,
							Message: "CYCLE clause can only be specified for recursive WITH queries",
						}
					}
					if cc.CycleColList != nil {
						seenCycle := make(map[string]bool, len(cc.CycleColList.Items))
						for _, colNode := range cc.CycleColList.Items {
							cName := stringVal(colNode)
							if !vAny && !stringSliceContains(cteOutputCols, cName) && !c.snapshot {
								return nil, nil, nil, nil, errUndefinedColumn(cName)
							}
							if seenCycle[cName] && !c.snapshot {
								return nil, nil, nil, nil, errDuplicateColumn(cName)
							}
							seenCycle[cName] = true
						}
					}
					if cc.CycleMarkColumn != "" {
						if !vAny && stringSliceContains(vOrder, cc.CycleMarkColumn) && !c.snapshot {
							return nil, nil, nil, nil, errDuplicateColumn(cc.CycleMarkColumn)
						}
						vOrder = append(vOrder, cc.CycleMarkColumn)
						baseTypes = append(baseTypes, BOOLOID)
					}
					if cc.CyclePathColumn != "" {
						if !vAny && stringSliceContains(vOrder, cc.CyclePathColumn) && !c.snapshot {
							return nil, nil, nil, nil, errDuplicateColumn(cc.CyclePathColumn)
						}
						vOrder = append(vOrder, cc.CyclePathColumn)
						baseTypes = append(baseTypes, RECORDARRAYOID)
					}
				}
				cteCols[cte.Ctename] = vOrder
				cteColTypes[cte.Ctename] = baseTypes
				cteAny[cte.Ctename] = vAny
			}
		}
	}
	return cteNames, cteCols, cteColTypes, cteAny, nil
}

// checkRuleRecursiveWith applies parse_cte.c's checks to a WITH RECURSIVE list:
// no cycle of references between two or more items, and an item that refers
// to itself has the form non-recursive-term UNION [ALL] recursive-term with no
// self reference in the non-recursive term.
func checkRuleRecursiveWith(w *nodes.WithClause) error {
	var ctes []*nodes.CommonTableExpr
	index := make(map[string]int)
	for _, item := range w.Ctes.Items {
		if cte, ok := item.(*nodes.CommonTableExpr); ok && cte.Ctename != "" {
			if _, dup := index[cte.Ctename]; !dup {
				index[cte.Ctename] = len(ctes)
				ctes = append(ctes, cte)
			}
		}
	}
	refs := make([]map[int]bool, len(ctes))
	for i, cte := range ctes {
		refs[i] = ruleCTERefs(cte.Ctequery, index)
	}
	// Iterative depth-first search for a cycle that does not run through a
	// self reference alone.
	state := make([]int, len(ctes))
	for start := range ctes {
		if state[start] != 0 {
			continue
		}
		type frame struct {
			node int
			next []int
		}
		succ := func(i int) []int {
			var out []int
			for j := range refs[i] {
				if j != i {
					out = append(out, j)
				}
			}
			return out
		}
		state[start] = 1
		stack := []frame{{start, succ(start)}}
		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			if len(top.next) == 0 {
				state[top.node] = 2
				stack = stack[:len(stack)-1]
				continue
			}
			j := top.next[0]
			top.next = top.next[1:]
			switch state[j] {
			case 1:
				return &Error{Code: CodeFeatureNotSupported, Message: "mutual recursion between WITH items is not implemented"}
			case 0:
				state[j] = 1
				stack = append(stack, frame{j, succ(j)})
			}
		}
	}
	for i, cte := range ctes {
		if !refs[i][i] {
			continue
		}
		sel, ok := cte.Ctequery.(*nodes.SelectStmt)
		if !ok || sel.Op != nodes.SETOP_UNION {
			return &Error{
				Code:    CodeInvalidRecursion,
				Message: fmt.Sprintf("recursive query %q does not have the form non-recursive-term UNION [ALL] recursive-term", cte.Ctename),
			}
		}
		if ruleCTERefs(sel.Larg, index)[i] {
			return &Error{
				Code:    CodeInvalidRecursion,
				Message: fmt.Sprintf("recursive reference to query %q must not appear within its non-recursive term", cte.Ctename),
			}
		}
	}
	return nil
}

// ruleCTERefs returns the items of index that n names as an unqualified
// relation. A nested WITH item of the same name shadows the outer one within
// its scope.
func ruleCTERefs(n nodes.Node, index map[string]int) map[int]bool {
	out := make(map[int]bool)
	nodes.Walk(ruleCTERefFinder{index: index, out: out}, n)
	return out
}

// ruleCTERefFinder walks a query for references to the items of index.
// shadow holds the nested WITH names in scope, and outer the names in scope
// before the WITH clause of the statement being walked.
type ruleCTERefFinder struct {
	index  map[string]int
	out    map[int]bool
	shadow map[string]bool
	outer  map[string]bool
}

func (f ruleCTERefFinder) Visit(n nodes.Node) nodes.Visitor {
	if n == nil {
		return nil
	}
	switch v := n.(type) {
	case *nodes.RangeVar:
		if v.Schemaname == "" && v.Catalogname == "" && !f.shadow[v.Relname] {
			if i, found := f.index[v.Relname]; found {
				f.out[i] = true
			}
		}
	case *nodes.WithClause:
		// A WITH RECURSIVE item sees every name of its list; another item
		// sees only the items before it.
		var names []string
		if v.Ctes != nil {
			for _, item := range v.Ctes.Items {
				if cte, ok := item.(*nodes.CommonTableExpr); ok {
					names = append(names, cte.Ctename)
				}
			}
			for i, item := range v.Ctes.Items {
				cte, ok := item.(*nodes.CommonTableExpr)
				if !ok {
					continue
				}
				visible := names[:i]
				if v.Recursive {
					visible = names
				}
				scope := withShadowNames(f.outer, visible)
				nodes.Walk(ruleCTERefFinder{index: f.index, out: f.out, shadow: scope, outer: scope}, cte.Ctequery)
			}
		}
		return nil
	default:
		if w := ruleStmtWithClause(n); w != nil && w.Ctes != nil {
			var names []string
			for _, item := range w.Ctes.Items {
				if cte, ok := item.(*nodes.CommonTableExpr); ok {
					names = append(names, cte.Ctename)
				}
			}
			return ruleCTERefFinder{index: f.index, out: f.out, shadow: withShadowNames(f.shadow, names), outer: f.shadow}
		}
	}
	return f
}

func withShadowNames(base map[string]bool, names []string) map[string]bool {
	out := make(map[string]bool, len(base)+len(names))
	for k := range base {
		out[k] = true
	}
	for _, name := range names {
		out[name] = true
	}
	return out
}

func ruleStmtWithClause(n nodes.Node) *nodes.WithClause {
	switch v := n.(type) {
	case *nodes.SelectStmt:
		return v.WithClause
	case *nodes.InsertStmt:
		return v.WithClause
	case *nodes.UpdateStmt:
		return v.WithClause
	case *nodes.DeleteStmt:
		return v.WithClause
	case *nodes.MergeStmt:
		return v.WithClause
	}
	return nil
}

func (c *Catalog) bindTargetColumns(b *ruleBindings, relOID uint32, list *nodes.List) error {
	if relOID == 0 || list == nil {
		return nil
	}
	rel := c.relationByOID[relOID]
	if rel == nil {
		return nil
	}
	seenDirectCol := make(map[int16]bool)
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok || rt.Name == "" {
			continue
		}
		idx, exists := rel.colByName[rt.Name]
		if !exists {
			if c.snapshot {
				continue
			}
			return errUndefinedColumn(rt.Name)
		}
		attNum := rel.Columns[idx].AttNum
		if (rt.Indirection == nil || len(rt.Indirection.Items) == 0) && !c.snapshot {
			if seenDirectCol[attNum] {
				return &Error{
					Code:    CodeSyntaxError,
					Message: fmt.Sprintf("column %q specified more than once", rt.Name),
				}
			}
			seenDirectCol[attNum] = true
		}
		b.colByRT[rt] = ruleColBinding{relOID: relOID, attNum: attNum}
		b.addColDep(relOID, attNum)
	}
	return nil
}

func (c *Catalog) checkRuleTypeAssignable(col *Column, valType uint32) error {
	if c.snapshot || col == nil || col.TypeOID == 0 || col.TypeOID == UNKNOWNOID || valType == 0 || valType == UNKNOWNOID {
		return nil
	}
	if !c.CanCoerce(valType, col.TypeOID, 'a') {
		return &Error{
			Code:    CodeDatatypeMismatch,
			Message: fmt.Sprintf("column %q is of type %s but expression is of type %s", col.Name, c.typeName(col.TypeOID), c.typeName(valType)),
		}
	}
	return nil
}

func (c *Catalog) checkRuleTargetListTypes(b *ruleBindings, relOID uint32, list *nodes.List, ctx *ruleEvalContext) error {
	if c.snapshot || relOID == 0 || list == nil {
		return nil
	}
	rel := c.relationByOID[relOID]
	if rel == nil {
		return nil
	}
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok || rt == nil || rt.Name == "" || (rt.Indirection != nil && len(rt.Indirection.Items) > 0) {
			continue
		}
		idx, exists := rel.colByName[rt.Name]
		if !exists {
			continue
		}
		valType := c.inferRuleExprType(b, rt.Val, ctx)
		if err := c.checkRuleTypeAssignable(rel.Columns[idx], valType); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) resolveInsertTargetCols(relOID uint32, cols *nodes.List) []*Column {
	if relOID == 0 {
		return nil
	}
	rel := c.relationByOID[relOID]
	if rel == nil {
		return nil
	}
	if cols == nil || len(cols.Items) == 0 {
		return rel.Columns
	}
	out := make([]*Column, len(cols.Items))
	for i, item := range cols.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok || rt == nil || rt.Name == "" || (rt.Indirection != nil && len(rt.Indirection.Items) > 0) {
			continue
		}
		if idx, exists := rel.colByName[rt.Name]; exists {
			out[i] = rel.Columns[idx]
		}
	}
	return out
}

func (c *Catalog) checkRuleOrdinalExpr(n nodes.Node, numCols int, hasUnknownCols bool, clause string) error {
	if c.snapshot || hasUnknownCols || n == nil {
		return nil
	}
	var iv *nodes.Integer
	switch v := n.(type) {
	case *nodes.A_Const:
		if !v.Isnull {
			iv, _ = v.Val.(*nodes.Integer)
		}
	case *nodes.Integer:
		iv = v
	}
	if iv == nil {
		return nil
	}
	if iv.Ival < 1 || int(iv.Ival) > numCols {
		return &Error{
			Code:    CodeInvalidColumnReference,
			Message: fmt.Sprintf("%s position %d is not in select list", clause, iv.Ival),
		}
	}
	return nil
}

func recordPositionalInsertTargetDeps(c *Catalog, b *ruleBindings, targetOID uint32, selectNode nodes.Node) error {
	if targetOID == 0 || selectNode == nil {
		return nil
	}
	rel := c.relationByOID[targetOID]
	if rel == nil {
		return nil
	}
	sel, ok := selectNode.(*nodes.SelectStmt)
	if !ok {
		return nil
	}
	maxCols := 0
	if sel.ValuesLists != nil {
		for _, rowNode := range sel.ValuesLists.Items {
			if row, ok := rowNode.(*nodes.List); ok {
				if len(row.Items) > len(rel.Columns) && !c.snapshot {
					return &Error{
						Code:    CodeSyntaxError,
						Message: "INSERT has more expressions than target columns",
					}
				}
				if len(row.Items) > maxCols {
					maxCols = len(row.Items)
				}
			}
		}
	} else if sel.Op != nodes.SETOP_NONE {
		qCols, qAny := c.extractRuleQueryCols(b, sel)
		if !qAny {
			if len(qCols) > len(rel.Columns) && !c.snapshot {
				return &Error{
					Code:    CodeSyntaxError,
					Message: "INSERT has more expressions than target columns",
				}
			}
			maxCols = len(qCols)
		} else {
			maxCols = len(rel.Columns)
		}
	} else if sel.TargetList != nil {
		qCols, qAny := c.extractRuleQueryCols(b, sel)
		if !qAny {
			if len(qCols) > len(rel.Columns) && !c.snapshot {
				return &Error{
					Code:    CodeSyntaxError,
					Message: "INSERT has more expressions than target columns",
				}
			}
			maxCols = len(qCols)
		} else {
			maxCols = len(rel.Columns)
		}
	}
	for i := 0; i < maxCols && i < len(rel.Columns); i++ {
		b.addColDep(targetOID, rel.Columns[i].AttNum)
	}
	return nil
}

func (c *Catalog) bindRuleStmt(b *ruleBindings, n nodes.Node, parent *ruleEvalContext) error {
	if n == nil {
		return nil
	}
	switch s := n.(type) {
	case *nodes.SelectStmt:
		if s.IntoClause != nil && !c.snapshot {
			return &Error{
				Code:    CodeSyntaxError,
				Message: "SELECT ... INTO is not allowed here",
			}
		}
		cteNames, cteCols, cteColTypes, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		var outer *ruleEvalContext
		if parent != nil {
			outer = parent.outer
			if len(parent.scope) > 0 || parent.hasOpenSource {
				outer = parent
			}
		}
		if s.Op != nodes.SETOP_NONE {
			branchCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteColTypes: cteColTypes, cteAny: cteAny}
			if s.Larg != nil {
				if err := c.bindRuleStmt(b, s.Larg, branchCtx); err != nil {
					return err
				}
			}
			if s.Rarg != nil {
				if err := c.bindRuleStmt(b, s.Rarg, branchCtx); err != nil {
					return err
				}
			}
			if !c.snapshot {
				lTypes, lAny := c.extractRuleQueryColTypes(b, s.Larg)
				rTypes, rAny := c.extractRuleQueryColTypes(b, s.Rarg)
				if !lAny && !rAny {
					if len(lTypes) != len(rTypes) {
						return &Error{
							Code:    CodeSyntaxError,
							Message: fmt.Sprintf("each %s query must have the same number of columns", setOpKeyword(s.Op)),
						}
					}
					for i := range lTypes {
						lt, rt := lTypes[i], rTypes[i]
						if lt != 0 && rt != 0 && lt != UNKNOWNOID && rt != UNKNOWNOID && lt != rt {
							if _, err := c.selectCommonType([]AnalyzedExpr{
								&ConstExpr{TypeOID: lt, TypeMod: -1},
								&ConstExpr{TypeOID: rt, TypeMod: -1},
							}, setOpKeyword(s.Op)); err != nil {
								return errDatatypeMismatch(err.Error())
							}
						}
					}
				}
			}
			qCols, qAny := c.extractRuleQueryCols(b, s)
			setOpAliases := make(map[string]bool, len(qCols))
			for _, col := range qCols {
				if col != "" {
					setOpAliases[col] = true
				}
			}
			sortCtx := &ruleEvalContext{
				outer:         outer,
				cteNames:      cteNames,
				cteCols:       cteCols,
				cteColTypes:   cteColTypes,
				cteAny:        cteAny,
				targetAliases: setOpAliases,
				hasOpenSource: qAny,
			}
			return c.bindSortAndLimit(b, s, sortCtx)
		}
		if s.ValuesLists != nil {
			valCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteColTypes: cteColTypes, cteAny: cteAny, disallowAgg: "VALUES", disallowWin: "VALUES", disallowSRF: "VALUES"}
			if parent != nil && parent.insertSource && len(s.ValuesLists.Items) == 1 {
				valCtx.disallowSRF = ""
			}
			firstRowLen := -1
			var colExprs [][]AnalyzedExpr
			for _, rowNode := range s.ValuesLists.Items {
				if row, ok := rowNode.(*nodes.List); ok {
					if firstRowLen < 0 {
						firstRowLen = len(row.Items)
						colExprs = make([][]AnalyzedExpr, firstRowLen)
					} else if len(row.Items) != firstRowLen && !c.snapshot {
						return &Error{
							Code:    CodeSyntaxError,
							Message: "VALUES lists must all be the same length",
						}
					}
					for i, expr := range row.Items {
						if err := c.bindRuleExpr(b, expr, valCtx); err != nil {
							return err
						}
						if i < len(colExprs) {
							if t := c.inferRuleExprType(b, expr, valCtx); t != 0 && t != UNKNOWNOID {
								colExprs[i] = append(colExprs[i], &ConstExpr{TypeOID: t, TypeMod: -1})
							}
						}
					}
				}
			}
			if !c.snapshot {
				for _, exprs := range colExprs {
					if len(exprs) > 1 {
						if _, err := c.selectCommonType(exprs, "VALUES"); err != nil {
							return errDatatypeMismatch(err.Error())
						}
					}
				}
			}
			qCols, qAny := c.extractRuleQueryCols(b, s)
			setOpAliases := make(map[string]bool, len(qCols))
			for _, col := range qCols {
				if col != "" {
					setOpAliases[col] = true
				}
			}
			sortCtx := &ruleEvalContext{
				outer:         outer,
				cteNames:      cteNames,
				cteCols:       cteCols,
				cteColTypes:   cteColTypes,
				cteAny:        cteAny,
				targetAliases: setOpAliases,
				hasOpenSource: qAny,
			}
			return c.bindSortAndLimit(b, s, sortCtx)
		}
		windowNames := make(map[string]bool)
		windowDefs := make(map[string]*nodes.WindowDef)
		if s.WindowClause != nil {
			for _, item := range s.WindowClause.Items {
				if wd, ok := item.(*nodes.WindowDef); ok && wd != nil && wd.Name != "" {
					windowDefs[wd.Name] = wd
					if windowNames[wd.Name] && !c.snapshot {
						return &Error{
							Code:    CodeWindowingError,
							Message: fmt.Sprintf("window %q is already defined", wd.Name),
						}
					}
					windowNames[wd.Name] = true
				}
			}
		}
		ctx := &ruleEvalContext{
			scope:       make(map[string]ruleScopeEntry),
			outer:       outer,
			cteNames:    cteNames,
			cteCols:     cteCols,
			cteColTypes: cteColTypes,
			cteAny:      cteAny,
			windowNames: windowNames,
			windowDefs:  windowDefs,
		}
		if err := c.collectFromScope(b, s.FromClause, ctx); err != nil {
			return err
		}
		if b != nil && b.stmtCtx != nil {
			b.stmtCtx[s] = ctx
		}
		targetAliases := make(map[string]bool)
		if s.TargetList != nil {
			for _, item := range s.TargetList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					if rt.Name != "" {
						targetAliases[rt.Name] = true
					} else if b != nil && b.targetNameByRT != nil {
						if inferred := inferRuleTargetColName(rt.Val); inferred != "" && inferred != "?column?" {
							b.targetNameByRT[rt] = inferred
						}
					}
					if err := c.bindRuleExpr(b, rt.Val, ctx); err != nil {
						return err
					}
					if b != nil && b.starExpandByRT != nil {
						if cr, isCR := rt.Val.(*nodes.ColumnRef); isCR && isStarColumnRef(cr) {
							if starCols, anyOpen := c.expandRuleStarTarget(b, cr, s.FromClause, 0, ctx); !anyOpen && len(starCols) > 0 {
								b.starExpandByRT[rt] = starCols
							}
						}
					}
				} else if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		if err := c.bindRuleBoolExpr(b, s.WhereClause, ctx, "WHERE"); err != nil {
			return err
		}
		aliasCtx := &ruleEvalContext{
			scope:         ctx.scope,
			scopeOrder:    ctx.scopeOrder,
			outer:         ctx.outer,
			cteNames:      ctx.cteNames,
			cteCols:       ctx.cteCols,
			cteColTypes:   ctx.cteColTypes,
			cteAny:        ctx.cteAny,
			usingCols:     ctx.usingCols,
			targetAliases: targetAliases,
			windowNames:   windowNames,
			windowDefs:    windowDefs,
			hasOpenSource: ctx.hasOpenSource,
		}
		if s.DistinctClause != nil {
			qCols, qAny := c.extractRuleQueryCols(b, s)
			for _, item := range s.DistinctClause.Items {
				if err := c.bindRuleExpr(b, item, aliasCtx); err != nil {
					return err
				}
				if item != nil {
					if err := c.checkRuleOrdinalExpr(item, len(qCols), qAny, "DISTINCT ON"); err != nil {
						return err
					}
				}
			}
		}
		if s.GroupClause != nil {
			qCols, qAny := c.extractRuleQueryCols(b, s)
			aliasCtx.disallowAgg, aliasCtx.disallowWin = "GROUP BY", "GROUP BY"
			for _, item := range s.GroupClause.Items {
				if err := c.bindRuleExpr(b, item, aliasCtx); err != nil {
					return err
				}
				if err := c.checkRuleOrdinalExpr(item, len(qCols), qAny, "GROUP BY"); err != nil {
					return err
				}
			}
			aliasCtx.disallowAgg, aliasCtx.disallowWin = "", ""
		}
		if err := c.bindRuleBoolExpr(b, s.HavingClause, ctx, "HAVING"); err != nil {
			return err
		}
		if s.WindowClause != nil {
			// A WINDOW entry may copy only an entry declared before it
			// (parse_clause.c transformWindowDefinitions).
			declared := make(map[string]bool)
			for _, item := range s.WindowClause.Items {
				if wd, ok := item.(*nodes.WindowDef); ok {
					if wd.Refname != "" && !declared[wd.Refname] && !c.snapshot {
						return &Error{
							Code:    CodeUndefinedObject,
							Message: fmt.Sprintf("window %q does not exist", wd.Refname),
						}
					}
					if err := c.bindWindowClauseDef(b, wd, ctx); err != nil {
						return err
					}
					if wd.Name != "" {
						declared[wd.Name] = true
					}
				}
			}
		}
		return c.bindSortAndLimit(b, s, aliasCtx)

	case *nodes.InsertStmt:
		cteNames, cteCols, cteColTypes, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		targetOID, _, err := c.bindRangeVar(b, s.Relation, cteNames, true)
		if err != nil {
			return err
		}
		if err := c.bindTargetColumns(b, targetOID, s.Cols); err != nil {
			return err
		}
		if s.SelectStmt != nil {
			if err := c.bindRuleStmt(b, s.SelectStmt, &ruleEvalContext{cteNames: cteNames, cteCols: cteCols, cteColTypes: cteColTypes, cteAny: cteAny, insertSource: true}); err != nil {
				return err
			}
			if s.Cols == nil {
				if err := recordPositionalInsertTargetDeps(c, b, targetOID, s.SelectStmt); err != nil {
					return err
				}
			} else if targetOID != 0 && !c.snapshot {
				targetCols := len(s.Cols.Items)
				if sel, ok := s.SelectStmt.(*nodes.SelectStmt); ok && sel != nil {
					if sel.ValuesLists != nil {
						for _, rowNode := range sel.ValuesLists.Items {
							if row, ok := rowNode.(*nodes.List); ok {
								if len(row.Items) > targetCols {
									return &Error{
										Code:    CodeSyntaxError,
										Message: "INSERT has more expressions than target columns",
									}
								}
								if len(row.Items) < targetCols {
									return &Error{
										Code:    CodeSyntaxError,
										Message: "INSERT has more target columns than expressions",
									}
								}
							}
						}
					} else {
						qCols, qAny := c.extractRuleQueryCols(b, sel)
						if !qAny {
							if len(qCols) > targetCols {
								return &Error{
									Code:    CodeSyntaxError,
									Message: "INSERT has more expressions than target columns",
								}
							}
							if len(qCols) < targetCols {
								return &Error{
									Code:    CodeSyntaxError,
									Message: "INSERT has more target columns than expressions",
								}
							}
						}
					}
				}
			}
			if targetOID != 0 && !c.snapshot {
				targetColDefs := c.resolveInsertTargetCols(targetOID, s.Cols)
				qTypes, qAny := c.extractRuleQueryColTypes(b, s.SelectStmt)
				if !qAny {
					for i, valType := range qTypes {
						if i < len(targetColDefs) && targetColDefs[i] != nil {
							if err := c.checkRuleTypeAssignable(targetColDefs[i], valType); err != nil {
								return err
							}
						}
					}
				}
			}
		}
		ctx := &ruleEvalContext{
			scope:       make(map[string]ruleScopeEntry),
			cteNames:    cteNames,
			cteCols:     cteCols,
			cteColTypes: cteColTypes,
			cteAny:      cteAny,
		}
		if targetOID != 0 && s.Relation != nil {
			if s.Relation.Alias != nil && s.Relation.Alias.Aliasname != "" {
				if err := ctx.setScope(s.Relation.Alias.Aliasname, ruleScopeEntry{relOID: targetOID, hasAlias: true, alias: s.Relation.Alias.Aliasname}, c.snapshot); err != nil {
					return err
				}
			} else {
				if err := ctx.setScope(s.Relation.Relname, ruleScopeEntry{relOID: targetOID}, c.snapshot); err != nil {
					return err
				}
			}
			if s.OnConflictClause != nil && s.OnConflictClause.TargetList != nil && len(s.OnConflictClause.TargetList.Items) > 0 {
				if err := ctx.setScope("excluded", ruleScopeEntry{relOID: targetOID, hasAlias: true, alias: "excluded"}, c.snapshot); err != nil {
					return err
				}
			}
		}
		if s.OnConflictClause != nil {
			oc := s.OnConflictClause
			isDoUpdate := oc.TargetList != nil && len(oc.TargetList.Items) > 0
			hasInferTarget := oc.Infer != nil && (oc.Infer.Conname != "" || (oc.Infer.IndexElems != nil && len(oc.Infer.IndexElems.Items) > 0))
			if isDoUpdate && !hasInferTarget && !c.snapshot {
				return &Error{
					Code:    CodeSyntaxError,
					Message: "ON CONFLICT DO UPDATE requires inference specification or constraint name",
				}
			}
			if oc.Infer != nil {
				if oc.Infer.Conname != "" && targetOID != 0 {
					var matchedCon *Constraint
					for _, con := range c.consByRel[targetOID] {
						if con != nil && con.Name == oc.Infer.Conname {
							matchedCon = con
							break
						}
					}
					relName := ""
					if rel := c.relationByOID[targetOID]; rel != nil {
						relName = rel.Name
					}
					if matchedCon == nil {
						if !c.snapshot {
							return &Error{
								Code:    CodeUndefinedObject,
								Message: fmt.Sprintf("constraint %q for table %q does not exist", oc.Infer.Conname, relName),
							}
						}
					} else {
						if matchedCon.Type != ConstraintPK && matchedCon.Type != ConstraintUnique && matchedCon.Type != ConstraintExclude && !c.snapshot {
							return &Error{
								Code:    CodeInvalidObjectDefinition,
								Message: fmt.Sprintf("constraint %q for table %q is not a unique, primary key, or exclusion constraint", oc.Infer.Conname, relName),
							}
						}
						b.conByInfer[oc.Infer] = matchedCon.OID
						b.addConDep(matchedCon.OID)
					}
				}
				if oc.Infer.IndexElems != nil {
					var rel *Relation
					if targetOID != 0 {
						rel = c.relationByOID[targetOID]
					}
					am := c.accessMethodByName["btree"]
					for _, ie := range oc.Infer.IndexElems.Items {
						elem, ok := ie.(*nodes.IndexElem)
						if !ok || elem == nil {
							continue
						}
						var colTypeOID uint32
						if elem.Name != "" && rel != nil {
							idx, exists := rel.colByName[elem.Name]
							if !exists {
								if !c.snapshot {
									return errUndefinedColumn(elem.Name)
								}
							} else {
								attNum := rel.Columns[idx].AttNum
								colTypeOID = rel.Columns[idx].TypeOID
								b.colByIE[elem] = ruleColBinding{relOID: targetOID, attNum: attNum}
								b.addColDep(targetOID, attNum)
							}
						}
						if elem.Expr != nil {
							if err := c.bindRuleExpr(b, elem.Expr, ctx); err != nil {
								return err
							}
						}
						if elem.Collation != nil && len(elem.Collation.Items) > 0 {
							_, uid, err := c.resolveCollationList(elem.Collation)
							if err != nil {
								if !c.snapshot {
									return err
								}
							} else if uid != 0 {
								b.collByIE[elem] = uid
								b.addCollDep(uid)
							}
						}
						if elem.Opclass != nil && len(elem.Opclass.Items) > 0 {
							var opClassParts []string
							for _, part := range elem.Opclass.Items {
								opClassParts = append(opClassParts, stringVal(part))
							}
							ref := opClassRef{
								schema:  strings.Join(opClassParts[:len(opClassParts)-1], "."),
								name:    opClassParts[len(opClassParts)-1],
								typeOID: colTypeOID,
							}
							if err := c.validateIndexOpClass(ref, am, elem.Opclassopts); err != nil {
								if !c.snapshot {
									return err
								}
							} else if am != nil {
								if opcOID := c.resolveOpClassOID(ref, am.OID); opcOID >= FirstNormalObjectId {
									b.opcByIE[elem] = opcOID
									b.addOpcDep(opcOID)
								}
							}
						}
					}
				}
				if err := c.bindRuleBoolExpr(b, oc.Infer.WhereClause, ctx, "WHERE"); err != nil {
					return err
				}
			}
			if err := c.bindTargetColumns(b, targetOID, oc.TargetList); err != nil {
				return err
			}
			if oc.TargetList != nil {
				for _, item := range oc.TargetList.Items {
					if rt, ok := item.(*nodes.ResTarget); ok {
						if err := c.bindRuleUpdateSource(b, rt.Val, ctx); err != nil {
							return err
						}
					}
				}
				if err := c.checkRuleTargetListTypes(b, targetOID, oc.TargetList, ctx); err != nil {
					return err
				}
			}
			if err := c.bindRuleBoolExpr(b, oc.WhereClause, ctx, "WHERE"); err != nil {
				return err
			}
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, targetOID, ctx); err != nil {
				return err
			}
		}
		return nil

	case *nodes.UpdateStmt:
		cteNames, cteCols, cteColTypes, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		ctx := &ruleEvalContext{
			scope:       make(map[string]ruleScopeEntry),
			cteNames:    cteNames,
			cteCols:     cteCols,
			cteColTypes: cteColTypes,
			cteAny:      cteAny,
		}
		targetOID, err := addScopeRangeVar(c, b, ctx, s.Relation, true)
		if err != nil {
			return err
		}
		if err := c.collectFromScope(b, s.FromClause, ctx); err != nil {
			return err
		}
		if err := c.bindTargetColumns(b, targetOID, s.TargetList); err != nil {
			return err
		}
		if s.TargetList != nil {
			for _, item := range s.TargetList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					if err := c.bindRuleUpdateSource(b, rt.Val, ctx); err != nil {
						return err
					}
				}
			}
			if err := c.checkRuleTargetListTypes(b, targetOID, s.TargetList, ctx); err != nil {
				return err
			}
		}
		if err := c.bindRuleBoolExpr(b, s.WhereClause, ctx, "WHERE"); err != nil {
			return err
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, targetOID, ctx); err != nil {
				return err
			}
		}
		return nil

	case *nodes.DeleteStmt:
		cteNames, cteCols, cteColTypes, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		ctx := &ruleEvalContext{
			scope:       make(map[string]ruleScopeEntry),
			cteNames:    cteNames,
			cteCols:     cteCols,
			cteColTypes: cteColTypes,
			cteAny:      cteAny,
		}
		targetOID, err := addScopeRangeVar(c, b, ctx, s.Relation, true)
		if err != nil {
			return err
		}
		if err := c.collectFromScope(b, s.UsingClause, ctx); err != nil {
			return err
		}
		if err := c.bindRuleBoolExpr(b, s.WhereClause, ctx, "WHERE"); err != nil {
			return err
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, targetOID, ctx); err != nil {
				return err
			}
		}
		return nil

	case *nodes.MergeStmt:
		cteNames, cteCols, cteColTypes, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		ctx := &ruleEvalContext{
			scope:       make(map[string]ruleScopeEntry),
			cteNames:    cteNames,
			cteCols:     cteCols,
			cteColTypes: cteColTypes,
			cteAny:      cteAny,
		}
		targetOID, err := addScopeRangeVar(c, b, ctx, s.Relation, true)
		if err != nil {
			return err
		}
		if err := c.collectTableRefScope(b, s.SourceRelation, ctx); err != nil {
			return err
		}
		if err := c.bindRuleBoolExpr(b, s.JoinCondition, ctx, "JOIN/ON"); err != nil {
			return err
		}
		if s.MergeWhenClauses != nil {
			for _, item := range s.MergeWhenClauses.Items {
				mwc, ok := item.(*nodes.MergeWhenClause)
				if !ok || mwc == nil {
					continue
				}
				if err := c.bindRuleBoolExpr(b, mwc.Condition, ctx, "WHEN"); err != nil {
					return err
				}
				if mwc.CommandType == nodes.CMD_INSERT {
					if err := c.bindTargetColumns(b, targetOID, mwc.TargetList); err != nil {
						return err
					}
					if mwc.Values != nil {
						for _, valExpr := range mwc.Values.Items {
							if err := c.bindRuleMergeValue(b, valExpr, ctx); err != nil {
								return err
							}
						}
						if targetOID != 0 && !c.snapshot {
							if mwc.TargetList != nil && len(mwc.TargetList.Items) > 0 {
								if len(mwc.Values.Items) > len(mwc.TargetList.Items) {
									return &Error{
										Code:    CodeSyntaxError,
										Message: "INSERT has more expressions than target columns",
									}
								}
								if len(mwc.Values.Items) < len(mwc.TargetList.Items) {
									return &Error{
										Code:    CodeSyntaxError,
										Message: "INSERT has more target columns than expressions",
									}
								}
							} else if rel := c.relationByOID[targetOID]; rel != nil && len(rel.Columns) > 0 && len(mwc.Values.Items) > len(rel.Columns) {
								return &Error{
									Code:    CodeSyntaxError,
									Message: "INSERT has more expressions than target columns",
								}
							}
							targetColDefs := c.resolveInsertTargetCols(targetOID, mwc.TargetList)
							for i, valExpr := range mwc.Values.Items {
								if i < len(targetColDefs) && targetColDefs[i] != nil {
									valType := c.inferRuleExprType(b, valExpr, ctx)
									if err := c.checkRuleTypeAssignable(targetColDefs[i], valType); err != nil {
										return err
									}
								}
							}
						}
					}
				} else if mwc.CommandType == nodes.CMD_UPDATE {
					if err := c.bindTargetColumns(b, targetOID, mwc.TargetList); err != nil {
						return err
					}
					if mwc.TargetList != nil {
						for _, tlItem := range mwc.TargetList.Items {
							if rt, ok := tlItem.(*nodes.ResTarget); ok {
								if err := c.bindRuleUpdateSource(b, rt.Val, ctx); err != nil {
									return err
								}
							}
						}
						if err := c.checkRuleTargetListTypes(b, targetOID, mwc.TargetList, ctx); err != nil {
							return err
						}
					}
				}
			}
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, targetOID, ctx); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (c *Catalog) bindAndValidateReturningList(b *ruleBindings, list *nodes.List, targetOID uint32, ctx *ruleEvalContext) error {
	if list == nil {
		return nil
	}
	if ctx != nil {
		// RETURNING refuses aggregate, window and set-returning calls.
		prevAgg, prevWin, prevSRF := ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF
		ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF = "RETURNING", "RETURNING", "RETURNING"
		defer func() { ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF = prevAgg, prevWin, prevSRF }()
	}
	expandedCols := 0
	var retTypes []uint32
	var retTypMods []int32
	unknownCount := false
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok {
			continue
		}
		if rt.Name == "" && b != nil && b.targetNameByRT != nil {
			if inferred := inferRuleTargetColName(rt.Val); inferred != "" && inferred != "?column?" {
				b.targetNameByRT[rt] = inferred
			}
		}
		if err := c.bindRuleExpr(b, rt.Val, ctx); err != nil {
			return err
		}
		if cr, isCR := rt.Val.(*nodes.ColumnRef); isCR && isStarColumnRef(cr) {
			starCols, anyOpen := c.expandRuleStarTarget(b, cr, nil, targetOID, ctx)
			if !anyOpen && len(starCols) > 0 && b != nil && b.starExpandByRT != nil {
				b.starExpandByRT[rt] = starCols
			}
			if anyOpen {
				unknownCount = true
			} else {
				expandedCols += len(starCols)
				for _, sc := range starCols {
					retTypes = append(retTypes, sc.typeOID)
					retTypMods = append(retTypMods, sc.typeMod)
				}
			}
			continue
		}
		expandedCols++
		retTypes = append(retTypes, c.inferRuleExprType(b, rt.Val, ctx))
		retTypMods = append(retTypMods, c.inferRuleExprTypMod(b, rt.Val, ctx))
	}
	if b.ownerRelOID != 0 && !unknownCount && !c.snapshot {
		if ownerRel := c.relationByOID[b.ownerRelOID]; ownerRel != nil {
			if expandedCols > len(ownerRel.Columns) {
				return &Error{
					Code:    CodeInvalidObjectDefinition,
					Message: "RETURNING list has too many entries",
				}
			}
			if expandedCols < len(ownerRel.Columns) {
				return &Error{
					Code:    CodeInvalidObjectDefinition,
					Message: "RETURNING list has too few entries",
				}
			}
			for i, col := range ownerRel.Columns {
				if i < len(retTypes) && retTypes[i] != 0 && retTypes[i] != UNKNOWNOID && retTypes[i] != col.TypeOID {
					return &Error{
						Code:    CodeDatatypeMismatch,
						Message: fmt.Sprintf("RETURNING list's entry %d has different type from column %q", i+1, col.Name),
					}
				}
				if i < len(retTypMods) && col.TypeMod >= 0 && retTypMods[i] >= 0 && retTypMods[i] != col.TypeMod {
					return &Error{
						Code:    CodeDatatypeMismatch,
						Message: fmt.Sprintf("RETURNING list's entry %d has different size from column %q", i+1, col.Name),
					}
				}
			}
		}
	}
	return nil
}

func (c *Catalog) inferRuleExprTypMod(b *ruleBindings, expr nodes.Node, ctx *ruleEvalContext) int32 {
	if expr == nil {
		return -1
	}
	switch v := expr.(type) {
	case *nodes.CollateClause:
		return c.inferRuleExprTypMod(b, v.Arg, ctx)
	case *nodes.TypeCast:
		if v.TypeName != nil {
			if _, typMod, err := c.resolveTypeName(v.TypeName); err == nil {
				return typMod
			}
		}
		return -1
	case *nodes.ColumnRef:
		if b != nil {
			if cb, ok := b.colByRef[v]; ok && cb.relOID != 0 && cb.attNum != 0 {
				if rel := c.relationByOID[cb.relOID]; rel != nil {
					for _, col := range rel.Columns {
						if col.AttNum == cb.attNum {
							return col.TypeMod
						}
					}
				}
			}
		}
		if v.Fields == nil || len(v.Fields.Items) == 0 || ctx == nil {
			return -1
		}
		items := v.Fields.Items
		if len(items) == 4 {
			items = items[1:]
		}
		switch len(items) {
		case 1:
			colName := stringVal(items[0])
			for cur := ctx; cur != nil; cur = cur.outer {
				for k, se := range cur.scope {
					if k == "excluded" || se.hiddenByJoinColAlias || se.isVirtual {
						continue
					}
					if rel := c.relationByOID[se.relOID]; rel != nil {
						if attNum, ok := resolveScopeColumn(rel, se, colName); ok {
							for _, col := range rel.Columns {
								if col.AttNum == attNum {
									return col.TypeMod
								}
							}
						}
					}
				}
			}
		case 2:
			qual := stringVal(items[0])
			colName := stringVal(items[1])
			if se, ok := lookupScopeEntry(ctx, qual); ok && !se.isVirtual {
				if rel := c.relationByOID[se.relOID]; rel != nil {
					if attNum, ok := resolveScopeColumn(rel, se, colName); ok {
						for _, col := range rel.Columns {
							if col.AttNum == attNum {
								return col.TypeMod
							}
						}
					}
				}
			}
		case 3:
			relName := stringVal(items[1])
			colName := stringVal(items[2])
			if se, ok := lookupScopeEntry(ctx, relName); ok && !se.isVirtual && !se.hasAlias {
				if rel := c.relationByOID[se.relOID]; rel != nil {
					if attNum, ok := resolveScopeColumn(rel, se, colName); ok {
						for _, col := range rel.Columns {
							if col.AttNum == attNum {
								return col.TypeMod
							}
						}
					}
				}
			}
		}
	}
	return -1
}

func isStarColumnRef(cr *nodes.ColumnRef) bool {
	if cr == nil || cr.Fields == nil || len(cr.Fields.Items) == 0 {
		return false
	}
	_, isStar := cr.Fields.Items[len(cr.Fields.Items)-1].(*nodes.A_Star)
	return isStar
}

func (c *Catalog) expandRuleStarTarget(b *ruleBindings, cr *nodes.ColumnRef, fromClause *nodes.List, targetOID uint32, ctx *ruleEvalContext) ([]ruleStarCol, bool) {
	if cr == nil || cr.Fields == nil || len(cr.Fields.Items) == 0 {
		return nil, true
	}
	fields := cr.Fields.Items
	if len(fields) == 1 {
		if fromClause != nil {
			var out []ruleStarCol
			for _, fromItem := range fromClause.Items {
				subCols, open := c.expandTableRefStarCols(b, fromItem, ctx)
				if open {
					return nil, true
				}
				out = append(out, subCols...)
			}
			return out, false
		}
		if ctx == nil || ctx.hasOpenSource {
			return nil, true
		}
		var out []ruleStarCol
		seenRel := map[uint32]bool{}
		for _, k := range ctx.scopeOrder {
			se, ok := ctx.scope[k]
			if !ok || k == "excluded" || se.hiddenByJoinColAlias {
				continue
			}
			if se.isVirtual {
				if se.isJoinAlias && !se.hasJoinColAliases {
					continue
				}
				if se.anyCol {
					return nil, true
				}
				for _, vCol := range se.virtualColOrder {
					var vType uint32
					if se.virtualColTypes != nil {
						vType = se.virtualColTypes[vCol]
					}
					out = append(out, ruleStarCol{
						prefix:   se.alias,
						colAlias: vCol,
						origName: vCol,
						typeOID:  vType,
						typeMod:  -1,
					})
				}
				continue
			}
			if seenRel[se.relOID] {
				continue
			}
			seenRel[se.relOID] = true
			rel := c.relationByOID[se.relOID]
			if rel == nil {
				return nil, true
			}
			for i, col := range rel.Columns {
				if i < len(se.colAliases) && se.colAliases[i] != "" {
					out = append(out, ruleStarCol{
						relOID:     se.relOID,
						prefix:     se.alias,
						useRelName: !se.hasAlias,
						colAlias:   se.colAliases[i],
						origName:   se.colAliases[i],
						typeOID:    col.TypeOID,
						typeMod:    col.TypeMod,
					})
				} else {
					out = append(out, ruleStarCol{
						relOID:     se.relOID,
						attNum:     col.AttNum,
						prefix:     se.alias,
						useRelName: !se.hasAlias,
						origName:   col.Name,
						typeOID:    col.TypeOID,
						typeMod:    col.TypeMod,
					})
				}
			}
		}
		return out, false
	}

	qual := stringVal(fields[len(fields)-2])
	if (strings.EqualFold(qual, "old") && b != nil && b.allowOld) || (strings.EqualFold(qual, "new") && b != nil && b.allowNew) {
		ownerRel := c.relationByOID[b.ownerRelOID]
		if ownerRel == nil {
			return nil, true
		}
		out := make([]ruleStarCol, 0, len(ownerRel.Columns))
		for _, col := range ownerRel.Columns {
			out = append(out, ruleStarCol{
				relOID:   b.ownerRelOID,
				attNum:   col.AttNum,
				prefix:   strings.ToLower(qual),
				origName: col.Name,
				typeOID:  col.TypeOID,
				typeMod:  col.TypeMod,
			})
		}
		return out, false
	}
	se, ok := lookupScopeEntry(ctx, qual)
	if !ok {
		return nil, true
	}
	return c.expandScopeEntryStarCols(se, qual)
}

func (c *Catalog) expandScopeEntryStarCols(se ruleScopeEntry, fallbackPrefix string) ([]ruleStarCol, bool) {
	if se.isVirtual {
		if se.anyCol {
			return nil, true
		}
		prefix := se.alias
		if prefix == "" {
			prefix = fallbackPrefix
		}
		out := make([]ruleStarCol, 0, len(se.virtualColOrder))
		for _, vCol := range se.virtualColOrder {
			var vType uint32
			if se.virtualColTypes != nil {
				vType = se.virtualColTypes[vCol]
			}
			out = append(out, ruleStarCol{
				prefix:   prefix,
				colAlias: vCol,
				origName: vCol,
				typeOID:  vType,
				typeMod:  -1,
			})
		}
		return out, false
	}
	rel := c.relationByOID[se.relOID]
	if rel == nil {
		return nil, true
	}
	out := make([]ruleStarCol, 0, len(rel.Columns))
	for i, col := range rel.Columns {
		if i < len(se.colAliases) && se.colAliases[i] != "" {
			out = append(out, ruleStarCol{
				relOID:     se.relOID,
				prefix:     se.alias,
				useRelName: !se.hasAlias,
				colAlias:   se.colAliases[i],
				origName:   se.colAliases[i],
				typeOID:    col.TypeOID,
				typeMod:    col.TypeMod,
			})
		} else {
			out = append(out, ruleStarCol{
				relOID:     se.relOID,
				attNum:     col.AttNum,
				prefix:     se.alias,
				useRelName: !se.hasAlias,
				origName:   col.Name,
				typeOID:    col.TypeOID,
				typeMod:    col.TypeMod,
			})
		}
	}
	return out, false
}

func (c *Catalog) expandTableRefStarCols(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext) ([]ruleStarCol, bool) {
	if n == nil || ctx == nil {
		return nil, true
	}
	switch v := n.(type) {
	case *nodes.RangeVar:
		if v == nil {
			return nil, true
		}
		name := v.Relname
		if v.Alias != nil && v.Alias.Aliasname != "" {
			name = v.Alias.Aliasname
		}
		se, ok := ctx.scope[name]
		if !ok {
			return nil, true
		}
		return c.expandScopeEntryStarCols(se, name)
	case *nodes.RangeSubselect:
		if v.Alias != nil && v.Alias.Aliasname != "" {
			if se, ok := ctx.scope[v.Alias.Aliasname]; ok {
				return c.expandScopeEntryStarCols(se, v.Alias.Aliasname)
			}
		}
	case *nodes.RangeFunction:
		scopeName := ""
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName = v.Alias.Aliasname
		} else if v.Functions != nil && len(v.Functions.Items) == 1 {
			scopeName = singleRangeFuncName(v.Functions.Items[0])
		}
		if scopeName != "" {
			if se, ok := ctx.scope[scopeName]; ok {
				return c.expandScopeEntryStarCols(se, scopeName)
			}
		}
	case *nodes.RangeTableFunc:
		scopeName := "xmltable"
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName = v.Alias.Aliasname
		}
		if se, ok := ctx.scope[scopeName]; ok {
			return c.expandScopeEntryStarCols(se, scopeName)
		}
	case *nodes.JsonTable:
		scopeName := "json_table"
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName = v.Alias.Aliasname
		}
		if se, ok := ctx.scope[scopeName]; ok {
			return c.expandScopeEntryStarCols(se, scopeName)
		}
	case *nodes.RangeTableSample:
		return c.expandTableRefStarCols(b, v.Relation, ctx)
	case *nodes.JoinExpr:
		if v.Alias != nil && v.Alias.Aliasname != "" && len(extractAliasColnames(v.Alias)) > 0 {
			if se, ok := ctx.scope[v.Alias.Aliasname]; ok {
				return c.expandScopeEntryStarCols(se, v.Alias.Aliasname)
			}
		}
		lStar, lOpen := c.expandTableRefStarCols(b, v.Larg, ctx)
		rStar, rOpen := c.expandTableRefStarCols(b, v.Rarg, ctx)
		if lOpen || rOpen {
			return nil, true
		}
		var commonNames []string
		commonSet := make(map[string]bool)
		if v.UsingClause != nil && len(v.UsingClause.Items) > 0 {
			for _, u := range v.UsingClause.Items {
				uCol := stringVal(u)
				if uCol != "" && !commonSet[uCol] {
					commonSet[uCol] = true
					commonNames = append(commonNames, uCol)
				}
			}
		} else if v.IsNatural {
			rNames := make(map[string]bool, len(rStar))
			for _, rc := range rStar {
				rNames[rc.origName] = true
			}
			for _, lc := range lStar {
				if lc.origName != "" && rNames[lc.origName] && !commonSet[lc.origName] {
					commonSet[lc.origName] = true
					commonNames = append(commonNames, lc.origName)
				}
			}
		}
		if len(commonNames) == 0 {
			return append(lStar, rStar...), false
		}
		var out []ruleStarCol
		for _, cName := range commonNames {
			for _, lc := range lStar {
				if lc.origName == cName {
					out = append(out, lc)
					break
				}
			}
		}
		for _, lc := range lStar {
			if !commonSet[lc.origName] {
				out = append(out, lc)
			}
		}
		for _, rc := range rStar {
			if !commonSet[rc.origName] {
				out = append(out, rc)
			}
		}
		return out, false
	}
	return nil, true
}

func (c *Catalog) bindRuleSortClause(b *ruleBindings, list *nodes.List, ctx *ruleEvalContext) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		sb, ok := item.(*nodes.SortBy)
		if !ok || sb == nil {
			continue
		}
		if err := c.bindRuleExpr(b, sb.Node, ctx); err != nil {
			return err
		}
		if sb.SortbyDir == nodes.SORTBY_USING && sb.UseOp != nil {
			sortType := c.inferRuleExprType(b, sb.Node, ctx)
			resType, err := c.bindRuleOperator(b, sb.UseOp, sortType, sortType, false, true)
			if err != nil {
				return err
			}
			if resType != 0 && resType != BOOLOID && !c.snapshot {
				return &Error{
					Code:    CodeWrongObjectType,
					Message: fmt.Sprintf("sort operator %s must return type boolean", deparseOperatorName(sb.UseOp)),
				}
			}
		}
	}
	return nil
}

func (c *Catalog) bindRuleOperator(b *ruleBindings, opNameList *nodes.List, lType, rType uint32, isPrefix, requireExists bool) (uint32, error) {
	if opNameList == nil || len(opNameList.Items) == 0 {
		return 0, nil
	}
	opSchema, opName := qualifiedName(opNameList)
	if opName == "" {
		return 0, nil
	}
	if opSchema != "" && c.schemaByName[opSchema] == nil && !c.snapshot {
		return 0, errUndefinedSchema(opSchema)
	}
	var searchNS []uint32
	if opSchema != "" {
		if s := c.schemaByName[opSchema]; s != nil {
			searchNS = []uint32{s.OID}
		}
	} else {
		searchNS = c.routineSearchPathWithCatalog()
	}

	opNamespace := func(op *BuiltinOperator) uint32 {
		if op.OID >= FirstNormalObjectId {
			return c.operatorSchemas[op.OID]
		}
		return PGCatalogNamespace
	}
	findExactInNS := func(left, right uint32) *BuiltinOperator {
		ops := c.operByKey[operKey{name: opName, left: left, right: right}]
		if len(ops) == 0 {
			return nil
		}
		for _, nsOID := range searchNS {
			for _, op := range ops {
				if (isPrefix && op.Kind != 'l') || (!isPrefix && op.Kind == 'l') {
					continue
				}
				if opNamespace(op) == nsOID {
					return op
				}
			}
		}
		return nil
	}
	recordOp := func(op *BuiltinOperator) uint32 {
		if op == nil {
			return 0
		}
		if op.OID >= FirstNormalObjectId && b != nil {
			b.addOperDep(op.OID)
			if op.Code >= FirstNormalObjectId {
				b.addFuncDep(op.Code)
			}
			if b.operByList != nil {
				b.operByList[opNameList] = op.OID
			}
		}
		return op.Result
	}

	if (isPrefix || lType != 0) && rType != 0 {
		if isPrefix {
			if op := findExactInNS(0, rType); op != nil {
				return recordOp(op), nil
			}
			if rType == UNKNOWNOID {
				for _, commonType := range []uint32{TEXTOID, INT4OID, FLOAT8OID} {
					if op := findExactInNS(0, commonType); op != nil {
						return recordOp(op), nil
					}
				}
			}
			type sigKey struct{ l, r uint32 }
			seenSig := make(map[sigKey]bool)
			var candidates []funcCandidate
			for _, nsOID := range searchNS {
				for _, op := range c.operByOID {
					if op.Name != opName || op.Kind != 'l' || opNamespace(op) != nsOID {
						continue
					}
					sk := sigKey{0, op.Right}
					if seenSig[sk] {
						continue
					}
					seenSig[sk] = true
					if c.canCoerceArgs([]uint32{rType}, []uint32{op.Right}) {
						candidates = append(candidates, funcCandidate{argTypes: []uint32{op.Right}, op: op})
					}
				}
			}
			if len(candidates) == 1 {
				return recordOp(candidates[0].op), nil
			}
			if len(candidates) > 1 {
				if best := c.funcSelectCandidate([]uint32{rType}, candidates); best != nil && best.op != nil {
					return recordOp(best.op), nil
				}
			}
			if !c.snapshot {
				return 0, errUndefinedFunction(opName, []uint32{rType})
			}
			return 0, nil
		}

		// Binary operator with both operand types known.
		if op := findExactInNS(lType, rType); op != nil {
			return recordOp(op), nil
		}
		if lType == UNKNOWNOID && rType != UNKNOWNOID {
			if op := findExactInNS(rType, rType); op != nil {
				return recordOp(op), nil
			}
			if pref := c.preferredTypeForCategory(rType); pref != 0 && pref != rType {
				if op := findExactInNS(pref, pref); op != nil {
					return recordOp(op), nil
				}
			}
		}
		if rType == UNKNOWNOID && lType != UNKNOWNOID {
			if op := findExactInNS(lType, lType); op != nil {
				return recordOp(op), nil
			}
			if pref := c.preferredTypeForCategory(lType); pref != 0 && pref != lType {
				if op := findExactInNS(pref, pref); op != nil {
					return recordOp(op), nil
				}
			}
		}
		if lType == UNKNOWNOID && rType == UNKNOWNOID {
			if op := findExactInNS(TEXTOID, TEXTOID); op != nil {
				return recordOp(op), nil
			}
		}

		inputTypes := []uint32{lType, rType}
		type sigKey struct{ l, r uint32 }
		seenSig := make(map[sigKey]bool)
		var coercible []funcCandidate
		for _, nsOID := range searchNS {
			for _, op := range c.operByOID {
				if op.Name != opName || op.Kind == 'l' || opNamespace(op) != nsOID {
					continue
				}
				sk := sigKey{op.Left, op.Right}
				if seenSig[sk] {
					continue
				}
				seenSig[sk] = true
				candTypes := []uint32{op.Left, op.Right}
				if c.canCoerceArgs(inputTypes, candTypes) {
					coercible = append(coercible, funcCandidate{argTypes: candTypes, op: op})
				}
			}
		}
		if len(coercible) == 1 {
			return recordOp(coercible[0].op), nil
		}
		if len(coercible) > 1 {
			if best := c.funcSelectCandidate(inputTypes, coercible); best != nil && best.op != nil {
				return recordOp(best.op), nil
			}
		}
		if !c.snapshot {
			return 0, errUndefinedFunction(opName, []uint32{lType, rType})
		}
		return 0, nil
	}

	// Fallback when an operand type is unknown (e.g., from an open virtual source).
	var matched []*BuiltinOperator
	for _, nsOID := range searchNS {
		for _, op := range c.operByOID {
			if op.Name != opName || opNamespace(op) != nsOID {
				continue
			}
			if (isPrefix && op.Kind != 'l') || (!isPrefix && op.Kind == 'l') {
				continue
			}
			if !isPrefix && lType != 0 && lType != UNKNOWNOID && !c.canCoerceArgs([]uint32{lType}, []uint32{op.Left}) {
				continue
			}
			if rType != 0 && rType != UNKNOWNOID && !c.canCoerceArgs([]uint32{rType}, []uint32{op.Right}) {
				continue
			}
			matched = append(matched, op)
		}
	}
	if len(matched) == 1 {
		return recordOp(matched[0]), nil
	}
	if requireExists && len(matched) == 0 && !c.snapshot {
		shown := opName
		if opSchema != "" {
			shown = opSchema + "." + opName
		}
		return 0, &Error{
			Code:    CodeUndefinedFunction,
			Message: fmt.Sprintf("operator does not exist: %s", shown),
		}
	}
	return 0, nil
}

// bindRuleUpdateSource binds an UPDATE assignment value, where aggregate,
// window and set-returning calls are refused (EXPR_KIND_UPDATE_SOURCE).
// bindRuleMergeValue binds a MERGE ... INSERT VALUES item, where aggregate
// and window calls are refused as in VALUES.
func (c *Catalog) bindRuleMergeValue(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext) error {
	if ctx != nil {
		prevAgg, prevWin := ctx.disallowAgg, ctx.disallowWin
		ctx.disallowAgg, ctx.disallowWin = "VALUES", "VALUES"
		defer func() { ctx.disallowAgg, ctx.disallowWin = prevAgg, prevWin }()
	}
	return c.bindRuleExpr(b, n, ctx)
}

func (c *Catalog) bindRuleUpdateSource(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext) error {
	if ctx != nil {
		prevAgg, prevWin, prevSRF := ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF
		ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF = "UPDATE", "UPDATE", "UPDATE"
		defer func() { ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF = prevAgg, prevWin, prevSRF }()
	}
	return c.bindRuleExpr(b, n, ctx)
}

// bindRuleLimitExpr binds a LIMIT or OFFSET argument, where set-returning
// calls are refused.
func (c *Catalog) bindRuleLimitExpr(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext, construct string) error {
	if n == nil {
		return nil
	}
	if ctx != nil {
		// LIMIT and OFFSET refuse aggregate, window and set-returning calls.
		prevAgg, prevWin, prevSRF := ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF
		ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF = construct, construct, construct
		defer func() { ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF = prevAgg, prevWin, prevSRF }()
	}
	return c.bindRuleExpr(b, n, ctx)
}

func (c *Catalog) bindSortAndLimit(b *ruleBindings, s *nodes.SelectStmt, ctx *ruleEvalContext) error {
	if err := c.bindRuleSortClause(b, s.SortClause, ctx); err != nil {
		return err
	}
	if s.SortClause != nil && !c.snapshot {
		qCols, qAny := c.extractRuleQueryCols(b, s)
		for _, item := range s.SortClause.Items {
			if sb, ok := item.(*nodes.SortBy); ok && sb != nil {
				if err := c.checkRuleOrdinalExpr(sb.Node, len(qCols), qAny, "ORDER BY"); err != nil {
					return err
				}
			}
		}
		if err := c.checkRuleDistinctOnMatchesOrderBy(b, s); err != nil {
			return err
		}
	}
	if !c.snapshot {
		if err := c.checkRuleGroupingFuncs(b, s); err != nil {
			return err
		}
		if err := c.checkRuleUngroupedColumns(b, s, ctx); err != nil {
			return err
		}
	}
	if s.LimitOption == nodes.LIMIT_OPTION_WITH_TIES && (s.SortClause == nil || len(s.SortClause.Items) == 0) && !c.snapshot {
		return &Error{Code: CodeSyntaxError, Message: "WITH TIES cannot be specified without ORDER BY clause"}
	}
	if err := c.bindRuleLimitExpr(b, s.LimitCount, ctx, "LIMIT"); err != nil {
		return err
	}
	if s.LimitCount != nil && !c.snapshot && c.ruleExprRefersToLevel(b, s.LimitCount, ctx) {
		return &Error{Code: CodeInvalidColumnReference, Message: "argument of LIMIT must not contain variables"}
	}
	if s.LimitCount != nil && !c.snapshot {
		if t := c.inferRuleExprType(b, s.LimitCount, ctx); t != 0 && t != UNKNOWNOID && !c.CanCoerce(t, INT8OID, 'a') {
			return &Error{
				Code:    CodeDatatypeMismatch,
				Message: fmt.Sprintf("argument of LIMIT must be type bigint, not type %s", c.typeName(t)),
			}
		}
	}
	if err := c.bindRuleLimitExpr(b, s.LimitOffset, ctx, "OFFSET"); err != nil {
		return err
	}
	if s.LimitOffset != nil && !c.snapshot && c.ruleExprRefersToLevel(b, s.LimitOffset, ctx) {
		return &Error{Code: CodeInvalidColumnReference, Message: "argument of OFFSET must not contain variables"}
	}
	if s.LimitOffset != nil && !c.snapshot {
		if t := c.inferRuleExprType(b, s.LimitOffset, ctx); t != 0 && t != UNKNOWNOID && !c.CanCoerce(t, INT8OID, 'a') {
			return &Error{
				Code:    CodeDatatypeMismatch,
				Message: fmt.Sprintf("argument of OFFSET must be type bigint, not type %s", c.typeName(t)),
			}
		}
	}
	if s.LockingClause != nil && len(s.LockingClause.Items) > 0 {
		if !c.snapshot {
			if err := c.checkRuleLockingClauseAllowed(b, s); err != nil {
				return err
			}
		}
		for _, item := range s.LockingClause.Items {
			lc, ok := item.(*nodes.LockingClause)
			if !ok || lc == nil || lc.LockedRels == nil {
				continue
			}
			for _, rItem := range lc.LockedRels.Items {
				if rv, ok := rItem.(*nodes.RangeVar); ok && rv != nil && rv.Relname != "" {
					if rv.Catalogname != "" {
						if err := c.checkDatabaseQualifier(rv.Catalogname); err != nil {
							return err
						}
					}
					if _, found := lookupScopeEntry(ctx, rv.Relname); !found && !ctx.hasOpenSource && !c.snapshot {
						return &Error{
							Code:    CodeUndefinedTable,
							Message: fmt.Sprintf("relation %q in FOR UPDATE/SHARE clause not found in FROM clause", rv.Relname),
						}
					}
				}
			}
		}
	}
	return nil
}

// ruleExprRefersToLevel reports whether n refers to a column of level's own
// FROM items, which checkExprIsVarFree (parse_clause.c) refuses in LIMIT and
// OFFSET even from inside a subquery. A subquery's own FROM items shadow
// level. NEW and OLD are not counted, so a reference to them is let through.
func (c *Catalog) ruleExprRefersToLevel(b *ruleBindings, n nodes.Node, level *ruleEvalContext) bool {
	if b == nil || level == nil || len(level.scope) == 0 {
		return false
	}
	found := false
	nodes.Walk(ruleLevelVarFinder{c: c, b: b, level: level, found: &found}, n)
	return found
}

type ruleLevelVarFinder struct {
	c     *Catalog
	b     *ruleBindings
	level *ruleEvalContext
	inner []*ruleEvalContext
	found *bool
}

func (f ruleLevelVarFinder) Visit(n nodes.Node) nodes.Visitor {
	if n == nil || *f.found {
		return nil
	}
	switch v := n.(type) {
	case *nodes.SelectStmt:
		if sctx := f.b.stmtCtx[v]; sctx != nil {
			f.inner = append(f.inner[:len(f.inner):len(f.inner)], sctx)
		}
	case *nodes.ColumnRef:
		if f.refersToLevel(v) {
			*f.found = true
			return nil
		}
	}
	return f
}

func (f ruleLevelVarFinder) refersToLevel(cr *nodes.ColumnRef) bool {
	if cr.Fields == nil || len(cr.Fields.Items) == 0 {
		return false
	}
	items := cr.Fields.Items
	if len(items) == 1 {
		name := stringVal(items[0])
		if name == "" {
			return false
		}
		for i := len(f.inner) - 1; i >= 0; i-- {
			if localScopeHasCol(f.c, f.inner[i], name) {
				return false
			}
		}
		return localScopeHasCol(f.c, f.level, name)
	}
	qual := stringVal(items[len(items)-2])
	if qual == "" || strings.EqualFold(qual, "old") || strings.EqualFold(qual, "new") {
		return false
	}
	for i := len(f.inner) - 1; i >= 0; i-- {
		if _, ok := f.inner[i].scope[qual]; ok {
			return false
		}
	}
	_, ok := f.level.scope[qual]
	return ok
}

func (c *Catalog) checkRuleDistinctOnMatchesOrderBy(b *ruleBindings, s *nodes.SelectStmt) error {
	if s == nil || s.DistinctClause == nil || s.SortClause == nil || len(s.SortClause.Items) == 0 {
		return nil
	}
	hasDistinctOn := false
	for _, item := range s.DistinctClause.Items {
		if item != nil {
			hasDistinctOn = true
			break
		}
	}

	type targetKeyEntry struct {
		alias string
		key   string
	}
	dep := ruleDeparser{cat: c, b: b}
	exprCanonicalKey := func(expr nodes.Node) string {
		if expr == nil {
			return ""
		}
		if cr, ok := expr.(*nodes.ColumnRef); ok && b != nil {
			if cb, found := b.colByRef[cr]; found && cb.relOID != 0 && cb.attNum != 0 {
				return fmt.Sprintf("col:%d:%d", cb.relOID, cb.attNum)
			}
		}
		return "expr:" + dep.deparseExpr(expr)
	}

	var targets []targetKeyEntry
	if s.TargetList != nil {
		for _, item := range s.TargetList.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok || rt == nil {
				continue
			}
			if cr, isCR := rt.Val.(*nodes.ColumnRef); isCR && isStarColumnRef(cr) {
				if b != nil && b.starExpandByRT != nil {
					if starCols, found := b.starExpandByRT[rt]; found && len(starCols) > 0 {
						for _, sc := range starCols {
							colName := sc.colAlias
							if colName == "" {
								colName = sc.origName
							}
							key := "expr:" + quoteIdent(colName)
							if sc.relOID != 0 && sc.attNum != 0 {
								key = fmt.Sprintf("col:%d:%d", sc.relOID, sc.attNum)
							}
							targets = append(targets, targetKeyEntry{alias: colName, key: key})
						}
						continue
					}
				}
				return nil
			}
			alias := rt.Name
			if alias == "" && b != nil && b.targetNameByRT != nil {
				alias = b.targetNameByRT[rt]
			}
			targets = append(targets, targetKeyEntry{
				alias: alias,
				key:   exprCanonicalKey(rt.Val),
			})
		}
	}

	// A bare name in ORDER BY or DISTINCT ON names an output column before an
	// input column (parse_clause.c findTargetlistEntrySQL92).
	resolveSortGroupKey := func(expr nodes.Node) string {
		var iv *nodes.Integer
		switch v := expr.(type) {
		case *nodes.A_Const:
			if !v.Isnull {
				iv, _ = v.Val.(*nodes.Integer)
			}
		case *nodes.Integer:
			iv = v
		}
		if iv != nil && iv.Ival >= 1 && int(iv.Ival) <= len(targets) {
			return targets[iv.Ival-1].key
		}
		if cr, ok := expr.(*nodes.ColumnRef); ok && cr != nil && cr.Fields != nil && len(cr.Fields.Items) == 1 {
			if str, isStr := cr.Fields.Items[0].(*nodes.String); isStr && str.Str != "" {
				for _, t := range targets {
					if t.alias == str.Str {
						return t.key
					}
				}
			}
		}
		return exprCanonicalKey(expr)
	}

	if !hasDistinctOn {
		// Plain DISTINCT: every ORDER BY expression must be an output column
		// (parse_clause.c transformDistinctClause). Expressions compare by
		// structure, with each column reduced to the column it binds to.
		targetKeys := make(map[string]bool, len(targets))
		for _, t := range targets {
			targetKeys[t.key] = true
		}
		targetPrints := make(map[string]bool)
		for _, item := range s.TargetList.Items {
			if rt, ok := item.(*nodes.ResTarget); ok && rt != nil {
				targetPrints[ruleExprFingerprint(b, rt.Val)] = true
			}
		}
		for _, item := range s.SortClause.Items {
			sb, ok := item.(*nodes.SortBy)
			if !ok || sb == nil || sb.Node == nil || targetKeys[resolveSortGroupKey(sb.Node)] || targetPrints[ruleExprFingerprint(b, sb.Node)] {
				continue
			}
			return &Error{
				Code:    CodeInvalidColumnReference,
				Message: "for SELECT DISTINCT, ORDER BY expressions must appear in select list",
			}
		}
		return nil
	}

	distinctKeys := make(map[string]bool)
	for _, item := range s.DistinctClause.Items {
		if item == nil {
			continue
		}
		k := resolveSortGroupKey(item)
		if k != "" {
			distinctKeys[k] = true
		}
	}
	if len(distinctKeys) == 0 {
		return nil
	}

	matched := make(map[string]bool)
	for _, item := range s.SortClause.Items {
		if len(matched) == len(distinctKeys) {
			break
		}
		sb, ok := item.(*nodes.SortBy)
		if !ok || sb == nil || sb.Node == nil {
			continue
		}
		k := resolveSortGroupKey(sb.Node)
		if !distinctKeys[k] {
			return &Error{
				Code:    CodeInvalidColumnReference,
				Message: "SELECT DISTINCT ON expressions must match initial ORDER BY expressions",
			}
		}
		matched[k] = true
	}
	return nil
}

// ruleExprFingerprint renders n's structure without source positions. A
// column bound to a table column renders as that column, and any other
// column as its last name, so a qualified and an unqualified reference to the
// same column compare equal.
func ruleExprFingerprint(b *ruleBindings, n nodes.Node) string {
	var sb strings.Builder
	locType := reflect.TypeOf(nodes.Loc{})
	var walk func(v reflect.Value)
	walk = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.Interface, reflect.Pointer:
			if v.IsNil() {
				sb.WriteString("nil")
				return
			}
			if cr, ok := v.Interface().(*nodes.ColumnRef); ok && cr.Fields != nil && len(cr.Fields.Items) > 0 {
				if b != nil {
					if cb, found := b.colByRef[cr]; found && cb.relOID != 0 && cb.attNum != 0 {
						fmt.Fprintf(&sb, "col:%d:%d", cb.relOID, cb.attNum)
						return
					}
				}
				if str, isStr := cr.Fields.Items[len(cr.Fields.Items)-1].(*nodes.String); isStr {
					sb.WriteString("name:" + str.Str)
					return
				}
			}
			walk(v.Elem())
		case reflect.Struct:
			if v.Type() == locType {
				return
			}
			sb.WriteString(v.Type().Name())
			sb.WriteByte('{')
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).IsExported() {
					walk(v.Field(i))
					sb.WriteByte(',')
				}
			}
			sb.WriteByte('}')
		case reflect.Slice:
			sb.WriteByte('[')
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i))
				sb.WriteByte(',')
			}
			sb.WriteByte(']')
		default:
			fmt.Fprint(&sb, v.Interface())
		}
	}
	walk(reflect.ValueOf(n))
	return sb.String()
}

func (c *Catalog) checkRuleLockingClauseAllowed(b *ruleBindings, s *nodes.SelectStmt) error {
	if s == nil {
		return nil
	}
	if s.Op != nodes.SETOP_NONE {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "FOR UPDATE/SHARE is not allowed with UNION/INTERSECT/EXCEPT",
		}
	}
	if s.ValuesLists != nil {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "FOR UPDATE/SHARE is not allowed with VALUES",
		}
	}
	if s.DistinctClause != nil && len(s.DistinctClause.Items) > 0 {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "FOR UPDATE/SHARE is not allowed with DISTINCT clause",
		}
	}
	if s.GroupClause != nil && len(s.GroupClause.Items) > 0 {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "FOR UPDATE/SHARE is not allowed with GROUP BY clause",
		}
	}
	if s.HavingClause != nil {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "FOR UPDATE/SHARE is not allowed with HAVING clause",
		}
	}
	hasAgg, hasWin := c.selectHasTopLevelAggOrWindow(b, s)
	if hasAgg {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "FOR UPDATE/SHARE is not allowed with aggregate functions",
		}
	}
	if hasWin {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "FOR UPDATE/SHARE is not allowed with window functions",
		}
	}
	return nil
}

func (c *Catalog) selectHasTopLevelAggOrWindow(b *ruleBindings, s *nodes.SelectStmt) (hasAgg bool, hasWin bool) {
	if s == nil {
		return false, false
	}
	if s.TargetList != nil {
		hasAgg, hasWin = c.ruleExprHasAggOrWindow(b, s.TargetList)
	}
	if s.SortClause != nil {
		agg, win := c.ruleExprHasAggOrWindow(b, s.SortClause)
		hasAgg, hasWin = hasAgg || agg, hasWin || win
	}
	return hasAgg, hasWin
}

// checkRuleNestedAggOrWindow rejects what parse_agg.c rejects inside a call's
// own query level: an aggregate or window call in an aggregate's arguments or
// ordering, and a window call in a window function's arguments. A subquery
// starts a new level, so its calls are not nested.
func (c *Catalog) checkRuleNestedAggOrWindow(b *ruleBindings, v *nodes.FuncCall) error {
	if v.Over != nil {
		if v.Args != nil {
			if _, win := c.ruleExprHasAggOrWindow(b, v.Args); win {
				return &Error{Code: CodeWindowingError, Message: "window function calls cannot be nested"}
			}
			if ruleExprHasSRF(b, v.Args) {
				return &Error{Code: CodeFeatureNotSupported, Message: "window function calls cannot contain set-returning function calls"}
			}
		}
		return nil
	}
	if !c.isRuleAggFuncCall(b, v) {
		return nil
	}
	for _, list := range []*nodes.List{v.Args, v.AggOrder} {
		if list != nil && ruleExprHasSRF(b, list) {
			return &Error{Code: CodeFeatureNotSupported, Message: "aggregate function calls cannot contain set-returning function calls"}
		}
	}
	if v.AggDistinct && v.AggOrder != nil && v.Args != nil {
		// pg: parse_agg.c transformAggregateCall — with DISTINCT, each
		// ORDER BY expression must be one of the arguments.
		argPrints := make(map[string]bool, len(v.Args.Items))
		for _, arg := range v.Args.Items {
			argPrints[ruleExprFingerprint(b, arg)] = true
		}
		for _, item := range v.AggOrder.Items {
			if sb, ok := item.(*nodes.SortBy); ok && sb != nil && !argPrints[ruleExprFingerprint(b, sb.Node)] {
				return &Error{Code: CodeInvalidColumnReference, Message: "in an aggregate with DISTINCT, ORDER BY expressions must appear in argument list"}
			}
		}
	}
	for _, list := range []*nodes.List{v.Args, v.AggOrder} {
		if list == nil {
			continue
		}
		agg, win := c.ruleExprHasAggOrWindow(b, list)
		if agg {
			return &Error{Code: CodeGroupingError, Message: "aggregate function calls cannot be nested"}
		}
		if win {
			return &Error{Code: CodeGroupingError, Message: "aggregate function calls cannot contain window function calls"}
		}
	}
	return nil
}

// ruleExprHasSRF reports whether n holds a set-returning call of its own
// query level; it does not descend into subqueries.
func ruleExprHasSRF(b *ruleBindings, n nodes.Node) bool {
	if b == nil || len(b.srfCalls) == 0 {
		return false
	}
	found := false
	nodes.Inspect(n, func(x nodes.Node) bool {
		if found {
			return false
		}
		switch v := x.(type) {
		case *nodes.SubLink, *nodes.SelectStmt:
			return false
		case *nodes.FuncCall:
			if b.srfCalls[v] {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// ruleExprHasAggOrWindow reports whether n holds an aggregate or a window
// call of its own query level; it does not descend into subqueries.
func (c *Catalog) ruleExprHasAggOrWindow(b *ruleBindings, root nodes.Node) (hasAgg bool, hasWin bool) {
	var walk func(n nodes.Node)
	walk = func(n nodes.Node) {
		if n == nil || (hasAgg && hasWin) {
			return
		}
		switch v := n.(type) {
		case *nodes.SubLink, *nodes.RangeSubselect, *nodes.SelectStmt, *nodes.JsonArrayQueryConstructor:
			return
		case *nodes.List:
			if v == nil {
				return
			}
			for _, item := range v.Items {
				walk(item)
			}
		case *nodes.ResTarget:
			walk(v.Val)
		case *nodes.SortBy:
			walk(v.Node)
		case *nodes.FuncCall:
			if v.Over != nil {
				hasWin = true
			} else if v.AggStar || v.AggDistinct || v.AggFilter != nil || (v.AggOrder != nil && len(v.AggOrder.Items) > 0) {
				hasAgg = true
			} else if b != nil && b.funcByFC != nil {
				if oid, ok := b.funcByFC[v]; ok && oid != 0 {
					if p := c.procByOID[oid]; p != nil {
						if p.Kind == 'a' {
							hasAgg = true
						} else if p.Kind == 'w' {
							hasWin = true
						}
					}
				}
			}
			if !hasAgg && !hasWin && v.Over == nil {
				_, fName := qualifiedName(v.Funcname)
				for _, p := range c.procByName[strings.ToLower(fName)] {
					if p != nil && p.Kind == 'a' {
						hasAgg = true
						break
					}
				}
			}
			if v.Args != nil {
				walk(v.Args)
			}
			walk(v.AggFilter)
			if v.AggOrder != nil {
				walk(v.AggOrder)
			}
		case *nodes.A_Expr:
			walk(v.Lexpr)
			walk(v.Rexpr)
		case *nodes.BoolExpr:
			if v.Args != nil {
				walk(v.Args)
			}
		case *nodes.NullTest:
			walk(v.Arg)
		case *nodes.BooleanTest:
			walk(v.Arg)
		case *nodes.TypeCast:
			walk(v.Arg)
		case *nodes.CollateClause:
			walk(v.Arg)
		case *nodes.CaseExpr:
			walk(v.Arg)
			if v.Args != nil {
				walk(v.Args)
			}
			walk(v.Defresult)
		case *nodes.CaseWhen:
			walk(v.Expr)
			walk(v.Result)
		case *nodes.CoalesceExpr:
			if v.Args != nil {
				walk(v.Args)
			}
		case *nodes.MinMaxExpr:
			if v.Args != nil {
				walk(v.Args)
			}
		case *nodes.RowExpr:
			if v.Args != nil {
				walk(v.Args)
			}
		case *nodes.A_ArrayExpr:
			if v.Elements != nil {
				walk(v.Elements)
			}
		case *nodes.ArrayExpr:
			if v.Elements != nil {
				walk(v.Elements)
			}
		case *nodes.A_Indirection:
			walk(v.Arg)
			if v.Indirection != nil {
				walk(v.Indirection)
			}
		case *nodes.A_Indices:
			walk(v.Lidx)
			walk(v.Uidx)
		case *nodes.NamedArgExpr:
			walk(v.Arg)
		case *nodes.GroupingFunc:
			hasAgg = true
		}
	}
	walk(root)
	return hasAgg, hasWin
}

func (c *Catalog) isRuleAggFuncCall(b *ruleBindings, v *nodes.FuncCall) bool {
	if v == nil || v.Over != nil {
		return false
	}
	if v.AggStar || v.AggDistinct || v.AggFilter != nil || (v.AggOrder != nil && len(v.AggOrder.Items) > 0) {
		return true
	}
	if b != nil && b.funcByFC != nil {
		if oid, ok := b.funcByFC[v]; ok && oid != 0 {
			if p := c.procByOID[oid]; p != nil {
				return p.Kind == 'a'
			}
		}
	}
	_, fName := qualifiedName(v.Funcname)
	for _, p := range c.procByName[strings.ToLower(fName)] {
		if p != nil && p.Kind == 'a' {
			return true
		}
	}
	return false
}

func localScopeHasCol(c *Catalog, ctx *ruleEvalContext, colName string) bool {
	if ctx == nil || colName == "" {
		return false
	}
	for k, se := range ctx.scope {
		if k == "excluded" || se.hiddenByJoinColAlias {
			continue
		}
		if se.isVirtual {
			if se.isJoinAlias && !se.hasJoinColAliases {
				continue
			}
			if se.virtualCols[colName] > 0 || se.anyCol {
				return true
			}
			continue
		}
		if rel := c.relationByOID[se.relOID]; rel != nil {
			if _, ok := resolveScopeColumn(rel, se, colName); ok {
				return true
			}
		}
	}
	if se, ok := ctx.scope[colName]; ok && colName != "excluded" {
		if !se.isVirtual && se.relOID != 0 {
			return true
		}
		if len(se.joinRels) > 0 {
			return true
		}
	}
	return false
}

func (c *Catalog) checkRuleUngroupedColumns(b *ruleBindings, s *nodes.SelectStmt, ctx *ruleEvalContext) error {
	if s == nil || s.Op != nodes.SETOP_NONE || s.ValuesLists != nil || ctx == nil || ctx.hasOpenSource {
		return nil
	}
	for k, se := range ctx.scope {
		if k != "excluded" && se.isVirtual && se.anyCol {
			return nil
		}
	}

	hasGroupClause := s.GroupClause != nil && len(s.GroupClause.Items) > 0
	hasHaving := s.HavingClause != nil
	hasAgg, _ := c.selectHasTopLevelAggOrWindow(b, s)
	if !hasAgg && s.HavingClause != nil {
		tmpStmt := &nodes.SelectStmt{TargetList: &nodes.List{Items: []nodes.Node{s.HavingClause}}}
		hasAgg, _ = c.selectHasTopLevelAggOrWindow(b, tmpStmt)
	}
	if !hasGroupClause && !hasHaving && !hasAgg {
		return nil
	}

	relOIDCount := make(map[uint32]int)
	for k, se := range ctx.scope {
		if k == "excluded" || se.isVirtual || se.relOID == 0 {
			continue
		}
		relOIDCount[se.relOID]++
	}

	colBindingKey := func(cb ruleColBinding) string {
		if cb.relOID == 0 {
			return ""
		}
		if relOIDCount[cb.relOID] > 1 && cb.prefix != "" {
			return fmt.Sprintf("col:%s:%d:%d", cb.prefix, cb.relOID, cb.attNum)
		}
		return fmt.Sprintf("col:%d:%d", cb.relOID, cb.attNum)
	}

	type targetKeyEntry struct {
		alias string
		key   string
		expr  nodes.Node
	}
	dep := ruleDeparser{cat: c, b: b}
	exprCanonicalKey := func(expr nodes.Node) string {
		if expr == nil {
			return ""
		}
		if cr, ok := expr.(*nodes.ColumnRef); ok && b != nil {
			if cb, found := b.colByRef[cr]; found && cb.relOID != 0 && cb.attNum != 0 {
				return colBindingKey(cb)
			}
		}
		return "expr:" + dep.deparseExpr(expr)
	}

	var targets []targetKeyEntry
	if s.TargetList != nil {
		for _, item := range s.TargetList.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok || rt == nil {
				continue
			}
			if cr, isCR := rt.Val.(*nodes.ColumnRef); isCR && isStarColumnRef(cr) {
				if b != nil && b.starExpandByRT != nil {
					if starCols, found := b.starExpandByRT[rt]; found && len(starCols) > 0 {
						for _, sc := range starCols {
							colName := sc.colAlias
							if colName == "" {
								colName = sc.origName
							}
							key := "expr:" + quoteIdent(colName)
							if sc.relOID != 0 && sc.attNum != 0 {
								key = colBindingKey(ruleColBinding{relOID: sc.relOID, attNum: sc.attNum, prefix: sc.prefix})
							}
							targets = append(targets, targetKeyEntry{alias: colName, key: key})
						}
						continue
					}
				}
				return nil
			}
			alias := rt.Name
			if alias == "" && b != nil && b.targetNameByRT != nil {
				alias = b.targetNameByRT[rt]
			}
			targets = append(targets, targetKeyEntry{
				alias: alias,
				key:   exprCanonicalKey(rt.Val),
				expr:  rt.Val,
			})
		}
	}

	resolveSortGroupKey := func(expr nodes.Node) (string, nodes.Node) {
		var iv *nodes.Integer
		switch v := expr.(type) {
		case *nodes.A_Const:
			if !v.Isnull {
				iv, _ = v.Val.(*nodes.Integer)
			}
		case *nodes.Integer:
			iv = v
		}
		if iv != nil && iv.Ival >= 1 && int(iv.Ival) <= len(targets) {
			t := targets[iv.Ival-1]
			return t.key, t.expr
		}
		if cr, ok := expr.(*nodes.ColumnRef); ok && cr != nil && cr.Fields != nil && len(cr.Fields.Items) == 1 {
			if b == nil || b.colByRef[cr].relOID == 0 {
				if str, isStr := cr.Fields.Items[0].(*nodes.String); isStr && str.Str != "" {
					for _, t := range targets {
						if t.alias == str.Str {
							return t.key, t.expr
						}
					}
				}
			}
		}
		return exprCanonicalKey(expr), expr
	}

	groupKeys := make(map[string]bool)
	pkCoveredRels := make(map[uint32]bool)
	var addGroupItem func(n nodes.Node)
	addGroupItem = func(n nodes.Node) {
		if n == nil {
			return
		}
		switch v := n.(type) {
		case *nodes.GroupingSet:
			if v.Content != nil {
				for _, item := range v.Content.Items {
					addGroupItem(item)
				}
			}
		case *nodes.RowExpr:
			if v.Args != nil {
				for _, item := range v.Args.Items {
					addGroupItem(item)
				}
			}
		default:
			k, resolvedExpr := resolveSortGroupKey(n)
			if k != "" {
				groupKeys[k] = true
			}
			for _, candidate := range []nodes.Node{n, resolvedExpr} {
				cr, ok := candidate.(*nodes.ColumnRef)
				if !ok || cr == nil || cr.Fields == nil || len(cr.Fields.Items) == 0 {
					continue
				}
				if b != nil {
					if cb, found := b.colByRef[cr]; found && cb.relOID != 0 && cb.attNum == 0 && relOIDCount[cb.relOID] == 1 {
						pkCoveredRels[cb.relOID] = true
					}
				}
				lastCol := stringVal(cr.Fields.Items[len(cr.Fields.Items)-1])
				if lastCol != "" {
					groupKeys["vcol:"+lastCol] = true
				}
			}
		}
	}
	if s.GroupClause != nil {
		for _, item := range s.GroupClause.Items {
			addGroupItem(item)
		}
	}

	for relOID, cnt := range relOIDCount {
		if cnt != 1 {
			continue
		}
		for _, con := range c.consByRel[relOID] {
			if con.Type != ConstraintPK || len(con.Columns) == 0 {
				continue
			}
			allGrouped := true
			for _, attNum := range con.Columns {
				if !groupKeys[fmt.Sprintf("col:%d:%d", relOID, attNum)] {
					allGrouped = false
					break
				}
			}
			if allGrouped {
				pkCoveredRels[relOID] = true
				break
			}
		}
	}

	checkColRef := func(cr *nodes.ColumnRef, inOrderByTop bool) error {
		if cr == nil || cr.Fields == nil || len(cr.Fields.Items) == 0 || isStarColumnRef(cr) {
			return nil
		}
		items := cr.Fields.Items
		if len(items) == 4 {
			items = items[1:]
		}
		switch len(items) {
		case 1:
			colName := stringVal(items[0])
			if colName == "" {
				return nil
			}
			if inOrderByTop && ctx.targetAliases[colName] {
				return nil
			}
			if !localScopeHasCol(c, ctx, colName) {
				return nil
			}
			if b != nil {
				if cb, ok := b.colByRef[cr]; ok && cb.relOID != 0 {
					if pkCoveredRels[cb.relOID] || groupKeys[colBindingKey(cb)] {
						return nil
					}
					return &Error{
						Code:    CodeGroupingError,
						Message: fmt.Sprintf("column %q must appear in the GROUP BY clause or be used in an aggregate function", dep.deparseExpr(cr)),
					}
				}
			}
			if groupKeys[exprCanonicalKey(cr)] || groupKeys["vcol:"+colName] {
				return nil
			}
			return &Error{
				Code:    CodeGroupingError,
				Message: fmt.Sprintf("column %q must appear in the GROUP BY clause or be used in an aggregate function", dep.deparseExpr(cr)),
			}
		case 2:
			qual := stringVal(items[0])
			colName := stringVal(items[1])
			if (strings.EqualFold(qual, "old") && b != nil && b.allowOld) || (strings.EqualFold(qual, "new") && b != nil && b.allowNew) {
				return nil
			}
			if _, ok := ctx.scope[qual]; !ok {
				return nil
			}
			if b != nil {
				if cb, ok := b.colByRef[cr]; ok && cb.relOID != 0 && cb.attNum != 0 {
					if (relOIDCount[cb.relOID] == 1 && pkCoveredRels[cb.relOID]) || groupKeys[colBindingKey(cb)] {
						return nil
					}
					return &Error{
						Code:    CodeGroupingError,
						Message: fmt.Sprintf("column %q must appear in the GROUP BY clause or be used in an aggregate function", dep.deparseExpr(cr)),
					}
				}
			}
			if groupKeys[exprCanonicalKey(cr)] || groupKeys["vcol:"+colName] {
				return nil
			}
			return &Error{
				Code:    CodeGroupingError,
				Message: fmt.Sprintf("column %q must appear in the GROUP BY clause or be used in an aggregate function", dep.deparseExpr(cr)),
			}
		case 3:
			relName := stringVal(items[1])
			colName := stringVal(items[2])
			se, ok := ctx.scope[relName]
			if !ok || se.hasAlias {
				return nil
			}
			if b != nil {
				if cb, ok := b.colByRef[cr]; ok && cb.relOID != 0 && cb.attNum != 0 {
					if (relOIDCount[cb.relOID] == 1 && pkCoveredRels[cb.relOID]) || groupKeys[colBindingKey(cb)] {
						return nil
					}
					return &Error{
						Code:    CodeGroupingError,
						Message: fmt.Sprintf("column %q must appear in the GROUP BY clause or be used in an aggregate function", dep.deparseExpr(cr)),
					}
				}
			}
			if groupKeys[exprCanonicalKey(cr)] || groupKeys["vcol:"+colName] {
				return nil
			}
			return &Error{
				Code:    CodeGroupingError,
				Message: fmt.Sprintf("column %q must appear in the GROUP BY clause or be used in an aggregate function", dep.deparseExpr(cr)),
			}
		}
		return nil
	}

	var walk func(n nodes.Node, inOrderByTop bool) error
	walk = func(n nodes.Node, inOrderByTop bool) error {
		if n == nil {
			return nil
		}
		if len(groupKeys) > 0 {
			if k := exprCanonicalKey(n); k != "" && groupKeys[k] {
				return nil
			}
		}
		switch v := n.(type) {
		case *nodes.SubLink:
			return walk(v.Testexpr, false)
		case *nodes.RangeSubselect, *nodes.SelectStmt, *nodes.JsonArrayQueryConstructor, *nodes.GroupingFunc:
			return nil
		case *nodes.List:
			for _, item := range v.Items {
				if err := walk(item, false); err != nil {
					return err
				}
			}
		case *nodes.ResTarget:
			if cr, isCR := v.Val.(*nodes.ColumnRef); isCR && isStarColumnRef(cr) {
				if b != nil && b.starExpandByRT != nil {
					if starCols, found := b.starExpandByRT[v]; found {
						for _, sc := range starCols {
							if sc.relOID != 0 && sc.attNum != 0 {
								key := colBindingKey(ruleColBinding{relOID: sc.relOID, attNum: sc.attNum, prefix: sc.prefix})
								if (relOIDCount[sc.relOID] == 1 && pkCoveredRels[sc.relOID]) || groupKeys[key] {
									continue
								}
								return &Error{
									Code:    CodeGroupingError,
									Message: fmt.Sprintf("column %q must appear in the GROUP BY clause or be used in an aggregate function", sc.origName),
								}
							}
						}
					}
				}
				return nil
			}
			return walk(v.Val, false)
		case *nodes.SortBy:
			return walk(v.Node, true)
		case *nodes.ColumnRef:
			return checkColRef(v, inOrderByTop)
		case *nodes.FuncCall:
			if c.isRuleAggFuncCall(b, v) {
				return nil
			}
			if v.Args != nil {
				if err := walk(v.Args, false); err != nil {
					return err
				}
			}
			if err := walk(v.AggFilter, false); err != nil {
				return err
			}
			if wd, ok := v.Over.(*nodes.WindowDef); ok && wd != nil {
				if wd.PartitionClause != nil {
					if err := walk(wd.PartitionClause, false); err != nil {
						return err
					}
				}
				if wd.OrderClause != nil {
					if err := walk(wd.OrderClause, false); err != nil {
						return err
					}
				}
			}
		case *nodes.A_Expr:
			if err := walk(v.Lexpr, false); err != nil {
				return err
			}
			return walk(v.Rexpr, false)
		case *nodes.BoolExpr:
			if v.Args != nil {
				return walk(v.Args, false)
			}
		case *nodes.NullTest:
			return walk(v.Arg, false)
		case *nodes.BooleanTest:
			return walk(v.Arg, false)
		case *nodes.TypeCast:
			return walk(v.Arg, false)
		case *nodes.CollateClause:
			return walk(v.Arg, false)
		case *nodes.CaseExpr:
			if err := walk(v.Arg, false); err != nil {
				return err
			}
			if v.Args != nil {
				if err := walk(v.Args, false); err != nil {
					return err
				}
			}
			return walk(v.Defresult, false)
		case *nodes.CaseWhen:
			if err := walk(v.Expr, false); err != nil {
				return err
			}
			return walk(v.Result, false)
		case *nodes.CoalesceExpr:
			if v.Args != nil {
				return walk(v.Args, false)
			}
		case *nodes.MinMaxExpr:
			if v.Args != nil {
				return walk(v.Args, false)
			}
		case *nodes.RowExpr:
			if v.Args != nil {
				return walk(v.Args, false)
			}
		case *nodes.A_ArrayExpr:
			if v.Elements != nil {
				return walk(v.Elements, false)
			}
		case *nodes.ArrayExpr:
			if v.Elements != nil {
				return walk(v.Elements, false)
			}
		case *nodes.A_Indirection:
			if err := walk(v.Arg, false); err != nil {
				return err
			}
			if v.Indirection != nil {
				return walk(v.Indirection, false)
			}
		case *nodes.A_Indices:
			if err := walk(v.Lidx, false); err != nil {
				return err
			}
			return walk(v.Uidx, false)
		}
		return nil
	}

	if s.TargetList != nil {
		if err := walk(s.TargetList, false); err != nil {
			return err
		}
	}
	if err := walk(s.HavingClause, false); err != nil {
		return err
	}
	if s.SortClause != nil {
		if err := walk(s.SortClause, false); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) bindWindowDef(b *ruleBindings, wd *nodes.WindowDef, ctx *ruleEvalContext) error {
	return c.bindWindowDefInternal(b, wd, ctx, false)
}

func (c *Catalog) bindWindowClauseDef(b *ruleBindings, wd *nodes.WindowDef, ctx *ruleEvalContext) error {
	return c.bindWindowDefInternal(b, wd, ctx, true)
}

func (c *Catalog) bindWindowDefInternal(b *ruleBindings, wd *nodes.WindowDef, ctx *ruleEvalContext, inWindowClause bool) error {
	if wd == nil {
		return nil
	}
	if !c.snapshot {
		if !inWindowClause && wd.Name != "" {
			if ctx == nil || !ctx.windowNames[wd.Name] {
				return &Error{
					Code:    CodeUndefinedObject,
					Message: fmt.Sprintf("window %q does not exist", wd.Name),
				}
			}
		}
		if wd.Refname != "" {
			if ctx == nil || !ctx.windowNames[wd.Refname] {
				return &Error{
					Code:    CodeUndefinedObject,
					Message: fmt.Sprintf("window %q does not exist", wd.Refname),
				}
			}
			if err := checkRuleWindowCopy(wd, ctx.windowDefs[wd.Refname]); err != nil {
				return err
			}
		}
	}
	if ctx != nil {
		// Window definitions refuse window calls; frame offsets also refuse
		// aggregates and set-returning calls (parse_agg.c, parse_func.c).
		prevAgg, prevWin, prevSRF := ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF
		defer func() { ctx.disallowAgg, ctx.disallowWin, ctx.disallowSRF = prevAgg, prevWin, prevSRF }()
		if ctx.disallowWin == "" {
			ctx.disallowWin = "window definitions"
		}
	}
	if wd.PartitionClause != nil {
		for _, item := range wd.PartitionClause.Items {
			if err := c.bindRuleExpr(b, item, ctx); err != nil {
				return err
			}
		}
	}
	if err := c.bindRuleSortClause(b, wd.OrderClause, ctx); err != nil {
		return err
	}
	if ctx != nil {
		mode := "RANGE"
		if wd.FrameOptions&nodes.FRAMEOPTION_ROWS != 0 {
			mode = "ROWS"
		} else if wd.FrameOptions&nodes.FRAMEOPTION_GROUPS != 0 {
			mode = "GROUPS"
		}
		if ctx.disallowAgg == "" {
			ctx.disallowAgg = "window " + mode
		}
		if ctx.disallowSRF == "" {
			ctx.disallowSRF = "window definitions"
		}
	}
	if err := c.bindRuleExpr(b, wd.StartOffset, ctx); err != nil {
		return err
	}
	if err := c.bindRuleExpr(b, wd.EndOffset, ctx); err != nil {
		return err
	}
	if !c.snapshot {
		// A frame offset may not use the query's own columns
		// (parse_clause.c transformFrameOffset, checkExprIsVarFree).
		for _, off := range []nodes.Node{wd.StartOffset, wd.EndOffset} {
			if off != nil && c.ruleExprRefersToLevel(b, off, ctx) {
				mode := "RANGE"
				if wd.FrameOptions&nodes.FRAMEOPTION_ROWS != 0 {
					mode = "ROWS"
				} else if wd.FrameOptions&nodes.FRAMEOPTION_GROUPS != 0 {
					mode = "GROUPS"
				}
				return &Error{Code: CodeInvalidColumnReference, Message: fmt.Sprintf("argument of %s must not contain variables", mode)}
			}
		}
	}
	return nil
}

// checkRuleWindowCopy applies parse_clause.c's rules for a window that
// copies an earlier one: it may not override PARTITION BY, may add ORDER BY
// only when the copied window has none, and may not copy a window that has a
// frame clause.
func checkRuleWindowCopy(wd, ref *nodes.WindowDef) error {
	if ref == nil {
		return nil
	}
	if wd.PartitionClause != nil && len(wd.PartitionClause.Items) > 0 {
		return &Error{Code: CodeWindowingError, Message: fmt.Sprintf("cannot override PARTITION BY clause of window %q", wd.Refname)}
	}
	if wd.OrderClause != nil && len(wd.OrderClause.Items) > 0 && ref.OrderClause != nil && len(ref.OrderClause.Items) > 0 {
		return &Error{Code: CodeWindowingError, Message: fmt.Sprintf("cannot override ORDER BY clause of window %q", wd.Refname)}
	}
	if ref.FrameOptions&nodes.FRAMEOPTION_NONDEFAULT != 0 {
		return &Error{Code: CodeWindowingError, Message: fmt.Sprintf("cannot copy window %q because it has a frame clause", wd.Refname)}
	}
	return nil
}

func (c *Catalog) checkRuleGroupingFuncs(b *ruleBindings, s *nodes.SelectStmt) error {
	if s == nil {
		return nil
	}
	var gfs []*nodes.GroupingFunc
	var collectGF func(n nodes.Node)
	collectGF = func(n nodes.Node) {
		if n == nil {
			return
		}
		switch v := n.(type) {
		case *nodes.SubLink, *nodes.RangeSubselect, *nodes.SelectStmt, *nodes.JsonArrayQueryConstructor:
			return
		case *nodes.GroupingFunc:
			gfs = append(gfs, v)
			if v.Args != nil {
				collectGF(v.Args)
			}
		case *nodes.List:
			for _, item := range v.Items {
				collectGF(item)
			}
		case *nodes.ResTarget:
			collectGF(v.Val)
		case *nodes.SortBy:
			collectGF(v.Node)
		case *nodes.FuncCall:
			if v.Args != nil {
				collectGF(v.Args)
			}
			collectGF(v.AggFilter)
			if v.AggOrder != nil {
				collectGF(v.AggOrder)
			}
		case *nodes.A_Expr:
			collectGF(v.Lexpr)
			collectGF(v.Rexpr)
		case *nodes.BoolExpr:
			if v.Args != nil {
				collectGF(v.Args)
			}
		case *nodes.NullTest:
			collectGF(v.Arg)
		case *nodes.BooleanTest:
			collectGF(v.Arg)
		case *nodes.TypeCast:
			collectGF(v.Arg)
		case *nodes.CollateClause:
			collectGF(v.Arg)
		case *nodes.CaseExpr:
			collectGF(v.Arg)
			if v.Args != nil {
				collectGF(v.Args)
			}
			collectGF(v.Defresult)
		case *nodes.CaseWhen:
			collectGF(v.Expr)
			collectGF(v.Result)
		case *nodes.CoalesceExpr:
			if v.Args != nil {
				collectGF(v.Args)
			}
		case *nodes.MinMaxExpr:
			if v.Args != nil {
				collectGF(v.Args)
			}
		case *nodes.RowExpr:
			if v.Args != nil {
				collectGF(v.Args)
			}
		case *nodes.A_ArrayExpr:
			if v.Elements != nil {
				collectGF(v.Elements)
			}
		case *nodes.ArrayExpr:
			if v.Elements != nil {
				collectGF(v.Elements)
			}
		}
	}
	if s.TargetList != nil {
		collectGF(s.TargetList)
	}
	collectGF(s.HavingClause)
	if s.SortClause != nil {
		collectGF(s.SortClause)
	}
	if len(gfs) == 0 {
		return nil
	}

	type targetKeyEntry struct {
		alias string
		key   string
	}
	dep := ruleDeparser{cat: c, b: b}
	exprCanonicalKey := func(expr nodes.Node) string {
		if expr == nil {
			return ""
		}
		if cr, ok := expr.(*nodes.ColumnRef); ok && b != nil {
			if cb, found := b.colByRef[cr]; found && cb.relOID != 0 && cb.attNum != 0 {
				return fmt.Sprintf("col:%d:%d", cb.relOID, cb.attNum)
			}
		}
		return "expr:" + dep.deparseExpr(expr)
	}
	var targets []targetKeyEntry
	if s.TargetList != nil {
		for _, item := range s.TargetList.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok || rt == nil {
				continue
			}
			if cr, isCR := rt.Val.(*nodes.ColumnRef); isCR && isStarColumnRef(cr) {
				if b != nil && b.starExpandByRT != nil {
					if starCols, found := b.starExpandByRT[rt]; found && len(starCols) > 0 {
						for _, sc := range starCols {
							colName := sc.colAlias
							if colName == "" {
								colName = sc.origName
							}
							key := "expr:" + quoteIdent(colName)
							if sc.relOID != 0 && sc.attNum != 0 {
								key = fmt.Sprintf("col:%d:%d", sc.relOID, sc.attNum)
							}
							targets = append(targets, targetKeyEntry{alias: colName, key: key})
						}
						continue
					}
				}
				return nil
			}
			alias := rt.Name
			if alias == "" && b != nil && b.targetNameByRT != nil {
				alias = b.targetNameByRT[rt]
			}
			targets = append(targets, targetKeyEntry{
				alias: alias,
				key:   exprCanonicalKey(rt.Val),
			})
		}
	}
	resolveSortGroupKey := func(expr nodes.Node) string {
		var iv *nodes.Integer
		switch v := expr.(type) {
		case *nodes.A_Const:
			if !v.Isnull {
				iv, _ = v.Val.(*nodes.Integer)
			}
		case *nodes.Integer:
			iv = v
		}
		if iv != nil && iv.Ival >= 1 && int(iv.Ival) <= len(targets) {
			return targets[iv.Ival-1].key
		}
		if cr, ok := expr.(*nodes.ColumnRef); ok && cr != nil && cr.Fields != nil && len(cr.Fields.Items) == 1 {
			if b == nil || b.colByRef[cr].relOID == 0 {
				if str, isStr := cr.Fields.Items[0].(*nodes.String); isStr && str.Str != "" {
					for _, t := range targets {
						if t.alias == str.Str {
							return t.key
						}
					}
				}
			}
		}
		return exprCanonicalKey(expr)
	}

	groupKeys := make(map[string]bool)
	var addGroupItem func(n nodes.Node)
	addGroupItem = func(n nodes.Node) {
		if n == nil {
			return
		}
		switch v := n.(type) {
		case *nodes.GroupingSet:
			if v.Content != nil {
				for _, item := range v.Content.Items {
					addGroupItem(item)
				}
			}
		case *nodes.RowExpr:
			if v.Args != nil {
				for _, item := range v.Args.Items {
					addGroupItem(item)
				}
			}
		default:
			if k := resolveSortGroupKey(n); k != "" {
				groupKeys[k] = true
			}
		}
	}
	if s.GroupClause != nil {
		for _, item := range s.GroupClause.Items {
			addGroupItem(item)
		}
	}
	for _, gf := range gfs {
		if gf.Args == nil || len(gf.Args.Items) == 0 {
			return &Error{
				Code:    CodeGroupingError,
				Message: "arguments to GROUPING must be grouping expressions of the associated query level",
			}
		}
		for _, arg := range gf.Args.Items {
			k := resolveSortGroupKey(arg)
			if !groupKeys[k] {
				return &Error{
					Code:    CodeGroupingError,
					Message: "arguments to GROUPING must be grouping expressions of the associated query level",
				}
			}
		}
	}
	return nil
}

func lookupScopeEntry(ctx *ruleEvalContext, name string) (ruleScopeEntry, bool) {
	for cur := ctx; cur != nil; cur = cur.outer {
		if cur.scope != nil {
			if se, ok := cur.scope[name]; ok {
				return se, true
			}
		}
	}
	return ruleScopeEntry{}, false
}

func resolveScopeColumn(rel *Relation, se ruleScopeEntry, colName string) (int16, bool) {
	if rel == nil {
		return 0, false
	}
	if len(se.colAliases) > 0 {
		for i, aliasName := range se.colAliases {
			if aliasName == colName && i < len(rel.Columns) {
				return rel.Columns[i].AttNum, true
			}
		}
	}
	if idx, exists := rel.colByName[colName]; exists && idx >= len(se.colAliases) {
		return rel.Columns[idx].AttNum, true
	}
	return 0, false
}

func recordAllRelColumnDeps(c *Catalog, b *ruleBindings, relOID uint32) {
	rel := c.relationByOID[relOID]
	if rel == nil {
		return
	}
	for _, col := range rel.Columns {
		b.addColDep(relOID, col.AttNum)
	}
}

func (c *Catalog) bindRuleExpr(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext) error {
	if n == nil {
		return nil
	}
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if v.Fields == nil || len(v.Fields.Items) == 0 {
			return nil
		}
		items := v.Fields.Items
		if len(items) == 1 {
			if _, isStar := items[0].(*nodes.A_Star); isStar {
				hasAny := false
				if ctx != nil {
					for k, se := range ctx.scope {
						if k == "excluded" {
							continue
						}
						hasAny = true
						if se.relOID != 0 {
							recordAllRelColumnDeps(c, b, se.relOID)
						}
					}
					if ctx.hasOpenSource {
						hasAny = true
					}
				}
				if !hasAny && !c.snapshot {
					return &Error{Code: CodeSyntaxError, Message: "SELECT * with no tables specified is not valid"}
				}
				return nil
			}
			colName := stringVal(items[0])
			if colName == "" {
				return nil
			}
			for cur := ctx; cur != nil; cur = cur.outer {
				var matchedOID uint32
				var matchedAtt int16
				var matchedType uint32
				var matchedSE ruleScopeEntry
				count := 0
				hasAnyOpen := cur.hasOpenSource
				seenRel := map[uint32]bool{}
				for k, se := range cur.scope {
					if k == "excluded" || se.hiddenByJoinColAlias {
						continue
					}
					if se.isVirtual {
						if se.isJoinAlias && !se.hasJoinColAliases {
							continue
						}
						if nMatches := se.virtualCols[colName]; nMatches > 0 {
							count += nMatches
							if matchedType == 0 && se.virtualColTypes != nil {
								matchedType = se.virtualColTypes[colName]
							}
							for _, rOID := range se.joinRels {
								if rel := c.relationByOID[rOID]; rel != nil {
									if idx, ok := rel.colByName[colName]; ok {
										b.addColDep(rOID, rel.Columns[idx].AttNum)
										if matchedType == 0 {
											matchedType = rel.Columns[idx].TypeOID
										}
									} else {
										recordAllRelColumnDeps(c, b, rOID)
									}
								}
							}
						} else if se.anyCol {
							hasAnyOpen = true
						}
						continue
					}
					if seenRel[se.relOID] {
						continue
					}
					seenRel[se.relOID] = true
					if rel := c.relationByOID[se.relOID]; rel != nil {
						if attNum, ok := resolveScopeColumn(rel, se, colName); ok {
							matchedOID = se.relOID
							matchedAtt = attNum
							matchedSE = se
							for _, col := range rel.Columns {
								if col.AttNum == attNum {
									matchedType = col.TypeOID
									break
								}
							}
							count++
							b.addColDep(se.relOID, attNum)
						}
					}
				}
				if count > 1 && cur.usingCols[colName] > 0 {
					count -= cur.usingCols[colName]
					if count < 1 {
						count = 1
					}
				}
				if count > 1 {
					return &Error{
						Code:    CodeAmbiguousColumn,
						Message: fmt.Sprintf("column reference %q is ambiguous", colName),
					}
				}
				if count == 1 {
					if matchedOID != 0 && len(cur.colAliasesFor(matchedOID)) == 0 && cur.usingCols[colName] == 0 {
						cb := ruleColBinding{relOID: matchedOID, attNum: matchedAtt}
						if b.fn != nil {
							cb.prefix, cb.useRelName = matchedSE.alias, !matchedSE.hasAlias
						}
						b.colByRef[v] = cb
					}
					if matchedType != 0 && b != nil && b.exprTypeByNode != nil {
						b.exprTypeByNode[v] = matchedType
					}
					return nil
				}
				if se, ok := cur.scope[colName]; ok && colName != "excluded" {
					if !se.isVirtual && se.relOID != 0 {
						b.addRelDep(se.relOID)
						b.colByRef[v] = ruleColBinding{
							relOID:     se.relOID,
							attNum:     0,
							prefix:     se.alias,
							useRelName: !se.hasAlias,
						}
						if rel := c.relationByOID[se.relOID]; rel != nil && rel.RowTypeOID != 0 && b != nil && b.exprTypeByNode != nil {
							b.exprTypeByNode[v] = rel.RowTypeOID
						}
						return nil
					}
					for _, rOID := range se.joinRels {
						b.addRelDep(rOID)
					}
					if b != nil && b.exprTypeByNode != nil {
						b.exprTypeByNode[v] = RECORDOID
					}
					return nil
				}
				if cur.targetAliases[colName] || hasAnyOpen {
					return nil
				}
			}
			if num := b.fn.lookup(colName); num > 0 {
				b.bindFnParam(v, num, ctx)
				return nil
			}
			isOld := strings.EqualFold(colName, "old")
			isNew := strings.EqualFold(colName, "new")
			if (isOld && b.allowOld) || (isNew && b.allowNew) {
				b.colByRef[v] = ruleColBinding{
					relOID: b.ownerRelOID,
					attNum: 0,
					prefix: strings.ToLower(colName),
				}
				if ownerRel := c.relationByOID[b.ownerRelOID]; ownerRel != nil && ownerRel.RowTypeOID != 0 && b != nil && b.exprTypeByNode != nil {
					b.exprTypeByNode[v] = ownerRel.RowTypeOID
				}
				return nil
			}
			if !c.snapshot {
				return errUndefinedColumn(colName)
			}
			return nil
		}

		if len(items) == 2 {
			qual := stringVal(items[0])
			_, isStar := items[1].(*nodes.A_Star)
			colName := stringVal(items[1])
			isOld := strings.EqualFold(qual, "old")
			isNew := strings.EqualFold(qual, "new")
			if (isOld && b.allowOld) || (isNew && b.allowNew) {
				ownerRel := c.relationByOID[b.ownerRelOID]
				if ownerRel == nil {
					return nil
				}
				if isStar {
					recordAllRelColumnDeps(c, b, b.ownerRelOID)
					return nil
				}
				idx, exists := ownerRel.colByName[colName]
				if !exists {
					if !c.snapshot {
						return errUndefinedColumn(colName)
					}
					return nil
				}
				attNum := ownerRel.Columns[idx].AttNum
				b.colByRef[v] = ruleColBinding{
					relOID: b.ownerRelOID,
					attNum: attNum,
					prefix: strings.ToLower(qual),
				}
				if b != nil && b.exprTypeByNode != nil {
					b.exprTypeByNode[v] = ownerRel.Columns[idx].TypeOID
				}
				b.addColDep(b.ownerRelOID, attNum)
				return nil
			}
			if se, ok := lookupScopeEntry(ctx, qual); ok {
				if se.isVirtual {
					if !isStar && se.virtualColTypes != nil && b != nil && b.exprTypeByNode != nil && b.exprTypeByNode[v] == 0 {
						b.exprTypeByNode[v] = se.virtualColTypes[colName]
					}
					for _, rOID := range se.joinRels {
						if isStar {
							recordAllRelColumnDeps(c, b, rOID)
						} else if rel := c.relationByOID[rOID]; rel != nil {
							if idx, exists := rel.colByName[colName]; exists {
								b.addColDep(rOID, rel.Columns[idx].AttNum)
								if b != nil && b.exprTypeByNode != nil && b.exprTypeByNode[v] == 0 {
									b.exprTypeByNode[v] = rel.Columns[idx].TypeOID
								}
							} else {
								recordAllRelColumnDeps(c, b, rOID)
							}
						}
					}
					if !isStar && !c.snapshot {
						if se.virtualCols[colName] > 1 {
							return &Error{
								Code:    CodeAmbiguousColumn,
								Message: fmt.Sprintf("column reference %q is ambiguous", qual+"."+colName),
							}
						}
						if !se.anyCol && se.virtualCols[colName] == 0 {
							return errUndefinedColumn(colName)
						}
					}
					return nil
				}
				rel := c.relationByOID[se.relOID]
				if rel == nil {
					return nil
				}
				if isStar {
					recordAllRelColumnDeps(c, b, se.relOID)
					return nil
				}
				attNum, exists := resolveScopeColumn(rel, se, colName)
				if !exists {
					if !c.snapshot {
						return errUndefinedColumn(colName)
					}
					return nil
				}
				if len(se.colAliases) == 0 {
					b.colByRef[v] = ruleColBinding{
						relOID:     se.relOID,
						attNum:     attNum,
						prefix:     se.alias,
						useRelName: !se.hasAlias,
					}
				}
				if b != nil && b.exprTypeByNode != nil {
					for _, col := range rel.Columns {
						if col.AttNum == attNum {
							b.exprTypeByNode[v] = col.TypeOID
							break
						}
					}
				}
				b.addColDep(se.relOID, attNum)
				return nil
			}
			if b.fn != nil && !isStar && qual == b.fn.name {
				if num := b.fn.lookup(colName); num > 0 {
					b.bindFnParam(v, num, ctx)
					return nil
				}
			}
			if !c.snapshot {
				if isOld && b.isOnSelect {
					return &Error{
						Code:    CodeInvalidObjectDefinition,
						Message: "ON SELECT rule cannot use OLD",
					}
				}
				if isOld || isNew {
					return &Error{
						Code:    CodeUndefinedTable,
						Message: fmt.Sprintf("invalid reference to FROM-clause entry for table %q", strings.ToLower(qual)),
					}
				}
				return &Error{
					Code:    CodeUndefinedTable,
					Message: fmt.Sprintf("missing FROM-clause entry for table %q", qual),
				}
			}
			return nil
		}

		if len(items) > 4 {
			if !c.snapshot {
				return &Error{Code: CodeSyntaxError, Message: "improper qualified name (too many dotted names)"}
			}
			return nil
		}
		if len(items) == 4 {
			if err := c.checkDatabaseQualifier(stringVal(items[0])); err != nil {
				return err
			}
			items = items[1:]
		}
		if len(items) == 3 {
			schName := stringVal(items[0])
			relName := stringVal(items[1])
			_, isStar := items[2].(*nodes.A_Star)
			colName := stringVal(items[2])
			se, ok := lookupScopeEntry(ctx, relName)
			if !ok || se.hasAlias {
				if !c.snapshot {
					return &Error{
						Code:    CodeUndefinedTable,
						Message: fmt.Sprintf("missing FROM-clause entry for table %q", relName),
					}
				}
				return nil
			}
			if se.isVirtual {
				return nil
			}
			rel := c.relationByOID[se.relOID]
			if rel == nil {
				return nil
			}
			if rel.Schema != nil && rel.Schema.Name != schName {
				if !c.snapshot {
					return &Error{
						Code:    CodeUndefinedTable,
						Message: fmt.Sprintf("missing FROM-clause entry for table %q", relName),
					}
				}
				return nil
			}
			if isStar {
				recordAllRelColumnDeps(c, b, se.relOID)
				return nil
			}
			attNum, exists := resolveScopeColumn(rel, se, colName)
			if !exists {
				if !c.snapshot {
					return errUndefinedColumn(colName)
				}
				return nil
			}
			if len(se.colAliases) == 0 {
				b.colByRef[v] = ruleColBinding{
					relOID:     se.relOID,
					attNum:     attNum,
					prefix:     se.alias,
					useRelName: !se.hasAlias,
				}
			}
			if b != nil && b.exprTypeByNode != nil {
				for _, col := range rel.Columns {
					if col.AttNum == attNum {
						b.exprTypeByNode[v] = col.TypeOID
						break
					}
				}
			}
			b.addColDep(se.relOID, attNum)
			return nil
		}
		return nil

	case *nodes.TypeCast:
		if err := c.bindRuleExpr(b, v.Arg, ctx); err != nil {
			return err
		}
		if v.TypeName != nil {
			oid, _, err := c.resolveTypeName(v.TypeName)
			if err != nil {
				if !c.snapshot {
					return err
				}
				return nil
			}
			if oid != 0 {
				isArr := v.TypeName.ArrayBounds != nil && len(v.TypeName.ArrayBounds.Items) > 0
				elemOID := oid
				if bt := c.typeByOID[oid]; bt != nil && bt.Elem != 0 {
					elemOID = bt.Elem
				}
				b.typeByTN[v.TypeName] = ruleTypeBinding{typeOID: elemOID, isArray: isArr}
				b.addTypeDep(elemOID)
				if oid == REGCLASSOID {
					if err := c.bindRegclassLiteral(b, v.Arg); err != nil {
						return err
					}
				}
			}
		}
		return nil

	case *nodes.CollateClause:
		if err := c.bindRuleExpr(b, v.Arg, ctx); err != nil {
			return err
		}
		_, uid, err := c.resolveCollationList(v.Collname)
		if err != nil {
			if !c.snapshot {
				return err
			}
			return nil
		}
		if uid != 0 {
			b.collByCC[v] = uid
			b.addCollDep(uid)
		}
		if !c.snapshot {
			argType := c.inferRuleExprType(b, v.Arg, ctx)
			if argType != 0 && argType != UNKNOWNOID && c.typeCollation(argType) == 0 {
				return &Error{
					Code:    CodeDatatypeMismatch,
					Message: fmt.Sprintf("collations are not supported by type %s", c.typeName(argType)),
				}
			}
		}
		return nil

	case *nodes.FuncCall:
		var argCount int
		if v.Args != nil {
			argCount = len(v.Args.Items)
			for _, arg := range v.Args.Items {
				if err := c.bindRuleExpr(b, arg, ctx); err != nil {
					return err
				}
			}
		}
		if err := c.bindRuleSortClause(b, v.AggOrder, ctx); err != nil {
			return err
		}
		if err := c.bindRuleBoolExpr(b, v.AggFilter, ctx, "FILTER"); err != nil {
			return err
		}
		if wd, ok := v.Over.(*nodes.WindowDef); ok {
			if err := c.bindWindowDef(b, wd, ctx); err != nil {
				return err
			}
		}
		if v.Funcname != nil {
			if err := c.bindRuleFuncCall(b, v, argCount, ctx); err != nil {
				return err
			}
		}
		if !c.snapshot {
			return c.checkRuleNestedAggOrWindow(b, v)
		}
		return nil

	case *nodes.NamedArgExpr:
		return c.bindRuleExpr(b, v.Arg, ctx)
	case *nodes.A_Indirection:
		if err := c.bindRuleExpr(b, v.Arg, ctx); err != nil {
			return err
		}
		if v.Indirection != nil {
			for _, item := range v.Indirection.Items {
				if idx, ok := item.(*nodes.A_Indices); ok {
					if err := c.bindRuleExpr(b, idx.Lidx, ctx); err != nil {
						return err
					}
					if err := c.bindRuleExpr(b, idx.Uidx, ctx); err != nil {
						return err
					}
				}
			}
		}
		return nil
	case *nodes.A_Expr:
		if err := c.bindRuleExpr(b, v.Lexpr, ctx); err != nil {
			return err
		}
		if list, ok := v.Rexpr.(*nodes.List); ok {
			for _, item := range list.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		} else if err := c.bindRuleExpr(b, v.Rexpr, ctx); err != nil {
			return err
		}
		if v.Name != nil {
			switch v.Kind {
			case nodes.AEXPR_OP:
				isPrefix := v.Lexpr == nil
				lType := c.inferRuleExprType(b, v.Lexpr, ctx)
				rType := c.inferRuleExprType(b, v.Rexpr, ctx)
				resType, err := c.bindRuleOperator(b, v.Name, lType, rType, isPrefix, true)
				if err != nil {
					return err
				}
				if resType != 0 && b != nil && b.exprTypeByNode != nil {
					b.exprTypeByNode[v] = resType
				}
			case nodes.AEXPR_OP_ANY, nodes.AEXPR_OP_ALL:
				lType := c.inferRuleExprType(b, v.Lexpr, ctx)
				rArrType := c.inferRuleExprType(b, v.Rexpr, ctx)
				var rElemType uint32
				if bt := c.typeByOID[rArrType]; bt != nil && bt.Elem != 0 {
					rElemType = bt.Elem
				}
				if _, err := c.bindRuleOperator(b, v.Name, lType, rElemType, false, true); err != nil {
					return err
				}
				if b != nil && b.exprTypeByNode != nil {
					b.exprTypeByNode[v] = BOOLOID
				}
			case nodes.AEXPR_IN:
				lType := c.inferRuleExprType(b, v.Lexpr, ctx)
				if list, ok := v.Rexpr.(*nodes.List); ok {
					for _, item := range list.Items {
						rType := c.inferRuleExprType(b, item, ctx)
						if _, err := c.bindRuleOperator(b, v.Name, lType, rType, false, true); err != nil {
							return err
						}
					}
				}
				if b != nil && b.exprTypeByNode != nil {
					b.exprTypeByNode[v] = BOOLOID
				}
			case nodes.AEXPR_DISTINCT, nodes.AEXPR_NOT_DISTINCT, nodes.AEXPR_NULLIF:
				lType := c.inferRuleExprType(b, v.Lexpr, ctx)
				rType := c.inferRuleExprType(b, v.Rexpr, ctx)
				if _, err := c.bindRuleOperator(b, v.Name, lType, rType, false, true); err != nil {
					return err
				}
			case nodes.AEXPR_BETWEEN, nodes.AEXPR_NOT_BETWEEN, nodes.AEXPR_BETWEEN_SYM, nodes.AEXPR_NOT_BETWEEN_SYM:
				lType := c.inferRuleExprType(b, v.Lexpr, ctx)
				if list, ok := v.Rexpr.(*nodes.List); ok && len(list.Items) == 2 {
					lowType := c.inferRuleExprType(b, list.Items[0], ctx)
					highType := c.inferRuleExprType(b, list.Items[1], ctx)
					if _, err := c.bindRuleOperator(b, &nodes.List{Items: []nodes.Node{&nodes.String{Str: ">="}}}, lType, lowType, false, true); err != nil {
						return err
					}
					if _, err := c.bindRuleOperator(b, &nodes.List{Items: []nodes.Node{&nodes.String{Str: "<="}}}, lType, highType, false, true); err != nil {
						return err
					}
				}
				if b != nil && b.exprTypeByNode != nil {
					b.exprTypeByNode[v] = BOOLOID
				}
			}
		}
		return nil
	case *nodes.BoolExpr:
		if v.Args != nil {
			construct := "AND"
			switch v.Boolop {
			case nodes.OR_EXPR:
				construct = "OR"
			case nodes.NOT_EXPR:
				construct = "NOT"
			}
			for _, item := range v.Args.Items {
				if err := c.bindRuleBoolExpr(b, item, ctx, construct); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.NullTest:
		return c.bindRuleExpr(b, v.Arg, ctx)
	case *nodes.BooleanTest:
		return c.bindRuleBoolExpr(b, v.Arg, ctx, "IS")
	case *nodes.CoalesceExpr:
		if ctx != nil && ctx.disallowSRF == "" {
			ctx.disallowSRF = "COALESCE"
			defer func() { ctx.disallowSRF = "" }()
		}
		if v.Args != nil {
			var exprs []AnalyzedExpr
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
				if t := c.inferRuleExprType(b, item, ctx); t != 0 && t != UNKNOWNOID {
					exprs = append(exprs, &ConstExpr{TypeOID: t, TypeMod: -1})
				}
			}
			if !c.snapshot && len(exprs) > 1 {
				if _, err := c.selectCommonType(exprs, "COALESCE"); err != nil {
					return errDatatypeMismatch(err.Error())
				}
			}
		}
		return nil
	case *nodes.MinMaxExpr:
		if v.Args != nil {
			var exprs []AnalyzedExpr
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
				if t := c.inferRuleExprType(b, item, ctx); t != 0 && t != UNKNOWNOID {
					exprs = append(exprs, &ConstExpr{TypeOID: t, TypeMod: -1})
				}
			}
			if !c.snapshot && len(exprs) > 1 {
				if _, err := c.selectCommonType(exprs, "GREATEST/LEAST"); err != nil {
					return errDatatypeMismatch(err.Error())
				}
			}
		}
		return nil
	case *nodes.CaseExpr:
		if ctx != nil && ctx.disallowSRF == "" {
			ctx.disallowSRF = "CASE"
			defer func() { ctx.disallowSRF = "" }()
		}
		if err := c.bindRuleExpr(b, v.Arg, ctx); err != nil {
			return err
		}
		argType := c.inferRuleExprType(b, v.Arg, ctx)
		var resExprs []AnalyzedExpr
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if cw, ok := item.(*nodes.CaseWhen); ok {
					if v.Arg == nil {
						if err := c.bindRuleBoolExpr(b, cw.Expr, ctx, "CASE/WHEN"); err != nil {
							return err
						}
					} else {
						if err := c.bindRuleExpr(b, cw.Expr, ctx); err != nil {
							return err
						}
						whenType := c.inferRuleExprType(b, cw.Expr, ctx)
						if _, err := c.bindRuleOperator(b, &nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}}, argType, whenType, false, true); err != nil {
							return err
						}
					}
					if err := c.bindRuleExpr(b, cw.Result, ctx); err != nil {
						return err
					}
					if t := c.inferRuleExprType(b, cw.Result, ctx); t != 0 && t != UNKNOWNOID {
						resExprs = append(resExprs, &ConstExpr{TypeOID: t, TypeMod: -1})
					}
				}
			}
		}
		if err := c.bindRuleExpr(b, v.Defresult, ctx); err != nil {
			return err
		}
		if v.Defresult != nil {
			if t := c.inferRuleExprType(b, v.Defresult, ctx); t != 0 && t != UNKNOWNOID {
				resExprs = append(resExprs, &ConstExpr{TypeOID: t, TypeMod: -1})
			}
		}
		if !c.snapshot && len(resExprs) > 1 {
			if _, err := c.selectCommonType(resExprs, "CASE"); err != nil {
				return errDatatypeMismatch(err.Error())
			}
		}
		return nil
	case *nodes.A_ArrayExpr:
		if v.Elements != nil {
			var exprs []AnalyzedExpr
			for _, item := range v.Elements.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
				if t := c.inferRuleExprType(b, item, ctx); t != 0 && t != UNKNOWNOID {
					exprs = append(exprs, &ConstExpr{TypeOID: t, TypeMod: -1})
				}
			}
			if !c.snapshot && len(exprs) > 1 {
				if _, err := c.selectCommonType(exprs, "ARRAY"); err != nil {
					return errDatatypeMismatch(err.Error())
				}
			}
		}
		return nil
	case *nodes.RowExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.SubLink:
		inRuleWhere := b != nil && b.inRuleWhere
		if inRuleWhere {
			b.inRuleWhere = false
			defer func() { b.inRuleWhere = true }()
		}
		if err := c.bindRuleExpr(b, v.Testexpr, ctx); err != nil {
			return err
		}
		subCtx := &ruleEvalContext{}
		if ctx != nil {
			subCtx.outer = ctx
			subCtx.cteNames = ctx.cteNames
			subCtx.cteCols = ctx.cteCols
			subCtx.cteColTypes = ctx.cteColTypes
			subCtx.cteAny = ctx.cteAny
		}
		if err := c.bindRuleStmt(b, v.Subselect, subCtx); err != nil {
			return err
		}
		qCols, qAny := c.extractRuleQueryColTypes(b, v.Subselect)
		if !c.snapshot && !qAny {
			switch nodes.SubLinkType(v.SubLinkType) {
			case nodes.EXPR_SUBLINK, nodes.ARRAY_SUBLINK:
				if len(qCols) > 1 {
					return &Error{
						Code:    CodeSyntaxError,
						Message: "subquery must return only one column",
					}
				}
				if len(qCols) < 1 {
					return &Error{
						Code:    CodeSyntaxError,
						Message: "subquery must return a column",
					}
				}
			case nodes.ANY_SUBLINK, nodes.ALL_SUBLINK:
				if row, isRow := v.Testexpr.(*nodes.RowExpr); isRow && row != nil && row.Args != nil {
					if len(row.Args.Items) > len(qCols) {
						return &Error{Code: CodeSyntaxError, Message: "subquery has too few columns"}
					}
					if len(row.Args.Items) < len(qCols) {
						return &Error{Code: CodeSyntaxError, Message: "subquery has too many columns"}
					}
				} else if len(qCols) > 1 {
					return &Error{
						Code:    CodeSyntaxError,
						Message: "subquery must return only one column",
					}
				}
			}
		}
		if v.OperName != nil {
			lType := c.inferRuleExprType(b, v.Testexpr, ctx)
			var rType uint32
			if !qAny && len(qCols) == 1 {
				rType = qCols[0]
			}
			if _, err := c.bindRuleOperator(b, v.OperName, lType, rType, false, true); err != nil {
				return err
			}
		} else if (nodes.SubLinkType(v.SubLinkType) == nodes.ANY_SUBLINK || nodes.SubLinkType(v.SubLinkType) == nodes.ALL_SUBLINK) && !qAny && len(qCols) == 1 {
			if _, isRow := v.Testexpr.(*nodes.RowExpr); !isRow {
				lType := c.inferRuleExprType(b, v.Testexpr, ctx)
				if _, err := c.bindRuleOperator(b, &nodes.List{Items: []nodes.Node{&nodes.String{Str: "="}}}, lType, qCols[0], false, true); err != nil {
					return err
				}
			}
		}
		if inRuleWhere && !c.snapshot {
			return &Error{
				Code:    CodeFeatureNotSupported,
				Message: "subqueries in rule WHERE clause are not supported",
			}
		}
		return nil
	case *nodes.MultiAssignRef:
		return c.bindRuleExpr(b, v.Source, ctx)
	case *nodes.ArrayExpr:
		if v.Elements != nil {
			for _, item := range v.Elements.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.GroupingFunc:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		if !c.snapshot {
			disallowAgg := ""
			if b != nil && b.inRuleWhere {
				disallowAgg = "WHERE"
			} else if ctx != nil {
				disallowAgg = ctx.disallowAgg
			}
			if disallowAgg != "" {
				return &Error{
					Code:    CodeGroupingError,
					Message: fmt.Sprintf("grouping operations are not allowed in %s", disallowAgg),
				}
			}
		}
		return nil
	case *nodes.GroupingSet:
		if v.Content != nil {
			for _, item := range v.Content.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.XmlExpr:
		if v.NamedArgs != nil {
			for _, item := range v.NamedArgs.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					if err := c.bindRuleExpr(b, rt.Val, ctx); err != nil {
						return err
					}
				} else if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.XmlSerialize:
		if err := c.bindRuleExpr(b, v.Expr, ctx); err != nil {
			return err
		}
		if v.TypeName != nil {
			oid, _, err := c.resolveTypeName(v.TypeName)
			if err != nil {
				if !c.snapshot {
					return err
				}
				return nil
			}
			if oid != 0 {
				isArr := v.TypeName.ArrayBounds != nil && len(v.TypeName.ArrayBounds.Items) > 0
				elemOID := oid
				if bt := c.typeByOID[oid]; bt != nil && bt.Elem != 0 {
					elemOID = bt.Elem
				}
				b.typeByTN[v.TypeName] = ruleTypeBinding{typeOID: elemOID, isArray: isArr}
				b.addTypeDep(elemOID)
			}
		}
		return nil
	case *nodes.JsonValueExpr:
		if v == nil {
			return nil
		}
		return c.bindRuleExpr(b, v.RawExpr, ctx)
	case *nodes.JsonIsPredicate:
		return c.bindRuleExpr(b, v.Expr, ctx)
	case *nodes.JsonParseExpr:
		if err := c.bindRuleExpr(b, v.Expr, ctx); err != nil {
			return err
		}
		return c.bindRuleJsonOutput(b, v.Output)
	case *nodes.JsonScalarExpr:
		if err := c.bindRuleExpr(b, v.Expr, ctx); err != nil {
			return err
		}
		return c.bindRuleJsonOutput(b, v.Output)
	case *nodes.JsonSerializeExpr:
		if err := c.bindRuleExpr(b, v.Expr, ctx); err != nil {
			return err
		}
		return c.bindRuleJsonOutput(b, v.Output)
	case *nodes.JsonObjectConstructor:
		if v.Exprs != nil {
			for _, item := range v.Exprs.Items {
				if kv, ok := item.(*nodes.JsonKeyValue); ok {
					if err := c.bindRuleExpr(b, kv.Key, ctx); err != nil {
						return err
					}
					if err := c.bindRuleExpr(b, kv.Value, ctx); err != nil {
						return err
					}
				}
			}
		}
		return c.bindRuleJsonOutput(b, v.Output)
	case *nodes.JsonArrayConstructor:
		if v.Exprs != nil {
			for _, item := range v.Exprs.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return c.bindRuleJsonOutput(b, v.Output)
	case *nodes.JsonArrayQueryConstructor:
		inRuleWhere := b != nil && b.inRuleWhere
		if inRuleWhere {
			b.inRuleWhere = false
			defer func() { b.inRuleWhere = true }()
		}
		subCtx := &ruleEvalContext{}
		if ctx != nil {
			subCtx.outer = ctx
			subCtx.cteNames = ctx.cteNames
			subCtx.cteCols = ctx.cteCols
			subCtx.cteAny = ctx.cteAny
		}
		if err := c.bindRuleStmt(b, v.Query, subCtx); err != nil {
			return err
		}
		if err := c.bindRuleJsonOutput(b, v.Output); err != nil {
			return err
		}
		if inRuleWhere && !c.snapshot {
			return &Error{
				Code:    CodeFeatureNotSupported,
				Message: "subqueries in rule WHERE clause are not supported",
			}
		}
		return nil
	case *nodes.JsonFuncExpr:
		if err := c.bindRuleExpr(b, v.ContextItem, ctx); err != nil {
			return err
		}
		if err := c.bindRuleExpr(b, v.Pathspec, ctx); err != nil {
			return err
		}
		if v.Passing != nil {
			for _, item := range v.Passing.Items {
				if ja, ok := item.(*nodes.JsonArgument); ok {
					if err := c.bindRuleExpr(b, ja.Val, ctx); err != nil {
						return err
					}
				}
			}
		}
		if v.OnEmpty != nil {
			if err := c.bindRuleExpr(b, v.OnEmpty.Expr, ctx); err != nil {
				return err
			}
		}
		if v.OnError != nil {
			if err := c.bindRuleExpr(b, v.OnError.Expr, ctx); err != nil {
				return err
			}
		}
		return c.bindRuleJsonOutput(b, v.Output)
	case *nodes.JsonObjectAgg:
		if v.Arg != nil {
			if err := c.bindRuleExpr(b, v.Arg.Key, ctx); err != nil {
				return err
			}
			if err := c.bindRuleExpr(b, v.Arg.Value, ctx); err != nil {
				return err
			}
		}
		return c.bindRuleJsonAggConstructor(b, v.Constructor, ctx)
	case *nodes.JsonArrayAgg:
		if err := c.bindRuleExpr(b, v.Arg, ctx); err != nil {
			return err
		}
		return c.bindRuleJsonAggConstructor(b, v.Constructor, ctx)
	case *nodes.ParamRef:
		if b != nil && b.fn != nil {
			if v.Number < 1 || v.Number > len(b.fn.types) {
				return &Error{
					Code:    CodeUndefinedParameter,
					Message: fmt.Sprintf("there is no parameter $%d", v.Number),
				}
			}
			b.bindFnParam(v, v.Number, ctx)
			return nil
		}
		if b != nil && b.ownerRelOID != 0 && !c.snapshot {
			return &Error{
				Code:    CodeUndefinedParameter,
				Message: fmt.Sprintf("there is no parameter $%d", v.Number),
			}
		}
		return nil
	case *nodes.CurrentOfExpr:
		if v != nil && v.CursorParam > 0 && b != nil && b.ownerRelOID != 0 && !c.snapshot {
			return &Error{
				Code:    CodeUndefinedParameter,
				Message: fmt.Sprintf("there is no parameter $%d", v.CursorParam),
			}
		}
		return nil
	}
	return nil
}

func (c *Catalog) bindRuleJsonOutput(b *ruleBindings, out *nodes.JsonOutput) error {
	if out == nil || out.TypeName == nil {
		return nil
	}
	oid, _, err := c.resolveTypeName(out.TypeName)
	if err != nil {
		if !c.snapshot {
			return err
		}
		return nil
	}
	if oid != 0 {
		isArr := out.TypeName.ArrayBounds != nil && len(out.TypeName.ArrayBounds.Items) > 0
		elemOID := oid
		if bt := c.typeByOID[oid]; bt != nil && bt.Elem != 0 {
			elemOID = bt.Elem
		}
		b.typeByTN[out.TypeName] = ruleTypeBinding{typeOID: elemOID, isArray: isArr}
		b.addTypeDep(elemOID)
	}
	return nil
}

func (c *Catalog) bindRuleTableFuncColDefs(b *ruleBindings, list *nodes.List) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		cd, ok := item.(*nodes.ColumnDef)
		if !ok || cd == nil {
			continue
		}
		if cd.TypeName != nil {
			oid, _, err := c.resolveTypeName(cd.TypeName)
			if err != nil {
				if !c.snapshot {
					return err
				}
			} else if oid != 0 {
				isArr := cd.TypeName.ArrayBounds != nil && len(cd.TypeName.ArrayBounds.Items) > 0
				elemOID := oid
				if bt := c.typeByOID[oid]; bt != nil && bt.Elem != 0 {
					elemOID = bt.Elem
				}
				b.typeByTN[cd.TypeName] = ruleTypeBinding{typeOID: elemOID, isArray: isArr}
				b.addTypeDep(elemOID)
			}
		}
		if cd.CollClause != nil {
			_, uid, err := c.resolveCollationList(cd.CollClause.Collname)
			if err != nil {
				if !c.snapshot {
					return err
				}
			} else if uid != 0 {
				b.collByCC[cd.CollClause] = uid
				b.addCollDep(uid)
			}
		}
	}
	return nil
}

func (c *Catalog) bindRuleJsonAggConstructor(b *ruleBindings, ac *nodes.JsonAggConstructor, ctx *ruleEvalContext) error {
	if ac == nil {
		return nil
	}
	if err := c.bindRuleSortClause(b, ac.Agg_order, ctx); err != nil {
		return err
	}
	if err := c.bindRuleBoolExpr(b, ac.Agg_filter, ctx, "FILTER"); err != nil {
		return err
	}
	if err := c.bindWindowDef(b, ac.Over, ctx); err != nil {
		return err
	}
	return c.bindRuleJsonOutput(b, ac.Output)
}

func (c *Catalog) extractRuleQueryColTypes(b *ruleBindings, n nodes.Node) ([]uint32, bool) {
	if n == nil {
		return nil, false
	}
	switch s := n.(type) {
	case *nodes.SelectStmt:
		if s.Op != nodes.SETOP_NONE {
			lTypes, lAny := c.extractRuleQueryColTypes(b, s.Larg)
			rTypes, rAny := c.extractRuleQueryColTypes(b, s.Rarg)
			if lAny {
				return nil, true
			}
			if !rAny && len(lTypes) == len(rTypes) {
				out := make([]uint32, len(lTypes))
				for i := range lTypes {
					if lTypes[i] != 0 && lTypes[i] != UNKNOWNOID {
						out[i] = lTypes[i]
					} else {
						out[i] = rTypes[i]
					}
				}
				return out, false
			}
			return lTypes, false
		}
		if s.ValuesLists != nil {
			if len(s.ValuesLists.Items) > 0 {
				if row, ok := s.ValuesLists.Items[0].(*nodes.List); ok {
					types := make([]uint32, len(row.Items))
					for _, rowNode := range s.ValuesLists.Items {
						rList, ok := rowNode.(*nodes.List)
						if !ok {
							continue
						}
						for i, cell := range rList.Items {
							if i < len(types) && (types[i] == 0 || types[i] == UNKNOWNOID) {
								types[i] = c.inferRuleExprType(b, cell, nil)
							}
						}
					}
					return types, false
				}
			}
			return nil, true
		}
		if s.TargetList == nil {
			return nil, false
		}
		var ctx *ruleEvalContext
		if b != nil {
			ctx = b.stmtCtx[s]
		}
		var types []uint32
		for _, item := range s.TargetList.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok || rt == nil {
				continue
			}
			if b != nil && b.starExpandByRT != nil {
				if starCols, hasStar := b.starExpandByRT[rt]; hasStar && len(starCols) > 0 {
					for _, sc := range starCols {
						types = append(types, sc.typeOID)
					}
					continue
				}
			}
			if cr, isCR := rt.Val.(*nodes.ColumnRef); isCR && isStarColumnRef(cr) {
				return nil, true
			}
			types = append(types, c.inferRuleExprType(b, rt.Val, ctx))
		}
		return types, false
	}
	return nil, true
}

func (c *Catalog) inferRuleExprType(b *ruleBindings, n nodes.Node, ctx *ruleEvalContext) uint32 {
	if n == nil {
		return 0
	}
	if b != nil && b.exprTypeByNode != nil {
		if t, ok := b.exprTypeByNode[n]; ok && t != 0 {
			return t
		}
	}
	switch v := n.(type) {
	case *nodes.NamedArgExpr:
		return c.inferRuleExprType(b, v.Arg, ctx)
	case *nodes.CollateClause:
		return c.inferRuleExprType(b, v.Arg, ctx)
	case *nodes.TypeCast:
		if v.TypeName != nil {
			if oid, _, err := c.resolveTypeName(v.TypeName); err == nil {
				return oid
			}
		}
		return 0
	case *nodes.A_Const:
		if v.Isnull {
			return UNKNOWNOID
		}
		return inferRuleRawConstType(v.Val)
	case *nodes.Integer, *nodes.Float, *nodes.String, *nodes.Boolean, *nodes.BitString:
		return inferRuleRawConstType(v)
	case *nodes.BoolExpr, *nodes.NullTest, *nodes.BooleanTest, *nodes.JsonIsPredicate:
		return BOOLOID
	case *nodes.A_Expr:
		switch v.Kind {
		case nodes.AEXPR_DISTINCT, nodes.AEXPR_NOT_DISTINCT, nodes.AEXPR_IN, nodes.AEXPR_LIKE, nodes.AEXPR_ILIKE, nodes.AEXPR_SIMILAR, nodes.AEXPR_BETWEEN, nodes.AEXPR_NOT_BETWEEN, nodes.AEXPR_BETWEEN_SYM, nodes.AEXPR_NOT_BETWEEN_SYM, nodes.AEXPR_OP_ANY, nodes.AEXPR_OP_ALL:
			return BOOLOID
		case nodes.AEXPR_NULLIF:
			return c.inferRuleExprType(b, v.Lexpr, ctx)
		}
		return 0
	case *nodes.CoalesceExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if t := c.inferRuleExprType(b, item, ctx); t != 0 && t != UNKNOWNOID {
					return t
				}
			}
			if len(v.Args.Items) > 0 {
				return TEXTOID
			}
		}
		return 0
	case *nodes.MinMaxExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if t := c.inferRuleExprType(b, item, ctx); t != 0 && t != UNKNOWNOID {
					return t
				}
			}
			if len(v.Args.Items) > 0 {
				return TEXTOID
			}
		}
		return 0
	case *nodes.CaseExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if cw, ok := item.(*nodes.CaseWhen); ok && cw != nil {
					if t := c.inferRuleExprType(b, cw.Result, ctx); t != 0 && t != UNKNOWNOID {
						return t
					}
				}
			}
		}
		if t := c.inferRuleExprType(b, v.Defresult, ctx); t != 0 && t != UNKNOWNOID {
			return t
		}
		return TEXTOID
	case *nodes.SQLValueFunction:
		switch v.Op {
		case nodes.SVFOP_CURRENT_DATE:
			return DATEOID
		case nodes.SVFOP_CURRENT_TIME, nodes.SVFOP_CURRENT_TIME_N:
			return TIMETZOID
		case nodes.SVFOP_CURRENT_TIMESTAMP, nodes.SVFOP_CURRENT_TIMESTAMP_N:
			return TIMESTAMPTZOID
		case nodes.SVFOP_LOCALTIME, nodes.SVFOP_LOCALTIME_N:
			return TIMEOID
		case nodes.SVFOP_LOCALTIMESTAMP, nodes.SVFOP_LOCALTIMESTAMP_N:
			return TIMESTAMPOID
		default:
			return NAMEOID
		}
	case *nodes.SubLink:
		switch nodes.SubLinkType(v.SubLinkType) {
		case nodes.EXISTS_SUBLINK, nodes.ANY_SUBLINK, nodes.ALL_SUBLINK:
			return BOOLOID
		case nodes.ARRAY_SUBLINK:
			if cols, any := c.extractRuleQueryColTypes(b, v.Subselect); !any && len(cols) == 1 && cols[0] != 0 {
				if bt := c.typeByOID[cols[0]]; bt != nil && bt.Array != 0 {
					return bt.Array
				}
			}
			return 0
		default:
			if cols, any := c.extractRuleQueryColTypes(b, v.Subselect); !any && len(cols) == 1 {
				return cols[0]
			}
			return 0
		}
	case *nodes.A_ArrayExpr:
		if v.Elements != nil {
			for _, item := range v.Elements.Items {
				if t := c.inferRuleExprType(b, item, ctx); t != 0 && t != UNKNOWNOID {
					if bt := c.typeByOID[t]; bt != nil && bt.Array != 0 {
						return bt.Array
					}
				}
			}
		}
		return TEXTARRAYOID
	case *nodes.ArrayExpr:
		if v.Elements != nil {
			for _, item := range v.Elements.Items {
				if t := c.inferRuleExprType(b, item, ctx); t != 0 && t != UNKNOWNOID {
					if bt := c.typeByOID[t]; bt != nil && bt.Array != 0 {
						return bt.Array
					}
				}
			}
		}
		return TEXTARRAYOID
	case *nodes.RowExpr:
		return RECORDOID
	case *nodes.GroupingFunc:
		return INT4OID
	case *nodes.XmlExpr:
		if v.Op == nodes.IS_DOCUMENT {
			return BOOLOID
		}
		if v.Op == nodes.IS_XMLSERIALIZE && v.Type != 0 {
			return uint32(v.Type)
		}
		return XMLOID
	case *nodes.XmlSerialize:
		if v.TypeName != nil {
			if oid, _, err := c.resolveTypeName(v.TypeName); err == nil {
				return oid
			}
		}
		return TEXTOID
	}
	return 0
}

// ruleIntegerConstType types an integer literal as make_const does: int4
// when it fits, otherwise int8.
func ruleIntegerConstType(v *nodes.Integer) uint32 {
	if v.Ival < math.MinInt32 || v.Ival > math.MaxInt32 {
		return INT8OID
	}
	return INT4OID
}

func inferRuleRawConstType(n nodes.Node) uint32 {
	switch v := n.(type) {
	case *nodes.Integer:
		return ruleIntegerConstType(v)
	case *nodes.Float:
		return NUMERICOID
	case *nodes.Boolean:
		return BOOLOID
	case *nodes.String:
		return UNKNOWNOID
	case *nodes.BitString:
		return BITOID
	default:
		return UNKNOWNOID
	}
}

func (ctx *ruleEvalContext) colAliasesFor(relOID uint32) []string {
	if ctx == nil {
		return nil
	}
	for k, se := range ctx.scope {
		if k != "excluded" && se.relOID == relOID && len(se.colAliases) > 0 {
			return se.colAliases
		}
	}
	return nil
}

func procMatchesRuleCallArity(p *BuiltinProc, fc *nodes.FuncCall) ([]uint32, bool) {
	if p == nil {
		return nil, false
	}
	if fc != nil && fc.AggStar {
		return nil, p.NArgs == 0
	}
	var argNames []string
	argCount := 0
	if fc != nil && fc.Args != nil {
		argCount = len(fc.Args.Items)
		for i, a := range fc.Args.Items {
			if na, ok := a.(*nodes.NamedArgExpr); ok && na != nil {
				if argNames == nil {
					argNames = make([]string, i, argCount)
				}
				argNames = append(argNames, na.Name)
			} else if argNames != nil {
				argNames = append(argNames, "")
			}
		}
	}
	if fc != nil && fc.AggWithinGroup && fc.AggOrder != nil {
		argCount += len(fc.AggOrder.Items)
		if argNames != nil {
			for range fc.AggOrder.Items {
				argNames = append(argNames, "")
			}
		}
	}
	variadic := fc != nil && fc.FuncVariadic
	if argNames != nil {
		positions, ok := namedCallPositions(p, argNames, variadic)
		if !ok {
			return nil, false
		}
		types := make([]uint32, len(positions))
		for i, pos := range positions {
			if pos >= 0 && pos < len(p.ArgTypes) {
				types[i] = p.ArgTypes[pos]
			}
		}
		return types, true
	}
	return candidateArgTypesForCall(p, argCount, variadic)
}

func (c *Catalog) bindRuleFuncCall(b *ruleBindings, v *nodes.FuncCall, argCount int, ctx *ruleEvalContext) error {
	fSchema, fName := qualifiedName(v.Funcname)
	if fName == "" {
		return nil
	}
	var searchNS []uint32
	var typeSearchNS []uint32
	if fSchema != "" {
		ns := c.schemaByName[fSchema]
		if ns == nil {
			if c.snapshot {
				return nil
			}
			return errUndefinedSchema(fSchema)
		}
		searchNS = []uint32{ns.OID}
		typeSearchNS = []uint32{ns.OID}
	} else {
		searchNS = c.routineSearchPathWithCatalog()
		typeSearchNS = c.searchPathWithCatalog()
	}

	var argTypes []uint32
	var argNames []string
	allKnown := true
	if !v.AggStar && v.Args != nil {
		argTypes = make([]uint32, len(v.Args.Items))
		for i, arg := range v.Args.Items {
			if na, ok := arg.(*nodes.NamedArgExpr); ok && na != nil {
				if argNames == nil {
					argNames = make([]string, i, len(v.Args.Items))
				}
				argNames = append(argNames, na.Name)
				argTypes[i] = c.inferRuleExprType(b, na.Arg, ctx)
			} else {
				if argNames != nil {
					argNames = append(argNames, "")
				}
				argTypes[i] = c.inferRuleExprType(b, arg, ctx)
			}
			if argTypes[i] == 0 {
				allKnown = false
			}
		}
	}
	if v.AggWithinGroup && v.AggOrder != nil {
		// An ordered-set aggregate takes its WITHIN GROUP ordering
		// expressions as its trailing arguments (parse_func.c).
		for _, item := range v.AggOrder.Items {
			sb, ok := item.(*nodes.SortBy)
			if !ok || sb == nil {
				continue
			}
			t := c.inferRuleExprType(b, sb.Node, ctx)
			if t == 0 {
				allKnown = false
			}
			argTypes = append(argTypes, t)
			if argNames != nil {
				argNames = append(argNames, "")
			}
		}
	}

	procs := c.procByName[fName]
	foundAny := false
	var matchedProc *BuiltinProc
	var matchedArgTypes []uint32
	for _, nsOID := range searchNS {
		var nsProcs []*BuiltinProc
		for _, p := range procs {
			if p.Kind == 'p' {
				continue
			}
			pNS := uint32(PGCatalogNamespace)
			if p.OID >= FirstNormalObjectId {
				up := c.userProcs[p.OID]
				if up == nil || up.Schema == nil {
					continue
				}
				pNS = up.Schema.OID
			}
			if pNS != nsOID {
				continue
			}
			nsProcs = append(nsProcs, p)
		}
		if len(nsProcs) == 0 {
			continue
		}
		if allKnown {
			bestProc, resolvedArgs, err := c.resolveFuncOverload(nsProcs, argTypes, argNames, v.AggStar, v.FuncVariadic)
			if err == nil && bestProc != nil {
				foundAny = true
				matchedProc = bestProc
				matchedArgTypes = resolvedArgs
				break
			}
			continue
		}
		// Partial argument types known (e.g. open virtual source).
		var partialMatches []*BuiltinProc
		var firstArgTypes []uint32
		for _, p := range nsProcs {
			candTypes, ok := procMatchesRuleCallArity(p, v)
			if !ok {
				continue
			}
			compatible := true
			for i, at := range argTypes {
				if at != 0 && at != UNKNOWNOID && i < len(candTypes) {
					if !c.canCoerceArgs([]uint32{at}, []uint32{candTypes[i]}) {
						compatible = false
						break
					}
				}
			}
			if !compatible {
				continue
			}
			partialMatches = append(partialMatches, p)
			if firstArgTypes == nil {
				firstArgTypes = candTypes
			}
		}
		if len(partialMatches) > 0 {
			foundAny = true
			matchedArgTypes = firstArgTypes
			if len(partialMatches) == 1 {
				matchedProc = partialMatches[0]
			}
			break
		}
	}
	if foundAny && v.Args != nil {
		for i, arg := range v.Args.Items {
			if i < len(matchedArgTypes) && matchedArgTypes[i] == REGCLASSOID {
				if err := c.bindRegclassLiteral(b, arg); err != nil {
					return err
				}
			}
		}
	}
	if matchedProc != nil {
		if matchedProc.RetSet && b != nil && b.srfCalls != nil {
			b.srfCalls[v] = true
		}
		if matchedProc.RetSet && !c.snapshot {
			where := ""
			if b != nil && b.inRuleWhere {
				where = "WHERE"
			} else if ctx != nil {
				where = ctx.disallowSRF
			}
			if where != "" {
				return &Error{Code: CodeFeatureNotSupported, Message: fmt.Sprintf("set-returning functions are not allowed in %s", where)}
			}
		}
		if !c.snapshot {
			disallowAgg := ""
			disallowWin := ""
			if b != nil && b.inRuleWhere {
				disallowAgg = "WHERE"
				disallowWin = "WHERE"
			} else if ctx != nil {
				disallowAgg = ctx.disallowAgg
				disallowWin = ctx.disallowWin
			}
			if v.Over != nil {
				if disallowWin != "" {
					return &Error{
						Code:    CodeWindowingError,
						Message: fmt.Sprintf("window functions are not allowed in %s", disallowWin),
					}
				}
				if matchedProc.Kind != 'a' && matchedProc.Kind != 'w' {
					return &Error{
						Code:    CodeWrongObjectType,
						Message: fmt.Sprintf("OVER specified, but %s is not a window function nor an aggregate function", fName),
					}
				}
			} else {
				if matchedProc.Kind == 'w' {
					return &Error{
						Code:    CodeWrongObjectType,
						Message: fmt.Sprintf("window function %s requires an OVER clause", fName),
					}
				}
				hasAggSyntax := v.AggStar || v.AggDistinct || v.AggFilter != nil || (v.AggOrder != nil && len(v.AggOrder.Items) > 0)
				if hasAggSyntax && matchedProc.Kind != 'a' {
					return &Error{
						Code:    CodeWrongObjectType,
						Message: fmt.Sprintf("%s is not an aggregate function", fName),
					}
				}
				if matchedProc.Kind == 'a' && disallowAgg != "" {
					return &Error{
						Code:    CodeGroupingError,
						Message: fmt.Sprintf("aggregate functions are not allowed in %s", disallowAgg),
					}
				}
			}
		}
		// A polymorphic result follows the actual argument types, not the
		// declared ones the overload matched.
		retArgs := matchedArgTypes
		if isPolymorphic(matchedProc.RetType) && argNames == nil {
			retArgs = slices.Clone(matchedArgTypes)
			for i := range retArgs {
				if i < len(argTypes) && argTypes[i] != 0 && argTypes[i] != UNKNOWNOID {
					retArgs[i] = argTypes[i]
				}
			}
		}
		retType := c.resolveReturnType(matchedProc, c.withDefaultArgTypes(matchedProc, retArgs))
		if b != nil && b.exprTypeByNode != nil && retType != 0 {
			b.exprTypeByNode[v] = retType
		}
		if matchedProc.OID >= FirstNormalObjectId {
			b.funcByFC[v] = matchedProc.OID
			b.addFuncDep(matchedProc.OID)
		}
		return nil
	}
	if foundAny {
		return nil
	}

	// PostgreSQL allows single-argument functional cast syntax: typename(expr).
	if argCount == 1 && !v.AggStar && !v.AggDistinct && v.Over == nil && argNames == nil && !v.FuncVariadic {
		for _, nsOID := range typeSearchNS {
			if bt := c.typeByName[typeKey{ns: nsOID, name: fName}]; bt != nil && bt.IsDefined {
				if len(argTypes) == 1 && argTypes[0] != 0 && argTypes[0] != UNKNOWNOID && !c.CanCoerce(argTypes[0], bt.OID, 'e') {
					break
				}
				elemOID := bt.OID
				if bt.Elem != 0 {
					elemOID = bt.Elem
				}
				if b.castTypeByFC != nil {
					b.castTypeByFC[v] = elemOID
				}
				if b.exprTypeByNode != nil {
					b.exprTypeByNode[v] = bt.OID
				}
				b.addTypeDep(elemOID)
				if bt.OID == REGCLASSOID && v.Args != nil && len(v.Args.Items) == 1 {
					if err := c.bindRegclassLiteral(b, v.Args.Items[0]); err != nil {
						return err
					}
				}
				return nil
			}
		}
	}

	if !c.snapshot {
		shown := fName
		if fSchema != "" {
			shown = fSchema + "." + fName
		}
		return errUndefinedFunction(shown, argTypes)
	}
	return nil
}

func (c *Catalog) bindRegclassLiteral(b *ruleBindings, arg nodes.Node) error {
	if arg == nil {
		return nil
	}
	if na, ok := arg.(*nodes.NamedArgExpr); ok && na != nil {
		arg = na.Arg
	}
	for {
		if tc, ok := arg.(*nodes.TypeCast); ok && tc != nil {
			arg = tc.Arg
			continue
		}
		break
	}
	var litNode nodes.Node
	var raw string
	switch v := arg.(type) {
	case *nodes.A_Const:
		if v == nil || v.Isnull {
			return nil
		}
		s, ok := v.Val.(*nodes.String)
		if !ok || s == nil {
			return nil
		}
		litNode = v
		raw = s.Str
	case *nodes.String:
		if v == nil {
			return nil
		}
		litNode = v
		raw = v.Str
	default:
		return nil
	}
	schemaName, objName := parseRegclassIdent(raw)
	if objName == "" {
		return nil
	}
	if seq, err := c.findSequence(schemaName, objName); err == nil && seq != nil {
		if b != nil && b.regclassByNode != nil {
			b.regclassByNode[litNode] = ruleRegclassBinding{seqOID: seq.OID}
		}
		if b != nil {
			b.addSeqDep(seq.OID)
		}
		return nil
	}
	if _, rel, err := c.findRelation(schemaName, objName); err == nil && rel != nil {
		if b != nil && b.regclassByNode != nil {
			b.regclassByNode[litNode] = ruleRegclassBinding{relOID: rel.OID}
		}
		if b != nil {
			b.addRelDep(rel.OID)
		}
		return nil
	}
	if idx := c.findRegclassIndex(schemaName, objName); idx != nil {
		if b != nil && b.regclassByNode != nil {
			b.regclassByNode[litNode] = ruleRegclassBinding{idxOID: idx.OID}
		}
		return nil
	}
	if !c.snapshot {
		return errUndefinedTable(raw)
	}
	return nil
}

func (c *Catalog) findRegclassIndex(schemaName, idxName string) *Index {
	if schemaName != "" {
		if s := c.schemaByName[schemaName]; s != nil {
			return s.Indexes[idxName]
		}
		return nil
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		if s := c.schemas[nsOID]; s != nil {
			if idx := s.Indexes[idxName]; idx != nil {
				return idx
			}
		}
	}
	return nil
}

func parseRegclassIdent(raw string) (string, string) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", ""
	}
	var parts []string
	for i := 0; i < len(s); {
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
		if i >= len(s) {
			break
		}
		if s[i] == '"' {
			i++
			var b strings.Builder
			for i < len(s) {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						b.WriteByte('"')
						i += 2
						continue
					}
					i++
					break
				}
				b.WriteByte(s[i])
				i++
			}
			parts = append(parts, b.String())
			for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
				i++
			}
			if i < len(s) && s[i] == '.' {
				i++
			}
		} else {
			start := i
			for i < len(s) && s[i] != '.' && s[i] != '"' {
				i++
			}
			part := strings.ToLower(strings.TrimSpace(s[start:i]))
			parts = append(parts, part)
			if i < len(s) && s[i] == '.' {
				i++
			}
		}
	}
	if len(parts) == 1 {
		return "", parts[0]
	}
	if len(parts) >= 2 {
		return parts[len(parts)-2], parts[len(parts)-1]
	}
	return "", ""
}

func inferRuleTargetColName(n nodes.Node) string {
	if cr, ok := n.(*nodes.ColumnRef); ok && cr != nil && cr.Fields != nil && len(cr.Fields.Items) > 0 {
		last := cr.Fields.Items[len(cr.Fields.Items)-1]
		if _, isStar := last.(*nodes.A_Star); isStar {
			return ""
		}
		return stringVal(last)
	}
	return figureColName(n)
}

type ruleDeparser struct {
	cat *Catalog
	b   *ruleBindings
}

func (c *Catalog) deparseRuleWhere(r *Rule) string {
	if r == nil || r.WhereClause == nil {
		return ""
	}
	d := ruleDeparser{cat: c, b: r.bindings}
	return d.deparseExpr(r.WhereClause)
}

func (c *Catalog) deparseRuleActions(r *Rule) []string {
	if r == nil || r.Actions == nil {
		return nil
	}
	d := ruleDeparser{cat: c, b: r.bindings}
	var out []string
	for _, item := range r.Actions.Items {
		if act := d.deparseActionStmt(item); act != "" {
			out = append(out, act)
		}
	}
	return out
}

func ruleEventKeyword(ev nodes.CmdType) string {
	switch ev {
	case nodes.CMD_SELECT:
		return "SELECT"
	case nodes.CMD_UPDATE:
		return "UPDATE"
	case nodes.CMD_DELETE:
		return "DELETE"
	case nodes.CMD_INSERT:
		return "INSERT"
	default:
		return "SELECT"
	}
}

func (d *ruleDeparser) resolveBoundColName(cb ruleColBinding, fallback string) string {
	if d.cat != nil && cb.relOID != 0 && cb.attNum > 0 {
		if rel := d.cat.relationByOID[cb.relOID]; rel != nil {
			for _, col := range rel.Columns {
				if col.AttNum == cb.attNum {
					return col.Name
				}
			}
		}
	}
	return fallback
}

func (d *ruleDeparser) resolveJoinUsingColName(uc ruleJoinUsingCol) string {
	for _, dep := range uc.deps {
		if name := d.resolveBoundColName(ruleColBinding{relOID: dep.relOID, attNum: dep.attNum}, ""); name != "" {
			return name
		}
	}
	return uc.name
}

func (d *ruleDeparser) deparseBoundRegclassLiteral(n nodes.Node) string {
	if d.cat == nil || d.b == nil || d.b.regclassByNode == nil || n == nil {
		return ""
	}
	rb, ok := d.b.regclassByNode[n]
	if !ok {
		return ""
	}
	if rb.seqOID != 0 {
		for _, s := range d.cat.schemas {
			if s == nil {
				continue
			}
			for _, seq := range s.Sequences {
				if seq != nil && seq.OID == rb.seqOID {
					if s.Name != "public" && s.Name != "pg_catalog" {
						return quoteRuleLiteral(quoteIdentifier(s.Name) + "." + quoteIdentifier(seq.Name))
					}
					return quoteRuleLiteral(quoteIdentifier(seq.Name))
				}
			}
		}
	}
	if rb.relOID != 0 {
		if rel := d.cat.relationByOID[rb.relOID]; rel != nil {
			if rel.Schema != nil && rel.Schema.Name != "public" && rel.Schema.Name != "pg_catalog" {
				return quoteRuleLiteral(quoteIdentifier(rel.Schema.Name) + "." + quoteIdentifier(rel.Name))
			}
			return quoteRuleLiteral(quoteIdentifier(rel.Name))
		}
	}
	if rb.idxOID != 0 {
		if idx := d.cat.indexes[rb.idxOID]; idx != nil {
			if idx.Schema != nil && idx.Schema.Name != "public" && idx.Schema.Name != "pg_catalog" {
				return quoteRuleLiteral(quoteIdentifier(idx.Schema.Name) + "." + quoteIdentifier(idx.Name))
			}
			return quoteRuleLiteral(quoteIdentifier(idx.Name))
		}
	}
	return ""
}

func (d *ruleDeparser) deparseIndexElem(elem *nodes.IndexElem) string {
	if elem == nil {
		return ""
	}
	var b strings.Builder
	if elem.Name != "" {
		colName := elem.Name
		if d.b != nil {
			if cb, ok := d.b.colByIE[elem]; ok {
				colName = d.resolveBoundColName(cb, colName)
			}
		}
		b.WriteString(quoteIdentifier(colName))
	} else if elem.Expr != nil {
		b.WriteString("(")
		b.WriteString(d.deparseExpr(elem.Expr))
		b.WriteString(")")
	}
	if elem.Collation != nil && len(elem.Collation.Items) > 0 {
		collSQL := deparseQualifiedIdentList(elem.Collation)
		if d.cat != nil && d.b != nil {
			if uid, ok := d.b.collByIE[elem]; ok && uid != 0 {
				for _, s := range d.cat.schemas {
					if s == nil {
						continue
					}
					for _, uc := range s.Collations {
						if uc != nil && uc.OID == uid {
							if s.Name != "public" && s.Name != "pg_catalog" {
								collSQL = quoteIdentifier(s.Name) + "." + quoteIdentifier(uc.Name)
							} else {
								collSQL = quoteIdentifier(uc.Name)
							}
							break
						}
					}
				}
			}
		}
		b.WriteString(" COLLATE ")
		b.WriteString(collSQL)
	}
	if elem.Opclass != nil && len(elem.Opclass.Items) > 0 {
		opcSQL := deparseFuncNameList(elem.Opclass)
		if d.cat != nil && d.b != nil {
			if opcOID, ok := d.b.opcByIE[elem]; ok && opcOID != 0 {
				if opc := d.cat.opClasses[opcOID]; opc != nil {
					if sc := d.cat.schemas[opc.Schema]; sc != nil && (len(elem.Opclass.Items) > 1 || (sc.Name != "public" && sc.Name != "pg_catalog")) {
						opcSQL = quoteIdentifier(sc.Name) + "." + quoteIdentifier(opc.Name)
					} else {
						opcSQL = quoteIdentifier(opc.Name)
					}
				}
			}
		}
		b.WriteString(" ")
		b.WriteString(opcSQL)
		if elem.Opclassopts != nil && len(elem.Opclassopts.Items) > 0 {
			b.WriteString(" (")
			var optParts []string
			for _, item := range elem.Opclassopts.Items {
				if de, ok := item.(*nodes.DefElem); ok && de != nil {
					if de.Arg != nil {
						optParts = append(optParts, quoteIdentifier(de.Defname)+" = "+d.deparseDefElemArg(de.Arg))
					} else {
						optParts = append(optParts, quoteIdentifier(de.Defname))
					}
				}
			}
			b.WriteString(strings.Join(optParts, ", "))
			b.WriteString(")")
		}
	}
	switch elem.Ordering {
	case nodes.SORTBY_ASC:
		b.WriteString(" ASC")
	case nodes.SORTBY_DESC:
		b.WriteString(" DESC")
	}
	switch elem.NullsOrdering {
	case nodes.SORTBY_NULLS_FIRST:
		b.WriteString(" NULLS FIRST")
	case nodes.SORTBY_NULLS_LAST:
		b.WriteString(" NULLS LAST")
	}
	return b.String()
}

func (d *ruleDeparser) deparseDefElemArg(n nodes.Node) string {
	switch v := n.(type) {
	case *nodes.String:
		return v.Str
	case *nodes.Integer, *nodes.Float, *nodes.Boolean, *nodes.BitString:
		return deparseRawConstVal(v)
	case *nodes.TypeName:
		return d.deparseTypeName(v)
	default:
		return d.deparseExpr(n)
	}
}

func (c *Catalog) inferPreparedDMLParamTypes(b *ruleBindings, n nodes.Node, paramTypes []uint32, fixedParams ...bool) error {
	if len(paramTypes) == 0 || n == nil {
		return nil
	}
	isFixed := len(fixedParams) > 0 && fixedParams[0]
	var firstErr error
	setParam := func(num int, typ uint32) {
		if firstErr != nil || num <= 0 || num > len(paramTypes) || typ == 0 || typ == UNKNOWNOID {
			return
		}
		if paramTypes[num-1] == 0 || paramTypes[num-1] == UNKNOWNOID {
			paramTypes[num-1] = typ
		} else if !isFixed && paramTypes[num-1] != typ {
			firstErr = &Error{
				Code:    CodeAmbiguousParameter,
				Message: fmt.Sprintf("could not determine data type of parameter $%d", num),
			}
		}
	}
	resolveExprType := func(expr nodes.Node) uint32 {
		switch ev := expr.(type) {
		case *nodes.ColumnRef:
			if b != nil {
				if cb, ok := b.colByRef[ev]; ok && cb.relOID != 0 {
					if rel := c.relationByOID[cb.relOID]; rel != nil {
						for _, col := range rel.Columns {
							if col.AttNum == cb.attNum {
								return col.TypeOID
							}
						}
					}
				}
			}
		case *nodes.TypeCast:
			if ev.TypeName != nil {
				if oid, _, err := c.resolveTypeName(ev.TypeName); err == nil {
					return oid
				}
			}
		case *nodes.A_Const:
			if !ev.Isnull {
				switch cv := ev.Val.(type) {
				case *nodes.Integer:
					return ruleIntegerConstType(cv)
				case *nodes.Float:
					return NUMERICOID
				case *nodes.Boolean:
					return BOOLOID
				case *nodes.String:
					return TEXTOID
				}
			}
		}
		if t := c.inferRuleExprType(b, expr, nil); t != 0 && t != UNKNOWNOID {
			return t
		}
		return 0
	}
	var walkExpr func(expr nodes.Node, expectedType uint32)
	walkExpr = func(expr nodes.Node, expectedType uint32) {
		if expr == nil {
			return
		}
		switch ev := expr.(type) {
		case *nodes.ParamRef:
			setParam(ev.Number, expectedType)
		case *nodes.TypeCast:
			var castOID uint32
			if ev.TypeName != nil {
				castOID, _, _ = c.resolveTypeName(ev.TypeName)
			}
			walkExpr(ev.Arg, castOID)
		case *nodes.CollateClause:
			walkExpr(ev.Arg, expectedType)
		case *nodes.A_Expr:
			lType := resolveExprType(ev.Lexpr)
			rType := resolveExprType(ev.Rexpr)
			walkExpr(ev.Lexpr, rType)
			if list, ok := ev.Rexpr.(*nodes.List); ok {
				for _, item := range list.Items {
					walkExpr(item, lType)
				}
			} else {
				walkExpr(ev.Rexpr, lType)
			}
		case *nodes.BoolExpr:
			if ev.Args != nil {
				for _, item := range ev.Args.Items {
					walkExpr(item, BOOLOID)
				}
			}
		case *nodes.BooleanTest:
			walkExpr(ev.Arg, BOOLOID)
		case *nodes.NullTest:
			walkExpr(ev.Arg, 0)
		case *nodes.CoalesceExpr:
			if ev.Args != nil {
				var common uint32
				for _, item := range ev.Args.Items {
					if t := resolveExprType(item); t != 0 {
						common = t
						break
					}
				}
				if common == 0 {
					common = expectedType
				}
				for _, item := range ev.Args.Items {
					walkExpr(item, common)
				}
			}
		case *nodes.MinMaxExpr:
			if ev.Args != nil {
				var common uint32
				for _, item := range ev.Args.Items {
					if t := resolveExprType(item); t != 0 {
						common = t
						break
					}
				}
				if common == 0 {
					common = expectedType
				}
				for _, item := range ev.Args.Items {
					walkExpr(item, common)
				}
			}
		case *nodes.CaseExpr:
			argType := resolveExprType(ev.Arg)
			walkExpr(ev.Arg, 0)
			var resCommon uint32
			if ev.Args != nil {
				for _, item := range ev.Args.Items {
					if cw, ok := item.(*nodes.CaseWhen); ok && cw != nil {
						if resCommon == 0 {
							resCommon = resolveExprType(cw.Result)
						}
					}
				}
			}
			if resCommon == 0 {
				resCommon = resolveExprType(ev.Defresult)
			}
			if resCommon == 0 {
				resCommon = expectedType
			}
			whenExpected := uint32(BOOLOID)
			if ev.Arg != nil {
				whenExpected = argType
			}
			if ev.Args != nil {
				for _, item := range ev.Args.Items {
					if cw, ok := item.(*nodes.CaseWhen); ok && cw != nil {
						walkExpr(cw.Expr, whenExpected)
						walkExpr(cw.Result, resCommon)
					}
				}
			}
			walkExpr(ev.Defresult, resCommon)
		case *nodes.A_ArrayExpr:
			if ev.Elements != nil {
				var elemCommon uint32
				for _, item := range ev.Elements.Items {
					if t := resolveExprType(item); t != 0 {
						elemCommon = t
						break
					}
				}
				for _, item := range ev.Elements.Items {
					walkExpr(item, elemCommon)
				}
			}
		case *nodes.ArrayExpr:
			if ev.Elements != nil {
				for _, item := range ev.Elements.Items {
					walkExpr(item, 0)
				}
			}
		case *nodes.RowExpr:
			if ev.Args != nil {
				for _, item := range ev.Args.Items {
					walkExpr(item, 0)
				}
			}
		case *nodes.SubLink:
			walkExpr(ev.Testexpr, 0)
			if sel, ok := ev.Subselect.(*nodes.SelectStmt); ok && sel != nil {
				var optParams []uint32
				if isFixed {
					optParams = paramTypes
				}
				if q, err := c.analyzeSelectStmtWithOptions(sel, AnalyzeOptions{ParamTypes: optParams}); err == nil {
					if err := collectQueryParamTypes(q, paramTypes, isFixed); err != nil && firstErr == nil {
						firstErr = err
					}
				}
				if err := c.inferPreparedDMLParamTypes(b, sel, paramTypes, isFixed); err != nil && firstErr == nil {
					firstErr = err
				}
			}
		case *nodes.FuncCall:
			if ev.Args != nil {
				_, fName := qualifiedName(ev.Funcname)
				var argTypes []uint32
				for _, p := range c.procByName[fName] {
					if int(p.NArgs) == len(ev.Args.Items) {
						argTypes = p.ArgTypes
						break
					}
				}
				for i, arg := range ev.Args.Items {
					var at uint32
					if i < len(argTypes) {
						at = argTypes[i]
					}
					walkExpr(arg, at)
				}
			}
			walkExpr(ev.AggFilter, BOOLOID)
		}
	}
	walkReturning := func(list *nodes.List) {
		if list == nil {
			return
		}
		for _, item := range list.Items {
			if rt, ok := item.(*nodes.ResTarget); ok && rt != nil {
				walkExpr(rt.Val, 0)
			} else {
				walkExpr(item, 0)
			}
		}
	}
	walkWithClause := func(wc *nodes.WithClause) {
		if wc == nil || wc.Ctes == nil {
			return
		}
		for _, item := range wc.Ctes.Items {
			if cte, ok := item.(*nodes.CommonTableExpr); ok && cte != nil && cte.Ctequery != nil {
				if sel, isSel := cte.Ctequery.(*nodes.SelectStmt); isSel && sel != nil {
					var optParams []uint32
					if isFixed {
						optParams = paramTypes
					}
					if q, err := c.analyzeSelectStmtWithOptions(sel, AnalyzeOptions{ParamTypes: optParams}); err == nil {
						if err := collectQueryParamTypes(q, paramTypes, isFixed); err != nil && firstErr == nil {
							firstErr = err
						}
					}
				}
				if err := c.inferPreparedDMLParamTypes(b, cte.Ctequery, paramTypes, isFixed); err != nil && firstErr == nil {
					firstErr = err
				}
			}
		}
	}
	var walkFromItem func(item nodes.Node)
	walkFromItem = func(item nodes.Node) {
		if item == nil {
			return
		}
		switch v := item.(type) {
		case *nodes.JoinExpr:
			walkFromItem(v.Larg)
			walkFromItem(v.Rarg)
			walkExpr(v.Quals, BOOLOID)
		case *nodes.RangeSubselect:
			if sel, ok := v.Subquery.(*nodes.SelectStmt); ok && sel != nil {
				var optParams []uint32
				if isFixed {
					optParams = paramTypes
				}
				if q, err := c.analyzeSelectStmtWithOptions(sel, AnalyzeOptions{ParamTypes: optParams}); err == nil {
					if err := collectQueryParamTypes(q, paramTypes, isFixed); err != nil && firstErr == nil {
						firstErr = err
					}
				}
			}
			if err := c.inferPreparedDMLParamTypes(b, v.Subquery, paramTypes, isFixed); err != nil && firstErr == nil {
				firstErr = err
			}
		case *nodes.RangeFunction:
			if v.Functions != nil {
				for _, fnItem := range v.Functions.Items {
					if fnList, ok := fnItem.(*nodes.List); ok && fnList != nil {
						for _, elem := range fnList.Items {
							walkExpr(elem, 0)
						}
					} else {
						walkExpr(fnItem, 0)
					}
				}
			}
		case *nodes.RangeTableSample:
			walkFromItem(v.Relation)
			if v.Args != nil {
				for _, arg := range v.Args.Items {
					walkExpr(arg, 0)
				}
			}
			walkExpr(v.Repeatable, FLOAT8OID)
		case *nodes.RangeTableFunc:
			walkExpr(v.Docexpr, 0)
			walkExpr(v.Rowexpr, 0)
			if v.Columns != nil {
				for _, colItem := range v.Columns.Items {
					if col, ok := colItem.(*nodes.RangeTableFuncCol); ok && col != nil {
						walkExpr(col.Colexpr, 0)
						walkExpr(col.Coldefexpr, 0)
					}
				}
			}
		}
	}
	walkFromList := func(list *nodes.List) {
		if list == nil {
			return
		}
		for _, item := range list.Items {
			walkFromItem(item)
		}
	}
	targetColType := func(rt *nodes.ResTarget, rel *Relation) uint32 {
		if rt == nil || rel == nil {
			return 0
		}
		if idx, ok := rel.colByName[rt.Name]; ok {
			return rel.Columns[idx].TypeOID
		}
		return 0
	}
	switch s := n.(type) {
	case *nodes.SelectStmt:
		walkWithClause(s.WithClause)
		if s.Larg != nil {
			if err := c.inferPreparedDMLParamTypes(b, s.Larg, paramTypes, isFixed); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		if s.Rarg != nil {
			if err := c.inferPreparedDMLParamTypes(b, s.Rarg, paramTypes, isFixed); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		walkFromList(s.FromClause)
		walkReturning(s.TargetList)
		walkExpr(s.WhereClause, BOOLOID)
		walkExpr(s.HavingClause, BOOLOID)
		walkExpr(s.LimitCount, INT8OID)
		walkExpr(s.LimitOffset, INT8OID)
	case *nodes.InsertStmt:
		walkWithClause(s.WithClause)
		var rel *Relation
		if b != nil && s.Relation != nil {
			rel = c.relationByOID[b.relByRV[s.Relation]]
		}
		var colTypes []uint32
		if rel != nil {
			if s.Cols != nil && len(s.Cols.Items) > 0 {
				for _, item := range s.Cols.Items {
					if rt, ok := item.(*nodes.ResTarget); ok {
						colTypes = append(colTypes, targetColType(rt, rel))
					}
				}
			} else {
				for _, col := range rel.Columns {
					colTypes = append(colTypes, col.TypeOID)
				}
			}
		}
		if sel, ok := s.SelectStmt.(*nodes.SelectStmt); ok && sel != nil {
			if sel.ValuesLists != nil {
				for _, rowNode := range sel.ValuesLists.Items {
					if row, ok := rowNode.(*nodes.List); ok {
						for i, cell := range row.Items {
							var ct uint32
							if i < len(colTypes) {
								ct = colTypes[i]
							}
							walkExpr(cell, ct)
						}
					}
				}
			} else {
				if sel.TargetList != nil {
					for i, item := range sel.TargetList.Items {
						if rt, ok := item.(*nodes.ResTarget); ok {
							var ct uint32
							if i < len(colTypes) {
								ct = colTypes[i]
							}
							walkExpr(rt.Val, ct)
						}
					}
				}
				var optParams []uint32
				if isFixed {
					optParams = paramTypes
				}
				if q, err := c.analyzeSelectStmtWithOptions(sel, AnalyzeOptions{ParamTypes: optParams}); err == nil {
					if err := collectQueryParamTypes(q, paramTypes, isFixed); err != nil && firstErr == nil {
						firstErr = err
					}
				}
				if err := c.inferPreparedDMLParamTypes(b, sel, paramTypes, isFixed); err != nil && firstErr == nil {
					firstErr = err
				}
			}
		}
		if s.OnConflictClause != nil {
			if s.OnConflictClause.Infer != nil {
				if s.OnConflictClause.Infer.IndexElems != nil {
					for _, ie := range s.OnConflictClause.Infer.IndexElems.Items {
						if elem, ok := ie.(*nodes.IndexElem); ok && elem != nil {
							walkExpr(elem.Expr, 0)
						}
					}
				}
				walkExpr(s.OnConflictClause.Infer.WhereClause, BOOLOID)
			}
			if rel != nil && s.OnConflictClause.TargetList != nil {
				for _, item := range s.OnConflictClause.TargetList.Items {
					if rt, ok := item.(*nodes.ResTarget); ok {
						walkExpr(rt.Val, targetColType(rt, rel))
					}
				}
			}
			walkExpr(s.OnConflictClause.WhereClause, BOOLOID)
		}
		walkReturning(s.ReturningList)
	case *nodes.UpdateStmt:
		walkWithClause(s.WithClause)
		var rel *Relation
		if b != nil && s.Relation != nil {
			rel = c.relationByOID[b.relByRV[s.Relation]]
		}
		if s.TargetList != nil {
			for _, item := range s.TargetList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					walkExpr(rt.Val, targetColType(rt, rel))
				}
			}
		}
		walkFromList(s.FromClause)
		walkExpr(s.WhereClause, BOOLOID)
		walkReturning(s.ReturningList)
	case *nodes.DeleteStmt:
		walkWithClause(s.WithClause)
		walkFromList(s.UsingClause)
		walkExpr(s.WhereClause, BOOLOID)
		walkReturning(s.ReturningList)
	case *nodes.MergeStmt:
		walkWithClause(s.WithClause)
		var rel *Relation
		if b != nil && s.Relation != nil {
			rel = c.relationByOID[b.relByRV[s.Relation]]
		}
		walkFromItem(s.SourceRelation)
		walkExpr(s.JoinCondition, BOOLOID)
		if s.MergeWhenClauses != nil {
			for _, item := range s.MergeWhenClauses.Items {
				mwc, ok := item.(*nodes.MergeWhenClause)
				if !ok || mwc == nil {
					continue
				}
				walkExpr(mwc.Condition, BOOLOID)
				if mwc.CommandType == nodes.CMD_INSERT && mwc.Values != nil && rel != nil {
					var colTypes []uint32
					if mwc.TargetList != nil && len(mwc.TargetList.Items) > 0 {
						for _, tlItem := range mwc.TargetList.Items {
							if rt, ok := tlItem.(*nodes.ResTarget); ok {
								colTypes = append(colTypes, targetColType(rt, rel))
							}
						}
					} else {
						for _, col := range rel.Columns {
							colTypes = append(colTypes, col.TypeOID)
						}
					}
					for i, valExpr := range mwc.Values.Items {
						var ct uint32
						if i < len(colTypes) {
							ct = colTypes[i]
						}
						walkExpr(valExpr, ct)
					}
				} else if mwc.CommandType == nodes.CMD_UPDATE && mwc.TargetList != nil && rel != nil {
					for _, tlItem := range mwc.TargetList.Items {
						if rt, ok := tlItem.(*nodes.ResTarget); ok {
							walkExpr(rt.Val, targetColType(rt, rel))
						}
					}
				}
			}
		}
		walkReturning(s.ReturningList)
	}
	return firstErr
}

func (d *ruleDeparser) deparseActionStmt(n nodes.Node) string {
	if n == nil {
		return ""
	}
	switch s := n.(type) {
	case *nodes.SelectStmt:
		return d.deparseSelectStmt(s)
	case *nodes.InsertStmt:
		return d.deparseInsertStmt(s)
	case *nodes.UpdateStmt:
		return d.deparseUpdateStmt(s)
	case *nodes.DeleteStmt:
		return d.deparseDeleteStmt(s)
	case *nodes.MergeStmt:
		return d.deparseMergeStmt(s)
	case *nodes.NotifyStmt:
		return deparseRawNotifyStmt(s)
	default:
		return ""
	}
}

func (d *ruleDeparser) deparseMergeStmt(s *nodes.MergeStmt) string {
	if s == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(d.deparseWithClause(s.WithClause))
	b.WriteString("MERGE INTO ")
	b.WriteString(d.deparseRangeVar(s.Relation, true))
	if s.SourceRelation != nil {
		b.WriteString(" USING ")
		b.WriteString(d.deparseTableRef(s.SourceRelation))
	}
	if s.JoinCondition != nil {
		b.WriteString(" ON ")
		b.WriteString(d.deparseExpr(s.JoinCondition))
	}
	if s.MergeWhenClauses != nil {
		for _, item := range s.MergeWhenClauses.Items {
			mwc, ok := item.(*nodes.MergeWhenClause)
			if !ok || mwc == nil {
				continue
			}
			switch mwc.Kind {
			case nodes.MERGE_WHEN_MATCHED:
				b.WriteString(" WHEN MATCHED")
			case nodes.MERGE_WHEN_NOT_MATCHED_BY_SOURCE:
				b.WriteString(" WHEN NOT MATCHED BY SOURCE")
			default:
				b.WriteString(" WHEN NOT MATCHED")
			}
			if mwc.Condition != nil {
				b.WriteString(" AND ")
				b.WriteString(d.deparseExpr(mwc.Condition))
			}
			b.WriteString(" THEN ")
			switch mwc.CommandType {
			case nodes.CMD_INSERT:
				b.WriteString("INSERT")
				if mwc.TargetList != nil && len(mwc.TargetList.Items) > 0 {
					b.WriteString(" (")
					for i, tlItem := range mwc.TargetList.Items {
						if i > 0 {
							b.WriteString(", ")
						}
						if rt, ok := tlItem.(*nodes.ResTarget); ok {
							b.WriteString(d.deparseUpdateTargetCol(rt))
						}
					}
					b.WriteString(")")
				}
				switch mwc.Override {
				case nodes.OVERRIDING_USER_VALUE:
					b.WriteString(" OVERRIDING USER VALUE")
				case nodes.OVERRIDING_SYSTEM_VALUE:
					b.WriteString(" OVERRIDING SYSTEM VALUE")
				}
				if mwc.Values == nil {
					b.WriteString(" DEFAULT VALUES")
				} else {
					b.WriteString(" VALUES (")
					b.WriteString(d.deparseExprList(mwc.Values))
					b.WriteString(")")
				}
			case nodes.CMD_UPDATE:
				b.WriteString("UPDATE SET ")
				b.WriteString(d.deparseUpdateTargetList(mwc.TargetList))
			case nodes.CMD_DELETE:
				b.WriteString("DELETE")
			default:
				b.WriteString("DO NOTHING")
			}
		}
	}
	if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
		b.WriteString(" RETURNING ")
		b.WriteString(d.deparseTargetList(s.ReturningList))
	}
	return b.String()
}

func deparseRawNotifyStmt(s *nodes.NotifyStmt) string {
	if s.Payload != "" {
		return fmt.Sprintf("NOTIFY %s, %s", quoteIdentifier(s.Conditionname), quoteRuleLiteral(s.Payload))
	}
	return "NOTIFY " + quoteIdentifier(s.Conditionname)
}

func (d *ruleDeparser) deparseWithClause(w *nodes.WithClause) string {
	if w == nil || w.Ctes == nil || len(w.Ctes.Items) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("WITH ")
	if w.Recursive {
		b.WriteString("RECURSIVE ")
	}
	for i, item := range w.Ctes.Items {
		cte, ok := item.(*nodes.CommonTableExpr)
		if !ok {
			continue
		}
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(quoteIdentifier(cte.Ctename))
		if cte.Aliascolnames != nil && len(cte.Aliascolnames.Items) > 0 {
			b.WriteString(" (")
			for j, col := range cte.Aliascolnames.Items {
				if j > 0 {
					b.WriteString(", ")
				}
				b.WriteString(quoteIdentifier(stringVal(col)))
			}
			b.WriteString(")")
		}
		b.WriteString(" AS ")
		switch nodes.CTEMaterialize(cte.Ctematerialized) {
		case nodes.CTEMaterializeAlways:
			b.WriteString("MATERIALIZED ")
		case nodes.CTEMaterializeNever:
			b.WriteString("NOT MATERIALIZED ")
		}
		b.WriteString("(")
		b.WriteString(d.deparseActionStmt(cte.Ctequery))
		b.WriteString(")")
		if sc, ok := cte.SearchClause.(*nodes.CTESearchClause); ok && sc != nil {
			if sc.SearchBreadthFirst {
				b.WriteString(" SEARCH BREADTH FIRST BY ")
			} else {
				b.WriteString(" SEARCH DEPTH FIRST BY ")
			}
			if sc.SearchColList != nil {
				for j, col := range sc.SearchColList.Items {
					if j > 0 {
						b.WriteString(", ")
					}
					b.WriteString(quoteIdentifier(stringVal(col)))
				}
			}
			b.WriteString(" SET ")
			b.WriteString(quoteIdentifier(sc.SearchSeqColumn))
		}
		if cc, ok := cte.CycleClause.(*nodes.CTECycleClause); ok && cc != nil {
			b.WriteString(" CYCLE ")
			if cc.CycleColList != nil {
				for j, col := range cc.CycleColList.Items {
					if j > 0 {
						b.WriteString(", ")
					}
					b.WriteString(quoteIdentifier(stringVal(col)))
				}
			}
			b.WriteString(" SET ")
			b.WriteString(quoteIdentifier(cc.CycleMarkColumn))
			if cc.CycleMarkValue != nil && cc.CycleMarkDefault != nil {
				b.WriteString(" TO ")
				b.WriteString(d.deparseExpr(cc.CycleMarkValue))
				b.WriteString(" DEFAULT ")
				b.WriteString(d.deparseExpr(cc.CycleMarkDefault))
			}
			b.WriteString(" USING ")
			b.WriteString(quoteIdentifier(cc.CyclePathColumn))
		}
	}
	b.WriteString(" ")
	return b.String()
}

func (d *ruleDeparser) deparseSelectStmt(s *nodes.SelectStmt) string {
	if s == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(d.deparseWithClause(s.WithClause))

	if s.Op != nodes.SETOP_NONE {
		if s.Larg != nil {
			b.WriteString("(")
			b.WriteString(d.deparseSelectStmt(s.Larg))
			b.WriteString(") ")
		}
		switch s.Op {
		case nodes.SETOP_UNION:
			b.WriteString("UNION ")
		case nodes.SETOP_INTERSECT:
			b.WriteString("INTERSECT ")
		case nodes.SETOP_EXCEPT:
			b.WriteString("EXCEPT ")
		}
		if s.All {
			b.WriteString("ALL ")
		}
		if s.Rarg != nil {
			b.WriteString("(")
			b.WriteString(d.deparseSelectStmt(s.Rarg))
			b.WriteString(")")
		}
	} else if s.ValuesLists != nil && len(s.ValuesLists.Items) > 0 {
		b.WriteString("VALUES ")
		for i, rowNode := range s.ValuesLists.Items {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("(")
			if row, ok := rowNode.(*nodes.List); ok {
				for j, v := range row.Items {
					if j > 0 {
						b.WriteString(", ")
					}
					b.WriteString(d.deparseExpr(v))
				}
			}
			b.WriteString(")")
		}
	} else {
		b.WriteString("SELECT")
		if s.DistinctClause != nil {
			if len(s.DistinctClause.Items) == 1 && s.DistinctClause.Items[0] == nil {
				b.WriteString(" DISTINCT")
			} else if len(s.DistinctClause.Items) > 0 {
				b.WriteString(" DISTINCT ON (")
				b.WriteString(d.deparseExprList(s.DistinctClause))
				b.WriteString(")")
			}
		}
		if s.TargetList != nil && len(s.TargetList.Items) > 0 {
			b.WriteString(" ")
			b.WriteString(d.deparseTargetList(s.TargetList))
		}
		if s.FromClause != nil && len(s.FromClause.Items) > 0 {
			b.WriteString(" FROM ")
			b.WriteString(d.deparseFromList(s.FromClause))
		}
		if s.WhereClause != nil {
			b.WriteString(" WHERE ")
			b.WriteString(d.deparseExpr(s.WhereClause))
		}
		if s.GroupClause != nil && len(s.GroupClause.Items) > 0 {
			b.WriteString(" GROUP BY ")
			if s.GroupDistinct {
				b.WriteString("DISTINCT ")
			}
			b.WriteString(d.deparseExprList(s.GroupClause))
		}
		if s.HavingClause != nil {
			b.WriteString(" HAVING ")
			b.WriteString(d.deparseExpr(s.HavingClause))
		}
		if s.WindowClause != nil && len(s.WindowClause.Items) > 0 {
			var winDefs []string
			for _, item := range s.WindowClause.Items {
				if wd, ok := item.(*nodes.WindowDef); ok && wd != nil && wd.Name != "" {
					winDefs = append(winDefs, fmt.Sprintf("%s AS (%s)", quoteIdentifier(wd.Name), d.deparseWindowSpec(wd)))
				}
			}
			if len(winDefs) > 0 {
				b.WriteString(" WINDOW ")
				b.WriteString(strings.Join(winDefs, ", "))
			}
		}
	}

	if s.SortClause != nil && len(s.SortClause.Items) > 0 {
		b.WriteString(" ORDER BY ")
		b.WriteString(d.deparseSortClause(s.SortClause))
	}
	if s.LimitOption == nodes.LIMIT_OPTION_WITH_TIES {
		if s.LimitOffset != nil {
			b.WriteString(" OFFSET ")
			b.WriteString(d.deparseExpr(s.LimitOffset))
		}
		if s.LimitCount != nil {
			b.WriteString(" FETCH FIRST ")
			b.WriteString(d.deparseExpr(s.LimitCount))
			b.WriteString(" ROWS WITH TIES")
		}
	} else {
		if s.LimitCount != nil {
			if ac, ok := s.LimitCount.(*nodes.A_Const); ok && ac.Isnull {
				b.WriteString(" LIMIT ALL")
			} else {
				b.WriteString(" LIMIT ")
				b.WriteString(d.deparseExpr(s.LimitCount))
			}
		}
		if s.LimitOffset != nil {
			b.WriteString(" OFFSET ")
			b.WriteString(d.deparseExpr(s.LimitOffset))
		}
	}
	if s.LockingClause != nil && len(s.LockingClause.Items) > 0 {
		for _, item := range s.LockingClause.Items {
			lc, ok := item.(*nodes.LockingClause)
			if !ok || lc == nil {
				continue
			}
			switch nodes.LockClauseStrength(lc.Strength) {
			case nodes.LCS_FORKEYSHARE:
				b.WriteString(" FOR KEY SHARE")
			case nodes.LCS_FORSHARE:
				b.WriteString(" FOR SHARE")
			case nodes.LCS_FORNOKEYUPDATE:
				b.WriteString(" FOR NO KEY UPDATE")
			case nodes.LCS_FORUPDATE:
				b.WriteString(" FOR UPDATE")
			default:
				continue
			}
			if lc.LockedRels != nil && len(lc.LockedRels.Items) > 0 {
				b.WriteString(" OF ")
				for i, rItem := range lc.LockedRels.Items {
					if i > 0 {
						b.WriteString(", ")
					}
					if rv, ok := rItem.(*nodes.RangeVar); ok && rv != nil {
						b.WriteString(d.deparseRangeVar(rv, false))
					}
				}
			}
			switch nodes.LockWaitPolicy(lc.WaitPolicy) {
			case nodes.LockWaitSkip:
				b.WriteString(" SKIP LOCKED")
			case nodes.LockWaitError:
				b.WriteString(" NOWAIT")
			}
		}
	}
	return b.String()
}

func (d *ruleDeparser) deparseInsertStmt(s *nodes.InsertStmt) string {
	var b strings.Builder
	b.WriteString(d.deparseWithClause(s.WithClause))
	b.WriteString("INSERT INTO ")
	b.WriteString(d.deparseRangeVar(s.Relation, false))
	if s.Cols != nil && len(s.Cols.Items) > 0 {
		b.WriteString(" (")
		for i, item := range s.Cols.Items {
			if i > 0 {
				b.WriteString(", ")
			}
			if rt, ok := item.(*nodes.ResTarget); ok {
				colName := rt.Name
				if d.b != nil {
					if cb, ok := d.b.colByRT[rt]; ok {
						colName = d.resolveBoundColName(cb, colName)
					}
				}
				b.WriteString(quoteIdentifier(colName))
				b.WriteString(d.deparseIndirection(rt.Indirection))
			}
		}
		b.WriteString(")")
	}
	switch s.Override {
	case nodes.OVERRIDING_USER_VALUE:
		b.WriteString(" OVERRIDING USER VALUE")
	case nodes.OVERRIDING_SYSTEM_VALUE:
		b.WriteString(" OVERRIDING SYSTEM VALUE")
	}
	if s.SelectStmt == nil {
		b.WriteString(" DEFAULT VALUES")
	} else if sel, ok := s.SelectStmt.(*nodes.SelectStmt); ok {
		b.WriteString(" ")
		b.WriteString(d.deparseSelectStmt(sel))
	}
	if s.OnConflictClause != nil {
		oc := s.OnConflictClause
		b.WriteString(" ON CONFLICT")
		if oc.Infer != nil {
			if oc.Infer.Conname != "" {
				conName := oc.Infer.Conname
				if d.cat != nil && d.b != nil {
					if conOID, ok := d.b.conByInfer[oc.Infer]; ok && conOID != 0 {
						if con := d.cat.constraints[conOID]; con != nil && con.Name != "" {
							conName = con.Name
						}
					}
				}
				b.WriteString(" ON CONSTRAINT ")
				b.WriteString(quoteIdentifier(conName))
			} else if oc.Infer.IndexElems != nil && len(oc.Infer.IndexElems.Items) > 0 {
				b.WriteString(" (")
				for i, ie := range oc.Infer.IndexElems.Items {
					if i > 0 {
						b.WriteString(", ")
					}
					if elem, ok := ie.(*nodes.IndexElem); ok && elem != nil {
						b.WriteString(d.deparseIndexElem(elem))
					}
				}
				b.WriteString(")")
				if oc.Infer.WhereClause != nil {
					b.WriteString(" WHERE ")
					b.WriteString(d.deparseExpr(oc.Infer.WhereClause))
				}
			}
		}
		if oc.TargetList != nil && len(oc.TargetList.Items) > 0 {
			b.WriteString(" DO UPDATE SET ")
			b.WriteString(d.deparseUpdateTargetList(oc.TargetList))
			if oc.WhereClause != nil {
				b.WriteString(" WHERE ")
				b.WriteString(d.deparseExpr(oc.WhereClause))
			}
		} else {
			b.WriteString(" DO NOTHING")
		}
	}
	if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
		b.WriteString(" RETURNING ")
		b.WriteString(d.deparseTargetList(s.ReturningList))
	}
	return b.String()
}

func (d *ruleDeparser) deparseUpdateStmt(s *nodes.UpdateStmt) string {
	var b strings.Builder
	b.WriteString(d.deparseWithClause(s.WithClause))
	b.WriteString("UPDATE ")
	b.WriteString(d.deparseRangeVar(s.Relation, true))
	b.WriteString(" SET ")
	b.WriteString(d.deparseUpdateTargetList(s.TargetList))
	if s.FromClause != nil && len(s.FromClause.Items) > 0 {
		b.WriteString(" FROM ")
		b.WriteString(d.deparseFromList(s.FromClause))
	}
	if s.WhereClause != nil {
		b.WriteString(" WHERE ")
		b.WriteString(d.deparseExpr(s.WhereClause))
	}
	if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
		b.WriteString(" RETURNING ")
		b.WriteString(d.deparseTargetList(s.ReturningList))
	}
	return b.String()
}

func (d *ruleDeparser) deparseDeleteStmt(s *nodes.DeleteStmt) string {
	var b strings.Builder
	b.WriteString(d.deparseWithClause(s.WithClause))
	b.WriteString("DELETE FROM ")
	b.WriteString(d.deparseRangeVar(s.Relation, true))
	if s.UsingClause != nil && len(s.UsingClause.Items) > 0 {
		b.WriteString(" USING ")
		b.WriteString(d.deparseFromList(s.UsingClause))
	}
	if s.WhereClause != nil {
		b.WriteString(" WHERE ")
		b.WriteString(d.deparseExpr(s.WhereClause))
	}
	if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
		b.WriteString(" RETURNING ")
		b.WriteString(d.deparseTargetList(s.ReturningList))
	}
	return b.String()
}

func (d *ruleDeparser) deparseTargetList(list *nodes.List) string {
	if list == nil {
		return ""
	}
	var parts []string
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok {
			parts = append(parts, d.deparseExpr(item))
			continue
		}
		if rt.Name == "" && d.b != nil && d.b.starExpandByRT != nil {
			if starCols, hasStar := d.b.starExpandByRT[rt]; hasStar && len(starCols) > 0 {
				for _, sc := range starCols {
					parts = append(parts, d.deparseStarCol(sc))
				}
				continue
			}
		}
		val := d.deparseExpr(rt.Val)
		if rt.Name != "" {
			parts = append(parts, val+" AS "+quoteIdentifier(rt.Name))
		} else if d.b != nil && d.b.targetNameByRT != nil {
			if origName, hasOrig := d.b.targetNameByRT[rt]; hasOrig && origName != "" && d.currentInferredTargetColName(rt.Val, origName) != origName {
				parts = append(parts, val+" AS "+quoteIdentifier(origName))
			} else {
				parts = append(parts, val)
			}
		} else {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, ", ")
}

func (d *ruleDeparser) deparseStarCol(sc ruleStarCol) string {
	colName := sc.colAlias
	if colName == "" {
		colName = d.resolveBoundColName(ruleColBinding{relOID: sc.relOID, attNum: sc.attNum}, sc.origName)
	}
	var exprSQL string
	if sc.prefix == "old" || sc.prefix == "new" {
		exprSQL = sc.prefix + "." + quoteIdentifier(colName)
	} else if sc.useRelName && d.cat != nil {
		if rel := d.cat.relationByOID[sc.relOID]; rel != nil {
			exprSQL = quoteIdentifier(rel.Name) + "." + quoteIdentifier(colName)
		} else {
			exprSQL = quoteIdentifier(colName)
		}
	} else if sc.prefix != "" {
		exprSQL = quoteIdentifier(sc.prefix) + "." + quoteIdentifier(colName)
	} else {
		exprSQL = quoteIdentifier(colName)
	}
	if sc.origName != "" && colName != sc.origName {
		return exprSQL + " AS " + quoteIdentifier(sc.origName)
	}
	return exprSQL
}

func (d *ruleDeparser) currentInferredTargetColName(n nodes.Node, fallback string) string {
	if d.b == nil || n == nil {
		return fallback
	}
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if cb, ok := d.b.colByRef[v]; ok {
			if cb.attNum == 0 {
				if cb.prefix == "old" || cb.prefix == "new" {
					return cb.prefix
				}
				if cb.useRelName && d.cat != nil {
					if rel := d.cat.relationByOID[cb.relOID]; rel != nil {
						return rel.Name
					}
				}
				if cb.prefix != "" {
					return cb.prefix
				}
				return fallback
			}
			return d.resolveBoundColName(cb, fallback)
		}
	case *nodes.FuncCall:
		if procOID, ok := d.b.funcByFC[v]; ok && d.cat != nil {
			if up := d.cat.userProcs[procOID]; up != nil && up.Name != "" {
				return up.Name
			}
		}
		if typeOID, ok := d.b.castTypeByFC[v]; ok && d.cat != nil {
			if bt := d.cat.typeByOID[typeOID]; bt != nil && bt.TypeName != "" {
				return bt.TypeName
			}
		}
	case *nodes.TypeCast:
		inner := inferRuleTargetColName(v.Arg)
		if inner != "" && inner != "?column?" {
			return d.currentInferredTargetColName(v.Arg, inner)
		}
		if v.TypeName != nil {
			if tb, ok := d.b.typeByTN[v.TypeName]; ok && d.cat != nil {
				if bt := d.cat.typeByOID[tb.typeOID]; bt != nil && bt.TypeName != "" {
					return bt.TypeName
				}
			}
		}
	case *nodes.CollateClause:
		return d.currentInferredTargetColName(v.Arg, fallback)
	}
	return fallback
}

func (d *ruleDeparser) deparseUpdateTargetCol(rt *nodes.ResTarget) string {
	colName := rt.Name
	if d.b != nil {
		if cb, ok := d.b.colByRT[rt]; ok {
			colName = d.resolveBoundColName(cb, colName)
		}
	}
	return quoteIdentifier(colName) + d.deparseIndirection(rt.Indirection)
}

func (d *ruleDeparser) deparseUpdateTargetList(list *nodes.List) string {
	if list == nil {
		return ""
	}
	var parts []string
	for i := 0; i < len(list.Items); i++ {
		rt, ok := list.Items[i].(*nodes.ResTarget)
		if !ok {
			continue
		}
		if mar, isMulti := rt.Val.(*nodes.MultiAssignRef); isMulti && mar.Ncolumns > 0 {
			var cols []string
			cols = append(cols, d.deparseUpdateTargetCol(rt))
			for k := 1; k < mar.Ncolumns && i+1 < len(list.Items); k++ {
				i++
				if nextRT, ok := list.Items[i].(*nodes.ResTarget); ok {
					cols = append(cols, d.deparseUpdateTargetCol(nextRT))
				}
			}
			parts = append(parts, "("+strings.Join(cols, ", ")+") = "+d.deparseExpr(mar.Source))
			continue
		}
		parts = append(parts, d.deparseUpdateTargetCol(rt)+" = "+d.deparseExpr(rt.Val))
	}
	return strings.Join(parts, ", ")
}

func (d *ruleDeparser) deparseFromList(list *nodes.List) string {
	if list == nil {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		parts = append(parts, d.deparseTableRef(item))
	}
	return strings.Join(parts, ", ")
}

func (d *ruleDeparser) deparseTableRef(n nodes.Node) string {
	switch v := n.(type) {
	case *nodes.RangeVar:
		return d.deparseRangeVar(v, true)
	case *nodes.RangeSubselect:
		var b strings.Builder
		if v.Lateral {
			b.WriteString("LATERAL ")
		}
		b.WriteString("(")
		b.WriteString(d.deparseActionStmt(v.Subquery))
		b.WriteString(")")
		b.WriteString(deparseRawAlias(v.Alias))
		return b.String()
	case *nodes.RangeFunction:
		return d.deparseRangeFunction(v)
	case *nodes.RangeTableFunc:
		return d.deparseRangeTableFunc(v)
	case *nodes.JsonTable:
		return d.deparseJsonTable(v)
	case *nodes.RangeTableSample:
		var b strings.Builder
		b.WriteString(d.deparseTableRef(v.Relation))
		b.WriteString(" TABLESAMPLE ")
		b.WriteString(d.deparseTablesampleMethod(v))
		b.WriteString("(")
		b.WriteString(d.deparseExprList(v.Args))
		b.WriteString(")")
		if v.Repeatable != nil {
			b.WriteString(" REPEATABLE (")
			b.WriteString(d.deparseExpr(v.Repeatable))
			b.WriteString(")")
		}
		return b.String()
	case *nodes.JoinExpr:
		var b strings.Builder
		hasAlias := v.Alias != nil && v.Alias.Aliasname != ""
		if hasAlias {
			b.WriteString("(")
		}
		b.WriteString(d.deparseTableRef(v.Larg))
		naturalResolved := d.b != nil && d.b.naturalJoined[v]
		var boundUsing []ruleJoinUsingCol
		if d.b != nil {
			boundUsing = d.b.joinUsing[v]
		}
		if v.IsNatural && !naturalResolved {
			b.WriteString(" NATURAL")
		}
		switch v.Jointype {
		case nodes.JOIN_LEFT:
			b.WriteString(" LEFT JOIN ")
		case nodes.JOIN_RIGHT:
			b.WriteString(" RIGHT JOIN ")
		case nodes.JOIN_FULL:
			b.WriteString(" FULL JOIN ")
		default:
			if v.Quals == nil && v.UsingClause == nil && (!v.IsNatural || (naturalResolved && len(boundUsing) == 0)) {
				b.WriteString(" CROSS JOIN ")
			} else {
				b.WriteString(" JOIN ")
			}
		}
		b.WriteString(d.deparseTableRef(v.Rarg))
		if len(boundUsing) > 0 {
			b.WriteString(" USING (")
			for i, uc := range boundUsing {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(quoteIdentifier(d.resolveJoinUsingColName(uc)))
			}
			b.WriteString(")")
			if v.JoinUsing != nil && v.JoinUsing.Aliasname != "" {
				b.WriteString(" AS ")
				b.WriteString(quoteIdentifier(v.JoinUsing.Aliasname))
			}
		} else if naturalResolved && v.Jointype != nodes.JOIN_INNER {
			b.WriteString(" ON true")
		} else if v.UsingClause != nil && len(v.UsingClause.Items) > 0 {
			b.WriteString(" USING (")
			for i, u := range v.UsingClause.Items {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(quoteIdentifier(stringVal(u)))
			}
			b.WriteString(")")
			if v.JoinUsing != nil && v.JoinUsing.Aliasname != "" {
				b.WriteString(" AS ")
				b.WriteString(quoteIdentifier(v.JoinUsing.Aliasname))
			}
		} else if v.Quals != nil {
			b.WriteString(" ON ")
			b.WriteString(d.deparseExpr(v.Quals))
		}
		if hasAlias {
			b.WriteString(")")
			b.WriteString(deparseRawAlias(v.Alias))
		}
		return b.String()
	default:
		return d.deparseExpr(n)
	}
}

func (d *ruleDeparser) deparseRangeVar(rv *nodes.RangeVar, includeOnly bool) string {
	if rv == nil {
		return ""
	}
	var b strings.Builder
	if includeOnly && !rv.Inh && (rv.Loc.Start != 0 || rv.Loc.End != 0) {
		b.WriteString("ONLY ")
	}
	schemaName := rv.Schemaname
	relName := rv.Relname
	if d.cat != nil && d.b != nil {
		if relOID, ok := d.b.relByRV[rv]; ok {
			if rel := d.cat.relationByOID[relOID]; rel != nil {
				relName = rel.Name
				if rel.Schema != nil {
					if schemaName != "" || rel.Schema.Name != "public" {
						schemaName = rel.Schema.Name
					}
				}
			}
		} else if seqOID, ok := d.b.seqByRV[rv]; ok {
			if seq := d.cat.sequenceByOID[seqOID]; seq != nil {
				relName = seq.Name
				if seq.Schema != nil {
					if schemaName != "" || seq.Schema.Name != "public" {
						schemaName = seq.Schema.Name
					}
				}
			}
		}
	}
	if schemaName != "" {
		b.WriteString(quoteIdentifier(schemaName))
		b.WriteString(".")
	}
	b.WriteString(quoteIdentifier(relName))
	b.WriteString(deparseRawAlias(rv.Alias))
	return b.String()
}

func (c *Catalog) resolveTablesampleHandler(method *nodes.List) (uint32, error) {
	if method == nil || len(method.Items) == 0 {
		return 0, nil
	}
	fSchema, fName := qualifiedName(method)
	if fName == "" {
		return 0, nil
	}
	var searchNS []uint32
	if fSchema != "" {
		ns := c.schemaByName[fSchema]
		if ns == nil {
			return 0, errUndefinedSchema(fSchema)
		}
		searchNS = []uint32{ns.OID}
	} else {
		searchNS = c.routineSearchPathWithCatalog()
	}
	shown := fName
	if fSchema != "" {
		shown = fSchema + "." + fName
	}
	for _, nsOID := range searchNS {
		for _, p := range c.procByName[fName] {
			if p.Kind == 'p' || len(p.ArgTypes) != 1 || p.ArgTypes[0] != INTERNALOID {
				continue
			}
			pNS := uint32(PGCatalogNamespace)
			if p.OID >= FirstNormalObjectId {
				up := c.userProcs[p.OID]
				if up == nil || up.Schema == nil {
					continue
				}
				pNS = up.Schema.OID
			}
			if pNS != nsOID {
				continue
			}
			if p.RetType != TSMHANDLEROID {
				return 0, &Error{
					Code:    CodeWrongObjectType,
					Message: fmt.Sprintf("function %s must return type tsm_handler", shown),
				}
			}
			return p.OID, nil
		}
	}
	return 0, &Error{
		Code:    CodeUndefinedFunction,
		Message: fmt.Sprintf("tablesample method %s does not exist", shown),
	}
}

func (d *ruleDeparser) deparseTablesampleMethod(v *nodes.RangeTableSample) string {
	if v == nil {
		return ""
	}
	if d.cat != nil && d.b != nil {
		if procOID, ok := d.b.tsmByRTS[v]; ok {
			if up := d.cat.userProcs[procOID]; up != nil {
				origSchema, _ := qualifiedName(v.Method)
				if up.Schema != nil && (origSchema != "" || up.Schema.Name != "public") {
					return quoteIdentifier(up.Schema.Name) + "." + quoteIdentifier(up.Name)
				}
				return quoteIdentifier(up.Name)
			}
		}
	}
	return deparseFuncNameList(v.Method)
}

func deparseRawAlias(a *nodes.Alias) string {
	if a == nil || a.Aliasname == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(" AS ")
	b.WriteString(quoteIdentifier(a.Aliasname))
	if a.Colnames != nil && len(a.Colnames.Items) > 0 {
		b.WriteString(" (")
		for i, c := range a.Colnames.Items {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(quoteIdentifier(stringVal(c)))
		}
		b.WriteString(")")
	}
	return b.String()
}

func (d *ruleDeparser) deparseRangeFunction(v *nodes.RangeFunction) string {
	if v == nil {
		return ""
	}
	var b strings.Builder
	if v.Lateral {
		b.WriteString("LATERAL ")
	}
	if v.IsRowsfrom {
		b.WriteString("ROWS FROM (")
	}
	if v.Functions != nil {
		for i, item := range v.Functions.Items {
			if i > 0 {
				b.WriteString(", ")
			}
			if pair, ok := item.(*nodes.List); ok && len(pair.Items) > 0 {
				b.WriteString(d.deparseExpr(pair.Items[0]))
				if len(pair.Items) > 1 {
					if cdList, ok := pair.Items[1].(*nodes.List); ok && len(cdList.Items) > 0 {
						b.WriteString(" AS (")
						b.WriteString(d.deparseTableFuncColDefs(cdList))
						b.WriteString(")")
					}
				}
			} else {
				b.WriteString(d.deparseExpr(item))
			}
		}
	}
	if v.IsRowsfrom {
		b.WriteString(")")
	}
	if v.Ordinality {
		b.WriteString(" WITH ORDINALITY")
	}
	if v.Alias != nil && v.Alias.Aliasname != "" {
		b.WriteString(" AS ")
		b.WriteString(quoteIdentifier(v.Alias.Aliasname))
		if v.Coldeflist != nil && len(v.Coldeflist.Items) > 0 {
			b.WriteString(" (")
			b.WriteString(d.deparseTableFuncColDefs(v.Coldeflist))
			b.WriteString(")")
		} else if v.Alias.Colnames != nil && len(v.Alias.Colnames.Items) > 0 {
			b.WriteString(" (")
			for i, c := range v.Alias.Colnames.Items {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(quoteIdentifier(stringVal(c)))
			}
			b.WriteString(")")
		}
	} else if v.Coldeflist != nil && len(v.Coldeflist.Items) > 0 {
		b.WriteString(" AS (")
		b.WriteString(d.deparseTableFuncColDefs(v.Coldeflist))
		b.WriteString(")")
	}
	return b.String()
}

func (d *ruleDeparser) deparseTableFuncColDefs(list *nodes.List) string {
	if list == nil || len(list.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		cd, ok := item.(*nodes.ColumnDef)
		if !ok || cd == nil {
			continue
		}
		var cb strings.Builder
		cb.WriteString(quoteIdentifier(cd.Colname))
		if cd.TypeName != nil {
			cb.WriteString(" ")
			cb.WriteString(d.deparseTypeName(cd.TypeName))
		}
		if cd.CollClause != nil {
			collSQL := deparseQualifiedIdentList(cd.CollClause.Collname)
			if d.cat != nil && d.b != nil {
				if uid, ok := d.b.collByCC[cd.CollClause]; ok && uid != 0 {
					for _, s := range d.cat.schemas {
						if s == nil {
							continue
						}
						for _, uc := range s.Collations {
							if uc != nil && uc.OID == uid {
								if s.Name != "public" && s.Name != "pg_catalog" {
									collSQL = quoteIdentifier(s.Name) + "." + quoteIdentifier(uc.Name)
								} else {
									collSQL = quoteIdentifier(uc.Name)
								}
								break
							}
						}
					}
				}
			}
			cb.WriteString(" COLLATE ")
			cb.WriteString(collSQL)
		}
		parts = append(parts, cb.String())
	}
	return strings.Join(parts, ", ")
}

func (d *ruleDeparser) deparseRangeTableFunc(v *nodes.RangeTableFunc) string {
	if v == nil {
		return ""
	}
	var b strings.Builder
	if v.Lateral {
		b.WriteString("LATERAL ")
	}
	b.WriteString("XMLTABLE(")
	if v.Namespaces != nil && len(v.Namespaces.Items) > 0 {
		b.WriteString("XMLNAMESPACES(")
		for i, item := range v.Namespaces.Items {
			if i > 0 {
				b.WriteString(", ")
			}
			if rt, ok := item.(*nodes.ResTarget); ok && rt != nil {
				if rt.Name == "" {
					b.WriteString("DEFAULT ")
					b.WriteString(d.deparseExpr(rt.Val))
				} else {
					b.WriteString(d.deparseExpr(rt.Val))
					b.WriteString(" AS ")
					b.WriteString(quoteIdentifier(rt.Name))
				}
			} else {
				b.WriteString(d.deparseExpr(item))
			}
		}
		b.WriteString("), ")
	}
	b.WriteString(d.deparseExpr(v.Rowexpr))
	b.WriteString(" PASSING ")
	b.WriteString(d.deparseExpr(v.Docexpr))
	if v.Columns != nil && len(v.Columns.Items) > 0 {
		b.WriteString(" COLUMNS ")
		for i, item := range v.Columns.Items {
			fc, ok := item.(*nodes.RangeTableFuncCol)
			if !ok || fc == nil {
				continue
			}
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(quoteIdentifier(fc.Colname))
			if fc.ForOrdinality {
				b.WriteString(" FOR ORDINALITY")
				continue
			}
			if fc.TypeName != nil {
				b.WriteString(" ")
				b.WriteString(d.deparseTypeName(fc.TypeName))
			}
			if fc.Colexpr != nil {
				b.WriteString(" PATH ")
				b.WriteString(d.deparseExpr(fc.Colexpr))
			}
			if fc.Coldefexpr != nil {
				b.WriteString(" DEFAULT ")
				b.WriteString(d.deparseExpr(fc.Coldefexpr))
			}
			if fc.IsNotNull {
				b.WriteString(" NOT NULL")
			}
		}
	}
	b.WriteString(")")
	b.WriteString(deparseRawAlias(v.Alias))
	return b.String()
}

func (d *ruleDeparser) deparseJsonTable(v *nodes.JsonTable) string {
	if v == nil {
		return ""
	}
	var b strings.Builder
	if v.Lateral {
		b.WriteString("LATERAL ")
	}
	b.WriteString("JSON_TABLE(")
	if v.ContextItem != nil {
		b.WriteString(d.deparseExpr(v.ContextItem))
	}
	if v.Pathspec != nil {
		b.WriteString(", ")
		b.WriteString(d.deparseExpr(v.Pathspec.String))
		if v.Pathspec.Name != "" {
			b.WriteString(" AS ")
			b.WriteString(quoteIdentifier(v.Pathspec.Name))
		}
	}
	if v.Passing != nil && len(v.Passing.Items) > 0 {
		b.WriteString(" PASSING ")
		for i, item := range v.Passing.Items {
			if ja, ok := item.(*nodes.JsonArgument); ok && ja != nil {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(d.deparseExpr(ja.Val))
				b.WriteString(" AS ")
				b.WriteString(quoteIdentifier(ja.Name))
			}
		}
	}
	if v.Columns != nil && len(v.Columns.Items) > 0 {
		b.WriteString(" COLUMNS (")
		b.WriteString(d.deparseJsonTableColumns(v.Columns))
		b.WriteString(")")
	}
	if v.OnError != nil {
		b.WriteString(" ")
		b.WriteString(d.deparseJsonBehavior(v.OnError))
		b.WriteString(" ON ERROR")
	}
	b.WriteString(")")
	b.WriteString(deparseRawAlias(v.Alias))
	return b.String()
}

func (d *ruleDeparser) deparseJsonTableColumns(list *nodes.List) string {
	if list == nil || len(list.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		jtc, ok := item.(*nodes.JsonTableColumn)
		if !ok || jtc == nil {
			continue
		}
		var cb strings.Builder
		switch jtc.Coltype {
		case nodes.JTC_FOR_ORDINALITY:
			cb.WriteString(quoteIdentifier(jtc.Name))
			cb.WriteString(" FOR ORDINALITY")
		case nodes.JTC_NESTED:
			cb.WriteString("NESTED PATH ")
			if jtc.Pathspec != nil {
				cb.WriteString(d.deparseExpr(jtc.Pathspec.String))
				if jtc.Pathspec.Name != "" {
					cb.WriteString(" AS ")
					cb.WriteString(quoteIdentifier(jtc.Pathspec.Name))
				}
			}
			cb.WriteString(" COLUMNS (")
			cb.WriteString(d.deparseJsonTableColumns(jtc.Columns))
			cb.WriteString(")")
		case nodes.JTC_EXISTS:
			cb.WriteString(quoteIdentifier(jtc.Name))
			if jtc.TypeName != nil {
				cb.WriteString(" ")
				cb.WriteString(d.deparseTypeName(jtc.TypeName))
			}
			cb.WriteString(" EXISTS")
			if jtc.Pathspec != nil && jtc.Pathspec.String != nil {
				cb.WriteString(" PATH ")
				cb.WriteString(d.deparseExpr(jtc.Pathspec.String))
			}
			if jtc.OnError != nil {
				cb.WriteString(" ")
				cb.WriteString(d.deparseJsonBehavior(jtc.OnError))
				cb.WriteString(" ON ERROR")
			}
		default:
			cb.WriteString(quoteIdentifier(jtc.Name))
			if jtc.TypeName != nil {
				cb.WriteString(" ")
				cb.WriteString(d.deparseTypeName(jtc.TypeName))
			}
			if f := deparseJsonFormat(jtc.Format); f != "" {
				cb.WriteString(" ")
				cb.WriteString(f)
			}
			if jtc.Pathspec != nil && jtc.Pathspec.String != nil {
				cb.WriteString(" PATH ")
				cb.WriteString(d.deparseExpr(jtc.Pathspec.String))
			}
			switch jtc.Wrapper {
			case nodes.JSW_NONE:
				cb.WriteString(" WITHOUT ARRAY WRAPPER")
			case nodes.JSW_CONDITIONAL:
				cb.WriteString(" WITH CONDITIONAL ARRAY WRAPPER")
			case nodes.JSW_UNCONDITIONAL:
				cb.WriteString(" WITH UNCONDITIONAL ARRAY WRAPPER")
			}
			switch jtc.Quotes {
			case nodes.JS_QUOTES_KEEP:
				cb.WriteString(" KEEP QUOTES")
			case nodes.JS_QUOTES_OMIT:
				cb.WriteString(" OMIT QUOTES")
			}
			if jtc.OnEmpty != nil {
				cb.WriteString(" ")
				cb.WriteString(d.deparseJsonBehavior(jtc.OnEmpty))
				cb.WriteString(" ON EMPTY")
			}
			if jtc.OnError != nil {
				cb.WriteString(" ")
				cb.WriteString(d.deparseJsonBehavior(jtc.OnError))
				cb.WriteString(" ON ERROR")
			}
		}
		parts = append(parts, cb.String())
	}
	return strings.Join(parts, ", ")
}

func (d *ruleDeparser) deparseSortClause(list *nodes.List) string {
	if list == nil {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		sb, ok := item.(*nodes.SortBy)
		if !ok {
			continue
		}
		s := d.deparseExpr(sb.Node)
		switch sb.SortbyDir {
		case nodes.SORTBY_ASC:
			s += " ASC"
		case nodes.SORTBY_DESC:
			s += " DESC"
		case nodes.SORTBY_USING:
			if sb.UseOp != nil {
				s += " USING " + d.deparseOperatorRef(sb.UseOp)
			}
		}
		switch sb.SortbyNulls {
		case nodes.SORTBY_NULLS_FIRST:
			s += " NULLS FIRST"
		case nodes.SORTBY_NULLS_LAST:
			s += " NULLS LAST"
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ", ")
}

func (d *ruleDeparser) deparseWindowSpec(wd *nodes.WindowDef) string {
	if wd == nil {
		return ""
	}
	var parts []string
	if wd.Refname != "" {
		parts = append(parts, quoteIdentifier(wd.Refname))
	}
	if wd.PartitionClause != nil && len(wd.PartitionClause.Items) > 0 {
		parts = append(parts, "PARTITION BY "+d.deparseExprList(wd.PartitionClause))
	}
	if wd.OrderClause != nil && len(wd.OrderClause.Items) > 0 {
		parts = append(parts, "ORDER BY "+d.deparseSortClause(wd.OrderClause))
	}
	fo := wd.FrameOptions
	if fo&nodes.FRAMEOPTION_NONDEFAULT != 0 {
		var fb strings.Builder
		switch {
		case fo&nodes.FRAMEOPTION_RANGE != 0:
			fb.WriteString("RANGE ")
		case fo&nodes.FRAMEOPTION_ROWS != 0:
			fb.WriteString("ROWS ")
		case fo&nodes.FRAMEOPTION_GROUPS != 0:
			fb.WriteString("GROUPS ")
		}
		startBound := d.deparseFrameStartBound(fo, wd.StartOffset)
		if fo&nodes.FRAMEOPTION_BETWEEN != 0 {
			fb.WriteString("BETWEEN ")
			fb.WriteString(startBound)
			fb.WriteString(" AND ")
			fb.WriteString(d.deparseFrameEndBound(fo, wd.EndOffset))
		} else {
			fb.WriteString(startBound)
		}
		switch {
		case fo&nodes.FRAMEOPTION_EXCLUDE_CURRENT_ROW != 0:
			fb.WriteString(" EXCLUDE CURRENT ROW")
		case fo&nodes.FRAMEOPTION_EXCLUDE_GROUP != 0:
			fb.WriteString(" EXCLUDE GROUP")
		case fo&nodes.FRAMEOPTION_EXCLUDE_TIES != 0:
			fb.WriteString(" EXCLUDE TIES")
		}
		parts = append(parts, fb.String())
	}
	return strings.Join(parts, " ")
}

func (d *ruleDeparser) deparseFrameStartBound(fo int, offset nodes.Node) string {
	switch {
	case fo&nodes.FRAMEOPTION_START_UNBOUNDED_PRECEDING != 0:
		return "UNBOUNDED PRECEDING"
	case fo&nodes.FRAMEOPTION_START_UNBOUNDED_FOLLOWING != 0:
		return "UNBOUNDED FOLLOWING"
	case fo&nodes.FRAMEOPTION_START_CURRENT_ROW != 0:
		return "CURRENT ROW"
	case fo&nodes.FRAMEOPTION_START_OFFSET_PRECEDING != 0:
		return d.deparseExpr(offset) + " PRECEDING"
	case fo&nodes.FRAMEOPTION_START_OFFSET_FOLLOWING != 0:
		return d.deparseExpr(offset) + " FOLLOWING"
	default:
		return "UNBOUNDED PRECEDING"
	}
}

func (d *ruleDeparser) deparseFrameEndBound(fo int, offset nodes.Node) string {
	switch {
	case fo&nodes.FRAMEOPTION_END_UNBOUNDED_PRECEDING != 0:
		return "UNBOUNDED PRECEDING"
	case fo&nodes.FRAMEOPTION_END_UNBOUNDED_FOLLOWING != 0:
		return "UNBOUNDED FOLLOWING"
	case fo&nodes.FRAMEOPTION_END_CURRENT_ROW != 0:
		return "CURRENT ROW"
	case fo&nodes.FRAMEOPTION_END_OFFSET_PRECEDING != 0:
		return d.deparseExpr(offset) + " PRECEDING"
	case fo&nodes.FRAMEOPTION_END_OFFSET_FOLLOWING != 0:
		return d.deparseExpr(offset) + " FOLLOWING"
	default:
		return "CURRENT ROW"
	}
}

func (d *ruleDeparser) deparseExprList(list *nodes.List) string {
	if list == nil {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		parts = append(parts, d.deparseExpr(item))
	}
	return strings.Join(parts, ", ")
}

func (d *ruleDeparser) deparseIndirection(ind *nodes.List) string {
	if ind == nil {
		return ""
	}
	var b strings.Builder
	for _, item := range ind.Items {
		switch v := item.(type) {
		case *nodes.A_Indices:
			b.WriteString("[")
			if v.IsSlice {
				if v.Lidx != nil {
					b.WriteString(d.deparseExpr(v.Lidx))
				}
				b.WriteString(":")
				if v.Uidx != nil {
					b.WriteString(d.deparseExpr(v.Uidx))
				}
			} else if v.Uidx != nil {
				b.WriteString(d.deparseExpr(v.Uidx))
			}
			b.WriteString("]")
		case *nodes.String:
			b.WriteString(".")
			b.WriteString(quoteIdentifier(v.Str))
		case *nodes.A_Star:
			b.WriteString(".*")
		}
	}
	return b.String()
}

func (d *ruleDeparser) deparseExpr(n nodes.Node) string {
	if n == nil {
		return ""
	}
	if d.b != nil {
		if pr, ok := d.b.paramByNode[n]; ok {
			return d.b.fn.deparse(d.cat, pr)
		}
	}
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if v.Fields == nil {
			return ""
		}
		if d.b != nil {
			if cb, ok := d.b.colByRef[v]; ok {
				if cb.attNum == 0 {
					if cb.prefix == "old" || cb.prefix == "new" {
						return cb.prefix
					}
					if cb.useRelName && d.cat != nil {
						if rel := d.cat.relationByOID[cb.relOID]; rel != nil {
							return quoteIdentifier(rel.Name)
						}
					}
					if cb.prefix != "" {
						return quoteIdentifier(cb.prefix)
					}
					return quoteIdentifier(stringVal(v.Fields.Items[len(v.Fields.Items)-1]))
				}
				colName := d.resolveBoundColName(cb, stringVal(v.Fields.Items[len(v.Fields.Items)-1]))
				if cb.prefix == "old" || cb.prefix == "new" {
					return cb.prefix + "." + quoteIdentifier(colName)
				}
				if cb.useRelName && d.cat != nil {
					if rel := d.cat.relationByOID[cb.relOID]; rel != nil {
						return quoteIdentifier(rel.Name) + "." + quoteIdentifier(colName)
					}
				}
				if cb.prefix != "" {
					return quoteIdentifier(cb.prefix) + "." + quoteIdentifier(colName)
				}
				return quoteIdentifier(colName)
			}
		}
		parts := make([]string, 0, len(v.Fields.Items))
		for i, f := range v.Fields.Items {
			switch fv := f.(type) {
			case *nodes.String:
				if i == 0 && (strings.EqualFold(fv.Str, "old") || strings.EqualFold(fv.Str, "new")) {
					parts = append(parts, strings.ToLower(fv.Str))
				} else {
					parts = append(parts, quoteIdentifier(fv.Str))
				}
			case *nodes.A_Star:
				parts = append(parts, "*")
			}
		}
		return strings.Join(parts, ".")

	case *nodes.A_Const:
		if v.Isnull {
			return "NULL"
		}
		if s := d.deparseBoundRegclassLiteral(v); s != "" {
			return s
		}
		return deparseRawConstVal(v.Val)

	case *nodes.Integer, *nodes.Float, *nodes.String, *nodes.Boolean, *nodes.BitString:
		if s := d.deparseBoundRegclassLiteral(v); s != "" {
			return s
		}
		return deparseRawConstVal(v)

	case *nodes.CurrentOfExpr:
		if v.CursorName != "" {
			return "CURRENT OF " + quoteIdentifier(v.CursorName)
		}
		if v.CursorParam > 0 {
			return fmt.Sprintf("CURRENT OF $%d", v.CursorParam)
		}
		return "CURRENT OF"

	case *nodes.ParamRef:
		return fmt.Sprintf("$%d", v.Number)

	case *nodes.SetToDefault:
		return "DEFAULT"

	case *nodes.A_Star:
		return "*"

	case *nodes.TypeCast:
		return fmt.Sprintf("(%s)::%s", d.deparseExpr(v.Arg), d.deparseTypeName(v.TypeName))

	case *nodes.CollateClause:
		collSQL := deparseQualifiedIdentList(v.Collname)
		if d.cat != nil && d.b != nil {
			if uid, ok := d.b.collByCC[v]; ok && uid != 0 {
				for _, s := range d.cat.schemas {
					if s == nil {
						continue
					}
					for _, uc := range s.Collations {
						if uc != nil && uc.OID == uid {
							if s.Name != "public" {
								collSQL = quoteIdentifier(s.Name) + "." + quoteIdentifier(uc.Name)
							} else {
								collSQL = quoteIdentifier(uc.Name)
							}
							break
						}
					}
				}
			}
		}
		return fmt.Sprintf("(%s COLLATE %s)", d.deparseExpr(v.Arg), collSQL)

	case *nodes.A_Indirection:
		return fmt.Sprintf("(%s)%s", d.deparseExpr(v.Arg), d.deparseIndirection(v.Indirection))

	case *nodes.A_Expr:
		op := d.deparseOperatorRef(v.Name)
		switch v.Kind {
		case nodes.AEXPR_OP:
			if v.Lexpr == nil {
				return fmt.Sprintf("(%s %s)", op, d.deparseExpr(v.Rexpr))
			}
			return fmt.Sprintf("(%s %s %s)", d.deparseExpr(v.Lexpr), op, d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_OP_ANY:
			return fmt.Sprintf("(%s %s ANY (%s))", d.deparseExpr(v.Lexpr), op, d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_OP_ALL:
			return fmt.Sprintf("(%s %s ALL (%s))", d.deparseExpr(v.Lexpr), op, d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_DISTINCT:
			return fmt.Sprintf("(%s IS DISTINCT FROM %s)", d.deparseExpr(v.Lexpr), d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_NOT_DISTINCT:
			return fmt.Sprintf("(%s IS NOT DISTINCT FROM %s)", d.deparseExpr(v.Lexpr), d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_NULLIF:
			return fmt.Sprintf("NULLIF(%s, %s)", d.deparseExpr(v.Lexpr), d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_IN:
			kw := "IN"
			if op == "<>" {
				kw = "NOT IN"
			}
			rhs, _ := v.Rexpr.(*nodes.List)
			return fmt.Sprintf("(%s %s (%s))", d.deparseExpr(v.Lexpr), kw, d.deparseExprList(rhs))
		case nodes.AEXPR_LIKE:
			kw := "LIKE"
			if op == "!~~" {
				kw = "NOT LIKE"
			}
			return fmt.Sprintf("(%s %s %s)", d.deparseExpr(v.Lexpr), kw, d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_ILIKE:
			kw := "ILIKE"
			if op == "!~~*" {
				kw = "NOT ILIKE"
			}
			return fmt.Sprintf("(%s %s %s)", d.deparseExpr(v.Lexpr), kw, d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_SIMILAR:
			kw := "SIMILAR TO"
			if op == "!~" {
				kw = "NOT SIMILAR TO"
			}
			return fmt.Sprintf("(%s %s %s)", d.deparseExpr(v.Lexpr), kw, d.deparseExpr(v.Rexpr))
		case nodes.AEXPR_BETWEEN, nodes.AEXPR_NOT_BETWEEN, nodes.AEXPR_BETWEEN_SYM, nodes.AEXPR_NOT_BETWEEN_SYM:
			kw := "BETWEEN"
			switch v.Kind {
			case nodes.AEXPR_NOT_BETWEEN:
				kw = "NOT BETWEEN"
			case nodes.AEXPR_BETWEEN_SYM:
				kw = "BETWEEN SYMMETRIC"
			case nodes.AEXPR_NOT_BETWEEN_SYM:
				kw = "NOT BETWEEN SYMMETRIC"
			}
			if pair, ok := v.Rexpr.(*nodes.List); ok && len(pair.Items) == 2 {
				return fmt.Sprintf("(%s %s %s AND %s)", d.deparseExpr(v.Lexpr), kw, d.deparseExpr(pair.Items[0]), d.deparseExpr(pair.Items[1]))
			}
		}
		return fmt.Sprintf("(%s %s %s)", d.deparseExpr(v.Lexpr), op, d.deparseExpr(v.Rexpr))

	case *nodes.BoolExpr:
		if v.Args == nil || len(v.Args.Items) == 0 {
			return ""
		}
		if v.Boolop == nodes.NOT_EXPR {
			return fmt.Sprintf("(NOT %s)", d.deparseExpr(v.Args.Items[0]))
		}
		sep := " AND "
		if v.Boolop == nodes.OR_EXPR {
			sep = " OR "
		}
		parts := make([]string, 0, len(v.Args.Items))
		for _, a := range v.Args.Items {
			parts = append(parts, d.deparseExpr(a))
		}
		return "(" + strings.Join(parts, sep) + ")"

	case *nodes.NullTest:
		if v.Nulltesttype == nodes.IS_NOT_NULL {
			return fmt.Sprintf("(%s IS NOT NULL)", d.deparseExpr(v.Arg))
		}
		return fmt.Sprintf("(%s IS NULL)", d.deparseExpr(v.Arg))

	case *nodes.BooleanTest:
		suffix := "IS TRUE"
		switch v.Booltesttype {
		case nodes.IS_NOT_TRUE:
			suffix = "IS NOT TRUE"
		case nodes.IS_FALSE:
			suffix = "IS FALSE"
		case nodes.IS_NOT_FALSE:
			suffix = "IS NOT FALSE"
		case nodes.IS_UNKNOWN:
			suffix = "IS UNKNOWN"
		case nodes.IS_NOT_UNKNOWN:
			suffix = "IS NOT UNKNOWN"
		}
		return fmt.Sprintf("(%s %s)", d.deparseExpr(v.Arg), suffix)

	case *nodes.FuncCall:
		var b strings.Builder
		fnSQL := deparseFuncNameList(v.Funcname)
		if d.cat != nil && d.b != nil {
			if procOID, ok := d.b.funcByFC[v]; ok {
				if up := d.cat.userProcs[procOID]; up != nil {
					if up.Schema != nil && (len(v.Funcname.Items) > 1 || up.Schema.Name != "public") {
						fnSQL = quoteIdentifier(up.Schema.Name) + "." + quoteIdentifier(up.Name)
					} else {
						fnSQL = quoteIdentifier(up.Name)
					}
				}
			} else if typeOID, ok := d.b.castTypeByFC[v]; ok && typeOID != 0 {
				if bt := d.cat.typeByOID[typeOID]; bt != nil {
					if sc := d.cat.schemas[bt.Namespace]; sc != nil && (len(v.Funcname.Items) > 1 || (sc.Name != "public" && sc.Name != "pg_catalog")) {
						fnSQL = quoteIdentifier(sc.Name) + "." + quoteIdentifier(bt.TypeName)
					} else {
						fnSQL = quoteIdentifier(bt.TypeName)
					}
				}
			}
		}
		b.WriteString(fnSQL)
		b.WriteString("(")
		if v.AggStar {
			b.WriteString("*")
		} else {
			if v.AggDistinct {
				b.WriteString("DISTINCT ")
			}
			if v.Args != nil {
				for i, arg := range v.Args.Items {
					if i > 0 {
						b.WriteString(", ")
					}
					if v.FuncVariadic && i == len(v.Args.Items)-1 {
						b.WriteString("VARIADIC ")
					}
					b.WriteString(d.deparseExpr(arg))
				}
			}
			if !v.AggWithinGroup && v.AggOrder != nil && len(v.AggOrder.Items) > 0 {
				b.WriteString(" ORDER BY ")
				b.WriteString(d.deparseSortClause(v.AggOrder))
			}
		}
		b.WriteString(")")
		if v.AggWithinGroup && v.AggOrder != nil && len(v.AggOrder.Items) > 0 {
			b.WriteString(" WITHIN GROUP (ORDER BY ")
			b.WriteString(d.deparseSortClause(v.AggOrder))
			b.WriteString(")")
		}
		if v.AggFilter != nil {
			b.WriteString(" FILTER (WHERE ")
			b.WriteString(d.deparseExpr(v.AggFilter))
			b.WriteString(")")
		}
		if wd, ok := v.Over.(*nodes.WindowDef); ok && wd != nil {
			if wd.Name != "" {
				b.WriteString(" OVER ")
				b.WriteString(quoteIdentifier(wd.Name))
			} else {
				b.WriteString(" OVER (")
				b.WriteString(d.deparseWindowSpec(wd))
				b.WriteString(")")
			}
		}
		return b.String()

	case *nodes.NamedArgExpr:
		return fmt.Sprintf("%s => %s", quoteIdentifier(v.Name), d.deparseExpr(v.Arg))

	case *nodes.CoalesceExpr:
		return fmt.Sprintf("COALESCE(%s)", d.deparseExprList(v.Args))

	case *nodes.MinMaxExpr:
		kw := "GREATEST"
		if v.Op == nodes.IS_LEAST {
			kw = "LEAST"
		}
		return fmt.Sprintf("%s(%s)", kw, d.deparseExprList(v.Args))

	case *nodes.CaseExpr:
		var b strings.Builder
		b.WriteString("CASE")
		if v.Arg != nil {
			b.WriteString(" ")
			b.WriteString(d.deparseExpr(v.Arg))
		}
		if v.Args != nil {
			for _, w := range v.Args.Items {
				if cw, ok := w.(*nodes.CaseWhen); ok {
					b.WriteString(" WHEN ")
					b.WriteString(d.deparseExpr(cw.Expr))
					b.WriteString(" THEN ")
					b.WriteString(d.deparseExpr(cw.Result))
				}
			}
		}
		if v.Defresult != nil {
			b.WriteString(" ELSE ")
			b.WriteString(d.deparseExpr(v.Defresult))
		}
		b.WriteString(" END")
		return b.String()

	case *nodes.A_ArrayExpr:
		return fmt.Sprintf("ARRAY[%s]", d.deparseExprList(v.Elements))

	case *nodes.RowExpr:
		if v.Args == nil || len(v.Args.Items) == 0 {
			return "ROW()"
		}
		return fmt.Sprintf("ROW(%s)", d.deparseExprList(v.Args))

	case *nodes.SubLink:
		sub := d.deparseActionStmt(v.Subselect)
		switch nodes.SubLinkType(v.SubLinkType) {
		case nodes.EXISTS_SUBLINK:
			return fmt.Sprintf("(EXISTS (%s))", sub)
		case nodes.ARRAY_SUBLINK:
			return fmt.Sprintf("ARRAY(%s)", sub)
		case nodes.ANY_SUBLINK:
			op := "="
			if v.OperName != nil {
				op = d.deparseOperatorRef(v.OperName)
			}
			if op == "=" {
				return fmt.Sprintf("(%s IN (%s))", d.deparseExpr(v.Testexpr), sub)
			}
			return fmt.Sprintf("(%s %s ANY (%s))", d.deparseExpr(v.Testexpr), op, sub)
		case nodes.ALL_SUBLINK:
			op := "="
			if v.OperName != nil {
				op = d.deparseOperatorRef(v.OperName)
			}
			return fmt.Sprintf("(%s %s ALL (%s))", d.deparseExpr(v.Testexpr), op, sub)
		default:
			return fmt.Sprintf("(%s)", sub)
		}

	case *nodes.SQLValueFunction:
		switch v.Op {
		case nodes.SVFOP_CURRENT_DATE:
			return "CURRENT_DATE"
		case nodes.SVFOP_CURRENT_TIME:
			return "CURRENT_TIME"
		case nodes.SVFOP_CURRENT_TIME_N:
			if v.Typmod >= 0 {
				return fmt.Sprintf("CURRENT_TIME(%d)", v.Typmod)
			}
			return "CURRENT_TIME"
		case nodes.SVFOP_CURRENT_TIMESTAMP:
			return "CURRENT_TIMESTAMP"
		case nodes.SVFOP_CURRENT_TIMESTAMP_N:
			if v.Typmod >= 0 {
				return fmt.Sprintf("CURRENT_TIMESTAMP(%d)", v.Typmod)
			}
			return "CURRENT_TIMESTAMP"
		case nodes.SVFOP_LOCALTIME:
			return "LOCALTIME"
		case nodes.SVFOP_LOCALTIME_N:
			if v.Typmod >= 0 {
				return fmt.Sprintf("LOCALTIME(%d)", v.Typmod)
			}
			return "LOCALTIME"
		case nodes.SVFOP_LOCALTIMESTAMP:
			return "LOCALTIMESTAMP"
		case nodes.SVFOP_LOCALTIMESTAMP_N:
			if v.Typmod >= 0 {
				return fmt.Sprintf("LOCALTIMESTAMP(%d)", v.Typmod)
			}
			return "LOCALTIMESTAMP"
		case nodes.SVFOP_CURRENT_ROLE:
			return "CURRENT_ROLE"
		case nodes.SVFOP_CURRENT_USER:
			return "CURRENT_USER"
		case nodes.SVFOP_USER:
			return "USER"
		case nodes.SVFOP_SESSION_USER:
			return "SESSION_USER"
		case nodes.SVFOP_CURRENT_CATALOG:
			return "CURRENT_CATALOG"
		case nodes.SVFOP_CURRENT_SCHEMA:
			return "CURRENT_SCHEMA"
		}

	case *nodes.ArrayExpr:
		return fmt.Sprintf("ARRAY[%s]", d.deparseExprList(v.Elements))

	case *nodes.GroupingFunc:
		return fmt.Sprintf("GROUPING(%s)", d.deparseExprList(v.Args))

	case *nodes.GroupingSet:
		switch v.Kind {
		case nodes.GROUPING_SET_EMPTY:
			return "()"
		case nodes.GROUPING_SET_SIMPLE:
			return "(" + d.deparseExprList(v.Content) + ")"
		case nodes.GROUPING_SET_ROLLUP:
			return "ROLLUP(" + d.deparseExprList(v.Content) + ")"
		case nodes.GROUPING_SET_CUBE:
			return "CUBE(" + d.deparseExprList(v.Content) + ")"
		case nodes.GROUPING_SET_SETS:
			return "GROUPING SETS (" + d.deparseExprList(v.Content) + ")"
		default:
			return "(" + d.deparseExprList(v.Content) + ")"
		}

	case *nodes.XmlExpr:
		return d.deparseXmlExpr(v)

	case *nodes.XmlSerialize:
		opt := "CONTENT"
		if v.Xmloption == nodes.XMLOPTION_DOCUMENT {
			opt = "DOCUMENT"
		}
		var b strings.Builder
		b.WriteString("XMLSERIALIZE(")
		b.WriteString(opt)
		b.WriteString(" ")
		b.WriteString(d.deparseExpr(v.Expr))
		b.WriteString(" AS ")
		b.WriteString(d.deparseTypeName(v.TypeName))
		if v.Indent {
			b.WriteString(" INDENT")
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonValueExpr:
		if v == nil {
			return ""
		}
		s := d.deparseExpr(v.RawExpr)
		if f := deparseJsonFormat(v.Format); f != "" {
			s += " " + f
		}
		return s

	case *nodes.JsonIsPredicate:
		var b strings.Builder
		b.WriteString("(")
		b.WriteString(d.deparseExpr(v.Expr))
		b.WriteString(" IS JSON")
		switch v.ItemType {
		case nodes.JS_TYPE_OBJECT:
			b.WriteString(" OBJECT")
		case nodes.JS_TYPE_ARRAY:
			b.WriteString(" ARRAY")
		case nodes.JS_TYPE_SCALAR:
			b.WriteString(" SCALAR")
		}
		if v.UniqueKeys {
			b.WriteString(" WITH UNIQUE KEYS")
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonParseExpr:
		var b strings.Builder
		b.WriteString("JSON(")
		b.WriteString(d.deparseExpr(v.Expr))
		if v.UniqueKeys {
			b.WriteString(" WITH UNIQUE KEYS")
		}
		if out := d.deparseJsonOutput(v.Output); out != "" {
			b.WriteString(" ")
			b.WriteString(out)
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonScalarExpr:
		var b strings.Builder
		b.WriteString("JSON_SCALAR(")
		b.WriteString(d.deparseExpr(v.Expr))
		if out := d.deparseJsonOutput(v.Output); out != "" {
			b.WriteString(" ")
			b.WriteString(out)
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonSerializeExpr:
		var b strings.Builder
		b.WriteString("JSON_SERIALIZE(")
		b.WriteString(d.deparseExpr(v.Expr))
		if out := d.deparseJsonOutput(v.Output); out != "" {
			b.WriteString(" ")
			b.WriteString(out)
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonObjectConstructor:
		var b strings.Builder
		b.WriteString("JSON_OBJECT(")
		if v.Exprs != nil {
			for i, item := range v.Exprs.Items {
				if kv, ok := item.(*nodes.JsonKeyValue); ok {
					if i > 0 {
						b.WriteString(", ")
					}
					b.WriteString(d.deparseExpr(kv.Key))
					b.WriteString(" : ")
					b.WriteString(d.deparseExpr(kv.Value))
				}
			}
		}
		if v.AbsentOnNull {
			b.WriteString(" ABSENT ON NULL")
		}
		if v.UniqueKeys {
			b.WriteString(" WITH UNIQUE KEYS")
		}
		if out := d.deparseJsonOutput(v.Output); out != "" {
			b.WriteString(" ")
			b.WriteString(out)
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonArrayConstructor:
		var b strings.Builder
		b.WriteString("JSON_ARRAY(")
		if v.Exprs != nil {
			for i, item := range v.Exprs.Items {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(d.deparseExpr(item))
			}
		}
		if v.AbsentOnNull {
			b.WriteString(" ABSENT ON NULL")
		}
		if out := d.deparseJsonOutput(v.Output); out != "" {
			b.WriteString(" ")
			b.WriteString(out)
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonArrayQueryConstructor:
		var b strings.Builder
		b.WriteString("JSON_ARRAY(")
		b.WriteString(d.deparseActionStmt(v.Query))
		if f := deparseJsonFormat(v.Format); f != "" {
			b.WriteString(" ")
			b.WriteString(f)
		}
		if out := d.deparseJsonOutput(v.Output); out != "" {
			b.WriteString(" ")
			b.WriteString(out)
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonFuncExpr:
		var b strings.Builder
		switch v.Op {
		case nodes.JSON_EXISTS_OP:
			b.WriteString("JSON_EXISTS(")
		case nodes.JSON_QUERY_OP:
			b.WriteString("JSON_QUERY(")
		default:
			b.WriteString("JSON_VALUE(")
		}
		b.WriteString(d.deparseExpr(v.ContextItem))
		if v.Pathspec != nil {
			b.WriteString(", ")
			b.WriteString(d.deparseExpr(v.Pathspec))
		}
		if v.Passing != nil && len(v.Passing.Items) > 0 {
			b.WriteString(" PASSING ")
			for i, item := range v.Passing.Items {
				if ja, ok := item.(*nodes.JsonArgument); ok {
					if i > 0 {
						b.WriteString(", ")
					}
					b.WriteString(d.deparseExpr(ja.Val))
					b.WriteString(" AS ")
					b.WriteString(quoteIdentifier(ja.Name))
				}
			}
		}
		if out := d.deparseJsonOutput(v.Output); out != "" {
			b.WriteString(" ")
			b.WriteString(out)
		}
		switch v.Wrapper {
		case nodes.JSW_NONE:
			b.WriteString(" WITHOUT ARRAY WRAPPER")
		case nodes.JSW_CONDITIONAL:
			b.WriteString(" WITH CONDITIONAL ARRAY WRAPPER")
		case nodes.JSW_UNCONDITIONAL:
			b.WriteString(" WITH UNCONDITIONAL ARRAY WRAPPER")
		}
		switch v.Quotes {
		case nodes.JS_QUOTES_KEEP:
			b.WriteString(" KEEP QUOTES")
		case nodes.JS_QUOTES_OMIT:
			b.WriteString(" OMIT QUOTES")
		}
		if v.OnEmpty != nil {
			b.WriteString(" ")
			b.WriteString(d.deparseJsonBehavior(v.OnEmpty))
			b.WriteString(" ON EMPTY")
		}
		if v.OnError != nil {
			b.WriteString(" ")
			b.WriteString(d.deparseJsonBehavior(v.OnError))
			b.WriteString(" ON ERROR")
		}
		b.WriteString(")")
		return b.String()

	case *nodes.JsonObjectAgg:
		var b strings.Builder
		b.WriteString("JSON_OBJECTAGG(")
		if v.Arg != nil {
			b.WriteString(d.deparseExpr(v.Arg.Key))
			b.WriteString(" : ")
			b.WriteString(d.deparseExpr(v.Arg.Value))
		}
		if v.AbsentOnNull {
			b.WriteString(" ABSENT ON NULL")
		}
		if v.UniqueKeys {
			b.WriteString(" WITH UNIQUE KEYS")
		}
		if v.Constructor != nil {
			if out := d.deparseJsonOutput(v.Constructor.Output); out != "" {
				b.WriteString(" ")
				b.WriteString(out)
			}
		}
		b.WriteString(")")
		if v.Constructor != nil {
			b.WriteString(d.deparseJsonAggSuffix(v.Constructor))
		}
		return b.String()

	case *nodes.JsonArrayAgg:
		var b strings.Builder
		b.WriteString("JSON_ARRAYAGG(")
		b.WriteString(d.deparseExpr(v.Arg))
		if v.Constructor != nil && v.Constructor.Agg_order != nil && len(v.Constructor.Agg_order.Items) > 0 {
			b.WriteString(" ORDER BY ")
			b.WriteString(d.deparseSortClause(v.Constructor.Agg_order))
		}
		if v.AbsentOnNull {
			b.WriteString(" ABSENT ON NULL")
		}
		if v.Constructor != nil {
			if out := d.deparseJsonOutput(v.Constructor.Output); out != "" {
				b.WriteString(" ")
				b.WriteString(out)
			}
		}
		b.WriteString(")")
		if v.Constructor != nil {
			b.WriteString(d.deparseJsonAggSuffix(v.Constructor))
		}
		return b.String()
	}
	return ""
}

func (d *ruleDeparser) deparseXmlExpr(v *nodes.XmlExpr) string {
	if v == nil {
		return ""
	}
	switch v.Op {
	case nodes.IS_XMLCONCAT:
		return fmt.Sprintf("XMLCONCAT(%s)", d.deparseExprList(v.Args))

	case nodes.IS_XMLELEMENT:
		var b strings.Builder
		b.WriteString("XMLELEMENT(NAME ")
		b.WriteString(quoteIdentifier(v.Name))
		if v.NamedArgs != nil && len(v.NamedArgs.Items) > 0 {
			b.WriteString(", XMLATTRIBUTES(")
			b.WriteString(d.deparseXmlNamedArgs(v.NamedArgs, v.ArgNames))
			b.WriteString(")")
		}
		if v.Args != nil && len(v.Args.Items) > 0 {
			b.WriteString(", ")
			b.WriteString(d.deparseExprList(v.Args))
		}
		b.WriteString(")")
		return b.String()

	case nodes.IS_XMLFOREST:
		return fmt.Sprintf("XMLFOREST(%s)", d.deparseXmlNamedArgs(v.NamedArgs, v.ArgNames))

	case nodes.IS_XMLPARSE:
		opt := "CONTENT"
		if v.Xmloption == nodes.XMLOPTION_DOCUMENT {
			opt = "DOCUMENT"
		}
		var argSQL string
		preserveWS := false
		if v.Args != nil && len(v.Args.Items) > 0 {
			argSQL = d.deparseExpr(v.Args.Items[0])
			if len(v.Args.Items) > 1 {
				if ac, ok := v.Args.Items[1].(*nodes.A_Const); ok && !ac.Isnull {
					if bv, ok := ac.Val.(*nodes.Boolean); ok && bv.Boolval {
						preserveWS = true
					}
				}
			}
		}
		if preserveWS {
			return fmt.Sprintf("XMLPARSE(%s %s PRESERVE WHITESPACE)", opt, argSQL)
		}
		return fmt.Sprintf("XMLPARSE(%s %s)", opt, argSQL)

	case nodes.IS_XMLPI:
		var b strings.Builder
		b.WriteString("XMLPI(NAME ")
		b.WriteString(quoteIdentifier(v.Name))
		if v.Args != nil && len(v.Args.Items) > 0 {
			b.WriteString(", ")
			b.WriteString(d.deparseExpr(v.Args.Items[0]))
		}
		b.WriteString(")")
		return b.String()

	case nodes.IS_XMLROOT:
		var b strings.Builder
		b.WriteString("XMLROOT(")
		if v.Args != nil && len(v.Args.Items) > 0 {
			b.WriteString(d.deparseExpr(v.Args.Items[0]))
		}
		if v.Args != nil && len(v.Args.Items) > 1 {
			b.WriteString(", VERSION ")
			if ac, ok := v.Args.Items[1].(*nodes.A_Const); ok && ac.Isnull {
				b.WriteString("NO VALUE")
			} else {
				b.WriteString(d.deparseExpr(v.Args.Items[1]))
			}
		}
		if v.Args != nil && len(v.Args.Items) > 2 {
			if ac, ok := v.Args.Items[2].(*nodes.A_Const); ok && !ac.Isnull {
				if iv, ok := ac.Val.(*nodes.Integer); ok {
					switch iv.Ival {
					case nodes.XML_STANDALONE_YES:
						b.WriteString(", STANDALONE YES")
					case nodes.XML_STANDALONE_NO:
						b.WriteString(", STANDALONE NO")
					case nodes.XML_STANDALONE_NO_VALUE:
						b.WriteString(", STANDALONE NO VALUE")
					}
				}
			}
		}
		b.WriteString(")")
		return b.String()

	case nodes.IS_XMLSERIALIZE:
		opt := "CONTENT"
		if v.Xmloption == nodes.XMLOPTION_DOCUMENT {
			opt = "DOCUMENT"
		}
		var argSQL string
		if v.Args != nil && len(v.Args.Items) > 0 {
			argSQL = d.deparseExpr(v.Args.Items[0])
		}
		typSQL := "text"
		if d.cat != nil && v.Type != 0 {
			typSQL = d.cat.FormatType(uint32(v.Type), v.Typmod)
		}
		if v.Indent {
			return fmt.Sprintf("XMLSERIALIZE(%s %s AS %s INDENT)", opt, argSQL, typSQL)
		}
		return fmt.Sprintf("XMLSERIALIZE(%s %s AS %s)", opt, argSQL, typSQL)

	case nodes.IS_DOCUMENT:
		var argSQL string
		if v.Args != nil && len(v.Args.Items) > 0 {
			argSQL = d.deparseExpr(v.Args.Items[0])
		}
		return fmt.Sprintf("(%s IS DOCUMENT)", argSQL)
	}
	return ""
}

func (d *ruleDeparser) deparseXmlNamedArgs(namedArgs, argNames *nodes.List) string {
	if namedArgs == nil || len(namedArgs.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(namedArgs.Items))
	for i, item := range namedArgs.Items {
		if rt, ok := item.(*nodes.ResTarget); ok {
			s := d.deparseExpr(rt.Val)
			if rt.Name != "" {
				s += " AS " + quoteIdentifier(rt.Name)
			}
			parts = append(parts, s)
			continue
		}
		s := d.deparseExpr(item)
		if argNames != nil && i < len(argNames.Items) {
			if name := stringVal(argNames.Items[i]); name != "" {
				s += " AS " + quoteIdentifier(name)
			}
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ", ")
}

func deparseJsonFormat(f *nodes.JsonFormat) string {
	if f == nil || f.FormatType == nodes.JS_FORMAT_DEFAULT {
		return ""
	}
	var b strings.Builder
	b.WriteString("FORMAT ")
	if f.FormatType == nodes.JS_FORMAT_JSONB {
		b.WriteString("JSONB")
	} else {
		b.WriteString("JSON")
	}
	switch f.Encoding {
	case nodes.JS_ENC_UTF8:
		b.WriteString(" ENCODING UTF8")
	case nodes.JS_ENC_UTF16:
		b.WriteString(" ENCODING UTF16")
	case nodes.JS_ENC_UTF32:
		b.WriteString(" ENCODING UTF32")
	}
	return b.String()
}

func (d *ruleDeparser) deparseJsonOutput(out *nodes.JsonOutput) string {
	if out == nil || out.TypeName == nil {
		return ""
	}
	s := "RETURNING " + d.deparseTypeName(out.TypeName)
	if out.Returning != nil {
		if f := deparseJsonFormat(out.Returning.Format); f != "" {
			s += " " + f
		}
	}
	return s
}

func (d *ruleDeparser) deparseJsonBehavior(jb *nodes.JsonBehavior) string {
	if jb == nil {
		return "NULL"
	}
	switch jb.Btype {
	case nodes.JSON_BEHAVIOR_ERROR:
		return "ERROR"
	case nodes.JSON_BEHAVIOR_EMPTY:
		return "EMPTY"
	case nodes.JSON_BEHAVIOR_TRUE:
		return "TRUE"
	case nodes.JSON_BEHAVIOR_FALSE:
		return "FALSE"
	case nodes.JSON_BEHAVIOR_UNKNOWN:
		return "UNKNOWN"
	case nodes.JSON_BEHAVIOR_EMPTY_ARRAY:
		return "EMPTY ARRAY"
	case nodes.JSON_BEHAVIOR_EMPTY_OBJECT:
		return "EMPTY OBJECT"
	case nodes.JSON_BEHAVIOR_DEFAULT:
		return "DEFAULT " + d.deparseExpr(jb.Expr)
	default:
		return "NULL"
	}
}

func (d *ruleDeparser) deparseJsonAggSuffix(ac *nodes.JsonAggConstructor) string {
	if ac == nil {
		return ""
	}
	var b strings.Builder
	if ac.Agg_filter != nil {
		b.WriteString(" FILTER (WHERE ")
		b.WriteString(d.deparseExpr(ac.Agg_filter))
		b.WriteString(")")
	}
	if ac.Over != nil {
		if ac.Over.Name != "" {
			b.WriteString(" OVER ")
			b.WriteString(quoteIdentifier(ac.Over.Name))
		} else {
			b.WriteString(" OVER (")
			b.WriteString(d.deparseWindowSpec(ac.Over))
			b.WriteString(")")
		}
	}
	return b.String()
}

func deparseRawConstVal(n nodes.Node) string {
	switch v := n.(type) {
	case *nodes.Integer:
		return fmt.Sprintf("%d", v.Ival)
	case *nodes.Float:
		return v.Fval
	case *nodes.String:
		return quoteRuleLiteral(v.Str)
	case *nodes.Boolean:
		if v.Boolval {
			return "true"
		}
		return "false"
	case *nodes.BitString:
		if strings.HasPrefix(v.Bsval, "x") || strings.HasPrefix(v.Bsval, "X") {
			return "X'" + v.Bsval[1:] + "'"
		}
		if strings.HasPrefix(v.Bsval, "b") || strings.HasPrefix(v.Bsval, "B") {
			return "B'" + v.Bsval[1:] + "'"
		}
		return "B'" + v.Bsval + "'"
	default:
		return "NULL"
	}
}

func quoteRuleLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (d *ruleDeparser) deparseOperatorRef(list *nodes.List) string {
	if list == nil || len(list.Items) == 0 {
		return ""
	}
	if d.cat != nil && d.b != nil && d.b.operByList != nil {
		if opOID, ok := d.b.operByList[list]; ok && opOID >= FirstNormalObjectId {
			if op := d.cat.operByOID[opOID]; op != nil {
				nsOID := d.cat.operatorSchemas[opOID]
				if sc := d.cat.schemas[nsOID]; sc != nil && (len(list.Items) > 1 || (sc.Name != "public" && sc.Name != "pg_catalog")) {
					return fmt.Sprintf("OPERATOR(%s.%s)", quoteIdentifier(sc.Name), op.Name)
				}
				return op.Name
			}
		}
	}
	return deparseOperatorName(list)
}

func deparseOperatorName(list *nodes.List) string {
	if list == nil || len(list.Items) == 0 {
		return ""
	}
	if len(list.Items) == 1 {
		return stringVal(list.Items[0])
	}
	return fmt.Sprintf("OPERATOR(%s.%s)", quoteIdentifier(stringVal(list.Items[0])), stringVal(list.Items[1]))
}

func deparseQualifiedIdentList(list *nodes.List) string {
	if list == nil || len(list.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		parts = append(parts, quoteIdentifier(stringVal(item)))
	}
	return strings.Join(parts, ".")
}

func deparseFuncNameList(list *nodes.List) string {
	if list == nil || len(list.Items) == 0 {
		return ""
	}
	if len(list.Items) == 2 && stringVal(list.Items[0]) == "pg_catalog" {
		return quoteIdentifier(stringVal(list.Items[1]))
	}
	return deparseQualifiedIdentList(list)
}

func (d *ruleDeparser) deparseTypeName(tn *nodes.TypeName) string {
	if tn == nil {
		return ""
	}
	if d.cat != nil && d.b != nil {
		if tb, ok := d.b.typeByTN[tn]; ok && tb.typeOID >= FirstNormalObjectId {
			if bt := d.cat.typeByOID[tb.typeOID]; bt != nil {
				var base string
				if sc := d.cat.schemas[bt.Namespace]; sc != nil && sc.Name != "public" && sc.Name != "pg_catalog" {
					base = quoteIdentifier(sc.Name) + "." + quoteIdentifier(bt.TypeName)
				} else {
					base = quoteIdentifier(bt.TypeName)
				}
				if tn.Typmods != nil && len(tn.Typmods.Items) > 0 {
					base += "(" + d.deparseExprList(tn.Typmods) + ")"
				}
				if tb.isArray && tn.ArrayBounds != nil {
					for range tn.ArrayBounds.Items {
						base += "[]"
					}
				}
				return base
			}
		}
	}
	schema, name := typeNameParts(tn)
	var base string
	if schema == "" || schema == "pg_catalog" {
		switch strings.ToLower(name) {
		case "int2":
			base = "smallint"
		case "int4":
			base = "integer"
		case "int8":
			base = "bigint"
		case "float4":
			base = "real"
		case "float8":
			base = "double precision"
		case "bool":
			base = "boolean"
		case "bpchar":
			base = "character"
		case "varchar":
			base = "character varying"
		case "varbit":
			base = "bit varying"
		case "timetz":
			base = "time with time zone"
		case "timestamptz":
			base = "timestamp with time zone"
		default:
			base = quoteIdentifier(name)
		}
	} else {
		base = quoteIdentifier(schema) + "." + quoteIdentifier(name)
	}
	if tn.Typmods != nil && len(tn.Typmods.Items) > 0 {
		base += "(" + d.deparseExprList(tn.Typmods) + ")"
	}
	if tn.ArrayBounds != nil {
		for range tn.ArrayBounds.Items {
			base += "[]"
		}
	}
	return base
}
