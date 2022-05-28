package main

import (
	"flag"
	"init-golang/libs/config"
	"init-golang/libs/model"
	"init-golang/src/example"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "net/http/pprof"
)

func main() {
	// parse params from command line when start
	flag.Parse()
	config.Init()

	// read config from yml files
	conf := config.ReadConf()

	if conf.Name == "" {
		log.Printf("invalid config name")
		return
	}

	config.InitLog(conf.Name, conf.LoggerCfg)
	logger := config.DefaultLogger(conf.Name)

	apiCfg := conf.MarketAPICfg

	// connect mysql
	mdb := &model.MarketDB{
		UserName: conf.MySQLCfg.UserName,
		Password: conf.MySQLCfg.Password,
		Host:     conf.MySQLCfg.Host,
		Port:     conf.MySQLCfg.Port,
		DBName:   conf.MySQLCfg.DBName,
	}
	if ok := mdb.Connect(); !ok {
		logger.Error("connect mysql database failed")
		return
	}
	logger.Info("connect database success %v", mdb.GetConnection())

	// connect redis
	mrds := &model.MarketRedis{
		Addr:     conf.RedisCfg.Addr,
		Database: conf.RedisCfg.Database,
		Pwd:      conf.RedisCfg.Pwd,
	}
	if ok := mrds.ConnectRedis(); !ok {
		logger.Error("connect redis database failed")
		return
	}
	logger.Info("connect redis success %v", mrds.GetRedisConnection())

	mtx := config.NewMetrics("example", "", conf.MonitorCfg.ErrorDetails)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("err", http.ListenAndServe(conf.MonitorCfg.GetMonitorAddress(), nil))
	}()

	strategy := example.NewStrategy()
	strategy.Configure(conf.Name, apiCfg, logger, mdb, mrds, mtx)
	strategy.Start()

	// 读取 MySQL 数据
	// var dbConfigs []activemarket.Cfg
	// mdb.GetConnection().Find(&dbConfigs)
	// 读取 Redis 数据
	// result := mredis.GetRedisConnection().Get(context.Background(), "justin")
	// logger.Info("redis values %v %s", result.Val(), result.Err())

	// ticker sync configs
	configTimer := time.NewTicker(time.Second * 5)
	defer configTimer.Stop()

	rand.Seed(time.Now().UnixNano())

	// recv system interrupt
	interruptSig := make(chan os.Signal, 1)
	signal.Notify(interruptSig, os.Interrupt)
	// recv system kill
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case sig := <-interruptSig:
			logger.Warn("interrupt signal %d recv", sig)

			strategy.Stop()

			mdb.Close()
			mrds.Close()

			os.Exit(0)
		case sig := <-killSig:
			logger.Warn("kill signal %d recv", sig)

			strategy.Stop()

			mdb.Close()
			mrds.Close()

			os.Exit(0)
		case <-configTimer.C:
			logger.Debug("configs %v", apiCfg)

			strategy.Update()
		}
	}
}
