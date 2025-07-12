package main

import (
	"github.com/rivaldiheriyan/managementsystem/config"
	routers "github.com/rivaldiheriyan/managementsystem/router"
)

func main(){
	config.LoadEnv()
	db := config.DBConnect()
	r := routers.SetupRouter(db)
	r.Run() // default di :8080
}