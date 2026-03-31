// cursor.projection() / cursor.project() - Specify fields to return

// Include specific fields
db.users.find().projection({ name: 1, email: 1 })
db.users.find().project({ name: 1, email: 1 })

// Exclude _id
db.users.find().projection({ name: 1, email: 1, _id: 0 })

// Exclude specific fields
db.users.find().projection({ password: 0, secretKey: 0 })

// With query filter
db.users.find({ status: "active" }).projection({ name: 1, status: 1 })
db.users.find({ age: { $gte: 18 } }).project({ name: 1, age: 1 })

// Nested field projection
db.users.find().projection({ "address.city": 1, "address.country": 1 })
db.users.find().project({ "profile.name": 1, "profile.avatar": 1 })

// Chained with other cursor methods
db.users.find().sort({ name: 1 }).projection({ name: 1 }).limit(10)
db.users.find({ status: "active" }).project({ name: 1, email: 1 }).skip(10).limit(10)

// Array projection operators
db.posts.find().projection({ comments: { $slice: 5 } })
db.posts.find().projection({ comments: { $elemMatch: { score: { $gt: 5 } } } })
