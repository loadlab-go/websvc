package resp

type Response struct {
	OK   bool        `json:"ok"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OkResponse(data interface{}) Response {
	return Response{
		true,
		"success",
		data,
	}
}

func ErrResponse(err error) Response {
	return Response{
		false,
		err.Error(),
		nil,
	}
}
