package api

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/handler/base"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

//该接口只对管理员开发
type ListWordRequest struct {
	base.ApiBaseRequest
	Word string //word

}

type ListWordResponse struct {
	Words []*WordInfo `json:"words"`
}

type WordInfo struct {
	Word string `json:"word"`
}

func (req *ListWordRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "ListGroups"
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
	log.Debug("session:", req.RequestUUID, ", Word:", req.Word)

	groupsListResp, errResp := req.ListGroups()
	if errResp != nil {
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	//log.Debugf("groups:%+v", groupsListResp.Groups)
	c.JSON(http.StatusOK, groupsListResp)
}

// ParseParameters 参数解析
func (req *ListWordRequest) ParseParameters(c *gin.Context, raw *base.ApiBaseRequest) (resp base.ApiBaseResponse) {
	req.Word = c.Query("name")

	return nil
}

func (req *ListWordRequest) ListGroups() (wordsListResp ListWordResponse, resp base.ApiBaseResponse) {
	var err error
	wordsListResp = ListWordResponse{
		Words: []*WordInfo{},
	}

	list := []*WordInfo{}

	DB := mysql.GetDB()
	err = DB.Table("word_table").Where("word like ?", fmt.Sprintf("%%%s%%", req.Word)).
		Find(&list).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("session:", req.RequestUUID, ", name:", req.Word, ", db found groups error:", err.Error())
		return wordsListResp, req.ErrorResponse("INTERNAL_ERROR")
	}
	log.Debugf("word_list:%+v", list)

	wordsListResp.Words = list

	return wordsListResp, nil
}
func ListWordHandler() *ListWordRequest {
	return &ListWordRequest{}
}
