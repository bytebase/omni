// KeyVault.addKeyAlternateName() - Add an alternate name to a data encryption key

// Add alternate name to key
db.getMongo().getKeyVault().addKeyAlternateName(UUID("12345678-1234-1234-1234-123456789012"), "myKeyAlias")

// Add another alternate name
db.getMongo().getKeyVault().addKeyAlternateName(UUID("abcd1234-5678-90ab-cdef-123456789012"), "production-key")

// Add alternate name with special characters
db.getMongo().getKeyVault().addKeyAlternateName(UUID("00000000-0000-0000-0000-000000000001"), "key_v2_backup")
