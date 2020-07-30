package interceptor

import (
	"github.com/labstack/echo"
	"taskmanage_api/src/data"
	"taskmanage_api/src/exception"
)

var User = data.User{}

func CsrfAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := data.RedisGet(c.Request().Header.Get("user_token"))
		if err != nil {
			return exception.TokenException(c)
		}
		User = user
		if err := next(c); err != nil {
			c.Error(err)
		}
		User = data.User{}
		return nil
	}
}
