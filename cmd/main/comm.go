package main

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Transfer struct {
	Amount    float64
	Time      time.Time
	TradeHash string
}

func getWd() string {
	path, err := os.Getwd()
	if err != nil {

		log.Fatalln(err.Error())
	}

	return path
}

func generateTradeNo() string {
	rand.Seed(time.Now().UnixMilli())

	var tradeNo = time.Now().Format("20060102150405")
	var randInt = 100000 + rand.Intn(999999)

	return fmt.Sprintf("%s%d", tradeNo, randInt)
}

func getTransferList(address string) ([]Transfer, error) {
	var list []Transfer
	var startTime, _ = time.ParseDuration(fmt.Sprintf("-%ds", MaxExpireTime))
	var api = fmt.Sprintf("%sapi/token_trc20/transfers?limit=300&start=0&direction"+
		"&relatedAddress=%s&start_timestamp=%d", TronScanApi, address, time.Now().Add(startTime).UnixMilli())
	resp, err := httpGet(api)
	if err != nil {

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {

		return nil, errors.New(fmt.Sprintf("getTransferList状态码错误：%d", resp.StatusCode))
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, err
	}

	raw := gjson.ParseBytes(all)
	raw.Get("token_transfers").ForEach(func(key, value gjson.Result) bool {
		list = append(list, Transfer{
			Amount:    value.Get("quant").Float() / 1000000,
			TradeHash: value.Get("transaction_id").String(),
			Time:      time.UnixMilli(value.Get("block_ts").Int()),
		})

		return true
	})

	return list, nil
}

func httpGet(url string) (resp *http.Response, err error) {
	client := http.Client{Timeout: TimeOut}

	return client.Get(url)
}
