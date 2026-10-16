package parser

import (
	"sync"
	"testing"
)

// TestKeywordHelpersConcurrent calls TokenName and IsReservedKeyword from
// several goroutines at once, as callers that parse or review in parallel do.
// Run it with -race: the keyword tables must be complete before the first
// call, so every goroutine sees the answers a later sequential call returns.
func TestKeywordHelpersConcurrent(t *testing.T) {
	words := make([]string, 0, len(oracleKeywords))
	for word := range oracleKeywords {
		words = append(words, word)
	}
	type answer struct {
		name     string
		reserved bool
	}
	const goroutines = 8
	got := make([][]answer, goroutines)
	var wg sync.WaitGroup
	for g := range goroutines {
		wg.Go(func() {
			answers := make([]answer, len(words))
			for i, word := range words {
				answers[i] = answer{TokenName(oracleKeywords[word]), IsReservedKeyword(word)}
			}
			got[g] = answers
		})
	}
	wg.Wait()

	if name := TokenName(SELECT); name != "SELECT" {
		t.Fatalf("TokenName(SELECT) = %q, want SELECT", name)
	}
	if !IsReservedKeyword("select") || IsReservedKeyword("employees") {
		t.Fatalf("IsReservedKeyword: select = %v, employees = %v; want true, false",
			IsReservedKeyword("select"), IsReservedKeyword("employees"))
	}
	for g, answers := range got {
		for i, word := range words {
			want := answer{TokenName(oracleKeywords[word]), IsReservedKeyword(word)}
			if answers[i] != want {
				t.Errorf("goroutine %d: %s = %+v, want %+v", g, word, answers[i], want)
				break
			}
		}
	}
}
