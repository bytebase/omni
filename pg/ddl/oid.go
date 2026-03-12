package ddl

// Well-known namespace OIDs matching PostgreSQL.
const (
	InvalidOid          uint32 = 0
	PGCatalogNamespace  uint32 = 11
	PGToastNamespace    uint32 = 99
	PublicNamespace     uint32 = 2200
	FirstNormalObjectId uint32 = 16384
)

// OIDGenerator allocates monotonically increasing OIDs starting at FirstNormalObjectId.
type OIDGenerator struct{ next uint32 }

// NewOIDGenerator returns a generator starting at FirstNormalObjectId.
func NewOIDGenerator() *OIDGenerator {
	return &OIDGenerator{next: FirstNormalObjectId}
}

// Next returns the next available OID.
func (g *OIDGenerator) Next() uint32 {
	oid := g.next
	g.next++
	return oid
}
