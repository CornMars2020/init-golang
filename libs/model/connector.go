package model

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MarketDB Market 数据库连接配置(mysql)
type MarketDB struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        int32  `json:"port"`
	DBName      string `json:"dbname"`
	Connections int    `json:"connections"`

	conn *gorm.DB
}

// GetConnection 获取数据库连接实例
func (mdb *MarketDB) GetConnection() *gorm.DB {
	return mdb.conn
}

func (mdb *MarketDB) Close() (err error) {
	sqlDB, err := mdb.conn.DB()
	if err != nil {
		return
	}
	err = sqlDB.Close()
	return
}

// Connect 链接数据库
func (mdb *MarketDB) Connect() bool {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		mdb.UserName, mdb.Password,
		mdb.Host, mdb.Port, mdb.DBName,
	)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("failed to connect database:", err)
		return false
	}

	sqlDB, err := conn.DB()

	if err != nil {
		log.Println("failed to get database:", err)
		return false
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	if mdb.Connections > 0 {
		sqlDB.SetMaxOpenConns(mdb.Connections)
		sqlDB.SetMaxIdleConns(mdb.Connections * 2)
	}

	mdb.conn = conn
	return true
}

// MarketMem Market 缓存链接配置(redis)
type MarketMem struct {
}

type MarketRedis struct {
	Addr     string `json:"addr"`
	Database int    `json:"database"`
	Pwd      string `json:"pwd"`

	conn *redis.Client
}

func (mredis *MarketRedis) ConnectRedis() bool {
	opt := redis.Options{
		Addr:     mredis.Addr,
		Password: mredis.Pwd,
		DB:       mredis.Database,
	}

	redisClient := redis.NewClient(&opt)
	mredis.conn = redisClient

	return true
}

func (mredis *MarketRedis) Close() (err error) {
	err = mredis.conn.Close()
	return
}

// GetConnection 获取数据库连接实例
func (mredis *MarketRedis) GetRedisConnection() *redis.Client {
	return mredis.conn
}
