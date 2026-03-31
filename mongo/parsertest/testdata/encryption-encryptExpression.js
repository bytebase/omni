// ClientEncryption.encryptExpression() - Encrypt a match expression for queryable encryption

// Encrypt a simple equality expression
db.getMongo().getClientEncryption().encryptExpression(
    UUID("12345678-1234-1234-1234-123456789012"),
    { $and: [{ ssn: { $eq: "123-45-6789" } }] }
)

// Encrypt expression with range query
db.getMongo().getClientEncryption().encryptExpression(
    UUID("abcd1234-5678-90ab-cdef-123456789012"),
    { $and: [{ salary: { $gte: 50000, $lte: 100000 } }] }
)

// Encrypt expression for encrypted field search
db.getMongo().getClientEncryption().encryptExpression(
    UUID("12345678-1234-1234-1234-123456789012"),
    { $and: [{ medicalRecordNumber: "MRN-12345" }] }
)

// Encrypt complex expression
db.getMongo().getClientEncryption().encryptExpression(
    UUID("12345678-1234-1234-1234-123456789012"),
    {
        $and: [
            { ssn: { $eq: "123-45-6789" } },
            { status: "active" }
        ]
    }
)
