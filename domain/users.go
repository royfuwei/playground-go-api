package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 基本資訊
type User struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Account         string             `json:"account" bson:"account"`
	UserName        string             `json:"username" bson:"username"`
	NickName        string             `json:"nickname" bson:"nickname"`
	Telephone       string             `json:"telephone" bson:"telephone"`
	TelephoneRegion string             `json:"telephoneRegion" bson:"telephoneRegion"`
	Roles           []string           `json:"roles" bson:"roles"`
	Password        string             `json:"password" bson:"password"`
	ExpiryDate      primitive.DateTime `json:"expiryDate" bson:"expiryDate"`
	ModiTime        primitive.DateTime `json:"modiTime" bson:"modiTime"`
}

type InitUserJson struct {
	Account         string             `json:"account" bson:"account"`
	UserName        string             `json:"username" bson:"username"`
	NickName        string             `json:"nickname" bson:"nickname"`
	Telephone       string             `json:"telephone" bson:"telephone"`
	TelephoneRegion string             `json:"telephoneRegion" bson:"telephoneRegion"`
	Roles           []string           `json:"roles" bson:"roles"`
	Mandatory       bool               `json:"mandatory" bson:"mandatory"`
	ExpiryDate      primitive.DateTime `json:"expiryDate" bson:"expiryDate"`
}

type ReqCreateUserByAccount struct {
	Account         string `json:"account"`
	Password        string `json:"password"`
	UserName        string `json:"username"`
	NickName        string `json:"nickname"`
	Telephone       string `json:"telephone"`
	TelephoneRegion string `json:"telephoneRegion"`
}

type ReqCreateUserByTelephone struct {
	Password        string `json:"password"`
	UserName        string `json:"username"`
	NickName        string `json:"nickname"`
	Telephone       string `json:"telephone"`
	TelephoneRegion string `json:"telephoneRegion"`
}

type UserData struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Account         string             `json:"account" bson:"account"`
	UserName        string             `json:"username" bson:"username"`
	NickName        string             `json:"nickname" bson:"nickname"`
	Telephone       string             `json:"telephone" bson:"telephone"`
	TelephoneRegion string             `json:"telephoneRegion" bson:"telephoneRegion"`
	Roles           []string           `json:"roles" bson:"roles"`
	ExpiryDate      primitive.DateTime `json:"expiryDate" bson:"expiryDate"`
	ModiTime        primitive.DateTime `json:"modiTime" bson:"modiTime"`
}

type FilterUsersByField struct {
	FieldName string
	Value     interface{}
}

type UsersRepository interface {
	FindOneByFilter(filter primitive.M) (*User, error)
	FindById(id string) (user *User, err error)
	// 檢查欄位的值是否存在
	IsExistByField(fieldName string, fieldValue string) (bool, error)
	IsExistByFilter(filter bson.M) (bool, error)
	Add(data *User) (*User, error)
	FindFilterAndUpdate(filter primitive.M, data *User) (*User, error)
}

type UsersService interface {
	// 檢查密碼是否正確
	CheckPassword(attemptPass string, user *User) bool
}

type UsersUsecase interface {
	InitUsers()
	CreateUser(data *User) (*UserData, *UCaseErr)
	CreateUserByAccount(data *ReqCreateUserByAccount) (*UserData, *UCaseErr)
	CreateUserByTelephone(data *ReqCreateUserByTelephone) (*UserData, *UCaseErr)
}
