package response

func Bearer(token string) *ResponseFetch {
	return &ResponseFetch{
		Status:      "Ok",
		StatusCode:  200,
		AccessToken: token,
	}
}

func Ok(a ...string) *ResponseFetch {
	// проверка на nil не нужна, так как поле id опционально
	return &ResponseFetch{
		Status:     "Ok",
		StatusCode: 200,
		Id:         a,
	}
}

func Error(errReq error) *ResponseFetch {
	return &ResponseFetch{
		Status:     "Error",
		StatusCode: 400,
		TextError:  errReq.Error(),
	}
}

func ServerError() *ResponseFetch {
	return &ResponseFetch{
		Status:     "Error",
		StatusCode: 500,
		TextError:  "internal server error",
	}
}
