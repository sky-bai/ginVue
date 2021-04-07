package main

import (
	"ginVue/model"
	"ginVue/routers"
)

func main() {
	model.InitDb()
	routers.InitRouter()
}