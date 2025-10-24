package endpoints

import (
	"encoding/json"
	"github.com/nats-io/nats.go/micro"
	"github.com/telemac/natsservice"
	"github.com/telemac/natsservice/examples/basic"
	"github.com/telemac/natsservice/pkg/counter"
)

/*
Add endpoint - performs integer addition
request : nats req math.add '{"a": 1, "b": 2}'
response : {"result":3}
*/

// AddRequest represents the request payload
type AddRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// AddResponse represents the response payload
type AddResponse struct {
	Result int `json:"result"`
}

// Add endpoint that sums two integers
type Add struct {
	natsservice.Endpoint
	counter counter.Counter
}

// Name returns the endpoint name
func (e *Add) Name() string {
	return "add"
}

// Metadata returns endpoint metadata
func (e *Add) Metadata() map[string]string {
	return map[string]string{
		"version": "1.0.0",
	}
}

// Handle processes add requests
func (e *Add) Handle(req micro.Request) {
	log := e.Service().Logger()

	// Increment service counter if available
	basicService, ok := e.Service().(*basic.BasicService)
	if ok {
		basicService.Counter.Increment(1)
	}

	// Increment endpoint counter
	e.counter.Increment(1)

	// Parse request
	var addReq AddRequest
	err := json.Unmarshal(req.Data(), &addReq)
	if err != nil {
		req.Error("500", err.Error(), nil)
		return
	}

	// Log and respond
	log.Info("add handler called", "count", e.counter.Counter())
	req.RespondJSON(AddResponse{
		Result: addReq.A + addReq.B,
	})
}
