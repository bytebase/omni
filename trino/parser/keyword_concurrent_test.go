package parser

import (
	"sync"
	"testing"
)

// TestTokenNameConcurrent calls TokenName from several goroutines at once, as
// callers that parse in parallel do. Run it with -race: the name table must be
// complete before the first call, so every goroutine sees the names a later
// sequential call returns.
func TestTokenNameConcurrent(t *testing.T) {
	kinds := make([]TokenKind, 0, len(keywordMap))
	for _, k := range keywordMap {
		kinds = append(kinds, k)
	}
	const goroutines = 8
	got := make([][]string, goroutines)
	var wg sync.WaitGroup
	for g := range goroutines {
		wg.Go(func() {
			names := make([]string, len(kinds))
			for i, k := range kinds {
				names[i] = TokenName(k)
			}
			got[g] = names
		})
	}
	wg.Wait()

	if name := TokenName(kwSELECT); name != "SELECT" {
		t.Fatalf("TokenName(kwSELECT) = %q, want SELECT", name)
	}
	for g, names := range got {
		for i, k := range kinds {
			if want := TokenName(k); names[i] != want {
				t.Errorf("goroutine %d: TokenName(%d) = %q, want %q", g, k, names[i], want)
				break
			}
		}
	}
}
