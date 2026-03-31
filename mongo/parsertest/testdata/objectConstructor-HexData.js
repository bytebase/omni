// HexData() - Create binary data from hex string with subtype

// Basic usage - subtype 0 (generic binary)
HexData(0, "48656c6c6f20576f726c6421")

// UUID subtype (4)
HexData(4, "550e8400e29b41d4a716446655440000")

// MD5 digest subtype (5)
HexData(5, "d41d8cd98f00b204e9800998ecf8427e")

// User-defined binary (128)
HexData(128, "637573746f6d2062696e617279")

// In document
db.files.insertOne({
    name: "checksum.txt",
    md5: HexData(5, "d41d8cd98f00b204e9800998ecf8427e")
})

// In query
db.files.find({
    md5: HexData(5, "d41d8cd98f00b204e9800998ecf8427e")
})

// Various subtypes with hex data
db.hexdata.insertMany([
    { type: "generic", data: HexData(0, "67656e65726963") },
    { type: "function", data: HexData(1, "66756e6374696f6e") },
    { type: "old_binary", data: HexData(2, "6f6c64") },
    { type: "uuid_old", data: HexData(3, "757569642d6f6c64") },
    { type: "uuid", data: HexData(4, "75756964") },
    { type: "md5", data: HexData(5, "6d6435") },
    { type: "encrypted", data: HexData(6, "656e63727970746564") },
    { type: "user_defined", data: HexData(128, "75736572") }
])

// Storing binary hashes
db.hashes.insertOne({
    algorithm: "sha256",
    hash: HexData(0, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
})
