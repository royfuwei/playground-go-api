package service

import "testing"

func TestGenCaptchaUpString(t *testing.T) {
	svc := NewCaptchaService()

	type args struct {
		length int
	}
	type want struct {
		vcode  string
		length int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test 6 length",
			args: args{
				length: 6,
			},
			want: want{
				vcode:  "",
				length: 6,
			},
		},
		{
			name: "test 7 length",
			args: args{
				length: 7,
			},
			want: want{
				vcode:  "",
				length: 7,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := svc.GenCaptchaUpString(test.args.length); len(got) != test.want.length {
				t.Errorf("GenCaptchaUpString() got: %v, length: %v, want length: %v", got, len(got), test.want.length)
			}
		})
	}
}

func TestGenCaptchaNumber(t *testing.T) {
	svc := NewCaptchaService()

	type args struct {
		length int
	}
	type want struct {
		vcode  string
		length int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test 6 length",
			args: args{
				length: 6,
			},
			want: want{
				vcode:  "",
				length: 6,
			},
		},
		{
			name: "test 7 length",
			args: args{
				length: 7,
			},
			want: want{
				vcode:  "",
				length: 7,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := svc.GenCaptchaNumber(test.args.length); len(got) != test.want.length {
				t.Errorf("GenCaptchaNumber() got: %v, length: %v, want length: %v", got, len(got), test.want.length)
			}
		})
	}
}
