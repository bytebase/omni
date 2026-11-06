package review

import (
	"testing"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestSetSwitch reads the value of a SET directly. The walk only asks whether the checks
// are on, so its oracle cannot tell Off from Unknown, and this is what holds them apart.
func TestSetSwitch(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want review.Switch
	}{
		{"set a = 1", review.SwitchOn},
		{"set a = 2", review.SwitchOn},
		{"set a = 0", review.SwitchOff},
		{"set a = +0", review.SwitchOff},
		{"set a = +1", review.SwitchOn},
		{"set a = -1", review.SwitchUnknown},
		{"set a = true", review.SwitchOn},
		{"set a = false", review.SwitchOff},
		{"set a = on", review.SwitchOn},
		{"set a = ON", review.SwitchOn},
		{"set a = 'on'", review.SwitchOn},
		{"set a = off", review.SwitchOff},
		{"set a = OFF", review.SwitchOff},
		{"set a = 'off'", review.SwitchOff},
		{"set a = binary", review.SwitchUnknown},
		{"set a = default", review.SwitchUnknown},
		{"set a = 'x'", review.SwitchUnknown},
		{"set a = b", review.SwitchUnknown},
		{"set a = b.off", review.SwitchUnknown},
		{"set a = @off", review.SwitchUnknown},
		{"set a = 1 + 1", review.SwitchUnknown},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := parser.Parse(c.sql)
			if err != nil {
				t.Fatal(err)
			}
			value := list.Items[0].(*ast.SetStmt).Assignments[0].Value
			if got := setSwitch(value); got != c.want {
				t.Errorf("setSwitch = %v, want %v", got, c.want)
			}
		})
	}
}
