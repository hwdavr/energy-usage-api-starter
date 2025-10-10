package http

// Import UserIDKey from userctx package
import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/yourname/energy-usage-api/internal/pkg/userctx"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				log.Println("Missing or malformed Authorization header:", auth)
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			}, jwt.WithValidMethods([]string{"HS256"}))
			if err != nil {
				log.Println("JWT parse error:", err)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				log.Println("JWT token is not valid")
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("JWT claims are not of type MapClaims")
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			uid, ok := claims["uid"].(string)
			if !ok || uid == "" {
				log.Println("UID claim missing or not a string in token:", claims)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userctx.UserIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
