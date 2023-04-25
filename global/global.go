package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ServerConfig Config
	Logger       *zap.Logger
	EasterDB     *gorm.DB
)

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Charset  string `yaml:"charset"`
}

type Config struct {
	Db     MysqlConfig  `yaml:"db"`
	Log    LogConfig    `yaml:"log"`
	Consul ConsulConfig `yaml:"consul"`
	Ip     string       ``
	Port   int          ``
}

type LogConfig struct {
	Path  string `yaml:"path"`
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

type ConsulConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
