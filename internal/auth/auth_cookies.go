package auth

import "net/http"

const AccessTokenCookieName = "access_token"
const RefreshTokenCookieName = "refresh_token"
const DefaultCookieMaxAge = 3600

func NewAccessTokenCookie(token string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func NewRefreshTokenCookie(token string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
