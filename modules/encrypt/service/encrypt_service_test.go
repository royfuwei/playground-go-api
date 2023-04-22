package service

import (
	"testing"
)

func TestNewCBCEncrypter(t *testing.T) {
	svc := NewEncryptService()
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "TestNewCBCEncrypter test",
			value: "test",
			want:  "b28fc8efa623d43cd9f1f3b2be72393e",
		},
		{
			name:  "TestNewCBCEncrypter encrypt",
			value: "encrypt",
			want:  "3899f1d9f4b8b5cb9efdcaf0b3bec9f1",
		},
		{
			name:  "TestNewCBCEncrypter test_encrypt",
			value: "test_encrypt",
			want:  "c0bd9f11c4c48847768e4b1b1b729750",
		},
		{
			name:  "TestNewCBCEncrypter TestNewCBCEncrypter",
			value: "TestNewCBCEncrypter",
			want:  "acfe70fbaf31e0193e7190729525e041ab571c74f24d5937b029a91d43abcc25",
		},
		{
			name:  "TestNewCBCEncrypter TestNewCBCEncrypter",
			value: "",
			want:  "bfd5a76d5f790e4a2c58b386ca107f1e",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := svc.NewCBCEncrypter(test.value)
			// assert.Equal(result, test.want)
			if result != test.want {
				t.Errorf("NewCBCEncrypter() %v, want %v", result, test.want)
			}
		})
	}
}
func TestNewCBCDecrypter(t *testing.T) {
	svc := NewEncryptService()
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "TestNewCBCDecrypter test",
			value: "b28fc8efa623d43cd9f1f3b2be72393e",
			want:  "test",
		},
		{
			name:  "TestNewCBCDecrypter encrypt",
			value: "3899f1d9f4b8b5cb9efdcaf0b3bec9f1",
			want:  "encrypt",
		},
		{
			name:  "TestNewCBCDecrypter test_encrypt",
			value: "c0bd9f11c4c48847768e4b1b1b729750",
			want:  "test_encrypt",
		},
		{
			name:  "TestNewCBCDecrypter TestNewCBCEncrypter",
			value: "acfe70fbaf31e0193e7190729525e041ab571c74f24d5937b029a91d43abcc25",
			want:  "TestNewCBCEncrypter",
		},
		{
			name: "TestNewCBCDecrypter error",
			// value: "acfe70fbaf31e0193e7190729525e041ab571c74f24d5937b029a91d43abcc22",
			value: "acfe70fbaf31e0193e7190729525e041ab571c74f24d5937b029a91d43abcc21",
			want:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := svc.NewCBCDecrypter(test.value)
			// if errCode != nil {
			// 	t.Log(errCode)
			// }
			// if result == test.want {
			// 	t.Logf("NewCBCDecrypter() %v, want %v", result, test.want)
			// }
			if result != test.want {
				t.Errorf("NewCBCDecrypter() %v, want %v", result, test.want)
			}
			// if result == "TestNewCBCEncrypter" {
			// 	t.Errorf("NewCBCDecrypter() %v, want %v", result, test.want)
			// }
		})
	}
}
