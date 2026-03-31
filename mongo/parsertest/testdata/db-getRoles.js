// db.getRoles() - Returns all roles in the database

// Basic usage
db.getRoles()

// Include built-in roles
db.getRoles({ showBuiltinRoles: true })

// Include privileges
db.getRoles({ showPrivileges: true })

// Include both
db.getRoles({ showBuiltinRoles: true, showPrivileges: true })

// Include authentication restrictions
db.getRoles({ showAuthenticationRestrictions: true })
