package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/deepgram/deepgram-go-sdk/pkg/api/prerecorded/v1"
	"github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded"
	"log"
	"net/http"
	"os"
)

type TranscriptionResult struct {
	Transcript string  `json:"transcript"`
	Confidence float64 `json:"confidence"`
}

// DeepGramTranscribeAudio transcribes the given audio and writes the result to the response writer
// The audio type is expected to be PCM format, and internal logic will convert it to WAV format
// The audio is expected to be 16-bit, 2-channel, 16kHz
func (s *Service) DeepGramTranscribeAudio(audio []byte, w http.ResponseWriter) {

	tmpfile, err := os.CreateTemp("", fmt.Sprintf("audio_*.wav"))
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	err = ConvertPCMToWAV(tmpfile, audio, 16000, 16, 2)
	if err != nil {
		log.Println("ConvertPCMToWAV: ", err)
		return
	}

	ctx := context.Background()

	options := interfaces.PreRecordedTranscriptionOptions{
		Model:       "nova-2",
		SmartFormat: true,
	}

	c := client.New(s.Conf.DeepgramApiKey, interfaces.ClientOptions{})
	dg := prerecorded.New(c)

	res, err := dg.FromFile(ctx, tmpfile.Name(), options)
	if err != nil {
		fmt.Printf("FromStream failed. Err: %v\n", err)
		os.Exit(1)
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		os.Exit(1)
	}

	responseStruct, err := ParseDeepgramResponse(data)
	if err != nil {
		log.Printf("ParseDeepgramResponse failed. Err: %v\n", err)
		WriteErrorAsJSON("Error parsing deepgram response", w, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseStruct)
	if err != nil {
		log.Println("Error encoding transcript JSON")
		WriteErrorAsJSON("Error encoding transcript JSON", w, http.StatusInternalServerError)
	}
}

// ParseDeepgramResponse parses the Deepgram response and returns the transcript and confidence
func ParseDeepgramResponse(res []byte) (TranscriptionResult, error) {

	var d struct {
		Results struct {
			Channels []struct {
				Alternatives []TranscriptionResult `json:"alternatives"`
			} `json:"channels"`
		} `json:"results"`
	}

	if err := json.Unmarshal(res, &d); err != nil {
		log.Println("Error unmarshalling Deepgram response")
		return TranscriptionResult{}, err
	}

	if len(d.Results.Channels) > 0 && len(d.Results.Channels[0].Alternatives) > 0 {
		transcript := d.Results.Channels[0].Alternatives[0].Transcript
		confidence := d.Results.Channels[0].Alternatives[0].Confidence

		return TranscriptionResult{
			Transcript: transcript,
			Confidence: confidence,
		}, nil
	} else {
		log.Println("Unable to parse transcript JSON")
		return TranscriptionResult{}, fmt.Errorf("unable to parse transcript JSON")
	}
}
