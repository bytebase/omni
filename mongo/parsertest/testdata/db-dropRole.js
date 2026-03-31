// db.dropRole() - Drop a user-defined role

// Basic usage
db.dropRole("tempRole")
db.dropRole("readWriteReports")
db.dropRole("appAdmin")

// With write concern
db.dropRole("customRole", { w: "majority" })
db.dropRole("deprecatedRole", { w: 1, j: true })
