package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Rule represents a rewrite rule defined via CREATE RULE on a table or view.
type Rule struct {
	OID         uint32
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

type ruleBindings struct {
	ownerRelOID uint32
	relByRV     map[*nodes.RangeVar]uint32
	colByRef    map[*nodes.ColumnRef]ruleColBinding
	colByRT     map[*nodes.ResTarget]ruleColBinding
	typeByTN    map[*nodes.TypeName]ruleTypeBinding
	collByCC    map[*nodes.CollateClause]uint32
	funcByFC    map[*nodes.FuncCall]uint32
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
	if !s.Replace && rel.Rules != nil && rel.Rules[s.Rulename] != nil {
		return errDuplicateObject("rule", s.Rulename)
	}
	if rel.Rules == nil {
		rel.Rules = make(map[string]*Rule)
	}

	mode := byte('O')
	oid := c.oidGen.Next()
	if existing := rel.Rules[s.Rulename]; existing != nil {
		mode = existing.Mode
		oid = existing.OID
	}

	bindings := c.buildRuleBindings(rel, s.WhereClause, s.Actions)
	rel.Rules[s.Rulename] = &Rule{
		OID:         oid,
		Name:        s.Rulename,
		Event:       s.Event,
		Instead:     s.Instead,
		Mode:        mode,
		WhereClause: s.WhereClause,
		Actions:     s.Actions,
		bindings:    bindings,
	}
	rel.HasRules = true
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
	relOID   uint32
	hasAlias bool
	alias    string
}

func (c *Catalog) buildRuleBindings(ownerRel *Relation, whereClause nodes.Node, actions *nodes.List) *ruleBindings {
	b := &ruleBindings{
		ownerRelOID: ownerRel.OID,
		relByRV:     make(map[*nodes.RangeVar]uint32),
		colByRef:    make(map[*nodes.ColumnRef]ruleColBinding),
		colByRT:     make(map[*nodes.ResTarget]ruleColBinding),
		typeByTN:    make(map[*nodes.TypeName]ruleTypeBinding),
		collByCC:    make(map[*nodes.CollateClause]uint32),
		funcByFC:    make(map[*nodes.FuncCall]uint32),
	}
	if whereClause != nil {
		c.bindRuleExpr(b, whereClause, nil)
	}
	if actions != nil {
		for _, item := range actions.Items {
			c.bindRuleStmt(b, item, nil)
		}
	}
	return b
}

func (c *Catalog) bindRangeVar(b *ruleBindings, rv *nodes.RangeVar, cteNames map[string]bool) uint32 {
	if rv == nil {
		return 0
	}
	if rv.Schemaname == "" && cteNames[rv.Relname] {
		return 0
	}
	if _, rel, err := c.findRelation(rv.Schemaname, rv.Relname); err == nil && rel != nil {
		b.relByRV[rv] = rel.OID
		return rel.OID
	}
	return 0
}

func addScopeRangeVar(c *Catalog, b *ruleBindings, scope map[string]ruleScopeEntry, rv *nodes.RangeVar, cteNames map[string]bool) uint32 {
	relOID := c.bindRangeVar(b, rv, cteNames)
	if relOID == 0 {
		return 0
	}
	if rv.Alias != nil && rv.Alias.Aliasname != "" {
		scope[rv.Alias.Aliasname] = ruleScopeEntry{relOID: relOID, hasAlias: true, alias: rv.Alias.Aliasname}
	} else {
		scope[rv.Relname] = ruleScopeEntry{relOID: relOID}
	}
	return relOID
}

func (c *Catalog) collectFromScope(b *ruleBindings, list *nodes.List, scope map[string]ruleScopeEntry, cteNames map[string]bool) {
	if list == nil {
		return
	}
	for _, item := range list.Items {
		c.collectTableRefScope(b, item, scope, cteNames)
	}
}

func (c *Catalog) collectTableRefScope(b *ruleBindings, n nodes.Node, scope map[string]ruleScopeEntry, cteNames map[string]bool) {
	switch v := n.(type) {
	case *nodes.RangeVar:
		addScopeRangeVar(c, b, scope, v, cteNames)
	case *nodes.RangeSubselect:
		c.bindRuleStmt(b, v.Subquery, cteNames)
	case *nodes.JoinExpr:
		c.collectTableRefScope(b, v.Larg, scope, cteNames)
		c.collectTableRefScope(b, v.Rarg, scope, cteNames)
		if v.Quals != nil {
			c.bindRuleExpr(b, v.Quals, scope)
		}
	}
}

func (c *Catalog) bindWithClause(b *ruleBindings, w *nodes.WithClause, parentCTEs map[string]bool) map[string]bool {
	cteNames := make(map[string]bool, len(parentCTEs))
	for k, v := range parentCTEs {
		cteNames[k] = v
	}
	if w == nil || w.Ctes == nil {
		return cteNames
	}
	if w.Recursive {
		for _, item := range w.Ctes.Items {
			if cte, ok := item.(*nodes.CommonTableExpr); ok && cte.Ctename != "" {
				cteNames[cte.Ctename] = true
			}
		}
	}
	for _, item := range w.Ctes.Items {
		if cte, ok := item.(*nodes.CommonTableExpr); ok {
			c.bindRuleStmt(b, cte.Ctequery, cteNames)
			if cte.Ctename != "" {
				cteNames[cte.Ctename] = true
			}
		}
	}
	return cteNames
}

func (c *Catalog) bindTargetColumns(b *ruleBindings, relOID uint32, list *nodes.List) {
	if relOID == 0 || list == nil {
		return
	}
	rel := c.relationByOID[relOID]
	if rel == nil {
		return
	}
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok || rt.Name == "" {
			continue
		}
		if idx, exists := rel.colByName[rt.Name]; exists {
			b.colByRT[rt] = ruleColBinding{relOID: relOID, attNum: rel.Columns[idx].AttNum}
		}
	}
}

