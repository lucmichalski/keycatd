package api

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/keydotcat/backend/util"
)

const CSRF_COOKIE_NAME = "4d018d7e07"

type CSRF struct {
	sc *securecookie.SecureCookie
}

func NewCSRF(hKey, bKey []byte) CSRF {
	return CSRF{securecookie.New(hKey, bKey)}
}

func (c CSRF) checkToken(w http.ResponseWriter, r *http.Request) bool {
	val, ok := r.Header["X-Csrf-Token"]
	if !ok {
		return false
	}
	if len(val) == 0 {
		return false
	}
	return c.getToken(w, r) == val[0]
}

func (c CSRF) getToken(w http.ResponseWriter, r *http.Request) string {
	csrfToken := ""
	if cookie, err := r.Cookie(CSRF_COOKIE_NAME); err == nil {
		if err = c.sc.Decode(CSRF_COOKIE_NAME, cookie.Value, &csrfToken); err == nil {
			return csrfToken
		}
	}
	csrfToken = util.GenerateRandomToken(32)
	if encoded, err := c.sc.Encode(CSRF_COOKIE_NAME, csrfToken); err == nil {
		cookie := &http.Cookie{
			Name:  CSRF_COOKIE_NAME,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	} else {
		panic(err)
	}
	return csrfToken
}
