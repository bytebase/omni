// rs.addArb() - Add an arbiter to the replica set

// Add arbiter by hostname
rs.addArb("arbiter:27017")

// Add arbiter with full hostname and port
rs.addArb("arbiter.example.com:27017")