func (c *Catalog) bindRuleStmt(b *ruleBindings, n nodes.Node, parentCTEs map[string]bool) {
	if n == nil {
		return
	}
	switch s := n.(type) {
	case *nodes.SelectStmt:
		cteNames := c.bindWithClause(b, s.WithClause, parentCTEs)
		if s.Op != nodes.SETOP_NONE {
			if s.Larg != nil {
				c.bindRuleStmt(b, s.Larg, cteNames)
			}
			if s.Rarg != nil {
				c.bindRuleStmt(b, s.Rarg, cteNames)
			}
			c.bindSortAndLimit(b, s, nil)
			return
		}
		if s.ValuesLists != nil {
			for _, rowNode := range s.ValuesLists.Items {
				if row, ok := rowNode.(*nodes.List); ok {
					for _, expr := range row.Items {
						c.bindRuleExpr(b, expr, nil)
					}
				}
			}
			c.bindSortAndLimit(b, s, nil)
			return
		}
		scope := make(map[string]ruleScopeEntry)
		c.collectFromScope(b, s.FromClause, scope, cteNames)
		if s.DistinctClause != nil {
			for _, item := range s.DistinctClause.Items {
				c.bindRuleExpr(b, item, scope)
			}
		}
		if s.TargetList != nil {
			for _, item := range s.TargetList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					c.bindRuleExpr(b, rt.Val, scope)
				} else {
					c.bindRuleExpr(b, item, scope)
				}
			}
		}
		c.bindRuleExpr(b, s.WhereClause, scope)
		if s.GroupClause != nil {
			for _, item := range s.GroupClause.Items {
				c.bindRuleExpr(b, item, scope)
			}
		}
		c.bindRuleExpr(b, s.HavingClause, scope)
		if s.WindowClause != nil {
			for _, item := range s.WindowClause.Items {
				if wd, ok := item.(*nodes.WindowDef); ok {
					c.bindWindowDef(b, wd, scope)
				}
			}
		}
		c.bindSortAndLimit(b, s, scope)

	case *nodes.InsertStmt:
		cteNames := c.bindWithClause(b, s.WithClause, parentCTEs)
		targetOID := c.bindRangeVar(b, s.Relation, cteNames)
		c.bindTargetColumns(b, targetOID, s.Cols)
		if s.SelectStmt != nil {
			c.bindRuleStmt(b, s.SelectStmt, cteNames)
		}
		scope := make(map[string]ruleScopeEntry)
		if targetOID != 0 && s.Relation != nil {
			if s.Relation.Alias != nil && s.Relation.Alias.Aliasname != "" {
				scope[s.Relation.Alias.Aliasname] = ruleScopeEntry{relOID: targetOID, hasAlias: true, alias: s.Relation.Alias.Aliasname}
			} else {
				scope[s.Relation.Relname] = ruleScopeEntry{relOID: targetOID}
			}
			scope["excluded"] = ruleScopeEntry{relOID: targetOID, hasAlias: true, alias: "excluded"}
		}
		if s.OnConflictClause != nil {
			oc := s.OnConflictClause
			if oc.Infer != nil {
				if oc.Infer.IndexElems != nil && targetOID != 0 {
					if rel := c.relationByOID[targetOID]; rel != nil {
						for _, ie := range oc.Infer.IndexElems.Items {
							if elem, ok := ie.(*nodes.IndexElem); ok {
								if elem.Expr != nil {
									c.bindRuleExpr(b, elem.Expr, scope)
								}
							}
						}
					}
				}
				c.bindRuleExpr(b, oc.Infer.WhereClause, scope)
			}
			c.bindTargetColumns(b, targetOID, oc.TargetList)
			if oc.TargetList != nil {
				for _, item := range oc.TargetList.Items {
					if rt, ok := item.(*nodes.ResTarget); ok {
						c.bindRuleExpr(b, rt.Val, scope)
					}
				}
			}
			c.bindRuleExpr(b, oc.WhereClause, scope)
		}
		if s.ReturningList != nil {
			for _, item := range s.ReturningList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					c.bindRuleExpr(b, rt.Val, scope)
				}
			}
		}

	case *nodes.UpdateStmt:
		cteNames := c.bindWithClause(b, s.WithClause, parentCTEs)
		scope := make(map[string]ruleScopeEntry)
		targetOID := addScopeRangeVar(c, b, scope, s.Relation, cteNames)
		c.collectFromScope(b, s.FromClause, scope, cteNames)
		c.bindTargetColumns(b, targetOID, s.TargetList)
		if s.TargetList != nil {
			for _, item := range s.TargetList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					c.bindRuleExpr(b, rt.Val, scope)
				}
			}
		}
		c.bindRuleExpr(b, s.WhereClause, scope)
		if s.ReturningList != nil {
			for _, item := range s.ReturningList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					c.bindRuleExpr(b, rt.Val, scope)
				}
			}
		}

	case *nodes.DeleteStmt:
		cteNames := c.bindWithClause(b, s.WithClause, parentCTEs)
		scope := make(map[string]ruleScopeEntry)
		addScopeRangeVar(c, b, scope, s.Relation, cteNames)
		c.collectFromScope(b, s.UsingClause, scope, cteNames)
		c.bindRuleExpr(b, s.WhereClause, scope)
		if s.ReturningList != nil {
			for _, item := range s.ReturningList.Items {
				if rt, ok := item.(*nodes.ResTarget); ok {
					c.bindRuleExpr(b, rt.Val, scope)
				}
			}
		}
	}
}

