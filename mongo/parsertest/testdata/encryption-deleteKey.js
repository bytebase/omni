// KeyVault.deleteKey() - Delete a data encryption key from the key vault

// Delete key by UUID
db.getMongo().getKeyVault().deleteKey(UUID("12345678-1234-1234-1234-123456789012"))

// Delete key with specific UUID
db.getMongo().getKeyVault().deleteKey(UUID("abcd1234-5678-90ab-cdef-123456789012"))

// Delete key
db.getMongo().getKeyVault().deleteKey(UUID("00000000-0000-0000-0000-000000000001"))
