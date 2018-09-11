package config

import (
	"os"
	"fmt"
	"os/exec"
	"path/filepath"
	"github.com/vaughan0/go-ini"
	"strings"
	"errors"
	)

// 这个文件主要是解析配置文件
type (
	FileParser struct {
		// 是否是调试模式
		Debug bool
		// 文件路径
		Path string
	}
)

var (
	// releaseSuffer 模式配置文件后缀
	releaseSuffer = ".ini"
	// debugSuffer 模式配置文件后缀
	debugSuffer = "_debug.ini"

	// 应用配置目录的环境变量
	envConfigDir = "APP_CONFIG_DIR"
)

func NewFileParser(debug bool, configPath ...string) FileParser {
	// 没有传递自己配置文件目录,使用默认的配置文件目录
	if len(configPath) == 0{
		configPath = make([]string, 1)
		configPath[0] = FileParser{}.defaultConfigPath()
	}

	return FileParser{
		Path:configPath[0],
		Debug:debug,
	}
}

// 实现 parse 方法来解析文件
func (f FileParser) Parse() (data map[string]ini.File, err error) {
	// 这个文件就是读取配置文件解析到变量中
	data = make(map[string]ini.File)
	var (
		tmpData ini.File

		tmpFileList = make([]string, 0)
		fileList    = make([]string, 0)
	)
	// 没有系统的配置路径，读取默认的
	if f.Path == ""{
		f.Path = f.defaultConfigPath()
	}
	// 判断是否是以结尾
	if !strings.HasSuffix(f.Path, "/"){
		f.Path = f.Path + "/"
	}

	tmpFileList, err = filepath.Glob(fmt.Sprintf("%s*%s", f.Path, releaseSuffer))
	if err != nil{
		return
	}
	//fmt.Println(tmpFileList)

	// 根据当前的环境来判断
	for _, v := range tmpFileList{
		if f.Debug && !strings.HasSuffix(v, debugSuffer){
			continue
		}
		if !f.Debug && strings.HasSuffix(v, debugSuffer){
			continue
		}
		fileList = append(fileList, v)
	}
	// 循环读取文件
	for _, v := range tmpFileList{
		if tmpData, err = f.parseFile(v); err != nil{
			return
		}
		data[ f.getConfigName(v) ] = tmpData
	}

	return
}

// 获取配置文件名称的键
func (f FileParser) getConfigName(filePath string) string {
	suffix := releaseSuffer
	if f.Debug {
		suffix = debugSuffer
	}
	return strings.Replace(filepath.Base(filePath), suffix, "", -1)
}

// 根据配置文件路径读取配置文件
func (f FileParser) parseFile(filePath string) (data ini.File, err error) {

	var (
		file      *os.File
		//fileBytes []byte
	)
	// 判断文件是否存在
	if file, err = os.Open(filePath); err != nil {
		if err == os.ErrNotExist {
			err = errors.New("配置文件不存在")
		}
		return
	}
	// 关闭文件流
	defer file.Close()

	// 解析配置文件
	data, err = ini.LoadFile(filePath)

	return
}


// 默认配置文件目录
func (f FileParser) defaultConfigPath() string {
	// 检测是否有配置文件环境变量，没有则从默认读取
	path := os.Getenv( envConfigDir )
	if path == ""{
		path = fmt.Sprintf("%s/conf/", f.runDir() )
	}
	return path
}

// 查找运行项目目录
func (f FileParser) runDir() string {
	rootDir, err := exec.LookPath(os.Args[0])
	if err != nil{
		panic(err)
	}
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		panic(err)
	}
	return filepath.Dir(rootDir)
}