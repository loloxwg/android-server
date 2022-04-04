package types

import "gorm.io/gorm"

type SignalInfo struct {
	gorm.Model
	// UUID       string `json:"uuid" gorm:"type:varchar(36);primaryKey"`
	// SignalName string `json:"signal_name" gorm:"type:varchar(32);not null"` //名字
	// LocationX  string `json:"locon_x" gorm:"type:varchar(32);not null"`
	// LocationY  string `json:"locon_y" gorm:"type:varchar(32);not null"`
	// Rssi       string `json:"rssi" gorm:"type:varchar(32);not null"`
	//
	RecordData string `json:"record_data" gorm:"type:varchar(255);not null"`
}

// TableName UserInfo绑定表名
func (v SignalInfo) TableName() string {
	return "signal_tab"
}
