package walktest

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/review/walk"
)

// fakeEngine runs statements of the form "name" (accepted) or
// "name!code" (rejected with code), a rejected one ending in "+" changing
// the database. Its walk reports the statement named in report, with the
// message message.
type fakeEngine struct {
	digest  int
	ran     []string
	report  string
	message string
}

func (e *fakeEngine) Statements(change string) []string { return strings.Split(change, ";") }

func (e *fakeEngine) Exec(_ testing.TB, stmt string) Outcome {
	e.ran = append(e.ran, stmt)
	name, code, ok := strings.Cut(stmt, "!")
	if !ok {
		e.digest++
		return Outcome{}
	}
	if strings.HasSuffix(code, "+") {
		e.digest++
		code = strings.TrimSuffix(code, "+")
	}
	return Outcome{Code: code, Message: name}
}

func (e *fakeEngine) Report(_ testing.TB, change string) *Finding {
	stmts := strings.Split(change, "\n")
	for i, s := range stmts {
		if s == e.report {
			return &Finding{Statement: i, Message: e.message}
		}
	}
	return nil
}

func (e *fakeEngine) Violation(testing.TB, string) *walk.Violation { return &walk.Violation{} }

func (e *fakeEngine) Codes(message string) []string { return []string{message} }

func (e *fakeEngine) Holds(testing.TB, *walk.Violation, string) (bool, string) {
	return false, "the fake holds no claim"
}

func (e *fakeEngine) Digest(testing.TB) string { return fmt.Sprint(e.digest) }

type partialEngine struct{ fakeEngine }

func (*partialEngine) RejectsPartially() {}

// recorder is a testing.TB that records failures; Fatalf ends the
// goroutine, as testing.T's does.
type recorder struct {
	testing.TB
	errors []string
	fatal  string
}

func (r *recorder) Helper() {}

func (r *recorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

func (r *recorder) Fatalf(format string, args ...any) {
	r.fatal = fmt.Sprintf(format, args...)
	runtime.Goexit()
}

func judge(t *testing.T, e Engine, change string) (r *recorder, found, rejected int) {
	r = &recorder{TB: t}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		found, rejected = Judge(r, e, change, true)
	}()
	wg.Wait()
	return r, found, rejected
}

func TestJudgeStopsAfterAPartialRejection(t *testing.T) {
	e := &partialEngine{fakeEngine{report: "b!1051+", message: "1051"}}
	r, found, rejected := judge(t, e, "a;b!1051+;c!1146")
	if r.fatal != "" || len(r.errors) != 0 || found != 1 || rejected != 1 {
		t.Errorf("fatal %q, errors %q, found %d of %d; want 1 of 1 and no failure", r.fatal, r.errors, found, rejected)
	}
	if got := strings.Join(e.ran, ";"); got != "a;b!1051+" {
		t.Errorf("ran %q, want the change to stop after the partial rejection", got)
	}
}

func TestJudgeHoldsAPartialRejectionToItsCode(t *testing.T) {
	e := &partialEngine{fakeEngine{report: "b!1051+", message: "1146"}}
	r, _, _ := judge(t, e, "a;b!1051+")
	if len(r.errors) != 1 || !strings.Contains(r.errors[0], "cannot be checked") {
		t.Errorf("errors %q, want one that the claim cannot be checked", r.errors)
	}
}

func TestJudgeFailsAPartialRejectionWithoutOptIn(t *testing.T) {
	e := &fakeEngine{report: "b!1051+", message: "1051"}
	r, _, _ := judge(t, e, "a;b!1051+;c")
	if !strings.Contains(r.fatal, "changed the database") {
		t.Errorf("fatal %q, want the change of the database reported", r.fatal)
	}
}

