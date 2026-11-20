package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Rule represents a rewrite rule defined via CREATE RULE on a table or view.
type Rule struct {
	OID        uint32
	Name       string
	Event      nodes.CmdType // CMD_SELECT, CMD_INSERT, CMD_UPDATE, CMD_DELETE
	Instead    bool
	Mode       byte     // 'O'=origin (enabled), 'D'=disabled, 'R'=replica, 'A'=always
	WhereSQL   string   // deparsed WHERE qualification (empty if none)
	ActionsSQL []string // deparsed action statements (empty = DO NOTHING)
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

	var whereSQL string
	if s.WhereClause != nil {
		whereSQL = deparseRuleExpr(s.WhereClause)
	}
	var actionsSQL []string
	if s.Actions != nil {
		for _, item := range s.Actions.Items {
			if act := deparseRuleActionStmt(item); act != "" {
				actionsSQL = append(actionsSQL, act)
			}
		}
	}

	rel.Rules[s.Rulename] = &Rule{
		OID:        oid,
		Name:       s.Rulename,
		Event:      s.Event,
		Instead:    s.Instead,
		Mode:       mode,
		WhereSQL:   whereSQL,
		ActionsSQL: actionsSQL,
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

func deparseRuleActionStmt(n nodes.Node) string {
	if n == nil {
		return ""
	}
	switch s := n.(type) {
	case *nodes.SelectStmt:
		return deparseRawSelectStmt(s)
	case *nodes.InsertStmt:
		return deparseRawInsertStmt(s)
	case *nodes.UpdateStmt:
		return deparseRawUpdateStmt(s)
	case *nodes.DeleteStmt:
		return deparseRawDeleteStmt(s)
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

func deparseRawWithClause(w *nodes.WithClause) string {
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
		b.WriteString(deparseRuleActionStmt(cte.Ctequery))
		b.WriteString(")")
	}
	b.WriteString(" ")
	return b.String()
}

func deparseRawSelectStmt(s *nodes.SelectStmt) string {
	if s == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(deparseRawWithClause(s.WithClause))

	if s.Op != nodes.SETOP_NONE {
		if s.Larg != nil {
			b.WriteString("(")
			b.WriteString(deparseRawSelectStmt(s.Larg))
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
			b.WriteString(deparseRawSelectStmt(s.Rarg))
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
					b.WriteString(deparseRuleExpr(v))
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
				b.WriteString(deparseRuleExprList(s.DistinctClause))
				b.WriteString(")")
			}
		}
		if s.TargetList != nil && len(s.TargetList.Items) > 0 {
			b.WriteString(" ")
			b.WriteString(deparseRawTargetList(s.TargetList))
		}
		if s.FromClause != nil && len(s.FromClause.Items) > 0 {
			b.WriteString(" FROM ")
			b.WriteString(deparseRawFromList(s.FromClause))
		}
		if s.WhereClause != nil {
			b.WriteString(" WHERE ")
			b.WriteString(deparseRuleExpr(s.WhereClause))
		}
		if s.GroupClause != nil && len(s.GroupClause.Items) > 0 {
			b.WriteString(" GROUP BY ")
			if s.GroupDistinct {
				b.WriteString("DISTINCT ")
			}
			b.WriteString(deparseRuleExprList(s.GroupClause))
		}
		if s.HavingClause != nil {
			b.WriteString(" HAVING ")
			b.WriteString(deparseRuleExpr(s.HavingClause))
		}
	}

	if s.SortClause != nil && len(s.SortClause.Items) > 0 {
		b.WriteString(" ORDER BY ")
		b.WriteString(deparseRawSortClause(s.SortClause))
	}
	if s.LimitCount != nil {
		b.WriteString(" LIMIT ")
		b.WriteString(deparseRuleExpr(s.LimitCount))
	}
	if s.LimitOffset != nil {
		b.WriteString(" OFFSET ")
		b.WriteString(deparseRuleExpr(s.LimitOffset))
	}
	return b.String()
}

func deparseRawInsertStmt(s *nodes.InsertStmt) string {
	var b strings.Builder
	b.WriteString(deparseRawWithClause(s.WithClause))
	b.WriteString("INSERT INTO ")
	b.WriteString(deparseRawRangeVar(s.Relation, false))
	if s.Cols != nil && len(s.Cols.Items) > 0 {
		b.WriteString(" (")
		for i, item := range s.Cols.Items {
			if i > 0 {
				b.WriteString(", ")
			}
			if rt, ok := item.(*nodes.ResTarget); ok {
				b.WriteString(quoteIdentifier(rt.Name))
				b.WriteString(deparseRawIndirection(rt.Indirection))
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
		b.WriteString(deparseRawSelectStmt(sel))
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
							b.WriteString(deparseRuleExpr(elem.Expr))
							b.WriteString(")")
						}
					}
				}
				b.WriteString(")")
				if oc.Infer.WhereClause != nil {
					b.WriteString(" WHERE ")
					b.WriteString(deparseRuleExpr(oc.Infer.WhereClause))
				}
			}
		}
		if oc.TargetList != nil && len(oc.TargetList.Items) > 0 {
			b.WriteString(" DO UPDATE SET ")
			b.WriteString(deparseRawUpdateTargetList(oc.TargetList))
			if oc.WhereClause != nil {
				b.WriteString(" WHERE ")
				b.WriteString(deparseRuleExpr(oc.WhereClause))
			}
		} else {
			b.WriteString(" DO NOTHING")
		}
	}
	if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
		b.WriteString(" RETURNING ")
		b.WriteString(deparseRawTargetList(s.ReturningList))
	}
	return b.String()
}

