package utils

import (
	"github.com/gin-gonic/gin"
)

type JSON struct {
	context    *gin.Context
	statusCode int
	data       map[string]interface{}
}

func NewJson(c *gin.Context, code string, message string) *JSON {
	json := JSON{}
	json.context = c
	json.statusCode = 200

	json.data = map[string]interface{}{}
	json.data["code"] = code
	json.data["message"] = message

	return &json
}

func (json *JSON) Info(info interface{}) *JSON {
	json.data["info"] = info

	return json
}

func (json *JSON) Data(data interface{}) *JSON {
	json.data["data"] = data

	return json
}

func (json *JSON) StatusCode(code int) *JSON {
	json.statusCode = code

	return json
}

func (json *JSON) End() {
	json.context.JSON(json.statusCode, json.data)
}

func Json(c *gin.Context, code, message string) {
	NewJson(c, code, message).End()
}

func JsonWithData(c *gin.Context, code, message string, data interface{}) {
	NewJson(c, code, message).Data(data).End()
}

func JsonWithDataInfo(c *gin.Context, code, message string, data, info interface{}) {
	NewJson(c, code, message).Data(data).Info(info).End()
}
