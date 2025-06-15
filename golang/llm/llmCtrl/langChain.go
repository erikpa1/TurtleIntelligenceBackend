package llmCtrl

import (
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"turtle/lg"
	"turtle/llm/llmModels"
)

func AskLangChainModel(c *gin.Context, model *llmModels.LLM, prompt string) string {

	ollmodel := ollama.WithModel("deepseek-coder-v2:latest")
	keepAlive := ollama.WithKeepAlive(model.Ttl)

	llm, err := ollama.New(ollmodel, keepAlive)

	if err == nil {
		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, prompt)

		if complErr == nil {
			return completion
		} else {
			lg.LogE(complErr)
			return completion
		}
	} else {
		lg.LogE(err)
	}

	return "--unanswered--"
}
