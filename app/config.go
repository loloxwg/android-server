package app

import (
	"androidServer/app/log"
	"androidServer/app/mysql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//middle "senseptrel-backend/modules/middlePlatform"
)

const (
	NotEnterString        string = "8e82665f-e6e8-4a13-a90c-59a4bf26cca8"
	NotEnterInt           int64  = 0xFFFF
	RegisterMySelfTTL     int64  = 30
	PreviewSizeWidth             = 38
	PreviewSizeHeight            = 38
	TimerMin                     = 60 //最小60秒更新一次，不可更短
	TemplateAddUserFile          = "add_user1.xlsx"
	TemplateLifecycleFile        = "lifecycle_template1.json"
	TemplatePolicyFile           = "policy_template1.json"

	PolicyAdmin       = "STORAGE_ADMINISTRATOR_POLICY"
	PolicyGroupLeader = "STORAGE_GROUP_LEADER_POLICY"
	PolicyGroupMember = "STORAGE_GROUP_MEMBER_POLICY"
)

var Conf Config

type Config struct {
	Name            string               `json:"register_name"`
	ListenIP        string               `json:"listen_ip"`
	TcpPort         uint16               `json:"tcp_port"`
	HttpPort        uint16               `json:"http_port"`
	HttpMonitorPort uint16               `json:"http_monitor_port"`
	InternalAPIURL  string               `json:"internal_api_url"`
	EtcdConfig      EtcdConfig           `json:"etcd_config"`
	LogConfig       log.LogConfig        `json:"log_config"`
	CrashLogConfig  CrashLog             `json:"crash_log_config"`
	LifecycleConfig LifecycleConfig      `json:"lifecycle_config"`
	TaskCycleTime   uint16               `json:"task_cycle_time"`
	MYSQL           MYSQLConfig          `json:"mysql"`
	SuperAdmin      []VirtualUserInfo    `json:"super_admin"`         //超级管理员账户
	SreAdmin        []VirtualUserInfo    `json:"sre_admin"`           //运维管理员
	VirtualBucket   VirtualBucketInfo    `json:"virutal_bucket_info"` //平台的虚拟bucket，每个集群都需要创建，并且没有配额限制，带有生命周期
	InitCluster     []ConfigCluster      `json:"init_cluster"`
	PathPrefix      string               `json:"path_prefix"`
	ClearBucket     ClearBucketConfig    `json:"clear_bucket"` //清空桶的配置
	Preview         PreviewPic           `json:"preview"`      //预览文件
	Timer           TimerInfo            `json:"timer"`
	DefaultMaxKeys  int64                `json:"default_max_keys"`
	JwtConfig       MiddlePlatformConfig `json:"jwt_config"`
	Email           EmailConfInfo        `json:"email"`
}

type MiddlePlatformConfig struct {
	Endpoint     string `json:"endpoint"`
	JwtAccessKey string `json:"jwt_access_key"`
	JwtSecretKey string `json:"jwt_secret_key"`
	ServiceName  string `json:"service_name"`
}

type TimerInfo struct {
	Bucket       int `json:"bucket"`        //bucket级别 单位分钟
	UserStorage  int `json:"user_storage"`  //用户存储量 单位分钟
	GroupStorage int `json:"group_storage"` //组级别的存储量 单位分钟
	Task         int `json:"task"`          //更新task的时间
}

type PreviewPic struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
type VirtualBucketInfo struct {
	Bucket             string `json:"bucket"`
	LifecycleLabelName string `json:"lifecycle_name"`
	DeleteDay          int    `json:"delete_day"`
	Path               string `json:"path"`
	LastPath           string `json:"last_path"` //持久化的路径，该路径不删除文件
}
type EmailConfInfo struct {
	Server   string `json:"server"`   //email邮箱地址
	Port     int    `json:"port"`     //端口号
	User     string `json:"user"`     //账户名
	Password string `json:"password"` //密码
}

