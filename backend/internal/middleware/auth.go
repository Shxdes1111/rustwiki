package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	JWTSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{JWTSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeError(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims := &UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			writeError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) AuthenticateOptional(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next(w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			next(w, r)
			return
		}

		tokenStr := parts[1]
		claims := &UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			next(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func RequireRole(next http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*UserClaims)
		if !ok {
			writeError(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if claims.Role != role {
			writeError(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func GetUserClaims(r *http.Request) *UserClaims {
	claims, _ := r.Context().Value(UserContextKey).(*UserClaims)
	return claims
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
