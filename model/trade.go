package model

import (
	"database/sql"
	"time"
)

const (
	TradeStateUnconfirmed  = -1
	TradeStateNotifyFailed = 0
	TradeStateComplete     = 1
)

type Trade struct {
	Id            uint `gorm:"primaryKey"`
	State         int
	Amount        float64
	Address       string
	TradeNo       string `gorm:"unique"`
	TradeHash     sql.NullString
	NotifyUrl     string
	NotifyRetry   int
	NotifyTime    sql.NullTime
	ExpireTime    time.Time
	CreateTime    time.Time    `gorm:"autoCreateTime"`
	UpdateTime    sql.NullTime `gorm:"autoUpdateTime"`
	transfersList interface{}
}
