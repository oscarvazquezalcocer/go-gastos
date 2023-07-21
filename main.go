package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	//"gorm.io/gorm"
)

type Bill struct {
	Id          int       `json:"id"`
	Date        time.Time `json:"date"`
	Concept     string    `json:"concept"`
	Price 	    float32   `json:"price"`
}

func billsHandler(c *gin.Context) {
	// Dummy los datos de una lista 
	// TODO fecthar los datos de la base de datos
	bills := []Bill{
		Bill{1, time.Now(), "Tacos", 200 },
		Bill{2, time.Now(), "Oxxo", 300},
		Bill{3, time.Now(), "Pan", 50},
	}

	// Set the "Access-Control-Allow-Origin" header to allow all origins (*)
	c.Header("Access-Control-Allow-Origin", "*")

	c.JSON(200, gin.H{
		"bills": bills,
	})
}

func newBillHandler(c *gin.Context) {
	var jsonBill Bill

	// Convierte el Json en el tipo de objeto que necesitamos
	if err := c.ShouldBindJSON(&jsonBill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	
	// Aqui guardaremos en la base de datos
	// TODO: Guardar base de datos
	concept := jsonBill.Concept

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Received JSON", "data": concept})
}

func main() {
	r := gin.Default()

	// Enable CORS middleware with permissive configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/bills", billsHandler)
	r.POST("/bills", newBillHandler)
	r.Run(":8502")
}
