package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rivaldiheriyan/managementsystem/controllers"
	"github.com/rivaldiheriyan/managementsystem/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine{
	router := gin.Default()
	LoginController := &controllers.LoginDB{DB: db}
	

	router.GET("/", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})


	router.POST("/login", LoginController.Login)
	router.POST("/register", LoginController.Register)

	router.Use(middlewares.AuthJWT())

	BookRouter(router,db)
	BorrowRouter(router,db)


	router.Run(":5022")
	return router
}