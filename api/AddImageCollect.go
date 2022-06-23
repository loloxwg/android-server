package api

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/common/types"
	"androidServer/handler/base"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"net/http"
	"time"
)

type AddImageCollectRequest struct {
	base.ApiBaseRequest
	ImageID string
	Url     string
	Collect int
}

type AddImageCollectParams struct {
	Collect string `json:"collect"`
}

type AddImageCollectResponse struct {
	UUID    string `json:"uuid"`
	Url     string `json:"url"`
	Collect int    `json:"collect"`
}

func (req *AddImageCollectRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "AddImageCollect"
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

	errResp = req.DBAddImageCollect(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	AddImageCollectResp := &AddImageCollectResponse{
		UUID:    req.ImageID,
		Url:     req.Url,
		Collect: req.Collect,
	}
	log.Debug("session:", req.RequestUUID, ", add Image success, resp:", AddImageCollectResp)
	c.JSON(http.StatusOK, AddImageCollectResp)

}

func (req *AddImageCollectRequest) ParseParameters(c *gin.Context) (params *AddImageCollectParams, resp base.ApiBaseResponse) {
	params = &AddImageCollectParams{}
	if len(c.Param("id")) <= 1 {
		resp = req.ErrorResponse("PARAMS_ERROR", "id or filename is null or empty")
		log.Error("session:", req.RequestUUID, ", parse input params failed request id:", c.Param("id"))
		return nil, resp
	}
	req.ImageID = c.Param("id")[1:]

	params.Collect = c.PostForm("collect")
	return params, nil
}

func (req *AddImageCollectRequest) DBAddImageCollect(params *AddImageCollectParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()

	err := db.Model(types.ImageInfo{}).Where("uuid=?", req.ImageID).Update("collect", params.Collect).Error
	if err != nil {
		log.Error("session:", req.RequestUUID, "image id", req.ImageID)
		return req.ErrorResponse("INTERNAL_ERROR", "update collect failed")
	}
	imageInfo := &types.ImageInfo{}
	err = db.Model(types.ImageInfo{}).Where("uuid=?", req.ImageID).First(imageInfo).Error
	if err != nil {
		log.Error("session:", req.RequestUUID, "image id", req.ImageID)
		return req.ErrorResponse("INTERNAL_ERROR", "update collect failed")
	}
	req.Url = imageInfo.Url
	req.Collect = imageInfo.Collect

	return nil
}
func AddImageCollectHandler() *AddImageCollectRequest {
	return &AddImageCollectRequest{}
}
