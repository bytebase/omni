// db.updateUser() - Update a user

// Update password
db.updateUser("appUser", {
    pwd: "newPassword123"
})

// Update roles
db.updateUser("appUser", {
    roles: [
        { role: "readWrite", db: "app" },
        { role: "read", db: "reporting" }
    ]
})

// Update custom data
db.updateUser("appUser", {
    customData: { department: "engineering", team: "backend", level: "senior" }
})

// Update authentication mechanisms
db.updateUser("secureUser", {
    mechanisms: ["SCRAM-SHA-256"]
})

// Update authentication restrictions
db.updateUser("restrictedUser", {
    authenticationRestrictions: [
        { clientSource: ["10.0.0.0/8"], serverAddress: ["10.0.0.1"] }
    ]
})

// Update multiple fields
db.updateUser("fullUser", {
    pwd: "newSecurePassword",
    roles: ["readWrite"],
    customData: { updated: true }
})

// With write concern
db.updateUser("criticalUser", {
    pwd: "verySecurePassword"
}, { w: "majority" })
