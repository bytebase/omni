package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// BenchmarkExecDMLScaling measures Exec over function-free scripts of
// increasing statement counts. Time per run should scale linearly with n.
func BenchmarkExecDMLScaling(b *testing.B) {
	for _, n := range []int{1000, 10000} {
		var sb strings.Builder
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "UPDATE t SET c = %d WHERE id = %d;\n", i, i)
		}
		sql := sb.String()
		b.Run(fmt.Sprintf("dml_%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := New()
				if _, err := c.Exec(sql, nil); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
