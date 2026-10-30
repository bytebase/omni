package catalog

import (
	"fmt"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/deparse"
	"github.com/bytebase/omni/mariadb/parser"
)

// nodeCounter is a Visitor that counts the nodes Walk visits.
type nodeCounter struct{ n int }

func (c *nodeCounter) Visit(n nodes.Node) nodes.Visitor {
	if n == nil {
		return nil
	}
	c.n++
	return c
}

// TestNestedSourcesWalkedOnce: a table value constructor nested in the ORDER
// BY of another, which MariaDB 11.8 accepts, doubled the work of every walk at
// each level while both the constructor and the SelectStmt over it held the
// ORDER BY: at 20 levels ast.Walk visited 8.4 million nodes and NodeToString
// wrote 293 MB. Each consumer must now finish 1,000 levels within the
// deadline.
func TestNestedSourcesWalkedOnce(t *testing.T) {
	const depth = 1000
	for _, open := range []string{"VALUES (1) ORDER BY ("} {
		query := "SELECT (" + strings.Repeat(open, depth) + "SELECT 1" + strings.Repeat(")", depth+1)
		list, err := parser.Parse(query)
		if err != nil {
			t.Fatalf("%s: %v", open, err)
		}
		sel := list.Items[0].(*nodes.SelectStmt)
		c := New()
		if _, err := c.Exec("CREATE DATABASE db; USE db; CREATE TABLE t (a INT)", nil); err != nil {
			t.Fatal(err)
		}
		steps := []struct {
			name string
			run  func() error
		}{
			{"ast.Walk", func() error {
				var v nodeCounter
				nodes.Walk(&v, sel)
				if v.n > 10*depth {
					return fmt.Errorf("visited %d nodes", v.n)
				}
				return nil
			}},
			{"NodeToString", func() error {
				if n := len(nodes.NodeToString(sel)); n > 1000*depth {
					return fmt.Errorf("wrote %d bytes", n)
				}
				return nil
			}},
			{"deparse", func() error {
				deparse.DeparseSelect(sel)
				return nil
			}},
			{"AnalyzeSelectStmt", func() error {
				_, err := c.AnalyzeSelectStmt(sel)
				return err
			}},
			{"CREATE VIEW", func() error {
				results, err := c.Exec("CREATE VIEW v AS "+query, nil)
				if err != nil {
					return err
				}
				return results[0].Error
			}},
		}
		for _, step := range steps {
			done := make(chan error, 1)
			go func() { done <- step.run() }()
			select {
			case err := <-done:
				if err != nil {
					t.Errorf("%s %s: %v", open, step.name, err)
				}
			case <-time.After(10 * time.Second):
				t.Fatalf("%s %s: still running after 10s", open, step.name)
			}
		}
	}
}
