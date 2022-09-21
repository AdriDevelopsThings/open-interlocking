package authorization

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/adridevelopsthings/open-interlocking/pkg/config"
	"github.com/golang-jwt/jwt/v4"
)

type Claim struct {
	Permissions []string `json:"permissions"`
	ComputerID  string   `json:"cid"`
	Version     string   `json:"version"`
	jwt.StandardClaims
}

func CheckPermission(permissions []string, permissionNeeded string) bool {
	for _, permission := range permissions {
		p, _ := regexp.Compile(permission)
		if p.MatchString(permissionNeeded) {
			return true
		}
	}
	return false
}

func CheckAuthorization(authorization_header string, permission string, checkPermissions bool) int {
	if !strings.HasPrefix(authorization_header, "Bearer") {
		return http.StatusUnauthorized
	}

	claims := &Claim{}
	token, err := jwt.ParseWithClaims(authorization_header[7:], claims, func(token *jwt.Token) (interface{}, error) {
		secret := config.JWTSecret
		return []byte(secret), nil
	})
	if err != nil {
		return http.StatusForbidden
	}
	if !token.Valid {
		return http.StatusForbidden
	}

	if checkPermissions && !CheckPermission(claims.Permissions, permission) {
		return http.StatusForbidden
	}
	return 0
}

func CreateToken(cid *string, permissions *[]string) (string, error) {
	claims := &Claim{Permissions: *permissions, ComputerID: *cid, Version: config.JWTVersion}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	return tokenString, err
}
