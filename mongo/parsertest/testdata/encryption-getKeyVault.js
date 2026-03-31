// KeyVault.getKeyVault() - Get the KeyVault object for encryption key management

// Get the KeyVault object from the current connection
db.getMongo().getKeyVault()

// Chain with other KeyVault methods
db.getMongo().getKeyVault().getKeys()

// Get KeyVault and list all keys
db.getMongo().getKeyVault().getKeys().toArray()
