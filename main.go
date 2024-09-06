package main

import (
	"github.com/HauKuen/Annals/model"
	"github.com/HauKuen/Annals/routes"
	"github.com/HauKuen/Annals/utils"
)

func main() {
	utils.InitLogger()
	model.InitDb()
	routes.InitRouter()

}
