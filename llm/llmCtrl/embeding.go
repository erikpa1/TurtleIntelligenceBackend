package llmCtrl

import (
	"context"
	"github.com/tmc/langchaingo/llms/ollama"
	"turtle/lgr"
	"turtle/llm/llmModels"
)

const LLM_EMBEDING = "nomic-embed-text"

func CreateStringEmbedding(ctx context.Context, embedingString string) (llmModels.Embedding, error) {
	embedder, llmErr := ollama.New(
		ollama.WithModel(LLM_EMBEDING), // You can change this to your preferred embedding model
	)

	if llmErr != nil {
		lgr.Error(llmErr.Error())
		return [][]float32{}, llmErr
	}

	embeddings, err := embedder.CreateEmbedding(ctx, []string{embedingString})

	if err != nil {
		lgr.Error(err.Error())
	}

	return embeddings, err

}

func ExampleEmbedding() {
	//
	//for _, query := range queries {
	//	lgr.Info("Going to find: ", query)
	//
	//	queryEmbeding, _ := embedder.CreateEmbedding(context.Background(), []string{query})
	//
	//	for _, docEmbedding := range docEmbeddings {
	//
	//		db.DB.VectorSearch(context.Background(), documents.CT_DOC_EMBEDDINGS, queryEmbeding[0], 5, 0.6)
	//
	//		for _, firstRow := range docEmbedding.Embedding {
	//			lgr.Ok(cosineSimilarity(queryEmbeding[0], firstRow))
	//		}
	//	}
	//}

}
