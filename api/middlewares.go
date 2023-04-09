package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/varjangn/urlserv/storage"
)

type Contextkey int

const authUserKey Contextkey = 1

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func JWTAuth(next http.HandlerFunc, s storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := ExtractJWTToken(r)
		if err != nil {
			log.Println("failed to extact token")
			WriteJSON(w, http.StatusUnauthorized,
				Map{"error": err.Error()})
			return
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError,
				Map{"error": err.Error()})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			WriteJSON(w, http.StatusUnauthorized,
				Map{"error": "invalid token"})
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			WriteJSON(w, http.StatusUnauthorized,
				Map{"error": "token expired"})
			return
		}
		userEmail := claims["sub"].(string)
		user, err := s.GetUser(userEmail)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized,
				Map{"error": "invalid token"})
			return
		}
		ctx := context.WithValue(r.Context(), authUserKey, user)
		next(w, r.WithContext(ctx))
	})
}
