package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"playground-go-api/config"
	"playground-go-api/domain"
	"playground-go-api/domain/category"
	"playground-go-api/domain/errcode"
	usersfield "playground-go-api/domain/users/field"
	"playground-go-api/infrastructures/tools"

	"github.com/golang/glog"
	"github.com/nyaruka/phonenumbers"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type usersUsecase struct {
	usersRepo   domain.UsersRepository
	configsRepo domain.ConfigsRepository
	usersSvc    domain.UsersService
	encryptSvc  domain.EncryptService
	smsUcase    domain.SmsUsecase
}

func NewUsersUsecase(
	usersRepo domain.UsersRepository,
	configsRepo domain.ConfigsRepository,
	userSvc domain.UsersService,
	encryptSvc domain.EncryptService,
	smsUcase domain.SmsUsecase,
) domain.UsersUsecase {
	return &usersUsecase{
		usersRepo,
		configsRepo,
		userSvc,
		encryptSvc,
		smsUcase,
	}
}

func (ucase *usersUsecase) CreateUser(user *domain.User) (*domain.UserData, *domain.UCaseErr) {
	password := user.Password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		glog.Error(err)
		uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	user.Password = string(hash)
	if user.Roles == nil {
		user.Roles = []string{}
	}

	if len(user.Telephone) != 0 {
		libphoneParse, err := phonenumbers.Parse(user.Telephone, user.TelephoneRegion)
		if err != nil {
			uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
			return nil, uCaseErr
		}
		user.Telephone = fmt.Sprintf("%v%v", *libphoneParse.CountryCode, *libphoneParse.NationalNumber)
	}
	newUser, err := ucase.usersRepo.Add(user)
	if err != nil {
		glog.Error(err)
		uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
		return nil, uCaseErr
	}
	result := &domain.UserData{
		ID:              newUser.ID,
		Account:         newUser.Account,
		NickName:        newUser.NickName,
		Telephone:       newUser.Telephone,
		TelephoneRegion: newUser.TelephoneRegion,
		Roles:           newUser.Roles,
		ExpiryDate:      newUser.ExpiryDate,
		ModiTime:        newUser.ModiTime,
	}
	return result, nil
}

func (ucase *usersUsecase) CreateUserByTelephone(data *domain.ReqCreateUserByTelephone) (*domain.UserData, *domain.UCaseErr) {
	if len(data.Telephone) != 0 {
		filter := bson.M{
			string(usersfield.Telephone):       data.Telephone,
			string(usersfield.TelephoneRegion): data.TelephoneRegion,
		}
		isExistUser, err := ucase.usersRepo.IsExistByFilter(filter)
		if err != nil {
			glog.Error(err)
			uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
			return nil, uCaseErr
		}
		if isExistUser {
			uCaseErr := ucase.genExistUserUCaseErr(string(usersfield.Telephone))
			return nil, uCaseErr
		}
	}
	user := &domain.User{
		UserName:        data.UserName,
		NickName:        data.NickName,
		Telephone:       data.Telephone,
		TelephoneRegion: data.TelephoneRegion,
		Password:        data.Password,
	}
	result, uCaseErr := ucase.CreateUser(user)
	if uCaseErr != nil {
		return nil, uCaseErr
	}
	return result, nil
}

func (ucase *usersUsecase) genExistUserUCaseErr(field string) *domain.UCaseErr {
	err := fmt.Errorf("field: %v, is exist user", field)
	uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
	return uCaseErr
}

func (ucase *usersUsecase) CreateUserByAccount(data *domain.ReqCreateUserByAccount) (*domain.UserData, *domain.UCaseErr) {
	isExistUser, uCaseErr := ucase.isExitAccount(data.Account)
	if uCaseErr != nil {
		glog.Error(uCaseErr.Err)
		return nil, uCaseErr
	}
	if isExistUser {
		uCaseErr := ucase.genExistUserUCaseErr(string(usersfield.Account))
		return nil, uCaseErr
	}

	if isExistUser, uCaseErr = ucase.isExitNickname(data.NickName); uCaseErr != nil {
		glog.Error(uCaseErr.Err)
		return nil, uCaseErr
	}
	if isExistUser {
		uCaseErr := ucase.genExistUserUCaseErr(string(usersfield.NickName))
		return nil, uCaseErr
	}

	if len(data.Telephone) != 0 {
		isExistUser, uCaseErr := ucase.isExitTelephone(data.Telephone, data.TelephoneRegion)
		if uCaseErr != nil {
			glog.Error(uCaseErr.Err)
			return nil, uCaseErr
		}
		if isExistUser {
			uCaseErr := ucase.genExistUserUCaseErr(string(usersfield.Telephone))
			return nil, uCaseErr
		}
	}
	user := &domain.User{
		UserName:        data.UserName,
		NickName:        data.NickName,
		Telephone:       data.Telephone,
		TelephoneRegion: data.TelephoneRegion,
		Password:        data.Password,
	}
	result, uCaseErr := ucase.CreateUser(user)
	if uCaseErr != nil {
		return nil, uCaseErr
	}
	return result, nil
}

