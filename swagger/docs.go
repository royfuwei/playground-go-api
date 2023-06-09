// Package swagger GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package swagger

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "取得App 名稱",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "取得App 名稱",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.App"
                        }
                    }
                }
            }
        },
        "/auth/jwt/decode": {
            "post": {
                "description": "JwtDecode 解析jwt 內容",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "JwtDecode",
                "parameters": [
                    {
                        "description": "jwt",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqAuthJwtDecodeDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/auth/jwt/verify": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "JwtVerify 驗證jwt，成功返回內容",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "JwtVerify",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/auth/jwt/verify-expired": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "JwtVerifyExpired 驗證過期jwt，成功返回內容",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "JwtVerifyExpired",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/auth/login/account": {
            "post": {
                "description": "Auth Account Login",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Auth Account Login",
                "parameters": [
                    {
                        "description": "account login 內容",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqLoginAccountDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success login response",
                        "schema": {
                            "$ref": "#/definitions/domain.ResLoginTokenDTO"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/auth/login/telephone": {
            "post": {
                "description": "Auth Telephone Login",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Auth Telephone Login",
                "parameters": [
                    {
                        "description": "telephone login 內容",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqLoginTelephoneDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success login response",
                        "schema": {
                            "$ref": "#/definitions/domain.ResLoginTokenDTO"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Auth 登出",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Auth 登出",
                "responses": {
                    "200": {
                        "description": "User has been successfully logout.",
                        "schema": {
                            "$ref": "#/definitions/domain.ResLogoutDTO"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/auth/refresh-access-token": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "當access token失效時，拿refresh token來重新產生一組",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "refresh token 產生token",
                "parameters": [
                    {
                        "description": "telephone login 內容",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqRefreshAccessTokenDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success login response",
                        "schema": {
                            "$ref": "#/definitions/domain.ResLoginTokenDTO"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/captcha/sms/send": {
            "post": {
                "description": "Send Captcha SMS",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "captcha"
                ],
                "summary": "Send Captcha SMS",
                "parameters": [
                    {
                        "description": "Send Captcha SMS 內容",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqSmsCaptchaSendDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Send Captcha SMS 成功",
                        "schema": {
                            "$ref": "#/definitions/domain.ResSmsCaptchaSendDTO"
                        }
                    },
                    "400": {
                        "description": "Send Captcha SMS 失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/captcha/sms/validate": {
            "post": {
                "description": "驗證 Captcha SMS，不會檢查是否已經有使用者及電話號碼",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "captcha"
                ],
                "summary": "驗證 Captcha SMS",
                "parameters": [
                    {
                        "description": "驗證 Captcha SMS 內容",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqSmsCaptchaValidateDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "驗證 Captcha SMS 成功",
                        "schema": {
                            "$ref": "#/definitions/domain.ResSmsCaptchaValidateDTO"
                        }
                    },
                    "400": {
                        "description": "驗證 Captcha SMS 失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    },
                    "401": {
                        "description": "驗證 Captcha SMS 失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResSmsCaptchaValidateDTO"
                        }
                    }
                }
            }
        },
        "/register/captcha/sms/send": {
            "post": {
                "description": "發送註冊的sms captcha",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "register"
                ],
                "summary": "send sms Register captcha",
                "parameters": [
                    {
                        "description": "發送的手機號碼",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqSmsCaptchaSendDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/domain.ResSmsCaptchaSendDTO"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        },
        "/register/telephone/user": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Register User By telephone",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "register"
                ],
                "summary": "Register User By telephone",
                "parameters": [
                    {
                        "description": "建立",
                        "name": "default",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.ReqCreateUserByTelephone"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/domain.ResLoginTokenDTO"
                        }
                    },
                    "400": {
                        "description": "請求的body、header驗證失敗",
                        "schema": {
                            "$ref": "#/definitions/domain.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.App": {
            "type": "object",
            "properties": {
                "app": {
                    "type": "string"
                }
            }
        },
        "domain.ReqAuthJwtDecodeDTO": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                }
            }
        },
        "domain.ReqCreateUserByTelephone": {
            "type": "object",
            "properties": {
                "nickname": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "telephone": {
                    "type": "string"
                },
                "telephoneRegion": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "domain.ReqLoginAccountDTO": {
            "type": "object",
            "required": [
                "account",
                "password"
            ],
            "properties": {
                "account": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "domain.ReqLoginTelephoneDTO": {
            "type": "object",
            "required": [
                "password",
                "telephone",
                "telephoneRegion"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "telephone": {
                    "type": "string"
                },
                "telephoneRegion": {
                    "type": "string",
                    "example": "TW"
                }
            }
        },
        "domain.ReqRefreshAccessTokenDTO": {
            "type": "object",
            "properties": {
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "domain.ReqSmsCaptchaSendDTO": {
            "type": "object",
            "required": [
                "telephone",
                "telephoneRegion"
            ],
            "properties": {
                "telephone": {
                    "description": "電話號碼",
                    "type": "string"
                },
                "telephoneRegion": {
                    "description": "國家地區代碼",
                    "type": "string",
                    "example": "TW"
                }
            }
        },
        "domain.ReqSmsCaptchaValidateDTO": {
            "type": "object",
            "required": [
                "captcha",
                "identifier",
                "telephone",
                "telephoneRegion"
            ],
            "properties": {
                "captcha": {
                    "description": "手機取得的驗證碼",
                    "type": "string"
                },
                "identifier": {
                    "description": "識別碼",
                    "type": "string"
                },
                "telephone": {
                    "description": "電話號碼",
                    "type": "string"
                },
                "telephoneRegion": {
                    "description": "國家地區代碼",
                    "type": "string",
                    "example": "TW"
                }
            }
        },
        "domain.ResLoginTokenDTO": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "accessTokenExp": {
                    "type": "integer"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "domain.ResLogoutDTO": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "isDelete": {
                    "type": "boolean"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "domain.ResSmsCaptchaSendDTO": {
            "type": "object",
            "properties": {
                "expiresTime": {
                    "description": "存在時長(秒)",
                    "type": "integer"
                },
                "expiryDate": {
                    "description": "截止時間",
                    "type": "string"
                },
                "identifier": {
                    "description": "辨識碼",
                    "type": "string"
                },
                "telephone": {
                    "description": "電話號碼",
                    "type": "string"
                },
                "telephoneRegion": {
                    "description": "國家地區代碼",
                    "type": "string",
                    "example": "TW"
                }
            }
        },
        "domain.ResSmsCaptchaValidateDTO": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "accessTokenExp": {
                    "description": "ExpiresAt   int64  ` + "`" + `json\"expiresAt,omitempty\"` + "`" + `",
                    "type": "integer"
                },
                "message": {
                    "description": "回傳訊息",
                    "type": "string"
                },
                "process": {
                    "description": "回傳狀態",
                    "type": "string"
                }
            }
        },
        "domain.ResponseError": {
            "type": "object",
            "properties": {
                "errorCode": {
                    "description": "自定義錯誤代碼",
                    "type": "string"
                },
                "message": {
                    "description": "錯誤訊息",
                    "type": "string"
                },
                "path": {
                    "description": "API路徑",
                    "type": "string"
                },
                "status": {
                    "description": "錯誤狀態碼",
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Seeks Auth Service",
	Description:      "Seeks Auth API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
