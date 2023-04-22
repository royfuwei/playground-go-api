package service

import (
	"fmt"
	"playground-go-api/domain"
	smsMsg "playground-go-api/domain/sms/msg"
	"strings"
	"time"
)

type smsService struct {
}

// 產生需要的SMS message
func NewSmsService() domain.SmsService {
	return &smsService{}
}

func (svc *smsService) GenSmsMsgByCountryCode(countryCode, captcha, identifier string, expiresTime time.Duration) (string, error) {
	fmtMsg := smsMsg.CountryCodeSmsMessage[strings.ToUpper(countryCode)]
	expiresTimeStr := svc.expiresTimeString(countryCode, expiresTime)
	if len(fmtMsg) == 0 {
		return "", domain.ErrCountryCodeNotExist
	}
	msg := fmt.Sprintf(string(fmtMsg), identifier, captcha, expiresTimeStr)
	return msg, nil
}

// 回傳expiresTime 顯示的訊息
func (svc *smsService) expiresTimeString(countryCode string, expiresTime time.Duration) string {
	var result string
	sec := int(expiresTime.Seconds())
	switch {
	case sec < 60:
		result = fmt.Sprintf("%v秒", sec)
	default:
		min := int(sec / 60)
		result = fmt.Sprintf("%v分鐘", min)

	}
	return result
}
