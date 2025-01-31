package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Println("Warning: JWT_SECRET environment variable is not set")
	} else {
		log.Println("JWT_SECRET loaded successfully")
	}
}

var roleHierarchy = map[string]int{
	"admin":   3,
	"manager": 2,
	"user":    1,
}

func AuthMiddleware(requiredRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				log.Println("Missing or invalid token")
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
			}

			token = strings.TrimPrefix(token, "Bearer ")

			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Println("Invalid token signing method")
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
				}
				return jwtSecret, nil
			})

			if err != nil || !parsedToken.Valid {
				log.Println("Failed to parse token:", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Failed to extract claims from token")
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}

			role, exists := claims["role"].(string)
			if !exists {
				log.Println("Role not found in token")
				return echo.NewHTTPError(http.StatusUnauthorized, "Role not found")
			}

			if !isRoleAllowed(role, requiredRoles) {
				log.Println("User does not have the required role")
				return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
			}

			c.Set("username", claims["username"])
			c.Set("role", role)

			return next(c)
		}
	}
}

func isRoleAllowed(userRole string, requiredRoles []string) bool {
	userLevel, userExists := roleHierarchy[userRole]
	for _, role := range requiredRoles {
		if level, exists := roleHierarchy[role]; exists && userExists {
			if userLevel >= level {
				return true
			}
		}
	}
	return false
}
