package main

type response struct {
	OK   bool        `json:"ok"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func okResponse(data interface{}) response {
	return response{
		true,
		"success",
		data,
	}
}

func errResponse(err error) response {
	return response{
		false,
		err.Error(),
		nil,
	}
}
