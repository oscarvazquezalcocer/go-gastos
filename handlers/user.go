package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bills    []Bill `json:"bills" gorm:"OnDelete:CASCADE"`
}

// Rutas Users
func (app *App) GetUserHandler(c *gin.Context) {
	var user User

	userID := c.Param("id")

	// Check if the "id" parameter is empty
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// find the first value in the data base with userID
	if err := app.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
func (app *App) UsersHandler(c *gin.Context) {
	// fecthar los datos de la base de datos
	var users []User
	// SELECT * FROM users;
	// SELECT * FROM bills WHERE user_id IN (1,2,3,4);
	app.DB.Preload("Bills").Find(&users)

	// Set the "Access-Control-Allow-Origin" header to allow all origins (*)
	c.Header("Access-Control-Allow-Origin", "*")

	c.JSON(200, gin.H{
		"users": users,
	})
}
func (app *App) NewUserHandler(c *gin.Context) {
	var newUser User

	// Convierte el Json en el tipo de objeto que necesitamos
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Aqui guardaremos en la base de datos
	if err := app.DB.Create(&newUser).Error; err != nil {
		//Verificamos si se cumple alguna otra condicion

		// Si no es un error de clave externa, devolver otro mensaje de error genérico
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create User"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Received JSON", "data": newUser.Name})
}
func (app *App) DeleteUserHandler(c *gin.Context) {
	var user User

	userID := c.Param("id")

	// Check if the "id" parameter is empty
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// find the first value in the data base with userID
	if err := app.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// delete the user and bills associated
	if err := app.DB.Select("Bill").Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete User"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User delete successfully"})
}
