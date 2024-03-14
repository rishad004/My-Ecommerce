package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var BlacklistedTokens = make(map[string]bool)

type Claims struct {
	UserID uint
	Email  string
	Role   string
	Log    bool
	jwt.StandardClaims
}

func JwtCreate(c *gin.Context, UserID uint, Email string, Role string) {
	claims := Claims{
		UserID: UserID,
		Email:  Email,
		Role:   Role,
		Log:    true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 4).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println("=======Error JWT Create", err)
		c.JSON(403, gin.H{
			"Error": "Failed to create Token",
		})
		return
	}
	c.JSON(201, gin.H{
		"Token": tokenString,
	})

}

func Auth(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------AUTH MIDDLEWARE----------------------")

	path := c.Request.URL.Path

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"Error": "No Authorization Token found."})
		c.Abort()
		return
	}
	if BlacklistedTokens[tokenString] {
		c.JSON(401, gin.H{"Error": "Revoked token"})
		c.Abort()
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid || !claims.Log {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if path[1] == 97 {
		if claims.Role != "Admin" {
			c.JSON(401, gin.H{"Error": "Not Authorized"})
			c.Abort()
			return
		}
	}
	if path[1] == 117 {
		if claims.Role != "User" {
			c.JSON(401, gin.H{"Error": "Not Authorized"})
			c.Abort()
			return
		}
	}

	c.Set("token", tokenString)
	c.Set("Id", claims.UserID)
	c.Set("Email", claims.Email)
	c.Set("Role", claims.Role)

	c.Next()
}
