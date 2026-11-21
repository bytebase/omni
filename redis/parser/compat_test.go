package parser

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/redis/internal/testcontainer"
)

// TestMain removes the pinned Redis container if a recording or verification
// started one.
func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.Terminate()
	os.Exit(code)
}

// TestCompatibility holds omni's redis/parser to testdata/compat.tsv: a
// statement for each alternative of ../testdata/RedisCommands.g4 (generated
// from Redis 7.4.12's own COMMAND table) and its near-misses, plus lexical
// edge-case suite rows, with omni's verdict and Redis 7.4.12's.
//
// In normal runs, this test reads testdata/compat.tsv offline and starts no
// Docker container. With OMNI_RECORD=1 it asks the pinned redis:7.4.12-alpine
// server via COMMAND INFO and writes verdicts; with OMNI_RECORD=verify it
// re-checks every committed verdict read-only against the pinned server.
func TestCompatibility(t *testing.T) {
	td := filepath.Join("..", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join(td, "pins.tsv"),
		Grammar: filepath.Join(td, "RedisCommands.g4"),
		Tokens:  filepath.Join(td, "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict returns omni's verdict on one Redis command line.
func omniVerdict(line string) string {
	if _, err := ParseCommand(line); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

const compatCallTimeout = 10 * time.Second

// compatJudge starts the pinned redis:7.4.12-alpine server and judges each
// command line against the live server's COMMAND INFO metadata (command name,
// subcommand table, and exact/minimum arity) without executing destructive,
// mutating, or blocking commands. On startup, it cross-checks COMMAND INFO
// verdicts against raw inline TCP execution on safe read-only / non-mutating
// commands.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	client := testcontainer.Judge(t, pin)
	verifyJudgeAgainstExecution(t, client)

	return func(stmt string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), compatCallTimeout)
		defer cancel()
		return judgeCommandLine(ctx, client, stmt)
	}
}

// judgeCommandLine asks the live Redis server `COMMAND INFO <cmd>` to decide
// whether line is a syntactically valid Redis command with a valid command
// name, subcommand (if container), and argument count (arity).
//
// It intentionally does not call parser.Tokenize or parser.matchArity so the
// judge stays independent of the parser under test.
func judgeCommandLine(ctx context.Context, client *redis.Client, line string) (string, error) {
	if strings.ContainsRune(line, '\n') {
		return compat.Reject, nil
	}
	fields, err := splitRedisInline(line)
	if err != nil || len(fields) == 0 {
		return compat.Reject, nil
	}
	// Redis's COMMAND INFO accepts "parent|sub" lookups, whereas live command
	// dispatch rejects a pipe in the first token with ERR unknown command.
	if strings.ContainsRune(fields[0], '|') {
		return compat.Reject, nil
	}

	raw, err := client.Do(ctx, "COMMAND", "INFO", fields[0]).Result()
	if err != nil {
		return "", fmt.Errorf("COMMAND INFO %q: %w", fields[0], err)
	}
	list, ok := raw.([]any)
	if !ok || len(list) != 1 {
		return "", fmt.Errorf("COMMAND INFO %q: unexpected reply %T (len %d)", fields[0], raw, len(list))
	}
	if list[0] == nil {
		return compat.Reject, nil
	}
	entry, ok := list[0].([]any)
	if !ok || len(entry) < 10 {
		return "", fmt.Errorf("COMMAND INFO %q: malformed entry %T", fields[0], list[0])
	}
	name, ok := entry[0].(string)
	if !ok {
		return "", fmt.Errorf("COMMAND INFO %q: malformed name %T", fields[0], entry[0])
	}
	name = strings.ToLower(name)
	if name != strings.ToLower(fields[0]) {
		return compat.Reject, nil
	}
	arity64, ok := entry[1].(int64)
	if !ok || arity64 == 0 {
		return "", fmt.Errorf("COMMAND INFO %q: malformed arity %T (%v)", fields[0], entry[1], entry[1])
	}
	var subs []any
	if entry[9] != nil {
		subs, ok = entry[9].([]any)
		if !ok {
			return "", fmt.Errorf("COMMAND INFO %q: malformed subcommands %T", fields[0], entry[9])
		}
	}

	if len(subs) > 0 {
		if len(fields) == 1 {
			if judgeMatchArity(int(arity64), 1) {
				return compat.Accept, nil
			}
			return compat.Reject, nil
		}
		wantSub := name + "|" + strings.ToLower(fields[1])
		for _, subRaw := range subs {
			subEntry, ok := subRaw.([]any)
			if !ok || len(subEntry) < 2 {
				return "", fmt.Errorf("COMMAND INFO %q: malformed subcommand entry %T", fields[0], subRaw)
			}
			subName, ok := subEntry[0].(string)
			if !ok {
				return "", fmt.Errorf("COMMAND INFO %q: malformed subcommand name %T", fields[0], subEntry[0])
			}
			if strings.ToLower(subName) != wantSub {
				continue
			}
			subArity64, ok := subEntry[1].(int64)
			if !ok || subArity64 == 0 {
				return "", fmt.Errorf("COMMAND INFO %q: malformed subcommand arity %T (%v)", subName, subEntry[1], subEntry[1])
			}
			if judgeMatchArity(int(subArity64), len(fields)) {
				return compat.Accept, nil
			}
			return compat.Reject, nil
		}
		return compat.Reject, nil
	}

	if judgeMatchArity(int(arity64), len(fields)) {
		return compat.Accept, nil
	}
	return compat.Reject, nil
}

func judgeMatchArity(arity, argc int) bool {
	if arity > 0 {
		return argc == arity
	}
	if arity < 0 {
		return argc >= -arity
	}
	return false
}

// splitRedisInline mirrors Redis's own sdssplitargs (src/sds.c), used by
// Redis's inline protocol and redis-cli, keeping the judge independent of
// parser.Tokenize.
func splitRedisInline(line string) ([]string, error) {
	var args []string
	i, n := 0, len(line)
	for i < n {
		for i < n && isRedisSpace(line[i]) {
			i++
		}
		if i >= n {
			break
		}
		var cur strings.Builder
		inDouble := false
		inSingle := false
		done := false
		for !done {
			switch {
			case inDouble:
				if i >= n {
					return nil, errors.New("unbalanced quotes")
				}
				if line[i] == '\\' && i+3 < n && line[i+1] == 'x' && isHex(line[i+2]) && isHex(line[i+3]) {
					cur.WriteByte((fromHex(line[i+2]) << 4) | fromHex(line[i+3]))
					i += 4
				} else if line[i] == '\\' && i+1 < n {
					i++
					switch line[i] {
					case 'n':
						cur.WriteByte('\n')
					case 'r':
						cur.WriteByte('\r')
					case 't':
						cur.WriteByte('\t')
					case 'b':
						cur.WriteByte('\b')
					case 'a':
						cur.WriteByte('\a')
					default:
						cur.WriteByte(line[i])
					}
					i++
				} else if line[i] == '"' {
					if i+1 < n && !isRedisSpace(line[i+1]) {
						return nil, errors.New("unbalanced quotes")
					}
					i++
					done = true
				} else {
					cur.WriteByte(line[i])
					i++
				}
			case inSingle:
				if i >= n {
					return nil, errors.New("unbalanced quotes")
				}
				if line[i] == '\\' && i+1 < n && line[i+1] == '\'' {
					cur.WriteByte('\'')
					i += 2
				} else if line[i] == '\'' {
					if i+1 < n && !isRedisSpace(line[i+1]) {
						return nil, errors.New("unbalanced quotes")
					}
					i++
					done = true
				} else {
					cur.WriteByte(line[i])
					i++
				}
			default:
				if i >= n || isRedisSpace(line[i]) {
					done = true
				} else if line[i] == '"' {
					inDouble = true
					i++
				} else if line[i] == '\'' {
					inSingle = true
					i++
				} else {
					cur.WriteByte(line[i])
					i++
				}
			}
		}
		args = append(args, cur.String())
	}
	return args, nil
}

func isRedisSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n' || b == '\v' || b == '\f'
}

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

