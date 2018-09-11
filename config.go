package config

import (
	"github.com/vaughan0/go-ini"
	)

// 这个是获取配置文件   以来就得加载所有的配置文件 暂时不区分版本
type (
	Parser interface {
		// 使用接口用户解析文件
		Parse() (map[string]ini.File,error)
	}

	config struct {
		Debug bool    				// 是否是debug模式
		data map[string]ini.File	// 存储每个配置文件的参数
		parse Parser
	}
)

var (
	// 对外提供 Config参数
	c = config{}
	Config = c.data
)

// 获取配置文件
func Data(filename string) ini.File {
	return c.data[filename]
}

// 初始化需要做的事情 设置debug 加载相应的配置文件
func Init(debug bool, parse ...Parser) (err error) {
	return c.init(debug, parse...)
}

// 使用封装特性来封装  第二个参数表示有人自己重写一个解析配置文件方法
func (c *config) init(debug bool, parse ...Parser) (err error) {
	c.Debug = debug
	// 没有自定义解析配置文件函数调用系统的函数
	if len(parse) == 0{
		parse = make([]Parser,1)
		parse[0] = NewFileParser(debug)
	}

	c.parse = parse[0]
	return c.load()
}

func (c *config) load() (err error) {
	if c.data, err = c.parse.Parse(); err != nil{
		return
	}
	Config = c.data
	return
}
