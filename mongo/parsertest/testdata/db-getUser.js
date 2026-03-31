// db.getUser() - Returns user information

// Basic usage
db.getUser("appUser")
db.getUser("admin")
db.getUser("reportUser")

// With additional options
db.getUser("appUser", { showCredentials: true })
db.getUser("appUser", { showPrivileges: true })
db.getUser("appUser", { showAuthenticationRestrictions: true })

// With multiple options
db.getUser("appUser", {
    showCredentials: true,
    showPrivileges: true,
    showAuthenticationRestrictions: true
})
