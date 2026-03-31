// KeyVault.getKey() - Get a specific data encryption key by its UUID

// Get key by UUID
db.getMongo().getKeyVault().getKey(UUID("12345678-1234-1234-1234-123456789012"))

// Get key with string UUID
db.getMongo().getKeyVault().getKey(UUID("abcd1234-5678-90ab-cdef-123456789012"))

// Get key and convert to document
db.getMongo().getKeyVault().getKey(UUID("00000000-0000-0000-0000-000000000001"))
