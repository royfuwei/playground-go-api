package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST 請求sms 辨識碼
type ReqSmsCaptchaSendDTO struct {
	// 國家地區代碼
	TelephoneRegion string `json:"telephoneRegion" bson:"telephoneRegion" binding:"required" example:"TW"`
	// 電話號碼
	Telephone string `json:"telephone" bson:"telephone" binding:"required"`
}

// POST 回傳成功 請求sms 辨識碼
type ResSmsCaptchaSendDTO struct {
	// 辨識碼
	Identifier string `json:"identifier" bson:"identifier"`
	// 國家地區代碼
	TelephoneRegion string `json:"telephoneRegion" bson:"telephoneRegion" example:"TW"`
	// 電話號碼
	Telephone string `json:"telephone" bson:"telephone"`
	// 存在時長(秒)
	ExpiresTime int `json:"expiresTime" bson:"expiresTime"`
	// 截止時間
	ExpiryDate primitive.DateTime `json:"expiryDate" bson:"expiryDate" swaggertype:"string"`
}

// POST 請求驗證sms captcha
type ReqSmsCaptchaValidateDTO struct {
	// 識別碼
	Identifier string `json:"identifier" bson:"identifier" binding:"required"`
	// 手機取得的驗證碼
	Captcha string `json:"captcha" bson:"captcha" binding:"required"`
	// 國家地區代碼
	TelephoneRegion string `json:"telephoneRegion" bson:"telephoneRegion" binding:"required" example:"TW"`
	// 電話號碼
	Telephone string `json:"telephone" bson:"telephone" binding:"required"`
}

type ResSmsCaptchaValidateDTO struct {
	// 回傳狀態
	Process string `json:"process"`
	// 回傳訊息
	Message     string `json:"message"`
	AccessToken string `json:"accessToken,omitempty"`
	// ExpiresAt   int64  `json"expiresAt,omitempty"`
	AccessTokenExp int64 `json:"accessTokenExp,omitempty"`
}

type BaseCaptcha struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ExpiryDate primitive.DateTime `json:"expiryDate" bson:"expiryDate"`
}

type SmsCaptcha struct {
	BaseCaptcha `bson:",inline"`
	// 手機取得的驗證碼
	Captcha string `json:"captcha" bson:"captcha" binding:"required"`
	// 識別碼
	Identifier string `json:"identifier" bson:"identifier" binding:"required"`
	// 國家地區代碼
	TelephoneRegion string `json:"telephoneRegion" bson:"telephoneRegion" binding:"required" example:"TW"`
	// 電話號碼
	Telephone string `json:"telephone" bson:"telephone" binding:"required"`
	// 存在時長
	// ExpiresTime time.Duration `json:"expiresTime" bson:"expiresTime"`
	// 回傳訊息
	Message string `json:"message" bson:"message"`
}

func (r *ResSmsCaptchaValidateDTO) Success(token string, expiresAt int64) {
	r.Process = "success"
	r.Message = "success validate sms captcha"
	r.AccessToken = token
	r.AccessTokenExp = expiresAt
}

func (r *ResSmsCaptchaValidateDTO) Failed(err error) {
	r.Process = "failed"
	r.Message = err.Error()
	// r.ExpiresAt = 0
}

type CaptchaRepository interface {
	AddSmsCaptcha(data *SmsCaptcha, expiresTime time.Duration) (*SmsCaptcha, error)
	AddBaseCaptcha(expiresTime time.Duration) (*BaseCaptcha, error)
	FindOneSmsCaptcha(captcha, identifier, telephone, telephoneRegion string) (*SmsCaptcha, error)
	FindBaseCaptchaById(id string) (baseCaptcha *BaseCaptcha, err error)
	FindSmsCaptchaById(id string) (smsCaptcha *SmsCaptcha, err error)
	DeleteById(id string) (bool, error)
}

type CaptchaService interface {
	GenCaptchaNumber(length int) string
	GenCaptchaUpString(length int) string
}

type CaptchaUsecase interface {
	SmsCaptchaSend(data *ReqSmsCaptchaSendDTO) (*ResSmsCaptchaSendDTO, *UCaseErr)
	SmsCaptchaValidate(data *ReqSmsCaptchaValidateDTO) (*ResSmsCaptchaValidateDTO, *UCaseErr)
	VaildateBaseJwtId(id string) (*BaseCaptcha, *UCaseErr)
}