func fromHex(b byte) byte {
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10
	default:
		return b - 'A' + 10
	}
}

// verifyJudgeAgainstExecution proves on the live pinned server that
// judgeCommandLine agrees 100% with Redis 7.4.12's actual inline protocol and
// command dispatcher across safe read-only / non-mutating commands covering
// quoting, escapes, fixed arity, variadic arity, container commands, COMMAND's
// dual standalone/container behavior, unknown commands/subcommands, and
// pipe-name rejection.
func verifyJudgeAgainstExecution(t testing.TB, client *redis.Client) {
	t.Helper()
	cases := []struct {
		line string
		want string
	}{
		{"GET", compat.Reject},
		{"GET a", compat.Accept},
		{"GET a b", compat.Reject},
		{"PING", compat.Accept},
		{"PING a", compat.Accept},
		{"DBSIZE", compat.Accept},
		{"DBSIZE a", compat.Reject},
		{"COMMAND", compat.Accept},
		{"COMMAND zzz", compat.Reject},
		{"COMMAND COUNT", compat.Accept},
		{"COMMAND COUNT a", compat.Reject},
		{"COMMAND INFO", compat.Accept},
		{"COMMAND INFO GET", compat.Accept},
		{"ACL", compat.Reject},
		{"ACL zzz", compat.Reject},
		{"ACL LIST", compat.Accept},
		{"ACL LIST a", compat.Reject},
		{"OBJECT", compat.Reject},
		{"OBJECT zzz", compat.Reject},
		{"OBJECT ENCODING", compat.Reject},
		{"OBJECT ENCODING a", compat.Accept},
		{"OBJECT ENCODING a b", compat.Reject},
		{"acl|list", compat.Reject},
		{"UNKNOWNCMD a", compat.Reject},
		{`GET "unclosed`, compat.Reject},
		{`GET 'unclosed`, compat.Reject},
		{`GET "a\"`, compat.Reject},
		{`GET "two words"`, compat.Accept},
		{`GET 'c:\path'`, compat.Accept},
		{`GET ""`, compat.Accept},
		{`GET ''`, compat.Accept},
		{"GET #channel", compat.Accept},
		{"# not a comment", compat.Reject},
		{"gEt a", compat.Accept},
		{"AcL lIsT", compat.Accept},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	addr := client.Options().Addr

	for _, tc := range cases {
		judged, err := judgeCommandLine(ctx, client, tc.line)
		if err != nil {
			t.Fatalf("judgeCommandLine(%q): %v", tc.line, err)
		}
		if judged != tc.want {
			t.Fatalf("judgeCommandLine(%q) = %s, want %s", tc.line, judged, tc.want)
		}
		execVerdict, replyLine, err := sendInlineCommand(addr, tc.line)
		if err != nil {
			t.Fatalf("sendInlineCommand(%q): %v", tc.line, err)
		}
		if execVerdict != judged {
			t.Fatalf("live inline execution of %q disagreed with COMMAND INFO judge: exec=%s (reply=%q), judge=%s",
				tc.line, execVerdict, replyLine, judged)
		}
	}
}

// sendInlineCommand sends one raw inline command line to Redis over TCP and
// classifies Redis's RESP2 first reply line as syntax accept or reject.
func sendInlineCommand(addr, line string) (verdict, firstLine string, err error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = conn.Close() }()
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return "", "", err
	}
	if _, err := conn.Write([]byte(line + "\r\n")); err != nil {
		return "", "", err
	}
	r := bufio.NewReader(conn)
	firstLine, err = r.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	firstLine = strings.TrimRight(firstLine, "\r\n")
	if msg, ok := strings.CutPrefix(firstLine, "-"); ok {
		if strings.HasPrefix(msg, "ERR Protocol error: unbalanced quotes") ||
			strings.HasPrefix(msg, "ERR unknown command ") ||
			strings.HasPrefix(msg, "ERR unknown subcommand ") ||
			strings.HasPrefix(msg, "ERR wrong number of arguments ") {
			return compat.Reject, firstLine, nil
		}
	}
	return compat.Accept, firstLine, nil
}
