package sharedserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/testcontainers/testcontainers-go"
)

// Hold starts a testcontainer of each server named, "mysql80" or "oceanbase",
// writes the variables that point this package at them to envFile once every one
// answers, and keeps them until ctx is done. CI runs it through cmd/hold in the
// background for the length of a job, so the job's packages share one server of
// each.
//
// The servers belong to this process's testcontainers reaper, which removes them
// about ten seconds after the process ends, however it ends. Servers CI started
// with docker run and removed in an always() step outlived every job that died:
// when the kernel kills a runner's process for memory, systemd stops the whole
// runner, so no later step runs. Two OceanBase servers left that way, at 3.5 GiB
// each, helped exhaust the runners' machine until it had to be reset
// (2026-11-06).
func Hold(ctx context.Context, envFile string, names ...string) error {
	if len(names) == 0 {
		return errors.New("no server named")
	}
	servers := make([]*MySQLServer, len(names))
	envs := make([][]string, len(names))
	for i, name := range names {
		s, env, ok := holdable(name)
		if !ok {
			return fmt.Errorf("unknown server %q", name)
		}
		servers[i], envs[i] = s, env
	}

	containers := make([]testcontainers.Container, len(servers))
	dsns := make([]string, len(servers))
	errs := make([]error, len(servers))
	var wg sync.WaitGroup
	for i, s := range servers {
		wg.Go(func() {
			containers[i], dsns[i], errs[i] = s.run(ctx)
			if errs[i] == nil {
				errs[i] = waitReady(ctx, dsns[i])
			}
			if errs[i] != nil {
				errs[i] = fmt.Errorf("%s: %w", names[i], errs[i])
			}
		})
	}
	wg.Wait()
	defer func() {
		for _, c := range containers {
			_ = testcontainers.TerminateContainer(c)
		}
	}()
	if err := errors.Join(errs...); err != nil {
		return err
	}

	var out strings.Builder
	for i, dsn := range dsns {
		for _, env := range envs[i] {
			fmt.Fprintf(&out, "%s=%s\n", env, dsn)
		}
	}
	// A reader polling for the file never sees it half written.
	if err := os.WriteFile(envFile+".tmp", []byte(out.String()), 0o644); err != nil {
		return err
	}
	if err := os.Rename(envFile+".tmp", envFile); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

// holdable returns the server Hold starts for name, with the variables that point
// tests at it.
func holdable(name string) (*MySQLServer, []string, bool) {
	switch name {
	case "mysql80":
		return MySQL80, []string{MySQL80.Env}, true
	case "oceanbase":
		// The oracles that run as root and the tests that keep to a database of
		// their own take the same server.
		return OceanBase, []string{OceanBase.Env, PrivateOceanBase.Env}, true
	}
	return nil, nil, false
}
