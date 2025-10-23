package endpoints

import (
	"encoding/json"
	"github.com/nats-io/nats.go/micro"
	"github.com/telemac/nats_service"
	"github.com/telemac/nats_service/examples/basic"
	"github.com/telemac/nats_service/pkg/counter"
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
	nats_service.Endpoint
	counter counter.Counter
}

func (e *Add) Name() string {
	return "add"
}

func (e *Add) Metadata() map[string]string {
	return map[string]string{
		"version": "1.0.0",
	}
}

func (e *Add) Handle(req micro.Request) {
	log := e.Service().Logger()
	basicService, ok := e.Service().(*basic.BasicService)
	if ok {
		basicService.Counter.Increment(1)
	}

	e.counter.Increment(1)

	var addReq AddRequest
	err := json.Unmarshal(req.Data(), &addReq)
	if err != nil {
		req.Error("500", err.Error(), nil)
		return
	}
	log.Info("add handler called", "count", e.counter.Counter())
	req.RespondJSON(AddResponse{
		Result: addReq.A + addReq.B,
	})
}
