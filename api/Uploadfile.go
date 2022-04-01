package api

import (
	"androidServer/app/log"
	"androidServer/handler/base"
	"github.com/gin-gonic/gin"
	uuidcommon "github.com/gofrs/uuid"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

type UploadFileRequest struct {
	base.ApiBaseRequest
	FileContent []byte
	Path        string
	Mode        fs.FileMode
}

type UploadFileResponse struct {
	Mode fs.FileMode `json:"mode"`
	Path string      `json:"path"`
}

func (req *UploadFileRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "UploadFile"
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

	errResp := req.ParseParameters(c)
	if errResp != nil {
		log.Error("session:", req.RequestUUID, " parse parameter error, ", errResp)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	uploadFileResponse, errResp := req.WriteFile()
	if errResp != nil {
		log.Error("session:", req.RequestUUID, " write file error, ", errResp)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", req.RequestUUID)
	c.JSON(http.StatusOK, uploadFileResponse)
}

func (req *UploadFileRequest) ParseParameters(c *gin.Context) (resp base.ApiBaseResponse) {
	multipartFileHeader, err := c.FormFile("file")
	if err != nil {
		log.Error("session:", req.RequestUUID, " FormFile error, ", err)
		return req.ErrorResponse("INTERNAL_ERROR", err)
	}

	openFile, err := multipartFileHeader.Open()
	if err != nil {
		log.Error("session:", req.RequestUUID, " open multipart file error, ", err)
		return req.ErrorResponse("INTERNAL_ERROR", err)
	}

	req.FileContent, err = io.ReadAll(openFile)
	if err != nil {
		log.Error("session:", req.RequestUUID, " read file content error, ", err)
		return req.ErrorResponse("INTERNAL_ERROR", err)
	}

	path := c.PostForm("path")
	if path == "" {
		log.Error("session:", req.RequestUUID, " params empty, path")
		return req.ErrorResponse("PARAMS_ERROR", "path")
	}

	mode := c.PostForm("mode")
	if mode == "" {
		req.Mode = 0700
	} else {
		v, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			log.Error("session:", req.RequestUUID, " parse mode error", err)
			return req.ErrorResponse("PARAMS_ERROR", "mode")
		}
		req.Mode = fs.FileMode(v)
	}

	return nil
}

// 1. 拿到文件的父目录, 如果该路径不存在,报错
func (req *UploadFileRequest) WriteFile() (uploadFileResponse *UploadFileResponse, resp base.ApiBaseResponse) {

	file, err := os.OpenFile(req.Path, os.O_RDWR|os.O_CREATE, req.Mode)
	if err != nil {
		log.Error("session:", req.RequestUUID, " open file error, ", err)
		return nil, req.ErrorResponse("INTERNAL_ERROR", err)
	}

	n, err := file.Write(req.FileContent)
	if err != nil {
		log.Error("session:", req.RequestUUID, " write file error, ", err)
		return nil, req.ErrorResponse("INTERNAL_ERROR", err)
	}
	if n != len(req.FileContent) {
		log.Error("session:", req.RequestUUID, " write file error, len(req.FileContent):", len(req.FileContent), " n:", n)
		return nil, req.ErrorResponse("INTERNAL_ERROR", "write length unexpected")
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error("session:", req.RequestUUID, " get file info error, ", err)
		return nil, req.ErrorResponse("INTERNAL_ERROR", "failed to get file info")
	}
	fileStat, _ := fileInfo.Sys().(*syscall.Stat_t)

	uploadFileResponse = &UploadFileResponse{
		Path: req.Path,
		Mode: fs.FileMode(fileStat.Mode) & 0777,
	}

	return uploadFileResponse, nil
}

func UploadFileHandler() *UploadFileRequest {
	return &UploadFileRequest{}
}
