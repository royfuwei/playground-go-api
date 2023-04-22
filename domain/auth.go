package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth represent auth data
type Auth struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	MemberId     string             `json:"memberId,omitempty" bson:"memberId,omitempty"`
	RefreshToken string             `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
	CreateTime   primitive.DateTime `json:"createTime,omitempty" bson:"createTime,omitempty" swaggertype:"string"`
}

type RefreshToken struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId       primitive.ObjectID `json:"userId" bson:"userId"`
	RefreshToken string             `json:"refreshToken" bson:"refreshToken"`
	CreateAt     primitive.DateTime `json:"createAt" bson:"createAt" swaggertype:"string"`
}

type ResLoginTokenDTO struct {
	AccessToken    string `json:"accessToken" bson:"accessToken"`
	RefreshToken   string `json:"refreshToken" bson:"refreshToken"`
	AccessTokenExp int    `json:"accessTokenExp" bson:"accessTokenExp"`
}

type ReqLoginAccountDTO struct {
	Account  string `form:"account" binding:"required" json:"account" bson:"account"`
	Password string `form:"password" binding:"required" json:"password" bson:"password"`
}

type ReqLoginTelephoneDTO struct {
	Telephone       string `json:"telephone" bson:"telephone" binding:"required"`
	TelephoneRegion string `json:"telephoneRegion" bson:"telephoneRegion" binding:"required" example:"TW"`
	Password        string `json:"password" bson:"password" binding:"required"`
}

type ReqRefreshAccessTokenDTO struct {
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}

type ReqAuthJwtDecodeDTO struct {
	AccessToken string `json:"accessToken"`
}

type ResLogoutDTO struct {
	Message     string `json:"message"`
	IsDelete    bool   `json:"isDelete"`
	DeleteCount int    `json:"count"`
}

func (dto *ResLogoutDTO) Gen() {
	dto.IsDelete = dto.DeleteCount > 0
	if dto.IsDelete {
		dto.Message = "success logout and delete refreshTokens"
	} else {
		dto.Message = "success logout"
	}
}

type RefreshTokensRepository interface {
	Add(data *RefreshToken) (*RefreshToken, error)
	DeleteById(id string) (bool, error)
	FindById(id string) (refreshToken *RefreshToken, err error)
	FindByIds(ids ...string) (refreshTokens []*RefreshToken, total int64, err error)
	FindByUids(uids ...string) (refreshTokens []*RefreshToken, total int64, err error)
	FindByUidsAndDeleteMany(uid ...string) (bool, int, error)
	FindOneByFilter(filter bson.M) (*RefreshToken, error)
}

type AuthUsecase interface {
	GenResLoginTokenDTO(userData *UserData) (*ResLoginTokenDTO, *UCaseErr)
	CreateAccessToken(userData *UserData, jwtId *string) (expiresAt int64, token string, uCaseErr *UCaseErr)
	CreateRefreshToken(uid primitive.ObjectID) (string, *UCaseErr)
	AuthJwtVerify(token string) (*Claims, *UCaseErr)
	AuthJwtVerifyExpired(token string) (*Claims, *UCaseErr)
	AuthJwtDecode(data *ReqAuthJwtDecodeDTO) (*TokenClaims, *UCaseErr)
	LoginAccount(data *ReqLoginAccountDTO) (*ResLoginTokenDTO, *UCaseErr)
	LoginTelephone(data *ReqLoginTelephoneDTO) (*ResLoginTokenDTO, *UCaseErr)
	// 登出
	// 1. 檢查帶過來的refresh token是否為我們系統所簽發的。
	// 2. 若否則reject操作。
	// 3. 若refresh token合法則從db中將該組refresh token移除。
	Logout(claims *Claims, data *ReqRefreshAccessTokenDTO) (*ResLogoutDTO, *UCaseErr)
	// 拿refresh token來重新產生一組token
	// 1. 檢查refresh token是否還在db裡面，確保這組token是由系統所簽發的。
	// 2. 檢查是否已經過期。
	// 3. 創建refresh token並記錄於db。
	// 4. 刪除舊有的refresh token。
	RefreshAccessToken(claims *Claims, data *ReqRefreshAccessTokenDTO) (*ResLoginTokenDTO, *UCaseErr)
}
