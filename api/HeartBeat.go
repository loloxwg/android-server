package api

import (
	"androidServer/app/log"
	"androidServer/handler/base"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuidcommon "github.com/gofrs/uuid"
)

type HeartBeatRequest struct {
	base.ApiBaseRequest
}

func (req *HeartBeatRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "HeartBeat"
	uv4, err := uuidcommon.NewV4()
	if err == nil {
		raw.RequestUUID = uv4.String()
	}
	req.SetRawReqData(raw.Action, raw.RequestUUID)

	// Begin Action
	beginTime := time.Now()
	var resp base.ApiBaseResponse
	log.Info("Begin Action",
		"Method", c.Request.Method,
		"Proto", c.Request.Proto,
		"RemoteAddr", c.Request.RemoteAddr,
		"Action", raw.Action,
		"request_uuid", raw.RequestUUID,
		"request_path", c.Request.URL.RawPath,
		"request_query", c.Request.URL.RawQuery)

	// 心跳检测
	resp = req.OkResponse()

	defer func() {
		// End Action
		endTime := time.Now()
		delay := endTime.UnixNano()/1e6 - beginTime.UnixNano()/1e6

		log.Info("End Action",
			"Action:", raw.Action,
			"request_uuid", raw.RequestUUID,
			"response", resp,
			"time_cost:", delay)

	}()

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	c.JSON(http.StatusOK, resp)
}

func HeartBeatHandler() *HeartBeatRequest {
	return &HeartBeatRequest{}
}
