package model

const (
	// Stopped 服务处于停止状态 not running
	Stopped = 0
	// StartPending 服务准备启动, 启动前执行一些指定逻辑 starting
	StartPending = 3
	// StopPending 服务即将进入停止状态, 停止前执行一些指定逻辑 stopping
	StopPending = 2
	// Running 服务处于开启状态 running
	Running = 1
	// ContinuePending 服务即将恢复运行状态, 恢复前执行一些指定逻辑 continue is pending
	ContinuePending = 6
	// PausePending 服务即将进入暂停状态, 暂停前执行一些指定逻辑 pause is pending
	PausePending = 5
	// Paused 服务处于暂停状态 paused
	Paused = 4
)

// Configs 服务配置接口
type Configs interface {
	GetUniqID() string
	GetSymbol() string
	Equals(cfg *Configs) bool
	Clone() *Configs
}
