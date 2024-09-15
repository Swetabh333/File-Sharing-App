package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	models "github.com/Swetabh333/trademarkia/models"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	bcrypt "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Struct for storing the request body for registration
type register struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct for storing the request body for login
type login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// function to encrypt the password before storing in database
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// function to decrypt and compare the password while loggin in
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Handler function for /register route
func HandleRegistration(db *gorm.DB) gin.HandlerFunc {
	//function to pass database to handler function
	return func(c *gin.Context) { //actual handler function
		var regBody register

		err := c.BindJSON(&regBody)
		if err != nil {
			fmt.Println("Error binding request body")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some internal error occured",
			})
			return
		}
		hashedPassword, err := HashPassword(regBody.Password)
		if err != nil {
			fmt.Println("Error hashing password")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some internal error occured",
			})
			return

		}
		user := models.User{
			ID:       uuid.New(),
			Name:     regBody.Name,
			Email:    regBody.Email,
			Password: hashedPassword,
		}
		err = db.Create(&user).Error
		if err != nil {
			fmt.Printf("Error storing user in database: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some internal error occured",
			})
			return
		}
		fmt.Println("User successfully registered")
		c.JSON(http.StatusOK, gin.H{
			"message": "User successfully created",
		})
	}
}

//function for creating your JWT Token

func generateJWT(username string, ID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = username
	claims["sub"] = ID
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Handler function for /login route
func HandleLogin(db *gorm.DB) gin.HandlerFunc {
	//function to pass database to the handler function
	return func(c *gin.Context) { //handler function

		login := login{}
		user := models.User{}
		err := c.BindJSON(&login)
		if err != nil {
			fmt.Println("Error binding request body")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some internal error occured",
			})
			return
		}
		err = db.Where("name = ?", login.Name).Find(&user).Error
		if err != nil {
			fmt.Println("User not found")
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Username does not exist",
			})
			return
		}
		check := CheckPasswordHash(login.Password, user.Password)
		if check {
			token, err := generateJWT(user.Name, user.ID.String())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error generating cookies",
				})
			}
			c.SetCookie("token", token, 3600, "/", "", false, true)
			fmt.Println("Logged in")
			c.JSON(http.StatusOK, gin.H{
				"message": "Logged in successfully",
			})
		} else {
			fmt.Println("Password did not match")
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Incorrect Password",
			})
		}
	}
}
