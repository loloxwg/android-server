package base

import (
	"fmt"
)

var responseCode = map[string]map[string]interface{}{
	// Public
	"RETCODE_NOT_REGIST":               {"RetCode": -1, "Message": "RetCode is not exists"},      //
	"SERVICE_ERROR":                    {"RetCode": 130, "Message": "Service error and break"},   //
	"SERVICE_UNAVAILABLE":              {"RetCode": 150, "Message": "Service unavailable"},       //
	"MISSING_ACTION":                   {"RetCode": 160, "Message": "MISSING_ACTION"},            //
	"MISSING_PARAMS":                   {"RetCode": 220, "Message": "Missing params [%s]"},       //
	"PARAMS_ERROR":                     {"RetCode": 230, "Message": "Params [%s] not available"}, //
	"PARAM_TOO_MANY":                   {"RetCode": 240, "Message": "Params [%s] too large"},     //
	"INTERNAL_ERROR":                   {"RetCode": 1000, "Message": "%+v"},                      //
	"CLUSTER_NOT_FOUND":                {"RetCode": 1001, "Message": "cluster not found"},
	"RESOURCES_NOT_CLEAR":              {"RetCode": 1002, "Message": "resources not clear"},   //资源不为空
	"GROUP_NOT_FOUND":                  {"RetCode": 1003, "Message": "group not found"},       //组没有找到
	"BUCKET_LIST_NOT_FOUND":            {"RetCode": 1004, "Message": "bucket list not found"}, //bucketlist 不存在
	"DB_FIND_ERROR":                    {"RetCode": 1005, "Message": "db find error"},         //mysql 查找错误
	"USER_NOT_FOUND":                   {"RetCode": 1006, "Message": "user not found"},
	"USER_ALREADY_IN_GROUP":            {"RetCode": 1007, "Message": "user already in group"},
	"ILLEGAL_USER_DATA":                {"RetCode": 1008, "Message": "Illegal user data"}, //user 参数校验不合法
	"USER_ALREADY_EXISTS":              {"RetCode": 1009, "Message": "user already exists"},
	"QUOTA_INADEQUATE":                 {"RetCode": 1010, "Message": "quota inadequate"},          //资源不足
	"CEPH_PROXY_INTERNAL_ERROR":        {"RetCode": 1011, "Message": "ceph proxy internal error"}, //ceph网管错误
	"BUCKET_NOT_FOUND":                 {"RetCode": 1012, "Message": "bucket not found"},
	"USER_NOT_HAVE_PERMISSION":         {"RetCode": 1013, "Message": "the user does not have permission"}, //更新组，但是未做实质修改; 用户没有权限删除 notice
	"AUTH_ERROR":                       {"RetCode": 1014, "Message": "check user auth"},                   //组长aksk有问题
	"EVENT_NOT_FOUND":                  {"RetCode": 1015, "Message": "event not found"},                   //事件未找到
	"EVENT_NOT_SUPPORT":                {"RetCode": 1016, "Message": "event not support"},
	"LABEL_NOT_FOUND":                  {"RetCode": 1017, "Message": "label not found"},
	"POLICY_NOT_FOUND":                 {"RetCode": 1018, "Message": "policy not found"},
	"ACCOUNT_ID_FROM_HEADER_ERROR":     {"RetCode": 1019, "Message": "account_id from header error"},
	"IMAGE_SCAlE_ERROR":                {"RetCode": 1020, "Message": "image scale error"},
	"EVENT_HAS_ALREADY_APPROVAL":       {"RetCode": 1021, "Message": "event has already approval"},
	"EVENT_USER_ERROR":                 {"RetCode": 1022, "Message": "event user error"}, //申请人身份错误
	"TASK_NOT_FOUND":                   {"RetCode": 1023, "Message": "task not found"},
	"TASK_ALREADY_EXISTS":              {"RetCode": 1024, "Message": "task already exists"},
	"BUCKET_STATUS_IN_CLEARING":        {"RetCode": 1025, "Message": "bucket status in clearing"}, //bucket状态在清理中，不可写入
	"TASK_NOT_FINISH":                  {"RetCode": 1026, "Message": "task not finish"},           //没有做完，不允许下载
	"LIFECYCLE_NOT_FOUND":              {"RetCode": 1027, "Message": "lifecycle not found"},
	"ENABLED_LIFECYCLE_ALREADY_EXISTS": {"RetCode": 1028, "Message": "enabled lifecycle already exists"},
	"JSON_ERROR":                       {"RetCode": 1029, "Message": "json error"},
	"ADMINS_NOT_FOUND":                 {"RetCode": 1030, "Message": "admins not found"},
	"SUPER_ADMIN_NOT_ALLOWED_DELETE":   {"RetCode": 1031, "Message": "super admin not allowed delete"},
	"TASK_TYPE_ILLEGAL":                {"RetCode": 1032, "Message": "task type Illegal"},
	"TASK_STATUS_ILLEGAL":              {"RetCode": 1033, "Message": "task status Illegal"},
	"EMAIL_SEND_ERROR":                 {"RetCode": 1034, "Message": "email send error"},
	"NOTICE_NOT_FOUND":                 {"RetCode": 1035, "Message": "notice not exist"},
	"LABEL_VOLUE_KEY":                  {"RetCode": 1036, "Message": "label_name_key is already exists"},
	"LABEL_NAME":                       {"RetCode": 1037, "Message": "label_name is already exists"},
}

func generateErrorCode(name string, params ...interface{}) map[string]interface{} {
	res := map[string]interface{}{"RetCode": -1, "Message": "RetCode is not exists"}

	errorCode, ok := responseCode[name]
	if !ok {
		return res
	}

	message, ok := errorCode["Message"].(string)
	if !ok {
		return res
	}

	res["Message"] = fmt.Sprintf(message, params...)
	res["RetCode"] = errorCode["RetCode"]
	return res
}
