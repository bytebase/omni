// Scratch: not committed.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bytebase/omni/enginetest/grammar"
)

func main() {
	tokens, err := grammar.ReadTokens(os.Args[1])
	if err != nil {
		panic(err)
	}
	var src strings.Builder
	for _, p := range os.Args[2:] {
		b, err := os.ReadFile(p)
		if err != nil {
			panic(err)
		}
		src.Write(b)
		src.WriteString("\n")
	}
	g, err := grammar.ReadANTLR(src.String(), tokens)
	if err != nil {
		fmt.Println("ERR", err)
		os.Exit(1)
	}
	rows, uncovered := g.Rows()
	rec := g.Recognizer()
	alts, nm, altAcc, nmAcc := 0, 0, 0, 0
	for _, r := range rows {
		ok := rec.Accepts(r.Statement)
		if grammar.IsNearMiss(r.Rule) {
			nm++
			if ok {
				nmAcc++
			}
		} else {
			alts++
			if ok {
				altAcc++
			} else {
				fmt.Println("NOT DERIVED", r.Rule, "\t", r.Statement)
			}
		}
		if os.Getenv("DUMP") != "" {
			fmt.Printf("%v\t%s\t%s\n", ok, r.Rule, r.Statement)
		}
	}
	fmt.Printf("alternatives %d (derived %d), near-misses %d (derived %d), uncovered %d: %v\n", alts, altAcc, nm, nmAcc, len(uncovered), uncovered)
}
