// cursor.sort() - Sort documents in the result set

// Basic ascending sort
db.users.find().sort({ name: 1 })

// Basic descending sort
db.users.find().sort({ age: -1 })

// Multiple field sort
db.users.find().sort({ status: 1, name: 1 })
db.users.find().sort({ createdAt: -1, name: 1 })

// With query filter
db.users.find({ status: "active" }).sort({ name: 1 })
db.users.find({ age: { $gt: 18 } }).sort({ age: -1 })

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).limit(10)
db.users.find().sort({ createdAt: -1 }).skip(20).limit(10)
db.users.find({ status: "active" }).sort({ name: 1 }).limit(100).skip(0)

// Sort on nested fields
db.users.find().sort({ "address.city": 1 })
db.users.find().sort({ "profile.score": -1, "profile.name": 1 })

// Sort with text score (for text search)
db.articles.find({ $text: { $search: "mongodb" } }).sort({ score: { $meta: "textScore" } })
