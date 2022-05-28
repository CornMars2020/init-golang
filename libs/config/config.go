package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var locations = []string{
	"/etc/",
	"$HOME/.etc/",
	".",
	"./etc/",
	"../etc/",
	"../../etc/",
}

var (
	name     string
	confPath string
	env      string
)

func init() {
	flag.StringVar(&name, "n", "name", "service name")
	flag.StringVar(&confPath, "c", "configs", "config path")
	flag.StringVar(&env, "e", "env", "environment name")
}

func Init() {
	if confPath != "" {
		locations = append([]string{confPath}, locations...)
	}
}

// MySQL mysql 链接配置
type MySQL struct {
	UserName    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int32  `mapstructure:"port" json:"port"`
	DBName      string `mapstructure:"dbname" json:"dbname"`
	Connections string `mapstructure:"connections" json:"connections"`
}

// redis 链接配置
type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Database int    `mapstructure:"database" json:"database"`
	Pwd      string `mapstructure:"pwd" json:"pwd"`
	PreFix   string `mapstructure:"prefix" json:"prefix"`
}

// MarketURL 平台服务配置细节
type MarketURL struct {
	PubURL      string `mapstructure:"pub_url,omitempty" json:"pub_url,omitempty"`
	IsPubCache  int    `mapstructure:"is_pub_cache,omitempty" json:"is_pub_cache,omitempty"`
	PrivURL     string `mapstructure:"priv_url,omitempty" json:"priv_url,omitempty"`
	IsPrivCache int    `mapstructure:"is_priv_cache,omitempty" json:"is_priv_cache,omitempty"`
}

// MarketAPI 平台服务配置
type MarketAPI struct {
	Spot    MarketURL `mapstructure:"spot" json:"spot"`
	Futures MarketURL `mapstructure:"futures" json:"futures"`
	Swap    MarketURL `mapstructure:"swap" json:"swap"`
	Cache   MarketURL `mapstructure:"cache" json:"cache"`
	Extra   MarketURL `mapstructure:"extra" json:"extra"`
}

// Logger 日志配置文件
type Logger struct {
	Path           string `mapstructure:"path" json:"path"`
	FileRotateMode string `mapstructure:"file_rotate_mode" json:"file_rotate_mode"`
	TimeFormat     string `mapstructure:"time_format" json:"time_format"`
	IsHideKey      bool   `mapstructure:"is_hide_key" json:"is_hide_key"`
	IsColor        bool   `mapstructure:"is_color" json:"is_color"`
	IsFieldsOrder  bool   `mapstructure:"is_fields_order" json:"is_fields_order"`
}

type Monitor struct {
	Address      string `mapstructure:"address" json:"address,omitempty"`
	Host         string `mapstructure:"host" json:"host,omitempty"`
	Port         int64  `mapstructure:"port" json:"port,omitempty"`
	ErrorDetails bool   `mapstructure:"error_details" json:"error_details,omitempty"`
}

func (monitor *Monitor) GetMonitorAddress() string {
	if monitor.Address != "" {
		return monitor.Address
	}

	if monitor.Port != 0 {
		return fmt.Sprintf("%s:%d", monitor.Host, monitor.Port)
	}

	// default metrics monitor port
	return ":9090"
}

// Config 程序配置文件结构
type Config struct {
	Name         string    `mapstructure:"name" json:"name"`
	MySQLCfg     MySQL     `mapstructure:"mysql" json:"mysql"`
	RedisCfg     Redis     `mapstructure:"redis" json:"redis"`
	MarketAPICfg MarketAPI `mapstructure:"marketapi" json:"marketapi"`
	LoggerCfg    Logger    `mapstructure:"logger" json:"logger"`
	MonitorCfg   Monitor   `mapstructure:"monitor" json:"monitor"`
}

func onConfigChange(in fsnotify.Event) {
	log.Println("Config file changed:", in.Name)
}

// ReadConf 读取配置文件
func ReadConf() *Config {
	// 配置文件名称, 默认扩展名
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	runpath, _ := os.Getwd()

	// 配置文件查找路径
	locations = append(locations, runpath+"/etc/")

	for _, location := range locations {
		log.Printf("config file search path %s", location)
		viper.AddConfigPath(location)
	}

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		// 配置文件读取出错
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 文件找不存在
			log.Printf("文件不存在")
		} else {
			// 文件存在, 其他错误
			log.Printf("文件存在, 其他错误: %s", err)
		}

		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	conf := Config{Name: name}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("fatal error parse config file: %s", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(onConfigChange)

	return &conf
}
