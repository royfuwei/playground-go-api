package domain

import (
	"errors"
	"playground-go-api/domain/category"
	"playground-go-api/domain/errcode"
	"time"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if any the Not Found Error happen
	ErrNotFound = errors.New("Not Found Error")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrNotExist will throw if the current action not exists
	ErrNotExist = errors.New("Your Item not exist")
	// ErrValidation will throw if the request body valid failed
	ErrValidation          = errors.New("request body validate failed")
	ErrCountryCodeNotExist = errors.New("Country code is not exist")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrBearerNotValid      = errors.New("bearer token not valid")
	ErrParseWithClaims     = errors.New("error ParseWithClaims")
	ErrAuthorizationEmpty  = errors.New("header Authorization empty")
	ErrForbidden           = errors.New("Forbidden Error")
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	// 錯誤狀態碼
	Status int `json:"status"`
	// 錯誤訊息
	Message string `json:"message"`
	// 自定義錯誤代碼
	ErrorCode errcode.ErrorCode `json:"errorCode"`
	// API路徑
	Path string `json:"path"`
}

// UCaseErr represent the use case error response
type UCaseErr struct {
	// 來源是出自系統的那一個功能模組
	Category category.Category `json:"category,omitempty" bson:"category"`
	// 自定義錯誤代碼
	ErrorCode errcode.ErrorCode `json:"errorCode,omitempty" bson:"errorCode"`
	// Error
	Err error `json:"error,omitempty" bson:"error,omitempty"`
	// 錯誤檔案位置
	Position string `json:"position,omitempty" bson:"position"`
	// 錯誤紀錄時間
	ModiTime time.Time `json:"modiTime" from:"modiTime" binding:"required" bson:"modiTime"`
	// 額外紀錄的錯誤內容
	Content interface{} `json:"content,omitempty" bson:"content"`
	// `json:"type,omitempty" bson:"type"`
	// `json:"message,omitempty" bson:"message"`
}
