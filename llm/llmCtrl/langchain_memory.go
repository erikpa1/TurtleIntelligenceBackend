package llmCtrl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"turtle/lgr"

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
		lgr.Ok("You:")

		userInput, err := reader.ReadString('\n')

		if err != nil {
			lgr.Error(err.Error())
			return
		}

		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			lgr.Error("Empty input")
		}

		lgr.Ok("Agent:")
		lgr.Info("thinking...")

		response, err := chains.Run(
			ctx,
			conversionChat,
			userInput)

		fmt.Println("\r")

		if err != nil {
			lgr.Error(err.Error())
			return
		}

		lgr.Ok(response)
	}

}
