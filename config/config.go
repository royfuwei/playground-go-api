package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DataType = string

type Configs struct {
	MongoAddr           string
	Port                string
	RetryPeriod         int
	MaxRetryCount       int
	MaxBatchSize        int
	Goroutine           int
	UsersInitPath       string
	ConfigsInitPath     string
	MgoDBName           string
	AdminPwd            string
	CreateUserDefPwd    string
	SuperAdminAccount   string
	AccessTokenExpired  int
	RefreshTokenExpired int
	CaptchaExpired      int
	PrivateKeyPath      string
	PublicKeyPath       string
	SmsAccount          string
	SmsPassword         string
}

var (
	Cfgs   *Configs
	String DataType = "string"
	Date   DataType = "date"
)

type InitialConfigsDto struct {
	Name       string
	Value      string
	UseEncrypt bool
	DataType   DataType // string | date
}

func initVariable() {
	flag.Set("alsologtostderr", "true")
	flag.CommandLine.Parse([]string{})
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.String("mongo_addr", "127.0.0.1:27017", "mongodb address")
	pflag.String("private_key_path", "jwt/jwt.rsa", "privkey path")
	pflag.String("public_key_path", "jwt/jwt.rsa.pub", "pubkey path")
	pflag.String("port", "5003", "serve port")
	pflag.Int("retry_period", 3, "second of retry period")
	pflag.Int("max_retry_count", 5, "max retry count")
	pflag.Int("max_batch_size", 5000, "es捲動的最大size，代表每一個批次最大返回數量")
	pflag.Int("goroutine", 1, "goroutine number")
	pflag.Int("access_token_expired", 1800, "AccessTokenExpired sec")
	pflag.Int("refresh_token_expired", 3600, "RefreshTokenExpired sec")
	pflag.Int("captcha_expired", 600, "CaptchaExpired")
	pflag.String("sms_account", "account", "SMS Account")
	pflag.String("sms_password", "password", "SMS Password")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

// NewConfig new config
func NewConfig() *Configs {
	initVariable()
	Cfgs = &Configs{
		MongoAddr:           viper.GetString("mongo_addr"),
		Port:                ":" + viper.GetString("port"),
		RetryPeriod:         viper.GetInt("retry_period"),
		MaxRetryCount:       viper.GetInt("max_retry_count"),
		MaxBatchSize:        viper.GetInt("max_batch_size"),
		Goroutine:           viper.GetInt("goroutine"),
		AccessTokenExpired:  viper.GetInt("access_token_expired"),
		RefreshTokenExpired: viper.GetInt("refresh_token_expired"),
		CaptchaExpired:      viper.GetInt("captcha_expired"),
		UsersInitPath:       "initial/users.json",
		ConfigsInitPath:     "initial/configs.json",
		MgoDBName:           "auth",
		AdminPwd:            "adminPwd",
		CreateUserDefPwd:    "createUserDefPwd",
		SuperAdminAccount:   "admin",
		PrivateKeyPath:      viper.GetString("private_key_path"),
		PublicKeyPath:       viper.GetString("public_key_path"),
		SmsAccount:          viper.GetString("sms_account"),
		SmsPassword:         viper.GetString("sms_password"),
	}
	// initConfigsJson(Cfgs.ConfigsInitPath)
	return Cfgs
}

// func initConfigsJson(configsInitPath string) {
// 	var initialConfigsDto []InitialConfigsDto
// 	initConfigsJsonByte, err := ioutil.ReadFile(configsInitPath)
// 	if err != nil {
// 		glog.Error(err)
// 	}
// 	fmt.Printf(string(initConfigsJsonByte))
// 	err = json.Unmarshal(initConfigsJsonByte, &initialConfigsDto)
// 	// fmt.Printf("initialConfigsDto: %v \n", initialConfigsDto)
// 	if err != nil {
// 		glog.Error(err)
// 	}
// }
