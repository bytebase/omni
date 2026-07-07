package main

import (
	"fmt"

	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// omniTiDBVerdict runs one statement through the omni tidb parser.
// A panic is treated as reject with a panic marker — corpus statements must
// never crash the sweep (and panics are themselves findings).
func omniTiDBVerdict(sql string) (v Verdict, errMsg string) {
	defer func() {
		if r := recover(); r != nil {
			v = VerdictReject
			errMsg = "PANIC: " + fmt.Sprint(r)
		}
	}()
	if _, err := tidbparser.Parse(sql); err != nil {
		return VerdictReject, err.Error()
	}
	return VerdictAccept, ""
}
