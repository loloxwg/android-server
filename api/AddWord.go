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

type AddWordRequest struct {
	base.ApiBaseRequest
}

type AddWordParams struct {
	UUID string `json:"uuid"`
	Word string `json:"word"`
}

type AddWordResponse struct {
	UUID string `json:"uuid"`
	Word string `json:"word"`
}

func (req *AddWordRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "AddWord"
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

	addWordResp := &AddWordResponse{
		UUID: req.RequestUUID,
		Word: params.Word,
	}

	errResp = req.DBAddWord(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)

	log.Debug("session:", req.RequestUUID, ", add user success, resp:", addWordResp)
	c.JSON(http.StatusOK, addWordResp)

}

func (req *AddWordRequest) ParseParameters(c *gin.Context) (params *AddWordParams, resp base.ApiBaseResponse) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("session:", req.RequestUUID, ", Read request body failed error:", err)
		resp = req.ErrorResponse("INTERNAL_ERROR", "read body error")
		return
	}
	params = &AddWordParams{}
	if err := json.Unmarshal(bs, params); err != nil {
		log.Error("session:", req.RequestUUID, ", Unmarshal json failed error:", err, ", request:", string(bs))
		resp = req.ErrorResponse("INTERNAL_ERROR", "body format error")
		return
	}
	if params.Word == "" {
		log.Error("session:", req.RequestUUID, ", parameter is invalid:", params)
		resp = req.ErrorResponse("PARAMS_ERROR", params)
		return
	}
	return params, nil
}
func (req *AddWordRequest) DBAddWord(params *AddWordParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()
	userInfo := types.EnglishInfo{
		UUID: req.RequestUUID,
		Word: params.Word,
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
func AddWordHandler() *AddWordRequest {
	return &AddWordRequest{}
}