func deparseRawUpdateStmt(s *nodes.UpdateStmt) string {
	var b strings.Builder
	b.WriteString(deparseRawWithClause(s.WithClause))
	b.WriteString("UPDATE ")
	b.WriteString(deparseRawRangeVar(s.Relation, true))
	b.WriteString(" SET ")
	b.WriteString(deparseRawUpdateTargetList(s.TargetList))
	if s.FromClause != nil && len(s.FromClause.Items) > 0 {
		b.WriteString(" FROM ")
		b.WriteString(deparseRawFromList(s.FromClause))
	}
	if s.WhereClause != nil {
		b.WriteString(" WHERE ")
		b.WriteString(deparseRuleExpr(s.WhereClause))
	}
	if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
		b.WriteString(" RETURNING ")
		b.WriteString(deparseRawTargetList(s.ReturningList))
	}
	return b.String()
}

func deparseRawDeleteStmt(s *nodes.DeleteStmt) string {
	var b strings.Builder
	b.WriteString(deparseRawWithClause(s.WithClause))
	b.WriteString("DELETE FROM ")
	b.WriteString(deparseRawRangeVar(s.Relation, true))
	if s.UsingClause != nil && len(s.UsingClause.Items) > 0 {
		b.WriteString(" USING ")
		b.WriteString(deparseRawFromList(s.UsingClause))
	}
	if s.WhereClause != nil {
		b.WriteString(" WHERE ")
		b.WriteString(deparseRuleExpr(s.WhereClause))
	}
	if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
		b.WriteString(" RETURNING ")
		b.WriteString(deparseRawTargetList(s.ReturningList))
	}
	return b.String()
}

func deparseRawTargetList(list *nodes.List) string {
	if list == nil {
		return ""
	}
	var parts []string
	for _, item := range list.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok {
			parts = append(parts, deparseRuleExpr(item))
			continue
		}
		val := deparseRuleExpr(rt.Val)
		if rt.Name != "" {
			parts = append(parts, val+" AS "+quoteIdentifier(rt.Name))
		} else {
			parts = append(parts, val)
		}
	}
	return strings.Join(parts, ", ")
}

func deparseRawUpdateTargetList(list *nodes.List) string {
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
			cols = append(cols, quoteIdentifier(rt.Name)+deparseRawIndirection(rt.Indirection))
			for k := 1; k < mar.Ncolumns && i+1 < len(list.Items); k++ {
				i++
				if nextRT, ok := list.Items[i].(*nodes.ResTarget); ok {
					cols = append(cols, quoteIdentifier(nextRT.Name)+deparseRawIndirection(nextRT.Indirection))
				}
			}
			parts = append(parts, "("+strings.Join(cols, ", ")+") = "+deparseRuleExpr(mar.Source))
			continue
		}
		parts = append(parts, quoteIdentifier(rt.Name)+deparseRawIndirection(rt.Indirection)+" = "+deparseRuleExpr(rt.Val))
	}
	return strings.Join(parts, ", ")
}

func deparseRawFromList(list *nodes.List) string {
	if list == nil {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		parts = append(parts, deparseRawTableRef(item))
	}
	return strings.Join(parts, ", ")
}

func deparseRawTableRef(n nodes.Node) string {
	switch v := n.(type) {
	case *nodes.RangeVar:
		return deparseRawRangeVar(v, true)
	case *nodes.RangeSubselect:
		var b strings.Builder
		if v.Lateral {
			b.WriteString("LATERAL ")
		}
		b.WriteString("(")
		b.WriteString(deparseRuleActionStmt(v.Subquery))
		b.WriteString(")")
		b.WriteString(deparseRawAlias(v.Alias))
		return b.String()
	case *nodes.JoinExpr:
		var b strings.Builder
		b.WriteString(deparseRawTableRef(v.Larg))
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
		b.WriteString(deparseRawTableRef(v.Rarg))
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
			b.WriteString(deparseRuleExpr(v.Quals))
		}
		return b.String()
	default:
		return deparseRuleExpr(n)
	}
}

