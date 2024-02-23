package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
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
// - resChan: chan string to send the response to the client
// - errChan: chan error to send the error to the client
func (s *Service) ProxyStreamChatGPTReq(msgCtx []openai.ChatCompletionMessage, msg string, resChan chan string, errChan chan error) {
	msgCtx = append(msgCtx, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})

	c := openai.NewClient(s.Conf.OpenaiApiKey)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 500,
		Messages:  msgCtx,
		Stream:    true,
	}

	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Println("ChatCompletionStream error: ", err)
		errChan <- err
		return
	}
	defer stream.Close()

	var responseBuffer bytes.Buffer

	log.Println("Stream started")

	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			log.Println("Stream closed")
			close(resChan)
			return
		}

		if err != nil {
			log.Println("Streaming error")
			errChan <- err
			return
		}

		responseBuffer.WriteString(response.Choices[0].Delta.Content)
		resChan <- response.Choices[0].Delta.Content
	}
}
