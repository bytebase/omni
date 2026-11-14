// Command scriptdomg4 converts the T-SQL grammar of Microsoft's ScriptDOM,
// which is ANTLR 2 (mssql/testdata/TSql160.g with the keywords of
// TSqlTokenTypes.g, as the pinned commit publishes them), into the ANTLR 4
// grammar the compatibility rows are generated from,
// mssql/testdata/TSql160.g4. See convert for what the conversion keeps and
// drops.
//
// After a pin bump replaces the vendor files, from the repo root:
//
//	go run ./mssql/parser/cmd/scriptdomg4
//
// TestTSql160G4IsTheConversion fails while TSql160.g4 is anything else.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// testdata is where the grammars are, from the repo root.
const testdata = "mssql/testdata"

// The vendor files, as ScriptDOM publishes them under
// SqlScriptDom/Parser/TSql/.
const (
	grammarFile    = "TSql160.g"
	tokenTypesFile = "TSqlTokenTypes.g"
	supporterFile  = "CodeGenerationSupporter.cs"
	output         = "TSql160.g4"
)

func main() {
	text, err := convertDir(testdata)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scriptdomg4:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(testdata, output), []byte(text), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "scriptdomg4:", err)
		os.Exit(1)
	}
}

// convertDir converts the vendor files in dir.
func convertDir(dir string) (string, error) {
	read := func(name string) (source, error) {
		b, err := os.ReadFile(filepath.Join(dir, name))
		return source{name, string(b)}, err
	}
	g, err := read(grammarFile)
	if err != nil {
		return "", err
	}
	tt, err := read(tokenTypesFile)
	if err != nil {
		return "", err
	}
	cs, err := read(supporterFile)
	if err != nil {
		return "", err
	}
	return convert(g, tt, cs)
}
