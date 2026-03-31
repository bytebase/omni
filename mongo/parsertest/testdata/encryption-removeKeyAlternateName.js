// KeyVault.removeKeyAlternateName() - Remove an alternate name from a data encryption key

// Remove alternate name from key
db.getMongo().getKeyVault().removeKeyAlternateName(UUID("12345678-1234-1234-1234-123456789012"), "myKeyAlias")

// Remove alternate name
db.getMongo().getKeyVault().removeKeyAlternateName(UUID("abcd1234-5678-90ab-cdef-123456789012"), "old-key-name")

// Remove alternate name
db.getMongo().getKeyVault().removeKeyAlternateName(UUID("00000000-0000-0000-0000-000000000001"), "deprecated-alias")
