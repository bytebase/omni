// BinData() - Create binary data with subtype

// Basic usage - subtype 0 (generic binary)
BinData(0, "SGVsbG8gV29ybGQh")

// UUID subtype (4)
BinData(4, "JDQ3MjA3M2Q0LTJhOTMtNGUxNy1hNzYyLTlkMTI0NWE5ZDRjMQ==")

// MD5 digest subtype (5)
BinData(5, "d41d8cd98f00b204e9800998ecf8427e")

// User-defined binary (128)
BinData(128, "Y3VzdG9tIGJpbmFyeQ==")

// In document
db.files.insertOne({
    name: "attachment.pdf",
    data: BinData(0, "JVBERi0xLjQKJcOkw7zDtsOf...")
})

// In query
db.files.find({
    data: BinData(0, "SGVsbG8gV29ybGQh")
})

// Various subtypes
db.binaries.insertMany([
    { type: "generic", data: BinData(0, "Z2VuZXJpYw==") },
    { type: "function", data: BinData(1, "ZnVuY3Rpb24=") },
    { type: "old_binary", data: BinData(2, "b2xk") },
    { type: "uuid_old", data: BinData(3, "dXVpZC1vbGQ=") },
    { type: "uuid", data: BinData(4, "dXVpZA==") },
    { type: "md5", data: BinData(5, "bWQ1") },
    { type: "encrypted", data: BinData(6, "ZW5jcnlwdGVk") },
    { type: "user_defined", data: BinData(128, "dXNlcg==") }
])
