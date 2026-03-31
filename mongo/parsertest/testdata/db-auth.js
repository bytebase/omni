// db.auth() - Authenticate a user to the database

// Basic authentication with username and password
db.auth("username", "password")

// Authentication with document
db.auth({
    user: "username",
    pwd: "password"
})

// Authentication with mechanism
db.auth({
    user: "username",
    pwd: "password",
    mechanism: "SCRAM-SHA-256"
})

// Authentication with digest password disabled
db.auth({
    user: "username",
    pwd: "password",
    digestPassword: false
})
