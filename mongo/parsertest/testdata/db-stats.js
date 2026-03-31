// db.stats() - Get database statistics

// Basic stats
db.stats()

// Stats with scale factor (convert bytes to KB)
db.stats(1024)

// Stats with options document
db.stats({ scale: 1024 })
db.stats({ scale: 1048576 })

// Stats with freeStorage option
db.stats({ freeStorage: 1 })

// Stats with multiple options
db.stats({ scale: 1024, freeStorage: 1 })
