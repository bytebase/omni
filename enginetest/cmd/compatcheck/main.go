// Command compatcheck fails a change to committed engine answers that only a
// recording may make. It compares each compat.tsv and simulate.tsv the
// commits since <base> change with the file at their merge base (see
// compat.Guard), and fails a pin bump that also changes non-test Go code: a
// pin bump's diff must be what the release changed, so omni's fixes go in
// later changes. A pin bump is a file recorded at a new pin, as its # pin:
// line shows; an edit to pins.tsv alone bumps nothing.
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
	"github.com/bytebase/omni/enginetest/simulate"
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
	var code, guarded []string
	for _, f := range changed {
		if strings.HasSuffix(f, ".go") && !strings.HasSuffix(f, "_test.go") {
			code = append(code, f)
		}
		guard := guards[path.Base(f)]
		if guard == nil {
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
		problems, repinned, err := guard(mergeBase+":"+f, before, f, after)
		if err != nil {
			fail(err)
		}
		for _, msg := range problems {
			guarded = append(guarded, f+": "+msg)
		}
		if repinned {
			bumped[top(f)] = true
		}
	}

	var problems []string
	if len(bumped) > 0 && len(code) > 0 {
		problems = append(problems, fmt.Sprintf("a pin bump changes the pin and the recorded answers only, but this change also edits %s; fix omni in a separate change",
			strings.Join(code, ", ")))
	}
	problems = append(problems, guarded...)
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "compatcheck: %d problems:\n%s\n", len(problems), strings.Join(problems, "\n"))
		os.Exit(1)
	}
	fmt.Println("compatcheck: committed answers change only as recordings change them")
}

// guards holds each kind of committed answer file to its rules: what a
// change may not do to it, and whether the change recorded it at a new pin.
var guards = map[string]func(beforeName string, before []byte, afterName string, after []byte) ([]string, bool, error){
	"compat.tsv": func(beforeName string, before []byte, afterName string, after []byte) ([]string, bool, error) {
		b, err := compat.Parse(beforeName, before)
		if err != nil {
			return nil, false, err
		}
		a, err := compat.Parse(afterName, after)
		if err != nil {
			return nil, false, err
		}
		return compat.Guard(b, a), compat.PinChanged(b, a), nil
	},
	"simulate.tsv": func(beforeName string, before []byte, afterName string, after []byte) ([]string, bool, error) {
		b, err := simulate.Parse(beforeName, before)
		if err != nil {
			return nil, false, err
		}
		a, err := simulate.Parse(afterName, after)
		if err != nil {
			return nil, false, err
		}
		return simulate.Guard(b, a), simulate.PinChanged(b, a), nil
	},
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
