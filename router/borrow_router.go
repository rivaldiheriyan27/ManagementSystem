package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rivaldiheriyan/managementsystem/controllers"
	"gorm.io/gorm"
)

func BorrowRouter(router *gin.Engine, db *gorm.DB){
	BorrowController := &controllers.BorrowDb{DB: db}

	router.GET("/list-borrow", BorrowController.ListBorrow)
	router.POST("/borrow-book", BorrowController.BorrowBook)
	router.POST("/return-book", BorrowController.ReturnBook)

}