func TestJudgeGoesOnAfterAnAtomicRejection(t *testing.T) {
	e := &partialEngine{fakeEngine{report: "b!1051", message: "1051"}}
	r, found, rejected := judge(t, e, "a;b!1051;c!1146")
	// c's rejection is missed, which complete reports.
	if r.fatal != "" || len(r.errors) != 1 || found != 1 || rejected != 2 {
		t.Errorf("fatal %q, errors %q, found %d of %d; want 1 of 2 and the miss of c", r.fatal, r.errors, found, rejected)
	}
}

// listEngine is a Lister over the fake's statements. "+name" creates the
// object name and "-name" drops it; "name!code" is rejected with code, and
// "+name!code+" creates name while rejecting. Its walk reports the
// statement named in report, claiming claimed is fact, of kind where set.
type listEngine struct {
	held    map[string]string // name -> kind
	asked   []string          // every request to the server, in order
	report  string
	message string
	claimed string
	fact    walk.Fact
	kind    string
}

func (e *listEngine) Statements(change string) []string { return strings.Split(change, ";") }

func (e *listEngine) Exec(_ testing.TB, stmt string) Outcome {
	e.asked = append(e.asked, "exec "+stmt)
	body, code, rejects := strings.Cut(stmt, "!")
	if !rejects || strings.HasSuffix(code, "+") {
		if name, ok := strings.CutPrefix(body, "+"); ok {
			e.held[name] = "table"
		} else if name, ok := strings.CutPrefix(body, "-"); ok {
			delete(e.held, name)
		}
	}
	if !rejects {
		return Outcome{}
	}
	return Outcome{Code: strings.TrimSuffix(code, "+"), Message: body}
}

func (e *listEngine) Report(_ testing.TB, change string) *Finding {
	for i, s := range strings.Split(change, "\n") {
		if s == e.report {
			return &Finding{Statement: i, Message: e.message}
		}
	}
	return nil
}

func (e *listEngine) Violation(testing.TB, string) *walk.Violation {
	return &walk.Violation{Check: walk.Check{Key: walk.Key{Space: "relation", Name: e.claimed}}, Got: e.fact, Kind: e.kind}
}

func (e *listEngine) Codes(message string) []string { return []string{message} }

func (e *listEngine) Inventory(_ testing.TB, change string) []Object {
	e.asked = append(e.asked, "inventory "+change)
	var out []Object
	for name, kind := range e.held {
		out = append(out, Object{Key: walk.Key{Space: "relation", Name: name}, Kind: kind})
	}
	return out
}

// Same folds case, as an engine's name folding would.
func (e *listEngine) Same(claimed, stored walk.Key) bool {
	return claimed.Space == stored.Space && claimed.Name == strings.ToLower(stored.Name)
}

type partialListEngine struct{ listEngine }

func (*partialListEngine) RejectsPartially() {}

func TestJudgeReadsAClaimOffTheList(t *testing.T) {
	for _, c := range []struct {
		name    string
		held    map[string]string
		claimed string
		fact    walk.Fact
		kind    string
		wrong   string // a part of the failure, "" for none
	}{
		{name: "present, and listed", held: map[string]string{"T": "table"}, claimed: "t", fact: walk.Present},
		{name: "present of its kind", held: map[string]string{"T": "table"}, claimed: "t", fact: walk.Present, kind: "table"},
		{name: "present of a kind the list gives none for", held: map[string]string{"T": ""}, claimed: "t", fact: walk.Present, kind: "table"},
		{name: "absent, and not listed", held: map[string]string{"u": "table"}, claimed: "t", fact: walk.Absent},
		{name: "unknown claims nothing", held: map[string]string{"t": "table"}, claimed: "t", fact: walk.Unknown},
		{name: "present, but not listed", held: map[string]string{"u": "table"}, claimed: "t", fact: walk.Present, wrong: "is present: false"},
		{name: "absent, but listed", held: map[string]string{"T": "table"}, claimed: "t", fact: walk.Absent, wrong: "is present: true"},
		{name: "present of another kind", held: map[string]string{"t": "view"}, claimed: "t", fact: walk.Present, kind: "table", wrong: "is a view, not a table"},
	} {
		t.Run(c.name, func(t *testing.T) {
			e := &listEngine{held: c.held, report: "x!1", message: "another code", claimed: c.claimed, fact: c.fact, kind: c.kind}
			r, found, rejected := judge(t, e, "x!1")
			if r.fatal != "" || found != 1 || rejected != 1 {
				t.Fatalf("fatal %q, found %d of %d; want 1 of 1", r.fatal, found, rejected)
			}
			switch {
			case c.wrong == "" && len(r.errors) != 0:
				t.Errorf("errors %q, want none", r.errors)
			case c.wrong != "" && (len(r.errors) != 1 || !strings.Contains(r.errors[0], c.wrong)):
				t.Errorf("errors %q, want one with %q", r.errors, c.wrong)
			}
		})
	}
}

