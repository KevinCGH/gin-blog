package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Mode          string // debug | release
		Port          string
		DbType        string // mysql | sqlite3 | postgres
		DbAutoMigrate bool
		DbLogLevel    string // silent | error | warn | info
	}
	Log struct {
		Level     string // debug | info | warn | error
		Prefix    string
		Format    string // text | json
		Directory string
	}
	JWT struct {
		Secret string
		Expire int64 // hour
		Issuer string
	}
	Mysql struct {
		Host     string // 服务器地址
		Port     string
		Config   string // 高级配置
		Dbname   string
		Username string
		Password string
	}

	Postgres struct {
		Host     string
		Port     string
		Dbname   string
		Username string
		Password string
	}
	SQLite struct {
		Dsn string // DataSource Name
	}
	Session struct {
		Name   string
		Salt   string
		MaxAge int
	}
}

var Conf *Config

func GetConfig() *Config {
	if Conf == nil {
		log.Panic("配置文件未初始化")
		return nil
	}
	return Conf
}

func InitConfig() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AutomaticEnv()
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // SERVER_PORT => SERVER.PORT

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("配置文件读取失败: %v", err)
	}

	if err := v.Unmarshal(&Conf); err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
	}
	log.Println("✅配置加载成功！")
}

func (*Config) DbType() string {
	if Conf.Server.DbType == "" {
		Conf.Server.DbType = "sqlite"
	}
	return Conf.Server.DbType
}

func (*Config) DbDSN() string {
	switch Conf.Server.DbType {
	case "mysql":
		conf := Conf.Mysql
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Config,
		)
	case "postgres":
		conf := Conf.Postgres
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			conf.Host, conf.Username, conf.Password, conf.Dbname, conf.Port,
		)
	case "sqlite":
		return Conf.SQLite.Dsn
	default:
		Conf.Server.DbType = "sqlite"
		if Conf.SQLite.Dsn == "" {
			Conf.SQLite.Dsn = "file::memory:"
		}
		return Conf.SQLite.Dsn
	}
}
