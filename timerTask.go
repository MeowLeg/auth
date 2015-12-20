package main

import (
	"database/sql"
	"encoding/json"
	// "flag"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AccessInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func timerTask() {
	db := connectDB("./middle.db")
	defer db.Close()

	rows, err := db.Query("select distinct weixin from weixin")
	perr(err, "无法读取微信号")
	var weixins []string
	for rows.Next() {
		var w string
		rows.Scan(&w)
		weixins = append(weixins, w)
	}

	for {
		for _, w := range weixins {
			getAccessToken(db, w)
		}
		t := time.NewTimer(time.Hour * 1)
		<-t.C
	}
}

func getAccessToken(db *sql.DB, weixin string) {
	var (
		appid     string
		appsecret string
	)
	err := db.QueryRow("select appid, appsecret from weixin where weixin = ?", weixin).Scan(
		&appid, &appsecret,
	)
	perr(err, "无法读取微信信息")
	log.Println(appid) // 测试

	resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + appsecret)
	defer resp.Body.Close()
	perr(err, "连接失败")

	body, err := ioutil.ReadAll(resp.Body)
	perr(err, "无法读取body")

	var inf AccessInfo
	err = json.Unmarshal(body, &inf)
	perr(err, "无法从json获取信息")
	log.Println(inf) // 测试

	_, err = db.Exec("update weixin set access_token = ? where weixin = ?", inf.AccessToken, weixin)
	perr(err, "更新微信数据失败")
}

func connectDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	perr(err, "数据库连接失败")
	return db
}

func perr(err error, msg string) {
	if err != nil {
		log.Println(err)
		panic(msg)
	}
}

func main() {
	timerTask()
}
