// Unquoted keys
db.users.find({ name: "alice", age: 25 })
db.users.find({ firstName: "Alice", lastName: "Smith" })

// Quoted keys
db.users.find({ "name": "alice", "age": 25 })
db.users.find({ 'name': 'alice', 'age': 25 })

// Mixed quoted and unquoted keys
db.users.find({ name: "alice", "special-field": "value" })

// Keys with $ prefix (operators)
db.users.find({ age: { $gt: 25, $lt: 65 } })
db.users.find({ $and: [{ status: "active" }, { age: { $gte: 18 } }] })

// Nested documents
db.users.find({
    profile: {
        name: "test",
        settings: {
            theme: "dark",
            notifications: true
        }
    }
})

// Arrays
db.users.find({ tags: ["mongodb", "database", "nosql"] })
db.users.find({ scores: [85, 90, 78, 92] })

// Nested arrays
db.users.find({
    matrix: [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9]
    ]
})

// Trailing commas (allowed in mongosh)
db.users.find({ name: "alice", age: 25, })
db.users.find({ tags: ["a", "b", "c",] })
db.users.find({
    name: "alice",
    age: 25,
    active: true,
})
