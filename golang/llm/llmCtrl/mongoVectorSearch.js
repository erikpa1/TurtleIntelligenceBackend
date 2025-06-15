// MongoDB Server-Side JavaScript for Vector Search
// These functions run inside MongoDB using db.eval() or as stored procedures

// =====================================
// COSINE SIMILARITY FUNCTION
// =====================================
function cosineSimilarity(vectorA, vectorB) {
    if (!vectorA || !vectorB || vectorA.length !== vectorB.length) {
        return 0;
    }

    let dotProduct = 0;
    let normA = 0;
    let normB = 0;

    for (let i = 0; i < vectorA.length; i++) {
        dotProduct += vectorA[i] * vectorB[i];
        normA += vectorA[i] * vectorA[i];
        normB += vectorB[i] * vectorB[i];
    }

    if (normA === 0 || normB === 0) {
        return 0;
    }

    return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB));
}

// =====================================
// EUCLIDEAN DISTANCE FUNCTION
// =====================================
function euclideanDistance(vectorA, vectorB) {
    if (!vectorA || !vectorB || vectorA.length !== vectorB.length) {
        return Infinity;
    }

    let sum = 0;
    for (let i = 0; i < vectorA.length; i++) {
        const diff = vectorA[i] - vectorB[i];
        sum += diff * diff;
    }

    return Math.sqrt(sum);
}

// =====================================
// VECTOR SEARCH FUNCTION
// =====================================
function vectorSearch(queryVector, limit = 10, similarityThreshold = 0.7) {
    // Get the collection
    const collection = db.documents; // Change 'documents' to your collection name

    // Initialize results array
    let results = [];

    // Iterate through all documents and calculate similarity
    collection.find({}).forEach(function (doc) {
        if (doc.embedding && Array.isArray(doc.embedding)) {
            const similarity = cosineSimilarity(queryVector, doc.embedding);

            // Only include results above threshold
            if (similarity >= similarityThreshold) {
                results.push({
                    _id: doc._id,
                    content: doc.content,
                    similarity: similarity,
                    embedding: doc.embedding
                });
            }
        }
    });

    // Sort by similarity (descending)
    results.sort(function (a, b) {
        return b.similarity - a.similarity;
    });

    // Return top results
    return results.slice(0, limit);
}

// =====================================
// AGGREGATION PIPELINE VERSION
// =====================================
function vectorSearchAggregation(queryVector, limit = 10) {
    const collection = db.documents;

    return collection.aggregate([
        {
            $addFields: {
                similarity: {
                    $function: {
                        body: function (embedding, queryVec) {
                            if (!embedding || !queryVec || embedding.length !== queryVec.length) {
                                return 0;
                            }

                            let dotProduct = 0;
                            let normA = 0;
                            let normB = 0;

                            for (let i = 0; i < embedding.length; i++) {
                                dotProduct += embedding[i] * queryVec[i];
                                normA += embedding[i] * embedding[i];
                                normB += queryVec[i] * queryVec[i];
                            }

                            if (normA === 0 || normB === 0) {
                                return 0;
                            }

                            return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB));
                        },
                        args: ["$embedding", queryVector],
                        lang: "js"
                    }
                }
            }
        },
        {
            $match: {
                similarity: {$gt: 0}
            }
        },
        {
            $sort: {
                similarity: -1
            }
        },
        {
            $limit: limit
        },
        {
            $project: {
                _id: 1,
                content: 1,
                similarity: 1
            }
        }
    ]).toArray();
}

// =====================================
// BATCH VECTOR SEARCH
// =====================================
function batchVectorSearch(queryVectors, limit = 5) {
    let allResults = [];

    queryVectors.forEach(function (queryVector, index) {
        const results = vectorSearch(queryVector, limit);
        allResults.push({
            queryIndex: index,
            results: results
        });
    });

    return allResults;
}

// =====================================
// FIND SIMILAR DOCUMENTS
// =====================================
function findSimilarDocuments(documentId, limit = 5) {
    const collection = db.documents;

    // Get the source document
    const sourceDoc = collection.findOne({_id: documentId});

    if (!sourceDoc || !sourceDoc.embedding) {
        return {error: "Document not found or has no embedding"};
    }

    // Find similar documents (excluding the source document)
    return vectorSearch(sourceDoc.embedding, limit + 1)
        .filter(function (doc) {
            return !doc._id.equals(documentId);
        })
        .slice(0, limit);
}

