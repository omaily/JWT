package response

type ResponseFetch struct {
	Status      string          `json:"status"`
	StatusCode  int             `json:"status_code"`
	TextError   string          `json:"error,omitempty"`
	Id          []string        `json:"id,omitempty"` //возвращается при Insert или Update
	AccessToken string          `json:"access_token,omitempty"`
	Valid       []ValidateError `json:"validateError,omitempty"`
}

type ValidateError struct {
	NameStruct    string `json:"name_struct"`
	Type          string `json:"type"`
	NameFieldJson string `json:"name_fieldJson"`
	ActualTag     string `json:"actual_tag"`
	Value         string `json:"value"`
	Message       string `json:"message"`
}

type InternalError struct{}

func (m *InternalError) Error() string {
	return "Internal error"
}
