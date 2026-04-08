package ast

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// NodeToString returns a deterministic textual dump of any AST node.
// The format is Go-struct-like for readability:
//
//	BinaryExpr{Op:+ Left:VarRef{Name:a} Right:NumberLit{Val:1}}
//
// Loc fields are not dumped (positions are tested separately by
// TestGetLoc). For nil input, returns "<nil>".
//
// Used by ast_test.go for snapshot-style assertions and by future
// parser golden tests.
func NodeToString(n Node) string {
	if n == nil {
		return "<nil>"
	}
	var sb strings.Builder
	writeNode(&sb, n)
	return sb.String()
}

// writeNode dispatches on the concrete node type and appends a
// textual representation to sb.
func writeNode(sb *strings.Builder, n Node) {
	if n == nil {
		sb.WriteString("<nil>")
		return
	}
	// Detect typed-nil pointers (e.g. (*VarRef)(nil) wrapped in a Node
	// interface). The plain `n == nil` check above only catches an
	// untyped nil interface.
	if rv := reflect.ValueOf(n); rv.Kind() == reflect.Ptr && rv.IsNil() {
		sb.WriteString("<nil>")
		return
	}
	switch v := n.(type) {

	// -----------------------------------------------------------------------
	// node.go
	// -----------------------------------------------------------------------
	case *List:
		sb.WriteString("List{Items:[")
		for i, item := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
		sb.WriteString("]}")

	// -----------------------------------------------------------------------
	// literals.go
	// -----------------------------------------------------------------------
	case *StringLit:
		fmt.Fprintf(sb, "StringLit{Val:%q}", v.Val)
	case *NumberLit:
		fmt.Fprintf(sb, "NumberLit{Val:%s}", v.Val)
	case *BoolLit:
		fmt.Fprintf(sb, "BoolLit{Val:%t}", v.Val)
	case *NullLit:
		sb.WriteString("NullLit{}")
	case *MissingLit:
		sb.WriteString("MissingLit{}")
	case *DateLit:
		fmt.Fprintf(sb, "DateLit{Val:%s}", v.Val)
	case *TimeLit:
		sb.WriteString("TimeLit{Val:")
		sb.WriteString(v.Val)
		if v.Precision != nil {
			fmt.Fprintf(sb, " Precision:%d", *v.Precision)
		}
		if v.WithTimeZone {
			sb.WriteString(" WithTimeZone:true")
		}
		sb.WriteString("}")
	case *IonLit:
		fmt.Fprintf(sb, "IonLit{Text:%q}", v.Text)

	// -----------------------------------------------------------------------
	// exprs.go — operators and predicates
	// -----------------------------------------------------------------------
	case *BinaryExpr:
		fmt.Fprintf(sb, "BinaryExpr{Op:%s Left:", v.Op)
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		sb.WriteString("}")
	case *UnaryExpr:
		fmt.Fprintf(sb, "UnaryExpr{Op:%s Operand:", v.Op)
		writeNode(sb, v.Operand)
		sb.WriteString("}")
	case *InExpr:
		sb.WriteString("InExpr{Expr:")
		writeNode(sb, v.Expr)
		if v.Subquery != nil {
			sb.WriteString(" Subquery:")
			writeNode(sb, v.Subquery)
		} else {
			sb.WriteString(" List:[")
			for i, e := range v.List {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, e)
			}
			sb.WriteString("]")
		}
		fmt.Fprintf(sb, " Not:%t}", v.Not)
	case *BetweenExpr:
		sb.WriteString("BetweenExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" Low:")
		writeNode(sb, v.Low)
		sb.WriteString(" High:")
		writeNode(sb, v.High)
		fmt.Fprintf(sb, " Not:%t}", v.Not)
	case *LikeExpr:
		sb.WriteString("LikeExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" Pattern:")
		writeNode(sb, v.Pattern)
		if v.Escape != nil {
			sb.WriteString(" Escape:")
			writeNode(sb, v.Escape)
		}
		fmt.Fprintf(sb, " Not:%t}", v.Not)
	case *IsExpr:
		sb.WriteString("IsExpr{Expr:")
		writeNode(sb, v.Expr)
		fmt.Fprintf(sb, " Type:%s Not:%t}", v.Type, v.Not)

	// -----------------------------------------------------------------------
	// exprs.go — special-form expressions
	// -----------------------------------------------------------------------
	case *FuncCall:
		fmt.Fprintf(sb, "FuncCall{Name:%s", v.Name)
		if v.Quantifier != QuantifierNone {
			fmt.Fprintf(sb, " Quantifier:%s", v.Quantifier)
		}
		if v.Star {
			sb.WriteString(" Star:true")
		}
		sb.WriteString(" Args:[")
		for i, a := range v.Args {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString("]")
		if v.Over != nil {
			sb.WriteString(" Over:")
			writeNode(sb, v.Over)
		}
		sb.WriteString("}")
	case *CaseExpr:
		sb.WriteString("CaseExpr{")
		if v.Operand != nil {
			sb.WriteString("Operand:")
			writeNode(sb, v.Operand)
			sb.WriteString(" ")
		}
		sb.WriteString("Whens:[")
		for i, w := range v.Whens {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, w)
		}
		sb.WriteString("]")
		if v.Else != nil {
			sb.WriteString(" Else:")
			writeNode(sb, v.Else)
		}
		sb.WriteString("}")
	case *CaseWhen:
		sb.WriteString("CaseWhen{When:")
		writeNode(sb, v.When)
		sb.WriteString(" Then:")
		writeNode(sb, v.Then)
		sb.WriteString("}")
	case *CastExpr:
		fmt.Fprintf(sb, "CastExpr{Kind:%s Expr:", v.Kind)
		writeNode(sb, v.Expr)
		sb.WriteString(" AsType:")
		writeNode(sb, v.AsType)
		sb.WriteString("}")
	case *ExtractExpr:
		fmt.Fprintf(sb, "ExtractExpr{Field:%s From:", v.Field)
		writeNode(sb, v.From)
		sb.WriteString("}")
	case *TrimExpr:
		sb.WriteString("TrimExpr{")
		if v.Spec != TrimSpecNone {
			fmt.Fprintf(sb, "Spec:%s ", v.Spec)
		}
		if v.Sub != nil {
			sb.WriteString("Sub:")
			writeNode(sb, v.Sub)
			sb.WriteString(" ")
		}
		sb.WriteString("From:")
		writeNode(sb, v.From)
		sb.WriteString("}")
	case *SubstringExpr:
		sb.WriteString("SubstringExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" From:")
		writeNode(sb, v.From)
		if v.For != nil {
			sb.WriteString(" For:")
			writeNode(sb, v.For)
		}
		sb.WriteString("}")
	case *CoalesceExpr:
		sb.WriteString("CoalesceExpr{Args:[")
		for i, a := range v.Args {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString("]}")
	case *NullIfExpr:
		sb.WriteString("NullIfExpr{Left:")
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		sb.WriteString("}")
	case *WindowSpec:
		sb.WriteString("WindowSpec{PartitionBy:[")
		for i, e := range v.PartitionBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, e)
		}
		sb.WriteString("] OrderBy:[")
		for i, o := range v.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
		sb.WriteString("]}")

	// -----------------------------------------------------------------------
	// exprs.go — paths, vars, params, subqueries
	// -----------------------------------------------------------------------
	case *PathExpr:
		sb.WriteString("PathExpr{Root:")
		writeNode(sb, v.Root)
		sb.WriteString(" Steps:[")
		for i, s := range v.Steps {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString("]}")
	case *VarRef:
		fmt.Fprintf(sb, "VarRef{Name:%s", v.Name)
		if v.AtPrefixed {
			sb.WriteString(" AtPrefixed:true")
		}
		if v.CaseSensitive {
			sb.WriteString(" CaseSensitive:true")
		}
		sb.WriteString("}")
	case *ParamRef:
		sb.WriteString("ParamRef{}")
	case *SubLink:
		sb.WriteString("SubLink{Stmt:")
		writeNode(sb, v.Stmt)
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// exprs.go — collection literals
	// -----------------------------------------------------------------------
	case *ListLit:
		sb.WriteString("ListLit{Items:[")
		for i, e := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, e)
		}
		sb.WriteString("]}")
	case *BagLit:
		sb.WriteString("BagLit{Items:[")
		for i, e := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, e)
		}
		sb.WriteString("]}")
	case *TupleLit:
		sb.WriteString("TupleLit{Pairs:[")
		for i, p := range v.Pairs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *TuplePair:
		sb.WriteString("TuplePair{Key:")
		writeNode(sb, v.Key)
		sb.WriteString(" Value:")
		writeNode(sb, v.Value)
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// exprs.go — path steps
	// -----------------------------------------------------------------------
	case *DotStep:
		fmt.Fprintf(sb, "DotStep{Field:%s", v.Field)
		if v.CaseSensitive {
			sb.WriteString(" CaseSensitive:true")
		}
		sb.WriteString("}")
	case *AllFieldsStep:
		sb.WriteString("AllFieldsStep{}")
	case *IndexStep:
		sb.WriteString("IndexStep{Index:")
		writeNode(sb, v.Index)
		sb.WriteString("}")
	case *WildcardStep:
		sb.WriteString("WildcardStep{}")

	// -----------------------------------------------------------------------
	// tableexprs.go
	// -----------------------------------------------------------------------
	case *TableRef:
		fmt.Fprintf(sb, "TableRef{Name:%s", v.Name)
		if v.Schema != "" {
			fmt.Fprintf(sb, " Schema:%s", v.Schema)
		}
		if v.CaseSensitive {
			sb.WriteString(" CaseSensitive:true")
		}
		sb.WriteString("}")
	case *AliasedSource:
		sb.WriteString("AliasedSource{Source:")
		writeNode(sb, v.Source)
		writeOptString(sb, " As:", v.As)
		writeOptString(sb, " At:", v.At)
		writeOptString(sb, " By:", v.By)
		sb.WriteString("}")
	case *JoinExpr:
		fmt.Fprintf(sb, "JoinExpr{Kind:%s Left:", v.Kind)
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		if v.On != nil {
			sb.WriteString(" On:")
			writeNode(sb, v.On)
		}
		sb.WriteString("}")
	case *UnpivotExpr:
		sb.WriteString("UnpivotExpr{Source:")
		writeNode(sb, v.Source)
		writeOptString(sb, " As:", v.As)
		writeOptString(sb, " At:", v.At)
		writeOptString(sb, " By:", v.By)
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// types.go
	// -----------------------------------------------------------------------
	case *TypeRef:
		fmt.Fprintf(sb, "TypeRef{Name:%s", v.Name)
		if len(v.Args) > 0 {
			sb.WriteString(" Args:[")
			for i, a := range v.Args {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(strconv.Itoa(a))
			}
			sb.WriteString("]")
		}
		if v.WithTimeZone {
			sb.WriteString(" WithTimeZone:true")
		}
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// stmts.go — top-level statements
	// -----------------------------------------------------------------------
	case *SelectStmt:
		writeSelectStmt(sb, v)
	case *SetOpStmt:
		fmt.Fprintf(sb, "SetOpStmt{Op:%s", v.Op)
		if v.Quantifier != QuantifierNone {
			fmt.Fprintf(sb, " Quantifier:%s", v.Quantifier)
		}
		if v.Outer {
			sb.WriteString(" Outer:true")
		}
		sb.WriteString(" Left:")
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		sb.WriteString("}")
	case *ExplainStmt:
		sb.WriteString("ExplainStmt{Inner:")
		writeNode(sb, v.Inner)
		sb.WriteString("}")
	case *InsertStmt:
		writeDmlStmt(sb, "InsertStmt", v.Target, v.AsAlias, v.Value, v.Pos, v.OnConflict, v.Returning)
	case *UpdateStmt:
		sb.WriteString("UpdateStmt{Source:")
		writeNode(sb, v.Source)
		sb.WriteString(" Sets:[")
		for i, s := range v.Sets {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString("]")
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		if v.Returning != nil {
			sb.WriteString(" Returning:")
			writeNode(sb, v.Returning)
		}
		sb.WriteString("}")
	case *DeleteStmt:
		sb.WriteString("DeleteStmt{Source:")
		writeNode(sb, v.Source)
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		if v.Returning != nil {
			sb.WriteString(" Returning:")
			writeNode(sb, v.Returning)
		}
		sb.WriteString("}")
	case *UpsertStmt:
		writeDmlStmt(sb, "UpsertStmt", v.Target, v.AsAlias, v.Value, nil, v.OnConflict, v.Returning)
	case *ReplaceStmt:
		writeDmlStmt(sb, "ReplaceStmt", v.Target, v.AsAlias, v.Value, nil, v.OnConflict, v.Returning)
	case *RemoveStmt:
		sb.WriteString("RemoveStmt{Path:")
		writeNode(sb, v.Path)
		sb.WriteString("}")
	case *CreateTableStmt:
		sb.WriteString("CreateTableStmt{Name:")
		writeNode(sb, v.Name)
		sb.WriteString("}")
	case *CreateIndexStmt:
		sb.WriteString("CreateIndexStmt{Table:")
		writeNode(sb, v.Table)
		sb.WriteString(" Paths:[")
		for i, p := range v.Paths {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *DropTableStmt:
		sb.WriteString("DropTableStmt{Name:")
		writeNode(sb, v.Name)
		sb.WriteString("}")
	case *DropIndexStmt:
		sb.WriteString("DropIndexStmt{Index:")
		writeNode(sb, v.Index)
		sb.WriteString(" Table:")
		writeNode(sb, v.Table)
		sb.WriteString("}")
	case *ExecStmt:
		sb.WriteString("ExecStmt{Name:")
		writeNode(sb, v.Name)
		sb.WriteString(" Args:[")
		for i, a := range v.Args {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString("]}")

	// -----------------------------------------------------------------------
	// stmts.go — clause and DML helpers
	// -----------------------------------------------------------------------
	case *TargetEntry:
		sb.WriteString("TargetEntry{Expr:")
		writeNode(sb, v.Expr)
		writeOptString(sb, " Alias:", v.Alias)
		sb.WriteString("}")
	case *PivotProjection:
		sb.WriteString("PivotProjection{Value:")
		writeNode(sb, v.Value)
		sb.WriteString(" At:")
		writeNode(sb, v.At)
		sb.WriteString("}")
	case *LetBinding:
		sb.WriteString("LetBinding{Expr:")
		writeNode(sb, v.Expr)
		fmt.Fprintf(sb, " Alias:%s}", v.Alias)
	case *GroupByClause:
		sb.WriteString("GroupByClause{")
		if v.Partial {
			sb.WriteString("Partial:true ")
		}
		sb.WriteString("Items:[")
		for i, it := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, it)
		}
		sb.WriteString("]")
		writeOptString(sb, " GroupAs:", v.GroupAs)
		sb.WriteString("}")
	case *GroupByItem:
		sb.WriteString("GroupByItem{Expr:")
		writeNode(sb, v.Expr)
		writeOptString(sb, " Alias:", v.Alias)
		sb.WriteString("}")
	case *OrderByItem:
		sb.WriteString("OrderByItem{Expr:")
		writeNode(sb, v.Expr)
		fmt.Fprintf(sb, " Desc:%t", v.Desc)
		if v.NullsExplicit {
			fmt.Fprintf(sb, " NullsFirst:%t", v.NullsFirst)
		}
		sb.WriteString("}")
	case *SetAssignment:
		sb.WriteString("SetAssignment{Target:")
		writeNode(sb, v.Target)
		sb.WriteString(" Value:")
		writeNode(sb, v.Value)
		sb.WriteString("}")
	case *OnConflict:
		sb.WriteString("OnConflict{")
		if v.Target != nil {
			sb.WriteString("Target:")
			writeNode(sb, v.Target)
			sb.WriteString(" ")
		}
		fmt.Fprintf(sb, "Action:%s", v.Action)
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		sb.WriteString("}")
	case *OnConflictTarget:
		sb.WriteString("OnConflictTarget{")
		if v.ConstraintName != "" {
			fmt.Fprintf(sb, "ConstraintName:%s", v.ConstraintName)
		} else {
			sb.WriteString("Cols:[")
			for i, c := range v.Cols {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, c)
			}
			sb.WriteString("]")
		}
		sb.WriteString("}")
	case *ReturningClause:
		sb.WriteString("ReturningClause{Items:[")
		for i, it := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, it)
		}
		sb.WriteString("]}")
	case *ReturningItem:
		fmt.Fprintf(sb, "ReturningItem{Status:%s Mapping:%s", v.Status, v.Mapping)
		if v.Star {
			sb.WriteString(" Star:true")
		}
		if v.Expr != nil {
			sb.WriteString(" Expr:")
			writeNode(sb, v.Expr)
		}
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// patterns.go
	// -----------------------------------------------------------------------
	case *MatchExpr:
		sb.WriteString("MatchExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" Patterns:[")
		for i, p := range v.Patterns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *GraphPattern:
		sb.WriteString("GraphPattern{")
		if v.Selector != nil {
			sb.WriteString("Selector:")
			writeNode(sb, v.Selector)
			sb.WriteString(" ")
		}
		if v.Restrictor != PatternRestrictorNone {
			fmt.Fprintf(sb, "Restrictor:%s ", v.Restrictor)
		}
		if v.Variable != nil {
			sb.WriteString("Variable:")
			writeNode(sb, v.Variable)
			sb.WriteString(" ")
		}
		sb.WriteString("Parts:[")
		for i, p := range v.Parts {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *NodePattern:
		sb.WriteString("NodePattern{")
		first := true
		sep := func() {
			if !first {
				sb.WriteString(" ")
			}
			first = false
		}
		if v.Variable != nil {
			sep()
			sb.WriteString("Variable:")
			writeNode(sb, v.Variable)
		}
		if len(v.Labels) > 0 {
			sep()
			fmt.Fprintf(sb, "Labels:[%s]", strings.Join(v.Labels, " "))
		}
		if v.Where != nil {
			sep()
			sb.WriteString("Where:")
			writeNode(sb, v.Where)
		}
		sb.WriteString("}")
	case *EdgePattern:
		fmt.Fprintf(sb, "EdgePattern{Direction:%s", v.Direction)
		if v.Variable != nil {
			sb.WriteString(" Variable:")
			writeNode(sb, v.Variable)
		}
		if len(v.Labels) > 0 {
			fmt.Fprintf(sb, " Labels:[%s]", strings.Join(v.Labels, " "))
		}
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		if v.Quantifier != nil {
			sb.WriteString(" Quantifier:")
			writeNode(sb, v.Quantifier)
		}
		sb.WriteString("}")
	case *PatternQuantifier:
		fmt.Fprintf(sb, "PatternQuantifier{Min:%d Max:%d}", v.Min, v.Max)
	case *PatternSelector:
		fmt.Fprintf(sb, "PatternSelector{Kind:%s", v.Kind)
		if v.Kind == SelectorKindShortestK {
			fmt.Fprintf(sb, " K:%d", v.K)
		}
		sb.WriteString("}")

	default:
		fmt.Fprintf(sb, "<unknown:%T>", v)
	}
}

// writeSelectStmt is split out because SelectStmt has 13 fields and would
// dominate the main switch arm.
func writeSelectStmt(sb *strings.Builder, s *SelectStmt) {
	sb.WriteString("SelectStmt{")
	first := true
	add := func(label string) {
		if !first {
			sb.WriteString(" ")
		}
		first = false
		sb.WriteString(label)
	}
	if s.Quantifier != QuantifierNone {
		add(fmt.Sprintf("Quantifier:%s", s.Quantifier))
	}
	if s.Star {
		add("Star:true")
	}
	if s.Value != nil {
		add("Value:")
		writeNode(sb, s.Value)
	}
	if s.Pivot != nil {
		add("Pivot:")
		writeNode(sb, s.Pivot)
	}
	if len(s.Targets) > 0 {
		add("Targets:[")
		for i, t := range s.Targets {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
		sb.WriteString("]")
	}
	if s.From != nil {
		add("From:")
		writeNode(sb, s.From)
	}
	if len(s.Let) > 0 {
		add("Let:[")
		for i, l := range s.Let {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, l)
		}
		sb.WriteString("]")
	}
	if s.Where != nil {
		add("Where:")
		writeNode(sb, s.Where)
	}
	if s.GroupBy != nil {
		add("GroupBy:")
		writeNode(sb, s.GroupBy)
	}
	if s.Having != nil {
		add("Having:")
		writeNode(sb, s.Having)
	}
	if len(s.OrderBy) > 0 {
		add("OrderBy:[")
		for i, o := range s.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
		sb.WriteString("]")
	}
	if s.Limit != nil {
		add("Limit:")
		writeNode(sb, s.Limit)
	}
	if s.Offset != nil {
		add("Offset:")
		writeNode(sb, s.Offset)
	}
	sb.WriteString("}")
}

// writeDmlStmt is shared by InsertStmt, UpsertStmt, and ReplaceStmt — they
// have nearly the same shape (target, alias, value, pos, on-conflict,
// returning). For UpsertStmt and ReplaceStmt, callers pass nil for pos
// because the grammar's `replaceCommand` (line 121) and `upsertCommand`
// (line 125) do not have an `AT pos` clause.
func writeDmlStmt(sb *strings.Builder, name string, target TableExpr, alias *string, value ExprNode, pos ExprNode, oc *OnConflict, ret *ReturningClause) {
	fmt.Fprintf(sb, "%s{Target:", name)
	writeNode(sb, target)
	writeOptString(sb, " AsAlias:", alias)
	sb.WriteString(" Value:")
	writeNode(sb, value)
	if pos != nil {
		sb.WriteString(" Pos:")
		writeNode(sb, pos)
	}
	if oc != nil {
		sb.WriteString(" OnConflict:")
		writeNode(sb, oc)
	}
	if ret != nil {
		sb.WriteString(" Returning:")
		writeNode(sb, ret)
	}
	sb.WriteString("}")
}

// writeOptString appends a label + value if the optional string is non-nil.
func writeOptString(sb *strings.Builder, label string, s *string) {
	if s == nil {
		return
	}
	sb.WriteString(label)
	sb.WriteString(*s)
}
