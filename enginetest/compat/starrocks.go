package compat

import (
	"fmt"
	"strings"
)

// StarRocksProbe is the statement a StarRocks judge sends first, in the
// same request as the statement it asks about: StarRocksProbe + ";\n" +
// statement. It parses, and fails when it runs, with an error that names it.
const StarRocksProbe = "SELECT omni_compat_probe()"

// StarRocksVerdict reads a StarRocks server's answer to a statement sent
// behind StarRocksProbe: whether StarRocks takes its syntax. err is what
// sending the request returned, and message returns the server's error
// message, or false for an error the server did not send (a dropped
// connection, a timeout).
//
// StarRocks parses a request whole, building every statement's AST, before
// it runs the first, and stops at the first that fails. So the probe's own
// error ("No matching function with signature: omni_compat_probe()") is an
// acceptance, and nothing of the statement runs. Every other error the
// server sends came from the parse, and is a rejection: a syntax error
// ("Getting syntax error at line 2, column 10"), and what the AST builder
// refuses or fails on ("not support.", a ClassCastException's message).
//
// A request that ran, which the probe never does, and an error the server
// did not send are no verdict, and come back as an error, so the caller
// fails rather than records agreement.
func StarRocksVerdict(err error, message func(error) (string, bool)) (string, error) {
	if err == nil {
		return "", fmt.Errorf("the request ran, and %s cannot", StarRocksProbe)
	}
	msg, ok := message(err)
	if !ok {
		return "", err
	}
	if strings.Contains(msg, starRocksProbeName) {
		return Accept, nil
	}
	return Reject, nil
}

// starRocksProbeName is what the probe's error names.
const starRocksProbeName = "omni_compat_probe"
