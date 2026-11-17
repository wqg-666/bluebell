package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 保存程序的所有配置信息。
var Conf = new(AppConfig)

// mapstructure,它告诉 Viper：将配置文件中的 port 映射到结构体的 Port 字段。
// config.yaml
// 具体配置数据
// Viper 读取、解析
// setting 统一管理配置结构体
// dao数据库配置logic业务配置main服务配置
type AppConfig struct {
	Name            string `mapstructure:"name"`
	Mode            string `mapstructure:"mode"`
	Version         string `mapstructure:"version"`
	Port            int    `mapstructure:"port"`
	StartTime       string `mapstructure:"start_time"`
	MachineID       int64  `mapstructure:"machine_id"`
	*LogConfig      `mapstructure:"log"`
	*MySQLConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*PostListConfig `mapstructure:"post_list"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DbName      string `mapstructure:"dbname"`
	Port        int    `mapstructure:"port"`
	MaxOpenConn int    `mapstructure:"max_open_conns"`
	MaxIdleConn int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type PostListConfig struct {
	Page  int64  `mapstructure:"page"`
	Size  int64  `mapstructure:"size"`
	Order string `mapstructure:"order"`
}

func Init(filePath string) (err error) {
	if err := viper.BindEnv("auth.jwt_secret", "BLUEBELL_JWT_SECRET"); err != nil {
		return err
	}
	if err := viper.BindEnv("auth.password_salt", "BLUEBELL_PASSWORD_SALT"); err != nil {
		return err
	}
	// 1. 指定配置文件路径
	viper.SetConfigFile(filePath)
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}
	//把读取到的配置信息反序列化到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
		return err
	}
	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
		}
	})
	return
}
