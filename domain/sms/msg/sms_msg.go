package smsMsg

type SmsMessage string

const (
	TW SmsMessage = "請確認識別碼為%s，輸入簡訊驗證碼%s，驗證碼%s後自動失效！"
)

var CountryCodeSmsMessage = map[string]SmsMessage{
	"TW": TW,
}
