// String literals - double quoted
db.users.find({ name: "alice" })
db.users.find({ message: "Hello, World!" })
db.users.find({ path: "C:\\Users\\test" })
db.users.find({ unicode: "\u0048\u0065\u006C\u006C\u006F" })

// String literals - single quoted
db.users.find({ name: 'alice' })
db.users.find({ message: 'Hello, World!' })
db.users.find({ escaped: 'It\'s a test' })

// Numbers - integers
db.users.find({ age: 25 })
db.users.find({ count: 0 })
db.users.find({ score: -10 })
db.users.find({ bigNum: 1000000 })

// Numbers - floats
db.users.find({ price: 19.99 })
db.users.find({ temperature: -40.5 })
db.users.find({ ratio: 0.5 })
db.users.find({ tiny: .001 })

// Numbers - scientific notation
db.users.find({ distance: 1.5e10 })
db.users.find({ small: 1e-6 })
db.users.find({ precise: 6.022e23 })
db.users.find({ negative: -1.23e-4 })

// Booleans
db.users.find({ active: true })
db.users.find({ deleted: false })
db.users.find({ $and: [{ active: true }, { verified: true }] })

// Null
db.users.find({ deletedAt: null })
db.users.find({ middleName: null })
