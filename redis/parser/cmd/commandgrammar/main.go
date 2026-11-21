// Command commandgrammar generates redis/testdata/RedisCommands.g4 and
// redis/parser/commands_gen.go from the Redis command table snapshot
// (redis/testdata/commands.tsv).
//
// With -from-live, it starts the pinned default server from
// redis/testdata/pins.tsv (or connects to -addr if provided), verifies
// redis_version against the pin, queries COMMAND, and rewrites
// redis/testdata/commands.tsv before generating the grammar and Go registry.
//
// From the repo root:
//
//	go run ./redis/parser/cmd/commandgrammar
//	go run ./redis/parser/cmd/commandgrammar -from-live
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
)

const (
	defaultPins        = "redis/testdata/pins.tsv"
	defaultCommandsTSV = "redis/testdata/commands.tsv"
	defaultGrammarOut  = "redis/testdata/RedisCommands.g4"
	defaultGoOut       = "redis/parser/commands_gen.go"
)

// commandEntry is one callable Redis command or subcommand call shape.
type commandEntry struct {
	Command    string // lowercase top-level command name (e.g. "get", "acl")
	Subcommand string // lowercase subcommand name (e.g. "list"), or "" for standalone
	Arity      int    // Redis arity: >0 exact token count, <0 minimum token count (-N)
}

func main() {
	pinsPath := flag.String("pins", defaultPins, "path to redis/testdata/pins.tsv")
	commandsPath := flag.String("commands", defaultCommandsTSV, "path to redis/testdata/commands.tsv")
	grammarOut := flag.String("grammar", defaultGrammarOut, "output path for RedisCommands.g4")
	goOut := flag.String("go-out", defaultGoOut, "output path for redis/parser/commands_gen.go")
	fromLive := flag.Bool("from-live", false, "query COMMAND from the pinned Redis server (starts Docker container unless -addr is set) and update -commands")
	addr := flag.String("addr", "", "existing Redis server host:port to query when -from-live is set")
	flag.Parse()

	pin, err := enginetest.LoadPin(*pinsPath, "default")
	if err != nil {
		fatalf("load pin: %v", err)
	}

	var entries []commandEntry
	if *fromLive {
		entries, err = fetchLiveCommands(context.Background(), pin, *addr)
		if err != nil {
			fatalf("fetch live commands: %v", err)
		}
		tsv := renderCommandsTSV(pin, entries)
		if err := writeFile(*commandsPath, []byte(tsv)); err != nil {
			fatalf("write %s: %v", *commandsPath, err)
		}
		fmt.Printf("wrote %s (%d entries)\n", *commandsPath, len(entries))
	} else {
		raw, err := os.ReadFile(*commandsPath)
		if err != nil {
			fatalf("read %s: %v", *commandsPath, err)
		}
		entries, err = parseCommandsTSV(string(raw))
		if err != nil {
			fatalf("parse %s: %v", *commandsPath, err)
		}
	}

	grammarText, err := renderGrammar(pin, entries)
	if err != nil {
		fatalf("render grammar: %v", err)
	}
	if err := writeFile(*grammarOut, []byte(grammarText)); err != nil {
		fatalf("write %s: %v", *grammarOut, err)
	}

	goSrc, err := renderGoRegistry(pin, entries)
	if err != nil {
		fatalf("render Go registry: %v", err)
	}
	if err := writeFile(*goOut, goSrc); err != nil {
		fatalf("write %s: %v", *goOut, err)
	}

	topCount, containerCount, subCount := summarize(entries)
	fmt.Printf("wrote %s and %s (%d top-level commands, %d containers, %d subcommands, %d callable shapes)\n",
		*grammarOut, *goOut, topCount, containerCount, subCount, len(entries))
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "commandgrammar: "+format+"\n", args...)
	os.Exit(1)
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func fetchLiveCommands(ctx context.Context, pin enginetest.Pin, addr string) ([]commandEntry, error) {
	if addr == "" {
		c, err := testcontainers.Run(ctx, pin.ImageRef(),
			testcontainers.WithCmd("redis-server", "--save", "", "--appendonly", "no"),
			testcontainers.WithExposedPorts("6379/tcp"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("Ready to accept connections").WithStartupTimeout(2*time.Minute)),
		)
		if err != nil {
			return nil, fmt.Errorf("start %s: %w", pin.ImageRef(), err)
		}
		defer func() { _ = testcontainers.TerminateContainer(c) }()
		if err := verifyContainerDigest(ctx, c, pin.Digest); err != nil {
			return nil, err
		}
		host, err := c.Host(ctx)
		if err != nil {
			return nil, err
		}
		port, err := c.MappedPort(ctx, "6379/tcp")
		if err != nil {
			return nil, err
		}
		addr = host + ":" + port.Port()
	}

	client := redis.NewClient(&redis.Options{Addr: addr, Protocol: 2})
	defer func() { _ = client.Close() }()

	if err := verifyVersion(ctx, client, pin.Tag); err != nil {
		return nil, err
	}

	res, err := client.Do(ctx, "COMMAND").Result()
	if err != nil {
		return nil, fmt.Errorf("COMMAND: %w", err)
	}
	return decodeCommandReply(res)
}

func verifyContainerDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect image: %w", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("container image digests %v do not match pinned %s", img.RepoDigests, digest)
}

func verifyVersion(ctx context.Context, client *redis.Client, want string) error {
	info, err := client.Info(ctx, "server").Result()
	if err != nil {
		return fmt.Errorf("INFO server: %w", err)
	}
	for _, line := range strings.Split(info, "\n") {
		if after, ok := strings.CutPrefix(line, "redis_version:"); ok {
			got := strings.TrimSpace(after)
			if got != want {
				return fmt.Errorf("server reports redis_version %q, want %q", got, want)
			}
			return nil
		}
	}
	return fmt.Errorf("INFO server did not report redis_version")
}

// decodeCommandReply parses the RESP2 array returned by Redis's COMMAND.
// Each element has the layout:
//
//	[0] name (string)
//	[1] arity (int64)
//	[2..8] flags, key positions, acl categories, tips, key specs
//	[9] subcommands ([]any), each with the same layout and name "<cmd>|<sub>"
func decodeCommandReply(raw any) ([]commandEntry, error) {
	topList, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("expected []any from COMMAND, got %T", raw)
	}
	var entries []commandEntry
	for _, item := range topList {
		cmd, err := decodeEntry(item)
		if err != nil {
			return nil, err
		}
		if len(cmd.subs) == 0 {
			entries = append(entries, commandEntry{
				Command: cmd.name,
				Arity:   cmd.arity,
			})
			continue
		}
		// Container command: standalone invocation is only valid when 1 token
		// satisfies the container's arity (e.g. COMMAND has arity -1, whereas
		// ACL/CLIENT/CLUSTER/... have arity -2 and require a subcommand).
		if cmd.arity == 1 || cmd.arity == -1 {
			entries = append(entries, commandEntry{
				Command: cmd.name,
				Arity:   cmd.arity,
			})
		}
		for _, sub := range cmd.subs {
			parent, subName, ok := strings.Cut(sub.name, "|")
			if !ok || parent != cmd.name {
				return nil, fmt.Errorf("unexpected subcommand name %q under %q", sub.name, cmd.name)
			}
			if len(sub.subs) > 0 {
				return nil, fmt.Errorf("unsupported nested subcommands under %q", sub.name)
			}
			entries = append(entries, commandEntry{
				Command:    cmd.name,
				Subcommand: subName,
				Arity:      sub.arity,
			})
		}
	}
	sortEntries(entries)
	return entries, nil
}

type decodedCommand struct {
	name  string
	arity int
	subs  []decodedCommand
}

