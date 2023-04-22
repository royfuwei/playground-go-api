package usecase

import (
	"fmt"
	"playground-go-api/config"
	"playground-go-api/domain"
	"playground-go-api/domain/category"
	"playground-go-api/domain/errcode"
	usersfield "playground-go-api/domain/users/field"
	"playground-go-api/infrastructures/tools"
	"time"

	"github.com/golang/glog"
	"github.com/nyaruka/phonenumbers"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type authUsecase struct {
	refreshTokensRepo domain.RefreshTokensRepository
	jwtSvc            domain.JwtService
	usersRepo         domain.UsersRepository
	usersSvc          domain.UsersService
	expiresTime       time.Duration
}

func NewAuthUsecase(
	refreshTokensRepo domain.RefreshTokensRepository,
	jwtSvc domain.JwtService,
	usersRepo domain.UsersRepository,
	usersSvc domain.UsersService,
) domain.AuthUsecase {
	expiresTime := time.Duration(config.Cfgs.AccessTokenExpired) * time.Second
	return &authUsecase{
		refreshTokensRepo: refreshTokensRepo,
		jwtSvc:            jwtSvc,
		usersRepo:         usersRepo,
		usersSvc:          usersSvc,
		expiresTime:       expiresTime,
	}
}

// 檢查使用者是否存在
func (ucase *authUsecase) checkUserExist(field string, value string) (*domain.User, *domain.UCaseErr) {
	filter := bson.M{
		field: value,
	}
	user, err := ucase.usersRepo.FindOneByFilter(filter)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	return user, nil
}

// 檢查使用者密碼
func (ucase *authUsecase) checkUserPassword(user *domain.User, password string) *domain.UCaseErr {
	/* check password */
	isPwd := ucase.usersSvc.CheckPassword(password, user)
	if !isPwd {
		err := fmt.Errorf("usersSvc.CheckPassword error")
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return uCaseErr
	}
	return nil
}

// 簽發access token
func (ucase *authUsecase) CreateAccessToken(userData *domain.UserData, jwtId *string) (expiresAt int64, token string, uCaseErr *domain.UCaseErr) {

	expiresAt, token, err := ucase.jwtSvc.JwtSign(ucase.expiresTime, userData, jwtId)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return expiresAt, token, uCaseErr
	}
	return expiresAt, token, nil
}

// create refresh token by userId
func (ucase *authUsecase) CreateRefreshToken(uid primitive.ObjectID) (string, *domain.UCaseErr) {
	uuidV4 := uuid.NewV4().String()
	data := &domain.RefreshToken{
		UserId:       uid,
		RefreshToken: uuidV4,
	}
	refreshToken, err := ucase.refreshTokensRepo.Add(data)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return uuidV4, uCaseErr
	}
	return refreshToken.RefreshToken, nil
}

// 用telephone 登入
func (ucase *authUsecase) LoginTelephone(data *domain.ReqLoginTelephoneDTO) (*domain.ResLoginTokenDTO, *domain.UCaseErr) {
	countryCode := data.TelephoneRegion
	libphoneParse, err := phonenumbers.Parse(data.Telephone, countryCode)
	if err != nil {
		glog.Error(err)
	}
	telephone := fmt.Sprintf("%v%v", *libphoneParse.CountryCode, *libphoneParse.NationalNumber)
	/* get user */
	user, uCaseErr := ucase.checkUserExist(string(usersfield.Telephone), telephone)
	if uCaseErr != nil {
		return nil, uCaseErr
	}
	/* check password */
	if uCaseErr = ucase.checkUserPassword(user, data.Password); uCaseErr != nil {
		return nil, uCaseErr
	}

	userData := &domain.UserData{
		ID:              user.ID,
		Account:         user.Account,
		UserName:        user.UserName,
		NickName:        user.NickName,
		Telephone:       user.Telephone,
		TelephoneRegion: user.TelephoneRegion,
		Roles:           user.Roles,
		ExpiryDate:      user.ExpiryDate,
		ModiTime:        user.ModiTime,
	}
	return ucase.GenResLoginTokenDTO(userData)
}

