package response

func Ok(a ...string) *RequestError {
	// проверка на nil не нужна, так как поле id опционально
	return &RequestError{
		Status:     "Ok",
		StatusCode: 200,
		Id:         a,
	}
}

func Error(errReq error) *RequestError {
	return &RequestError{
		Status:     "Error",
		StatusCode: 400,
		TextError:  errReq.Error(),
	}
}

func ServerError() *RequestError {
	return &RequestError{
		Status:     "Error",
		StatusCode: 500,
		TextError:  "internal server error",
	}
}
