// db.collection.dropSearchIndex() - Drop an Atlas Search index

// Drop search index by name
db.movies.dropSearchIndex("default")

// Drop text search index
db.products.dropSearchIndex("product_search")

// Drop vector search index
db.embeddings.dropSearchIndex("vector_index")

// Drop search index with various names
db.articles.dropSearchIndex("english_search")
db.articles.dropSearchIndex("french_search")
db.users.dropSearchIndex("user_profile_search")

// Collection access patterns
db["movies"].dropSearchIndex("default")
db.getCollection("movies").dropSearchIndex("default")