func decodeEntry(raw any) (decodedCommand, error) {
	arr, ok := raw.([]any)
	if !ok || len(arr) < 2 {
		return decodedCommand{}, fmt.Errorf("malformed COMMAND entry %T", raw)
	}
	name, ok := arr[0].(string)
	if !ok || name == "" {
		return decodedCommand{}, fmt.Errorf("malformed command name %T", arr[0])
	}
	arity64, ok := arr[1].(int64)
	if !ok || arity64 == 0 {
		return decodedCommand{}, fmt.Errorf("malformed arity %T (%v) for %q", arr[1], arr[1], name)
	}
	out := decodedCommand{
		name:  strings.ToLower(name),
		arity: int(arity64),
	}
	if len(arr) >= 10 && arr[9] != nil {
		subList, ok := arr[9].([]any)
		if !ok {
			return decodedCommand{}, fmt.Errorf("malformed subcommands %T for %q", arr[9], name)
		}
		for _, subRaw := range subList {
			sub, err := decodeEntry(subRaw)
			if err != nil {
				return decodedCommand{}, err
			}
			out.subs = append(out.subs, sub)
		}
	}
	return out, nil
}

func sortEntries(entries []commandEntry) {
	slices.SortFunc(entries, func(a, b commandEntry) int {
		if a.Command != b.Command {
			return strings.Compare(a.Command, b.Command)
		}
		return strings.Compare(a.Subcommand, b.Subcommand)
	})
}

func renderCommandsTSV(pin enginetest.Pin, entries []commandEntry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Redis %s command table snapshot from COMMAND on the pinned default server.\n", pin.Tag)
	fmt.Fprintf(&b, "# Source: %s@%s (%s)\n", pin.Image, pin.Digest, pin.Tag)
	b.WriteString("# Regenerate with: go run ./redis/parser/cmd/commandgrammar -from-live\n")
	b.WriteString("#\n")
	b.WriteString("# Arity follows Redis's COMMAND convention: positive N means exact total token\n")
	b.WriteString("# count N (command + subcommand + arguments); negative -N means at least N.\n")
	b.WriteString("# Container commands that require a subcommand (arity -2) appear only via their\n")
	b.WriteString("# subcommand rows; COMMAND (arity -1) also has a standalone row with empty\n")
	b.WriteString("# subcommand.\n")
	b.WriteString("# columns: command, subcommand, arity\n")
	for _, e := range entries {
		fmt.Fprintf(&b, "%s\t%s\t%d\n", e.Command, e.Subcommand, e.Arity)
	}
	return b.String()
}

func parseCommandsTSV(src string) ([]commandEntry, error) {
	var entries []commandEntry
	seen := map[string]bool{}
	for i, line := range strings.Split(src, "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cols := strings.Split(line, "\t")
		if len(cols) != 3 {
			return nil, fmt.Errorf("line %d: want 3 tab-separated columns (command, subcommand, arity), got %d", i+1, len(cols))
		}
		cmd := strings.ToLower(cols[0])
		sub := strings.ToLower(cols[1])
		if cmd == "" {
			return nil, fmt.Errorf("line %d: empty command name", i+1)
		}
		arity, err := strconv.Atoi(cols[2])
		if err != nil || arity == 0 {
			return nil, fmt.Errorf("line %d: invalid arity %q", i+1, cols[2])
		}
		key := cmd + "\x00" + sub
		if seen[key] {
			return nil, fmt.Errorf("line %d: duplicate entry %q %q", i+1, cmd, sub)
		}
		seen[key] = true
		entries = append(entries, commandEntry{
			Command:    cmd,
			Subcommand: sub,
			Arity:      arity,
		})
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no command entries found")
	}
	sortEntries(entries)
	return entries, nil
}

func summarize(entries []commandEntry) (topCount, containerCount, subCount int) {
	tops := map[string]bool{}
	containers := map[string]bool{}
	for _, e := range entries {
		tops[e.Command] = true
		if e.Subcommand != "" {
			containers[e.Command] = true
			subCount++
		}
	}
	return len(tops), len(containers), subCount
}

