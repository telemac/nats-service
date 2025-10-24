package endpoints

import (
	"encoding/json"

	"github.com/nats-io/nats.go/micro"
	natsservice "github.com/telemac/natsservice"

	"os/exec"
)

/*
Say endpoint - converts text to speech using macOS say command
request : nats req tts.say '{"phrase": "Hello World"}'
response : {"output": "..."}
*/

// SayRequest represents the request payload
type SayRequest struct {
	Phrase string `json:"phrase"`
}

// SayResponse represents the response payload
type SayResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Say endpoint that uses macOS text-to-speech
type Say struct {
	natsservice.Endpoint
}

// Name returns the endpoint name
func (e *Say) Name() string {
	return "say"
}

// Metadata returns endpoint metadata
func (e *Say) Metadata() map[string]string {
	return map[string]string{
		"version": "1.0.0",
	}
}

// Handle processes text-to-speech requests
func (e *Say) Handle(req micro.Request) {
	log := e.Service().Logger()
	log.Info("Handle request", "req", req)

	// Parse request
	var sayRequest SayRequest
	err := json.Unmarshal(req.Data(), &sayRequest)
	if err != nil {
		log.Error("Unmarshal failed", "err", err)
		req.Error("500", err.Error(), nil)
		return
	}

	// Execute say command
	output, err := say(sayRequest.Phrase)
	if err != nil {
		log.Error("Say command failed", "err", err)
		req.Error("500", err.Error(), nil)
		return
	}

	// Respond with output
	req.RespondJSON(SayResponse{
		Output: string(output),
	})
}

// say executes the macOS say command
func say(text string) ([]byte, error) {
	cmd := exec.Command("say", text)
	output, err := cmd.CombinedOutput()
	return output, err
}
