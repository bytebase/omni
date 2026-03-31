// db.createRole() - Create a new role

// Basic role with privileges
db.createRole({
    role: "readWriteReports",
    privileges: [
        { resource: { db: "reporting", collection: "" }, actions: ["find", "insert", "update"] }
    ],
    roles: []
})

// Role inheriting from other roles
db.createRole({
    role: "appAdmin",
    privileges: [],
    roles: [
        { role: "readWrite", db: "app" },
        { role: "read", db: "logs" }
    ]
})

// Role with specific collection privileges
db.createRole({
    role: "orderManager",
    privileges: [
        { resource: { db: "sales", collection: "orders" }, actions: ["find", "insert", "update", "remove"] },
        { resource: { db: "sales", collection: "customers" }, actions: ["find"] }
    ],
    roles: []
})

// Role with write concern
db.createRole({
    role: "criticalWriter",
    privileges: [
        { resource: { db: "critical", collection: "" }, actions: ["insert", "update"] }
    ],
    roles: []
}, {
    w: "majority",
    wtimeout: 5000
})
