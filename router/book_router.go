package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rivaldiheriyan/managementsystem/controllers"
	"github.com/rivaldiheriyan/managementsystem/middlewares"
	"gorm.io/gorm"
)

func BookRouter(router *gin.Engine, db *gorm.DB){
	BookController := &controllers.BookDb{DB: db}

	router.GET("/book", BookController.ListBook)
	router.GET("/book/book_id", BookController.DetailBook)

	admin := router.Group("/")
	admin.Use(middlewares.AuthAdmin())
	{
		admin.POST("/create-book", BookController.CreateBook)
		admin.PUT("/edit-book/book_id", BookController.EditBook)
		admin.DELETE("/delete-book/book_id", BookController.DeleteBook)
	}
}