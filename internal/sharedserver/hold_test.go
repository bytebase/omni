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
	exited := make(chan error, 1)
	go func() { exited <- holder.Wait() }()
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
		case err := <-exited:
			t.Fatalf("holder exited before its server answered: %v\n%s", err, holderLog())
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
	exited <- nil // for the cleanup
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
		if time.Since(killed) > 2*time.Minute {
			t.Fatalf("server still answers 2 minutes after its holder was killed:\n%s", holderLog())
		}
		time.Sleep(2 * time.Second)
	}
}
