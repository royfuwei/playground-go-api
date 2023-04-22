package domain

type RegisterUsecase interface {
	RegisterSmsCaptchaSend(data *ReqSmsCaptchaSendDTO) (*ResSmsCaptchaSendDTO, *UCaseErr)
	RegisterUserByTelephone(claims *Claims, data *ReqCreateUserByTelephone) (*ResLoginTokenDTO, *UCaseErr)
}
