// https://claude.ai/chat/723a5c8f-4824-4c40-b32d-66b57573a269 create from
package llmCtrl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"io"
	"net/http"
	"strings"
	"time"
	"turtle/lg"
	"turtle/llm/llmModels"
)

type StreamConfig struct {
	ChunkSize     int           // Number of tokens to buffer before sending
	FlushInterval time.Duration // Maximum time to wait before flushing buffer
	Timeout       time.Duration // Overall timeout for the request
}

// DefaultStreamConfig returns sensible defaults
func DefaultStreamConfig() StreamConfig {
	return StreamConfig{
		ChunkSize:     10,                     // Buffer 10 tokens
		FlushInterval: 100 * time.Millisecond, // Flush every 100ms max
		Timeout:       5 * time.Minute,        // 5 minute timeout
	}
}

func AskLangChainModel(c *gin.Context, model *llmModels.LLM, prompt string) string {

	ollmodel := ollama.WithModel("deepseek-coder-v2:latest")

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
			lg.LogE(complErr)
			return completion
		}
	} else {
		lg.LogE(err)
	}

	return "--unanswered--"
}

// StreamLangChainModel streams the LangChain response to frontend using SSE
func AskLangChainModelStream(c *gin.Context, model *llmModels.LLM, prompt string) {
	config := DefaultStreamConfig()

	// Set headers for Server-Sent Events BEFORE any writes
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Write headers immediately
	c.Writer.WriteHeader(http.StatusOK)

	// Get flusher to ensure data is sent immediately
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		fmt.Fprintf(c.Writer, "event: error\ndata: Streaming not supported\n\n")
		return
	}

	ollmodel := ollama.WithModel("deepseek-coder-v2:latest")

	if model.Ttl == "-1" {
		model.Ttl = "-1s"
	}

	keepAlive := ollama.WithKeepAlive(model.Ttl)

	llm, err := ollama.New(ollmodel, keepAlive)
	if err != nil {
		sendSSEError(c.Writer, flusher, fmt.Sprintf("Failed to initialize LLM: %v", err))
		return
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
	defer cancel()

	// Send start event with config info
	startMsg := fmt.Sprintf("Starting generation (chunk_size: %d, flush_interval: %dms)",
		config.ChunkSize, config.FlushInterval.Milliseconds())
	sendSSEEvent(c.Writer, flusher, "start", startMsg)

	// Create token buffer and chunking logic
	tokenBuffer := make([]string, 0, config.ChunkSize)
	lastFlush := time.Now()

	// Ticker for periodic flushing
	flushTicker := time.NewTicker(config.FlushInterval)
	defer flushTicker.Stop()

	// Channel to receive tokens from streaming
	tokenChan := make(chan string, 100)
	errorChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	// Start LLM generation in goroutine
	go func() {
		defer close(tokenChan)

		lg.LogOk("Going to stream")

		_, err := llms.GenerateFromSinglePrompt(
			ctx,
			llm,
			prompt,
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				token := string(chunk)
				if token != "" {
					select {
					case tokenChan <- token:
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				return nil
			}),
		)

		if err != nil {
			errorChan <- err
		} else {
			doneChan <- true
		}
	}()

	// Function to flush current buffer
	flushBuffer := func() {
		if len(tokenBuffer) > 0 {
			chunk := strings.Join(tokenBuffer, "")
			sendSSEEvent(c.Writer, flusher, "chunk", chunk)
			tokenBuffer = tokenBuffer[:0] // Clear buffer
			lastFlush = time.Now()

			lg.LogE("Flushing")
		}
	}

	// Main streaming loop
	for {
		select {
		case token, ok := <-tokenChan:
			if !ok {
				// No more tokens, flush remaining buffer and complete
				flushBuffer()
				sendSSEEvent(c.Writer, flusher, "complete", "Generation completed")
				return
			}

			// Add token to buffer
			tokenBuffer = append(tokenBuffer, token)

			// Check if buffer is full
			if len(tokenBuffer) >= config.ChunkSize {
				flushBuffer()
			}

		case <-flushTicker.C:
			// Periodic flush if buffer has content and enough time has passed
			if len(tokenBuffer) > 0 && time.Since(lastFlush) >= config.FlushInterval {
				flushBuffer()
			}

		case err := <-errorChan:
			// Flush any remaining tokens before sending error
			flushBuffer()
			sendSSEError(c.Writer, flusher, fmt.Sprintf("Generation error: %v", err))
			return

		case <-doneChan:
			// Generation completed successfully
			flushBuffer()
			sendSSEEvent(c.Writer, flusher, "complete", "Generation completed")
			return

		case <-ctx.Done():
			// Timeout or client disconnected
			flushBuffer()
			sendSSEError(c.Writer, flusher, "Request timeout or connection closed")
			return
		}
	}
}

