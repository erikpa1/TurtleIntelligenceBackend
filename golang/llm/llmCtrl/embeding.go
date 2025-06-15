package llmCtrl

import (
	"context"
	"github.com/tmc/langchaingo/llms/ollama"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"turtle/ctrlApp"
	"turtle/db"
	"turtle/lg"
	"turtle/modelsApp"
	"turtle/tools"
)

func ExampleEmbedding() {

	documents := []string{
		"The quick brown fox jumps over the lazy dog",
		"Machine learning is a subset of artificial intelligence",
		"Go is a programming language developed by Google",
		"MongoDB is a NoSQL document database",
		"Ollama allows you to run large language models locally",
		"Vector embeddings represent text as numerical vectors",
		"Cosine similarity measures the angle between two vectors",
		"Natural language processing helps computers understand human language",
	}

	embedder, err := ollama.New(
		ollama.WithModel("nomic-embed-text"), // You can change this to your preferred embedding model
	)

	lg.LogE(err)

	for _, document := range documents {
		ojbId, _ := tools.StringToObjectID(document)
		embeddings, _ := embedder.CreateEmbedding(context.Background(), []string{document})
		ctrlApp.AddDocumentEmbedding(ojbId, embeddings)

		lg.LogOk(embeddings)
	}

	// Search examples
	queries := []string{
		"lazy dog",
		"programming languages",
		"artificial intelligence",
		"database storage",
		"text processing",
	}

	docEmbeddings := db.QueryEntities[modelsApp.DocumentEmbedding](ctrlApp.CT_DOC_EMBEDDINGS, bson.M{})

	for _, query := range queries {
		lg.LogI("Going to find: ", query)

		queryEmbeding, _ := embedder.CreateEmbedding(context.Background(), []string{query})

		for _, docEmbedding := range docEmbeddings {

			db.DB.VectorSearch(context.Background(), ctrlApp.CT_DOC_EMBEDDINGS, queryEmbeding[0], 5, 0.6)

			for _, firstRow := range docEmbedding.Embedding {
				lg.LogOk(cosineSimilarity(queryEmbeding[0], firstRow))
			}
		}
	}

}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
