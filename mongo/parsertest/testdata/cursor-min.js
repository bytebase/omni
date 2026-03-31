// cursor.min() - Specify inclusive lower bound for index-based query

// Basic usage with index key pattern
db.users.find({ age: { $lte: 65 } }).min({ age: 18 })

// With compound index
db.users.find().min({ status: "a", name: "a" })
db.products.find({ category: "electronics" }).min({ category: "electronics", price: 0 })

// Chained with max() for range
db.users.find().min({ age: 18 }).max({ age: 65 })
db.products.find().min({ price: 10 }).max({ price: 100 })

// With hint (required for min to work correctly)
db.users.find().min({ age: 21 }).hint({ age: 1 })
db.products.find().min({ price: 50 }).hint({ price: 1 })

// Chained with other cursor methods
db.users.find().min({ age: 18 }).hint({ age: 1 }).sort({ name: 1 })
db.users.find().min({ age: 21 }).hint({ age: 1 }).limit(10)

// With multiple fields
db.events.find().min({ date: ISODate("2024-01-01"), priority: 1 })

// Range query with both bounds
db.scores.find().min({ score: 60 }).max({ score: 100 }).hint({ score: 1 })