func deparseRawRangeVar(rv *nodes.RangeVar, includeOnly bool) string {
	if rv == nil {
		return ""
	}
	var b strings.Builder
	if includeOnly && !rv.Inh && (rv.Loc.Start != 0 || rv.Loc.End != 0) {
		b.WriteString("ONLY ")
	}
	if rv.Schemaname != "" {
		b.WriteString(quoteIdentifier(rv.Schemaname))
		b.WriteString(".")
	}
	b.WriteString(quoteIdentifier(rv.Relname))
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

func deparseRawSortClause(list *nodes.List) string {
	if list == nil {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		sb, ok := item.(*nodes.SortBy)
		if !ok {
			continue
		}
		s := deparseRuleExpr(sb.Node)
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

func deparseRuleExprList(list *nodes.List) string {
	if list == nil {
		return ""
	}
	parts := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		parts = append(parts, deparseRuleExpr(item))
	}
	return strings.Join(parts, ", ")
}

func deparseRawIndirection(ind *nodes.List) string {
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
					b.WriteString(deparseRuleExpr(v.Lidx))
				}
				b.WriteString(":")
				if v.Uidx != nil {
					b.WriteString(deparseRuleExpr(v.Uidx))
				}
			} else if v.Uidx != nil {
				b.WriteString(deparseRuleExpr(v.Uidx))
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

func deparseRuleExpr(n nodes.Node) string {
	if n == nil {
		return ""
	}
	switch v := n.(type) {
	case *nodes.ColumnRef:
		if v.Fields == nil {
			return ""
		}
		parts := make([]string, 0, len(v.Fields.Items))
		for i, f := range v.Fields.Items {
			switch fv := f.(type) {
			case *nodes.String:
				// Preserve OLD and NEW correlation names in canonical lowercase unquoted form.
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
		return fmt.Sprintf("(%s)::%s", deparseRuleExpr(v.Arg), deparseRawTypeName(v.TypeName))

	case *nodes.CollateClause:
		return fmt.Sprintf("(%s COLLATE %s)", deparseRuleExpr(v.Arg), deparseQualifiedIdentList(v.Collname))

	case *nodes.A_Indirection:
		return fmt.Sprintf("(%s)%s", deparseRuleExpr(v.Arg), deparseRawIndirection(v.Indirection))

	case *nodes.A_Expr:
		op := deparseOperatorName(v.Name)
		switch v.Kind {
		case nodes.AEXPR_OP:
			if v.Lexpr == nil {
				return fmt.Sprintf("(%s %s)", op, deparseRuleExpr(v.Rexpr))
			}
			return fmt.Sprintf("(%s %s %s)", deparseRuleExpr(v.Lexpr), op, deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_OP_ANY:
			return fmt.Sprintf("(%s %s ANY (%s))", deparseRuleExpr(v.Lexpr), op, deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_OP_ALL:
			return fmt.Sprintf("(%s %s ALL (%s))", deparseRuleExpr(v.Lexpr), op, deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_DISTINCT:
			return fmt.Sprintf("(%s IS DISTINCT FROM %s)", deparseRuleExpr(v.Lexpr), deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_NOT_DISTINCT:
			return fmt.Sprintf("(%s IS NOT DISTINCT FROM %s)", deparseRuleExpr(v.Lexpr), deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_NULLIF:
			return fmt.Sprintf("NULLIF(%s, %s)", deparseRuleExpr(v.Lexpr), deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_IN:
			kw := "IN"
			if op == "<>" {
				kw = "NOT IN"
			}
			rhs, _ := v.Rexpr.(*nodes.List)
			return fmt.Sprintf("(%s %s (%s))", deparseRuleExpr(v.Lexpr), kw, deparseRuleExprList(rhs))
		case nodes.AEXPR_LIKE:
			kw := "LIKE"
			if op == "!~~" {
				kw = "NOT LIKE"
			}
			return fmt.Sprintf("(%s %s %s)", deparseRuleExpr(v.Lexpr), kw, deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_ILIKE:
			kw := "ILIKE"
			if op == "!~~*" {
				kw = "NOT ILIKE"
			}
			return fmt.Sprintf("(%s %s %s)", deparseRuleExpr(v.Lexpr), kw, deparseRuleExpr(v.Rexpr))
		case nodes.AEXPR_SIMILAR:
			kw := "SIMILAR TO"
			if op == "!~" {
				kw = "NOT SIMILAR TO"
			}
			return fmt.Sprintf("(%s %s %s)", deparseRuleExpr(v.Lexpr), kw, deparseRuleExpr(v.Rexpr))
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
				return fmt.Sprintf("(%s %s %s AND %s)", deparseRuleExpr(v.Lexpr), kw, deparseRuleExpr(pair.Items[0]), deparseRuleExpr(pair.Items[1]))
			}
		}
		return fmt.Sprintf("(%s %s %s)", deparseRuleExpr(v.Lexpr), op, deparseRuleExpr(v.Rexpr))

	case *nodes.BoolExpr:
		if v.Args == nil || len(v.Args.Items) == 0 {
			return ""
		}
		if v.Boolop == nodes.NOT_EXPR {
			return fmt.Sprintf("(NOT %s)", deparseRuleExpr(v.Args.Items[0]))
		}
		sep := " AND "
		if v.Boolop == nodes.OR_EXPR {
			sep = " OR "
		}
		parts := make([]string, 0, len(v.Args.Items))
		for _, a := range v.Args.Items {
			parts = append(parts, deparseRuleExpr(a))
		}
		return "(" + strings.Join(parts, sep) + ")"

	case *nodes.NullTest:
		if v.Nulltesttype == nodes.IS_NOT_NULL {
			return fmt.Sprintf("(%s IS NOT NULL)", deparseRuleExpr(v.Arg))
		}
		return fmt.Sprintf("(%s IS NULL)", deparseRuleExpr(v.Arg))

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
		return fmt.Sprintf("(%s %s)", deparseRuleExpr(v.Arg), suffix)

	case *nodes.FuncCall:
		var b strings.Builder
		b.WriteString(deparseFuncNameList(v.Funcname))
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
					b.WriteString(deparseRuleExpr(arg))
				}
			}
			if !v.AggWithinGroup && v.AggOrder != nil && len(v.AggOrder.Items) > 0 {
				b.WriteString(" ORDER BY ")
				b.WriteString(deparseRawSortClause(v.AggOrder))
			}
		}
		b.WriteString(")")
		if v.AggWithinGroup && v.AggOrder != nil && len(v.AggOrder.Items) > 0 {
			b.WriteString(" WITHIN GROUP (ORDER BY ")
			b.WriteString(deparseRawSortClause(v.AggOrder))
			b.WriteString(")")
		}
		if v.AggFilter != nil {
			b.WriteString(" FILTER (WHERE ")
			b.WriteString(deparseRuleExpr(v.AggFilter))
			b.WriteString(")")
		}
		return b.String()

	case *nodes.NamedArgExpr:
		return fmt.Sprintf("%s => %s", quoteIdentifier(v.Name), deparseRuleExpr(v.Arg))

	case *nodes.CoalesceExpr:
		return fmt.Sprintf("COALESCE(%s)", deparseRuleExprList(v.Args))

	case *nodes.MinMaxExpr:
		kw := "GREATEST"
		if v.Op == nodes.IS_LEAST {
			kw = "LEAST"
		}
		return fmt.Sprintf("%s(%s)", kw, deparseRuleExprList(v.Args))

	case *nodes.CaseExpr:
		var b strings.Builder
		b.WriteString("CASE")
		if v.Arg != nil {
			b.WriteString(" ")
			b.WriteString(deparseRuleExpr(v.Arg))
		}
		if v.Args != nil {
			for _, w := range v.Args.Items {
				if cw, ok := w.(*nodes.CaseWhen); ok {
					b.WriteString(" WHEN ")
					b.WriteString(deparseRuleExpr(cw.Expr))
					b.WriteString(" THEN ")
					b.WriteString(deparseRuleExpr(cw.Result))
				}
			}
		}
		if v.Defresult != nil {
			b.WriteString(" ELSE ")
			b.WriteString(deparseRuleExpr(v.Defresult))
		}
		b.WriteString(" END")
		return b.String()

	case *nodes.A_ArrayExpr:
		return fmt.Sprintf("ARRAY[%s]", deparseRuleExprList(v.Elements))

	case *nodes.RowExpr:
		if v.Args == nil || len(v.Args.Items) == 0 {
			return "ROW()"
		}
		return fmt.Sprintf("ROW(%s)", deparseRuleExprList(v.Args))

	case *nodes.SubLink:
		sub := deparseRuleActionStmt(v.Subselect)
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
				return fmt.Sprintf("(%s IN (%s))", deparseRuleExpr(v.Testexpr), sub)
			}
			return fmt.Sprintf("(%s %s ANY (%s))", deparseRuleExpr(v.Testexpr), op, sub)
		case nodes.ALL_SUBLINK:
			op := "="
			if v.OperName != nil {
				op = deparseOperatorName(v.OperName)
			}
			return fmt.Sprintf("(%s %s ALL (%s))", deparseRuleExpr(v.Testexpr), op, sub)
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

func deparseRawTypeName(tn *nodes.TypeName) string {
	if tn == nil {
		return ""
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
		base += "(" + deparseRuleExprList(tn.Typmods) + ")"
	}
	if tn.ArrayBounds != nil {
		for range tn.ArrayBounds.Items {
			base += "[]"
		}
	}
	return base
}
