package user

import (
	"errors"
	"net/http"

	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(id uuid.UUID, expire time.Time) (string, error) {
	claim := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := tokens.SignedString(secretKey)

	if err != nil {
		slog.Error("GenerateJwtToken err:", slog.Any("err", err))
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")

}
func AddCookieTokens(id uuid.UUID, w http.ResponseWriter, domain string) error {
	expirationTimeAccess := time.Now().Add(4 * time.Hour)
	expirationTimeRefresh := time.Now().Add(14 * 24 * time.Hour)
	refreshToken, err := GenerateJwtToken(id, expirationTimeRefresh)
	if err != nil {
		return err
	}
	accessToken, err := GenerateJwtToken(id, expirationTimeAccess)
	if err != nil {
		return err
	}
	http.SetCookie(w, GenerateCookie("accessToken", expirationTimeAccess, false, accessToken, domain))
	http.SetCookie(w, GenerateCookie("refreshToken", expirationTimeRefresh, true, refreshToken, domain))

	return nil
}
func NewAccessToken(tokenString string, threshold time.Duration, w http.ResponseWriter, domain string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime > threshold {
		return tokenString, nil
	}

	newExpire := time.Now().Add(4 * time.Hour)
	newToken, err := GenerateJwtToken(claims.ID, newExpire)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, GenerateCookie("accessToken", newExpire, false, newToken, domain))
	return newToken, nil
}
func GenerateCookie(name string, expire time.Time, httpOnly bool, value string, domain string) *http.Cookie {
	cookie := &http.Cookie{
		Name:        name,
		Value:       value,
		Expires:     expire,
		Partitioned: true,
		Path:        "/",
		Secure:      true,
		HttpOnly:    httpOnly,
		SameSite:    http.SameSiteLaxMode,
	}
	if domain := os.Getenv("DOMAIN_PROD"); domain != "" {
		cookie.Domain = domain
	}
	return cookie
}
