package mysql

import (
	"androidServer/app/log"
	"androidServer/common/types"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
负责mysql初始化的相关工作
建表顺序：
1.user_table
2.group_table
3.group_member_table
4.cluster_table
5.bucket_table
6.outside_group_table
7.bucket_storage_table
8.user_storage_table
9.group_cluster_storage_table
10.total_storage_table
11.label_table
12.lifecycle_table
13 .policy_table
14.event_table
15.event_result_table
16.notice
17.log
18.task
19.role
20 object_tag_table
*/
//测试数据
const (
	USERNAME = "wangjingwen"
	PASSWORD = "WAng80488"
	NETWORK  = "tcp"
	SERVER   = "sh.paas.sensetime.com"
	PORT     = 36010
	DATABASE = "web_portal"
)

var _db *gorm.DB

//GetDB 获取db连接池句柄
func GetDB() *gorm.DB {
	return _db
}

//InitMysql 连接数据库，并检测数据库表结构，如果没有表的话，自动构建
func InitMysql(dsn string, maxOpentCons int, maxIdleConns int) error {
	//dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("open mysql error:", err.Error())
		return err
	}
	err = checkTables(_db)
	if err != nil {
		log.Error("check and build mysql failed,err:", err.Error())
		return err
	}
	sqlDB, _ := _db.DB()
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	return nil
}

//checkTables 检查表关系
func checkTables(DB *gorm.DB) error {
	exist := DB.Migrator().HasTable(&types.UserInfo{})
	if !exist {
		err := DB.Migrator().CreateTable(&types.UserInfo{})
		if err != nil {
			fmt.Println("create UserInfo error:", err.Error())
			return err
		}
	}
	exist = DB.Migrator().HasTable(&types.EnglishInfo{})
	if !exist {
		err := DB.Migrator().CreateTable(&types.EnglishInfo{})
		if err != nil {
			fmt.Println("create UserInfo error:", err.Error())
			return err
		}
	}
	//exist = DB.Migrator().HasTable(&types.SignalInfo{})
	//if !exist {
	//	err := DB.Migrator().CreateTable(&types.SignalInfo{})
	//	if err != nil {
	//		fmt.Println("create UserInfo error:", err.Error())
	//		return err
	//	}
	//}
	return nil
}
