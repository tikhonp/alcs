package util

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// in seconds
const sessionMaxAge = 86400 * 7

func GetValue(key interface{}, c echo.Context) (interface{}, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}
	return sess.Values[key], nil
}

func SetValue(key, value interface{}, c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		HttpOnly: true,
	}
	sess.Values[key] = value
	return sess.Save(c.Request(), c.Response())
}

