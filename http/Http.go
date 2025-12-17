package http

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewBaseResponse() *BaseResponse {
	return &BaseResponse{
		Code:    0,
		Message: "",
		Data:    nil,
	}
}
