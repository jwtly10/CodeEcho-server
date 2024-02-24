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
	personalisation := "[AI]: (Thought: I need to remember, I have the knowledge of a senior software engineer and am skilled " +
		"in multiple languages and frameworks, I help the user with their coding project, provide guidance and share best practises." +
		"\nThe user is also a professional. When the user asks me to write code, I only output the code without any explanation needed. " +
		"\nOnly add explanation for non-obvious things about the code. Always output production ready quality code, not code examples. " +
		")\n"

	msg = personalisation + msg

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

		log.Printf("DEBUG: Stream item: %s\n", response.Choices[0].Delta.Content)
		responseBuffer.WriteString(response.Choices[0].Delta.Content)
		resChan <- response.Choices[0].Delta.Content
	}
}