func (c *Catalog) bindSortAndLimit(b *ruleBindings, s *nodes.SelectStmt, scope map[string]ruleScopeEntry) {
	if s.SortClause != nil {
		for _, item := range s.SortClause.Items {
			if sb, ok := item.(*nodes.SortBy); ok {
				c.bindRuleExpr(b, sb.Node, scope)
			}
		}
	}
	c.bindRuleExpr(b, s.LimitCount, scope)
	c.bindRuleExpr(b, s.LimitOffset, scope)
}

func (c *Catalog) bindWindowDef(b *ruleBindings, wd *nodes.WindowDef, scope map[string]ruleScopeEntry) {
	if wd == nil {
		return
	}
	if wd.PartitionClause != nil {
		for _, item := range wd.PartitionClause.Items {
			c.bindRuleExpr(b, item, scope)
		}
	}
	if wd.OrderClause != nil {
		for _, item := range wd.OrderClause.Items {
			if sb, ok := item.(*nodes.SortBy); ok {
				c.bindRuleExpr(b, sb.Node, scope)
			}
		}
	}
	c.bindRuleExpr(b, wd.StartOffset, scope)
	c.bindRuleExpr(b, wd.EndOffset, scope)
}

func (c *Catalog) bindRuleExpr(b *ruleBindings, n nodes.Node, scope map[string]ruleScopeEntry) {
	if n == nil {
		return
	}
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if v.Fields == nil {
			return
		}
		items := v.Fields.Items
		if len(items) == 2 {
			qual := stringVal(items[0])
			colName := stringVal(items[1])
			if colName == "" {
				return
			}
			if strings.EqualFold(qual, "old") || strings.EqualFold(qual, "new") {
				if ownerRel := c.relationByOID[b.ownerRelOID]; ownerRel != nil {
					if idx, exists := ownerRel.colByName[colName]; exists {
						b.colByRef[v] = ruleColBinding{
							relOID: b.ownerRelOID,
							attNum: ownerRel.Columns[idx].AttNum,
							prefix: strings.ToLower(qual),
						}
					}
				}
				return
			}
			if se, ok := scope[qual]; ok {
				if rel := c.relationByOID[se.relOID]; rel != nil {
					if idx, exists := rel.colByName[colName]; exists {
						b.colByRef[v] = ruleColBinding{
							relOID:     se.relOID,
							attNum:     rel.Columns[idx].AttNum,
							prefix:     se.alias,
							useRelName: !se.hasAlias,
						}
					}
				}
			}
		} else if len(items) == 1 {
			colName := stringVal(items[0])
			if colName == "" || len(scope) == 0 {
				return
			}
			var matchedOID uint32
			var matchedAtt int16
			count := 0
			seenRel := map[uint32]bool{}
			for k, se := range scope {
				if k == "excluded" || seenRel[se.relOID] {
					continue
				}
				seenRel[se.relOID] = true
				if rel := c.relationByOID[se.relOID]; rel != nil {
					if idx, exists := rel.colByName[colName]; exists {
						matchedOID = se.relOID
						matchedAtt = rel.Columns[idx].AttNum
						count++
					}
				}
			}
			if count == 1 {
				b.colByRef[v] = ruleColBinding{relOID: matchedOID, attNum: matchedAtt}
			}
		}

	case *nodes.TypeCast:
		c.bindRuleExpr(b, v.Arg, scope)
		if v.TypeName != nil {
			if oid, _, err := c.resolveTypeName(v.TypeName); err == nil && oid != 0 {
				isArr := v.TypeName.ArrayBounds != nil && len(v.TypeName.ArrayBounds.Items) > 0
				if isArr {
					if bt := c.typeByOID[oid]; bt != nil && bt.Elem != 0 {
						oid = bt.Elem
					}
				}
				b.typeByTN[v.TypeName] = ruleTypeBinding{typeOID: oid, isArray: isArr}
			}
		}

	case *nodes.CollateClause:
		c.bindRuleExpr(b, v.Arg, scope)
		if _, uid, err := c.resolveCollationList(v.Collname); err == nil && uid != 0 {
			b.collByCC[v] = uid
		}

	case *nodes.FuncCall:
		if v.Funcname != nil {
			fSchema, fName := qualifiedName(v.Funcname)
			if fSchema != "pg_catalog" {
				var matched uint32
				count := 0
				for oid, up := range c.userProcs {
					if up.Name != fName || up.Schema == nil {
						continue
					}
					if fSchema != "" && up.Schema.Name != fSchema {
						continue
					}
					matched = oid
					count++
				}
				if count == 1 {
					b.funcByFC[v] = matched
				}
			}
		}
		if v.Args != nil {
			for _, arg := range v.Args.Items {
				c.bindRuleExpr(b, arg, scope)
			}
		}
		if v.AggOrder != nil {
			for _, item := range v.AggOrder.Items {
				if sb, ok := item.(*nodes.SortBy); ok {
					c.bindRuleExpr(b, sb.Node, scope)
				}
			}
		}
		c.bindRuleExpr(b, v.AggFilter, scope)
		if wd, ok := v.Over.(*nodes.WindowDef); ok {
			c.bindWindowDef(b, wd, scope)
		}

	case *nodes.NamedArgExpr:
		c.bindRuleExpr(b, v.Arg, scope)
	case *nodes.A_Indirection:
		c.bindRuleExpr(b, v.Arg, scope)
		if v.Indirection != nil {
			for _, item := range v.Indirection.Items {
				if idx, ok := item.(*nodes.A_Indices); ok {
					c.bindRuleExpr(b, idx.Lidx, scope)
					c.bindRuleExpr(b, idx.Uidx, scope)
				}
			}
		}
	case *nodes.A_Expr:
		c.bindRuleExpr(b, v.Lexpr, scope)
		if list, ok := v.Rexpr.(*nodes.List); ok {
			for _, item := range list.Items {
				c.bindRuleExpr(b, item, scope)
			}
		} else {
			c.bindRuleExpr(b, v.Rexpr, scope)
		}
	case *nodes.BoolExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				c.bindRuleExpr(b, item, scope)
			}
		}
	case *nodes.NullTest:
		c.bindRuleExpr(b, v.Arg, scope)
	case *nodes.BooleanTest:
		c.bindRuleExpr(b, v.Arg, scope)
	case *nodes.CoalesceExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				c.bindRuleExpr(b, item, scope)
			}
		}
	case *nodes.MinMaxExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				c.bindRuleExpr(b, item, scope)
			}
		}
	case *nodes.CaseExpr:
		c.bindRuleExpr(b, v.Arg, scope)
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if cw, ok := item.(*nodes.CaseWhen); ok {
					c.bindRuleExpr(b, cw.Expr, scope)
					c.bindRuleExpr(b, cw.Result, scope)
				}
			}
		}
		c.bindRuleExpr(b, v.Defresult, scope)
	case *nodes.A_ArrayExpr:
		if v.Elements != nil {
			for _, item := range v.Elements.Items {
				c.bindRuleExpr(b, item, scope)
			}
		}
	case *nodes.RowExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				c.bindRuleExpr(b, item, scope)
			}
		}
	case *nodes.SubLink:
		c.bindRuleExpr(b, v.Testexpr, scope)
		c.bindRuleStmt(b, v.Subselect, nil)
	case *nodes.MultiAssignRef:
		c.bindRuleExpr(b, v.Source, scope)
	}
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
	if s.LimitCount != nil {
		b.WriteString(" LIMIT ")
		b.WriteString(d.deparseExpr(s.LimitCount))
	}
	if s.LimitOffset != nil {
		b.WriteString(" OFFSET ")
		b.WriteString(d.deparseExpr(s.LimitOffset))
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
				b.WriteString(" ON CONSTRAINT ")
				b.WriteString(quoteIdentifier(oc.Infer.Conname))
			} else if oc.Infer.IndexElems != nil && len(oc.Infer.IndexElems.Items) > 0 {
				b.WriteString(" (")
				for i, ie := range oc.Infer.IndexElems.Items {
					if i > 0 {
						b.WriteString(", ")
					}
					if elem, ok := ie.(*nodes.IndexElem); ok {
						if elem.Name != "" {
							b.WriteString(quoteIdentifier(elem.Name))
						} else if elem.Expr != nil {
							b.WriteString("(")
							b.WriteString(d.deparseExpr(elem.Expr))
							b.WriteString(")")
						}
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
		} else {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, ", ")
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
	case *nodes.JoinExpr:
		var b strings.Builder
		b.WriteString(d.deparseTableRef(v.Larg))
		if v.IsNatural {
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
			if v.Quals == nil && v.UsingClause == nil && !v.IsNatural {
				b.WriteString(" CROSS JOIN ")
			} else {
				b.WriteString(" JOIN ")
			}
		}
		b.WriteString(d.deparseTableRef(v.Rarg))
		if v.UsingClause != nil && len(v.UsingClause.Items) > 0 {
			b.WriteString(" USING (")
			for i, u := range v.UsingClause.Items {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(quoteIdentifier(stringVal(u)))
			}
			b.WriteString(")")
		} else if v.Quals != nil {
			b.WriteString(" ON ")
			b.WriteString(d.deparseExpr(v.Quals))
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
				s += " USING " + deparseOperatorName(sb.UseOp)
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
		return deparseRawConstVal(v.Val)

	case *nodes.Integer, *nodes.Float, *nodes.String, *nodes.Boolean, *nodes.BitString:
		return deparseRawConstVal(v)

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
		op := deparseOperatorName(v.Name)
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
				op = deparseOperatorName(v.OperName)
			}
			if op == "=" {
				return fmt.Sprintf("(%s IN (%s))", d.deparseExpr(v.Testexpr), sub)
			}
			return fmt.Sprintf("(%s %s ANY (%s))", d.deparseExpr(v.Testexpr), op, sub)
		case nodes.ALL_SUBLINK:
			op := "="
			if v.OperName != nil {
				op = deparseOperatorName(v.OperName)
			}
			return fmt.Sprintf("(%s %s ALL (%s))", d.deparseExpr(v.Testexpr), op, sub)
		default:
			return fmt.Sprintf("(%s)", sub)
		}

	case *nodes.SQLValueFunction:
		switch v.Op {
		case nodes.SVFOP_CURRENT_DATE:
			return "CURRENT_DATE"
		case nodes.SVFOP_CURRENT_TIME, nodes.SVFOP_CURRENT_TIME_N:
			return "CURRENT_TIME"
		case nodes.SVFOP_CURRENT_TIMESTAMP, nodes.SVFOP_CURRENT_TIMESTAMP_N:
			return "CURRENT_TIMESTAMP"
		case nodes.SVFOP_LOCALTIME, nodes.SVFOP_LOCALTIME_N:
			return "LOCALTIME"
		case nodes.SVFOP_LOCALTIMESTAMP, nodes.SVFOP_LOCALTIMESTAMP_N:
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
	}
	return ""
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
