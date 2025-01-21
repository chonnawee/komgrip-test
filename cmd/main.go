package main

import (
	"encoding/json"
	"fmt"
	"komgrip-test/adapters"
	"komgrip-test/entities"
	"komgrip-test/usecases"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var port string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("cannot read environment variables")
	}
}

func initializeConfig() {
	port = os.Getenv("APP_PORT")
}

func main() {
	initializeConfig()
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db connected")

	db.AutoMigrate(entities.Beers{})

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	uow := adapters.NewUnitOfWorkDB(db)
	usecase := usecases.NewBeersService(uow)
	handler := adapters.NewBeersHandler(usecase)

	app.Post("/beer", handler.CreateBeer)
	app.Get("/beer", handler.GetBeers)
	app.Put("/beer/:id", handler.UpdateBeer)
	app.Delete("/beer/:id", handler.DeleteBeer)
	app.Listen(":" + port)
}
