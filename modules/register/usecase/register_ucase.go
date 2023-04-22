package usecase

import (
	"errors"
	"fmt"
	"playground-go-api/domain"
	"playground-go-api/domain/category"
	"playground-go-api/domain/errcode"
	usersfield "playground-go-api/domain/users/field"
	"playground-go-api/infrastructures/tools"

	"github.com/golang/glog"
	"github.com/nyaruka/phonenumbers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type registerUcase struct {
	usersUcase   domain.UsersUsecase
	authUcase    domain.AuthUsecase
	usersRepo    domain.UsersRepository
	captchaUcase domain.CaptchaUsecase
}

func NewRegisterUsecase(
	usersUcase domain.UsersUsecase,
	authUcase domain.AuthUsecase,
	usersRepo domain.UsersRepository,
	captchaUcase domain.CaptchaUsecase,
) domain.RegisterUsecase {
	return &registerUcase{
		usersUcase,
		authUcase,
		usersRepo,
		captchaUcase,
	}
}

func (ucase *registerUcase) RegisterUserByTelephone(claims *domain.Claims, data *domain.ReqCreateUserByTelephone) (*domain.ResLoginTokenDTO, *domain.UCaseErr) {
	_, uCaseErr := ucase.captchaUcase.VaildateBaseJwtId(claims.Id)
	if uCaseErr != nil {
		return nil, uCaseErr
	}
	if claims.Telephone != data.Telephone || claims.TelephoneRegion != data.TelephoneRegion {
		// uCaseErr := tools.NewUCaseErr(category.Register, errcode.Default, domain.ErrUnauthorized, nil)
		// return nil, uCaseErr

	}
	userData, uCaseErr := ucase.usersUcase.CreateUserByTelephone(data)
	if uCaseErr != nil {
		return nil, uCaseErr
	}

	result, uCaseErr := ucase.authUcase.GenResLoginTokenDTO(userData)
	if uCaseErr != nil {
		return nil, uCaseErr
	}
	return result, nil
}

func (ucase *registerUcase) RegisterSmsCaptchaSend(data *domain.ReqSmsCaptchaSendDTO) (*domain.ResSmsCaptchaSendDTO, *domain.UCaseErr) {
	countryCode := data.TelephoneRegion
	libphoneParse, err := phonenumbers.Parse(data.Telephone, countryCode)
	if err != nil {
		glog.Error(err)
	}
	telephone := fmt.Sprintf("%v%v", *libphoneParse.CountryCode, *libphoneParse.NationalNumber)

	filter := bson.M{
		string(usersfield.Telephone):       telephone,
		string(usersfield.TelephoneRegion): countryCode,
	}

	user, err := ucase.usersRepo.FindOneByFilter(filter)
	if err != nil && err != mongo.ErrNoDocuments {
		uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	if user != nil {
		err := errors.New("telephone number is exist user")
		uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
		return nil, uCaseErr
	}

	return ucase.captchaUcase.SmsCaptchaSend(data)
}
