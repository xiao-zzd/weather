package main
import (
	"encoding/json"
   "database/sql"
   _ "github.com/go-sql-driver/mysql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)
//设置结构体用于存放数据
type Forecast struct {
	Date     string `json:"date"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Ymd      string `json:"ymd"`
	Week     string `json:"week"`
	Sunrise  string `json:"sunrise"`
	Sunset   string `json:"sunset"`
	Aqi      int    `json:"aqi"`
	Fx       string `json:"fx"`
	Fl       string `json:"fl"`
	Type     string `json:"type"`
	Notice   string `json:"notice"`
}

type WeatherResponse struct {
	Data struct {
		Forecast []Forecast `json:"forecast"`
	} `json:"data"`
}
var db *sql.DB
func main() {
	// 发起 HTTP 请求获取天气数据
	response, err := http.Get("http://t.weather.sojson.com/api/weather/city/101030100")
	if err != nil {
		log.Fatal("HTTP 请求失败：", err)
	}
	defer response.Body.Close()

	// 读取响应的 JSON 数据
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("读取响应数据失败：", err)
	}

	// 解析 JSON 数据
	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		log.Fatal("解析 JSON 数据失败：", err)
	}

	// 打印 forecast 中的 low 和 high 字段数据到终端
	for _, forecast := range weatherResp.Data.Forecast {
		fmt.Printf("Low: %s, High: %s\n", forecast.Low, forecast.High)
	}

   //初始化数据库
   initDB()

   for _, forecast := range weatherResp.Data.Forecast {
		_, err = db.Exec("INSERT INTO forecast (date, high, low, ymd, week, sunrise, sunset, aqi, fx, fl, wtype, notice) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
      forecast.Date, forecast.High, forecast.Low, forecast.Ymd, forecast.Week, forecast.Sunrise, forecast.Sunset, forecast.Aqi, forecast.Fx, forecast.Fl, forecast.Type, forecast.Notice)
		if err != nil {
			log.Fatal("插入数据失败：", err)
		}
}
}
func initDB() {
	var err error
    //创建与mysql的连接
    //username:password@tcp(127.0.0.1:3306)/database_name中的username、password和database_name替换为你的MySQL数据库的实际连接信息。
	db, err = sql.Open("mysql", "root:zhanzhaodong@tcp(127.0.0.1:3306)/zzd")
	if err != nil {
		log.Fatal(err)
	}
	// 创建表 ，如果表已经存在，则IF NOT EXISTS语句将防止重新创建该表。
	createTableSQL := `CREATE TABLE IF NOT EXISTS forecast (
		date VARCHAR(10),
		high VARCHAR(10),
		low VARCHAR(10),
		ymd VARCHAR(10),
		week VARCHAR(10),
		sunrise VARCHAR(10),
		sunset VARCHAR(10),
		aqi INT,
		fx VARCHAR(10),
		fl VARCHAR(10),
		wtype VARCHAR(10),
		notice TEXT
	)`
    //执行createTableSQL的创建表sql语法
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}