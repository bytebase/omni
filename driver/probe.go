package driver

import "context"

// PrivateDatabaseLinkOwner is implemented by a driver whose engine lets a
// statement reach another database over a link the session's account owns,
// without the statement naming it (Oracle: a private database link). A caller
// that must know whether such a link could carry data out asks before it runs
// the statement.
type PrivateDatabaseLinkOwner interface {
	// OwnsPrivateDatabaseLink reports whether the driver's account owns a
	// private database link.
	OwnsPrivateDatabaseLink(ctx context.Context) (bool, error)
}

// AllNodesQuerier is implemented by a Session whose engine can run a script on
// every node of a cluster as well as on the one node a command routes to
// (Redis: a cluster's shards). Bytebase's query API offers a per-request "run
// on every node" choice, which Session.Query cannot honour because it routes a
// command the way the server does, so a caller that offers that choice asks the
// session it got from Querier.Session for this interface. Only a clustered
// topology answers it; a session of one server refuses.
type AllNodesQuerier interface {
	// QueryAllNodes returns one Result per command of script, whose rows are
	// every node's answer, ordered by node address.
	QueryAllNodes(ctx context.Context, script string, opts QueryOptions) ([]*Result, error)
}
