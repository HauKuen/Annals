package main

import (
	"github.com/HauKuen/Annals/model"
	"github.com/HauKuen/Annals/routes"
)

func main() {
	model.InitDb()
	routes.InitRouter()

}