//虚拟用户
type VirtualUserInfo struct {
	UID         string `json:"uid"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	KeyType     string `json:"key_type"`
	AccessKey   string `json:"access-key"`
	SecretKey   string `json:"secret-key"`
	LDAP        string `json:"ldap"`
	UserCaps    string `json:"user-caps"`
	GenerateKey bool   `json:"generate-key"`
	MaxBuckets  int    `json:"max-buckets"`
	Suspended   bool   `json:"suspended"`
	Tenant      string `json:"tenant"`
	AccountID   string `json:"account_id"`
	CephUid     string `json:"ceph_uid"`
}

type ConfigCluster struct {
	ClusterName   string `json:"cluster_name"`
	Region        string `json:"region"`
	ClusterType   string `json:"cluster_type"`
	StandardQuota int64  `json:"standard_quota"`
	QuickQuota    int64  `json:"quick_quota"`
	InEndpoint    string `json:"in_endpoint"`
	OutEndpoint   string `json:"out_endpoint"`
}

//mysql相关数据
type MYSQLConfig struct {
	Password      string `json:"password"`
	UserName      string `json:"username"`
	Network       string `json:"network"`
	Server        string `json:"server"`
	Port          int    `json:"port"`
	DataBase      string `json:"database"`
	MaxOpentConns int    `json:"max_open_conns"`
	MaxIdleConns  int    `json:"max_idle_conns"`
}

type CrashLog struct {
	LogDir    string `json:"log_dir"`
	LogPrefix string `json:"log_prefix"`
	LogSuffix string `json:"log_suffix"`
}

type EtcdConfig struct {
	EtcdNames  map[string]string `json:"etcd_names"`  //"short_name":"long_name"
	EtcdServer []string          `json:"etcd_server"` //"ip:port"
	CertFile   string            `json:"cert_file"`
	KeyFile    string            `json:"key_file"`
	CAFile     string            `json:"ca_file"`

	PolicyTemplate    string `json:"policy_template"` //批量导入policy模板
	LifecycleTemplate string `json:"lifecycle_template"`
	AddUserTemplate   string `json:"add_user_template"`
}

type LifecycleConfig struct {
	DeleteDaysLimit   Range `json:"delete_days_limit"`
	ArchivalDaysLimit Range `json:"archival_days_limit"`
	IaDaysLimit       Range `json:"ia_days_limit"`
}

type ClearBucketConfig struct {
	ThreadNum int   `json:"thread_num"`
	Limit     int64 `json:"limit"` //每次拉列表的次数
}
type Range struct {
	Lower int64 `json:"lower"`
	Upper int64 `json:"upper"`
}

func LoadConf(config string) error {

	// 配置读取
	bytes, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Printf("Read senseptrel-backend config file %v failed: %v\n", config, err)
		panic(err)
	}

	//设置默认值
	Conf.TaskCycleTime = 30
	Conf.PathPrefix = "/api/storage"
	Conf.DefaultMaxKeys = 100

	//解析
	if err = json.Unmarshal(bytes, &Conf); err != nil {
		fmt.Printf("Parse senseptrel-backend config json failed: %v\n", err)
		panic(err)
	}

	if Conf.ClearBucket.ThreadNum <= 0 || Conf.ClearBucket.ThreadNum >= 5 {
		Conf.ClearBucket.ThreadNum = 1 // 默认限定
	}
	if Conf.ClearBucket.Limit <= 0 || Conf.ClearBucket.Limit >= 1000 {
		Conf.ClearBucket.Limit = 100 //默认限定
	}
	if Conf.Preview.Width <= 0 || Conf.Preview.Width >= 200 {
		Conf.Preview.Width = 38
	}
	if Conf.Preview.Height <= 0 || Conf.Preview.Height >= 200 {
		Conf.Preview.Height = 38
	}
	if Conf.Timer.Bucket <= 0 || Conf.Timer.Bucket >= 100 {
		Conf.Timer.Bucket = 5
	}
	if Conf.Timer.UserStorage <= 0 || Conf.Timer.UserStorage >= 100 {
		Conf.Timer.UserStorage = 5
	}
	if Conf.Timer.GroupStorage <= 0 || Conf.Timer.GroupStorage >= 100 {
		Conf.Timer.GroupStorage = 5
	}
	if Conf.Timer.Bucket <= TimerMin {
		Conf.Timer.Bucket = TimerMin
	}
	if Conf.Timer.GroupStorage <= TimerMin {
		Conf.Timer.GroupStorage = TimerMin
	}
	if Conf.Timer.Task <= TimerMin {
		Conf.Timer.Task = TimerMin
	}
	if Conf.Timer.UserStorage <= TimerMin {
		Conf.Timer.UserStorage = TimerMin
	}

	fmt.Println("load config finished")
	return nil
}

func InitMysql() error {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", Conf.MYSQL.UserName, Conf.MYSQL.Password, Conf.MYSQL.Network, Conf.MYSQL.Server, Conf.MYSQL.Port, Conf.MYSQL.DataBase)
	fmt.Println("dsn:", dsn)
	err := mysql.InitMysql(dsn, Conf.MYSQL.MaxOpentConns, Conf.MYSQL.MaxIdleConns)
	if err != nil {
		fmt.Println("init mysql error,err:", err.Error())
		return err
	}
	fmt.Println("mysql init ok")
	return nil
}

//func InitJWT() error {
//	_, err := middle.NewJWT(Conf.JwtConfig.Endpoint, Conf.JwtConfig.JwtAccessKey, Conf.JwtConfig.JwtSecretKey, Conf.JwtConfig.ServiceName)
//	if err != nil {
//		fmt.Println("init jwt error,err:", err.Error())
//		return err
//	}
//	fmt.Println("jwt init ok")
//	return nil
//}
