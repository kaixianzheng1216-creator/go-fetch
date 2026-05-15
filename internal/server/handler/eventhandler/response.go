package eventhandler

type EventOK struct {
	OK bool `json:"ok"`
}

type okResponse struct {
	Body EventOK
}

func newOKResponse() *okResponse {
	return &okResponse{Body: EventOK{OK: true}}
}
