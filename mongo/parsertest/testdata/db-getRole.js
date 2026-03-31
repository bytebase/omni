// db.getRole() - Returns role information

// Basic usage
db.getRole("read")
db.getRole("readWrite")
db.getRole("dbAdmin")

// With privilege information
db.getRole("customRole", { showPrivileges: true })
db.getRole("appAdmin", { showPrivileges: true })

// With built-in roles
db.getRole("read", { showBuiltinRoles: true })

// With authentication restrictions
db.getRole("restrictedRole", { showAuthenticationRestrictions: true })
