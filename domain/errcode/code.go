package errcode

// ErrorCode 定義錯誤代碼型態
type ErrorCode string

const (
	// Default default error code
	Default                       ErrorCode = "default"
	None                          ErrorCode = "none"
	EncryptAesNewCipherError      ErrorCode = "encrypt_aes_new_encrypt"
	EncryptCipherTextToShortError ErrorCode = "encrypt_ciphertext_too_short"
	ErrAuthorizationEmpty         ErrorCode = "header_authorization_empty"
	ErrBearerNotValid             ErrorCode = "header_bearer_not_valid"
	JwtVerifyError                ErrorCode = "jwt_verify_error"
)
