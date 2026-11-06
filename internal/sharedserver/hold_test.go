package sharedserver

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// holderEnv names the env file for which a re-run of this test binary holds a
// server instead of running tests.
const holderEnv = "OMNI_SHAREDSERVER_HOLDER_ENV"

// holdForTest is the holder TestHoldServerGoesWithItsHolder kills: it holds a
// MySQL 8.0 until the process dies.
func holdForTest(envFile string) {
	if err := Hold(context.Background(), envFile, "mysql80"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

// TestHoldServerGoesWithItsHolder kills a holder as the kernel's OOM killer and
// systemd do, with SIGKILL, so the holder removes nothing itself, and waits for
// its server to stop answering. The reaper removing it is what keeps a CI job
// that dies from leaving its servers on the runners' machine.
func TestHoldServerGoesWithItsHolder(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, "shared.env")
	logFile, err := os.Create(filepath.Join(dir, "holder.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer logFile.Close()
	holderLog := func() string {
		b, _ := os.ReadFile(logFile.Name())
		return string(b)
	}

	holder := exec.Command(os.Args[0], "-test.run=^$")
	holder.Env = append(os.Environ(), holderEnv+"="+envFile)
	holder.Stdout, holder.Stderr = logFile, logFile
	if err := holder.Start(); err != nil {
		t.Fatal(err)
	}
	// exited closes when the holder ends, so the startup loop, the kill and the
	// cleanup can each wait for it.
	exited := make(chan struct{})
	var waitErr error
	go func() {
		waitErr = holder.Wait()
		close(exited)
	}()
	t.Cleanup(func() {
		_ = holder.Process.Kill()
		<-exited
	})

	var dsn string
	for deadline := time.Now().Add(5 * time.Minute); dsn == ""; {
		if b, err := os.ReadFile(envFile); err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				if v, ok := strings.CutPrefix(line, MySQL80.Env+"="); ok {
					dsn = v
				}
			}
			if dsn == "" {
				t.Fatalf("%s names no %s:\n%s", envFile, MySQL80.Env, b)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("holder wrote no env file in 5 minutes:\n%s", holderLog())
		}
		select {
		case <-exited:
			t.Fatalf("holder exited before its server answered: %v\n%s", waitErr, holderLog())
		case <-time.After(time.Second):
		}
	}
	ctx := context.Background()
	db, _, err := ping(ctx, dsn)
	if err != nil {
		t.Fatalf("held server does not answer: %v", err)
	}
	db.Close()

	if err := holder.Process.Signal(syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	<-exited
	killed := time.Now()
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		db, _, err := ping(pingCtx, dsn)
		cancel()
		if err != nil {
			t.Logf("server gone %v after its holder was killed", time.Since(killed).Round(time.Second))
			return
		}
		db.Close()
		// The reaper takes about 12 s. One run on a machine with its swap full,
		// which is when this matters, failed a 2-minute window (2026-11-06).
		if time.Since(killed) > 5*time.Minute {
			t.Fatalf("server still answers 5 minutes after its holder was killed:\n%s", holderLog())
		}
		time.Sleep(2 * time.Second)
	}
}
