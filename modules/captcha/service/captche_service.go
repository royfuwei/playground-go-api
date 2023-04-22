package service

import (
	"bytes"
	"fmt"
	"math/rand"
	"playground-go-api/domain"
	"time"

	"github.com/golang/glog"
)

type captchaService struct{}

func NewCaptchaService() domain.CaptchaService {
	return &captchaService{}
}

func (svc *captchaService) GenCaptchaNumber(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	lenFmt := fmt.Sprint("%0", length, "v")
	vcode := fmt.Sprintf(lenFmt, rnd.Int31n(1000000))
	glog.Infof("vcode: %v", vcode)
	return vcode
}

func (svc *captchaService) GenCaptchaUpString(length int) string {
	var result bytes.Buffer
	var temp byte
	for i := 0; i < length; {
		if svc.randInt(65, 91) != temp {
			temp = svc.randInt(65, 91)
			result.WriteByte(temp)
			i++
		}
	}
	return result.String()
}

func (svc *captchaService) randInt(min int, max int) byte {
	rand.Seed(time.Now().UnixNano())
	return byte(min + rand.Intn(max-min))
}
