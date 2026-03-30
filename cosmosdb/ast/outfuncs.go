package ast

import (
	"fmt"
	"strings"
)

// NodeToString converts a Node to its string representation.
func NodeToString(node Node) string {
	if node == nil {
		return "<>"
	}
	var sb strings.Builder
	writeNode(&sb, node, 0)
	return sb.String()
}

func indent(sb *strings.Builder, level int) {
	for range level {
		sb.WriteString("  ")
	}
}

func writeLoc(sb *strings.Builder, loc Loc) {
	fmt.Fprintf(sb, " :loc_start %d :loc_end %d", loc.Start, loc.End)
}

func writeNode(sb *strings.Builder, node Node, level int) {
	if node == nil {
		sb.WriteString("<>")
		return
	}

	switch n := node.(type) {
	case *List:
		sb.WriteString("(")
		for i, item := range n.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item, level)
		}
		sb.WriteString(")")

	// Literals
	case *StringLit:
		fmt.Fprintf(sb, "{STRLIT :val %q", n.Val)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *NumberLit:
		fmt.Fprintf(sb, "{NUMLIT :val %q", n.Val)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BoolLit:
		fmt.Fprintf(sb, "{BOOLLIT :val %t", n.Val)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *NullLit:
		sb.WriteString("{NULL")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *UndefinedLit:
		sb.WriteString("{UNDEFINED")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *InfinityLit:
		sb.WriteString("{INFINITY")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *NanLit:
		sb.WriteString("{NAN")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	// Wrapper
	case *RawStmt:
		sb.WriteString("{RAWSTMT")
		fmt.Fprintf(sb, " :stmt_location %d :stmt_len %d", n.StmtLocation, n.StmtLen)
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":stmt ")
		writeNode(sb, n.Stmt, level+1)
		sb.WriteString("}")

	// Statement
	case *SelectStmt:
		writeSelectStmt(sb, n, level)

	// Clause nodes
	case *TargetEntry:
		sb.WriteString("{TARGET :expr ")
		writeNode(sb, n.Expr, level)
		if n.Alias != nil {
			fmt.Fprintf(sb, " :alias %q", *n.Alias)
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *SortExpr:
		sb.WriteString("{SORT :expr ")
		writeNode(sb, n.Expr, level)
		if n.Desc {
			sb.WriteString(" :desc true")
		}
		if n.Rank {
			sb.WriteString(" :rank true")
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *OffsetLimitClause:
		sb.WriteString("{OFFSETLIMIT :offset ")
		writeNode(sb, n.Offset, level)
		sb.WriteString(" :limit ")
		writeNode(sb, n.Limit, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *JoinExpr:
		sb.WriteString("{JOIN :source ")
		writeNode(sb, n.Source, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	// Table expressions
	case *ContainerRef:
		if n.Root {
			sb.WriteString("{CONTAINER :root true")
		} else {
			fmt.Fprintf(sb, "{CONTAINER :name %q", n.Name)
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *AliasedTableExpr:
		sb.WriteString("{ALIASED :source ")
		writeNode(sb, n.Source, level)
		fmt.Fprintf(sb, " :alias %q", n.Alias)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ArrayIterationExpr:
		fmt.Fprintf(sb, "{ARRAYITER :alias %q :source ", n.Alias)
		writeNode(sb, n.Source, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *SubqueryExpr:
		sb.WriteString("{SUBQUERY\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	// Expression nodes
	case *ColumnRef:
		fmt.Fprintf(sb, "{COLREF :name %q", n.Name)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *DotAccessExpr:
		sb.WriteString("{DOT :expr ")
		writeNode(sb, n.Expr, level)
		fmt.Fprintf(sb, " :property %q", n.Property)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BracketAccessExpr:
		sb.WriteString("{BRACKET :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :index ")
		writeNode(sb, n.Index, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BinaryExpr:
		fmt.Fprintf(sb, "{BINEXPR :op %q :left ", n.Op)
		writeNode(sb, n.Left, level)
		sb.WriteString(" :right ")
		writeNode(sb, n.Right, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *UnaryExpr:
		fmt.Fprintf(sb, "{UNARYEXPR :op %q :operand ", n.Op)
		writeNode(sb, n.Operand, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *TernaryExpr:
		sb.WriteString("{TERNARY :cond ")
		writeNode(sb, n.Cond, level)
		sb.WriteString(" :then ")
		writeNode(sb, n.Then, level)
		sb.WriteString(" :else ")
		writeNode(sb, n.Else, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *InExpr:
		sb.WriteString("{IN")
		if n.Not {
			sb.WriteString(" :not true")
		}
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :list (")
		for i, item := range n.List {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item, level)
		}
		sb.WriteString(")")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *BetweenExpr:
		sb.WriteString("{BETWEEN")
		if n.Not {
			sb.WriteString(" :not true")
		}
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :low ")
		writeNode(sb, n.Low, level)
		sb.WriteString(" :high ")
		writeNode(sb, n.High, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *LikeExpr:
		sb.WriteString("{LIKE")
		if n.Not {
			sb.WriteString(" :not true")
		}
		sb.WriteString(" :expr ")
		writeNode(sb, n.Expr, level)
		sb.WriteString(" :pattern ")
		writeNode(sb, n.Pattern, level)
		if n.Escape != nil {
			sb.WriteString(" :escape ")
			writeNode(sb, n.Escape, level)
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *FuncCall:
		fmt.Fprintf(sb, "{FUNCCALL :name %q", n.Name)
		if n.Star {
			sb.WriteString(" :star true")
		}
		if len(n.Args) > 0 {
			sb.WriteString(" :args (")
			for i, arg := range n.Args {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, arg, level)
			}
			sb.WriteString(")")
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *UDFCall:
		fmt.Fprintf(sb, "{UDFCALL :name %q", n.Name)
		if len(n.Args) > 0 {
			sb.WriteString(" :args (")
			for i, arg := range n.Args {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, arg, level)
			}
			sb.WriteString(")")
		}
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ExistsExpr:
		sb.WriteString("{EXISTS\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ArrayExpr:
		sb.WriteString("{ARRAY\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *CreateArrayExpr:
		sb.WriteString("{CREATEARRAY :elements (")
		for i, elem := range n.Elements {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, elem, level)
		}
		sb.WriteString(")")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *CreateObjectExpr:
		sb.WriteString("{CREATEOBJECT :fields (")
		for i, f := range n.Fields {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, f, level)
		}
		sb.WriteString(")")
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ObjectFieldPair:
		fmt.Fprintf(sb, "{FIELD :key %q :value ", n.Key)
		writeNode(sb, n.Value, level)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *ParamRef:
		fmt.Fprintf(sb, "{PARAM :name %q", n.Name)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")
	case *SubLink:
		sb.WriteString("{SUBLINK\n")
		indent(sb, level+1)
		sb.WriteString(":select ")
		writeNode(sb, n.Select, level+1)
		writeLoc(sb, n.Loc)
		sb.WriteString("}")

	default:
		fmt.Fprintf(sb, "{UNKNOWN %T}", node)
	}
}

func writeSelectStmt(sb *strings.Builder, n *SelectStmt, level int) {
	sb.WriteString("{SELECT")
	if n.Top != nil {
		fmt.Fprintf(sb, " :top %d", *n.Top)
	}
	if n.Distinct {
		sb.WriteString(" :distinct true")
	}
	if n.Value {
		sb.WriteString(" :value true")
	}
	if n.Star {
		sb.WriteString(" :star true")
	}
	if len(n.Targets) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":targets (")
		for i, t := range n.Targets {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t, level+1)
		}
		sb.WriteString(")")
	}
	if n.From != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":from ")
		writeNode(sb, n.From, level+1)
	}
	if len(n.Joins) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":joins (")
		for i, j := range n.Joins {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, j, level+1)
		}
		sb.WriteString(")")
	}
	if n.Where != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":where ")
		writeNode(sb, n.Where, level+1)
	}
	if len(n.GroupBy) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":group_by (")
		for i, g := range n.GroupBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, g, level+1)
		}
		sb.WriteString(")")
	}
	if n.Having != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":having ")
		writeNode(sb, n.Having, level+1)
	}
	if len(n.OrderBy) > 0 {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":order_by (")
		for i, o := range n.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o, level+1)
		}
		sb.WriteString(")")
	}
	if n.OffsetLimit != nil {
		sb.WriteString("\n")
		indent(sb, level+1)
		sb.WriteString(":offset_limit ")
		writeNode(sb, n.OffsetLimit, level+1)
	}
	writeLoc(sb, n.Loc)
	sb.WriteString("}")
}
