package llmCtrl

import (
	"encoding/json"
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

const CT_LLM_AGENT_TOOLS = "llm_agent_tools"
const CT_LLM_AGENTS = "llm_agents"
const CT_LLM_AGENT_TESTS = "llm_agent_tests"

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

func GetOverallAgentsPrompt(user *models.User, userQuery string) string {

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

func COULLMAgent(user *models.User, agent *llmModels.LLMAgent) {
	if user.IsAdminWithError() {
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

func GetAgent(org primitive.ObjectID, uid primitive.ObjectID) *llmModels.LLMAgent {

	return db.QueryEntity[llmModels.LLMAgent](CT_LLM_AGENTS,
		bson.M{
			"_id": uid,
			"org": org,
		})
}

func FindAgent(c *gin.Context, user *models.User, text string) llmModels.AgentTestResponse {

	models := ListAgenticOrNormalModels()

	for _, model := range models {
		response := ChatModel(c, user, model, text)

		if response != nil {
			db.InsertEntity(CT_LLM_AGENT_TESTS, response)

			return *response
		}
	}

	return llmModels.AgentTestResponse{}

}

func ChatModel(c *gin.Context, user *models.User, model *llmModels.LLM, text string) *llmModels.AgentTestResponse {

	result := llmModels.AgentTestResponse{}

	ollmodel := ollama.WithModel(model.ModelVersion)
	keepAlive := ollama.WithKeepAlive(model.Ttl)

	llm, err := ollama.New(ollmodel, keepAlive)

	if err == nil {

		prompt := GetOverallAgentsPrompt(user, text)

		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, prompt)
		result.ResultRaw = completion

		if complErr == nil {

			resultBson := bson.M{}

			serializationErr := json.Unmarshal([]byte(completion), &resultBson)

			if serializationErr == nil {

				uid, uuidOk := primitive.ObjectIDFromHex(resultBson["selected_agent"].(string))

				if uuidOk == nil {
					result.Result.SelectedAgent = uid
				}

				result.Result.Confidence = float32(resultBson["confidence"].(float64))
				result.Result.Parameters = bson.M(resultBson["parameters"].(map[string]interface{}))
				result.Result.Reasoning = resultBson["reasoning"].(string)

				result.AgentUid = result.Result.SelectedAgent

			} else {
				result.Error = serializationErr.Error()
				result.State = 0
			}

		} else {
			result.Error = complErr.Error()
			result.State = 0
			lg.LogE(completion)

		}
	} else {
		result.Error = err.Error()
		result.State = 0
		lg.LogE(err)
	}

	return nil
}

func AddToolToAgent(user *models.User, agentUid primitive.ObjectID, toolUid primitive.ObjectID) {

	db.InsertEntity(CT_LLM_AGENT_TOOLS, llmModels.LLMAgentTool{
		Agent: agentUid,
		Tool:  toolUid,
		Org:   user.Org,
	})
}

func DeleteTool(user *models.User, relationUid primitive.ObjectID) {
	db.DeleteEntity(CT_LLM_AGENT_TOOLS, bson.M{
		"_id": relationUid,
		"org": user.Org,
	})
}

func DeleteToolsOfAgent(user *models.User, agentUid primitive.ObjectID) {
	db.DeleteEntity(CT_LLM_AGENT_TOOLS, bson.M{
		"agent": agentUid,
		"org":   user.Org,
	})
}

func AskAgents(c *gin.Context, user *models.User, text string) llmModels.AgentTestResponse {

	return llmModels.AgentTestResponse{}

}

func ChatAgent(c *gin.Context, user *models.User, agentUid primitive.ObjectID, text string) llmModels.AgentTestResponse {

	result := llmModels.AgentTestResponse{}
	result.AgentUid = agentUid

	ollmodel := ollama.WithModel("mistral:7b")
	keepAlive := ollama.WithKeepAlive("10h")

	llm, err := ollama.New(ollmodel, keepAlive)

	if err == nil {

		prompt := GetOverallAgentsPrompt(user, text)

		lg.LogI("Going to ask llm")
		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, prompt)
		lg.LogOk("LLM responded")
		result.ResultRaw = completion

		if complErr == nil {

			resultBson := bson.M{}

			serializationErr := json.Unmarshal([]byte(completion), &resultBson)

			if serializationErr == nil {

				uid, uuidOk := primitive.ObjectIDFromHex(resultBson["selected_agent"].(string))

				if uuidOk == nil {
					result.Result.SelectedAgent = uid
				}

				result.Result.Confidence = float32(resultBson["confidence"].(float64))
				result.Result.Parameters = bson.M(resultBson["parameters"].(map[string]interface{}))
				result.Result.Reasoning = resultBson["reasoning"].(string)

				result.AgentUid = result.Result.SelectedAgent

				if result.Result.SelectedAgent != agentUid {
					result.State = 2
					//In this case there was an bad agent selected
				} else {
					result.State = 1
					//In this case agent was selected right
				}
			} else {
				result.Error = serializationErr.Error()
				result.State = 0
			}

		} else {
			result.Error = complErr.Error()
			result.State = 0
			lg.LogE(completion)

		}
	} else {
		result.Error = err.Error()
		result.State = 0
		lg.LogE(err)
	}

	db.InsertEntity(CT_LLM_AGENT_TESTS, &result)

	return result
}

func DeleteLLMAgent(user *models.User, uid primitive.ObjectID) {
	if user.IsAdminWithError() {

		DeleteAgentTestHistory(user, uid)
		DeleteToolsOfAgent(user, uid)

		db.DeleteEntity(CT_LLM_AGENTS, bson.M{
			"_id": uid,
			"org": user.Org,
		})
	}
}

func DeleteAgentTestHistory(user *models.User, uid primitive.ObjectID) {
	if user.IsAdminWithError() {
		db.DeleteEntity(CT_LLM_AGENT_TESTS, bson.M{
			"agentUid": uid,
		})
	}
}
