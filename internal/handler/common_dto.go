package handler

type OK struct {
	OK bool `json:"ok"`
}

type OKOutput struct {
	Body OK
}

type emptyRequest struct{}

func NewOKOutput() *OKOutput {
	return &OKOutput{Body: OK{OK: true}}
}

func enumValues(values []string) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}
