// Command compatcheck fails a change to committed engine answers that only a
// recording may make. It compares each compat.tsv the commits since <base>
// change with the file at their merge base (see compat.Guard), and fails a
// pin bump that also changes non-test Go code: a pin bump's diff must be
// what the release changed, so omni's fixes go in later changes.
//
// CI runs it on every pull request. From the repo root:
//
//	go run ./enginetest/cmd/compatcheck origin/main
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/bytebase/omni/enginetest/compat"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: compatcheck <base-ref>")
		os.Exit(2)
	}
	base := os.Args[1]
	mergeBase := strings.TrimSpace(git("merge-base", base, "HEAD"))
	changed := strings.Fields(git("diff", "--name-only", base+"...HEAD"))

	bumped := map[string]bool{}
	var code []string
	for _, f := range changed {
		if strings.HasSuffix(f, "/testdata/pins.tsv") {
			// A pins.tsv that is new introduces the engine's pin; only
			// changing one that existed bumps it.
			if _, existed := show(mergeBase, f); existed {
				bumped[top(f)] = true
			}
		}
		if strings.HasSuffix(f, ".go") && !strings.HasSuffix(f, "_test.go") {
			code = append(code, f)
		}
	}

	var problems []string
	if len(bumped) > 0 && len(code) > 0 {
		problems = append(problems, fmt.Sprintf("a pin bump changes the pin and the recorded answers only, but this change also edits %s; fix omni in a separate change",
			strings.Join(code, ", ")))
	}
	for _, f := range changed {
		if path.Base(f) != "compat.tsv" {
			continue
		}
		before, ok := show(mergeBase, f)
		if !ok {
			continue // a new file: every row is new
		}
		after, ok := show("HEAD", f)
		if !ok {
			continue // a deleted file
		}
		b, err := compat.Parse(mergeBase+":"+f, before)
		if err != nil {
			fail(err)
		}
		a, err := compat.Parse(f, after)
		if err != nil {
			fail(err)
		}
		for _, p := range compat.Guard(b, a, bumped[top(f)]) {
			problems = append(problems, f+": "+p)
		}
	}
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "compatcheck: %d problems:\n%s\n", len(problems), strings.Join(problems, "\n"))
		os.Exit(1)
	}
	fmt.Println("compatcheck: committed answers change only as recordings change them")
}

// top returns the top-level directory of a repo path: its engine.
func top(f string) string {
	dir, _, _ := strings.Cut(f, "/")
	return dir
}

// show returns a file's content at a revision, and false if it has none there.
func show(rev, f string) ([]byte, bool) {
	out, err := exec.Command("git", "show", rev+":"+f).Output()
	if err != nil {
		return nil, false
	}
	return out, true
}

func git(args ...string) string {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		fail(fmt.Errorf("git %s: %w", strings.Join(args, " "), err))
	}
	return string(out)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "compatcheck:", err)
	os.Exit(1)
}
