// db.createUser() - Create a new user

// Basic user with role
db.createUser({
    user: "appUser",
    pwd: "password123",
    roles: ["readWrite"]
})

// User with multiple roles
db.createUser({
    user: "adminUser",
    pwd: "securePassword",
    roles: [
        { role: "readWrite", db: "app" },
        { role: "read", db: "reporting" }
    ]
})

// User with custom data
db.createUser({
    user: "developer",
    pwd: "devPassword",
    roles: ["readWrite"],
    customData: { department: "engineering", team: "backend" }
})

// User with authentication restrictions
db.createUser({
    user: "restrictedUser",
    pwd: "password",
    roles: ["read"],
    authenticationRestrictions: [
        { clientSource: ["192.168.1.0/24"], serverAddress: ["192.168.1.100"] }
    ]
})

// User with SCRAM-SHA-256 mechanism
db.createUser({
    user: "secureUser",
    pwd: "strongPassword",
    roles: ["readWrite"],
    mechanisms: ["SCRAM-SHA-256"]
})

// User with passwordDigestor
db.createUser({
    user: "clientDigestUser",
    pwd: "password",
    roles: ["read"],
    passwordDigestor: "client"
})
