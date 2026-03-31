// KeyVault.getKeys() - Get all data encryption keys in the key vault

// Get all keys
db.getMongo().getKeyVault().getKeys()

// Get keys and convert to array
db.getMongo().getKeyVault().getKeys().toArray()

// Count all keys
db.getMongo().getKeyVault().getKeys().count()
