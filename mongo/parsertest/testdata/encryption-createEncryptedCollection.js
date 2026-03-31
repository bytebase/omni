// ClientEncryption.createEncryptedCollection() - Create a collection with encrypted fields

// Create encrypted collection with basic configuration
db.getMongo().getClientEncryption().createEncryptedCollection("myDatabase", "myCollection", {
    provider: "local",
    createCollectionOptions: {
        encryptedFields: {
            fields: [
                {
                    path: "ssn",
                    bsonType: "string",
                    queries: { queryType: "equality" }
                }
            ]
        }
    }
})

// Create with AWS KMS
db.getMongo().getClientEncryption().createEncryptedCollection("hr", "employees", {
    provider: "aws",
    createCollectionOptions: {
        encryptedFields: {
            fields: [
                {
                    path: "ssn",
                    bsonType: "string",
                    queries: { queryType: "equality" }
                },
                {
                    path: "salary",
                    bsonType: "int"
                }
            ]
        }
    },
    masterKey: {
        region: "us-east-1",
        key: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
    }
})

// Create with multiple encrypted fields
db.getMongo().getClientEncryption().createEncryptedCollection("medical", "patients", {
    provider: "local",
    createCollectionOptions: {
        encryptedFields: {
            fields: [
                {
                    path: "medicalRecordNumber",
                    bsonType: "string",
                    queries: { queryType: "equality" }
                },
                {
                    path: "diagnosis",
                    bsonType: "string"
                },
                {
                    path: "insuranceNumber",
                    bsonType: "string",
                    queries: { queryType: "equality" }
                }
            ]
        }
    }
})
