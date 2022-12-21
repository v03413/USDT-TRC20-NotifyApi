package main

import (
	"USDT-TRC20-NotifyApi/model"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const (
	TimeOut           = time.Second * 5
	TronScanApi       = "https://apilist.tronscan.org/"
	MaxExpireTime     = 1200 * 3 * 24
	MinExpireTime     = 180
	DefaultExpireTime = 600
)

var dbPath = getWd() + "/main.db"

func main() {
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
						go toNotify(trade, itm.TradeHash)
					}
				}
			}
		}(v.Address)
	}
}

func init() {
	_, err := os.Stat(dbPath)
	if err != nil {

		log.Fatalln("数据文件丢失，请尝试重新安装！")
	}
}