// Helper function to send SSE events
func sendSSEEvent(w io.Writer, flusher http.Flusher, event, data string) {
	// Escape newlines in data
	data = strings.ReplaceAll(data, "\n", "\\n")
	data = strings.ReplaceAll(data, "\r", "\\r")

	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
	flusher.Flush()
}

// Helper function to send SSE errors
func sendSSEError(w io.Writer, flusher http.Flusher, errorMsg string) {
	sendSSEEvent(w, flusher, "error", errorMsg)
}

// Alternative streaming function with word-based chunking
func StreamLangChainModelWordChunks(c *gin.Context, model *llmModels.LLM, prompt string, wordCount int) {
	config := StreamConfig{
		ChunkSize:     wordCount,
		FlushInterval: 150 * time.Millisecond,
		Timeout:       5 * time.Minute,
	}

	// Set headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		sendSSEError(c.Writer, flusher, "Streaming not supported")
		return
	}

	ollmodel := ollama.WithModel("deepseek-coder-v2:latest")
	if model.Ttl == "-1" {
		model.Ttl = "-1s"
	}
	keepAlive := ollama.WithKeepAlive(model.Ttl)

	llm, err := ollama.New(ollmodel, keepAlive)
	if err != nil {
		sendSSEError(c.Writer, flusher, fmt.Sprintf("Failed to initialize LLM: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
	defer cancel()

	sendSSEEvent(c.Writer, flusher, "start", fmt.Sprintf("Starting generation (word chunks: %d)", wordCount))

	var wordBuffer []string
	var currentText strings.Builder

	_, err = llms.GenerateFromSinglePrompt(
		ctx,
		llm,
		prompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			token := string(chunk)
			currentText.WriteString(token)

			// Split by spaces to get words
			if strings.Contains(token, " ") || strings.Contains(token, "\n") {
				words := strings.Fields(currentText.String())

				if len(words) > 0 {
					wordBuffer = append(wordBuffer, words[:len(words)-1]...)

					// Keep the last incomplete word
					if len(words) > 0 {
						currentText.Reset()
						currentText.WriteString(words[len(words)-1])
					}

					// Send chunk if buffer is full
					if len(wordBuffer) >= wordCount {
						chunk := strings.Join(wordBuffer[:wordCount], " ") + " "
						sendSSEEvent(c.Writer, flusher, "chunk", chunk)
						wordBuffer = wordBuffer[wordCount:]
					}
				}
			}

			return nil
		}),
	)

	// Send remaining words
	if len(wordBuffer) > 0 {
		chunk := strings.Join(wordBuffer, " ")
		sendSSEEvent(c.Writer, flusher, "chunk", chunk)
	}

	// Send any remaining text
	if currentText.Len() > 0 {
		sendSSEEvent(c.Writer, flusher, "chunk", currentText.String())
	}

	if err != nil {
		sendSSEError(c.Writer, flusher, fmt.Sprintf("Generation error: %v", err))
	} else {
		sendSSEEvent(c.Writer, flusher, "complete", "Generation completed")
	}
}
