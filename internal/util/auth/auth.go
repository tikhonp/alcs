package auth

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/db/models/auth"
	"github.com/tikhonp/alcs/internal/util"
)

var ErrNotAuthentificated = errors.New("user arent authentificated for given request")

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

// GetUser retriesves user id from the request session
// and fetches an object from db.
func GetUser(ctx echo.Context, users auth.Users) (*auth.User, error) {
	userId, ok := ctx.Get("userId").(int)
    if !ok {
        return nil, ErrNotAuthentificated
    } 
    return users.GetById(userId)
}
