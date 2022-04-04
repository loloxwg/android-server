package api

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"androidServer/common/types"
	"androidServer/handler/base"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AddSignalRequest struct {
	base.ApiBaseRequest
}

type AddSignalParams struct {
	// UUID       string `json:"uuid"`
	// SignalName string `json:"signal_name" ` //名字
	// LocationX  string `json:"location_x" `  //性别
	// LocationY  string `json:"location_y" `  //民族
	// Rssi       string `json:"rssi"`         //

	RecordData string `json:"recordData"` //ibeacon记录
}

type AddSignalResponse struct { //返回结构体
	// UUID       string `json:"uuid"`
	// SignalName string `json:"signal_name" ` //名字
	// LocationX  string `json:"locon_x" `     //性别
	// LocationY  string `json:"locon_y" `     //民族
	// Rssi       string `json:"rssi"`         //
	// Message    string `json:"message"`

	RecordData string `json:"RecordData"`
}

// c.Writer.Header().Set("location","")

func (req *AddSignalRequest) Process(c *gin.Context) {
	var raw = new(base.ApiBaseRequest)
	raw.Action = "AddSignal"
	uv4, err := uncommon.NewV4()
	if err == nil {
		raw.RequestUUID = uv4.String()
	}
	req.SetRawReqData(raw.Action, raw.RequestUUID)

	// Begin Action
	beginTime := time.Now()
	var resp base.ApiBaseResponse
	log.Info(
		"Begin Action",
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
	//	注释了DB操作代码
	// errResp = req.DBAddUser(params)
	// if errResp != nil {
	// 	c.JSON(http.StatusBadRequest, errResp)
	// 	return
	// }

	addSignalResp := &AddSignalResponse{ //返回报文结构
		// UUID:       req.RequestUUID,
		// SignalName: params.SignalName,
		// LocationY:  params.LocationY,
		// LocationX:  params.LocationX,
		// Rssi:       params.Rssi,
		RecordData: params.RecordData,
	}
	c.Writer.Header().Set("request_uuid", raw.RequestUUID)
	c.Writer.Header().Set("location", "location")
	log.Debug("session:", req.RequestUUID, ", add user success, resp:", addSignalResp)
	c.JSON(http.StatusOK, addSignalResp)

}

func (req *AddSignalRequest) ParseParameters(c *gin.Context) (params *AddSignalParams, resp base.ApiBaseResponse) {
	params = &AddSignalParams{}
	// params.SignalName = c.PostForm("signal_name")
	// params.LocationX = c.PostForm("location_x")
	// params.LocationY = c.PostForm("location_y")
	// params.Rssi = c.PostForm("rssi")
	params.RecordData = c.PostForm("recordData")
	log.Debug("session: ", req.RequestUUID, "record_data:", params.RecordData)

	return params, nil
}
func (req *AddSignalRequest) DBAddUser(params *AddSignalParams) (errResp base.ApiBaseResponse) {
	db := mysql.GetDB()
	signalInfo := types.SignalInfo{
		// UUID:       req.RequestUUID,
		// SignalName: params.SignalName,
		// LocationY:  params.LocationY,
		// LocationX:  params.LocationX,
		// Rssi:       params.Rssi,
		RecordData: params.RecordData,
	}

	err := db.Transaction(func(db *gorm.DB) error {
		//user_table增加记录
		if err := db.Create(&signalInfo).Error; err != nil {
			log.Error("session:", req.RequestUUID, ", db create userInfo err:", err, ", userInfo:", signalInfo)
			return err
		}
		return nil
	})
	if err != nil {
		log.Error("session:", req.RequestUUID, ",err")
		return req.ErrorResponse("INTERNAL_ERROR", "db create signalInfo err:", err.Error())
	}
	return nil
}
func AddSignalHandler() *AddSignalRequest {
	return &AddSignalRequest{}
}

//func (req *AddSignalRequest) ReturnXY(params *AddSignalParams) (x, y string) {
//	///
//	//types.SignalInfo{}
//	db := mysql.GetDB()
//	db.Select("*").Find(types.SignalInfo{})
//
//	return
//}
