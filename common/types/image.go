package types

import "gorm.io/gorm"

type ImageInfo struct {
	gorm.Model
	UUID      string `json:"uuid" gorm:"type:varchar(36);primaryKey"`
	Url       string `json:"url" gorm:"type:varchar(32);not null"`
	ImageName string `json:"image_name" gorm:"type:varchar(255);not null"`
	Collect   int    `json:"collect" gorm:"type:int;not null"`
}

// TableName UserInfo绑定表名
func (v ImageInfo) TableName() string {
	return "image_table"
}
