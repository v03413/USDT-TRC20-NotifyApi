package main

import (
	"USDT-TRC20-NotifyApi/model"
	"math"
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

		return setNotifyFailed(trade)
	}

	if resp.StatusCode == http.StatusOK {

		return setNotifySucc(trade)
	}

	return setNotifyFailed(trade)
}

func toNotifyRetry() {
	var trades []model.Trade
	db.Where("state = ?", model.TradeStateNotifyFailed, time.Now()).Find(&trades)
	for _, trade := range trades {
		s := TimeOut.Seconds() * math.Pow(2, float64(trade.NotifyRetry))
		LastNotifyTime := trade.NotifyTime.Time.Unix()
		if time.Now().Unix() < int64(s)+LastNotifyTime {

			continue
		}

		go toNotify(trade, "")
	}
}

func setNotifySucc(trade model.Trade) error {
	trade.State = model.TradeStateComplete
	trade.NotifyRetry += 1
	trade.NotifyTime.Time = time.Now()
	trade.NotifyTime.Valid = true

	db.Save(&trade)

	return db.Error
}

func setNotifyFailed(trade model.Trade) error {
	trade.State = model.TradeStateNotifyFailed
	trade.NotifyRetry += 1
	trade.NotifyTime.Time = time.Now()
	trade.NotifyTime.Valid = true

	db.Save(&trade)

	return db.Error
}
