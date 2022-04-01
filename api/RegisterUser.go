package api

import (
	"androidServer/api/apibase"
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/common/types"
	"androidServer/handler/base"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type RegisterUserRequest struct {
	base.ApiBaseRequest
}

type RegisterUserParams struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username" ` //名字
	PassWord string `json:"password"`
}

type RegisterUserResponse struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username" `
	Message  string `json:"message"`
}

func (req *RegisterUserRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "RegisterUser"
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

	errResp = req.MapAddUser(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	registerResp := &RegisterUserResponse{
		UUID:     req.RequestUUID,
		UserName: params.UserName,
		Message:  "register success",
	}
	log.Debug("session:", req.RequestUUID, ", add user success, resp:", registerResp)
	c.JSON(http.StatusOK, registerResp)

}

func (req *RegisterUserRequest) ParseParameters(c *gin.Context) (params *RegisterUserParams, resp base.ApiBaseResponse) {

	params = &RegisterUserParams{}

	params.UserName = c.PostForm("username")
	params.PassWord = c.PostForm("password")

	if params.UserName == "" || params.PassWord == "" {
		log.Error("session:", req.RequestUUID, "ParseParameters error: missing username or password")
		return nil, req.ErrorResponse("PARAMS_ERROR", "user name or password")
	}
	return params, nil
}

func (req *RegisterUserRequest) MapAddUser(params *RegisterUserParams) (errResp base.ApiBaseResponse) {
	if _, ok := apibase.NameMap[params.UserName]; ok {
		log.Error("session:", req.RequestUUID, "user already registered")
		return req.ErrorResponse("INTERNAL_ERROR", "user already registered")
	}
	apibase.NameMap[params.UserName] = params.PassWord
	return nil
}
func (req *RegisterUserRequest) DBAddUser(params *RegisterUserParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()
	userInfo := types.UserInfo{
		UUID:     req.RequestUUID,
		Username: params.UserName,
	}
	err := db.Transaction(func(db *gorm.DB) error {
		//user_table增加记录
		if err := db.Create(&userInfo).Error; err != nil {
			log.Error("session:", req.RequestUUID, ", db create userInfo err:", err, ", userInfo:", userInfo)
			return err
		}
		return nil
	})
	if err != nil {
		return req.ErrorResponse("INTERNAL_ERROR", err)
	}
	return nil
}
func RegisterUserHandler() *RegisterUserRequest {
	return &RegisterUserRequest{}
}
