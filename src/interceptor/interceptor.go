package interceptor

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"taskmanage_api/src/domain"
)

var User = domain.User{}

type intercept struct {
	rr domain.RedisRepository
}

func NewIntercept(rr domain.RedisRepository) *intercept {
	return &intercept{
		rr: rr,
	}
}

//func (i *intercept)CsrfAuth(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		user, err := i.rr.RedisGet(c.Request().Header.Get("user_token"))
//		if err != nil {
//			return exception.TokenException(c)
//		}
//		User = user
//		if err := next(c); err != nil {
//			c.Error(err)
//		}
//		User = domain.User{}
//		return nil
//	}
//}

func (i *intercept) CsrfAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["user_id"]
		u := domain.User{
			ID: int(userId.(float64)),
		}
		User = u
		if err := next(c); err != nil {
			c.Error(err)
		}
		User = domain.User{}
		return nil
	}
}