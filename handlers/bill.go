package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	Concept string  `json:"concept"`
	Price   float32 `json:"price"`
}

// Rutas Bills
func (app *App) GetBillHandler(c *gin.Context) {
	var bill Bill

	billID := c.Param("id")

	// Check if the "id" parameter is empty
	if billID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	// find the first value in the data base with billID
	if err := app.DB.First(&bill, billID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	c.JSON(http.StatusOK, bill)
}
func (app *App) BillsHandler(c *gin.Context) {
	// fecthar los datos de la base de datos
	var bills []Bill
	app.DB.Find(&bills)

	// Set the "Access-Control-Allow-Origin" header to allow all origins (*)
	c.Header("Access-Control-Allow-Origin", "*")

	c.JSON(200, gin.H{
		"bills": bills,
	})
}
func (app *App) NewBillHandler(c *gin.Context) {
	var newBill Bill

	// Convierte el Json en el tipo de objeto que necesitamos
	if err := c.ShouldBindJSON(&newBill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Aqui guardaremos en la base de datos
	if err := app.DB.Create(&newBill).Error; err != nil {
		//Verificamos si se cumple alguna otra condicion
		if isForeignKeyConstraintError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Foreign key constraint failed"})
			return
		}
		// Si no es un error de clave externa, devolver otro mensaje de error genérico
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Bill"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Received JSON", "data": newBill.Concept})
}
func (app *App) DeleteBillHandler(c *gin.Context) {
	var bill Bill

	billID := c.Param("id")

	// Check if the "id" parameter is empty
	if billID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	// find the first value in the data base with billID
	if err := app.DB.First(&bill, billID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	// delete the bill is found
	if err := app.DB.Delete(&bill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Bill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill delete successfully"})
}

func isForeignKeyConstraintError(err error) bool {
	sqliteErr, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}

	return sqliteErr.Code == sqlite3.ErrConstraint
}
