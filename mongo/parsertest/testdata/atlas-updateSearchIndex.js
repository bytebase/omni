// db.collection.updateSearchIndex() - Update an Atlas Search index

// Update search index definition
db.movies.updateSearchIndex("default", {
    definition: {
        mappings: { dynamic: true }
    }
})

// Update search index with specific field mappings
db.products.updateSearchIndex("product_search", {
    definition: {
        mappings: {
            dynamic: false,
            fields: {
                name: { type: "string", analyzer: "lucene.standard" },
                description: { type: "string", analyzer: "lucene.english" },
                category: { type: "token" },
                price: { type: "number" }
            }
        }
    }
})

// Update vector search index
db.embeddings.updateSearchIndex("vector_index", {
    definition: {
        fields: [{
            type: "vector",
            path: "embedding",
            numDimensions: 1536,
            similarity: "euclidean"
        }]
    }
})

// Update search index with new analyzer
db.articles.updateSearchIndex("article_search", {
    definition: {
        analyzer: "lucene.english",
        searchAnalyzer: "lucene.english",
        mappings: { dynamic: true }
    }
})

// Update search index with synonyms
db.products.updateSearchIndex("synonym_search", {
    definition: {
        mappings: { dynamic: true },
        synonyms: [{
            name: "updated_synonyms",
            source: { collection: "new_synonyms" },
            analyzer: "lucene.standard"
        }]
    }
})

// Collection access patterns
db["movies"].updateSearchIndex("default", { definition: { mappings: { dynamic: true } } })
db.getCollection("movies").updateSearchIndex("default", { definition: { mappings: { dynamic: true } } })
