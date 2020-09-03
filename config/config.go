// 根据文件进行初始化配置，例如设定日志的配置，并监控文件的变化
package config

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/team4yf/yf-fpm-server-go/pkg/log"
	"github.com/team4yf/yf-fpm-server-go/pkg/utils"
)

// Config 读取配置
type Config struct {
	Name string
}

// Init 初始化配置，默认读取config.local.yaml
func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}
	viper.SetDefault("mode", "debug")
	viper.SetDefault("log.writers", "stdout")
	viper.SetDefault("log.level", "DEBUG")
	viper.SetDefault("log.logger_file", "./logs/app.log")
	viper.SetDefault("log.logger_warn_file", "./logs/app.warn.log")
	viper.SetDefault("log.logger_error_file", "./logs/app.err.log")
	viper.SetDefault("log.log_format_text", true)
	viper.SetDefault("log.rollingPolicy", "daily")
	viper.SetDefault("log.log_rotate_date", 1)
	viper.SetDefault("log.log_rotate_size", 1)
	viper.SetDefault("log.log_backup_count", 7)
	// 初始化日志包
	c.initLog()
	return nil
}

func (cfg *Config) initConfig() error {
	viper.AutomaticEnv()      // 读取匹配的环境变量
	viper.SetEnvPrefix("FPM") // 读取环境变量的前缀为 FPM
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	deployMode := viper.GetString("deploy.mode")
	if cfg.Name != "" {
		viper.SetConfigFile(cfg.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		if deployMode == "" {
			deployMode = "local"
		}
		deployMode = strings.ToLower(deployMode)
		//read config file by FPM_DEPLOY_MODE=PROD
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config." + deployMode)
	}
	//读取配置目录下的文件列表以及其后缀
	//是否存在配置目录
	if !utils.IsDir("conf") {
		//没有，则直接返回
		return nil
	}

	list, err := ioutil.ReadDir("conf") //要读取的目录地址DIR，得到列表
	if err != nil {
		fmt.Printf("Reade conf folder error: %v", err)
		// 读取失败，直接返回
		return nil
	}
	for _, info := range list {
		//遍历目录下的内容，获取文件详情，同os.Stat(filename)获取的信息
		confFileName := info.Name() //文件名
		if strings.HasPrefix(confFileName, "config."+deployMode) {
			ext := path.Ext(confFileName)
			viper.SetConfigType(ext[1:]) // 设置配置文件格式
		}

	}

	//
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return errors.WithStack(err)
	}
	return nil
}

func (cfg *Config) initLog() {
	config := log.Config{
		Writers:         viper.GetString("log.writers"),
		LoggerLevel:     viper.GetString("log.level"),
		LoggerFile:      viper.GetString("log.logger_file"),
		LoggerWarnFile:  viper.GetString("log.logger_warn_file"),
		LoggerErrorFile: viper.GetString("log.logger_error_file"),
		LogFormatText:   viper.GetBool("log.log_format_text"),
		RollingPolicy:   viper.GetString("log.rollingPolicy"),
		LogRotateDate:   viper.GetInt("log.log_rotate_date"),
		LogRotateSize:   viper.GetInt("log.log_rotate_size"),
		LogBackupCount:  viper.GetInt("log.log_backup_count"),
	}
	err := log.NewLogger(&config, log.InstanceZapLogger)
	if err != nil {
		fmt.Printf("InitWithConfig err: %v", err)
	}
}
