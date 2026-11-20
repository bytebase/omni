package catalog

import (
	"fmt"
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

type ruleBindings struct {
	ownerRelOID    uint32
	isOnSelect     bool
	allowOld       bool
	allowNew       bool
	relByRV        map[*nodes.RangeVar]uint32
	colByRef       map[*nodes.ColumnRef]ruleColBinding
	colByRT        map[*nodes.ResTarget]ruleColBinding
	targetNameByRT map[*nodes.ResTarget]string
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
	joinUsing      map[*nodes.JoinExpr][]ruleJoinUsingCol
	naturalJoined  map[*nodes.JoinExpr]bool
	stmtCtx        map[*nodes.SelectStmt]*ruleEvalContext

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
		relByRV:        make(map[*nodes.RangeVar]uint32, len(b.relByRV)),
		colByRef:       make(map[*nodes.ColumnRef]ruleColBinding, len(b.colByRef)),
		colByRT:        make(map[*nodes.ResTarget]ruleColBinding, len(b.colByRT)),
		targetNameByRT: make(map[*nodes.ResTarget]string, len(b.targetNameByRT)),
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
	for k, v := range b.colByRef {
		cp.colByRef[k] = v
	}
	for k, v := range b.colByRT {
		cp.colByRT[k] = v
	}
	for k, v := range b.targetNameByRT {
		cp.targetNameByRT[k] = v
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
	if !s.Replace && rel.Rules != nil && rel.Rules[s.Rulename] != nil {
		return errDuplicateObject("rule", s.Rulename)
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
			if te.Expr.exprType() != col.TypeOID || te.Expr.exprTypMod() != col.TypeMod {
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
	return nil
}

func newRuleBindings(ownerRelOID uint32, event nodes.CmdType) *ruleBindings {
	return &ruleBindings{
		ownerRelOID:    ownerRelOID,
		isOnSelect:     event == nodes.CMD_SELECT,
		allowOld:       event == nodes.CMD_UPDATE || event == nodes.CMD_DELETE,
		allowNew:       event == nodes.CMD_INSERT || event == nodes.CMD_UPDATE,
		relByRV:        make(map[*nodes.RangeVar]uint32),
		colByRef:       make(map[*nodes.ColumnRef]ruleColBinding),
		colByRT:        make(map[*nodes.ResTarget]ruleColBinding),
		targetNameByRT: make(map[*nodes.ResTarget]string),
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
		joinUsing:      make(map[*nodes.JoinExpr][]ruleJoinUsingCol),
		naturalJoined:  make(map[*nodes.JoinExpr]bool),
		stmtCtx:        make(map[*nodes.SelectStmt]*ruleEvalContext),
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

func (c *Catalog) validatePreparedDMLStmt(n nodes.Node, paramTypes []uint32) error {
	b := newRuleBindings(0, 0)
	if err := c.bindRuleStmt(b, n, &ruleEvalContext{}); err != nil {
		return err
	}
	c.inferPreparedDMLParamTypes(b, n, paramTypes)
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
	virtualColOrder      []string
	anyCol               bool
	isJoinAlias          bool
	hasJoinColAliases    bool
	hiddenByJoinColAlias bool
	joinRels             []uint32
}

type ruleEvalContext struct {
	scope         map[string]ruleScopeEntry
	outer         *ruleEvalContext
	cteNames      map[string]bool
	cteCols       map[string][]string
	cteAny        map[string]bool
	usingCols     map[string]int
	targetAliases map[string]bool
	hasOpenSource bool
}

func (c *Catalog) buildRuleBindings(ownerRel *Relation, event nodes.CmdType, whereClause nodes.Node, actions *nodes.List) (*ruleBindings, error) {
	b := newRuleBindings(ownerRel.OID, event)
	if whereClause != nil {
		if err := c.bindRuleExpr(b, whereClause, &ruleEvalContext{}); err != nil {
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

func applyVirtualAliasCols(baseCols []string, baseAny bool, aliasCols []string) (map[string]int, []string, bool) {
	if len(baseCols) == 0 && len(aliasCols) == 0 {
		return nil, nil, baseAny
	}
	n := len(baseCols)
	if len(aliasCols) > n {
		n = len(aliasCols)
	}
	outMap := make(map[string]int, n)
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
		}
		outOrder = append(outOrder, name)
	}
	return outMap, outOrder, baseAny
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
				b.addSeqDep(seq.OID)
				return 0, "seq", nil
			}
		}
		if c.snapshot {
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
		baseAny := false
		switch virtualKind {
		case "cte":
			baseCols = ctx.cteCols[rv.Relname]
			baseAny = ctx.cteAny[rv.Relname]
		case "seq":
			baseCols = []string{"last_value", "log_cnt", "is_called"}
		default:
			baseAny = true
		}
		vCols, vOrder, vAny := applyVirtualAliasCols(baseCols, baseAny, extractAliasColnames(rv.Alias))
		ctx.scope[name] = ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           alias,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}
		if vAny {
			ctx.hasOpenSource = true
		}
		return 0, nil
	}
	if relOID != 0 {
		ctx.scope[name] = ruleScopeEntry{
			relOID:     relOID,
			hasAlias:   hasAlias,
			alias:      alias,
			colAliases: extractAliasColnames(rv.Alias),
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
		subCtx := &ruleEvalContext{cteNames: ctx.cteNames, cteCols: ctx.cteCols, cteAny: ctx.cteAny}
		if v.Lateral {
			subCtx.outer = ctx
		}
		if err := c.bindRuleStmt(b, v.Subquery, subCtx); err != nil {
			return err
		}
		baseCols, baseAny := c.extractRuleQueryCols(b, v.Subquery)
		vCols, vOrder, vAny := applyVirtualAliasCols(baseCols, baseAny, extractAliasColnames(v.Alias))
		if v.Alias != nil && v.Alias.Aliasname != "" {
			ctx.scope[v.Alias.Aliasname] = ruleScopeEntry{
				hasAlias:        true,
				alias:           v.Alias.Aliasname,
				isVirtual:       true,
				virtualCols:     vCols,
				virtualColOrder: vOrder,
				anyCol:          vAny,
			}
		}
		if vAny {
			ctx.hasOpenSource = true
		}
	case *nodes.RangeFunction:
		if v.Functions != nil {
			fnCtx := ctx
			if !v.Lateral {
				fnCtx = &ruleEvalContext{outer: ctx.outer, cteNames: ctx.cteNames, cteCols: ctx.cteCols, cteAny: ctx.cteAny}
			}
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
		vCols, vOrder, vAny := applyVirtualAliasCols(baseCols, baseAny, extractAliasColnames(v.Alias))
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
		if scopeName == "" {
			scopeName = fmt.Sprintf("\x00rowsfrom_%d", len(ctx.scope))
		}
		ctx.scope[scopeName] = ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           aliasName,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}
		if vAny {
			ctx.hasOpenSource = true
		}
	case *nodes.RangeTableFunc:
		fnCtx := ctx
		if !v.Lateral {
			fnCtx = &ruleEvalContext{outer: ctx.outer, cteNames: ctx.cteNames, cteCols: ctx.cteCols, cteAny: ctx.cteAny}
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
		if v.Columns != nil {
			for _, item := range v.Columns.Items {
				fc, ok := item.(*nodes.RangeTableFuncCol)
				if !ok || fc == nil {
					continue
				}
				if fc.Colname != "" {
					baseCols = append(baseCols, fc.Colname)
				}
				if fc.TypeName != nil {
					if err := c.bindRuleTypeName(b, fc.TypeName); err != nil {
						return err
					}
				}
				if err := c.bindRuleExpr(b, fc.Colexpr, fnCtx); err != nil {
					return err
				}
				if err := c.bindRuleExpr(b, fc.Coldefexpr, fnCtx); err != nil {
					return err
				}
			}
		}
		vCols, vOrder, vAny := applyVirtualAliasCols(baseCols, false, extractAliasColnames(v.Alias))
		scopeName, hasAlias, aliasName := "xmltable", false, ""
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName, hasAlias, aliasName = v.Alias.Aliasname, true, v.Alias.Aliasname
		}
		ctx.scope[scopeName] = ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           aliasName,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}
	case *nodes.JsonTable:
		fnCtx := ctx
		if !v.Lateral {
			fnCtx = &ruleEvalContext{outer: ctx.outer, cteNames: ctx.cteNames, cteCols: ctx.cteCols, cteAny: ctx.cteAny}
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
		vCols, vOrder, vAny := applyVirtualAliasCols(baseCols, false, extractAliasColnames(v.Alias))
		scopeName, hasAlias, aliasName := "json_table", false, ""
		if v.Alias != nil && v.Alias.Aliasname != "" {
			scopeName, hasAlias, aliasName = v.Alias.Aliasname, true, v.Alias.Aliasname
		}
		ctx.scope[scopeName] = ruleScopeEntry{
			hasAlias:        hasAlias,
			alias:           aliasName,
			isVirtual:       true,
			virtualCols:     vCols,
			virtualColOrder: vOrder,
			anyCol:          vAny,
		}
	case *nodes.RangeTableSample:
		if err := c.collectTableRefScope(b, v.Relation, ctx); err != nil {
			return err
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
				if len(aliasCols) == 0 && !seenUsing[uCol] {
					seenUsing[uCol] = true
					ctx.usingCols[uCol]++
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
				ctx.scope[v.JoinUsing.Aliasname] = ruleScopeEntry{
					hasAlias:        true,
					alias:           v.JoinUsing.Aliasname,
					isVirtual:       true,
					virtualCols:     juCols,
					virtualColOrder: juOrder,
					isJoinAlias:     true,
					joinRels:        append(append([]uint32(nil), lRels...), rRels...),
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
			if err := c.bindRuleExpr(b, v.Quals, ctx); err != nil {
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
			vCols, vOrder, vAny := applyVirtualAliasCols(combinedCols, lAny || rAny, aliasCols)
			joinRels := append(append([]uint32(nil), lRels...), rRels...)
			ctx.scope[v.Alias.Aliasname] = ruleScopeEntry{
				hasAlias:          true,
				alias:             v.Alias.Aliasname,
				isVirtual:         true,
				virtualCols:       vCols,
				virtualColOrder:   vOrder,
				anyCol:            vAny,
				isJoinAlias:       true,
				hasJoinColAliases: len(aliasCols) > 0,
				joinRels:          joinRels,
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
		searchNS = c.searchPathWithCatalog()
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

func (c *Catalog) bindWithClause(b *ruleBindings, w *nodes.WithClause, parent *ruleEvalContext) (map[string]bool, map[string][]string, map[string]bool, error) {
	cteNames := make(map[string]bool)
	cteCols := make(map[string][]string)
	cteAny := make(map[string]bool)
	if parent != nil {
		for k, v := range parent.cteNames {
			cteNames[k] = v
		}
		for k, v := range parent.cteCols {
			cteCols[k] = v
		}
		for k, v := range parent.cteAny {
			cteAny[k] = v
		}
	}
	if w == nil || w.Ctes == nil {
		return cteNames, cteCols, cteAny, nil
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
	for _, item := range w.Ctes.Items {
		if cte, ok := item.(*nodes.CommonTableExpr); ok {
			var outer *ruleEvalContext
			if parent != nil {
				outer = parent
			}
			cteCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteAny: cteAny}
			if err := c.bindRuleStmt(b, cte.Ctequery, cteCtx); err != nil {
				return nil, nil, nil, err
			}
			if cc, ok := cte.CycleClause.(*nodes.CTECycleClause); ok && cc != nil {
				if err := c.bindRuleExpr(b, cc.CycleMarkValue, cteCtx); err != nil {
					return nil, nil, nil, err
				}
				if err := c.bindRuleExpr(b, cc.CycleMarkDefault, cteCtx); err != nil {
					return nil, nil, nil, err
				}
			}
			if cte.Ctename != "" {
				cteNames[cte.Ctename] = true
				baseCols, baseAny := c.extractRuleQueryCols(b, cte.Ctequery)
				_, vOrder, vAny := applyVirtualAliasCols(baseCols, baseAny, extractNodeListStrings(cte.Aliascolnames))
				if sc, ok := cte.SearchClause.(*nodes.CTESearchClause); ok && sc != nil && sc.SearchSeqColumn != "" {
					vOrder = append(vOrder, sc.SearchSeqColumn)
				}
				if cc, ok := cte.CycleClause.(*nodes.CTECycleClause); ok && cc != nil {
					if cc.CycleMarkColumn != "" {
						vOrder = append(vOrder, cc.CycleMarkColumn)
					}
					if cc.CyclePathColumn != "" {
						vOrder = append(vOrder, cc.CyclePathColumn)
					}
				}
				cteCols[cte.Ctename] = vOrder
				cteAny[cte.Ctename] = vAny
			}
		}
	}
	return cteNames, cteCols, cteAny, nil
}

func (c *Catalog) bindTargetColumns(b *ruleBindings, relOID uint32, list *nodes.List) error {
	if relOID == 0 || list == nil {
		return nil
	}
	rel := c.relationByOID[relOID]
	if rel == nil {
		return nil
	}
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
		b.colByRT[rt] = ruleColBinding{relOID: relOID, attNum: attNum}
		b.addColDep(relOID, attNum)
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
		maxCols = len(rel.Columns)
	} else if sel.TargetList != nil {
		hasStar := false
		for _, item := range sel.TargetList.Items {
			if rt, ok := item.(*nodes.ResTarget); ok {
				if cr, ok := rt.Val.(*nodes.ColumnRef); ok && cr.Fields != nil && len(cr.Fields.Items) > 0 {
					if _, isStar := cr.Fields.Items[len(cr.Fields.Items)-1].(*nodes.A_Star); isStar {
						hasStar = true
						break
					}
				}
			}
		}
		if hasStar {
			maxCols = len(rel.Columns)
		} else {
			if len(sel.TargetList.Items) > len(rel.Columns) && !c.snapshot {
				return &Error{
					Code:    CodeSyntaxError,
					Message: "INSERT has more expressions than target columns",
				}
			}
			maxCols = len(sel.TargetList.Items)
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
		cteNames, cteCols, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
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
			branchCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteAny: cteAny}
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
			sortCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteAny: cteAny, hasOpenSource: true}
			return c.bindSortAndLimit(b, s, sortCtx)
		}
		if s.ValuesLists != nil {
			valCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteAny: cteAny}
			for _, rowNode := range s.ValuesLists.Items {
				if row, ok := rowNode.(*nodes.List); ok {
					for _, expr := range row.Items {
						if err := c.bindRuleExpr(b, expr, valCtx); err != nil {
							return err
						}
					}
				}
			}
			sortCtx := &ruleEvalContext{outer: outer, cteNames: cteNames, cteCols: cteCols, cteAny: cteAny, hasOpenSource: true}
			return c.bindSortAndLimit(b, s, sortCtx)
		}
		ctx := &ruleEvalContext{
			scope:    make(map[string]ruleScopeEntry),
			outer:    outer,
			cteNames: cteNames,
			cteCols:  cteCols,
			cteAny:   cteAny,
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
				} else if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		if err := c.bindRuleExpr(b, s.WhereClause, ctx); err != nil {
			return err
		}
		aliasCtx := &ruleEvalContext{
			scope:         ctx.scope,
			outer:         ctx.outer,
			cteNames:      ctx.cteNames,
			cteCols:       ctx.cteCols,
			cteAny:        ctx.cteAny,
			usingCols:     ctx.usingCols,
			targetAliases: targetAliases,
			hasOpenSource: ctx.hasOpenSource,
		}
		if s.DistinctClause != nil {
			for _, item := range s.DistinctClause.Items {
				if err := c.bindRuleExpr(b, item, aliasCtx); err != nil {
					return err
				}
			}
		}
		if s.GroupClause != nil {
			for _, item := range s.GroupClause.Items {
				if err := c.bindRuleExpr(b, item, aliasCtx); err != nil {
					return err
				}
			}
		}
		if err := c.bindRuleExpr(b, s.HavingClause, ctx); err != nil {
			return err
		}
		if s.WindowClause != nil {
			for _, item := range s.WindowClause.Items {
				if wd, ok := item.(*nodes.WindowDef); ok {
					if err := c.bindWindowDef(b, wd, ctx); err != nil {
						return err
					}
				}
			}
		}
		return c.bindSortAndLimit(b, s, aliasCtx)

	case *nodes.InsertStmt:
		cteNames, cteCols, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
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
			if err := c.bindRuleStmt(b, s.SelectStmt, &ruleEvalContext{cteNames: cteNames, cteCols: cteCols, cteAny: cteAny}); err != nil {
				return err
			}
			if s.Cols == nil {
				if err := recordPositionalInsertTargetDeps(c, b, targetOID, s.SelectStmt); err != nil {
					return err
				}
			} else if targetOID != 0 {
				if sel, ok := s.SelectStmt.(*nodes.SelectStmt); ok && sel.ValuesLists != nil {
					for _, rowNode := range sel.ValuesLists.Items {
						if row, ok := rowNode.(*nodes.List); ok && len(row.Items) > len(s.Cols.Items) && !c.snapshot {
							return &Error{
								Code:    CodeSyntaxError,
								Message: "INSERT has more expressions than target columns",
							}
						}
					}
				}
			}
		}
		ctx := &ruleEvalContext{
			scope:    make(map[string]ruleScopeEntry),
			cteNames: cteNames,
			cteCols:  cteCols,
			cteAny:   cteAny,
		}
		if targetOID != 0 && s.Relation != nil {
			if s.Relation.Alias != nil && s.Relation.Alias.Aliasname != "" {
				ctx.scope[s.Relation.Alias.Aliasname] = ruleScopeEntry{relOID: targetOID, hasAlias: true, alias: s.Relation.Alias.Aliasname}
			} else {
				ctx.scope[s.Relation.Relname] = ruleScopeEntry{relOID: targetOID}
			}
			ctx.scope["excluded"] = ruleScopeEntry{relOID: targetOID, hasAlias: true, alias: "excluded"}
		}
		if s.OnConflictClause != nil {
			oc := s.OnConflictClause
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
				if err := c.bindRuleExpr(b, oc.Infer.WhereClause, ctx); err != nil {
					return err
				}
			}
			if err := c.bindTargetColumns(b, targetOID, oc.TargetList); err != nil {
				return err
			}
			if oc.TargetList != nil {
				for _, item := range oc.TargetList.Items {
					if rt, ok := item.(*nodes.ResTarget); ok {
						if err := c.bindRuleExpr(b, rt.Val, ctx); err != nil {
							return err
						}
					}
				}
			}
			if err := c.bindRuleExpr(b, oc.WhereClause, ctx); err != nil {
				return err
			}
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, ctx); err != nil {
				return err
			}
		}
		return nil

	case *nodes.UpdateStmt:
		cteNames, cteCols, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		ctx := &ruleEvalContext{
			scope:    make(map[string]ruleScopeEntry),
			cteNames: cteNames,
			cteCols:  cteCols,
			cteAny:   cteAny,
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
					if err := c.bindRuleExpr(b, rt.Val, ctx); err != nil {
						return err
					}
				}
			}
		}
		if err := c.bindRuleExpr(b, s.WhereClause, ctx); err != nil {
			return err
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, ctx); err != nil {
				return err
			}
		}
		return nil

	case *nodes.DeleteStmt:
		cteNames, cteCols, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		ctx := &ruleEvalContext{
			scope:    make(map[string]ruleScopeEntry),
			cteNames: cteNames,
			cteCols:  cteCols,
			cteAny:   cteAny,
		}
		if _, err := addScopeRangeVar(c, b, ctx, s.Relation, true); err != nil {
			return err
		}
		if err := c.collectFromScope(b, s.UsingClause, ctx); err != nil {
			return err
		}
		if err := c.bindRuleExpr(b, s.WhereClause, ctx); err != nil {
			return err
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, ctx); err != nil {
				return err
			}
		}
		return nil

	case *nodes.MergeStmt:
		cteNames, cteCols, cteAny, err := c.bindWithClause(b, s.WithClause, parent)
		if err != nil {
			return err
		}
		ctx := &ruleEvalContext{
			scope:    make(map[string]ruleScopeEntry),
			cteNames: cteNames,
			cteCols:  cteCols,
			cteAny:   cteAny,
		}
		targetOID, err := addScopeRangeVar(c, b, ctx, s.Relation, true)
		if err != nil {
			return err
		}
		if err := c.collectTableRefScope(b, s.SourceRelation, ctx); err != nil {
			return err
		}
		if err := c.bindRuleExpr(b, s.JoinCondition, ctx); err != nil {
			return err
		}
		if s.MergeWhenClauses != nil {
			for _, item := range s.MergeWhenClauses.Items {
				mwc, ok := item.(*nodes.MergeWhenClause)
				if !ok || mwc == nil {
					continue
				}
				if err := c.bindRuleExpr(b, mwc.Condition, ctx); err != nil {
					return err
				}
				if mwc.CommandType == nodes.CMD_INSERT {
					if err := c.bindTargetColumns(b, targetOID, mwc.TargetList); err != nil {
						return err
					}
					if mwc.Values != nil {
						for _, valExpr := range mwc.Values.Items {
							if err := c.bindRuleExpr(b, valExpr, ctx); err != nil {
								return err
							}
						}
						if targetOID != 0 && !c.snapshot {
							maxCols := 0
							if mwc.TargetList != nil && len(mwc.TargetList.Items) > 0 {
								maxCols = len(mwc.TargetList.Items)
							} else if rel := c.relationByOID[targetOID]; rel != nil {
								maxCols = len(rel.Columns)
							}
							if maxCols > 0 && len(mwc.Values.Items) > maxCols {
								return &Error{
									Code:    CodeSyntaxError,
									Message: "INSERT has more expressions than target columns",
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
								if err := c.bindRuleExpr(b, rt.Val, ctx); err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}
		if s.ReturningList != nil {
			if err := c.bindAndValidateReturningList(b, s.ReturningList, ctx); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (c *Catalog) bindAndValidateReturningList(b *ruleBindings, list *nodes.List, ctx *ruleEvalContext) error {
	if list == nil {
		return nil
	}
	expandedCols := 0
	unknownCount := false
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok {
			continue
		}
		if err := c.bindRuleExpr(b, rt.Val, ctx); err != nil {
			return err
		}
		if cr, isCR := rt.Val.(*nodes.ColumnRef); isCR && cr.Fields != nil && len(cr.Fields.Items) > 0 {
			if _, isStar := cr.Fields.Items[len(cr.Fields.Items)-1].(*nodes.A_Star); isStar {
				if len(cr.Fields.Items) == 1 {
					if ctx != nil && ctx.hasOpenSource {
						unknownCount = true
					} else if ctx != nil {
						seenRel := map[uint32]bool{}
						for k, se := range ctx.scope {
							if k == "excluded" || se.isVirtual || seenRel[se.relOID] {
								continue
							}
							seenRel[se.relOID] = true
							if rel := c.relationByOID[se.relOID]; rel != nil {
								expandedCols += len(rel.Columns)
							}
						}
					}
				} else {
					qual := stringVal(cr.Fields.Items[len(cr.Fields.Items)-2])
					if (strings.EqualFold(qual, "old") && b.allowOld) || (strings.EqualFold(qual, "new") && b.allowNew) {
						if ownerRel := c.relationByOID[b.ownerRelOID]; ownerRel != nil {
							expandedCols += len(ownerRel.Columns)
						}
					} else if se, ok := lookupScopeEntry(ctx, qual); ok {
						if se.isVirtual {
							unknownCount = true
						} else if rel := c.relationByOID[se.relOID]; rel != nil {
							expandedCols += len(rel.Columns)
						}
					}
				}
				continue
			}
		}
		expandedCols++
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
		}
	}
	return nil
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
			if err := c.bindRuleOperatorName(b, sb.UseOp, true); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Catalog) bindRuleOperatorName(b *ruleBindings, opNameList *nodes.List, requireExists bool) error {
	if opNameList == nil || len(opNameList.Items) == 0 {
		return nil
	}
	opSchema, opName := qualifiedName(opNameList)
	if opName == "" {
		return nil
	}
	if opSchema != "" && c.schemaByName[opSchema] == nil && !c.snapshot {
		return errUndefinedSchema(opSchema)
	}
	var searchNS []uint32
	if opSchema != "" {
		if s := c.schemaByName[opSchema]; s != nil {
			searchNS = []uint32{s.OID}
		}
	} else {
		searchNS = c.searchPathWithCatalog()
	}
	foundAny := false
	var firstUserOID uint32
	for _, nsOID := range searchNS {
		for oid, op := range c.operByOID {
			if op.Name != opName {
				continue
			}
			opNS := uint32(PGCatalogNamespace)
			if oid >= FirstNormalObjectId {
				opNS = c.operatorSchemas[oid]
			}
			if opNS != nsOID {
				continue
			}
			foundAny = true
			if oid >= FirstNormalObjectId {
				if firstUserOID == 0 {
					firstUserOID = oid
				}
				b.addOperDep(oid)
				if op.Code >= FirstNormalObjectId {
					b.addFuncDep(op.Code)
				}
			}
		}
	}
	if firstUserOID != 0 && b != nil && b.operByList != nil {
		b.operByList[opNameList] = firstUserOID
	}
	if requireExists && !foundAny && !c.snapshot {
		shown := opName
		if opSchema != "" {
			shown = opSchema + "." + opName
		}
		return &Error{
			Code:    CodeUndefinedFunction,
			Message: fmt.Sprintf("operator does not exist: %s", shown),
		}
	}
	return nil
}

func (c *Catalog) bindSortAndLimit(b *ruleBindings, s *nodes.SelectStmt, ctx *ruleEvalContext) error {
	if err := c.bindRuleSortClause(b, s.SortClause, ctx); err != nil {
		return err
	}
	if err := c.bindRuleExpr(b, s.LimitCount, ctx); err != nil {
		return err
	}
	if err := c.bindRuleExpr(b, s.LimitOffset, ctx); err != nil {
		return err
	}
	if s.LockingClause != nil {
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

func (c *Catalog) bindWindowDef(b *ruleBindings, wd *nodes.WindowDef, ctx *ruleEvalContext) error {
	if wd == nil {
		return nil
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
	if err := c.bindRuleExpr(b, wd.StartOffset, ctx); err != nil {
		return err
	}
	return c.bindRuleExpr(b, wd.EndOffset, ctx)
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
	if idx, exists := rel.colByName[colName]; exists {
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
							for _, rOID := range se.joinRels {
								if rel := c.relationByOID[rOID]; rel != nil {
									if idx, ok := rel.colByName[colName]; ok {
										b.addColDep(rOID, rel.Columns[idx].AttNum)
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
						b.colByRef[v] = ruleColBinding{relOID: matchedOID, attNum: matchedAtt}
					}
					return nil
				}
				if cur.targetAliases[colName] || hasAnyOpen {
					return nil
				}
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
				b.addColDep(b.ownerRelOID, attNum)
				return nil
			}
			if se, ok := lookupScopeEntry(ctx, qual); ok {
				if se.isVirtual {
					for _, rOID := range se.joinRels {
						if isStar {
							recordAllRelColumnDeps(c, b, rOID)
						} else if rel := c.relationByOID[rOID]; rel != nil {
							if idx, exists := rel.colByName[colName]; exists {
								b.addColDep(rOID, rel.Columns[idx].AttNum)
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
				b.addColDep(se.relOID, attNum)
				return nil
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
		if err := c.bindRuleExpr(b, v.AggFilter, ctx); err != nil {
			return err
		}
		if wd, ok := v.Over.(*nodes.WindowDef); ok {
			if err := c.bindWindowDef(b, wd, ctx); err != nil {
				return err
			}
		}
		if v.Funcname != nil {
			if err := c.bindRuleFuncCall(b, v, argCount); err != nil {
				return err
			}
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
		if v.Name != nil && (v.Kind == nodes.AEXPR_OP || v.Kind == nodes.AEXPR_OP_ANY || v.Kind == nodes.AEXPR_OP_ALL) {
			if err := c.bindRuleOperatorName(b, v.Name, false); err != nil {
				return err
			}
		}
		return nil
	case *nodes.BoolExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.NullTest:
		return c.bindRuleExpr(b, v.Arg, ctx)
	case *nodes.BooleanTest:
		return c.bindRuleExpr(b, v.Arg, ctx)
	case *nodes.CoalesceExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.MinMaxExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
				}
			}
		}
		return nil
	case *nodes.CaseExpr:
		if err := c.bindRuleExpr(b, v.Arg, ctx); err != nil {
			return err
		}
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if cw, ok := item.(*nodes.CaseWhen); ok {
					if err := c.bindRuleExpr(b, cw.Expr, ctx); err != nil {
						return err
					}
					if err := c.bindRuleExpr(b, cw.Result, ctx); err != nil {
						return err
					}
				}
			}
		}
		return c.bindRuleExpr(b, v.Defresult, ctx)
	case *nodes.A_ArrayExpr:
		if v.Elements != nil {
			for _, item := range v.Elements.Items {
				if err := c.bindRuleExpr(b, item, ctx); err != nil {
					return err
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
		if err := c.bindRuleExpr(b, v.Testexpr, ctx); err != nil {
			return err
		}
		if v.OperName != nil {
			if err := c.bindRuleOperatorName(b, v.OperName, false); err != nil {
				return err
			}
		}
		subCtx := &ruleEvalContext{}
		if ctx != nil {
			subCtx.outer = ctx
			subCtx.cteNames = ctx.cteNames
			subCtx.cteCols = ctx.cteCols
			subCtx.cteAny = ctx.cteAny
		}
		return c.bindRuleStmt(b, v.Subselect, subCtx)
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
		return c.bindRuleJsonOutput(b, v.Output)
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
	if err := c.bindRuleExpr(b, ac.Agg_filter, ctx); err != nil {
		return err
	}
	if err := c.bindWindowDef(b, ac.Over, ctx); err != nil {
		return err
	}
	return c.bindRuleJsonOutput(b, ac.Output)
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

func (c *Catalog) bindRuleFuncCall(b *ruleBindings, v *nodes.FuncCall, argCount int) error {
	fSchema, fName := qualifiedName(v.Funcname)
	if fName == "" {
		return nil
	}
	var searchNS []uint32
	if fSchema != "" {
		ns := c.schemaByName[fSchema]
		if ns == nil {
			if c.snapshot {
				return nil
			}
			return errUndefinedSchema(fSchema)
		}
		searchNS = []uint32{ns.OID}
	} else {
		searchNS = c.searchPathWithCatalog()
	}

	procs := c.procByName[fName]
	foundAny := false
	var matchedUserOID uint32
	var matchedArgTypes []uint32
	for _, nsOID := range searchNS {
		var exactUserOID uint32
		var nsFound bool
		var firstArgTypes []uint32
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
			candTypes, ok := procMatchesRuleCallArity(p, v)
			if !ok {
				continue
			}
			nsFound = true
			if firstArgTypes == nil {
				firstArgTypes = candTypes
			}
			if p.OID >= FirstNormalObjectId && exactUserOID == 0 {
				exactUserOID = p.OID
				firstArgTypes = candTypes
			}
		}
		if nsFound {
			foundAny = true
			matchedUserOID = exactUserOID
			matchedArgTypes = firstArgTypes
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
	if matchedUserOID != 0 {
		b.funcByFC[v] = matchedUserOID
		b.addFuncDep(matchedUserOID)
		return nil
	}
	if foundAny {
		return nil
	}

	// PostgreSQL allows single-argument functional cast syntax: typename(expr).
	if argCount == 1 && !v.AggStar && !v.AggDistinct && v.Over == nil {
		for _, nsOID := range searchNS {
			if bt := c.typeByName[typeKey{ns: nsOID, name: fName}]; bt != nil && bt.IsDefined {
				elemOID := bt.OID
				if bt.Elem != 0 {
					elemOID = bt.Elem
				}
				if b.castTypeByFC != nil {
					b.castTypeByFC[v] = elemOID
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
		return &Error{
			Code:    CodeUndefinedFunction,
			Message: fmt.Sprintf("function %s does not exist", shown),
		}
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

func (c *Catalog) inferPreparedDMLParamTypes(b *ruleBindings, n nodes.Node, paramTypes []uint32) {
	if len(paramTypes) == 0 || n == nil {
		return
	}
	setParam := func(num int, typ uint32) {
		if num > 0 && num <= len(paramTypes) && typ != 0 && typ != UNKNOWNOID && (paramTypes[num-1] == 0 || paramTypes[num-1] == UNKNOWNOID) {
			paramTypes[num-1] = typ
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
				switch ev.Val.(type) {
				case *nodes.Integer:
					return INT4OID
				case *nodes.Float:
					return NUMERICOID
				case *nodes.Boolean:
					return BOOLOID
				case *nodes.String:
					return TEXTOID
				}
			}
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
	case *nodes.InsertStmt:
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
				if q, err := c.analyzeSelectStmtWithOptions(sel, AnalyzeOptions{ParamTypes: paramTypes}); err == nil {
					collectQueryParamTypes(q, paramTypes)
				}
			}
		}
		if s.OnConflictClause != nil && rel != nil {
			if s.OnConflictClause.TargetList != nil {
				for _, item := range s.OnConflictClause.TargetList.Items {
					if rt, ok := item.(*nodes.ResTarget); ok {
						walkExpr(rt.Val, targetColType(rt, rel))
					}
				}
			}
			walkExpr(s.OnConflictClause.WhereClause, BOOLOID)
		}
	case *nodes.UpdateStmt:
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
		walkExpr(s.WhereClause, BOOLOID)
	case *nodes.DeleteStmt:
		walkExpr(s.WhereClause, BOOLOID)
	case *nodes.MergeStmt:
		var rel *Relation
		if b != nil && s.Relation != nil {
			rel = c.relationByOID[b.relByRV[s.Relation]]
		}
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
	}
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
	case *nodes.NotifyStmt:
		return deparseRawNotifyStmt(s)
	default:
		return ""
	}
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

func (d *ruleDeparser) currentInferredTargetColName(n nodes.Node, fallback string) string {
	if d.b == nil || n == nil {
		return fallback
	}
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if cb, ok := d.b.colByRef[v]; ok {
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
		b.WriteString(deparseFuncNameList(v.Method))
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
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if v.Fields == nil {
			return ""
		}
		if d.b != nil {
			if cb, ok := d.b.colByRef[v]; ok {
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
