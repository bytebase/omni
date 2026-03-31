// db.collection.renameCollection() - Rename a collection

// Basic rename
db.users.renameCollection("users_old")
db.orders.renameCollection("archived_orders")

// Rename with dropTarget option (drops target collection if exists)
db.users.renameCollection("users_backup", true)
db.orders.renameCollection("orders_v2", false)

// Rename to collection in same database
db.tempData.renameCollection("permanentData")
db.staging.renameCollection("production")

// Collection access patterns
db["users"].renameCollection("users_backup")
db.getCollection("users").renameCollection("users_archive")
db["old-data"].renameCollection("archived-data")
