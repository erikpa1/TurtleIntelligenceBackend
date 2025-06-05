package llmCtrl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"turtle/db"
	"turtle/lg"
	"turtle/llm/llmModels"
	"turtle/models"
	"turtle/tools"
)

const CT_LLM_AGENTS = "llm_agents"

func ExampleAgent() {
	agent_hierarchy := bson.M{
		"coordinator": bson.M{
			"name":          "main_coordinator",
			"role":          "Task Coordinator",
			"system_prompt": "You coordinate tasks between specialized agents and ensure workflow completion.",
			"subordinates":  bson.A{"research_agent", "analysis_agent", "writer_agent"},
		},
		"specialists": bson.A{
			bson.M{
				"name":           "research_agent",
				"specialization": "information_gathering",
				"capabilities":   bson.A{"web_search", "data_extraction"},
			},
		},
	}

	lg.LogI(agent_hierarchy)
}

func GetSuitableAgentPrompt(user *models.User, userQuery string) string {

	//Mistral DOC https://ollama.com/library/mistral

	finalPrompt := `
SYSTEM: You are an AI assistant that can route user queries to specialized agents. Based on the user's question, select the most appropriate agent and provide the necessary parameters.
AVAILABLE AGENTS:
`
	agents := ListLLMAgents(user)

	for i, agent := range agents {
		//Claude finnish
		if user.Type < agent.UserLevel {
			continue
		}

		finalPrompt += fmt.Sprintf(`%d. **%s**`, i+1, agent.Uid.Hex())
		finalPrompt += fmt.Sprintf(`- Description: %s`, agent.Description)

		//1. **stock_price_agent**
		//   - Description: Retrieves current and historical stock prices, market data, and financial metrics
		//   - Required parameters:
		//     - symbol (string): Stock ticker symbol (e.g., "AAPL", "TSLA")
		//     - timeframe (string): "current", "1d", "1w", "1m", "1y"
		//   - Optional parameters:
		//     - include_metrics (boolean): Include P/E ratio, market cap, etc.
		//

		if len(agent.AgentProps.RequiredParameters) > 0 {
			finalPrompt += fmt.Sprintf(`- Required parameters: %s`, strings.Join(agent.AgentProps.RequiredParameters, ","))
		}

		if len(agent.AgentProps.OptionalParameters) > 0 {
			finalPrompt += fmt.Sprintf(`- Optional parameters: %s`, strings.Join(agent.AgentProps.RequiredParameters, ","))
		}

		finalPrompt += "\n"
	}

	finalPrompt += fmt.Sprintf(`
USER QUERY: {%s}

INSTRUCTIONS:
1. Analyze the user query and determine which agent is most appropriate
2. Extract the necessary parameters from the query
3. Respond in the following JSON format:

{
  "selected_agent": "agent_name",
  "confidence": 0.95,
  "parameters": {
    "param1": "value1",
    "param2": "value2"
  },
  "reasoning": "Brief explanation of why this agent was chosen"
}
`, userQuery)

	return finalPrompt
}

func ListLLMAgents(user *models.User) []*llmModels.LLMAgent {
	return db.QueryEntities[llmModels.LLMAgent](CT_LLM_AGENTS, bson.M{
		"org": user.Org,
	})
}

func DeleteLLMAgent(user *models.User, uid primitive.ObjectID) {
	if user.IsAdmin() {
		db.DeleteEntity(CT_LLM_AGENTS, bson.M{
			"_id": uid,
			"org": user.Org,
		})
	}
}

func COULLMAgent(user *models.User, agent *llmModels.LLMAgent) {
	if user.IsAdmin() {
		agent.Org = user.Org
		agent.UpdatedAt = tools.GetTimeNowMillis()

		if agent.Uid.IsZero() {
			agent.CreatedAt = tools.GetTimeNowMillis()
			db.InsertEntity(CT_LLM_AGENTS, agent)
		} else {
			db.UpdateOneCustom(CT_LLM_AGENTS, bson.M{
				"_id": agent.Uid,
				"org": user.Org,
			},
				bson.M{
					"$set": agent,
				},
			)

		}
	}
}

func buildMistralPrompt(agent *llmModels.LLMAgent, userText string) string {
	prompt := fmt.Sprintf(`[INST] %s
%s
User request: %s`, agent.AgentProps.Role, agent.AgentProps.SystemPrompt, userText)

	if agent.AgentProps.AnswerFormat != "" {
		prompt += fmt.Sprintf("\n\nPlease format your response as: %s", agent.AgentProps.AnswerFormat)
	}
	prompt += " [/INST]"

	return prompt
}

func TestLLMAgent(c *gin.Context, user *models.User, text string) {

	ollmodel := ollama.WithModel("mistral:7b")
	keepAlive := ollama.WithKeepAlive("10h")

	llm, err := ollama.New(ollmodel, keepAlive)

	if err == nil {

		prompt := GetSuitableAgentPrompt(user, text)

		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, prompt)
		if complErr == nil {
			lg.LogI(completion)
		} else {
			lg.LogE(completion)
		}
	} else {
		lg.LogE(err)
	}

}
