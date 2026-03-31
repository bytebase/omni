// cursor.max() - Specify exclusive upper bound for index-based query

// Basic usage with index key pattern
db.users.find({ age: { $gte: 18 } }).max({ age: 65 })

// With compound index
db.users.find().max({ status: "z", name: "z" })
db.products.find({ category: "electronics" }).max({ category: "electronics", price: 1000 })

// Chained with min() for range
db.users.find().min({ age: 18 }).max({ age: 65 })
db.products.find().min({ price: 10 }).max({ price: 100 })

// With hint (required for max to work correctly)
db.users.find().max({ age: 30 }).hint({ age: 1 })
db.products.find().min({ price: 50 }).max({ price: 200 }).hint({ price: 1 })

// Chained with other cursor methods
db.users.find().max({ age: 50 }).hint({ age: 1 }).sort({ name: 1 })
db.users.find().max({ age: 30 }).hint({ age: 1 }).limit(10)

// With multiple fields
db.events.find().max({ date: ISODate("2024-12-31"), priority: 10 })
