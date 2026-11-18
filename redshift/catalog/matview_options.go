package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// deparseMatViewOptions renders the storage a CREATE MATERIALIZED VIEW named,
// so a migration re-creates the view as it was written: the USING access
// method, then BACKUP, DISTSTYLE, DISTKEY, SORTKEY and ENCODE, AUTO REFRESH
// with its interval, and the Iceberg form's LOCATION, PARTITIONED BY and
// TABLE PROPERTIES.
//
// Grammar reference: materializedViewCreateSqlCommand in
// redshift/testdata/RedshiftDocs.g4, which is what redshift/parser reads. A
// name the parser does not write there is left out rather than guessed at.
func deparseMatViewOptions(rel *Relation) string {
	var sb strings.Builder
	if rel.MatViewAccessMethod != "" {
		fmt.Fprintf(&sb, " USING %s", strings.ToUpper(rel.MatViewAccessMethod))
	}
	elems := rel.MatViewOptions
	each := func(name string, render func(elem *nodes.DefElem) string) {
		for _, item := range defElems(elems) {
			elem, ok := item.(*nodes.DefElem)
			if !ok || elem.Defname != name {
				continue
			}
			sb.WriteString(render(elem))
		}
	}

	each("backup", func(elem *nodes.DefElem) string {
		return " BACKUP " + strings.ToUpper(stringVal(elem.Arg))
	})
	each("diststyle", func(elem *nodes.DefElem) string {
		return " DISTSTYLE " + strings.ToUpper(stringVal(elem.Arg))
	})
	each("distkey", func(elem *nodes.DefElem) string {
		// The attribute parser takes a bare DISTKEY as well as DISTKEY (col).
		list, ok := elem.Arg.(*nodes.List)
		if !ok {
			return " DISTKEY"
		}
		return " DISTKEY (" + strings.Join(quotedIdentItems(list), ", ") + ")"
	})
	// SORTKEY's style is a DefElem of its own, written before it.
	each("sortstyle", func(elem *nodes.DefElem) string {
		return " " + strings.ToUpper(stringVal(elem.Arg))
	})
	each("sortkey", func(elem *nodes.DefElem) string {
		if list, ok := elem.Arg.(*nodes.List); ok {
			return " SORTKEY (" + strings.Join(quotedIdentItems(list), ", ") + ")"
		}
		return " SORTKEY " + strings.ToUpper(stringVal(elem.Arg))
	})
	each("encode", func(elem *nodes.DefElem) string {
		return " ENCODE " + strings.ToUpper(stringVal(elem.Arg))
	})
	each("auto_refresh", func(elem *nodes.DefElem) string {
		if b, ok := elem.Arg.(*nodes.Boolean); ok && b.Boolval {
			return " AUTO REFRESH YES"
		}
		return " AUTO REFRESH NO"
	})
	each("auto_refresh_interval", func(elem *nodes.DefElem) string {
		return fmt.Sprintf(" INTERVAL %d", intVal(elem.Arg))
	})
	each("location", func(elem *nodes.DefElem) string {
		return " LOCATION " + quoteLiteral(stringVal(elem.Arg))
	})
	each("partitioned_by", func(elem *nodes.DefElem) string {
		list, ok := elem.Arg.(*nodes.List)
		if !ok {
			return ""
		}
		transforms := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			transforms = append(transforms, deparsePartitionTransform(item))
		}
		return " PARTITIONED BY (" + strings.Join(transforms, ", ") + ")"
	})
	each("table_properties", func(elem *nodes.DefElem) string {
		list, ok := elem.Arg.(*nodes.List)
		if !ok {
			return ""
		}
		pairs := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			prop, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			pairs = append(pairs, quoteLiteral(prop.Defname)+" = "+quoteLiteral(stringVal(prop.Arg)))
		}
		return " TABLE PROPERTIES (" + strings.Join(pairs, ", ") + ")"
	})
	return sb.String()
}

// deparsePartitionTransform renders one partition_transform back: a column, or
// the YEAR, MONTH, DAY, HOUR, BUCKET or TRUNCATE call the docs list.
func deparsePartitionTransform(node nodes.Node) string {
	switch n := node.(type) {
	case *nodes.ColumnRef:
		if n.Fields != nil && len(n.Fields.Items) == 1 {
			return quoteIdentAlways(stringVal(n.Fields.Items[0]))
		}
	case *nodes.FuncCall:
		args := make([]string, 0, len(n.Args.Items))
		for _, arg := range n.Args.Items {
			if ref, ok := arg.(*nodes.ColumnRef); ok && ref.Fields != nil && len(ref.Fields.Items) == 1 {
				args = append(args, quoteIdentAlways(stringVal(ref.Fields.Items[0])))
				continue
			}
			args = append(args, deparseAConst(arg))
		}
		return strings.ToUpper(stringVal(n.Funcname.Items[0])) + "(" + strings.Join(args, ", ") + ")"
	}
	return ""
}

// quotedIdentItems renders a name list as the quoted identifiers migration DDL
// writes, so DISTKEY ("UserID") keeps the column's spelling.
func quotedIdentItems(list *nodes.List) []string {
	names := stringListItems(list)
	for i, name := range names {
		names[i] = quoteIdentAlways(name)
	}
	return names
}

// defElems returns options' items, or nil when it holds none.
func defElems(options *nodes.List) []nodes.Node {
	if options == nil {
		return nil
	}
	return options.Items
}
