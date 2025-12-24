package llmCtrl

import (
	"fmt"
	"turtle/core/users"
	"turtle/lg"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func DescribeDocument(c *gin.Context, user *users.User, text string) string {

	finalPrompt := fmt.Sprintf(`
SYSTEM: You are document analyst, you describe document with 200 words maximum
INSTRUCTIONS:
1. Analyze the user query and determine which agent is most appropriate
2. Extract the necessary parameters from the query
3. Respond in the following JSON format:

{
  "description": "description",
}

USER QUERY: {%s}


`, text)

	for _, model := range ListLLMModels(user.Org) {

		ollmodel := ollama.WithModel(model.ModelVersion)

		llm, err := ollama.New(ollmodel)

		if err == nil {
			completion, complErr := llms.GenerateFromSinglePrompt(c, llm, finalPrompt)

			if complErr == nil {
				lg.LogOk(completion)
			} else {
				lg.LogE(complErr)
			}
		} else {
			lg.LogE(err)
		}
	}

	return ""
}
