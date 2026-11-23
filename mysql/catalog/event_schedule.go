package catalog

import (
	"strings"
	"time"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// eventTimeLayout is how SHOW CREATE EVENT prints a time.
const eventTimeLayout = "2006-01-02 15:04:05"

// eventLiteralTime evaluates a schedule time the catalog can: a date or
// datetime string literal, plus or minus INTERVAL n of a single unit. MySQL
// evaluates AT, STARTS and ENDS when the statement runs and stores the time.
// loc is the session's time zone, which a literal is read in.
func eventLiteralTime(expr nodes.ExprNode, loc *time.Location) (time.Time, bool) {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return eventLiteralTime(e.Expr, loc)
	case *nodes.StringLit:
		for _, layout := range []string{eventTimeLayout, "2006-01-02"} {
			if t, err := time.ParseInLocation(layout, strings.TrimSpace(e.Value), loc); err == nil {
				return t, true
			}
		}
	case *nodes.BinaryExpr:
		if e.Op != nodes.BinOpAdd && e.Op != nodes.BinOpSub {
			return time.Time{}, false
		}
		base, ok := eventLiteralTime(e.Left, loc)
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
		switch strings.ToUpper(iv.Unit) {
		case "SECOND":
			return base.Add(time.Duration(n) * time.Second), true
		case "MINUTE":
			return base.Add(time.Duration(n) * time.Minute), true
		case "HOUR":
			return base.Add(time.Duration(n) * time.Hour), true
		case "DAY":
			return base.AddDate(0, 0, int(n)), true
		case "WEEK":
			return base.AddDate(0, 0, 7*int(n)), true
		}
	}
	return time.Time{}, false
}

func eventIntervalCount(expr nodes.ExprNode) (int64, bool) {
	if lit, ok := expr.(*nodes.IntLit); ok {
		return lit.Value, true
	}
	return 0, false
}

// validateEventTimes refuses a schedule whose literal times MySQL refuses: an
// EVERY interval that is not positive (1542), and ENDS before STARTS (1543),
// STARTS being the statement's time when it is not written (mysql:8.0.46).
// loc is the session's time zone; with known false, the catalog cannot tell
// it, and a time is compared with now only when the other is written too.
func validateEventTimes(schedule *nodes.EventSchedule, now time.Time, loc *time.Location, known bool) error {
	if schedule == nil || schedule.Every == nil {
		return nil
	}
	if n, ok := eventIntervalCount(schedule.Every); ok && n <= 0 {
		return &Error{Code: 1542, SQLState: "HY000", Message: "INTERVAL is either not positive or too big"}
	}
	if schedule.Ends == nil {
		return nil
	}
	ends, ok := eventLiteralTime(schedule.Ends, loc)
	if !ok {
		return nil
	}
	starts := now
	if schedule.Starts != nil {
		if starts, ok = eventLiteralTime(schedule.Starts, loc); !ok {
			return nil
		}
	} else if !known {
		return nil
	}
	if ends.Before(starts) {
		return &Error{Code: 1543, SQLState: "HY000", Message: "ENDS is either invalid or before STARTS"}
	}
	return nil
}

// eventScheduleText is the schedule SHOW CREATE EVENT prints for one a
// statement wrote: an AT time the catalog can evaluate is printed as the time
// it is, so AT '2030-01-01 00:00:00' + INTERVAL 1 DAY reads back as AT
// '2030-01-02 00:00:00' (mysql:8.0.46). It also returns that time, and
// whether there is one.
func eventScheduleText(schedule *nodes.EventSchedule, loc *time.Location) (string, time.Time, bool) {
	if schedule == nil {
		return "", time.Time{}, false
	}
	if schedule.At != nil {
		if at, ok := eventLiteralTime(schedule.At, loc); ok {
			return "AT '" + at.Format(eventTimeLayout) + "'", at, true
		}
	}
	return schedule.RawText, time.Time{}, false
}

// sessionTimeLocation is the session's time zone, which MySQL reads a
// schedule's literals in, and whether the catalog can tell it.
func (c *Catalog) sessionTimeLocation() (*time.Location, bool) {
	tz := c.resolveEffectiveTimeZone()
	switch {
	case !tz.valid:
		return time.UTC, false
	case tz.loc != nil:
		return tz.loc, true
	case tz.hasFixed:
		return time.FixedZone("", tz.fixedMin*60), true
	}
	return time.UTC, false
}
