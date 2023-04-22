package service

import (
	"fmt"
	"playground-go-api/domain"
	smsMsg "playground-go-api/domain/sms/msg"
	"testing"
	"time"
)

func TestGenSmsMsgByCountryCode(t *testing.T) {
	svc := NewSmsService()
	type args struct {
		locate      string
		captcha     string
		identifier  string
		expiresTime time.Duration
	}
	type want struct {
		result string
		err    error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test tw captcha",
			args: args{
				locate:      "tw",
				captcha:     "123456",
				identifier:  "BBBB",
				expiresTime: time.Duration(120 * time.Second),
			},
			want: want{
				result: fmt.Sprintf(string(smsMsg.TW), "BBBB", "123456", "2分鐘"),
				err:    nil,
			},
		},
		{
			name: "test TW captcha",
			args: args{
				locate:      "TW",
				captcha:     "24680123",
				identifier:  "AAAA",
				expiresTime: time.Duration(50 * time.Second),
			},
			want: want{
				result: fmt.Sprintf(string(smsMsg.TW), "AAAA", "24680123", "50秒"),
				err:    nil,
			},
		},
		{
			name: "test not exist locate",
			args: args{
				locate:      "EOF",
				captcha:     "24680123",
				identifier:  "AAAA",
				expiresTime: time.Duration(60 * time.Second),
			},
			want: want{
				result: "",
				err:    domain.ErrCountryCodeNotExist,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := svc.GenSmsMsgByCountryCode(test.args.locate, test.args.captcha, test.args.identifier, test.args.expiresTime)
			if got != test.want.result {
				t.Errorf("GenSmsMsgByCountryCode() %v, want: %v", got, test.want.result)
			}
			if err != test.want.err {
				t.Errorf("GenSmsMsgByCountryCode() err: %v, want: %v", err, test.want.err)
			}
		})
	}
}
