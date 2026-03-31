// cursor.collation() - Specify collation for string comparison

// Basic collation with locale
db.users.find().collation({ locale: "en" })
db.users.find().collation({ locale: "fr" })
db.users.find().collation({ locale: "de" })

// Case-insensitive collation
db.users.find({ name: "alice" }).collation({ locale: "en", strength: 2 })

// Numeric ordering for strings
db.products.find().collation({ locale: "en", numericOrdering: true })

// With sort for locale-aware sorting
db.users.find().sort({ name: 1 }).collation({ locale: "en" })
db.users.find().sort({ name: 1 }).collation({ locale: "de", caseLevel: true })

// Full collation options
db.users.find().collation({
    locale: "en",
    strength: 2,
    caseLevel: false,
    caseFirst: "off",
    numericOrdering: false,
    alternate: "non-ignorable",
    maxVariable: "punct",
    backwards: false
})

// Chained with other cursor methods
db.users.find({ status: "active" }).collation({ locale: "en" }).sort({ name: 1 }).limit(10)
