// ClientEncryption.getClientEncryption() - Get the ClientEncryption object for field-level encryption

// Get the ClientEncryption object from the current connection
db.getMongo().getClientEncryption()

// Chain with encryption methods
db.getMongo().getClientEncryption().encrypt(UUID("12345678-1234-1234-1234-123456789012"), "sensitiveData", "AEAD_AES_256_CBC_HMAC_SHA_512-Random")

// Get ClientEncryption for various operations
db.getMongo().getClientEncryption().decrypt("encryptedValue")
