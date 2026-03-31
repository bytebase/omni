// KeyVault.rewrapManyDataKey() - Rewrap (re-encrypt) multiple data encryption keys

// Rewrap all keys with no filter
db.getMongo().getKeyVault().rewrapManyDataKey({})

// Rewrap keys matching a filter
db.getMongo().getKeyVault().rewrapManyDataKey({ masterKey: { provider: "aws" } })

// Rewrap with new master key
db.getMongo().getKeyVault().rewrapManyDataKey({}, {
    provider: "aws",
    masterKey: {
        region: "us-east-1",
        key: "arn:aws:kms:us-east-1:123456789012:key/new-key-id"
    }
})

// Rewrap keys from one provider to another
db.getMongo().getKeyVault().rewrapManyDataKey({ "masterKey.provider": "local" }, {
    provider: "gcp",
    masterKey: {
        projectId: "my-project",
        location: "us-east1",
        keyRing: "my-keyring",
        keyName: "my-key"
    }
})

// Rewrap with Azure KMS
db.getMongo().getKeyVault().rewrapManyDataKey({ keyAltNames: "production" }, {
    provider: "azure",
    masterKey: {
        keyVaultEndpoint: "https://my-vault.vault.azure.net",
        keyName: "new-master-key"
    }
})
