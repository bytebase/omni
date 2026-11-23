package analysis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bytebase/omni/enginetest/span"
)

// planNode is one operator of a logical plan (EXPLAIN LOGICAL).
type planNode struct {
	indent   int
	text     string // the operator's line, after "- "
	outputs  []int  // the slots after "=>"
	children []*planNode
	// defs are the slots the operator defines, "n:name := expression".
	defs map[int]string
}

var (
	// slotRefRe matches a slot reference, "3:amount" or "3: amount".
	slotRefRe = regexp.MustCompile(`\b(\d+): ?([A-Za-z_][A-Za-z0-9_]*)`)
	// slotOnlyRe matches an expression that is one slot reference.
	slotOnlyRe = regexp.MustCompile(`^(\d+): ?[A-Za-z_][A-Za-z0-9_]*$`)
	// scanRe matches a table's scan, "SCAN [orders] => [...]" or
	// "SCHEMA-SCAN [tables] => [...]".
	scanRe = regexp.MustCompile(`^(?:SCHEMA-)?SCAN \[([^\]]+)\]`)
)

// parsePlan reads a logical plan's lines into its operators.
func parsePlan(lines []string) (*planNode, error) {
	root := &planNode{indent: -1}
	stack := []*planNode{root}
	var last *planNode
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmed)
		if strings.HasPrefix(trimmed, "- ") {
			n := &planNode{indent: indent, text: strings.TrimPrefix(trimmed, "- "), defs: map[int]string{}}
			if i := strings.LastIndex(n.text, "=> ["); i >= 0 && strings.HasSuffix(n.text, "]") {
				for _, item := range splitSlots(n.text[i+4 : len(n.text)-1]) {
					m := regexp.MustCompile(`^(\d+): ?`).FindStringSubmatch(item)
					if m == nil {
						// A table function lists its slots as {1,6}; the
						// top Output list is the one that must be read.
						if strings.HasPrefix(n.text, "Output =>") {
							return nil, fmt.Errorf("an output %q", item)
						}
						n.outputs = nil
						break
					}
					var slot int
					fmt.Sscan(m[1], &slot)
					n.outputs = append(n.outputs, slot)
				}
			}
			for len(stack) > 1 && stack[len(stack)-1].indent >= indent {
				stack = stack[:len(stack)-1]
			}
			parent := stack[len(stack)-1]
			parent.children = append(parent.children, n)
			stack = append(stack, n)
			last = n
			continue
		}
		if i := strings.Index(trimmed, " := "); i > 0 && last != nil {
			m := regexp.MustCompile(`^(\d+):`).FindStringSubmatch(trimmed[:i])
			if m != nil {
				var slot int
				fmt.Sscan(m[1], &slot)
				last.defs[slot] = trimmed[i+4:]
			}
		}
	}
	if len(root.children) != 1 || !strings.HasPrefix(root.children[0].text, "Output =>") {
		return nil, fmt.Errorf("no Output operator")
	}
	return root.children[0], nil
}

// splitSlots splits an output list on the commas between its items, not on
// those inside an item's parentheses.
func splitSlots(list string) []string {
	var out []string
	depth, start := 0, 0
	for i, c := range list {
		switch c {
		case '(', '[':
			depth++
		case ')', ']':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(list[start:i]))
				start = i + 1
			}
		}
	}
	if s := strings.TrimSpace(list[start:]); s != "" {
		out = append(out, s)
	}
	return out
}

// slotSource is what a slot reads: table columns, as database.table.column,
// and whether that is all of it and it is one of them unchanged.
type slotSource struct {
	columns  map[string]bool
	plain    bool
	complete bool
}

