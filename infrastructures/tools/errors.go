package tools

import (
	"fmt"
	"net/http"
	"playground-go-api/domain"
	"playground-go-api/domain/category"
	"playground-go-api/domain/errcode"
	"runtime"
	"strings"
	"time"
)

// ParseErrStatus 解析錯誤的response status code
func ParseErrStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrNotExist:
		return http.StatusConflict
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrValidation:
		return http.StatusBadRequest
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// NewResponseErr new response error
func NewResponseErr(err error, errorCode errcode.ErrorCode, path string, status *int) domain.ResponseError {
	if status == nil {
		status = Int(ParseErrStatus(err))
	}
	return domain.ResponseError{
		Message:   err.Error(),
		Status:    *status,
		ErrorCode: errorCode,
		Path:      path,
	}
}

// NewUCaseErr new usecase error
func NewUCaseErr(category category.Category, errorCode errcode.ErrorCode, err error, content interface{}) *domain.UCaseErr {
	_, fn, line, _ := runtime.Caller(1)
	slash := strings.LastIndex(fn, "/")
	file := fn[slash+1:]
	position := fmt.Sprintf("%s:%d", file, line)
	return &domain.UCaseErr{
		Category:  category,
		ErrorCode: errorCode,
		Err:       err,
		Position:  position,
		ModiTime:  time.Now(),
		Content:   content,
	}
}
