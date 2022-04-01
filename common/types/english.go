package types

import "gorm.io/gorm"

type EnglishInfo struct {
	gorm.Model
	UUID string `json:"uuid" gorm:"type:varchar(36);primaryKey"`
	Word string `json:"user_name" gorm:"type:varchar(32);not null"` //名字
}

// TableName UserInfo绑定表名
func (v EnglishInfo) TableName() string {
	return "word_table"
}
