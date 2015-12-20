package switcher

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// var redirect_url string = "http://develop.zsgd.com:11003/index.html"
// var redirect_url string = "http://develop.zsgd.com:11002/votes/index.html?vote_id=11"

type Xl map[string]func(http.ResponseWriter, *http.Request) (string, interface{})

type AccessInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

type UserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
}

/*
var (
	zsgdAppId     = "wxee6284b0c702c21e"
	zsgdAppSecret = "86707f6633745e51c4639f11c82abed2"
	jbntAppId = "wx8e3e7383c989c94a"
	jbntAppSecret = "6547db3da6a229a659bda118f48f8038"
)
*/

func Dispatch(db *sql.DB) Xl {
	return Xl{
		"check": func(rw http.ResponseWriter, r *http.Request) (string, interface{}) {
			var nickname string
			log.Println(GetParameter(r, "openid"))
			err := db.QueryRow("select nickname from register where openid = ?", GetParameter(r, "openid")).Scan(&nickname)
			perror(err, "查询签章信息失败")

			return "获取用户注册信息成功", nickname
		},

		"auth": func(rw http.ResponseWriter, r *http.Request) (string, interface{}) {
			defer func() {
				if err := recover(); err != nil {
					// http.Redirect(rw, r, redirect_url, http.StatusSeeOther)
					panic("重定向")
				}
			}()

			state := GetParameter(r, "state")
			var weixinStr string
			var r_url string
			err := db.QueryRow("select weixin, url from project where key = ?", state).Scan(
				&weixinStr, &r_url,
			)
			perror(err, "没有该工程")
			if weixinStr == "" {
				weixinStr = "zsgd93"
			}
			log.Println("weixinStr", weixinStr)
			code := GetParameter(r, "code")
			var appid string
			var appsecret string
			err = db.QueryRow("select appid, appsecret from weixin where weixin = ?", weixinStr).Scan(
				&appid, &appsecret,
			)
			perror(err, "无法获取微信信息")
			accessInfo := getAccessInfo(code, appid, appsecret)
			openid := accessInfo.Openid
			// log.Println(accessInfo) // DEL

			userInfo := getUserInfo(accessInfo.AccessToken, openid)
			log.Println(userInfo) // DEL

			db.Exec("insert into register values(?,?)", openid, userInfo.Nickname)

			log.Println(r_url)
			prefix := "?"
			if strings.Index(r_url, "?") > -1 {
				prefix = "&"
			}
			if ifSubscribe(db, weixinStr, openid) {
				http.Redirect(rw, r, r_url+prefix+"openid="+openid+"&nickname="+userInfo.Nickname+"&headimgurl="+userInfo.Headimgurl, http.StatusSeeOther)
			} else {
				http.Redirect(rw, r, "http://develop.zsgd.com:11010/error.html", http.StatusSeeOther)
			}
			return "正常重定向", nil
		},
	}
}

func getAccessInfo(code string, appid string, appsecret string) *AccessInfo {
	log.Println(appid, " => ", appsecret)
	u := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appid + "&secret=" + appsecret + "&code=" + code + "&grant_type=authorization_code"
	v := getWeixinInfo(u, new(AccessInfo)).(*AccessInfo)
	log.Println("accessInfo:", v)
	return v
}

func getUserInfo(accessToken string, openid string) *UserInfo {
	u := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openid + "&lang=zh_CN"
	log.Println(u)
	return getWeixinInfo(u, new(UserInfo)).(*UserInfo)
}

func getWeixinInfo(u string, inf interface{}) interface{} {
	resp, err := http.Get(u)
	defer resp.Body.Close()
	perror(err, "连接失败")

	body, err := ioutil.ReadAll(resp.Body)
	perror(err, "无法读取body")

	err = json.Unmarshal(body, inf)
	perror(err, "无法从json获取信息")

	return inf
}

func perror(e error, errMsg string) {
	if e != nil {
		log.Println(e)
		panic(errMsg)
	}
}

func GetParameter(r *http.Request, key string) string {
	s := r.URL.Query().Get(key)
	if s == "" {
		panic("没有参数" + key)
	}
	return s
}

// 判断是否是公众号用户
type SubscribeInfo struct {
	Subscribe int `json:"subscribe"`
}

func ifSubscribe(db *sql.DB, weixin string, openid string) bool {
	if weixin == "zsgd93" {
		return true
	}
	var accessToken string
	err := db.QueryRow("select access_token from weixin where weixin = ?", weixin).Scan(&accessToken)
	perror(err, "无法获取access_token")
	v := getWeixinInfo("https://api.weixin.qq.com/cgi-bin/user/info?access_token="+accessToken+"&openid="+openid+"&lang=zh_CN", new(SubscribeInfo)).(*SubscribeInfo)
	return v.Subscribe == 1
}