// =====================================
// VECTOR CLUSTERING (K-MEANS STYLE)
// =====================================
function simpleVectorClustering(k = 3, maxIterations = 10) {
    const collection = db.documents;

    // Get all embeddings
    const docs = collection.find({embedding: {$exists: true}}).toArray();

    if (docs.length < k) {
        return {error: "Not enough documents for clustering"};
    }

    // Initialize centroids randomly
    let centroids = [];
    for (let i = 0; i < k; i++) {
        const randomDoc = docs[Math.floor(Math.random() * docs.length)];
        centroids.push([...randomDoc.embedding]); // Copy the embedding
    }

    let clusters = [];

    for (let iteration = 0; iteration < maxIterations; iteration++) {
        // Assign documents to clusters
        clusters = Array(k).fill().map(() => []);

        docs.forEach(function (doc) {
            let bestCluster = 0;
            let bestSimilarity = cosineSimilarity(doc.embedding, centroids[0]);

            for (let j = 1; j < k; j++) {
                const similarity = cosineSimilarity(doc.embedding, centroids[j]);
                if (similarity > bestSimilarity) {
                    bestSimilarity = similarity;
                    bestCluster = j;
                }
            }

            clusters[bestCluster].push({
                _id: doc._id,
                content: doc.content,
                similarity: bestSimilarity
            });
        });

        // Update centroids
        for (let j = 0; j < k; j++) {
            if (clusters[j].length > 0) {
                const clusterDocs = clusters[j].map(item =>
                    docs.find(doc => doc._id.equals(item._id))
                );

                // Calculate mean of embeddings
                const embeddingLength = clusterDocs[0].embedding.length;
                const newCentroid = new Array(embeddingLength).fill(0);

                clusterDocs.forEach(function (doc) {
                    for (let i = 0; i < embeddingLength; i++) {
                        newCentroid[i] += doc.embedding[i];
                    }
                });

                for (let i = 0; i < embeddingLength; i++) {
                    newCentroid[i] /= clusterDocs.length;
                }

                centroids[j] = newCentroid;
            }
        }
    }

    return {
        clusters: clusters,
        centroids: centroids
    };
}

// =====================================
// USAGE EXAMPLES
// =====================================

// Example: Search for similar vectors
// db.eval(function() {
//     const queryVector = [0.1, 0.2, 0.3, ...]; // Your query embedding
//     return vectorSearch(queryVector, 5, 0.8);
// });

// Example: Using aggregation pipeline
// db.eval(function() {
//     const queryVector = [0.1, 0.2, 0.3, ...];
//     return vectorSearchAggregation(queryVector, 10);
// });

// Example: Find documents similar to a specific document
// db.eval(function() {
//     return findSimilarDocuments(ObjectId("your_document_id"), 5);
// });

// Example: Cluster documents
// db.eval(function() {
//     return simpleVectorClustering(3, 15);
// });

// =====================================
// UTILITY FUNCTIONS
// =====================================

// Function to add a new document with embedding
function addDocumentWithEmbedding(content, embedding) {
    const collection = db.documents;

    const result = collection.insertOne({
        content: content,
        embedding: embedding,
        createdAt: new Date(),
        metadata: {
            embeddingModel: "nomic-embed-text",
            vectorLength: embedding.length
        }
    });

    return result;
}

// Function to update embedding for existing document
function updateDocumentEmbedding(documentId, embedding) {
    const collection = db.documents;

    const result = collection.updateOne(
        {_id: documentId},
        {
            $set: {
                embedding: embedding,
                updatedAt: new Date(),
                "metadata.vectorLength": embedding.length
            }
        }
    );

    return result;
}

// Function to get vector statistics
function getVectorStats() {
    const collection = db.documents;

    return collection.aggregate([
        {
            $match: {
                embedding: {$exists: true}
            }
        },
        {
            $group: {
                _id: null,
                totalDocuments: {$sum: 1},
                avgVectorLength: {$avg: {$size: "$embedding"}},
                maxVectorLength: {$max: {$size: "$embedding"}},
                minVectorLength: {$min: {$size: "$embedding"}}
            }
        }
    ]).toArray()[0];
}

// Function to create a basic search index on content
function createTextIndex() {
    const collection = db.documents;
    return collection.createIndex({content: "text"});
}

// =====================================
// EXAMPLE SHELL COMMANDS TO RUN
// =====================================

/*
// 1. Add sample documents with embeddings
db.documents.insertMany([
    {
        content: "Machine learning is a subset of artificial intelligence",
        embedding: [0.1, 0.2, 0.3, 0.4, 0.5]
    },
    {
        content: "Deep learning uses neural networks with multiple layers",
        embedding: [0.15, 0.25, 0.35, 0.45, 0.55]
    },
    {
        content: "Natural language processing helps computers understand text",
        embedding: [0.2, 0.3, 0.4, 0.5, 0.6]
    }
]);

// 2. Search for similar documents
db.eval(function() {
    const queryVector = [0.12, 0.22, 0.32, 0.42, 0.52];
    return vectorSearch(queryVector, 3, 0.5);
});

// 3. Get vector statistics
db.eval(function() {
    return getVectorStats();
});

// 4. Find similar documents to a specific one
db.eval(function() {
    const docId = db.documents.findOne()._id;
    return findSimilarDocuments(docId, 2);
});

// 5. Cluster documents
db.eval(function() {
    return simpleVectorClustering(2, 5);
});
*/