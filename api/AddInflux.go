package api

import (
	"androidServer/app/log"
	"androidServer/handler/base"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	uncommon "github.com/gofrs/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"io/ioutil"
	"net/http"
	"time"
)

type AddInfluxRequest struct {
	base.ApiBaseRequest
}
type AddInfluxParams struct {
	Humidity    float64 `json:"humidity"`
	Location    string  `json:"location"`
	Pm25        float64 `json:"pm2_5"`
	Temperature int     `json:"temperature"`
}
type AddInfluxResponse struct {
	Humidity    float64 `json:"humidity"`
	Location    string  `json:"location"`
	Pm25        float64 `json:"pm2_5"`
	Temperature int     `json:"temperature"`
}

func (req *AddInfluxRequest) Process(c *gin.Context) {
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

	addInfluxResponse := &AddInfluxResponse{
		Humidity:    params.Humidity,
		Location:    params.Location,
		Pm25:        params.Pm25,
		Temperature: params.Temperature,
	}

	errResp = req.DBAddUser(params)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	c.Writer.Header().Set("request_uuid", raw.RequestUUID)

	log.Debug("session:", req.RequestUUID, ", add user success, resp:", addInfluxResponse)
	c.JSON(http.StatusOK, addInfluxResponse)

}

func (req *AddInfluxRequest) ParseParameters(c *gin.Context) (params *AddInfluxParams, resp base.ApiBaseResponse) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("session:", req.RequestUUID, ", Read request body failed error:", err)
		resp = req.ErrorResponse("INTERNAL_ERROR", "read body error")
		return
	}
	params = &AddInfluxParams{}
	if err := json.Unmarshal(bs, params); err != nil {
		log.Error("session:", req.RequestUUID, ", Unmarshal json failed error:", err, ", request:", string(bs))
		resp = req.ErrorResponse("INTERNAL_ERROR", "body format error")
		return
	}
	//if params.Location == "" {
	//	log.Error("session:", req.RequestUUID, ", parameter is invalid:", params)
	//	resp = req.ErrorResponse("PARAMS_ERROR", params)
	//	return
	//}
	return params, nil
}
func (req *AddInfluxRequest) DBAddUser(params *AddInfluxParams) (errResp base.ApiBaseResponse) {
	const token = "F-QFQpmCL9UkR3qyoXnLkzWj03s6m4eCvYgDl1ePfHBf9ph7yxaSgQ6WN0i9giNgRTfONwVMK1f977r_g71oNQ=="
	const bucket = "users_business_events"
	const org = "iot"

	client := influxdb2.NewClient("http://124.222.47.219:8086", token)
	// always close client at the end
	defer client.Close()
	writeAPI := client.WriteAPI(org, bucket)
	p := influxdb2.NewPointWithMeasurement("sensor").
		AddTag("unit", "th").
		AddField("temperature", float64(params.Temperature)).
		AddField("humidity", int(params.Humidity)).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()
	fmt.Println("finished writing point")
	return nil
}
func AddInfluxHandler() *AddInfluxRequest {
	return &AddInfluxRequest{}
}
