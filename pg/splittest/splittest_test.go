package splittest

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// Scale knobs. PR runs use the defaults; nightly overrides via env:
//
//	SPLITTEST_N=1000000 SPLITTEST_SEED=<date-derived> go test ./pg/splittest/
const (
	defaultN    = 10000
	defaultSeed = 20260716
)

func scale() (n int, seed int64) {
	n, seed = defaultN, defaultSeed
	if v := os.Getenv("SPLITTEST_N"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			n = p
		}
	}
	if v := os.Getenv("SPLITTEST_SEED"); v != "" {
		if p, err := strconv.ParseInt(v, 10, 64); err == nil {
			seed = p
		}
	}
	return
}

// TestConstructive is the generator gate: every case carries its own
// expected segmentation (truth_source=construct). Failures dump the
// full script for reproduction; the seed is printed so any run can be
// replayed exactly.
func TestConstructive(t *testing.T) {
	n, seed := scale()
	r := rand.New(rand.NewSource(seed))
	atoms := EnabledAtoms( /* enable "D2","D3","D5" as fixes land */ )
	t.Logf("constructive: n=%d seed=%d atoms=%d", n, seed, len(atoms))
	fails := 0
	for i := 0; i < n; i++ {
		script := Compose(r, atoms, 1+r.Intn(5))
		if err := CheckScript(script); err != nil {
			fails++
			t.Errorf("case %d (seed %d): %v", i, seed, err)
			if fails >= 10 {
				t.Fatalf("too many failures, aborting at case %d", i)
			}
		}
	}
}

// TestPairCoverage guarantees every ordered pair of enabled atom
// classes appears within one statement, several times, per the
// combination-layer quota (design §2.2).
func TestPairCoverage(t *testing.T) {
	_, seed := scale()
	r := rand.New(rand.NewSource(seed + 1))
	atoms := EnabledAtoms()
	const perPair = 5
	for ai := range atoms {
		for bi := range atoms {
			for k := 0; k < perPair; k++ {
				a, b := atoms[ai], atoms[bi]
				stmt := glue(r) + " " + a.Gen(r) + " " + glue(r) + " " + b.Gen(r) + " " + glue(r)
				script := Script{SQL: stmt + ";", Want: []string{stmt + ";"}}
				if err := CheckScript(script); err != nil {
					t.Errorf("pair %s×%s: %v", a.Class, b.Class, err)
				}
			}
		}
	}
}

// TestInvariantsOverPgregress is the S1 hard gate: the truth-free
// invariants must hold over the full vendored PG regression corpus
// (truth_source=invariant). Any violation is a P0 by definition —
// lossless reconstruction failed on real-world input.
func TestInvariantsOverPgregress(t *testing.T) {
	root := filepath.Join("..", "pgregress", "testdata", "sql")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Skipf("pgregress corpus not found: %v", err)
	}
	files := 0
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		content, err := os.ReadFile(filepath.Join(root, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := CheckInvariants(string(content)); err != nil {
			t.Errorf("%s: %v", e.Name(), err)
		}
		files++
	}
	t.Logf("invariants over %d pgregress files", files)
	if files < 200 {
		t.Errorf("expected >=200 corpus files, found %d — corpus path drift?", files)
	}
}

// TestByteMutation fuzzes with quote/semicolon/backslash/dollar
// insertions into generated scripts. Mutated inputs are usually
// invalid SQL — only the truth-free invariants apply, which is exactly
// their job (truth_source=invariant).
func TestByteMutation(t *testing.T) {
	n, seed := scale()
	r := rand.New(rand.NewSource(seed + 2))
	atoms := EnabledAtoms()
	hot := []byte{'\'', '"', ';', '\\', '$', '-', '*', '/', 'E'}
	for i := 0; i < n/2; i++ {
		script := Compose(r, atoms, 1+r.Intn(3))
		b := []byte(script.SQL)
		for m := 0; m < 1+r.Intn(4); m++ {
			switch r.Intn(3) {
			case 0: // insert
				pos := r.Intn(len(b) + 1)
				b = append(b[:pos], append([]byte{hot[r.Intn(len(hot))]}, b[pos:]...)...)
			case 1: // delete
				if len(b) > 0 {
					pos := r.Intn(len(b))
					b = append(b[:pos], b[pos+1:]...)
				}
			default: // replace
				if len(b) > 0 {
					b[r.Intn(len(b))] = hot[r.Intn(len(hot))]
				}
			}
		}
		if err := CheckInvariants(string(b)); err != nil {
			t.Errorf("mutation case %d (seed %d): %v\ninput: %q", i, seed, err, string(b))
		}
	}
}

// FuzzSplitInvariants hooks the invariants into Go native fuzzing so
// `go test -fuzz` can explore beyond the structured generators. Seeds
// cover every audit defect class.
func FuzzSplitInvariants(f *testing.F) {
	seeds := []string{
		`SELECT E'a\';b'; SELECT 2;`,
		"CREATE RULE r AS ON UPDATE TO t DO ALSO (SELECT 1; SELECT 2);",
		"SELECT 1 AS abc$x$y; SELECT 2;",
		"CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT t4.end FROM t4; END; SELECT 1;",
		"COPY t FROM stdin;\na;b\n\\.\nSELECT 1;",
		"/* /* nested; */ ; */ SELECT $tag$;$tag$;",
		"SELECT U&'d\\0061t;a'; SELECT 'x''y;z';",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, sql string) {
		if err := CheckInvariants(sql); err != nil {
			t.Fatalf("%v", err)
		}
	})
}
