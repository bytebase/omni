// db.collection.createSearchIndex() - Create an Atlas Search index

// Basic search index
db.movies.createSearchIndex({
    name: "default",
    definition: { mappings: { dynamic: true } }
})

// Search index with specific mappings
db.products.createSearchIndex({
    name: "product_search",
    definition: {
        mappings: {
            dynamic: false,
            fields: {
                name: { type: "string", analyzer: "lucene.standard" },
                description: { type: "string" },
                category: { type: "token" }
            }
        }
    }
})

// Search index with analyzers
db.articles.createSearchIndex({
    name: "article_search",
    definition: {
        analyzer: "lucene.standard",
        searchAnalyzer: "lucene.standard",
        mappings: {
            dynamic: true
        }
    }
})

// Search index with synonyms
db.products.createSearchIndex({
    name: "synonym_search",
    definition: {
        mappings: { dynamic: true },
        synonyms: [{
            name: "product_synonyms",
            source: { collection: "synonyms" },
            analyzer: "lucene.standard"
        }]
    }
})

// Vector search index
db.embeddings.createSearchIndex({
    name: "vector_index",
    type: "vectorSearch",
    definition: {
        fields: [{
            type: "vector",
            path: "embedding",
            numDimensions: 1536,
            similarity: "cosine"
        }]
    }
})

// Collection access patterns
db["movies"].createSearchIndex({ name: "search" })
db.getCollection("movies").createSearchIndex({ name: "search" })
