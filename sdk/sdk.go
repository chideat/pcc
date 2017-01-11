package sdk

type ResponseData struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data"`
	Info    map[string]interface{} `json:"info"`
}