// renderGrammar emits redis/testdata/RedisCommands.g4.
func renderGrammar(pin enginetest.Pin, entries []commandEntry) (string, error) {
	containers := map[string][]commandEntry{}
	var containerNames []string
	var standalone []commandEntry

	hasSub := map[string]bool{}
	for _, e := range entries {
		if e.Subcommand != "" {
			hasSub[e.Command] = true
		}
	}
	for _, e := range entries {
		if hasSub[e.Command] {
			if _, ok := containers[e.Command]; !ok {
				containerNames = append(containerNames, e.Command)
			}
			containers[e.Command] = append(containers[e.Command], e)
		} else {
			standalone = append(standalone, e)
		}
	}
	slices.Sort(containerNames)

	var b strings.Builder
	b.WriteString("// Redis command-line call-shape grammar as ANTLR 4, generated from the\n")
	fmt.Fprintf(&b, "// official Redis %s COMMAND table by redis/parser/cmd/commandgrammar.\n", pin.Tag)
	b.WriteString("// Do not edit by hand: rerun `go run ./redis/parser/cmd/commandgrammar`.\n")
	b.WriteString("//\n")
	b.WriteString("// Scope (§2.2 of docs/specs/2026-11-11-engine-compatibility-without-servers-design.md):\n")
	b.WriteString("// covers every Redis command name, container subcommand, exact/minimum arity\n")
	b.WriteString("// call shape, and shell-like command-line argument token form.\n")
	b.WriteString("//\n")
	fmt.Fprintf(&b, "// pin: %s\n\n", pin.String())
	b.WriteString("grammar RedisCommands;\n\n")

	b.WriteString("commandLine\n    : command EOF\n    ;\n\n")

	b.WriteString("command\n    : standaloneCommand\n")
	for _, c := range containerNames {
		fmt.Fprintf(&b, "    | %sCommand\n", ruleIdent(c))
	}
	b.WriteString("    ;\n\n")

	// Standalone commands.
	b.WriteString("// Top-level Redis commands without subcommands.\n")
	b.WriteString("standaloneCommand\n")
	for i, e := range standalone {
		prefix := "    : "
		if i > 0 {
			prefix = "    | "
		}
		alt, err := formatCommandAlt(e, false)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "%s%s // COMMAND %s (arity %d)\n", prefix, alt, e.Command, e.Arity)
	}
	b.WriteString("    ;\n\n")

	// Container commands, one rule per container.
	for _, c := range containerNames {
		fmt.Fprintf(&b, "// Container command %s and its subcommands.\n", strings.ToUpper(c))
		fmt.Fprintf(&b, "%sCommand\n", ruleIdent(c))
		for i, e := range containers[c] {
			prefix := "    : "
			if i > 0 {
				prefix = "    | "
			}
			alt, err := formatCommandAlt(e, true)
			if err != nil {
				return "", err
			}
			label := e.Command
			if e.Subcommand != "" {
				label += "|" + e.Subcommand
			}
			fmt.Fprintf(&b, "%s%s // COMMAND %s (arity %d)\n", prefix, alt, label, e.Arity)
		}
		b.WriteString("    ;\n\n")
	}

	b.WriteString("// Zero or more extra arguments for variadic (negative-arity) commands.\n")
	b.WriteString("variadicArgs\n    :\n    | arg variadicArgs\n    ;\n\n")

	b.WriteString("// Command-line argument token forms (see redis/testdata/tokens.tsv).\n")
	b.WriteString("arg\n")
	b.WriteString("    : WORD\n")
	b.WriteString("    | SINGLE_QUOTED\n")
	b.WriteString("    | DOUBLE_QUOTED\n")
	b.WriteString("    | HASH_ARG\n")
	b.WriteString("    | INT_LITERAL\n")
	b.WriteString("    | FLOAT_LITERAL\n")
	b.WriteString("    ;\n\n")

	b.WriteString("// Lexer token classes for command-line arguments.\n")
	b.WriteString("SINGLE_QUOTED : '\\'' (~['\\r\\n])* '\\'' ;\n")
	b.WriteString("DOUBLE_QUOTED : '\"' ('\\\\' . | ~[\"\\\\\\r\\n])* '\"' ;\n")
	b.WriteString("HASH_ARG      : '#' ~[ \\t\\r\\n]+ ;\n")
	b.WriteString("FLOAT_LITERAL : '-'? [0-9]+ '.' [0-9]+ ;\n")
	b.WriteString("INT_LITERAL   : '-'? [0-9]+ ;\n")
	b.WriteString("WORD          : ~[ \\t\\r\\n'\"#]+ ;\n")
	b.WriteString("WS            : [ \\t\\r]+ -> skip ;\n")

	return b.String(), nil
}

