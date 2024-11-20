package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/util"
)

// Saves user id to the session
func Login(ctx echo.Context, userId int) error {
	err := util.SetValue("userId", userId, ctx)
	ctx.Set("userId", userId)
	return err
}

// Logout logs out the user.
// Removes the user id from the session.
func Logout(ctx echo.Context) error {
	err := util.SetValue("userId", nil, ctx)
	ctx.Set("userId", nil)
	return err
}
