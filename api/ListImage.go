package api

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/handler/base"
	"errors"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type ListImageRequest struct {
	base.ApiBaseRequest
}

type ListImageResponse struct {
	Images []*ImageInfo `json:"Images"`
}

type ImageInfo struct {
	FileName string `json:"file_name"`
	Url      string `json:"url"`
	Collect  int    `json:"collect"`
}

func (req *ListImageRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "ListImages"
	uv4, err := uncommon.NewV4()
	if err == nil {
		raw.RequestUUID = uv4.String()
	}
	req.SetRawReqData(raw.Action, raw.RequestUUID)

	// Begin Action
	beginTime := time.Now()
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
			" response:", "",
			" time_cost:", delay)
	}()

	errResp := req.ParseParameters(c, raw)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	groupsListResp, errResp := req.ListImages()
	if errResp != nil {
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	//log.Debugf("groups:%+v", groupsListResp.Groups)
	c.JSON(http.StatusOK, groupsListResp)
}

// ParseParameters 参数解析
func (req *ListImageRequest) ParseParameters(c *gin.Context, raw *base.ApiBaseRequest) (resp base.ApiBaseResponse) {

	return nil
}

func (req *ListImageRequest) ListImages() (ImagesListResp ListImageResponse, resp base.ApiBaseResponse) {
	var err error
	ImagesListResp = ListImageResponse{
		Images: []*ImageInfo{},
	}

	list := []*ImageInfo{}

	DB := mysql.GetDB()
	err = DB.Table("Image_table").Find(&list).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("session:", req.RequestUUID, "db found groups error:", err.Error())
		return ImagesListResp, req.ErrorResponse("INTERNAL_ERROR")
	}
	log.Debugf("Image_list:%+v", list)

	ImagesListResp.Images = list

	return ImagesListResp, nil
}
func ListImageHandler() *ListImageRequest {
	return &ListImageRequest{}
}
