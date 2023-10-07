package response

type RequestError struct {
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Err        error           `json:"err,omitempty"`
	TextError  string          `json:"error,omitempty"`
	Id         []string        `json:"id,omitempty"` //возвращается при Insert или Update
	Valid      []ValidateError `json:"validateError,omitempty"`
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

type ValidateError struct {
	NameStruct    string `json:"nameStruct"`
	Type          string `json:"type"`
	NameFieldJson string `json:"nameFieldJson"`
	ActualTag     string `json:"actualTag"`
	Value         string `json:"value"`
	Message       string `json:"message"`
}

type InternalError struct{}

func (m *InternalError) Error() string {
	return "Internal error"
}
