// db.getUsers() - Returns all users in the database

// Basic usage
db.getUsers()

// With options
db.getUsers({ showCredentials: true })
db.getUsers({ showPrivileges: true })

// With filter
db.getUsers({ filter: { roles: { $elemMatch: { role: "readWrite" } } } })

// With multiple options
db.getUsers({
    showCredentials: false,
    showPrivileges: true,
    filter: { "customData.department": "engineering" }
})
