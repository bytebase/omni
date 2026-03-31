// Regex literals
db.users.find({ name: /alice/ })
db.users.find({ name: /^alice/i })
db.users.find({ email: /.*@example\.com$/ })
db.users.find({ username: /user\d+/gm })

// RegExp constructor
db.users.find({ name: RegExp("alice") })
db.users.find({ name: RegExp("^alice", "i") })
db.users.find({ name: RegExp("test", "gi") })

// Regex with query operators
db.users.find({ name: { $regex: /^test/i } })
db.logs.find({ message: { $regex: /error/i, $not: /warning/ } })
