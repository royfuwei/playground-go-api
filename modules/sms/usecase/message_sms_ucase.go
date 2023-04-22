package usecase

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"playground-go-api/config"
	"playground-go-api/domain"
)

type messageSmsUsecase struct {
	account  string
	password string
	encoding string
	url      string
}

// http://api.message.net.tw/query.php?id=&password=&columns=prms;tel;mstat&sindex=0&mcount=100
// id 使用者帳號
// password 密碼
// sdate 發送時間(直接發送則不用設) def:即時
// tel 電話一;電話二;電話三 max:100
// msg 簡訊內容 若使用URL編碼,參考附表二
// mtype 簡訊種類 def:G G:一般簡訊(G為大寫)
// encoding 簡訊內容的編碼方式 def:big5 utf8:簡訊內容採用UTF-8編碼 urlencode:簡訊內容採用URL編碼 urlencode_utf8:簡訊內容採用URL與UTF-8編碼
// 附註:
// 1.每次遞交簡訊之內容字數，純英文時最多160個字元，含中文時最多70個字元。 2.每次遞交簡訊之tel數目(以分號隔開)，最少1個號碼、最多100個號碼。 3.國內簡訊發送，發一則扣一通，國際簡訊發一則扣三通，國際簡訊發送不完全支援mtype指令。 4.以上發送皆提供來源IP限制功能。
func NewMessageSmsUsecase() domain.SmsUsecase {
	return &messageSmsUsecase{
		account:  config.Cfgs.SmsAccount,
		password: config.Cfgs.SmsPassword,
		encoding: "urlencode_utf8",
		url:      "http://api.message.net.tw/query.php",
	}
}

func (ucase *messageSmsUsecase) SendMsgToTelephone(tel string, msg string) (*domain.SendMsgToTelDTO, error) {
	url, err := url.Parse(ucase.url)
	if err != nil {
		// glog.Fatal(err)
		return nil, err
	}
	query := url.Query()
	query.Add("id", ucase.account)
	query.Add("password", ucase.password)
	query.Add("encoding", ucase.encoding)
	query.Add("tel", tel)
	query.Add("msg", msg)
	url.RawQuery = query.Encode()
	fmt.Printf("url: %v \n", url.String())

	res, err := http.Get(url.String())
	if err != nil {
		// glog.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("ioutil.ReadAll error: %v", err)
		return nil, err
	}
	if res.StatusCode >= 400 {
		err = fmt.Errorf("statusCode: %d, status: %s \n", res.StatusCode, res.Status)
		return nil, err
	}

	result := &domain.SendMsgToTelDTO{}
	result.Success()
	return result, nil

}
