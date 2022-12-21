package main

import (
	"USDT-TRC20-NotifyApi/model"
	"net/http"
	"net/url"
	"time"
)

func toNotify(trade model.Trade, tradeHash string) error {
	if trade.TradeHash.String == "" {
		trade.TradeHash.String = tradeHash
		trade.TradeHash.Valid = true
		db.Save(&trade)
	}

	var data = make(url.Values)
	data.Add("hash", trade.TradeHash.String)
	data.Add("trade_no", trade.TradeNo)

	client := http.Client{Timeout: TimeOut}
	resp, err := client.PostForm(trade.NotifyUrl, data)
	if err != nil {

		return toNotifyFailed(trade)
	}

	if resp.StatusCode == http.StatusOK {

		return toNotifySucc(trade)
	}

	return toNotifyFailed(trade)
}

func toNotifySucc(trade model.Trade) error {
	trade.State = model.TradeStateComplete
	trade.NotifyRetry += 1
	trade.NotifyTime.Time = time.Now()
	trade.NotifyTime.Valid = true

	db.Save(&trade)

	return db.Error
}

func toNotifyFailed(trade model.Trade) error {
	trade.State = model.TradeStateNotifyFailed
	trade.NotifyRetry += 1
	trade.NotifyTime.Time = time.Now()
	trade.NotifyTime.Valid = true

	db.Save(&trade)

	return db.Error
}
