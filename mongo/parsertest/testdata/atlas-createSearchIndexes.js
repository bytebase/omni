// db.collection.createSearchIndexes() - Create multiple Atlas Search indexes

// Create multiple search indexes at once
db.movies.createSearchIndexes([
    {
        name: "default",
        definition: { mappings: { dynamic: true } }
    },
    {
        name: "title_search",
        definition: {
            mappings: {
                dynamic: false,
                fields: {
                    title: { type: "string", analyzer: "lucene.standard" }
                }
            }
        }
    }
])

// Create text and vector search indexes together
db.products.createSearchIndexes([
    {
        name: "text_search",
        definition: {
            mappings: {
                dynamic: false,
                fields: {
                    name: { type: "string" },
                    description: { type: "string" }
                }
            }
        }
    },
    {
        name: "vector_search",
        type: "vectorSearch",
        definition: {
            fields: [{
                type: "vector",
                path: "embedding",
                numDimensions: 768,
                similarity: "dotProduct"
            }]
        }
    }
])

// Create search indexes with different analyzers
db.articles.createSearchIndexes([
    {
        name: "english_search",
        definition: {
            analyzer: "lucene.english",
            mappings: { dynamic: true }
        }
    },
    {
        name: "french_search",
        definition: {
            analyzer: "lucene.french",
            mappings: { dynamic: true }
        }
    }
])

// Collection access patterns
db["movies"].createSearchIndexes([{ name: "idx1" }])
db.getCollection("movies").createSearchIndexes([{ name: "idx1" }])
