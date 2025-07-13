package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rivaldiheriyan/managementsystem/controllers"
	"github.com/rivaldiheriyan/managementsystem/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine{
	router := gin.Default()
	LoginController := &controllers.LoginDB{DB: db}
	
	// Tambahkan CORS Middleware di sini
	router.Use(middlewares.CORSMiddleware())


	router.POST("/login", LoginController.Login)
	router.POST("/register", LoginController.Register)

	router.Use(middlewares.AuthJWT())

	BookRouter(router,db)
	BorrowRouter(router,db)


	router.Run(":5022")
	return router
}