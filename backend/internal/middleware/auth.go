package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
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
	Logger    *logger.Logger
}

func NewAuthMiddleware(jwtSecret string, log *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{JWTSecret: jwtSecret, Logger: log}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.Logger.WithFields(logrus.Fields{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			}).Warn("missing auth header")
			writeError(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			m.Logger.WithFields(logrus.Fields{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			}).Warn("invalid auth header format")
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
			m.Logger.WithFields(logrus.Fields{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err,
			}).Warn("invalid or expired token")
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

func (m *AuthMiddleware) RequireRole(next http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*UserClaims)
		if !ok {
			m.Logger.WithFields(logrus.Fields{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			}).Warn("auth required for role check")
			writeError(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if claims.Role != role {
			m.Logger.WithFields(logrus.Fields{
				"ip":       r.RemoteAddr,
				"path":     r.URL.Path,
				"username": claims.Username,
				"user_id":  claims.UserID,
				"have":     claims.Role,
				"need":     role,
			}).Warn("insufficient role")
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
