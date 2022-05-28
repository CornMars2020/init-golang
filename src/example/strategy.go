package example

import (
	"init-golang/libs/config"
	"init-golang/libs/model"
	"log"
	"os"
	"time"
)

////////////////////////////// 策略 //////////////////////////////

// 格子策略
type Strategy struct {
	ServiceID string

	ApiCfg config.MarketAPI
	Logger *config.CFLogger
	Mdb    *model.ExampleDB
	Mrds   *model.MarketRedis
	Mtx    *config.Metrics

	status        ServerStatus
	serviceConfig model.ExampleSettings
}

func NewStrategy() *Strategy {
	ins := &Strategy{}

	ins.status = STATUS_INITTED

	return ins
}

func (ins *Strategy) Configure(serviceID string, apiCfg config.MarketAPI, logger *config.CFLogger, mdb *model.MarketDB, mrds *model.MarketRedis, mtx *config.Metrics) {
	if ins.status != STATUS_INITTED {
		log.Printf("status not initted\n")
		return
	}

	ins.ServiceID = serviceID
	ins.ApiCfg = apiCfg
	ins.Logger = logger

	ins.Mdb = &model.ExampleDB{MarketDB: mdb}
	ins.Mrds = mrds
	ins.Mtx = mtx

	ins.SetStatus(STATUS_PENDING)
}

func (ins *Strategy) SetStatus(status ServerStatus) {
	ins.status = status

	ins.Mtx.SetSvcValue(ins.ServiceID, "status", "service_status", float64(ins.status))
}

func (ins *Strategy) Start() {
	if ins.status != STATUS_PENDING && ins.status != STATUS_STOPPED {
		log.Printf("status not pending or stopped\n")
		return
	}
	ins.SetStatus(STATUS_STARTING)

	if !ins.Update() {
		os.Exit(1)
	}

	// ticker := time.NewTicker(time.Second * 5)
	// defer ticker.Stop()
	// go func() {
	// 	for {
	// 		if ins.status == STATUS_STOPPED {
	// 			break
	// 		}
	// 		<-ticker.C
	// 		ins.Update()
	// 	}
	// }()

	ins.SetStatus(STATUS_STARTED)

	ins.Mtx.SetSvcValue(ins.ServiceID, "process", "start_time", float64(time.Now().Unix()))

	go func() {
		for {
			if ins.status == STATUS_STARTED && ins.serviceConfig.Status == 1 {
				ins.main()
				time.Sleep(time.Millisecond * time.Duration(ins.serviceConfig.TickTime))
			}
		}
	}()
}

func (ins *Strategy) Update() bool {
	var err error
	serviceConfig, err := ins.Mdb.GetConfig(ins.ServiceID)
	if err != nil || (serviceConfig == model.ExampleSettings{}) {
		ins.Logger.Error("get service config failed")
		return false
	}

	if (ins.serviceConfig == model.ExampleSettings{}) || ins.serviceConfig.LogLevel != serviceConfig.LogLevel {
		ins.Logger.SetLevel(serviceConfig.LogLevel)
	}
	ins.serviceConfig = serviceConfig

	ins.Logger.Debug("service config loaded %v", ins.serviceConfig)

	ins.Mtx.SetSvcValue(ins.ServiceID, "process", "update_time", float64(time.Now().Unix()))

	return true
}

func (ins *Strategy) Stop() {
	if ins.status != STATUS_STARTED {
		log.Printf("status not started\n")
		return
	}
	ins.SetStatus(STATUS_STOPPING)

	ins.Mtx.SetSvcValue(ins.ServiceID, "process", "stop_time", float64(time.Now().Unix()))

	ins.SetStatus(STATUS_STOPPED)
}

// 业务主逻辑
func (ins *Strategy) main() {
	ins.Logger.Trace("run main logic")
}
