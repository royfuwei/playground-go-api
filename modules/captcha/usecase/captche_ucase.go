package usecase

import (
	"errors"
	"fmt"
	"playground-go-api/config"
	"playground-go-api/domain"
	"playground-go-api/domain/category"
	"playground-go-api/domain/errcode"
	"playground-go-api/infrastructures/tools"
	"time"

	"github.com/golang/glog"
	"github.com/nyaruka/phonenumbers"
	"go.mongodb.org/mongo-driver/mongo"
)

type captchaUsecase struct {
	captchaSvc       domain.CaptchaService
	captchaRepo      domain.CaptchaRepository
	usersUcase       domain.UsersUsecase
	smsSvc           domain.SmsService
	jwtSvc           domain.JwtService
	smsUcase         domain.SmsUsecase
	expiresTime      time.Duration
	captchaNumberLen int
}

func NewCaptchaUsecase(
	captchaSvc domain.CaptchaService,
	captchaRepo domain.CaptchaRepository,
	usersUcase domain.UsersUsecase,
	smsSvc domain.SmsService,
	jwtSvc domain.JwtService,
	smsUcase domain.SmsUsecase,
) domain.CaptchaUsecase {
	return &captchaUsecase{
		captchaSvc:       captchaSvc,
		captchaRepo:      captchaRepo,
		usersUcase:       usersUcase,
		smsSvc:           smsSvc,
		jwtSvc:           jwtSvc,
		smsUcase:         smsUcase,
		expiresTime:      time.Duration(config.Cfgs.CaptchaExpired) * time.Second,
		captchaNumberLen: 6,
	}
}

func (ucase *captchaUsecase) VaildateBaseJwtId(id string) (*domain.BaseCaptcha, *domain.UCaseErr) {
	if id == "" {
		err := errors.New("jwt id is not exist")
		uCaseErr := tools.NewUCaseErr(category.Captcha, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	captcha, err := ucase.captchaRepo.FindBaseCaptchaById(id)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Captcha, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	return captcha, nil
}

// libphoneParse: country_code:886 national_number:966000996
// RFC3966: tel:+886-966-000-996
// INTERNATIONAL: +886 966 000 996
// NATIONAL: 0966 000 996
// E164: +886966000996
// 寄出sms 驗證碼
// 1. gen captcha code
// 2. save captcha db
// 3. gen sms message
// 4. send sms
func (ucase *captchaUsecase) SmsCaptchaSend(data *domain.ReqSmsCaptchaSendDTO) (*domain.ResSmsCaptchaSendDTO, *domain.UCaseErr) {
	countryCode := data.TelephoneRegion
	libphoneParse, err := phonenumbers.Parse(data.Telephone, countryCode)
	if err != nil {
		glog.Error(err)
	}
	telephone := fmt.Sprintf("%v%v", *libphoneParse.CountryCode, *libphoneParse.NationalNumber)

	captcha := ucase.captchaSvc.GenCaptchaNumber(ucase.captchaNumberLen)
	identifier := ucase.captchaSvc.GenCaptchaUpString(4)
	msg, err := ucase.smsSvc.GenSmsMsgByCountryCode(countryCode, captcha, identifier, ucase.expiresTime)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Captcha, errcode.Default, err, nil)
		return nil, uCaseErr
	}

	if _, err := ucase.smsUcase.SendMsgToTelephone(telephone, msg); err != nil {
		uCaseErr := tools.NewUCaseErr(category.Captcha, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	expiresTime := ucase.expiresTime
	insertData := &domain.SmsCaptcha{
		Identifier:      identifier,
		TelephoneRegion: countryCode,
		Telephone:       telephone,
		Message:         msg,
		Captcha:         captcha,
	}
	smsCaptcha, err := ucase.captchaRepo.AddSmsCaptcha(insertData, expiresTime)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Captcha, errcode.Default, err, nil)
		return nil, uCaseErr
	}

	return &domain.ResSmsCaptchaSendDTO{
		TelephoneRegion: smsCaptcha.TelephoneRegion,
		Identifier:      smsCaptcha.Identifier,
		Telephone:       smsCaptcha.Telephone,
		ExpiryDate:      smsCaptcha.ExpiryDate,
		ExpiresTime:     int(expiresTime.Seconds()),
	}, nil
}

func (ucase *captchaUsecase) SmsCaptchaValidate(data *domain.ReqSmsCaptchaValidateDTO) (*domain.ResSmsCaptchaValidateDTO, *domain.UCaseErr) {
	countryCode := data.TelephoneRegion
	libphoneParse, err := phonenumbers.Parse(data.Telephone, countryCode)
	if err != nil {
		glog.Error(err)
	}
	telephone := fmt.Sprintf("%v%v", *libphoneParse.CountryCode, *libphoneParse.NationalNumber)

	smsCaptcha, err := ucase.captchaRepo.FindOneSmsCaptcha(data.Captcha, data.Identifier, telephone, countryCode)
	result := &domain.ResSmsCaptchaValidateDTO{}
	fmt.Println(err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			result.Failed(domain.ErrUnauthorized)
			uCaseErr := tools.NewUCaseErr(category.Captcha, errcode.Default, domain.ErrUnauthorized, nil)
			return result, uCaseErr
		}
		uCaseErr := tools.NewUCaseErr(category.Captcha, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	ucase.captchaRepo.DeleteById(smsCaptcha.ID.Hex())

	baseCaptcha, err := ucase.captchaRepo.AddBaseCaptcha(3600 * time.Second)

	jwtId := baseCaptcha.ID.Hex()
	userData := &domain.UserData{
		Telephone:       telephone,
		TelephoneRegion: countryCode,
	}
	expiresAt, token, err := ucase.jwtSvc.JwtSign(ucase.expiresTime, userData, &jwtId)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return nil, uCaseErr
	}

	result.Success(token, expiresAt)
	return result, nil
}
