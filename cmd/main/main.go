package main

import (
	"USDT-TRC20-NotifyApi/log"
	"USDT-TRC20-NotifyApi/model"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const (
	Version           = "0.1"
	LogLevel          = "debug"
	LogOutput         = "/var/log/main.log"
	TimeOut           = time.Second * 5
	TronScanApi       = "https://apilist.tronscan.org/"
	MaxExpireTime     = 1200
	MinExpireTime     = 180
	DefaultExpireTime = 600
)

var dbPath = getWd() + "/main.db"

func init() {
	log.Init(LogLevel, LogOutput)

	_, err := os.Stat(dbPath)
	if err != nil {

		log.Fatal("数据文件丢失，请尝试重新安装！")
	}
}

func main() {
	log.Println(fmt.Sprintf("当前版本：%s", Version))
	log.Println("开源地址：https://github.com/v03413/USDT-TRC20-NotifyApi")

	go heartbeat()

	{
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		runtime.GC()
	}
}

func heartbeat() {
	go start()

	for range time.Tick(time.Second * 3) {
		go toNotifyRetry()
		go dealWith()
	}
}

func dealWith() {
	var rows []model.Trade
	db.Distinct("address").Where("state = ? and expire_time >= ?", model.TradeStateUnconfirmed, time.Now()).Find(&rows)
	for _, v := range rows {
		go func(address string) {
			list, err := getTransferList(address)
			if err != nil {

				log.Println(err.Error())
				return
			}
			if len(list) == 0 {

				return
			}

			// 获取交易订单
			var trades []model.Trade
			db.Where("address = ? and state = ? and expire_time >= ?", address, model.TradeStateUnconfirmed, time.Now()).Find(&trades)
			for _, trade := range trades {
				for _, itm := range list {
					if trade.Amount == itm.Amount && trade.ExpireTime.UnixMilli() >= itm.Time.UnixMilli() {
						go func() {
							err := toNotify(trade, itm.TradeHash)
							if err != nil {

								log.Println(err.Error())
							}
						}()
					}
				}
			}
		}(v.Address)
	}
}