func formatCommandAlt(e commandEntry, isContainer bool) (string, error) {
	var parts []string
	parts = append(parts, "'"+strings.ToUpper(e.Command)+"'")
	fixedTokens := 1
	if e.Subcommand != "" {
		parts = append(parts, "'"+strings.ToUpper(e.Subcommand)+"'")
		fixedTokens = 2
	}

	// Standalone invocation of a container command (e.g. COMMAND with arity -1):
	// Redis dispatches any second token as a subcommand, so the standalone
	// alternative takes no trailing arguments.
	if isContainer && e.Subcommand == "" {
		return strings.Join(parts, " "), nil
	}

	var requiredArgs int
	var variadic bool
	if e.Arity > 0 {
		requiredArgs = e.Arity - fixedTokens
	} else {
		requiredArgs = -e.Arity - fixedTokens
		variadic = true
	}
	if requiredArgs < 0 {
		return "", fmt.Errorf("invalid arity %d for %s %s", e.Arity, e.Command, e.Subcommand)
	}
	for range requiredArgs {
		parts = append(parts, "arg")
	}
	if variadic {
		parts = append(parts, "variadicArgs")
	}
	return strings.Join(parts, " "), nil
}

func ruleIdent(cmd string) string {
	return strings.ReplaceAll(cmd, "-", "_")
}

// renderGoRegistry emits redis/parser/commands_gen.go.
func renderGoRegistry(pin enginetest.Pin, entries []commandEntry) ([]byte, error) {
	type spec struct {
		standaloneArity int
		subs            []commandEntry
	}
	byCmd := map[string]*spec{}
	var cmdNames []string
	for _, e := range entries {
		s, ok := byCmd[e.Command]
		if !ok {
			s = &spec{}
			byCmd[e.Command] = s
			cmdNames = append(cmdNames, e.Command)
		}
		if e.Subcommand == "" {
			s.standaloneArity = e.Arity
		} else {
			s.subs = append(s.subs, e)
		}
	}
	slices.Sort(cmdNames)

	var buf bytes.Buffer
	buf.WriteString("// Code generated by redis/parser/cmd/commandgrammar; DO NOT EDIT.\n")
	fmt.Fprintf(&buf, "// Source: %s@%s (%s)\n\n", pin.Image, pin.Digest, pin.Tag)
	buf.WriteString("package parser\n\n")
	buf.WriteString("// CommandSpec describes one top-level Redis command's arity and subcommands.\n")
	buf.WriteString("// Arity follows Redis's COMMAND table: positive N requires exact token count\n")
	buf.WriteString("// N, negative -N requires at least N tokens, and 0 means the command cannot\n")
	buf.WriteString("// be invoked without a subcommand.\n")
	buf.WriteString("type CommandSpec struct {\n")
	buf.WriteString("\tName            string\n")
	buf.WriteString("\tStandaloneArity int\n")
	buf.WriteString("\tSubcommands     map[string]int\n")
	buf.WriteString("}\n\n")
	buf.WriteString("var commandRegistry = map[string]CommandSpec{\n")
	for _, name := range cmdNames {
		s := byCmd[name]
		if len(s.subs) == 0 {
			fmt.Fprintf(&buf, "\t%q: {Name: %q, StandaloneArity: %d},\n", name, name, s.standaloneArity)
			continue
		}
		fmt.Fprintf(&buf, "\t%q: {\n", name)
		fmt.Fprintf(&buf, "\t\tName: %q,\n", name)
		if s.standaloneArity != 0 {
			fmt.Fprintf(&buf, "\t\tStandaloneArity: %d,\n", s.standaloneArity)
		}
		buf.WriteString("\t\tSubcommands: map[string]int{\n")
		for _, sub := range s.subs {
			fmt.Fprintf(&buf, "\t\t\t%q: %d,\n", sub.Subcommand, sub.Arity)
		}
		buf.WriteString("\t\t},\n")
		buf.WriteString("\t},\n")
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("gofmt generated registry: %w", err)
	}
	return formatted, nil
}