// TestJudgeAsksAListerTheSameWhateverTheWalkReports holds Judge to what a
// tape needs: the requests a Lister's server gets follow from the change
// and the server's answers alone. The list is read before every statement
// but one after a rejection, whose list is the one read after it.
func TestJudgeAsksAListerTheSameWhateverTheWalkReports(t *testing.T) {
	const change = "+a;x!1;y!2;+b;z!3"
	var want []string
	for _, e := range []*listEngine{
		{},
		{report: "x!1", message: "1"},
		{report: "y!2", message: "no code", claimed: "a", fact: walk.Present},
		{report: "z!3", message: "no code", claimed: "a", fact: walk.Absent},
		{report: "+b", message: "a statement the server accepts"},
	} {
		e.held = map[string]string{}
		judge(t, e, change)
		if want == nil {
			want = e.asked
			if got, list := strings.Join(want, "; "), "inventory "+change; got != strings.Join([]string{
				list, "exec +a", list, "exec x!1", list, "exec y!2", list, "exec +b", list, "exec z!3", list,
			}, "; ") {
				t.Fatalf("asked %q", got)
			}
			continue
		}
		if !slices.Equal(e.asked, want) {
			t.Errorf("a walk reporting %q asked %q, want %q", e.report, e.asked, want)
		}
	}
}

func TestJudgeFailsAListerWhoseRejectedStatementChangesTheList(t *testing.T) {
	e := &listEngine{held: map[string]string{}, report: "+a!1+", message: "1"}
	r, _, _ := judge(t, e, "+a!1+;b")
	if !strings.Contains(r.fatal, "changed the database") || !strings.Contains(r.fatal, "+ relation a (table)") {
		t.Errorf("fatal %q, want the change of the database reported with the object it added", r.fatal)
	}

	p := &partialListEngine{listEngine{held: map[string]string{}, report: "+a!1+", message: "1"}}
	r, found, rejected := judge(t, p, "+a!1+;b!2")
	if r.fatal != "" || len(r.errors) != 0 || found != 1 || rejected != 1 {
		t.Errorf("fatal %q, errors %q, found %d of %d; want a partial rejecter to judge it by its code and stop", r.fatal, r.errors, found, rejected)
	}
}

func TestJudgeNeedsAListerOrAHolder(t *testing.T) {
	r, _, _ := judge(t, bareEngine{}, "a")
	if !strings.Contains(r.fatal, "neither lists") {
		t.Errorf("fatal %q, want an engine that shows no catalog refused", r.fatal)
	}
}

type bareEngine struct{}

func (bareEngine) Statements(change string) []string            { return []string{change} }
func (bareEngine) Exec(testing.TB, string) Outcome              { return Outcome{} }
func (bareEngine) Report(testing.TB, string) *Finding           { return nil }
func (bareEngine) Violation(testing.TB, string) *walk.Violation { return nil }
func (bareEngine) Codes(string) []string                        { return nil }

