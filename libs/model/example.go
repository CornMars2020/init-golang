package model

type ExampleSettings struct {
	// ServiceID 唯一 ID
	ServiceID string `gorm:"column:service_id;VARCHAR(32);NOT NULL;DEFAULT '';" json:"service_id" mapstructure:"service_id"`
	// Status 状态: 0 关闭, 1 开启, 2 关闭并撤单
	Status int16 `gorm:"column:status;TINYINT(4);NOT NULL;DEFAULT 0;" json:"status" mapstructure:"status"`
	// TickTime 执行周期
	TickTime int64 `gorm:"column:tick_time;INT(11);NOT NULL;DEFAULT 1000;" json:"tick_time" mapstructure:"tick_time"`
	// LogLevel 日志等级: 0 关键数据 + Error 级别, 1 Info 级别, 2 Debug 级别, 3 Trace 级别, 4 Everything 级别
	LogLevel int16 `gorm:"column:log_level;TINYINT(4);NOT NULL;DEFAULT 0;" json:"log_level" mapstructure:"log_level"`
}

// ExampleDB
type ExampleDB struct {
	MarketDB *MarketDB
}

func (ins *ExampleDB) GetConfig(serviceID string) (ExampleSettings, error) {
	var config ExampleSettings

	ins.MarketDB.GetConnection().
		Table("grid_orders_configs").
		Where("service_id = ?", serviceID).
		First(&config)

	return config, nil
}
