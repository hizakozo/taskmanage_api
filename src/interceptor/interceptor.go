package interceptor

import (
	"github.com/labstack/echo"
	"taskmanage_api/src/domain"
	"taskmanage_api/src/exception"
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

func (i *intercept)CsrfAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := i.rr.RedisGet(c.Request().Header.Get("user_token"))
		if err != nil {
			return exception.TokenException(c)
		}
		User = user
		if err := next(c); err != nil {
			c.Error(err)
		}
		User = domain.User{}
		return nil
	}
}
