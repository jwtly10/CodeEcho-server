package main

import (
	"encoding/json"
	"github.com/sashabaranov/go-openai"
	"net/http"
)

type Handlers struct {
	Conf Config
}

func NewHandlers(config *Config) *Handlers {
	return &Handlers{
		Conf: *config,
	}
}

type ChatGPTReq struct {
	Ctx []openai.ChatCompletionMessage `json:"messages"`
	Msg string                         `json:"msg"`
}

func (h *Handlers) ChatGPTStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorAsJSON("Method not allowed", w, http.StatusMethodNotAllowed)
		return
	}

	var req ChatGPTReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorAsJSON("Invalid request, unable to decode body", w, http.StatusBadRequest)
		return
	}

	ProxyStreamChatGPTReq(h.Conf.OpenaiApiKey, req.Ctx, req.Msg, w)
}

type ErrorResp struct {
	Error string `json:"error"`
}

func WriteErrorAsJSON(msg string, w http.ResponseWriter, code int) {
	errRes := ErrorResp{
		Error: msg,
	}

	jsonResp, err := json.Marshal(errRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResp)
}
