package handlers

import (
	// data base sqlite3
	//"gorm.io/driver/sqlite"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

// Conexion Base De datos
const (
	DBHost     = "localhost"
	DBPort     = 5432
	DBUser     = "postgres"
	DBPassword = ""
	DBName     = "postgres"
)

func (app *App) ConnectDB() {

	var err error
	//Conexion sqlite3
	//app.DB, err = gorm.Open(sqlite.Open("bills.sqlite"), &gorm.Config{})

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)
	app.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("error al conectar ala base de datos")
	}

	// AutoMigrate intenta crear la tabala si no existe
	err = app.DB.AutoMigrate(&User{}, &Bill{})

	// Habilitamos la funcion para poder usar claves foraneas
	app.DB.Exec("PRAGMA foreign_keys = ON;")

	if err != nil {
		panic("Error al crear la tabla en la base de datos")
	}
}
