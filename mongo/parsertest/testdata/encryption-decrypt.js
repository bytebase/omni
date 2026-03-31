// ClientEncryption.decrypt() - Decrypt an encrypted value

// Decrypt with string representation
db.getMongo().getClientEncryption().decrypt("encryptedBinaryData")

// Decrypt placeholder encrypted data
db.getMongo().getClientEncryption().decrypt("base64EncodedEncryptedValue")

// Decrypt another encrypted value
db.getMongo().getClientEncryption().decrypt("anotherEncryptedValue")
