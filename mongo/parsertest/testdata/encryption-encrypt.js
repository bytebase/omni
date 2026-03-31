// ClientEncryption.encrypt() - Encrypt a value using client-side field level encryption

// Encrypt a string value with deterministic encryption
db.getMongo().getClientEncryption().encrypt(
    UUID("12345678-1234-1234-1234-123456789012"),
    "sensitive-data",
    "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic"
)

// Encrypt with random encryption (default)
db.getMongo().getClientEncryption().encrypt(
    UUID("abcd1234-5678-90ab-cdef-123456789012"),
    "secret-value",
    "AEAD_AES_256_CBC_HMAC_SHA_512-Random"
)

// Encrypt a number
db.getMongo().getClientEncryption().encrypt(
    UUID("12345678-1234-1234-1234-123456789012"),
    123456789,
    "AEAD_AES_256_CBC_HMAC_SHA_512-Random"
)

// Encrypt with options document
db.getMongo().getClientEncryption().encrypt(
    UUID("12345678-1234-1234-1234-123456789012"),
    { ssn: "123-45-6789" },
    "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic"
)

// Encrypt using key alternate name
db.getMongo().getClientEncryption().encrypt(
    "myDataKey",
    "confidential",
    "AEAD_AES_256_CBC_HMAC_SHA_512-Random"
)
