// db.collection.dropIndexes() - Drop all indexes or specified indexes from a collection

// Drop all indexes (except _id)
db.users.dropIndexes()
db.orders.dropIndexes()

// Drop specific index by name
db.users.dropIndexes("email_1")

// Drop multiple indexes by name
db.users.dropIndexes(["email_1", "username_1"])

// Drop index by key pattern
db.orders.dropIndexes({ customerId: 1 })

// Collection access patterns
db["users"].dropIndexes()
db.getCollection("users").dropIndexes()
db["orders"].dropIndexes("status_1")
db.getCollection("orders").dropIndexes(["idx1", "idx2"])
