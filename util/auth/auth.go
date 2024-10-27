package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/db/models/auth"
	"github.com/tikhonp/alcs/util"
)

// Login logs in the user with the given login and password.
// Returns the user id if the login and password are correct.
// Saves the User in the session.
func Login(ctx echo.Context, login, password string) (*auth.User, error) {
	panic("not implemented")
}

// Logout logs out the user.
// Removes the user id from the session.
func Logout(ctx echo.Context) error {
	err := util.SetValue("userId", nil, ctx)
	ctx.Set("userId", nil)
	return err
}

