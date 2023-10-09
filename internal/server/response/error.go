package response

type ResponseFetch struct {
	Status      string   `json:"status"`
	StatusCode  int      `json:"status_code"`
	TextError   string   `json:"error,omitempty"`
	Id          []string `json:"id,omitempty"` //возвращается при Insert или Update
	AccessToken string   `json:"access_token,omitempty"`
}

func Ok(a ...string) *ResponseFetch {
	// проверка на nil не нужна, так как поле id опционально
	return &ResponseFetch{
		Status:     "Ok",
		StatusCode: 200,
		Id:         a,
	}
}

func Bearer(token string) *ResponseFetch {
	return &ResponseFetch{
		Status:      "Ok",
		StatusCode:  200,
		AccessToken: token,
	}
}
