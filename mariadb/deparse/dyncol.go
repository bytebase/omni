package deparse

import (
	"strings"

	ast "github.com/bytebase/omni/mariadb/ast"
)

// writeDynColCall writes a dynamic column function the way MariaDB stores it.
// Measured on mariadb:11.8.8 through SHOW CREATE VIEW:
//
//	COLUMN_CREATE('a', 1 AS INT, 'b', 2)   column_create('a',1 AS int,'b',2)
//	COLUMN_ADD(b, 'x', 1 AS REAL)          column_add(`t`.`b`,'x',1 AS double)
//	COLUMN_DELETE(b, 'a', 'b')             column_add(`t`.`b`,'a',NULL AS int,'b',NULL AS int)
//	COLUMN_GET(b, 'a' AS INT)              column_get(`t`.`b`,'a' as signed)
//
// COLUMN_GET's type is a CAST type and is written as CAST writes it.
func writeDynColCall(b *strings.Builder, n *ast.FuncCallExpr) {
	name := strings.ToUpper(n.Name)
	if name == "COLUMN_DELETE" {
		b.WriteString("column_add(")
		writeExpr(b, n.Args[0])
		for _, key := range n.Args[1:] {
			b.WriteString(",")
			writeExpr(b, key)
			b.WriteString(",NULL AS int")
		}
		b.WriteString(")")
		return
	}
	b.WriteString(strings.ToLower(name))
	b.WriteString("(")
	for i, arg := range n.Args {
		if i > 0 {
			b.WriteString(",")
		}
		writeExpr(b, arg)
		dt := n.ArgTypes[i]
		if dt == nil {
			continue
		}
		if name == "COLUMN_GET" {
			b.WriteString(" as ")
			b.WriteString(deparseDataType(dt))
		} else {
			b.WriteString(" AS ")
			b.WriteString(deparseDynColType(dt))
		}
	}
	b.WriteString(")")
}

// deparseDynColType writes a dynamic column's value type as the server
// stores it, which keeps the type and drops its size: INT and UNSIGNED INT
// stay, DOUBLE, REAL and FLOAT are double, DECIMAL(10,2) is decimal, TIME(3)
// and DATETIME(3) lose the scale. A string type names its character set and
// collation and ends in a space: mariadb:11.8.8 stores 'x' AS CHAR as
// 'x' AS char charset utf8mb4 collate utf8mb4_uca1400_ai_ci , taking the
// connection's character set, and NCHAR as utf8mb3 with utf8mb3_general_ci.
func deparseDynColType(dt *ast.DataType) string {
	switch strings.ToUpper(dt.Name) {
	case "INT":
		if dt.Unsigned {
			return "unsigned int"
		}
		return "int"
	case "DOUBLE", "REAL", "FLOAT":
		return "double"
	case "NCHAR":
		return "char charset utf8mb3 collate utf8mb3_general_ci "
	case "CHAR":
		cs := strings.ToLower(dt.Charset)
		switch cs {
		case "":
			cs = "utf8mb4"
		case "utf8":
			cs = "utf8mb3"
		}
		if cs == "binary" {
			// CHAR BYTE: char charset binary , with no collation.
			return "char charset binary "
		}
		coll := strings.ToLower(dt.Collate)
		switch {
		case coll != "":
		case dt.Binary:
			coll = cs + "_bin"
		default:
			coll = DefaultCollations[cs]
		}
		return "char charset " + cs + " collate " + coll + " "
	default:
		// DECIMAL, DATE, TIME and DATETIME.
		return strings.ToLower(dt.Name)
	}
}

// writePercentileCall writes PERCENTILE_CONT, PERCENTILE_DISC and MEDIAN as
// MariaDB stores them. Measured on mariadb:11.8.8 through SHOW CREATE VIEW:
//
//	PERCENTILE_DISC(0.5) WITHIN GROUP (ORDER BY a DESC) OVER (PARTITION BY b)
//	    percentile_disc(0.5) within group ( order by `t`.`a` desc) over ( partition by `t`.`b`)
//	MEDIAN(a) OVER ()
//	    percentile_cont(0.5) within group ( order by `t`.`a` desc) over ()
//
// MEDIAN is PERCENTILE_CONT(0.5) with its argument for the order, which the
// server writes descending.
func writePercentileCall(b *strings.Builder, n *ast.FuncCallExpr) {
	var item *ast.OrderByItem
	if n.WithinGroup {
		b.WriteString(strings.ToLower(n.Name))
		b.WriteString("(")
		writeExpr(b, n.Args[0])
		b.WriteString(")")
		item = n.OrderBy[0]
	} else {
		b.WriteString("percentile_cont(0.5)")
		item = &ast.OrderByItem{Expr: n.Args[0], Desc: true}
	}
	b.WriteString(" within group ( order by ")
	writeExpr(b, item.Expr)
	if item.Desc {
		b.WriteString(" desc")
	}
	b.WriteString(") ")
	writeWindowDef(b, n.Over)
}
