package domain

import "time"

type TwSmsSendApiRes struct {
	// Code
	//
	// 00000 完成 00001 狀態尚未回復 00010 帳號或密碼格式錯誤 00011 帳號錯誤 00012 密碼錯誤
	// 00020 通數不足 00030 IP 無使用權限 00040 帳號已停用 00050 sendtime 格式錯誤 00060 expirytime 格式錯誤
	// 00100 手機號碼格式錯誤 00110 沒有簡訊內容 00120 長簡訊不支援國際門號 00130 簡訊內容超過長度 00140 drurl 格式錯誤
	// 00150 sendtime 預約的時間已經超過 00160 drurl 帶入的網址無法連線 00300 找不到 msgid 00310 預約尚未送出 00400 找不到 snumber 辨識碼
	// 00410 沒有任何 mo 資料 00420 smsQuery 指定查詢的格式錯誤 00430 moQuery 指定查詢的格式錯誤 99998 資料處理異常，請重新發送 99999 系統錯誤，請通知系統廠商
	Code string `json:"code"`
	Text string `json:"text"`
	// 發送簡訊後 API 給予的序號(長度為 8~11)
	MsgId int `json:"msgId"`
}

type SendMsgToTelDTO struct {
	process string `json:"process"`
	Message string `json:"message"`
	MsgId   string `json:"messageId,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (dto *SendMsgToTelDTO) Success() {
	dto.process = "success"
}
func (dto *SendMsgToTelDTO) Failed() {
	dto.process = "failed"
}

type SmsService interface {
	GenSmsMsgByCountryCode(countryCode, captcha, identifier string, expiresTime time.Duration) (string, error)
}

type SmsUsecase interface {
	// SendMsgToTelephone Send Message to telephone, 回傳http response byte, error
	SendMsgToTelephone(tel string, msg string) (*SendMsgToTelDTO, error)
}
