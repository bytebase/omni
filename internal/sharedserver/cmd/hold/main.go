// Command hold starts the shared servers that CI's test packages use and keeps
// them until it is stopped; sharedserver.Hold says why. CI runs it in the
// background and appends the file it writes to $GITHUB_ENV once it appears:
//
//	hold -env "$RUNNER_TEMP/shared.env" mysql80 oceanbase > hold.log 2>&1 &
//
// SIGTERM or SIGINT removes the servers and exits. Any other end leaves them to
// the testcontainers reaper, which removes them about ten seconds later.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytebase/omni/internal/sharedserver"
)

func main() {
	env := flag.String("env", "", "file to write the servers' DSN variables to once they answer")
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: hold -env FILE mysql80|oceanbase...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *env == "" || flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := sharedserver.Hold(ctx, *env, flag.Args()...)
	stop()
	if err != nil {
		log.Fatal(err)
	}
}
