package themes

import (
	"turtle/auth"
	"turtle/lg"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"

	"github.com/gin-gonic/gin"
)

func _ChatCopilot(c *gin.Context) {

	user := auth.GetUserFromContext(c)

	request := llmModels.ChatRequestParams{}
	request.SystemPrompt = llmModels.MCPSystemPrefix(`application theme.

You manage following items:

# Structures
## Theme stricture - this is theme user see
{
    "name": "",
    "uid": "",
    "topBarHeightBig": "45px",
    "bigPadding": "15px",
    "headingFontColor": "#069AF3",
    "iconPrimaryColor": "#ff0000",
    "iconSecondaryColor": "#333333",
    "borderColor": "#069AF3",
    "borderHoverColor": "#52b7f8"
}

# Commands
new_theme - creates new theme definned by user


Please format your response as JSON:

{
"confidence": 0.2,
"command": "commandType",
"data": {}
}

Example:
{
"confidence": 0.9,
"command": "new_theme",
"data": {
    "iconPrimaryColor": "#ffff00",
}
}
`)

	request.UserPrompt = "User want to create theme in colors of sunset"

	model := llmModels.LLM{}
	model.ModelVersion = "gemma3:1b"

	response := llmCtrl.ChatModelWithSystem(c, user, &model, &request)

	lg.LogOkson(response.Result.Parameters)

}

func _ListCopilotExamples(c *gin.Context) {
	c.String(200, `
	<h1>Examples</h1>
	<p>Create sunset theme</p>
`)
}
