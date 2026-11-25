package review

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	obcatalog "github.com/bytebase/omni/oceanbase/catalog"
	"github.com/bytebase/omni/review"
)

// spaceSequence names a sequence's name for its length limit; the ledger keeps
// no sequence.
const spaceSequence = "sequence"

// sequence lowers CREATE, ALTER or DROP SEQUENCE, false for OceanBase's other
// objects, which the walk does not follow. A sequence's name is apart from
// every name the walk keeps, and Bytebase's OceanBase sync reads no sequence,
// so the walk holds a sequence statement only to what OceanBase refuses of it
// whatever the database holds (obcatalog.CheckSequenceStatement), and goes on
// past it.
func (b *builder) sequence(v *ast.ObjectDDLStmt) bool {
	if !strings.HasSuffix(v.Action, " SEQUENCE") {
		return false
	}
	rng := b.anchor()
	if b.tooLong(spaceSequence, v.Name, rng) {
		return true
	}
	var refusal *obcatalog.Error
	if !errors.As(obcatalog.CheckSequenceStatement(v), &refusal) {
		return true
	}
	code, message := review.UnsupportedForm, refusal.Message
	switch refusal.Code {
	case 1059:
		code = review.NameTooLong
	case 1149:
		message = "only ALTER SEQUENCE takes RESTART"
	}
	b.fail(rng, nil, code, "Sequence "+review.Code(ident(v.Name))+": "+message+".")
	return true
}
