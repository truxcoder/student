package middleware

import (
	"net/http"
	"sync"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	var once sync.Once
	once.Do(func() {
		tokenAuth = jwtauth.New("HS256", []byte("jfaeb3!$lfu24&4hh4k322gri65DdfeY2lkiu32!234"), nil)
		token, _ := GenerateToken("admin")
		println("token: ", token)
	})
}

func JWTVerifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

func JWTAuthenticator() func(http.Handler) http.Handler {
	return jwtauth.Authenticator(tokenAuth)
}

func GenerateToken(userID string) (string, error) {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": userID})
	return tokenString, err
}

func GetClaims(r *http.Request) map[string]interface{} {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return claims
}
