package session

import (
	"fmt"
	"net/http"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func IsAuth(c echo.Context) (*domain.UserInfo, error) {
	sess, err := session.Get("session", c)

	if err != nil {
		return nil, err
	}

	refreshToken, ok := sess.Values["refresh_token"].(string)

	if !ok {
		return nil, fmt.Errorf("no refresh_token in cookie session")
	}

	claims, err := jwt.ValidateRefreshToken(refreshToken)

	if err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		return nil, err
	}

	ui := &domain.UserInfo{
		ID:       uid,
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}

	return ui, nil
}

func DefaultSaveSession(c echo.Context, RefrestToken *string) error {
	// This cookie to store refresh token
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	// This cookie to store expiration time
	exp, err := session.Get("expire", c)
	if err != nil {
		return err
	}

	// session
	sess.Options = &sessions.Options{
		Secure:   true,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   2592000,
		SameSite: http.SameSiteNoneMode,
	}

	if RefrestToken != nil {
		sess.Values["refresh_token"] = *RefrestToken
	}

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		return err
	}

	// expiration
	exp.Options = &sessions.Options{
		Secure:   true,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   2592000,
		SameSite: http.SameSiteNoneMode,
	}

	err = exp.Save(c.Request(), c.Response())

	if err != nil {
		return err
	}

	return nil
}

func DeleteSession(c echo.Context) error {
	// session
	sess, err := session.Get("session", c)

	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/session",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}

	delete(sess.Values, "refresh_token")

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		return err
	}

	// expiration
	exp, err := session.Get("expire", c)

	if err != nil {
		return err
	}

	exp.Options = &sessions.Options{
		Path:     "/expire",
		HttpOnly: false,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
	}

	err = exp.Save(c.Request(), c.Response())

	if err != nil {
		return err
	}

	return nil
}
