package api

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/common/types"
	"androidServer/handler/base"
	"encoding/json"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

type AddUserRequest struct {
	base.ApiBaseRequest
}

type AddUserParams struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username" ` //名字
	Sex      string `json:"sex" `      //性别
	Nation   string `json:"nation" `   //民族
	Eid      string `json:"eid" `      //民族
	Address  string `json:"address"`   //地址
	Birthday string `json:"birthday"`  //生日
}

type AddUserResponse struct {
	UUID     string `json:"uuid"`
	UserName string `json:"username" ` //名字
	Sex      string `json:"sex" `      //性别
	Nation   string `json:"nation" `   //民族
	Eid      string `json:"eid" `      //民族
	Address  string `json:"address"`   //地址
	Birthday string `json:"birthday"`  //生日
}

func (req *AddUserRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "AddUser"
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

	addUserResp := &AddUserResponse{
		UUID:     req.RequestUUID,
		UserName: params.UserName,
		Sex:      params.Sex,
		Nation:   params.Nation,
		Eid:      params.Eid,
		Address:  params.Address,
		Birthday: params.Birthday,
	}

	errResp = req.DBAddUser(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)

	log.Debug("session:", req.RequestUUID, ", add user success, resp:", addUserResp)
	c.JSON(http.StatusOK, addUserResp)

}

func (req *AddUserRequest) ParseParameters(c *gin.Context) (params *AddUserParams, resp base.ApiBaseResponse) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("session:", req.RequestUUID, ", Read request body failed error:", err)
		resp = req.ErrorResponse("INTERNAL_ERROR", "read body error")
		return
	}
	params = &AddUserParams{}
	if err := json.Unmarshal(bs, params); err != nil {
		log.Error("session:", req.RequestUUID, ", Unmarshal json failed error:", err, ", request:", string(bs))
		resp = req.ErrorResponse("INTERNAL_ERROR", "body format error")
		return
	}
	if params.UserName == "" || params.Eid == "" || params.Address == "" || params.Nation == "" || params.Sex == "" {
		log.Error("session:", req.RequestUUID, ", parameter is invalid:", params)
		resp = req.ErrorResponse("PARAMS_ERROR", params)
		return
	}
	params.Birthday = params.Eid[6:14]
	return params, nil
}
func (req *AddUserRequest) DBAddUser(params *AddUserParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()
	userInfo := types.UserInfo{
		UUID:     req.RequestUUID,
		Username: params.UserName,
		Sex:      params.Sex,
		Nation:   params.Nation,
		Eid:      params.Eid,
		Address:  params.Address,
		Birthday: params.Birthday,
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
func AddUserHandler() *AddUserRequest {
	return &AddUserRequest{}
}
