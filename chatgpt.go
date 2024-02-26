package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jwtly10/CodeEcho-Server/logger"
	"github.com/sashabaranov/go-openai"
	"io"
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
	log := logger.Get()
	personalisation := "[AI]: (Thought: I need to remember, I have the knowledge of a senior software engineer and am skilled " +
		"in multiple languages and frameworks, I will mainly help the user with their coding project, providing guidance and sharing best practises." +
		"\nThe user is also a professional. I am precise with my responses, When the user asks me to write code, I only output the code without any explanation needed." +
		"\nOnly add explanation for non-obvious things about the code. Always output production ready quality code, not just code examples. " +
		"\n[IMPORTANT] Responses should be in markdown format where possible. If you include any code blocks you MUST use markdown code sections " +
		"for this.)\n"

	log.Debug().Msg(fmt.Sprintf("Context before manipulation: %v\n", msgCtx))
	// Drop the context for anything beyond the last 3 messages
	if len(msgCtx) > 4 {
		msgCtx = msgCtx[len(msgCtx)-4:]
	}
	log.Debug().Msg(fmt.Sprintf("Context after manipulation: %v\n", msgCtx))

	msg = personalisation + msg
	log.Debug().Msg(fmt.Sprintf("Message after manipulation: %s\n", msg))

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
		log.Error().Err(err).Msg("Error creating ChatGPT stream")
		errChan <- err
		return
	}
	defer stream.Close()

	var responseBuffer bytes.Buffer

	log.Debug().Msg("ChatGPT stream created")

	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			log.Debug().Msg("ChatGPT stream closed")
			close(resChan)
			return
		}

		if err != nil {
			log.Error().Err(err).Msg("Error receiving ChatGPT stream")
			errChan <- err
			return
		}

		log.Debug().Msg(fmt.Sprintf("Stream item: %s\n", response.Choices[0].Delta.Content))
		responseBuffer.WriteString(response.Choices[0].Delta.Content)
		resChan <- response.Choices[0].Delta.Content
	}
}
