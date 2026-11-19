package catalog

import nodes "github.com/bytebase/omni/redshift/ast"

// cloneNodeList copies a list and the nodes in it, so a clone of a catalog
// shares no storage option with its source.
func cloneNodeList(list *nodes.List) *nodes.List {
	if list == nil {
		return nil
	}
	items := make([]nodes.Node, len(list.Items))
	for i, item := range list.Items {
		items[i] = cloneNode(item)
	}
	return &nodes.List{Items: items}
}

// cloneNode copies the node kinds a materialized view's options hold: the
// DefElems and the strings, numbers, lists and expressions beside them. A node
// of any other kind is shared, as the catalog treats its options as data.
func cloneNode(node nodes.Node) nodes.Node {
	switch n := node.(type) {
	case nil:
		return nil
	case *nodes.List:
		return cloneNodeList(n)
	case *nodes.DefElem:
		clone := *n
		clone.Arg = cloneNode(n.Arg)
		return &clone
	case *nodes.String:
		clone := *n
		return &clone
	case *nodes.Integer:
		clone := *n
		return &clone
	case *nodes.Float:
		clone := *n
		return &clone
	case *nodes.Boolean:
		clone := *n
		return &clone
	case *nodes.A_Const:
		clone := *n
		clone.Val = cloneNode(n.Val)
		return &clone
	case *nodes.ColumnRef:
		clone := *n
		clone.Fields = cloneNodeList(n.Fields)
		return &clone
	case *nodes.FuncCall:
		clone := *n
		clone.Funcname = cloneNodeList(n.Funcname)
		clone.Args = cloneNodeList(n.Args)
		return &clone
	}
	return node
}
