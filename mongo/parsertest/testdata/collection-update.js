// db.collection.update() - Update documents (deprecated - use updateOne or updateMany)

// Basic update with $set
db.users.update({ name: "alice" }, { $set: { age: 26 } })
db.users.update({ _id: 1 }, { $set: { status: "active" } })

// Update with multiple operators
db.users.update({ name: "bob" }, {
    $set: { status: "active" },
    $inc: { loginCount: 1 },
    $currentDate: { lastLogin: true }
})

// Update with $inc
db.products.update({ sku: "abc123" }, { $inc: { qty: -1 } })
db.users.update({ _id: ObjectId("507f1f77bcf86cd799439011") }, { $inc: { score: 10, attempts: 1 } })

// Update with $unset
db.users.update({ name: "charlie" }, { $unset: { tempField: "" } })

// Update with $push
db.users.update({ name: "dave" }, { $push: { tags: "premium" } })
db.users.update({ name: "eve" }, { $push: { scores: { $each: [85, 90, 78] } } })

// Update with $pull
db.users.update({ name: "frank" }, { $pull: { tags: "inactive" } })
db.users.update({ name: "grace" }, { $pull: { items: { qty: { $lte: 0 } } } })

// Update with $addToSet
db.users.update({ name: "henry" }, { $addToSet: { roles: "editor" } })
db.users.update({ name: "ivy" }, { $addToSet: { tags: { $each: ["a", "b", "c"] } } })

// Update multiple documents (multi option)
db.users.update({ status: "pending" }, { $set: { status: "active" } }, { multi: true })
db.orders.update({ shipped: false }, { $set: { shipped: true, shippedAt: Date() } }, { multi: true })

// Update with upsert
db.users.update({ email: "new@example.com" }, { $set: { name: "New User" } }, { upsert: true })
db.settings.update({ key: "theme" }, { $set: { value: "dark" } }, { upsert: true })

// Update with multi and upsert
db.counters.update({ name: "pageViews" }, { $inc: { count: 1 } }, { upsert: true, multi: false })

// Replace document (no update operators)
db.users.update({ _id: 1 }, { name: "replaced", status: "new" })

// Update with arrayFilters
db.inventory.update(
    { item: "abc123" },
    { $set: { "sizes.$[elem].qty": 0 } },
    { arrayFilters: [{ "elem.size": "small" }] }
)

// Update with write concern
db.users.update({ name: "jack" }, { $set: { verified: true } }, { writeConcern: { w: "majority" } })

// Update with collation
db.products.update({ name: "cafe" }, { $set: { available: true } }, { collation: { locale: "fr", strength: 1 } })

// Update with hint
db.users.update({ status: "active" }, { $set: { lastChecked: Date() } }, { hint: { status: 1 } })

// Update with let variables
db.users.update(
    { $expr: { $eq: ["$_id", "$$targetId"] } },
    { $set: { found: true } },
    { let: { targetId: ObjectId("507f1f77bcf86cd799439011") } }
)

// Collection access patterns
db["users"].update({ x: 1 }, { $set: { y: 2 } })
db["users"].update({ a: 1 }, { $inc: { b: 1 } }, { multi: true })
db.getCollection("users").update({ status: "old" }, { $set: { status: "new" } })
db.getCollection("config").update({ key: "value" }, { $set: { data: {} } }, { upsert: true })
db["user-preferences"].update({ userId: 1 }, { $set: { theme: "dark" } })
