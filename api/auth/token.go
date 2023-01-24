package auth

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/PixDale/sh-code-challenge/api/responses"
)

// BearerTokenSize represents the size of a slice containing a header bearer token
const bearerTokenSize = 2

// CreateToken creates a token for a given user and role
func CreateToken(userID uint32, role uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// TokenValid validates a token in a given request
func TokenValid(c *fiber.Ctx) error {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return responses.ErrorFailParseClaims
	}
	return nil
}

// ExtractToken extracts the token string from a given request
// the function checks for the token inside the query string and inside the headers
func ExtractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	splittedToken := strings.Split(bearerToken, " ")
	if len(splittedToken) == bearerTokenSize {
		return splittedToken[1]
	}
	return ""
}

// ExtractTokenID extracts the User ID from the token of a given request
func ExtractTokenID(c *fiber.Ctx) (uint32, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}

// HasRoleManager checks if the given token has a Manager Role
func HasRoleManager(c *fiber.Ctx) bool {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		role, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["role"]), 10, 32)
		if err != nil {
			return false
		}
		return uint32(role) == ManagerRole
	}
	return false
}

// HasRoleTechnician checks if the given token has a Technician Role
func HasRoleTechnician(c *fiber.Ctx) bool {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		role, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["role"]), 10, 32)
		if err != nil {
			return false
		}
		return uint32(role) == TechnicianRole
	}
	return false
}
