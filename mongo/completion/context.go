package completion

import "github.com/bytebase/omni/mongo/parser"

// completionContext identifies the kind of completion expected.
type completionContext int

const (
	contextTopLevel      completionContext = iota // start of input or after semicolon
	contextAfterDbDot                             // db.|
	contextAfterCollDot                           // db.users.|
	contextAfterBracket                           // db[|
	contextInsideArgs                             // db.users.find(|
	contextDocumentKey                            // {| or {age: 1, |
	contextQueryOperator                          // {age: {$|
	contextAggStage                               // [{$|
	contextCursorChain                            // db.users.find().|
	contextShowTarget                             // show |
	contextAfterRsDot                             // rs.|
	contextAfterShDot                             // sh.|
)

// detectContext analyzes the token sequence to determine the completion context.
func detectContext(tokens []parser.Token) completionContext {
	return contextTopLevel
}
