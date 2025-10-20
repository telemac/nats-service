package endpoints

import (
	"encoding/json"
	"github.com/nats-io/nats.go/micro"
	"github.com/telemac/nats_service"
)

/*
request : nats req math.add '{"a": 1, "b": 2}'
response : {"result":3}
*/

type AddRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type AddResponse struct {
	Result int `json:"result"`
}

type Add struct {
	nats_service.NatsEndpoint
	Count int
}

func (e *Add) Name() string {
	return "add"
}

func (e *Add) Handle(req micro.Request) {
	log := e.Service().Logger()

	e.Count++

	var addReq AddRequest
	err := json.Unmarshal(req.Data(), &addReq)
	if err != nil {
		req.Error("500", err.Error(), nil)
		return
	}
	log.Info("add handler called", "count", e.Count)
	req.RespondJSON(AddResponse{
		Result: addReq.A + addReq.B,
	})
}
