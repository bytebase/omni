// KeyVault.createKey() - Create a new data encryption key

// Basic key creation with local KMS
db.getMongo().getKeyVault().createKey("local")

// Key with AWS KMS
db.getMongo().getKeyVault().createKey("aws", {
    region: "us-east-1",
    key: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
})

// Key with alternate names
db.getMongo().getKeyVault().createKey("local", "", ["myKey", "keyAlias"])

// Key with GCP KMS
db.getMongo().getKeyVault().createKey("gcp", {
    projectId: "my-project",
    location: "us-east1",
    keyRing: "my-keyring",
    keyName: "my-key"
})

// Key with Azure KMS
db.getMongo().getKeyVault().createKey("azure", {
    keyVaultEndpoint: "https://my-vault.vault.azure.net",
    keyName: "my-key"
})

// Key with KMIP provider
db.getMongo().getKeyVault().createKey("kmip")

// Key with alternate key names array
db.getMongo().getKeyVault().createKey("local", null, ["key1", "key2", "key3"])
