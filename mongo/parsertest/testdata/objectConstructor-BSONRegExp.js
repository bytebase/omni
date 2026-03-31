// BSONRegExp() - Create BSON regular expression

// Basic usage with pattern only
BSONRegExp("pattern")

// With flags
BSONRegExp("pattern", "i")
BSONRegExp("pattern", "im")
BSONRegExp("pattern", "imxs")

// In document
db.rules.insertOne({
    name: "email_validation",
    pattern: BSONRegExp("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", "i")
})

// In query - match documents with specific regex pattern
db.rules.find({
    pattern: BSONRegExp("pattern", "i")
})

// Complex patterns
db.validations.insertMany([
    { name: "phone", regex: BSONRegExp("^\\+?[1-9]\\d{1,14}$") },
    { name: "url", regex: BSONRegExp("^https?://", "i") },
    { name: "ipv4", regex: BSONRegExp("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$") },
    { name: "hex_color", regex: BSONRegExp("^#[0-9A-Fa-f]{6}$", "i") }
])

// Using in aggregation pipeline
db.logs.aggregate([
    {
        $match: {
            message: { $regex: BSONRegExp("error", "i") }
        }
    }
])

// Various flag combinations
db.patterns.insertMany([
    { flags: "none", pattern: BSONRegExp("test") },
    { flags: "case_insensitive", pattern: BSONRegExp("test", "i") },
    { flags: "multiline", pattern: BSONRegExp("^line", "m") },
    { flags: "dotall", pattern: BSONRegExp(".*", "s") },
    { flags: "extended", pattern: BSONRegExp("test", "x") },
    { flags: "combined", pattern: BSONRegExp("test", "imsx") }
])
