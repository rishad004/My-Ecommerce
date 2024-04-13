package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserID uint
	Email  string
	Role   string
	Log    bool
	jwt.StandardClaims
}

func JwtCreate(c *gin.Context, UserID uint, Email string, Role string) (string, error) {
	claims := Claims{
		UserID: UserID,
		Email:  Email,
		Role:   Role,
		Log:    true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	return tokenString, err
}

func Auth(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------AUTH MIDDLEWARE----------------------")

	var otp []models.Otp
	er := database.Db.Find(&otp, "Expr <?", time.Now()).Error
	if er != nil {
		fmt.Println("")
		fmt.Println("Error on deleting otps......................")
	} else {
		fmt.Println("")
		for _, v := range otp {
			fmt.Println("Deleting the expired otps..................")
			database.Db.Delete(&v)
		}
	}

	path := c.Request.URL.Path
	var tokenString string

	if path[1] == 97 {
		tokenString, _ = c.Cookie("Jwt-Admin")
	}
	if path[1] == 117 {
		tokenString, _ = c.Cookie("Jwt-User")
	}

	if tokenString == "" {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "No Authorization Token found!",
			"Data":    gin.H{},
		})
		c.Abort()
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid || !claims.Log {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Invalid token!",
			"Data":    gin.H{},
		})
		c.Abort()
		return
	}

	if path[1] == 117 {
		var user models.Users
		eror := database.Db.First(&user, "ID=?", claims.UserID).Error
		if eror != nil {
			c.JSON(404, gin.H{
				"Status":  "Error!",
				"Code":    404,
				"Message": "User not found!",
				"Data":    gin.H{},
			})
			c.Abort()
			return
		}

		if !user.Blocking {
			c.JSON(401, gin.H{
				"Status":  "Error!",
				"Code":    401,
				"Message": "You are blocked!",
				"Data":    gin.H{},
			})
			c.Abort()
			return
		}
		if claims.Role != "User" {
			c.JSON(401, gin.H{
				"Status":  "Error!",
				"Code":    401,
				"Message": "Not Authorized!",
				"Data":    gin.H{},
			})
			c.Abort()
			return
		}
	}

	if path[1] == 97 {
		if claims.Role != "Admin" {
			c.JSON(401, gin.H{
				"Status":  "Error!",
				"Code":    401,
				"Message": "Not Authorized!",
				"Data":    gin.H{},
			})
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
