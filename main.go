package main

import (
	"github.com/aysf/gojwt/config"
	"github.com/aysf/gojwt/middlewares"
	"github.com/aysf/gojwt/routes"
)

func main() {

	config.InitDB()

	e := routes.New()
	middlewares.LogMiddlewares(e)
	e.Logger.Fatal(e.Start(":8080"))
}
