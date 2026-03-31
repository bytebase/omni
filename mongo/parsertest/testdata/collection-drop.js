// db.collection.drop() - Drop a collection from the database

// Basic drop
db.users.drop()
db.orders.drop()
db.tempData.drop()

// Drop with write concern
db.logs.drop({ writeConcern: { w: "majority" } })

// Collection access patterns
db["users"].drop()
db.getCollection("users").drop()
db["temp-data"].drop()
db.getCollection("archived.orders").drop()
