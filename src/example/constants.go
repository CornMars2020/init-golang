package example

import "github.com/shopspring/decimal"

////////////////////////////// 常量 //////////////////////////////

var ZERO decimal.Decimal = decimal.NewFromInt(0)
var TWO decimal.Decimal = decimal.NewFromInt(2)

// 服务状态
type ServerStatus int64

const (
	STATUS_INITTED  ServerStatus = iota // 默认状态
	STATUS_PENDING                      // 初始化状态
	STATUS_STARTING                     // 启动中
	STATUS_STARTED                      // 已启动
	STATUS_STOPPING                     // 停止中
	STATUS_STOPPED                      // 已停止
)
