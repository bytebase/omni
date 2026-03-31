// db.collection.compactStructuredEncryptionData() - Compact encrypted data structures

// Basic compaction
db.users.compactStructuredEncryptionData()
db.orders.compactStructuredEncryptionData()
db.medicalRecords.compactStructuredEncryptionData()

// Compact with options
db.financialData.compactStructuredEncryptionData({})
db.sensitiveData.compactStructuredEncryptionData({})

// Compact encrypted collections
db.encryptedPII.compactStructuredEncryptionData()
db.patientRecords.compactStructuredEncryptionData()
db.creditCards.compactStructuredEncryptionData()

// Collection access patterns
db["users"].compactStructuredEncryptionData()
db.getCollection("users").compactStructuredEncryptionData()
db["encrypted-data"].compactStructuredEncryptionData()
db.getCollection("sensitive.records").compactStructuredEncryptionData()
