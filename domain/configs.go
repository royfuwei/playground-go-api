package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Config 基本資訊
type Config struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Value      interface{}        `json:"value" bson:"value"`
	DataType   string             `json:"dataType" bson:"dataType"`
	UseEncrypt bool               `json:"useEncrypt" bson:"useEncrypt"`
	ModiTime   primitive.DateTime `json:"modiTime" bson:"modiTime"`
}

type ConfigsRepository interface {
	// 初始化 initial/configs.json 寫入db
	InitFindNameAndUpdate(data *Config) (*Config, error)
	FindByName(name string) (*Config, error)
	FindOneByFilter(filter primitive.M) (*Config, error)
}

type ConfigsUsecase interface {
	// 初始化 initial/configs.json
	InitConfigs()
	GetCreateUserDefPwd() (string, *UCaseErr)
	GetAccountExpireAt() (time.Time, *UCaseErr)
}
