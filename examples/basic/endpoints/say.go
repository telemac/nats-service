package endpoints

import (
	"encoding/json"
	"github.com/nats-io/nats.go/micro"
	"github.com/telemac/nats_service"
	"os/exec"
)

/*
request : nats req tts.say '{"a": 1, "b": 2}'
response : {"result":3}
*/

type SayRequest struct {
	Phrase string `json:"phrase"`
}

type SayResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Say struct {
	nats_service.Endpoint
}

func (e *Say) Name() string {
	return "say"
}

func (e *Say) Metadata() map[string]string {
	return map[string]string{
		"version": "1.0.0",
	}
}

func (e *Say) Handle(req micro.Request) {
	log := e.Service().Logger()
	log.Info("Handle request", "req", req)

	var sayRequest SayRequest
	err := json.Unmarshal(req.Data(), &sayRequest)
	if err != nil {
		log.Error("Unmarshal failed", "err", err)
		req.Error("500", err.Error(), nil)
		return
	}
	output, err := say(sayRequest.Phrase)
	if err != nil {
		log.Error("Say command failed", "err", err)
		req.Error("500", err.Error(), nil)
		return
	}
	req.RespondJSON(SayResponse{
		Output: string(output),
	})
}

func say(text string) ([]byte, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("say", text)
	output, err := cmd.CombinedOutput()
	return output, err
}
