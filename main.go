package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// need install GCC compiler
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Bill struct {
	Id      int       `json:"id"`
	Date    time.Time `json:"date"`
	Concept string    `json:"concept"`
	Price   float32   `json:"price"`
}

func billsHandler(c *gin.Context) {
	// Dummy los datos de una lista
	// TODO fecthar los datos de la base de datos
	var bills []Bill
	db.Find(&bills)
	// bills := []Bill{
	// 	Bill{1, time.Now(), "Tacos", 200},
	// 	Bill{2, time.Now(), "Oxxo", 300},
	// 	Bill{3, time.Now(), "Pan", 50},
	// }

	// Set the "Access-Control-Allow-Origin" header to allow all origins (*)
	c.Header("Access-Control-Allow-Origin", "*")

	c.JSON(200, gin.H{
		"bills": bills,
	})
}

func newBillHandler(c *gin.Context) {
	var newBill Bill

	// Convierte el Json en el tipo de objeto que necesitamos
	if err := c.ShouldBindJSON(&newBill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Aqui guardaremos en la base de datos
	db.Create(&newBill)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Received JSON", "data": newBill.Concept})
}

var db *gorm.DB

func connectDB() {

	var err error
	db, err = gorm.Open(sqlite.Open("bills.sqlite"), &gorm.Config{})
	if err != nil {
		panic("error al conectar ala base de datos")
	}

	// AutoMigrate intenta crear la tabala si no existe
	db.AutoMigrate(&Bill{})
}

func main() {
	r := gin.Default()

	// Conexion ala base de datos
	connectDB()
	// Enable CORS middleware with permissive configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/bills", billsHandler)
	r.POST("/bills", newBillHandler)
	r.Run(":8502")
}
