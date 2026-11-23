package catalog

import (
	"fmt"
	"math"
	"strings"
	"time"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// eventTimeLayout is how SHOW CREATE EVENT prints a time.
const eventTimeLayout = "2006-01-02 15:04:05"

// eventZoneMargin is longer than any time zone is ahead of or behind UTC. The
// catalog does not model the session's time_zone, which MariaDB reads a
// schedule's literals in, so it compares a literal with the present only when
// the two are further apart than this, which no time zone can change.
const eventZoneMargin = 26 * time.Hour

// eventLiteralTime evaluates a schedule time the catalog can: a date or
// datetime string literal, plus or minus INTERVAL n of a single unit, read as
// UTC. MariaDB evaluates AT, STARTS and ENDS when the statement runs and stores
// the time.
func eventLiteralTime(expr nodes.ExprNode) (time.Time, bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return eventLiteralTime(e.Expr)
	case *nodes.StringLit:
		for _, layout := range []string{eventTimeLayout, "2006-01-02"} {
			if t, err := time.ParseInLocation(layout, strings.TrimSpace(e.Value), time.UTC); err == nil {
				return t, true
			}
		}
	case *nodes.BinaryExpr:
		if e.Op != nodes.BinOpAdd && e.Op != nodes.BinOpSub {
			return time.Time{}, false
		}
		base, ok := eventLiteralTime(e.Left)
		iv, isInterval := e.Right.(*nodes.IntervalExpr)
		if !ok || !isInterval {
			return time.Time{}, false
		}
		n, ok := eventIntervalCount(iv.Value)
		if !ok {
			return time.Time{}, false
		}
		if e.Op == nodes.BinOpSub {
			n = -n
		}
		// Seconds, minutes and hours are added as whole days and the rest,
		// which a Duration holds whatever the count.
		var unit int64
		switch strings.ToUpper(iv.Unit) {
		case "SECOND":
			unit = 1
		case "MINUTE":
			unit = 60
		case "HOUR":
			unit = 3600
		}
		if unit > 0 {
			if n > math.MaxInt64/unit || n < -math.MaxInt64/unit {
				return time.Time{}, false
			}
			secs := n * unit
			return base.AddDate(0, 0, int(secs/86400)).Add(time.Duration(secs%86400) * time.Second), true
		}
		switch strings.ToUpper(iv.Unit) {
		case "DAY":
			return base.AddDate(0, 0, int(n)), true
		case "WEEK":
			return base.AddDate(0, 0, 7*int(n)), true
		}
	}
	return time.Time{}, false
}

func eventIntervalCount(expr nodes.ExprNode) (int64, bool) {
	switch e := expr.(type) {
	case *nodes.IntLit:
		return e.Value, true
	case *nodes.UnaryExpr:
		if n, ok := eventIntervalCount(e.Operand); ok && e.Op == nodes.UnaryMinus {
			return -n, true
		}
	}
	return 0, false
}

// validateEventTimes refuses a schedule whose literal times MariaDB refuses,
// measured on mariadb:11.8.8, in the order Event_parse_data::check_parse_data
// reads them: a time string that names no time (1525), an EVERY interval in
// microseconds (1235) or one that is not positive (1542), and ENDS before
// STARTS (1543), STARTS being the statement's time when it is not written. A
// zero now compares no time with the present: a declarative schema states the
// event it wants, which must not stop loading as time passes.
func validateEventTimes(schedule *nodes.EventSchedule, now time.Time) error {
	if schedule == nil {
		return nil
	}
	if err := refuseEventTimeText("AT", schedule.At); err != nil {
		return err
	}
	if schedule.Every == nil {
		return nil
	}
	if strings.Contains(schedule.Unit, "MICROSECOND") {
		return &Error{Code: 1235, SQLState: "42000", Message: "This version of MariaDB doesn't yet support 'MICROSECOND'"}
	}
	if n, ok := eventIntervalCount(schedule.Every); ok && n <= 0 {
		return &Error{Code: 1542, SQLState: "HY000", Message: "INTERVAL is either not positive or too big"}
	}
	if err := refuseEventTimeText("STARTS", schedule.Starts); err != nil {
		return err
	}
	if schedule.Ends == nil {
		return nil
	}
	// An ENDS that names no time is 1543 as well (mariadb:11.8.8: ENDS 'zzz').
	if namesNoTime(schedule.Ends) {
		return errEventEndsBeforeStarts()
	}
	ends, ok := eventLiteralTime(schedule.Ends)
	if !ok {
		return nil
	}
	if schedule.Starts == nil {
		if !now.IsZero() && ends.Before(now.Add(-eventZoneMargin)) {
			return errEventEndsBeforeStarts()
		}
		return nil
	}
	starts, ok := eventLiteralTime(schedule.Starts)
	if ok && ends.Before(starts) {
		return errEventEndsBeforeStarts()
	}
	return nil
}

// refuseEventTimeText is 1525 for an AT or STARTS time that names none:
// mariadb:11.8.8 answers AT 'zzz' and STARTS 'zzz' with it.
func refuseEventTimeText(clause string, expr nodes.ExprNode) error {
	if !namesNoTime(expr) {
		return nil
	}
	return &Error{Code: 1525, SQLState: "HY000", Message: fmt.Sprintf("Incorrect %s value: '%s'", clause, expr.(*nodes.StringLit).Value)}
}

// namesNoTime reports a schedule time written as a string with no digit in it.
// A string that holds digits is left to the server, whose reading of a
// datetime is lenient.
func namesNoTime(expr nodes.ExprNode) bool {
	lit, ok := expr.(*nodes.StringLit)
	return ok && !strings.ContainsAny(lit.Value, "0123456789")
}

func errEventEndsBeforeStarts() error {
	return &Error{Code: 1543, SQLState: "HY000", Message: "ENDS is either invalid or before STARTS"}
}

// eventScheduleText is the schedule SHOW CREATE EVENT prints for one a
// statement wrote: an AT time the catalog can evaluate is printed as the time
// it is, so AT '2030-01-01 00:00:00' + INTERVAL 1 DAY reads back as AT
// '2030-01-02 00:00:00' (mariadb:11.8.8). It also returns that time, and
// whether there is one.
func eventScheduleText(schedule *nodes.EventSchedule) (string, time.Time, bool) {
	if schedule == nil {
		return "", time.Time{}, false
	}
	if schedule.At != nil {
		if at, ok := eventLiteralTime(schedule.At); ok {
			return "AT '" + at.Format(eventTimeLayout) + "'", at, true
		}
	}
	return schedule.RawText, time.Time{}, false
}
