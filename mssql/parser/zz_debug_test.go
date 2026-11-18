package parser

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestZZDebug parses statements from $ZZ_STMTS (one per line) and prints the verdict.
func TestZZDebug(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("ZZ_STMTS"))
	if err != nil {
		t.Skip(err)
	}
	for _, sql := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
		if strings.TrimSpace(sql) == "" {
			continue
		}
		_, perr := Parse(sql)
		if perr != nil {
			fmt.Printf("REJECT %-60s | %v\n", sql, perr)
		} else {
			fmt.Printf("accept %-60s\n", sql)
		}
	}
}
