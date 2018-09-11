package config

import (
	"testing"
		"os"
	"fmt"
)

// 测试是否是能够成功
func TestConfig_Init(t *testing.T) {
	appDir := os.Getenv("GOPATH") + "/src/config/conf/"
	fileParser := NewFileParser(true, appDir)
	fileParser.Debug = true
	if err := Init(true, fileParser); err != nil {
		t.Error(err)
	}

	//version := String(Config["mysql"].Get("app", "version"))
	version := String(Data("mysql").Get("app", "version"))

	if "v2.0" != version {
		t.Error("debug 模式下读取的 version 不对,version:", version)
	}
	fmt.Println(version)
}