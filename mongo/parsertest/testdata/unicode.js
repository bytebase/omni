// Unicode support - emoji characters
db.posts.find({ reaction: "👍" })
db.posts.find({ emoji: "🎉🎊🎁" })
db.posts.insertOne({ title: "Hello 🌍", content: "Testing emoji support 😀" })
db.messages.updateOne({ _id: 1 }, { $set: { status: "✅ completed" } })

// Unicode - emoji in field names (quoted)
db.reactions.find({ "👍": { $gt: 10 } })
db.reactions.insertOne({ "👍": 100, "👎": 5, "❤️": 50 })

// Unicode - Hindi (Devanagari script)
db.users.find({ name: "नमस्ते" })
db.users.find({ greeting: "आपका स्वागत है" })
db.products.insertOne({ name: "चाय", description: "भारतीय मसाला चाय" })
db.articles.find({ title: "हिंदी में लेख" })

// Unicode - Chinese (Simplified)
db.users.find({ name: "你好世界" })
db.products.find({ category: "电子产品" })
db.orders.insertOne({ item: "笔记本电脑", price: 5999 })

// Unicode - Chinese (Traditional)
db.users.find({ name: "繁體中文" })
db.products.find({ description: "傳統工藝品" })

// Unicode - Japanese (Hiragana, Katakana, Kanji)
db.users.find({ greeting: "こんにちは" })
db.products.find({ name: "カタカナ" })
db.articles.find({ title: "日本語の記事" })

// Unicode - Korean
db.users.find({ name: "안녕하세요" })
db.products.find({ category: "한국 음식" })

// Unicode - Arabic (right-to-left)
db.users.find({ name: "مرحبا بالعالم" })
db.articles.find({ title: "مقال باللغة العربية" })

// Unicode - Thai
db.users.find({ greeting: "สวัสดี" })
db.products.find({ name: "อาหารไทย" })

// Unicode - Russian (Cyrillic)
db.users.find({ name: "Привет мир" })
db.articles.find({ title: "Русский текст" })

// Unicode - Greek
db.users.find({ name: "Γειά σου κόσμε" })
db.math.find({ symbol: "π", value: 3.14159 })

// Unicode - Hebrew
db.users.find({ greeting: "שלום עולם" })

// Unicode - mixed scripts in single document
db.international.insertOne({
    english: "Hello",
    hindi: "नमस्ते",
    chinese: "你好",
    japanese: "こんにちは",
    korean: "안녕하세요",
    arabic: "مرحبا",
    russian: "Привет",
    emoji: "👋🌍"
})

// Unicode - special characters and symbols
db.math.find({ formula: "∑∏∫∂√∞" })
db.currency.find({ symbols: "€£¥₹₽₿" })
db.symbols.find({ arrows: "←→↑↓↔↕" })
db.music.find({ notes: "♪♫♬♩" })

// Unicode - combining characters and diacritics
db.users.find({ name: "José García" })
db.users.find({ name: "François Müller" })
db.users.find({ name: "Søren Kierkegaard" })
db.users.find({ name: "Nguyễn Văn A" })

// Unicode in aggregation pipelines
db.posts.aggregate([
    { $match: { reaction: "👍" } },
    { $group: { _id: "$emoji", count: { $sum: 1 } } }
])

db.international.aggregate([
    { $project: { greeting: { $concat: ["Hello ", "$hindi", " ", "$emoji"] } } }
])

// Unicode in array values
db.tags.find({ labels: { $in: ["重要", "紧急", "🔥"] } })
db.posts.insertOne({ tags: ["日本語", "한국어", "中文", "हिंदी"] })

// Unicode escape sequences (existing support)
db.users.find({ escaped: "\u0048\u0065\u006C\u006C\u006F" })
db.users.find({ mixed: "Hello \u4E16\u754C" })

// Unicode in regex patterns
db.articles.find({ title: /नमस्ते/ })
db.products.find({ name: /^你好/ })

// Unicode in index and collection names (via bracket notation)
db["文章"].find({})
db["статьи"].insertOne({ title: "Тест" })
db["記事"].aggregate([{ $match: { published: true } }])
