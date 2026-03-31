// db.setWriteConcern() - Set the default write concern

// Set write concern with w value
db.setWriteConcern({ w: 1 })
db.setWriteConcern({ w: "majority" })
db.setWriteConcern({ w: 2 })

// Set write concern with journal
db.setWriteConcern({ w: 1, j: true })
db.setWriteConcern({ w: "majority", j: true })

// Set write concern with timeout
db.setWriteConcern({ w: "majority", wtimeout: 5000 })

// Combined options
db.setWriteConcern({ w: "majority", j: true, wtimeout: 10000 })