func TestShards(t *testing.T) {
	// A sweep of the default size keeps to one server; the sizes CI runs
	// spread out, up to the cap.
	for _, c := range []struct{ n, want int }{
		{0, 1}, {40, 1}, {60, 1}, {199, 1}, {200, 2}, {399, 3}, {400, 4}, {500, 4}, {5000, 4},
	} {
		if got := Shards(c.n); got != c.want {
			t.Errorf("Shards(%d) = %d, want %d", c.n, got, c.want)
		}
	}
}

// TestSweepChangesComeFromOneSequence holds SweepChanges to the reason it
// exists: the changes are generated before any of them runs, from one
// generator over one rng, so a sweep covers the same changes however many
// shards it deals them over. A generator moved back into a shard's loop
// would give every shard the sequence from its start, and fail here.
func TestSweepChangesComeFromOneSequence(t *testing.T) {
	const n = 37
	generator := func() (gen func(*rand.Rand) (string, string), calls *int, rngs *[]*rand.Rand) {
		calls, rngs = new(int), new([]*rand.Rand)
		return func(rng *rand.Rand) (string, string) {
			*calls++
			*rngs = append(*rngs, rng)
			return fmt.Sprintf("base%d", rng.Int()), fmt.Sprintf("change%d", rng.Int())
		}, calls, rngs
	}

	gen, calls, rngs := generator()
	changes := SweepChanges(n, 20261104, gen)
	if len(changes) != n {
		t.Fatalf("got %d changes, want %d", len(changes), n)
	}
	if *calls != n {
		t.Errorf("the generator ran %d times, want %d", *calls, n)
	}
	for _, rng := range *rngs {
		if rng != (*rngs)[0] {
			t.Fatal("the generator ran on more than one rng, so which changes a sweep covers depends on how they are dealt")
		}
	}

	again, _, _ := generator()
	for i, c := range SweepChanges(n, 20261104, again) {
		if c != changes[i] {
			t.Fatalf("change %d of a second sweep of the same seed is %q, want %q", i, c, changes[i])
		}
	}
}

// TestShardsDealEveryChangeOnce holds the arithmetic every sweep deals
// its changes by: shard s takes s, s+shards, ..., and the shards together
// take every change exactly once.
func TestShardsDealEveryChangeOnce(t *testing.T) {
	for n := range 40 {
		for shards := 1; shards <= 5; shards++ {
			dealt := make([]int, n)
			for shard := range shards {
				for i := shard; i < n; i += shards {
					dealt[i]++
				}
			}
			for i, times := range dealt {
				if times != 1 {
					t.Fatalf("%d changes over %d shards: change %d is dealt %d times", n, shards, i, times)
				}
			}
		}
	}
}

