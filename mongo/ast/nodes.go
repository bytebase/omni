// Package ast defines MongoDB (mongosh) parse tree node types.
package ast

// Loc represents a source location range (byte offsets).
type Loc struct {
	Start int // inclusive start byte offset
	End   int // exclusive end byte offset
}

// Node is the interface implemented by all mongosh parse tree nodes.
type Node interface {
	GetLoc() Loc
}

// ---------------------------------------------------------------------------
// Statement nodes — one per mongosh command family
// ---------------------------------------------------------------------------

// ShowCommand represents show dbs, show collections, show profile, etc.
type ShowCommand struct {
	Command string // e.g. "dbs", "collections", "profile"
	Loc     Loc
}

func (n *ShowCommand) GetLoc() Loc { return n.Loc }

// CollectionStatement represents db.collection.method().cursor()...
// This is the most complex and common statement type.
type CollectionStatement struct {
	Database       string         // explicit db name or "" for implicit
	Collection     string         // collection name
	Method         *MethodCall    // primary method (find, insert, etc.)
	CursorMethods  []CursorMethod // chained cursor modifiers (.sort, .limit, etc.)
	Loc            Loc
}

func (n *CollectionStatement) GetLoc() Loc { return n.Loc }

// DatabaseStatement represents db.method() calls (e.g. db.getCollectionNames()).
type DatabaseStatement struct {
	Database string      // explicit db name or "" for implicit
	Method   *MethodCall // method call
	Loc      Loc
}

func (n *DatabaseStatement) GetLoc() Loc { return n.Loc }

// BulkStatement represents db.collection.initializeOrderedBulkOp()/initializeUnorderedBulkOp() chains.
type BulkStatement struct {
	Database   string          // explicit db name or "" for implicit
	Collection string          // collection name
	Ordered    bool            // true for initializeOrderedBulkOp
	Operations []BulkOperation // chained operations
	Execute    bool            // whether .execute() was called
	Loc        Loc
}

func (n *BulkStatement) GetLoc() Loc { return n.Loc }

// ConnectionStatement represents Mongo(), connect(), db.getMongo() chains.
type ConnectionStatement struct {
	Method *MethodCall   // the connection method
	Chain  []*MethodCall // chained methods
	Loc    Loc
}

func (n *ConnectionStatement) GetLoc() Loc { return n.Loc }

// RsStatement represents rs.method() calls (replica set).
type RsStatement struct {
	Method *MethodCall // method call
	Loc    Loc
}

func (n *RsStatement) GetLoc() Loc { return n.Loc }

// ShStatement represents sh.method() calls (sharding).
type ShStatement struct {
	Method *MethodCall // method call
	Loc    Loc
}

func (n *ShStatement) GetLoc() Loc { return n.Loc }

// EncryptionStatement represents db.getMongo().getKeyVault()/getClientEncryption() chains.
type EncryptionStatement struct {
	Database string        // explicit db name or "" for implicit
	Chain    []*MethodCall // method chain
	Loc      Loc
}

func (n *EncryptionStatement) GetLoc() Loc { return n.Loc }

// PlanCacheStatement represents db.collection.getPlanCache() chains.
type PlanCacheStatement struct {
	Database   string      // explicit db name or "" for implicit
	Collection string      // collection name
	Method     *MethodCall // plan cache method (clear, list, etc.)
	Loc        Loc
}

func (n *PlanCacheStatement) GetLoc() Loc { return n.Loc }

// SpStatement represents sp.method() or sp.x.method() calls (query planner).
type SpStatement struct {
	Namespace string      // sub-namespace (e.g. "x" in sp.x.method())
	Method    *MethodCall // method call
	Loc       Loc
}

func (n *SpStatement) GetLoc() Loc { return n.Loc }

// NativeFunctionCall represents top-level mongosh functions: sleep(), load(), etc.
type NativeFunctionCall struct {
	Name string // function name
	Args []Node // arguments
	Loc  Loc
}

func (n *NativeFunctionCall) GetLoc() Loc { return n.Loc }

// ---------------------------------------------------------------------------
// Supporting types
// ---------------------------------------------------------------------------

// CursorMethod represents a chained cursor modifier like .sort({x:1}).
type CursorMethod struct {
	Name string // method name (sort, limit, skip, etc.)
	Args []Node // arguments
	Loc  Loc
}

// MethodCall represents a method call with name and arguments.
type MethodCall struct {
	Name string // method name
	Args []Node // arguments
	Loc  Loc
}

func (n *MethodCall) GetLoc() Loc { return n.Loc }

// BulkOperation represents a single operation in a bulk chain.
type BulkOperation struct {
	Name string // operation name (insert, find, update, etc.)
	Args []Node // arguments
	Loc  Loc
}

// ---------------------------------------------------------------------------
// Value / Expression nodes
// ---------------------------------------------------------------------------

// Document represents a JSON/BSON document literal { key: value, ... }.
type Document struct {
	Pairs []KeyValue
	Loc   Loc
}

func (n *Document) GetLoc() Loc { return n.Loc }

// KeyValue represents a single key-value pair in a document.
type KeyValue struct {
	Key    string // field name
	KeyLoc Loc    // location of the key
	Value  Node   // value expression
	Loc    Loc
}

func (n *KeyValue) GetLoc() Loc { return n.Loc }

// Array represents a JSON/BSON array literal [ ... ].
type Array struct {
	Elements []Node
	Loc      Loc
}

func (n *Array) GetLoc() Loc { return n.Loc }

// StringLiteral represents a quoted string value.
type StringLiteral struct {
	Value string // the unescaped string content
	Loc   Loc
}

func (n *StringLiteral) GetLoc() Loc { return n.Loc }

// NumberLiteral represents a numeric value (integer or float).
type NumberLiteral struct {
	Value string // raw text of the number
	Loc   Loc
}

func (n *NumberLiteral) GetLoc() Loc { return n.Loc }

// BoolLiteral represents true or false.
type BoolLiteral struct {
	Value bool
	Loc   Loc
}

func (n *BoolLiteral) GetLoc() Loc { return n.Loc }

// NullLiteral represents null.
type NullLiteral struct {
	Loc Loc
}

func (n *NullLiteral) GetLoc() Loc { return n.Loc }

// RegexLiteral represents a /pattern/flags regular expression.
type RegexLiteral struct {
	Pattern string
	Flags   string
	Loc     Loc
}

func (n *RegexLiteral) GetLoc() Loc { return n.Loc }

// HelperCall represents a BSON helper like ObjectId("..."), NumberLong(1), etc.
type HelperCall struct {
	Name string // helper name (ObjectId, NumberLong, etc.)
	Args []Node // arguments
	Loc  Loc
}

func (n *HelperCall) GetLoc() Loc { return n.Loc }

// Identifier represents an unquoted identifier or variable reference.
type Identifier struct {
	Name string
	Loc  Loc
}

func (n *Identifier) GetLoc() Loc { return n.Loc }
