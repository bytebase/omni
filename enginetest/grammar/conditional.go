package grammar

import (
	"fmt"
	"slices"
	"strings"
)

// Conditional resolves the %ifdef lines of a grammar file that holds more
// than one grammar, and returns the one built with the names in defined.
// They are not bison's: MariaDB keeps its default and its Oracle-mode
// grammar in one sql_yacc.yy, and its build writes each out before bison
// runs (sql/gen_yy_files.cmake):
//
//	%ifdef MARIADB
//	%expect 62
//	%else
//	%expect 63
//	%endif
//
// The lines from "%ifdef NAME" to "%else" or "%endif" are kept when NAME is
// in defined, and the lines from "%else" to "%endif" when it is not. Blocks
// do not nest. A dropped line and the directives themselves are left as
// blank lines, so every other line keeps its number. A file with no
// directive comes back as it is.
func Conditional(src string, defined ...string) (string, error) {
	const (
		outside = iota
		kept
		dropped
	)
	lines := strings.Split(src, "\n")
	state, opened := outside, 0
	for i, line := range lines {
		word, rest, _ := strings.Cut(strings.TrimRight(line, " \t\r"), " ")
		switch word {
		case "%ifdef":
			if state != outside {
				return "", fmt.Errorf("line %d: %%ifdef inside the %%ifdef of line %d", i+1, opened)
			}
			name := strings.TrimSpace(rest)
			if name == "" {
				return "", fmt.Errorf("line %d: %%ifdef names no condition", i+1)
			}
			state, opened = dropped, i+1
			if slices.Contains(defined, name) {
				state = kept
			}
		case "%else":
			if state == outside {
				return "", fmt.Errorf("line %d: %%else outside an %%ifdef", i+1)
			}
			state = kept + dropped - state
		case "%endif":
			// MariaDB writes the condition again after some: %endif ORACLE.
			if state == outside {
				return "", fmt.Errorf("line %d: %%endif outside an %%ifdef", i+1)
			}
			state = outside
		default:
			if state != dropped {
				continue
			}
		}
		lines[i] = ""
	}
	if state != outside {
		return "", fmt.Errorf("line %d: %%ifdef is never closed", opened)
	}
	return strings.Join(lines, "\n"), nil
}
