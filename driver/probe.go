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
