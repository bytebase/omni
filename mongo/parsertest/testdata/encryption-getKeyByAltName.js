// KeyVault.getKeyByAltName() - Get a data encryption key by its alternate name

// Get key by alternate name
db.getMongo().getKeyVault().getKeyByAltName("myKeyAlias")

// Get key by production name
db.getMongo().getKeyVault().getKeyByAltName("production-key")

// Get key by backup name
db.getMongo().getKeyVault().getKeyByAltName("backup-key")

// Get key with underscore name
db.getMongo().getKeyVault().getKeyByAltName("data_encryption_key_v1")
