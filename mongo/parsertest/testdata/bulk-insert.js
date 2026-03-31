// Bulk.insert() - Queue an insert operation

// Basic insert
db.users.initializeUnorderedBulkOp().insert({ name: "alice", age: 25 })

// Multiple inserts
db.users.initializeOrderedBulkOp().insert({ name: "bob" }).insert({ name: "charlie" })

// Insert with nested document
db.users.initializeUnorderedBulkOp().insert({ name: "dave", address: { city: "NYC", zip: "10001" } })

// Insert with array
db.users.initializeOrderedBulkOp().insert({ name: "eve", tags: ["admin", "user"] })

// Insert with helper functions
db.users.initializeUnorderedBulkOp().insert({ name: "frank", createdAt: ISODate(), id: UUID("550e8400-e29b-41d4-a716-446655440000") })

// Insert and execute
db.users.initializeOrderedBulkOp().insert({ name: "grace" }).execute()
