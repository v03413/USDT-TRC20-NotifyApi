package main

import (
	"USDT-TRC20-NotifyApi/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
	"strconv"
	"time"
)

var db *gorm.DB

type respJson struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	TradeNo string `json:"trade_no"`
}

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("./main.db"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "mm_",
			SingularTable: true,
		},
	})
	if err != nil {

		log.Fatalf("数据库错误：%s\n", err.Error())
	}
}

func start() {
	r := gin.Default()
	r.GET("/", index)
	r.GET("/submit", submit)
	err := r.Run("127.0.0.1:8080")
	if err != nil {

		log.Fatalln(err.Error())
	}
}

func index(c *gin.Context) {
	var url = fmt.Sprintf("%sapi/system/status", TronScanApi)
	get, err := httpGet(url)
	if err != nil {
		c.String(200, fmt.Sprintf("运行异常：%s(请检查网络！)", err.Error()))
		return
	}

	if get.StatusCode != http.StatusOK {
		c.String(200, fmt.Sprintf("运行异常：错误的HTTP状态码(%d)", get.StatusCode))
		return
	}

	c.String(200, fmt.Sprintf("[运行正常] %s", time.Now().Format("2006-01-02 15:04:05")))
}

func submit(c *gin.Context) {
	tradeNo := generateTradeNo()
	address := c.Query("address")
	if address[:1] != "T" || len(address) != 34 {
		c.JSON(200, respJson{
			Code: 202,
			Msg:  "Tron钱包地址应该是一个T开头的34位长度字符串",
		})
		return
	}

	amount, err := strconv.ParseFloat(c.Query("amount"), 64)
	if err != nil {
		c.JSON(200, respJson{
			Code: 201,
			Msg:  fmt.Sprintf("amount 参数错误%s", err.Error()),
		})
		return
	}

	expire, err := strconv.ParseUint(c.Query("expire"), 10, 32)
	if err != nil {
		c.JSON(200, respJson{
			Code: 201,
			Msg:  fmt.Sprintf("expire 参数错误(只能是180～1200之间的整数)：%s", err.Error()),
		})
		return
	}
	if expire < MinExpireTime || expire > MaxExpireTime {

		expire = DefaultExpireTime
	}

	var row model.Trade
	db.Where("amount = ? and address = ? and expire_time >= ?", amount, address, time.Now()).Find(&row)
	if row.Id != 0 {
		c.JSON(200, respJson{
			Code: 203,
			Msg:  "有效期内已存在同地址金额的监控记录，请更换钱包或金额再重试！",
		})
		return
	}

	db.Create(&model.Trade{
		State:      model.TradeStateUnconfirmed,
		Amount:     amount,
		Address:    address,
		TradeNo:    tradeNo,
		NotifyUrl:  c.Query("notify_url"),
		ExpireTime: time.Now().Add(time.Second * time.Duration(expire)),
	})

	c.JSON(200, respJson{Code: 200, Msg: "success", TradeNo: tradeNo})
}
