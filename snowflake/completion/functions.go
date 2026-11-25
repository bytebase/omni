package completion

import (
	_ "embed"
	"strings"
	"sync"
)

// functionsTSV is Snowflake's function index as
// snowflake/completion/cmd/docsfunctions recorded it.
//
//go:embed functions.tsv
var functionsTSV string

type builtin struct {
	name, comment string
}

// builtinFunctions returns the functions of functionsTSV.
var builtinFunctions = sync.OnceValue(func() []builtin {
	var out []builtin
	for _, line := range strings.Split(functionsTSV, "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, comment, _ := strings.Cut(line, "\t")
		out = append(out, builtin{name: name, comment: comment})
	}
	return out
})
