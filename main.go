package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// need install GCC compiler
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	//Id      int       `json:"id"`
	//Date    time.Time `json:"date"`
	Concept string  `json:"concept"`
	Price   float32 `json:"price"`
}

var db *gorm.DB

func main() {
	r := gin.Default()

	// Conexion ala base de datos
	connectDB()
	// Enable CORS middleware with permissive configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// routes
	r.GET("/bills", billsHandler)
	r.POST("/bills", newBillHandler)
	r.GET("/bills/:id", getBillHandler)
	r.DELETE("/bills/:id", deleteBillHandler)

	r.Run(":8502")
}

func getBillHandler(c *gin.Context) {
	var bill Bill

	billID := c.Param("id")

	// Check if the "id" parameter is empty
	if billID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	// find the first value in the data base with billID
	if err := db.First(&bill, billID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

func billsHandler(c *gin.Context) {
	// fecthar los datos de la base de datos
	var bills []Bill
	db.Find(&bills)

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

func deleteBillHandler(c *gin.Context) {
	var bill Bill

	billID := c.Param("id")

	// Check if the "id" parameter is empty
	if billID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	// find the first value in the data base with billID
	if err := db.First(&bill, billID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	// delete the bill is found
	if err := db.Delete(&bill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Bill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill delete successfully"})
}

func connectDB() {

	var err error
	db, err = gorm.Open(sqlite.Open("bills.sqlite"), &gorm.Config{})
	if err != nil {
		panic("error al conectar ala base de datos")
	}

	// AutoMigrate intenta crear la tabala si no existe
	db.AutoMigrate(&Bill{})
}
