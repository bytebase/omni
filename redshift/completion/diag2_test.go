package completion

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/redshift/parser"
)

func TestDiagUnion(t *testing.T) {
	// (SELECT c1 FROM t1) UNION (SELECT | FROM t2)
	// cursor should be right after "SELECT " inside the second parens
	sql := "(SELECT c1 FROM t1) UNION (SELECT  FROM t2)"
	//       0123456789...
	// Let me find the offset: "(SELECT c1 FROM t1) UNION (SELECT " = 34
	offset := 34

	fmt.Printf("sql[%d:]=%q\n", offset, sql[offset:])

	cs := parser.Collect(sql, offset)
	if cs == nil {
		fmt.Println("Collect returned nil")
		return
	}
	fmt.Printf("Rules: %v\n", cs.Rules)
	fmt.Printf("Tokens count: %d\n", len(cs.Tokens))
	for _, tok := range cs.Tokens {
		name := parser.TokenName(tok)
		if name != "" {
			fmt.Printf("  token: %s (%d)\n", name, tok)
		}
	}

	// Also try as a statement (not wrapped in parens)
	sql2 := "SELECT c1 FROM t1 UNION SELECT  FROM t2"
	offset2 := 31
	fmt.Printf("\nsql2[%d:]=%q\n", offset2, sql2[offset2:])
	cs2 := parser.Collect(sql2, offset2)
	if cs2 == nil {
		fmt.Println("Collect2 returned nil")
		return
	}
	fmt.Printf("Rules2: %v\n", cs2.Rules)
	fmt.Printf("HasColumnref2: %v\n", cs2.HasRule("columnref"))
}