func (ucase *usersUsecase) isExitNickname(nickname string) (bool, *domain.UCaseErr) {
	filter := bson.M{
		string(usersfield.NickName): nickname,
	}
	isExistUser, err := ucase.usersRepo.IsExistByFilter(filter)
	if err != nil {
		glog.Error(err)
		uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
		return false, uCaseErr
	}
	if isExistUser {
		return true, nil
	}
	return false, nil
}

func (ucase *usersUsecase) isExitAccount(account string) (bool, *domain.UCaseErr) {
	filter := bson.M{
		string(usersfield.Account): account,
	}
	isExistUser, err := ucase.usersRepo.IsExistByFilter(filter)
	if err != nil {
		glog.Error(err)
		uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
		return false, uCaseErr
	}
	if isExistUser {
		return true, nil
	}
	return false, nil
}

func (ucase *usersUsecase) isExitTelephone(telephone, telephoneRegion string) (bool, *domain.UCaseErr) {
	filter := bson.M{
		string(usersfield.Telephone):       telephone,
		string(usersfield.TelephoneRegion): telephoneRegion,
	}
	isExistUser, err := ucase.usersRepo.IsExistByFilter(filter)
	if err != nil {
		glog.Error(err)
		uCaseErr := tools.NewUCaseErr(category.User, errcode.Default, err, nil)
		return false, uCaseErr
	}
	if isExistUser {
		return true, nil
	}
	return false, nil
}

// initial users
func (ucase *usersUsecase) InitUsers() {
	fmt.Printf("NewUsersUsecase InitUsers \n")
	var initJson []domain.InitUserJson

	byteData, err := ioutil.ReadFile(config.Cfgs.UsersInitPath)
	if err != nil {
		glog.Errorf("NewConfigsUsecase InitConfigs ioutil.ReadFile error: %v \n", err)
		return
	}

	err = json.Unmarshal(byteData, &initJson)
	if err != nil {
		glog.Errorf("NewConfigsUsecase InitConfigs json.Unmarshal error: %v \n", err)
		return
	}

	// adminConfigs, err := ucase.configsRepo.FindByName("")
	for _, initUser := range initJson {
		isExistUser, err := ucase.usersRepo.IsExistByField(string(usersfield.Account), initUser.Account)
		if err != nil {
			glog.Error(err)
			continue
		}
		user := &domain.User{
			Account:         initUser.Account,
			UserName:        initUser.UserName,
			NickName:        initUser.NickName,
			Telephone:       initUser.Telephone,
			TelephoneRegion: initUser.TelephoneRegion,
			Roles:           initUser.Roles,
			ExpiryDate:      initUser.ExpiryDate,
		}

		adminPwdConfig, err := ucase.configsRepo.FindByName(config.Cfgs.AdminPwd)
		if err != nil {
			glog.Error(err)
			continue
		}
		adminPwd := ucase.encryptSvc.NewCBCDecrypter(adminPwdConfig.Value.(string))
		createUserDefPwdConfig, err := ucase.configsRepo.FindByName(config.Cfgs.CreateUserDefPwd)
		if err != nil {
			glog.Error(err)
			continue
		}
		createUserDefPwd := ucase.encryptSvc.NewCBCDecrypter(createUserDefPwdConfig.Value.(string))
		var pwd string
		if initUser.Account == config.Cfgs.SuperAdminAccount {
			pwd = adminPwd
		} else {
			pwd = createUserDefPwd
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			glog.Error(err)
			continue
		}
		user.Password = string(hash)
		if !isExistUser {
			if _, err := ucase.usersRepo.Add(user); err != nil {
				glog.Error(err)
				continue
			}
		} else if initUser.Mandatory {
			filter := bson.M{
				"account": initUser.Account,
			}
			if _, err := ucase.usersRepo.FindFilterAndUpdate(filter, user); err != nil {
				glog.Error(err)
				continue
			}
		}
	}

}
