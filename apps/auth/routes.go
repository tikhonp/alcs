package auth

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/tikhonp/alcs/apps/auth/handlers"
	"github.com/tikhonp/alcs/config"
	"github.com/tikhonp/alcs/db"
	"github.com/tikhonp/alcs/util/auth"
)

func ConfigureAuthGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory) {
	ah := handlers.AuthHandler{}

	goth.UseProviders(
		google.New(cfg.Auth.GoogleKey, cfg.Auth.GoogleSecret, fmt.Sprintf("%s/auth/google/callback", cfg.BaseHost)),
	)

	g.GET("/login", ah.Login)

	g.GET("/:provider/callback", func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))

		guser, err := gothic.CompleteUserAuth(c.Response(), c.Request().WithContext(ctx))
		if err != nil {
			return err
		}

        user, err := modelsFactory.AuthUsers().FromOAuth(&guser)
        if err != nil {
            return err
        }

        auth.Login(c, user.Id)

		t, _ := template.New("foo").Parse(userTemplate)
		return t.Execute(c.Response(), guser)
	})

	g.GET("/logout/:provider", func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))

		gothic.Logout(c.Response(), c.Request().WithContext(ctx))
		c.Response().Header().Set("Location", "/")
		c.Response().WriteHeader(http.StatusTemporaryRedirect)
		return nil
	})

	g.GET("/:provider", func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, c.Param("provider"))

		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request().WithContext(ctx)); err == nil {
			t, _ := template.New("foo").Parse("") // Template string
			return t.Execute(c.Response(), gothUser)
		}

		gothic.BeginAuthHandler(c.Response(), c.Request().WithContext(ctx))

		return nil
	})

}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
