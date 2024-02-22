package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
)

type Service struct {
	Conf *Config
}

func NewService(config *Config) *Service {
	return &Service{Conf: config}
}

// ProxyStreamChatGPTReq handles the stream ChatGPT request
// The request body is expected to be a JSON object with the following fields:
// - messages: []openai.ChatCompletionMessage the context of the conversation
// - msg: string
func (s *Service) ProxyStreamChatGPTReq(msgCtx []openai.ChatCompletionMessage, msg string, w http.ResponseWriter) {
	msgCtx = append(msgCtx, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})

	c := openai.NewClient(s.Conf.OpenaiApiKey)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 20,
		Messages:  msgCtx,
		Stream:    true,
	}

	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		WriteErrorAsJSON("Internal server error", w, http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Print("Stream closed")
			return
		}

		if err != nil {
			log.Printf("Stream error: %v", err)
			WriteErrorAsJSON("Internal server error", w, http.StatusInternalServerError)
			return
		}

		if _, writeErr := fmt.Fprintf(w, response.Choices[0].Delta.Content); writeErr != nil {
			log.Printf("Write error: %v", writeErr)
			WriteErrorAsJSON("Internal server error", w, http.StatusInternalServerError)
			return
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