// planLineage traces each output slot of a logical plan to the scanned
// columns it reads. tables are the database.table names the query dump
// lists, which name each scan's database: a scan names its table alone.
func planLineage(lines []string, tables []string) ([][]string, error) {
	out, err := parsePlan(lines)
	if err != nil {
		return nil, err
	}
	databases := map[string][]string{}
	for _, t := range tables {
		db, name, _ := strings.Cut(t, ".")
		databases[name] = append(databases[name], db)
	}
	// What defines each slot: a scan's column, an expression, or a set
	// operation's arms.
	scanned := map[int]string{}
	defs := map[int]string{}
	setOps := map[int][]int{}
	var walk func(n *planNode) error
	walk = func(n *planNode) error {
		if m := scanRe.FindStringSubmatch(n.text); m != nil {
			table := m[1]
			db := ""
			if dbs := databases[table]; len(dbs) == 1 {
				db = dbs[0]
			}
			// A scan's columns are its outputs and the slots its own
			// expressions read that nothing else defines.
			name := func(slot int, column string) {
				if db == "" {
					scanned[slot] = ""
					return
				}
				scanned[slot] = db + "." + table + "." + column
			}
			_, after, _ := strings.Cut(n.text, "=> [")
			for _, m := range slotRefRe.FindAllStringSubmatch(after, -1) {
				var slot int
				fmt.Sscan(m[1], &slot)
				if _, defined := n.defs[slot]; !defined {
					name(slot, m[2])
				}
			}
			for _, expr := range n.defs {
				for _, m := range slotRefRe.FindAllStringSubmatch(expr, -1) {
					var slot int
					fmt.Sscan(m[1], &slot)
					if _, defined := n.defs[slot]; !defined {
						name(slot, m[2])
					}
				}
			}
		}
		for slot, expr := range n.defs {
			defs[slot] = expr
		}
		for _, op := range []string{"UNION", "INTERSECT", "EXCEPT"} {
			if strings.HasPrefix(n.text, op+" ") {
				for i, slot := range n.outputs {
					for _, c := range n.children {
						if i < len(c.outputs) {
							setOps[slot] = append(setOps[slot], c.outputs[i])
						}
					}
				}
			}
		}
		for _, c := range n.children {
			if err := walk(c); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(out); err != nil {
		return nil, err
	}

	var resolve func(slot int, seen map[int]bool) slotSource
	resolve = func(slot int, seen map[int]bool) slotSource {
		if seen[slot] {
			return slotSource{columns: map[string]bool{}}
		}
		seen[slot] = true
		defer delete(seen, slot)
		if column, ok := scanned[slot]; ok {
			if column == "" {
				return slotSource{columns: map[string]bool{}}
			}
			return slotSource{columns: map[string]bool{column: true}, plain: true, complete: true}
		}
		if expr, ok := defs[slot]; ok {
			src := slotSource{columns: map[string]bool{}, complete: true}
			refs := slotRefRe.FindAllStringSubmatch(expr, -1)
			for _, m := range refs {
				var ref int
				fmt.Sscan(m[1], &ref)
				r := resolve(ref, seen)
				for c := range r.columns {
					src.columns[c] = true
				}
				src.complete = src.complete && r.complete
			}
			if slotOnlyRe.MatchString(expr) && len(refs) == 1 {
				var ref int
				fmt.Sscan(refs[0][1], &ref)
				src.plain = resolve(ref, seen).plain
			}
			return src
		}
		if arms, ok := setOps[slot]; ok {
			src := slotSource{columns: map[string]bool{}, complete: true}
			for _, arm := range arms {
				r := resolve(arm, seen)
				for c := range r.columns {
					src.columns[c] = true
				}
				src.complete = src.complete && r.complete
			}
			return src
		}
		// Nothing in the plan defines it: an operator the optimizer
		// folded away.
		return slotSource{columns: map[string]bool{}}
	}

	var lineage [][]string
	for _, slot := range out.outputs {
		r := resolve(slot, map[int]bool{})
		var sources []string
		for c := range r.columns {
			sources = append(sources, c)
		}
		if !r.plain || !r.complete {
			sources = append(sources, span.Unstated)
		}
		lineage = append(lineage, sources)
	}
	return lineage, nil
}
