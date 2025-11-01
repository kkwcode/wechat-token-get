package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// AccessTokenResponse 微信access_token响应结构
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// AccessTokenResult 返回给客户端的access_token结果
type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Timestamp   int64  `json:"timestamp"`
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
}

// AccessTokenHandler access_token接口处理器
func AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	if r.Method == http.MethodGet {
		tokenResult, err := getAccessToken()
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = tokenResult
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持，请使用GET方法", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// getAccessToken 获取微信access_token
func getAccessToken() (*AccessTokenResult, error) {
	// 从环境变量获取appid和secret
	appid := os.Getenv("WECHAT_APPID")
	secret := os.Getenv("WECHAT_SECRET")

	// 检查环境变量是否配置
	if appid == "" || secret == "" {
		return &AccessTokenResult{
			Success: false,
			Message: "请配置微信公众平台的WECHAT_APPID和WECHAT_SECRET环境变量",
		}, nil
	}

	// 调用微信API获取access_token
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, secret)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var wechatResp AccessTokenResponse
	err = json.Unmarshal(body, &wechatResp)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查微信API返回的错误
	if wechatResp.ErrCode != 0 {
		return &AccessTokenResult{
			Success: false,
			Message: fmt.Sprintf("微信API错误: %s (错误码: %d)", wechatResp.ErrMsg, wechatResp.ErrCode),
		}, nil
	}

	// 返回成功的access_token结果
	return &AccessTokenResult{
		AccessToken: wechatResp.AccessToken,
		ExpiresIn:   wechatResp.ExpiresIn,
		Timestamp:   time.Now().Unix(),
		Success:     true,
		Message:     "获取access_token成功",
	}, nil
}