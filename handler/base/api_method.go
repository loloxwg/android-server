package base

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
)

type ApiBaseRequest struct {
	Action      string `json:"Action"`
	RequestUUID string `json:"request_uuid"`
}

func (req *ApiBaseRequest) SetRawReqData(action string, uuid string) {
	if action != "" {
		req.Action = action
	}
	if uuid != "" {
		req.RequestUUID = uuid
	}
}

type (
	ApiBaseResponse map[string]interface{}
)

func (req *ApiBaseRequest) FieldErrorResponse(field string, tag string, params ...interface{}) ApiBaseResponse {
	newParams := make([]interface{}, 0)
	if field != "" {
		newParams = append(newParams, field)
	}
	newParams = append(newParams, params...)
	var resp ApiBaseResponse
	switch tag {
	case "required":
		resp = generateErrorCode("MISSING_PARAMETER_ERROR", newParams...)
	case "notnull":
		resp = generateErrorCode("PARAMS_ERROR", newParams...)
	case "int":
		resp = generateErrorCode("INVALID_PARAMETER_ERROR", newParams...)
	case "number":
		resp = generateErrorCode("INVALID_PARAMETER_ERROR", newParams...)
	case "string":
		resp = generateErrorCode("INVALID_PARAMETER_ERROR", newParams...)
	case "map":
		resp = generateErrorCode("INVALID_PARAMETER_ERROR", newParams...)
	case "boolean":
		resp = generateErrorCode("INVALID_PARAMETER_ERROR", newParams...)
	case "regex":
		resp = generateErrorCode("INVALID_PARAMETER_ERROR", newParams...)
	case "length":
		resp = generateErrorCode("PARAMETER_RANG_ERROR", newParams...)
	case "range":
		resp = generateErrorCode("PARAMETER_RANG_ERROR", newParams...)
	case "enumgroup":
		resp = generateErrorCode("PARAMETER_RANG_ERROR", newParams...)
	case "enumcustom":
		resp = generateErrorCode("PARAMETER_RANG_ERROR", newParams...)
	case "atLeast":
		resp = generateErrorCode("MISSING_PARAMETER_ERROR", newParams...)
	default:
		resp = generateErrorCode("RETCODE_NOT_REGIST", newParams...)
	}
	resp["Action"] = req.Action + "Response"
	return resp
}

func (req *ApiBaseRequest) ErrorResponse(name string, params ...interface{}) ApiBaseResponse {
	resp := generateErrorCode(name, params...)
	resp["Action"] = req.Action + "Response"
	return resp
}

func (req *ApiBaseRequest) OkResponse() ApiBaseResponse {
	resp := ApiBaseResponse{
		"Action":  req.Action + "Response",
		"RetCode": 0,
	}
	return resp
}

type ApiHandlerInterface interface {
	Process(c *gin.Context) ApiBaseResponse
	FieldErrorResponse(field string, tag string, params ...interface{}) ApiBaseResponse
	ErrorResponse(name string, params ...interface{}) ApiBaseResponse
	OkResponse() ApiBaseResponse
	SetRawReqData(action string, uuid string)
}

type ApiHandler func() ApiHandlerInterface

var handlerMap map[string]ApiHandler
var gl sync.Mutex

func GetHandler(action string) (ApiHandler, error) {
	gl.Lock()
	h, ok := handlerMap[action]
	gl.Unlock()
	if !ok {
		return nil, errors.New("MISSING_ACTION")
	}

	return h, nil
}

func RegisterHandler(action string, h ApiHandler) {
	registerHandler(action, h)
}

func registerHandler(action string, h ApiHandler) {
	gl.Lock()
	handlerMap[action] = h
	gl.Unlock()
}

func init() {
	handlerMap = map[string]ApiHandler{}
}
