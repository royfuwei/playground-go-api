package service

import (
	"crypto/rsa"
	"io/ioutil"
	"playground-go-api/config"
	"playground-go-api/domain"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

/* rs256 test */
// https://github.com/dgrijalva/jwt-go/blob/master/http_example_test.go
// location of the files used for signing and verification
type jwtService struct {
	privateKeyPath string
	publicKeyPath  string
	verifyKey      *rsa.PublicKey  // openssl genpkey -algorithm RSA -out jwt.rsa -pkeyopt rsa_keygen_bits:2048
	signKey        *rsa.PrivateKey // openssl rsa -in jwt.rsa -pubout > jwt.rsa.pub
}

func NewJwtService() domain.JwtService {
	a := &jwtService{
		privateKeyPath: config.Cfgs.PrivateKeyPath,
		publicKeyPath:  config.Cfgs.PublicKeyPath,
	}
	a.setRsaKeys()
	return a
}

func (svc *jwtService) JwtSign(expiresTime time.Duration, user *domain.UserData, jwtId *string) (expiresAt int64, token string, err error) {
	now := time.Now()
	expiresAt = now.Add(expiresTime).Unix()
	uid := user.ID.Hex()
	account := user.Account
	telephone := user.Telephone
	telephoneRegion := user.TelephoneRegion
	roles := user.Roles
	id := ""
	if jwtId != nil {
		id = *jwtId
	}
	// set claims and sign
	claims := domain.Claims{
		Uid:             uid,
		Account:         account,
		Roles:           roles,
		Telephone:       telephone,
		TelephoneRegion: telephoneRegion,
		StandardClaims: jwt.StandardClaims{
			Id:        id,
			Issuer:    "playground-go-api",
			IssuedAt:  now.Unix(),
			ExpiresAt: expiresAt,
			// Audience:  account,
			// Subject:   account,
			// NotBefore: now.Add(10 * time.Second).Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err = tokenClaims.SignedString(svc.signKey)
	return expiresAt, token, err
}

func (svc *jwtService) JwtVerify(token string) (*domain.Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &domain.Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return svc.verifyKey, nil
	})
	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				message = "token is malformed"
			case ve.Errors&jwt.ValidationErrorUnverifiable != 0:
				message = "token could not be verified because of signing problems"
			case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				message = "signature validation failed"
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				message = "token is expired"
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				message = "token is not yet valid before sometime"
			default:
				message = "can not handle this token"
			}
		}
		glog.Errorf("jwt.ParseWithClaims error message: %v \n", message)
	}
	if claims, ok := tokenClaims.Claims.(*domain.Claims); ok && tokenClaims.Valid {
		return claims, nil
	} else {
		return nil, domain.ErrUnauthorized
	}
}

func (svc *jwtService) JwtVerifyExpired(token string) (*domain.Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &domain.Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return svc.verifyKey, nil
	})
	isErr := false
	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				message = "token is malformed"
				isErr = true
			case ve.Errors&jwt.ValidationErrorUnverifiable != 0:
				message = "token could not be verified because of signing problems"
				isErr = true
			case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				message = "signature validation failed"
				isErr = true
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				message = "token is expired"
				isErr = false
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				message = "token is not yet valid before sometime"
				isErr = true
			default:
				message = "can not handle this token"
				isErr = true
			}
		}
		if isErr {
			glog.Errorf("jwt.ParseWithClaims error message: %v \n", message)
		}
	}
	if claims, ok := tokenClaims.Claims.(*domain.Claims); ok && !isErr {
		return claims, nil
	} else {
		return nil, domain.ErrUnauthorized
	}
}

func (svc *jwtService) JwtDecode(token string) (domain.TokenClaims, error) {
	tokenClaims, _ := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		return svc.verifyKey, nil
	})
	return tokenClaims, nil
}

func (svc *jwtService) setRsaKeys() {
	signBytes, err := ioutil.ReadFile(svc.privateKeyPath)
	svc.fatal(err)
	verifyBytes, err := ioutil.ReadFile(svc.publicKeyPath)
	svc.fatal(err)
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	svc.fatal(err)
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	svc.fatal(err)
	svc.verifyKey = verifyKey
	svc.signKey = signKey
}

func (svc *jwtService) fatal(err error) {
	if err != nil {
		glog.Fatal(err)
	}
}
