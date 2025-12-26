package llmCtrl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"turtle/lg"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
)

func RunMemoryTest() {

	ctx := context.Background()

	llm, _ := ollama.New(ollama.WithModel("gemma3:1b"))

	bufferMemory := memory.NewConversationBuffer()
	conversionChat := chains.NewConversation(llm, bufferMemory)

	reader := bufio.NewReader(os.Stdin)

	for {
		lg.LogOk("You:")

		userInput, err := reader.ReadString('\n')

		if err != nil {
			lg.LogE(err)
			return
		}

		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			lg.LogE("Empty input")
		}

		lg.LogOk("Agent:")
		lg.LogI("thinking...")

		response, err := chains.Run(
			ctx,
			conversionChat,
			userInput)

		fmt.Println("\r")

		if err != nil {
			lg.LogE(err)
			return
		}

		lg.LogOk(response)
	}

}
