package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// need install GCC compiler

	"go-test/handlers"
)

func main() {
	r := gin.Default()

	// Conexion ala base de datos
	app := &handlers.App{}
	app.ConnectDB()
	// Enable CORS middleware with permissive configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// routes Bills
	r.GET("/bills", app.BillsHandler)
	r.POST("/bills", app.NewBillHandler)
	r.GET("/bills/:id", app.GetBillHandler)
	r.DELETE("/bills/:id", app.DeleteBillHandler)

	// routes Users
	r.GET("/users", app.UsersHandler)
	r.POST("/users", app.NewUserHandler)
	r.GET("/users/:id", app.GetUserHandler)
	r.DELETE("/users/:id", app.DeleteUserHandler)

	r.Run(":8502")
}
