package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"playground-go-api/config"
	"playground-go-api/domain"

	"github.com/golang/glog"
)

type twsmsSmsUsecase struct {
	account  string
	password string
	encoding string
	url      string
}

// http://api.twsms.com/json/sms_send.php?username=test&password=1234&mobile=0911222333&message
// username 使用者帳號
// password 密碼
// mobile 電話 09xxxxxxxx(台灣門號) or 國碼+門號 (國際門號)
// message 可使用 UTF8 碼或 BIG5 碼（BIG5 碼部分特殊字會變成?號）中文最長 70 個字，扣 1 點數
// sendtime 格式：YYYYMMDDHHII （請使用 24 小時制）預約時間，例如 201203121830
func NewTwsmsSmsUsecase() domain.SmsUsecase {
	return &twsmsSmsUsecase{
		account:  config.Cfgs.SmsAccount,
		password: config.Cfgs.SmsPassword,
		encoding: "urlencode_utf8",
		url:      "http://api.twsms.com/json/sms_send.php",
	}
}

// SendMsgToTelephone Send Message to telephone
func (ucase *twsmsSmsUsecase) SendMsgToTelephone(tel string, msg string) (*domain.SendMsgToTelDTO, error) {
	url, err := url.Parse(ucase.url)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	query := url.Query()
	query.Add("username", ucase.account)
	query.Add("password", ucase.password)
	query.Add("mobile", tel)
	query.Add("message", msg)
	url.RawQuery = query.Encode()

	res, err := http.Get(url.String())
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("ioutil.ReadAll error: %v", err)
		return nil, err
	}
	if res.StatusCode >= 400 {
		err = fmt.Errorf("statusCode: %d, status: %s \n", res.StatusCode, res.Status)
		return nil, err
	}

	var twSmsSendApiRes domain.TwSmsSendApiRes
	if err := json.Unmarshal(body, &twSmsSendApiRes); err != nil {
		return nil, err
	}
	return ucase.resSendMsgToTelDTO(twSmsSendApiRes)
}

// Code
//
// 00000 完成 00001 狀態尚未回復 00010 帳號或密碼格式錯誤 00011 帳號錯誤 00012 密碼錯誤
// 00020 通數不足 00030 IP 無使用權限 00040 帳號已停用 00050 sendtime 格式錯誤 00060 expirytime 格式錯誤
// 00100 手機號碼格式錯誤 00110 沒有簡訊內容 00120 長簡訊不支援國際門號 00130 簡訊內容超過長度 00140 drurl 格式錯誤
// 00150 sendtime 預約的時間已經超過 00160 drurl 帶入的網址無法連線 00300 找不到 msgid 00310 預約尚未送出 00400 找不到 snumber 辨識碼
// 00410 沒有任何 mo 資料 00420 smsQuery 指定查詢的格式錯誤 00430 moQuery 指定查詢的格式錯誤 99998 資料處理異常，請重新發送 99999 系統錯誤，請通知系統廠商
func (ucase *twsmsSmsUsecase) resSendMsgToTelDTO(twSmsSendApiRes domain.TwSmsSendApiRes) (*domain.SendMsgToTelDTO, error) {
	result := &domain.SendMsgToTelDTO{
		Message: twSmsSendApiRes.Text,
		MsgId:   string(twSmsSendApiRes.MsgId),
		Code:    twSmsSendApiRes.Code,
	}
	if twSmsSendApiRes.Code != "00000" {
		err := fmt.Errorf("sms errorCode: %v, msg: %v", twSmsSendApiRes.Code, twSmsSendApiRes.Text)
		result.Failed()
		return result, err
	}
	result.Success()
	return result, nil
}
