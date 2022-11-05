package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"fetch-app/conf"
	"fetch-app/constants"
	"fetch-app/datatransfers"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	jwt.StandardClaims
	UID      string `json:"uid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

type JWTConfig struct {
	JWTSignatureKey   string
	JWTPublicKey      string
	JWTExpirationTime time.Duration
}

type UserData struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

func VerifyTokenAdmin(ctx *context.Context) {
	userData, err := doVerifyToken(ctx.Request)
	if err != nil || !userData.IsAdmin {
		errAuth(ctx)
		return
	}
}

func VerifyToken(ctx *context.Context) {
	userData, err := doVerifyToken(ctx.Request)
	if err != nil {
		errAuth(ctx)
		return
	}

	ctx.Input.SetData("uid", userData.UID)
	ctx.Input.SetData("isAdmin", userData.IsAdmin)
}

func doVerifyToken(r *http.Request) (result *UserData, err error) {

	token, err := getToken(r)
	if err != nil {
		return
	}

	isVerified, claims, err := parseTokenJWT(token)
	if err != nil {
		return
	}

	if !isVerified {
		return
	}

	result = &UserData{
		UID:     claims.UID,
		IsAdmin: claims.IsAdmin,
	}
	return

}

func getToken(r *http.Request) (token string, err error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err = errors.New("token is empty")
		return
	}

	s := strings.Split(authHeader, " ")
	if len(s) != 2 {
		err = errors.New("token is invalid")
		return
	}

	token = s[1]
	return
}

func parseTokenJWT(token string) (isVerified bool, result *JWTClaims, err error) {
	result = &JWTClaims{}
	jwtClaims, err := jwt.ParseWithClaims(token, result, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.AppConfig.JWTConfig.JWTSignatureKey), nil
	})

	if result == nil || jwtClaims == nil || !jwtClaims.Valid || err != nil {
		return
	}
	isVerified = true
	return
}

func errAuth(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusUnauthorized)
	errResponse := &datatransfers.CustomError{
		Code:    constants.NotAuthorizedErrCode,
		Status:  http.StatusUnauthorized,
		Message: "UNAUTHORIZED",
	}

	resBody, err := json.Marshal(errResponse)
	if err != nil {
		panic(err)
	}
	ctx.Output.Body(resBody)
}
