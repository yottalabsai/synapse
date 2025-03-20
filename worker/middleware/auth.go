package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

const (
	authenticationHeader = "Authentication"
	acceptLanguage       = "Accept-Language"
	cfIpCountry          = "cf-ipcountry"
)

var (
	ErrBadAuth      = fmt.Errorf("malformed authentication")
	ErrAuthNotFound = fmt.Errorf("authentication info not found")
)

func MustGetUserInfo(ctx context.Context) *ContextUserInfo {
	if info, ok := GetUserInfo(ctx); ok {
		return info
	}
	panic("can't get user info")
}

func GetUserInfo(ctx context.Context) (*ContextUserInfo, bool) {
	v := ctx.Value(authenticationHeader)
	userInfo, ok := v.(*ContextUserInfo)
	return userInfo, ok
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement authentication
		//err := setUserInfo(c, true)
		//if err != nil {
		//	common.JSON(c, http.StatusOK, common.ErrUnauthorized)
		//	c.Abort()
		//}
		c.Next()
	}
}

type ContextUserInfo struct {
	Token string `json:"token"`
	User  struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
	} `json:"user"`
}

func (c ContextUserInfo) UserID() string {
	return c.User.UserID
}

func setUserInfo(c *gin.Context, required bool) error {
	rawUserInfo := c.Request.Header.Get(authenticationHeader)
	if len(rawUserInfo) != 0 {
		userInfoBytes, err := base64.StdEncoding.DecodeString(rawUserInfo)
		if err != nil {
			return errors.WithMessagef(err, "can't decode user info")
		}
		ui := new(ContextUserInfo)
		if err := jsoniter.Unmarshal(userInfoBytes, ui); err != nil {
			return ErrBadAuth
		}
		c.Set(authenticationHeader, ui)
		return nil
	}

	if required {
		return ErrAuthNotFound
	}
	return nil
}

func RequestHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		SetAcceptLanguage(c)
		setCfIpCountry(c)
		c.Next()
	}
}

func SetAcceptLanguage(c *gin.Context) {
	c.Set(acceptLanguage, c.Request.Header.Get(acceptLanguage))
}

func setCfIpCountry(c *gin.Context) {
	c.Set(cfIpCountry, c.Request.Header.Get(cfIpCountry))
}

func GetAcceptLanguage(ctx context.Context) string {
	v := ctx.Value(acceptLanguage)
	if typ, ok := v.(string); ok {
		if len(typ) == 0 {
			return language.English.String()
		}
		return typ
	}
	return language.English.String()
}

func GetCfIpCountry(ctx context.Context) string {
	v := ctx.Value(cfIpCountry)
	if typ, ok := v.(string); ok {
		return typ
	}
	return ""
}
