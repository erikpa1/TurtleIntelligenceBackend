// https://claude.ai/chat/723a5c8f-4824-4c40-b32d-66b57573a269 create from
package llmCtrl

import (
	"context"
	"fmt"
	"github.com/erikpa1/turtle/llm/llmModels"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"net/http"
)

func AskLangChainModel(c *gin.Context, model *llmModels.LLM, prompt string) string {

	ollmodel := ollama.WithModel(model.ModelVersion)

	if model.Ttl == "-1" {
		model.Ttl = "-1s"
	}

	keepAlive := ollama.WithKeepAlive(model.Ttl)

	llm, err := ollama.New(ollmodel, keepAlive)

	if err == nil {
		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, prompt)

		if complErr == nil {
			return completion
		} else {
			return complErr.Error()
		}
	} else {
		return err.Error()
	}

	return "--unanswered--"
}

// StreamLangChainModel streams the LangChain response to frontend using SSE
func AskLangChainModelStream(c *gin.Context, model *llmModels.LLM, prompt string) {
	// Set headers for Server-Sent Events
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Get flusher to ensure data is sent immediately
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	ollmodel := ollama.WithModel(model.ModelVersion)

	if model.Ttl == "-1" {
		model.Ttl = "-1s"
	}

	keepAlive := ollama.WithKeepAlive(model.Ttl)

	llm, err := ollama.New(ollmodel, keepAlive)
	if err != nil {
		// Send error event
		fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", fmt.Sprintf("Failed to initialize LLM: %v", err))
		flusher.Flush()
		return
	}

	// Create a context that can be cancelled if client disconnects
	ctx := c.Request.Context()

	// Send start event
	fmt.Fprintf(c.Writer, "event: start\ndata: Starting generation...\n\n")
	flusher.Flush()

	// Use streaming generate
	_, err = llms.GenerateFromSinglePrompt(
		ctx,
		llm,
		prompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			// Send each chunk as SSE
			fmt.Fprintf(c.Writer, "event: token\ndata: %s\n\n", string(chunk))
			flusher.Flush()

			// Check if client disconnected
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		}),
	)

	if err != nil {
		fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", fmt.Sprintf("Generation error: %v", err))
	} else {
		fmt.Fprintf(c.Writer, "event: complete\ndata: Generation completed\n\n")
	}

	flusher.Flush()
}
