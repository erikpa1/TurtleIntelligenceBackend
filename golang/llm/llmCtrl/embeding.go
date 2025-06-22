package llmCtrl

import (
	"context"
	"github.com/tmc/langchaingo/llms/ollama"
	"math"
)

func CreateStringEmbedding(ctx context.Context, embedingString string) ([][]float32, error) {
	embedder, llmErr := ollama.New(
		ollama.WithModel("nomic-embed-text"), // You can change this to your preferred embedding model
	)

	if llmErr != nil {
		return [][]float32{}, llmErr
	}

	embeddings, err := embedder.CreateEmbedding(ctx, []string{embedingString})
	return embeddings, err

}

func ExampleEmbedding() {
	//
	//for _, query := range queries {
	//	lg.LogI("Going to find: ", query)
	//
	//	queryEmbeding, _ := embedder.CreateEmbedding(context.Background(), []string{query})
	//
	//	for _, docEmbedding := range docEmbeddings {
	//
	//		db.DB.VectorSearch(context.Background(), documents.CT_DOC_EMBEDDINGS, queryEmbeding[0], 5, 0.6)
	//
	//		for _, firstRow := range docEmbedding.Embedding {
	//			lg.LogOk(cosineSimilarity(queryEmbeding[0], firstRow))
	//		}
	//	}
	//}

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
