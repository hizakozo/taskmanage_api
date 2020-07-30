package exception

import (
	"github.com/labstack/echo"
	"net/http"
	"taskmanage_api/src/response"
)

func PermissionException(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Permission Exception"})
}

func TokenException(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, "Token exception")
}

func NotFoundData(c echo.Context) error {
	return c.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Not found data"})
}

func InputFailed(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Input valid failed"})
}

func FormBindException(c echo.Context) error {
	return c.JSON(http.StatusServiceUnavailable, "Input bind failed")
}

func FileUploadFailed(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: "File upload failed"})
}

func DataAlreadyExists(c echo.Context, data string) error {
	return c.JSON(http.StatusConflict, response.ErrorResponse{Message: data + "already exists"})
}
