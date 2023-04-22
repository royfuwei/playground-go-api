package domain

type EncryptService interface {
	// AEC CBC 編碼
	NewCBCDecrypter(value string) string
	// AEC CBC 解碼
	NewCBCEncrypter(value string) string
}
