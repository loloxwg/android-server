package types

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	UUID     string `json:"uuid" gorm:"type:varchar(36);primaryKey"`
	Username string `json:"user_name" gorm:"type:varchar(32);not null"` //名字
	Sex      string `json:"sex" gorm:"type:varchar(32);not null"`       //性别
	Nation   string `json:"nation" gorm:"type:varchar(32);not null"`    //民族
	Eid      string `json:"eid" gorm:"type:varchar(32);not null"`       //民族
	Address  string `json:"address" gorm:"type:varchar(1024);not null"` //地址
	Birthday string `json:"birthday" gorm:"type:varchar(32)"`           //生日
	PassWord string `json:"password" gorm:"type:varchar(32)"`
}

// TableName UserInfo绑定表名
func (v UserInfo) TableName() string {
	return "user_table"
}
