package middlewares

import (
	"errors"
	"fmt"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rivaldiheriyan/managementsystem/config"
	"github.com/rivaldiheriyan/managementsystem/helpers"
	"github.com/rivaldiheriyan/managementsystem/models"
	"gorm.io/gorm"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			fmt.Println("Bad Authorization")
			c.AbortWithError(http.StatusBadRequest, errors.New("Bad Authorization"))
			c.JSON(401, gin.H{
				"message": "Unautorized",
			})
			return
		}

		authSplit := strings.Split(auth, " ")
		if len(authSplit) != 2 {
			fmt.Println("Bad Authorization")
			c.AbortWithError(http.StatusBadRequest, errors.New("Bad Authorization"))
			c.JSON(401, gin.H{
				"message": "Unautorized",
			})
			return
		}
		if authSplit[0] != "Bearer" {
			fmt.Println("Bad Authorization")
			c.AbortWithError(http.StatusBadRequest, errors.New("Bad Authorization"))
			c.JSON(401, gin.H{
				"message": "Unautorized",
			})
			return
		}

		token, err := helpers.VerifyToken(authSplit[1])
		if err != nil {
			fmt.Println("Bad Authorization")
			c.AbortWithError(http.StatusBadRequest, errors.New("Bad Authorization"))
			c.JSON(401, gin.H{
				"message": "Unautorized",
			})
			return
		}

		dbResult := models.User{}
		username := token["username"]
		err = config.DBConnect().Debug().Where("username = ?", username).Last(&dbResult).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("Username not found")
				c.AbortWithError(http.StatusBadRequest, errors.New("Username not found"))
				c.JSON(401, gin.H{
					"message": "Unautorized",
				})
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			c.JSON(401, gin.H{
				"message": "Unautorized",
			})
			return
		}

		// Ini set usernya jangan jadiin int64
		// c.Set("userId", strconv.FormatInt(int64(dbResult.ID), 10))
		c.Set("userUUID", dbResult.UserID.String())
		c.Set("role", dbResult.Role)
		
		c.Next()
	}
}
