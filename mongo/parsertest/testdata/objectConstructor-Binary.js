// Binary() - Create binary data (BSON Binary type)

// In document - basic usage with base64 string
db.files.insertOne({
    name: "document.bin",
    content: Binary("Y29udGVudA==")
})

// In query
db.files.find({
    content: Binary("Y29udGVudA==")
})

// Binary with subtype
db.files.insertOne({
    data: Binary("SGVsbG8gV29ybGQh", 0)
})

// With different subtypes
db.binaries.insertMany([
    { type: "generic", data: Binary("Z2VuZXJpYw==", 0) },
    { type: "uuid", data: Binary("dXVpZA==", 4) },
    { type: "md5", data: Binary("bWQ1", 5) }
])

// Using Binary.createFromBase64() static method
db.files.insertOne({
    fromBase64: Binary.createFromBase64("SGVsbG8=")
})

// Binary.createFromBase64() with subtype
db.files.insertOne({
    fromBase64: Binary.createFromBase64("SGVsbG8=", 0)
})

// Using Binary.createFromHexString() static method
db.files.insertOne({
    fromHex: Binary.createFromHexString("48656c6c6f")
})

// Binary.createFromHexString() with subtype
db.files.insertOne({
    fromHex: Binary.createFromHexString("48656c6c6f", 0)
})

// Using all static methods in a single document
db.files.insertOne({
    direct: Binary("Y29udGVudA=="),
    fromBase64: Binary.createFromBase64("SGVsbG8="),
    fromHex: Binary.createFromHexString("48656c6c6f")
})

// In aggregation pipeline
db.files.aggregate([
    {
        $match: {
            data: Binary("dGVzdA==")
        }
    }
])
