package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/jwtly10/CodeEcho-Server/logger"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
)

type Handlers struct {
	C Config
	S Service
}

func NewHandlers(conf *Config, service *Service) *Handlers {
	return &Handlers{
		C: *conf,
		S: *service,
	}
}

type ChatGPTReq struct {
	Ctx []openai.ChatCompletionMessage `json:"messages"`
	Msg string                         `json:"msg"`
}

// DeepGramTranscribeHandler handles the transcription of audio
// The audio req body is expected to be base64 encoded of []bytes PCM audio

func (h *Handlers) Test(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	log.Info().Msg("incoming request")
	log.Debug().Msg("incoming request")
	log.Warn().Msg("incoming request")

	w.Write([]byte("Hello world"))
}

func (h *Handlers) DeepGramTranscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorAsJSON("Method not allowed", w, http.StatusMethodNotAllowed)
		return
	}

	audio, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorAsJSON("Unable to read body", w, http.StatusBadRequest)
		return
	}

	decodedAudio, err := base64.StdEncoding.DecodeString(string(audio))
	if err != nil {
		WriteErrorAsJSON("Unable to decode audio", w, http.StatusBadRequest)
		return
	}

	// Transcribes the audio and write the result to the response writer
	h.S.DeepGramTranscribeAudio(decodedAudio, w)
}

// ChatGPTStreamHandler handles the stream ChatGPT request
func (h *Handlers) ChatGPTStreamHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
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

	// Proxy the stream ChatGPT request and write the result to the response writer
	resChan := make(chan string)
	errChan := make(chan error)
	go h.S.ProxyStreamChatGPTReq(req.Ctx, req.Msg, resChan, errChan)
	for {
		select {
		case response, ok := <-resChan:
			if !ok {
				return
			}
			w.Write([]byte(response))
			w.(http.Flusher).Flush()
		case err := <-errChan:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Error().Err(err).Msg("Error proxying ChatGPT stream")
			return
		}
	}
}

type ErrorResp struct {
	Error string `json:"error"`
}

// WriteErrorAsJSON writes an error message to the response writer as JSON
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
