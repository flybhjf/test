package config

import (
	"fmt"
	"log"

	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config struct to read configuration from YAML file
type Config struct {
	MySQL struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
	} `yaml:"mysql"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

var (
	DB        *gorm.DB
	RedisPool *redis.Pool
	Conn      redis.Conn
)

func Connect() {
	config, err := loadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// MySQL Connection
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.MySQL.Username, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.DBName)
	// DB, err = sql.Open("mysql", mysqlDSN)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to MySQL: %v", err)
	// }
	// defer DB.Close()
	DB, err = gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	// Redis Connection
	redisAddr := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
	RedisPool = newRedisPool(redisAddr, config.Redis.Password)
	Conn = RedisPool.Get()
	// You can now use the 'db' and 'redisPool' objects to interact with MySQL and Redis.
}

func loadConfig(filename string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func newRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		MaxActive:   0,
		IdleTimeout: 30,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}
