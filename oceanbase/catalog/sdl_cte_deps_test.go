package catalog

import "testing"

// TestSDLViewReadThroughSameNamedCTE holds LoadSDL's order tolerance where a view's
// CTE is named like an object the same CTE reads. Without RECURSIVE a CTE's body
// stands outside its own name, so WITH a AS (SELECT b FROM a) reads the view a and
// must be loaded after it, whichever order the statements are written in.
func TestSDLViewReadThroughSameNamedCTE(t *testing.T) {
	for _, sdl := range []string{
		"CREATE DATABASE d; USE d; CREATE VIEW v AS WITH a AS (SELECT b FROM a) SELECT b FROM a; CREATE VIEW a AS SELECT 1 AS b",
		"CREATE DATABASE d; USE d; CREATE VIEW a AS SELECT 1 AS b; CREATE VIEW v AS WITH a AS (SELECT b FROM a) SELECT b FROM a",
		// A recursive CTE reads itself, so it depends on no object of that name.
		"CREATE DATABASE d; USE d; CREATE VIEW v AS WITH RECURSIVE r AS (SELECT 1 AS x UNION ALL SELECT x+1 FROM r WHERE x < 3) SELECT x FROM r",
		// A later CTE reads an earlier one, which is no object either.
		"CREATE DATABASE d; USE d; CREATE VIEW v AS WITH c AS (SELECT 1 AS x), e AS (SELECT x FROM c) SELECT x FROM e",
	} {
		if _, err := LoadSDL(sdl); err != nil {
			t.Errorf("%s\n  %v", sdl, err)
		}
	}
}
