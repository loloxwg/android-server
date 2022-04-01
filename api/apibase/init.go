package apibase

import (
	"errors"
	"fmt"
)

var (
	NameMap = make(map[string]string)
)

func InitCache() error {
	fmt.Println(NameMap)
	NameMap["master"] = "master"
	if NameMap == nil {
		return errors.New("init userMaps failed")
	}
	fmt.Println("init userMaps successfully...")
	return nil
}
