package middleware

import (
	"os"
	//	"log"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//Verifies and returns the JWT Token that is stored on the cookie

func verifyToken(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

//Extracts the JWT Token used as cookie and verifies if the user is authenticated or not, also sets the ID of the user as part of the context

func CheckAuthentication(c *gin.Context) {

	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "c.Cookie me hi error aagya",
		})
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return

	}
	token, err := verifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "failed to verify",
		})
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return

	}
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["sub"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no ID claim",
		})
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return

	}
	c.Set("ID", uuid.MustParse(userID))
	c.Next()
}
