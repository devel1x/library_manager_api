package main

import (
	_ "template/docs"
	"template/internal/app"
)

// @title Library Manager API
// @version 1.0
// @description	This is a sample implementation of a RESTful API for a library.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
func main() {
	app.Start()
}