func (ucase *authUsecase) GenResLoginTokenDTO(userData *domain.UserData) (*domain.ResLoginTokenDTO, *domain.UCaseErr) {
	expiresAt, token, uCaseErr := ucase.CreateAccessToken(userData, nil)
	if uCaseErr != nil {
		return nil, uCaseErr
	}

	refreshToken, uCaseErr := ucase.CreateRefreshToken(userData.ID)
	if uCaseErr != nil {
		return nil, uCaseErr
	}

	result := &domain.ResLoginTokenDTO{
		AccessToken:    token,
		RefreshToken:   refreshToken,
		AccessTokenExp: int(expiresAt),
	}
	return result, nil
}

// 用account 登入
func (ucase *authUsecase) LoginAccount(data *domain.ReqLoginAccountDTO) (*domain.ResLoginTokenDTO, *domain.UCaseErr) {
	/* get user */
	user, uCaseErr := ucase.checkUserExist(string(usersfield.Account), data.Account)
	if uCaseErr != nil {
		return nil, uCaseErr
	}
	/* check password */
	if uCaseErr = ucase.checkUserPassword(user, data.Password); uCaseErr != nil {
		return nil, uCaseErr
	}

	userData := &domain.UserData{
		ID:              user.ID,
		Account:         user.Account,
		UserName:        user.UserName,
		NickName:        user.NickName,
		Telephone:       user.Telephone,
		TelephoneRegion: user.TelephoneRegion,
		Roles:           user.Roles,
		ExpiryDate:      user.ExpiryDate,
		ModiTime:        user.ModiTime,
	}
	return ucase.GenResLoginTokenDTO(userData)
}

func (ucase *authUsecase) AuthJwtVerify(token string) (*domain.Claims, *domain.UCaseErr) {
	claims, err := ucase.jwtSvc.JwtVerify(token)
	if err != nil {
		return nil, tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
	}
	return claims, nil
}
func (ucase *authUsecase) AuthJwtVerifyExpired(token string) (*domain.Claims, *domain.UCaseErr) {
	claims, err := ucase.jwtSvc.JwtVerifyExpired(token)
	if err != nil {
		return nil, tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
	}
	return claims, nil
}
func (ucase *authUsecase) AuthJwtDecode(data *domain.ReqAuthJwtDecodeDTO) (*domain.TokenClaims, *domain.UCaseErr) {
	claims, err := ucase.jwtSvc.JwtDecode(data.AccessToken)
	if err != nil {
		return nil, tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
	}
	return &claims, nil
}

func (ucase *authUsecase) Logout(claims *domain.Claims, data *domain.ReqRefreshAccessTokenDTO) (*domain.ResLogoutDTO, *domain.UCaseErr) {
	uid := claims.Uid
	_, deleteCount, err := ucase.refreshTokensRepo.FindByUidsAndDeleteMany(uid)
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	result := &domain.ResLogoutDTO{
		DeleteCount: deleteCount,
	}
	result.Gen()
	return result, nil
}

func (ucase *authUsecase) RefreshAccessToken(claims *domain.Claims, data *domain.ReqRefreshAccessTokenDTO) (*domain.ResLoginTokenDTO, *domain.UCaseErr) {
	uid := claims.Uid
	filter := bson.M{
		"refreshToken": data.RefreshToken,
	}
	refreshToken, err := ucase.refreshTokensRepo.FindOneByFilter(filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrUnauthorized
		}
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	if uid != refreshToken.UserId.Hex() {
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, domain.ErrForbidden, nil)
		return nil, uCaseErr
	}

	user, err := ucase.usersRepo.FindById(uid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = domain.ErrUnauthorized
		}
		uCaseErr := tools.NewUCaseErr(category.Auth, errcode.Default, err, nil)
		return nil, uCaseErr
	}

	userData := &domain.UserData{
		ID:              user.ID,
		Account:         user.Account,
		UserName:        user.UserName,
		NickName:        user.NickName,
		Telephone:       user.Telephone,
		TelephoneRegion: user.TelephoneRegion,
		Roles:           user.Roles,
		ExpiryDate:      user.ExpiryDate,
		ModiTime:        user.ModiTime,
	}
	result, uCaseErr := ucase.GenResLoginTokenDTO(userData)
	if uCaseErr != nil {
		return nil, uCaseErr
	}

	ucase.refreshTokensRepo.DeleteById(refreshToken.ID.Hex())
	return result, nil
}
