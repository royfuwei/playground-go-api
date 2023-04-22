package usersfield

// filter 需要用到的欄位
type UsersField string

const (
	Account         UsersField = "account"
	Telephone       UsersField = "telephone"
	TelephoneRegion UsersField = "telephoneRegion"
	NickName        UsersField = "nickname"
)
