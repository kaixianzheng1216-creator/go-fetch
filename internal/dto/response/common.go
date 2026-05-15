package response

type OK struct {
	OK bool `json:"ok"`
}

type OKOutput struct {
	Body OK
}

func NewOKOutput() *OKOutput {
	return &OKOutput{Body: OK{OK: true}}
}
