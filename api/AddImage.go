package api

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/common/types"
	"androidServer/handler/base"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AddImageRequest struct {
	base.ApiBaseRequest
	Url string
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
		UUID: req.RequestUUID,
	}
	log.Debug("session:", req.RequestUUID, ", add Image success, resp:", addImageResp)
	c.JSON(http.StatusOK, addImageResp)

}

func (req *AddImageRequest) ParseParameters(c *gin.Context) (params *AddImageParams, resp base.ApiBaseResponse) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("session:", req.RequestUUID, ", Read request body failed error:", err)
		resp = req.ErrorResponse("INTERNAL_ERROR", "read body error")
		return
	}
	params = &AddImageParams{}
	if err := json.Unmarshal(bs, params); err != nil {
		log.Error("session:", req.RequestUUID, ", Unmarshal json failed error:", err, ", request:", string(bs))
		resp = req.ErrorResponse("INTERNAL_ERROR", "body format error")
		return
	}

	return params, nil
}
func (req *AddImageRequest) OSSAddImage(params *AddImageParams) (errResp base.ApiBaseResponse) {
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
	name := params.FileName
	f := strings.NewReader(string(params.FileContent))

	_, err := c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		log.Error("session:", req.RequestUUID, "INTERNAL_ERROR", err.Error())
		return
	}
	objurl := c.Object.GetObjectURL(params.FileName)
	req.Url = objurl.String()
	return nil
}
func (req *AddImageRequest) DBAddImage(params *AddImageParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()
	ImageInfo := types.ImageInfo{
		UUID:      req.RequestUUID,
		ImageName: params.FileName,
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
