package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"playground-go-api/config"
	"playground-go-api/domain"
	"playground-go-api/domain/category"
	"playground-go-api/domain/errcode"
	"playground-go-api/infrastructures/tools"
	"time"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type configsUsecase struct {
	configsRepo domain.ConfigsRepository
	encryptSvc  domain.EncryptService
}

func NewConfigsUsecase(
	configsRepo domain.ConfigsRepository,
	encryptSvc domain.EncryptService,
) domain.ConfigsUsecase {
	return &configsUsecase{
		configsRepo,
		encryptSvc,
	}
}

func (ucase *configsUsecase) GetCreateUserDefPwd() (string, *domain.UCaseErr) {
	config, err := ucase.configsRepo.FindByName("createUserDefPwd")
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Config, errcode.Default, err, nil)
		return "", uCaseErr
	}
	valueStr := fmt.Sprintf("%v", config.Value)
	result := ucase.encryptSvc.NewCBCDecrypter(valueStr)
	// fmt.Printf("%v \n", result)
	return result, nil
}

func (ucase *configsUsecase) GetAccountExpireAt() (time.Time, *domain.UCaseErr) {
	var accountExpireAt time.Time
	config, err := ucase.configsRepo.FindByName("accountExpireAt")
	if err != nil {
		uCaseErr := tools.NewUCaseErr(category.Config, errcode.Default, err, nil)
		return accountExpireAt, uCaseErr
	}
	dateTime := config.Value.(primitive.DateTime)
	accountExpireAt = dateTime.Time()
	// fmt.Printf("%v \n", accountExpireAt)
	return accountExpireAt, nil
}

func (ucase *configsUsecase) InitConfigs() {
	fmt.Printf("NewConfigsUsecase InitConfigs \n")
	var initJson []domain.Config

	byteData, err := ioutil.ReadFile(config.Cfgs.ConfigsInitPath)
	if err != nil {
		glog.Errorf("NewConfigsUsecase InitConfigs ioutil.ReadFile error: %v \n", err)
		return
	}

	err = json.Unmarshal(byteData, &initJson)
	if err != nil {
		glog.Errorf("NewConfigsUsecase InitConfigs json.Unmarshal error: %v \n", err)
		return
	}

	for _, config := range initJson {
		data := ucase.valueNormalizer(&config)

		_, err := ucase.configsRepo.InitFindNameAndUpdate(data)
		if err != nil {
			glog.Errorf("NewConfigsUsecase InitConfigs configsRepo.FindOneAndUpdate error: %v \n", err)
			continue
		}
	}

}

// valueNormalizer value配置轉換
func (ucase *configsUsecase) valueNormalizer(data *domain.Config) *domain.Config {
	valueStr := fmt.Sprintf("%v", data.Value)
	valueDate, err := tools.TimeParseByFormats(valueStr)
	if err == nil {
		data.Value = primitive.NewDateTimeFromTime(valueDate)
	}
	if data.UseEncrypt {
		data.Value = ucase.encryptSvc.NewCBCEncrypter(valueStr)
	}
	return data
}
