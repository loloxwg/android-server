package api

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/common/types"
	"androidServer/handler/base"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AddImageRequest struct {
	base.ApiBaseRequest
	Url         string
	FileContent []byte
	FileName    string
}

type AddImageParams struct {
	FileContent []byte `json:"file_content"`
	FileName    string `json:"file_name"`
}

type AddImageResponse struct {
	UUID     string `json:"uuid"`
	Url      string `json:"url"`
	FileName string `json:"file_name"`
}

func (req *AddImageRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "AddImage"
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
	errResp = req.OSSAddImage(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	errResp = req.DBAddImage(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	addImageResp := &AddImageResponse{
		UUID:     req.RequestUUID,
		Url:      req.Url,
		FileName: req.FileName,
	}
	log.Debug("session:", req.RequestUUID, ", add Image success, resp:", addImageResp)
	c.JSON(http.StatusOK, addImageResp)

}

func (req *AddImageRequest) ParseParameters(c *gin.Context) (params *AddImageParams, resp base.ApiBaseResponse) {
	multipartFileHeader, err := c.FormFile("file_content")
	if err != nil {
		log.Error("session:", req.RequestUUID, " FormFile error, ", err.Error())
		return nil, req.ErrorResponse("INTERNAL_ERROR", err)
	}

	openFile, err := multipartFileHeader.Open()
	if err != nil {
		log.Error("session:", req.RequestUUID, " open multipart file error, ", err)
		return nil, req.ErrorResponse("INTERNAL_ERROR", err)
	}

	req.FileContent, err = io.ReadAll(openFile)
	if err != nil {
		log.Error("session:", req.RequestUUID, " read file content error, ", err)
		return nil, req.ErrorResponse("INTERNAL_ERROR", err)
	}

	name := c.PostForm("file_name")
	if name == "" {
		log.Error("session:", req.RequestUUID, " params empty, file_name")
		return nil, req.ErrorResponse("PARAMS_ERROR", "file_name")
	}
	req.FileName = name

	return params, nil
}
func (req *AddImageRequest) OSSAddImage(params *AddImageParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()

	testImageInfo := &types.ImageInfo{}
	err := db.Table("image_table").Where("image_name=?", req.FileName).First(testImageInfo).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("session:", req.RequestUUID, ",db found bucket by uuid error:", err)
		return req.ErrorResponse("INTERNAL_ERROR", err.Error())
	} else if testImageInfo.ImageName != "" {
		log.Error("session:", req.RequestUUID, ", image_name already exist")
		return req.ErrorResponse("INTERNAL_ERROR", "image_name_already_exist")
	}

	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://zxwzs-1311854978.cos.ap-nanjing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: "AKIDnQntiGYZJmzvay6xWopa4GBqELBk23oH",
			// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: "FDtJTLwGJ5anjPdeOr3yc18yXprc5HCr",
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	// Case1 上传对象
	name := req.FileName
	f := strings.NewReader(string(req.FileContent))
	filename := strings.Split(req.FileName, ".")
	ext := filename[len(filename)-1]
	imgMime := ""
	switch ext {
	case "jpg":
		imgMime = "image/jpg"
		break
	case "jpeg":
		imgMime = "image/jpeg"
		break
	case "png":
		imgMime = "image/png"
		break
	case "svg":
		imgMime = "image/svg+xml"
		break
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: imgMime,
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
			//XCosACL: "private",
		},
	}
	_, err = c.Object.Put(context.Background(), name, f, opt)
	if err != nil {
		log.Error("session:", req.RequestUUID, "INTERNAL_ERROR", err.Error())
		return
	}
	objurl := c.Object.GetObjectURL(req.FileName)
	req.Url = objurl.String()
	return nil
}
func (req *AddImageRequest) DBAddImage(params *AddImageParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()
	ImageInfo := types.ImageInfo{
		UUID:      req.RequestUUID,
		ImageName: req.FileName,
		Url:       req.Url,
	}

	err := db.Transaction(func(db *gorm.DB) error {

		//Image_table增加记录
		if err := db.Create(&ImageInfo).Error; err != nil {
			log.Error("session:", req.RequestUUID, ", db create ImageInfo err:", err, ", ImageInfo:", ImageInfo)
			return err
		}
		return nil
	})
	if err != nil {
		return req.ErrorResponse("INTERNAL_ERROR", err)
	}
	return nil
}
func AddImageHandler() *AddImageRequest {
	return &AddImageRequest{}
}
