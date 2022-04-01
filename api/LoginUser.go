package api

import (
	"androidServer/api/apibase"
	"androidServer/app/log"
	"androidServer/handler/base"
	"fmt"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"net/http"
	"time"
)

type LoginUserRequest struct {
	base.ApiBaseRequest
}

type LoginUserParams struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username" ` //名字
	PassWord string `json:"password"`
}

type LoginUserResponse struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username" `
	Message  string `json:"message"`
}

func (req *LoginUserRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "LoginUser"
	uv4, err := uncommon.NewV4()
	if err == nil {
		raw.RequestUUID = uv4.String()
	}
	req.SetRawReqData(raw.Action, raw.RequestUUID)

	// Begin Action
	beginTime := time.Now()
	var resp base.ApiBaseResponse
	log.Info("Begin Action",
		" Method:", c.Request.Method,
		" Proto:", c.Request.Proto,
		" RemoteAddr:", c.Request.RemoteAddr,
		" Action:", raw.Action,
		" request_uuid:", raw.RequestUUID,
		" request_path:", c.Request.URL.RawPath,
		" request_query:", c.Request.URL.RawQuery)

	defer func() {
		// End Action
		endTime := time.Now()
		delay := endTime.UnixNano()/1e6 - beginTime.UnixNano()/1e6

		log.Info("End Action",
			" Action:", raw.Action,
			" request_uuid:", raw.RequestUUID,
			" response:", resp,
			" time_cost:", delay)
	}()
	params, errResp := req.ParseParameters(c)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	errResp = req.MapCheckUser(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	loginResp := &RegisterUserResponse{
		UUID:     req.RequestUUID,
		UserName: params.UserName,
		Message:  "login successful",
	}
	log.Debug("session:", req.RequestUUID, ", login user success, resp:", loginResp)
	c.JSON(http.StatusOK, loginResp)

}

func (req *LoginUserRequest) ParseParameters(c *gin.Context) (params *LoginUserParams, resp base.ApiBaseResponse) {

	params = &LoginUserParams{}

	params.UserName = c.PostForm("username")
	params.PassWord = c.PostForm("password")

	if params.UserName == "" || params.PassWord == "" {
		log.Error("session:", req.RequestUUID, "ParseParameters error: missing username or password")
		return nil, req.ErrorResponse("PARAMS_ERROR", "user name or password")
	}
	return params, nil
}

func (req *LoginUserRequest) MapCheckUser(params *LoginUserParams) (errResp base.ApiBaseResponse) {

	if apibase.NameMap[params.UserName] != params.PassWord {
		fmt.Println("password error")
		return req.ErrorResponse("INTERNAL_ERROR", "user password error")
	}
	return nil
}

func LoginUserHandler() *LoginUserRequest {
	return &LoginUserRequest{}
}