func writeSimFiles(t *testing.T, tsv string, fixtures map[string]string) (tsvPath, fxDir string) {
	t.Helper()
	dir := t.TempDir()
	fxDir = filepath.Join(dir, "fixtures")
	if err := os.MkdirAll(fxDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, sql := range fixtures {
		if err := os.WriteFile(filepath.Join(fxDir, name+".sql"), []byte(sql), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	tsvPath = filepath.Join(dir, "simulate.tsv")
	header := "# pin: default test:1 sha256:0000 1\n# columns: omni, expected, source, fixture, rule, statement\n"
	if err := os.WriteFile(tsvPath, []byte(header+tsv), 0o644); err != nil {
		t.Fatal(err)
	}
	return tsvPath, fxDir
}

func runCheckSimulation(t *testing.T, s Simulation) (*recorder, SimulationResult) {
	r := &recorder{TB: t}
	var res SimulationResult
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		res = CheckSimulation(r, s)
	}()
	wg.Wait()
	return r, res
}

func TestCheckSimulationSoundness(t *testing.T) {
	fixtures := map[string]string{
		"fx_a": "CREATE TABLE a (id INT);",
		"fx_b": "CREATE TABLE t (id INT);",
	}
	same := func(claimed, stored walk.Key) bool {
		return claimed.Space == stored.Space &&
			claimed.Parent == stored.Parent &&
			strings.EqualFold(claimed.Name, stored.Name)
	}

	t.Run("passing mix of ok, matched code, true fact, and silence", func(t *testing.T) {
		tsv := strings.Join([]string{
			"ok:0123456789ab\tok:0123456789ab\tengine\tfx_a\tr/ok\tSELECT 1",
			"1050\t1050\tengine\tfx_a\tr/code\tCREATE TABLE a (id INT)",
			"1060\t1060\tengine\tfx_a\tr/fact\tALTER TABLE a ADD COLUMN x INT",
			"1146\t1146\tengine\tfx_b\tr/silent\tDROP TABLE missing",
			"42P07\t42P07\tengine\tfx_a\tr/pgstate\tCREATE TABLE a (id INT)",
		}, "\n") + "\n"
		tsvPath, fxDir := writeSimFiles(t, tsv, fixtures)

		setupCalls := map[string]int{}
		r, res := runCheckSimulation(t, Simulation{
			Path:     tsvPath,
			Fixtures: fxDir,
			Setup: func(_ testing.TB, fx simulate.Fixture) SimFixture {
				setupCalls[fx.Name]++
				var objs []Object
				if fx.Name == "fx_a" {
					objs = []Object{{Key: walk.Key{Space: "relation", Name: "A"}, Kind: "table"}}
				} else {
					objs = []Object{{Key: walk.Key{Space: "relation", Name: "t"}, Kind: "table"}}
				}
				return SimFixture{
					Objects: objs,
					Report: func(_ testing.TB, stmt string) *Finding {
						switch stmt {
						case "CREATE TABLE a (id INT)":
							return &Finding{Statement: 0, Message: "exists"}
						case "ALTER TABLE a ADD COLUMN x INT":
							return &Finding{Statement: 0, Message: "other-reason"}
						default:
							return nil
						}
					},
					Violation: func(_ testing.TB, stmt string) *walk.Violation {
						if stmt == "ALTER TABLE a ADD COLUMN x INT" {
							return &walk.Violation{
								Check: walk.Check{Key: walk.Key{Space: "relation", Name: "a"}},
								Got:   walk.Present,
								Kind:  "table",
							}
						}
						return nil
					},
				}
			},
			Codes: func(msg, _ string) []string {
				if msg == "exists" {
					return []string{"1050", "42P07"}
				}
				return nil
			},
			Same: same,
		})
		if r.fatal != "" || len(r.errors) != 0 {
			t.Fatalf("unexpected failure: fatal=%q errors=%q", r.fatal, r.errors)
		}
		want := SimulationResult{Total: 5, OK: 1, Refused: 4, Found: 3}
		if res != want {
			t.Errorf("result = %+v, want %+v", res, want)
		}
		if setupCalls["fx_a"] != 1 || setupCalls["fx_b"] != 1 {
			t.Errorf("setup calls = %v, want each fixture loaded once", setupCalls)
		}
	})

	t.Run("any finding on an ok row fails", func(t *testing.T) {
		tsv := "ok:0123456789ab\tok:0123456789ab\tengine\tfx_a\tr/ok\tALTER TABLE a ADD COLUMN y INT\n"
		tsvPath, fxDir := writeSimFiles(t, tsv, fixtures)

		r, res := runCheckSimulation(t, Simulation{
			Path:     tsvPath,
			Fixtures: fxDir,
			Setup: func(_ testing.TB, _ simulate.Fixture) SimFixture {
				return SimFixture{
					Report: func(_ testing.TB, _ string) *Finding {
						return &Finding{Statement: 0, Message: "bogus finding"}
					},
					Violation: func(_ testing.TB, _ string) *walk.Violation { return nil },
				}
			},
			Same: same,
		})
		if r.fatal != "" || len(r.errors) != 1 || !strings.Contains(r.errors[0], "engine accepts") {
			t.Errorf("fatal=%q errors=%q; want one error about engine accepting the row", r.fatal, r.errors)
		}
		if res.OK != 1 || res.Refused != 0 {
			t.Errorf("result = %+v, want OK=1 Refused=0", res)
		}
	})

	t.Run("wrong code with false fact on refused row fails", func(t *testing.T) {
		tsv := "1146\t1146\tengine\tfx_b\tr/bad\tALTER TABLE a ADD COLUMN y INT\n"
		tsvPath, fxDir := writeSimFiles(t, tsv, fixtures)

		r, _ := runCheckSimulation(t, Simulation{
			Path:     tsvPath,
			Fixtures: fxDir,
			Setup: func(_ testing.TB, _ simulate.Fixture) SimFixture {
				return SimFixture{
					Objects: []Object{{Key: walk.Key{Space: "relation", Name: "t"}, Kind: "table"}},
					Report: func(_ testing.TB, _ string) *Finding {
						return &Finding{Statement: 0, Message: "wrong code"}
					},
					Violation: func(_ testing.TB, _ string) *walk.Violation {
						// Claims table a is present, but fx_b only has t.
						return &walk.Violation{
							Check: walk.Check{Key: walk.Key{Space: "relation", Name: "a"}},
							Got:   walk.Present,
							Kind:  "table",
						}
					},
				}
			},
			Codes: func(string, string) []string { return []string{"1050"} },
			Same:  same,
		})
		if r.fatal != "" || len(r.errors) != 1 || !strings.Contains(r.errors[0], "is present: false") {
			t.Errorf("fatal=%q errors=%q; want one error saying claim does not hold", r.fatal, r.errors)
		}
	})

	t.Run("wrong code with wrong kind on refused row fails", func(t *testing.T) {
		tsv := "1347\t1347\tengine\tfx_a\tr/kind\tDROP VIEW a\n"
		tsvPath, fxDir := writeSimFiles(t, tsv, fixtures)

		r, _ := runCheckSimulation(t, Simulation{
			Path:     tsvPath,
			Fixtures: fxDir,
			Setup: func(_ testing.TB, _ simulate.Fixture) SimFixture {
				return SimFixture{
					Objects: []Object{{Key: walk.Key{Space: "relation", Name: "a"}, Kind: "table"}},
					Report: func(_ testing.TB, _ string) *Finding {
						return &Finding{Statement: 0, Message: "claims view"}
					},
					Violation: func(_ testing.TB, _ string) *walk.Violation {
						return &walk.Violation{
							Check: walk.Check{Key: walk.Key{Space: "relation", Name: "a"}},
							Got:   walk.Present,
							Kind:  "view",
						}
					},
				}
			},
			Codes: func(string, string) []string { return []string{"1051"} },
			Same:  same,
		})
		if r.fatal != "" || len(r.errors) != 1 || !strings.Contains(r.errors[0], "is a table, not a view") {
			t.Errorf("fatal=%q errors=%q; want one error about wrong kind", r.fatal, r.errors)
		}
	})

	t.Run("unjudged row fails", func(t *testing.T) {
		tsv := "?\t?\t-\tfx_a\tr/unj\tCREATE TABLE x (id INT)\n"
		tsvPath, fxDir := writeSimFiles(t, tsv, fixtures)

		r, _ := runCheckSimulation(t, Simulation{
			Path:     tsvPath,
			Fixtures: fxDir,
			Setup: func(_ testing.TB, _ simulate.Fixture) SimFixture {
				return SimFixture{
					Report:    func(_ testing.TB, _ string) *Finding { return nil },
					Violation: func(_ testing.TB, _ string) *walk.Violation { return nil },
				}
			},
			Same: same,
		})
		if r.fatal != "" || len(r.errors) != 1 || !strings.Contains(r.errors[0], "unjudged") {
			t.Errorf("fatal=%q errors=%q; want one error about unjudged row", r.fatal, r.errors)
		}
	})
}
