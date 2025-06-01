package llmApi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xyproto/ollamaclient/v2"
	"github.com/xyproto/usermodel"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/auth"
	"turtle/llm/llmCtrl"
	"turtle/tools"
)

func _GetChatHistory(c *gin.Context) {
	query := tools.QueryHeader[bson.M](c)
	tools.AutoReturn(c, llmCtrl.QueryChatHistory(query))

}

func _ChatAsk(c *gin.Context) {

	oc := ollamaclient.New(usermodel.GetTextGenerationModel())
	oc.Verbose = true
	if err := oc.PullIfNeeded(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	prompt := "Write a haiku about the color of cows."
	output, err := oc.GetOutput(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("\n%s\n", output)
}

func InitLLMChatApi(r *gin.Engine) {
	r.GET("/api/llm/chat-history", auth.LoginRequired, _GetChatHistory)
	r.POST("/api/llm/chat-ask", auth.LoginRequired, _ChatAsk)
}
