package handler

import "github.com/labstack/echo/v4"

type UserHandler struct{}

func (u UserHandler) HandleUserShow(c echo.Context) error {
	return nil
}
